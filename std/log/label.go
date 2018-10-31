package main

import "log"

func main()  {

	log.Println(1111)

	//var a string

	var (
		c = 1

	)
LABEL:
	//d := 4
	for {
			for {
				if c == 5{
					break LABEL
				}
				log.Println(c)
				c++
			}
		}
		e := 3
		//log.Println(d)
		log.Println("*******")
		log.Println(e)
		log.Println("=======")

}
