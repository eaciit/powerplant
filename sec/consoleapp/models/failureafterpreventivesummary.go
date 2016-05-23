package models

import (
	"sync"
	"time"

	"github.com/eaciit/orm"
)

type FailureAfterPreventiveSummary struct {
	sync.RWMutex
	orm.ModelBase `bson:"-" json:"-"`
	// Id                      int64     `bson:"Id" json:"Id"`
	EquipmentType                        string    `bson:"EquipmentType" json:"EquipmentType"`
	Plant                                string    `bson:"Plant" json:"Plant"`
	OrderType                            string    `bson:"OrderType" json:"OrderType"`
	FunctionalLocation                   string    `bson:"FunctionalLocation" json:"FunctionalLocation"`
	FLDescription                        string    `bson:"FLDescription" json:"FLDescription"`
	NextPlannedMaintenanceDate           time.Time `bson:"NextPlannedMaintenanceDate" json:"NextPlannedMaintenanceDate"`
	CountOfMaintenanceAfter1Week         int       `bson:"CountOfMaintenanceAfter1Week" json:"CountOfMaintenanceAfter1Week"`
	CountOfMaintenanceAfter1Month        int       `bson:"CountOfMaintenanceAfter1Month" json:"CountOfMaintenanceAfter1Month"`
	CountOfMaintenanceAfter3Month        int       `bson:"CountOfMaintenanceAfter3Month" json:"CountOfMaintenanceAfter3Month"`
	CountOfMaintenanceAfterAnnual        int       `bson:"CountOfMaintenanceAfterAnnual" json:"CountOfMaintenanceAfterAnnual"`
	CountOfMaintenanceAfter6Month        int       `bson:"CountOfMaintenanceAfter6Month" json:"CountOfMaintenanceAfter6Month"`
	CountOfMaintenanceAfterAnnual6Month  int       `bson:"CountOfMaintenanceAfterAnnual6Month" json:"CountOfMaintenanceAfterAnnual6Month"`
	CountOfMaintenanceAfter9Month        int       `bson:"CountOfMaintenanceAfter9Month" json:"CountOfMaintenanceAfter9Month"`
	CountOfMaintenanceAfter12Month       int       `bson:"CountOfMaintenanceAfter12Month" json:"CountOfMaintenanceAfter12Month"`
	CountOfMaintenanceAfterAnnual12Month int       `bson:"CountOfMaintenanceAfterAnnual12Month" json:"CountOfMaintenanceAfterAnnual12Month"`
	CostOfMaintenanceAfter1Week          float64   `bson:"CostOfMaintenanceAfter1Week" json:"CostOfMaintenanceAfter1Week"`
	CostOfMaintenanceAfter1Month         float64   `bson:"CostOfMaintenanceAfter1Month" json:"CostOfMaintenanceAfter1Month"`
	CostOfMaintenanceAfter3Month         float64   `bson:"CostOfMaintenanceAfter3Month" json:"CostOfMaintenanceAfter3Month"`
	CostOfMaintenanceAfterAnnual         float64   `bson:"CostOfMaintenanceAfterAnnual" json:"CostOfMaintenanceAfterAnnual"`
	CostOfMaintenanceAfter6Month         float64   `bson:"CostOfMaintenanceAfter6Month" json:"CostOfMaintenanceAfter6Month"`
	CostOfMaintenanceAfterAnnual6Month   float64   `bson:"CostOfMaintenanceAfterAnnual6Month" json:"CostOfMaintenanceAfterAnnual6Month"`
	CostOfMaintenanceAfter9Month         float64   `bson:"CostOfMaintenanceAfter9Month" json:"CostOfMaintenanceAfter9Month"`
	CostOfMaintenanceAfter12Month        float64   `bson:"CostOfMaintenanceAfter12Month" json:"CostOfMaintenanceAfter12Month"`
	CostOfMaintenanceAfterAnnual12Month  float64   `bson:"CostOfMaintenanceAfterAnnual12Month" json:"CostOfMaintenanceAfterAnnual12Month"`
}

func (m *FailureAfterPreventiveSummary) TableName() string {
	return "FailureAfterPreventiveSummary"
}
