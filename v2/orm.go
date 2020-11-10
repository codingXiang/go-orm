package orm

import (
	"fmt"
	"github.com/codingXiang/configer/v2"
	"github.com/codingXiang/go-logger/v2"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"    // mysql
	_ "github.com/jinzhu/gorm/dialects/postgres" //postgresql
	"github.com/spf13/viper"
	"time"
)

//Orm
type Orm struct {
	*gorm.DB
	config *viper.Viper
}

var (
	DatabaseORM *Orm
)

//NewOrm : 新增 ORM 實例
func NewOrm(database *Database) (*Orm, error) {
	c := viper.New()
	c.Set(configer.GetConfigPath(DB, Url), database.GetURL())
	c.Set(configer.GetConfigPath(DB, LogMode), database.GetLogMode())
	c.Set(configer.GetConfigPath(DB, Port), database.GetPort())
	c.Set(configer.GetConfigPath(DB, Name), database.GetName())
	c.Set(configer.GetConfigPath(DB, Username), database.GetUsername())
	c.Set(configer.GetConfigPath(DB, Password), database.GetPassword())
	c.Set(configer.GetConfigPath(DB, Type), database.GetType())
	c.Set(configer.GetConfigPath(DB, TablePrefix), database.GetTablePrefix())
	c.Set(configer.GetConfigPath(DB, MaxIdleConns), database.GetMaxIdelConns())
	c.Set(configer.GetConfigPath(DB, MaxLifeTime), database.GetMaxLifeTime())
	c.Set(configer.GetConfigPath(DB, MaxOpenConns), database.GetMaxOpenConns())

	return New(c)
}

func New(config *viper.Viper) (*Orm, error) {
	var (
		err error
		orm = new(Orm)
	)
	if configer.Config == nil {
		//初始化 configer
		configer.Config = configer.NewConfiger()
	}
	var (
		logMode     = config.GetBool(configer.GetConfigPath(DB, LogMode))
		tablePrefix = config.GetString(configer.GetConfigPath(DB, TablePrefix))
	)
	//設定資料庫型態 (MySQL, PostgreSQL) 與連線資訊
	logger.Log.Debug("setup database type")
	err = orm.setDBType(config)

	//設定資料庫參數
	logger.Log.Debug("setup database config")
	orm.setDbConfig(config)

	//設定是否開啟 Log 模式
	logger.Log.Debug("setup log mode =", logMode)
	if logMode {
		orm.LogMode(true)
		orm.DB = orm.Debug()
	} else {
		orm.LogMode(false)
	}

	//設定預設 Table 前綴字
	logger.Log.Debug("setup table name prefix", tablePrefix)
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return tablePrefix + defaultTableName
	}

	return orm, err
}

func (orm *Orm) ShowVersion() (version string) {
	orm.DB.DB().QueryRow("SELECT VERSION()").Scan(&version)
	return
}

func (orm *Orm) CloseDB() {
	logger.Log.Debug("close database connection")
	defer orm.Close()
}

//GetTableName : 透過傳入 struct 回傳 table 名稱
func (orm *Orm) GetTableName(tb interface{}) string {
	return orm.NewScope(tb).TableName()
}

//CheckTable : 檢查 Table 是否存在，不存在建立並回傳 false, 反之回傳 true
func (orm *Orm) CheckTable(migrate bool, value interface{}) error {
	var (
		err error
		tx  = orm.Begin()
	)
	if !orm.HasTable(value) {
		if err = tx.Set("gorm:table_options", "ENGINE=InnoDB  DEFAULT CHARSET=utf8 AUTO_INCREMENT=1;").CreateTable(value).Error; err != nil {
			fmt.Println(err)
			tx.Rollback()
			return err
		}
		tx.Commit()
	} else {
		if migrate {
			if err = tx.Set("gorm:table_options", "ENGINE=InnoDB  DEFAULT CHARSET=utf8 AUTO_INCREMENT=1;").AutoMigrate(value).Error; err != nil {
				tx.Rollback()
				return err
			}
			tx.Commit()
		}
	}
	return nil
}

//setDBType : 設定資料庫型態
func (orm *Orm) setDBType(config *viper.Viper) error {
	var (
		err      error
		url      = config.GetString(configer.GetConfigPath(DB, Url))
		port     = config.GetInt(configer.GetConfigPath(DB, Port))
		dbName   = config.GetString(configer.GetConfigPath(DB, Name))
		username = config.GetString(configer.GetConfigPath(DB, Username))
		password = config.GetString(configer.GetConfigPath(DB, Password))
		_type    = NewDatabaseType(config.GetString(configer.GetConfigPath(DB, Type)))
	)
	logger.Log.Debug("set database type = ", _type)
	switch _type {
	case MySQL:
		var connectStr = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true",
			username, password, url, port, dbName)
		logger.Log.Debug("connection string = ", connectStr)
		orm.DB, err = gorm.Open(_type.String(), connectStr)
		if err != nil {
			logger.Log.Error("connect to database error", err)
			return err
		}
		break
	case PostgreSQL:
		var connectStr = fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
			url, port, dbName, username, password)
		logger.Log.Debug("connection string = ", connectStr)
		orm.DB, err = gorm.Open(_type.String(), connectStr)
		if err != nil {
			logger.Log.Error("connect to database error", err)
			return err
		}
		break
	}
	return nil
}

// 設定資料庫組態
func (orm *Orm) setDbConfig(config *viper.Viper) {
	var (
		maxOpenConns = config.GetInt("database.maxOpenConns")
		maxIdleConns = config.GetInt("database.maxIdleConns")
		maxLifeTime  = config.GetInt("database.maxLifeTime")
	)
	logger.Log.Debug("set max idle connections", maxIdleConns)
	orm.DB.DB().SetMaxIdleConns(maxIdleConns)
	logger.Log.Debug("set max open connections", maxOpenConns)
	orm.DB.DB().SetMaxOpenConns(maxOpenConns)
	logger.Log.Debug("set connection max life time", maxLifeTime)
	orm.DB.DB().SetConnMaxLifetime(time.Duration(maxLifeTime))
}
