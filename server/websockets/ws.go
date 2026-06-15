package websockets

import (
	"encoding/json"
	"log"
	"net/http"
	"study4cash/DB/models"
	"study4cash/auth"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var hub = NewHub()

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var chatsCollection *mongo.Collection
var usersCollection *mongo.Collection

// TODO: This is a MESS
func WebsocketsHandler(c *echo.Context, db *mongo.Database) error {
	log.Println("WEBSOCKET CONNECTION REQUEST")
	if chatsCollection == nil || usersCollection == nil {
		chatsCollection = db.Collection("chats")
		usersCollection = db.Collection("users")
	}

	// Parse Query parameters
	token := c.QueryParam("token")
	payload, err := auth.ValidateJWT(token)
	if err != nil {
		log.Println(err.Error())
		return c.String(http.StatusUnauthorized, "Invalid token")
	}
	userID := payload.UserID
	userid, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err.Error())
		return c.String(http.StatusBadRequest, "user_id and chat_id are required")
	}

	chatID := c.QueryParam("chat_id")
	chatid, err := bson.ObjectIDFromHex(chatID)
	if err != nil {
		log.Println(err.Error())
		return c.String(http.StatusBadRequest, "user_id and chat_id are required")
	}

	// Check if user can access the chat
	result := chatsCollection.FindOne(c.Request().Context(), bson.M{"_id": chatid, "members": userid})
	if result.Err() != nil {
		return c.String(http.StatusUnauthorized, "Unauthorized")
	}

	// Get user data for name and surname
	userData := models.User{}
	result = usersCollection.FindOne(c.Request().Context(), bson.M{"_id": userid})
	if result.Err() != nil {
		log.Println(result.Err().Error())
		return c.String(http.StatusInternalServerError, result.Err().Error())
	}
	err = result.Decode(&userData)
	if err != nil {
		log.Println(result.Err().Error())
		return c.String(http.StatusInternalServerError, result.Err().Error())
	}

	// Upgrade connection to Websocket protocol
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Println("failed to upgrade connection", err.Error())
		return c.String(http.StatusInternalServerError, "failed to upgrade connection")
	}
	defer conn.Close()

	// Register the user in connections storage
	_ = hub.Register(chatID, userID, conn)
	defer hub.Unregister(chatID, userID)

	// Read loop
	for {
		// Read the actual message
		mt, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Stupid stuff", err.Error())
			break
		}

		// Put the message into db
		msgStr := string(msg)
		msgStruct := models.ChatMessage{
			Author:  userid,
			Message: msgStr,
			Name:    userData.Name,
			Surname: userData.Surname,
		}
		filter := bson.M{"_id": chatid}
		update := bson.M{
			"$push": bson.M{
				"messages": msgStruct,
			},
		}
		result, err := chatsCollection.UpdateOne(c.Request().Context(), filter, update)
		if err != nil {
			log.Println(err.Error())
		}
		log.Println(result.ModifiedCount)

		// Send the message to all relevant users
		jsonData, err := json.Marshal(msgStruct)
		if err != nil {
			log.Println(err.Error())
		}
		log.Print(string(jsonData))
		hub.SendToIDs(chatID, mt, jsonData)
	}

	return nil
}
