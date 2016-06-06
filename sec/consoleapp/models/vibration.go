package models

import (
	"sync"

	"github.com/eaciit/orm"
)

type Vibration struct {
	sync.RWMutex
	orm.ModelBase `bson:"-" json:"-"`
	Id            string  `bson:"_id" json:"Id"`
	Plant         string  `bson:"Plant" json:"Plant"`
	WeekNo        int     `bson:"WeekNo" json:"WeekNo"`
	AmbiantTemp   float64 `bson:"AmbiantTemp" json:"AmbiantTemp"`
	AlarmLimit    float64 `bson:"AlarmLimit" json:"AlarmLimit"`
	TripLimit     float64 `bson:"TripLimit" json:"TripLimit"`
	UnitNo        string  `bson:"UnitNo" json:"UnitNo"`
	LoadMW        float64 `bson:"Load" json:"Load"`
	Mvar          float64 `bson:"Mvar" json:"Mvar"`
	StatorTemp    float64 `bson:"StatorTemp" json:"StatorTemp"`
	FieldTemp     float64 `bson:"FieldTemp" json:"FieldTemp"`
	MaxVib        float64 `bson:"MaxVib" json:"MaxVib"`
	Remark        string  `bson:"Remark" json:"Remark"`
	Year          int     `bson:"Year" json:"Year"`
	B1            []Bearing1
	B2            []Bearing2
	B3            []Bearing3
	B4            []Bearing4
	B5            []Bearing5
}

type BearingDetail struct {
	sync.RWMutex
	orm.ModelBase `bson:"-" json:"-"`
	// Id            int64   `bson:"Id" json:"Id"`
	VibrationId string  `bson:"VibrationId" json:"VibrationId"`
	Type        string  `bson:"Type" json:"Type"`
	Value       float64 `bson:"Value" json:"Value"`
	Name        string  `bson:"Name" json:"Name"`
}

func (m *BearingDetail) TableName() string {
	return "BearingDetail"
}

type Bearing1 struct {
	sync.RWMutex
	orm.ModelBase `bson:"-" json:"-"`
	// Id              int64   `bson:"Id" json:"Id"`
	Value string `bson:"Value" json:"Value"`
	Name  string `bson:"Name" json:"Name"`
}

type Bearing2 struct {
	sync.RWMutex
	orm.ModelBase `bson:"-" json:"-"`
	// Id              int64   `bson:"Id" json:"Id"`
	Value string `bson:"Value" json:"Value"`
	Name  string `bson:"Name" json:"Name"`
}

type Bearing3 struct {
	sync.RWMutex
	orm.ModelBase `bson:"-" json:"-"`
	// Id              int64   `bson:"Id" json:"Id"`
	Value string `bson:"Value" json:"Value"`
	Name  string `bson:"Name" json:"Name"`
}

type Bearing4 struct {
	sync.RWMutex
	orm.ModelBase `bson:"-" json:"-"`
	// Id              int64   `bson:"Id" json:"Id"`
	Value string `bson:"Value" json:"Value"`
	Name  string `bson:"Name" json:"Name"`
}

type Bearing5 struct {
	sync.RWMutex
	orm.ModelBase `bson:"-" json:"-"`
	// Id              int64   `bson:"Id" json:"Id"`
	Value string `bson:"Value" json:"Value"`
	Name  string `bson:"Name" json:"Name"`
}

func (m *Vibration) TableName() string {
	return "Vibration"
}
