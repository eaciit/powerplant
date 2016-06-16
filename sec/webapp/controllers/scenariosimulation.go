package controllers

import (
	"github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	. "github.com/eaciit/powerplant/sec/library/models"
	tk "github.com/eaciit/toolkit"
	//"gopkg.in/mgo.v2/bson"
	"strconv"
	"time"
)

type ScenarioSimulation struct {
	*BaseController
}

type Sim struct {
	Revenue         string
	LaborCost       string
	MaterialCost    string
	ServiceCost     string
	OperatingCost   string
	MaintenanceCost string
	ValueEquation   string
}

type Simfloat struct {
	Revenue         float64
	LaborCost       float64
	MaterialCost    float64
	ServiceCost     float64
	OperatingCost   float64
	MaintenanceCost float64
	ValueEquation   float64
}

type SelScenario struct {
	ID    string
	Name  string
	Value float64
}

type SeleScenario struct {
	ID    string
	Name  string
	Value string
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
	//csr, e := s.DB().Find(new(PlantModel), tk.M{}.Set("skip", 0).Set("limit", 0))
	csr, e := s.DB().Connection.NewQuery().Select("Plant").From("MasterPlant").Cursor(nil)
	PlantList := make([]PlantModel, 0)
	e = csr.Fetch(&PlantList, 0, true)
	if e != nil {
		return e.Error()
	}

	/*csr, e = s.DB().Find(new(UnitModel), tk.M{}.Set("skip", 0).Set("limit", 0))*/
	csr, e = s.DB().Connection.NewQuery().Select("Unit").From("MasterUnit").Cursor(nil)
	UnitList := make([]UnitModel, 0)
	e = csr.Fetch(&UnitList, 0, true)
	if e != nil {
		return e.Error()
	}

	defer csr.Close()

	result := tk.M{}
	result.Set("PlantList", PlantList)
	result.Set("UnitList", UnitList)
	return ResultInfo(result, e)
}

func (s *ScenarioSimulation) GetData(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	status := ""
	msg := ""
	res := []*ValueEquation{}

	payloads := struct {
		StartPeriod         string
		EndPeriod           string
		SelectedPlant       []string
		SelectedUnit        []string
		SelectedSimulation  string
		SelectedDescription string
	}{}

	e := k.GetPayload(&payloads)
	if e != nil {
		s.WriteLog(e)
	}

	start, _ := time.Parse(time.RFC3339, payloads.StartPeriod)
	end, _ := time.Parse(time.RFC3339, payloads.EndPeriod)

	query := []*dbox.Filter{}
	query = append(query, dbox.Gte("Dates", start))
	query = append(query, dbox.Lte("Dates", end))

	if payloads.SelectedPlant != nil && len(payloads.SelectedPlant) > 0 {
		query = append(query, dbox.In("Plant", payloads.SelectedPlant))
	}

	if payloads.SelectedUnit != nil && len(payloads.SelectedUnit) > 0 {
		query = append(query, dbox.In("Unit", payloads.SelectedUnit))
	}
	csr, err := s.DB().Connection.NewQuery().From("ValueEquation").Select("Plant", "AvgNetGeneration", "VOMR", "TotalOutageDuration", "TotalDuration", "OperatingCost", "Revenue", "TotalLabourCost", "TotalMaterialCost", "TotalServicesCost", "TotalFuelCost").Where(query...).Cursor(nil)
	if err != nil {
		msg = err.Error()
		status = "NOK"
	} else {
		status = "OK"
	}
	err = csr.Fetch(&res, 0, true)
	defer csr.Close()

	result := tk.M{}
	result.Set("Status", status)
	result.Set("Message", msg)
	result.Set("Data", res)

	return result
}

type SimulationData struct {
	Id               int
	Start_Period     string
	End_Period       string
	Name             string
	Description      string
	SelectedPlant    []string
	SelectedUnit     []string
	SelectedScenario []SelScenario
	HistoricResult   Sim
	FutureResult     Sim
	Differential     Sim
	Last_Update      string
}

func (s *ScenarioSimulation) GetDataSimulation(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	result, status, msg := tk.M{}, "", ""

	simDataTmp := []tk.M{}
	filter := tk.M{}
	filter.Set("name", "GetAllScenarioSimulation")
	csr, err := s.DB().Connection.NewQuery().Command("procedure", filter).Cursor(nil)
	defer csr.Close()

	err = csr.Fetch(&simDataTmp, 0, true)

	if err != nil {
		msg = err.Error()
		status = "NOK"
	} else {
		status = "OK"
	}

	SimulationDatas := constructScenarionSimulation(simDataTmp)

	result.Set("Status", status)
	result.Set("Message", msg)
	result.Set("Data", SimulationDatas)

	return result
}

func constructScenarionSimulation(tmp []tk.M) (result []SimulationData) {

	for _, val := range tmp {
		simData := SimulationData{}
		differential := Sim{}
		historicres := Sim{}
		futureres := Sim{}
		existSimData := false

		id := val.GetInt("id")
		desc := val.GetString("description")
		end := val.GetString("endperiod")
		start := val.GetString("startperiod")
		name := val.GetString("name")
		last := val.GetString("lastupdate")

		for _, res := range result {
			if res.Id == id {
				existSimData = true
				break
			}
		}

		if !existSimData {
			simData.Id = id
			simData.Description = desc
			simData.End_Period = end
			simData.Start_Period = start
			simData.Name = name
			simData.Last_Update = last
			differential.Revenue = strconv.FormatFloat(val.GetFloat64("differentialrevenue"), 'f', 6, 64)
			differential.LaborCost = strconv.FormatFloat(val.GetFloat64("differentiallaborcost"), 'f', 6, 64)
			differential.MaintenanceCost = strconv.FormatFloat(val.GetFloat64("differentialmaintenacecost"), 'f', 6, 64)
			differential.OperatingCost = strconv.FormatFloat(val.GetFloat64("differentialoperationcost"), 'f', 6, 64)
			differential.ServiceCost = strconv.FormatFloat(val.GetFloat64("differentialservicecost"), 'f', 6, 64)
			differential.ValueEquation = strconv.FormatFloat(val.GetFloat64("differentialvalueequation"), 'f', 6, 64)
			differential.MaterialCost = strconv.FormatFloat(val.GetFloat64("differentialmaterialcost"), 'f', 6, 64)
			simData.Differential = differential
			historicres.Revenue = strconv.FormatFloat(val.GetFloat64("historicresultrevenue"), 'f', 6, 64)
			historicres.LaborCost = strconv.FormatFloat(val.GetFloat64("historicresultlaborcost"), 'f', 6, 64)
			historicres.MaintenanceCost = strconv.FormatFloat(val.GetFloat64("historicresultmaintenacecost"), 'f', 6, 64)
			historicres.OperatingCost = strconv.FormatFloat(val.GetFloat64("historicresultoperationcost"), 'f', 6, 64)
			historicres.ServiceCost = strconv.FormatFloat(val.GetFloat64("historicresultservicecost"), 'f', 6, 64)
			historicres.ValueEquation = strconv.FormatFloat(val.GetFloat64("historicresultvalueequation"), 'f', 6, 64)
			historicres.MaterialCost = strconv.FormatFloat(val.GetFloat64("historicresultmaterialcost"), 'f', 6, 64)
			simData.HistoricResult = historicres
			futureres.Revenue = strconv.FormatFloat(val.GetFloat64("futureresultrevenue"), 'f', 6, 64)
			futureres.LaborCost = strconv.FormatFloat(val.GetFloat64("futureresultlaborcost"), 'f', 6, 64)
			futureres.MaintenanceCost = strconv.FormatFloat(val.GetFloat64("futureresultmaintenacecost"), 'f', 6, 64)
			futureres.OperatingCost = strconv.FormatFloat(val.GetFloat64("futureresultoperationcost"), 'f', 6, 64)
			futureres.ServiceCost = strconv.FormatFloat(val.GetFloat64("futureresultservicecost"), 'f', 6, 64)
			futureres.ValueEquation = strconv.FormatFloat(val.GetFloat64("futureresultvalueequation"), 'f', 6, 64)
			futureres.MaterialCost = strconv.FormatFloat(val.GetFloat64("futureresultmaterialcost"), 'f', 6, 64)
			simData.FutureResult = futureres
			result = append(result, simData)
		}
	}

	for idx, res := range result {

		for _, val := range tmp {
			id := val.GetInt("id")

			if res.Id == id {
				plant := val.Get("selectedplant")
				unit := val.Get("selectedunit")
				scenario := val.Get("scenarioid")

				if plant != nil {
					plantName := val.GetString("selectedplant")
					tmpSelectedPlant := ""
					existPlant := false

					for _, sPlant := range result[idx].SelectedPlant {
						if plantName == sPlant {
							existPlant = true
						}
					}

					if !existPlant {
						tmpSelectedPlant = plantName
						result[idx].SelectedPlant = append(result[idx].SelectedPlant, tmpSelectedPlant)
					}
				}

				if unit != nil {
					unitName := val.GetString("selectedunit")
					tmpSelectedUnit := ""
					existUnit := false

					for _, sUnit := range result[idx].SelectedUnit {
						if unitName == sUnit {
							existUnit = true
						}
					}

					if !existUnit {
						tmpSelectedUnit = unitName
						result[idx].SelectedUnit = append(result[idx].SelectedUnit, tmpSelectedUnit)
					}
				}

				if scenario != nil {
					scenarioid := val.GetString("scenarioid")
					scenarioname := val.GetString("scenarioname")
					scenarioval := val.GetFloat64("scenariovalue")
					tmpSelectedScenario := SelScenario{}
					existscenario := false

					for _, sScenario := range result[idx].SelectedScenario {
						if scenarioid == sScenario.ID {
							existscenario = true
						}
					}

					if !existscenario {
						tmpSelectedScenario.ID = scenarioid
						tmpSelectedScenario.Name = scenarioname
						tmpSelectedScenario.Value = scenarioval
						result[idx].SelectedScenario = append(result[idx].SelectedScenario, tmpSelectedScenario)
					}
				}
			}
		}

	}

	return
}

func (s *ScenarioSimulation) RemoveData(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	status := ""
	msg := ""
	result := tk.M{}

	p := struct {
		SelectedSimulation  string
		SelectedDescription string
	}{}

	e := k.GetPayload(&p)
	if e != nil {
		s.WriteLog(e)
	}

	Name := p.SelectedSimulation
	Desc := p.SelectedDescription

	csr, err := s.DB().Connection.NewQuery().Command("procedure", tk.M{}.Set("name", "DeleteScenarioSimulation").Set("parms", tk.M{}.Set("@NAME", Name).Set("@DESC", Desc))).Cursor(nil)
	del := tk.M{}
	err = csr.Fetch(&del, 0, false)
	defer csr.Close()

	if err != nil {
		status = "OK"
		msg = "Delete Success"
	} else {
		msg = err.Error()
		status = "NOK"
	}

	result.Set("Status", status)
	result.Set("Message", msg)

	return result
}

func (s *ScenarioSimulation) SaveData(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	status := ""
	msg := ""
	result := tk.M{}

	p := struct {
		StartPeriod           string
		EndPeriod             string
		SelectedPlant         []string
		SelectedUnit          []string
		SimulationName        string
		SimulationDescription string
		SelectedScenario      []SeleScenario
		HistoricData          Simfloat
		FutureData            Simfloat
		Differential          Simfloat
	}{}

	e := k.GetPayload(&p)
	if e != nil {
		s.WriteLog(e)
	}

	start, _ := time.Parse(time.RFC3339, p.StartPeriod)
	end, _ := time.Parse(time.RFC3339, p.EndPeriod)
	name := p.SimulationName
	description := p.SimulationDescription
	historicrevenue := p.HistoricData.Revenue
	historiclabor := p.HistoricData.LaborCost
	historicmaterial := p.HistoricData.MaterialCost
	historicservice := p.HistoricData.ServiceCost
	historicoperation := p.HistoricData.OperatingCost
	historicmaintenance := p.HistoricData.MaintenanceCost
	historicvalueequation := p.HistoricData.ValueEquation
	futurerevenue := p.FutureData.Revenue
	futurelabor := p.FutureData.LaborCost
	futurematerial := p.FutureData.MaterialCost
	futureservice := p.FutureData.ServiceCost
	futureoperation := p.FutureData.OperatingCost
	futuremaintenance := p.FutureData.MaintenanceCost
	futurevalueequation := p.FutureData.ValueEquation
	differentialrevenue := p.Differential.Revenue
	differentiallabor := p.Differential.LaborCost
	differentialmaterial := p.Differential.MaterialCost
	differentialservice := p.Differential.ServiceCost
	differentialoperation := p.Differential.OperatingCost
	differentialmaintenance := p.Differential.MaintenanceCost
	differentialvalueequation := p.Differential.ValueEquation

	filter := tk.M{}
	params := tk.M{}
	filter.Set("name", "SaveScenarioSimulation")
	params.Set("@START_PERIOD", start)
	params.Set("@END_PERIOD", end)
	params.Set("@NAME", name)
	params.Set("@DESC", description)
	params.Set("@HISTORIC_REVENUE", historicrevenue)
	params.Set("@HISTORIC_LABOR", historiclabor)
	params.Set("@HISTORIC_MATERIAL", historicmaterial)
	params.Set("@HISTORIC_SERVICE", historicservice)
	params.Set("@HISTORIC_OPERATING", historicoperation)
	params.Set("@HISTORIC_MAINTENANCE", historicmaintenance)
	params.Set("@HISTORIC_VALUE_EQUATION", historicvalueequation)
	params.Set("@FUTURE_REVENUE", futurerevenue)
	params.Set("@FUTURE_LABOR", futurelabor)
	params.Set("@FUTURE_MATERIAL", futurematerial)
	params.Set("@FUTURE_SERVICE", futureservice)
	params.Set("@FUTURE_OPERATING", futureoperation)
	params.Set("@FUTURE_MAINTENANCE", futuremaintenance)
	params.Set("@FUTURE_VALUE_EQUATION", futurevalueequation)
	params.Set("@DIFFERENTIAL_REVENUE", differentialrevenue)
	params.Set("@DIFFERENTIAL_LABOR", differentiallabor)
	params.Set("@DIFFERENTIAL_MATERIAL", differentialmaterial)
	params.Set("@DIFFERENTIAL_SERVICE", differentialservice)
	params.Set("@DIFFERENTIAL_OPERATING", differentialoperation)
	params.Set("@DIFFERENTIAL_MAINTENANCE", differentialmaintenance)
	params.Set("@DIFFERENTIAL_VALUE_EQUATION", differentialvalueequation)
	filter.Set("parms", params)
	csr, err := s.DB().Connection.NewQuery().Command("procedure", filter).Cursor(nil)
	res := tk.M{}
	err = csr.Fetch(&res, 0, false)
	if err != nil {
		for _, plant := range p.SelectedPlant {
			csr, err = s.DB().Connection.NewQuery().Command("procedure", tk.M{}.Set("name", "SaveSelectedPlant").Set("parms", tk.M{}.Set("@Plant", plant))).Cursor(nil)
			plants := tk.M{}
			err = csr.Fetch(&plants, 0, true)
		}
		for _, unit := range p.SelectedUnit {
			csr, err = s.DB().Connection.NewQuery().Command("procedure", tk.M{}.Set("name", "SaveSelectedUnit").Set("parms", tk.M{}.Set("@Unit", unit))).Cursor(nil)
			units := tk.M{}
			err = csr.Fetch(&units, 0, true)
		}
		for _, scenario := range p.SelectedScenario {
			ID := scenario.ID
			Name := scenario.Name
			Value, e := strconv.ParseFloat(scenario.Value, 64)
			if e != nil {
				tk.Println(e.Error())
			}
			csr, err = s.DB().Connection.NewQuery().Command("procedure", tk.M{}.Set("name", "SaveSelectedScenario").Set("parms", tk.M{}.Set("@ID", ID).Set("@NAME", Name).Set("@VALUE", Value))).Cursor(nil)
			scenarios := tk.M{}
			err = csr.Fetch(&scenarios, 0, true)
		}
		status = "OK"
		msg = "Save Complete"
	} else {
		status = "NOK"
		msg = err.Error()
	}

	defer csr.Close()

	result.Set("Status", status)
	result.Set("Message", msg)

	return result
}
