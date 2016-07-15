package models

import (
	"github.com/eaciit/orm"
	"sync"
)

// MasterUnitPlant
type MasterUnitPlant struct {
	sync.RWMutex
	orm.ModelBase `bson:"-",json:"-"`
	Plant         string `bson:"Plant" json:"Plant"`
	Unit          string `bson:"Unit" json:"Unit"`
}

func NewMasterUnitPlant() *MasterUnitPlant {
	m := new(MasterUnitPlant)
	return m
}

func (m *MasterUnitPlant) TableName() string {
	return "MasterUnitPlant"
}

// MasterActivityType
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

// MasterPhase
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

// MasterUnit
type MasterUnit struct {
	sync.RWMutex
	orm.ModelBase `bson:"-" json:"-"`
	Unit          string `bson:"Unit" json:"Unit"`
}

func (m *MasterUnit) TableName() string {
	return "MasterUnit"
}
