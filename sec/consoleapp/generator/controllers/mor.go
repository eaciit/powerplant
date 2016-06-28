package controllers

import (
	"strings"
	"time"

	"github.com/eaciit/crowd"
	"github.com/eaciit/dbox"
	. "github.com/eaciit/powerplant/sec/consoleapp/generator/helpers"
	. "github.com/eaciit/powerplant/sec/library/models"
	tk "github.com/eaciit/toolkit"
)

// GenMOR
type GenMOR struct {
	*BaseController
}

// Generate
func (m *GenMOR) Generate(base *BaseController) {
	if base != nil {
		m.BaseController = base
	}

	var e error
	tk.Println("##Generating MOR..")

	e = m.generateMORSummary()
	if e != nil {
		tk.Println(e)
	}

	e = m.generateMORFlatCalculationSummary()
	if e != nil {
		tk.Println(e)
	}

	tk.Println("##MOR Data : DONE\n")
}

// generateMORSummary
func (m *GenMOR) generateMORSummary() error {
	tk.Println("Generating MOR Summary..")
	ctx := m.BaseController.Ctx
	c := ctx.Connection

	PowerPlantInfos := []PowerPlantInfo{}
	csr, e := c.NewQuery().From(new(PowerPlantInfo).TableName()).Cursor(nil)

	if e != nil {
		return e
	}
	e = csr.Fetch(&PowerPlantInfos, 0, false)
	if e != nil {
		return e
	}
	csr.Close()

	OperationalDatas := []OperationalData{}
	csr, e = c.NewQuery().From(new(OperationalData).TableName()).Cursor(nil)
	if e != nil {
		return e
	}
	e = csr.Fetch(&OperationalDatas, 0, false)
	if e != nil {
		return e
	}
	csr.Close()

	MaintenanceCostFLList := []MaintenanceCostFL{}
	csr, e = c.NewQuery().From(new(MaintenanceCostFL).TableName()).Take(5).Cursor(nil)
	if e != nil {
		return e
	}
	e = csr.Fetch(&MaintenanceCostFLList, 0, false)
	if e != nil {
		return e
	}
	csr.Close()
	// Maintenance
	for _, cost := range MaintenanceCostFLList {

		Infos := crowd.From(&PowerPlantInfos).Where(func(x interface{}) interface{} {
			return strings.Contains(strings.ToLower(x.(PowerPlantInfo).Name), strings.ToLower(PlantNormalization(cost.Plant)))
		}).Exec().Result.Data().([]PowerPlantInfo)
		Province := ""
		Region := ""
		City := ""
		if Infos != nil && len(Infos) > 0 {
			info := Infos[0]
			Province = info.Province
			Region = info.Region
			City = info.City
		}

		OpDatas := crowd.From(&OperationalDatas).Where(func(x interface{}) interface{} {
			op := x.(OperationalData)
			return op.Year == cost.Period.Year() &&
				strings.Contains(strings.ToLower(x.(OperationalData).Plant), strings.ToLower(PlantNormalization(cost.Plant)))
		})

		NetGen := OpDatas.Sum(func(x interface{}) interface{} {
			return x.(OperationalData).GenerationNet
		}).Exec().Result.Sum

		ServiceHours := OpDatas.Sum(func(x interface{}) interface{} {
			return x.(OperationalData).ServiceHours
		}).Exec().Result.Sum

		// Internal Labor
		d := new(MORSummary)
		d.Period = time.Date(cost.Period.Year(), cost.Period.Month(), 1, 0, 0, 0, 0, time.UTC)
		d.Plant = PlantNormalization(cost.Plant)
		d.TopElement = "Maintenance"
		d.Element = "Internal Labor"
		d.SubElement = "Internal Labor"
		d.Value = cost.InternalLaborActual
		d.Province = Province
		d.Region = Region
		d.City = City
		d.NetGeneration = NetGen / 12
		d.ServiceHours = ServiceHours / 12
		if d.NetGeneration != 0 {
			d.ValueByNetGeneration = d.Value / d.NetGeneration
		}
		if d.ServiceHours != 0 {
			d.ValueByServiceHours = d.Value / d.ServiceHours
		}
		_, e := ctx.InsertOut(d)
		if e != nil {
			tk.Println(e)
			break
		}

		// Internal Material
		d = new(MORSummary)
		d.Period = time.Date(cost.Period.Year(), cost.Period.Month(), 1, 0, 0, 0, 0, time.UTC)
		d.Plant = PlantNormalization(cost.Plant)
		d.TopElement = "Maintenance"
		d.Element = "Internal Material"
		d.SubElement = "Internal Material"
		d.Value = cost.InternalMaterialActual
		d.Province = Province
		d.Region = Region
		d.City = City
		d.NetGeneration = NetGen / 12
		d.ServiceHours = ServiceHours / 12
		if d.NetGeneration != 0 {
			d.ValueByNetGeneration = d.Value / d.NetGeneration
		}
		if d.ServiceHours != 0 {
			d.ValueByServiceHours = d.Value / d.ServiceHours
		}
		_, e = ctx.InsertOut(d)
		if e != nil {
			tk.Println(e)
			break
		}

		// Direct Material
		d = new(MORSummary)
		d.Period = time.Date(cost.Period.Year(), cost.Period.Month(), 1, 0, 0, 0, 0, time.UTC)
		d.Plant = PlantNormalization(cost.Plant)
		d.TopElement = "Maintenance"
		d.Element = "Direct Material"
		d.SubElement = "Direct Material"
		d.Value = cost.DirectMaterialActual
		d.Province = Province
		d.Region = Region
		d.City = City
		d.NetGeneration = NetGen / 12
		d.ServiceHours = ServiceHours / 12
		if d.NetGeneration != 0 {
			d.ValueByNetGeneration = d.Value / d.NetGeneration
		}
		if d.ServiceHours != 0 {
			d.ValueByServiceHours = d.Value / d.ServiceHours
		}
		_, e = ctx.InsertOut(d)
		if e != nil {
			tk.Println(e)
			break
		}

		// External Service
		d = new(MORSummary)
		d.Period = time.Date(cost.Period.Year(), cost.Period.Month(), 1, 0, 0, 0, 0, time.UTC)
		d.Plant = PlantNormalization(cost.Plant)
		d.TopElement = "Maintenance"
		d.Element = "External Service"
		d.SubElement = "External Service"
		d.Value = cost.ExternalServiceActual
		d.Province = Province
		d.Region = Region
		d.City = City
		d.NetGeneration = NetGen / 12
		d.ServiceHours = ServiceHours / 12
		if d.NetGeneration != 0 {
			d.ValueByNetGeneration = d.Value / d.NetGeneration
		}
		if d.ServiceHours != 0 {
			d.ValueByServiceHours = d.Value / d.ServiceHours
		}
		_, e = ctx.InsertOut(d)
		if e != nil {
			tk.Println(e)
			break
		}

	}

	// Repair
	for _, cost := range MaintenanceCostFLList {

		Infos := crowd.From(&PowerPlantInfos).Where(func(x interface{}) interface{} {
			return strings.Contains(strings.ToLower(x.(PowerPlantInfo).Name), strings.ToLower(PlantNormalization(cost.Plant)))
		}).Exec().Result.Data().([]PowerPlantInfo)
		Province := ""
		Region := ""
		City := ""
		if Infos != nil && len(Infos) > 0 {
			info := Infos[0]
			Province = info.Province
			Region = info.Region
			City = info.City
		}

		OpDatas := crowd.From(&OperationalDatas).Where(func(x interface{}) interface{} {
			op := x.(OperationalData)
			return op.Year == cost.Period.Year() &&
				strings.Contains(strings.ToLower(x.(OperationalData).Plant), strings.ToLower(PlantNormalization(cost.Plant)))
		})

		NetGen := OpDatas.Sum(func(x interface{}) interface{} {
			return x.(OperationalData).GenerationNet
		}).Exec().Result.Sum

		ServiceHours := OpDatas.Sum(func(x interface{}) interface{} {
			return x.(OperationalData).ServiceHours
		}).Exec().Result.Sum

		// Internal Labor
		d := new(MORSummary)
		d.Period = time.Date(cost.Period.Year(), cost.Period.Month(), 1, 0, 0, 0, 0, time.UTC)
		d.Plant = PlantNormalization(cost.Plant)
		d.TopElement = "Repair"
		d.Element = "Internal Labor"
		d.SubElement = "Internal Labor"
		d.Value = cost.InternalLaborActual
		d.Province = Province
		d.Region = Region
		d.City = City
		d.NetGeneration = NetGen / 12
		d.ServiceHours = ServiceHours / 12
		if d.NetGeneration != 0 {
			d.ValueByNetGeneration = d.Value / d.NetGeneration
		}
		if d.ServiceHours != 0 {
			d.ValueByServiceHours = d.Value / d.ServiceHours
		}
		_, e := ctx.InsertOut(d)
		if e != nil {
			tk.Println(e)
			break
		}

		// Internal Material
		d = new(MORSummary)
		d.Period = time.Date(cost.Period.Year(), cost.Period.Month(), 1, 0, 0, 0, 0, time.UTC)
		d.Plant = PlantNormalization(cost.Plant)
		d.TopElement = "Repair"
		d.Element = "Internal Material"
		d.SubElement = "Internal Material"
		d.Value = cost.InternalMaterialActual
		d.Province = Province
		d.Region = Region
		d.City = City
		d.NetGeneration = NetGen / 12
		d.ServiceHours = ServiceHours / 12
		if d.NetGeneration != 0 {
			d.ValueByNetGeneration = d.Value / d.NetGeneration
		}
		if d.ServiceHours != 0 {
			d.ValueByServiceHours = d.Value / d.ServiceHours
		}
		_, e = ctx.InsertOut(d)
		if e != nil {
			tk.Println(e)
			break
		}

		// Direct Material
		d = new(MORSummary)
		d.Period = time.Date(cost.Period.Year(), cost.Period.Month(), 1, 0, 0, 0, 0, time.UTC)
		d.Plant = PlantNormalization(cost.Plant)
		d.TopElement = "Repair"
		d.Element = "Direct Material"
		d.SubElement = "Direct Material"
		d.Value = cost.DirectMaterialActual
		d.Province = Province
		d.Region = Region
		d.City = City
		d.NetGeneration = NetGen / 12
		d.ServiceHours = ServiceHours / 12
		if d.NetGeneration != 0 {
			d.ValueByNetGeneration = d.Value / d.NetGeneration
		}
		if d.ServiceHours != 0 {
			d.ValueByServiceHours = d.Value / d.ServiceHours
		}
		_, e = ctx.InsertOut(d)
		if e != nil {
			tk.Println(e)
			break
		}

		// External Service
		d = new(MORSummary)
		d.Period = time.Date(cost.Period.Year(), cost.Period.Month(), 1, 0, 0, 0, 0, time.UTC)
		d.Plant = PlantNormalization(cost.Plant)
		d.TopElement = "Repair"
		d.Element = "External Service"
		d.SubElement = "External Service"
		d.Value = cost.ExternalServiceActual
		d.Province = Province
		d.Region = Region
		d.City = City
		d.NetGeneration = NetGen / 12
		d.ServiceHours = ServiceHours / 12
		if d.NetGeneration != 0 {
			d.ValueByNetGeneration = d.Value / d.NetGeneration
		}
		if d.ServiceHours != 0 {
			d.ValueByServiceHours = d.Value / d.ServiceHours
		}
		_, e = ctx.InsertOut(d)
		if e != nil {
			tk.Println(e)
			break
		}

	}

	// Operation
	UnitCostList := []UnitCost{}
	csr, e = c.NewQuery().From(new(UnitCost).TableName()).Take(5).Cursor(nil)
	if e != nil {
		return e
	}
	e = csr.Fetch(&UnitCostList, 0, false)
	csr.Close()
	if e != nil {
		return e
	}
	Year := 2014
	for _, cost := range UnitCostList {

		Infos := crowd.From(&PowerPlantInfos).Where(func(x interface{}) interface{} {
			return strings.Contains(strings.ToLower(x.(PowerPlantInfo).Name), strings.ToLower(PlantNormalization(cost.Plant)))
		}).Exec().Result.Data().([]PowerPlantInfo)
		Province := ""
		Region := ""
		City := ""
		if Infos != nil && len(Infos) > 0 {
			info := Infos[0]
			Province = info.Province
			Region = info.Region
			City = info.City
		}

		OpDatas := crowd.From(&OperationalDatas).Where(func(x interface{}) interface{} {
			op := x.(OperationalData)
			return op.Year == Year &&
				strings.Contains(strings.ToLower(x.(OperationalData).Plant), strings.ToLower(PlantNormalization(cost.Plant)))
		})

		NetGen := OpDatas.Sum(func(x interface{}) interface{} {
			return x.(OperationalData).GenerationNet
		}).Exec().Result.Sum

		ServiceHours := OpDatas.Sum(func(x interface{}) interface{} {
			return x.(OperationalData).ServiceHours
		}).Exec().Result.Sum
		var Month time.Month = 0
		for Month < 12 {
			Month += 1
			d := new(MORSummary)
			d.Period = time.Date(Year, Month, 1, 0, 0, 0, 0, time.UTC)
			d.Plant = PlantNormalization(cost.Plant)
			d.TopElement = "Operation"
			d.Element = "Material"
			d.SubElement = "Storehouse"
			d.Value = cost.MaterialCostStorehouse / 12
			d.Province = Province
			d.Region = Region
			d.City = City
			d.NetGeneration = NetGen / 12
			d.ServiceHours = ServiceHours / 12
			if d.NetGeneration != 0 {
				d.ValueByNetGeneration = d.Value / d.NetGeneration
			}
			if d.ServiceHours != 0 {
				d.ValueByServiceHours = d.Value / d.ServiceHours
			}
			_, e := ctx.InsertOut(d)
			if e != nil {
				tk.Println(e)
				break
			}

			d = new(MORSummary)
			d.Period = time.Date(Year, Month, 1, 0, 0, 0, 0, time.UTC)
			d.Plant = PlantNormalization(cost.Plant)
			d.TopElement = "Operation"
			d.Element = "Material"
			d.SubElement = "Direct Charge"
			d.Value = cost.MaterialCostDirectCharge / 12
			d.Province = Province
			d.Region = Region
			d.City = City
			d.NetGeneration = NetGen / 12
			d.ServiceHours = ServiceHours / 12
			if d.NetGeneration != 0 {
				d.ValueByNetGeneration = d.Value / d.NetGeneration
			}
			if d.ServiceHours != 0 {
				d.ValueByServiceHours = d.Value / d.ServiceHours
			}
			_, e = ctx.InsertOut(d)
			if e != nil {
				tk.Println(e)
				break
			}

			d = new(MORSummary)
			d.Period = time.Date(Year, Month, 1, 0, 0, 0, 0, time.UTC)
			d.Plant = PlantNormalization(cost.Plant)
			d.TopElement = "Operation"
			d.Element = "Material"
			d.SubElement = "Direct Purchase b"
			d.Value = cost.MaterialCostDirectPurchasesb / 12
			d.Province = Province
			d.Region = Region
			d.City = City
			d.NetGeneration = NetGen / 12
			d.ServiceHours = ServiceHours / 12
			if d.NetGeneration != 0 {
				d.ValueByNetGeneration = d.Value / d.NetGeneration
			}
			if d.ServiceHours != 0 {
				d.ValueByServiceHours = d.Value / d.ServiceHours
			}
			_, e = ctx.InsertOut(d)
			if e != nil {
				tk.Println(e)
				break
			}

			d = new(MORSummary)
			d.Period = time.Date(Year, Month, 1, 0, 0, 0, 0, time.UTC)
			d.Plant = PlantNormalization(cost.Plant)
			d.TopElement = "Operation"
			d.Element = "Contract Invoices"
			d.SubElement = "Professional Service"
			d.Value = cost.ProfessionalServices / 12
			d.Province = Province
			d.Region = Region
			d.City = City
			d.NetGeneration = NetGen / 12
			d.ServiceHours = ServiceHours / 12
			if d.NetGeneration != 0 {
				d.ValueByNetGeneration = d.Value / d.NetGeneration
			}
			if d.ServiceHours != 0 {
				d.ValueByServiceHours = d.Value / d.ServiceHours
			}
			_, e = ctx.InsertOut(d)
			if e != nil {
				tk.Println(e)
				break
			}

			d = new(MORSummary)
			d.Period = time.Date(Year, Month, 1, 0, 0, 0, 0, time.UTC)
			d.Plant = PlantNormalization(cost.Plant)
			d.TopElement = "Operation"
			d.Element = "Contract Invoices"
			d.SubElement = "Contract Maintenance"
			d.Value = cost.ContractMaintenance / 12
			d.Province = Province
			d.Region = Region
			d.City = City
			d.NetGeneration = NetGen / 12
			d.ServiceHours = ServiceHours / 12
			if d.NetGeneration != 0 {
				d.ValueByNetGeneration = d.Value / d.NetGeneration
			}
			if d.ServiceHours != 0 {
				d.ValueByServiceHours = d.Value / d.ServiceHours
			}
			_, e = ctx.InsertOut(d)
			if e != nil {
				tk.Println(e)
				break
			}

			d = new(MORSummary)
			d.Period = time.Date(Year, Month, 1, 0, 0, 0, 0, time.UTC)
			d.Plant = PlantNormalization(cost.Plant)
			d.TopElement = "Operation"
			d.Element = "Contract Invoices"
			d.SubElement = "Construction Equipment - Rental"
			d.Value = cost.ConstructionEquipmentRental / 12
			d.Province = Province
			d.Region = Region
			d.City = City
			d.NetGeneration = NetGen / 12
			d.ServiceHours = ServiceHours / 12
			if d.NetGeneration != 0 {
				d.ValueByNetGeneration = d.Value / d.NetGeneration
			}
			if d.ServiceHours != 0 {
				d.ValueByServiceHours = d.Value / d.ServiceHours
			}
			_, e = ctx.InsertOut(d)
			if e != nil {
				tk.Println(e)
				break
			}

			d = new(MORSummary)
			d.Period = time.Date(Year, Month, 1, 0, 0, 0, 0, time.UTC)
			d.Plant = PlantNormalization(cost.Plant)
			d.TopElement = "Operation"
			d.Element = "Fuel Invoices"
			d.SubElement = "Light Crude Oil"
			d.Value = cost.FuelInvoiceCostLightCrudeOil / 12
			d.Province = Province
			d.Region = Region
			d.City = City
			d.NetGeneration = NetGen / 12
			d.ServiceHours = ServiceHours / 12
			if d.NetGeneration != 0 {
				d.ValueByNetGeneration = d.Value / d.NetGeneration
			}
			if d.ServiceHours != 0 {
				d.ValueByServiceHours = d.Value / d.ServiceHours
			}
			_, e = ctx.InsertOut(d)
			if e != nil {
				tk.Println(e)
				break
			}

			d = new(MORSummary)
			d.Period = time.Date(Year, Month, 1, 0, 0, 0, 0, time.UTC)
			d.Plant = PlantNormalization(cost.Plant)
			d.TopElement = "Operation"
			d.Element = "Fuel Invoices"
			d.SubElement = "Gas"
			d.Value = cost.FuelInvoiceCostGas / 12
			d.Province = Province
			d.Region = Region
			d.City = City
			d.NetGeneration = NetGen / 12
			d.ServiceHours = ServiceHours / 12
			if d.NetGeneration != 0 {
				d.ValueByNetGeneration = d.Value / d.NetGeneration
			}
			if d.ServiceHours != 0 {
				d.ValueByServiceHours = d.Value / d.ServiceHours
			}
			_, e = ctx.InsertOut(d)
			if e != nil {
				tk.Println(e)
				break
			}

			d = new(MORSummary)
			d.Period = time.Date(Year, Month, 1, 0, 0, 0, 0, time.UTC)
			d.Plant = PlantNormalization(cost.Plant)
			d.TopElement = "Operation"
			d.Element = "Fuel Invoices"
			d.SubElement = "Diesel"
			d.Value = cost.FuelInvoiceCostDeisel / 12
			d.Province = Province
			d.Region = Region
			d.City = City
			d.NetGeneration = NetGen / 12
			d.ServiceHours = ServiceHours / 12
			if d.NetGeneration != 0 {
				d.ValueByNetGeneration = d.Value / d.NetGeneration
			}
			if d.ServiceHours != 0 {
				d.ValueByServiceHours = d.Value / d.ServiceHours
			}
			_, e = ctx.InsertOut(d)
			if e != nil {
				tk.Println(e)
				break
			}

			d = new(MORSummary)
			d.Period = time.Date(Year, Month, 1, 0, 0, 0, 0, time.UTC)
			d.Plant = PlantNormalization(cost.Plant)
			d.TopElement = "Operation"
			d.Element = "Fuel Invoices"
			d.SubElement = "Havey"
			d.Value = cost.FuelInvoiceCostHavey / 12
			d.Province = Province
			d.Region = Region
			d.City = City
			d.NetGeneration = NetGen / 12
			d.ServiceHours = ServiceHours / 12
			if d.NetGeneration != 0 {
				d.ValueByNetGeneration = d.Value / d.NetGeneration
			}
			if d.ServiceHours != 0 {
				d.ValueByServiceHours = d.Value / d.ServiceHours
			}
			_, e = ctx.InsertOut(d)
			if e != nil {
				tk.Println(e)
				break
			}

			d = new(MORSummary)
			d.Period = time.Date(Year, Month, 1, 0, 0, 0, 0, time.UTC)
			d.Plant = PlantNormalization(cost.Plant)
			d.TopElement = "Operation"
			d.Element = "Fuel Invoices"
			d.SubElement = "Light Crude Oil Filteration"
			d.Value = cost.FuelInvoiceCostLightCrudeOilfilteration / 12
			d.Province = Province
			d.Region = Region
			d.City = City
			d.NetGeneration = NetGen / 12
			d.ServiceHours = ServiceHours / 12
			if d.NetGeneration != 0 {
				d.ValueByNetGeneration = d.Value / d.NetGeneration
			}
			if d.ServiceHours != 0 {
				d.ValueByServiceHours = d.Value / d.ServiceHours
			}
			_, e = ctx.InsertOut(d)
			if e != nil {
				tk.Println(e)
				break
			}

			d = new(MORSummary)
			d.Period = time.Date(Year, Month, 1, 0, 0, 0, 0, time.UTC)
			d.Plant = PlantNormalization(cost.Plant)
			d.TopElement = "Operation"
			d.Element = "Fuel Invoices"
			d.SubElement = "Havey Crude Oil FIlteration"
			d.Value = cost.FuelInvoiceCostHaveyCrudeOilfilteration / 12
			d.Province = Province
			d.Region = Region
			d.City = City
			d.NetGeneration = NetGen / 12
			d.ServiceHours = ServiceHours / 12
			if d.NetGeneration != 0 {
				d.ValueByNetGeneration = d.Value / d.NetGeneration
			}
			if d.ServiceHours != 0 {
				d.ValueByServiceHours = d.Value / d.ServiceHours
			}
			_, e = ctx.InsertOut(d)
			if e != nil {
				tk.Println(e)
				break
			}

			d = new(MORSummary)
			d.Period = time.Date(Year, Month, 1, 0, 0, 0, 0, time.UTC)
			d.Plant = PlantNormalization(cost.Plant)
			d.TopElement = "Operation"
			d.Element = "Shared Services"
			d.SubElement = "Material Service Cost Allocated"
			d.Value = cost.MaterialsServiceCostAllocated / 12
			d.Province = Province
			d.Region = Region
			d.City = City
			d.NetGeneration = NetGen / 12
			d.ServiceHours = ServiceHours / 12
			if d.NetGeneration != 0 {
				d.ValueByNetGeneration = d.Value / d.NetGeneration
			}
			if d.ServiceHours != 0 {
				d.ValueByServiceHours = d.Value / d.ServiceHours
			}
			_, e = ctx.InsertOut(d)
			if e != nil {
				tk.Println(e)
				break
			}

			d = new(MORSummary)
			d.Period = time.Date(Year, Month, 1, 0, 0, 0, 0, time.UTC)
			d.Plant = PlantNormalization(cost.Plant)
			d.TopElement = "Operation"
			d.Element = "Shared Services"
			d.SubElement = "Human Resource Services Cost"
			d.Value = cost.HumanResourcesServicesCost / 12
			d.Province = Province
			d.Region = Region
			d.City = City
			d.NetGeneration = NetGen / 12
			d.ServiceHours = ServiceHours / 12
			if d.NetGeneration != 0 {
				d.ValueByNetGeneration = d.Value / d.NetGeneration
			}
			if d.ServiceHours != 0 {
				d.ValueByServiceHours = d.Value / d.ServiceHours
			}
			_, e = ctx.InsertOut(d)
			if e != nil {
				tk.Println(e)
				break
			}

			d = new(MORSummary)
			d.Period = time.Date(Year, Month, 1, 0, 0, 0, 0, time.UTC)
			d.Plant = PlantNormalization(cost.Plant)
			d.TopElement = "Operation"
			d.Element = "Shared Services Differential"
			d.SubElement = "Material Service Cost Differential"
			d.Value = cost.MaterialServicesCostDifferentia / 12
			d.Province = Province
			d.Region = Region
			d.City = City
			d.NetGeneration = NetGen / 12
			d.ServiceHours = ServiceHours / 12
			if d.NetGeneration != 0 {
				d.ValueByNetGeneration = d.Value / d.NetGeneration
			}
			if d.ServiceHours != 0 {
				d.ValueByServiceHours = d.Value / d.ServiceHours
			}
			_, e = ctx.InsertOut(d)
			if e != nil {
				tk.Println(e)
				break
			}

			d = new(MORSummary)
			d.Period = time.Date(Year, Month, 1, 0, 0, 0, 0, time.UTC)
			d.Plant = PlantNormalization(cost.Plant)
			d.TopElement = "Operation"
			d.Element = "Human Resources"
			d.SubElement = "Training Development"
			d.Value = cost.TrainingDevelopmentCostAllocat / 12
			d.Province = Province
			d.Region = Region
			d.City = City
			d.NetGeneration = NetGen / 12
			d.ServiceHours = ServiceHours / 12
			if d.NetGeneration != 0 {
				d.ValueByNetGeneration = d.Value / d.NetGeneration
			}
			if d.ServiceHours != 0 {
				d.ValueByServiceHours = d.Value / d.ServiceHours
			}
			_, e = ctx.InsertOut(d)
			if e != nil {
				tk.Println(e)
				break
			}

			d = new(MORSummary)
			d.Period = time.Date(Year, Month, 1, 0, 0, 0, 0, time.UTC)
			d.Plant = PlantNormalization(cost.Plant)
			d.TopElement = "Operation"
			d.Element = "Business Unit Overhead"
			d.SubElement = "Power Plant Operation"
			d.Value = cost.PowerPlantOperationalCost / 12
			d.Province = Province
			d.Region = Region
			d.City = City
			d.NetGeneration = NetGen / 12
			d.ServiceHours = ServiceHours / 12
			if d.NetGeneration != 0 {
				d.ValueByNetGeneration = d.Value / d.NetGeneration
			}
			if d.ServiceHours != 0 {
				d.ValueByServiceHours = d.Value / d.ServiceHours
			}
			_, e = ctx.InsertOut(d)
			if e != nil {
				tk.Println(e)
				break
			}

			d = new(MORSummary)
			d.Period = time.Date(Year, Month, 1, 0, 0, 0, 0, time.UTC)
			d.Plant = PlantNormalization(cost.Plant)
			d.TopElement = "Operation"
			d.Element = "Business Unit Overhead"
			d.SubElement = "Power Plant Maintenance"
			d.Value = cost.PowerPlantMaintenanceCost / 12
			d.Province = Province
			d.Region = Region
			d.City = City
			d.NetGeneration = NetGen / 12
			d.ServiceHours = ServiceHours / 12
			if d.NetGeneration != 0 {
				d.ValueByNetGeneration = d.Value / d.NetGeneration
			}
			if d.ServiceHours != 0 {
				d.ValueByServiceHours = d.Value / d.ServiceHours
			}
			_, e = ctx.InsertOut(d)
			if e != nil {
				tk.Println(e)
				break
			}

			d = new(MORSummary)
			d.Period = time.Date(Year, Month, 1, 0, 0, 0, 0, time.UTC)
			d.Plant = PlantNormalization(cost.Plant)
			d.TopElement = "Operation"
			d.Element = "Business Unit Overhead"
			d.SubElement = "Power Plant Technical Support"
			d.Value = cost.PowerPlantTechnicalSupport / 12
			d.Province = Province
			d.Region = Region
			d.City = City
			d.NetGeneration = NetGen / 12
			d.ServiceHours = ServiceHours / 12
			if d.NetGeneration != 0 {
				d.ValueByNetGeneration = d.Value / d.NetGeneration
			}
			if d.ServiceHours != 0 {
				d.ValueByServiceHours = d.Value / d.ServiceHours
			}
			_, e = ctx.InsertOut(d)
			if e != nil {
				tk.Println(e)
				break
			}

			d = new(MORSummary)
			d.Period = time.Date(Year, Month, 1, 0, 0, 0, 0, time.UTC)
			d.Plant = PlantNormalization(cost.Plant)
			d.TopElement = "Operation"
			d.Element = "Secondary Allocated"
			d.SubElement = "Secondary Labor - Saudi"
			d.Value = cost.SecondaryLaborCostCharedSaudi / 12
			d.Province = Province
			d.Region = Region
			d.City = City
			d.NetGeneration = NetGen / 12
			d.ServiceHours = ServiceHours / 12
			if d.NetGeneration != 0 {
				d.ValueByNetGeneration = d.Value / d.NetGeneration
			}
			if d.ServiceHours != 0 {
				d.ValueByServiceHours = d.Value / d.ServiceHours
			}
			_, e = ctx.InsertOut(d)
			if e != nil {
				tk.Println(e)
				break
			}

			d = new(MORSummary)
			d.Period = time.Date(Year, Month, 1, 0, 0, 0, 0, time.UTC)
			d.Plant = PlantNormalization(cost.Plant)
			d.TopElement = "Operation"
			d.Element = "Secondary Labor Cost Allocated"
			d.SubElement = "Secondary Labor - Non Saudi"
			d.Value = cost.SecondaryLaborCostChargedNonSau / 12
			d.Province = Province
			d.Region = Region
			d.City = City
			d.NetGeneration = NetGen / 12
			d.ServiceHours = ServiceHours / 12
			if d.NetGeneration != 0 {
				d.ValueByNetGeneration = d.Value / d.NetGeneration
			}
			if d.ServiceHours != 0 {
				d.ValueByServiceHours = d.Value / d.ServiceHours
			}
			_, e = ctx.InsertOut(d)
			if e != nil {
				tk.Println(e)
				break
			}

			d = new(MORSummary)
			d.Period = time.Date(Year, Month, 1, 0, 0, 0, 0, time.UTC)
			d.Plant = PlantNormalization(cost.Plant)
			d.TopElement = "Operation"
			d.Element = "Secondary Labor Cost Allocated"
			d.SubElement = "Secondary Labor - Saudi"
			d.Value = cost.SecondaryAllocated / 12
			d.Province = Province
			d.Region = Region
			d.City = City
			d.NetGeneration = NetGen / 12
			d.ServiceHours = ServiceHours / 12
			if d.NetGeneration != 0 {
				d.ValueByNetGeneration = d.Value / d.NetGeneration
			}
			if d.ServiceHours != 0 {
				d.ValueByServiceHours = d.Value / d.ServiceHours
			}
			_, e = ctx.InsertOut(d)
			if e != nil {
				tk.Println(e)
				break
			}

			d = new(MORSummary)
			d.Period = time.Date(Year, Month, 1, 0, 0, 0, 0, time.UTC)
			d.Plant = PlantNormalization(cost.Plant)
			d.TopElement = "Operation"
			d.Element = "Secondary OT Cost Allocated"
			d.SubElement = "Secondary OT Labor - Saudi"
			d.Value = cost.SecondaryOtCostChargedSaudi / 12
			d.Province = Province
			d.Region = Region
			d.City = City
			d.NetGeneration = NetGen / 12
			d.ServiceHours = ServiceHours / 12
			if d.NetGeneration != 0 {
				d.ValueByNetGeneration = d.Value / d.NetGeneration
			}
			if d.ServiceHours != 0 {
				d.ValueByServiceHours = d.Value / d.ServiceHours
			}
			_, e = ctx.InsertOut(d)
			if e != nil {
				tk.Println(e)
				break
			}
		}

	}

	return nil
}

// generateMORFlatCalculationSummary
func (m *GenMOR) generateMORFlatCalculationSummary() error {
	ctx := m.BaseController.Ctx
	c := ctx.Connection
	var (
		query []*dbox.Filter
	)
	tk.Println("Generating MOR Flat Calculation Summary..")
	Years := []int{2013, 2014, 2015}

	query = []*dbox.Filter{}
	query = append(query, dbox.Gte("TopElement", "Maintenance"))

	MORSummaryList := []MORSummary{}
	csr, e := c.NewQuery().Select("Element").From(new(MORSummary).TableName()).Where(query...).Group("Element").Cursor(nil)

	if e != nil {
		return e
	}

	e = csr.Fetch(&MORSummaryList, 0, false)
	if e != nil {
		return e
	}
	csr.Close()

	for _, year := range Years {
		query = []*dbox.Filter{}
		query = append(query, dbox.Gte("Period", time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)))
		query = append(query, dbox.Lt("Period", time.Date((year+1), 1, 1, 0, 0, 0, 0, time.UTC)))

		MaintenanceCostList := []MaintenanceCost{}
		csr, e := c.NewQuery().From(new(MaintenanceCost).TableName()).Where(query...).Cursor(nil)

		if e != nil {
			return e
		}

		e = csr.Fetch(&MaintenanceCostList, 0, false)
		if e != nil {
			return e
		}
		csr.Close()

		Plants := crowd.From(&MaintenanceCostList).Group(func(x interface{}) interface{} {
			return x.(MaintenanceCost).Plant
		}, nil).Exec().Result.Data().([]crowd.KV)

		for _, p := range Plants {
			plant := p.Key.(string)
			EqType := crowd.From(&MaintenanceCostList).Where(func(x interface{}) interface{} {
				return x.(MaintenanceCost).Plant == plant
			}).Group(func(x interface{}) interface{} {
				return x.(MaintenanceCost).EquipmentType
			}, nil).Exec().Result.Data().([]crowd.KV)

			for _, eqt := range EqType {
				eq := eqt.Key.(string)
				ActType := crowd.From(&MaintenanceCostList).Where(func(x interface{}) interface{} {
					o := x.(MaintenanceCost)
					return o.Plant == plant && o.EquipmentType == eq
				}).Group(func(x interface{}) interface{} {
					return x.(MaintenanceCost).MaintenanceActivityType
				}, nil).Exec().Result.Data().([]crowd.KV)

				for _, a := range ActType {
					act := a.Key.(string)
					OrderType := crowd.From(&MaintenanceCostList).Where(func(x interface{}) interface{} {
						o := x.(MaintenanceCost)
						return o.Plant == plant && o.EquipmentType == eq && o.MaintenanceActivityType == act
					}).Group(func(x interface{}) interface{} {
						return x.(MaintenanceCost).OrderType
					}, nil).Exec().Result.Data().([]crowd.KV)

					for _, o := range OrderType {
						order := o.Key.(string)
						for _, mor := range MORSummaryList {
							d := new(MORCalculationFlatSummary)
							d.PeriodYear = year
							d.OrderType = order
							if len(eq) == 1 {
								d.EquipmentType = "Other"
							} else {
								d.EquipmentType = eq
							}
							if len(eq) == 1 {
								d.EquipmentTypeDescription = "Other"
							} else {
								EqTypeDesc := crowd.From(&MaintenanceCostList).Where(func(x interface{}) interface{} {
									o := x.(MaintenanceCost)
									return o.Plant == plant && o.EquipmentType == eq
								}).Exec().Result.Data().([]MaintenanceCost)
								if len(EqTypeDesc) > 0 {
									d.EquipmentTypeDescription = EqTypeDesc[0].EquipmentTypeDesc
								}
							}

							d.ActivityType = act
							d.Plant = PlantNormalization(plant)
							d.Element = mor.Element
							d.MOCount = len(OrderType)

							switch d.Element {
							case "Internal Labor":
								d.Value = crowd.From(&MaintenanceCostList).Where(func(x interface{}) interface{} {
									o := x.(MaintenanceCost)
									return o.Plant == plant && o.EquipmentType == eq && o.MaintenanceActivityType == act
								}).Sum(func(x interface{}) interface{} {
									return x.(MaintenanceCost).InternalLaborActual
								}).Exec().Result.Sum
								break
							case "Internal Material":
								d.Value = crowd.From(&MaintenanceCostList).Where(func(x interface{}) interface{} {
									o := x.(MaintenanceCost)
									return o.Plant == plant && o.EquipmentType == eq && o.MaintenanceActivityType == act
								}).Sum(func(x interface{}) interface{} {
									return x.(MaintenanceCost).InternalMaterialActual
								}).Exec().Result.Sum
								break
							case "Direct Material":
								d.Value = crowd.From(&MaintenanceCostList).Where(func(x interface{}) interface{} {
									o := x.(MaintenanceCost)
									return o.Plant == plant && o.EquipmentType == eq && o.MaintenanceActivityType == act
								}).Sum(func(x interface{}) interface{} {
									return x.(MaintenanceCost).DirectMaterialActual
								}).Exec().Result.Sum
								break
							case "External Service":
								d.Value = crowd.From(&MaintenanceCostList).Where(func(x interface{}) interface{} {
									o := x.(MaintenanceCost)
									return o.Plant == plant && o.EquipmentType == eq && o.MaintenanceActivityType == act
								}).Sum(func(x interface{}) interface{} {
									return x.(MaintenanceCost).ExternalServiceActual
								}).Exec().Result.Sum
								break
							default:
								break
							}

							_, e := ctx.InsertOut(d)
							if e != nil {
								tk.Println(e)
								break
							}

						}
					}

				}
			}
		}

	}
	tk.Println(len(MORSummaryList))
	return nil
}
