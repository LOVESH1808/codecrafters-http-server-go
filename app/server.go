package main

import (
	//"fmt"
	"net"
	"os"
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
	_, err = conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	if err != nil {

		return
	}
}
