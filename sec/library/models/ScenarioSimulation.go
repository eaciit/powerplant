package models

import (
	"github.com/eaciit/orm"
	//"gopkg.in/mgo.v2/bson"
	"sync"
	"time"
)

type ScenarioSimulation struct {
	sync.RWMutex
	orm.ModelBase                 `bson:"-" json:"-"`
	Id                            string    `bson:"_id"`
	StartPeriod                   time.Time `bson:"Start_Period" json:"StartPeriod"`
	EndPeriod                     time.Time `bson:"End_Period" json:"EndPeriod"`
	Name                          string    `bson:"Name" json:"Name"`
	Description                   string    `bson:"Description" json:"Description"`
	SelectedPlant                 []ScenarioSimulationSelectedPlant
	SelectedUnit                  []ScenarioSimulationSelectedUnit
	SelectedScenario              []ScenarioSimulationSelectedScenario
	HistoricResultRevenue         float64   `bson:"HistoricResultRevenue" json:"HistoricResultRevenue"`
	HistoricResultLaborCost       float64   `bson:"HistoricResultLaborCost" json:"HistoricResultLaborCost"`
	HistoricResultMaterialCost    float64   `bson:"HistoricResultMaterialCost" json:"HistoricResultMaterialCost"`
	HistoricResultServiceCost     float64   `bson:"HistoricResultServiceCost" json:"HistoricResultServiceCost"`
	HistoricResultOperationCost   float64   `bson:"HistoricResultOperationCost" json:"HistoricResultOperationCost"`
	HistoricResultMaintenanceCost float64   `bson:"HistoricResultMaintenanceCost" json:"HistoricResultMaintenanceCost"`
	HistoricResultValueEquation   float64   `bson:"HistoricResultValueEquation" json:"HistoricResultValueEquation"`
	FutureResultRevenue           float64   `bson:"FutureResultRevenue" json:"FutureResultRevenue"`
	FutureResultLaborCost         float64   `bson:"FutureResultLaborCost" json:"FutureResultLaborCost"`
	FutureResultMaterialCost      float64   `bson:"FutureResultMaterialCost" json:"FutureResultMaterialCost"`
	FutureResultServiceCost       float64   `bson:"FutureResultServiceCost" json:"FutureResultServiceCost"`
	FutureResultOperationCost     float64   `bson:"FutureResultOperationCost" json:"FutureResultOperationCost"`
	FutureResultMaintenanceCost   float64   `bson:"FutureResultMaintenanceCost" json:"FutureResultMaintenanceCost"`
	FutureResultValueEquation     float64   `bson:"FutureResultValueEquation" json:"FutureResultValueEquation"`
	DifferentialRevenue           float64   `bson:"DifferentialRevenue" json:"DifferentialRevenue"`
	DifferentialLaborCost         float64   `bson:"DifferentialLaborCost" json:"DifferentialLaborCost"`
	DifferentialMaterialCost      float64   `bson:"DifferentialMaterialCost" json:"DifferentialMaterialCost"`
	DifferentialServiceCost       float64   `bson:"DifferentialServiceCost" json:"DifferentialServiceCost"`
	DifferentialOperationCost     float64   `bson:"DifferentialOperationCost" json:"DifferentialOperationCost"`
	DifferentialMaintenanceCost   float64   `bson:"DifferentialMaintenanceCost" json:"DifferentialMaintenanceCost"`
	DifferentialValueEquation     float64   `bson:"DifferentialValueEquation" json:"DifferentialValueEquation"`
	LastUpdate                    time.Time `bson:"Last_Update" json:"LastUpdate"`
}

func (s *ScenarioSimulation) TableName() string {
	return "ScenarioSimulation"
}

type ScenarioSimulationSelectedPlant struct {
	sync.RWMutex
	orm.ModelBase `bson:"-" json:"-"`
	SSId          string `bson:"SSId" json:"SSId"`
	Plant         string `bson:"Plant" json:"Plant"`
}

func (sp *ScenarioSimulationSelectedPlant) TableName() string {
	return "ScenarioSimulationSelectedPlant"
}

type ScenarioSimulationSelectedUnit struct {
	sync.RWMutex
	orm.ModelBase `bson:"-" json:"-"`
	SSId          string `bson:"SSId" json:"SSId"`
	Unit          string `bson:"Unit" json:"Unit"`
}

func (su *ScenarioSimulationSelectedUnit) TableName() string {
	return "ScenarioSimulationSelectedUnit"
}

type ScenarioSimulationSelectedScenario struct {
	sync.RWMutex
	orm.ModelBase `bson:"-" json:"-"`
	SSId          string  `bson:"SSId" json:"SSId"`
	ID            string  `bson:"ID" json:"ID"`
	Name          string  `bson:"Name" json:"Name"`
	Value         float64 `bson:"Value" json:"Value"`
}

func (ss *ScenarioSimulationSelectedScenario) TableName() string {
	return "ScenarioSimulationSelectedScenario"
}
