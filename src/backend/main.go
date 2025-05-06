package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// Configure the upgrader
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	// Start sending messages every 400ms
	ticker := time.NewTicker(400 * time.Millisecond)
	defer ticker.Stop()

	done := make(chan struct{})

	// Read loop
	go func() {
		for {
			messageType, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println("Read error:", err)
				close(done)
				return
			}
			log.Printf("Received: %s\n", msg)

			response := fmt.Sprintf("Echo: %s", msg)
			if err = conn.WriteMessage(messageType, []byte(response)); err != nil {
				log.Println("Write error:", err)
				close(done)
				return
			}
		}
	}()

	// Write loop
	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			message := fmt.Sprintf("Tick at %s", t.Format("15:04:05.000"))
			if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
				log.Println("Tick send error:", err)
				return
			}
		}
	}
}

func main() {
	http.HandleFunc("/ws", wsHandler)
	fmt.Println("Server started on :4000")
	log.Fatal(http.ListenAndServe("127.0.0.1:4000", nil))
}
