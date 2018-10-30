package utils

import (
	"log"
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

type CustomerConfig struct {

	//Product
}


func YmlParse()  {
	config := map[string]map[string]map[string]string{}

	c,err := ioutil.ReadFile("../config/config.yaml")
	//c,err := ioutil.ReadFile("./config/config.yaml")
	if err != nil{
		log.Println(err)
	}

	err = yaml.Unmarshal(c,&config)
	if err != nil{
		log.Println(err)
	}

	log.Println(config)

	host := config["redis"]["test"]["host"]
	port := config["redis"]["test"]["port"]
	passwd := config["redis"]["test"]["password"]
	db := config["redis"]["test"]["db"]

	log.Println("===================")
	log.Println(host,port,passwd,db)
}