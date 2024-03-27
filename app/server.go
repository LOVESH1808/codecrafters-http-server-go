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
	headers := make(map[string]string)
	request := strings.Split(string(buff), "\n")
	head := strings.Split(request[0], " ")
	for _, val := range request[1:] {
		if !strings.Contains(val, ":") {
			break
		}
		print(strings.Split(val, " ")[0])
		headers[strings.Split(val, " ")[0]] = strings.Trim(strings.Split(val, " ")[1], "\r")

	}
	var response []byte
	if len(head[1]) > 5 && head[1][:6] == "/echo/" {
		str := head[1][6:]
		length := len(str)
		rsp := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s\r\n", length, str)
		response = []byte(rsp)
	} else if head[1] == "/user-agent" {
		str := headers["User-Agent:"]
		length := len(str)
		rsp := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s\r\n", length, str)

		response = []byte(rsp)
	} else if head[1] == "/" {
		response = []byte("HTTP/1.1 200 OK\r\n\r\n")
	} else {
		response = []byte("HTTP/1.1 404 Not Found\r\n\r\n")
	}
	_, err = conn.Write(response)
	if err != nil {
		return
	}
}
