package models

import (
	"sync"
	"time"

	"github.com/eaciit/orm"
)

type PlannedMaintenance struct {
	sync.RWMutex
	orm.ModelBase `bson:"-",json:"-"`
	// Id                    int64     `bson:"id",json:"id"`
	Equipment             string    `bson:"Equipment",json:"Equipment"`
	MaintenancePlan       string    `bson:"MaintenancePlan",json:"MaintenancePlan"`
	MainWorkCtr           string    `bson:"MainWorkCtr",json:"MainWorkCtr"`
	MaintenanceItem       string    `bson:"MaintenanceItem",json:"MaintenanceItem"`
	CostCtr               string    `bson:"CostCtr",json:"CostCtr"`
	FunctionalLocation    string    `bson:"FunctionalLocation",json:"FunctionalLocation"`
	Description           string    `bson:"Description",json:"Description"`
	MaintenanceItemText   string    `bson:"MaintenanceItemText",json:"MaintenanceItemText"`
	CreatedBy             string    `bson:"CreatedBy",json:"CreatedBy"`
	CreatedOn             time.Time `bson:"CreatedOn",json:"CreatedOn"`
	Priority              string    `bson:"Priority",json:"Priority"`
	ABCIndicator          string    `bson:"ABCIndicator",json:"ABCIndicator"`
	ChangedOn             time.Time `bson:"ChangedOn",json:"ChangedOn"`
	ChangedBy             string    `bson:"ChangedBy",json:"ChangedBy"`
	Asset                 string    `bson:"Asset",json:"Asset"`
	SubNumber             string    `bson:"SubNumber",json:"SubNumber"`
	WorkCtr               string    `bson:"WorkCtr",json:"WorkCtr"`
	OrderType             string    `bson:"OrderType",json:"OrderType"`
	OrderNo               string    `bson:"OrderNo",json:"OrderNo"`
	Assembl               string    `bson:"Assembl",json:"Assembl"`
	PlantSection          string    `bson:"PlantSection",json:"PlantSection"`
	PurchaseOrder         string    `bson:"PurchaseOrder",json:"PurchaseOrder"`
	Item1                 string    `bson:"Item1",json:"Item1"`
	CompanyCode           string    `bson:"CompanyCode",json:"CompanyCode"`
	CycleSetSequence      string    `bson:"CycleSetSequence",json:"CycleSetSequence"`
	StandingOrder         string    `bson:"StandingOrder",json:"StandingOrder"`
	ChangeAuthentication  string    `bson:"ChangeAuthentication",json:"ChangeAuthentication"`
	LogicalSystem1        string    `bson:"LogicalSystem1",json:"LogicalSystem1"`
	LogicalSystem2        string    `bson:"LogicalSystem2",json:"LogicalSystem2"`
	SortField             string    `bson:"SortField",json:"SortField"`
	BusinessArea          string    `bson:"BusinessArea",json:"BusinessArea"`
	SettlementOrder       string    `bson:"SettlementOrder",json:"SettlementOrder"`
	MaintenanceActiveType string    `bson:"MaintenanceActiveType",json:"MaintenanceActiveType"`
	ILOAIndividual        string    `bson:"ILOAIndividual",json:"ILOAIndividual"`
	LOCACCAssmt           string    `bson:"LOC_ACCAssmt",json:"LOC_ACCAssmt"`
	SettlementRule        string    `bson:"SettlementRule",json:"SettlementRule"`
	PlanningPlant         string    `bson:"PlanningPlant",json:"PlanningPlant"`
	SalesDocument         string    `bson:"SalesDocument",json:"SalesDocument"`
	Item2                 string    `bson:"Item2",json:"Item2"`
	COArea                string    `bson:"COArea",json:"COArea"`
	Language              string    `bson:"Language",json:"Language"`
	LastOrder             string    `bson:"LastOrder",json:"LastOrder"`
	LTextIndicator        string    `bson:"LTextIndicator",json:"LTextIndicator"`
	Client                string    `bson:"Client",json:"Client"`
	MPlanCategory         string    `bson:"MPlanCategory",json:"MPlanCategory"`
	Room                  string    `bson:"Room",json:"Room"`
	ObjectList            int       `bson:"ObjectList",json:"ObjectList"`
	GroupCounter          string    `bson:"GroupCounter",json:"GroupCounter"`
	Grup                  string    `bson:"Grup",json:"Grup"`
	TaskListType          string    `bson:"TaskListType",json:"TaskListType"`
	FLDescription         string    `bson:"FLDescription",json:"FLDescription"`
	WBSElement            string    `bson:"WBSElement",json:"WBSElement"`
	NotificationType      string    `bson:"NotificationType",json:"NotificationType"`
	SerialNumber          string    `bson:"SerialNumber",json:"SerialNumber"`
	Material              string    `bson:"Material",json:"Material"`
	Division              string    `bson:"Division",json:"Division"`
	Status                string    `bson:"Status",json:"Status"`
	Location              string    `bson:"Location",json:"Location"`
	MaintenancePlant      string    `bson:"MaintenancePlant",json:"MaintenancePlant"`
	SalesOrg              string    `bson:"SalesOrg",json:"SalesOrg"`
	DistributionChannel   string    `bson:"DistributionChannel",json:"DistributionChannel"`
	PlannerGroup          string    `bson:"PlannerGroup",json:"PlannerGroup"`
	ItemNumber            string    `bson:"ItemNumber",json:"ItemNumber"`
	Strategy              string    `bson:"Strategy",json:"Strategy"`
}

func (m *PlannedMaintenance) TableName() string {
	return "PlannedMaintenance"
}
