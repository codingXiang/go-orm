package orm

import (
	"context"
	"fmt"
	"github.com/codingXiang/configer/v2"
	"github.com/codingXiang/go-logger/v2"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

//Orm
type Orm struct {
	*gorm.DB
	ctx    context.Context
	dbConf *gorm.Config
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
	c.Set(configer.GetConfigPath(DB, log, debug), database.Logger.IsDebug())
	c.Set(configer.GetConfigPath(DB, log, slowThreshold), database.Logger.GetSlowThreshold())
	c.Set(configer.GetConfigPath(DB, log, logLevel), database.Logger.GetLevel())
	c.Set(configer.GetConfigPath(DB, log, ignoreRecordNotFoundError), database.Logger.IsIgnoreRecordNotFoundError())
	c.Set(configer.GetConfigPath(DB, log, colorful), database.Logger.IsColorful())
	return New(c)
}

func New(config *viper.Viper, ctx ...context.Context) (*Orm, error) {
	var (
		err error
		orm = new(Orm)
	)
	orm.dbConf = &gorm.Config{}
	if len(ctx) > 0 {
		orm.ctx = ctx[0]
	} else {
		orm.ctx = context.Background()
	}

	if configer.Config == nil {
		//初始化 configer
		configer.Config = configer.NewConfiger()
	}
	var (
		debugMode = config.GetBool(configer.GetConfigPath(DB, log, debug))
	)
	//設定資料庫型態 (MySQL, PostgreSQL) 與連線資訊
	logger.Log.Debug("setup database type")
	err = orm.setDBType(config)

	//設定資料庫參數
	logger.Log.Debug("setup database config")
	orm.setDbConfig(config)

	logger.Log.Debug("setup database logger")
	orm.setLogger(config)

	//設定是否開啟 Log 模式
	logger.Log.Debug("enable debug mode =", debug)
	if debugMode {
		orm.DB = orm.Debug()
	}

	return orm, err
}

func (orm *Orm) ShowVersion() (version string) {
	db, _ := orm.DB.DB()
	db.QueryRow("SELECT VERSION()").Scan(&version)
	return
}

func (orm *Orm) CloseDB() {
	logger.Log.Debug("close database connection")
	defer orm.CloseDB()
}

//CheckTable : 檢查 Table 是否存在，不存在建立並回傳 false, 反之回傳 true
func (orm *Orm) CheckTable(migrate bool, value interface{}) error {
	var (
		err error
		tx  = orm.Begin()
	)
	if !orm.Migrator().HasTable(value) && migrate {
		if err = tx.Set("gorm:table_options", "ENGINE=InnoDB  DEFAULT CHARSET=utf8 AUTO_INCREMENT=1;").AutoMigrate(value); err != nil {
			fmt.Println(err)
			tx.Rollback()
			return err
		}
		tx.Commit()
	} else {
		if migrate {
			if err = tx.Set("gorm:table_options", "ENGINE=InnoDB  DEFAULT CHARSET=utf8 AUTO_INCREMENT=1;").AutoMigrate(value); err != nil {
				tx.Rollback()
				return err
			}
			tx.Commit()
		}
	}
	return nil
}

func (orm *Orm) setLogger(config *viper.Viper) error {
	orm.dbConf.Logger = NewLogger(config).ToGormLogger()
	return nil
}

func (orm *Orm) Info(msg string) {
	orm.Logger.Info(orm.ctx, msg)
}

func (orm *Orm) Warn(msg string) {
	orm.Logger.Warn(orm.ctx, msg)
}

func (orm *Orm) Error(msg string) {
	orm.Logger.Error(orm.ctx, msg)
}

//setDBType : 設定資料庫型態
func (orm *Orm) setDBType(config *viper.Viper) error {
	var (
		err        error
		url        = config.GetString(configer.GetConfigPath(DB, Url))
		port       = config.GetInt(configer.GetConfigPath(DB, Port))
		dbName     = config.GetString(configer.GetConfigPath(DB, Name))
		username   = config.GetString(configer.GetConfigPath(DB, Username))
		password   = config.GetString(configer.GetConfigPath(DB, Password))
		_type      = NewDatabaseType(config.GetString(configer.GetConfigPath(DB, Type)))
		connectStr string
	)
	logger.Log.Debug("set database type = ", _type)
	switch _type {
	case MySQL:
		connectStr = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true",
			username, password, url, port, dbName)
		break
	case PostgreSQL:
		connectStr = fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
			url, port, dbName, username, password)
		break
	}
	logger.Log.Debug("connection string = ", connectStr)
	orm.DB, err = gorm.Open(mysql.Open(connectStr), orm.dbConf)
	if err != nil {
		logger.Log.Error("connect to database error", err)
		return err
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
	db, _ := orm.DB.DB()
	db.SetMaxIdleConns(maxIdleConns)
	logger.Log.Debug("set max open connections", maxOpenConns)
	db.SetMaxOpenConns(maxOpenConns)
	logger.Log.Debug("set connection max life time", maxLifeTime)
	db.SetConnMaxLifetime(time.Duration(maxLifeTime))
}
