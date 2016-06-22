package models

import (
	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
)

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

type PhaseModel struct {
	orm.ModelBase `bson:"-",json:"-"`
	Id            int    ` bson:"Id" json:"Id" `
	Phase         string ` bson:"Phase" json:"Phase" `
}

func NewPhaseModel() *PhaseModel {
	m := new(PhaseModel)
	// m.Id = bson.NewObjectId()
	return m
}

func (e *PhaseModel) RecordID() interface{} {
	return e.Id
}

func (m *PhaseModel) TableName() string {
	return "MasterPhase"
}
