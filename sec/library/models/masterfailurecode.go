package models

import (
	"sync"

	"github.com/eaciit/orm"
)

type MasterFailureCode struct {
	sync.RWMutex
	orm.ModelBase `bson:"-" json:"-"`
	Id            string `bson:"_id" json:"Id"`
	Text          string `bson:"text" json:"Text"`
}

func (e *MasterFailureCode) RecordID() interface{} {
	return e.Id
}

func (m *MasterFailureCode) TableName() string {
	return "MasterFailureCode"
}
