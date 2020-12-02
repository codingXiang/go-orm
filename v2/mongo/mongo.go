package mongo

import (
	"encoding/json"
	"github.com/codingXiang/configer/v2"
	"github.com/codingXiang/go-orm/v2"
	"github.com/spf13/viper"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"sync"
	"time"
)

type Client struct {
	*mgo.Session
	*mgo.Database
	*mgo.Collection
	lock sync.Mutex
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
	c.Set(configer.GetConfigPath(MONGO, DATABASE), r.Database)
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
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.Collection.Insert(data)
}

func (c *Client) First(selector bson.M) (*RawData, error) {
	target := new(RawData)
	err := c.Collection.Find(selector).One(&target)
	return target, err
}

func (c *Client) Find(selector bson.M, queryCondition *QueryCondition) ([]*RawData, error) {
	target := make([]*RawData, 0)
	var q = c.Collection.Find(selector)
	if queryCondition.Sort != nil {
		q = q.Sort(queryCondition.Sort...)
	}
	if queryCondition.Limit != 0 {
		q = q.Limit(queryCondition.Limit)
	}
	err := q.All(&target)
	return target, err
}

func (c *Client) Update(selector bson.M, data interface{}) (*mgo.ChangeInfo, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.Collection.Upsert(selector, data)
}

func (c *Client) Delete(selector bson.M) error {
	c.lock.Lock()
	defer c.lock.Unlock()
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

func (c *Client) WaitForChange(collection string, selector bson.M, onChange func(data *RawData) bool, onDelete func()) error {
	data, err := c.C(collection).First(selector)
	if err != nil {
		return err
	}
	tmp, _ := json.Marshal(data.Tag)
	originTag := string(tmp)
	tmp, _ = json.Marshal(data.Raw)
	originRaw := string(tmp)
CHECK:
	for {
		select {
		case <-time.Tick(1 * time.Second):
			check, err1 := c.First(selector)
			if check.Identity == "" {
				onDelete()
				break CHECK
			} else {
				tmp, _ = json.Marshal(check.Tag)
				checkTag := string(tmp)
				tmp, _ = json.Marshal(check.Raw.(bson.M))
				checkRaw := string(tmp)
				if err1 == nil {
					if originTag != checkTag || originRaw != checkRaw {
						if onChange(check) {
							break CHECK
						} else {
							originRaw = checkRaw
							originTag = checkTag
						}
					} else {
						originRaw = checkRaw
						originTag = checkTag
					}
				} else {
					err = err1
					break CHECK
				}
			}

		}
	}
	return err
}
