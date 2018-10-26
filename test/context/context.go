package main


import (
	"net/http"
	"context"
	"time"
	"io/ioutil"
	"fmt"
	"log"
	"io"
	"os"
	"sync/atomic"
	"sync"
)
var (
	i = 100
	f = 3.14
	b = true
	s = "Clear is better than clever"
	p = point{252, 101}
)
type point struct{ x, y int64 }

var bi int64 = -922337203685
var ui uint64 = 1844679551615

func outputToWriter() (string, error) {
	file, err := ioutil.TempFile("", "example")
	if err != nil {
		return "", err
	}
	defer file.Close()
	fmt.Fprintf(file, "i = %#v\nf = %#v\nb = %#v\ns = %#v\nbi = %#v\nui = %#v\np = %#v\n", i, f, b, s, bi, ui, p)
	return file.Name(), nil
}


func demo(res http.ResponseWriter,request *http.Request)  {

	ctx := request.Context()

	context.WithTimeout(ctx,time.Second * 2)

}
// pipe
func PipeExample() error {
	r, w := io.Pipe()

	go func() {
		w.Write([]byte("test\n"))
		w.Close()
	}()

	if _, err := io.Copy(os.Stdout, r); err != nil {
		return err
	}

	return nil
}

func todo()  {

	tod := context.TODO()
	back := context.Background()

	<- tod.Done()
	_ := context.Canceled
	_ := context.DeadlineExceeded
	//context.
	//back.Err()

	ctx,cancelFunc := context.WithTimeout(tod,time.Second)
	ctx,cancelFunc = context.WithDeadline(tod,time.Now().Add(time.Second))
	ctx,cancelFunc = context.WithCancel(back)




	ctx = context.WithValue(tod,"key","value")

	ctx.Value("key")
	ctx.Deadline()
	//r := http.Request{}
	//r.Context()

	log.Printf("%+v,%+v",tod,back)
	log.Printf("%#v,%#v",tod,back)
	log.Printf("%v,%v",tod,back)


	//atomic.Value{}

	pool := sync.Pool{
		New: func() interface{} {
			log.Println(345)
			return 3
		},
	}
	pool.Get()
	pool.Put(1)

	fmt.Println(1111)


	m := sync.Map{}
	m.Load(1)



}


func main()  {

	todo()
	os.Exit(1)
	PipeExample()

	n,_ := outputToWriter()
	log.Print(n)
}