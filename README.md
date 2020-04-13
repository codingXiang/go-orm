# 封裝 GORM 的 ORM
## 如何使用
### 載入模組
```shell
go get -u --insecure 10.40.42.38/BP05G0/go-orm
```

### 設定 ORM 實例
```
/*
	設定讀取並轉換 yaml 檔案
 */
var config *model.Database
file, err := ioutil.ReadFile("/Users/user/go/src/pkg/orm/example/config.yaml")

if err != nil {
	log.Fatalln("讀取 yaml 檔發生錯誤", err)
}

fmt.Println(string(file))

err = yaml.Unmarshal(file, &config)
if err != nil {
	log.Fatalln("轉換 yaml 檔發生錯誤", err)
}
//設定實例
var ORM = orm.NewOrm(GetConfig(new(model.Database)))

//版本更新
ORM.Upgrade()

//取得 gorm 實例
ORM.GetInstance()

```


## 參數設定
可以參考 example 裡面的 config.yaml，此格式可對照 model 裡面的 database

## 版本更新
一開始會在 database 中建立名為 versions 的 table，並且將 Version 參數寫入，後續只要有透過 `Upgrade` 方法更新版本，就會到 `UpgradeFilePath` 參數設定的路徑底下尋找舊版本版號的 SQL 檔案進行版本更新