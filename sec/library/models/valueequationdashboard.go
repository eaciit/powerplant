package models

import (
	"github.com/eaciit/orm"
	"sync"
	"time"
)

type ValueEquationDashboard struct {
	sync.RWMutex
	orm.ModelBase `bson:"-" json:"-"`
	// Id                                 int64    `bson:"Id" json:"Id"`
	Year                               int       `bson:"Period.Year" json:"Year"`
	Month                              int       `bson:"Period.Month" json:"Month"`
	MonthYear                          int       `bson:"Period.MonthYear" json:"MonthYear"`
	Quarter                            int       `bson:"Period.Quarter" json:"Quarter"`
	QuarterYear                        int       `bson:"Period.QuarterYear" json:"QuarterYear"`
	Dates                              time.Time `bson:"Period.Dates" json:"Dates"`
	Plant                              string    `bson:"Plant" json:"Plant"`
	PlantCode                          string    `bson:"PlantDetail.PlantCode" json:"PlantCode"`
	PlantType                          string    `bson:"PlantDetail.PlantType" json:"PlantType"`
	Province                           string    `bson:"PlantDetail.Province" json:"Province"`
	Region                             string    `bson:"PlantDetail.Region" json:"Region"`
	City                               string    `bson:"PlantDetail.City" json:"City"`
	FuelTypes_Crude                    bool      `bson:"PlantDetail.FuelTypes_Crude" json:"FuelTypes_Crude"`
	FuelTypes_Heavy                    bool      `bson:"PlantDetail.FuelTypes_Heavy" json:"FuelTypes_Heavy"`
	FuelTypes_Diesel                   bool      `bson:"PlantDetail.FuelTypes_Diesel" json:"FuelTypes_Diesel"`
	FuelTypes_Gas                      bool      `bson:"PlantDetail.FuelTypes_Gas" json:"FuelTypes_Gas"`
	GasTurbineUnit                     float64   `bson:"PlantDetail.GasTurbineUnit" json:"GasTurbineUnit"`
	GasTurbineCapacity                 float64   `bson:"PlantDetail.GasTurbineCapacity" json:"GasTurbineCapacity"`
	SteamUnit                          float64   `bson:"PlantDetail.SteamUnit" json:"SteamUnit"`
	SteamCapacity                      float64   `bson:"PlantDetail.SteamCapacity" json:"SteamCapacity"`
	DieselUnit                         float64   `bson:"PlantDetail.DieselUnit" json:"DieselUnit"`
	DieselCapacity                     float64   `bson:"PlantDetail.DieselCapacity" json:"DieselCapacity"`
	CombinedCycleUnit                  float64   `bson:"PlantDetail.CombinedCycleUnit" json:"CombinedCycleUnit"`
	CombinedCycleCapacity              float64   `bson:"PlantDetail.CombinedCycleCapacity" json:"CombinedCycleCapacity"`
	Longitude                          float64   `bson:"PlantDetail.Longitude" json:"Longitude"`
	Latitude                           float64   `bson:"PlantDetail.Latitude" json:"Latitude"`
	Unit                               string    `bson:"Unit" json:"Unit"`
	UnitGroup                          string    `bson:"UnitGroup" json:"UnitGroup"`
	ShortName                          string    `bson:"TurbineInfos.ShortName" json:"ShortName"`
	Manufacturer                       string    `bson:"TurbineInfos.Manufacturer" json:"Manufacturer"`
	Model                              string    `bson:"TurbineInfos.Model" json:"Model"`
	UnitType                           string    `bson:"TurbineInfos.UnitType" json:"UnitType"`
	InstalledCapacity                  float64   `bson:"TurbineInfos.InstalledCapacity" json:"InstalledCapacity"`
	OperationalCapacity                float64   `bson:"TurbineInfos.OperationalCapacity" json:"OperationalCapacity"`
	PrimaryFuel                        string    `bson:"TurbineInfos.PrimaryFuel" json:"PrimaryFuel"`
	PrimaryFuel2                       string    `bson:"TurbineInfos.PrimaryFuel2" json:"PrimaryFuel2"`
	BackupFuel                         string    `bson:"TurbineInfos.BackupFuel" json:"BackupFuel"`
	HeatRate                           float64   `bson:"TurbineInfos.HeatRate" json:"HeatRate"`
	Efficiency                         float64   `bson:"TurbineInfos.Efficiency" json:"Efficiency"`
	CommisioningDate                   time.Time `bson:"TurbineInfos.CommisioningDate" json:"CommisioningDate"`
	RetirementPlan                     string    `bson:"TurbineInfos.RetirementPlan" json:"RetirementPlan"`
	InstalledMWH                       float64   `bson:"TurbineInfos.InstalledMWH" json:"InstalledMWH"`
	ActualEnergyGeneration             float64   `bson:"TurbineInfos.ActualEnergyGeneration" json:"ActualEnergyGeneration"`
	ActualFuelConsumption_GASMMSCF     float64   `bson:"TurbineInfos.ActualFuelConsumption_GASMMSCF" json:"ActualFuelConsumption_GASMMSCF"`
	ActualFuelConsumption_CrudeBarrel  float64   `bson:"TurbineInfos.ActualFuelConsumption_CrudeBarrel" json:"ActualFuelConsumption_CrudeBarrel"`
	ActualFuelConsumption_HFOBarrel    float64   `bson:"TurbineInfos.ActualFuelConsumption_HFOBarrel" json:"ActualFuelConsumption_HFOBarrel"`
	ActualFuelConsumption_DieselBarrel float64   `bson:"TurbineInfos.ActualFuelConsumption_DieselBarrel" json:"ActualFuelConsumption_DieselBarrel"`
	CapacityFactor                     float64   `bson:"TurbineInfos.CapacityFactor" json:"CapacityFactor"`
	UpdatedEnergyGeneration            float64   `bson:"TurbineInfos.UpdatedEnergyGeneration" json:"UpdatedEnergyGeneration"`
	UpdatedFuelConsumption             float64   `bson:"TurbineInfos.UpdatedFuelConsumption" json:"UpdatedFuelConsumption"`
	Phase                              string    `bson:"Phase" json:"Phase"`
	Capacity                           float64   `bson:"Capacity" json:"Capacity"`
	NetGeneration                      float64   `bson:"NetGeneration" json:"NetGeneration"`
	AvgNetGeneration                   float64   `bson:"AvgNetGeneration" json:"AvgNetGeneration"`
	PrctWAF                            float64   `bson:"PrctWAF" json:"PrctWAF"`
	PrctWUF                            float64   `bson:"PrctWUF" json:"PrctWUF"`
	MaxCapacity                        float64   `bson:"MaxCapacity" json:"MaxCapacity"`
	MaxPowerGeneration                 float64   `bson:"MaxPowerGeneration" json:"MaxPowerGeneration"`
	PotentialRevenue                   float64   `bson:"PotentialRevenue" json:"PotentialRevenue"`
	ForcedOutages                      float64   `bson:"ForcedOutages" json:"ForcedOutages"`
	ForcedOutagesLoss                  float64   `bson:"ForcedOutagesLoss" json:"ForcedOutagesLoss"`
	UnforcedOutages                    float64   `bson:"UnforcedOutages" json:"UnforcedOutages"`
	UnforcedOutagesLoss                float64   `bson:"UnforcedOutagesLoss" json:"UnforcedOutagesLoss"`
	CapacityPayment                    float64   `bson:"CapacityPayment" json:"CapacityPayment"`
	EnergyPayment                      float64   `bson:"EnergyPayment" json:"EnergyPayment"`
	PenaltyAmount                      float64   `bson:"PenaltyAmount" json:"PenaltyAmount"`
	SRF                                float64   `bson:"SRF" json:"SRF"`
	VOMR                               float64   `bson:"VOMR" json:"VOMR"`
	UnplannedOutages                   float64   `bson:"UnplannedOutages" json:"UnplannedOutages"`
	TotalOutageDuration                float64   `bson:"TotalOutageDuration" json:"TotalOutageDuration"`
	StartupPayment                     float64   `bson:"StartupPayment" json:"StartupPayment"`
	Incentive                          float64   `bson:"Incentive" json:"Incentive"`
	Revenue                            float64   `bson:"Revenue" json:"Revenue"`
	PrimaryFuelTotalCost               float64   `bson:"PrimaryFuelTotalCost" json:"PrimaryFuelTotalCost"`
	BackupFuelTotalCost                float64   `bson:"BackupFuelTotalCost" json:"BackupFuelTotalCost"`
	TotalFuelCost                      float64   `bson:"TotalFuelCost" json:"TotalFuelCost"`
	FuelTransportCost                  float64   `bson:"FuelTransportCost" json:"FuelTransportCost"`
	OperatingCost                      float64   `bson:"OperatingCost" json:"OperatingCost"`
	TotalLabourCost                    float64   `bson:"TotalLabourCost" json:"TotalLabourCost"`
	TotalMaterialCost                  float64   `bson:"TotalMaterialCost" json:"TotalMaterialCost"`
	TotalServicesCost                  float64   `bson:"TotalServicesCost" json:"TotalServicesCost"`
	TotalDuration                      float64   `bson:"TotalDuration" json:"TotalDuration"`
	MaintenanceCost                    float64   `bson:"MaintenanceCost" json:"MaintenanceCost"`
	ValueEquationCost                  float64   `bson:"ValueEquationCost" json:"ValueEquationCost"`
}

func (v *ValueEquationDashboard) TableName() string {
	return "ValueEquation_DashboardTest"
}

// Fuel
type VEDFuel struct {
	sync.RWMutex
	orm.ModelBase `bson:"-" json:"-"`
	// Id            int64  `bson:"Id" json:"Id"`
	VEId                  int64   `bson:"VEId" json:"VEId"`
	IsPrimaryFuel         bool    `bson:"isPrimaryFuel" json:"IsPrimaryFuel"`
	FuelType              string  `bson:"FuelType" json:"FuelType"`
	FuelCostPerUnit       float64 `bson:"FuelCostPerUnit" json:"FuelCostPerUnit"`
	FuelConsumed          float64 `bson:"FuelConsumed" json:"FuelConsumed"`
	ConvertedFuelConsumed float64 `bson:"ConvertedFuelConsumed" json:"ConvertedFuelConsumed"`
	FuelCost              float64 `bson:"FuelCost" json:"FuelCost"`
}

func (m *VEDFuel) TableName() string {
	return "VEDFuelTest"
}

// Detail
type VEDDetail struct {
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

func (m *VEDDetail) TableName() string {
	return "VEDDetailTest"
}

// VEDTop10
type VEDTop10 struct {
	sync.RWMutex
	orm.ModelBase `bson:"-" json:"-"`
	// Id            int64  `bson:"Id" json:"Id"`
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

func (m *VEDTop10) TableName() string {
	return "VEDTop10Test"
}
