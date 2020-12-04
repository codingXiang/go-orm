package mongo

import (
	"context"
	"errors"
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

func (c *Client) check(ctx context.Context, in *RawData, listenFunc func() (*RawData, error), ch chan *CheckStatus) error {
	checker := NewChecker(in)
	for {
		select {
		case <-ctx.Done():
			return errors.New("execute timeout")
		case <-time.Tick(1 * time.Second):
			target, err := listenFunc()
			if err != nil {
				ch <- NewCheckStatus(Delete, nil)
			}
			if checker.Check(target) {
				checker = NewChecker(target)
				ch <- NewCheckStatus(Change, target)
			} else {
				ch <- NewCheckStatus(Same, target)
			}
		}
	}
}

func (c *Client) WaitForChange(ctx context.Context, listenFunc func() (*RawData, error), onChange func(data *RawData) (stop bool, err error), onDelete func()) error {
	ct, cancel := context.WithCancel(ctx)
	defer cancel()
	data, err := listenFunc()
	if err != nil {
		return err
	}
	ch := make(chan *CheckStatus)
	go c.check(ct, data, listenFunc, ch)
	for {
		select {
		case <-ct.Done():
			return errors.New("execute timeout")
		case s := <-ch:
			switch s.Status {
			case Delete:
				onDelete()
				return errors.New("object has been delete")
			case Change:
				if stop, e := onChange(s.Data); stop {
					return e
				}
			}
		}
	}
	//originTag := interface2String(data.Tag)
	//originRaw := interface2String(data.Raw)
	//for {
	//	select {
	//	case <-ctx.Done():
	//		context.WithValue(ctx, "err", ctx.Err())
	//		log.Println("timeout")
	//		return ctx.Err()
	//	case <-time.Tick(1 * time.Second):
	//		check, err1 := listenFunc()
	//		if check.Identity == "" {
	//			onDelete()
	//			return nil
	//		} else {
	//			checkTag := interface2String(check.Tag)
	//			checkRaw := interface2String(check.Raw)
	//			if err1 == nil {
	//				if originTag != checkTag || originRaw != checkRaw {
	//					stop, e := onChange(check)
	//					if stop {
	//						return e
	//					} else {
	//						originRaw = checkRaw
	//						originTag = checkTag
	//					}
	//				} else {
	//					originRaw = checkRaw
	//					originTag = checkTag
	//				}
	//			} else {
	//				return err1
	//			}
	//		}
	//	}
	//}
}
