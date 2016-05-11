package models

import (
	"github.com/eaciit/orm"
)

type MasterActivityType struct {
	orm.ModelBase `bson:"-",json:"-"`
	Id            string ` bson:"_id" , json:"_id" `
}

func (e *MasterActivityType) RecordID() interface{} {
	return e.Id
}

func (m *MasterActivityType) TableName() string {
	return "MasterActivityType"
}
