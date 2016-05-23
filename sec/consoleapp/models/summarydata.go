package models

import (
	"sync"

	"github.com/eaciit/orm"
)

type SummaryData struct {
	sync.RWMutex
	orm.ModelBase `bson:"-" json:"-"`
	// Id                        int64   `bson:"Id" json:"Id"`
	FunctionalLocation        string  `bson:"FunctionalLocation" json:"FunctionalLocation"`
	FLDescription             string  `bson:"FLDescription" json:"FLDescription"`
	SortField                 string  `bson:"SortField" json:"SortField"`
	ParentFL                  string  `bson:"ParentFL" json:"ParentFL"`
	HasChild                  bool    `bson:"HasChild" json:"HasChild"`
	Province                  string  `bson:"Province" json:"Province"`
	Region                    string  `bson:"Region" json:"Region"`
	City                      string  `bson:"City" json:"City"`
	GasTurbineUnit            int     `bson:"GasTurbineUnit" json:"GasTurbineUnit"`
	GasTurbineCapacity        float64 `bson:"GasTurbineCapacity" json:"GasTurbineCapacity"`
	SteamUnit                 int     `bson:"SteamUnit" json:"SteamUnit"`
	SteamUnitCapacity         float64 `bson:"SteamUnitCapacity" json:"SteamUnitCapacity"`
	DieselUnit                int     `bson:"DieselUnit" json:"DieselUnit"`
	DieselUnitCapacity        float64 `bson:"DieselUnitCapacity" json:"DieselUnitCapacity"`
	CombinedCycleUnit         int     `bson:"CombinedCycleUnit" json:"CombinedCycleUnit"`
	CombinedCycleUnitCapacity float64 `bson:"CombinedCycleUnitCapacity" json:"CombinedCycleUnitCapacity"`
	MaintenancePlan           []MaintenancePlan
	MaintenanceWorkOrder      []MaintenanceWorkOrder
	RPPCloseWO                []RPPCloseWO
	OperationalData           []OperationalData
	Availability              []Availability
}

func (m *SummaryData) TableName() string {
	return "SummaryData"
}
