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
	logger.Log = logger.NewLogger(logger.Logger{
		Level:  "debug",
		Format: "json",
	})
	databaseConfig := configer.NewConfigerCore("yaml", "config", "./example")
	if orm.DatabaseORM, err = orm.NewOrm("database", databaseConfig); err != nil {
		panic(err)
	}
	if err = orm.DatabaseORM.Upgrade(&Test{}); err != nil {
		panic(err.Error())
	}
	orm.DatabaseORM.GetInstance()
}
