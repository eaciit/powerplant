package models

import (
	"github.com/eaciit/orm"
	"sync"
)

type ValueEquation struct {
	orm.ModelBase `bson:"-",json:"-"`
	Id            string ` bson:"_id" , json:"_id" `
}

func (m *ValueEquation) TableName() string {
	return "ValueEquation"
}

type ValueEquationFuel struct {
	sync.RWMutex
	orm.ModelBase         `bson:"-" json:"-"`
	VEId                  string  `bson:"VEId" json:"VEId"`
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
	VEId          string  `bson:"VEId" json:"VEId"`
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
	orm.ModelBase            `bson:"-" json:"-"`
	VEId                     string  `bson:"VEId" json:"VEId"`
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
