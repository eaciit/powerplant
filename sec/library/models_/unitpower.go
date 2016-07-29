package models

import (
	"sync"

	"github.com/eaciit/orm"
	//"gopkg.in/mgo.v2/bson"
)

type UnitPower struct {
	sync.RWMutex
	orm.ModelBase `bson:"-" json:"-"`
	//Id            bson.ObjectId `bson:"_id" json:"id"`
	Plant    string  `bson:"Plant" json:"Plant"`
	Year     int     `bson:"Year" json:"Year"`
	Unit     string  `bson:"Unit" json:"Unit"`
	MaxPower float64 `bson:"MaxPower" json:"MaxPower"`
}

func (m *UnitPower) TableName() string {
	return "UnitPower"
}
