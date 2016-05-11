package models

import (
	. "github.com/eaciit/orm"
)

type PlantData struct {
	ModelBase             `bson:"base"`
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

type PlantCapacity struct {
	ModelBase     `bson:"base"`
	PlantCode     string  `bson:"_id"`
	TotalCapacity float64 `bson:"TotalCapacity"`
}

func (c *PlantCapacity) TableName() string {
	return "DataBrowser"
}
