package errors

import (
	"fmt"
)

type WORK_STATUS int16

const (
	_  WORK_STATUS = iota
	WORK_ON_STATUS
	WORK_OFF_STATUS
)


var StatusScopePrefix = "xxx-system:"

//状态码
const G_STATUS_SCOPE_BASE  = 10000

//code 1
const (
	USER_STATUS_INVALID WORK_STATUS  = 1* G_STATUS_SCOPE_BASE + iota

	USER_STATUS_NOT_EXISTS

	USER_STATUS_TOO_LONG

)

//code 2
const (
	ORDER_STATUS_INVALID WORK_STATUS  =  2 * G_STATUS_SCOPE_BASE + iota

	ORDER_STATUS_NOT_EXISTS

	ORDER_STATUS_TOO_LONG
)


//code 3
const (
	GOODS_STATUS_INVALID WORK_STATUS  =  3 * G_STATUS_SCOPE_BASE + iota

	GOODS_STATUS_NOT_EXISTS

	GOODS_STATUS_TOO_LONG

)

var GStatusCode = map[WORK_STATUS]string{

	//code1
	USER_STATUS_INVALID:"无效的用户状态",
	USER_STATUS_NOT_EXISTS:"用户不存在",
	USER_STATUS_TOO_LONG:"xxx太长",

	//code2
	ORDER_STATUS_INVALID:"",
	ORDER_STATUS_NOT_EXISTS:"",
	ORDER_STATUS_TOO_LONG:"",

	//code3
	GOODS_STATUS_INVALID:"",
	GOODS_STATUS_NOT_EXISTS:"",
	GOODS_STATUS_TOO_LONG:"",
}

func (w WORK_STATUS) String() string  {
	if v,ok := GStatusCode[w]; ok {
		return v
	}
	return "nothing"
}

func (w WORK_STATUS)Error() string{
	return fmt.Sprintf("%s,%d:%s",StatusScopePrefix,w,w.String())
}