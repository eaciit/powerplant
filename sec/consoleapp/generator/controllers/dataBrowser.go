package controllers

import (
	"strings"
	"time"

	"github.com/eaciit/crowd"
	"github.com/eaciit/dbox"
	. "github.com/eaciit/powerplant/sec/library/models"
	. "github.com/eaciit/powerplant/sec/webapp/helpers"
	tk "github.com/eaciit/toolkit"
)

var (
	anNewEquipmentType = tk.M{}
)

// GenValueEquation ...
type GenDataBrowser struct {
	*BaseController
}

// Generate ...
func (d *GenDataBrowser) Generate(base *BaseController) {
	var (
		e error
	)
	if base != nil {
		d.BaseController = base
	}

	years := []int{2011, 2012, 2013, 2014, 2015}
	conditions := tk.M{}
	conditions.Set("2110", tk.M{}.Set("length", 9).Set("det", []tk.M{
		tk.M{}.Set("desc", "STEAM TURBINE"),
		tk.M{}.Set("desc", "GAS TURBINE"),
	}))
	conditions.Set("2210", tk.M{}.Set("length", 11).Set("det", []tk.M{
		tk.M{}.Set("desc", "STEAM TURBINE"),
		tk.M{}.Set("desc", "GAS TURBINE"),
	}))
	conditions.Set("2220", tk.M{}.Set("length", 9).Set("det", []tk.M{
		tk.M{}.Set("desc", "Unit"),
	}))
	conditions.Set("2310", tk.M{}.Set("length", 9).Set("det", []tk.M{
		tk.M{}.Set("desc", "UNIT #"),
	}))
	conditions.Set("2320", tk.M{}.Set("length", 9).Set("det", []tk.M{
		tk.M{}.Set("desc", "Unit"),
	}))
	conditions.Set("2325", tk.M{}.Set("length", 11).Set("det", []tk.M{
		tk.M{}.Set("desc", "STEAM TURBINE"),
		tk.M{}.Set("desc", "Gas Turbine"),
	}))

	e = d.generateMaintenanceDataBrowser(years, conditions)
	ErrorHandler(e, "Generate")

	tk.Println("##Data Browser Data : DONE\n")
}

func (d *GenDataBrowser) generateMaintenanceDataBrowser(years []int, conditions tk.M) (e error) {
	ctx := d.BaseController.Ctx
	c := ctx.Connection
	// dataToSave := []DataBrowser{}

	tk.Println("Generating Maintenance Data Browser..")

	// Get fuelCost
	csr, e := c.NewQuery().From(new(FuelCost).TableName()).Cursor(nil)
	fuelCosts := []FuelCost{}
	e = csr.Fetch(&fuelCosts, 0, false)
	ErrorHandler(e, "generateMaintenanceDataBrowser")
	csr.Close()

	// Get plants
	csr, e = c.NewQuery().From(new(PowerPlantCoordinates).TableName()).Cursor(nil)
	plants := []PowerPlantCoordinates{}
	e = csr.Fetch(&plants, 0, false)
	ErrorHandler(e, "generateMaintenanceDataBrowser")
	csr.Close()

	// Get generalInfo
	csr, e = c.NewQuery().From(new(GeneralInfo).TableName()).Cursor(nil)
	generalInfos := []GeneralInfo{}
	e = csr.Fetch(&generalInfos, 0, false)
	ErrorHandler(e, "generateMaintenanceDataBrowser")
	csr.Close()

	// Get GeneralInfoDetails
	csr, e = c.NewQuery().From(new(GeneralInfoDetails).TableName()).Cursor(nil)
	generalInfoDetails := []GeneralInfoDetails{}
	e = csr.Fetch(&generalInfoDetails, 0, false)
	ErrorHandler(e, "generateMaintenanceDataBrowser")
	csr.Close()

	// Get GeneralInfoActualFuelConsumption
	csr, e = c.NewQuery().From(new(GeneralInfoActualFuelConsumption).TableName()).Cursor(nil)
	generalInfoActualFuelConsumption := []GeneralInfoActualFuelConsumption{}
	e = csr.Fetch(&generalInfoActualFuelConsumption, 0, false)
	ErrorHandler(e, "generateMaintenanceDataBrowser")
	csr.Close()

	// Get NewEquipmentType
	csr, e = c.NewQuery().Select("EquipmentType, NewEquipmentGroup").From(new(NewEquipmentType).TableName()).Where(dbox.Ne("NewEquipmentGroup", "Disregard")).Cursor(nil)
	newEquipmentType := []NewEquipmentType{}
	e = csr.Fetch(&newEquipmentType, 0, false)
	ErrorHandler(e, "generateMaintenanceDataBrowser")
	csr.Close()

	constructEquipmentType(newEquipmentType)

	/*wg := &sync.WaitGroup{}
	wg.Add(len(plants))*/

	for _, plant := range plants {
		go func() {
			plantCodeStr := plant.PlantCode

			tmpPlantCondition := conditions.Get(plantCodeStr)
			if tmpPlantCondition != nil {
				plantCondition := tmpPlantCondition.(tk.M)

				length := plantCondition.GetInt("length")
				dets := plantCondition.Get("det").([]tk.M)
				dets = append(dets, tk.M{})

				turbinesCodes := []string{}

				for i, det := range dets {
					assets := []FunctionalLocation{}
					systemAssets := []FunctionalLocation{}

					query := []*dbox.Filter{}
					query = append(query, dbox.Contains("FunctionalLocationCode", plantCodeStr))
					query = append(query, dbox.Eq("PIPI", plantCodeStr))

					if i == 1 {
						query = append(query, dbox.Contains("Description", det.GetString("desc")))
						if plantCodeStr == "2220" {
							query = append(query, dbox.Lte("LEN(FunctionalLocationCode)", length))
						} else {
							query = append(query, dbox.Eq("LEN(FunctionalLocationCode)", length))
						}

					} else if i == len(dets)-1 {
						query = append(query, dbox.Contains("Description", det.GetString("desc")))
						if plantCodeStr == "2220" {
							query = append(query, dbox.Lte("LEN(FunctionalLocationCode)", length))
						} else {
							query = append(query, dbox.Eq("LEN(FunctionalLocationCode)", length))
						}
					} else {
						if len(turbinesCodes) > 1 {
							query = append(query, dbox.Nin("FunctionalLocationCode", turbinesCodes))
						}
					}

					csrDet, e := c.NewQuery().From(new(FunctionalLocation).TableName()).Where(query...).Cursor(nil)
					ErrorHandler(e, "generateMaintenanceDataBrowser")
					e = csrDet.Fetch(&assets, 0, false)
					ErrorHandler(e, "generateMaintenanceDataBrowser")
					csrDet.Close()

					if len(assets) > 0 {
						relatedAssets := []FunctionalLocation{}

						for _, asset := range assets {

							if plantCodeStr == "2110" {
								if len(asset.FunctionalLocationCode) <= 13 {
									query = []*dbox.Filter{}
									query = append(query,
										dbox.And(
											dbox.Eq("PG", "MP1"),
											dbox.Eq("PIPI", plantCodeStr),
											dbox.Contains("FunctionalLocationCode", asset.FunctionalLocationCode),
											dbox.And(
												dbox.Lte("LEN(FunctionalLocationCode)", 13),
												dbox.Gte("LEN(FunctionalLocationCode)", 12),
											),
										),
									)

									csrDet, e := c.NewQuery().From(new(FunctionalLocation).TableName()).Where(query...).Cursor(nil)
									ErrorHandler(e, "generateMaintenanceDataBrowser")
									e = csrDet.Fetch(&systemAssets, 0, false)
									/*if e != nil {
										tk.Printf("%v | %v \n", asset.FunctionalLocationCode, plantCodeStr)
										tk.Println("1-2")
									}*/
									ErrorHandler(e, "generateMaintenanceDataBrowser")
									csrDet.Close()
								}

								if systemAssets == nil || len(systemAssets) == 0 {
									query = []*dbox.Filter{}

									if i != 2 {
										query = append(query, dbox.Contains("FunctionalLocationCode", asset.FunctionalLocationCode))
									} else {
										query = append(query, dbox.Eq("FunctionalLocationCode", asset.FunctionalLocationCode))
									}

									csrDet, e := c.NewQuery().From(new(FunctionalLocation).TableName()).Where(query...).Cursor(nil)
									ErrorHandler(e, "generateMaintenanceDataBrowser")
									e = csrDet.Fetch(&relatedAssets, 0, false)
									ErrorHandler(e, "generateMaintenanceDataBrowser")
									csrDet.Close()

									for _, related := range relatedAssets {
										isTurbineSystem := false
										if related.FunctionalLocationCode == asset.FunctionalLocationCode && i != 2 {
											isTurbineSystem = true
										}

										newEquipment := d.getNewEquipmentType(related.ObjectType, isTurbineSystem)

										if newEquipment != "" {

											if i != 2 {
												turbinesCodes = append(turbinesCodes, related.FunctionalLocationCode)
											}

											for _, year := range years {
												data := DataBrowser{}
												data.PeriodYear = year
												data.FunctionalLocation = related.FunctionalLocationCode
												data.FLDescription = related.Description
												data.IsTurbine = false

												tk.Printf("%v | %v | ", related.FunctionalLocationCode, i)

												if related.FunctionalLocationCode == asset.FunctionalLocationCode && i != 2 {
													data.IsTurbine = true
													tk.Println(" isTurbine: TRUE")
												} else {
													tk.Println(" isTurbine: FALSE")
													data.TurbineParent = asset.FunctionalLocationCode
												}

												data.AssetType = "Other"

												if i == 0 {
													data.AssetType = "Steam"
												} else if i == 1 {
													data.AssetType = "Gas"
												}

												data.EquipmentType = newEquipment
												data.EquipmentTypeDescription = newEquipment
												data.Plant = plant
												data.PlantCode = plant.PlantCode

												if data.IsTurbine {
													info := GeneralInfo{}
													substr := ""
													substrValInt := 0
													if data.AssetType == "Steam" {
														substrValInt = 1
														substr = "ST"
													} else if data.AssetType == "Gas" {
														substrValInt = 2
													}

													if substrValInt != 0 {
														tmpInfo := crowd.From(&generalInfos).Where(func(x interface{}) interface{} {
															y := x.(GeneralInfo)
															substr = substr + data.FunctionalLocation[len(data.FunctionalLocation)-substrValInt:]
															return strings.Contains(strings.ToLower(strings.Trim(y.Plant, " ")), strings.ToLower(plant.PlantName)) && y.Unit == substr
														}).Exec().Result.Data().([]GeneralInfo)
														if len(tmpInfo) > 0 {
															info = tmpInfo[0]
														}

														if info.Id != "" {
															data.TInfShortName = info.Unit
															data.TInfManufacturer = info.Manufacturer
															data.TInfModel = info.Model
															data.TInfUnitType = info.UnitType
															data.TInfInstalledCapacity = info.InstalledCapacity
															data.TInfOperationalCapacity = info.OperationalCapacity
															data.TInfPrimaryFuel = info.PrimaryFuel1
															data.TInfPrimaryFuel2 = info.PrimaryFuel2Startup
															data.TInfBackupFuel = info.BackupFuel
															data.TInfHeatRate = info.HeatRate
															data.TInfEfficiency = info.Efficiency

															commDate, e := time.Parse("01/02/2006", "01/01"+tk.ToString(info.CommissioningDate))
															ErrorHandler(e, "generateMaintenanceDataBrowser")
															data.TInfCommisioningDate = commDate

															if info.RetirementPlan != "" {
																retirementPlanStr := strings.Split(info.RetirementPlan, "(")[0]
																retirementPlan, e := time.Parse("01/02/2006", "01/01"+retirementPlanStr)
																ErrorHandler(e, "generateMaintenanceDataBrowser")
																data.TInfRetirementPlan = retirementPlan
															}

															installedMWH := crowd.From(&generalInfoDetails).Where(func(x interface{}) interface{} {
																y := x.(GeneralInfoDetails)
																return y.GenID == info.Id && y.Type == "InstalledMWH" && y.Year == year
															}).Exec().Result.Data().([]GeneralInfoDetails)[0]

															data.TInfInstalledMWH = installedMWH.Value

															actualEnergyGeneration := crowd.From(&generalInfoDetails).Where(func(x interface{}) interface{} {
																y := x.(GeneralInfoDetails)
																return y.GenID == info.Id && y.Type == "ActualEnergyGeneration" && y.Year == year
															}).Exec().Result.Data().([]GeneralInfoDetails)[0]

															data.TInfActualEnergyGeneration = actualEnergyGeneration.Value

															capacityFactor := crowd.From(&generalInfoDetails).Where(func(x interface{}) interface{} {
																y := x.(GeneralInfoDetails)
																return y.GenID == info.Id && y.Type == "CapacityFactor" && y.Year == year
															}).Exec().Result.Data().([]GeneralInfoDetails)[0]

															data.TInfCapacityFactor = capacityFactor.Value

															actualFuelConsumption := crowd.From(&generalInfoActualFuelConsumption).Where(func(x interface{}) interface{} {
																y := x.(GeneralInfoActualFuelConsumption)
																return y.GenID == info.Id && y.Year == year
															}).Exec().Result.Data().([]GeneralInfoActualFuelConsumption)[0]

															data.TInfActualFuelConsumption_CrudeBarrel = actualFuelConsumption.CrudeBarrel
															data.TInfActualFuelConsumption_DieselBarrel = actualFuelConsumption.DieselBarrel
															data.TInfActualFuelConsumption_GASMMSCF = actualFuelConsumption.GASMMSCF
															data.TInfActualFuelConsumption_HFOBarrel = actualFuelConsumption.HFOBarrel

															fuelCostCrowd := crowd.From(&fuelCosts).Where(func(x interface{}) interface{} {
																y := x.(FuelCost)
																unitID := strings.Replace(
																	strings.Replace(
																		strings.Replace(
																			strings.Replace(y.UnitId, ".", "", -1), " 0", "", -1), " ", "", -1), "C.C ", "", -1)

																return y.Year == year && y.Plant == data.Plant.PlantName && unitID == data.TInfShortName
															})

															data.TInfUpdateEnergyGeneration = fuelCostCrowd.Sum(func(x interface{}) interface{} {
																y := x.(FuelCost)
																return y.EnergyNetProduction
															}).Exec().Result.Data().(float64)

															data.TInfUpdateFuelConsumption = fuelCostCrowd.Sum(func(x interface{}) interface{} {
																y := x.(FuelCost)
																return y.PrimaryFuelConsumed
															}).Exec().Result.Data().(float64)
														}
													}
													// Vibrations handled by sql query
												}

												// Maintenance handled by sql query
												// FailureNotifications handled by sql query
												// MROElements handled by sql query
												// Operationals handled by sql query

												e = ctx.Insert(&data)
												ErrorHandler(e, "generateMaintenanceDataBrowser")
												tk.Println("save 1")
											}
										}
									}

								} else {
									// with system assets
									if i != 2 {
										turbinesCodes = append(turbinesCodes, asset.FunctionalLocationCode)
									}

									for _, year := range years {
										data := DataBrowser{}
										data.PeriodYear = year
										data.FunctionalLocation = asset.FunctionalLocationCode
										data.FLDescription = asset.Description
										data.IsTurbine = false
										data.IsSystem = false

										tk.Printf("%v | %v | ", asset.FunctionalLocationCode, i)

										if asset.FunctionalLocationCode == asset.FunctionalLocationCode && i != 2 {
											data.IsTurbine = true
											data.IsSystem = true
											tk.Println(" isTurbine: TRUE")
										} else {
											tk.Println(" isTurbine: FALSE")
											data.TurbineParent = asset.FunctionalLocationCode
											data.SystemParent = asset.FunctionalLocationCode
										}

										data.AssetType = "Other"

										if i == 0 {
											data.AssetType = "Steam"
										} else if i == 1 {
											data.AssetType = "Gas"
										}

										data.EquipmentType = "System"
										data.EquipmentTypeDescription = "System"
										data.Plant = plant
										data.PlantCode = plant.PlantCode

										if data.IsTurbine {
											info := GeneralInfo{}
											substr := ""
											substrValInt := 0
											if data.AssetType == "Steam" {
												substrValInt = 1
												substr = "ST"
											} else if data.AssetType == "Gas" {
												substrValInt = 2
											}

											if substrValInt != 0 {
												tmpInfo := crowd.From(&generalInfos).Where(func(x interface{}) interface{} {
													y := x.(GeneralInfo)
													substr = substr + data.FunctionalLocation[len(data.FunctionalLocation)-substrValInt:]
													return strings.Contains(strings.ToLower(strings.Trim(y.Plant, " ")), strings.ToLower(plant.PlantName)) && y.Unit == substr
												}).Exec().Result.Data().([]GeneralInfo)
												if len(tmpInfo) > 0 {
													info = tmpInfo[0]
												}

												if info.Id != "" {
													data.TInfShortName = info.Unit
													data.TInfManufacturer = info.Manufacturer
													data.TInfModel = info.Model
													data.TInfUnitType = info.UnitType
													data.TInfInstalledCapacity = info.InstalledCapacity
													data.TInfOperationalCapacity = info.OperationalCapacity
													data.TInfPrimaryFuel = info.PrimaryFuel1
													data.TInfPrimaryFuel2 = info.PrimaryFuel2Startup
													data.TInfBackupFuel = info.BackupFuel
													data.TInfHeatRate = info.HeatRate
													data.TInfEfficiency = info.Efficiency

													commDate, e := time.Parse("01/02/2006", "01/01"+tk.ToString(info.CommissioningDate))
													ErrorHandler(e, "generateMaintenanceDataBrowser")
													data.TInfCommisioningDate = commDate

													if info.RetirementPlan != "" {
														retirementPlanStr := strings.Split(info.RetirementPlan, "(")[0]
														retirementPlan, e := time.Parse("01/02/2006", "01/01"+retirementPlanStr)
														ErrorHandler(e, "generateMaintenanceDataBrowser")
														data.TInfRetirementPlan = retirementPlan
													}

													installedMWH := crowd.From(&generalInfoDetails).Where(func(x interface{}) interface{} {
														y := x.(GeneralInfoDetails)
														return y.GenID == info.Id && y.Type == "InstalledMWH" && y.Year == year
													}).Exec().Result.Data().([]GeneralInfoDetails)[0]

													data.TInfInstalledMWH = installedMWH.Value

													actualEnergyGeneration := crowd.From(&generalInfoDetails).Where(func(x interface{}) interface{} {
														y := x.(GeneralInfoDetails)
														return y.GenID == info.Id && y.Type == "ActualEnergyGeneration" && y.Year == year
													}).Exec().Result.Data().([]GeneralInfoDetails)[0]

													data.TInfActualEnergyGeneration = actualEnergyGeneration.Value

													capacityFactor := crowd.From(&generalInfoDetails).Where(func(x interface{}) interface{} {
														y := x.(GeneralInfoDetails)
														return y.GenID == info.Id && y.Type == "CapacityFactor" && y.Year == year
													}).Exec().Result.Data().([]GeneralInfoDetails)[0]

													data.TInfCapacityFactor = capacityFactor.Value

													actualFuelConsumption := crowd.From(&generalInfoActualFuelConsumption).Where(func(x interface{}) interface{} {
														y := x.(GeneralInfoActualFuelConsumption)
														return y.GenID == info.Id && y.Year == year
													}).Exec().Result.Data().([]GeneralInfoActualFuelConsumption)[0]

													data.TInfActualFuelConsumption_CrudeBarrel = actualFuelConsumption.CrudeBarrel
													data.TInfActualFuelConsumption_DieselBarrel = actualFuelConsumption.DieselBarrel
													data.TInfActualFuelConsumption_GASMMSCF = actualFuelConsumption.GASMMSCF
													data.TInfActualFuelConsumption_HFOBarrel = actualFuelConsumption.HFOBarrel

													fuelCostCrowd := crowd.From(&fuelCosts).Where(func(x interface{}) interface{} {
														y := x.(FuelCost)
														unitID := strings.Replace(
															strings.Replace(
																strings.Replace(
																	strings.Replace(y.UnitId, ".", "", -1), " 0", "", -1), " ", "", -1), "C.C ", "", -1)

														return y.Year == year && y.Plant == data.Plant.PlantName && unitID == data.TInfShortName
													})

													data.TInfUpdateEnergyGeneration = fuelCostCrowd.Sum(func(x interface{}) interface{} {
														y := x.(FuelCost)
														return y.EnergyNetProduction
													}).Exec().Result.Data().(float64)

													data.TInfUpdateFuelConsumption = fuelCostCrowd.Sum(func(x interface{}) interface{} {
														y := x.(FuelCost)
														return y.PrimaryFuelConsumed
													}).Exec().Result.Data().(float64)
												}
											}
											// Vibrations handled by sql query
										}

										// Maintenance handled by sql query
										// FailureNotifications handled by sql query
										// MROElements handled by sql query
										// Operationals handled by sql query

										e = ctx.Insert(&data)
										ErrorHandler(e, "generateMaintenanceDataBrowser")
										tk.Println("save 2")

									}

									for _, sysAsset := range systemAssets {
										query = []*dbox.Filter{}

										if i != 2 {
											query = append(query, dbox.Contains("FunctionalLocationCode", sysAsset.FunctionalLocationCode))
										} else {
											query = append(query, dbox.Eq("FunctionalLocationCode", sysAsset.FunctionalLocationCode))
										}

										csrDet, e := c.NewQuery().From(new(FunctionalLocation).TableName()).Where(query...).Cursor(nil)
										ErrorHandler(e, "generateMaintenanceDataBrowser")
										e = csrDet.Fetch(&relatedAssets, 0, false)
										ErrorHandler(e, "generateMaintenanceDataBrowser")
										csrDet.Close()

										for _, related := range relatedAssets {
											isTurbineSystem := false
											if (related.FunctionalLocationCode == asset.FunctionalLocationCode || related.FunctionalLocationCode == sysAsset.FunctionalLocationCode) && i != 2 {
												isTurbineSystem = true
											}

											newEquipment := d.getNewEquipmentType(related.ObjectType, isTurbineSystem)

											if newEquipment != "" {

												if i != 2 {
													turbinesCodes = append(turbinesCodes, related.FunctionalLocationCode)
												}

												for _, year := range years {
													data := DataBrowser{}
													data.PeriodYear = year
													data.FunctionalLocation = related.FunctionalLocationCode
													data.FLDescription = related.Description
													data.IsTurbine = false
													data.IsSystem = false

													tk.Printf("%v | %v | ", related.FunctionalLocationCode, i)

													if related.FunctionalLocationCode == sysAsset.FunctionalLocationCode && i != 2 {
														data.IsTurbine = true
														data.IsSystem = true
														tk.Println(" isTurbine: TRUE")
													} else {
														tk.Println(" isTurbine: FALSE")
														data.TurbineParent = asset.FunctionalLocationCode
														data.SystemParent = asset.FunctionalLocationCode
													}

													data.AssetType = "Other"

													if i == 0 {
														data.AssetType = "Steam"
													} else if i == 1 {
														data.AssetType = "Gas"
													}

													data.EquipmentType = newEquipment
													data.EquipmentTypeDescription = newEquipment
													data.Plant = plant
													data.PlantCode = plant.PlantCode

													if data.IsTurbine {
														info := GeneralInfo{}
														substr := ""
														substrValInt := 0
														if data.AssetType == "Steam" {
															substrValInt = 1
															substr = "ST"
														} else if data.AssetType == "Gas" {
															substrValInt = 2
														}

														if substrValInt != 0 {
															tmpInfo := crowd.From(&generalInfos).Where(func(x interface{}) interface{} {
																y := x.(GeneralInfo)
																substr = substr + data.FunctionalLocation[len(data.FunctionalLocation)-substrValInt:]
																return strings.Contains(strings.ToLower(strings.Trim(y.Plant, " ")), strings.ToLower(plant.PlantName)) && y.Unit == substr
															}).Exec().Result.Data().([]GeneralInfo)
															if len(tmpInfo) > 0 {
																info = tmpInfo[0]
															}

															if info.Id != "" {
																data.TInfShortName = info.Unit
																data.TInfManufacturer = info.Manufacturer
																data.TInfModel = info.Model
																data.TInfUnitType = info.UnitType
																data.TInfInstalledCapacity = info.InstalledCapacity
																data.TInfOperationalCapacity = info.OperationalCapacity
																data.TInfPrimaryFuel = info.PrimaryFuel1
																data.TInfPrimaryFuel2 = info.PrimaryFuel2Startup
																data.TInfBackupFuel = info.BackupFuel
																data.TInfHeatRate = info.HeatRate
																data.TInfEfficiency = info.Efficiency

																commDate, e := time.Parse("01/02/2006", "01/01"+tk.ToString(info.CommissioningDate))
																ErrorHandler(e, "generateMaintenanceDataBrowser")
																data.TInfCommisioningDate = commDate

																if info.RetirementPlan != "" {
																	retirementPlanStr := strings.Split(info.RetirementPlan, "(")[0]
																	retirementPlan, e := time.Parse("01/02/2006", "01/01"+retirementPlanStr)
																	ErrorHandler(e, "generateMaintenanceDataBrowser")
																	data.TInfRetirementPlan = retirementPlan
																}

																installedMWH := crowd.From(&generalInfoDetails).Where(func(x interface{}) interface{} {
																	y := x.(GeneralInfoDetails)
																	return y.GenID == info.Id && y.Type == "InstalledMWH" && y.Year == year
																}).Exec().Result.Data().([]GeneralInfoDetails)[0]

																data.TInfInstalledMWH = installedMWH.Value

																actualEnergyGeneration := crowd.From(&generalInfoDetails).Where(func(x interface{}) interface{} {
																	y := x.(GeneralInfoDetails)
																	return y.GenID == info.Id && y.Type == "ActualEnergyGeneration" && y.Year == year
																}).Exec().Result.Data().([]GeneralInfoDetails)[0]

																data.TInfActualEnergyGeneration = actualEnergyGeneration.Value

																capacityFactor := crowd.From(&generalInfoDetails).Where(func(x interface{}) interface{} {
																	y := x.(GeneralInfoDetails)
																	return y.GenID == info.Id && y.Type == "CapacityFactor" && y.Year == year
																}).Exec().Result.Data().([]GeneralInfoDetails)[0]

																data.TInfCapacityFactor = capacityFactor.Value

																actualFuelConsumption := crowd.From(&generalInfoActualFuelConsumption).Where(func(x interface{}) interface{} {
																	y := x.(GeneralInfoActualFuelConsumption)
																	return y.GenID == info.Id && y.Year == year
																}).Exec().Result.Data().([]GeneralInfoActualFuelConsumption)[0]

																data.TInfActualFuelConsumption_CrudeBarrel = actualFuelConsumption.CrudeBarrel
																data.TInfActualFuelConsumption_DieselBarrel = actualFuelConsumption.DieselBarrel
																data.TInfActualFuelConsumption_GASMMSCF = actualFuelConsumption.GASMMSCF
																data.TInfActualFuelConsumption_HFOBarrel = actualFuelConsumption.HFOBarrel

																fuelCostCrowd := crowd.From(&fuelCosts).Where(func(x interface{}) interface{} {
																	y := x.(FuelCost)
																	unitID := strings.Replace(
																		strings.Replace(
																			strings.Replace(
																				strings.Replace(y.UnitId, ".", "", -1), " 0", "", -1), " ", "", -1), "C.C ", "", -1)

																	return y.Year == year && y.Plant == data.Plant.PlantName && unitID == data.TInfShortName
																})

																data.TInfUpdateEnergyGeneration = fuelCostCrowd.Sum(func(x interface{}) interface{} {
																	y := x.(FuelCost)
																	return y.EnergyNetProduction
																}).Exec().Result.Data().(float64)

																data.TInfUpdateFuelConsumption = fuelCostCrowd.Sum(func(x interface{}) interface{} {
																	y := x.(FuelCost)
																	return y.PrimaryFuelConsumed
																}).Exec().Result.Data().(float64)
															}
														}
														// Vibrations handled by sql query
													}

													// Maintenance handled by sql query
													// FailureNotifications handled by sql query
													// MROElements handled by sql query
													// Operationals handled by sql query

													e = ctx.Insert(&data)
													ErrorHandler(e, "generateMaintenanceDataBrowser")
													tk.Println("save 3")
												}
											}
										}
									}
								}
							} else {
								// another plant(s)

								query = []*dbox.Filter{}

								if i != len(dets) {
									query = append(query, dbox.Contains("FunctionalLocationCode", asset.FunctionalLocationCode))
								} else {
									query = append(query, dbox.Eq("FunctionalLocationCode", asset.FunctionalLocationCode))
								}

								csrDet, e := c.NewQuery().From(new(FunctionalLocation).TableName()).Where(query...).Cursor(nil)
								ErrorHandler(e, "generateMaintenanceDataBrowser")
								e = csrDet.Fetch(&relatedAssets, 0, false)
								ErrorHandler(e, "generateMaintenanceDataBrowser")
								csrDet.Close()

								for _, related := range relatedAssets {
									isTurbineSystem := false
									if related.FunctionalLocationCode == asset.FunctionalLocationCode && i != len(dets) {
										isTurbineSystem = true
									}

									newEquipment := d.getNewEquipmentType(related.ObjectType, isTurbineSystem)

									if newEquipment != "" {

										if i != len(dets) {
											turbinesCodes = append(turbinesCodes, related.FunctionalLocationCode)
										}

										for _, year := range years {
											data := DataBrowser{}
											data.PeriodYear = year
											data.FunctionalLocation = related.FunctionalLocationCode
											data.FLDescription = related.Description
											data.IsTurbine = false

											if related.FunctionalLocationCode == asset.FunctionalLocationCode && i != len(dets) {
												data.IsTurbine = true
											} else {
												data.TurbineParent = asset.FunctionalLocationCode
											}

											data.AssetType = "Other"

											if i == 0 {
												data.AssetType = "Steam"
											} else if i == 1 && len(dets) > 1 {
												data.AssetType = "Gas"
											}

											data.EquipmentType = newEquipment
											data.EquipmentTypeDescription = newEquipment
											data.Plant = plant
											data.PlantCode = plant.PlantCode

											if data.IsTurbine {
												info := GeneralInfo{}
												substr := ""
												substrValInt := 0
												if data.AssetType == "Steam" {
													// substrValInt = 1
													substrValInt = 2
													substr = "ST" + substr
												} else if data.AssetType == "Gas" {
													substrValInt = 2
													substr = "GT" + substr
												}

												if substrValInt != 0 {
													tmpInfo := crowd.From(&generalInfos).Where(func(x interface{}) interface{} {
														y := x.(GeneralInfo)
														substr = data.FunctionalLocation[len(data.FunctionalLocation)-substrValInt:]
														return strings.Contains(strings.ToLower(strings.Trim(y.Plant, " ")), strings.ToLower(plant.PlantName)) && y.Unit == substr
													}).Exec().Result.Data().([]GeneralInfo)
													if len(tmpInfo) > 0 {
														info = tmpInfo[0]
													}

													if info.Id != "" {
														data.TInfShortName = info.Unit
														data.TInfManufacturer = info.Manufacturer
														data.TInfModel = info.Model
														data.TInfUnitType = info.UnitType
														data.TInfInstalledCapacity = info.InstalledCapacity
														data.TInfOperationalCapacity = info.OperationalCapacity
														data.TInfPrimaryFuel = info.PrimaryFuel1
														data.TInfPrimaryFuel2 = info.PrimaryFuel2Startup
														data.TInfBackupFuel = info.BackupFuel
														data.TInfHeatRate = info.HeatRate
														data.TInfEfficiency = info.Efficiency

														commDate, e := time.Parse("01/02/2006", "01/01"+tk.ToString(info.CommissioningDate))
														ErrorHandler(e, "generateMaintenanceDataBrowser")

														data.TInfCommisioningDate = commDate

														if info.RetirementPlan != "" {
															retirementPlanStr := strings.Split(info.RetirementPlan, "(")[0]
															retirementPlan, e := time.Parse("01/02/2006", "01/01"+retirementPlanStr)
															ErrorHandler(e, "generateMaintenanceDataBrowser")
															data.TInfRetirementPlan = retirementPlan
														}

														installedMWH := crowd.From(&generalInfoDetails).Where(func(x interface{}) interface{} {
															y := x.(GeneralInfoDetails)
															return y.GenID == info.Id && y.Type == "InstalledMWH" && y.Year == year
														}).Exec().Result.Data().([]GeneralInfoDetails)[0]

														data.TInfInstalledMWH = installedMWH.Value

														actualEnergyGeneration := crowd.From(&generalInfoDetails).Where(func(x interface{}) interface{} {
															y := x.(GeneralInfoDetails)
															return y.GenID == info.Id && y.Type == "ActualEnergyGeneration" && y.Year == year
														}).Exec().Result.Data().([]GeneralInfoDetails)[0]

														data.TInfActualEnergyGeneration = actualEnergyGeneration.Value

														capacityFactor := crowd.From(&generalInfoDetails).Where(func(x interface{}) interface{} {
															y := x.(GeneralInfoDetails)
															return y.GenID == info.Id && y.Type == "CapacityFactor" && y.Year == year
														}).Exec().Result.Data().([]GeneralInfoDetails)[0]

														data.TInfCapacityFactor = capacityFactor.Value

														actualFuelConsumption := crowd.From(&generalInfoActualFuelConsumption).Where(func(x interface{}) interface{} {
															y := x.(GeneralInfoActualFuelConsumption)
															return y.GenID == info.Id && y.Year == year
														}).Exec().Result.Data().([]GeneralInfoActualFuelConsumption)[0]

														data.TInfActualFuelConsumption_CrudeBarrel = actualFuelConsumption.CrudeBarrel
														data.TInfActualFuelConsumption_DieselBarrel = actualFuelConsumption.DieselBarrel
														data.TInfActualFuelConsumption_GASMMSCF = actualFuelConsumption.GASMMSCF
														data.TInfActualFuelConsumption_HFOBarrel = actualFuelConsumption.HFOBarrel

														fuelCostCrowd := crowd.From(&fuelCosts).Where(func(x interface{}) interface{} {
															y := x.(FuelCost)
															unitID := strings.Replace(
																strings.Replace(
																	strings.Replace(
																		strings.Replace(y.UnitId, ".", "", -1), " 0", "", -1), " ", "", -1), "C.C ", "", -1)

															return y.Year == year && y.Plant == data.Plant.PlantName && unitID == data.TInfShortName
														})

														data.TInfUpdateEnergyGeneration = fuelCostCrowd.Sum(func(x interface{}) interface{} {
															y := x.(FuelCost)
															return y.EnergyNetProduction
														}).Exec().Result.Data().(float64)

														data.TInfUpdateFuelConsumption = fuelCostCrowd.Sum(func(x interface{}) interface{} {
															y := x.(FuelCost)
															return y.PrimaryFuelConsumed
														}).Exec().Result.Data().(float64)
													}
												}
												// Vibrations handled by sql query
											}

											// Maintenance handled by sql query
											// FailureNotifications handled by sql query
											// MROElements handled by sql query
											// Operationals handled by sql query

											e = ctx.Insert(&data)
											ErrorHandler(e, "generateMaintenanceDataBrowser")
											tk.Println("save 4")
										}
									}
								}
							}

						}
					}
				}
			}
			// wg.Done()
		}()
	}

	// wg.Wait()
	return
}

func constructEquipmentType(list []NewEquipmentType) {
	for _, item := range list {
		anNewEquipmentType.Set(item.EquipmentType, item.NewEquipmentGroup)
	}
}

func (d *GenDataBrowser) getNewEquipmentType(equipmentType string, isSystem bool) (retVal string) {
	retVal = anNewEquipmentType.GetString(equipmentType)

	if ((retVal == "Default" || retVal == "System") && isSystem) || (retVal != "" && retVal != "Default" && retVal != "System") {
		retVal = retVal
	} else {
		retVal = ""
	}

	return
}

/*func (d *GenDataBrowser) insertMultiple(list []DataBrowser){
	ctx := d.BaseController.Ctx

    if len(list)
}*/

/*func (d *GenDataBrowser) getNewEquipmentType(equipmentType string, isSystem bool) (retVal string) {
	ctx := d.BaseController.Ctx
	c := ctx.Connection

	res := []tk.M{}

	csr, e := c.NewQuery().Select("NewEquipmentGroup").From(new(NewEquipmentType).TableName()).
		Where(dbox.And(
			dbox.Eq("EquipmentType", equipmentType),
			dbox.Ne("NewEquipmentGroup", "Disregard"),
		),
		).
		Cursor(nil)

	if e != nil {
		tk.Println(e.Error())
	}

	e = csr.Fetch(&res, 0, false)

	if e != nil {
		tk.Println(e.Error())
	}

	csr.Close()

	if len(res) > 0 {
		retVal = res[0].GetString("NewEquipmentGroup")

		if ((retVal == "Default" || retVal == "System") && isSystem) || (retVal != "" && retVal != "Default" && retVal != "System") {
			retVal = retVal
		} else {
			retVal = ""
		}
	}

	return
}*/
