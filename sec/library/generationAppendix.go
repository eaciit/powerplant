package models

import (
	"sync"

	"github.com/eaciit/orm"
)

type GenerationAppendix struct {
	sync.RWMutex
	orm.ModelBase      `bson:"-" json:"-"`
	Plant              string  `bson:"PowerPlant" json:"Plant"`
	UnitType           string  `bson:"UnitType" json:"UnitType"`
	Units              int     `bson:"Units" json:"Units"`
	ContractedCapacity float64 `bson:"ContractedCapacity" json:"ContractedCapacity"`
	CCR                float64 `bson:"CCR" json:"CCR"`
	FOMR               float64 `bson:"FOMR" json:"FOMR"`
	VOMR               float64 `bson:"VOMR" json:"VOMR"`
	Startup            float64 `bson:"Startup" json:"Startup"`
	Deduct             float64 `bson:"Deduct" json:"Deduct"`
}

func (m *GenerationAppendix) TableName() string {
	return "GenerationAppendix"
}
