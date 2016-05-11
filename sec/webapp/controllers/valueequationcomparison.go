package controllers

import (
	"github.com/eaciit/knot/knot.v1"
	. "github.com/eaciit/powerplant/sec/webapp/models"
	tk "github.com/eaciit/toolkit"
)

type ValueEquationComparisonController struct {
	*BaseController
}

func (c *ValueEquationComparisonController) Default(k *knot.WebContext) interface{} {
	c.LoadPartial(k)
	k.Config.OutputType = knot.OutputTemplate

	infos := PageInfo{}
	infos.PageId = "ValueEquationComparison"
	infos.PageTitle = "Value Equation Comparison"
	infos.Breadcrumbs = make(map[string]string, 0)
	tk.Println("testse")
	return infos
}

func (c *ValueEquationComparisonController) GetData(k *knot.WebContext) interface{} {
	vec := ValueEquationComparison{}
	result, e := vec.GetData(c.Ctx, k)
	return ResultInfo(result, e)
}
func (c *ValueEquationComparisonController) GetUnitList(k *knot.WebContext) interface{} {
	vec := ValueEquationComparison{}
	result, e := vec.GetUnitList(c.Ctx, k)
	return ResultInfo(result, e)
}
