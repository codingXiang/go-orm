package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/codingXiang/go-orm/v2/mongo"
	"time"
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
	data := mongo.NewRawData("1234", nil, raw)
	data.AddTag("namespace", "nms")
	if err := client.C(collection).Insert(data); err != nil {
		panic(err)
	}
	//fmt.Println(data.Tag)
	//搜尋
	selector := mongo.NewSearchCondition("", "1234", nil, nil)
	//if _, err := client.C(collection).First(selector); err != nil {
	//	panic(err)
	//}
	//

	c := context.Background()
	ctx, _ := context.WithTimeout(c, 5 * time.Minute)
	err = client.WaitForChange(ctx, func() (*mongo.RawData, error) {
		return client.C(collection).First(selector)
	}, func(data *mongo.RawData) (bool, error) {
		fmt.Println(data.Tag)
		return true, nil
	}, func() {
		fmt.Println("delete")
	})
	fmt.Println(err)
	//刪除
	//if err := client.C(collection).Delete(selector); err != nil {
	//	panic(err)
	//}
}
