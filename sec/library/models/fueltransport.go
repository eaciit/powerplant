package models

import (
	"sync"

	"github.com/eaciit/orm"
)

type FuelTransport struct {
	sync.RWMutex
	orm.ModelBase `bson:"-" json:"-"`
	Plant         string  `bson:"Plant" json:"Plant"`
	Year          int     `bson:"Year" json:"Year"`
	TransportCost float64 `bson:"TransportCost" json:"TransportCost"`
}

func (m *FuelTransport) TableName() string {
	return "FuelTransport"
}
