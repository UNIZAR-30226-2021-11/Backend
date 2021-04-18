package main

import (
	"Backend/pkg/events"
	"bufio"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"os"
	"time"
)

func main() {
	u := url.URL{Scheme: "ws", Host: ":9000", Path: "/simulation"}

	// Establish connection
	c, _, _ := websocket.DefaultDialer.Dial(u.String(), nil)
	defer c.Close()

	// Receive messages
	go func() {
		for {
			_, message, _ := c.ReadMessage()
			log.Printf("Message received: %s", message)
		}
	}()

	//Create a game
	time.Sleep(2 * time.Second)
	event := events.Event{
		GameID:    1,
		PlayerID:  1,
		EventType: 0,
	}
	_  = c.WriteJSON(event)

	var i uint32
	//Join the game
	for i = 1; i < 5; i++ {
			time.Sleep(2 * time.Second)
			event := events.Event{
				GameID:    1,
				PlayerID:  i,
				EventType: 1,
			}
			_  = c.WriteJSON(event)
	}


	// Read from stdin and send through websocket
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		_ = c.WriteMessage(websocket.TextMessage, []byte(scanner.Text()))
		log.Printf("Message sent: %s", scanner.Text())
	}
}
