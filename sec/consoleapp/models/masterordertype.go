package models

import "github.com/eaciit/orm"

type MasterOrderType struct {
	orm.ModelBase `bson:"-",json:"-"`
	Id            string `bson:"id",json:"id"`
	OrderTypeDesc string `bson:"OrderTypeDesc",json:"OrderTypeDesc"`
}

func (m *MasterOrderType) TableName() string {
	return "MasterOrderType"
}
