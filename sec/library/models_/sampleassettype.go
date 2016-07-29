package models

import (
	"sync"

	"github.com/eaciit/orm"
)

type SampleAssetType struct {
	sync.RWMutex
	orm.ModelBase `bson:"-" json:"-"`
	//Id            string `bson:"_id" json:"_id"`
	Code     string `bson:"code" json:"code"`
	Typename string `bson:"typename" json:"typename"`
}

func (m *SampleAssetType) TableName() string {
	return "SampleAssetType"
}
