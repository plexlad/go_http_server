package main

import (
	"fmt"
	"net/http"
	"net"
	"os"
)

func httpResponse(version string, statusCode int, optionalMessage string) []byte {
  // TODO: Proper response HERE
  responseString := fmt.Sprintf("HTTP/%s %d %s\r\n\r\n", version, statusCode, optionalMessage)
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
	
  conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
  
  fmt.Println("Accepted")

  conn.Write(httpResponse("1.1", http.StatusOK, "OK"))

  fmt.Println("Responded")
}
