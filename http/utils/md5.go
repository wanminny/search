package utils

import (
	"crypto/md5"
	"encoding/hex"
)

/***
	整个项目通用的utils
 */
func MD5(text string) string{
	ctx := md5.New()
	ctx.Write([]byte(text))
	return hex.EncodeToString(ctx.Sum(nil))
}


