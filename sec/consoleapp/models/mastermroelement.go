package models

import "github.com/eaciit/orm"

type MasterMROElement struct {
	orm.ModelBase `bson:"-",json:"-"`
	Element       string `bson:"Element",json:"Element"`
}

func (m *MasterMROElement) TableName() string {
	return "MasterMROElement"
}
