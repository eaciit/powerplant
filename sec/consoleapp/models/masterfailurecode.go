package models

import (
	"github.com/eaciit/orm"
	"sync"
)

type MasterFailureCode struct {
	sync.RWMutex
	orm.ModelBase `bson:"-" json:"-"`
	Id            string `bson:"_id" json:"Id"`
	Text          string `bson:"text" json:"Text"`
}

func (m *MasterFailureCode) TableName() string {
	return "MasterFailureCode"
}
