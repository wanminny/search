package main

import (
	"log"
)

func init()  {
	log.SetPrefix("Trace: ")
	//log.SetFlags(log.Ldate | log.Lmicroseconds | log.Llongfile)
	log.SetFlags(log.Llongfile )
}

func test()  {
	//log.Logger{}
}
func main()  {

	log.Println("start")
	// 1 Trace: /gopath/src/gobible/logmanager/cli/std/log/log.go:17: AA
	// 2 Trace: /usr/local/Cellar/go/1.10.1/libexec/src/runtime/proc.go:198: AA
	log.Output(2,"AA")

	log.Println(log.Flags(),log.Prefix())
	//log.Printf()
	//fmt.Fprintf()
	//log.Fatal("fatal o ooo ")
	log.Panicln("ooo panic")
	log.Println("end")
}
