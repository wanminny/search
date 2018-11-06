package tcpserver

import (
	"log"
	"net"
)

// TCP server
type server struct {
	address                  string // Address to open connection: localhost:9999

	//应该把客户端的连接进行管理？？
	// conn ?

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

	// 可以将  listener 保存起来；同下面的conn

	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		log.Fatal("Error starting TCP server.")
	}
	defer listener.Close()

	for {
		conn, _ := listener.Accept()

		//终于；在此处进行了连接的 管理！
		client := &Client{
			listener:listener,
			conn:   conn,
			Server: s,
		}
		go client.Process()
	}
}

