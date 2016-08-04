package controllers

import (
	"github.com/eaciit/crowd"
	"github.com/eaciit/dbox"
	"github.com/eaciit/orm"
	. "github.com/eaciit/powerplant/sec/consoleapp/generator/helpers"
	. "github.com/eaciit/powerplant/sec/library/models"
	tk "github.com/eaciit/toolkit"
	// "log"
	"strconv"
	"strings"
	"time"
)

// GenValueEquationDashboard ...
type GenValueEquationDashboard struct {
	*BaseController
}

// Generate ...
func (d *GenValueEquationDashboard) Generate(base *BaseController) {
	var (
		e error
	)
	if base != nil {
		d.BaseController = base
	}

	// Plant := []string{"Qurayyah CC", "Qurayyah", "Rabigh", "Ghazlan", "PP9", "Shoaiba"}
	// Year := []int{2011, 2012, 2013, 2014, 2015}
	// e = d.generateValueEquationAllMaintenanceRedoDashboard(Year, Plant)

	Years := []int{2011, 2012, 2013, 2014, 2015}
	e = d.generateOtherMaintenanceData(Years)

	if e != nil {
		tk.Println(e)
	}
	tk.Println("##Value Equation DashboardData : DONE\n")
}

func (d *GenValueEquationDashboard) generateValueEquationAllMaintenanceRedoDashboard(Years []int, Plants []string) error {
	ctx := d.BaseController.Ctx
	c := ctx.Connection

	var (
		query []*dbox.Filter
		e     error
	)

	for _, Year := range Years {
		YearFirst := strconv.Itoa(Year) + "-01-01"
		YearLast := strconv.Itoa(Year+1) + "-01-01"

		for _, Plant := range Plants {
			query = append(query[0:0], dbox.And(dbox.Eq("Plant", Plant), dbox.Eq("Year", Year)))
			csr, _ := c.NewQuery().From(new(PerformanceFactors).TableName()).Where(query...).Cursor(nil)
			pfs := []tk.M{}
			e = csr.Fetch(&pfs, 0, false)
			csr.Close()

			query = append(query[0:0], dbox.Eq("Plant", Plant))
			query = append(query, dbox.And(dbox.Gte("ConsolidatedDate", YearFirst), dbox.Lt("ConsolidatedDate", YearLast)))
			csr, _ = c.NewQuery().From(new(Consolidated).TableName()).Where(query...).Cursor(nil)
			cons := []tk.M{}
			e = csr.Fetch(&cons, 0, false)
			csr.Close()

			query = append(query[0:0], dbox.Eq("Plant", Plant))
			query = append(query, dbox.And(dbox.Gte("DatePerformed", YearFirst), dbox.Lt("DatePerformed", YearLast)))
			csr, _ = c.NewQuery().From(new(PrevMaintenanceValueEquation).TableName()).Where(query...).Cursor(nil)
			lists := []tk.M{}
			e = csr.Fetch(&lists, 0, false)
			csr.Close()

			outages := []tk.M{}

			if Plant == "Qurayyah" || Plant == "Qurayyah CC" {
				query = append(query[0:0], dbox.Eq("Plant", "Qurayyah"))
				query = append(query, dbox.Eq("Year", Year))
			} else {
				query = append(query[0:0], dbox.Eq("Plant", Plant))
				query = append(query, dbox.Eq("Year", Year))
			}

			csr, e = c.NewQuery().From(new(PowerPlantOutages).TableName()).Where(query...).Cursor(nil)
			e = csr.Fetch(&outages, 0, false)
			csr.Close()

			query = append(query[0:0], dbox.Eq("Plant", Plant))
			query = append(query, dbox.Eq("Year", Year))
			csr, e = c.NewQuery().From(new(FuelCost).TableName()).Where(query...).Cursor(nil)
			fuelcosts := []tk.M{}
			e = csr.Fetch(&fuelcosts, 0, false)
			csr.Close()

			csr, e = c.NewQuery().From(new(FuelTransport).TableName()).Where(query...).Cursor(nil)
			trans := []tk.M{}
			e = csr.Fetch(&trans, 0, false)
			csr.Close()

			query = append(query[0:0], dbox.Eq("Plant", Plant))
			query = append(query, dbox.And(dbox.Gte("ScheduledStart", YearFirst), dbox.Lt("ScheduledStart", YearLast)))
			csr, _ = c.NewQuery().From(new(SyntheticPM).TableName()).Where(query...).Cursor(nil)
			syn := []tk.M{}
			e = csr.Fetch(&syn, 0, false)
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

			units := crowd.From(&fuelcosts).Group(func(x interface{}) interface{} {
				return strings.TrimSpace(strings.Replace(x.(tk.M).GetString("unitid"), " ", "", -1))
			}, nil).Exec().Result.Data().([]crowd.KV)

			Units := []string{}
			for _, unit := range units {
				Units = append(Units, unit.Key.(string))
			}

			tempFuelCosts1 := crowd.From(&fuelcosts).Where(func(x interface{}) interface{} {
				return x.(tk.M).GetString("primaryfueltype") == "DIESEL"
			}).Exec().Result.Data().([]tk.M)
			tempFuelCosts2 := crowd.From(&fuelcosts).Where(func(x interface{}) interface{} {
				return x.(tk.M).GetString("backupfueltype") == "DIESEL"
			}).Exec().Result.Data().([]tk.M)

			DieselConsumptions := crowd.From(&tempFuelCosts1).Sum(func(x interface{}) interface{} {
				return x.(tk.M).GetFloat64("primaryfuelconsumed")
			}).Exec().Result.Sum + crowd.From(&tempFuelCosts2).Sum(func(x interface{}) interface{} {
				return x.(tk.M).GetFloat64("backupfuelconsumed")
			}).Exec().Result.Sum*1000

			TransportCosts := 0.0
			if DieselConsumptions != 0 || len(trans) > 0 {
				TransportCosts = trans[0].GetFloat64("transportcost") / DieselConsumptions
			}

			UnitsList := crowd.From(&Units).Where(func(x interface{}) interface{} {
				return !strings.Contains(x.(string), "CS")
			}).Exec().Result.Data().([]string)

			for _, unit := range UnitsList {
				NormalizedUnit := ""
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

				val := new(ValueEquationDashboard)
				val.Plant = Plant
				val.Dates = time.Date(Year, 1, 1, 0, 0, 0, 0, time.UTC)
				val.Month = 1
				val.Year = Year
				val.Unit = strings.Replace(strings.Replace(NormalizedUnit, ".", "", -1), " ", "", -1)
				val.UnitGroup = val.Unit[0:2]

				query = append(query[0:0], dbox.Eq("Plant", val.Plant))
				csr, _ := c.NewQuery().From("GeneralInfo").Where(query...).Cursor(nil)
				infoList := []tk.M{}
				e = csr.Fetch(&infoList, 0, false)
				infos := crowd.From(&infoList).Where(func(x interface{}) interface{} {
					return strings.Replace(x.(tk.M).GetString("unit"), "GT", "", -1) == strings.Replace(val.Unit, "GT", "", -1)
				}).Exec().Result.Data().([]tk.M)

				_ = infos

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

				if len(tempCons) > 0 {

					val.AvgNetGeneration = crowd.From(&tempCons).Avg(func(x interface{}) interface{} {
						return x.(tk.M).GetFloat64("energynet")
					}).Exec().Result.Avg
				}

				// #region Availability
				if Plant == "PP9" || Plant == "Qurayyah" || Plant == "Qurayyah CC" {
					tempAvail := crowd.From(&avail).Where(func(x interface{}) interface{} {
						return x.(tk.M).GetString("plant") == Plant && x.(tk.M).GetString("turbine") == strings.Replace(strings.Replace(val.Unit, "GT0", "GT", -1), "ST0", "ST", -1)
					}).Exec().Result.Data().([]tk.M)
					if len(tempAvail) > 0 {
						val.PrctWAF = tempAvail[0].GetFloat64("prctwaf")
						val.PrctWUF = tempAvail[0].GetFloat64("prctwuf")
					}
				} else if Plant == "Rabigh" {
					if strings.Contains(val.Unit, "GT") {
						tempAvail := crowd.From(&avail).Where(func(x interface{}) interface{} {
							return strings.Contains(x.(tk.M).GetString("plant"), Plant) && x.(tk.M).GetString("turbine") == strings.Replace(strings.Replace(val.Unit, "GT0", "GT", -1), "ST0", "ST", -1)
						}).Exec().Result.Data().([]tk.M)
						if len(tempAvail) > 0 {
							val.PrctWAF = tempAvail[0].GetFloat64("prctwaf")
							val.PrctWUF = tempAvail[0].GetFloat64("prctwuf")
						}
					} else if strings.Contains(val.Unit, "GT") {
						tempAvail := crowd.From(&avail).Where(func(x interface{}) interface{} {
							return x.(tk.M).GetString("plant") == "Rabigh Steam" && x.(tk.M).GetString("turbine") == strings.Replace(strings.Replace(val.Unit, "GT0", "GT", -1), "ST0", "ST", -1)
						}).Exec().Result.Data().([]tk.M)
						if len(tempAvail) > 0 {
							val.PrctWAF = tempAvail[0].GetFloat64("prctwaf")
							val.PrctWUF = tempAvail[0].GetFloat64("prctwuf")
						}
					}
				} else if Plant == "Shoaiba" || Plant == "Ghazlan" {
					tempAvail := crowd.From(&avail).Where(func(x interface{}) interface{} {
						return strings.Contains(x.(tk.M).GetString("plant"), Plant) && x.(tk.M).GetString("turbine") == strings.Replace(strings.Replace(val.Unit, "GT0", "GT", -1), "ST0", "ST", -1)
					}).Exec().Result.Data().([]tk.M)
					if len(tempAvail) > 0 {
						val.PrctWAF = tempAvail[0].GetFloat64("prctwaf")
						val.PrctWUF = tempAvail[0].GetFloat64("prctwuf")
					}
				}
				// #endregion
				// #region Revenue
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
						unittemp, _ := strconv.Atoi(strings.Replace(val.Unit, "ST", "", -1))

						if unittemp <= 4 {
							tempappendix = crowd.From(&genA).Where(func(x interface{}) interface{} {
								return x.(tk.M).GetString("plant") == "Ghazlan I (1-4)"
							}).Exec().Result.Data().([]tk.M)
						} else if unittemp <= 8 {
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
							return x.(tk.M).GetString("plant") == "Rabigh Steam"
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
						unittemp, _ := strconv.Atoi(strings.Replace(val.Unit, "GT", "", -1))

						if unittemp <= 12 {
							tempappendix = crowd.From(&genA).Where(func(x interface{}) interface{} {
								return x.(tk.M).GetString("plant") == "Rabigh Combined"
							}).Exec().Result.Data().([]tk.M)
						} else if unittemp <= 40 {
							tempappendix = crowd.From(&genA).Where(func(x interface{}) interface{} {
								return x.(tk.M).GetString("plant") == "Rabigh Gas" && x.(tk.M).GetFloat64("units") == 28
							}).Exec().Result.Data().([]tk.M)
						}
					} else if Plant == "PP9" {
						unittemp, _ := strconv.Atoi(strings.Replace(val.Unit, "GT", "", -1))

						if unittemp <= 16 {
							tempappendix = crowd.From(&genA).Where(func(x interface{}) interface{} {
								return x.(tk.M).GetString("plant") == "PP9 CC"
							}).Exec().Result.Data().([]tk.M)
						} else if unittemp <= 24 {
							tempappendix = crowd.From(&genA).Where(func(x interface{}) interface{} {
								return x.(tk.M).GetString("plant") == "PPEXT" && x.(tk.M).GetFloat64("units") == 8
							}).Exec().Result.Data().([]tk.M)
						} else if unittemp <= 56 {
							tempappendix = crowd.From(&genA).Where(func(x interface{}) interface{} {
								return x.(tk.M).GetString("plant") == "PPEXT" && x.(tk.M).GetFloat64("units") == 32
							}).Exec().Result.Data().([]tk.M)
						}
					}

					if len(tempappendix) > 0 {
						apendixResult := tempappendix[0].GetFloat64("contractedcapacity") * (tempappendix[0].GetFloat64("fomr") + tempappendix[0].GetFloat64("ccr"))
						val.CapacityPayment = apendixResult * 365 * 10
						val.VOMR = tempappendix[0].GetFloat64("vomr")
					}

					if len(cons) > 0 {
						consResult := crowd.From(&cons).Where(func(x interface{}) interface{} {
							return strings.Replace(x.(tk.M).GetString("unit"), "ST0", "ST", -1) == strings.Replace(val.Unit, "ST0", "ST", -1)
						}).Exec().Result.Data().([]tk.M)

						val.EnergyPayment = crowd.From(&consResult).Sum(func(x interface{}) interface{} {
							return x.(tk.M).GetFloat64("energynet") * tempappendix[0].GetFloat64("vomr") * 10
						}).Exec().Result.Sum
					}

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
						if len(tempappendix) > 0 {
							val.StartupPayment = tempappendix[0].GetFloat64("startup")
						}
						val.PenaltyAmount = 0
					} else {
						val.StartupPayment = 0
						if len(tempappendix) > 0 {
							val.PenaltyAmount = tempappendix[0].GetFloat64("deduct")
						}
					}

					if len(tempappendix) > 0 {
						val.PenaltyAmount += tempappendix[0].GetFloat64("deduct") * val.UnplannedOutages
					}

					val.Incentive = 0
					val.Revenue = val.CapacityPayment + val.EnergyPayment + val.Incentive + val.StartupPayment - val.PenaltyAmount
					// #endregion

					// #region OperatingCost
					// #region Primary Fuel
					Fuels := []VEDFuel{}
					tempFuelCosts := crowd.From(&fuelcosts).Where(func(x interface{}) interface{} {
						return strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(x.(tk.M).GetString("unitid"), " ", "", -1), ".", "", -1), "ST0", "ST", -1), "GT0", "", -1), "GT", "", -1) == strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(val.Unit, " ", "", -1), ".", "", -1), "ST0", "ST", -1), "GT0", "", -1), "GT", "", -1)
					}).Exec().Result.Data().([]tk.M)

					PrimaryFuelType := ""
					if len(tempFuelCosts) > 0 {
						PrimaryFuelType = tempFuelCosts[0].GetString("primaryfueltype")
					}

					if strings.TrimSpace(strings.ToLower(PrimaryFuelType)) == "hfo" {
						// #region hfo
						PrimaryFuelConsumed := crowd.From(&tempFuelCosts).Sum(func(x interface{}) interface{} {
							return x.(tk.M).GetFloat64("primaryfuelconsumed")
						}).Exec().Result.Sum

						if strings.TrimSpace(strings.ToLower(val.Plant)) == "shoaiba" {
							fuelconsumption := VEDFuel{}
							fuelconsumption.IsPrimaryFuel = true
							fuelconsumption.FuelType = "CRUDE"
							fuelconsumption.FuelCostPerUnit = 0.1
							fuelconsumption.FuelConsumed = PrimaryFuelConsumed / 3
							fuelconsumption.ConvertedFuelConsumed = fuelconsumption.FuelConsumed * 1000
							fuelconsumption.FuelCost = fuelconsumption.FuelCostPerUnit * fuelconsumption.ConvertedFuelConsumed
							Fuels = append(Fuels, fuelconsumption)
							val.PrimaryFuelTotalCost += fuelconsumption.FuelCost

							fuelconsumption = VEDFuel{}
							fuelconsumption.IsPrimaryFuel = true
							fuelconsumption.FuelType = "CRUDE HEAVY"
							fuelconsumption.FuelCostPerUnit = 0.049
							fuelconsumption.FuelConsumed = PrimaryFuelConsumed / 3
							fuelconsumption.ConvertedFuelConsumed = fuelconsumption.FuelConsumed * 1000
							fuelconsumption.FuelCost = fuelconsumption.FuelCostPerUnit * fuelconsumption.ConvertedFuelConsumed
							Fuels = append(Fuels, fuelconsumption)
							val.PrimaryFuelTotalCost += fuelconsumption.FuelCost

							fuelconsumption = VEDFuel{}
							fuelconsumption.IsPrimaryFuel = true
							fuelconsumption.FuelType = "DIESEL"
							fuelconsumption.FuelCostPerUnit = 0.085
							fuelconsumption.FuelConsumed = PrimaryFuelConsumed / 3
							fuelconsumption.ConvertedFuelConsumed = fuelconsumption.FuelConsumed * 1000
							fuelconsumption.FuelCost = fuelconsumption.FuelCostPerUnit * fuelconsumption.ConvertedFuelConsumed
							Fuels = append(Fuels, fuelconsumption)
							val.PrimaryFuelTotalCost += fuelconsumption.FuelCost
						} else if strings.TrimSpace(strings.ToLower(Plant)) == "rabigh" {

							fuelconsumption := VEDFuel{}
							fuelconsumption.IsPrimaryFuel = true
							fuelconsumption.FuelType = "CRUDE"
							fuelconsumption.FuelCostPerUnit = 0.1
							fuelconsumption.FuelConsumed = PrimaryFuelConsumed / 3
							fuelconsumption.ConvertedFuelConsumed = fuelconsumption.FuelConsumed * 1000
							fuelconsumption.FuelCost = fuelconsumption.FuelCostPerUnit * fuelconsumption.ConvertedFuelConsumed
							Fuels = append(Fuels, fuelconsumption)
							val.PrimaryFuelTotalCost += fuelconsumption.FuelCost

							fuelconsumption = VEDFuel{}
							fuelconsumption.IsPrimaryFuel = true
							fuelconsumption.FuelType = "DIESEL"
							fuelconsumption.FuelCostPerUnit = 0.085
							fuelconsumption.FuelConsumed = PrimaryFuelConsumed / 3
							fuelconsumption.ConvertedFuelConsumed = fuelconsumption.FuelConsumed * 1000
							fuelconsumption.FuelCost = fuelconsumption.FuelCostPerUnit * fuelconsumption.ConvertedFuelConsumed
							Fuels = append(Fuels, fuelconsumption)
							val.PrimaryFuelTotalCost += fuelconsumption.FuelCost
						}
						//#endregion
					} else {
						//#region not hfo
						fuelconsumption := VEDFuel{}
						fuelconsumption.IsPrimaryFuel = true
						fuelconsumption.FuelType = PrimaryFuelType
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
						Fuels = append(Fuels, fuelconsumption)
						val.PrimaryFuelTotalCost += fuelconsumption.FuelCost

						if Plant == "Qurayyah" {
							fuelconsumption = VEDFuel{}
							fuelconsumption.IsPrimaryFuel = true
							if len(tempFuelCosts) > 0 {
								fuelconsumption.FuelType = tempFuelCosts[0].GetString("primary2fueltype")
							}

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
							Fuels = append(Fuels, fuelconsumption)
							val.PrimaryFuelTotalCost += fuelconsumption.FuelCost
						}
						//#endregion
					}
					//#endregion
					//#region backup fuel
					BackupFuelType := ""
					if len(tempFuelCosts) > 0 {
						BackupFuelType = tempFuelCosts[0].GetString("backupfueltype")
					}

					if strings.TrimSpace(strings.ToLower(BackupFuelType)) == "hfo" {
						// #region hfo
						BackupFuelConsumed := crowd.From(&tempFuelCosts).Sum(func(x interface{}) interface{} {
							return x.(tk.M).GetFloat64("backupfuelconsumed")
						}).Exec().Result.Sum

						if strings.TrimSpace(strings.ToLower(val.Plant)) == "shoaiba" {
							fuelconsumption := VEDFuel{}
							fuelconsumption.IsPrimaryFuel = false
							fuelconsumption.FuelType = "CRUDE"
							fuelconsumption.FuelCostPerUnit = 0.1
							fuelconsumption.FuelConsumed = BackupFuelConsumed / 3
							fuelconsumption.ConvertedFuelConsumed = fuelconsumption.FuelConsumed * 1000
							fuelconsumption.FuelCost = fuelconsumption.FuelCostPerUnit * fuelconsumption.ConvertedFuelConsumed
							Fuels = append(Fuels, fuelconsumption)
							val.BackupFuelTotalCost += fuelconsumption.FuelCost

							fuelconsumption = VEDFuel{}
							fuelconsumption.IsPrimaryFuel = false
							fuelconsumption.FuelType = "CRUDE HEAVY"
							fuelconsumption.FuelCostPerUnit = 0.049
							fuelconsumption.FuelConsumed = BackupFuelConsumed / 3
							fuelconsumption.ConvertedFuelConsumed = fuelconsumption.FuelConsumed * 1000
							fuelconsumption.FuelCost = fuelconsumption.FuelCostPerUnit * fuelconsumption.ConvertedFuelConsumed
							Fuels = append(Fuels, fuelconsumption)
							val.BackupFuelTotalCost += fuelconsumption.FuelCost

							fuelconsumption = VEDFuel{}
							fuelconsumption.IsPrimaryFuel = false
							fuelconsumption.FuelType = "DIESEL"
							fuelconsumption.FuelCostPerUnit = 0.085
							fuelconsumption.FuelConsumed = BackupFuelConsumed / 3
							fuelconsumption.ConvertedFuelConsumed = fuelconsumption.FuelConsumed * 1000
							fuelconsumption.FuelCost = fuelconsumption.FuelCostPerUnit * fuelconsumption.ConvertedFuelConsumed
							Fuels = append(Fuels, fuelconsumption)
							val.BackupFuelTotalCost += fuelconsumption.FuelCost
						} else if strings.TrimSpace(strings.ToLower(val.Plant)) == "rabigh" {

							fuelconsumption := VEDFuel{}
							fuelconsumption.IsPrimaryFuel = false
							fuelconsumption.FuelType = "CRUDE"
							fuelconsumption.FuelCost = 0.1
							fuelconsumption.FuelConsumed = BackupFuelConsumed / 3
							fuelconsumption.ConvertedFuelConsumed = fuelconsumption.FuelConsumed * 1000
							Fuels = append(Fuels, fuelconsumption)
							val.BackupFuelTotalCost += fuelconsumption.FuelCost

							fuelconsumption = VEDFuel{}
							fuelconsumption.IsPrimaryFuel = true
							fuelconsumption.FuelType = "DIESEL"
							fuelconsumption.FuelCost = 0.085
							fuelconsumption.FuelConsumed = BackupFuelConsumed / 3
							fuelconsumption.ConvertedFuelConsumed = fuelconsumption.FuelConsumed * 1000
							Fuels = append(Fuels, fuelconsumption)
							val.BackupFuelTotalCost += fuelconsumption.FuelCost
						}
						//#endregion
					} else {
						//#region not hfo
						fuelconsumption := VEDFuel{}
						fuelconsumption.IsPrimaryFuel = false
						fuelconsumption.FuelType = BackupFuelType
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
						Fuels = append(Fuels, fuelconsumption)
						val.BackupFuelTotalCost += fuelconsumption.FuelCost
						//#endregion
					}
					//#endregion
					totaldieselconsumed := 0.0
					fuelconsumptionFilter := crowd.From(&Fuels).Where(func(x interface{}) interface{} {
						return strings.ToLower(strings.Trim(x.(VEDFuel).FuelType, " ")) == "diesel"
					}).Exec().Result.Data().([]VEDFuel)

					totaldieselconsumed = crowd.From(&fuelconsumptionFilter).Sum(func(x interface{}) interface{} {
						return x.(VEDFuel).ConvertedFuelConsumed
					}).Exec().Result.Sum

					val.FuelTransportCost = TransportCosts * totaldieselconsumed
					val.TotalFuelCost = val.PrimaryFuelTotalCost + val.BackupFuelTotalCost
					val.OperatingCost = val.FuelTransportCost + val.TotalFuelCost
					// #endregion
					// #region Maintenance
					tempLists = crowd.From(&lists).Where(func(x interface{}) interface{} {
						return x.(tk.M).GetString("unit") == tempunit
					}).Exec().Result.Data().([]tk.M)
					val.TotalLabourCost = crowd.From(&tempLists).Sum(func(x interface{}) interface{} {
						return x.(tk.M).GetFloat64("skilledlabour") + x.(tk.M).GetFloat64("unskilledlabour")
					}).Exec().Result.Sum
					val.TotalMaterialCost = crowd.From(&tempLists).Sum(func(x interface{}) interface{} {
						return x.(tk.M).GetFloat64("materials")
					}).Exec().Result.Sum
					val.TotalMaterialCost = crowd.From(&tempLists).Sum(func(x interface{}) interface{} {
						return x.(tk.M).GetFloat64("contractmaintenance")
					}).Exec().Result.Sum

					Details := []VEDDetail{}
					Top10s := []VEDTop10{}

					tempGT := tempLists

					if len(tempGT) > 0 {
						for _, gt := range tempGT {
							det := VEDDetail{}
							det.DataSource = "Paper Records"
							det.WorkOrderType = gt.GetString("wotype")
							det.LaborCost = gt.GetFloat64("skilledlabour") + gt.GetFloat64("unskilledlabour")
							det.MaterialCost = gt.GetFloat64("materials")
							det.ServiceCost = gt.GetFloat64("contractmaintenance")
							Details = append(Details, det)
						}
					}

					tempsyn := crowd.From(&syn).Where(func(x interface{}) interface{} {
						unitDB := x.(tk.M).GetString("unit")
						return strings.Replace(strings.Replace(strings.Replace(strings.Replace(unitDB, "GT0", "", -1), "GT", "", -1), "ST0", "S", -1), "ST", "S", -1) == strings.Replace(strings.Replace(strings.Replace(strings.Replace(unit, "GT0", "", -1), "GT", "", -1), "ST0", "S", -1), "ST", "S", -1)
					}).Exec().Result.Data().([]tk.M)

					if len(tempsyn) > 0 {
						for _, pm := range tempsyn {
							det := VEDDetail{}
							det.DataSource = "SAP PM"
							det.WorkOrderType = pm.GetString("wotype")
							det.LaborCost = pm.GetFloat64("plannedlaborcost")
							det.MaterialCost = pm.GetFloat64("actualmaterialcost")
							det.ServiceCost = 0

							Details = append(Details, det)

							val.TotalLabourCost += pm.GetFloat64("plannedlaborcost")
							val.TotalMaterialCost += pm.GetFloat64("actualmaterialcost")
						}
					}

					tempbrowser := crowd.From(&databrowser).Where(func(x interface{}) interface{} {
						return x.(tk.M).Get("isturbine").(bool) && strings.Replace(strings.Replace(strings.Replace(x.(tk.M).GetString("tinfshortname"), "GT0", "", -1), "GT", "", -1), "ST0", "ST", -1) == strings.Replace(strings.Replace(strings.Replace(val.Unit, "GT0", "", -1), "GT", "", -1), "ST0", "ST", -1)
					}).Exec().Result.Data().([]tk.M)

					if len(tempbrowser) > 0 {
						tempDataBrowser := crowd.From(&databrowser).Where(func(x interface{}) interface{} {
							return x.(tk.M).GetString("turbineparent") == tempbrowser[0].GetString("functionallocation") || x.(tk.M).GetString("functionallocation") == tempbrowser[0].GetString("functionallocation")
						}).Exec().Result.Data().([]tk.M)

						databrowse := []string{}
						if len(tempDataBrowser) > 0 {
							DataBrowserGroup := crowd.From(&tempDataBrowser).Group(func(x interface{}) interface{} {
								return strings.TrimSpace(x.(tk.M).GetString("functionallocation"))
							}, nil).Exec().Result.Data().([]crowd.KV)

							for _, temp := range DataBrowserGroup {
								databrowse = append(databrowse, temp.Key.(string))
							}
						}

						if len(databrowse) > 0 {
							query = append(query[0:0], dbox.In("FunctionalLocation", databrowse))
							csr, e = c.NewQuery().From(new(WOList).TableName()).Where(query...).Cursor(nil)
							tempWoList := []tk.M{}
							e = csr.Fetch(&tempWoList, 0, false)
							csr.Close()

							WoListTemp := crowd.From(&tempWoList).Group(func(x interface{}) interface{} {
								return x.(tk.M).GetString("ordercode")
							}, nil).Exec().Result.Data().([]crowd.KV)

							MaintenanceOrderList := []string{}
							for _, wo := range WoListTemp {
								MaintenanceOrderList = append(MaintenanceOrderList, wo.Key.(string))
							}

							maintCost := []tk.M{}
							if len(MaintenanceOrderList) > 0 {
								query = append(query[0:0], dbox.And(dbox.In("MaintenanceOrder", MaintenanceOrderList), dbox.Eq("Period", YearFirst)))
								csr, e = c.NewQuery().From(new(MaintenanceCost).TableName()).Where(query...).Cursor(nil)
								e = csr.Fetch(&maintCost, 0, false)
								csr.Close()
							}
							maintHour := []tk.M{}
							if len(MaintenanceOrderList) > 0 {
								query = append(query[0:0], dbox.And(dbox.In("MaintenanceOrder", MaintenanceOrderList), dbox.Eq("Period", YearFirst)))
								csr, e = c.NewQuery().From(new(MaintenanceCostByHour).TableName()).Where(query...).Cursor(nil)
								e = csr.Fetch(&maintHour, 0, false)
								csr.Close()
							}

							if len(maintCost) > 0 {
								tempMaintCost := crowd.From(&maintCost).Group(func(x interface{}) interface{} {
									return x.(tk.M).GetString("ordertype")
								}, nil).Exec().Result.Data().([]crowd.KV)

								MROTypes := []string{}
								for _, tempMaint := range tempMaintCost {
									MROTypes = append(MROTypes, tempMaint.Key.(string))
								}

								for _, types := range MROTypes {
									det := VEDDetail{}
									det.DataSource = "SAP PM"
									det.WorkOrderType = types
									tempMaintHour := crowd.From(&maintHour).Where(func(x interface{}) interface{} {
										return x.(tk.M).GetString("ordertype") == types
									}).Exec().Result.Data().([]tk.M)

									tempMaintCost := crowd.From(&maintHour).Where(func(x interface{}) interface{} {
										return x.(tk.M).GetString("ordertype") == types
									}).Exec().Result.Data().([]tk.M)

									det.Duration = crowd.From(&tempMaintHour).Sum(func(x interface{}) interface{} {
										return x.(tk.M).GetFloat64("actual")
									}).Exec().Result.Sum
									det.LaborCost = crowd.From(&tempMaintCost).Sum(func(x interface{}) interface{} {
										return x.(tk.M).GetFloat64("internallaboractual")
									}).Exec().Result.Sum
									det.MaterialCost = crowd.From(&tempMaintCost).Sum(func(x interface{}) interface{} {
										return x.(tk.M).GetFloat64("directmaterialactual") + x.(tk.M).GetFloat64("internalmaterialactual")
									}).Exec().Result.Sum
									det.ServiceCost = crowd.From(&tempMaintCost).Sum(func(x interface{}) interface{} {
										return x.(tk.M).GetFloat64("externalserviceactual")
									}).Exec().Result.Sum
									Details = append(Details, det)
									val.TotalLabourCost += det.LaborCost
									val.TotalMaterialCost += det.MaterialCost
									val.TotalServicesCost += det.ServiceCost
									val.TotalDuration += det.Duration
								}
								//#region Top10
								for _, fl := range MaintenanceOrderList {
									db := crowd.From(&maintCost).Where(func(x interface{}) interface{} {
										return x.(tk.M).GetString("maintenanceorder") == fl
									}).Exec().Result.Data().([]tk.M)

									if len(db) > 0 {
										top10 := VEDTop10{}
										top10.WorkOrderID = db[0].GetString("maintenanceorder")
										top10.WorkOrderDescription = db[0].GetString("maintenanceorderdesc")
										top10.WorkOrderType = db[0].GetString("ordertype")
										top10.WorkOrderTypeDescription = db[0].GetString("ordertypedesc")
										top10.EquipmentType = db[0].GetString("equipmenttype")
										top10.EquipmentTypeDescription = db[0].GetString("equipmenttypedesc")
										top10.MaintenanceActivity = db[0].GetString("mainactivitytype")
										tempMaintHour := crowd.From(&maintHour).Where(func(x interface{}) interface{} {
											return x.(tk.M).GetString("maintenanceorder") == fl
										}).Exec().Result.Data().([]tk.M)
										top10.Duration = crowd.From(&tempMaintHour).Sum(func(x interface{}) interface{} {
											return x.(tk.M).GetFloat64("actual")
										}).Exec().Result.Sum

										tempSyn := crowd.From(&tempsyn).Where(func(x interface{}) interface{} {
											return x.(tk.M).GetString("woid") == fl
										}).Exec().Result.Data().([]tk.M)
										jumPlanned := crowd.From(&tempSyn).Sum(func(x interface{}) interface{} {
											return x.(tk.M).GetFloat64("plannedlaborcost")
										}).Exec().Result.Sum
										jumActual := crowd.From(&tempSyn).Sum(func(x interface{}) interface{} {
											return x.(tk.M).GetFloat64("actualmaterialcost")
										}).Exec().Result.Sum
										top10.LaborCost = crowd.From(&db).Sum(func(x interface{}) interface{} {
											return x.(tk.M).GetFloat64("internallaboractual") + jumPlanned
										}).Exec().Result.Sum
										top10.MaterialCost = crowd.From(&db).Sum(func(x interface{}) interface{} {
											return x.(tk.M).GetFloat64("internalmaterialactual") + x.(tk.M).GetFloat64("directmaterialactual") + jumActual
										}).Exec().Result.Sum
										top10.ServiceCost = crowd.From(&db).Sum(func(x interface{}) interface{} {
											return x.(tk.M).GetFloat64("externalserviceactual")
										}).Exec().Result.Sum
										top10.MaintenanceCost = top10.LaborCost + top10.MaterialCost + top10.ServiceCost

										Top10s = append(Top10s, top10)
									}
								}
								//#endregion
							}
						}

						val.MaintenanceCost = val.TotalLabourCost + val.TotalMaterialCost + val.TotalServicesCost

						//#endregion
						//#region New Report
						if Plant == "PP9" || Plant == "Qurayyah" || Plant == "Qurayyah CC" {
							tempUnitPower := crowd.From(&unitpower).Where(func(x interface{}) interface{} {
								return x.(tk.M).GetString("plant") == Plant && x.(tk.M).GetString("unit") == strings.Replace(strings.Replace(val.Unit, "GT0", "GT", -1), "ST0", "ST", -1)
							}).Exec().Result.Data().([]tk.M)
							if len(tempUnitPower) > 0 {
								val.MaxCapacity = tempUnitPower[0].GetFloat64("maxpower")
							}
						} else if Plant == "Rabigh" {
							if strings.Contains(val.Unit, "GT") {
								tempUnitPower := crowd.From(&unitpower).Where(func(x interface{}) interface{} {
									return strings.Contains(x.(tk.M).GetString("plant"), Plant) && x.(tk.M).GetString("unit") == strings.Replace(strings.Replace(val.Unit, "GT0", "GT", -1), "ST0", "ST", -1)
								}).Exec().Result.Data().([]tk.M)
								if len(tempUnitPower) > 0 {
									val.MaxCapacity = tempUnitPower[0].GetFloat64("maxpower")
								}
							} else if strings.Contains(val.Unit, "ST") {
								tempUnitPower := crowd.From(&unitpower).Where(func(x interface{}) interface{} {
									return x.(tk.M).GetString("plant") == "Rabigh Steam" && x.(tk.M).GetString("unit") == strings.Replace(strings.Replace(val.Unit, "GT0", "GT", -1), "ST0", "ST", -1)
								}).Exec().Result.Data().([]tk.M)
								if len(tempUnitPower) > 0 {
									val.MaxCapacity = tempUnitPower[0].GetFloat64("maxpower")
								}
							}
						} else if Plant == "Shoaiba" || Plant == "Ghazlan" {
							tempUnitPower := crowd.From(&unitpower).Where(func(x interface{}) interface{} {
								return strings.Contains(x.(tk.M).GetString("plant"), Plant) && x.(tk.M).GetString("unit") == strings.Replace(strings.Replace(val.Unit, "GT0", "GT", -1), "ST0", "ST", -1)
							}).Exec().Result.Data().([]tk.M)
							if len(tempUnitPower) > 0 {
								val.MaxCapacity = tempUnitPower[0].GetFloat64("maxpower")
							}
						}

						val.MaxPowerGeneration = val.MaxCapacity * 24 * 365
						val.PotentialRevenue = val.CapacityPayment + (val.MaxPowerGeneration * val.VOMR * 10) + val.Incentive + val.StartupPayment - val.PenaltyAmount

						if Plant == "Rabigh" {
							sintax := "select * from PowerPlantOutagesDetails inner join PowerPlantOutages on PowerPlantOutagesDetails.POId = PowerPlantOutages.Id where PowerPlantOutagesDetails.UnitNo = '" + val.Unit + "' and PowerPlantOutages.Plant = 'Rabigh Steam'"
							csr, e = c.NewQuery().Command("freequery", tk.M{}.Set("syntax", sintax)).Cursor(nil)
							outagesDetails := []tk.M{}
							e = csr.Fetch(&outagesDetails, 0, false)
							csr.Close()

							tempOutagesDetails := crowd.From(&outagesDetails).Where(func(x interface{}) interface{} {
								return strings.Contains(x.(tk.M).GetString("outagetype"), "FO")
							}).Exec().Result.Data().([]tk.M)

							val.ForcedOutages = crowd.From(&tempOutagesDetails).Sum(func(x interface{}) interface{} {
								return x.(tk.M).GetFloat64("totalhours")
							}).Exec().Result.Sum

							tempOutagesDetails = crowd.From(&outagesDetails).Where(func(x interface{}) interface{} {
								return !strings.Contains(x.(tk.M).GetString("outagetype"), "FO")
							}).Exec().Result.Data().([]tk.M)
							val.UnforcedOutages = crowd.From(&tempOutagesDetails).Sum(func(x interface{}) interface{} {
								return x.(tk.M).GetFloat64("totalhours")
							}).Exec().Result.Sum
						} else if Plant == "Qurayyah" || Plant == "Qurayyah CC" {
							sintax := "select * from PowerPlantOutagesDetails inner join PowerPlantOutages on PowerPlantOutagesDetails.POId = PowerPlantOutages.Id where PowerPlantOutagesDetails.UnitNo = '" + val.Unit + "' and PowerPlantOutages.Plant = '" + Plant + "'"
							csr, e = c.NewQuery().Command("freequery", tk.M{}.Set("syntax", sintax)).Cursor(nil)
							outagesDetails := []tk.M{}
							e = csr.Fetch(&outagesDetails, 0, false)
							csr.Close()

							tempOutagesDetails := crowd.From(&outagesDetails).Where(func(x interface{}) interface{} {
								return strings.Contains(x.(tk.M).GetString("outagetype"), "FO")
							}).Exec().Result.Data().([]tk.M)

							val.ForcedOutages = crowd.From(&tempOutagesDetails).Sum(func(x interface{}) interface{} {
								return x.(tk.M).GetFloat64("totalhours")
							}).Exec().Result.Sum

							tempOutagesDetails = crowd.From(&outagesDetails).Where(func(x interface{}) interface{} {
								return !strings.Contains(x.(tk.M).GetString("outagetype"), "FO")
							}).Exec().Result.Data().([]tk.M)
							val.UnforcedOutages = crowd.From(&tempOutagesDetails).Sum(func(x interface{}) interface{} {
								return x.(tk.M).GetFloat64("totalhours")
							}).Exec().Result.Sum
						} else {
							sintax := "select * from PowerPlantOutagesDetails inner join PowerPlantOutages on PowerPlantOutagesDetails.POId = PowerPlantOutages.Id where PowerPlantOutagesDetails.UnitNo = '" + val.Unit + "'"
							csr, e = c.NewQuery().Command("freequery", tk.M{}.Set("syntax", sintax)).Cursor(nil)
							outagesDetails := []tk.M{}
							e = csr.Fetch(&outagesDetails, 0, false)
							csr.Close()

							tempOutagesDetails := crowd.From(&outagesDetails).Where(func(x interface{}) interface{} {
								return strings.Contains(x.(tk.M).GetString("outagetype"), "FO")
							}).Exec().Result.Data().([]tk.M)

							val.ForcedOutages = crowd.From(&tempOutagesDetails).Sum(func(x interface{}) interface{} {
								return x.(tk.M).GetFloat64("totalhours")
							}).Exec().Result.Sum

							tempOutagesDetails = crowd.From(&outagesDetails).Where(func(x interface{}) interface{} {
								return !strings.Contains(x.(tk.M).GetString("outagetype"), "FO")
							}).Exec().Result.Data().([]tk.M)
							val.UnforcedOutages = crowd.From(&tempOutagesDetails).Sum(func(x interface{}) interface{} {
								return x.(tk.M).GetFloat64("totalhours")
							}).Exec().Result.Sum
						}

						val.ForcedOutagesLoss = (val.PotentialRevenue / (24 * 365)) * val.ForcedOutages
						val.UnforcedOutagesLoss = (val.PotentialRevenue / (24 * 365)) * val.UnforcedOutages

						//#endregion
						//
						val.ValueEquationCost = val.Revenue - val.OperatingCost - val.MaintenanceCost

						id, _ := ctx.InsertOut(val)
						if len(Fuels) > 0 {
							for _, data := range Fuels {
								data.VEId = id

								_, e = ctx.InsertOut(&data)
							}
						}

						if len(Details) > 0 {
							for _, data := range Details {
								data.VEId = id

								_, e = ctx.InsertOut(&data)
							}
						}

						if len(Top10s) > 0 {
							for _, data := range Top10s {
								data.VEId = id

								_, e = ctx.InsertOut(&data)
							}
						}
					}
				}
			}
		}
	}
	return e
}

func (d *GenValueEquationDashboard) generateOtherMaintenanceData(Years []int) error {
	ctx := d.BaseController.Ctx
	c := ctx.Connection

	var (
		query []*dbox.Filter
		e     error
	)

	for _, Year := range Years {
		YearFirst := strconv.Itoa(Year) + "-01-01"
		YearLast := strconv.Itoa(Year+1) + "-01-01"

		sintax := "select VEDTop10.WorkOrderID from ValueEquation_Dashboard inner join VEDTop10 on ValueEquation_Dashboard.Id = VEDTop10.VEId where ValueEquation_Dashboard.Year = " + strconv.Itoa(Year) + " Group By VEDTop10.WorkOrderID"
		csr, _ := c.NewQuery().Command("freequery", tk.M{}.Set("syntax", sintax)).Cursor(nil)
		result := []tk.M{}
		e = csr.Fetch(&result, 0, false)
		csr.Close()

		codes := []string{}
		for _, x := range result {
			codes = append(codes, x.GetString("workorderid"))
		}

		maxLength := 1000
		maxLoop := tk.ToInt(tk.ToFloat64(len(codes)/maxLength, 0, tk.RoundingUp), tk.RoundingUp)

		for i := 0; i <= maxLoop; i++ {
			codeTemp := []string{}

			if i != maxLoop {
				codeTemp = codes[i*maxLength : (i*maxLength)+maxLength]

			} else {
				codeTemp = codes[i*maxLength:]
			}

			tmpNotIn := []orm.IModel{}
			for _, val := range codeTemp {
				tmpData := new(GenVEDNotInTmp)
				tmpData.WorkOrderID = val
				tmpNotIn = append(tmpNotIn, tmpData)
			}
			// tk.Printf("tmpNotIn: %v \n", len(tmpNotIn))
			e = ctx.InsertBulk(tmpNotIn)
			ErrorHandler(e, "generateMaintenanceDataBrowser")
		}

		e = c.NewQuery().Delete().From(new(GenVEDNotInTmp).TableName()).SetConfig("multiexec", true).Exec(nil)
		/*query = append(query, dbox.Nin("MaintenanceOrder", codes))
		query = append(query, dbox.Ne("MaintenanceOrder", ""))
		query = append(query, dbox.Eq("Period", YearFirst))

		csr, e = c.NewQuery().From(new(MaintenanceCost).TableName()).Cursor(nil)
		maintCost := []tk.M{}
		e = csr.Fetch(&maintCost, 0, false)
		csr.Close()*/

		_ = YearFirst
		_ = YearLast
	}

	_ = query
	return e
}
