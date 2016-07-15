package controllers

import (
	"strings"

	"github.com/eaciit/crowd"
	"github.com/eaciit/dbox"
	. "github.com/eaciit/powerplant/sec/consoleapp/generator/helpers"
	. "github.com/eaciit/powerplant/sec/library/models"
	tk "github.com/eaciit/toolkit"
)

// GenPlantMaster struct
type GenPlantMaster struct {
	*BaseController
}

// Generate Do the generate process
func (d *GenPlantMaster) Generate(base *BaseController) {
	if base != nil {
		d.BaseController = base
	}

	tk.Println("##Generating Plant Master..")
	e := d.generatePlantMaster()
	if e != nil {
		tk.Println(e)
	}

	tk.Println("##Plant Master : DONE\n")
}

// GeneratePlantMaster To Generate Master Plant
func (d *GenPlantMaster) generatePlantMaster() error {
	tk.Println("Generating..")
	ctx := d.BaseController.Ctx
	c := ctx.Connection
	query := []*dbox.Filter{}
	FunctionalLocationList := []FunctionalLocation{}

	query = append(query, dbox.Eq("LEN(FunctionalLocationCode)", 4))
	csr, e := c.NewQuery().From(new(FunctionalLocation).TableName()).Where(query...).Cursor(nil)

	defer csr.Close()

	if e != nil {
		return e
	}

	e = csr.Fetch(&FunctionalLocationList, 0, false)

	if e != nil {
		return e
	}

	PowerPlantInfoList := []PowerPlantInfo{}
	csr, e = c.NewQuery().From(new(PowerPlantInfo).TableName()).Cursor(nil)
	if e != nil {
		return e
	}

	e = csr.Fetch(&PowerPlantInfoList, 0, false)

	if e != nil {
		return e
	}

	RegenMP := new(RegenMasterPlant)
	temp := []PowerPlantInfo{}

	for _, plant := range FunctionalLocationList {

		RegenMP = new(RegenMasterPlant)
		temp = []PowerPlantInfo{}
		var tempValue float64

		plantDes := PlantNormalization(plant.Description)

		datas := crowd.From(&PowerPlantInfoList).Where(func(x interface{}) interface{} {
			return strings.Contains(strings.ToLower(x.(PowerPlantInfo).Name), strings.ToLower(plantDes))
		}).Exec().Result.Data().([]PowerPlantInfo)

		RegenMP.PlantCode = plant.FunctionalLocationCode
		RegenMP.PlantName = plantDes
		RegenMP.City = datas[0].City
		RegenMP.Province = datas[0].Province
		RegenMP.Region = datas[0].Region
		temp = crowd.From(&datas).Where(func(x interface{}) interface{} {
			return x.(PowerPlantInfo).FuelTypes_Crude
		}).Exec().Result.Data().([]PowerPlantInfo)

		if len(temp) > 0 {
			RegenMP.FuelTypes_Crude = true
		} else {
			RegenMP.FuelTypes_Crude = false
		}
		temp = crowd.From(&datas).Where(func(x interface{}) interface{} {
			return x.(PowerPlantInfo).FuelTypes_Diesel
		}).Exec().Result.Data().([]PowerPlantInfo)

		if len(temp) > 0 {
			RegenMP.FuelTypes_Diesel = true
		} else {
			RegenMP.FuelTypes_Diesel = false
		}

		temp = crowd.From(&datas).Where(func(x interface{}) interface{} {
			return x.(PowerPlantInfo).FuelTypes_Gas
		}).Exec().Result.Data().([]PowerPlantInfo)

		if len(temp) > 0 {
			RegenMP.FuelTypes_Gas = true
		} else {
			RegenMP.FuelTypes_Gas = false
		}

		temp = crowd.From(&datas).Where(func(x interface{}) interface{} {
			return x.(PowerPlantInfo).FuelTypes_Heavy
		}).Exec().Result.Data().([]PowerPlantInfo)

		if len(temp) > 0 {
			RegenMP.FuelTypes_Heavy = true
		} else {
			RegenMP.FuelTypes_Heavy = false
		}

		tempValue = crowd.From(&datas).Sum(func(x interface{}) interface{} {
			return x.(PowerPlantInfo).GasTurbineUnit
		}).Exec().Result.Sum

		RegenMP.GasTurbineUnit = int(tempValue)

		RegenMP.GasTurbineCapacity = crowd.From(&datas).Sum(func(x interface{}) interface{} {
			return x.(PowerPlantInfo).GasTurbineCapacity
		}).Exec().Result.Sum

		tempValue = crowd.From(&datas).Sum(func(x interface{}) interface{} {
			return x.(PowerPlantInfo).SteamUnit
		}).Exec().Result.Sum

		RegenMP.SteamUnit = int(tempValue)

		RegenMP.SteamCapacity = crowd.From(&datas).Sum(func(x interface{}) interface{} {
			return x.(PowerPlantInfo).SteamCapacity
		}).Exec().Result.Sum

		tempValue = crowd.From(&datas).Sum(func(x interface{}) interface{} {
			return x.(PowerPlantInfo).DieselUnit
		}).Exec().Result.Sum

		RegenMP.DieselUnit = int(tempValue)

		RegenMP.DieselCapacity = crowd.From(&datas).Sum(func(x interface{}) interface{} {
			return x.(PowerPlantInfo).DieselCapacity
		}).Exec().Result.Sum

		tempValue = crowd.From(&datas).Sum(func(x interface{}) interface{} {
			return x.(PowerPlantInfo).CombinedCycleUnit
		}).Exec().Result.Sum

		RegenMP.CombinedCycleUnit = int(tempValue)

		RegenMP.CombinedCycleCapacity = crowd.From(&datas).Sum(func(x interface{}) interface{} {
			return x.(PowerPlantInfo).CombinedCycleCapacity
		}).Exec().Result.Sum

		e := ctx.Insert(RegenMP)

		if e != nil {
			tk.Println(e.Error())
		}
	}

	return e
}
