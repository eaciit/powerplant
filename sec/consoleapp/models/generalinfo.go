package models

import (
	"sync"
	"time"

	"github.com/eaciit/orm"
)

type GeneralInfo struct {
	sync.RWMutex
	orm.ModelBase       `bson:"-" json:"-"`
	Id                  string               `bson:"_id" json:"Id"`
	Company             string               `bson:"Company" json:"Company"`
	Area                string               `bson:"Area" json:"Area"`
	Plant               string               `bson:"Plant" json:"Plant"`
	Unit                string               `bson:"Unit" json:"Unit"`
	Manufacturer        string               `bson:"Manufacturer" json:"Manufacturer"`
	Model               string               `bson:"Model" json:"Model"`
	UnitType            string               `bson:"UnitType" json:"UnitType"`
	InstalledCapacity   float64              `bson:"InstalledCapacity" json:"InstalledCapacity"`
	OperationalCapacity float64              `bson:"OperationalCapacity" json:"OperationalCapacity"`
	PrimaryFuel1        string               `bson:"PrimaryFuel1" json:"PrimaryFuel1"`
	PrimaryFuel2Startup string               `bson:"PrimaryFuel2Startup" json:"PrimaryFuel2Startup"`
	BackupFuel          string               `bson:"BackupFuel" json:"BackupFuel"`
	DutyCycle           string               `bson:"DutyCycle" json:"DutyCycle"`
	HeatRate            float64              `bson:"HeatRate" json:"HeatRate"`
	Efficiency          float64              `bson:"Efficiency" json:"Efficiency"`
	CommisioningDate    float64              `bson:"CommisioningDate" json:"CommisioningDate"`
	RetirementPlant     string               `bson:"RetirementPlant" json:"RetirementPlant"`
	Update              time.Time            `bson:"Update" json:"Update"`
	Details             []GeneralInfoDetails `bson:"Details" json:"Details"`
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
