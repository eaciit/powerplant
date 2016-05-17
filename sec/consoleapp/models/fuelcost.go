package models

import (
	"sync"

	"github.com/eaciit/orm"
)

type FuelCost struct {
	sync.RWMutex
	orm.ModelBase `bson:"-",json:"-"`
	// Id            int64     `bson:"id",json:"id"`
	Plant                 string  `bson:"Plant",json:"Plant"`
	UnitId                string  `bson:"UnitId",json:"UnitId"`
	Year                  int     `bson:"Year",json:"Year"`
	Month                 int     `bson:"Month",json:"Month"`
	DependableCapacity    float64 `bson:"DependableCapacity",json:"DependableCapacity"`
	EnergyGrossProduction float64 `bson:"EnergyGrossProduction",json:"EnergyGrossProduction"`
	EnergyNetProduction   float64 `bson:"EnergyNetProduction",json:"EnergyNetProduction"`
	PrimaryFuelType       string  `bson:"PrimaryFuelType",json:"PrimaryFuelType"`
	PrimaryFuelGCV        float64 `bson:"PrimaryFuelGCV",json:"PrimaryFuelGCV"`
	PrimaryFuelConsumed   float64 `bson:"PrimaryFuelConsumed",json:"PrimaryFuelConsumed"`
	Primary2FuelType      string  `bson:"Primary2FuelType",json:"Primary2FuelType"`
	Primary2FuelGCV       float64 `bson:"Primary2FuelGCV",json:"Primary2FuelGCV"`
	Primary2FuelConsumed  float64 `bson:"Primary2FuelConsumed",json:"Primary2FuelConsumed"`
	BackupFuelType        string  `bson:"BackupFuelType",json:"BackupFuelType"`
	BackupFuelGCV         float64 `bson:"BackupFuelGCV",json:"BackupFuelGCV"`
	BackupFuelConsumed    float64 `bson:"BackupFuelConsumed",json:"BackupFuelConsumed"`
	AverageHeatRate       float64 `bson:"AverageHeatRate",json:"AverageHeatRate"`
}

func (m *FuelCost) TableName() string {
	return "FuelCost"
}
