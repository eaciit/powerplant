package models

import (
	"github.com/eaciit/orm"
	"sync"
)

type PowerPlantInfo struct {
	sync.RWMutex
	orm.ModelBase         `bson:"-" json:"-"`
	Name                  string  `bson:"Name" json:"Name"`
	System                string  `bson:"System" json:"System"`
	Province              string  `bson:"Province" json:"Province"`
	Region                string  `bson:"Region" json:"Region"`
	City                  string  `bson:"City" json:"City"`
	FuelTypes_Crude       bool    `bson:"FuelTypes_Crude" json:"FuelTypes_Crude"`
	FuelTypes_Heavy       bool    `bson:"FuelTypes_Heavy" json:"FuelTypes_Heavy"`
	FuelTypes_Diesel      bool    `bson:"FuelTypes_Diesel" json:"FuelTypes_Diesel"`
	FuelTypes_Gas         bool    `bson:"FuelTypes_Gas" json:"FuelTypes_Gas"`
	GasTurbineUnit        int     `bson:"GasTurbineUnit" json:"GasTurbineUnit"`
	GasTurbineCapacity    float64 `bson:"GasTurbineCapacity" json:"GasTurbineCapacity"`
	SteamUnit             int     `bson:"SteamUnit" json:"SteamUnit"`
	SteamCapacity         float64 `bson:"SteamCapacity" json:"SteamCapacity"`
	DieselUnit            int     `bson:"DieselUnit" json:"DieselUnit"`
	DieselCapacity        float64 `bson:"DieselCapacity" json:"DieselCapacity"`
	CombinedCycleUnit     int     `bson:"CombinedCycleUnit" json:"CombinedCycleUnit"`
	CombinedCycleCapacity float64 `bson:"CombinedCycleCapacity" json:"CombinedCycleCapacity"`
}

func (m *PowerPlantInfo) TableName() string {
	return "PowerPlantInfo"
}
