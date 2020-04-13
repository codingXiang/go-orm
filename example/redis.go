package main

import (
	"10.40.42.38/BP05G0/go-logger"
	"10.40.42.38/BP05G0/go-orm"
	"10.40.42.38/BP05G0/go-orm/model"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

func main() {
	/*
		設定 Logger
	 */
	logger.Log = logger.NewLogger(logger.Logger{
		Level:  "debug",
		Format: "json",
	})

	/*
		設定參數
	 */
	var config = new(model.Redis)
	file, err := ioutil.ReadFile("/Users/user/go/src/pkg/orm/example/redis-config.yaml")

	if err != nil {
		log.Fatalln("讀取 yaml 檔發生錯誤", err)
	}

	fmt.Println(string(file))

	err = yaml.Unmarshal(file, &config)
	if err != nil {
		log.Fatalln("轉換 yaml 檔發生錯誤", err)
	}

	/*
		建立實例
	 */
	orm.NewRedisClient(config)
	//上傳 key
	orm.RedisORM.SetKeyValue("test", "test", 0)
}
