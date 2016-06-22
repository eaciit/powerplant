package models

import (
	"sync"

	"github.com/eaciit/orm"
)

type MasterUnit struct {
	sync.RWMutex
	orm.ModelBase `bson:"-" json:"-"`
	Unit          string `bson:"Unit" json:"Unit"`
}

func (m *MasterUnit) TableName() string {
	return "MasterUnit"
}
