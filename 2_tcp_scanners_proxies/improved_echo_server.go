package main

import (
	"bufio"
	"log"
	"net"
)

// echo is a handler function that simply echoes received data
func echo(conn net.Conn) {
	defer conn.Close()

	for {
		reader := bufio.NewReader(conn)
		s, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalln("Unable to read data")
		}
		// This doesn't have the same problem with the extra newline as echo_server does
		log.Printf("Read %d bytes: %s", len(s), s)

		log.Println("Writing data")
		writer := bufio.NewWriter(conn)
		if _, err := writer.WriteString(s); err != nil {
			log.Fatalln("Unable to write data")
		}
		writer.Flush()
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
