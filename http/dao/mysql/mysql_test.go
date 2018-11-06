package mysql

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"testing"
)

func init() {
	log.SetFlags(log.Llongfile | log.Ltime)
}

//方法 的调用与普通函数的调用方式 不同。
func TestMySQLClient_Init(t *testing.T) {

	config := map[string]map[string]map[string]string{}
	c, err := ioutil.ReadFile("../../config/config.yaml")
	//c,err := ioutil.ReadFile("./config/config.yaml")
	if err != nil {
		log.Println(err)
	}
	err = yaml.Unmarshal(c, &config)
	if err != nil {
		log.Println(err)
	}
	//log.Println(config)

	host := config["mysql"]["test"]["host"]
	port := config["mysql"]["test"]["port"]
	uname := config["mysql"]["test"]["username"]
	passwd := config["mysql"]["test"]["password"]
	db := config["mysql"]["test"]["database"]
	idle := config["mysql"]["test"]["maxIdle"]
	open := config["mysql"]["test"]["maxOpen"]

	log.Println("==============================")
	log.Println(host, port, uname, passwd, db, idle, open)

	var mysqlC = MySQLClient{
		Host:    host,
		Port:    port,
		User:    uname,
		Pwd:     passwd,
		DB:      db,
		MaxIdle: idle,
		MaxOpen: open,
	}
	mysqlC.Init()

	rows, err := mysqlC.Pool.Query("select * from api")

	if err != nil {
		log.Println(err)
	}
	var id = 0
	var aid = ""
	var num = ""

	//	`url` varchar(240) DEFAULT NULL COMMENT '请求地址',
	//	`name` varchar(100) DEFAULT NULL COMMENT '接口名',
	//	`des` varchar(300) DEFAULT NULL COMMENT '接口描述',
	//	`parameter` text COMMENT '请求参数{所有的主求参数,以json格式在此存放}',
	//	`memo` text COMMENT '备注',
	//	`re` text COMMENT '返回值',
	//	`lasttime` int(11) unsigned DEFAULT NULL COMMENT '提后操作时间',
	//	`lastuid` int(11) unsigned DEFAULT NULL COMMENT '最后修改uid',
	//	`isdel` tinyint(4) unsigned DEFAULT '0' COMMENT '{0:正常,1:删除}',
	//`type` char(11) DEFAULT NULL COMMENT '请求方式',
	//`ord` int(11) DEFAULT '0' COMMENT '排序(值越大,越靠前)',

	var url = ""
	var name = ""
	var des = ""
	var parameter = ""
	var demo = ""
	var re = ""
	var lasttime int
	var lastuid int
	var isdel int
	var typea string
	var ord int

	//或者使用结构体方式 要方便一些 但是要先定义结构体！

	log.Println("%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%")

	//必须所有的字段写全才可以！！
	for rows.Next() {
		rows.Scan(&id, &aid, &num, &url, &name, &des, &parameter, &demo, &re, &lasttime, &lastuid, &isdel, &typea, &ord)
		log.Println(id, aid, num, url, name, des, parameter, demo, re, lasttime, lastuid, isdel, typea, ord)
		break
	}

	mysqlC.Pool.Prepare("")
	//log.Printf("%#v",rows)
}
