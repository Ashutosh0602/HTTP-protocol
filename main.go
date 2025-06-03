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

var routes = map[string]func(net.Conn){
	"/api/hello":   handleHello,
	"/api/goodbye": handleGoodbye,
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
	// conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\nHello from TCP server!\n"))
	fmt.Println("Headers:", headers)

	// Handle the request based on the path
	handleRequest(conn, path)

	// -------------- Day 2: Ends ----------------

	html := `
<html>
	<head><title>Mini Go Server</title></head>
	<body><h1>Hello from Go!</h1></body>
</html>
`

	switch path {
	case "/":
		writeResponse(conn, 200, "text/plain", "Welcome to the home page!")
	case "/about":
		writeResponse(conn, 200, "text/plain", "This is the about page.")
	case "/home":
		writeResponse(conn, 200, "text/plain", "Welcome to the home page!")
	case "/html":
		writeResponse(conn, 200, "text/html", html)
	case "/json":
		writeResponse(conn, 200, "application/json", `{"message": "Hello from Go!"}`)
	default:
		// writeResponse(conn, 404, "text/plain", "Page not found.")
		serveStaticFile(conn, path)
	}

}

// -------------- Day 3: Start ----------------
func writeResponse(conn net.Conn, statusCode int, contentType, body string) {
	statusText := map[int]string{
		200: "OK",
		404: "Not Found",
		500: "Internal Server Error",
	}[statusCode]

	if statusText == "" {
		statusText = "Unknown Status"
	}

	statusLine := fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, statusText)
	headers := fmt.Sprintf("Content-Type: %s\r\nContent-Length: %d\r\n\r\n", contentType, len(body))

	response := statusLine + headers + body
	_, err := conn.Write([]byte(response))
	if err != nil {
		fmt.Println("Error writing response:", err)
		return
	}
	fmt.Println("Response sent:", response)
}

// -------------- Day 3: Ends ----------------

// -------------- Day 4: Start ----------------
func serveStaticFile(conn net.Conn, path string) {
	filePath := "./public" + path
	data, err := os.ReadFile(filePath)
	if err != nil {
		writeResponse(conn, 404, "text/plain", "404 - File not found")
		return
	}
	contentType := getContentType(filePath)
	writeResponse(conn, 200, contentType, string(data))
}

func getContentType(filePath string) string {
	switch {
	case strings.HasSuffix(filePath, ".html"):
		return "text/html"
	case strings.HasSuffix(filePath, ".css"):
		return "text/css"
	case strings.HasSuffix(filePath, ".js"):
		return "application/javascript"
	case strings.HasSuffix(filePath, ".json"):
		return "application/json"
	case strings.HasSuffix(filePath, ".png"):
		return "image/png"
	case strings.HasSuffix(filePath, ".jpg") || strings.HasSuffix(filePath, ".jpeg"):
		return "image/jpeg"
	default:
		return "application/octet-stream"
	}
}

func handleRequest(conn net.Conn, path string) {
	if handler, exists := routes[path]; exists {
		handler(conn)
	} else {
		writeResponse(conn, 404, "text/html", "404 - Not Found")
	}
}
func handleHello(conn net.Conn) {
	response := "Hello from the /api/hello endpoint!\n"
	writeResponse(conn, 200, "text/plain", response)
}
func handleGoodbye(conn net.Conn) {
	response := "Goodbye from the /api/goodbye endpoint!\n"
	writeResponse(conn, 200, "text/plain", response)
}
