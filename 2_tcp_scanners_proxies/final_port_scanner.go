package main

import (
	"fmt"
	"net"
	"sort"
)

//noinspection GoDuplicate
func worker(ports, results chan int) {
	for p := range ports {
		address := fmt.Sprintf("scanme.nmap.org:%d", p)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			// Send a zero if the port is closed
			results <- 0
			continue
		}
		conn.Close()
		results <- p
	}
}

func main() {
	ports := make(chan int, 100)
	// Create a separate channel to communicate the results from the worker to the main thread
	results := make(chan int)
	var openPorts []int
	for i := 0; i < cap(ports); i++ {
		go worker(ports, results)
	}
	// Send to the workers in a separate goroutine
	go func() {
		for i := 1;i <= 1024; i++ {
			ports <- i
		}
	}()
	// Results gathering goroutine
	for i := 0; i < 1024; i++ {
		port := <-results
		if port != 0 {
			openPorts = append(openPorts, port)
		}
	}

	close(ports)
	close(results)
	// Sort the results for better printing
	sort.Ints(openPorts)
	for _, port := range openPorts {
		fmt.Printf("%d open\n", port)
	}
}
