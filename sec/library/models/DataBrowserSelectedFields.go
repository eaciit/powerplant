package models

import (
	"sync"

	"github.com/eaciit/orm"
)

type DataBrowserSelectedFields struct {
	sync.RWMutex
	orm.ModelBase   `bson:"-" json:"-"`
	FieldsReference string `bson:"FieldsReference" json:"FieldsReference"`
	Hypothesis      string `bson:"Hypothesis" json:"Hypothesis"`
	SelectedFields  string `bson:"SelectedFields" json:"SelectedFields"`
}

func (m *DataBrowserSelectedFields) TableName() string {
	return "DataBrowserSelectedFields"
}
