package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
)

const (
  serverAddress = "0.0.0.0"
  serverPort = "4221"
)

func httpResponse(statusCode int, optionalMessage string) []byte {
  // TODO: Proper response HERE
  responseString := fmt.Sprintf("HTTP/1.1 %d %s\r\n\r\n", statusCode, optionalMessage)
  return []byte(responseString)
}

// TODO: Finish
func httpResponseWithData(statusCode int, optionalMessage, headers, data string) []byte {
  responseString := fmt.Sprintf("HTTP/1.1 %d %s\r\n%s\r\n%s", statusCode, optionalMessage, headers, data)
  return []byte(responseString)
}

func main() {
  fmt.Printf("Serving on %s:%s\n", serverAddress, serverPort)
  l, err := net.Listen("tcp", serverAddress + ":" + serverPort)
	if err != nil {
		fmt.Println("Failed to bind to port " + serverPort)
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
  requestString := string(requestBuffer)
  parsedRequest := strings.Split(requestString, " ") // In format [info, headers, body]

  fmt.Println(requestString)
  
  // Checks for echo
  // TODO: Make an endpoint handler (map that refers to function?)
  if strings.HasPrefix(requestString, "GET /echo/") {
    // Status is ok
    // Convoluted code that works
    // Route data is the data of the root
    routeData := strings.Split(parsedRequest[1], "/")
    if len(routeData) > 1 {
      dataToEcho := routeData[2]
      responseHeaders := fmt.Sprintf("Content-Type: %s\r\nContent-Length: %d\r\n", "text/plain", len(routeData[2]))
      connection.Write(httpResponseWithData(http.StatusOK, "OK", responseHeaders, dataToEcho))
      fmt.Println("Correct link")
    } else {
      connection.Write(httpResponse(http.StatusNotFound, "Not Found"))
      fmt.Println("Incorrect link")
    }
  } else {
    // Status is incorrect
    connection.Write(httpResponse(http.StatusNotFound, "Not Found"))
    fmt.Println("Incorrect link")
  }

  fmt.Println("Responded")
}
