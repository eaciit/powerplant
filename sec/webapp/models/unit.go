package models

import (
	"github.com/eaciit/orm"
	// "gopkg.in/mgo.v2/bson"
)

type UnitModel struct {
	orm.ModelBase `bson:"-",json:"-"`
	Id            int    `bson:"Id"  json:"Id" `
	Unit          string `bson:"Unit"  json:"Unit" `
}

func NewUnitModel() *UnitModel {
	m := new(UnitModel)
	// m.Id = bson.NewObjectId()
	return m
}

func (e *UnitModel) RecordID() interface{} {
	return e.Id
}

func (m *UnitModel) TableName() string {
	return "MasterUnit"
}
