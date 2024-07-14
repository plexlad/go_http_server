package main

import (
	"fmt"
	"net"
	"net/http"
	"strings"
)

// This code is for a handler and http request. Makes adding endpoints easier.
// Not very complex, but works

type HttpRequest struct {
  Connection net.Conn
  Method     string
  Route      string
  RouteData  string // The part at the end of the route
  Protocol   string
  Header     map[string]string 
  Body       string
}

// Creates an HttpRequest from raw byte data. Parses into a convenient form factor
func ParseHttpRequest(rawData []byte, connection net.Conn) HttpRequest {
  stringData := string(rawData)

  // Parses the CR LF characters to read the request
  parsedLines := strings.Split(stringData, "\r\n")
  
  // Gets the request info and parses it into a slice
  requestInfo := strings.Split(parsedLines[0], " ")
  
  // Parses the route to get the data at the end
  parsedRoute := strings.Split(requestInfo[1], "/")
  routeData := parsedRoute[len(parsedRoute)-1]
  
  // Set the header key value map
  header := make(map[string]string)
  for _, value := range parsedLines[1:len(parsedLines)-1] {
    headerKeyAndValue := strings.Split(value, ": ")
    if len(headerKeyAndValue) > 1 {
      header[headerKeyAndValue[0]] = headerKeyAndValue[1]
    }
  }

  // TODO: Manage different types of requests
  return HttpRequest{
    connection,
    requestInfo[0], // Method
    requestInfo[1], // Route
    routeData,      // Route data
    requestInfo[2], // Protocol
    
    header, // Header
    parsedLines[len(parsedLines)-1], // Body
  }
}

// Manages a set of routes and requests.
type HttpHandler struct {
  Routes map[string]func(HttpRequest) error // Route handlers assigned to a string root
}

// TODO: Add more functionality for more convenience (RegisterRoute for HttpHandler)
// The added functionality can include glob support, etc.
// Consider using a binary tree for the route handling.

// Registers a route and its function handler
func (h *HttpHandler) RegisterRoute(route string, handler func(HttpRequest) error) error {
  // Removes slash from the start and beginning for easier parsing
  parsedRoute := strings.Trim(route, "/")

  // Make a new hashmap if it does not exist
  if h.Routes == nil {
    h.Routes = make(map[string]func(HttpRequest) error)
  }
  h.Routes[parsedRoute] = handler

  return nil
}

func (h *HttpHandler) HandleRequest(request HttpRequest) error {
  // Trims and seperates the route into individual words
  parsedRequestRoute := strings.Split(strings.Trim(request.Route, "/"), "/")

  // Gets the base of the route to find the handler
  baseRequestRoute := ""
  for _, routeSection := range parsedRequestRoute {
    baseRequestRoute = fmt.Sprintf("%s/%s", baseRequestRoute, routeSection)
  }

  // Remove the slashes for use with Routes
  baseRequestRoute = strings.Trim(baseRequestRoute, "/")

  fmt.Printf("baseRequestRoute = '%s'\n", baseRequestRoute)

  // Checks both route and route data
  handlerFunc, ok := h.Routes[baseRequestRoute]
  if !ok {
    for _, routeSection := range parsedRequestRoute[:len(parsedRequestRoute)-1] {
      baseRequestRoute = fmt.Sprintf("%s/%s", baseRequestRoute, routeSection)
    }

    fmt.Printf("baseRequestRoute = '%s'\n", baseRequestRoute)

    handlerFunc, ok = h.Routes[baseRequestRoute]
    if !ok {
      request.Connection.Write(httpResponse(http.StatusNotFound, "Not Found"))
      return NewError(RouteNotFoundError, "Invalid route in HttpRequest")
    }
  }

  return handlerFunc(request)
}
