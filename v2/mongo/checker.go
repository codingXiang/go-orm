package mongo

import (
	"encoding/json"
	"sync"
)

type Status int

const (
	Same Status = iota
	Change
	Delete
)

type CheckStatus struct {
	Status Status
	Data   *RawData
}

func NewCheckStatus(s Status, data *RawData) *CheckStatus {
	return &CheckStatus{
		Status: s,
		Data:   data,
	}
}

type Checker struct {
	tag string
	raw string
	mux sync.Mutex
}

func NewChecker(in *RawData) *Checker {
	c := &Checker{
		tag: interface2String(in.Tag),
		raw: interface2String(in.Raw),
	}
	return c
}

func (c *Checker) Check(target *RawData) bool {
	c.mux.Lock()
	defer c.mux.Unlock()
	checkTag := interface2String(target.Tag)
	checkRaw := interface2String(target.Raw)
	if c.tag != checkTag || c.raw != checkRaw {
		return true
	}
	return false
}

func interface2String(in interface{}) string {
	tmp, _ := json.Marshal(in)
	return string(tmp)
}
