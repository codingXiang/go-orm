package main

import (
	"github.com/codingXiang/configer/v2"
	"github.com/codingXiang/go-logger/v2"
	"github.com/codingXiang/go-orm/v2"
)

type Test struct {
	ID   int
	Name string
}

func main() {
	//初始化 logger
	logger.Log = logger.Default()
	db := configer.NewCore(configer.YAML, "config", "./example", ".")
	//建立 orm instance
	if c, err := db.ReadConfig(); err == nil {
		if orm.DatabaseORM, err = orm.New(c); err != nil {
			panic(err)
		}
	} else {
		panic(err)
	}
	orm.DatabaseORM.CheckTable(true, &Test{})
	orm.DatabaseORM.Info("info")
	orm.DatabaseORM.Warn("warn")
	orm.DatabaseORM.Error("error")

	if err := orm.DatabaseORM.Create(&Test{
		ID:   1111,
		Name: "hi",
	}).Error; err != nil {
		logger.Log.Error(err.Error)
	}
}
