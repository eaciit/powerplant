package models

import (
	"github.com/eaciit/orm"
)

type StartupPaymentAndPenalty struct {
	sync.RWMutex
	orm.ModelBase `bson:"-",json:"-"`
	// id             int64   `bson:"id",json:"id"`
	Plant          string  `bson:"Plant",json:"Plant"`
	Year           int     `bson:"Year",json:"Year"`
	Unit           string  `bson:"Unit",json:"Unit"`
	StartupPayment float64 `bson:"StartupPayment",json:"StartupPayment"`
	Penalty        float64 `bson:"Penalty",json:"Penalty"`
}

func (m *StartupPaymentAndPenalty) TableName() string {
	return "StartupPaymentAndPenalty"
}
