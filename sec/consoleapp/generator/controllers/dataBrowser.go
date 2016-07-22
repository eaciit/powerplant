package controllers

import (
	"strings"
	"sync"
	"time"

	"github.com/eaciit/crowd"
	"github.com/eaciit/dbox"
	"github.com/eaciit/orm"
	. "github.com/eaciit/powerplant/sec/library/models"
	. "github.com/eaciit/powerplant/sec/webapp/helpers"
	tk "github.com/eaciit/toolkit"
)

var (
	anNewEquipmentType               = tk.M{}
	dataCounts                       = tk.M{}
	mugen                            = &sync.Mutex{}
	fuelCosts                        = []FuelCost{}
	plants                           = []PowerPlantCoordinates{}
	generalInfos                     = []GeneralInfo{}
	generalInfoDetails               = []GeneralInfoDetails{}
	generalInfoActualFuelConsumption = []GeneralInfoActualFuelConsumption{}
	newEquipmentType                 = []NewEquipmentType{}
	conditions                       = tk.M{}
	years                            = []int{2011, 2012, 2013, 2014, 2015}
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

	conditions.Set("2110", tk.M{}.Set("length", 9).Set("det", []tk.M{
		tk.M{}.Set("desc", "STEAM TURBINE"),
		tk.M{}.Set("desc", "GAS TURBINE"),
	}))
	/*conditions.Set("2210", tk.M{}.Set("length", 11).Set("det", []tk.M{
		tk.M{}.Set("desc", "STEAM TURBINE"),
		tk.M{}.Set("desc", "GAS TURBINE"),
	}))*/
	/*conditions.Set("2220", tk.M{}.Set("length", 9).Set("det", []tk.M{
		tk.M{}.Set("desc", "Unit"),
	}))*/
	/*conditions.Set("2310", tk.M{}.Set("length", 9).Set("det", []tk.M{
		tk.M{}.Set("desc", "UNIT #"),
	}))*/
	/*conditions.Set("2320", tk.M{}.Set("length", 9).Set("det", []tk.M{
		tk.M{}.Set("desc", "Unit"),
	}))*/
	/*conditions.Set("2325", tk.M{}.Set("length", 11).Set("det", []tk.M{
		tk.M{}.Set("desc", "STEAM TURBINE"),
		tk.M{}.Set("desc", "Gas Turbine"),
	}))*/

	tStart := time.Now()

	e = d.generateMaintenanceDataBrowser()
	ErrorHandler(e, "Generate")

	tk.Printf("##Data Browser Data : DONE | %v | count: \n %#v \n", time.Since(tStart), dataCounts)
}

func (d *GenDataBrowser) generateMaintenanceDataBrowser() (e error) {
	ctx := d.BaseController.Ctx
	c := ctx.Connection
	// dataToSave := []DataBrowser{}
	// dataCount = 0

	tk.Println("Generating Maintenance Data Browser..")

	// Get fuelCost
	csr, e := c.NewQuery().From(new(FuelCost).TableName()).Cursor(nil)
	e = csr.Fetch(&fuelCosts, 0, false)
	ErrorHandler(e, "generateMaintenanceDataBrowser")
	csr.Close()

	// Get plants
	csr, e = c.NewQuery().From(new(PowerPlantCoordinates).TableName()).Cursor(nil)
	e = csr.Fetch(&plants, 0, false)
	ErrorHandler(e, "generateMaintenanceDataBrowser")
	csr.Close()

	// Get generalInfo
	csr, e = c.NewQuery().From(new(GeneralInfo).TableName()).Cursor(nil)
	e = csr.Fetch(&generalInfos, 0, false)
	ErrorHandler(e, "generateMaintenanceDataBrowser")
	csr.Close()

	// Get GeneralInfoDetails
	csr, e = c.NewQuery().From(new(GeneralInfoDetails).TableName()).Cursor(nil)
	e = csr.Fetch(&generalInfoDetails, 0, false)
	ErrorHandler(e, "generateMaintenanceDataBrowser")
	csr.Close()

	// Get GeneralInfoActualFuelConsumption
	csr, e = c.NewQuery().From(new(GeneralInfoActualFuelConsumption).TableName()).Cursor(nil)
	e = csr.Fetch(&generalInfoActualFuelConsumption, 0, false)
	ErrorHandler(e, "generateMaintenanceDataBrowser")
	csr.Close()

	// Get NewEquipmentType
	csr, e = c.NewQuery().Select("EquipmentType, NewEquipmentGroup").From(new(NewEquipmentType).TableName()).Where(dbox.Ne("NewEquipmentGroup", "Disregard")).Cursor(nil)
	e = csr.Fetch(&newEquipmentType, 0, false)
	ErrorHandler(e, "generateMaintenanceDataBrowser")
	csr.Close()

	constructEquipmentType(newEquipmentType)

	wg := &sync.WaitGroup{}
	wg.Add(len(plants))

	for _, plant := range plants {
		go doGenerateMaintenanceDataBrowser(wg, d, plant)
	}

	wg.Wait()

	return
}

func doGenerateMaintenanceDataBrowser(wg *sync.WaitGroup, d *GenDataBrowser, plant PowerPlantCoordinates) {
	var e error
	ctx := d.BaseController.Ctx
	c := ctx.Connection

	saveDatas := []orm.IModel{}
	plantCodeStr := plant.PlantCode
	dataCount := 0

	tmpPlantCondition := conditions.Get(plantCodeStr)
	if tmpPlantCondition != nil {
		plantCondition := tmpPlantCondition.(tk.M)

		length := plantCondition.GetInt("length")
		dets := plantCondition.Get("det").([]tk.M)
		dets = append(dets, tk.M{})

		turbinesCodes := []string{}
		genIDTempTable := tk.GenerateRandomString("", 20)

		detsLen := len(dets) - 1

		tk.Printf("detsLen: %v \n", detsLen)

		for i, det := range dets {
			freeQuery := false
			tk.Println(i)
			tk.Printf("dets: %#v \n", det)
			assets := []FunctionalLocation{}
			systemAssets := []FunctionalLocation{}
			desc := det.GetString("desc")

			query := []*dbox.Filter{}
			query = append(query, dbox.Contains("FunctionalLocationCode", plantCodeStr))
			query = append(query, dbox.Eq("PIPI", plantCodeStr))

			queryStr := "select * from FunctionalLocation where FunctionalLocationCode like('%" + plantCodeStr + "%') and PIPI = '" + plantCodeStr + "' "

			if i == 0 {
				query = append(query, dbox.Contains("Description", desc))
				if plantCodeStr == "2220" {
					query = append(query, dbox.Lte("LEN(FunctionalLocationCode)", length))
				} else {
					query = append(query, dbox.Eq("LEN(FunctionalLocationCode)", length))
				}
			} else if i == 1 && detsLen == 2 {
				query = append(query, dbox.Contains("Description", desc))
				if plantCodeStr == "2220" {
					query = append(query, dbox.Lte("LEN(FunctionalLocationCode)", length))
				} else {
					query = append(query, dbox.Eq("LEN(FunctionalLocationCode)", length))
				}
			} else {
				tk.Printf("turbinesCodes: %v \n", len(turbinesCodes))
				if len(turbinesCodes) > 1 {
					maxLength := 1000
					maxLoop := tk.ToInt(tk.ToFloat64(len(turbinesCodes)/maxLength, 0, tk.RoundingUp), tk.RoundingUp)

					for i := 0; i <= maxLoop; i++ {
						datas := []string{}

						if i != maxLoop {
							datas = turbinesCodes[i*maxLength : (i*maxLength)+maxLength]

						} else {
							datas = turbinesCodes[i*maxLength:]
						}

						tmpNotIn := []orm.IModel{}
						for _, val := range datas {
							tmpData := new(GenDataBrowserNotInTmp)
							tmpData.ID = genIDTempTable
							tmpData.FLCode = val
							tmpNotIn = append(tmpNotIn, tmpData)
						}
						// tk.Printf("tmpNotIn: %v \n", len(tmpNotIn))
						e = ctx.InsertBulk(tmpNotIn)
						ErrorHandler(e, "generateMaintenanceDataBrowser")
					}

					queryStr = queryStr + " and FunctionalLocationCode not in(select flcode from GenDataBrowserNotInTmp where ID = '" + genIDTempTable + "')"
					freeQuery = true
				}
			}

			if freeQuery {
				csrDet, e := c.NewQuery().Command("freequery", tk.M{}.Set("syntax", queryStr)).Cursor(nil)
				ErrorHandler(e, "generateMaintenanceDataBrowser")

				e = csrDet.Fetch(&assets, 0, false)
				ErrorHandler(e, "generateMaintenanceDataBrowser")
				csrDet.Close()
			} else {
				csrDet, e := c.NewQuery().From(new(FunctionalLocation).TableName()).Where(dbox.And(query...)).Cursor(nil)
				ErrorHandler(e, "generateMaintenanceDataBrowser")

				e = csrDet.Fetch(&assets, 0, false)
				ErrorHandler(e, "generateMaintenanceDataBrowser")
				csrDet.Close()
			}

			tk.Printf("-- assets: %v \n", len(assets))

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

							for _, relasset := range relatedAssets {
								isTurbineSystem := false
								if relasset.FunctionalLocationCode == asset.FunctionalLocationCode && i != 2 {
									isTurbineSystem = true
								}

								newEquipment := d.getNewEquipmentType(relasset.ObjectType, isTurbineSystem)

								if newEquipment != "" {

									if i != 2 {
										turbinesCodes = append(turbinesCodes, relasset.FunctionalLocationCode)
									}

									for _, year := range years {
										data := DataBrowser{}
										data.PeriodYear = year
										data.FunctionalLocation = relasset.FunctionalLocationCode
										data.FLDescription = relasset.Description
										data.IsTurbine = false

										// tk.Printf("%v | %v | ", relasset.FunctionalLocationCode, i)

										if relasset.FunctionalLocationCode == asset.FunctionalLocationCode && i != 2 {
											data.IsTurbine = true
											// tk.Println(" isTurbine: TRUE")
										} else {
											// tk.Println(" isTurbine: FALSE")
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

										// e = ctx.Insert(&data)
										saveDatas = append(saveDatas, &data)
										dataCount++
										if len(saveDatas) == 1000 {
											mugen.Lock()
											e = ctx.InsertBulk(saveDatas)
											mugen.Unlock()
											ErrorHandler(e, "generateMaintenanceDataBrowser")
											tk.Printf("%v : %v \n", plantCodeStr, dataCount)

											saveDatas = []orm.IModel{}
										}

										/*ErrorHandler(e, "generateMaintenanceDataBrowser")
										tk.Println("save 1")*/
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

								// tk.Printf("%v | %v | ", asset.FunctionalLocationCode, i)

								if asset.FunctionalLocationCode == asset.FunctionalLocationCode && i != 2 {
									data.IsTurbine = true
									data.IsSystem = true
									// tk.Println(" isTurbine: TRUE")
								} else {
									// tk.Println(" isTurbine: FALSE")
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

								// e = ctx.Insert(&data)
								// dataCount++
								saveDatas = append(saveDatas, &data)
								dataCount++
								if len(saveDatas) == 1000 {
									mugen.Lock()
									e = ctx.InsertBulk(saveDatas)
									mugen.Unlock()
									ErrorHandler(e, "generateMaintenanceDataBrowser")
									tk.Printf("%v : %v \n", plantCodeStr, dataCount)

									saveDatas = []orm.IModel{}
								}

								/*ErrorHandler(e, "generateMaintenanceDataBrowser")
								tk.Println("save 2")*/

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

								for _, relasset := range relatedAssets {
									isTurbineSystem := false
									if (relasset.FunctionalLocationCode == asset.FunctionalLocationCode || relasset.FunctionalLocationCode == sysAsset.FunctionalLocationCode) && i != 2 {
										isTurbineSystem = true
									}

									newEquipment := d.getNewEquipmentType(relasset.ObjectType, isTurbineSystem)

									if newEquipment != "" {

										if i != 2 {
											turbinesCodes = append(turbinesCodes, relasset.FunctionalLocationCode)
										}

										for _, year := range years {
											data := DataBrowser{}
											data.PeriodYear = year
											data.FunctionalLocation = relasset.FunctionalLocationCode
											data.FLDescription = relasset.Description
											data.IsTurbine = false
											data.IsSystem = false

											// tk.Printf("%v | %v | ", relasset.FunctionalLocationCode, i)

											if relasset.FunctionalLocationCode == sysAsset.FunctionalLocationCode && i != 2 {
												data.IsTurbine = true
												data.IsSystem = true
												// tk.Println(" isTurbine: TRUE")
											} else {
												// tk.Println(" isTurbine: FALSE")
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

											// e = ctx.Insert(&data)
											saveDatas = append(saveDatas, &data)
											dataCount++
											if len(saveDatas) == 1000 {
												mugen.Lock()
												e = ctx.InsertBulk(saveDatas)
												mugen.Unlock()
												ErrorHandler(e, "generateMaintenanceDataBrowser")
												tk.Printf("%v : %v \n", plantCodeStr, dataCount)

												saveDatas = []orm.IModel{}
											}

											/*ErrorHandler(e, "generateMaintenanceDataBrowser")
											tk.Println("save 3")*/
										}
									}
								}
							}
						}
					} else {
						// another plant(s)

						query = []*dbox.Filter{}

						if i != detsLen {
							query = append(query, dbox.Contains("FunctionalLocationCode", asset.FunctionalLocationCode))
						} else {
							query = append(query, dbox.Eq("FunctionalLocationCode", asset.FunctionalLocationCode))
						}

						csrDet, e := c.NewQuery().From(new(FunctionalLocation).TableName()).Where(query...).Cursor(nil)
						ErrorHandler(e, "generateMaintenanceDataBrowser")
						e = csrDet.Fetch(&relatedAssets, 0, false)
						ErrorHandler(e, "generateMaintenanceDataBrowser")
						csrDet.Close()
						tk.Printf("-- related assets: %v \n", len(relatedAssets))

						for _, relasset := range relatedAssets {
							isTurbineSystem := false
							if relasset.FunctionalLocationCode == asset.FunctionalLocationCode && i != detsLen {
								isTurbineSystem = true
							}

							newEquipment := ""
							newEquipment = d.getNewEquipmentType(relasset.ObjectType, isTurbineSystem)

							if newEquipment != "" {

								if i != detsLen {
									turbinesCodes = append(turbinesCodes, relasset.FunctionalLocationCode)
								}

								for _, year := range years {
									_ = year
									data := DataBrowser{}
									data.PeriodYear = year
									data.FunctionalLocation = relasset.FunctionalLocationCode
									data.FLDescription = relasset.Description
									data.IsTurbine = false

									if relasset.FunctionalLocationCode == asset.FunctionalLocationCode && i != detsLen {
										data.IsTurbine = true
									} else {
										data.TurbineParent = asset.FunctionalLocationCode
									}

									data.AssetType = "Other"

									if i == 0 {
										data.AssetType = "Steam"
									} else if i == 1 && detsLen > 1 {
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

									// e = ctx.Insert(&data)
									saveDatas = append(saveDatas, &data)
									dataCount++
									if len(saveDatas) == 1000 {
										mugen.Lock()
										e = ctx.InsertBulk(saveDatas)
										mugen.Unlock()
										ErrorHandler(e, "generateMaintenanceDataBrowser")
										tk.Printf("%v : %v \n", plantCodeStr, dataCount)

										saveDatas = []orm.IModel{}
									}

									/*ErrorHandler(e, "generateMaintenanceDataBrowser")
									tk.Println("save 4")*/
								}
							}
						}
					}

				}
			}
		}

		e = c.NewQuery().Delete().From(new(GenDataBrowserNotInTmp).TableName()).SetConfig("multiexec", true).Where(dbox.Eq("ID", genIDTempTable)).Exec(nil)

		if len(saveDatas) > 0 {
			e = ctx.InsertBulk(saveDatas)
			ErrorHandler(e, "generateMaintenanceDataBrowser")
			tk.Printf("%v : %v \n", plantCodeStr, dataCount)
		}

	}

	mugen.Lock()
	dataCounts.Set(plantCodeStr, dataCount)
	mugen.Unlock()

	wg.Done()
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
