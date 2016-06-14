package models

import (
	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
)

type PlantModel struct {
	orm.ModelBase `bson:"-",json:"-"`
	Id            bson.ObjectId ` bson:"_id" , json:"Id" `
	Plant         string        ` bson:"Plant" , json:"Plant" `
}

func NewPlantModel() *PlantModel {
	m := new(PlantModel)
	m.Id = bson.NewObjectId()
	return m
}

func (e *PlantModel) RecordID() interface{} {
	return e.Id
}

func (m *PlantModel) TableName() string {
	return "MasterPlant"
}
