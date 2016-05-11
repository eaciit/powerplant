package models

import (
	"github.com/eaciit/orm"
)

type ValueEquationDataQuality struct {
	orm.ModelBase `bson:"-",json:"-"`
}

func (m *ValueEquationDataQuality) TableName() string {
	return "ValueEquationDataQuality"
}
