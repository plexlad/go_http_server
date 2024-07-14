package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
//  "strings"
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
//  requestString := string(requestBuffer)
//  parsedRequest := strings.Split(requestString, " ") // In format [info, headers, body]

  // Create a new handler as well as a request to handle
  handler := new(HttpHandler)
  request := ParseHttpRequest(requestBuffer, connection)

  // Returns OK status if there is no route data
  // Register the server route
  handler.RegisterRoute("/", func(request HttpRequest) error {
    if request.RouteData == "" && request.Method == "GET" {
      request.Connection.Write(httpResponse(http.StatusOK, "OK"))
      fmt.Println("Valid link")
      return nil
    }

    request.Connection.Write(httpResponse(http.StatusNotFound, "Not Found"))
    fmt.Println("Invalid link")
    return nil
  })

  // Register echo server
  handler.RegisterRoute("/echo/", func(request HttpRequest) error {
    if request.Method == "GET" {
      // Create headers
      responseHeader := fmt.Sprintf("Content-Type: %s\r\nContent-Length: %d\r\n", "text/plain", len(request.RouteData))
      // Send the response
      request.Connection.Write(httpResponseWithData(http.StatusOK, "OK", responseHeader, request.RouteData))
      fmt.Println("Valid link")
      return nil
    }

    // Send the incorrect response if method is not GET
    request.Connection.Write(httpResponse(http.StatusNotFound, "Not Found"))
    fmt.Println("Invalid link")
    return nil
  })

  // Registers the user agent header example endpoint
  // Sends the user agent back to the user
  handler.RegisterRoute("/user-agent/", func(request HttpRequest) error {
    if request.Method == "GET" {
      responseHeader := fmt.Sprintf("Content-Type: %s\r\nContent-Length: %d\r\n", "text/plain", len(request.RouteData))
      request.Connection.Write(httpResponseWithData(http.StatusOK, "OK", responseHeader, request.Header["User-Agent"]))
      fmt.Println("user-agent: Valid link")
      return nil
    }
    
    request.Connection.Write(httpResponse(http.StatusNotFound, "Not Found"))
    fmt.Println("user-agent: Invalid link")
    return nil
  })

  // Handle the request with the handler!
  err = handler.HandleRequest(request)
  if err != nil {
    fmt.Println("Error: " + err.Error())
    return
  }

  fmt.Println("Responded")
}
