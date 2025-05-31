package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting TCP server:", err)
		os.Exit(1)
	}

	defer listener.Close()
	fmt.Println("TCP server listening on :8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// -------------- Day 1: Start ----------------
	// buffer := make([]byte, 10)

	// for {
	// 	n, err := conn.Read(buffer)
	// 	if err != nil {
	// 		fmt.Println("Error reading from connection:", err)
	// 		return
	// 	}

	// 	request := string(buffer[:n])
	// 	fmt.Println("Received request:", request)
	// 	response := "Hello from TCP server!\n"
	// 	_, err = conn.Write([]byte(response))
	// }

	// -------------- Day 1: Ends ----------------

	// -------------- Day 2: Start ----------------
	reader := bufio.NewReader(conn)
	requestLine, err := reader.ReadString('\n')

	if err != nil {
		fmt.Println("Error reading from connection:", err)
		return
	}
	parts := strings.Split(strings.TrimSpace(requestLine), " ")
	if len(parts) != 3 {
		fmt.Println("Malformed request line")
		return
	}
	method, path, version := parts[0], parts[1], parts[2]
	fmt.Println("Method:", method)
	fmt.Println("Path:", path)
	fmt.Println("Version:", version)

	headers := make(map[string]string)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading header:", err)
			return
		}
		line = strings.TrimSpace(line)
		if line == "" {
			break // End of headers
		}

		headerParts := strings.SplitN(line, ":", 2)
		if len(headerParts) == 2 {
			headers[strings.TrimSpace(headerParts[0])] = strings.TrimSpace(headerParts[1])
		}

	}
	conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\nHello from TCP server!\n"))
	fmt.Println("Headers:", headers)

}
