package models

import (
	"github.com/eaciit/crowd"
	"github.com/eaciit/dbox"
	"github.com/eaciit/orm"
	. "github.com/eaciit/powerplant/sec/library/models"
	tk "github.com/eaciit/toolkit"
	"log"
	"strings"
)

type RegenMasterPlant1 struct {
}

func (r *RegenMasterPlant1) GeneratePlantMaster(ctx *orm.DataContext) error {
	tk.Println("Generating Plant Master..")
	c := ctx.Connection

	query := []*dbox.Filter{}
	FunctionalLocationList := []FunctionalLocation{}

	query = append(query, dbox.Eq("LEN(FunctionalLocationCode)", 4))
	csr, e := c.NewQuery().From(new(FunctionalLocation).TableName()).Where(query...).Cursor(nil)

	if e != nil {
		return e
	} else {
		defer csr.Close()
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
		var temp1 float64

		plantDes := r.PlantNormalization(plant.Description)

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
			return x.(PowerPlantInfo).FuleTypes_Diesel
		}).Exec().Result.Data().([]PowerPlantInfo)

		if len(temp) > 0 {
			RegenMP.FuelTypes_Diesel = true
		} else {
			RegenMP.FuelTypes_Diesel = false
		}

		temp = crowd.From(&datas).Where(func(x interface{}) interface{} {
			return x.(PowerPlantInfo).FuleTypes_Gas
		}).Exec().Result.Data().([]PowerPlantInfo)

		if len(temp) > 0 {
			RegenMP.FuelTypes_Gas = true
		} else {
			RegenMP.FuelTypes_Gas = false
		}

		temp = crowd.From(&datas).Where(func(x interface{}) interface{} {
			return x.(PowerPlantInfo).FuleTypes_Heavy
		}).Exec().Result.Data().([]PowerPlantInfo)

		if len(temp) > 0 {
			RegenMP.FuelTypes_Heavy = true
		} else {
			RegenMP.FuelTypes_Heavy = false
		}

		temp1 = crowd.From(&datas).Sum(func(x interface{}) interface{} {
			return x.(PowerPlantInfo).GasTurbineUnit
		}).Exec().Result.Sum

		RegenMP.GasTurbineUnit = int(temp1)

		RegenMP.GasTurbineCapacity = crowd.From(&datas).Sum(func(x interface{}) interface{} {
			return x.(PowerPlantInfo).GasTurbineCapacity
		}).Exec().Result.Sum

		temp1 = crowd.From(&datas).Sum(func(x interface{}) interface{} {
			return x.(PowerPlantInfo).SteamUnit
		}).Exec().Result.Sum

		RegenMP.SteamUnit = int(temp1)

		RegenMP.SteamCapacity = crowd.From(&datas).Sum(func(x interface{}) interface{} {
			return x.(PowerPlantInfo).SteamCapacity
		}).Exec().Result.Sum

		temp1 = crowd.From(&datas).Sum(func(x interface{}) interface{} {
			return x.(PowerPlantInfo).DieselUnit
		}).Exec().Result.Sum

		RegenMP.DieselUnit = int(temp1)

		RegenMP.DieselCapacity = crowd.From(&datas).Sum(func(x interface{}) interface{} {
			return x.(PowerPlantInfo).DieselCapacity
		}).Exec().Result.Sum

		temp1 = crowd.From(&datas).Sum(func(x interface{}) interface{} {
			return x.(PowerPlantInfo).CombinedCycleUnit
		}).Exec().Result.Sum

		RegenMP.CombinedCycleUnit = int(temp1)

		RegenMP.CombinedCycleCapacity = crowd.From(&datas).Sum(func(x interface{}) interface{} {
			return x.(PowerPlantInfo).CombinedCycleCapacity
		}).Exec().Result.Sum

		e := ctx.Insert(RegenMP)

		if e != nil {
			log.Println("kkkkkkkk " + e.Error())
		}
		break
	}

	return e
}

func (r *RegenMasterPlant1) PlantNormalization(PlantName string) string {
	retVal := ""
	//switch (PlantName)
	//{
	//    case "Rabigh PP": retVal = "Rabigh"; break;
	//    case "QURAYYAH": retVal = "Qurayyah"; break;
	//    case "GHZLAN": retVal = "Ghazlan"; break;
	//    case "QURAYYAH CC": retVal = "Qurayyah CC"; break;
	//    case "Shuaiba Power Plant": retVal = "Shoaiba"; break;
	//    case "RABIGH POWER PLANT": retVal = "Rabigh"; break;
	//    case "Qurayyah Power Plant": retVal = "Qurayyah"; break;
	//    case "Qurayyah Steam": retVal = "Qurayyah"; break;
	//    case "GHAZLAN POWER PLANT": retVal = "Ghazlan"; break;
	//    default: retVal = PlantName; break;
	//}

	switch PlantName {
	case "POWER PLANT #9":
		retVal = "PP9"
	case "RABIGH POWER PLANT":
		retVal = "Rabigh"
	case "Rabigh 2":
		retVal = "Rabigh"
	case "Rabigh PP":
		retVal = "Rabigh"
	case "Shuaiba Power Plant":
		retVal = "Shoaiba"
	case "Sha'iba (CC)":
		retVal = "Shoaiba"
	case "Sha'iba (SEC)":
		retVal = "Shoaiba"
	case "GHAZLAN POWER PLANT":
		retVal = "Ghazlan"
	case "GHZLAN":
		retVal = "Ghazlan"
	case "Qurayyah Power Plant":
		retVal = "Qurayyah"
	case "Qurayyah -Steam":
		retVal = "Qurayyah"
	case "Qurayyah Combined Cycle Power Plant":
		retVal = "Qurayyah CC"
	case "Qurayyah- Combined Cycle":
		retVal = "Qurayyah CC"
	case "QurayyahCC":
		retVal = "Qurayyah CC"
	case "QURAYYAH CC":
		retVal = "Qurayyah CC"
	case "QurayyahPP":
		retVal = "Qurayyah"
	case "Qurayyah Steam":
		retVal = "Qurayyah"
	case "QURAYYAH":
		retVal = "Qurayyah"
	default:
		retVal = PlantName
	}

	return retVal
}
