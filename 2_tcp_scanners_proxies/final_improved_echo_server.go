package main

import (
	"io"
	"log"
	"net"
)

// echo is a handler function that simply echoes received data
func echo(conn net.Conn) {
	defer conn.Close()

	// we are just echoing here, we're not logging anything
	for {
		if _, err := io.Copy(conn, conn); err != nil {
			log.Fatalln("Unable to read/write data")
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
