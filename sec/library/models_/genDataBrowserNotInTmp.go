package models

import (
	"sync"

	"github.com/eaciit/orm"
)

type GenDataBrowserNotInTmp struct {
	sync.RWMutex
	orm.ModelBase `bson:"-" json:"-"`
	ID            string `bson:"ID" json:"ID"`
	FLCode        string `bson:"FLCode" json:"FLCode"`
}

func (m *GenDataBrowserNotInTmp) TableName() string {
	return "GenDataBrowserNotInTmp"
}
