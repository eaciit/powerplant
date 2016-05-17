package models

import (
	"time"

	"github.com/eaciit/orm"
)


type MaintenancePlan struct {
	orm.ModelBase         `bson:"-",json:"-"`
	MaintenanceItem       string    `bson:"MaintenanceItem",json:"MaintenanceItem"`
	MaintenancePlanCode   string    `bson:"MaintenancePlanCode",json:"MaintenancePlanCode"`
	FunctionalLocation    string    `bson:"FunctionalLocation",json:"FunctionalLocation"`
	PWC_TLO               string    `bson:"PWC_TLO",json:"PWC_TLO"`
	PG_PLN                string    `bson:"PG_PLN",json:"PG_PLN"`
	TaskListType          string    `bson:"TaskListType",json:"TaskListType"`
	TaskListGroup         string    `bson:"TaskListGroup",json:"TaskListGroup"`
	TaskListCounter       string    `bson:"TaskListCounter",json:"TaskListCounter"`
	TLCounterDescription  string    `bson:"TLCounterDescription",json:"TLCounterDescription"`
	Next1stScheduledDate  time.Time `bson:"Next1stScheduledDate",json:"Next1stScheduledDate"`
	Next2ndScheduledDate  time.Time `bson:"Next2ndScheduledDate",json:"Next2ndScheduledDate"`
	Next3rdScheduledDate  time.Time `bson:"Next3rdScheduledDate",json:"Next3rdScheduledDate"`
	Next4thScheduledDate  time.Time `bson:"Next4thScheduledDate",json:"Next4thScheduledDate"`
	Equipment             string    `bson:"Equipment",json:"Equipment"`
	EquipmentDescription  string    `bson:"EquipmentDescription",json:"EquipmentDescription"`
	EquipmentType         string    `bson:"EquipmentType",json:"EquipmentType"`
	EquipmentStatus       string    `bson:"EquipmentStatus",json:"EquipmentStatus"`
	EquipmentStatusText   string    `bson:"EquipmentStatusText",json:"EquipmentStatusText"`
	Class                 string    `bson:"Class",json:"Class"`
	Building              string    `bson:"Building",json:"Building"`
	BuildingText          string    `bson:"BuildingText",json:"BuildingText"`
	PlanningPlant         string    `bson:"PlanningPlant",json:"PlanningPlant"`
	MaintenancePlant      string    `bson:"MaintenancePlant",json:"MaintenancePlant"`
	Location              string    `bson:"Location",json:"Location"`
	PlantSection          string    `bson:"PlantSection",json:"PlantSection"`
	CompanyCode           string    `bson:"CompanyCode",json:"CompanyCode"`
	CostCtr               string    `bson:"CostCtr",json:"CostCtr"`
	MWC_PLN               string    `bson:"MWC_PLN",json:"MWC_PLN"`
	MWC_PLNDescription    string    `bson:"MWC_PLNDescription",json:"MWC_PLNDescription"`
	MWC_TLH               string    `bson:"MWC_TLH",json:"MWC_TLH"`
	MWC_TLHDescription    string    `bson:"MWC_TLHDescription",json:"MWC_TLHDescription"`
	MWC_TLODescription    string    `bson:"MWC_TLODescription",json:"MWC_TLODescription"`
	MWC_EQ                string    `bson:"MWC_EQ",json:"MWC_EQ"`
	MWC_FL                string    `bson:"MWC_FL",json:"MWC_FL"`
	PG_PLNDescription     string    `bson:"PG_PLNDescription",json:"PG_PLNDescription"`
	PG_FL                 string    `bson:"PG_FL",json:"PG_FL"`
	PG_EQ                 string    `bson:"PG_EQ",json:"PG_EQ"`
	StdTextKey            string    `bson:"StdTextKey",json:"StdTextKey"`
	StdTextKeyDescription string    `bson:"StdTextKeyDescription",json:"StdTextKeyDescription"`
	OperationShortText    string    `bson:"OperationShortText",json:"OperationShortText"`
	DurationOfActivity    int       `bson:"DurationOfActivity",json:"DurationOfActivity"`
	UoMDuration           string    `bson:"UoMDuration",json:"UoMDuration"`
	TotalWorkInActivity   int       `bson:"TotalWorkInActivity",json:"TotalWorkInActivity"`
	UoMWork               string    `bson:"UoMWork",json:"UoMWork"`
	NoOfPerson            int       `bson:"NoOfPerson",json:"NoOfPerson"`
	ActivityType          string    `bson:"ActivityType",json:"ActivityType"`
	MaintenanceStrategy   string    `bson:"MaintenanceStrategy",json:"MaintenanceStrategy"`
	Package1              float64   `bson:"Package1",json:"Package1"`
	Package1Text          string    `bson:"Package1Text",json:"Package1Text"`
	Package2              float64   `bson:"Package2",json:"Package2"`
	Package2Text          string    `bson:"Package2Text",json:"Package2Text"`
	Package3              float64   `bson:"Package3",json:"Package3"`
	Package3Text          string    `bson:"Package3Text",json:"Package3Text"`
	Package4              float64   `bson:"Package4",json:"Package4"`
	Package4Text          string    `bson:"Package4Text",json:"Package4Text"`
	Package5              float64   `bson:"Package5",json:"Package5"`
	Package5Text          string    `bson:"Package5Text",json:"Package5Text"`
	Package6              float64   `bson:"Package6",json:"Package6"`
	Package6Text          string    `bson:"Package6Text",json:"Package6Text"`
	Package7              float64   `bson:"Package7",json:"Package7"`
	Package7Text          string    `bson:"Package7Text",json:"Package7Text"`
	PlanCreatedBy         string    `bson:"PlanCreatedBy",json:"PlanCreatedBy"`
	PlanCreatedOn         time.Time `bson:"PlanCreatedOn",json:"PlanCreatedOn"`
	PlanChangedBy         string    `bson:"PlanChangedBy",json:"PlanChangedBy"`
	PlanChangedOn         time.Time `bson:"PlanChangedOn",json:"PlanChangedOn"`
	Manufacturer          string    `bson:"Manufacturer",json:"Manufacturer"`
	ManufacturerSerialNo  string    `bson:"ManufacturerSerialNo",json:"ManufacturerSerialNo"`
	ManufacturerModelNo   string    `bson:"ManufacturerModelNo",json:"ManufacturerModelNo"`
	ManufacturerPartNo    string    `bson:"ManufacturerPartNo",json:"ManufacturerPartNo"`
}

func (m *MaintenancePlan) TableName() string {
	return "MaintenancePlan"
}
