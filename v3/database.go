package orm

import (
	"github.com/codingXiang/configer/v2"
	"github.com/spf13/viper"
	"gorm.io/gorm/logger"
	_log "log"
	"os"
	"time"
)

type Logger struct {
	Debug                     bool            `json:"debug"`
	SlowThreshold             time.Duration   `json:"slowThreshold"`
	Level                     logger.LogLevel `json:"level"`
	IgnoreRecordNotFoundError bool            `json:"ignoreRecordNotFoundError"`
	Colorful                  bool            `json:"colorful"`
}

func NewLogger(config *viper.Viper) *Logger {
	return &Logger{
		Debug:                     config.GetBool(configer.GetConfigPath(DB, log, debug)),
		SlowThreshold:             config.GetDuration(configer.GetConfigPath(DB, log, slowThreshold)),
		Level:                     logger.LogLevel(config.GetInt(configer.GetConfigPath(DB, log, LEVEL))),
		IgnoreRecordNotFoundError: config.GetBool(configer.GetConfigPath(DB, log, ignoreRecordNotFoundError)),
		Colorful:                  config.GetBool(configer.GetConfigPath(DB, log, colorful)),
	}
}

func (l *Logger) IsDebug() bool {
	return l.Debug
}

//GetSlowThreshold 取得慢查詢的門檻值
func (l *Logger) GetSlowThreshold() time.Duration {
	return l.SlowThreshold
}

//GetLevel 取得 log 等級
func (l *Logger) GetLevel() logger.LogLevel {
	return l.Level
}

//IsIgnoreRecordNotFoundError 是否要隱藏查詢不到紀錄的錯誤
func (l *Logger) IsIgnoreRecordNotFoundError() bool {
	return l.IgnoreRecordNotFoundError
}

//IsColorful log 是否要有顏色
func (l *Logger) IsColorful() bool {
	return l.Colorful
}

func (l *Logger) ToGormLogger() logger.Interface {
	return logger.New(
		_log.New(os.Stdout, "\r\n", _log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             l.GetSlowThreshold(),            // Slow SQL threshold
			LogLevel:                  l.GetLevel(),                    // Log level
			IgnoreRecordNotFoundError: l.IsIgnoreRecordNotFoundError(), // Ignore ErrRecordNotFound error for logger
			Colorful:                  l.IsColorful(),                  // Disable color
		},
	)
}

type Database struct {
	URL          string  `yaml:"url"`          //Server 的位置
	Name         string  `yaml:"name"`         //名稱
	Port         int     `yaml:"port"`         //Port
	LogMode      bool    `yaml:"logMode"`      //Log模式
	Username     string  `yaml:"username"`     //使用者名稱
	Password     string  `yaml:"password"`     //密碼
	Type         string  `yaml:"type"`         //類型（例如 mysql、postgre、sqlite等)
	TablePrefix  string  `yaml:"tablePrefix"`  //table前綴字
	MaxOpenConns int     `yaml:"maxOpenConns"` //最大開啟連線數
	MaxIdleConns int     `yaml:"maxIdleConns"` //最大連線
	MaxLifeTime  int     `yaml:"maxLifeTime"`  //最長連線時間
	Logger       *Logger `yaml:"logger"`       //log 設定
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
		Logger: &Logger{
			Debug:                     true,
			SlowThreshold:             time.Second,
			Level:                     logger.Silent,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
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

//GetLogger 取得 Database logger 設定
func (db *Database) GetLogger() *Logger {
	return db.Logger
}
