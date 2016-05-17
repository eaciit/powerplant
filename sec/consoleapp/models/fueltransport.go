package models

import "github.com/eaciit/orm"


type FuelTransport struct {
	orm.ModelBase `bson:"-",json:"-"`
	Plant         string `bson:"Plant",json:"Plant"`
	Year          int    `bson:"Year",json:"Year"`
	TransportCost float  `bson:"TransportCost",json:"TransportCost"`
}

func (m *FuelTransport) TableName() string {
	return "FuelTransport"
}
