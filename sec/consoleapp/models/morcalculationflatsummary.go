package models

import (
	"sync"

	"github.com/eaciit/orm"
)

type MORCalculationFlatSummary struct {
	sync.RWMutex
	orm.ModelBase `bson:"-" json:"-"`
	// Id int64 `bson:"Id" json:"Id"`
	OrderType                string  `bson:"OrderType" json:"OrderType"`
	EquipmentType            string  `bson:"EquipmentType" json:"EquipmentType"`
	EquipmentTypeDescription string  `bson:"EquipmentTypeDescription" json:"EquipmentTypeDescription"`
	ActivityType             string  `bson:"ActivityType" json:"ActivityType"`
	PeriodYear               int     `bson:"Period.Year" json:"PeriodYear"`
	Plant                    string  `bson:"Plant" json:"Plant"`
	Element                  string  `bson:"Element" json:"Element"`
	Value                    float64 `bson:"Value" json:"Value"`
	MOCount                  int     `bson:"MOCount" json:"MOCount"`
}

func (m *MORCalculationFlatSummary) TableName() string {
	return "MORCalculationFlatSummary"
}
