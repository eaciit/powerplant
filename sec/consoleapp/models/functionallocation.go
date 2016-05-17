package models

import (
	"sync"
	"time"

	"github.com/eaciit/orm"
)

type FunctionalLocation struct {
	sync.RWMutex
	orm.ModelBase          `bson:"-",json:"-"`
	FunctionalLocationCode string    `bson:"FunctionalLocationCode",json:"FunctionalLocationCode"`
	Str                    string    `bson:"FLDescription",json:"Description"`
	CostCtr                string    `bson:"CostCtr",json:"CostCtr"`
	Location               string    `bson:"Location",json:"Location"`
	PIPI                   string    `bson:"FLDescription",json:"PIPI"`
	PInt                   string    `bson:"PInt",json:"PInt"`
	MainWorkCtr            string    `bson:"MainWorkCtr",json:"MainWorkCtr"`
	CatProf                string    `bson:"CatProf",json:"CatProf"`
	SortField              string    `bson:"SortField",json:"SortField"`
	ModelNo                string    `bson:"ModelNo",json:"ModelNo"`
	SerNo                  string    `bson:"SerNo",json:"SerNo"`
	UserStatus             string    `bson:"UserStatus",json:"UserStatus"`
	A                      string    `bson:"A",json:"A"`
	ObjectType             string    `bson:"ObjectType",json:"ObjectType"`
	PG                     string    `bson:"PG",json:"PG"`
	ManParNo               string    `bson:"ManParNo",json:"ManParNo"`
	Asset                  string    `bson:"Asset",json:"Asset"`
	Date                   time.Time `bson:"Date",json:"Date"`
	AcqValue               string    `bson:"AcqValue",json:"AcqValue"`
	InvNo                  string    `bson:"InvNo",json:"InvNo"`
	ConstType              string    `bson:"ConstType",json:"ConstType"`
	StartFrom              time.Time `bson:"StartFrom",json:"StartFrom"`
	CreatedOn              time.Time `bson:"CreatedOn",json:"CreatedOn"`
	SupFunctionalLocation  string    `bson:"SupFunctionalLocation",json:"SupFunctionalLocation"`
}

func (m *FunctionalLocation) TableName() string {
	return "FunctionalLocation"
}

type AnomaliesFunctionalLocation struct {
	sync.RWMutex
	orm.ModelBase          `bson:"-",json:"-"`
	FunctionalLocationCode string    `bson:"FunctionalLocationCode",json:"FunctionalLocationCode"`
	Str                    string    `bson:"FLDescription",json:"Description"`
	CostCtr                string    `bson:"CostCtr",json:"CostCtr"`
	Location               string    `bson:"Location",json:"Location"`
	PIPI                   string    `bson:"FLDescription",json:"PIPI"`
	PInt                   string    `bson:"PInt",json:"PInt"`
	MainWorkCtr            string    `bson:"MainWorkCtr",json:"MainWorkCtr"`
	CatProf                string    `bson:"CatProf",json:"CatProf"`
	SortField              string    `bson:"SortField",json:"SortField"`
	ModelNo                string    `bson:"ModelNo",json:"ModelNo"`
	SerNo                  string    `bson:"SerNo",json:"SerNo"`
	UserStatus             string    `bson:"UserStatus",json:"UserStatus"`
	A                      string    `bson:"A",json:"A"`
	ObjectType             string    `bson:"ObjectType",json:"ObjectType"`
	PG                     string    `bson:"PG",json:"PG"`
	ManParNo               string    `bson:"ManParNo",json:"ManParNo"`
	Asset                  string    `bson:"Asset",json:"Asset"`
	Date                   time.Time `bson:"Date",json:"Date"`
	AcqValue               string    `bson:"AcqValue",json:"AcqValue"`
	InvNo                  string    `bson:"InvNo",json:"InvNo"`
	ConstType              string    `bson:"ConstType",json:"ConstType"`
	StartFrom              time.Time `bson:"StartFrom",json:"StartFrom"`
	CreatedOn              time.Time `bson:"CreatedOn",json:"CreatedOn"`
	SupFunctionalLocation  string    `bson:"SupFunctionalLocation",json:"SupFunctionalLocation"`
}

func (m *AnomaliesFunctionalLocation) TableName() string {
	return "Anomalies_FunctionalLocation"
}
