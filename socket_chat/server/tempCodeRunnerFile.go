package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	listen, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer listen.Close()

	conn, err := listen.Accept()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	for {
		msg, _ := reader.ReadString('\n')
		fmt.Println("Received:", msg)
		writer.WriteString("Server received: " + msg)
		writer.Flush()
	}
}
