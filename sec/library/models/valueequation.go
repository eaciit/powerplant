package models

import (
	"github.com/eaciit/orm"
	//"gopkg.in/mgo.v2/bson"
	"sync"
	"time"
)

type ValueEquation struct {
	sync.RWMutex
	orm.ModelBase `bson:"-" json:"-"`
	//Id                   int64 `bson:"Id" json:"Id"`
	Dates                time.Time `bson:"PeriodDates" json:"Dates"`
	Year                 int       `bson:"PeriodYear" json:"Year"`
	Month                int       `bson:"PeriodMonth" json:"Month"`
	Plant                string    `bson:"Plant" json:"Plant"`
	Unit                 string    `bson:"Unit" json:"Unit"`
	UnitGroup            string    `bson:"UnitGroup" json:"UnitGroup"`
	Phase                string    `bson:"Phase" json:"Phase"`
	Capacity             float64   `bson:"Capacity" json:"Capacity"`
	NetGeneration        float64   `bson:"NetGeneration" json:"NetGeneration"`
	AvgNetGeneration     float64   `bson:"AvgNetGeneration" json:"AvgNetGeneration"`
	PrctWAF              float64   `bson:"PrctWAF" json:"PrctWAF"`
	PrctWUF              float64   `bson:"PrctWUF" json:"PrctWUF"`
	MaxCapacity          float64   `bson:"MaxCapacity" json:"MaxCapacity"`
	MaxPowerGeneration   float64   `bson:"MaxPowerGeneration" json:"MaxPowerGeneration"`
	PotentialRevenue     float64   `bson:"PotentialRevenue" json:"PotentialRevenue"`
	ForcedOutages        float64   `bson:"ForcedOutages" json:"ForcedOutages"`
	ForcedOutagesLoss    float64   `bson:"ForcedOutagesLoss" json:"ForcedOutagesLoss"`
	UnforcedOutages      float64   `bson:"UnforcedOutages" json:"UnforcedOutages"`
	UnforcedOutagesLoss  float64   `bson:"UnforcedOutagesLoss" json:"UnforcedOutagesLoss"`
	CapacityPayment      float64   `bson:"CapacityPayment" json:"CapacityPayment"`
	EnergyPayment        float64   `bson:"EnergyPayment" json:"EnergyPayment"`
	PenaltyAmount        float64   `bson:"PenaltyAmount" json:"PenaltyAmount"`
	SRF                  float64   `bson:"SRF" json:"SRF"`
	VOMR                 float64   `bson:"VOMR" json:"VOMR"`
	UnplannedOutages     float64   `bson:"UnplannedOutages" json:"UnplannedOutages"`
	TotalOutageDuration  float64   `bson:"TotalOutageDuration" json:"TotalOutageDuration"`
	StartupPayment       float64   `bson:"StartupPayment" json:"StartupPayment"`
	Incentive            float64   `bson:"Incentive" json:"Incentive"`
	Revenue              float64   `bson:"Revenue" json:"Revenue"`
	Fuel                 []ValueEquationFuel
	PrimaryFuelTotalCost float64 `bson:"PrimaryFuelTotalCost" json:"PrimaryFuelTotalCost"`
	BackupFuelTotalCost  float64 `bson:"BackupFuelTotalCost" json:"BackupFuelTotalCost"`
	TotalFuelCost        float64 `bson:"TotalFuelCost" json:"TotalFuelCost"`
	FuelTransportCost    float64 `bson:"FuelTransportCost" json:"FuelTransportCost"`
	OperatingCost        float64 `bson:"OperatingCost" json:"OperatingCost"`
	Detail               []ValueEquationDetails
	Top10                []ValueEquationTop10
	TotalLabourCost      float64 `bson:"TotalLabourCost" json:"TotalLabourCost"`
	TotalMaterialCost    float64 `bson:"TotalMaterialCost" json:"TotalMaterialCost"`
	TotalServicesCost    float64 `bson:"TotalServicesCost" json:"TotalServicesCost"`
	TotalDuration        float64 `bson:"TotalDuration" json:"TotalDuration"`
	MaintenanceCost      float64 `bson:"MaintenanceCost" json:"MaintenanceCost"`
	ValueEquationCost    float64 `bson:"ValueEquationCost" json:"ValueEquationCost"`
}

func (v *ValueEquation) TableName() string {
	return "ValueEquation"
}

/*type ValueEquationPeriod struct {
	sync.RWMutex
	orm.ModelBase `bson:"-" json:"-"`
	Id            string    `bson:"Id" json:"Id"`
	Year          int       `bson:"Year" json:"Year"`
	Month         int       `bson:"Month" json:"Month"`
	MonthYear     int       `bson:"MonthYear" json:"MonthYear"`
	Quarter       int       `bson:"Quarter" json:"Quarter"`
	QuarterYear   int       `bson:"QuarterYear" json:"QuarterYear"`
	Dates         time.Time `bson:"DatesDates" json:"Dates"`
}

func (vp *ValueEquationPeriod) TableName() string {
	return "ValueEquationPeriod"
}*/

type ValueEquationFuel struct {
	sync.RWMutex
	orm.ModelBase `bson:"-" json:"-"`
	// Id            int64  `bson:"Id" json:"Id"`
	VEId                  int64   `bson:"VEId" json:"VEId"`
	IsPrimaryFuel         bool    `bson:"isPrimaryFuel" json:"isPrimaryFuel"`
	FuelType              string  `bson:"FuelType" json:"FuelType"`
	FuelCostPerUnit       float64 `bson:"FuelCostPerUnit" json:"FuelCostPerUnit"`
	FuelConsumed          float64 `bson:"FuelConsumed" json:"FuelConsumed"`
	ConvertedFuelConsumed float64 `bson:"ConvertedFuelConsumed" json:"ConvertedFuelConsumed"`
	FuelCost              float64 `bson:"FuelCost" json:"FuelCost"`
}

func (vf *ValueEquationFuel) TableName() string {
	return "ValueEquationFuelData"
}

type ValueEquationDetails struct {
	sync.RWMutex
	orm.ModelBase `bson:"-" json:"-"`
	// Id            int64  `bson:"Id" json:"Id"`
	VEId          int64   `bson:"VEId" json:"VEId"`
	DataSource    string  `bson:"DataSource" json:"DataSource"`
	WorkOrderType string  `bson:"WorkOrderType" json:"WorkOrderType"`
	Duration      float64 `bson:"Duration" json:"Duration"`
	LaborCost     float64 `bson:"LaborCost" json:"LaborCost"`
	MaterialCost  float64 `bson:"MaterialCost" json:"MaterialCost"`
	ServiceCost   float64 `bson:"ServiceCost" json:"ServiceCost"`
}

func (vd *ValueEquationDetails) TableName() string {
	return "ValueEquationDetails"
}

type ValueEquationTop10 struct {
	sync.RWMutex
	orm.ModelBase `bson:"-" json:"-"`
	// Id            			int64  `bson:"Id" json:"Id"`
	VEId                     int64   `bson:"VEId" json:"VEId"`
	WorkOrderID              string  `bson:"WorkOrderID" json:"WorkOrderID"`
	WorkOrderDescription     string  `bson:"WorkOrderDescription" json:"WorkOrderDescription"`
	EquipmentType            string  `bson:"EquipmentType" json:"EquipmentType"`
	EquipmentTypeDescription string  `bson:"EquipmentTypeDescription" json:"EquipmentTypeDescription"`
	WorkOrderType            string  `bson:"WorkOrderType" json:"WorkOrderType"`
	WorkOrderTypeDescription string  `bson:"WorkOrderTypeDescription" json:"WorkOrderTypeDescription"`
	MaintenanceActivity      string  `bson:"MaintenanceActivity" json:"MaintenanceActivity"`
	Duration                 float64 `bson:"Duration" json:"Duration"`
	LaborCost                float64 `bson:"LaborCost" json:"LaborCost"`
	MaterialCost             float64 `bson:"MaterialCost" json:"MaterialCost"`
	ServiceCost              float64 `bson:"ServiceCost" json:"ServiceCost"`
	MaintenanceCost          float64 `bson:"MaintenanceCost" json:"MaintenanceCost"`
}

func (vt *ValueEquationTop10) TableName() string {
	return "ValueEquationTop10"
}
