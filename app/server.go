package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	fmt.Println("Listening on port 4221")
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handleConnections(conn)
	}
}
func handleConnections(conn net.Conn) {
	fmt.Println("Established connection with", conn.RemoteAddr())
	defer conn.Close()
	buff := make([]byte, 1024)
	_, err := conn.Read(buff)
	if err != nil {
		fmt.Println(err.Error())
	}
	request := string(buff)
	handleRequest(request, conn)
}
func handleRequest(request string, conn net.Conn) {
	path := strings.Split(request, " ")[1]
	switch {
	case path == "/":
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	case strings.HasPrefix(path, "/echo/"):
		text := strings.Split(path, "/echo/")[1]
		response := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %v\r\n\r\n%v\r\n\r\n", len(text), text)
		conn.Write([]byte(response))
	case path == "/user-agent":
		temp := strings.Split(request, "User-Agent: ")[1]
		user_agent := strings.Split(temp, "\r")[0]
		response := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %v\r\n\r\n%v\r\n\r\n", len(user_agent), user_agent)
		conn.Write([]byte(response))
	default:
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}
}
