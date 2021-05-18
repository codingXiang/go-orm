package orm

import (
	"github.com/8treenet/gcache"
	"github.com/codingXiang/configer/v2"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)

type CacheHandler struct {
	DB    *gorm.DB
	//redis *redis.Redis
	gcache.Plugin
}

func NewCacheHandler(db *gorm.DB, config *viper.Viper) *CacheHandler {
	c := new(CacheHandler)
	opt := c.initOption(config)
	c.Plugin = gcache.AttachDB(c.DB, &opt, &gcache.RedisOption{Addr: "localhost:6379"})

	return c
}

func (c *CacheHandler) initOption(config *viper.Viper) gcache.DefaultOption {
	var (
		exipres         = config.GetInt(configer.GetConfigPath(DB, CACHE, EXPIRE))
		asyncWrite      = config.GetBool(configer.GetConfigPath(DB, CACHE, ASYNC_WRITE))
		//debug           = config.GetBool(configer.GetConfigPath(DB, CACHE, DEBUG))
		//level           = config.GetString(configer.GetConfigPath(DB, CACHE, LEVEL))
		penetrationSafe = config.GetBool(configer.GetConfigPath(DB, CACHE, PENETRATIONSAFE))
	)
	opt := gcache.DefaultOption{}
	opt.Expires = exipres                 //缓存时间，默认60秒。范围 30-900
	opt.Level = gcache.LevelDisable       //缓存级别，默认LevelSearch。LevelDisable:关闭缓存，LevelModel:模型缓存， LevelSearch:查询缓存
	opt.AsyncWrite = asyncWrite           //异步缓存更新, 默认false。 insert update delete 成功后是否异步更新缓存
	opt.PenetrationSafe = penetrationSafe //开启防穿透, 默认false。
	return opt
}