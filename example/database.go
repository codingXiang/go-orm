package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

func main() {

	logger.Log = logger.NewLogger(logger.Logger{
		Level: "debug",
		Format: "json",
	})

	orm.NewOrm(GetDatqbaseConfig(new(model.Database)))
	orm.DatabaseORM.Upgrade()
	orm.DatabaseORM.GetInstance()
}

func GetDatqbaseConfig(config *model.Database) model.DatabaseInterface {
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
