package models

import (
	"sync"

	"github.com/eaciit/orm"
)

type MappedEquipmentType struct {
	sync.RWMutex
	orm.ModelBase            `bson:"-",json:"-"`
	EquipmentType            string `bson:"EquipmentType",json:"EquipmentType"`
	EquipmentText            string `bson:"EquipmentText",json:"EquipmentText"`
	EquipmentTypeMapped3Char string `bson:"EquipmentTypeMapped3Char",json:"EquipmentTypeMapped3Char"`
	EquipmentTypeMapped4Char string `bson:"EquipmentTypeMapped4Char",json:"EquipmentTypeMapped4Char"`
	EquipmentTypeMappedText  string `bson:"EquipmentTypeMappedText",json:"EquipmentTypeMappedText"`
}

func (m *MappedEquipmentType) TableName() string {
	return "MappedEquipmentType"
}
