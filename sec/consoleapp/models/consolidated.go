package models

import (
	"time"

	"github.com/eaciit/orm"
)

type Consolidated struct {
	orm.ModelBase `bson:"-",json:"-"`
	// Id            int64     `bson:"id",json:"id"`
	Plant                  string    `bson:"Plant",json:"Plant"`
	Unit                   string    `bson:"Unit",json:"Unit"`
	ConsolidatedDate       time.Time `bson:"ConsolidatedDate",json:"ConsolidatedDate"`
	Load0                  float     `bson:"Load0",json:"Load0"`
	Load1                  float     `bson:"Load1",json:"Load"`
	Load2                  float     `bson:"Load2",json:"Load"`
	Load3                  float     `bson:"Load3",json:"Load"`
	Load4                  float     `bson:"Load4",json:"Load"`
	Load5                  float     `bson:"Load5",json:"Load"`
	Load6                  float     `bson:"Load6",json:"Load"`
	Load7                  float     `bson:"Load7",json:"Load"`
	Load8                  float     `bson:"Load8",json:"Load"`
	Load9                  float     `bson:"Load9",json:"Load"`
	Load10                 float     `bson:"Load10",json:"Load10"`
	Load11                 float     `bson:"Load11",json:"Load11"`
	Load12                 float     `bson:"Load12",json:"Load12"`
	Load13                 float     `bson:"Load13",json:"Load13"`
	Load14                 float     `bson:"Load14",json:"Load14"`
	Load15                 float     `bson:"Load15",json:"Load15"`
	Load16                 float     `bson:"Load16",json:"Load16"`
	Load17                 float     `bson:"Load17",json:"Load17"`
	Load18                 float     `bson:"Load18",json:"Load18"`
	Load19                 float     `bson:"Load19",json:"Load19"`
	Load20                 float     `bson:"Load20",json:"Load20"`
	Load21                 float     `bson:"Load21",json:"Load21"`
	Load22                 float     `bson:"Load22",json:"Load22"`
	Load23                 float     `bson:"Load23",json:"Load23"`
	EnergyGross            float     `bson:"EnergyGross",json:"EnergyGross"`
	EnergyNet              float     `bson:"EnergyNet",json:"EnergyNet"`
	Capacity               float     `bson:"Capacity",json:"Capacity"`
	FuelType               string    `bson:"FuelType",json:"FuelType"`
	FuelConsumption_Gas    float     `bson:"FuelConsumption_Gas",json:"FuelConsumption_Gas"`
	FuelConsumption_Diesel float     `bson:"FuelConsumption_Diesel",json:"FuelConsumption_Diesel"`
	FuelConsumption_Crude  float     `bson:"FuelConsumption_Crude",json:"FuelConsumption_Crude"`
	TotalCapacity          float     `bson:"TotalCapacity",json:"TotalCapacity"`
	CapacityPayment        float     `bson:"CapacityPayment",json:"CapacityPayment"`
	EnergyPayment          float     `bson:"EnergyPayment",json:"EnergyPayment"`
	StartupPayment         float     `bson:"StartupPayment",json:"StartupPayment"`
	Penalty                float     `bson:"Penalty",json:"Penalty"`
	Incentive              float     `bson:"Incentive",json:"Incentive"`
}

func (m *Consolidated) TableName() string {
	return "Consolidated"
}
