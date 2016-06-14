package models

import (
	"github.com/eaciit/orm"
	"sync"
)

type MasterUnit struct {
	sync.RWMutex
	orm.ModelBase `bson:"-" json:"-"`
	Unit          string `bson:"Unit" json:"Unit"`
}

func (m *MasterUnit) TableName() string {
	return "MasterUnit"
}
