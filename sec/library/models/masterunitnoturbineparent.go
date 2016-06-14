package models

import (
	"github.com/eaciit/orm"
	"sync"
)

type MasterUnitNoTurbineParent struct {
	sync.RWMutex
	orm.ModelBase `bson:"-" json:"-"`
	// Id            string `bson:"Id" json:"Id"`
	UnitNo        string `bson:"UnitNo" json:"UnitNo"`
	TurbineParent string `bson:"TurbineParent" json:"TurbineParent"`
}

func (m *MasterUnitNoTurbineParent) TableName() string {
	return "MasterUnitNoTurbineParent"
}
