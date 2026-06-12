package routes

import (
	"net/http"
	"study4cash/DB/models"
	"study4cash/auth"

	"github.com/labstack/echo/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var chatsCollection *mongo.Collection

func RouteChats(prefix string, e *echo.Echo, db *mongo.Database) {
	chatsCollection = db.Collection("chats")
	// Entire group is protected
	chatsRouter := e.Group(prefix)
	chatsRouter.Use(auth.JWTMiddleware)

	chatsRouter.GET("/", getUsersChats)
	// chatsRouter.DELETE("/", DeleteAccount)
	// chatsRouter.PUT("/setName", UpdateAccount)
}

type GetUsersChatsResponse struct {
	Chats []models.Chat `json:"chats"`
}

func getUsersChats(c *echo.Context) error {
	userID := c.Get("userID").(string)
	userIDObj, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	result, err := chatsCollection.Find(c.Request().Context(), map[string]interface{}{"members": userIDObj})
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	var chats []models.Chat
	err = result.All(c.Request().Context(), &chats)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, GetUsersChatsResponse{Chats: chats})
}
