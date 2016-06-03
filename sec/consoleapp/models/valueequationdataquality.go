package models

import (
	"github.com/eaciit/orm"
	"sync"
	"time"
)

type ValueEquationDataQuality struct {
	sync.RWMutex
	orm.ModelBase `bson:"-" json:"-"`
	// Id            int64     `bson:"Id" json:"Id"`
	Year        int       `bson:"Period.Year" json:"Year"`
	Month       int       `bson:"Period.Month" json:"Month"`
	MonthYear   int       `bson:"Period.MonthYear" json:"MonthYear"`
	Quarter     int       `bson:"Period.Quarter" json:"Quarter"`
	QuarterYear int       `bson:"Period.QuarterYear" json:"QuarterYear"`
	Dates       time.Time `bson:"Period.Dates" json:"Dates"`

	Plant                    string  `bson:"Plant" json:"Plant"`
	Unit                     string  `bson:"Unit" json:"Unit"`
	Appendix_Data            float64 `bson:"Appendix_Data" json:"Appendix_Data"`
	Consolidated_Data        float64 `bson:"Consolidated_Data" json:"Consolidated_Data"`
	Synthetic_Data           float64 `bson:"Synthetic_Data" json:"Synthetic_Data"`
	PerformanceFactor_Data   float64 `bson:"PerformanceFactor_Data" json:"PerformanceFactor_Data"`
	FuelTransport_Data       float64 `bson:"FuelTransport_Data" json:"FuelTransport_Data"`
	Outages_Data             float64 `bson:"Outages_Data" json:"Outages_Data"`
	CapacityPayment_Data     float64 `bson:"CapacityPayment_Data" json:"CapacityPayment_Data"`
	EnergyPayment_Data       float64 `bson:"EnergyPayment_Data" json:"EnergyPayment_Data"`
	StartupPayment_Data      float64 `bson:"StartupPayment_Data" json:"StartupPayment_Data"`
	Penalty_Data             float64 `bson:"Penalty_Data" json:"Penalty_Data"`
	Incentive_Data           float64 `bson:"Incentive_Data" json:"Incentive_Data"`
	PrimaryFuel1st_Data      float64 `bson:"PrimaryFuel1st_Data" json:"PrimaryFuel1st_Data"`
	PrimaryFuel2nd_Data      float64 `bson:"PrimaryFuel2nd_Data" json:"PrimaryFuel2nd_Data"`
	BackupFuel_Data          float64 `bson:"BackupFuel_Data" json:"BackupFuel_Data"`
	MaintenanceCost_Data     float64 `bson:"MaintenanceCost_Data" json:"MaintenanceCost_Data"`
	MaintenanceDuration_Data float64 `bson:"MaintenanceDuration_Data" json:"MaintenanceDuration_Data"`
}

func (m *ValueEquationDataQuality) TableName() string {
	return "ValueEquationDataQuality"
}

// CapacityPaymentDocuments
type VEDQCapacityPaymentDocuments struct {
	sync.RWMutex
	orm.ModelBase `bson:"-" json:"-"`
	// Id            int64  `bson:"Id" json:"Id"`
	VEId         int64  `bson:"VEId" json:"VEId"`
	DocName      string `bson:"DocName" json:"DocName"`
	Availability bool   `bson:"Availability" json:"Availability"`
}

func (m *VEDQCapacityPaymentDocuments) TableName() string {
	return "VEDQCapacityPaymentDocuments"
}

// EnergyPaymentDocuments
type VEDQEnergyPaymentDocuments struct {
	sync.RWMutex
	orm.ModelBase `bson:"-" json:"-"`
	// Id            int64  `bson:"Id" json:"Id"`
	VEId         int64  `bson:"VEId" json:"VEId"`
	DocName      string `bson:"DocName" json:"DocName"`
	Availability bool   `bson:"Availability" json:"Availability"`
}

func (m *VEDQEnergyPaymentDocuments) TableName() string {
	return "VEDQEnergyPaymentDocuments"
}

// StartupPaymentDocuments
type VEDQStartupPaymentDocuments struct {
	sync.RWMutex
	orm.ModelBase `bson:"-" json:"-"`
	// Id            int64  `bson:"Id" json:"Id"`
	VEId         int64  `bson:"VEId" json:"VEId"`
	DocName      string `bson:"DocName" json:"DocName"`
	Availability bool   `bson:"Availability" json:"Availability"`
}

func (m *VEDQStartupPaymentDocuments) TableName() string {
	return "VEDQStartupPaymentDocuments"
}

// PenaltyDocuments
type VEDQPenaltyDocuments struct {
	sync.RWMutex
	orm.ModelBase `bson:"-" json:"-"`
	// Id            int64  `bson:"Id" json:"Id"`
	VEId         int64  `bson:"VEId" json:"VEId"`
	DocName      string `bson:"DocName" json:"DocName"`
	Availability bool   `bson:"Availability" json:"Availability"`
}

func (m *VEDQPenaltyDocuments) TableName() string {
	return "VEDQPenaltyDocuments"
}

// PrimaryFuel1stDocuments
type VEDQPrimaryFuel1stDocuments struct {
	sync.RWMutex
	orm.ModelBase `bson:"-" json:"-"`
	// Id            int64  `bson:"Id" json:"Id"`
	VEId         int64  `bson:"VEId" json:"VEId"`
	DocName      string `bson:"DocName" json:"DocName"`
	Availability bool   `bson:"Availability" json:"Availability"`
}

func (m *VEDQPrimaryFuel1stDocuments) TableName() string {
	return "VEDQPrimaryFuel1stDocuments"
}

// PrimaryFuel2ndDocuments
type VEDQPrimaryFuel2ndDocuments struct {
	sync.RWMutex
	orm.ModelBase `bson:"-" json:"-"`
	// Id            int64  `bson:"Id" json:"Id"`
	VEId         int64  `bson:"VEId" json:"VEId"`
	DocName      string `bson:"DocName" json:"DocName"`
	Availability bool   `bson:"Availability" json:"Availability"`
}

func (m *VEDQPrimaryFuel2ndDocuments) TableName() string {
	return "VEDQPrimaryFuel2ndDocuments"
}

// BackupFuelDocuments
type VEDQBackupFuelDocuments struct {
	sync.RWMutex
	orm.ModelBase `bson:"-" json:"-"`
	// Id            int64  `bson:"Id" json:"Id"`
	VEId         int64  `bson:"VEId" json:"VEId"`
	DocName      string `bson:"DocName" json:"DocName"`
	Availability bool   `bson:"Availability" json:"Availability"`
}

func (m *VEDQBackupFuelDocuments) TableName() string {
	return "VEDQBackupFuelDocuments"
}

// MaintenanceCostDocuments
type VEDQMaintenanceCostDocuments struct {
	sync.RWMutex
	orm.ModelBase `bson:"-" json:"-"`
	// Id            int64  `bson:"Id" json:"Id"`
	VEId         int64  `bson:"VEId" json:"VEId"`
	DocName      string `bson:"DocName" json:"DocName"`
	Availability bool   `bson:"Availability" json:"Availability"`
}

func (m *VEDQMaintenanceCostDocuments) TableName() string {
	return "VEDQMaintenanceCostDocuments"
}

// MaintenanceDurationDocuments
type VEDQMaintenanceDurationDocuments struct {
	sync.RWMutex
	orm.ModelBase `bson:"-" json:"-"`
	// Id            int64  `bson:"Id" json:"Id"`
	VEId         int64  `bson:"VEId" json:"VEId"`
	DocName      string `bson:"DocName" json:"DocName"`
	Availability bool   `bson:"Availability" json:"Availability"`
}

func (m *VEDQMaintenanceDurationDocuments) TableName() string {
	return "VEDQMaintenanceDurationDocuments"
}
