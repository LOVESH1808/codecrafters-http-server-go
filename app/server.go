package main

import (
	//"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	server, err := net.Listen("tcp", "localhost:4221") // change this for the real one
	if err != nil {
		os.Exit(1)
	}

	var conn, err2 = server.Accept()
	if err2 != nil {
		return
	}

	defer server.Close()
	buff := make([]byte, 1024)
	_, err = conn.Read(buff)
	if err != nil {
		return
	}

	defer conn.Close()
	request := string(buff)
	path := strings.Split(request, " ")[1]
	if path == "/" {
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	} else {
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))

	}
}
