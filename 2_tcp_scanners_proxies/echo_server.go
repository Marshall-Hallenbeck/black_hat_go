package main

import (
	"io"
	"log"
	"net"
)

// echo is a handler function that simply echoes received data
func echo(conn net.Conn) {
	defer conn.Close()

	for {
		// Create a buffer to store received data
		// We define the buffer inside the for loop because it's not cleared between requests
		b := make([]byte, 512)
		// Receive data via conn.Read into a buffer
		size, err := conn.Read(b[0:])
		if err == io.EOF {
			log.Println("Client disconnected")
			break
		}
		if err != nil {
			log.Println("Unexpected error")
			break
		}
		// Getting a new line because this is how the logger function works
		// https://golang.org/pkg/log/#Logger.Output
		log.Printf("Received %d bytes: %s", size, string(b))

		// Send data via conn.Write
		log.Print("Writing data")
		if _, err := conn.Write(b[0:size]); err != nil {
			log.Fatalln("Unable to write data")
		}
	}
}

func main() {
	// Bind to TCP port 20080 on all interfaces
	listener, err := net.Listen("tcp", ":20080")
	if err != nil {
		log.Fatalln("Unable to bind to port")
	}
	log.Print("Listening on 0.0.0.0:20080")
	for {
		// Wait for connection and create net.Conn on connection established
		conn, err := listener.Accept()
		log.Print("Received connection")
		if err != nil {
			log.Fatalln("Unable to accept connection")
		}
		// Handle the connection using goroutine for concurrency
		go echo(conn)
	}
}
