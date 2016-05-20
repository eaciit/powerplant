package models

import (
	"sync"
	"time"

	"github.com/eaciit/orm"
)

type PrevMaintenanceValueEquation struct {
	sync.RWMutex
	orm.ModelBase `bson:"-",json:"-"`
	// _Id                 int64     `bson:"_id",json:"_id"`
	Plant               string    `bson:"Plant",json:"Plant"`
	Phase               string    `bson:"Phase",json:"Phase"`
	Block               string    `bson:"Block",json:"Block"`
	Unit                string    `bson:"Unit",json:"Unit"`
	Id                  string    `bson:"Id",json:"Id"`
	DatePerformed       time.Time `bson:"DatePerformed",json:"DatePerformed"`
	WOType              string    `bson:"WOType",json:"WOType"`
	UserStatus          string    `bson:"UserStatus",json:"UserStatus"`
	Description         string    `bson:"Description",json:"Description"`
	Days                int       `bson:"Days",json:"Days"`
	Materials           float64   `bson:"Materials",json:"Materials"`
	SkilledLabour       float64   `bson:"SkilledLabour",json:"SkilledLabour"`
	UnSkilledLabour     float64   `bson:"UnSkilledLabour",json:"UnSkilledLabour"`
	ExtraCost           float64   `bson:"ExtraCost",json:"ExtraCost"`
	ContractMaintenance float64   `bson:"ContractMaintenance",json:"ContractMaintenance"`
	TotalCost           float64   `bson:"TotalCost",json:"TotalCost"`
}

func (m *PrevMaintenanceValueEquation) TableName() string {
	return "PrevMaintenanceValueEquation"
}
