package models

import (
	"sync"

	"github.com/eaciit/orm"
)

type GenerationAppendix struct {
	sync.RWMutex
	orm.ModelBase      `bson:"-",json:"-"`
	Plant              string  `bson:"Plant",json:"Plant"`
	Type               string  `bson:"Type",json:"Type"`
	Units              int     `bson:"Units",json:"Units"`
	ContractedCapacity float64 `bson:"ContractedCapacity",json:"ContractedCapacity"`
	CCR                float64 `bson:"CCR",json:"CCR"`
	FOMR               float64 `bson:"FOMR",json:"FOMR"`
	VOMR               float64 `bson:"VOMR",json:"VOMR"`
	AGP                float64 `bson:"AGP",json:"AGP"`
	LCSummer           float64 `bson:"LCSummer",json:"LCSummer"`
	LCWinter           float64 `bson:"LCWinter",json:"LCWinter"`
	LCTotal            float64 `bson:"LCTotal",json:"LCTotal"`
	Startup            float64 `bson:"Startup",json:"Startup"`
	Deduct             float64 `bson:"Deduct",json:"Deduct"`
}

func (m *GenerationAppendix) TableName() string {
	return "GenerationAppendix"
}
