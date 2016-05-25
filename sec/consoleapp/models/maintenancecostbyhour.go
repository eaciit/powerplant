package models

import (
	"sync"
	"time"

	"github.com/eaciit/orm"
)

type MaintenanceCostByHour struct {
	sync.RWMutex
	orm.ModelBase           `bson:"-" json:"-"`
	Id                      string    `bson:"_id" json:"id"`
	OrderTypeDesc           string    `bson:"OrderTypeDesc" json:"OrderTypeDesc"`
	OrderType               string    `bson:"OrderType" json:"OrderType"`
	MaintenanceOrder        string    `bson:"MaintenanceOrder" json:"MaintenanceOrder"`
	MaintenanceOrderDesc    string    `bson:"MaintenanceOrderDesc" json:"MaintenanceOrderDesc"`
	EquipmentType           string    `bson:"EquipmentType" json:"EquipmentType"`
	EquipmentTypeDesc       string    `bson:"EquipmentTypeDesc" json:"EquipmentTypeDesc"`
	Equipment               string    `bson:"Equipment" json:"Equipment"`
	EquipmentDesc           string    `bson:"EquipmentDesc" json:"EquipmentDesc"`
	MaintenanceActivityType string    `bson:"MaintenanceActivityType" json:"MaintenanceActivityType"`
	PlanVal                 float64   `bson:"Plan" json:"PlanVal"`
	Period                  time.Time `bson:"Period" json:"Period"`
	Plant                   string    `bson:"Plant" json:"Plant"`
	Actual                  float64   `bson:"Actual" json:"Actual"`
	VarPerc                 float64   `bson:"VarPerc" json:"VarPerc"`
}

func (m *MaintenanceCostByHour) TableName() string {
	return "MaintenanceCostByHour"
}
