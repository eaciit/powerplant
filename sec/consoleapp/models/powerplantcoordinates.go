package models

import (
	"github.com/eaciit/orm"
)

type PowerPlantCoordinates struct {
	orm.ModelBase         `bson:"-",json:"-"`
	PlantCode             string  `bson:"PlantCode",json:"PlantCode"`
	PlantName             string  `bson:"PlantName",json:"PlantName"`
	PlantType             string  `bson:"PlantType",json:"PlantType"`
	Province              string  `bson:"Province",json:"Province"`
	Region                string  `bson:"Region",json:"Region"`
	City                  string  `bson:"City",json:"City"`
	FuelTypes_Crude       bool    `bson:"FuelTypes_Crude",json:"FuelTypes_Crude"`
	FuleTypes_Heavy       bool    `bson:"FuleTypes_Heavy",json:"FuleTypes_Heavy"`
	FuleTypes_Diesel      bool    `bson:"FuleTypes_Diesel",json:"FuleTypes_Diesel"`
	FuleTypes_Gas         bool    `bson:"FuleTypes_Gas",json:"FuleTypes_Gas"`
	GasTurbineUnit        float64 `bson:"GasTurbineUnit",json:"GasTurbineUnit"`
	GasTurbineCapacity    float64 `bson:"GasTurbineCapacity",json:"GasTurbineCapacity"`
	SteamUnit             int     `bson:"SteamUnit",json:"SteamUnit"`
	SteamCapacity         float64 `bson:"SteamCapacity",json:"SteamCapacity"`
	DieselUnit            int     `bson:"DieselUnit",json:"DieselUnit"`
	DieselCapacity        float64 `bson:"DieselCapacity",json:"DieselCapacity"`
	CombinedCycleUnit     int     `bson:"CombinedCycleUnit",json:"CombinedCycleUnit"`
	CombinedCycleCapacity float64 `bson:"CombinedCycleCapacity",json:"CombinedCycleCapacity"`
	Longitude             float64 `bson:"Longitude",json:"Longitude"`
	Latitude              float64 `bson:"Latitude",json:"Latitude"`
}

func (m *PowerPlantCoordinates) TableName() string {
	return "PowerPlantCoordinates"
}
