package controllers

import (
	"github.com/eaciit/crowd"
	"github.com/eaciit/dbox"
	// . "github.com/eaciit/powerplant/sec/consoleapp/generator/helpers"
	. "github.com/eaciit/powerplant/sec/library/models"
	tk "github.com/eaciit/toolkit"
	"log"
	"strconv"
	"strings"
	"time"
)

// GenValueEquation ...
type GenValueEquation struct {
	*BaseController
}

// Generate ...
func (d *GenValueEquation) Generate(base *BaseController) {
	var (
		e error
	)
	if base != nil {
		d.BaseController = base
	}

	/*e = d.generateValueEquation()
	if e != nil {
		tk.Println(e)
	}

	e = d.generateValueEquationDataQuality(2014, "Qurayyah CC")
	if e != nil {
		tk.Println(e)
	}
	*/

	Plant := []string{"PP9"}
	e = d.generateValueEquationAllMaintenance(2014, Plant)
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

func (d *GenValueEquation) generateValueEquation(Year int, Plant string) error {
	var e error

	ctx := d.BaseController.Ctx
	c := ctx.Connection

	YearFirst := strconv.Itoa(Year) + "-01-01"
	YearLast := strconv.Itoa(Year+1) + "-01-01"

	query := []*dbox.Filter{}

	query = append(query, dbox.Eq("Plant", Plant))
	csr, e := c.NewQuery().From(new(PerformanceFactors).TableName()).Where(query...).Cursor(nil)
	pfs := []tk.M{}
	e = csr.Fetch(&pfs, 0, false)
	csr.Close()

	csr, e = c.NewQuery().From(new(Consolidated).TableName()).Where(query...).Cursor(nil)
	cons := []tk.M{}
	e = csr.Fetch(&cons, 0, false)
	csr.Close()

	query = append(query, dbox.And(dbox.Gte("DatePerformed", YearFirst), dbox.Lt("DatePerformed", YearLast)))
	csr, e = c.NewQuery().From(new(PrevMaintenanceValueEquation).TableName()).Where(query...).Cursor(nil)
	lists := []tk.M{}
	e = csr.Fetch(&lists, 0, false)
	csr.Close()

	query = append(query[0:0], dbox.Eq("Plant", Plant))
	query = append(query, dbox.Eq("Year", Year))

	csr, e = c.NewQuery().From(new(PowerPlantOutages).TableName()).Where(query...).Cursor(nil)
	outages := []tk.M{}
	e = csr.Fetch(&outages, 0, false)
	csr.Close()

	csr, e = c.NewQuery().From(new(StartupPaymentAndPenalty).TableName()).Where(query...).Cursor(nil)
	start := []tk.M{}
	e = csr.Fetch(&start, 0, false)
	csr.Close()

	csr, e = c.NewQuery().From(new(FuelCost).TableName()).Where(query...).Cursor(nil)
	fuelcosts := []tk.M{}
	e = csr.Fetch(&fuelcosts, 0, false)
	csr.Close()

	query = append(query[0:0], dbox.Eq("Plant", Plant))
	query = append(query, dbox.And(dbox.Gte("ScheduledStart", YearFirst), dbox.Lt("ScheduledStart", YearLast)))

	csr, e = c.NewQuery().From(new(SyntheticPM).TableName()).Where(query...).Cursor(nil)
	syn := []tk.M{}
	e = csr.Fetch(&syn, 0, false)
	csr.Close()

	query = append(query[0:0], dbox.Eq("Plant", Plant))
	query = append(query, dbox.Eq("Year", Year))
	csr, e = c.NewQuery().From(new(FuelTransport).TableName()).Where(query...).Cursor(nil)
	trans := []tk.M{}
	e = csr.Fetch(&trans, 0, false)
	csr.Close()

	sintax := "select * from DataBrowser inner join PowerPlantCoordinates on DataBrowser.PlantCode = PowerPlantCoordinates.PlantCode where PeriodYear = " + strconv.Itoa(Year) + " and PowerPlantCoordinates.PlantName = '" + Plant + "'"
	csr, e = c.NewQuery().Command("freequery", tk.M{}.Set("syntax", sintax)).Cursor(nil)
	databrowser := []tk.M{}
	e = csr.Fetch(&databrowser, 0, false)
	csr.Close()

	csr, e = c.NewQuery().From(new(GenerationAppendix).TableName()).Cursor(nil)
	genA := []tk.M{}
	e = csr.Fetch(&genA, 0, false)
	csr.Close()

	csr, e = c.NewQuery().From(new(ValueEquationFuel).TableName()).Cursor(nil)
	valueequationfuel := []tk.M{}
	e = csr.Fetch(&valueequationfuel, 0, false)
	csr.Close()

	UnitData := crowd.From(&pfs).Group(func(x interface{}) interface{} {
		return x.(tk.M).GetString("unit")
	}, nil).Exec().Result.Data().([]crowd.KV)

	Units := []string{}

	for _, u := range UnitData {
		Units = append(Units, u.Key.(string))
	}

	DieselConsumptionsTemp := crowd.From(&fuelcosts).Where(func(x interface{}) interface{} {
		return x.(tk.M).GetString("primaryfueltype") == "DIESEL"
	}).Exec().Result.Data().([]tk.M)

	DieselConsumptions1 := crowd.From(&DieselConsumptionsTemp).Sum(func(x interface{}) interface{} {
		return x.(tk.M).GetFloat64("primaryfuelconsumed")
	}).Exec().Result.Sum

	DieselConsumptionsTemp = crowd.From(&fuelcosts).Where(func(x interface{}) interface{} {
		return x.(tk.M).GetString("backupfueltype") == "DIESEL"
	}).Exec().Result.Data().([]tk.M)

	DieselConsumptions2 := crowd.From(&DieselConsumptionsTemp).Sum(func(x interface{}) interface{} {
		return x.(tk.M).GetFloat64("backupfuelconsumed")
	}).Exec().Result.Sum

	DieselConsumptions := (DieselConsumptions1 + DieselConsumptions2) * 1000

	var TransportCosts float64
	if DieselConsumptions == 0.0 {
		TransportCosts = 0.0
	} else {
		TransportCosts = trans[0].GetFloat64("transportcost") / DieselConsumptions
	}

	_ = TransportCosts
	if len(Units) > 0 {
		for _, unit := range Units {
			tempunit := unit
			if len(tempunit) == 3 && strings.Contains(tempunit, "ST") {
				tempunit = "GT0" + strings.Replace(tempunit, "GT", "", -1)
			}

			val := new(ValueEquation)
			val.Plant = Plant
			val.Dates = time.Date(Year, 1, 1, 0, 0, 0, 0, time.UTC)
			val.Month = 1
			val.Year = Year
			val.Unit = tempunit
			val.UnitGroup = tempunit[0:2]

			phases := crowd.From(&lists).Where(func(x interface{}) interface{} {
				return x.(tk.M).GetString("unit") == tempunit
			}).Exec().Result.Data().([]tk.M)

			if len(phases) > 0 {
				val.Phase = phases[0].GetString("phase")
			}

			CapacityList := crowd.From(&cons).Where(func(x interface{}) interface{} {
				temp1 := strings.Replace(x.(tk.M).GetString("unit"), "ST0", "ST", -1)
				temp2 := strings.Replace(tempunit, "ST0", "ST", -1)

				return temp1 == temp2
			}).Exec().Result.Data().([]tk.M)

			val.Capacity = crowd.From(&CapacityList).Sum(func(x interface{}) interface{} {
				return x.(tk.M).GetFloat64("capacity")
			}).Exec().Result.Sum

			val.NetGeneration = crowd.From(&CapacityList).Sum(func(x interface{}) interface{} {
				return x.(tk.M).GetFloat64("energynet")
			}).Exec().Result.Sum

			//region Revenue
			tempappendix := []tk.M{}
			if strings.Contains(val.Unit, "ST") {
				if Plant == "Qurayyah" {
					tempappendix = crowd.From(&genA).Where(func(x interface{}) interface{} {
						return x.(tk.M).GetString("plant") == "QRPP"
					}).Exec().Result.Data().([]tk.M)
				} else if Plant == "Qurayyah CC" {
					tempappendix = crowd.From(&genA).Where(func(x interface{}) interface{} {
						return x.(tk.M).GetString("plant") == "QCCP"
					}).Exec().Result.Data().([]tk.M)
				} else if Plant == "Ghazlan" {
					unittemp1, _ := strconv.Atoi(strings.Replace(val.Unit, "ST", "", -1))
					if unittemp1 <= 4 {
						tempappendix = crowd.From(&genA).Where(func(x interface{}) interface{} {
							return x.(tk.M).GetString("plant") == "Ghazlan I (1-4)"
						}).Exec().Result.Data().([]tk.M)
					} else if unittemp1 <= 8 {
						tempappendix = crowd.From(&genA).Where(func(x interface{}) interface{} {
							return x.(tk.M).GetString("plant") == "Ghazlan II (5-8)"
						}).Exec().Result.Data().([]tk.M)
					}
				} else if Plant == "Shoaiba" {
					tempappendix = crowd.From(&genA).Where(func(x interface{}) interface{} {
						return x.(tk.M).GetString("plant") == "Shoiba Steam"
					}).Exec().Result.Data().([]tk.M)
				} else if Plant == "Rabigh" {
					tempappendix = crowd.From(&genA).Where(func(x interface{}) interface{} {
						return x.(tk.M).GetString("plant") == "Rabigh"
					}).Exec().Result.Data().([]tk.M)
				} else if Plant == "PP9" {
					tempappendix = crowd.From(&genA).Where(func(x interface{}) interface{} {
						return x.(tk.M).GetString("plant") == "PP9 CC"
					}).Exec().Result.Data().([]tk.M)
				}
			} else {
				if Plant == "Qurayyah" {
					tempappendix = crowd.From(&genA).Where(func(x interface{}) interface{} {
						return x.(tk.M).GetString("plant") == "QRPP"
					}).Exec().Result.Data().([]tk.M)
				} else if Plant == "Qurayyah CC" {
					tempappendix = crowd.From(&genA).Where(func(x interface{}) interface{} {
						return x.(tk.M).GetString("plant") == "QCCP"
					}).Exec().Result.Data().([]tk.M)
				} else if Plant == "Rabigh" {
					unittemp1, _ := strconv.Atoi(strings.Replace(val.Unit, "GT", "", -1))
					if unittemp1 <= 12 {
						tempappendix = crowd.From(&genA).Where(func(x interface{}) interface{} {
							return x.(tk.M).GetString("plant") == "Rabigh Combined"
						}).Exec().Result.Data().([]tk.M)
					} else if unittemp1 <= 40 {
						tempappendix = crowd.From(&genA).Where(func(x interface{}) interface{} {
							return x.(tk.M).GetString("plant") == "Rabigh Gas" && x.(tk.M).GetInt("units") == 28
						}).Exec().Result.Data().([]tk.M)
					}
				} else if Plant == "PP9" {
					unittemp1, _ := strconv.Atoi(strings.Replace(val.Unit, "GT", "", -1))
					if unittemp1 <= 16 {
						tempappendix = crowd.From(&genA).Where(func(x interface{}) interface{} {
							return x.(tk.M).GetString("plant") == "PP9 CC"
						}).Exec().Result.Data().([]tk.M)
					} else if unittemp1 <= 24 {
						tempappendix = crowd.From(&genA).Where(func(x interface{}) interface{} {
							return x.(tk.M).GetString("plant") == "PPEXT" && x.(tk.M).GetInt("units") == 8
						}).Exec().Result.Data().([]tk.M)
					} else if unittemp1 <= 56 {
						tempappendix = crowd.From(&genA).Where(func(x interface{}) interface{} {
							return x.(tk.M).GetString("plant") == "PPEXT" && x.(tk.M).GetInt("units") == 32
						}).Exec().Result.Data().([]tk.M)
					}
				}
			}

			if len(tempappendix) > 0 {
				apendixResult := tempappendix[0].GetFloat64("contractedcapacity") * (tempappendix[0].GetFloat64("fomr") + tempappendix[0].GetFloat64("ccr"))
				totalDays := (time.Date(2014, 12, 31, 0, 0, 0, 0, time.UTC).Sub(time.Date(2013, 12, 31, 0, 0, 0, 0, time.UTC)).Seconds()) / 86400
				val.CapacityPayment = apendixResult * totalDays * 10
			}

			if len(cons) > 0 {
				consResult := crowd.From(&cons).Where(func(x interface{}) interface{} {
					unitCons1 := strings.Replace(x.(tk.M).GetString("unit"), "GT", "", -1)
					unitCons2 := strings.Replace(tempunit, "STO", "ST", -1)

					return unitCons1 == unitCons2
				}).Exec().Result.Data().([]tk.M)

				val.EnergyPayment = crowd.From(&consResult).Sum(func(x interface{}) interface{} {
					return x.(tk.M).GetFloat64("energynet") * tempappendix[0].GetFloat64("vomr") * 10
				}).Exec().Result.Sum
			}

			if len(pfs) > 0 {
				pfsList := crowd.From(&pfs).Where(func(x interface{}) interface{} {
					unitCons1 := strings.Replace(tempunit, "STO", "ST", -1)

					return x.(tk.M).GetString("unit") == unitCons1
				}).Exec().Result.Data().([]tk.M)

				if len(pfsList) > 0 {
					val.SRF = pfsList[0].GetFloat64("srf")
				}
			}

			if Plant == "Rabigh" {
				if len(outages) == 0 {
					val.UnplannedOutages = 0
				} else {
					POId := outages[0].GetString("id")
					query = append(query[0:0], dbox.Eq("POId", POId))
					query = append(query, dbox.Eq("UnitNo", tempunit))
					query = append(query, dbox.Ne("OutageType", "PO"))
					query = append(query, dbox.Eq("PlantName", "Rabigh Steam"))
					csr, e = c.NewQuery().From(new(StartupPaymentAndPenalty).TableName()).Where(query...).Cursor(nil)
					val.UnplannedOutages = float64(csr.Count())
					csr.Close()
				}
			} else if Plant == "Qurayyah" || Plant == "Qurayyah CC" {
				if len(outages) == 0 {
					val.UnplannedOutages = 0
				} else {
					POId := outages[0].GetString("id")

					query = append(query[0:0], dbox.Eq("POId", POId))
					query = append(query, dbox.Eq("UnitNo", tempunit))
					query = append(query, dbox.Ne("OutageType", "PO"))
					query = append(query, dbox.Eq("PlantName", Plant))
					csr, e = c.NewQuery().From(new(StartupPaymentAndPenalty).TableName()).Where(query...).Cursor(nil)
					val.UnplannedOutages = float64(csr.Count())
					csr.Close()
				}
			} else {
				if len(outages) == 0 {
					val.UnplannedOutages = 0
				} else {
					POId := outages[0].GetString("id")

					query = append(query[0:0], dbox.Eq("POId", POId))
					query = append(query, dbox.Eq("UnitNo", tempunit))
					query = append(query, dbox.Ne("OutageType", "PO"))
					csr, e = c.NewQuery().From(new(StartupPaymentAndPenalty).TableName()).Where(query...).Cursor(nil)
					val.UnplannedOutages = float64(csr.Count())
					csr.Close()
				}
			}

			if val.SRF == 100 {
				val.StartupPayment = tempappendix[0].GetFloat64("startup")
				val.PenaltyAmount = 0
			} else {
				val.StartupPayment = 0
				val.PenaltyAmount = tempappendix[0].GetFloat64("deduct")
			}

			val.PenaltyAmount += tempappendix[0].GetFloat64("deduct") * val.UnplannedOutages
			val.Incentive = 0
			val.Revenue = val.CapacityPayment + val.EnergyPayment + val.Incentive + val.StartupPayment - val.PenaltyAmount

			//endregion
			//region OperatingCost
			//region Primary Fuel

			tempResult := crowd.From(&fuelcosts).Where(func(x interface{}) interface{} {
				unitid := x.(tk.M).GetString("unitid")
				unitid = strings.Replace(unitid, " ", "", -1)
				unitid = strings.Replace(unitid, ".", "", -1)
				unitid = strings.Replace(unitid, "ST0", "ST", -1)
				unitid = strings.Replace(unitid, "GT0", "", -1)
				unitid = strings.Replace(unitid, "GT", "", -1)

				tempunit = strings.Replace(tempunit, " ", "", -1)
				tempunit = strings.Replace(tempunit, ".", "", -1)
				tempunit = strings.Replace(tempunit, "ST0", "ST", -1)
				tempunit = strings.Replace(tempunit, "GT0", "", -1)
				tempunit = strings.Replace(tempunit, "GT", "", -1)

				return unitid == tempunit
			}).Exec().Result.Data().([]tk.M)

			PrimaryFuelType := tempResult[0].GetString("primaryfueltype")

			var fuelconsumptionArray []ValueEquationFuel

			if strings.ToLower(strings.Trim(PrimaryFuelType, " ")) == "hfo" {
				//region hfo
				PrimaryFuelConsumed := crowd.From(&tempResult).Sum(func(x interface{}) interface{} {
					return x.(tk.M).GetFloat64("primaryfuelconsumed")
				}).Exec().Result.Sum

				var fuelconsumption ValueEquationFuel

				if strings.ToLower(strings.Trim(val.Plant, " ")) == "shoaiba" {
					fuelconsumption.IsPrimaryFuel = true
					fuelconsumption.FuelType = "CRUDE"
					fuelconsumption.FuelCostPerUnit = 0.1
					fuelconsumption.FuelConsumed = PrimaryFuelConsumed / 3
					fuelconsumption.ConvertedFuelConsumed = fuelconsumption.FuelConsumed * 1000
					fuelconsumption.FuelCost = fuelconsumption.FuelCostPerUnit * fuelconsumption.ConvertedFuelConsumed

					fuelconsumptionArray = append(fuelconsumptionArray, fuelconsumption)

					val.PrimaryFuelTotalCost += fuelconsumption.FuelCost

					fuelconsumption.IsPrimaryFuel = true
					fuelconsumption.FuelType = "CRUDE HEAVY"
					fuelconsumption.FuelCostPerUnit = 0.049
					fuelconsumption.FuelConsumed = PrimaryFuelConsumed / 3
					fuelconsumption.ConvertedFuelConsumed = fuelconsumption.FuelConsumed * 1000
					fuelconsumption.FuelCost = fuelconsumption.FuelCostPerUnit * fuelconsumption.ConvertedFuelConsumed

					fuelconsumptionArray = append(fuelconsumptionArray, fuelconsumption)

					val.PrimaryFuelTotalCost += fuelconsumption.FuelCost

					fuelconsumption.IsPrimaryFuel = true
					fuelconsumption.FuelType = "DIESEL"
					fuelconsumption.FuelCostPerUnit = 0.085
					fuelconsumption.FuelConsumed = PrimaryFuelConsumed / 3
					fuelconsumption.ConvertedFuelConsumed = fuelconsumption.FuelConsumed * 1000
					fuelconsumption.FuelCost = fuelconsumption.FuelCostPerUnit * fuelconsumption.ConvertedFuelConsumed

					fuelconsumptionArray = append(fuelconsumptionArray, fuelconsumption)

					val.PrimaryFuelTotalCost += fuelconsumption.FuelCost
				} else if strings.ToLower(strings.Trim(val.Plant, " ")) == "rabigh" {
					fuelconsumption.IsPrimaryFuel = true
					fuelconsumption.FuelType = "CRUDE"
					fuelconsumption.FuelCostPerUnit = 0.1
					fuelconsumption.FuelConsumed = PrimaryFuelConsumed / 3
					fuelconsumption.ConvertedFuelConsumed = fuelconsumption.FuelConsumed * 1000
					fuelconsumption.FuelCost = fuelconsumption.FuelCostPerUnit * fuelconsumption.ConvertedFuelConsumed

					fuelconsumptionArray = append(fuelconsumptionArray, fuelconsumption)

					val.PrimaryFuelTotalCost += fuelconsumption.FuelCost

					fuelconsumption.IsPrimaryFuel = true
					fuelconsumption.FuelType = "DIESEL"
					fuelconsumption.FuelCostPerUnit = 0.085
					fuelconsumption.FuelConsumed = PrimaryFuelConsumed / 3
					fuelconsumption.ConvertedFuelConsumed = fuelconsumption.FuelConsumed * 1000
					fuelconsumption.FuelCost = fuelconsumption.FuelCostPerUnit * fuelconsumption.ConvertedFuelConsumed

					fuelconsumptionArray = append(fuelconsumptionArray, fuelconsumption)

					val.PrimaryFuelTotalCost += fuelconsumption.FuelCost
				}
				//endregion
			} else {
				//region not hfo
				//var fuelconsumptionArray []ValueEquationFuel
				var fuelconsumption ValueEquationFuel
				fuelconsumption.IsPrimaryFuel = true

				fuelconsumption.FuelType = tempResult[0].GetString("primaryfueltype")
				if strings.Contains(strings.ToLower(fuelconsumption.FuelType), "gas") {
					fuelconsumption.FuelCostPerUnit = 2.813
				} else if strings.Contains(strings.ToLower(fuelconsumption.FuelType), "crude") {
					fuelconsumption.FuelCostPerUnit = 0.1
				} else if strings.Contains(strings.ToLower(fuelconsumption.FuelType), "diesel") {
					fuelconsumption.FuelCostPerUnit = 0.085
				}

				fuelconsumption.FuelConsumed = crowd.From(&tempResult).Sum(func(x interface{}) interface{} {
					return x.(tk.M).GetFloat64("primaryfuelconsumed")
				}).Exec().Result.Sum

				if strings.Contains(strings.ToLower(fuelconsumption.FuelType), "gas") {
					fuelconsumption.ConvertedFuelConsumed = fuelconsumption.FuelConsumed * 0.0353
				} else {
					fuelconsumption.ConvertedFuelConsumed = fuelconsumption.FuelConsumed * 1000
				}

				fuelconsumption.FuelCost = fuelconsumption.FuelCostPerUnit * fuelconsumption.ConvertedFuelConsumed

				fuelconsumptionArray = append(fuelconsumptionArray, fuelconsumption)

				val.PrimaryFuelTotalCost += fuelconsumption.FuelCost

				//endregion
			}
			//endregion

			//region backup fuel
			FuelCostData := crowd.From(&fuelcosts).Where(func(x interface{}) interface{} {
				unitid := x.(tk.M).GetString("unitid")
				unitid = strings.Replace(unitid, " ", "", -1)
				unitid = strings.Replace(unitid, ".", "", -1)
				unitid = strings.Replace(unitid, "ST0", "ST", -1)
				unitid = strings.Replace(unitid, "GT0", "", -1)
				unitid = strings.Replace(unitid, "GT", "", -1)

				unit = strings.Replace(unit, " ", "", -1)
				unit = strings.Replace(unit, ".", "", -1)
				unit = strings.Replace(unit, "ST0", "ST", -1)
				unit = strings.Replace(unit, "GT0", "", -1)
				unit = strings.Replace(unit, "GT", "", -1)

				return unitid == tempunit
			}).Exec().Result.Data().([]tk.M)

			BackupFuelType := FuelCostData[0].GetString("backupfueltype")

			BackupFuelConsumed := crowd.From(&FuelCostData).Sum(func(x interface{}) interface{} {
				return x.(tk.M).GetFloat64("backupfuelconsumed")
			}).Exec().Result.Sum

			var fuelconsumption ValueEquationFuel

			if strings.ToLower(strings.Trim(BackupFuelType, " ")) == "hfo" {
				//#region hfo
				BackupFuelConsumed := crowd.From(&FuelCostData).Sum(func(x interface{}) interface{} {
					return x.(tk.M).GetFloat64("backupfuelconsumed")
				}).Exec().Result.Sum

				if strings.ToLower(strings.Trim(val.Plant, " ")) == "shoaiba" {
					fuelconsumption.IsPrimaryFuel = false
					fuelconsumption.FuelType = "CRUDE"
					fuelconsumption.FuelCostPerUnit = 0.1
					fuelconsumption.FuelConsumed = BackupFuelConsumed / 3
					fuelconsumption.ConvertedFuelConsumed = fuelconsumption.FuelConsumed * 1000
					fuelconsumption.FuelCost = fuelconsumption.FuelCostPerUnit * fuelconsumption.ConvertedFuelConsumed

					fuelconsumptionArray = append(fuelconsumptionArray, fuelconsumption)

					val.BackupFuelTotalCost += fuelconsumption.FuelCost

					fuelconsumption.IsPrimaryFuel = false
					fuelconsumption.FuelType = "CRUDE HEAVY"
					fuelconsumption.FuelCostPerUnit = 0.049
					fuelconsumption.FuelConsumed = BackupFuelConsumed / 3
					fuelconsumption.ConvertedFuelConsumed = fuelconsumption.FuelConsumed * 1000
					fuelconsumption.FuelCost = fuelconsumption.FuelCostPerUnit * fuelconsumption.ConvertedFuelConsumed

					fuelconsumptionArray = append(fuelconsumptionArray, fuelconsumption)

					val.BackupFuelTotalCost += fuelconsumption.FuelCost

					fuelconsumption.IsPrimaryFuel = false
					fuelconsumption.FuelType = "DIESEL"
					fuelconsumption.FuelCostPerUnit = 0.085
					fuelconsumption.FuelConsumed = BackupFuelConsumed / 3
					fuelconsumption.ConvertedFuelConsumed = fuelconsumption.FuelConsumed * 1000
					fuelconsumption.FuelCost = fuelconsumption.FuelCostPerUnit * fuelconsumption.ConvertedFuelConsumed

					fuelconsumptionArray = append(fuelconsumptionArray, fuelconsumption)

					val.BackupFuelTotalCost += fuelconsumption.FuelCost
				} else if strings.ToLower(strings.Trim(val.Plant, " ")) == "Rabigh" {
					fuelconsumption.IsPrimaryFuel = false
					fuelconsumption.FuelType = "CRUDE"
					fuelconsumption.FuelCost = 0.1
					fuelconsumption.FuelConsumed = BackupFuelConsumed / 3
					fuelconsumption.ConvertedFuelConsumed = fuelconsumption.FuelConsumed * 1000

					fuelconsumptionArray = append(fuelconsumptionArray, fuelconsumption)

					val.BackupFuelTotalCost += fuelconsumption.FuelCost

					fuelconsumption.IsPrimaryFuel = false
					fuelconsumption.FuelType = "DIESEL"
					fuelconsumption.FuelCost = 0.085
					fuelconsumption.FuelConsumed = BackupFuelConsumed / 3
					fuelconsumption.ConvertedFuelConsumed = fuelconsumption.FuelConsumed * 1000

					fuelconsumptionArray = append(fuelconsumptionArray, fuelconsumption)

					val.BackupFuelTotalCost += fuelconsumption.FuelCost
				}
				//#endregion
			} else {
				//#region not hfo
				fuelconsumption.IsPrimaryFuel = false
				fuelconsumption.FuelType = BackupFuelType

				fuelconsumption.FuelConsumed = BackupFuelConsumed

				if strings.Contains(strings.ToLower(fuelconsumption.FuelType), "gas") {
					fuelconsumption.FuelCostPerUnit = 2.813
					fuelconsumption.ConvertedFuelConsumed = fuelconsumption.FuelConsumed * 0.0353
				} else if strings.Contains(strings.ToLower(fuelconsumption.FuelType), "crude") {
					fuelconsumption.FuelCostPerUnit = 0.1
				} else if strings.Contains(strings.ToLower(fuelconsumption.FuelType), "diesel") {
					fuelconsumption.FuelCostPerUnit = 0.085
				}

				if !strings.Contains(strings.ToLower(fuelconsumption.FuelType), "gas") {
					fuelconsumption.ConvertedFuelConsumed = fuelconsumption.FuelConsumed * 1000
				}

				fuelconsumption.FuelCost = fuelconsumption.FuelCostPerUnit * fuelconsumption.ConvertedFuelConsumed

				fuelconsumptionArray = append(fuelconsumptionArray, fuelconsumption)

				val.BackupFuelTotalCost += fuelconsumption.FuelCost
				//#endregion
			}
			//#endregion
			//var totaldieselconsumed float64

			fuelconsumptionFilter := crowd.From(&fuelconsumptionArray).Where(func(x interface{}) interface{} {
				return strings.ToLower(strings.Trim(x.(ValueEquationFuel).FuelType, " ")) == "diesel"
			}).Exec().Result.Data().([]ValueEquationFuel)

			totaldieselconsumed := crowd.From(&fuelconsumptionFilter).Sum(func(x interface{}) interface{} {
				return x.(ValueEquationFuel).ConvertedFuelConsumed
			}).Exec().Result.Sum

			val.FuelTransportCost = TransportCosts * totaldieselconsumed
			val.TotalFuelCost = val.PrimaryFuelTotalCost + val.BackupFuelTotalCost
			val.OperatingCost = val.FuelTransportCost + val.TotalFuelCost
			//#endregion

			//#region Maintenance
			val.TotalLabourCost = crowd.From(&phases).Sum(func(x interface{}) interface{} {
				return x.(tk.M).GetFloat64("skilledlabour") + x.(tk.M).GetFloat64("unskilledlabour")
			}).Exec().Result.Sum

			val.TotalMaterialCost = crowd.From(&phases).Sum(func(x interface{}) interface{} {
				return x.(tk.M).GetFloat64("materials")
			}).Exec().Result.Sum

			val.TotalServicesCost = crowd.From(&phases).Sum(func(x interface{}) interface{} {
				return x.(tk.M).GetFloat64("contractmaintenance")
			}).Exec().Result.Sum

			details := []ValueEquationDetails{}

			tempGT := crowd.From(&lists).Where(func(x interface{}) interface{} {
				return x.(tk.M).GetString("unit") == tempunit
			}).Exec().Result.Data().([]tk.M)

			if len(tempGT) > 0 {
				for _, gt := range tempGT {
					det := ValueEquationDetails{}
					det.DataSource = "Paper Records"
					det.WorkOrderType = gt.GetString("wotype")
					det.LaborCost = gt.GetFloat64("skilledlabour") + gt.GetFloat64("unskilledlabour")
					det.MaterialCost = gt.GetFloat64("materials")
					det.ServiceCost = gt.GetFloat64("contractmaintenance")

					details = append(details, det)
				}
			}

			tempsyn := crowd.From(&syn).Where(func(x interface{}) interface{} {
				unitDB := x.(tk.M).GetString("unit")
				return strings.Replace(strings.Replace(strings.Replace(strings.Replace(unitDB, "GT0", "", -1), "GT", "", -1), "ST0", "S", -1), "ST", "S", -1) == strings.Replace(strings.Replace(strings.Replace(strings.Replace(unit, "GT0", "", -1), "GT", "", -1), "ST0", "S", -1), "ST", "S", -1)
			}).Exec().Result.Data().([]tk.M)

			if len(tempsyn) > 0 {
				for _, pm := range tempsyn {
					det := ValueEquationDetails{}
					det.DataSource = "SAP PM"
					det.WorkOrderType = pm.GetString("wotype")
					det.LaborCost = pm.GetFloat64("plannedlaborcost")
					det.MaterialCost = pm.GetFloat64("actualmaterialcost")
					det.ServiceCost = 0

					details = append(details, det)

					val.TotalLabourCost += pm.GetFloat64("plannedlaborcost")
					val.TotalMaterialCost += pm.GetFloat64("actualmaterialcost")
				}
			}

			tempbrowser := crowd.From(&databrowser).Where(func(x interface{}) interface{} {
				isturbine := x.(tk.M).Get("isturbine").(bool)
				turbineInfos := x.(tk.M).GetString("tinfshortname") != ""
				tiShortName := strings.Replace(strings.Replace(strings.Replace(x.(tk.M).GetString("shortname"), "GTO", "", -1), "GT", "", -1), "ST0", "ST", -1)
				unitTemp := strings.Replace(strings.Replace(strings.Replace(unit, "GTO", "", -1), "GT", "", -1), "ST0", "ST", -1)

				return isturbine && turbineInfos && tiShortName == unitTemp
			}).Exec().Result.Data().([]tk.M)

			veTop10s := []ValueEquationTop10{}

			if len(tempbrowser) > 0 {
				sintax := "select distinct(OrderType) from MaintenanceCost"
				csr, e = c.NewQuery().Command("freequery", tk.M{}.Set("syntax", sintax)).Cursor(nil)
				MROTypes := []tk.M{}
				e = csr.Fetch(&MROTypes, 0, false)
				csr.Close()

				for _, types := range MROTypes {
					det := ValueEquationDetails{}
					det.DataSource = "SAP PM"
					det.WorkOrderType = types.GetString("ordertype")

					actualDuration := []tk.M{}
					csr, e = c.NewQuery().Command("procedure", tk.M{}.Set("name", "spDataBrowserGetActualDuration").Set("parms", tk.M{}.Set("@WOType", det.WorkOrderType))).Cursor(nil)
					e = csr.Fetch(&actualDuration, 0, false)
					csr.Close()
					det.Duration = crowd.From(&actualDuration).Sum(func(x interface{}) interface{} {
						return x.(tk.M).GetFloat64("actualduration")
					}).Exec().Result.Sum

					query = append(query[0:0], dbox.Eq("OrderType", det.WorkOrderType))
					csr, e = c.NewQuery().From(new(MaintenanceCost).TableName()).Select("InternalLaborActual", "DirectMaterialActual", "InternalMaterialActual", "ExternalServiceActual").Cursor(nil)
					maintenanceCostData := []tk.M{}
					e = csr.Fetch(&maintenanceCostData, 0, false)

					det.LaborCost = crowd.From(&maintenanceCostData).Sum(func(x interface{}) interface{} {
						return x.(tk.M).GetFloat64("internallaboractual")
					}).Exec().Result.Sum

					det.MaterialCost = crowd.From(&maintenanceCostData).Sum(func(x interface{}) interface{} {
						return x.(tk.M).GetFloat64("directmaterialactual") + x.(tk.M).GetFloat64("internalmaterialactual")
					}).Exec().Result.Sum

					det.ServiceCost = crowd.From(&maintenanceCostData).Sum(func(x interface{}) interface{} {
						return x.(tk.M).GetFloat64("externalserviceactual")
					}).Exec().Result.Sum

					details = append(details, det)

					val.TotalLabourCost += det.LaborCost
					val.TotalMaterialCost += det.MaterialCost
					val.TotalServicesCost += det.ServiceCost
					val.TotalDuration += det.Duration
				}

				sintax = "select distinct(MaintenanceOrder) from MaintenanceCost"
				csr, e = c.NewQuery().Command("freequery", tk.M{}.Set("syntax", sintax)).Cursor(nil)
				maintenanceOrders := []tk.M{}
				e = csr.Fetch(&maintenanceOrders, 0, false)
				csr.Close()

				for _, fl := range maintenanceOrders {
					query = append(query[0:0], dbox.Eq("MaintenanceOrder", fl.GetString("maintenanceorder")))
					csr, e = c.NewQuery().From(new(MaintenanceCost).TableName()).Where(query...).Cursor(nil)
					db := []tk.M{}
					csr.Fetch(&db, 0, false)
					csr.Close()

					top10 := ValueEquationTop10{}
					top10.WorkOrderID = db[0].GetString("maintenanceOrder")
					top10.WorkOrderType = db[0].GetString("ordertype")
					top10.EquipmentType = db[0].GetString("equipmenttype")
					actualDuration := []tk.M{}
					csr, e = c.NewQuery().Command("procedure", tk.M{}.Set("name", "spDataBrowserGetActualDuration").Set("parms", tk.M{}.Set("@WOType", fl.GetString("maintenanceorder")))).Cursor(nil)
					e = csr.Fetch(&actualDuration, 0, false)
					csr.Close()
					top10.Duration = crowd.From(&actualDuration).Sum(func(x interface{}) interface{} {
						return x.(tk.M).GetFloat64("actualduration")
					}).Exec().Result.Sum
					top10.LaborCost = crowd.From(&db).Sum(func(x interface{}) interface{} {
						return x.(tk.M).GetFloat64("internallaboractual")
					}).Exec().Result.Sum
					top10.MaterialCost = crowd.From(&db).Sum(func(x interface{}) interface{} {
						return x.(tk.M).GetFloat64("internalmaterialactual") + x.(tk.M).GetFloat64("directmaterialactual")
					}).Exec().Result.Sum
					top10.ServiceCost = crowd.From(&db).Sum(func(x interface{}) interface{} {
						return x.(tk.M).GetFloat64("externalserviceactual")
					}).Exec().Result.Sum

					veTop10s = append(veTop10s, top10)

				}
				//#endregion
			}

			val.MaintenanceCost = val.TotalLabourCost + val.TotalMaterialCost + val.TotalServicesCost

			//#endregion

			val.ValueEquationCost = val.Revenue - val.OperatingCost - val.MaintenanceCost

			id, _ := ctx.InsertOut(val)

			if len(details) > 0 {
				for _, data := range details {
					data.VEId = id

					_, e = ctx.InsertOut(&data)
				}
			}

			if len(veTop10s) > 0 {
				for _, data := range details {
					data.VEId = id

					_, e = ctx.InsertOut(&data)
				}
			}
		}
	}

	return e
}

func (d *GenValueEquation) generateValueEquationAllMaintenance(Year int, Plants []string) error {
	ctx := d.BaseController.Ctx
	c := ctx.Connection
	var (
		query []*dbox.Filter
		e     error
	)

	YearFirst := strconv.Itoa(Year) + "-01-01"
	YearLast := strconv.Itoa(Year+1) + "-01-01"

	for _, Plant := range Plants {
		query = append(query, dbox.Eq("Plant", Plant))
		csr, _ := c.NewQuery().From(new(PerformanceFactors).TableName()).Where(query...).Cursor(nil)
		pfs := []tk.M{}
		e = csr.Fetch(&pfs, 0, false)
		csr.Close()

		csr, _ = c.NewQuery().From(new(Consolidated).TableName()).Where(query...).Cursor(nil)
		cons := []tk.M{}
		e = csr.Fetch(&cons, 0, false)
		csr.Close()

		query = append(query, dbox.And(dbox.Gte("DatePerformed", YearFirst), dbox.Lt("DatePerformed", YearLast)))
		csr, _ = c.NewQuery().From(new(PrevMaintenanceValueEquation).TableName()).Where(query...).Cursor(nil)
		lists := []tk.M{}
		e = csr.Fetch(&lists, 0, false)
		csr.Close()

		if Plant == "Qurayyah" || Plant == "Qurayyah CC" {
			query = append(query[0:0], dbox.And(dbox.Eq("Plant", "Qurayyah"), dbox.Eq("Year", Year)))
		} else {
			query = append(query[0:0], dbox.And(dbox.Eq("Plant", Plant), dbox.Eq("Year", Year)))
		}

		csr, _ = c.NewQuery().From(new(PowerPlantOutages).TableName()).Where(query...).Cursor(nil)
		outages := []tk.M{}
		e = csr.Fetch(&outages, 0, false)
		csr.Close()

		query = append(query[0:0], dbox.And(dbox.Eq("Plant", Plant), dbox.Eq("Year", Year)))
		csr, _ = c.NewQuery().From(new(StartupPaymentAndPenalty).TableName()).Where(query...).Cursor(nil)
		start := []tk.M{}
		e = csr.Fetch(&start, 0, false)
		csr.Close()

		query = append(query[0:0], dbox.And(dbox.Eq("Plant", Plant), dbox.Eq("Year", Year)))
		csr, _ = c.NewQuery().From(new(FuelCost).TableName()).Where(query...).Cursor(nil)
		fuelcosts := []tk.M{}
		e = csr.Fetch(&fuelcosts, 0, false)
		csr.Close()

		query = append(query[0:0], dbox.Eq("Plant", Plant))
		query = append(query, dbox.And(dbox.Gte("ScheduledStart", YearFirst), dbox.Lt("ScheduledStart", YearLast)))
		csr, _ = c.NewQuery().From(new(SyntheticPM).TableName()).Where(query...).Cursor(nil)
		syn := []tk.M{}
		e = csr.Fetch(&syn, 0, false)
		csr.Close()

		query = append(query[0:0], dbox.And(dbox.Eq("Plant", Plant), dbox.Eq("Year", Year)))
		csr, _ = c.NewQuery().From(new(FuelTransport).TableName()).Where(query...).Cursor(nil)
		trans := []tk.M{}
		e = csr.Fetch(&trans, 0, false)
		csr.Close()

		sintax := "select * from DataBrowser inner join PowerPlantCoordinates on DataBrowser.PlantCode = PowerPlantCoordinates.PlantCode where PeriodYear = " + strconv.Itoa(Year) + " and PowerPlantCoordinates.PlantName = '" + Plant + "'"
		csr, e = c.NewQuery().Command("freequery", tk.M{}.Set("syntax", sintax)).Cursor(nil)
		databrowser := []tk.M{}
		e = csr.Fetch(&databrowser, 0, false)
		csr.Close()

		csr, _ = c.NewQuery().From(new(GenerationAppendix).TableName()).Cursor(nil)
		genA := []tk.M{}
		e = csr.Fetch(&genA, 0, false)
		csr.Close()

		csr, _ = c.NewQuery().From(new(Availability).TableName()).Cursor(nil)
		avail := []tk.M{}
		e = csr.Fetch(&avail, 0, false)
		csr.Close()

		csr, _ = c.NewQuery().From(new(UnitPower).TableName()).Cursor(nil)
		unitpower := []tk.M{}
		e = csr.Fetch(&unitpower, 0, false)
		csr.Close()

		UnitsData := crowd.From(&fuelcosts).Group(func(x interface{}) interface{} {
			unitId := x.(tk.M).GetString("unitid")
			return strings.Replace(strings.TrimSpace(unitId), " ", "", -1)
		}, nil).Exec().Result.Data().([]crowd.KV)

		var Units []string
		for _, unit := range UnitsData {
			Units = append(Units, unit.Key.(string))
		}

		DieselData := crowd.From(&fuelcosts).Where(func(x interface{}) interface{} {
			return x.(tk.M).GetString("primaryfueltype") == "DIESEL"
		}).Exec().Result.Data().([]tk.M)

		DieselConsumptions := 0.0
		if len(DieselData) > 0 {
			DieselConsumptions = crowd.From(&DieselData).Sum(func(x interface{}) interface{} {
				return x.(tk.M).GetFloat64("primaryfuelconsumed")
			}).Exec().Result.Sum

			DieselConsumptions = DieselConsumptions * 1000
		}

		TransportCosts := 0.0
		if DieselConsumptions != 0.0 {
			TransportCosts = trans[0].GetFloat64("transportcost") / DieselConsumptions
		}

		UnitsTemp := crowd.From(&Units).Where(func(x interface{}) interface{} {
			return !strings.Contains(x.(string), "CS")
		}).Exec().Result.Data().([]string)

		for _, unit := range UnitsTemp {
			var NormalizedUnit string

			if len(unit) < 3 {
				if Plant == "PP9" {
					NormalizedUnit = "GT" + unit
				}
			} else {
				NormalizedUnit = strings.Replace(strings.Replace(strings.Replace(strings.Replace(unit, ".", "", -1), " ", "", -1), "GT0", "GT", -1), "ST0", "ST", -1)
			}

			tempunit := strings.Replace(strings.Replace(NormalizedUnit, ".", "", -1), " ", "", -1)

			if len(tempunit) == 3 && !strings.Contains(tempunit, "ST") {
				tempunit = "GT0" + strings.Replace(tempunit, "GT", "", -1)
			}

			val := new(ValueEquation)
			val.Plant = Plant
			val.Dates = time.Date(Year, 1, 1, 0, 0, 0, 0, time.UTC)
			val.Month = 1
			val.Year = Year
			val.Unit = strings.Replace(strings.Replace(NormalizedUnit, ".", "", -1), " ", "", -1)
			val.UnitGroup = val.Unit[0:2]

			tempLists := crowd.From(&lists).Where(func(x interface{}) interface{} {
				return x.(tk.M).GetString("unit") == tempunit
			}).Exec().Result.Data().([]tk.M)

			if len(tempLists) > 0 {
				val.Phase = tempLists[0].GetString("phase")
			}

			tempCons := crowd.From(&cons).Where(func(x interface{}) interface{} {
				return strings.Replace(x.(tk.M).GetString("unit"), "ST0", "ST", -1) == strings.Replace(val.Unit, "ST0", "ST", -1)
			}).Exec().Result.Data().([]tk.M)

			val.Capacity = crowd.From(&tempCons).Sum(func(x interface{}) interface{} {
				return x.(tk.M).GetFloat64("capacity")
			}).Exec().Result.Sum

			val.NetGeneration = crowd.From(&tempCons).Sum(func(x interface{}) interface{} {
				return x.(tk.M).GetFloat64("energynet")
			}).Exec().Result.Sum

			val.AvgNetGeneration = crowd.From(&tempCons).Avg(func(x interface{}) interface{} {
				return x.(tk.M).GetFloat64("energynet")
			}).Exec().Result.Avg

			if Plant == "PP9" || Plant == "Qurayyah" || Plant == "Qurayyah CC" {
				tempAvail := crowd.From(&avail).Where(func(x interface{}) interface{} {
					return x.(tk.M).GetString("powerplant") == Plant && x.(tk.M).GetString("turbine") == strings.Replace(strings.Replace(val.Unit, "GT0", "GT", -1), "ST0", "ST", -1)
				}).Exec().Result.Data().([]tk.M)
				if len(tempAvail) > 0 {
					val.PrctWAF = tempAvail[0].GetFloat64("prctwaf")
					val.PrctWUF = tempAvail[0].GetFloat64("prctwuf")
				}
			} else if Plant == "Rabigh" {
				if strings.Contains(val.Unit, "GT") {
					tempAvail := crowd.From(&avail).Where(func(x interface{}) interface{} {
						return strings.Contains(x.(tk.M).GetString("powerplant"), Plant) && x.(tk.M).GetString("turbine") == strings.Replace(strings.Replace(val.Unit, "GT0", "GT", -1), "ST0", "ST", -1)
					}).Exec().Result.Data().([]tk.M)

					if len(tempAvail) > 0 {
						val.PrctWAF = tempAvail[0].GetFloat64("prctwaf")
						val.PrctWUF = tempAvail[0].GetFloat64("prctwuf")
					}
				} else if strings.Contains(val.Unit, "ST") {
					tempAvail := crowd.From(&avail).Where(func(x interface{}) interface{} {
						return x.(tk.M).GetString("powerplant") == "Rabigh Steam" && x.(tk.M).GetString("turbine") == strings.Replace(strings.Replace(val.Unit, "GT0", "GT", -1), "ST0", "ST", -1)
					}).Exec().Result.Data().([]tk.M)

					if len(tempAvail) > 0 {
						val.PrctWAF = tempAvail[0].GetFloat64("prctwaf")
						val.PrctWUF = tempAvail[0].GetFloat64("prctwuf")
					}
				}
			} else if Plant == "Shoaiba" || Plant == "Ghazlan" {
				tempAvail := crowd.From(&avail).Where(func(x interface{}) interface{} {
					return strings.Contains(x.(tk.M).GetString("powerplant"), Plant) && x.(tk.M).GetString("turbine") == strings.Replace(strings.Replace(val.Unit, "GT0", "GT", -1), "ST0", "ST", -1)
				}).Exec().Result.Data().([]tk.M)

				if len(tempAvail) > 0 {
					val.PrctWAF = tempAvail[0].GetFloat64("prctwaf")
					val.PrctWUF = tempAvail[0].GetFloat64("prctwuf")
				}
			}

			//#endregion

			//#region Revenue
			tempAppendix := []tk.M{}
			if strings.Contains(val.Unit, "ST") {
				if Plant == "Qurayyah" {
					tempAppendix = crowd.From(&genA).Where(func(x interface{}) interface{} {
						return x.(tk.M).GetString("plant") == "QRPP"
					}).Exec().Result.Data().([]tk.M)
				} else if Plant == "Qurayyah CC" {
					tempAppendix = crowd.From(&genA).Where(func(x interface{}) interface{} {
						return x.(tk.M).GetString("plant") == "QCCP"
					}).Exec().Result.Data().([]tk.M)
				} else if Plant == "Ghazlan" {
					unittemp, _ := strconv.Atoi(strings.Replace(val.Unit, "ST", "", -1))
					if unittemp <= 4 {
						tempAppendix = crowd.From(&genA).Where(func(x interface{}) interface{} {
							return x.(tk.M).GetString("plant") == "Ghazlan I (1-4)"
						}).Exec().Result.Data().([]tk.M)
					} else if unittemp <= 8 {
						tempAppendix = crowd.From(&genA).Where(func(x interface{}) interface{} {
							return x.(tk.M).GetString("plant") == "Ghazlan II (5-8)"
						}).Exec().Result.Data().([]tk.M)
					}
				} else if Plant == "Shoaiba" {
					tempAppendix = crowd.From(&genA).Where(func(x interface{}) interface{} {
						return x.(tk.M).GetString("plant") == "Shoaiba Steam"
					}).Exec().Result.Data().([]tk.M)
				} else if Plant == "Rabigh" {
					tempAppendix = crowd.From(&genA).Where(func(x interface{}) interface{} {
						return x.(tk.M).GetString("plant") == "Rabigh Steam"
					}).Exec().Result.Data().([]tk.M)
				} else if Plant == "PP9" {
					tempAppendix = crowd.From(&genA).Where(func(x interface{}) interface{} {
						return x.(tk.M).GetString("plant") == "PP9 CC"
					}).Exec().Result.Data().([]tk.M)
				}
			} else {
				if Plant == "Qurayyah" {
					tempAppendix = crowd.From(&genA).Where(func(x interface{}) interface{} {
						return x.(tk.M).GetString("plant") == "QRPP"
					}).Exec().Result.Data().([]tk.M)
				} else if Plant == "Qurayyah CC" {
					tempAppendix = crowd.From(&genA).Where(func(x interface{}) interface{} {
						return x.(tk.M).GetString("plant") == "QCCP"
					}).Exec().Result.Data().([]tk.M)
				} else if Plant == "Rabigh" {
					unittemp, _ := strconv.Atoi(strings.Replace(val.Unit, "GT", "", -1))
					if unittemp <= 12 {
						tempAppendix = crowd.From(&genA).Where(func(x interface{}) interface{} {
							return x.(tk.M).GetString("plant") == "Rabigh Combined"
						}).Exec().Result.Data().([]tk.M)
					} else if unittemp <= 40 {
						tempAppendix = crowd.From(&genA).Where(func(x interface{}) interface{} {
							return x.(tk.M).GetString("plant") == "Rabigh Gas" && x.(tk.M).GetFloat64("units") == 28
						}).Exec().Result.Data().([]tk.M)
					}
				} else if Plant == "PP9" {
					unittemp, _ := strconv.Atoi(strings.Replace(val.Unit, "GT", "", -1))
					if unittemp <= 16 {
						tempAppendix = crowd.From(&genA).Where(func(x interface{}) interface{} {
							return x.(tk.M).GetString("plant") == "PP9 CC"
						}).Exec().Result.Data().([]tk.M)
					} else if unittemp <= 24 {
						tempAppendix = crowd.From(&genA).Where(func(x interface{}) interface{} {
							return x.(tk.M).GetString("plant") == "PPEXT" && x.(tk.M).GetFloat64("units") == 8
						}).Exec().Result.Data().([]tk.M)
					} else if unittemp <= 56 {
						tempAppendix = crowd.From(&genA).Where(func(x interface{}) interface{} {
							return x.(tk.M).GetString("plant") == "PPEXT" && x.(tk.M).GetFloat64("units") == 32
						}).Exec().Result.Data().([]tk.M)
					}
				}
			}

			totalDays := (time.Date(Year, 12, 31, 0, 0, 0, 0, time.UTC).Sub(time.Date(Year-1, 12, 31, 0, 0, 0, 0, time.UTC)).Seconds()) / 86400
			if len(tempAppendix) > 0 {
				val.CapacityPayment = tempAppendix[0].GetFloat64("contractedcapacity") * (tempAppendix[0].GetFloat64("fomr") + tempAppendix[0].GetFloat64("ccr")) * totalDays * 10
			}

			tempCons = crowd.From(&cons).Where(func(x interface{}) interface{} {
				return strings.Replace(x.(tk.M).GetString("unit"), "STO", "ST", -1) == strings.Replace(val.Unit, "STO", "ST", -1)
			}).Exec().Result.Data().([]tk.M)

			val.EnergyPayment = crowd.From(&tempCons).Sum(func(x interface{}) interface{} {
				return x.(tk.M).GetFloat64("energynet")
			}).Exec().Result.Sum * tempAppendix[0].GetFloat64("vomr") * 10

			val.VOMR = tempAppendix[0].GetFloat64("vomr")

			tempPfs := crowd.From(&pfs).Where(func(x interface{}) interface{} {
				return x.(tk.M).GetString("unit") == strings.Replace(val.Unit, "ST0", "ST", -1)
			}).Exec().Result.Data().([]tk.M)

			if len(tempPfs) > 0 {
				val.SRF = tempPfs[0].GetFloat64("srf")
			}

			if Plant == "Rabigh" {
				if len(outages) > 0 {
					sintax := "select count(*) as Count from PowerPlantOutagesDetails inner join PowerPlantOutages on PowerPlantOutagesDetails.POId = PowerPlantOutages.Id where PowerPlantOutagesDetails.UnitNo = '" + val.Unit + "' and OutageType != 'PO' and PowerPlantOutages.Plant = 'Rabigh Steam'"
					csr, e = c.NewQuery().Command("freequery", tk.M{}.Set("syntax", sintax)).Cursor(nil)
					count := []tk.M{}
					e = csr.Fetch(&count, 0, false)
					csr.Close()

					if len(count) > 0 {
						val.UnplannedOutages = count[0].GetFloat64("count")
					}

					sintax = "select * from PowerPlantOutagesDetails inner join PowerPlantOutages on PowerPlantOutagesDetails.POId = PowerPlantOutages.Id where PowerPlantOutagesDetails.UnitNo = '" + val.Unit + "' and PowerPlantOutages.Plant = 'Rabigh Steam'"
					csr, e = c.NewQuery().Command("freequery", tk.M{}.Set("syntax", sintax)).Cursor(nil)
					tempOutages := []tk.M{}
					e = csr.Fetch(&tempOutages, 0, false)
					csr.Close()

					if len(tempOutages) > 0 {
						val.TotalOutageDuration = crowd.From(&tempOutages).Sum(func(x interface{}) interface{} {
							return x.(tk.M).GetFloat64("totalhours")
						}).Exec().Result.Sum
					}
				}
			} else if Plant == "Qurayyah" || Plant == "Qurayyah CC" {
				if len(outages) > 0 {
					sintax := "select count(*) as Count from PowerPlantOutagesDetails inner join PowerPlantOutages on PowerPlantOutagesDetails.POId = PowerPlantOutages.Id where PowerPlantOutagesDetails.UnitNo = '" + val.Unit + "' and OutageType != 'PO' and PowerPlantOutages.Plant = '" + Plant + "'"
					csr, e = c.NewQuery().Command("freequery", tk.M{}.Set("syntax", sintax)).Cursor(nil)
					count := []tk.M{}
					e = csr.Fetch(&count, 0, false)
					csr.Close()

					if len(count) > 0 {
						val.UnplannedOutages = count[0].GetFloat64("count")
					}

					sintax = "select * from PowerPlantOutagesDetails inner join PowerPlantOutages on PowerPlantOutagesDetails.POId = PowerPlantOutages.Id where PowerPlantOutagesDetails.UnitNo = '" + val.Unit + "' and PowerPlantOutages.Plant = '" + Plant + "'"
					csr, e = c.NewQuery().Command("freequery", tk.M{}.Set("syntax", sintax)).Cursor(nil)
					tempOutages := []tk.M{}
					e = csr.Fetch(&tempOutages, 0, false)
					csr.Close()

					if len(tempOutages) > 0 {
						val.TotalOutageDuration = crowd.From(&tempOutages).Sum(func(x interface{}) interface{} {
							return x.(tk.M).GetFloat64("totalhours")
						}).Exec().Result.Sum
					}
				}
			} else {
				if len(outages) > 0 {
					sintax := "select count(*) as Count from PowerPlantOutagesDetails inner join PowerPlantOutages on PowerPlantOutagesDetails.POId = PowerPlantOutages.Id where PowerPlantOutagesDetails.UnitNo = '" + val.Unit + "' and OutageType != 'PO'"
					csr, e = c.NewQuery().Command("freequery", tk.M{}.Set("syntax", sintax)).Cursor(nil)
					count := []tk.M{}
					e = csr.Fetch(&count, 0, false)
					csr.Close()

					if len(count) > 0 {
						val.UnplannedOutages = count[0].GetFloat64("count")
					}

					sintax = "select * from PowerPlantOutagesDetails inner join PowerPlantOutages on PowerPlantOutagesDetails.POId = PowerPlantOutages.Id where PowerPlantOutagesDetails.UnitNo = '" + val.Unit + "'"
					csr, e = c.NewQuery().Command("freequery", tk.M{}.Set("syntax", sintax)).Cursor(nil)
					tempOutages := []tk.M{}
					e = csr.Fetch(&tempOutages, 0, false)
					csr.Close()

					if len(tempOutages) > 0 {
						val.TotalOutageDuration = crowd.From(&tempOutages).Sum(func(x interface{}) interface{} {
							return x.(tk.M).GetFloat64("totalhours")
						}).Exec().Result.Sum
					}
				}
			}

			if val.SRF == 100 {
				tempStart := crowd.From(&start).Where(func(x interface{}) interface{} {
					return strings.Replace(x.(tk.M).GetString("unit"), "ST0", "ST", -1)
				}).Exec().Result.Data().([]tk.M)

				if len(tempStart) > 0 {
					val.StartupPayment = tempStart[0].GetFloat64("startuppayment")
					val.PenaltyAmount = 0
				}
			} else {
				val.StartupPayment = 0
				if len(tempAppendix) > 0 {
					val.PenaltyAmount = tempAppendix[0].GetFloat64("deduct")
				}
			}

			val.PenaltyAmount += tempAppendix[0].GetFloat64("deduct") * val.UnplannedOutages
			val.Incentive = 0
			val.Revenue = val.CapacityPayment + val.EnergyPayment + val.Incentive + val.StartupPayment - val.PenaltyAmount
			//#endregion
			//#region OperatingCost
			//#region Primary Fuel
			valueequationFuels := []ValueEquationFuel{}
			tempFuelCosts := crowd.From(&fuelcosts).Where(func(x interface{}) interface{} {
				return strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(x.(tk.M).GetString("unitid"), " ", "", -1), ".", "", -1), "ST0", "ST", -1), "GT0", "", -1), "GT", "", -1) == strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(val.Unit, " ", "", -1), ".", "", -1), "ST0", "ST", -1), "GT0", "", -1), "GT", "", -1)
			}).Exec().Result.Data().([]tk.M)

			PrimaryFuelType := tempFuelCosts[0].GetString("primaryfueltype")
			if strings.TrimSpace(strings.ToLower(PrimaryFuelType)) == "hfo" {
				//#region hfo
				PrimaryFuelConsumed := crowd.From(&tempFuelCosts).Sum(func(x interface{}) interface{} {
					return x.(tk.M).GetFloat64("primaryfuelconsumed")
				}).Exec().Result.Sum

				if strings.TrimSpace(strings.ToLower(val.Plant)) == "shoaiba" {
					fuelconsumption := ValueEquationFuel{}
					fuelconsumption.IsPrimaryFuel = true
					fuelconsumption.FuelType = "CRUDE"
					fuelconsumption.FuelCostPerUnit = 0.1
					fuelconsumption.FuelConsumed = PrimaryFuelConsumed / 3
					fuelconsumption.ConvertedFuelConsumed = fuelconsumption.FuelConsumed * 1000
					fuelconsumption.FuelCost = fuelconsumption.FuelCostPerUnit * fuelconsumption.ConvertedFuelConsumed
					valueequationFuels = append(valueequationFuels, fuelconsumption)

					val.PrimaryFuelTotalCost += fuelconsumption.FuelCost

					fuelconsumption.IsPrimaryFuel = true
					fuelconsumption.FuelType = "CRUDE HEAVY"
					fuelconsumption.FuelCostPerUnit = 0.049
					fuelconsumption.FuelConsumed = PrimaryFuelConsumed / 3
					fuelconsumption.ConvertedFuelConsumed = fuelconsumption.FuelConsumed * 1000
					fuelconsumption.FuelCost = fuelconsumption.FuelCostPerUnit * fuelconsumption.ConvertedFuelConsumed
					valueequationFuels = append(valueequationFuels, fuelconsumption)

					val.PrimaryFuelTotalCost += fuelconsumption.FuelCost

					fuelconsumption.IsPrimaryFuel = true
					fuelconsumption.FuelType = "DIESEL"
					fuelconsumption.FuelCostPerUnit = 0.085
					fuelconsumption.FuelConsumed = PrimaryFuelConsumed / 3
					fuelconsumption.ConvertedFuelConsumed = fuelconsumption.FuelConsumed * 1000
					fuelconsumption.FuelCost = fuelconsumption.FuelCostPerUnit * fuelconsumption.ConvertedFuelConsumed
					valueequationFuels = append(valueequationFuels, fuelconsumption)

					val.PrimaryFuelTotalCost += fuelconsumption.FuelCost
				} else if strings.TrimSpace(strings.ToLower(val.Plant)) == "Rabigh" {
					fuelconsumption := ValueEquationFuel{}
					fuelconsumption.IsPrimaryFuel = true
					fuelconsumption.FuelType = "CRUDE"
					fuelconsumption.FuelCostPerUnit = 0.1
					fuelconsumption.FuelConsumed = PrimaryFuelConsumed / 3
					fuelconsumption.ConvertedFuelConsumed = fuelconsumption.FuelConsumed * 1000
					fuelconsumption.FuelCost = fuelconsumption.FuelCostPerUnit * fuelconsumption.ConvertedFuelConsumed
					valueequationFuels = append(valueequationFuels, fuelconsumption)

					val.PrimaryFuelTotalCost += fuelconsumption.FuelCost

					fuelconsumption.IsPrimaryFuel = true
					fuelconsumption.FuelType = "DIESEL"
					fuelconsumption.FuelCostPerUnit = 0.085
					fuelconsumption.FuelConsumed = PrimaryFuelConsumed / 3
					fuelconsumption.ConvertedFuelConsumed = fuelconsumption.FuelConsumed * 1000
					fuelconsumption.FuelCost = fuelconsumption.FuelCostPerUnit * fuelconsumption.ConvertedFuelConsumed
					valueequationFuels = append(valueequationFuels, fuelconsumption)

					val.PrimaryFuelTotalCost += fuelconsumption.FuelCost
				}
				//#endregion
			} else {
				fuelconsumption := ValueEquationFuel{}
				fuelconsumption.IsPrimaryFuel = true
				fuelconsumption.FuelType = tempFuelCosts[0].GetString("primaryfueltype")
				if strings.Contains(strings.ToLower(fuelconsumption.FuelType), "gas") {
					fuelconsumption.FuelCostPerUnit = 2.813
				} else if strings.Contains(strings.ToLower(fuelconsumption.FuelType), "crude") {
					fuelconsumption.FuelCostPerUnit = 0.1
				} else if strings.Contains(strings.ToLower(fuelconsumption.FuelType), "diesel") {
					fuelconsumption.FuelCostPerUnit = 0.085
				}

				fuelconsumption.FuelConsumed = crowd.From(&tempFuelCosts).Sum(func(x interface{}) interface{} {
					return x.(tk.M).GetFloat64("primaryfuelconsumed")
				}).Exec().Result.Sum

				if strings.Contains(strings.ToLower(fuelconsumption.FuelType), "gas") {
					fuelconsumption.ConvertedFuelConsumed = fuelconsumption.FuelConsumed * 0.0353
				} else {
					fuelconsumption.ConvertedFuelConsumed = fuelconsumption.FuelConsumed * 1000
				}

				fuelconsumption.FuelCost = fuelconsumption.FuelCostPerUnit * fuelconsumption.ConvertedFuelConsumed
				valueequationFuels = append(valueequationFuels, fuelconsumption)

				val.PrimaryFuelTotalCost += fuelconsumption.FuelCost

				if Plant == "Qurayyah" {
					fuelconsumption := ValueEquationFuel{}
					fuelconsumption.IsPrimaryFuel = true
					fuelconsumption.FuelType = tempFuelCosts[0].GetString("primary2fueltype")

					if strings.Contains(strings.ToLower(fuelconsumption.FuelType), "gas") {
						fuelconsumption.FuelCostPerUnit = 2.813
					} else if strings.Contains(strings.ToLower(fuelconsumption.FuelType), "crude") {
						fuelconsumption.FuelCostPerUnit = 0.1
					} else if strings.Contains(strings.ToLower(fuelconsumption.FuelType), "diesel") {
						fuelconsumption.FuelCostPerUnit = 0.085
					}

					fuelconsumption.FuelConsumed = crowd.From(&tempFuelCosts).Sum(func(x interface{}) interface{} {
						return x.(tk.M).GetFloat64("primary2fuelconsumed")
					}).Exec().Result.Sum

					if strings.Contains(strings.ToLower(fuelconsumption.FuelType), "gas") {
						fuelconsumption.ConvertedFuelConsumed = fuelconsumption.FuelConsumed * 0.0353
					} else {
						fuelconsumption.ConvertedFuelConsumed = fuelconsumption.FuelConsumed * 1000
					}

					fuelconsumption.FuelCost = fuelconsumption.FuelCostPerUnit * fuelconsumption.ConvertedFuelConsumed
					valueequationFuels = append(valueequationFuels, fuelconsumption)

					val.PrimaryFuelTotalCost += fuelconsumption.FuelCost
				}
				//#endregion
			}
			//#endregion
			//#region backup fuel
			BackupFuelType := tempFuelCosts[0].GetString("backupfueltype")

			if strings.TrimSpace(strings.ToLower(BackupFuelType)) == "hfo" {
				//#region hfo
				BackupFuelConsumed := crowd.From(&tempFuelCosts).Sum(func(x interface{}) interface{} {
					return x.(tk.M).GetFloat64("backupfuelconsumed")
				}).Exec().Result.Sum

				if strings.TrimSpace(strings.ToLower(val.Plant)) == "shoaiba" {
					fuelconsumption := ValueEquationFuel{}
					fuelconsumption.IsPrimaryFuel = false
					fuelconsumption.FuelType = "CRUDE"
					fuelconsumption.FuelCostPerUnit = 0.1
					fuelconsumption.FuelConsumed = BackupFuelConsumed / 3
					fuelconsumption.ConvertedFuelConsumed = fuelconsumption.FuelConsumed * 1000
					fuelconsumption.FuelCost = fuelconsumption.FuelCostPerUnit * fuelconsumption.ConvertedFuelConsumed

					valueequationFuels = append(valueequationFuels, fuelconsumption)
					val.BackupFuelTotalCost += fuelconsumption.FuelCost

					fuelconsumption = ValueEquationFuel{}
					fuelconsumption.IsPrimaryFuel = false
					fuelconsumption.FuelType = "CRUDE HEAVY"
					fuelconsumption.FuelCostPerUnit = 0.049
					fuelconsumption.FuelConsumed = BackupFuelConsumed / 3
					fuelconsumption.ConvertedFuelConsumed = fuelconsumption.FuelConsumed * 1000
					fuelconsumption.FuelCost = fuelconsumption.FuelCostPerUnit * fuelconsumption.ConvertedFuelConsumed

					valueequationFuels = append(valueequationFuels, fuelconsumption)
					val.BackupFuelTotalCost += fuelconsumption.FuelCost

					fuelconsumption = ValueEquationFuel{}
					fuelconsumption.IsPrimaryFuel = false
					fuelconsumption.FuelType = "DIESEL"
					fuelconsumption.FuelCostPerUnit = 0.085
					fuelconsumption.FuelConsumed = BackupFuelConsumed / 3
					fuelconsumption.ConvertedFuelConsumed = fuelconsumption.FuelConsumed * 1000
					fuelconsumption.FuelCost = fuelconsumption.FuelCostPerUnit * fuelconsumption.ConvertedFuelConsumed

					valueequationFuels = append(valueequationFuels, fuelconsumption)
					val.BackupFuelTotalCost += fuelconsumption.FuelCost
				} else if strings.TrimSpace(strings.ToLower(val.Plant)) == "rabigh" {
					fuelconsumption := ValueEquationFuel{}
					fuelconsumption.IsPrimaryFuel = false
					fuelconsumption.FuelType = "CRUDE"
					fuelconsumption.FuelCost = 0.1
					fuelconsumption.FuelConsumed = BackupFuelConsumed / 3
					fuelconsumption.ConvertedFuelConsumed = fuelconsumption.FuelConsumed * 1000

					valueequationFuels = append(valueequationFuels, fuelconsumption)
					val.BackupFuelTotalCost += fuelconsumption.FuelCost

					fuelconsumption = ValueEquationFuel{}
					fuelconsumption.IsPrimaryFuel = false
					fuelconsumption.FuelType = "DIESEL"
					fuelconsumption.FuelCost = 0.085
					fuelconsumption.FuelConsumed = BackupFuelConsumed / 3
					fuelconsumption.ConvertedFuelConsumed = fuelconsumption.FuelConsumed * 1000
					fuelconsumption.FuelCost = fuelconsumption.FuelCostPerUnit * fuelconsumption.ConvertedFuelConsumed

					valueequationFuels = append(valueequationFuels, fuelconsumption)
					val.BackupFuelTotalCost += fuelconsumption.FuelCost
				}
				//#endregion
			} else {
				//#region not hfo
				fuelconsumption := ValueEquationFuel{}
				fuelconsumption.IsPrimaryFuel = false
				fuelconsumption.FuelType = tempFuelCosts[0].GetString("backupfueltype")

				if strings.Contains(strings.ToLower(fuelconsumption.FuelType), "gas") {
					fuelconsumption.FuelCostPerUnit = 2.813
				} else if strings.Contains(strings.ToLower(fuelconsumption.FuelType), "crude") {
					fuelconsumption.FuelCostPerUnit = 0.1
				} else if strings.Contains(strings.ToLower(fuelconsumption.FuelType), "diesel") {
					fuelconsumption.FuelCostPerUnit = 0.085
				}

				fuelconsumption.FuelConsumed = crowd.From(&tempFuelCosts).Sum(func(x interface{}) interface{} {
					return x.(tk.M).GetFloat64("backupfuelconsumed")
				}).Exec().Result.Sum

				if strings.Contains(strings.ToLower(fuelconsumption.FuelType), "gas") {
					fuelconsumption.ConvertedFuelConsumed = fuelconsumption.FuelConsumed * 0.0353
				} else {
					fuelconsumption.ConvertedFuelConsumed = fuelconsumption.FuelConsumed * 1000
				}

				fuelconsumption.FuelCost = fuelconsumption.FuelCostPerUnit * fuelconsumption.ConvertedFuelConsumed

				valueequationFuels = append(valueequationFuels, fuelconsumption)
				val.BackupFuelTotalCost += fuelconsumption.FuelCost
				//#endregion
			}
			//#endregion
			totaldieselconsumed := 0.0
			tempValueEquationFuels := crowd.From(&valueequationFuels).Where(func(x interface{}) interface{} {
				return strings.TrimSpace(strings.ToLower(x.(ValueEquationFuel).FuelType)) == "diesel"
			}).Exec().Result.Data().([]ValueEquationFuel)

			totaldieselconsumed = crowd.From(&tempValueEquationFuels).Sum(func(x interface{}) interface{} {
				return x.(ValueEquationFuel).ConvertedFuelConsumed
			}).Exec().Result.Sum

			val.FuelTransportCost = TransportCosts * totaldieselconsumed
			val.TotalFuelCost = val.PrimaryFuelTotalCost + val.BackupFuelTotalCost
			val.OperatingCost = val.FuelTransportCost + val.TotalFuelCost

			//#endregion
			//#region Maintenance
			tempLists = crowd.From(&lists).Where(func(x interface{}) interface{} {
				return x.(tk.M).GetString("unit") == tempunit
			}).Exec().Result.Data().([]tk.M)

			val.TotalLabourCost = crowd.From(&tempLists).Sum(func(x interface{}) interface{} {
				return x.(tk.M).GetFloat64("skilledlabour") + x.(tk.M).GetFloat64("unskilledlabour")
			}).Exec().Result.Sum

			val.TotalMaterialCost = crowd.From(&tempLists).Sum(func(x interface{}) interface{} {
				return x.(tk.M).GetFloat64("materials")
			}).Exec().Result.Sum

			val.TotalServicesCost = crowd.From(&tempLists).Sum(func(x interface{}) interface{} {
				return x.(tk.M).GetFloat64("ContractMaintenance")
			}).Exec().Result.Sum

			details := []ValueEquationDetails{}
			top10s := []ValueEquationTop10{}

			tempGT := crowd.From(&lists).Where(func(x interface{}) interface{} {
				return x.(tk.M).GetString("unit") == tempunit
			}).Exec().Result.Data().([]tk.M)

			if len(tempGT) > 0 {
				for _, gt := range tempGT {
					det := ValueEquationDetails{}
					det.DataSource = "Paper Records"
					det.WorkOrderType = gt.GetString("wotype")
					det.LaborCost = gt.GetFloat64("skilledlabour") + gt.GetFloat64("unskilledlabour")
					det.MaterialCost = gt.GetFloat64("materials")
					det.ServiceCost = gt.GetFloat64("contractmaintenance")

					details = append(details, det)

					top10 := ValueEquationTop10{}
					top10.WorkOrderID = gt.GetString("id")
					top10.WorkOrderDescription = gt.GetString("description")
					top10.EquipmentType = "Not Available"
					top10.EquipmentTypeDescription = "Not Available"
					top10.MaintenanceActivity = gt.GetString("description")
					top10.Duration = gt.GetFloat64("days") * 24
					top10.LaborCost = det.LaborCost
					top10.MaterialCost = det.MaterialCost
					top10.ServiceCost = det.ServiceCost
					top10.MaintenanceCost = top10.LaborCost + top10.MaterialCost + top10.ServiceCost

					top10s = append(top10s, top10)
				}
			}

			tempbrowser := crowd.From(&databrowser).Where(func(x interface{}) interface{} {
				//isTurbine, _ :=
				return x.(tk.M).Get("isturbine").(bool) && strings.Replace(strings.Replace(strings.Replace(x.(tk.M).GetString("tinfshortname"), "GT0", "", -1), "GT", "", -1), "ST0", "ST", -1) == strings.Replace(strings.Replace(strings.Replace(val.Unit, "GT0", "", -1), "GT", "", -1), "ST0", "ST", -1)
			}).Exec().Result.Data().([]tk.M)

			var databrowse []interface{}

			if len(tempbrowser) > 0 {
				tempDataBrowser := crowd.From(&databrowser).Where(func(x interface{}) interface{} {
					return x.(tk.M).GetString("turbineparent") == tempbrowser[0].GetString("functionallocation") || x.(tk.M).GetString("functionallocation") == tempbrowser[0].GetString("functionallocation")
				}).Exec().Result.Data().([]tk.M)

				tempDataBrowse := crowd.From(&tempDataBrowser).Group(func(x interface{}) interface{} {
					return strings.TrimSpace(x.(tk.M).GetString("functionallocation"))
				}, nil).Exec().Result.Data().([]crowd.KV)

				if len(tempDataBrowse) > 0 {
					for _, brow := range tempDataBrowse {
						databrowse = append(databrowse, brow.Key.(string))
					}
				}
			}

			tempWoList := []tk.M{}

			if len(databrowse) > 0 {
				query = append(query[0:0], dbox.In("FunctionalLocation", databrowse...))
				csr, e = c.NewQuery().From("WOList").Where(query...).Cursor(nil)
				e = csr.Fetch(&tempWoList, 0, false)
			}

			tempWoList1 := []crowd.KV{}

			if len(tempWoList) > 0 {
				tempWoList1 = crowd.From(&tempWoList).Group(func(x interface{}) interface{} {
					return x.(tk.M).GetString("ordercode")
				}, nil).Exec().Result.Data().([]crowd.KV)
			}

			MaintenanceOrderList := []string{}
			if len(tempWoList1) > 0 {
				for _, wo := range tempWoList1 {
					MaintenanceOrderList = append(MaintenanceOrderList, wo.Key.(string))
				}
			}

			tempsyn := crowd.From(&syn).Where(func(x interface{}) interface{} {
				if len(MaintenanceOrderList) > 0 {
					tempMain := crowd.From(&MaintenanceOrderList).Where(func(y interface{}) interface{} {
						return strings.Contains(y.(string), x.(tk.M).GetString("woid"))
					}).Exec().Result.Data().([]string)

					if len(tempMain) > 0 {
						if val.Unit != "" {
							return strings.Replace(strings.Replace(strings.Replace(strings.Replace(x.(tk.M).GetString("unit"), "GT0", "", -1), "GT", "", -1), "ST0", "S", -1), "ST", "S", -1) == strings.Replace(strings.Replace(strings.Replace(strings.Replace(val.Unit, "GT0", "", -1), "GT", "", -1), "ST0", "S", -1), "ST", "S", -1)
						} else {
							return false
						}
					} else {
						return false
					}
				} else {
					return false
				}
			}).Exec().Result.Data().([]tk.M)

			if len(tempsyn) > 0 {
				for _, pm := range tempsyn {
					det := ValueEquationDetails{}
					det.DataSource = "SAP PM"
					det.WorkOrderType = pm.GetString("wotype")
					det.LaborCost = pm.GetFloat64("plannedlaborcost")
					det.MaterialCost = pm.GetFloat64("actualmaterialcost")
					det.ServiceCost = 0

					details = append(details, det)

					val.TotalLabourCost += pm.GetFloat64("plannedlaborcost")
					val.TotalMaterialCost += pm.GetFloat64("actualmaterialcost")
				}
			}

			if len(tempbrowser) > 0 {
				query = append(query[0:0], dbox.And(dbox.In("MaintenanceOrder", MaintenanceOrderList), dbox.Ne("MaintenanceOrder", ""), dbox.Eq("Period", YearFirst)))
				csr, e = c.NewQuery().From("MaintenanceCost").Where(query...).Cursor(nil)
				maintCost := []tk.M{}
				e = csr.Fetch(&maintCost, 0, false)
				csr.Close()

				query = append(query[0:0], dbox.And(dbox.In("MaintenanceOrder", MaintenanceOrderList), dbox.Ne("MaintenanceOrder", ""), dbox.Eq("Period", YearFirst)))
				csr, e = c.NewQuery().From("MaintenanceCostByHour").Where(query...).Cursor(nil)
				maintHour := []tk.M{}
				e = csr.Fetch(&maintHour, 0, false)
				csr.Close()
				log.Println(maintHour)
				//IList<MaintenanceCost> maintCost = DataHelper.Populate<MaintenanceCost>("MaintenanceCost", Query.And(Query.In("MaintenanceOrder", new BsonArray(MaintenanceOrderList)), Query.NE("MaintenanceOrder",string.Empty))).Where(x => x.Period.Value.Year == Year).ToList();
			}
		}
	}
	return e
}
