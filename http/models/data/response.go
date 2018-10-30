package data

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
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

func SetCROS(res http.ResponseWriter)  {

	res.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	res.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	res.Header().Set("content-type", "application/json")             //返回数据格式是json
}