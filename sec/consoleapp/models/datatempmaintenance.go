package models

import (
	"github.com/eaciit/orm"
	"sync"
	"time"
)

type DataTempMaintenance struct {
	sync.RWMutex
	orm.ModelBase `bson:"-" json:"-"`
	// Id                       int       `bson:"Id" json:"Id"`
	Year                     int       `bson:"Period.Year" json:"Year"`
	Month                    int       `bson:"Period.Month" json:"Month"`
	MonthYear                int       `bson:"Period.MonthYear" json:"MonthYear"`
	Quarter                  int       `bson:"Period.Quarter" json:"Quarter"`
	QuarterYear              int       `bson:"Period.QuarterYear" json:"QuarterYear"`
	Dates                    time.Time `bson:"Period.Dates" json:"Dates"`
	FunctionalLocation       string    `bson:"FunctionalLocation" json:"FunctionalLocation"`
	FLDescription            string    `bson:"FLDescription" json:"FLDescription"`
	IsTurbine                bool      `bson:"IsTurbine" json:"IsTurbine"`
	TurbineInfos             string    `bson:"TurbineInfos" json:"TurbineInfos"`
	TurbineParent            string    `bson:"TurbineParent" json:"TurbineParent"`
	AssetType                string    `bson:"AssetType" json:"AssetType"`
	EquipmentType            string    `bson:"EquipmentType" json:"EquipmentType"`
	EquipmentTypeDescription string    `bson:"EquipmentTypeDescription" json:"EquipmentTypeDescription"`
	PlantCode                string    `bson:"Plant.PlantCode" json:"PlantCode"`
	PlantName                string    `bson:"Plant.PlantName" json:"PlantName"`
	PlantType                string    `bson:"Plant.PlantType" json:"PlantType"`
	Province                 string    `bson:"Plant.Province" json:"Province"`
	Region                   string    `bson:"Plant.Region" json:"Region"`
	City                     string    `bson:"Plant.City" json:"City"`
	FuelTypes_Crude          bool      `bson:"Plant.FuelTypes_Crude" json:"FuelTypes_Crude"`
	FuelTypes_Heavy          bool      `bson:"Plant.FuelTypes_Heavy" json:"FuelTypes_Heavy"`
	FuelTypes_Diesel         bool      `bson:"Plant.FuelTypes_Diesel" json:"FuelTypes_Diesel"`
	FuelTypes_Gas            bool      `bson:"Plant.FuelTypes_Gas" json:"FuelTypes_Gas"`
	GasTurbineUnit           float64   `bson:"Plant.GasTurbineUnit" json:"GasTurbineUnit"`
	GasTurbineCapacity       float64   `bson:"Plant.GasTurbineCapacity" json:"GasTurbineCapacity"`
	SteamUnit                float64   `bson:"Plant.SteamUnit" json:"SteamUnit"`
	SteamCapacity            float64   `bson:"Plant.SteamCapacity" json:"SteamCapacity"`
	DieselUnit               float64   `bson:"Plant.DieselUnit" json:"DieselUnit"`
	DieselCapacity           float64   `bson:"Plant.DieselCapacity" json:"DieselCapacity"`
	CombinedCycleUnit        float64   `bson:"Plant.CombinedCycleUnit" json:"CombinedCycleUnit"`
	CombinedCycleCapacity    float64   `bson:"Plant.CombinedCycleCapacity" json:"CombinedCycleCapacity"`
	WorkOrderType            string    `bson:"Maintenance.WorkOrderType" json:"WorkOrderType"`
	MaintenanceOrder         string    `bson:"Maintenance.MaintenanceOrder" json:"MaintenanceOrder"`
	MaintenanceDescription   string    `bson:"Maintenance.MaintenanceDescription" json:"MaintenanceDescription"`
	ActivityType             string    `bson:"Maintenance.ActivityType" json:"ActivityType"`
	PlanStart                time.Time `bson:"Maintenance.PlanStart" json:"PlanStart"`
	PlanFinish               time.Time `bson:"Maintenance.PlanFinish" json:"PlanFinish"`
	PlanDuration             float64   `bson:"Maintenance.PlanDuration" json:"PlanDuration"`
	ActualStart              time.Time `bson:"Maintenance.ActualStart" json:"ActualStart"`
	ActualFinish             time.Time `bson:"Maintenance.ActualFinish" json:"ActualFinish"`
	ActualDuration           float64   `bson:"Maintenance.ActualDuration" json:"ActualDuration"`
	LastMaintenanceDate      time.Time `bson:"Maintenance.LastMaintenanceDate" json:"LastMaintenanceDate"`
	MaintenanceInterval      float64   `bson:"Maintenance.MaintenanceInterval" json:"MaintenanceInterval"`
	MaintenanceCost          float64   `bson:"Maintenance.MaintenanceCost" json:"MaintenanceCost"`
}

func (m *DataTempMaintenance) TableName() string {
	return "DataTempMaintenance"
}
