package models

import (
	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type SelectedScenario struct {
	ID    string  `bson:"ID"`
	Name  string  `bson:"Name"`
	Value float64 `bson:"Value"`
}

type SimulationResult struct {
	Revenue         float64 `bson:"Revenue"`
	LaborCost       float64 `bson:"LaborCost"`
	MaterialCost    float64 `bson:"MaterialCost"`
	ServiceCost     float64 `bson:"ServiceCost"`
	OperatingCost   float64 `bson:"OperatingCost"`
	MaintenanceCost float64 `bson:"MaintenanceCost"`
	ValueEquation   float64 `bson:"ValueEquation"`
}

type ScenarioSimulationModel struct {
	orm.ModelBase    `bson:"-",json:"-"`
	ID               bson.ObjectId      `bson:"_id" , json:"_id"`
	Start_Period     time.Time          `bson:"Start_Period"`
	End_Period       time.Time          `bson:"End_Period"`
	Name             string             `bson:"Name"`
	Description      string             `bson:"Description"`
	SelectedPlant    []string           `bson:"SelectedPlant"`
	SelectedUnit     []string           `bson:"SelectedUnit"`
	SelectedScenario []SelectedScenario `bson:"SelectedScenario"`
	HistoricResult   SimulationResult   `bson:"HistoricResult"`
	FutureResult     SimulationResult   `bson:"FutureResult"`
	Differential     SimulationResult   `bson:"Differential"`
	Last_Update      time.Time          `bson:"Last_Update"`
}

func (a *ScenarioSimulationModel) RecordID() interface{} {
	return a.ID
}

func NewScenarioSimulation() *ScenarioSimulationModel {
	e := new(ScenarioSimulationModel)
	return e
}

func (c *ScenarioSimulationModel) TableName() string {
	return "ScenarioSimulation"
}
