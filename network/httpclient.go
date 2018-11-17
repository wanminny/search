package main

import (
	"net/http"
	"log"
	"io"
	"os"
)

func main()  {


	client := http.Client{}
	rs,err := client.Get("http://www.baidu.com")

	if err != nil{
		log.Fatal(err)
	}

	io.Copy(os.Stdout,rs.Body)


	tran := 	http.Transport{}

	//tran.
}
