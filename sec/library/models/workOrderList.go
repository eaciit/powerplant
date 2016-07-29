package models

import (
	"sync"
	"time"

	"github.com/eaciit/orm"
)

type WorkOrderList struct {
	sync.RWMutex
	orm.ModelBase      `bson:"-" json:"-"`
	Id                 int64     `bson:"Id" json:"Id"`
	Plant              string    `bson:"Plant" json:"Plant"`
	FunctionalLocation string    `bson:"FunctionalLocation" json:"FunctionalLocation"`
	Type               string    `bson:"Type" json:"Type"`
	OrderCode          string    `bson:"OrderCode" json:"OrderCode"`
	Description        string    `bson:"Description" json:"Description"`
	ActualStart        time.Time `bson:"ActualStart" json:"ActualStart"`
	ActualFinish       time.Time `bson:"ActualFinish" json:"ActualFinish"`
	ActualCost         float64   `bson:"ActualCost" json:"ActualCost"`
}

func (m *WorkOrderList) TableName() string {
	return "WorkOrderList"
}
