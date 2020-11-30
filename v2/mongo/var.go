package mongo

import (
	"github.com/hashicorp/go-uuid"
	"gopkg.in/mgo.v2/bson"
)

const (
	MONGO    = "mongo"
	MODE     = "mode"
	DATABASE = "database"
)

const (
	IDENTITY = "identity"
	RAW      = "raw"
	TAG      = "tag"
)

const (
	Primary            = "Primary"            // Default mode. All operations read from the current replica set primary.
	PrimaryPreferred   = "PrimaryPreferred"   // Read from the primary if available. Read from the secondary otherwise.
	Secondary          = "Secondary"          // Read from one of the nearest secondary members of the replica set.
	SecondaryPreferred = "SecondaryPreferred" // Read from one of the nearest secondaries if available. Read from primary otherwise.
	Nearest            = "Nearest"            // Read from one of the nearest members, irrespective of it being primary or secondary.
	Eventual           = "Eventual"           // Same as Nearest, but may change servers between reads.
	Monotonic          = "Monotonic"          // Same as SecondaryPreferred before first write. Same as Primary after first write.
	Strong             = "Strong"             // Same as Primary.
)

type Mongo struct {
	URL        string
	Port       int
	Username   string
	Password   string
	Mode       string
	Database   string
	Collection string
}

func NewSearchCondition(id string, identity string, tag bson.M, data bson.M) bson.M {
	if data == nil {
		data = make(bson.M)
	}
	if id != "" {
		data[_id] = bson.ObjectIdHex(id)
	}

	if identity != "" {
		data[IDENTITY] = identity
	}

	if tag != nil {
		data[TAG] = tag
	}
	return data
}

type RawData struct {
	Identity string      `json:"identity"`
	Tag      bson.M      `json:"tag"`
	Raw      interface{} `json:"raw"`
}

func NewRawData(id string, tag map[string]interface{}, raw interface{}) *RawData {
	if id == "" {
		id, _ = uuid.GenerateUUID()
	}
	if tag == nil {
		tag = make(map[string]interface{})
	}
	out := new(RawData)
	out.Identity = id
	out.Tag = tag
	out.Raw = raw
	return out
}

func (r *RawData) AddTag(key string, value interface{}) *RawData {
	r.Tag[key] = value
	return r
}

func (r *RawData) GetTag() map[string]interface{} {
	return r.Tag
}

func (r *RawData) GetIdentity() string {
	return r.Identity
}

func (r *RawData) GetRaw() interface{} {
	return r.Raw
}

func Default() *Mongo {
	return &Mongo{
		URL:      "127.0.0.1",
		Mode:     "Primary",
		Username: "",
		Password: "",
		Database: "admin",
		Port:     27017,
	}
}
