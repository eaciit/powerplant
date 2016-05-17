package models

import "github.com/eaciit/orm"

type MasterEquipmentType struct {
	orm.ModelBase     `bson:"-",json:"-"`
	EquipmentType     string `bson:"EquipmentType",json:"EquipmentType"`
	EquipmentTypeDesc string `bson:"EquipmentTypeDesc",json:"EquipmentTypeDesc"`
}

func (m *MasterEquipmentType) TableName() string {
	return "MasterEquipmentType"
}
