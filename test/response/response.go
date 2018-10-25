package response

import (
	"github.com/sirupsen/logrus"
	"encoding/json"
)


// 返回Result
type Result struct {
	Code int `json:"code"`
	Message string `json:"message"`
	Data interface{} `json:"data"`
}


// 返回Json
func NewJson(code int,msg string,data interface{}) string {

	if data == nil{
		data = make([]string,0)
	}
	rlt,err := json.Marshal(Result{Code:code,Message:msg,Data:data})
	if err != nil{
		logrus.Println(err)
	}
	return string(rlt)
}

