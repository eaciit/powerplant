package models

import (
	"sync"

	"github.com/eaciit/orm"
)

type MasterOrderType struct {
	sync.RWMutex
	orm.ModelBase `bson:"-" json:"-"`
	Type          string `bson:"Type" json:"Type"`
	Description   string `bson:"FLDescription" json:"Description"`
}

func (m *nc) TableName() nc {
	return "MasterOrderType"
}
