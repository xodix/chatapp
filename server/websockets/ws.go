package websockets

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v5"
)

var hub = NewHub()

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

/*
	TODO:
	1. Get userID (only for the author in message Response)
	2. Get GroupID
	4. Get chatID
	5. Authenticate if the user can Have access to the Chat
	5. Send only if the message is to the same chat
*/

func WebsocketsHandler(c *echo.Context) error {
	userID := c.Get("userID").(string)

	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := hub.Register(userID, conn)
	defer hub.Unregister(userID)

	// Read loop — blocks until client disconnects
	for {
		mt, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}

		// Example: broadcast received message to a target list
		recipients := []string{"user-2", "user-3"}
		hub.SendToIDs(recipients, mt, msg)
	}

	return nil
}
