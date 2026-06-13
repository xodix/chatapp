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
	chatsRouter.POST("/", createChat)
	chatsRouter.GET("/messages", getMessages)
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

	if len(chats) == 0 {
		return c.JSON(http.StatusOK, GetUsersChatsResponse{Chats: []models.Chat{}})
	}

	return c.JSON(http.StatusOK, GetUsersChatsResponse{Chats: chats})
}

type CreateChatRequest struct {
	Name    string   `json:"name" validate:"required,min=3,max=60"`
	Members []string `json:"members" validate:"required,min=1"`
}

func createChat(c *echo.Context) error {
	var request CreateChatRequest
	err := c.Bind(&request)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
	}

	ids := []bson.ObjectID{}
	for _, id := range request.Members {
		idObj, err := bson.ObjectIDFromHex(id)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
		}
		ids = append(ids, idObj)
	}

	chatsCollection.InsertOne(c.Request().Context(), models.Chat{
		Name:     request.Name,
		Messages: []models.ChatMessage{},
		Members:  ids,
	})

	return c.String(http.StatusOK, "Chat created")
}

func getMessages(c *echo.Context) error {
	id := c.Get("userID").(string)
	userID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	chatID := c.QueryParam("id")
	objID, err := bson.ObjectIDFromHex(chatID)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	result := chatsCollection.FindOne(c.Request().Context(), map[string]interface{}{"_id": objID, "members": userID})
	if result.Err() != nil {
		return c.String(http.StatusBadRequest, result.Err().Error())
	}
	var chat models.Chat
	err = result.Decode(&chat)

	return c.JSON(200, chat.Messages)
}
