package redis

import (
	"fmt"
	"github.com/codingXiang/configer/v2"
	"github.com/codingXiang/go-logger/v2"
	"github.com/codingXiang/go-orm/v2"
	"github.com/ghodss/yaml"
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"strings"
	"time"
)

type RedisClient struct {
	*redis.Client
	prefix string
}

var (
	RedisORM *RedisClient
)

func NewRedis(r *Redis) *RedisClient {
	c := viper.New()
	c.Set(configer.GetConfigPath(REDIS, orm.Url), r.URL)
	c.Set(configer.GetConfigPath(REDIS, orm.Port), r.Port)
	c.Set(configer.GetConfigPath(REDIS, DB), r.DB)
	c.Set(configer.GetConfigPath(REDIS, orm.Password), r.Password)
	c.Set(configer.GetConfigPath(REDIS, Prefix), r.Prefix)
	return New(c)
}

//New : 建立 Redis Client 實例
func New(config *viper.Viper) *RedisClient {
	var rc = new(RedisClient)
	//讀取 config
	var (
		url      = config.GetString(configer.GetConfigPath(REDIS, orm.Url))
		port     = config.GetInt(configer.GetConfigPath(REDIS, orm.Port))
		password = config.GetString(configer.GetConfigPath(REDIS, orm.Password))
		db       = config.GetInt(configer.GetConfigPath(REDIS, DB))
		prefix   = config.GetString(configer.GetConfigPath(REDIS, Prefix))
	)
	//設定連線資訊
	option := &redis.Options{
		Addr: fmt.Sprintf("%s:%d", url, port),
		DB:   db,
	}
	rc.prefix = prefix
	if password != "" {
		option.Password = password
	}
	rc.Client = redis.NewClient(option)
	logger.Log.Debug("check redis ...", rc.Client)
	return rc
	//_, err = rc.GetInfo()
	//if err != nil {
	//	errMsg := "redis connect error"
	//	logger.Log.Error(errMsg, err)
	//	return nil, err
	//} else {
	//	logger.Log.Info("redis connect success")
	//	return rc, nil
	//}
}

//GetRedisInfo 取得 Redis 資訊
func (r *RedisClient) GetInfo() (string, error) {
	return r.Client.Ping().Result()
}

//SetKeyValueWithExpire : 設定 Key 與 Value
func (r *RedisClient) SetKeyValue(key string, value interface{}, expiration time.Duration) error {
	err := r.Client.Set(r.prefix+key, value, expiration).Err()
	return err
}

//GetValue : 取得 Key 的 Value
func (r *RedisClient) GetValue(key string) (string, error) {
	val := r.Client.Get(r.prefix + key)
	return val.Val(), val.Err()
}

//RemoveKey : 刪除 Key
func (r *RedisClient) RemoveKey(key string) error {
	return r.Client.Del(r.prefix + key).Err()
}

//Publish : 發佈
func (r *RedisClient) Publish(channel string, message interface{}) error {
	return r.Client.Publish(channel, message).Err()
}

//Subscribe : 訂閱
func (r *RedisClient) Subscribe(channel string) *redis.PubSub {
	return r.Client.Subscribe(channel)
}

func (r *RedisClient) Info(sections ...string) map[string]map[string]interface{} {
	var (
		result = make(map[string]map[string]interface{})
	)
	for _, section := range sections {
		var (
			sectionData string
			sectionMap  map[string]interface{}
			err         error
		)
		r.Client.Info(section).Scan(&sectionData)
		if sectionMap, err = r.parseSection(sectionData); err == nil {
			result[section] = sectionMap
		}
	}
	return result
}

func (r *RedisClient) parseSection(data string) (result map[string]interface{}, err error) {
	data = strings.ReplaceAll(data, ":", ": ")
	dataTmp := strings.Split(data, "\r\n")
	data = strings.Join(dataTmp[1:], "\n")
	err = yaml.Unmarshal([]byte(data), &result)
	fmt.Println("")
	return
}
