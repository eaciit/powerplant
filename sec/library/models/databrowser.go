package models

import (
	"sync"
	"time"

	"github.com/eaciit/orm"
)

type DataBrowser struct {
	sync.RWMutex
	orm.ModelBase `bson:"-" json:"-"`
	// Id int64 `bson:"Id" json:"Id"`
	PeriodYear                             int    `bson:"Period.Year" json:"PeriodYear"`
	FunctionalLocation                     string `bson:"FunctionalLocation" json:"FunctionalLocation"`
	FLDescription                          string `bson:"FLDescription" json:"FLDescription"`
	IsTurbine                              bool   `bson:"IsTurbine" json:"IsTurbine"`
	IsSystem                               bool   `bson:"IsSystem" json:"IsSystem"`
	TurbineParent                          string `bson:"TurbineParent" json:"TurbineParent"`
	SystemParent                           string `bson:"SystemParent" json:"SystemParent"`
	AssetType                              string `bson:"AssetType" json:"AssetType"`
	EquipmentType                          string `bson:"EquipmentType" json:"EquipmentType"`
	EquipmentTypeDescription               string `bson:"EquipmentTypeDescription" json:"EquipmentTypeDescription"`
	PlantCode                              string `bson:"Plant.PlantCode" json:"PlantCode"`
	Plant                                  PowerPlantCoordinates
	TInfShortName                          string    `bson:"TurbineInfos.ShortName" json:"TInfShortName"`
	TInfManufacturer                       string    `bson:"TurbineInfos.Manufacturer" json:"TInfManufacturer"`
	TInfModel                              string    `bson:"TurbineInfos.Model" json:"TInfModel"`
	TInfUnitType                           string    `bson:"TurbineInfos.UnitType" json:"TInfUnitType"`
	TInfInstalledCapacity                  float64   `bson:"TurbineInfos.InstalledCapacity" json:"TInfInstalledCapacity"`
	TInfOperationalCapacity                float64   `bson:"TurbineInfos.OperationalCapacity" json:"TInfOperationalCapacity"`
	TInfPrimaryFuel                        string    `bson:"TurbineInfos.PrimaryFuel" json:"TInfPrimaryFuel"`
	TInfPrimaryFuel2                       string    `bson:"TurbineInfos.PrimaryFuel2" json:"TInfPrimaryFuel2"`
	TInfBackupFuel                         string    `bson:"TurbineInfos.BackupFuel" json:"TInfBackupFuel"`
	TInfHeatRate                           float64   `bson:"TurbineInfos.HeatRate" json:"TInfHeatRate"`
	TInfEfficiency                         float64   `bson:"TurbineInfos.Efficiency" json:"TInfEfficiency"`
	TInfCommisioningDate                   time.Time `bson:"TurbineInfos.CommisioningDate" json:"TInfCommisioningDate"`
	TInfRetirementPlan                     time.Time `bson:"TurbineInfos.RetirementPlan" json:"TInfRetirementPlan"`
	TInfInstalledMWH                       float64   `bson:"TurbineInfos.InstalledMWH" json:"TInfInstalledMWH"`
	TInfActualEnergyGeneration             float64   `bson:"TurbineInfos.ActualEnergyGeneration" json:"TInfActualEnergyGeneration"`
	TInfActualFuelConsumption_GASMMSCF     float64   `bson:"TurbineInfos.ActualFuelConsumption_GASMMSCF" json:"TInfActualFuelConsumption_GASMMSCF"`
	TInfActualFuelConsumption_CrudeBarrel  float64   `bson:"TurbineInfos.ActualFuelConsumption_CrudeBarrel" json:"TInfActualFuelConsumption_CrudeBarrel"`
	TInfActualFuelConsumption_HFOBarrel    float64   `bson:"TurbineInfos.ActualFuelConsumption_HFOBarrel" json:"TInfActualFuelConsumption_HFOBarrel"`
	TInfActualFuelConsumption_DieselBarrel float64   `bson:"TurbineInfos.ActualFuelConsumption_DieselBarrel" json:"TInfActualFuelConsumption_DieselBarrel"`
	TInfCapacityFactor                     float64   `bson:"TurbineInfos.CapacityFactor" json:"TInfCapacityFactor"`
	TInfUpdateEnergyGeneration             float64   `bson:"TurbineInfos.UpdateEnergyGeneration" json:"TInfUpdateEnergyGeneration"`
	TInfUpdateFuelConsumption              float64   `bson:"TurbineInfos.UpdateFuelConsumption" json:"TInfUpdateFuelConsumption"`
	TurbineVibrations                      []Vibration

	Maintenances         []AssetMaintenance
	FailureNotifications []NotificationFailureNoYear
	MROElements          []MaintenanceCost

	Operationals []OperationalData
	// Plant Plant
	// Outages [] ---> cannot determine yet from which table
}

func (m *DataBrowser) TableName() string {
	return "DataBrowser"
}

type AssetMaintenance struct {
	sync.RWMutex
	orm.ModelBase          `bson:"-" json:"-"`
	Id                     int64     `bson:"Id" json:"Id"`
	WorkOrderType          string    `bson:"WorkOrderType" json:"WorkOrderType"`
	MaintenanceOrder       string    `bson:"MaintenanceOrder" json:"MaintenanceOrder"`
	MaintenanceDescription string    `bson:"MaintenanceDescription" json:"MaintenanceDescription"`
	ActivityType           string    `bson:"ActivityType" json:"ActivityType"`
	PlanStart              time.Time `bson:"PlanStart" json:"PlanStart"`
	PlanFinish             time.Time `bson:"PlanFinish" json:"PlanFinish"`
	PlanDuration           int       `bson:"PlanDuration" json:"PlanDuration"`
	ActualStart            time.Time `bson:"ActualStart" json:"ActualStart"`
	ActualFinish           time.Time `bson:"ActualFinish" json:"ActualFinish"`
	ActualDuration         int       `bson:"ActualDuration" json:"ActualDuration"`
	LastMaintenanceDate    time.Time `bson:"LastMaintenanceDate" json:"LastMaintenanceDate"`
	MaintenanceInterval    int       `bson:"MaintenanceInterval" json:"MaintenanceInterval"`
	MaintenanceCost        float64   `bson:"MaintenanceCost" json:"MaintenanceCost"`
	DataBrowserId          int64     `bson:"DataBrowserId" json:"DataBrowserId"`
}

func (m *AssetMaintenance) TableName() string {
	return "AssetMaintenance"
}
