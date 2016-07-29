package models

import (
	"sync"

	"github.com/eaciit/orm"
)

type SampleAssetLevel struct {
	sync.RWMutex
	orm.ModelBase `bson:"-" json:"-"`
	//Id            string `bson:"_id" json:"_id"`
	Code      string `bson:"code" json:"code"`
	Levelname string `bson:"levelname" json:"levelname"`
}

func (m *SampleAssetLevel) TableName() string {
	return "SampleAssetLevel"
}
