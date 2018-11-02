package main

import (
	"gobible/logmanager/cli/tcpserver/server"
	"log"
)

func init()  {
	log.SetFlags(log.Llongfile | log.Ltime)
}

func main() {


	server := tcpserver.New("localhost:8881")

	server.OnClient(func(c *tcpserver.Client) {
		log.Println("onclinent.")
		// new client connected
		// lets send some message
		c.Send("Hello")
	})
	server.OnMessage(func(c *tcpserver.Client, message string) {
		// new message received
		log.Println("on message.")
		log.Println(message)
	})
	server.OnDisConnection(func(c *tcpserver.Client, err error) {
		// connection with client lost
		log.Println("on disconection.")
		log.Println(err)
	})
	server.Listen()
}
