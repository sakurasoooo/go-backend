package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
)

type Message struct {
	Username   string `json:"username"`
	MsgType    string `json:"msgType"`
	MsgContent string `json:"msgContent"`
}

type Client struct {
	Conn     net.Conn
	Username string
}

func main() {
	listener, _ := net.Listen("tcp", "localhost:8080")
	defer listener.Close()

	clients := make(map[net.Conn]Client)

	for {
		conn, _ := listener.Accept()
		go handleClient(conn, clients)
	}
}

func handleClient(conn net.Conn, clients map[net.Conn]Client) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		var msg Message
		err := json.NewDecoder(reader).Decode(&msg)
		if err != nil {
			fmt.Println("Client disconnected")
			delete(clients, conn)
			break
		}

		if _, ok := clients[conn]; !ok {
			clients[conn] = Client{Conn: conn, Username: msg.Username}
		}

		if msg.MsgType == "chat" {
			for _, client := range clients {
				json.NewEncoder(client.Conn).Encode(msg)
			}
		}

		fmt.Printf("Received from client %s: %s\n", msg.Username, msg.MsgContent)
	}
}
