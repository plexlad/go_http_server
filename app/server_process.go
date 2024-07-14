package main

import (
	"fmt"
	"net"
)

// The goal is for the server to be interacted with externally.
// No manual initialization

type Server struct {
  address      string
  port         string
  routeHandler HttpHandler
  listener     net.Listener
}

// Just a wrapper. Look at the comments for RegisterRoute
func (s *Server) RegisterRoute(route string, handler func(HttpRequest) error) error {
  return s.routeHandler.RegisterRoute(route, handler)
}

func (s *Server) HandleConnection(connection net.Conn) {
  // Close the connection after it is handled
  defer connection.Close()
  fmt.Println("Accepted request")
  
  // Buffer with request information
  requestBuffer := make([]byte, 1024)
  connection.Read(requestBuffer)

  request := ParseHttpRequest(requestBuffer, connection)
  err := s.routeHandler.HandleRequest(request)
  if err != nil {
    fmt.Println("Error: " + err.Error())
    return
  }

  fmt.Println("Responded")
}

// Starts the server. Some of the only initialization that is needed.
func (s *Server) Start() {
  // Close the listener when this is all done
  fmt.Printf("Serving on %s:%s\n", s.address, s.port)
  fmt.Println("Listening")

  defer s.listener.Close()
  for {
    connection, err := s.listener.Accept()
    if err != nil {
      fmt.Println("Error accepting connection: ", err.Error())
    }

    go s.HandleConnection(connection)
  }
}

// Use this function to make a new server! Initializing things manually can really be a pain
func NewServer(address string, port string) (*Server, error) {
  listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", address, port))
  if err != nil {
    fmt.Println("Failed to bind to port " + port)
    fmt.Println("Error" + err.Error())
    return nil, err
  }

  server := Server{
    address,
    port,
    *new(HttpHandler),
    listener,
  }
 
  return &server, nil
}
