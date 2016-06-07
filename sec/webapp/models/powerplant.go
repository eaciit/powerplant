package models

import (
	"github.com/eaciit/orm"
)

type PlantData struct {
	orm.ModelBase         `bson:"-",json:"-"`
	PlantName             string  `bson:"PlantName"`
	PlantCode             string  `bson:"PlantCode"`
	Province              string  `bson:"Province"`
	Region                string  `bson:"Region"`
	City                  string  `bson:"City"`
	Longitude             float64 `bson:"Longitude"`
	Latitude              float64 `bson:"Latitude"`
	SteamCapacity         float64 `bson:"SteamCapacity"`
	DieselCapacity        float64 `bson:"DieselCapacity"`
	GasTurbineCapacity    float64 `bson:"GasTurbineCapacity"`
	CombinedCycleCapacity float64 `bson:"CombinedCycleCapacity"`
	TotalCapacity         float64 `bson:"TotalCapacity"`
}

func (c *PlantData) TableName() string {
	return "PowerPlantCoordinates"
}
