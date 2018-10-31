package main

import (
	"net/http"
	"os"
	"log"
	"io"
)

func main()  {

	res,err := http.Get(os.Args[1])

	if err != nil{
		log.Println(err)
	}

	defer res.Body.Close()

	f,err := os.Create(os.Args[2])
	if err != nil{
		log.Fatal(err)
	}
	defer f.Close()
	w := io.MultiWriter(f,os.Stdout)

	io.Copy(w,res.Body)

}
