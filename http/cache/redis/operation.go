package redis

import (
	"time"
	"github.com/garyburd/redigo/redis"
	"log"
	"io/ioutil"
	"encoding/json"
	"github.com/sirupsen/logrus"
)


var (
	RedisClient *redis.Pool
)

type ConfigContent struct {
	RedisHost string `json:"redis_host"`
	RedisPort string `json:"redis_port"`
	RedisDb string `json:"redis_db"`
	RedisPassWd string `json:"redis_pass_wd"`
	Down string `json:"down"`
}


func CheckRedisConf() (config * ConfigContent) {

	configC := &ConfigContent{}
	byteC,err := ioutil.ReadFile("./config/config.json")
	//byteC,err := ioutil.ReadFile("../../config/config.json")

	if err != nil{
		log.Fatal(err)
		return
	}
	err = json.Unmarshal(byteC,configC)
	if err != nil{
		log.Fatal(err)
		return
	}

	return configC
}

func InitPool()  {

	configC := CheckRedisConf()

	RedisClient = &redis.Pool{
		MaxIdle:     10,
		MaxActive:   0,
		IdleTimeout: 180 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", configC.RedisHost + ":"+ configC.RedisPort)
			if err != nil {
				logrus.Println("Redis Dial 异常: ",err)
				return nil, err
			}

			//自由填写了密码才去拨号
			if len(configC.RedisPassWd) >0 {
				if _, err := c.Do("AUTH", configC.RedisPassWd); err != nil {
					c.Close()
					logrus.Println(" redis AUTH 异常 :",err)
					return nil, err
				}
			}

			// 选择db
			c.Do("SELECT", configC.RedisDb)
			return c, nil
		},
	}
}


func init()  {

	InitPool()
}


func PING(){
	//获取值
	conn := RedisClient.Get()
	_,err := redis.String(conn.Do("PING"))
	if err != nil{
		log.Fatal(err)
	}
}

//获取值
func GetValue(key string) (value string,err error) {

	conn := RedisClient.Get()

	return redis.String(conn.Do("GET",key))
}

//设置值
func SetValue(key string,value string)  {

	conn := RedisClient.Get()

	_, err := conn.Do("SET", key, value)

	if err != nil{
		logrus.Println(err)
	}

}

func DelKey(key string)  {

	conn := RedisClient.Get()
	_, err := conn.Do("DEL", key)
	if err != nil{
		logrus.Println(err)
	}
}

func AllKeys() (kvs map[string]string) {

	kvs = make(map[string]string)
	conn := RedisClient.Get()

	reply ,_ := redis.Values(conn.Do("keys","*"))

	for _,v := range reply{

		k := string(v.([]byte))
		value ,_ := redis.String(conn.Do("GET",k))
		kvs[k] = value
	}
	return
}

//入队列 左边进；
func LPush(key,value string,) (err error) {
	conn := RedisClient.Get()
	_, err = conn.Do("LPUSH", key,value)
	if err != nil{
		logrus.Println(err)
	}
	return
}

//右边出
func Rop(key,value string)  (err error)  {
	conn := RedisClient.Get()
	_, err = conn.Do("RPOP", key)
	if err != nil{
		logrus.Println(err)
	}
	return
}
//获取队列
func LRange(key string,start,end int) (err error)  {
	conn := RedisClient.Get()
	_, err = conn.Do("lrange", key,start,end)
	if err != nil{
		logrus.Println(err)
	}
	return
}

// 设置值;
func HSet(key,field string,value interface{}) (err error) {
	conn := RedisClient.Get()
	_, err = conn.Do("hset", key,field,value)
	if err != nil{
		logrus.Println(err)
	}
	return
}

//查询状态;
func HGet(key,field string)  (err error) {
	conn := RedisClient.Get()
	_, err = conn.Do("hget", key,field)
	if err != nil{
		logrus.Println(err)
	}
	return

}

func HMSet(key string,m map[string]interface{}) (err error)  {

	conn := RedisClient.Get()
	//m := map[string]string{
	//	"title":  "Example2",
	//	"author": "Steve",
	//	"body":   "Map",
	//}
	if _, err := conn.Do("HMSET", redis.Args{}.Add(key).AddFlat(m)...); err != nil {
		panic(err)
	}
	return
}

//func Trans()  {
//	conn := RedisClient.Get()
//	conn.Send("MULTI")
//	conn.Send()
//	conn.Do("EXEC")
//}