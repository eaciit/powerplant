package models

import (
	"sync"

	"github.com/eaciit/orm"
)

type PreventiveCorrectiveSummary struct {
	sync.RWMutex
	orm.ModelBase `bson:"-" json:"-"`
	// Id                       int64   `bson:"Id" json:"Id"`
	OrderType                string  `bson:"OrderType" json:"OrderType"`
	EquipmentNo              string  `bson:"EquipmentNo" json:"EquipmentNo"`
	EquipmentDescription     string  `bson:"EquipmentDescription" json:"EquipmentDescription"`
	EquipmentType            string  `bson:"EquipmentType" json:"EquipmentType"`
	EquipmentTypeDescription string  `bson:"EquipmentTypeDescription" json:"EquipmentTypeDescription"`
	ActivityType             string  `bson:"ActivityType" json:"ActivityType"`
	PeriodYear               int     `bson:"PeriodYear" json:"PeriodYear"`
	Plant                    string  `bson:"Plant" json:"Plant"`
	Element                  string  `bson:"Element" json:"Element"`
	Value                    float64 `bson:"Value" json:"Value"`
	MOCount                  int     `bson:"MOCount" json:"MOCount"`
}

func (m *PreventiveCorrectiveSummary) TableName() string {
	return "PreventiveCorrectiveSummary"
}
