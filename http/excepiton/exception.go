package excepiton

import "github.com/sirupsen/logrus"

type Exception interface {
	Error()
}


// 包的方式代用没有效果！ 注意panic and recover是配对的！也就是说recover捕获的是panic其他不行

func Finally()  {

	defer func() {
		if err := recover(); err != nil{
			logrus.Println(err)
		}
	}()
}


