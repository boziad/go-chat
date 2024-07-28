package main

import (
	"fmt"
	"net"
	"sync"
)

var (
	connections = make(map[net.Conn]bool)
	mutex       sync.Mutex
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server is listening on port 8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		mutex.Lock()
		connections[conn] = true
		mutex.Unlock()
		fmt.Println("New client connected")
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer func() {
		mutex.Lock()
		delete(connections, conn)
		mutex.Unlock()
		conn.Close()
	}()

	buffer := make([]byte, 1024)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		message := buffer[:n]
		fmt.Printf("Received: %s\n", message)
		broadcastMessage(message, conn)
	}
}

func broadcastMessage(message []byte, sender net.Conn) {
	mutex.Lock()
	defer mutex.Unlock()
	for conn := range connections {
		if conn != sender {
			_, err := conn.Write(message)
			if err != nil {
				fmt.Println("Error sending message:", err)
				conn.Close()
				delete(connections, conn)
			}
		}
	}
}
