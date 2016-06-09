package controllers

import (
	"github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	//. "github.com/eaciit/powerplant/sec/webapp/models"
	. "github.com/eaciit/powerplant/sec/library/models"
	tk "github.com/eaciit/toolkit"
)

type HypothesisController struct {
	*BaseController
}

func (c *HypothesisController) Initiate(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	type ReturnValue struct {
		PlantList        []MasterPlant
		PlantListH2      []MasterPlant
		EQTypeList       []MasterEquipmentType
		MROElementList   []MasterMROElement
		OrderTypeList    []MasterOrderType
		ActivityTypeList []MasterActivityType
		FailureCode      []MasterFailureCode
		UnitList         []tk.M
	}

	var (
		Result ReturnValue
		e      error
	)

	d := struct {
		selectedPlant []string
	}{}

	e = k.GetPayload(&d)

	// tk.Printf("d: %#v \n\n", d)

	curr, _ := c.DB().Find(&MasterPlant{}, nil)

	defer curr.Close()

	e = curr.Fetch(&Result.PlantList, 0, false)

	// tk.Printf("Result.PlantList: %#v \n\n", Result.PlantList)

	filter := tk.M{}

	filter.Set("where", dbox.Nin("Plant", "Qurayyah CC", "PP9"))

	curr = nil
	curr, e = c.DB().Find(&MasterPlant{}, filter)
	e = curr.Fetch(&Result.PlantListH2, 0, false)

	// tk.Printf("Result.PlantListH2: %#v \n\n", Result.PlantListH2)

	curr = nil
	curr, e = c.DB().Find(&MasterEquipmentType{}, nil)
	e = curr.Fetch(&Result.EQTypeList, 0, false)

	// tk.Printf("Result.EQTypeList: %#v \n\n", Result.EQTypeList)

	curr = nil
	curr, e = c.DB().Find(&MasterMROElement{}, nil)
	e = curr.Fetch(&Result.MROElementList, 0, false)

	// tk.Printf("Result.MROElementList: %#v \n\n", Result.MROElementList)

	curr = nil
	curr, e = c.DB().Find(&MasterOrderType{}, nil)
	e = curr.Fetch(&Result.OrderTypeList, 0, false)

	// tk.Printf("Result.OrderTypeList: %#v \n\n", Result.OrderTypeList)

	curr = nil
	curr, e = c.DB().Find(&MasterActivityType{}, nil)
	e = curr.Fetch(&Result.ActivityTypeList, 0, false)

	// tk.Printf("Result.ActivityTypeList: %#v \n\n", Result.ActivityTypeList)

	curr = nil
	curr, e = c.DB().Find(&MasterFailureCode{}, nil)
	e = curr.Fetch(&Result.FailureCode, 0, false)

	// tk.Printf("Result.FailureCode: %#v \n\n", Result.FailureCode)

	var filter1 []*dbox.Filter

	if len(d.selectedPlant) != 0 {
		filter1 = append(filter1, dbox.In("Plant", d.selectedPlant))
	} else {
		filter1 = append(filter1, dbox.Eq("1", "1"))
	}

	/*pipes = append(pipes, tk.M{"$group": tk.M{
		"_id": "$unit",
	}})*/

	curr = nil

	curr, e = c.DB().Connection.NewQuery().
		Select("Unit").
		From(new(MasterUnitPlant).TableName()).
		Where(filter1...).
		Group("Unit").
		Cursor(nil)

	e = curr.Fetch(&Result.UnitList, 0, true)

	tk.Printf("Result.UnitList: %#v \n\n", Result.UnitList)

	return ResultInfo(Result, e)
}
