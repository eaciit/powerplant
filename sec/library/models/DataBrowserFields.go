package models

import (
	"sync"

	"github.com/eaciit/orm"
)

type DataBrowserFields struct {
	sync.RWMutex
	orm.ModelBase   `bson:"-" json:"-"`
	FieldsReference string `bson:"FieldsReference" json:"FieldsReference"`
	Alias           string `bson:"Alias" json:"Alias"`
	Field           string `bson:"Field" json:"Field"`
	Type            string `bson:"Type" json:"Type"`
}

func (m *DataBrowserFields) TableName() string {
	return "DataBrowserFields"
}
