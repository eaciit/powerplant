package models

import (
	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
)

type MasterEquipmentType struct {
	orm.ModelBase     `bson:"-",json:"-"`
	Id                bson.ObjectId ` bson:"_id" , json:"_id" `
	EquipmentType     string        ` bson:"EquipmentType" , json:"EquipmentType" `
	EquipmentTypeDesc string        ` bson:"EquipmentTypeDesc" , json:"EquipmentTypeDesc" `
}

func (e *MasterEquipmentType) RecordID() interface{} {
	return e.Id
}

func (m *MasterEquipmentType) TableName() string {
	return "MasterEquipmentType"
}
