package models

import (
	"github.com/eaciit/orm"
)

type MasterPlant struct {
	orm.ModelBase `bson:"-",json:"-"`
	Plant         string `bson:"Plant",json:"Plant"`
}

func (m *MasterPlant) TableName() string {
	return "MasterPlant"
}
