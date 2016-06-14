package controllers

/*import (
	"github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	. "github.com/eaciit/powerplant/sec/webapp/models"
	tk "github.com/eaciit/toolkit"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"time"
)

/*type ScenarioSimulation struct {
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

type SelScenario struct {
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

func (s *ScenarioSimulation) GetData(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	status := ""
	msg := ""
	res := []tk.M{}

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
	query = append(query, dbox.Gte("Period.Dates", start))
	query = append(query, dbox.Lte("Period.Dates", end))

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

	tk.Println(res)

	defer csr.Close()

	result := tk.M{}
	result.Set("Status", status)
	result.Set("Message", msg)
	result.Set("Data", res)

	return result
}

func (s *ScenarioSimulation) GetDataSimulation(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	data := []tk.M{}
	result := tk.M{}
	status := ""
	msg := ""
	csr, err := s.DB().Connection.NewQuery().From("ScenarioSimulation").Select().Cursor(nil)

	if err != nil {
		msg = err.Error()
		status = "NOK"
	} else {
		status = "OK"
	}

	err = csr.Fetch(&data, 0, true)

	defer csr.Close()

	result.Set("Status", status)
	result.Set("Message", msg)
	result.Set("Data", data)

	return result
}

func (s *ScenarioSimulation) RemoveData(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	status := ""
	msg := ""
	query := []*dbox.Filter{}
	result := tk.M{}

	p := struct {
		SelectedSimulation  string
		SelectedDescription string
	}{}

	e := k.GetPayload(&p)
	if e != nil {
		s.WriteLog(e)
	}

	query = append(query, dbox.Eq("Name", p.SelectedSimulation))
	query = append(query, dbox.Eq("Description", p.SelectedDescription))

	err := s.DB().Connection.NewQuery().Delete().From("ScenarioSimulation").Where(query...).Exec(nil)
	if err != nil {
		msg = err.Error()
		status = "NOK"
	} else {
		status = "OK"
		msg = "Delete Success"
	}

	result.Set("Status", status)
	result.Set("Message", msg)

	return result
}

func (s *ScenarioSimulation) SaveData(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	status := ""
	msg := ""
	query := []*dbox.Filter{}
	result := tk.M{}

	p := struct {
		StartPeriod            string
		EndPeriod              string
		Name                   string
		Description            string
		SelectedPlant          []string
		SelectedUnit           []string
		SelectedSimulation     string
		SelectedDescription    string
		SelectedScenario       []SelScenario
		HistoricData           Sim
		FutureData             Sim
		Differential           Sim
		SelectedScenarioLength string
	}{}

	e := k.GetPayload(&p)
	if e != nil {
		s.WriteLog(e)
	}

	start, _ := time.Parse(time.RFC3339, p.StartPeriod)
	end, _ := time.Parse(time.RFC3339, p.EndPeriod)

	if p.SelectedSimulation != "" {
		query = append(query, dbox.Eq("Name", p.SelectedSimulation))
		query = append(query, dbox.Eq("Description", p.SelectedDescription))
		err := s.DB().Connection.NewQuery().Delete().From("ScenarioSimulation").Where(query...).Exec(nil)
		if err != nil {
			tk.Println(err.Error())
		}
	}

	data := NewScenarioSimulation()
	data.ID = bson.NewObjectId()
	data.Start_Period = start
	data.End_Period = end
	data.Name = p.SelectedSimulation
	data.Description = p.SelectedDescription
	data.SelectedPlant = p.SelectedPlant
	data.SelectedUnit = p.SelectedUnit
	SelectedScenarioLength, _ := strconv.Atoi(p.SelectedScenarioLength)
	data.SelectedScenario = []SelectedScenario{}
	for i := 0; i < SelectedScenarioLength; i++ {
		scenario := SelectedScenario{}
		scenario.ID = p.SelectedScenario[i].ID
		scenario.Name = p.SelectedScenario[i].Name
		scenario.Value, _ = strconv.ParseFloat(p.SelectedScenario[i].Value, 64)
		data.SelectedScenario = append(data.SelectedScenario, scenario)
	}

	data.HistoricResult.Revenue, _ = strconv.ParseFloat(p.HistoricData.Revenue, 64)
	data.HistoricResult.LaborCost, _ = strconv.ParseFloat(p.HistoricData.LaborCost, 64)
	data.HistoricResult.MaterialCost, _ = strconv.ParseFloat(p.HistoricData.MaterialCost, 64)
	data.HistoricResult.ServiceCost, _ = strconv.ParseFloat(p.HistoricData.ServiceCost, 64)
	data.HistoricResult.OperatingCost, _ = strconv.ParseFloat(p.HistoricData.OperatingCost, 64)
	data.HistoricResult.MaintenanceCost, _ = strconv.ParseFloat(p.HistoricData.MaintenanceCost, 64)
	data.HistoricResult.ValueEquation, _ = strconv.ParseFloat(p.HistoricData.ValueEquation, 64)

	data.FutureResult.Revenue, _ = strconv.ParseFloat(p.FutureData.Revenue, 64)
	data.FutureResult.LaborCost, _ = strconv.ParseFloat(p.FutureData.LaborCost, 64)
	data.FutureResult.MaterialCost, _ = strconv.ParseFloat(p.FutureData.MaterialCost, 64)
	data.FutureResult.ServiceCost, _ = strconv.ParseFloat(p.FutureData.ServiceCost, 64)
	data.FutureResult.OperatingCost, _ = strconv.ParseFloat(p.FutureData.OperatingCost, 64)
	data.FutureResult.MaintenanceCost, _ = strconv.ParseFloat(p.FutureData.MaintenanceCost, 64)
	data.FutureResult.ValueEquation, _ = strconv.ParseFloat(p.FutureData.ValueEquation, 64)

	data.Differential.Revenue, _ = strconv.ParseFloat(p.Differential.Revenue, 64)
	data.Differential.LaborCost, _ = strconv.ParseFloat(p.Differential.LaborCost, 64)
	data.Differential.MaterialCost, _ = strconv.ParseFloat(p.Differential.MaterialCost, 64)
	data.Differential.ServiceCost, _ = strconv.ParseFloat(p.Differential.ServiceCost, 64)
	data.Differential.OperatingCost, _ = strconv.ParseFloat(p.Differential.OperatingCost, 64)
	data.Differential.MaintenanceCost, _ = strconv.ParseFloat(p.Differential.MaintenanceCost, 64)
	data.Differential.ValueEquation, _ = strconv.ParseFloat(p.Differential.ValueEquation, 64)

	data.Last_Update = time.Now()

	err := s.DB().Save(data)
	if err != nil {
		msg = err.Error()
		status = "NOK"
	} else {
		status = "OK"
		msg = "Save Complete"
	}

	result.Set("Status", status)
	result.Set("Message", msg)

	return result
}*/
