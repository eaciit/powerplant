package models

import (
	"github.com/eaciit/orm"
	// "gopkg.in/mgo.v2/bson"
)

type PhaseModel struct {
	orm.ModelBase `bson:"-",json:"-"`
	Id            string ` bson:"_id" , json:"_id" `
	Phase  string
}

func NewPhaseModel() *PhaseModel {
	m := new(PhaseModel)
	// m.Id = bson.NewObjectId()
	return m
}

func (e *PhaseModel) RecordID() interface{} {
	return e.Id
}

func (m *PhaseModel) TableName() string {
	return "MasterPhase"
}
