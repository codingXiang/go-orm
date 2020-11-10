package main

import (
	"github.com/codingXiang/go-logger/v2"
	"github.com/codingXiang/go-orm/v2"
)

type Test struct {
	ID   int
	Name string
}

func main() {
	var err error
	//初始化 logger
	logger.Log = logger.Default()
	dbConfig := orm.Default()
	dbConfig.Password = "a12345"
	//建立 orm instance
	if orm.DatabaseORM, err = orm.NewOrm(dbConfig); err != nil {
		panic(err)
	} else {
		logger.Log.Info(orm.DatabaseORM.ShowVersion())
	}
	orm.DatabaseORM.CheckTable(true, &Test{})
	if err := orm.DatabaseORM.Create(&Test{
		ID:   1111,
		Name: "hi",
	}).Error; err != nil {
		logger.Log.Error(err.Error)
	}
}
