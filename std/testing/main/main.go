package main

import (
	"gobible/logmanager/cli/std/testing"
	"log"
	//_ "expvar"

	//"expvar"
)


// expvar  http://localhost:8083/debug/vars

func main()  {

	testing.RunHeartbeatService(":8083")

	log.Println(4444)
	//select {
	//
	//}
}
