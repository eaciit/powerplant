package models

import (
	"github.com/eaciit/orm"
	"sync"
)

type MasterUnit struct {
	sync.RWMutex
	orm.ModelBase `bson:"-" json:"-"`
	// Id            string `bson:"Id" json:"Id"`
	Unit string `bson:"Unit" json:"Unit"`
}

func (m *MasterUnit) TableName() string {
	return "MasterUnit"
}
