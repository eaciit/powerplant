package controllers

import (
	"github.com/eaciit/crowd"
	"github.com/eaciit/dbox"
	// . "github.com/eaciit/powerplant/sec/consoleapp/generator/helpers"
	. "github.com/eaciit/powerplant/sec/library/models"
	tk "github.com/eaciit/toolkit"
	"strconv"
	"strings"
	"time"
)

// GenValueEquation
type GenValueEquation struct {
	*BaseController
}

// Generate
func (d *GenValueEquation) Generate(base *BaseController) {
	var (
		e error
	)
	if base != nil {
		d.BaseController = base
	}

	e = d.generateValueEquationDataQuality(2014, "Qurayyah CC")
	if e != nil {
		tk.Println(e)
	}
	tk.Println("##Value Equation Data : DONE\n")
}

func (d *GenValueEquation) generateValueEquationDataQuality(Year int, Plant string) error {
	ctx := d.BaseController.Ctx
	c := ctx.Connection
	var (
		query []*dbox.Filter
		e     error
	)

	tk.Println("Generating Value Equation Data Quality..")

	// Get Performance Factors
	query = append(query, dbox.Eq("Year", Year))
	query = append(query, dbox.Eq("Plant", Plant))
	csr, e := c.NewQuery().From(new(PerformanceFactors).TableName()).Where(query...).Cursor(nil)
	PerformanceFactorsData := []PerformanceFactors{}
	e = csr.Fetch(&PerformanceFactorsData, 0, false)
	csr.Close()

	// Get Consolidated Data
	query = append(query[0:0], dbox.Eq("Plant", Plant))
	csr, e = c.NewQuery().From(new(Consolidated).TableName()).Where(query...).Cursor(nil)
	ConsolidatedData := []Consolidated{}
	e = csr.Fetch(&ConsolidatedData, 0, false)
	csr.Close()

	// Get PrevMaintenanceValueEquation
	query = append(query[0:0], dbox.Eq("Plant", Plant))
	csr, e = c.NewQuery().From(new(PrevMaintenanceValueEquation).TableName()).Where(query...).Cursor(nil)
	PrevMaintenanceData := []PrevMaintenanceValueEquation{}
	e = csr.Fetch(&PrevMaintenanceData, 0, false)
	csr.Close()

	// Get PowerPlantOutages
	query = append(query[0:0], dbox.Eq("Id", strconv.Itoa(Year)+Plant))
	csr, e = c.NewQuery().From(new(PowerPlantOutagesDetails).TableName()).Where(query...).Cursor(nil)
	PowerPlantOutagesData := []PowerPlantOutagesDetails{}
	e = csr.Fetch(&PowerPlantOutagesData, 0, false)
	csr.Close()

	// Get FuelCost
	query = append(query[0:0], dbox.Eq("Plant", Plant))
	query = append(query, dbox.Eq("Year", Year))
	csr, e = c.NewQuery().From(new(FuelCost).TableName()).Where(query...).Cursor(nil)
	FuelCostData := []FuelCost{}
	e = csr.Fetch(&FuelCostData, 0, false)
	csr.Close()

	// // Get SyntheticPM
	query = append(query[0:0], dbox.Eq("Plant", Plant))
	csr, e = c.NewQuery().From(new(SyntheticPM).TableName()).Where(query...).Cursor(nil)
	SyntheticPMData := []SyntheticPM{}
	e = csr.Fetch(&SyntheticPMData, 0, false)
	csr.Close()

	// // Get SyntheticPM
	query = append(query[0:0], dbox.Eq("Plant", Plant))
	query = append(query, dbox.Eq("Year", Year))
	csr, e = c.NewQuery().From(new(FuelTransport).TableName()).Where(query...).Cursor(nil)
	FuelTransportData := []FuelTransport{}
	e = csr.Fetch(&FuelTransportData, 0, false)
	csr.Close()

	// // Get 	Plant Code List
	// query = append(query[0:0], dbox.Eq("PlantName", Plant))
	// csr, e = c.NewQuery().From(new(PowerPlantCoordinates).TableName()).Where(query...).Cursor(nil)
	// PowerPlantCoordinatesData := []PowerPlantCoordinates{}
	// e = csr.Fetch(&PowerPlantCoordinatesData, 0, false)
	// csr.Close()
	// PlantCodeList := []interface{}{}
	// for _, i := range PowerPlantCoordinatesData {
	// 	PlantCodeList = append(PlantCodeList, i.PlantCode)
	// }
	// // Get Data Browser
	// query = append(query[0:0], dbox.Eq("PeriodYear", Year))
	// query = append(query, dbox.In("PlantCode", PlantCodeList...))
	// csr, e = c.NewQuery().From(new(DataBrowser).TableName()).Where(query...).Cursor(nil)
	// DataBrowserData := []DataBrowser{}
	// e = csr.Fetch(&DataBrowserData, 0, false)
	// csr.Close()

	// // GEt Generation APpendix
	csr, e = c.NewQuery().From(new(GenerationAppendix).TableName()).Cursor(nil)
	GenerationAppendixData := []GenerationAppendix{}
	e = csr.Fetch(&GenerationAppendixData, 0, false)
	csr.Close()

	UnitData := crowd.From(&FuelCostData).Group(func(x interface{}) interface{} {
		return strings.Replace(x.(FuelCost).UnitId, " ", "", -1)
	}, nil).Exec().Result.Data().([]crowd.KV)
	Units := []string{}
	for _, u := range UnitData {
		Units = append(Units, u.Key.(string))
	}
	// Declare DocumentName
	// AppendixFile := "Appendix File"
	// ConsolidatedFile := "Consolidated File"
	// SyntheticFile := "Synthetic File"
	// PerformanceFile := "Performance Factor File"
	// FuelTransportFile := "Fuel Transport File"
	// OutagesFile := "Outages File"
	// FuelFile := "Fuel File"
	// MaintenanceFile := "Maintenance File"
	// DerivedMaintenanceFile := "Derived Maintenance File"

	tempFuelCost := crowd.From(&FuelCostData).Where(func(x interface{}) interface{} {
		return x.(FuelCost).PrimaryFuelType == "DIESEL"
	}).Exec().Result.Data().([]FuelCost)
	DieselPrimaryFuelConsumption := crowd.From(&tempFuelCost).Sum(func(x interface{}) interface{} {
		return x.(FuelCost).PrimaryFuelConsumed
	}).Exec().Result.Sum

	tempFuelCost = crowd.From(&FuelCostData).Where(func(x interface{}) interface{} {
		return x.(FuelCost).BackupFuelType == "DIESEL"
	}).Exec().Result.Data().([]FuelCost)
	DieselBackupFuelConsumption := crowd.From(&tempFuelCost).Sum(func(x interface{}) interface{} {
		return x.(FuelCost).BackupFuelConsumed
	}).Exec().Result.Sum

	DieselConsumptions := DieselPrimaryFuelConsumption + DieselBackupFuelConsumption*1000
	var TransportCosts float64 = 0
	if DieselConsumptions != 0 {
		for _, i := range FuelTransportData {
			TransportCosts = i.TransportCost
		}
		TransportCosts = TransportCosts / DieselConsumptions
	}

	for _, unit := range Units {
		NormalizedUnit := ""
		if len(unit) < 3 {
			if Plant == "PP9" {
				NormalizedUnit = "GT" + unit
			}
		} else {
			NormalizedUnit = strings.Replace(strings.Replace(unit, ".", "", -1), " ", "", -1)
		}
		tempunit := strings.Replace(strings.Replace(NormalizedUnit, ".", "", -1), " ", "", -1)

		if len(tempunit) == 3 && !strings.ContainsAny(tempunit, "ST") {
			tempunit = "GT0" + strings.Replace(tempunit, "GT", "", -1)
		}

		data := new(ValueEquationDataQuality)
		data.Plant = Plant
		data.Dates = time.Date(Year, 1, 1, 0, 0, 0, 0, time.UTC)
		data.Month = 1
		data.Year = Year
		data.Unit = strings.Replace(strings.Replace(NormalizedUnit, ".", "", -1), " ", "", -1)

		tempAppendix := []GenerationAppendix{}
		if strings.ContainsAny(data.Unit, "ST") {
			if Plant == "Qurayyah" {
				tempAppendix = crowd.From(&GenerationAppendixData).Where(func(x interface{}) interface{} {
					return x.(GenerationAppendix).Plant == "QRPP"
				}).Exec().Result.Data().([]GenerationAppendix)
			} else if Plant == "Qurayyah CC" {
				tempAppendix = crowd.From(&GenerationAppendixData).Where(func(x interface{}) interface{} {
					return x.(GenerationAppendix).Plant == "QCCP"
				}).Exec().Result.Data().([]GenerationAppendix)
			} else if Plant == "Ghazlan" {
				unitnumber, _ := strconv.Atoi(strings.Replace(data.Unit, "ST", "", -1))
				if unitnumber <= 4 {
					tempAppendix = crowd.From(&GenerationAppendixData).Where(func(x interface{}) interface{} {
						return x.(GenerationAppendix).Plant == "Ghazlan I (1-4)"
					}).Exec().Result.Data().([]GenerationAppendix)
				} else if unitnumber <= 8 {
					tempAppendix = crowd.From(&GenerationAppendixData).Where(func(x interface{}) interface{} {
						return x.(GenerationAppendix).Plant == "Ghazlan II (5-8)"
					}).Exec().Result.Data().([]GenerationAppendix)
				}
			} else if Plant == "Shoaiba" {
				tempAppendix = crowd.From(&GenerationAppendixData).Where(func(x interface{}) interface{} {
					return x.(GenerationAppendix).Plant == "Shoiba Steam"
				}).Exec().Result.Data().([]GenerationAppendix)
			} else if Plant == "Rabigh" {
				tempAppendix = crowd.From(&GenerationAppendixData).Where(func(x interface{}) interface{} {
					return x.(GenerationAppendix).Plant == "Rabigh Steam"
				}).Exec().Result.Data().([]GenerationAppendix)
			} else if Plant == "PP9" {
				tempAppendix = crowd.From(&GenerationAppendixData).Where(func(x interface{}) interface{} {
					return x.(GenerationAppendix).Plant == "PP9 CC"
				}).Exec().Result.Data().([]GenerationAppendix)
			}
		} else {
			if Plant == "Qurayyah" {
				tempAppendix = crowd.From(&GenerationAppendixData).Where(func(x interface{}) interface{} {
					return x.(GenerationAppendix).Plant == "QRPP"
				}).Exec().Result.Data().([]GenerationAppendix)
			} else if Plant == "Qurayyah CC" {
				tempAppendix = crowd.From(&GenerationAppendixData).Where(func(x interface{}) interface{} {
					return x.(GenerationAppendix).Plant == "QCCP"
				}).Exec().Result.Data().([]GenerationAppendix)
			} else if Plant == "Rabigh" {
				unitnumber, _ := strconv.Atoi(strings.Replace(data.Unit, "GT", "", -1))
				if unitnumber <= 12 {
					tempAppendix = crowd.From(&GenerationAppendixData).Where(func(x interface{}) interface{} {
						return x.(GenerationAppendix).Plant == "Rabigh Combined"
					}).Exec().Result.Data().([]GenerationAppendix)
				} else if unitnumber <= 40 {
					tempAppendix = crowd.From(&GenerationAppendixData).Where(func(x interface{}) interface{} {
						return x.(GenerationAppendix).Plant == "Rabigh Gas"
					}).Exec().Result.Data().([]GenerationAppendix)
				}
			} else if Plant == "PP9" {
				unitnumber, _ := strconv.Atoi(strings.Replace(data.Unit, "GT", "", -1))
				if unitnumber <= 16 {
					tempAppendix = crowd.From(&GenerationAppendixData).Where(func(x interface{}) interface{} {
						return x.(GenerationAppendix).Plant == "PP9 CC"
					}).Exec().Result.Data().([]GenerationAppendix)
				} else if unitnumber <= 24 {
					tempAppendix = crowd.From(&GenerationAppendixData).Where(func(x interface{}) interface{} {
						return x.(GenerationAppendix).Plant == "PPEXT" && x.(GenerationAppendix).Units == 8
					}).Exec().Result.Data().([]GenerationAppendix)
				} else if unitnumber <= 56 {
					tempAppendix = crowd.From(&GenerationAppendixData).Where(func(x interface{}) interface{} {
						return x.(GenerationAppendix).Plant == "PPEXT" && x.(GenerationAppendix).Units == 32
					}).Exec().Result.Data().([]GenerationAppendix)
				}
			}
		}
		if len(tempAppendix) > 0 {
			data.Appendix_Data = 1
		} else {
			data.Appendix_Data = 3
		}

		tempConsolidatedData := crowd.From(&ConsolidatedData).Where(func(x interface{}) interface{} {
			return strings.Replace(x.(Consolidated).Unit, "ST0", "ST", -1) == strings.Replace(data.Unit, "ST0", "ST", -1)
		}).Exec().Result.Data().([]Consolidated)
		if len(tempConsolidatedData) > 0 {
			data.Consolidated_Data = 1
		} else {
			data.Consolidated_Data = 3
		}

		tempSyntheticPMData := crowd.From(&SyntheticPMData).Where(func(x interface{}) interface{} {
			return strings.Replace(strings.Replace(strings.Replace(x.(SyntheticPM).Unit, "GT0", "", -1), "GT", "", -1), "ST", "S", -1) == strings.Replace(strings.Replace(strings.Replace(data.Unit, "GT0", "", -1), "GT", "", -1), "ST", "S", -1)
		}).Exec().Result.Data().([]SyntheticPM)
		if len(tempSyntheticPMData) > 0 {
			data.Synthetic_Data = 1
		} else {
			data.Synthetic_Data = 3
		}

		tempPerformanceFactorsData := crowd.From(&PerformanceFactorsData).Where(func(x interface{}) interface{} {
			return x.(PerformanceFactors).Unit == strings.Replace(data.Unit, "ST0", "ST", -1)
		}).Exec().Result.Data().([]PerformanceFactors)
		if len(tempPerformanceFactorsData) > 0 {
			data.PerformanceFactor_Data = 1
		} else {
			data.PerformanceFactor_Data = 3
		}
		if TransportCosts > 0 {
			data.FuelTransport_Data = 1
		} else {
			data.FuelTransport_Data = 3
		}

		if Plant == "Rabigh" {
			if len(PowerPlantOutagesData) == 0 {
				data.Outages_Data = 3
			} else {
				tempPowerPlantOutagesData := crowd.From(&PowerPlantOutagesData).Where(func(x interface{}) interface{} {
					return x.(PowerPlantOutagesDetails).UnitNo == data.Unit && x.(PowerPlantOutagesDetails).PlantName == "Rabigh Steam"
				}).Exec().Result.Data().([]PowerPlantOutagesDetails)
				if len(tempPowerPlantOutagesData) > 0 {
					data.Outages_Data = 1
				} else {
					data.Outages_Data = 3
				}
			}
		} else if Plant == "Qurayyah" || Plant == "Qrayyah CC" {
			if len(PowerPlantOutagesData) == 0 {
				data.Outages_Data = 3
			} else {

				tempPowerPlantOutagesData := crowd.From(&PowerPlantOutagesData).Where(func(x interface{}) interface{} {
					return x.(PowerPlantOutagesDetails).UnitNo == data.Unit && x.(PowerPlantOutagesDetails).PlantName == Plant
				}).Exec().Result.Data().([]PowerPlantOutagesDetails)
				if len(tempPowerPlantOutagesData) > 0 {
					data.Outages_Data = 1
				} else {
					data.Outages_Data = 3
				}
			}
		} else {
			if len(PowerPlantOutagesData) == 0 {
				data.Outages_Data = 3
			} else {
				tempPowerPlantOutagesData := crowd.From(&PowerPlantOutagesData).Where(func(x interface{}) interface{} {
					return x.(PowerPlantOutagesDetails).UnitNo == data.Unit
				}).Exec().Result.Data().([]PowerPlantOutagesDetails)
				if len(tempPowerPlantOutagesData) > 0 {
					data.Outages_Data = 1
				} else {
					data.Outages_Data = 3
				}
			}
		}
		if data.Appendix_Data == 1 {
			data.CapacityPayment_Data = 1
		} else {
			data.CapacityPayment_Data = 3
		}
		if data.Consolidated_Data == 1 && data.Appendix_Data == 1 {
			data.EnergyPayment_Data = 1
		} else {
			if data.Consolidated_Data == 3 && data.Appendix_Data == 3 {
				data.EnergyPayment_Data = 3
			} else {
				data.EnergyPayment_Data = 2
			}
		}
		if data.Appendix_Data == 1 && data.Outages_Data == 1 {
			data.StartupPayment_Data = 1
		} else {
			if data.Appendix_Data == 3 && data.Outages_Data == 3 {
				data.StartupPayment_Data = 3
			} else {
				data.StartupPayment_Data = 2
			}
		}
		data.Incentive_Data = 3

		tempFuelCost := crowd.From(&FuelCostData).Where(func(x interface{}) interface{} {
			return strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(x.(FuelCost).UnitId, " ", "", -1), ".", "", -1), "ST0", "ST", -1), "GT0", "", -1), "GT", "", -1) == strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(data.Unit, " ", "", -1), ".", "", -1), "ST0", "ST", -1), "GT0", "", -1), "GT", "", -1)
		}).Exec().Result.Data().([]FuelCost)
		if len(tempFuelCost) > 0 && tempFuelCost[0].PrimaryFuelType != "" {
			data.PrimaryFuel1st_Data = 1
		} else {
			data.PrimaryFuel1st_Data = 3
		}
		if len(tempFuelCost) > 0 && tempFuelCost[0].Primary2FuelType != "" {
			data.PrimaryFuel2nd_Data = 1
		} else {
			data.PrimaryFuel2nd_Data = 3
		}
		// tk.Println(data)
		_, e := ctx.InsertOut(data)
		tk.Println("#")
		if e != nil {
			tk.Println(e)
			break
		}
	}
	return e
}
