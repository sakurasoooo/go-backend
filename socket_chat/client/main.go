package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"
)

type Message struct {
	Username   string `json:"username"`
	MsgType    string `json:"msgType"`
	MsgContent string `json:"msgContent"`
}

func readHandler(conn net.Conn, username string, messageChan chan Message, quitChan chan bool) {
	serverReader := bufio.NewReader(conn)
	for {
		var msg Message
		err := json.NewDecoder(serverReader).Decode(&msg)
		if err != nil {
			fmt.Println("Server closed connection. Attempting to reconnect...")
			quitChan <- true
			conn = retryConnect(conn, username, messageChan)
			serverReader = bufio.NewReader(conn)
		}
		fmt.Printf("Received from %s: %s\n", msg.Username, msg.MsgContent)
	}
}

func sendHeartbeat(conn net.Conn, username string) {
	ticker := time.NewTicker(10 * time.Second)
	for range ticker.C {
		msg := Message{
			Username:   username,
			MsgType:    "heartbeat",
			MsgContent: "ping",
		}
		json.NewEncoder(conn).Encode(msg)
	}
}

func retryConnect(conn net.Conn, username string, messageChan chan Message) net.Conn {
	var err error
	for i := 0; i < 3; i++ {
		conn, err = net.Dial("tcp", "localhost:8080")
		if err != nil {
			fmt.Printf("Reconnection attempt %d failed. Retrying in 10 seconds...\n", i+1)
			time.Sleep(10 * time.Second)
		} else {
			go sendHeartbeat(conn, username)
			quitChan := make(chan bool)
			go sendMessage(conn, messageChan, quitChan)
			go readHandler(conn, username, messageChan, quitChan)
			return conn
		}
	}
	fmt.Println("Failed to reconnect to server after 3 attempts. Exiting...")
	os.Exit(1)
	return nil
}

func sendMessage(conn net.Conn, messageChan chan Message, quitChan chan bool) {
	for {
		select {
		case msg := <-messageChan:
			json.NewEncoder(conn).Encode(msg)
		case <-quitChan:
			return
		}
	}
}

func main() {
	fmt.Print("Enter your username: ")
	reader := bufio.NewReader(os.Stdin)
	username, _ := reader.ReadString('\n')
	username = username[:len(username)-1] // remove trailing newline

	var conn net.Conn
	var err error
	for i := 0; i < 3; i++ {
		conn, err = net.Dial("tcp", "localhost:8080")
		if err != nil {
			fmt.Println("Failed to connect to server, retrying in 10 seconds...")
			time.Sleep(10 * time.Second)
		} else {
			break
		}
	}

	if err != nil {
		fmt.Println("Failed to connect to server after 3 attempts. Exiting...")
		os.Exit(1)
	}

	defer conn.Close()
	messageChan := make(chan Message)
	quitChan := make(chan bool)

	go sendHeartbeat(conn, username)
	go readHandler(conn, username, messageChan, quitChan)
	go sendMessage(conn, messageChan, quitChan)

	for {
		fmt.Print("Enter text: ")
		text, _ := reader.ReadString('\n')
		msg := Message{
			Username:   username,
			MsgType:    "chat",
			MsgContent: text,
		}
		messageChan <- msg
	}
}
