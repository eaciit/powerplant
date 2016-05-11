package models

import (
	"github.com/eaciit/orm"
)

type ValueEquation struct {
	orm.ModelBase `bson:"-",json:"-"`
	Id            string ` bson:"_id" , json:"_id" `
}

func (m *ValueEquation) TableName() string {
	return "ValueEquation"
}
