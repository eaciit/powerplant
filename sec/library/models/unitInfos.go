package models

import (
	"sync"

	"github.com/eaciit/orm"
)

type UnitInfos struct {
	sync.RWMutex
	orm.ModelBase        `bson:"-" json:"-"`
	Plant                string  `bson:"Plant" json:"Plant"`
	Unit                 string  `bson:"Unit" json:"Unit"`
	Month                int     `bson:"Month" json:"Month"`
	Year                 int     `bson:"Year" json:"Year"`
	SRFPercentage        float64 `bson:"SRFPercentage" json:"SRFPercentage"`
	UnitMaxPower         float64 `bson:"UnitMaxPower" json:"UnitMaxPower"`
	WAFPercentage        float64 `bson:"WAFPercentage" json:"WAFPercentage"`
	WUFPercentage        float64 `bson:"WUFPercentage" json:"WUFPercentage"`
	StartupPayment       float64 `bson:"StartupPayment" json:"StartupPayment"`
	Penalty              float64 `bson:"Penalty" json:"Penalty"`
	PrimaryFuelType      string  `bson:"PrimaryFuelType" json:"PrimaryFuelType"`
	PrimaryFuelConsumed  float64 `bson:"PrimaryFuelConsumed" json:"PrimaryFuelConsumed"`
	Primary2FuelType     string  `bson:"Primary2FuelType" json:"Primary2FuelType"`
	Primary2FuelConsumed float64 `bson:"Primary2FuelConsumed" json:"Primary2FuelConsumed"`
	BackupFuelType       string  `bson:"BackupFuelType" json:"BackupFuelType"`
	BackupFuelConsumed   float64 `bson:"BackupFuelConsumed" json:"BackupFuelConsumed"`
}

func (m *UnitInfos) TableName() string {
	return "UnitInfos"
}
