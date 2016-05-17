package models

import (
	"time"

	"github.com/eaciit/orm"
)


type MaintenanceWorkOrder struct {
	orm.ModelBase       `bson:"-",json:"-"`
	OrderCode           string    `bson:"OrderCode",json:"OrderCode"`
	OrderDescription    string    `bson:"OrderDescription",json:"OrderDescription"`
	FunctionalLocation  string    `bson:"FunctionalLocation",json:"FunctionalLocation"`
	UserStatus          string    `bson:"UserStatus",json:"UserStatus"`
	CreatedOn           time.Time `bson:"CreatedOn",json:"CreatedOn"`
	BasStartDate        time.Time `bson:"BasStartDate",json:"BasStartDate"`
	ActualRelease       time.Time `bson:"ActualRelease",json:"ActualRelease"`
	OrderType           string    `bson:"OrderType",json:"OrderType"`
	PlantWorkCtr        string    `bson:"PlantWorkCtr",json:"PlantWorkCtr"`
	CompanyCode         string    `bson:"CompanyCode",json:"CompanyCode"`
	SortField           string    `bson:"SortField",json:"SortField"`
	Description         string    `bson:"Description",json:"Description"`
	Equipment           string    `bson:"Equipment",json:"Equipment"`
	MainWorkCtr         string    `bson:"MainWorkCtr",json:"MainWorkCtr"`
	PlanningPlant       string    `bson:"PlanningPlant",json:"PlanningPlant"`
	CostCtr             string    `bson:"CostCtr",json:"CostCtr"`
	RespCostCtr         string    `bson:"RespCostCtr",json:"RespCostCtr"`
	ObjectNumber        string    `bson:"ObjectNumber",json:"ObjectNumber"`
	ProfitCtr           string    `bson:"ProfitCtr",json:"ProfitCtr"`
	Priority            string    `bson:"Priority",json:"Priority"`
	PriorityDescription string    `bson:"PriorityDescription",json:"PriorityDescription"`
	Notification        string    `bson:"Notification",json:"Notification"`
	Location            string    `bson:"Location",json:"Location"`
	SystemStatus        string    `bson:"SystemStatus",json:"SystemStatus"`
	MainPlant           string    `bson:"MainPlant",json:"MainPlant"`
}

func (m *MaintenanceWorkOrder) TableName() string {
	return "MaintenanceWorkOrder"
}
