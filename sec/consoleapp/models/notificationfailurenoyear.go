package models

import (
	"sync"
	"time"

	"github.com/eaciit/orm"
)

type NotificationFailureNoYear struct {
	sync.RWMutex
	orm.ModelBase `bson:"-" json:"-"`
	// Id	int64	`bson:"id" json:"id"`
	Plant                string    `bson:"Plant" json:"Plant"`
	ObjectType           string    `bson:"ObjectType" json:"ObjectType"`
	ObjectTypeDesc       string    `bson:"ObjectTypeDesc" json:"ObjectTypeDesc"`
	Equipment            string    `bson:"Equipment" json:"Equipment"`
	EquipmentDesc        string    `bson:"EquipmentDesc" json:"EquipmentDesc"`
	MaintenanceOrder     string    `bson:"MaintenanceOrder" json:"MaintenanceOrder"`
	MaintenanceOrderDesc string    `bson:"MaintenanceOrderDesc" json:"MaintenanceOrderDesc"`
	WorkCenter           string    `bson:"WorkCenter" json:"WorkCenter"`
	WorkCenterDesc       string    `bson:"WorkCenterDesc" json:"WorkCenterDesc"`
	PMUserStatus         string    `bson:"PMUserStatus" json:"PMUserStatus"`
	Notification         string    `bson:"Notification" json:"Notification"`
	NotifComplDate       time.Time `bson:"NotifComplDate" json:"NotifComplDate"`
	SystemStatus         string    `bson:"SystemStatus" json:"SystemStatus"`
	SystemStatusDesc     string    `bson:"SystemStatusDesc" json:"SystemStatusDesc"`
	ObjectPart           string    `bson:"ObjectPart" json:"ObjectPart"`
	ObjectPartDesc       string    `bson:"ObjectPartDesc" json:"ObjectPartDesc"`
	Damage               string    `bson:"Damage" json:"Damage"`
	DamageDesc           string    `bson:"DamageDesc" json:"DamageDesc"`
	Cause                string    `bson:"Cause" json:"Cause"`
	CauseDesc            string    `bson:"CauseDesc" json:"CauseDesc"`
	Task                 string    `bson:"Task" json:"Task"`
	TaskDesc             string    `bson:"TaskDesc" json:"TaskDesc"`
	Activity             string    `bson:"Activity" json:"Activity"`
	ActivityDesc         string    `bson:"ActivityDesc" json:"ActivityDesc"`
	FailureCode          string    `bson:"FailureCode" json:"FailureCode"`
}

func (m *NotificationFailureNoYear) TableName() string {
	return "NotificationFailureNoYear"
}
