package models

import (
	"sync"

	"github.com/eaciit/orm"
)

type WODurationSummary struct {
	sync.RWMutex
	orm.ModelBase `bson:"-" json:"-"`
	// Id                       int64   `bson:"Id" json:"Id"`
	OrderType                string  `bson:"OrderType" json:"OrderType"`
	EquipmentType            string  `bson:"EquipmentType" json:"EquipmentType"`
	EquipmentTypeDescription string  `bson:"EquipmentTypeDescription" json:"EquipmentTypeDescription"`
	ActivityType             string  `bson:"ActivityType" json:"ActivityType"`
	PeriodYear               int     `bson:"Period.Year" json:"PeriodYear"`
	Plant                    string  `bson:"Plant" json:"Plant"`
	ActualValue              float64 `bson:"ActualValue" json:"ActualValue"`
	PlanValue                float64 `bson:"PlanValue" json:"PlanValue"`
	MaxActualValue           float64 `bson:"MaxActualValue" json:"MaxActualValue"`
	MinActualValue           float64 `bson:"MinActualValue" json:"MinActualValue"`
	MaxPlanValue             float64 `bson:"MaxPlanValue" json:"MaxPlanValue"`
	MinPlanValue             float64 `bson:"MinPlanValue" json:"MinPlanValue"`
	Cost                     float64 `bson:"Cost" json:"Cost"`
	WOCount                  int     `bson:"WOCount" json:"WOCount"`
}

func (m *WODurationSummary) TableName() string {
	return "WODurationSummaryTest"
}
