package tcpserver

import (
	"log"
	"net"
)

// TCP server
type server struct {
	address                  string // Address to open connection: localhost:9999

	//回调函数
	onClientCallback      func(c *Client)
	onClientConnectionClosed func(c *Client, err error)
	onNewMessage             func(c *Client, message string)
}

// Called right after server starts listening new client
func (s *server) OnClient(callback func(c *Client)) {
	s.onClientCallback = callback
}

// Called right after connection closed
func (s *server) OnDisConnection(callback func(c *Client, err error)) {
	s.onClientConnectionClosed = callback
}

// Called when Client receives new message
func (s *server) OnMessage(callback func(c *Client, message string)) {
	s.onNewMessage = callback
}

// Start network server
func (s *server) Listen() {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		log.Fatal("Error starting TCP server.")
	}
	defer listener.Close()

	for {
		conn, _ := listener.Accept()
		client := &Client{
			conn:   conn,
			Server: s,
		}
		go client.Process()
	}
}

