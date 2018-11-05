package main

import "log"

func init()  {
	log.SetFlags(log.Llongfile | log.Ltime)
}

func test()  {
	(*struct{ error }).Error(nil)
}


//
func SliceCopyImplement1(destSlice []string, srcSlice []string)  {

	log.Println(destSlice,srcSlice)

	n := copy(destSlice,srcSlice)
	log.Println(n)

	log.Println(destSlice,srcSlice)
}


func SliceCopyImplement(destSlice *[]string, srcSlice []string)  {


	//dest := make([]string,len(srcSlice))
	n := copy(*destSlice,srcSlice)
	log.Println(n)

	log.Println(destSlice,srcSlice)
}

func main()  {

	src := []string{"golang","interpreter","compiler","network"}
	dest := []string{"abc","123","456"}

	SliceCopyImplement(&dest,src)


	log.Println(src,dest)

}
