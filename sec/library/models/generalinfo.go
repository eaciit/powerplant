package models

import (
	"sync"
	"time"

	"github.com/eaciit/orm"
)

type GeneralInfo struct {
	sync.RWMutex
	orm.ModelBase          `bson:"-" json:"-"`
	Id                     string                             `bson:"_id" json:"Id"`
	Company                string                             `bson:"COMPANY" json:"Company"`
	Area                   string                             `bson:"AREA" json:"Area"`
	Plant                  string                             `bson:"PLANT" json:"Plant"`
	Unit                   string                             `bson:"UNIT" json:"Unit"`
	Manufacturer           string                             `bson:"MANUFACTURER" json:"Manufacturer"`
	Model                  string                             `bson:"MODEL" json:"Model"`
	UnitType               string                             `bson:"UnitTYPE" json:"UnitType"`
	InstalledCapacity      float64                            `bson:"InstalledCapacity" json:"InstalledCapacity"`
	OperationalCapacity    float64                            `bson:"OperationalCapacity" json:"OperationalCapacity"`
	PrimaryFuel1           string                             `bson:"PRIMARYFuel1" json:"PrimaryFuel1"`
	PrimaryFuel2Startup    string                             `bson:"PRIMARYFuel2StartUp" json:"PrimaryFuel2Startup"`
	BackupFuel             string                             `bson:"BACKUPFuel" json:"BackupFuel"`
	DutyCycle              string                             `bson:"DutyCycle" json:"DutyCycle"`
	HeatRate               float64                            `bson:"HeatRate" json:"HeatRate"`
	Efficiency             float64                            `bson:"Efficiency" json:"Efficiency"`
	CommissioningDate      float64                            `bson:"CommissioningDate" json:"CommissioningDate"`
	RetirementPlan         string                             `bson:"RetirementPlan" json:"RetirementPlan"`
	UpdDate                time.Time                          `bson:"Update" json:"UpdDate"`
	InstalledMWH           []GeneralInfoDetails               `bson:"InstalledMWH" json:"InstalledMWH"`
	ActualEnergyGeneration []GeneralInfoDetails               `bson:"ActualEnergyGeneration" json:"ActualEnergyGeneration"`
	ActualFuelConsumption  []GeneralInfoActualFuelConsumption `bson:"GeneralInfoActualFuelConsumption" json:"GeneralInfoActualFuelConsumption"`
	CapacityFactor         []GeneralInfoDetails               `bson:"CapacityFactor" json:"CapacityFactor"`
}

func (m *GeneralInfo) TableName() string {
	return "GeneralInfo"
}

type GeneralInfoDetails struct {
	sync.RWMutex
	orm.ModelBase `bson:"-" json:"-"`
	GenID         string  `bson:"GenID" json:"GenID"`
	Type          string  `bson:"Type" json:"Type"`
	Year          int     `bson:"Year" json:"Year"`
	Value         float64 `bson:"Value" json:"Value"`
}

func (m *GeneralInfoDetails) TableName() string {
	return "GeneralInfoDetails"
}

type GeneralInfoActualFuelConsumption struct {
	sync.RWMutex
	orm.ModelBase `bson:"-" json:"-"`
	// Id            int64   `bson:"Id" json:"Id"`
	GenID        string  `bson:"GenID" json:"GenID"`
	Year         int     `bson:"Year" json:"Year"`
	GASMMSCF     float64 `bson:"GASMMSCF" json:"GASMMSCF"`
	CrudeBarrel  float64 `bson:"CrudeBarrel" json:"CrudeBarrel"`
	HFOBarrel    float64 `bson:"HFOBarrel" json:"HFOBarrel"`
	DieselBarrel float64 `bson:"DieselBarrel" json:"DieselBarrel"`
}

func (m *GeneralInfoActualFuelConsumption) TableName() string {
	return "GeneralInfoActualFuelConsumption"
}
