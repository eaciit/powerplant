package models

import (
	"sync"

	"github.com/eaciit/orm"
)

type MasterPlant struct {
	sync.RWMutex
	orm.ModelBase `bson:"-",json:"-"`
	Plant         string `bson:"Plant",json:"Plant"`
}

func (m *MasterPlant) TableName() string {
	return "MasterPlant"
}
