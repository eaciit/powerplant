package models

import (
	"sync"
	"time"

	"github.com/eaciit/orm"
)

type RPPCloseWO struct {
	sync.RWMutex
	orm.ModelBase `bson:"-",json:"-"`
	// id	int64	`bson:"id",json:"id"`
	Notification       string    `bson:"Notification",json:"Notification"`
	OrderCode          string    `bson:"OrderCode",json:"OrderCode"`
	UserStatus         string    `bson:"UserStatus",json:"UserStatus"`
	FunctionalLocation string    `bson:"FunctionalLocation",json:"FunctionalLocation"`
	MainWorkCtr        string    `bson:"MainWorkCtr",json:"MainWorkCtr"`
	Description        string    `bson:"Description",json:"Description"`
	OrderType          string    `bson:"OrderType",json:"OrderType"`
	ReferenceDate      time.Time `bson:"ReferenceDate",json:"ReferenceDate"`
	WorkCenter         string    `bson:"WorkCenter",json:"WorkCenter"`
	Equipment          string    `bson:"Equipment",json:"Equipment"`
	TotalActCost       float64   `bson:"TotalActCost",json:"TotalActCost"`
	ActualStart        time.Time `bson:"ActualStart",json:"ActualStart"`
	CostCenter         string    `bson:"CostCenter",json:"CostCenter"`
	TotalSettlement    float64   `bson:"TotalSettlement",json:"TotalSettlement"`
	PlannerGroup       string    `bson:"PlannerGroup",json:"PlannerGroup"`
	WBSElement         string    `bson:"WBSElement",json:"WBSElement"`
	SystemStatus       string    `bson:"SystemStatus",json:"SystemStatus"`
}

func (m *RPPCloseWO) TableName() string {
	return "RPPCloseWO"
}
