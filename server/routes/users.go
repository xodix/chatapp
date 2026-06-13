package routes

import (
	"net/http"
	"study4cash/DB/models"
	"study4cash/auth"
	"time"

	"github.com/alexedwards/argon2id"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/labstack/echo/v5"
)

var userCollection *mongo.Collection

func RouteUsers(prefix string, e *echo.Echo, db *mongo.Database) {
	userCollection = db.Collection("users")

	userRouter := e.Group(prefix)
	userRouter.POST("/register", register)
	userRouter.POST("/login", login)

	protected := userRouter.Group("")
	protected.Use(auth.JWTMiddleware)
	protected.GET("/", details)
	protected.DELETE("/", deleteAccount)
	protected.PUT("/", updateAccount)
	protected.GET("/all", getAllUsers)
}

type RegisterRequest struct {
	Email     string    `json:"email"  validate:"required,email"`
	Password  string    `json:"password"  validate:"required,min=8,max=32"`
	Name      string    `json:"name"  validate:"required,min=3,max=60"`
	Surname   string    `json:"surname"  validate:"required,min=3,max=60"`
	Birthdate time.Time `json:"birthdate"  validate:"required"`
}
type RegisterResponse struct {
	JWTToken string `json:"token"`
	UserID   string `json:"user_id"`
}

func register(c *echo.Context) error {
	registerData := new(RegisterRequest)
	if err := c.Bind(registerData); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(registerData); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	response := userCollection.FindOne(c.Request().Context(), map[string]interface{}{"email": registerData.Email})
	if response.Err() != nil {
		if response.Err() != mongo.ErrNoDocuments {
			return c.String(http.StatusInternalServerError, response.Err().Error())
		}
	} else {
		return c.String(http.StatusBadRequest, "User already exists")
	}

	password, err := argon2id.CreateHash(registerData.Password, argon2id.DefaultParams)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	userModel := models.User{
		Email:       registerData.Email,
		Password:    password,
		Name:        registerData.Name,
		Surname:     registerData.Surname,
		Birthdate:   registerData.Birthdate,
		Active:      true,
		LastLoginAt: time.Now(),
	}
	result, err := userCollection.InsertOne(c.Request().Context(), userModel)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	token, err := auth.GenerateJWT(result.InsertedID.(bson.ObjectID).Hex())
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, RegisterResponse{
		JWTToken: *token,
		UserID:   result.InsertedID.(bson.ObjectID).Hex(),
	})
}

type LoginRequest struct {
	Email    string `json:"email"  validate:"required,email"`
	Password string `json:"password"  validate:"required,min=8,max=32"`
}
type LoginResponse struct {
	JWTToken string `json:"token"`
	UserID   string `json:"user_id"`
}

func login(c *echo.Context) error {
	loginData := new(LoginRequest)
	if err := c.Bind(loginData); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(loginData); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	var user models.User
	response := userCollection.FindOne(c.Request().Context(), map[string]interface{}{"email": loginData.Email})
	if response.Err() != nil {
		if response.Err() == mongo.ErrNoDocuments {
			return c.String(http.StatusBadRequest, "User not found")
		}
		return c.String(http.StatusInternalServerError, response.Err().Error())
	}
	err := response.Decode(&user)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	match, err := argon2id.ComparePasswordAndHash(loginData.Password, user.Password)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Could not verify password")
	}
	if !match {
		return c.String(http.StatusUnauthorized, "Invalid password")
	}

	token, err := auth.GenerateJWT(user.ID.Hex())

	return c.JSON(http.StatusOK, LoginResponse{
		JWTToken: *token,
		UserID:   user.ID.Hex(),
	})
}

type DetailsResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"  validate:"required,email"`
	Name      string    `json:"name"  validate:"required,min=3,max=60"`
	Surname   string    `json:"surname"  validate:"required,min=3,max=60"`
	Birthdate time.Time `json:"birthdate"  validate:"required"`
}

func details(c *echo.Context) error {
	userID := c.Get("userID").(string)
	userIDObj, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	// users, err := gorm.G[models.User](db).Where("id = ?", userID).Find(c.Request().Context())
	var user models.User
	result := userCollection.FindOne(c.Request().Context(), map[string]interface{}{"_id": userIDObj})
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return c.String(http.StatusNotFound, "User not found")
		}
		return c.String(http.StatusInternalServerError, result.Err().Error())
	}
	err = result.Decode(&user)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, DetailsResponse{
		ID:        user.ID.String(),
		Email:     user.Email,
		Name:      user.Name,
		Surname:   user.Surname,
		Birthdate: user.Birthdate,
	})
}

func deleteAccount(c *echo.Context) error {
	userID := c.Get("userID").(string)
	userIDObj, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	result, err := userCollection.DeleteOne(c.Request().Context(), map[string]interface{}{"_id": userIDObj})
	if result.DeletedCount != 1 {
		return c.String(http.StatusNotFound, "Could not delete the user")
	}
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, "Account deleted")
}

type UpdateRequest struct {
	Email     string    `json:"email"  validate:"required,email"`
	Name      string    `json:"name"  validate:"required,min=3,max=60"`
	Surname   string    `json:"surname"  validate:"required,min=3,max=60"`
	Birthdate time.Time `json:"birthdate"  validate:"required"`
}

func updateAccount(c *echo.Context) error {
	userID := c.Get("userID").(string)
	userIDObj, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	var request UpdateRequest
	err = c.Bind(&request)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	updateResult, err := userCollection.UpdateOne(c.Request().Context(), map[string]interface{}{"_id": userIDObj}, map[string]interface{}{
		"$set": map[string]interface{}{
			"email":   request.Email,
			"name":    request.Name,
			"surname": request.Surname,
			"date":    request.Birthdate,
		},
	})
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	if updateResult.ModifiedCount != 1 {
		return c.String(http.StatusBadRequest, "Could not modify the user")
	}

	return c.String(http.StatusOK, "Updated account data")
}

type GetAllUsersUser struct {
	ID      bson.ObjectID `json:"id" bson:"_id,omitempty"`
	Name    string        `json:"name" bson:"name"`
	Surname string        `json:"surname" bson:"surname"`
}

type GetAllUsersResponse struct {
	Users []GetAllUsersUser `json:"users"`
}

func getAllUsers(c *echo.Context) error {
	result, err := userCollection.Find(c.Request().Context(), map[string]interface{}{})
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	var users []GetAllUsersUser
	err = result.All(c.Request().Context(), &users)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	var response GetAllUsersResponse
	for _, user := range users {
		response.Users = append(response.Users, GetAllUsersUser{
			ID:      user.ID,
			Name:    user.Name,
			Surname: user.Surname,
		})
	}

	return c.JSON(http.StatusOK, response)
}
