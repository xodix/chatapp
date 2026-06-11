package routes

import (
	"net/http"
	"study4cash/DB/models"
	"study4cash/auth"
	"time"

	"github.com/alexedwards/argon2id"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/labstack/echo/v5"
)

var userCollection *mongo.Collection

func RouteUsers(prefix string, e *echo.Echo, db *mongo.Database) {
	userCollection = db.Collection("users")

	userRouter := e.Group(prefix)
	userRouter.POST("/register", Register)
	userRouter.POST("/login", Login)

	protected := userRouter.Group("")
	protected.Use(auth.JWTMiddleware)
	protected.GET("/", Details)
	protected.DELETE("/", DeleteAccount)
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

func Register(c *echo.Context) error {
	registerData := new(RegisterRequest)
	if err := c.Bind(registerData); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(registerData); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	// users, err := gorm.G[models.User](db).Where("email = ?", registerData.Email).Find(c.Request().Context())
	var user models.User
	response := userCollection.FindOne(c.Request().Context(), map[string]interface{}{"email": registerData.Email})
	if response.Err() != nil {
		return c.String(http.StatusInternalServerError, response.Err().Error())
	}
	err := response.Decode(&user)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
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

	token, err := auth.GenerateJWT(result.InsertedID.(string))
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, RegisterResponse{
		JWTToken: *token,
		UserID:   result.InsertedID.(string),
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

func Login(c *echo.Context) error {
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

	token, err := auth.GenerateJWT(*user.ID)

	return c.JSON(http.StatusOK, LoginResponse{
		JWTToken: *token,
		UserID:   *user.ID,
	})
}

type DetailsResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"  validate:"required,email"`
	Name      string    `json:"name"  validate:"required,min=3,max=60"`
	Surname   string    `json:"surname"  validate:"required,min=3,max=60"`
	Birthdate time.Time `json:"birthdate"  validate:"required"`
}

func Details(c *echo.Context) error {
	userID := c.Get("userID").(uint)

	// users, err := gorm.G[models.User](db).Where("id = ?", userID).Find(c.Request().Context())
	var user models.User
	result := userCollection.FindOne(c.Request().Context(), map[string]interface{}{"_id": userID})
	if result.Err() != nil {
		return c.String(http.StatusInternalServerError, result.Err().Error())
	}
	err := result.Decode(&user)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, DetailsResponse{
		ID:        *user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Surname:   user.Surname,
		Birthdate: user.Birthdate,
	})
}

func DeleteAccount(c *echo.Context) error {
	userID := c.Get("userID").(uint)

	// tx := db.Delete(&models.User{}, userID)
	result, err := userCollection.DeleteOne(c.Request().Context(), map[string]interface{}{"_id": userID})
	if result.DeletedCount != 1 {
		return c.String(http.StatusNotFound, "Could not delete the user")
	}
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, "Account deleted")
}
