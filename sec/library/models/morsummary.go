package models

import (
	"sync"
	"time"

	"github.com/eaciit/orm"
)

type MORSummary struct {
	sync.RWMutex
	orm.ModelBase `bson:"-" json:"-"`
	// Id                   int64     `bson:"Id" json:"Id"`
	Period               time.Time `bson:"Period" json:"Period"`
	Plant                string    `bson:"Plant" json:"Plant"`
	Province             string    `bson:"Province" json:"Province"`
	Region               string    `bson:"Region" json:"Region"`
	City                 string    `bson:"City" json:"City"`
	TopElement           string    `bson:"TopElement" json:"TopElement"`
	Element              string    `bson:"Element" json:"Element"`
	SubElement           string    `bson:"SubElement" json:"SubElement"`
	Value                float64   `bson:"Value" json:"Value"`
	NetGeneration        float64   `bson:"NetGeneration" json:"NetGeneration"`
	ServiceHours         float64   `bson:"ServiceHours" json:"ServiceHours"`
	ValueByNetGeneration float64   `bson:"ValueByNetGeneration" json:"ValueByNetGeneration"`
	ValueByServiceHours  float64   `bson:"ValueByServiceHours" json:"ValueByServiceHours"`
}

func (m *MORSummary) TableName() string {
	return "MORSummary"
}
