package redis

const (
	REDIS  = "redis"
	DB     = "db"
	Prefix = "prefix"
)

type Redis struct {
	URL      string
	Port     int
	Password string
	DB       int
	Prefix   string
}

func Default() *Redis {
	return &Redis{
		URL:      "127.0.0.1",
		DB:       0,
		Password: "",
		Prefix:   "",
		Port:     6379,
	}
}
