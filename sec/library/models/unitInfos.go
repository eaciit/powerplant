package models

import (
	"sync"

	"github.com/eaciit/orm"
)

type UnitInfos struct {
	sync.RWMutex
	orm.ModelBase            `bson:"-" json:"-"`
	Plant                    string  `bson:"Plant" json:"Plant"`
	Unit                     string  `bson:"Unit" json:"Unit"`
	Month                    int     `bson:"Month" json:"Month"`
	Year                     int     `bson:"Year" json:"Year"`
	SRFPercentage            float64 `bson:"SRFPercentage" json:"SRFPercentage"`
	UnitMaxPower             float64 `bson:"UnitMaxPower" json:"UnitMaxPower"`
	WAFPercentage            float64 `bson:"WAFPercentage" json:"WAFPercentage"`
	WUFPercentage            float64 `bson:"WUFPercentage" json:"WUFPercentage"`
	StartupPayment           float64 `bson:"StartupPayment" json:"StartupPayment"`
	Penalty                  float64 `bson:"Penalty" json:"Penalty"`
	PrimaryFuelType          string  `bson:"PrimaryFuelType" json:"PrimaryFuelType"`
	PrimaryFuelConsumedCost  float64 `bson:"PrimaryFuelConsumedCost" json:"PrimaryFuelConsumedCost"`
	Primary2FuelType         string  `bson:"Primary2FuelType" json:"Primary2FuelType"`
	Primary2FuelConsumedCost float64 `bson:"Primary2FuelConsumedCost" json:"Primary2FuelConsumedCost"`
	BackupFuelType           string  `bson:"BackupFuelType" json:"BackupFuelType"`
	BackupFuelConsumedCost   float64 `bson:"BackupFuelConsumedCost" json:"BackupFuelConsumedCost"`
}

func (m *UnitInfos) TableName() string {
	return "UnitInfos"
}
