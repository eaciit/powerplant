package models

import (
	"sync"
	"time"

	"github.com/eaciit/orm"
)

type PowerPlantOutages struct {
	sync.RWMutex
	orm.ModelBase      `bson:"-" json:"-"`
	Id                 string    `bson:"Id" json:"Id"`
	Plant              string    `bson:"Plant" json:"Plant"`
	Unit               string    `bson:"Unit" json:"Unit"`
	FunctionalLocation string    `bson:"FunctionalLocation" json:"FunctionalLocation"`
	Component          string    `bson:"Component" json:"Component"`
	OutageType         string    `bson:"OutageType" json:"OutageType"`
	StartDate          time.Time `bson:"StartDate" json:"StartDate"`
	FinishDate         time.Time `bson:"FinishDate" json:"FinishDate"`
	TotalHours         float64   `bson:"TotalHours" json:"TotalHours"`
}

func (p *PowerPlantOutages) TableName() string {
	return "PowerPlantOutages"
}
