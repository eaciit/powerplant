package models

import (
	"sync"

	"github.com/eaciit/orm"
)

type MasterEquipmentType struct {
	sync.RWMutex
	orm.ModelBase     `bson:"-" json:"-"`
	EquipmentType     string `bson:"EquipmentType" json:"EquipmentType"`
	EquipmentTypeDesc string `bson:"EquipmentTypeDesc" json:"EquipmentTypeDesc"`
}

func (m *MasterEquipmentType) TableName() string {
	return "MasterEquipmentType"
}
