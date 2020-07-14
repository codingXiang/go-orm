package main

import (
	"github.com/codingXiang/configer"
	"github.com/codingXiang/go-logger"
	"github.com/codingXiang/go-orm"
)

var (
	Publisher  orm.RedisClientInterface
	Subscriber orm.RedisClientInterface
)

func main() {
	/*
		設定 Logger
	*/
	logger.Log = logger.NewLogger(logger.Logger{
		Level:  "debug",
		Format: "json",
	})
	config := configer.NewConfigerCore("yaml", "redis-config", "./example")

	/*
		建立實例
	*/
	var err error
	if Publisher, err = orm.NewRedisClient("redis", config); err != nil {
		panic(err.Error())
	}
	if Subscriber, err = orm.NewRedisClient("redis", config); err != nil {
		panic(err.Error())
	}
	////上傳 key
	//orm.RedisORM.SetKeyValue("test", "test", 0)
	//發佈
}
