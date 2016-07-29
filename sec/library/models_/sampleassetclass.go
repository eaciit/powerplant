package models

import (
	"sync"

	"github.com/eaciit/orm"
)

type SampleAssetClass struct {
	sync.RWMutex
	orm.ModelBase `bson:"-" json:"-"`
	Id            int64  `bson:"Id" json:"Id"`
	Code          string `bson:"code" json:"code"`
	Classname     string `bson:"classname" json:"classname"`
}

func (m *SampleAssetClass) TableName() string {
	return "SampleAssetClass"
}
