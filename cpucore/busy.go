package main

import (
	"log"
	"runtime"
)

func init()  {

	//runtime.GOMAXPROCS(2)
}

func main() {

	log.Println("start ")
	log.Println("set cpu core :",runtime.NumCPU())

	//启动一个go程 和两个后 分别观察 cpu
	// 注意这里是 计算型的任务 内存基本不算太多。
	// summary : one goroutine 90+%  two go routine 190% three goroutine 270%
	// four : 330%
	// five : 我本机是四核的;基本是 和four差不多 330 左右不等 CPU使用率

	for i := 0 ;i < 5; i++{
		go func() {
			for{
				// compute cpu-consume task
				for i := 1;i <10000;i++ {
					_ = i * i
				}
			}
		}()

	}

	log.Println("one core is busy now ?")
	select {

	}
}
