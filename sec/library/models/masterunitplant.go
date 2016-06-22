package models

import (
	"sync"

	"github.com/eaciit/orm"
	tk "github.com/eaciit/toolkit"
)

type Master struct {
}

func (m *Master) GetMasterPlant() tk.M {
	return tk.M{}
}

// MasterUnitPlant
type MasterUnitPlant struct {
	sync.RWMutex
	orm.ModelBase `bson:"-",json:"-"`
	Plant         string `bson:"Plant",json:"Plant"`
	Unit          string `bson:"Unit",json:"Unit"`
}

func NewMasterUnitPlant() *MasterUnitPlant {
	m := new(MasterUnitPlant)
	return m
}

func (m *MasterUnitPlant) TableName() string {
	return "MasterUnitPlant"
}
