package models

import (
	"github.com/eaciit/orm"
	"time"
)

type PowerPlantOutagesDetails struct {
	Id            string    `bson:"_id"`
	UnitNo        string    `bson:"UnitNo"`
	DepCap        string    `bson:"DepCap"`
	MFR           string    `bson:"MFR"`
	Model         string    `bson:"Model"`
	Component     string    `bson:"Component"`
	OutageType    string    `bson:"OutageType"`
	OutageReason  string    `bson:"OutageReason"`
	StartDate     time.Time `bson:"StartDate"`
	FinishDate    time.Time `bson:"FinishDate"`
	TotalHours    float64   `bson:"TotalHours"`
	ExpFinishDate time.Time `bson:"ExpFinishDate"`
	PlantName     string    `bson:"PlantName"`
}

type PowerPlantOutages struct {
	orm.ModelBase `bson:"-" json:"-"`
	Id            string                     `bson:"_id"`
	Plant         string                     `bson:"Plant"`
	Year          int                        `bson:"Year"`
	Details       []PowerPlantOutagesDetails `bson:"Details"`
}

func (p *PowerPlantOutages) TableName() string {
	return "PowerPlantOutages"
}
