package models

import (
	"sync"

	"github.com/eaciit/orm"
)

type MasterFunctionalLocation struct {
	sync.RWMutex
	orm.ModelBase              `bson:"-" json:"-"`
	Plant                      string `bson:"Plant" json:"Plant"`
	FunctionalLocationCode     string `bson:"FunctionalLocationCode" json:"FunctionalLocationCode"`
	Description                string `bson:"FLDescription" json:"Description"`
	SuperiorFunctionalLocation string `bson:"SupFunctionalLocation" json:"SupFunctionalLocation"`
	EquipmentType              string `bson:"EquipmentType" json:"EquipmentType"`
	IsTurbine                  bool   `bson:"IsTurbine" json:"IsTurbine"`
}

func (m *MasterFunctionalLocation) TableName() string {
	return "MasterFunctionalLocation"
}
