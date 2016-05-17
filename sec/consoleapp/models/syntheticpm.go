package models

import (
	"github.com/eaciit/orm"
	"time"
)

type SyntheticPM struct {
	orm.ModelBase `bson:"-",json:"-"`
	// Id                 int64    `bson:"id",json:"id"`
	Plant              string    `bson:"Plant",json:"Plant"`
	Unit               string    `bson:"Unit",json:"Unit"`
	ScheduledStart     time.Time `bson:"ScheduledStart",json:"ScheduledStart"`
	WOID               string    `bson:"WOID",json:"WOID"`
	WOType             string    `bson:"WOType",json:"WOType"`
	Description        string    `bson:"Description",json:"Description"`
	PlannedLaborHours  int       `bson:"PlannedLaborHours",json:"PlannedLaborHours"`
	PlannedLaborCost   float64   `bson:"PlannedLaborCost",json:"PlannedLaborCost"`
	ActualMaterialCost float64   `bson:"ActualMaterialCost",json:"ActualMaterialCost"`
	Total              float64   `bson:"Total",json:"Total"`
}

func (m *SyntheticPM) TableName() string {
	return "SyntheticPM"
}
