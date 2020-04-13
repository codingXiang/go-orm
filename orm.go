package orm

import (
	. "10.40.42.38/BP05G0/go-orm/model"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"    // mysql
	_ "github.com/jinzhu/gorm/dialects/postgres" //postgresql
	"io/ioutil"
	"log"
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

//NewOrm : 新增 ORM 實例
func NewOrm(config DatabaseInterface) OrmInterface {
	var o = &Orm{
		config: config,
		Error:  nil,
	}
	return o.Init(config)
}

//Init : 初始化 ORM
func (this *Orm) Init(config DatabaseInterface) OrmInterface {
	//設定資料庫型態 (MySQL, PostgreSQL) 與連線資訊
	this.Error = this.setDatabaseType(config)
	//設定資料庫參數
	this.setDbConfig(config)

	//設定是否開啟 Log 模式
	if config.GetLogMode() {
		this.GetInstance().LogMode(true)
		this.SetInstance(this.GetInstance().Debug())
	} else {
		this.GetInstance().LogMode(false)
	}

	//設定預設 Table 前綴字
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return config.GetTablePrefix() + defaultTableName
	}

	//設定版本控制 Schema
	this.Error = this.setVersion(config.GetVersion())
	return this
}

func (this *Orm) CloseDB() {
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
	switch config.GetType() {
	case "mysql":
		var connectStr = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true",
			config.GetUsername(), config.GetPassword(), config.GetURL(), config.GetPort(), config.GetName())
		this.db, err = gorm.Open(config.GetType(), connectStr)
		if err != nil {
			log.Println(err)
			return err
		}
		break
	case "postgres":
		var connectStr = fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
			config.GetURL(), config.GetPort(), config.GetName(), config.GetUsername(), config.GetPassword())
		this.db, err = gorm.Open(config.GetType(), connectStr)
		if err != nil {
			return err
		}
		break
	}
	return nil
}

// 設定資料庫組態
func (this *Orm) setDbConfig(config DatabaseInterface) {
	this.GetInstance().DB().SetMaxIdleConns(config.GetMaxIdelConns())
	this.GetInstance().DB().SetMaxOpenConns(config.GetMaxOpenConns())
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
		if err := tx.Model(&Version{}).Create(&version).Error; err != nil {
			tx.Rollback()
			return err
		} else {
			tx.Commit()
		}
		return nil
	}

	var oldVersion = vs[0]
	if oldVersion.GetVersion() != version.GetVersion() {
		var (
			ov  int
			nv  int
			sql []byte
		)
		//判斷版本是否高於現有版本
		if ov, err = this.transformVersion(oldVersion.GetVersion()); err != nil {
			return err
		}
		if nv, err = this.transformVersion(version.GetVersion()); err != nil {
			return err
		}
		if nv > ov {
			if sql, err = this.loadSQLFile(oldVersion.GetVersion()); err != nil {
				return err
			}
			/*
				執行版本更新 SQL
			 */
			var tx1 = this.GetInstance().Begin()
			if err := tx1.Exec(string(sql)).Error; err != nil {
				tx1.Rollback()
				return err
			} else {
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
			var tx2 = this.GetInstance().Begin()
			if err := tx2.Table(this.GetTableName(version)).Update(data).Error; err != nil {
				tx2.Rollback()
				return err
			} else {
				tx2.Commit()
			}
			return nil
		}
		return nil
	}

	return nil
}

//loadSQLFile 讀取更新的 SQL 檔案
func (this *Orm) loadSQLFile(version string) ([]byte, error) {
	var (
		sql []byte
		err error
	)
	if sql, err = ioutil.ReadFile(this.config.GetUpgradeFilePath() + version + ".sql"); err != nil {
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
