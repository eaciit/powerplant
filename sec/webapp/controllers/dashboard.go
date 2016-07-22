package controllers

import (
	"time"

	"strings"

	"github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	. "github.com/eaciit/powerplant/sec/library/models"
	tk "github.com/eaciit/toolkit"
)

type DashboardController struct {
	*BaseController
}

func (c *DashboardController) CheckNotError(e error) error {
	if e != nil && (!strings.EqualFold(e.Error(), "no more data to fetched!")) {
		return e
	}
	return nil
}

func (c *DashboardController) Default(k *knot.WebContext) interface{} {
	if k.Session("userid") == nil {
		c.Redirect(k, "login", "default")
	}
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
		AssetClass []SampleAssetClass
		AssetLevel []SampleAssetLevel
		AssetType  []SampleAssetType
	}

	var (
		Result ReturnValue
	)

	filter := tk.M{}
	curr, e := c.DB().Find(&SampleAssetClass{}, filter)

	if e != nil {
	}

	e = curr.Fetch(&Result.AssetClass, 0, false)
	if e != nil {
		return e.Error()
	}

	curr, e = c.DB().Find(&SampleAssetLevel{}, filter)

	if e != nil {
	}

	e = curr.Fetch(&Result.AssetLevel, 0, false)
	if e != nil {
		return e.Error()
	}

	curr, e = c.DB().Find(&SampleAssetType{}, filter)

	if e != nil {
	}

	e = curr.Fetch(&Result.AssetType, 0, false)
	if e != nil {
		return e.Error()
	}

	defer curr.Close()

	selectedPeriod := time.Now().Year() - 1

	return (tk.M{}).Set("success", true).Set("data", Result).Set("message", "").Set("selected Period", selectedPeriod)
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
		PlantList         []PowerPlantCoordinates
		PlantCapacityList []tk.M
	}

	var (
		Result ReturnValue
	)

	filter := tk.M{}
	curr, e := c.DB().Find(&PowerPlantCoordinates{}, filter)

	if e != nil {
	}

	e = curr.Fetch(&Result.PlantList, 0, false)

	if e != nil {
		return e.Error()
	}

	defer curr.Close()

	//selectedPeriod := d.StartDate
	selectedPeriod := time.Now().Year() - 1

	cursor, _ := c.DB().Connection.NewQuery().
		Select("PlantCode as _id").
		From("ValueEquation_Dashboard").
		Where(dbox.Eq("Year", selectedPeriod)).
		Aggr(dbox.AggrSum, "InstalledCapacity", "TotalCapacity").
		Group("PlantCode").
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

		filter = append(filter, dbox.Eq("Year", selectedPeriod))
		filter = append(filter, dbox.Ne("UnitType", ""))

		if len(d.Plant) != 0 {
			filter = append(filter, dbox.Eq("Plant", d.Plant[0]))
		}

		result := make([]tk.M, 0)

		cursor, _ := c.DB().Connection.NewQuery().
			Select("UnitType as _id").
			From("ValueEquation_Dashboard").
			Where(filter...).
			Group("UnitType").
			Aggr(dbox.AggrSum, 1, "count").
			Order("count").
			Cursor(nil)

		defer cursor.Close()
		e = cursor.Fetch(&result, 0, true)

		e = c.CheckNotError(e)
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
			filter []*dbox.Filter
		)

		selectedPeriod := d.Period
		filter = append(filter, dbox.Eq("Year", selectedPeriod))

		if len(d.Plant) > 0 {
			filter = append(filter, dbox.Eq("Plant", d.Plant[0]))
		}

		result := make([]tk.M, 0)

		cursor, e := c.DB().Connection.NewQuery().
			Select("Plant as _id").
			From("ValueEquation_Dashboard").
			Where(filter...).
			Group("Plant").
			Aggr(dbox.AggrSum, "UpdatedFuelConsumption", "FuelConsumtion").
			Aggr(dbox.AggrSum, "NetGeneration", "Power").
			Order("_id").
			Cursor(nil)

		defer cursor.Close()

		e = cursor.Fetch(&result, 0, true)

		e = c.CheckNotError(e)

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

		filter := ""
		sintax := ""

		if len(d.Plant) > 0 {
			filter = d.Plant[0]
		}

		result := make([]tk.M, 0)

		if filter == "" {
			sintax = "select dbo.ValueEquation_Dashboard.Year as period, dbo.VEDTop10.WorkOrderType, count(*) as count, sum(dbo.VEDTop10.MaintenanceCost) as cost from dbo.ValueEquation_Dashboard inner join dbo.VEDTop10 on dbo.ValueEquation_Dashboard.Id = dbo.VEDTop10.VEId group by dbo.ValueEquation_Dashboard.Year, dbo.VEDTop10.WorkOrderType order by period asc, cost asc"
		} else {
			sintax = "select dbo.ValueEquation_Dashboard.Year as period, dbo.VEDTop10.WorkOrderType, count(*) as count, sum(dbo.VEDTop10.MaintenanceCost) as cost from dbo.ValueEquation_Dashboard inner join dbo.VEDTop10 on dbo.ValueEquation_Dashboard.Id = dbo.VEDTop10.VEId where dbo.ValueEquation_Dashboard.Plant = '" + filter + "' group by dbo.ValueEquation_Dashboard.Year, dbo.VEDTop10.WorkOrderType order by period asc, cost asc"
		}

		cursor, e := c.DB().Connection.NewQuery().
			Command("freequery", tk.M{}.
				Set("syntax", sintax)).
			Cursor(nil)

		_ = filter
		defer cursor.Close()
		e = cursor.Fetch(&result, 0, true)

		e = c.CheckNotError(e)

		return result, e
	}, nil)

	return ResultInfo(r, e)
}
