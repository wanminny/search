package main

import "gobible/logmanager/cli/http/hotcompiler"


/// 热加载 ok ！
// 注意路径
// 本代码采用的是放在 http/ 目录下面的 hot [也即将本代码编译成功后 放入该目录]
// 然后修改 dao/mysql 目录可以有效果！


// [重点] 这个是另外作为一个进程运行的；不可以直接使用go func(){} 附在http 进程中！
//
func main()  {

	//fsnotify
	go hotcompiler.HotCompile()

	select {

	}
}
