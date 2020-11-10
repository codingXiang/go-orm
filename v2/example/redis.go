package main

import (
	"fmt"
	"github.com/codingXiang/go-logger/v2"
	"github.com/codingXiang/go-orm/v2/redis"
	"time"
)

var (
	Publisher  *redis.RedisClient
	Subscriber *redis.RedisClient
)

func main() {
	/*
		設定 Logger
	*/
	logger.Log = logger.Default()
	//var config *viper.Viper
	//core := configer.NewCore(configer.YAML, "redis-config", "./example")

	/*
		建立實例
	*/
	//if c, err := core.ReadConfig(); err == nil {
	//	config = c
	//} else {
	//	logger.Log.Error(err)
	//}
	r := redis.Default()
	r.Password = "a12345"
	client := redis.NewRedis(r)
	if _, err := client.GetInfo(); err == nil {
		redis.RedisORM = client
		if info := redis.RedisORM.Info("server"); info != nil {
			fmt.Println(info["server"]["redis_version"])
		}

	} else {
		panic(err.Error())

	}
	//if Subscriber, err = orm.NewRedisClient("redis", config); err != nil {
	//	panic(err.Error())
	//}
	////上傳 key
	err := redis.RedisORM.SetKeyValue("test", "test", 30 * time.Second)
	logger.Log.Debug(err)
	result, err := redis.RedisORM.GetValue("test")
	logger.Log.Debug(result)
	//發佈
}
