package models

import (
	"sync"

	"github.com/eaciit/orm"
)

type UnitCost struct {
	sync.RWMutex
	orm.ModelBase                   `bson:"-" json:"-"`
	Id                              string  `bson:"_id" json:"id"`
	Plant                           string  `bson:"Plant" json:"Plant"`
	Turbine                         string  `bson:"Turbine" json:"Turbine"`
	MaterialCostStorehouse          float64 `bson:"MaterialCostStorehouse" json:"MaterialCostStorehouse"`
	MaterialCostDirectCharge        float64 `bson:"MaterialCostDirectCharge" json:"MaterialCostDirectCharge"`
	MaterialCostDirectPurchasesb    float64 `bson:"MaterialCostDirectPurchasesb" json:"MaterialCostDirectPurchasesb"`
	MaterialsCost                   float64 `bson:"MaterialsCost" json:"MaterialsCost"`
	ProfessionalServices            float64 `bson:"ProfessionalServices" json:"ProfessionalServices"`
	ContractMaintenance             float64 `bson:"ContractMaintenance" json:"ContractMaintenance"`
	ConstructionEquipmentRental     float64 `bson:"ConstructionEquipmentRental" json:"ConstructionEquipmentRental"`
	ContractInvoicesCost            float64 `bson:"ContractInvoicesCost" json:"ContractInvoicesCost"`
	FuelInvoiceCostLightCrudeOil    float64 `bson:"FuelInvoiceCostLightCrudeOil" json:"FuelInvoiceCostLightCrudeOil"`
	FuelInvoiceCostGas              float64 `bson:"FuelInvoiceCostGas" json:"FuelInvoiceCostGas"`
	ExpenseCost                     float64 `bson:"ExpenseCost" json:"ExpenseCost"`
	MaterialsServiceCostAllocated   float64 `bson:"MaterialsServiceCostAllocated" json:"MaterialsServiceCostAllocated"`
	HumanResourcesServicesCost      float64 `bson:"HumanResourcesServicesCost" json:"HumanResourcesServicesCost"`
	SharedServicesCost              float64 `bson:"SharedServicesCost" json:"SharedServicesCost"`
	MaterialServicesCostDifferentia float64 `bson:"MaterialServicesCostDifferentia" json:"MaterialServicesCostDifferentia"`
	SharedServicesCostDifferentia   float64 `bson:"SharedServicesCostDifferentia" json:"SharedServicesCostDifferentia"`
	TrainingDevelopmentCostAllocat  float64 `bson:"TrainingDevelopmentCostAllocat" json:"TrainingDevelopmentCostAllocat"`
	HumanResourcesCost              float64 `bson:"HumanResourcesCost" json:"HumanResourcesCost"`
	PowerPlantManagementCost        float64 `bson:"PowerPlantManagementCost" json:"PowerPlantManagementCost"`
	PowerPlantOperationalCost       float64 `bson:"PowerPlantOperationalCost" json:"PowerPlantOperationalCost"`
	PowerPlantMaintenanceCost       float64 `bson:"PowerPlantMaintenanceCost" json:"PowerPlantMaintenanceCost"`
	PowerPlantTechnicalSupport      float64 `bson:"PowerPlantTechnicalSupport" json:"PowerPlantTechnicalSupport"`
	PowerPlantCommonFacilitiest     float64 `bson:"PowerPlantCommonFacilitiest" json:"PowerPlantCommonFacilitiest"`
	BusinessUnitOverHeadCostsAlloca float64 `bson:"BusinessUnitOverHeadCostsAlloca" json:"BusinessUnitOverHeadCostsAlloca"`
	SecondaryLaborCostCharedSaudi   float64 `bson:"SecondaryLaborCostCharedSaudi" json:"SecondaryLaborCostCharedSaudi"`
	SecondaryLaborCostChargedNonSau float64 `bson:"SecondaryLaborCostChargedNonSau" json:"SecondaryLaborCostChargedNonSau"`
	SecondaryLaborCostAllocated     float64 `bson:"SecondaryLaborCostAllocated" json:"SecondaryLaborCostAllocated"`
	SecondaryLaborOtCostAllocated   float64 `bson:"SecondaryLaborOtCostAllocated" json:"SecondaryLaborOtCostAllocated"`
	SecondaryOtCostChargedSaudi     float64 `bson:"SecondaryOtCostChargedSaudi" json:"SecondaryOtCostChargedSaudi"`
	SecondaryAllocated              float64 `bson:"SecondaryAllocated" json:"SecondaryAllocated"`
	TotalCost                       float64 `bson:"TotalCost" json:"TotalCost"`
}

func (m *UnitCost) TableName() string {
	return "UnitCost"
}
