package config

import "errors"

var HashNotExists  = errors.New("该任务不存在！")

var AuthFail  = errors.New("auth认证失败！")
