package models

import (
	"sync"
	"time"

	"github.com/eaciit/orm"
)

type Consolidated struct {
	sync.RWMutex
	orm.ModelBase `bson:"-" json:"-"`
	// Id            int64     `bson:"id" json:"id"`
	Plant                  string    `bson:"Plant" json:"Plant"`
	Unit                   string    `bson:"Unit" json:"Unit"`
	ConsolidatedDate       time.Time `bson:"ConsolidatedDate" json:"ConsolidatedDate"`
	Load0                  float64   `bson:"Load0" json:"Load0"`
	Load1                  float64   `bson:"Load1" json:"Load"`
	Load2                  float64   `bson:"Load2" json:"Load"`
	Load3                  float64   `bson:"Load3" json:"Load"`
	Load4                  float64   `bson:"Load4" json:"Load"`
	Load5                  float64   `bson:"Load5" json:"Load"`
	Load6                  float64   `bson:"Load6" json:"Load"`
	Load7                  float64   `bson:"Load7" json:"Load"`
	Load8                  float64   `bson:"Load8" json:"Load"`
	Load9                  float64   `bson:"Load9" json:"Load"`
	Load10                 float64   `bson:"Load10" json:"Load10"`
	Load11                 float64   `bson:"Load11" json:"Load11"`
	Load12                 float64   `bson:"Load12" json:"Load12"`
	Load13                 float64   `bson:"Load13" json:"Load13"`
	Load14                 float64   `bson:"Load14" json:"Load14"`
	Load15                 float64   `bson:"Load15" json:"Load15"`
	Load16                 float64   `bson:"Load16" json:"Load16"`
	Load17                 float64   `bson:"Load17" json:"Load17"`
	Load18                 float64   `bson:"Load18" json:"Load18"`
	Load19                 float64   `bson:"Load19" json:"Load19"`
	Load20                 float64   `bson:"Load20" json:"Load20"`
	Load21                 float64   `bson:"Load21" json:"Load21"`
	Load22                 float64   `bson:"Load22" json:"Load22"`
	Load23                 float64   `bson:"Load23" json:"Load23"`
	EnergyGross            float64   `bson:"EnergyGross" json:"EnergyGross"`
	EnergyNet              float64   `bson:"EnergyNet" json:"EnergyNet"`
	Capacity               float64   `bson:"Capacity" json:"Capacity"`
	FuelType               string    `bson:"FuelType" json:"FuelType"`
	FuelConsumption_Gas    float64   `bson:"FuelConsumption_Gas" json:"FuelConsumption_Gas"`
	FuelConsumption_Diesel float64   `bson:"FuelConsumption_Diesel" json:"FuelConsumption_Diesel"`
	FuelConsumption_Crude  float64   `bson:"FuelConsumption_Crude" json:"FuelConsumption_Crude"`
	TotalCapacity          float64   `bson:"TotalCapacity" json:"TotalCapacity"`
	CapacityPayment        float64   `bson:"CapacityPayment" json:"CapacityPayment"`
	EnergyPayment          float64   `bson:"EnergyPayment" json:"EnergyPayment"`
	StartupPayment         float64   `bson:"StartupPayment" json:"StartupPayment"`
	Penalty                float64   `bson:"Penalty" json:"Penalty"`
	Incentive              float64   `bson:"Incentive" json:"Incentive"`
}

func (m *Consolidated) TableName() string {
	return "Consolidated"
}
