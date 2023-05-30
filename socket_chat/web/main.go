package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil) // Error is not ignored anymore
		if err != nil {
			log.Printf("Failed to set websocket upgrade: %+v", err)
			return
		}

		for {
			// Read message from browser
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				log.Printf("Error during message reading: %+v", err)
				return
			}

			// Print the message to the console
			fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(msg))

			// Write message back to browser
			if err = conn.WriteMessage(msgType, msg); err != nil {
				log.Printf("Error during message writing: %+v", err)
				return
			}
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "websockets.html")
	})

	http.ListenAndServe(":8080", nil)
}
