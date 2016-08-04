package models

import (
	"sync"

	"github.com/eaciit/orm"
)

type GenVEDNotInTmp struct {
	sync.RWMutex
	orm.ModelBase `bson:"-" json:"-"`
	//Id            string `bson:"Id" json:"Id"`
	WorkOrderID string `bson:"WorkOrderID" json:"WorkOrderID"`
}

func (m *GenVEDNotInTmp) TableName() string {
	return "GenVEDNotInTmp"
}
