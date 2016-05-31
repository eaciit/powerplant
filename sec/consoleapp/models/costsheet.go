package models

import (
	"sync"

	"github.com/eaciit/orm"
)

type CostSheet struct {
	sync.RWMutex
	orm.ModelBase `bson:"-" json:"-"`
	Id            string             `bson:"_id" json:"Id"`
	Year          string             `bson:"Year" json:"Year"`
	Plant         string             `bson:"Plant" json:"Plant"`
	Details       []CostSheetDetails `bson:"Details" json:"Details"`
}

func (m *CostSheet) TableName() string {
	return "CostSheet"
}

type CostSheetDetails struct {
	sync.RWMutex
	orm.ModelBase `bson:"-" json:"-"`
	// Id              int64   `bson:"Id" json:"Id"`
	CostSheet       string  `bson:"CostSheet" json:"CostSheet"`
	CostElementCode string  `bson:"CostElementCode" json:"CostElementCode"`
	CostElementDesc string  `bson:"CostElementDesc" json:"CostElementDesc"`
	ActualCost      float64 `bson:"ActualCosts" json:"ActualCost"`
	PlannedCost     float64 `bson:"PlannedCosts" json:"PlannedCost"`
	VarAbs          float64 `bson:"VarAbs" json:"VarAbs"`
	VarPerc         float64 `bson:"VarPerc" json:"VarPerc"`
}

func (m *CostSheetDetails) TableName() string {
	return "CostSheetDetails"
}
