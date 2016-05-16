package controllers

import (
	"github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	. "github.com/eaciit/powerplant/sec/webapp/models"
	tk "github.com/eaciit/toolkit"
)

type HypothesisController struct {
	*BaseController
}

func (c *HypothesisController) Initiate(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	type ReturnValue struct {
		PlantList        []PlantModel
		PlantListH2      []PlantModel
		EQTypeList       []MasterEquipmentType
		MROElementList   []MasterMROElement
		OrderTypeList    []MasterOrderType
		ActivityTypeList []MasterActivityType
		FailureCode      []MasterFailureCode
		UnitList         []tk.M
	}

	var (
		Result ReturnValue
	)

	d := struct {
		selectedPlant []string
	}{}

	e := k.GetPayload(&d)

	curr, _ := c.DB().Find(&PlantModel{}, nil)
	defer curr.Close()
	e = curr.Fetch(&Result.PlantList, 0, false)

	filter := tk.M{}

	filter.Set("where", dbox.Nin("Plant", "Qurayyah CC", "PP9"))

	curr = nil
	curr, e = c.DB().Find(&PlantModel{}, filter)
	e = curr.Fetch(&Result.PlantListH2, 0, false)

	curr = nil
	curr, e = c.DB().Find(&MasterEquipmentType{}, nil)
	e = curr.Fetch(&Result.EQTypeList, 0, false)

	curr = nil
	curr, e = c.DB().Find(&MasterMROElement{}, nil)
	e = curr.Fetch(&Result.MROElementList, 0, false)

	curr = nil
	curr, e = c.DB().Find(&MasterOrderType{}, nil)
	e = curr.Fetch(&Result.OrderTypeList, 0, false)

	curr = nil
	curr, e = c.DB().Find(&MasterActivityType{}, nil)
	e = curr.Fetch(&Result.ActivityTypeList, 0, false)

	curr = nil
	curr, e = c.DB().Find(&MasterFailureCode{}, nil)
	e = curr.Fetch(&Result.FailureCode, 0, false)

	var filter1 []*dbox.Filter

	var pipes []tk.M

	if len(d.selectedPlant) != 0 {
		filter1 = append(filter1, dbox.In("Plant", d.selectedPlant))
	}

	pipes = append(pipes, tk.M{"$group": tk.M{
		"_id": "$unit",
	}})

	curr = nil
	curr, e = c.DB().Connection.NewQuery().
		From("MasterUnitPlant").
		Where(filter1...).
		Group("Unit").
		Cursor(nil)

	e = curr.Fetch(&Result.UnitList, 0, true)
	_ = e

	return ResultInfo(Result, e)
}