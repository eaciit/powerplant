package models

import (
	"sync"
	"time"

	"github.com/eaciit/orm"
)

type MaintenanceCostFL struct {
	sync.RWMutex
	orm.ModelBase `bson:"-" json:"-"`
	// Id                     int64     `bson:"_id" json:"id"`
	OrderType              string    `bson:"OrderType" json:"OrderType"`
	OrderTypeDesc          string    `bson:"OrderTypeDesc" json:"OrderTypeDesc"`
	FunctionalLocation     string    `bson:"FunctionalLocation" json:"FunctionalLocation"`
	InternalLaborPlan      float64   `bson:"InternalLaborPLAN" json:"InternalLaborPlan"`
	Period                 time.Time `bson:"Period" json:"Period"`
	Plant                  string    `bson:"Plant" json:"Plant"`
	InternalLaborActual    float64   `bson:"InternalLaborActual" json:"InternalLaborActual"`
	InternalMaterialPlan   float64   `bson:"InternalMaterialPlan" json:"InternalMaterialPlan"`
	InternalMaterialActual float64   `bson:"InternalMaterialActual" json:"InternalMaterialActual"`
	DirectMaterialPlan     float64   `bson:"DirectMaterialPlan" json:"DirectMaterialPlan"`
	DirectMaterialActual   float64   `bson:"DirectMaterialActual" json:"DirectMaterialActual"`
	ExternalServicePlan    float64   `bson:"ExternalServicePlan" json:"ExternalServicePlan"`
	ExternalServiceActual  float64   `bson:"ExternalServiceActual" json:"ExternalServiceActual"`
	PeriodTotalPlan        float64   `bson:"PeriodTotalPlan" json:"PeriodTotalPlan"`
	PeriodTotalActual      float64   `bson:"PeriodTotalActual" json:"PeriodTotalActual"`
}

func (m *MaintenanceCostFL) TableName() string {
	return "MaintenanceCostFL"
}
