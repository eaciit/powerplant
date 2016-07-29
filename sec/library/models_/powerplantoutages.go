package models

import (
	"sync"
	"time"

	"github.com/eaciit/orm"
)

type PowerPlantOutages struct {
	sync.RWMutex
	orm.ModelBase `bson:"-" json:"-"`
	Id            string                     `bson:"_id" json:"Id"`
	Plant         string                     `bson:"Plant" json:"Plant"`
	Year          int                        `bson:"Year" json:"Year"`
	Details       []PowerPlantOutagesDetails `bson:"Details" json:"Details"`
}

func (p *PowerPlantOutages) TableName() string {
	return "PowerPlantOutages"
}

type PowerPlantOutagesDetails struct {
	sync.RWMutex
	orm.ModelBase `bson:"-" json:"-"`
	// Id            int64     `bson:"Id" json:"Id"`
	POId          string    `bson:"POId" json:"POId"`
	UnitNo        string    `bson:"UnitNo" json:"UnitNo"`
	DepCap        string    `bson:"DepCap" json:"DepCap"`
	MFR           string    `bson:"MFR" json:"MFR"`
	Model         string    `bson:"Model" json:"Model"`
	Component     string    `bson:"Component" json:"Component"`
	OutageType    string    `bson:"OutageType" json:"OutageType"`
	OutageReason  string    `bson:"OutageReason" json:"OutageReason"`
	StartDate     time.Time `bson:"StartDate" json:"StartDate"`
	FinishDate    time.Time `bson:"FinishDate" json:"FinishDate"`
	TotalHours    float64   `bson:"TotalHours" json:"TotalHours"`
	ExpFinishDate time.Time `bson:"ExpFinishDate" json:"ExpFinishDate"`
	PlantName     string    `bson:"PlantName" json:"PlantName"`
}

func (d *PowerPlantOutagesDetails) TableName() string {
	return "PowerPlantOutagesDetails"
}
