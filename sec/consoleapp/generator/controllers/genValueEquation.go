package controllers

import (
	"github.com/eaciit/dbox"
	// . "github.com/eaciit/powerplant/sec/consoleapp/generator/helpers"

	"strconv"
	"strings"

	"github.com/eaciit/crowd"
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

	plants := []string{
		"2110",
		// "2210",
		// "2220",
		// "2310",
		// "2320",
		// "2325",
	}

	for _, plant := range plants {
		e = d.generateValueEquation(2014, plant)

		if e != nil {
			tk.Println(e)
		}
	}

	/*
		e = d.generateValueEquationDataQuality(2014, "Qurayyah CC")
		if e != nil {
			tk.Println(e)
		}
	*/
	tk.Println("##Value Equation Data : DONE\n")
}

func (d *GenValueEquation) generateValueEquation(year int, plantCode string) error {
	var e error

	ctx := d.BaseController.Ctx
	c := ctx.Connection

	yearStart := strconv.Itoa(year) + "-01-01"
	yearEnd := strconv.Itoa(year+1) + "-01-01"
	// daysInYear := 0.0

	fuelCostPerUnit := tk.M{
		"CRUDE":       0.1,
		"CRUDE HEAVY": 0.049,
		"DIESEL":      0.085,
		"GAS":         2.813,
	}

	// MasterPowerPlant
	csr, e := c.NewQuery().From(new(MasterPowerPlant).TableName()).Where(dbox.Eq("PlantCode", plantCode)).Cursor(nil)
	plants := []MasterPowerPlant{}
	e = csr.Fetch(&plants, 0, false)
	csr.Close()

	var totalData int64

	if len(plants) > 0 {
		plant := plants[0]

		_ = plant
		// tk.Printf("Plant: %#v \n", plant.PlantCode)
		// GenerationAppendix
		csr, e = c.NewQuery().From(new(GenerationAppendix).TableName()).Where(dbox.Eq("Plant", plantCode)).Cursor(nil)
		generationAppendix := []GenerationAppendix{}
		e = csr.Fetch(&generationAppendix, 0, false)
		csr.Close()

		// tk.Printf("Appendix: %#v \n", len(generationAppendix))

		for _, appendix := range generationAppendix {
			unit := appendix.Unit

			valueEquation := new(ValueEquation)

			// UnitInfos
			csr, e = c.NewQuery().
				From(new(UnitInfos).TableName()).
				Where(
					dbox.Eq("Plant", plantCode),
					dbox.Eq("Year", year),
					dbox.Eq("Unit", unit)).
				Cursor(nil)
			unitInfos := []UnitInfos{}
			e = csr.Fetch(&unitInfos, 0, false)
			csr.Close()

			var dieselConsumptions float64
			var transportCosts float64

			// tk.Printf("appendix: %v \n", appendix.Plant+" "+appendix.UnitType+" "+appendix.Unit)
			// tk.Printf("UnitInfos: %#v \n", len(unitInfos))

			if len(unitInfos) > 0 {
				unitInfo := unitInfos[0]

				// FuelTransport
				csr, e = c.NewQuery().
					From(new(FuelTransport).TableName()).
					Where(
						dbox.Eq("Plant", plantCode),
						dbox.Eq("Year", year)).
					Cursor(nil)
				fuelTransports := []FuelTransport{}
				e = csr.Fetch(&fuelTransports, 0, false)
				csr.Close()

				fuelTransport := FuelTransport{}

				if len(fuelTransports) > 0 {
					fuelTransport = fuelTransports[0]
				}

				dieselConsumptions = crowd.From(&unitInfos).Where(func(x interface{}) interface{} {
					return strings.ToLower(x.(UnitInfos).PrimaryFuelType) == "diesel"
				}).Sum(func(x interface{}) interface{} {
					return x.(UnitInfos).PrimaryFuelConsumed
				}).Exec().Result.Sum

				dieselConsumptionsBackup := crowd.From(&unitInfos).Where(func(x interface{}) interface{} {
					return strings.ToLower(x.(UnitInfos).BackupFuelType) == "diesel"
				}).Sum(func(x interface{}) interface{} {
					return x.(UnitInfos).BackupFuelConsumed
				}).Exec().Result.Sum

				dieselConsumptions += dieselConsumptionsBackup
				dieselConsumptions = dieselConsumptions * 1000

				if dieselConsumptions != 0 && fuelTransport.TransportCost != 0 {
					transportCosts = fuelTransport.TransportCost / dieselConsumptions
				}

				valueEquation.Plant = plantCode
				valueEquation.Year = year
				valueEquation.Month = 1
				valueEquation.Unit = unit
				valueEquation.UnitGroup = unit[:2]

				// Consolidated
				csr, e = c.NewQuery().
					From(new(Consolidated).TableName()).
					Where(
						dbox.Eq("Plant", plantCode),
						dbox.Eq("Unit", unit),
						dbox.And(dbox.Gte("ConsolidatedDate", yearStart), dbox.Lt("ConsolidatedDate", yearEnd))).
					Cursor(nil)
				consolidateds := []Consolidated{}
				e = csr.Fetch(&consolidateds, 0, false)
				csr.Close()

				valueEquation.Capacity = crowd.From(&consolidateds).Sum(func(x interface{}) interface{} {
					return x.(Consolidated).Capacity
				}).Exec().Result.Sum
				valueEquation.NetGeneration = crowd.From(&consolidateds).Sum(func(x interface{}) interface{} {
					return x.(Consolidated).EnergyNet
				}).Exec().Result.Sum
				valueEquation.AvgNetGeneration = crowd.From(&consolidateds).Avg(func(x interface{}) interface{} {
					return x.(Consolidated).EnergyNet
				}).Exec().Result.Avg

				valueEquation.WAFPercentage = unitInfo.WAFPercentage
				valueEquation.WUFPercentage = unitInfo.WUFPercentage

				valueEquation.VOMR = appendix.VOMR
				valueEquation.CapacityPayment = (appendix.ContractedCapacity * (appendix.FOMR + appendix.CCR)) //* daysInYear * 10
				valueEquation.EnergyPayment = valueEquation.NetGeneration * valueEquation.VOMR * 10
				valueEquation.SRF = unitInfo.SRFPercentage

				// PowerPlantOutages
				csr, e = c.NewQuery().
					From(new(PowerPlantOutages).TableName()).
					Where(
						dbox.Eq("Plant", plantCode),
						dbox.Eq("Unit", unit),
						dbox.Ne("OutageType", "PO"),
						dbox.And(dbox.Gte("StartDate", yearStart), dbox.Lt("StartDate", yearEnd))).
					Cursor(nil)
				powerPlantOutages := []PowerPlantOutages{}
				e = csr.Fetch(&powerPlantOutages, 0, false)
				csr.Close()

				valueEquation.UnplannedOutages = len(powerPlantOutages)
				valueEquation.TotalOutageDuration = crowd.From(&powerPlantOutages).Sum(func(x interface{}) interface{} {
					return x.(PowerPlantOutages).TotalHours
				}).Exec().Result.Sum

				if valueEquation.SRF == 100 {
					valueEquation.StartupPayment = appendix.Startup
					valueEquation.PenaltyAmount = 0
				} else {
					valueEquation.StartupPayment = 0
					valueEquation.PenaltyAmount = appendix.Deduct
				}

				valueEquation.PenaltyAmount += (appendix.Deduct) * tk.ToFloat64(valueEquation.UnplannedOutages, 1, tk.RoundingAuto)
				valueEquation.Incentive = 0
				valueEquation.Revenue = valueEquation.CapacityPayment + valueEquation.EnergyPayment + valueEquation.Incentive + valueEquation.StartupPayment - valueEquation.PenaltyAmount

				// Fuels
				fuels := []ValueEquationFuelData{}
				for typeCode := 1; typeCode <= 2; typeCode++ {
					// primary fuels
					isPrimary := true
					fuelType := unitInfo.PrimaryFuelType
					fuelConsumed := crowd.From(&unitInfos).Sum(func(x interface{}) interface{} {
						return x.(UnitInfos).PrimaryFuelConsumed
					}).Exec().Result.Sum

					if typeCode == 2 {
						// backup fuels
						isPrimary = false
						fuelType = unitInfo.BackupFuelType
						fuelConsumed = crowd.From(&unitInfos).Sum(func(x interface{}) interface{} {
							return x.(UnitInfos).BackupFuelConsumed
						}).Exec().Result.Sum
					}

					if strings.ToLower(fuelType) == "hfo" {
						if plantCode == "2220" {
							// shoaiba
							for i := 0; i < 3; i++ {
								fuel := ValueEquationFuelData{}
								fuel.IsPrimaryFuel = isPrimary

								if i == 0 {
									fuel.FuelType = "CRUDE"
									fuel.FuelCostPerUnit = fuelCostPerUnit.GetFloat64(fuel.FuelType)
								} else if i == 1 {
									fuel.FuelType = "CRUDE HEAVY"
									fuel.FuelCostPerUnit = fuelCostPerUnit.GetFloat64(fuel.FuelType)
								} else if i == 2 {
									fuel.FuelType = "DIESEL"
									fuel.FuelCostPerUnit = fuelCostPerUnit.GetFloat64(fuel.FuelType)
								}

								fuel.FuelConsumed = fuelConsumed / 3
								fuel.ConvertedFuelConsumed = fuel.FuelConsumed * 1000
								fuel.FuelCost = fuel.FuelCostPerUnit * fuel.ConvertedFuelConsumed

								fuels = append(fuels, fuel)

								if typeCode == 1 {
									valueEquation.PrimaryFuelTotalCost += fuel.FuelCost
								} else {
									valueEquation.BackupFuelTotalCost += fuel.FuelCost
								}

							}
						} else if plantCode == "2210" {
							// Rabigh
							for i := 0; i < 2; i++ {
								fuel := ValueEquationFuelData{}
								fuel.IsPrimaryFuel = isPrimary

								if i == 0 {
									fuel.FuelType = "CRUDE"
									fuel.FuelCostPerUnit = fuelCostPerUnit.GetFloat64(fuel.FuelType)
								} else if i == 1 {
									fuel.FuelType = "DIESEL"
									fuel.FuelCostPerUnit = fuelCostPerUnit.GetFloat64(fuel.FuelType)
								}

								fuel.FuelConsumed = fuelConsumed / 3
								fuel.ConvertedFuelConsumed = fuel.FuelConsumed * 1000
								fuel.FuelCost = fuel.FuelCostPerUnit * fuel.ConvertedFuelConsumed

								fuels = append(fuels, fuel)
								if typeCode == 1 {
									valueEquation.PrimaryFuelTotalCost += fuel.FuelCost
								} else {
									valueEquation.BackupFuelTotalCost += fuel.FuelCost
								}
							}
						}
					} else if fuelType != "" {
						max := 0

						if plantCode == "2320" && typeCode == 1 {
							max = 1
						}

						for i := 0; i <= max; i++ {
							if i == 1 {
								fuelType = unitInfo.Primary2FuelType
								fuelConsumed = crowd.From(&unitInfos).Sum(func(x interface{}) interface{} {
									return x.(UnitInfos).Primary2FuelConsumed
								}).Exec().Result.Sum
							}

							fuel := ValueEquationFuelData{}
							fuel.IsPrimaryFuel = isPrimary
							fuel.FuelType = fuelType
							fuel.FuelCostPerUnit = fuelCostPerUnit.GetFloat64(fuel.FuelType)
							fuel.FuelConsumed = fuelConsumed

							if strings.ToLower(fuel.FuelType) == "gas" {
								fuel.ConvertedFuelConsumed = fuel.FuelConsumed * 0.353
							} else {
								fuel.ConvertedFuelConsumed = fuel.FuelConsumed * 1000
							}

							fuel.FuelCost = fuel.FuelCostPerUnit * fuel.ConvertedFuelConsumed

							fuels = append(fuels, fuel)
							if typeCode == 1 {
								valueEquation.PrimaryFuelTotalCost += fuel.FuelCost
							} else {
								valueEquation.BackupFuelTotalCost += fuel.FuelCost
							}
						}
					}
				}

				valueEquation.Fuels = fuels

				totalDieselConsumed := crowd.From(&fuels).Where(func(x interface{}) interface{} {
					return strings.ToLower(x.(ValueEquationFuelData).FuelType) == "diesel"
				}).Sum(func(x interface{}) interface{} {
					return x.(ValueEquationFuelData).ConvertedFuelConsumed
				}).Exec().Result.Sum

				valueEquation.FuelTransportCost = transportCosts * totalDieselConsumed
				valueEquation.TotalFuelCost = valueEquation.PrimaryFuelTotalCost + valueEquation.BackupFuelTotalCost
				valueEquation.OperatingCost = valueEquation.FuelTransportCost + valueEquation.TotalFuelCost

				// PreventiveMaintenance
				csr, e = c.NewQuery().
					From(new(PreventiveMaintenance).TableName()).
					Where(
						dbox.Eq("Plant", plantCode),
						dbox.Eq("Unit", unit),
						dbox.And(dbox.Gte("DatePerformed", yearStart), dbox.Lt("DatePerformed", yearEnd))).
					Cursor(nil)
				preventiveMaintenances := []PreventiveMaintenance{}
				e = csr.Fetch(&preventiveMaintenances, 0, false)
				csr.Close()

				details := []ValueEquationDetails{}

				for _, prev := range preventiveMaintenances {
					detail := ValueEquationDetails{}
					detail.DataSource = "Paper Records"
					detail.WorkOrderType = prev.WOType
					detail.LaborCost = prev.SkilledLabourSAR + prev.UnSkilledLabourSAR
					detail.MaterialCost = prev.MaterialsSAR
					detail.ServiceCost = prev.ContractMaintenanceSAR

					details = append(details, detail)

					valueEquation.TotalLabourCost += (prev.SkilledLabourSAR + prev.UnSkilledLabourSAR)
					valueEquation.TotalMaterialCost += prev.MaterialsSAR
					valueEquation.TotalServicesCost += prev.ContractMaintenanceSAR
				}

				// MasterFunctionalLocation / Databrowser
				csr, e = c.NewQuery().
					Select("FunctionalLocationCode").
					From(new(MasterFunctionalLocation).TableName()).
					Where(
						dbox.Eq("IsTurbine", 1),
						dbox.Eq("Unit", unit)).
					Cursor(nil)
				selectedFunctionalLocations := []MasterFunctionalLocation{}
				e = csr.Fetch(&selectedFunctionalLocations, 0, false)
				csr.Close()

				csr, e = c.NewQuery().
					From(new(MasterFunctionalLocation).TableName()).
					Where(
						dbox.Or(
							dbox.In("SuperiorFunctionalLocation", selectedFunctionalLocations),
							dbox.In("FunctionalLocation", selectedFunctionalLocations))).
					Cursor(nil)
				dBrowser := []MasterFunctionalLocation{}
				e = csr.Fetch(&dBrowser, 0, false)
				csr.Close()

				valueEquation.Details = details

				if e != nil {
					tk.Println(e.Error())
				}

				// if unit == "ST1" {
				if unit == "GT11" {
					// tk.Printf("--> %#v \n", unitInfos)
					tk.Printf("--> %#v \n", valueEquation)
					/*tk.Println("-----------------------------")
					tk.Println(appendix)
					tk.Println(powerPlantOutages)
					tk.Println("-----------------------------")*/
				}

				// tk.Printf("--> %#v \n", valueEquation)

				totalData += 1
			}
		}
	}

	tk.Printf("TotalData: %v \n", totalData)

	return e
}
