package controllers

import (
	"github.com/eaciit/dbox"
	. "github.com/eaciit/powerplant/sec/library/models"
	tk "github.com/eaciit/toolkit"
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
	if e != nil {
		tk.Println(e)
	}

	tk.Println("##Value Equation Data : DONE\n")
}

func (d *GenDataBrowser) generateMaintenanceDataBrowser(years []int, conditions tk.M) (e error) {
	ctx := d.BaseController.Ctx
	c := ctx.Connection

	tk.Println("Generating Maintenance Data Browser..")

	// Get fuelCost
	csr, e := c.NewQuery().From(new(FuelCost).TableName()).Cursor(nil)
	fuelCosts := []FuelCost{}
	e = csr.Fetch(&fuelCosts, 0, false)
	csr.Close()

	// Get plants
	csr, e = c.NewQuery().From(new(PowerPlantCoordinates).TableName()).Cursor(nil)
	plants := []PowerPlantCoordinates{}
	e = csr.Fetch(&plants, 0, false)
	csr.Close()

	// Get generalInfo
	csr, e = c.NewQuery().From(new(GeneralInfo).TableName()).Cursor(nil)
	generalInfos := []GeneralInfo{}
	e = csr.Fetch(&generalInfos, 0, false)
	csr.Close()

	for _, plant := range plants {
		assets := []FunctionalLocation{}
		systemAssets := []FunctionalLocation{}
		plantCodeStr := plant.PlantCode

		tmpPlantCondition := conditions.Get(plantCodeStr)
		if tmpPlantCondition != nil {
			plantCondition := tmpPlantCondition.(tk.M)

			length := plantCondition.GetInt("length")
			dets := plantCondition.Get("det").([]tk.M)
			dets = append(dets, tk.M{})

			turbinesCodes := []string{}

			for i, det := range dets {
				query := []*dbox.Filter{}
				query = append(query, dbox.Contains("FunctionalLocationCode", plantCodeStr))
				query = append(query, dbox.Eq("PIPI", plantCodeStr))

				if i == 1 {
					query = append(query, dbox.Contains("Description", det.GetString("desc")))
					query = append(query, dbox.Eq("LEN(FunctionalLocationCode)", length))
				} else if i == len(dets)-1 {
					query = append(query, dbox.Contains("Description", det.GetString("desc")))
					query = append(query, dbox.Eq("LEN(FunctionalLocationCode)", length))
				} else {
					if len(turbinesCodes) > 1 {
						query = append(query, dbox.Nin("FunctionalLocationCode", turbinesCodes))
					}
				}

				csrDet, err := c.NewQuery().From(new(FunctionalLocation).TableName()).Where(query...).Cursor(nil)

				if err != nil {
					e = err
				}
				e = csrDet.Fetch(&assets, 0, false)
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

								csrDet, err := c.NewQuery().From(new(FunctionalLocation).TableName()).Where(query...).Cursor(nil)

								if err != nil {
									e = err
								}
								e = csrDet.Fetch(&systemAssets, 0, false)
								csrDet.Close()
							}

							if systemAssets != nil || len(systemAssets) == 0 {
								query = []*dbox.Filter{}

								if i != 2 {
									query = append(query, dbox.Contains("FunctionalLocationCode", asset.FunctionalLocationCode))
								} else {
									query = append(query, dbox.Eq("FunctionalLocationCode", asset.FunctionalLocationCode))
								}

								csrDet, err := c.NewQuery().From(new(FunctionalLocation).TableName()).Where(query...).Cursor(nil)

								if err != nil {
									e = err
								}

								e = csrDet.Fetch(&relatedAssets, 0, false)
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

											if related.FunctionalLocationCode == asset.FunctionalLocationCode && i != 2 {
												data.IsTurbine = true
												data.TurbineParent = asset.FunctionalLocationCode
											}

											data.AssetType = "Other"

											if i == 0 {
												data.AssetType = "Steam"
											} else if 1 == 1 {
												data.AssetType = "Gas"
											}

											data.EquipmentType = newEquipment
											data.EquipmentTypeDescription = newEquipment
											data.Plant = plant

											if data.IsTurbine {
												info := GeneralInfo{}

											}
										}

									}
								}

							} else {

							}
						} else {

						}

					}
				}
			}
		}
	}

	return
}

func (d *GenDataBrowser) getNewEquipmentType(equipmentType string, isSystem bool) (retVal string) {
	ctx := d.BaseController.Ctx
	c := ctx.Connection

	res := []tk.M{}

	csr, err := c.NewQuery().Select("NewEquipmentGroup").From(new(NewEquipmentType).TableName()).
		Where(dbox.And(
			dbox.Eq("EquipmentType", equipmentType),
			dbox.Ne("NewEquipmentGroup", "Disregard"),
		),
		).
		Cursor(nil)

	e = csr.Fetch(&res, 0, false)
	csr.Close()

	if len(res) > 0 {
		retVal = res[0].GetString("NewEquipmentGroup")

		if ((retVal == "Default" || retVal == "System") && isSystem) || (retval != "" && retVal != "Default" && retVal != "System") {
			retVal = retVal
		} else {
			retval = ""
		}
	}

	return
}
