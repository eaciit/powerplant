package models

import (
	"sync"
	"time"

	"github.com/eaciit/orm"
)

type PreventiveMaintenance struct {
	sync.RWMutex
	orm.ModelBase          `bson:"-" json:"-"`
	Plant                  string    `bson:"Plant" json:"Plant"`
	Unit                   string    `bson:"Unit" json:"Unit"`
	FunctionalLocation     string    `bson:"FunctionalLocation" json:"FunctionalLocation"`
	DatePerformed          time.Time `bson:"DatePerformed" json:"DatePerformed"`
	WOType                 string    `bson:"WOType" json:"WOType"`
	Description            string    `bson:"Description" json:"Description"`
	Days                   int       `bson:"Days" json:"Days"`
	MaterialsSAR           float64   `bson:"MaterialsSAR" json:"MaterialsSAR"`
	SkilledLabourSAR       float64   `bson:"SkilledLabourSAR" json:"SkilledLabourSAR"`
	UnSkilledLabourSAR     float64   `bson:"UnSkilledLabourSAR" json:"UnSkilledLabourSAR"`
	ExtraCostSAR           float64   `bson:"ExtraCostSAR" json:"ExtraCostSAR"`
	ContractMaintenanceSAR float64   `bson:"ContractMaintenanceSAR" json:"ContractMaintenanceSAR"`
	TotalCostSAR           float64   `bson:"TotalCostSAR" json:"TotalCostSAR"`
}

func (m *PreventiveMaintenance) TableName() string {
	return "PreventiveMaintenance"
}
