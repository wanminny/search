package main

import (
	"net"
	"log"
	"time"
)

func CheckError(err error)  {

	if err != nil{
		//log.Println(err)
		//return
		log.Fatal(err)
	}
}

func main()  {


	conn,err := net.Dial("tcp","127.0.0.1:12345")

	CheckError(err)

	var buf = make([]byte,512)

	for {

		n, err := conn.Read(buf)
		CheckError(err)

		log.Printf("received %d bytes contents is :%s",n,string(buf[:n]))

		conn.Write([]byte("client miss you!"))

		time.Sleep(time.Second *2)
	}

}
