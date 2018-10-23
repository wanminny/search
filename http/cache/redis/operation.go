package redis

import (
	"time"
	"github.com/garyburd/redigo/redis"
	"github.com/sirupsen/logrus"
)

//const REDIS_HOST  =  "127.0.0.1:6379"

const REDIS_HOST  =  "10.0.3.11:6379"

const PASSWD  = "export-redis-01:3Ndy4tcWv3"

const REDIS_DB = 6

var (
	RedisClient *redis.Pool
)

func init()  {
	RedisClient = &redis.Pool{
		MaxIdle:     10,
		MaxActive:   0,
		IdleTimeout: 180 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", REDIS_HOST)
			if err != nil {
				return nil, err
			}

			if _, err := c.Do("AUTH", PASSWD); err != nil {
				c.Close()
				return nil, err
			}

			// 选择db
			c.Do("SELECT", REDIS_DB)
			return c, nil
		},
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