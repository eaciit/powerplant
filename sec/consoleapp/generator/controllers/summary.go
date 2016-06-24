package controllers

import (
	"strings"

	. "github.com/eaciit/powerplant/sec/library/models"
	tk "github.com/eaciit/toolkit"
)

// GenSummaryData
type GenSummaryData struct {
	*BaseController
}

// Generate
func (s *GenSummaryData) Generate(base *BaseController) {
	if base != nil {
		s.BaseController = base
	}

	tk.Println("##Generating Summary Data..")
	e := s.generateSummaryData()
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
