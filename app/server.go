package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

func buildResponse(code, headers, body string) []byte {
	clrf := "\r\n"
	code200 := "HTTP/1.1 200 OK\r\n"
	code201 := "HTTP/1.1 201 OK\r\n"
	code404 := "HTTP/1.1 404 Not Found\r\n"
	code500 := "HTTP/1.1 500 Internal Server Error\r\n"
	if headers == "" && body == "" {
		switch code {
		case "200":
			return []byte(code200 + clrf)
		case "201":
			return []byte(code201 + clrf)
		case "404":
			return []byte(code404 + clrf)
		case "500":
			return []byte(code500 + clrf)
		}
	}
	switch code {
	case "200":
		return []byte(code200 + headers + clrf + clrf + body)
	case "201":
		return []byte(code201 + headers + clrf + clrf + body)
	default:
		return nil
	}

}
func handleRequest(conn net.Conn, filesDir string) {
	defer conn.Close()
	buff := make([]byte, 1024)
	_, err := conn.Read(buff)
	if err != nil {
		fmt.Println("Error reading from connection: ", err.Error())
		os.Exit(1)
	}
	var res []byte
	path := strings.Split(string(buff), " ")[1]
	userAgent := strings.TrimPrefix(strings.Split(string(buff), "\r\n")[2], "User-Agent: ")
	if path == "/" {
		res = buildResponse("200", "", "")
	} else if strings.Contains(path, "echo") {
		body := strings.TrimPrefix(path, "/echo/")
		header := fmt.Sprintf("Content-Type: text/plain\r\nContent-Length: %d", len(body))

		res = buildResponse("200", header, body)
	} else if path == "/user-agent" {
		header := fmt.Sprintf("Content-Type: text/plain\r\nContent-Length: %d", len(userAgent))

		res = buildResponse("200", header, userAgent)
	} else if strings.Contains(path, "files") {
		filePath := strings.TrimPrefix(path, "/files/")
		method := strings.Split(string(buff), " ")[0]
		if method == "GET" {
			if file, err := os.ReadFile(filesDir + filePath); err == nil {
				content := string(file)
				header := fmt.Sprintf("Content-Type: application/octet-stream\r\nContent-Length: %d", len(content))
				res = buildResponse("200", header, content)
			} else {
				res = buildResponse("404", "", "")
			}
		} else if method == "POST" {
			body := strings.Split(string(buff), "\r\n\r\n")[1]
			err := os.WriteFile(filesDir+filePath, []byte(strings.Trim(body, "\x00")), 0644)
			if err != nil {
				res = buildResponse("500", "", "")
			} else {
				res = buildResponse("201", "", "")
			}

		}
	} else {
		res = buildResponse("404", "", "")
	}
	_, err = conn.Write(res)
	if err != nil {
		fmt.Println("Error writing to connection: ", err.Error())
		os.Exit(1)
	}
}
func main() {
	fmt.Println("Logs from your program will appear here!")
	filesDir := flag.String("directory", "", "Directory to serve files from")
	flag.Parse()
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handleRequest(conn, *filesDir)
	}
}
