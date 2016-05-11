package models

import (
	"github.com/eaciit/orm"
	"time"
)

type SelectedScenario struct {
	ID    string
	Name  string
	Value float64
}

type SimulationResult struct {
	Revenue         float64
	LaborCost       float64
	MaterialCost    float64
	ServiceCost     float64
	OperatingCost   float64
	MaintenanceCost float64
	ValueEquation   float64
}

type ScenarioSimulationModel struct {
	orm.ModelBase    `bson:"base"`
	Start_Period     time.Time
	End_Period       time.Time
	Name             string
	Description      string
	SelectedPlant    []string
	SelectedUnit     []string
	SelectedScenario []SelectedScenario
	HistoricResult   SimulationResult
	FutureResult     SimulationResult
	Differential     SimulationResult
	Last_Update      time.Time
}

func (c *ScenarioSimulationModel) TableName() string {
	return "ScenarioSimulation"
}
