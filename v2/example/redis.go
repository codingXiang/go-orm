package main

import (
	"fmt"
	"github.com/codingXiang/configer/v2"
	"github.com/codingXiang/go-logger/v2"
	"github.com/codingXiang/go-orm/v2"
	"github.com/spf13/viper"
	"time"
)

var (
	Publisher  *orm.RedisClient
	Subscriber *orm.RedisClient
)

func main() {
	/*
		設定 Logger
	*/
	logger.Log = logger.Default()
	var config *viper.Viper
	core := configer.NewCore(configer.YAML, "redis-config", "./example")

	/*
		建立實例
	*/
	if c, err := core.ReadConfig(); err == nil {
		config = c
	} else {
		logger.Log.Error(err)
	}
	if client, err := orm.NewRedisClient("redis", config); err == nil {
		orm.RedisORM = client
		if info := client.Info("server"); info != nil {
			fmt.Println(info["server"]["redis_version"])
		}

	} else {
		panic(err.Error())

	}
	//if Subscriber, err = orm.NewRedisClient("redis", config); err != nil {
	//	panic(err.Error())
	//}
	////上傳 key
	err := orm.RedisORM.SetKeyValue("test", "test", 30 * time.Second)
	logger.Log.Debug(err)
	result, err := orm.RedisORM.GetValue("test")
	logger.Log.Debug(result)
	//發佈
}
