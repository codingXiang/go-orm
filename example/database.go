package main

import (
	"github.com/codingXiang/configer"
	"github.com/codingXiang/go-logger"
	"github.com/codingXiang/go-orm"
)

type Test struct {
	ID   int
	Name string
}

func main() {
	var err error
	//初始化 logger
	logger.Log = logger.NewLogger(logger.Logger{
		Level:  "debug",
		Format: "json",
	})
	//設定 configer
	databaseConfig := configer.NewConfigerCore("yaml", "config", "./example")
	//建立 orm instance
	if orm.DatabaseORM, err = orm.NewOrm("database", databaseConfig); err != nil {
		panic(err)
	}
	//取得實例
	orm.DatabaseORM.GetInstance()
	//版本更新
	if err = orm.DatabaseORM.Upgrade(&Test{}); err != nil {
		panic(err.Error())
	}
}
