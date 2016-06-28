package controllers

import (
	"github.com/eaciit/crowd"
	"github.com/eaciit/dbox"
	. "github.com/eaciit/powerplant/sec/consoleapp/generator/helpers"
	. "github.com/eaciit/powerplant/sec/library/models"
	tk "github.com/eaciit/toolkit"
	"strings"
	"time"
)

// GenSummaryData
type GenSummaryData struct {
	*BaseController
}

// Generate
func (s *GenSummaryData) Generate(base *BaseController) {
	var (
		e error
	)
	if base != nil {
		s.BaseController = base
	}

	tk.Println("##Generating Summary Data..")
	// e = s.generateSummaryData()
	// if e != nil {
	// 	tk.Println(e)
	// }
	e = s.generateDurationCostWorkOrderSummary()
	if e != nil {
		tk.Println(e)
	}
	tk.Println("##Summary Data : DONE\n")
}

// GenerateSummaryData
func (s *GenSummaryData) generateSummaryData() error {
	ctx := s.BaseController.Ctx
	c := ctx.Connection

	FunctionLocationList := []FunctionalLocation{}
	PowerPlantInfoList := []PowerPlantInfo{}
	PowerPlantInfos := []PowerPlantInfo{}
	//PlantInfo := PowerPlantInfo{}
	SummaryInfo := SummaryData{}

	//FunctionalLocationList
	csr, err := c.NewQuery().Command("procedure", tk.M{}.Set("name", "GetFunctionalLocation")).Cursor(nil)
	defer csr.Close()

	err = csr.Fetch(&FunctionLocationList, 0, false)

	if err != nil {
		tk.Println(err.Error())
	}

	//PowerPlantInfo
	csr, err = c.NewQuery().Select().From(new(PowerPlantInfo).TableName()).Cursor(nil)
	err = csr.Fetch(&PowerPlantInfos, 0, false)
	if err != nil {
		tk.Println(err.Error())
	}

	for _, loc := range FunctionLocationList {
		SummaryInfo.FunctionalLocation = loc.FunctionalLocationCode
		SummaryInfo.FLDescription = loc.Description
		SummaryInfo.SortField = loc.SortField
		SummaryInfo.ParentFL = loc.SupFunctionalLocation

		for _, checks := range FunctionLocationList {
			if strings.Contains(checks.FunctionalLocationCode, loc.FunctionalLocationCode) && checks.FunctionalLocationCode != loc.FunctionalLocationCode && loc.FunctionalLocationCode != "" {
				SummaryInfo.HasChild = true
				break
			}
		}

		csr, err = c.NewQuery().Command("procedure", tk.M{}.Set("name", "GetPowerPlantInfoBySortField").Set("parms", tk.M{}.Set("@SortField", SummaryInfo.SortField))).Cursor(nil)
		err = csr.Fetch(&PowerPlantInfoList, 0, false)
		if err != nil {
			//tk.Println(err.Error())
		}

		if SummaryInfo.SortField == "PP9" {
			tk.Printf("%#v,", SummaryInfo.SortField)
		}

		if len(PowerPlantInfoList) == 0 {
			for _, plantinfos := range PowerPlantInfos {
				Name := plantinfos.Name
				SplitedName := strings.Split(Name, " ")
				Desc := SummaryInfo.FLDescription
				SplitedDesc := strings.Split(Desc, " ")

				if SplitedName[0] == SplitedDesc[0] {
					SummaryInfo.Province = plantinfos.Province
					SummaryInfo.Region = plantinfos.Region
					SummaryInfo.City = plantinfos.City
					SummaryInfo.GasTurbineUnit = plantinfos.GasTurbineUnit
					SummaryInfo.GasTurbineCapacity = plantinfos.GasTurbineCapacity
					SummaryInfo.SteamUnit = plantinfos.SteamUnit
					SummaryInfo.SteamUnitCapacity = plantinfos.SteamCapacity
					SummaryInfo.DieselUnit = plantinfos.DieselUnit
					SummaryInfo.DieselUnitCapacity = plantinfos.DieselCapacity
					SummaryInfo.CombinedCycleUnit = plantinfos.CombinedCycleUnit
					SummaryInfo.CombinedCycleUnitCapacity = plantinfos.CombinedCycleCapacity
				}
			}
		} else if len(PowerPlantInfoList) > 0 {
			for _, plantinfos := range PowerPlantInfoList {
				SummaryInfo.Province = plantinfos.Province
				SummaryInfo.Region = plantinfos.Region
				SummaryInfo.City = plantinfos.City
				SummaryInfo.GasTurbineUnit = plantinfos.GasTurbineUnit
				SummaryInfo.GasTurbineCapacity = plantinfos.GasTurbineCapacity
				SummaryInfo.SteamUnit = plantinfos.SteamUnit
				SummaryInfo.SteamUnitCapacity = plantinfos.SteamCapacity
				SummaryInfo.DieselUnit = plantinfos.DieselUnit
				SummaryInfo.DieselUnitCapacity = plantinfos.DieselCapacity
				SummaryInfo.CombinedCycleUnit = plantinfos.CombinedCycleUnit
				SummaryInfo.CombinedCycleUnitCapacity = plantinfos.CombinedCycleCapacity
			}

			PowerPlantInfoList = []PowerPlantInfo{}
		}

		if SummaryInfo.Province != "" {
			tk.Printf("----------- Summary Data -----------\nProvince : %#v \n", SummaryInfo.Province)
			tk.Printf("Region : %#v \n", SummaryInfo.Region)
			tk.Printf("City : %#v \n", SummaryInfo.City)
			tk.Printf("GasTurbineUnit : %#v \n", SummaryInfo.GasTurbineUnit)
			tk.Printf("GasTurbineCapacity : %#v \n", SummaryInfo.GasTurbineCapacity)
			tk.Printf("SteamUnit : %#v \n", SummaryInfo.SteamUnit)
			tk.Printf("SteamUnitCapacity : %#v \n", SummaryInfo.SteamUnitCapacity)
			tk.Printf("DieselUnit : %#v \n", SummaryInfo.DieselUnit)
			tk.Printf("DieselUnitCapacity : %#v \n", SummaryInfo.DieselUnitCapacity)
			tk.Printf("CombinedCycleUnit : %#v \n", SummaryInfo.CombinedCycleUnit)
			tk.Printf("CombinedCycleUnitCapacity : %#v \n----------------------------------\n", SummaryInfo.CombinedCycleUnitCapacity)
		} else {
			tk.Printf("#")
		}

		query := tk.M{}
		params := tk.M{}
		query.Set("name", "SaveSummaryData")
		params.Set("@FunctionalLocation", SummaryInfo.FunctionalLocation)
		params.Set("@FLDescription", SummaryInfo.FLDescription)
		params.Set("@SortField", SummaryInfo.SortField)
		params.Set("@ParentFL", SummaryInfo.ParentFL)
		if SummaryInfo.HasChild {
			params.Set("@HasChild", 1)
		} else {
			params.Set("@HasChild", 0)
		}
		params.Set("@Province", SummaryInfo.Province)
		params.Set("@Region", SummaryInfo.Region)
		params.Set("@City", SummaryInfo.City)
		params.Set("@GasTurbineUnit", SummaryInfo.GasTurbineUnit)
		params.Set("@GasTurbineCapacity", SummaryInfo.GasTurbineCapacity)
		params.Set("@SteamUnit", SummaryInfo.SteamUnit)
		params.Set("@SteamUnitCapacity", SummaryInfo.SteamUnitCapacity)
		params.Set("@DieselUnit", SummaryInfo.DieselUnit)
		params.Set("@DieselUnitCapacity", SummaryInfo.DieselUnitCapacity)
		params.Set("@CombinedCycleUnit", SummaryInfo.CombinedCycleUnit)
		params.Set("@CombinedCycleUnitCapacity", SummaryInfo.CombinedCycleUnitCapacity)
		query.Set("parms", params)
		csr, err = c.NewQuery().Command("procedure", query).Cursor(nil)
		res := tk.M{}
		err = csr.Fetch(&res, 0, false)
		if err != nil {
			tk.Println(err.Error())
		}

		SummaryInfo = SummaryData{}

	}

	return nil
}

func (s *GenSummaryData) generateDurationCostWorkOrderSummary() error {
	ctx := s.BaseController.Ctx
	c := ctx.Connection
	var (
		query []*dbox.Filter
	)
	tk.Println("Generating Duration Cost Work Order Summary..")
	Years := []int{2013, 2014, 2015}
	query = []*dbox.Filter{}

	EqTypes := []MappedEquipmentType{}
	csr, e := c.NewQuery().From(new(MappedEquipmentType).TableName()).Cursor(nil)
	defer csr.Close()

	if e != nil {
		return e
	}

	e = csr.Fetch(&EqTypes, 0, false)
	if e != nil {
		return e
	}
	for _, year := range Years {
		query = []*dbox.Filter{}
		query = append(query, dbox.Gte("Period", time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)))
		query = append(query, dbox.Lt("Period", time.Date((year+1), 1, 1, 0, 0, 0, 0, time.UTC)))

		MaintenanceCostByHours := []MaintenanceCostByHour{}
		csr, e := c.NewQuery().From(new(MaintenanceCostByHour).TableName()).Where(query...).Cursor(nil)
		if e != nil {
			return e
		}
		e = csr.Fetch(&MaintenanceCostByHours, 0, false)
		if e != nil {
			return e
		}
		csr.Close()

		MaintenanceCostList := []MaintenanceCost{}
		csr, e = c.NewQuery().From(new(MaintenanceCost).TableName()).Where(query...).Cursor(nil)
		if e != nil {
			return e
		}
		e = csr.Fetch(&MaintenanceCostList, 0, false)
		if e != nil {
			return e
		}
		csr.Close()

		Plants := crowd.From(&MaintenanceCostByHours).Group(func(x interface{}) interface{} {
			return x.(MaintenanceCostByHour).Plant
		}, nil).Exec().Result.Data().([]crowd.KV)

		for _, p := range Plants {
			plant := p.Key.(string)
			EqType := crowd.From(&MaintenanceCostByHours).Where(func(x interface{}) interface{} {
				return x.(MaintenanceCostByHour).Plant == plant
			}).Group(func(x interface{}) interface{} {
				return x.(MaintenanceCostByHour).EquipmentType
			}, nil).Exec().Result.Data().([]crowd.KV)

			for _, eqt := range EqType {
				eq := eqt.Key.(string)
				ActType := crowd.From(&MaintenanceCostByHours).Where(func(x interface{}) interface{} {
					o := x.(MaintenanceCostByHour)
					return o.Plant == plant && o.EquipmentType == eq
				}).Group(func(x interface{}) interface{} {
					return x.(MaintenanceCostByHour).MaintenanceActivityType
				}, nil).Exec().Result.Data().([]crowd.KV)

				for _, a := range ActType {
					act := a.Key.(string)
					OrderType := crowd.From(&MaintenanceCostByHours).Where(func(x interface{}) interface{} {
						o := x.(MaintenanceCostByHour)
						return o.Plant == plant && o.EquipmentType == eq && o.MaintenanceActivityType == act
					}).Group(func(x interface{}) interface{} {
						return x.(MaintenanceCostByHour).OrderType
					}, nil).Exec().Result.Data().([]crowd.KV)

					for _, o := range OrderType {
						order := o.Key.(string)
						d := new(WODurationSummary)
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
							EqTypeDesc := crowd.From(&MaintenanceCostByHours).Where(func(x interface{}) interface{} {
								o := x.(MaintenanceCostByHour)
								return o.Plant == plant && o.EquipmentType == eq
							}).Exec().Result.Data().([]MaintenanceCostByHour)
							if len(EqTypeDesc) > 0 {
								d.EquipmentTypeDescription = EqTypeDesc[0].EquipmentTypeDesc
							}
						}

						d.ActivityType = act
						d.Plant = PlantNormalization(plant)
						d.PlanValue = crowd.From(&MaintenanceCostByHours).Where(func(x interface{}) interface{} {
							o := x.(MaintenanceCostByHour)
							return o.Plant == plant && o.EquipmentType == eq && o.MaintenanceActivityType == act
						}).Sum(func(x interface{}) interface{} {
							return x.(MaintenanceCostByHour).PlanVal
						}).Exec().Result.Sum
						d.ActualValue = crowd.From(&MaintenanceCostByHours).Where(func(x interface{}) interface{} {
							o := x.(MaintenanceCostByHour)
							return o.Plant == plant && o.EquipmentType == eq && o.MaintenanceActivityType == act
						}).Sum(func(x interface{}) interface{} {
							return x.(MaintenanceCostByHour).Actual
						}).Exec().Result.Sum
						d.WOCount = len(OrderType)
						d.Cost = crowd.From(&MaintenanceCostList).Where(func(x interface{}) interface{} {
							o := x.(MaintenanceCost)
							return o.Plant == plant && o.EquipmentType == eq && o.MaintenanceActivityType == act
						}).Sum(func(x interface{}) interface{} {
							return x.(MaintenanceCost).PeriodTotalActual
						}).Exec().Result.Sum
						_, e := ctx.InsertOut(d)
						tk.Println("#")
						if e != nil {
							tk.Println(e)
							break
						}

					}

				}
			}
		}
	}
	return e
}
