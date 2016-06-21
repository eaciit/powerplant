package models

import (
	"github.com/eaciit/orm"
	tk "github.com/eaciit/toolkit"
	"gopkg.in/mgo.v2/bson"
)

type Master struct {
}

func (m *Master) GetMasterPlant() tk.M {
	return tk.M{}
}

// MasterUnitPlant
type MasterUnitPlant struct {
	orm.ModelBase `bson:"-" json:"-"`
	Id            bson.ObjectId ` bson:"_id" json:"_id" `
	Plant         string        `json:Plant`
	Unit          string        `json:Unit`
}

func NewMasterUnitPlant() *MasterUnitPlant {
	m := new(MasterUnitPlant)
	m.Id = bson.NewObjectId()
	return m
}

func (e *MasterUnitPlant) RecordID() interface{} {
	return e.Id
}

func (m *MasterUnitPlant) TableName() string {
	return "MasterUnitPlant"
}
