package main

import (
	"encoding/json"
	"fmt"
	"github.com/codingXiang/go-orm/v2/mongo"
)

func main() {

	const (
		collection = "test"
	)
	m := mongo.Default()
	client := mongo.NewMongo(m)
	var raw interface{}
	err := json.Unmarshal([]byte(`{"test": "test"}`), &raw)
	if err != nil {
		panic(err)
	}
	// 新增
	data := mongo.NewRawData(raw)
	if err := client.C(collection).Insert(data); err != nil {
		panic(err)
	}

	//搜尋
	selector := mongo.NewSearchCondition("", data.GetIdentity(), nil)
	if d, err := client.C(collection).First(selector); err != nil {
		panic(err)
	} else {
		out, _ := json.Marshal(d)
		fmt.Println(string(out))
	}

	//刪除
	if err := client.C(collection).Delete(selector); err != nil {
		panic(err)
	}
}
