package main

import (
	//"fmt"
	"fmt"
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
	request := strings.Split(strings.Split(string(buff), "\n")[0], " ")
	var response []byte
	if len(request[1]) > 5 && request[1][:6] == "/echo/" {
		str := request[1][6:]
		length := len(str)
		rsp := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s\r\n", length, str)
		response = []byte(rsp)

	} else if request[1] == "/" {
		response = []byte("HTTP/1.1 200 OK\r\n\r\n")
	} else {
		response = []byte("HTTP/1.1 404 Not Found\r\n\r\n")
	}
	_, err = conn.Write(response)
	if err != nil {
		return
	}
}
