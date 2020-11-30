package mongo

import (
	"github.com/codingXiang/configer/v2"
	"github.com/codingXiang/go-orm/v2"
	"github.com/spf13/viper"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Client struct {
	*mgo.Session
	*mgo.Database
	*mgo.Collection
}

const (
	_id = "_id"
)

var (
	MongoClient *Client
)

func NewMongo(r *Mongo) *Client {
	c := viper.New()
	c.Set(configer.GetConfigPath(MONGO, orm.Url), r.URL)
	c.Set(configer.GetConfigPath(MONGO, orm.Port), r.Port)
	c.Set(configer.GetConfigPath(MONGO, orm.Username), r.Username)
	c.Set(configer.GetConfigPath(MONGO, orm.Password), r.Password)
	c.Set(configer.GetConfigPath(MONGO, MODE), r.Mode)
	return New(c)
}

func New(config *viper.Viper) *Client {
	var (
		url      = config.GetString(configer.GetConfigPath(MONGO, orm.Url))
		port     = config.GetString(configer.GetConfigPath(MONGO, orm.Port))
		username = config.GetString(configer.GetConfigPath(MONGO, orm.Username))
		password = config.GetString(configer.GetConfigPath(MONGO, orm.Password))
		mode     = config.GetString(configer.GetConfigPath(MONGO, MODE))
		database = config.GetString(configer.GetConfigPath(MONGO, DATABASE))
		err      error
	)
	client := new(Client)
	client.Session, err = mgo.Dial(url + ":" + port)
	if err != nil {
		panic(err)
	}
	cred := client.getCredential(username, password)
	if cred != nil {
		err = client.Session.Login(cred)
		if err != nil {
			panic(err)
		}
	}

	client.SetMode(client.getSessionMode(mode), true)
	client.SetDB(database)
	return client
}

func (c *Client) SetDB(name string) *mgo.Database {
	c.Database = c.Session.DB(name)
	return c.Database
}

func (c *Client) C(name string) *Client {
	c.Collection = c.Database.C(name)
	return c
}

func (c *Client) Insert(data *RawData) error {
	return c.Collection.Insert(data)
}

func (c *Client) First(selector SearchCondition) (*RawData, error) {
	target := new(RawData)
	err := c.Collection.Find(selector).One(&target)
	return target, err
}

func (c *Client) Find(selector SearchCondition) ([]*RawData, error) {
	target := make([]*RawData, 0)
	err := c.Collection.Find(selector).All(&target)
	return target, err
}

func (c *Client) Update(selector bson.M, data interface{}) (*mgo.ChangeInfo, error) {
	return c.Collection.Upsert(selector, data)
}

func (c *Client) Delete(selector SearchCondition) error {
	return c.Collection.Remove(selector)
}

func (c *Client) getCredential(username, password string) *mgo.Credential {
	if username == "" || password == "" {
		return nil
	}
	return &mgo.Credential{
		Username: username,
		Password: password,
	}
}

func (c *Client) getSessionMode(mode string) mgo.Mode {
	switch mode {
	case Primary:
		return mgo.Primary
	case PrimaryPreferred:
		return mgo.PrimaryPreferred
	case Secondary:
		return mgo.Secondary
	case SecondaryPreferred:
		return mgo.SecondaryPreferred
	case Nearest:
		return mgo.Nearest
	case Eventual:
		return mgo.Eventual
	case Monotonic:
		return mgo.Monotonic
	case Strong:
		return mgo.Strong
	default:
		return mgo.Primary
	}
}
