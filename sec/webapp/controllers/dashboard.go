package controllers

import (
	"time"

	"github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	. "github.com/eaciit/powerplant/sec/webapp/models"
	tk "github.com/eaciit/toolkit"
)

type DashboardController struct {
	*BaseController
}

func (c *DashboardController) Default(k *knot.WebContext) interface{} {
	//c.LoadBase(k)
	c.LoadPartial(k, "dashboard/numberofturbines.html", "dashboard/powervsfuelconsumtion.html", "dashboard/numberofworkorders.html", "dashboard/maintenancecost.html")

	k.Config.OutputType = knot.OutputTemplate

	infos := PageInfo{}
	infos.PageId = "Dashboard"
	infos.PageTitle = "Dashboard"
	infos.Breadcrumbs = make(map[string]string, 0)

	return infos
}

func (c *DashboardController) Initiate(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	type ReturnValue struct {
		AssetClass []AssetClass
		AssetLevel []AssetLevel
		AssetType  []AssetType
	}

	var (
		Result ReturnValue
	)

	filter := tk.M{}
	curr, e := c.DB().Find(&AssetClass{}, filter)

	if e != nil {
	}

	e = curr.Fetch(&Result.AssetClass, 0, false)
	if e != nil {
		return e.Error()
	}

	curr, e = c.DB().Find(&AssetLevel{}, filter)

	if e != nil {
	}

	e = curr.Fetch(&Result.AssetLevel, 0, false)
	if e != nil {
		return e.Error()
	}

	curr, e = c.DB().Find(&AssetType{}, filter)

	if e != nil {
	}

	e = curr.Fetch(&Result.AssetType, 0, false)
	if e != nil {
		return e.Error()
	}

	defer curr.Close()

	return (tk.M{}).Set("success", true).Set("data", Result).Set("message", "")
}

func (c *DashboardController) GetData(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	d := struct {
		StartDate string
		EndDate   string
		Plant     []string
	}{}

	e := k.GetPayload(&d)

	type ReturnValue struct {
		PlantList         []PlantData
		PlantCapacityList []tk.M
	}

	var (
		Result       ReturnValue
		query, pipes []tk.M
	)

	filter := tk.M{}
	curr, e := c.DB().Find(&PlantData{}, filter)

	if e != nil {
	}

	e = curr.Fetch(&Result.PlantList, 0, false)

	if e != nil {
		return e.Error()
	}

	defer curr.Close()

	//selectedPeriod := d.StartDate
	selectedPeriod := time.Now().Year() - 1

	query = append(query, tk.M{"Period.Year": tk.M{"$eq": selectedPeriod}})
	pipes = append(pipes, tk.M{"$match": tk.M{"$and": query}})
	pipes = append(pipes, tk.M{"$group": tk.M{
		"_id":           "$PlantDetail.PlantCode",
		"TotalCapacity": tk.M{"$sum": "$TurbineInfos.InstalledCapacity"},
	}})

	cursor, _ := c.DB().Connection.NewQuery().
		From("ValueEquation_Dashboard").
		Command("pipe", pipes).
		Cursor(nil)

	defer cursor.Close()
	e = cursor.Fetch(&Result.PlantCapacityList, 0, false)

	result := tk.M{}
	result.Set("success", true)
	result.Set("Data", Result)
	return ResultInfo(result, e)
}

func (c *DashboardController) GetNumberOfTurbines(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	var e error
	r := new(tk.Result)

	d := struct {
		StartDate string
		EndDate   string
		Plant     []string
	}{}

	e = k.GetPayload(&d)

	r.Run(func(in interface{}) (interface{}, error) {
		var filter []*dbox.Filter

		selectedPeriod := time.Now().Year() - 1

		filter = append(filter, dbox.Eq("Period.Year", selectedPeriod))
		filter = append(filter, dbox.Ne("TurbineInfos.UnitType", ""))
		filter = append(filter, dbox.Ne("TurbineInfos.UnitType", nil))

		if len(d.Plant) != 0 {
			filter = append(filter, dbox.Eq("Plant", d.Plant[0]))
		}

		result := make([]tk.M, 0)

		cursor, _ := c.DB().Connection.NewQuery().
			From("ValueEquation_Dashboard").
			Where(filter...).
			Group("TurbineInfos.UnitType").
			Aggr(dbox.AggrSum, 1, "count").
			Cursor(nil)

		defer cursor.Close()
		e = cursor.Fetch(&result, 0, true)

		return result, e
	}, nil)

	return ResultInfo(r, e)
}

func (c *DashboardController) GetPowerVsFuelConsumtion(k *knot.WebContext) interface{} {
	//c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson

	d := struct {
		StartDate string
		EndDate   string
		Period    int
		Plant     []string
	}{}

	e := k.GetPayload(&d)

	r := new(tk.Result)
	r.Run(func(in interface{}) (interface{}, error) {

		var (
			pipes  []tk.M
			filter []*dbox.Filter
		)

		selectedPeriod := d.Period
		filter = append(filter, dbox.Eq("Period.Year", selectedPeriod))

		if len(d.Plant) > 0 {
			filter = append(filter, dbox.In("Plant", d.Plant[0]))
		}

		result := make([]tk.M, 0)

		pipes = append(pipes, tk.M{"$group": tk.M{
			"_id":            "$Plant",
			"FuelConsumtion": tk.M{"$sum": "$TurbineInfos.UpdatedFuelConsumption"},
			"Power":          tk.M{"$sum": "$NetGeneration"},
		}})

		pipes = append(pipes, tk.M{"$sort": tk.M{
			"_id": 1,
		}})

		cursor, e := c.DB().Connection.NewQuery().
			Command("pipe", pipes).
			Where(filter...).
			From("ValueEquation_Dashboard").
			Cursor(nil)

		defer cursor.Close()

		e = cursor.Fetch(&result, 0, true)
		return result, e
	}, nil)

	return ResultInfo(r, e)
}

func (c *DashboardController) GetNumberOfWorkOrder(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	var e error
	r := new(tk.Result)

	d := struct {
		EndDate   string
		Plant     []string
		StartDate string
	}{}

	e = k.GetPayload(&d)

	r.Run(func(in interface{}) (interface{}, error) {

		var (
			pipes  []tk.M
			filter []*dbox.Filter
		)

		if len(d.Plant) > 0 {
			filter = append(filter, dbox.Eq("Plant", d.Plant[0]))
		}

		pipes = append(pipes, tk.M{"$unwind": "$Top10"})
		pipes = append(pipes, tk.M{"$group": tk.M{
			"_id":   tk.M{"period": "$Period.Year", "workorder": "$Top10.WorkOrderType"},
			"count": tk.M{"$sum": 1},
			"cost":  tk.M{"$sum": "$Top10.MaintenanceCost"},
		}})

		result := make([]tk.M, 0)

		cursor, e := c.DB().Connection.NewQuery().
			From("ValueEquation_Dashboard").
			Command("pipe", pipes).
			Where(filter...).
			Order("Plant").
			Cursor(nil)

		defer cursor.Close()
		e = cursor.Fetch(&result, 0, true)

		return result, e
	}, nil)

	return ResultInfo(r, e)
}
