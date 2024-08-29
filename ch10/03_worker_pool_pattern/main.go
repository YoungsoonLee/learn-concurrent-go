package main

import (
	"fmt"
	"net"
	"os"
	"regexp"
)

var r, _ = regexp.Compile("GET (.+) HTTP/1.1\r\n")

func handleHTTPRequest(conn net.Conn) {
	buff := make([]byte, 1024)
	size, _ := conn.Read(buff)

	if r.Match(buff[:size]) {
		file, err := os.ReadFile(fmt.Sprintf("../resource/%s", r.FindSubmatch(buff[:size])[1])) // If the request is a valid one, reads the request file from the resources directory
		if err == nil {                                                                         // If the file exists, responds to the client with the HTTP header and the file content
			conn.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Length: %d\r\n\r\n", len(file))))
			conn.Write(file)
		} else { // If the file does not exist, responds with an error
			conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n<html>Not Found</html>"))
		}
	} else {
		conn.Write([]byte("HTTP/1.1 400 Bad Request\r\n\r\n<html>Bad Request</html>"))
	}
	conn.Close()
}

func startHTTPWorkers(n int, incomingConnections <-chan net.Conn) {
	for i := 0; i < n; i++ {
		go func() {
			for conn := range incomingConnections {
				handleHTTPRequest(conn)
			}
		}()
	}
}

func main() {
	incomingConnections := make(chan net.Conn)
	startHTTPWorkers(3, incomingConnections) // Starts the worker pool with three goroutines

	server, _ := net.Listen("tcp", ":8080")
	defer server.Close()

	for {
		conn, _ := server.Accept()

		//incomingConnections <- conn // Passes the connection on the work queue channel

		select {
		case incomingConnections <- conn:
		default:
			fmt.Println("Server is busy")
			conn.Write([]byte("HTTP/1.1 429 Too Many Requests\r\n\r\n<html>Busy</html>\n"))
			conn.Close()
		}
	}
}
