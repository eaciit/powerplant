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
	LoadMW        float64 `bson:"LoadMW" json:"LoadMW"`
	Mvar          float64 `bson:"Mvar" json:"Mvar"`
	StatorTemp    float64 `bson:"StatorTemp" json:"StatorTemp"`
	FieldTemp     float64 `bson:"FieldTemp" json:"FieldTemp"`
	BB1           float64 `bson:"BB1" json:"BB1"`
	BB2           float64 `bson:"BB2" json:"BB2"`
	BB3           float64 `bson:"BB3" json:"BB3"`
	BB4           float64 `bson:"BB4" json:"BB4"`
	BB5           float64 `bson:"BB5" json:"BB5"`
	MaxVib        float64 `bson:"MaxVib" json:"MaxVib"`
	Remark        string  `bson:"Remark" json:"Remark"`
	Year          int     `bson:"Year" json:"Year"`
}

func (m *Vibration) TableName() string {
	return "Vibration"
}
