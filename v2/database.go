package orm

type Database struct {
	URL          string `yaml:"url"`          //Server 的位置
	Name         string `yaml:"name"`         //名稱
	Port         int    `yaml:"port"`         //Port
	LogMode      bool   `yaml:"logMode"`      //Log模式
	Username     string `yaml:"username"`     //使用者名稱
	Password     string `yaml:"password"`     //密碼
	Type         string `yaml:"type"`         //類型（例如 mysql、postgre、sqlite等)
	TablePrefix  string `yaml:"tablePrefix"`  //table前綴字
	MaxOpenConns int    `yaml:"maxOpenConns"` //最大開啟連線數
	MaxIdleConns int    `yaml:"maxIdleConns"` //最大連線
	MaxLifeTime  int    `yaml:"maxLifeTime"`  //最長連線時間
}

func Default() *Database {
	return &Database{
		URL:          "127.0.0.1",
		Name:         "test",
		Username:     "root",
		Password:     "",
		Port:         3306,
		LogMode:      true,
		Type:         MySQL.String(),
		MaxOpenConns: 1000,
		MaxIdleConns: 5,
		MaxLifeTime:  1000,
	}
}

//GetURL 取得 Database Server 位置
func (db *Database) GetURL() string {
	return db.URL
}

//GetName 取得 Database 名稱
func (db *Database) GetName() string {
	return db.Name
}

//GetPort 取得 Database Port
func (db *Database) GetPort() int {
	return db.Port
}

//GetType 取得 Database 類型
func (db *Database) GetType() string {
	return db.Type
}

//GetLogMode 取得 Database Log 模式
func (db *Database) GetLogMode() bool {
	return db.LogMode
}

//GetUsername 取得 Database 使用者名稱
func (db *Database) GetUsername() string {
	return db.Username
}

//GetPassword 取得 Database 密碼
func (db *Database) GetPassword() string {
	return db.Password
}

//GetTablePrefix 取得 Database Schema 前綴字
func (db *Database) GetTablePrefix() string {
	return db.TablePrefix
}

//GetURL 取得 Database Server maxOpenConns
func (db *Database) GetMaxOpenConns() int {
	return db.MaxOpenConns
}

//GetMaxIdleConns 取得 Database Server maxIdleConns
func (db *Database) GetMaxIdelConns() int {
	return db.MaxIdleConns
}

//GetMaxLifeTime 取得 Database Server maxLifeTime
func (db *Database) GetMaxLifeTime() int {
	return db.MaxLifeTime
}
