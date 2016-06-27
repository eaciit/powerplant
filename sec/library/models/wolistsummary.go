package models

import (
	"sync"
	"time"

	"github.com/eaciit/orm"
)

type WOListSummary struct {
	sync.RWMutex
	orm.ModelBase `bson:"-" json:"-"`
	// Id                      int64     `bson:"Id" json:"Id"`
	PeriodYear              int       `bson:"Period.Year" json:"PeriodYear"`
	OrderType               string    `bson:"OrderType" json:"OrderType"`
	MainenanceOrderCode     string    `bson:"MainenanceOrderCode" json:"MainenanceOrderCode"`
	NotificationCode        string    `bson:"NotificationCode" json:"NotificationCode"`
	EquipmentType           string    `bson:"EquipmentType" json:"EquipmentType"`
	Plant                   string    `bson:"Plant" json:"Plant"`
	PlantType               string    `bson:"PlantType" json:"PlantType"`
	FunctionalLocation      string    `bson:"FunctionalLocation" json:"FunctionalLocation"`
	PlanStart               time.Time `bson:"PlanStart" json:"PlanStart"`
	PlanEnd                 time.Time `bson:"PlanEnd" json:"PlanEnd"`
	PlanDuration            float64   `bson:"PlanDuration" json:"PlanDuration"`
	ActualStart             time.Time `bson:"ActualStart" json:"ActualStart"`
	ActualEnd               time.Time `bson:"ActualEnd" json:"ActualEnd"`
	ActualDuration          float64   `bson:"ActualDuration" json:"ActualDuration"`
	LastMaintenanceEnd      time.Time `bson:"LastMaintenanceEnd" json:"LastMaintenanceEnd"`
	Cost                    float64   `bson:"Cost" json:"Cost"`
	LastMaintenanceInterval float64   `bson:"LastMaintenanceInterval" json:"LastMaintenanceInterval"`
}

func (m *WOListSummary) TableName() string {
	return "WOListSummaryTest"
}
