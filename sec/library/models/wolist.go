package models

import (
	"sync"
	"time"

	"github.com/eaciit/orm"
)

type WOList struct {
	sync.RWMutex
	orm.ModelBase      `bson:"-" json:"-"`
	UserStatus         string    `bson:"UserStatus" json:"UserStatus"`
	SystemStatus       string    `bson:"SystemStatus" json:"SystemStatus"`
	Type               string    `bson:"Type" json:"Type"`
	OrderCode          string    `bson:"OrderCode" json:"OrderCode"`
	NotificationCode   string    `bson:"NotificationCode" json:"NotificationCode"`
	A                  string    `bson:"A" json:"A"`
	EnteredBy          string    `bson:"EnteredBy" json:"EnteredBy"`
	Description        string    `bson:"Description" json:"Description"`
	Plant              string    `bson:"Plant" json:"Plant"`
	FunctionalLocation string    `bson:"FunctionalLocation" json:"FunctionalLocation"`
	EquipmentCode      string    `bson:"EquipmentCode" json:"EquipmentCode"`
	SortField          string    `bson:"SortField" json:"SortField"`
	PriorityText       string    `bson:"PriorityText" json:"PriorityText"`
	WorkCtr            string    `bson:"WorkCtr" json:"WorkCtr"`
	CostCtr            string    `bson:"CostCtr" json:"CostCtr"`
	MAT                string    `bson:"MAT" json:"MAT"`
	CreatedOn          time.Time `bson:"CreatedOn" json:"CreatedOn"`
	BasicStart         time.Time `bson:"BasicStart" json:"BasicStart"`
	BasicFinish        time.Time `bson:"BasicFinish" json:"BasicFinish"`
	ScheduledStart     time.Time `bson:"ScheduledStart" json:"ScheduledStart"`
	ScheduledFinish    time.Time `bson:"ScheduledFinish" json:"ScheduledFinish"`
	ActualStart        time.Time `bson:"ActualStart" json:"ActualStart"`
	ActualFinish       time.Time `bson:"ActualFinish" json:"ActualFinish"`
	RefDate            time.Time `bson:"RefDate" json:"RefDate"`
	RespCCTR           string    `bson:"RespCCTR" json:"RespCCTR"`
	ReleaseDate        time.Time `bson:"ReleaseDate" json:"ReleaseDate"`
	AvailFrom          time.Time `bson:"AvailFrom" json:"AvailFrom"`
	AvailTo            time.Time `bson:"AvailTo" json:"AvailTo"`
	ActualCost         float64   `bson:"ActualCost" json:"ActualCost"`
	DF                 string    `bson:"DF" json:"DF"`
}

func (m *WOList) TableName() string {
	return "WOList"
}
