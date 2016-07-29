package models

import "time"

type SPDataBrowser struct {
	Id                 int64  `bson:"Id" json:"Id"`
	PeriodYear         int    `bson:"Period.Year" json:"PeriodYear"`
	FunctionalLocation string `bson:"FunctionalLocation" json:"FunctionalLocation"`
	FLDescription      string `bson:"FLDescription" json:"FLDescription"`
	/*IsTurbine                              bool      `bson:"IsTurbine" json:"IsTurbine"`
	IsSystem                               bool      `bson:"IsSystem" json:"IsSystem"`
	TurbineParent                          string    `bson:"TurbineParent" json:"TurbineParent"`
	SystemParent                           string    `bson:"SystemParent" json:"SystemParent"`
	AssetType                              string    `bson:"AssetType" json:"AssetType"`*/
	EquipmentType            string `bson:"EquipmentType" json:"EquipmentType"`
	EquipmentTypeDescription string `bson:"EquipmentTypeDescription" json:"EquipmentTypeDescription"`
	PlantCode                string `bson:"Plant.PlantCode" json:"PlantCode"`
	/*TInfShortName                          string    `bson:"TurbineInfos.ShortName" json:"TInfShortName"`
	TInfManufacturer                       string    `bson:"TurbineInfos.Manufacturer" json:"TInfManufacturer"`
	TInfModel                              string    `bson:"TurbineInfos.Model" json:"TInfModel"`
	TInfUnitType                           string    `bson:"TurbineInfos.UnitType" json:"TInfUnitType"`
	TInfInstalledCapacity                  float64   `bson:"TurbineInfos.InstalledCapacity" json:"TInfInstalledCapacity"`
	TInfOperationalCapacity                float64   `bson:"TurbineInfos.OperationalCapacity" json:"TInfOperationalCapacity"`
	TInfPrimaryFuel                        string    `bson:"TurbineInfos.PrimaryFuel" json:"TInfPrimaryFuel"`
	TInfPrimaryFuel2                       string    `bson:"TurbineInfos.PrimaryFuel2" json:"TInfPrimaryFuel2"`
	TInfBackupFuel                         string    `bson:"TurbineInfos.BackupFuel" json:"TInfBackupFuel"`
	TInfHeatRate                           float64   `bson:"TurbineInfos.HeatRate" json:"TInfHeatRate"`
	TInfEfficiency                         float64   `bson:"TurbineInfos.Efficiency" json:"TInfEfficiency"`
	TInfCommisioningDate                   time.Time `bson:"TurbineInfos.CommisioningDate" json:"TInfCommisioningDate"`
	TInfRetirementPlan                     time.Time `bson:"TurbineInfos.RetirementPlan" json:"TInfRetirementPlan"`
	TInfInstalledMWH                       float64   `bson:"TurbineInfos.InstalledMWH" json:"TInfInstalledMWH"`
	TInfActualEnergyGeneration             float64   `bson:"TurbineInfos.ActualEnergyGeneration" json:"TInfActualEnergyGeneration"`
	TInfActualFuelConsumption_GASMMSCF     float64   `bson:"TurbineInfos.ActualFuelConsumption_GASMMSCF" json:"TInfActualFuelConsumption_GASMMSCF"`
	TInfActualFuelConsumption_CrudeBarrel  float64   `bson:"TurbineInfos.ActualFuelConsumption_CrudeBarrel" json:"TInfActualFuelConsumption_CrudeBarrel"`
	TInfActualFuelConsumption_HFOBarrel    float64   `bson:"TurbineInfos.ActualFuelConsumption_HFOBarrel" json:"TInfActualFuelConsumption_HFOBarrel"`
	TInfActualFuelConsumption_DieselBarrel float64   `bson:"TurbineInfos.ActualFuelConsumption_DieselBarrel" json:"TInfActualFuelConsumption_DieselBarrel"`
	TInfCapacityFactor                     float64   `bson:"TurbineInfos.CapacityFactor" json:"TInfCapacityFactor"`
	TInfUpdateEnergyGeneration             float64   `bson:"TurbineInfos.UpdateEnergyGeneration" json:"TInfUpdateEnergyGeneration"`
	TInfUpdateFuelConsumption              float64   `bson:"TurbineInfos.UpdateFuelConsumption" json:"TInfUpdateFuelConsumption"`*/

	PlantPlantCode string `bson:"PlantPlantCode" json:"PlantPlantCode"`
	PlantPlantName string `bson:"PlantPlantName" json:"PlantPlantName"`
	/*PlantPlantType             string  `bson:"PlantPlantType" json:"PlantPlantType"`
	PlantProvince              string  `bson:"PlantProvince" json:"PlantProvince"`
	PlantRegion                string  `bson:"PlantRegion" json:"PlantRegion"`
	PlantCity                  string  `bson:"PlantCity" json:"PlantCity"`
	PlantFuelTypes_Crude       bool    `bson:"PlantFuelTypes_Crude" json:"PlantFuelTypes_Crude"`
	PlantFuelTypes_Heavy       bool    `bson:"PlantFuelTypes_Heavy" json:"PlantFuelTypes_Heavy"`
	PlantFuelTypes_Diesel      bool    `bson:"PlantFuelTypes_Diesel" json:"PlantFuelTypes_Diesel"`
	PlantFuelTypes_Gas         bool    `bson:"PlantFuelTypes_Gas" json:"PlantFuelTypes_Gas"`
	PlantGasTurbineUnit        int     `bson:"PlantGasTurbineUnit" json:"PlantGasTurbineUnit"`
	PlantGasTurbineCapacity    float64 `bson:"PlantGasTurbineCapacity" json:"PlantGasTurbineCapacity"`
	PlantSteamUnit             int     `bson:"PlantSteamUnit" json:"PlantSteamUnit"`
	PlantSteamCapacity         float64 `bson:"PlantSteamCapacity" json:"PlantSteamCapacity"`
	PlantDieselUnit            int     `bson:"PlantDieselUnit" json:"PlantDieselUnit"`
	PlantDieselCapacity        float64 `bson:"PlantDieselCapacity" json:"PlantDieselCapacity"`
	PlantCombinedCycleUnit     int     `bson:"PlantCombinedCycleUnit" json:"PlantCombinedCycleUnit"`
	PlantCombinedCycleCapacity float64 `bson:"PlantCombinedCycleCapacity" json:"PlantCombinedCycleCapacity"`
	PlantLongitude             float64 `bson:"PlantLongitude" json:"PlantLongitude"`
	PlantLatitude              float64 `bson:"PlantLatitude" json:"PlantLatitude"`*/

	// MaintenanceUserStatus         string    `bson:"MaintenanceUserStatus" json:"MaintenanceUserStatus"`
	// MaintenanceSystemStatus       string    `bson:"MaintenanceSystemStatus" json:"MaintenanceSystemStatus"`
	// MaintenanceType string `bson:"MaintenanceType" json:"MaintenanceType"`
	// MaintenanceOrderCode          string    `bson:"MaintenanceOrderCode" json:"MaintenanceOrderCode"`
	// MaintenanceNotificationCode   string    `bson:"MaintenanceNotificationCode" json:"MaintenanceNotificationCode"`
	// MaintenanceA                  string    `bson:"MaintenanceA" json:"MaintenanceA"`
	// MaintenanceEnteredBy          string    `bson:"MaintenanceEnteredBy" json:"MaintenanceEnteredBy"`
	// MaintenanceDescription        string    `bson:"MaintenanceDescription" json:"MaintenanceDescription"`
	// MaintenancePlant              string    `bson:"MaintenancePlant" json:"MaintenancePlant"`
	MaintenanceFunctionalLocation string `bson:"MaintenanceFunctionalLocation" json:"MaintenanceFunctionalLocation"`
	// MaintenanceEquipmentCode      string    `bson:"MaintenanceEquipmentCode" json:"MaintenanceEquipmentCode"`
	// MaintenanceSortField          string    `bson:"MaintenanceSortField" json:"MaintenanceSortField"`
	// MaintenancePriorityText       string    `bson:"MaintenancePriorityText" json:"MaintenancePriorityText"`
	// MaintenanceWorkCtr            string    `bson:"MaintenanceWorkCtr" json:"MaintenanceWorkCtr"`
	// MaintenanceCostCtr            string    `bson:"MaintenanceCostCtr" json:"MaintenanceCostCtr"`
	// MaintenanceMAT                string    `bson:"MaintenanceMAT" json:"MaintenanceMAT"`
	// MaintenanceCreatedOn          time.Time `bson:"MaintenanceCreatedOn" json:"MaintenanceCreatedOn"`
	// MaintenanceBasicStart         time.Time `bson:"MaintenanceBasicStart" json:"MaintenanceBasicStart"`
	// MaintenanceBasicFinish        time.Time `bson:"MaintenanceBasicFinish" json:"MaintenanceBasicFinish"`
	// MaintenanceScheduledStart     time.Time `bson:"MaintenanceScheduledStart" json:"MaintenanceScheduledStart"`
	// MaintenanceScheduledFinish    time.Time `bson:"MaintenanceScheduledFinish" json:"MaintenanceScheduledFinish"`
	// MaintenanceActualStart        time.Time `bson:"MaintenanceActualStart" json:"MaintenanceActualStart"`
	// MaintenanceActualFinish       time.Time `bson:"MaintenanceActualFinish" json:"MaintenanceActualFinish"`
	// MaintenanceRefDate            time.Time `bson:"MaintenanceRefDate" json:"MaintenanceRefDate"`
	// MaintenanceRespCCTR           string    `bson:"MaintenanceRespCCTR" json:"MaintenanceRespCCTR"`
	// MaintenanceReleaseDate        time.Time `bson:"MaintenanceReleaseDate" json:"MaintenanceReleaseDate"`
	// MaintenanceAvailFrom          time.Time `bson:"MaintenanceAvailFrom" json:"MaintenanceAvailFrom"`
	// MaintenanceAvailTo            time.Time `bson:"MaintenanceAvailTo" json:"MaintenanceAvailTo"`
	// MaintenanceActualCost         float64   `bson:"MaintenanceActualCost" json:"MaintenanceActualCost"`
	// MaintenanceDF                 string    `bson:"MaintenanceDF" json:"MaintenanceDF"`

	WorkOrderType          string    `bson:"WorkOrderType" json:"WorkOrderType"`
	MaintenanceOrder       string    `bson:"MaintenanceOrder" json:"MaintenanceOrder"`
	MaintenanceDescription string    `bson:"MaintenanceDescription" json:"MaintenanceDescription"`
	ActivityType           string    `bson:"ActivityType" json:"ActivityType"`
	PlanStart              time.Time `bson:"PlanStart" json:"PlanStart"`
	PlanFinish             time.Time `bson:"PlanFinish" json:"PlanFinish"`
	PlanDuration           float64   `bson:"PlanDuration" json:"PlanDuration"`
	ActualStart            time.Time `bson:"ActualStart" json:"ActualStart"`
	ActualFinish           time.Time `bson:"ActualFinish" json:"ActualFinish"`
	ActualDuration         float64   `bson:"ActualDuration" json:"ActualDuration"`
	LastMaintenanceDate    time.Time `bson:"LastMaintenanceDate" json:"LastMaintenanceDate"`
	MaintenanceInterval    float64   `bson:"MaintenanceInterval" json:"MaintenanceInterval"`
	// MaintenanceIntervalTmp float64   `bson:"MaintenanceIntervalTmp" json:"MaintenanceIntervalTmp"`
	MaintenanceCost float64 `bson:"MaintenanceCost" json:"MaintenanceCost"`
}

/*func (m *SPDataBrowser) TableName() string {
	return ""
}

func (m *SPDataBrowser) SPName() string {
	return "SP_DataBrowser"
}*/
