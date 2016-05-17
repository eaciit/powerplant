package models

import (
	"github.com/eaciit/orm"
)

type OperationalData struct {
	base.ConvertMGOToSQLServer(new(WOList))
	orm.ModelBase                 `bson:"-",json:"-"`
	Name                          string  `bson:"Name",json:"Name"`
	Year                          int     `bson:"Year",json:"Year"`
	GenerationGross               float64 `bson:"GenerationGross",json:"GenerationGross"`
	GenerationAux                 float64 `bson:"GenerationAux",json:"GenerationAux"`
	GenerationNet                 float64 `bson:"GenerationNet",json:"GenerationNet"`
	ServiceHours                  float64 `bson:"ServiceHours",json:"ServiceHours"`
	ReserveShutdownHours          float64 `bson:"ReserveShutdownHours",json:"ReserveShutdownHours"`
	MaintenanceOutageHours        float64 `bson:"MaintenanceOutageHours",json:"MaintenanceOutageHours"`
	ExtendedOutageMaitenanceHours float64 `bson:"ExtendedOutageMaitenanceHours",json:"ExtendedOutageMaitenanceHours"`
	PlantOutageHours              float64 `bson:"PlantOutageHours",json:"PlantOutageHours"`
	ExtendedPlanHours             float64 `bson:"ExtendedPlanHours",json:"ExtendedPlanHours"`
	ForcedOutageHours             float64 `bson:"ForcedOutageHours",json:"ForcedOutageHours"`
	OutOfManagementControl        float64 `bson:"OutOfManagementControl",json:"OutOfManagementControl"`
	MonthBall                     float64 `bson:"MonthBall",json:"MonthBall"`
	UnderCommisionFO              float64 `bson:"UnderCommisionFO",json:"UnderCommisionFO"`
	UnderCommisionIS              float64 `bson:"UnderCommisionIS",json:"UnderCommisionIS"`
	UnderCommisionMO              float64 `bson:"UnderCommisionMO",json:"UnderCommisionMO"`
	UnderCommisionRS              float64 `bson:"UnderCommisionRS",json:"UnderCommisionRS"`
	IR                            float64 `bson:"IR",json:"IR"`
	PH                            float64 `bson:"PH",json:"PH"`
	NoOfStartAttempt              int     `bson:"NoOfStartAttempt",json:"NoOfStartAttempt"`
	NoOfStartActual               int     `bson:"NoOfStartActual",json:"NoOfStartActual"`
	FuelDiesel                    float64 `bson:"FuelDiesel",json:"FuelDiesel"`
	FuelCrude                     float64 `bson:"FuelCrude",json:"FuelCrude"`
	FuelHeavy                     float64 `bson:"FuelHeavy",json:"FuelHeavy"`
	FuelGas                       float64 `bson:"FuelGas",json:"FuelGas"`
}

func (m *OperationalData) TableName() string {
	return "OperationalData"
}
