package models

import (
	"sync"

	"github.com/eaciit/orm"
)

type PerformanceFactors struct {
	sync.RWMutex
	orm.ModelBase `bson:"-",json:"-"`
	id            int64   `bson:"id",json:"id"`
	Plant         string  `bson:"Plant",json:"Plant"`
	Year          int     `bson:"Year",json:"Year"`
	Unit          string  `bson:"Unit",json:"Unit"`
	GSHR          float64 `bson:"GSHR",json:"GSHR"`
	NSHR          float64 `bson:"NSHR",json:"NSHR"`
	GTEF          float64 `bson:"GTEF",json:"GTEF"`
	NTEF          float64 `bson:"NTEF",json:"NTEF"`
	GCF           float64 `bson:"GCF",json:"GCF"`
	NCF           float64 `bson:"NCF",json:"NCF"`
	SRF           float64 `bson:"SRF",json:"SRF"`
	ORF           float64 `bson:"ORF",json:"ORF"`
	ART           float64 `bson:"ART",json:"ART"`
	EAF           float64 `bson:"EAF",json:"EAF"`
	EFOF          float64 `bson:"EFOF",json:"EFOF"`
	EUOF          float64 `bson:"EUOF",json:"EUOF"`
}

func (m *PerformanceFactors) PerformanceFactors() string {
	return "PerformanceFactors"
}
