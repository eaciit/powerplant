package controllers

import (
	. "github.com/eaciit/powerplant/sec/webapp/models"
	//"github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
)

type ScenarioSimulation struct {
	*BaseController
}

func (s *ScenarioSimulation) Default(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputTemplate
	s.LoadPartial(k)
	infos := PageInfo{}
	infos.PageId = "ScenarioSimulation"
	infos.PageTitle = "Scenario Simulation"
	infos.Breadcrumbs = make(map[string]string, 0)

	return infos
}

func (s *ScenarioSimulation) Initiate(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	csr, e := s.DB().Find(new(PlantModel), tk.M{}.Set("skip", 0).Set("limit", 0))
	PlantList := make([]PlantModel, 0)
	e = csr.Fetch(&PlantList, 0, false)
	if e != nil {
		return e.Error()
	}

	csr, e = s.DB().Find(new(UnitModel), tk.M{}.Set("skip", 0).Set("limit", 0))
	UnitList := make([]UnitModel, 0)
	e = csr.Fetch(&UnitList, 0, false)
	if e != nil {
		return e.Error()
	}

	defer csr.Close()

	result := tk.M{}
	result.Set("PlantList", PlantList)
	result.Set("UnitList", UnitList)
	return ResultInfo(result, e)
}

/*func (s *ScenarioSimulation) GetData(k *knot.WebContext) interface{} {

}

func (s *ScenarioSimulation) GetDataSimulation(k *knot.WebContext) interface{} {

}

func (s *ScenarioSimulation) RemoveData(k *knot.WebContext) interface{} {

}

func (s *ScenarioSimulation) SaveData(k *knot.WebContext) interface{} {

}*/
