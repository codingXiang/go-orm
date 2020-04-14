package orm

import (
	"encoding/json"
	"fmt"
	"github.com/codingXiang/go-logger"
	. "github.com/codingXiang/go-orm/model"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"    // mysql
	_ "github.com/jinzhu/gorm/dialects/postgres" //postgresql
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

//Orm
type (
	OrmInterface interface {
		Init(config DatabaseInterface) OrmInterface
		CloseDB()
		GetTableName(value interface{}) string
		CheckTable(migrate bool, value interface{}) error
		GetInstance() *gorm.DB
		SetInstance(db *gorm.DB)
		Upgrade() error
	}
	Orm struct {
		db     *gorm.DB
		config DatabaseInterface
		Error  error
	}
)

var (
	DatabaseORM OrmInterface
)

func InterfaceToDatabase(data interface{}) DatabaseInterface {
	var result = &Database{}
	if jsonStr, err := json.Marshal(data); err == nil {
		json.Unmarshal(jsonStr, &result)
	}
	return result
}

//NewOrm : 新增 ORM 實例
func NewOrm(config DatabaseInterface) {
	var o = &Orm{
		config: config,
		Error:  nil,
	}
	DatabaseORM = o.Init(config)
}

//Init : 初始化 ORM
func (this *Orm) Init(config DatabaseInterface) OrmInterface {
	//設定資料庫型態 (MySQL, PostgreSQL) 與連線資訊
	logger.Log.Debug("setup database type")
	this.Error = this.setDatabaseType(config)

	//設定資料庫參數
	logger.Log.Debug("setup database config")
	this.setDbConfig(config)

	//設定是否開啟 Log 模式
	logger.Log.Debug("setup log mode =", config.GetLogMode())
	if config.GetLogMode() {
		this.GetInstance().LogMode(true)
		this.SetInstance(this.GetInstance().Debug())
	} else {
		this.GetInstance().LogMode(false)
	}

	//設定預設 Table 前綴字
	logger.Log.Debug("setup table name prefix", config.GetTablePrefix())
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return config.GetTablePrefix() + defaultTableName
	}

	//設定版本控制 Schema
	logger.Log.Debug("setup version")
	this.Error = this.setVersion(config.GetVersion())
	return this
}

func (this *Orm) CloseDB() {
	logger.Log.Debug("close database connection")
	defer this.GetInstance().Close()
}

//GetTableName : 透過傳入 struct 回傳 table 名稱
func (this *Orm) GetTableName(tb interface{}) string {
	return this.GetInstance().NewScope(tb).TableName()
}

//CheckTable : 檢查 Table 是否存在，不存在建立並回傳 false, 反之回傳 true
func (this *Orm) CheckTable(migrate bool, value interface{}) error {
	var (
		err error
		tx  = this.GetInstance().Begin()
	)
	if !this.GetInstance().HasTable(value) {
		if err = tx.CreateTable(value).Error; err != nil {
			fmt.Println(err)
			tx.Rollback()
			return err
		}
		tx.Commit()
	} else {
		if migrate {
			if err = tx.AutoMigrate(value).Error; err != nil {
				tx.Rollback()
				return err
			}
			tx.Commit()
		}
	}
	return nil
}

func (orm *Orm) GetInstance() *gorm.DB {
	return orm.db
}

func (this *Orm) SetInstance(db *gorm.DB) {
	this.db = db
}

//setDatabaseType : 設定資料庫型態
func (this *Orm) setDatabaseType(config DatabaseInterface) error {
	var err error
	logger.Log.Debug("set database type = ", config.GetType())
	switch config.GetType() {
	case "mysql":
		var connectStr = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true",
			config.GetUsername(), config.GetPassword(), config.GetURL(), config.GetPort(), config.GetName())
		logger.Log.Debug("connection string = ", connectStr)
		this.db, err = gorm.Open(config.GetType(), connectStr)
		if err != nil {
			logger.Log.Error("connect to database error", err)
			return err
		}
		break
	case "postgres":
		var connectStr = fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
			config.GetURL(), config.GetPort(), config.GetName(), config.GetUsername(), config.GetPassword())
		logger.Log.Debug("connection string = ", connectStr)
		this.db, err = gorm.Open(config.GetType(), connectStr)
		if err != nil {
			logger.Log.Error("connect to database error", err)
			return err
		}
		break
	}
	return nil
}

// 設定資料庫組態
func (this *Orm) setDbConfig(config DatabaseInterface) {
	logger.Log.Debug("set max idle connections", config.GetMaxIdelConns())
	this.GetInstance().DB().SetMaxIdleConns(config.GetMaxIdelConns())
	logger.Log.Debug("set max open connections", config.GetMaxOpenConns())
	this.GetInstance().DB().SetMaxOpenConns(config.GetMaxOpenConns())
	logger.Log.Debug("set connection max life time", config.GetMaxLifeTime())
	this.GetInstance().DB().SetConnMaxLifetime(time.Duration(config.GetMaxLifeTime()))
}

//設定資料庫版本
func (this *Orm) setVersion(version *Version) (error) {
	var (
		err error
		vs  []*Version
		tx  = this.GetInstance().Begin()
	)

	if err = tx.Model(&Version{}).Find(&vs).Error; err != nil {
		logger.Log.Debug("not found version table")
		if err = this.CheckTable(false, &version); err != nil {
			fmt.Println(err)
			return err
		}
		if err := tx.Model(&Version{}).Create(&version).Error; err != nil {
			tx.Rollback()
			return err
		} else {
			tx.Commit()
		}
		return nil
	} else if len(vs) < 1 {
		logger.Log.Debug("not found version record")
		if err := tx.Model(&Version{}).Create(&version).Error; err != nil {
			tx.Rollback()
			return err
		} else {
			tx.Commit()
		}
		return nil
	}
	return err
}

func (this *Orm) Upgrade() error {
	var (
		err     error
		vs      []*Version
		version = this.config.GetVersion()
		tx      = this.GetInstance().Begin()
	)
	/*
		檢查是否可以更新
	 */
	if err = tx.Model(&Version{}).Find(&vs).Error; err != nil {
		logger.Log.Debug("not found version table")
		if err = this.CheckTable(false, &version); err != nil {
			fmt.Println(err)
			return err
		}
		if err := tx.Model(&Version{}).Create(&version).Error; err != nil {
			tx.Rollback()
			return err
		} else {
			tx.Commit()
		}
		return nil
	} else if len(vs) < 1 {
		logger.Log.Debug("not found version record")
		if err := tx.Model(&Version{}).Create(&version).Error; err != nil {
			tx.Rollback()
			return err
		} else {
			tx.Commit()
		}
		return nil
	}
	var oldVersion = vs[0]
	logger.Log.Debug("found version", oldVersion.GetVersion())
	if oldVersion.GetVersion() != version.GetVersion() {
		var (
			ov  int
			nv  int
			sql []byte
		)
		//判斷版本是否高於現有版本
		if ov, err = this.transformVersion(oldVersion.GetVersion()); err != nil {
			logger.Log.Error("transform old version error", err)
			return err
		}
		if nv, err = this.transformVersion(version.GetVersion()); err != nil {
			logger.Log.Error("transform new version error", err)
			return err
		}
		if nv > ov {
			logger.Log.Debug("can upgrade")
			if sql, err = this.loadSQLFile(oldVersion.GetVersion()); err != nil {
				return err
			}
			/*
				執行版本更新 SQL
			 */
			logger.Log.Debug("start execute upgrade sql")
			var tx1 = this.GetInstance().Begin()
			if err := tx1.Exec(string(sql)).Error; err != nil {
				logger.Log.Error("upgrade sql execute failed", err)
				tx1.Rollback()
				return err
			} else {
				logger.Log.Debug("upgrade sql execute success")
				tx1.Commit()
			}

			/*
				執行 Version Table 更新
			 */
			var (
				data map[string]interface{} = map[string]interface{}{
					"version": version.GetVersion(),
				}
			)
			logger.Log.Debug("start upgrade version record")
			var tx2 = this.GetInstance().Begin()
			if err := tx2.Table(this.GetTableName(version)).Update(data).Error; err != nil {
				logger.Log.Error("upgrade version record failed", err)
				tx2.Rollback()
				return err
			} else {
				logger.Log.Debug("upgrade version record success")
				tx2.Commit()
			}
			return nil
		}
		return nil
	}
	logger.Log.Debug("can not upgrade")

	return nil
}

//loadSQLFile 讀取更新的 SQL 檔案
func (this *Orm) loadSQLFile(version string) ([]byte, error) {
	var (
		file = version + ".sql"
		sql  []byte
		err  error
	)
	logger.Log.Debug("read sql file", file)
	if sql, err = ioutil.ReadFile(this.config.GetUpgradeFilePath() + file); err != nil {
		return nil, err
	}
	return sql, nil
}

//translformVersion 轉換版本權重
func (this *Orm) transformVersion(version string) (int, error) {
	var (
		tmp        []string
		err        error
		v1, v2, v3 int
		result     int = 0
	)
	tmp = strings.Split(version, ".")
	if v1, err = strconv.Atoi(tmp[0]); err != nil {
		return 0, err
	}
	if v1, err = strconv.Atoi(tmp[1]); err != nil {
		return 0, err
	}
	if v1, err = strconv.Atoi(tmp[2]); err != nil {
		return 0, err
	}

	result = v1*100 + v2*10 + v3

	return result, err
}
