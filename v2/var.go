package orm

type DatabaseType int

const (
	MySQL DatabaseType = iota
	PostgreSQL
)

const (
	_mysql      = "mysql"
	_postgrseql = "postgresql"
)

func NewDatabaseType(t string) DatabaseType {
	switch t {
	case _mysql:
		return MySQL
	case _postgrseql:
		return PostgreSQL
	default:
		return MySQL
	}
}

func (t DatabaseType) String() string {
	switch t {
	case MySQL:
		return _mysql
	case PostgreSQL:
		return _postgrseql
	default:
		return _mysql
	}
}

const (
	DB           = "database"
	Url          = "url"
	Port         = "port"
	Name         = "name"
	Username     = "username"
	Password     = "password"
	Type         = "type"
	LogMode      = "logMode"
	TablePrefix  = "tablePrefix"
	Version      = "version"
	MaxOpenConns = "maxOpenConns"
	MaxIdleConns = "maxIdleConns"
	MaxLifeTime  = "maxLifeTime"
)

const (
	CACHE           = "cache"
	EXPIRE          = "expire"
	LEVEL           = "level"
	DEBUG           = "debug"
	ASYNC_WRITE     = "asyncWrite"
	PENETRATIONSAFE = "penetrationSafe"
)
