package models

import (
	"github.com/eaciit/orm"
	"sync"
)

type TempMstPlant struct {
	sync.RWMutex
	orm.ModelBase `bson:"-" json:"-"`
	//Id                     int64                             `bson:"Id" json:"Id"`
	PlantCode string `bson:"PlantCode" json:"PlantCode"`
	PlantName string `bson:"PlantName" json:"PlantName"`
	PlantType string `bson:"PlantType" json:"PlantType"`
}

func (m *TempMstPlant) TableName() string {
	return "TempMstPlant"
}
