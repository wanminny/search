package utils

import (
	"time"
	"os"
	"log"
	"io/ioutil"
	"gobible/logmanager/cli/http/cache/redis"
)

// 超过半个月就删除文件 注意是秒数
const OVERTIME  = time.Hour * 24 * 15 / 1e9

const CHECKDELTA  = time.Hour * 8

// 超过指定的时间的文件要删除 下载目录中的所有文件遍历；处理超期的文件
func DeleteOverSomeTime()  {

	ticker := time.NewTicker(CHECKDELTA)

	configFile := redis.CheckRedisConf()
	configDown := configFile.Down
	//configDown := "/temp/php-cp"

	for {
		select {
			case <- ticker.C:
				infos,err := ioutil.ReadDir(configDown)
				if err != nil{
					log.Println(err)
					return
				}
				for _,info := range infos{
					if info.IsDir(){
						continue
					}else{
						log.Println(info.ModTime())
						if time.Now().Unix() > info.ModTime().Unix()  + int64(OVERTIME) {
							fileName := configDown + "/" + info.Name()
							err = os.Remove(fileName)
							if err != nil{
								log.Println(err)
							}
						}
						log.Println(info.Name())
					}
				}
		}
	}
}

