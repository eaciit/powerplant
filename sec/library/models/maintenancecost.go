package models

import (
	"sync"

	"github.com/eaciit/orm"
)

type MaintenanceCost struct {
	sync.RWMutex
	orm.ModelBase           `bson:"-" json:"-"`
	Id                      int64   `bson:"_id" json:"id"`
	FunctionalLocation      string  `bson:"FunctionalLocation" json:"FunctionalLocation"`
	Month                   int     `bson:"Month" json:"Month"`
	Year                    int     `bson:"Year" json:"Year"`
	MaintenanceOrder        string  `bson:"MaintenanceOrder" json:"MaintenanceOrder"`
	MaintenanceOrderDesc    string  `bson:"MaintenanceOrderDesc" json:"MaintenanceOrderDesc"`
	MaintenanceActivityType string  `bson:"MaintActivityType" json:"MaintenanceActivityType"`
	OrderType               string  `bson:"OrderType" json:"OrderType"`
	InternalLaborActual     float64 `bson:"InternalLaborActual" json:"InternalLaborActual"`
	InternalMaterialActual  float64 `bson:"InternalMaterialActual" json:"InternalMaterialActual"`
	DirectMaterialActual    float64 `bson:"DirectMaterialActual" json:"DirectMaterialActual"`
	ExternalServiceActual   float64 `bson:"ExternalServiceActual" json:"ExternalServiceActual"`
	PeriodTotalActual       float64 `bson:"PeriodTotalActual" json:"PeriodTotalActual"`
}

func (m *MaintenanceCost) TableName() string {
	return "MaintenanceCost"
}
