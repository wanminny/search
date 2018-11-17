package main

import (
	"net"
	"log"
	"time"
)



func main()  {

	for{
		conn,err := net.Dial("tcp",":50000")
		if err != nil{
			log.Println(err)
		}else{
			go func() {
				for {
					buf := make([]byte,64)
					conn.Read(buf)
					log.Println(buf)
				}
			}()
			go func() {
				for{
					conn.Write([]byte("self-conn"))
				}
			}()

			log.Println("ok ,entry this !")
			time.Sleep(time.Second * 2)
			break
		}
		//time.Sleep(time.Microsecond)
	}

}
