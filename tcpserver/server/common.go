package tcpserver

import "log"

// Creates new tcp server instance
func New(address string) *server {
	log.Println("Creating server with address", address)
	server := &server{
		address: address,
	}

	// 错误？
	//server.OnClient(func(c *Client) {})
	//server.OnMessage(func(c *Client, message string) {})
	//server.OnDisConnection(func(c *Client, err error) {})

	return server
}



