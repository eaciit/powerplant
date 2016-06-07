package models

import (
	"sync"

	"github.com/eaciit/orm"
)

type MasterMROElement struct {
	sync.RWMutex
	orm.ModelBase `bson:"-" json:"-"`
	Element       string `bson:"Element" json:"Element"`
}

func (m *MasterMROElement) TableName() string {
	return "MasterMROElement"
}
