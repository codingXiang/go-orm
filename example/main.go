package main

import (
	"10.40.42.38/BP05G0/go-orm"
	"10.40.42.38/BP05G0/go-orm/model"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

func main() {
	var (
		ORM = orm.NewOrm(GetConfig(new(model.Database)))
	)
	ORM.Upgrade()
}

func GetConfig(config *model.Database) model.DatabaseInterface {
	file, err := ioutil.ReadFile("/Users/user/go/src/pkg/orm/example/config.yaml")

	if err != nil {
		log.Fatalln("讀取 yaml 檔發生錯誤", err)
	}

	fmt.Println(string(file))

	err = yaml.Unmarshal(file, &config)
	if err != nil {
		log.Fatalln("轉換 yaml 檔發生錯誤", err)
	}

	return config
}
