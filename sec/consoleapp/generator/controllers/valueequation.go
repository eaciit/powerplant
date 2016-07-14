package controllers

import (
	"github.com/eaciit/crowd"
	"github.com/eaciit/dbox"
	// . "github.com/eaciit/powerplant/sec/consoleapp/generator/helpers"
	"log"
	"strconv"
	"strings"
	"time"

	. "github.com/eaciit/powerplant/sec/library/models"
	tk "github.com/eaciit/toolkit"
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

	e = d.generateValueEquation()
	if e != nil {
		tk.Println(e)
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

func (d *GenValueEquation) generateValueEquation() error {
	var e error
	Year := 2014
	YearFirst := strconv.Itoa(Year) + "-01-01 00:00:00.000"
	YearLast := strconv.Itoa(Year+1) + "-01-01 00:00:00.000"

	Plant := "Qurayyah CC"

	ctx := d.BaseController.Ctx
	c := ctx.Connection

	query := []*dbox.Filter{}

	query = append(query, dbox.Eq("Plant", Plant))
	csr, e := c.NewQuery().From(new(PerformanceFactors).TableName()).Where(query...).Cursor(nil)

	if e != nil {
		return e
	} else {
		defer csr.Close()
	}

	pfs := []tk.M{}
	e = csr.Fetch(&pfs, 0, false)

	csr1, e := c.NewQuery().From(new(Consolidated).TableName()).Where(query...).Cursor(nil)

	if e != nil {
		return e
	} else {
		defer csr1.Close()
	}

	cons := []tk.M{}
	e = csr1.Fetch(&cons, 0, false)

	query = append(query, dbox.And(dbox.Gte("DatePerformed", YearFirst), dbox.Lt("DatePerformed", YearLast)))
	csr2, e := c.NewQuery().From(new(PrevMaintenanceValueEquation).TableName()).Where(query...).Cursor(nil)

	if e != nil {
		return e
	} else {
		defer csr2.Close()
	}

	lists := []tk.M{}
	e = csr2.Fetch(&lists, 0, false)

	query = nil
	query = append(query, dbox.Eq("Plant", Plant))
	query = append(query, dbox.Eq("Year", Year))
	csr3, e := c.NewQuery().From(new(PowerPlantOutages).TableName()).Where(query...).Cursor(nil)

	if e != nil {
		return e
	} else {
		defer csr3.Close()
	}

	outages := []tk.M{}
	e = csr3.Fetch(&outages, 0, false)

	csr4, e := c.NewQuery().From(new(StartupPaymentAndPenalty).TableName()).Where(query...).Cursor(nil)

	if e != nil {
		return e
	} else {
		defer csr4.Close()
	}

	start := []tk.M{}
	e = csr4.Fetch(&start, 0, false)

	csr5, e := c.NewQuery().From(new(FuelCost).TableName()).Where(query...).Cursor(nil)

	if e != nil {
		return e
	} else {
		defer csr5.Close()
	}

	fuelcosts := []tk.M{}
	e = csr5.Fetch(&fuelcosts, 0, false)

	query = nil
	query = append(query, dbox.Eq("Plant", Plant))
	query = append(query, dbox.And(dbox.Gte("ScheduledStart", YearFirst), dbox.Lt("ScheduledStart", YearLast)))
	csr6, e := c.NewQuery().From(new(SyntheticPM).TableName()).Where(query...).Cursor(nil)

	if e != nil {
		return e
	} else {
		defer csr6.Close()
	}

	syn := []tk.M{}
	e = csr6.Fetch(&syn, 0, false)

	query = nil
	query = append(query, dbox.Eq("Plant", Plant))
	query = append(query, dbox.Eq("Year", Year))
	csr7, e := c.NewQuery().From(new(FuelTransport).TableName()).Where(query...).Cursor(nil)

	if e != nil {
		return e
	} else {
		defer csr7.Close()
	}

	trans := []tk.M{}
	e = csr7.Fetch(&trans, 0, false)

	sintax := "select * from DataBrowser inner join PowerPlantCoordinates on DataBrowser.PlantCode = PowerPlantCoordinates.PlantCode where PeriodYear = " + strconv.Itoa(Year) + " and PowerPlantCoordinates.PlantName = '" + Plant + "'"
	csr8, e := c.NewQuery().Command("freequery", tk.M{}.Set("syntax", sintax)).Cursor(nil)

	if e != nil {
		return e
	} else {
		defer csr8.Close()
	}

	databrowser := []tk.M{}
	e = csr8.Fetch(&databrowser, 0, false)

	csr9, e := c.NewQuery().From(new(GenerationAppendix).TableName()).Cursor(nil)

	if e != nil {
		return e
	} else {
		defer csr9.Close()
	}

	genA := []tk.M{}
	e = csr9.Fetch(&genA, 0, false)

	Units := crowd.From(&pfs).Group(func(x interface{}) interface{} {
		return x.(tk.M).GetString("unit")
	}, nil).Exec().Result.Data().([]crowd.KV)

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
			tempunit := unit.Key.(string)
			if len(tempunit) == 3 && strings.Contains(tempunit, "ST") {
				tempunit = "GT0" + strings.Replace(tempunit, "GT", "", -1)
			}

			val := new(ValueEquation)
			val.Plant = Plant
			val.Dates = time.Date(2014, 1, 1, 0, 0, 0, 0, time.UTC)
			val.Month = 1
			val.Year = 2014
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
					query = nil
					query = append(query, dbox.Eq("POId", POId))
					query = append(query, dbox.Eq("UnitNo", tempunit))
					query = append(query, dbox.Ne("OutageType", "PO"))
					query = append(query, dbox.Eq("PlantName", "Rabigh Steam"))
					csr10, e := c.NewQuery().From(new(StartupPaymentAndPenalty).TableName()).Where(query...).Cursor(nil)

					if e != nil {
						return e
					} else {
						defer csr10.Close()
					}

					val.UnplannedOutages = float64(csr10.Count())
				}
			} else if Plant == "Qurayyah" || Plant == "Qurayyah CC" {
				if len(outages) == 0 {
					val.UnplannedOutages = 0
				} else {
					POId := outages[0].GetString("id")
					query = nil
					query = append(query, dbox.Eq("POId", POId))
					query = append(query, dbox.Eq("UnitNo", tempunit))
					query = append(query, dbox.Ne("OutageType", "PO"))
					query = append(query, dbox.Eq("PlantName", Plant))
					csr11, e := c.NewQuery().From(new(StartupPaymentAndPenalty).TableName()).Where(query...).Cursor(nil)

					if e != nil {
						return e
					} else {
						defer csr11.Close()
					}

					val.UnplannedOutages = float64(csr11.Count())
				}
			} else {
				if len(outages) == 0 {
					val.UnplannedOutages = 0
				} else {
					POId := outages[0].GetString("id")
					query = nil
					query = append(query, dbox.Eq("POId", POId))
					query = append(query, dbox.Eq("UnitNo", tempunit))
					query = append(query, dbox.Ne("OutageType", "PO"))
					csr12, e := c.NewQuery().From(new(StartupPaymentAndPenalty).TableName()).Where(query...).Cursor(nil)

					if e != nil {
						return e
					} else {
						defer csr12.Close()
					}

					val.UnplannedOutages = float64(csr12.Count())
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
				log.Println(val.PrimaryFuelTotalCost)
			}
			//endregion

			//region backup fuel
		}
	}

	return e
}
