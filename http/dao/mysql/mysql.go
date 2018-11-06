package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/url"
	"strconv"
)

type MySQLClient struct {
	Host    string  `json:"host"`
	MaxIdle string  `json:"max_idle"`
	MaxOpen string  `json:"max_open"`
	User    string  `json:"username"`
	Pwd     string  `json:"password"`
	DB      string  `json:"database"`
	Port    string  `json:"port"`
	Pool    *sql.DB `json:"pool"`
}

func (mc *MySQLClient) Init() (err error) {

	// 构建 DSN 时尤其注意 loc 和 parseTime 正确设置
	// 东八区，允许解析时间字段
	uri := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&loc=%s&parseTime=true",
		mc.User,
		mc.Pwd,
		mc.Host,
		mc.Port,
		mc.DB,
		url.QueryEscape("Asia/Shanghai"))

	// Open 全局一个实例只需调用一次
	mc.Pool, err = sql.Open("mysql", uri)
	if err != nil {
		log.Fatal(err)
		return err
	}
	//使用前 Ping, 确保 DB 连接正常
	err = mc.Pool.Ping()
	if err != nil {
		log.Fatal(err)
		return err
	}

	// 设置最大连接数，一定要设置 MaxOpen
	maxIdle, _ := strconv.Atoi(mc.MaxIdle)
	maxOpen, _ := strconv.Atoi(mc.MaxOpen)

	mc.Pool.SetMaxIdleConns(maxIdle)
	mc.Pool.SetMaxOpenConns(maxOpen)
	return nil
}
