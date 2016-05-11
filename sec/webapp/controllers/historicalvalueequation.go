package controllers

import (
	"github.com/eaciit/knot/knot.v1"
	. "github.com/eaciit/powerplant/sec/webapp/models"
	// tk "github.com/eaciit/toolkit"
	// "strconv"
	// "gopkg.in/mgo.v2/bson"
	// "time"
)

type HistoricalValueEquationController struct {
	*BaseController
}

func (c *HistoricalValueEquationController) GetSummaryData(k *knot.WebContext) interface{} {
	hve := HistoricalValueEquation{}
	result, e := hve.GetSummaryData(c.Ctx, k)
	return ResultInfo(result, e)
}

func (c *HistoricalValueEquationController) GetMaintenance(k *knot.WebContext) interface{} {
	hve := HistoricalValueEquation{}
	result, e := hve.GetMaintenanceData(c.Ctx, k)
	return ResultInfo(result, e)
}

func (c *HistoricalValueEquationController) GetOperating(k *knot.WebContext) interface{} {
	hve := HistoricalValueEquation{}
	result, e := hve.GetOperatingData(c.Ctx, k)
	return ResultInfo(result, e)
}
func (c *HistoricalValueEquationController) GetRevenue(k *knot.WebContext) interface{} {
	hve := HistoricalValueEquation{}
	result, e := hve.GetRevenueData(c.Ctx, k)
	return ResultInfo(result, e)
}
func (c *HistoricalValueEquationController) GetDataQuality(k *knot.WebContext) interface{} {
	hve := HistoricalValueEquation{}
	result, e := hve.GetDataQuality(c.Ctx, k)
	return ResultInfo(result, e)
}

func (c *HistoricalValueEquationController) GetPerformance(k *knot.WebContext) interface{} {
	hve := HistoricalValueEquation{}
	result, e := hve.GetPerformanceData(c.Ctx, k)
	return ResultInfo(result, e)
}
func (c *HistoricalValueEquationController) GetAssetWork(k *knot.WebContext) interface{} {
	hve := HistoricalValueEquation{}
	result, e := hve.GetAssetWorkData(c.Ctx, k)
	return ResultInfo(result, e)
}
