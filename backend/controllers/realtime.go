package controllers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func TicketUpdates(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("Failed to set WebSocket upgrade:", err)
		return
	}
	defer conn.Close()

	for {
		// Send real-time ticket updates
		// Use a channel or database polling mechanism to send updates here
		ticketUpdate := map[string]interface{}{
			"event_id":          1,
			"available_tickets": 100,
		}
		if err := conn.WriteJSON(ticketUpdate); err != nil {
			fmt.Println("Error writing to WebSocket:", err)
			break
		}
	}
}
