package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

const (
  serverAddress = "0.0.0.0"
  serverPort = "4221"
)

func httpResponse(statusCode int, optionalMessage string) []byte {
  responseString := fmt.Sprintf("HTTP/1.1 %d %s\r\n\r\n", statusCode, optionalMessage)
  return []byte(responseString)
}

func httpResponseWithData(statusCode int, optionalMessage, headers, data string) []byte {
  responseString := fmt.Sprintf("HTTP/1.1 %d %s\r\n%s\r\n%s", statusCode, optionalMessage, headers, data)
  return []byte(responseString)
}

func main() {
  // TODO: Add get and post handlers! Make the syntax way shorter
  server, err := NewServer("0.0.0.0", "4221")
  if err != nil {
    fmt.Println(err.Error())
    os.Exit(1)
  }

  // Returns OK status if there is no route data
  // Register the server route
  server.routeHandler.RegisterRoute("/", func(request HttpRequest) error {
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
  server.RegisterRoute("/echo/", func(request HttpRequest) error {
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

  server.RegisterRoute("/files/", func(request HttpRequest) error {
    if request.Method == "GET" && request.RouteData != "" {
      // Serves files from a directory
      dir := os.Args[2]
      fileAddress := fmt.Sprintf("%s%s", dir, request.RouteData)
      // Reads data from the file
      data, err := os.ReadFile(fileAddress)
      if err != nil {
        fmt.Println("Issue serving file, error: ", err.Error())
        request.Connection.Write(httpResponse(http.StatusNotFound, "Not Found"))
        return nil
      }

      responseHeader := fmt.Sprintf("Content-Type: %s\r\nContent-Length: %d\r\n", "application/octet-stream", len(data))
      request.Connection.Write(httpResponseWithData(http.StatusOK, "OK", responseHeader, string(data)))
      fmt.Println("file: Valid link")
      return nil

    } else if request.Method == "POST" && request.RouteData != "" {
      // Creates files from a directory
      fmt.Println("POST")
      dir := os.Args[2]
      fileAddress := fmt.Sprintf("%s%s", dir, request.RouteData)
      err := os.WriteFile(fileAddress, []byte(strings.Trim(request.Body, "\x00")), 0644)
      if err != nil {
        fmt.Println("Issue writing file, error: ", err.Error())
        request.Connection.Write(httpResponse(http.StatusInternalServerError, "File Write Error"))
        return nil
      }
      request.Connection.Write(httpResponse(http.StatusCreated, "Created"))
      return nil
    }

    request.Connection.Write(httpResponse(http.StatusNotFound, "Not Found"))
    fmt.Println("Invalid link")
    return nil
  })

  // Registers the user agent header example endpoint
  // Sends the user agent back to the user
  server.RegisterRoute("/user-agent/", func(request HttpRequest) error {
    if request.Method == "GET" {
      responseHeader := fmt.Sprintf("Content-Type: %s\r\nContent-Length: %d\r\n", "text/plain", len(request.Header["User-Agent"]))
      request.Connection.Write(httpResponseWithData(http.StatusOK, "OK", responseHeader, request.Header["User-Agent"]))
      fmt.Println("user-agent: Valid link")
      return nil
    }
    
    request.Connection.Write(httpResponse(http.StatusNotFound, "Not Found"))
    fmt.Println("user-agent: Invalid link")
    return nil
  })

  server.Start()
}
