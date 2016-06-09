package models

import (
	"sync"

	"github.com/eaciit/orm"
)

type SampleAssetClass struct {
	sync.RWMutex
	orm.ModelBase `bson:"-" json:"-"`
	//Id            string `bson:"_id" json:"_id"`
	Code      string `bson:"code" json:"code"`
	Classname string `bson:"classname" json:"classname"`
}

func (m *SampleAssetClass) TableName() string {
	return "SampleAssetClass"
}
