package main

import (
	"fmt"
	"github.com/codingXiang/go-logger"
	"github.com/codingXiang/go-orm"
	"github.com/codingXiang/go-orm/model"
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
