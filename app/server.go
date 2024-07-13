package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
)

func httpResponse(statusCode int, optionalMessage string) []byte {
  // TODO: Proper response HERE
  responseString := fmt.Sprintf("HTTP/1.1 %d %s\r\n\r\n", statusCode, optionalMessage)
  return []byte(responseString)
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

  fmt.Println("Listening...")
	
  connection, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

  defer connection.Close() // Adds the connection close to the stack
  defer fmt.Println("Closing connection")
  defer l.Close() // Closes the listener after stack is finished
  defer fmt.Println("Closing server")
  
  fmt.Println("Accepted")
  
  requestBuffer := make([]byte, 4096) // big buffer

  connection.Read(requestBuffer)
  if strings.HasPrefix(string(requestBuffer), "GET / HTTP/1.1") {
    connection.Write(httpResponse(http.StatusOK, "OK"))
    fmt.Println("Correct link")
  } else {
    connection.Write(httpResponse(http.StatusNotFound, "Not found"))
    fmt.Println("Incorrect link")
  }

  fmt.Println("Responded")
}
