package controllers

import (
	"github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	. "github.com/eaciit/powerplant/sec/library/models"
	tk "github.com/eaciit/toolkit"
)

type ValueEquationController struct {
	*BaseController
}

func (c *ValueEquationController) Default(k *knot.WebContext) interface{} {
	if k.Session("userid") == nil {
		c.Redirect(k, "login", "default")
	}
	c.LoadPartial(k, "valueequation/browse.html",
		"valueequation/historicalvalueequation/index.html",
		"valueequation/historicalvalueequation/maintenance.html",
		"valueequation/historicalvalueequation/operating.html",
		"valueequation/historicalvalueequation/revenue.html",
		"valueequation/historicalvalueequation/availability.html",
		"valueequation/historicalvalueequation/outages.html",
		"valueequation/historicalvalueequation/summary.html",
		"valueequation/historicalvalueequation/dataquality.html",
		"valueequation/historicalvalueequation/performance.html",
		"valueequation/historicalvalueequation/assetworksummary.html")

	k.Config.OutputType = knot.OutputTemplate

	infos := PageInfo{}
	infos.PageId = "ValueEquation"
	infos.PageTitle = "Value Equation"
	infos.Breadcrumbs = make(map[string]string, 0)

	return infos
}

func (c *ValueEquationController) Initiate(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	var e error
	csr, e := c.Ctx.Find(new(PlantModel), tk.M{}.Set("skip", 0).Set("limit", 0))
	defer csr.Close()
	PlantList := make([]PlantModel, 0)
	e = csr.Fetch(&PlantList, 0, false)
	if e != nil {
		return e.Error()
	}
	csr, e = c.Ctx.Find(new(PhaseModel), tk.M{}.Set("skip", 0).Set("limit", 0))
	PhaseList := make([]PhaseModel, 0)
	e = csr.Fetch(&PhaseList, 0, false)
	if e != nil {
		return e.Error()
	}
	csr, e = c.Ctx.Find(new(UnitModel), tk.M{}.Set("skip", 0).Set("limit", 0))
	UnitList := make([]UnitModel, 0)
	e = csr.Fetch(&UnitList, 0, false)
	if e != nil {
		return e.Error()
	}
	result := tk.M{}
	result.Set("PlantList", PlantList)
	result.Set("PhaseList", PhaseList)
	result.Set("UnitList", UnitList)
	return ResultInfo(result, e)
}

func (c *ValueEquationController) GetUnitList(k *knot.WebContext) interface{} {
	d := struct {
		SelectedPlant string
	}{}
	e := k.GetPayload(&d)
	csr, e := c.Ctx.Find(new(MasterUnitPlant), tk.M{}.Set("where", dbox.Eq("Plant", d.SelectedPlant)))
	defer csr.Close()
	UnitPlantList := make([]MasterUnitPlant, 0)
	e = csr.Fetch(&UnitPlantList, 0, false)
	if e != nil {
		return e.Error()
	}
	return ResultInfo(UnitPlantList, e)
}
