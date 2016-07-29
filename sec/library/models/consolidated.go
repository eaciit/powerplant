package models

import (
	"sync"
	"time"

	"github.com/eaciit/orm"
)

type Consolidated struct {
	sync.RWMutex
	orm.ModelBase    `bson:"-" json:"-"`
	Plant            string    `bson:"Plant" json:"Plant"`
	Unit             string    `bson:"Unit" json:"Unit"`
	ConsolidatedDate time.Time `bson:"ConsolidatedDate" json:"ConsolidatedDate"`
	EnergyNet        float64   `bson:"EnergyNet" json:"EnergyNet"`
	Capacity         float64   `bson:"Capacity" json:"Capacity"`
	FuelType         string    `bson:"FuelType" json:"FuelType"`
}

func (m *Consolidated) TableName() string {
	return "Consolidated"
}
