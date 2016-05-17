package models

import (
	"github.com/eaciit/orm"
)

type NewEquipmentType struct {
	orm.ModelBase     `bson:"-",json:"-"`
	EquipmentType     string `bson:"EquipmentType",json:"EquipmentType"`
	EquipmentText     string `bson:"EquipmentText",json:"EquipmentText"`
	NewEquipmentGroup string `bson:"NewEquipmentGroup",json:"NewEquipmentGroup"`
}

func (m *NewEquipmentType) TableName() string {
	return "NewEquipmentType"
}
