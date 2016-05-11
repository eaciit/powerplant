package models

import (
	"github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	"github.com/eaciit/orm"
	tk "github.com/eaciit/toolkit"
	"time"
)

type HistoricalValueEquation struct {
	StartPeriod   time.Time
	EndPeriod     time.Time
	Scope         string
	Selected      []string
	SelectedPlant string
	OrderType     []string
	DF            string
	WOTOP         int
}

func (m *HistoricalValueEquation) SetPayLoad(k *knot.WebContext) error {
	d := struct {
		StartPeriod   string
		EndPeriod     string
		Scope         string
		Selected      []string
		SelectedPlant string
		OrderType     []string
		DF            string
		WOTOP         int
	}{}
	e := k.GetPayload(&d)
	m.StartPeriod, _ = time.Parse(time.RFC3339, d.StartPeriod)
	m.EndPeriod, _ = time.Parse(time.RFC3339, d.EndPeriod)
	m.Scope = d.Scope
	m.Selected = d.Selected

	m.SelectedPlant = d.SelectedPlant
	m.OrderType = d.OrderType
	m.DF = d.DF
	m.WOTOP = d.WOTOP

	return e
}

func (m *HistoricalValueEquation) GetSummaryData(ctx *orm.DataContext, k *knot.WebContext) (interface{}, error) {
	m.SetPayLoad(k)
	result, DataChart, DataDetail := tk.M{}, []tk.M{}, []tk.M{}
	c := ctx.Connection
	ve := ValueEquation{}
	query := []*dbox.Filter{}
	query = append(query, dbox.Gte("Period.Dates", m.StartPeriod))
	query = append(query, dbox.Lte("Period.Dates", m.EndPeriod))
	csr, e := c.NewQuery().
		Where(query...).
		Aggr(dbox.AggrSum, "$Revenue", "Amount").
		Aggr(dbox.AggrSum, "$MaintenanceCost", "MaintenanceCost").
		Aggr(dbox.AggrSum, "$OperatingCost", "OperatingCost").
		From(ve.TableName()).Group("Plant").Cursor(nil)
	e = csr.Fetch(&DataChart, 0, false)
	csr.Close()
	if e != nil {
		return nil, e
	}
	for _, i := range DataChart {
		i.Set("_id", i.Get("_id").(tk.M).Get("Plant"))
	}
	result.Set("DataChart", DataChart)

	groupBy := "Plant"
	if m.Scope == "Plant" {
		groupBy = "Unit"
		if m.Selected != nil && len(m.Selected) > 0 {
			query = append(query, dbox.In("Plant", m.Selected))
		}
	}

	if m.Scope != "Unit" {
		csr, e := c.NewQuery().
			Where(query...).
			Aggr(dbox.AggrSum, "$Revenue", "Revenue").
			Aggr(dbox.AggrSum, "$MaintenanceCost", "MaintenanceCost").
			Aggr(dbox.AggrSum, "$OperatingCost", "OperatingCost").
			From(ve.TableName()).Group(groupBy).Cursor(nil)
		e = csr.Fetch(&DataDetail, 0, false)
		csr.Close()
		if e != nil {
			return nil, e
		}
		for _, i := range DataDetail {
			i.Set("_id", i.Get("_id").(tk.M).Get(groupBy))
		}
	}

	result.Set("DataDetail", DataDetail)
	return result, nil
}

func (m *HistoricalValueEquation) GetMaintenanceData(ctx *orm.DataContext, k *knot.WebContext) (interface{}, error) {
	m.SetPayLoad(k)
	result, DataMainEx, DataOrder, DataChart, DataTable, pipes := tk.M{}, []tk.M{}, []tk.M{}, []tk.M{}, []tk.M{}, []tk.M{}
	c := ctx.Connection
	ve := ValueEquation{}
	query := []*dbox.Filter{}
	query = append(query, dbox.Gte("Period.Dates", m.StartPeriod))
	query = append(query, dbox.Lte("Period.Dates", m.EndPeriod))
	groupBy := "$Plant"
	switch m.Scope {
	case "Kingdom":
		break
	case "Plant":
		if m.Selected != nil && len(m.Selected) > 0 {
			query = append(query, dbox.In("Plant", m.Selected))
		}
		groupBy = "$Unit"
		break
	case "Phase":
		if m.Selected != nil && len(m.Selected) > 0 {
			query = append(query, dbox.In("Phase", m.Selected))
		}
		break
	case "Unit":
		query = append(query, dbox.Eq("Plant", m.SelectedPlant))
		if m.Selected != nil && len(m.Selected) > 0 {
			query = append(query, dbox.In("Unit", m.Selected))
		}
		groupBy = "$Unit"
		break
	default:
		break
	}

	// Get DataTable
	pipes = append(pipes, tk.M{"$unwind": "$Detail"})
	pipes = append(pipes, tk.M{"$group": tk.M{
		"_id": tk.M{"DataSource": "$Detail.DataSource",
			"WorkOrderType": "$Detail.WorkOrderType"},
		"LaborCost":    tk.M{"$sum": "$Detail.LaborCost"},
		"MaterialCost": tk.M{"$sum": "$Detail.MaterialCost"},
		"ServiceCost":  tk.M{"$sum": "$Detail.ServiceCost"},
	}})
	csr, e := c.NewQuery().Command("pipe", pipes).Where(query...).From(ve.TableName()).Cursor(nil)
	e = csr.Fetch(&DataTable, 0, false)
	csr.Close()
	if e != nil {
		return nil, e
	}

	// Get DataMainEx
	pipes = append(pipes[0:0], tk.M{"$group": tk.M{
		"_id":          groupBy,
		"LaborCost":    tk.M{"$sum": "$TotalLabourCost"},
		"MaterialCost": tk.M{"$sum": "$TotalMaterialCost"},
		"ServiceCost":  tk.M{"$sum": "$TotalServicesCost"},
	}})
	csr, e = c.NewQuery().Command("pipe", pipes).Where(query...).From(ve.TableName()).Cursor(nil)
	e = csr.Fetch(&DataMainEx, 0, false)
	csr.Close()
	if e != nil {
		return nil, e
	}

	// Get DataOrder
	pipes = append(pipes[0:0], tk.M{"$unwind": "$Top10"})
	pipes = append(pipes, tk.M{"$group": tk.M{
		"_id":           "$Top10.WorkOrderType",
		"WorkOrderType": tk.M{"$first": "$Top10.WorkOrderType"},
	}})
	csr, e = c.NewQuery().Command("pipe", pipes).Where(query...).From(ve.TableName()).Cursor(nil)
	e = csr.Fetch(&DataOrder, 0, false)
	csr.Close()
	if e != nil {
		return nil, e
	}

	// Data Chart - set for empty, if you want to add it. its available on the previous version while its still using beego.

	result.Set("DataMainEx", DataMainEx).Set("DataOrder", DataOrder).Set("DataChart", DataChart).Set("DataTable", DataTable)
	return result, nil
}

func (m *HistoricalValueEquation) GetOperatingData(ctx *orm.DataContext, k *knot.WebContext) (interface{}, error) {
	m.SetPayLoad(k)
	result, DataChart, DataTable, DataTotal, DataDetail := tk.M{}, []tk.M{}, []tk.M{}, []tk.M{}, []tk.M{}
	c := ctx.Connection
	ve := ValueEquation{}
	query := []*dbox.Filter{}
	query = append(query, dbox.Gte("Period.Dates", m.StartPeriod))
	query = append(query, dbox.Lte("Period.Dates", m.EndPeriod))
	groupBy := "Plant"
	switch m.Scope {
	case "Kingdom":
		break
	case "Plant":
		if m.Selected != nil && len(m.Selected) > 0 {
			query = append(query, dbox.In("Plant", m.Selected))
		}
		groupBy = "Unit"
		break
	case "Phase":
		if m.Selected != nil && len(m.Selected) > 0 {
			query = append(query, dbox.In("Phase", m.Selected))
		}
		break
	case "Unit":
		query = append(query, dbox.Eq("Plant", m.SelectedPlant))
		if m.Selected != nil && len(m.Selected) > 0 {
			query = append(query, dbox.In("Unit", m.Selected))
		}
		groupBy = "Unit"
		break
	default:
		break
	}

	// Getting "DataTable" - skip, need $unwind [ $Fuel ]
	csr, e := c.NewQuery().
		Where(query...).
		Aggr(dbox.AggrSum, "$OperatingCost", "OperatingCost").
		Aggr(dbox.AggrSum, "$FuelTransportCost", "FuelTransportCost").
		From(ve.TableName()).Group("").Cursor(nil)
	e = csr.Fetch(&DataTotal, 0, false)
	csr.Close()
	if e != nil {
		return nil, e
	}

	if m.Scope != "Unit" {
		csr, e = c.NewQuery().
			Where(query...).
			Aggr(dbox.AggrSum, "$Capacity", "Capacity").
			Aggr(dbox.AggrSum, "$NetGeneration", "NetGeneration").
			Aggr(dbox.AggrSum, "$OperatingCost", "OperatingCost").
			From(ve.TableName()).Group(groupBy).Cursor(nil)
		e = csr.Fetch(&DataDetail, 0, false)
		csr.Close()
		if e != nil {
			return nil, e
		}
		for _, i := range DataDetail {
			i.Set("_id", i.Get("_id").(tk.M).Get(groupBy))
		}
	}

	result.Set("DataChart", DataChart).Set("DataTable", DataTable).Set("DataTotal", DataTotal).Set("DataDetail", DataDetail)
	return result, e
}

func (m *HistoricalValueEquation) GetRevenueData(ctx *orm.DataContext, k *knot.WebContext) (interface{}, error) {
	m.SetPayLoad(k)
	result, DataChartRevenue, DataChartRevenueEx := tk.M{}, []tk.M{}, []tk.M{}
	c := ctx.Connection
	ve := ValueEquation{}
	query := []*dbox.Filter{}
	query = append(query, dbox.Gte("Period.Dates", m.StartPeriod))
	query = append(query, dbox.Lte("Period.Dates", m.EndPeriod))
	groupBy := "Plant"
	switch m.Scope {
	case "Kingdom":
		break
	case "Plant":
		if m.Selected != nil && len(m.Selected) > 0 {
			query = append(query, dbox.In("Plant", m.Selected))
		}
		groupBy = "Unit"
		break
	case "Phase":
		if m.Selected != nil && len(m.Selected) > 0 {
			query = append(query, dbox.In("Phase", m.Selected))
		}
		break
	case "Unit":
		query = append(query, dbox.Eq("Plant", m.SelectedPlant))
		if m.Selected != nil && len(m.Selected) > 0 {
			query = append(query, dbox.In("Unit", m.Selected))
		}
		groupBy = "Unit"
		break
	default:
		break
	}

	csr, e := c.NewQuery().
		Where(query...).
		Aggr(dbox.AggrSum, "$CapacityPayment", "CapacityPayment").
		Aggr(dbox.AggrSum, "$EnergyPayment", "EnergyPayment").
		Aggr(dbox.AggrSum, "$StartupPayment", "StartupPayment").
		Aggr(dbox.AggrSum, "$PenaltyAmount", "PenaltyAmount").
		Aggr(dbox.AggrSum, "$Incentive", "Incentive").
		Aggr(dbox.AggrSum, "$Revenue", "Revenue").
		From(ve.TableName()).Group("").Cursor(nil)
	e = csr.Fetch(&DataChartRevenue, 0, false)
	csr.Close()
	if e != nil {
		return nil, e
	}

	csr, e = c.NewQuery().
		Where(query...).
		Aggr(dbox.AggrSum, "$CapacityPayment", "CapacityPayment").
		Aggr(dbox.AggrSum, "$EnergyPayment", "EnergyPayment").
		Aggr(dbox.AggrSum, "$StartupPayment", "StartupPayment").
		Aggr(dbox.AggrSum, "$PenaltyAmount", "PenaltyAmount").
		Aggr(dbox.AggrSum, "$Incentive", "Incentive").
		Aggr(dbox.AggrSum, "$Revenue", "Revenue").
		From(ve.TableName()).Group(groupBy).Cursor(nil)
	e = csr.Fetch(&DataChartRevenueEx, 0, false)
	csr.Close()
	if e != nil {
		return nil, e
	}
	for _, i := range DataChartRevenueEx {
		i.Set("_id", i.Get("_id").(tk.M).Get(groupBy))
	}
	// Remaining : sort by _id. descending for "DataChartRevenueEx"
	return result.Set("DataChartRevenue", DataChartRevenue).Set("DataChartRevenueEx", DataChartRevenueEx), e
}

func (m *HistoricalValueEquation) GetDataQuality(ctx *orm.DataContext, k *knot.WebContext) (interface{}, error) {
	m.SetPayLoad(k)
	var e error = nil
	result := []tk.M{}
	c := ctx.Connection
	vedq := ValueEquationDataQuality{}
	query := []*dbox.Filter{}
	query = append(query, dbox.Gte("Period.Dates", m.StartPeriod))
	query = append(query, dbox.Lte("Period.Dates", m.EndPeriod))
	groupBy := "Plant"
	switch m.Scope {
	case "Kingdom":
		break
	case "Plant":
		if m.Selected != nil && len(m.Selected) > 0 {
			query = append(query, dbox.In("Plant", m.Selected))
		}
		groupBy = "Unit"
		break
	case "Phase":
		if m.Selected != nil && len(m.Selected) > 0 {
			query = append(query, dbox.In("Phase", m.Selected))
		}
		break
	case "Unit":
		query = append(query, dbox.Eq("Plant", m.SelectedPlant))
		if m.Selected != nil && len(m.Selected) > 0 {
			query = append(query, dbox.In("Unit", m.Selected))
		}
		groupBy = "Unit"
		break
	default:
		break
	}

	if m.Scope == "Unit" || (m.Scope == "Plant" && m.Selected != nil && len(m.Selected) == 1) {
		query := []*dbox.Filter{}
		query = append(query, dbox.Gte("Period.Dates", m.StartPeriod))
		query = append(query, dbox.Lte("Period.Dates", m.EndPeriod))
		if m.Scope == "Unit" {
			query = append(query, dbox.Eq("Plant", m.SelectedPlant))
			if m.Selected != nil && len(m.Selected) > 0 {
				query = append(query, dbox.In("Unit", m.Selected))
			}
		} else {
			query = append(query, dbox.Eq("Plant", m.Selected[0]))
		}

		csr, e := ctx.Find(new(ValueEquationDataQuality), tk.M{}.Set("where", dbox.And(query...)))
		e = csr.Fetch(&result, 0, false)
		csr.Close()
		if e != nil {
			return nil, e
		}
	} else {
		csr, e := c.NewQuery().
			Where(query...).
			Aggr(dbox.AggrSum, "1", "Count").
			Aggr(dbox.AggrSum, "$CapacityPayment_Data", "CapacityPayment_Data").
			Aggr(dbox.AggrSum, "$EnergyPayment_Data", "EnergyPayment_Data").
			Aggr(dbox.AggrSum, "$StartupPayment_Data", "StartupPayment_Data").
			Aggr(dbox.AggrSum, "$Penalty_Data", "Penalty_Data").
			Aggr(dbox.AggrSum, "$Incentive_Data", "Incentive_Data").
			Aggr(dbox.AggrSum, "$MaintenanceCost_Data", "MaintenanceCost_Data").
			Aggr(dbox.AggrSum, "$MaintenanceDuration_Data", "MaintenanceDuration_Data").
			Aggr(dbox.AggrSum, "$PrimaryFuel1st_Data", "PrimaryFuel1st_Data").
			Aggr(dbox.AggrSum, "$PrimaryFuel2nd_Data", "PrimaryFuel2nd_Data").
			Aggr(dbox.AggrSum, "$BackupFuel_Data", "BackupFuel_Data").
			Aggr(dbox.AggrSum, "$FuelTransport_Data", "FuelTransport_Data").
			From(vedq.TableName()).Group(groupBy).Cursor(nil)
		e = csr.Fetch(&result, 0, false)
		csr.Close()
		if e != nil {
			return nil, e
		}
		for _, i := range result {
			i.Set("_id", i.Get("_id").(tk.M).Get(groupBy))
		}
	}
	return result, e
}

func (m *HistoricalValueEquation) GetPerformanceData(ctx *orm.DataContext, k *knot.WebContext) (interface{}, error) {
	m.SetPayLoad(k)
	var e error = nil
	result := []tk.M{}
	c := ctx.Connection
	ve := ValueEquation{}
	query := []*dbox.Filter{}
	query = append(query, dbox.Gte("Period.Dates", m.StartPeriod))
	query = append(query, dbox.Lte("Period.Dates", m.EndPeriod))
	groupBy := "Plant"
	switch m.Scope {
	case "Kingdom":
		break
	case "Plant":
		if m.Selected != nil && len(m.Selected) > 0 {
			query = append(query, dbox.In("Plant", m.Selected))
		}
		groupBy = "Unit"
		break
	case "Phase":
		if m.Selected != nil && len(m.Selected) > 0 {
			query = append(query, dbox.In("Phase", m.Selected))
		}
		break
	case "Unit":
		query = append(query, dbox.Eq("Plant", m.SelectedPlant))
		if m.Selected != nil && len(m.Selected) > 0 {
			query = append(query, dbox.In("Unit", m.Selected))
		}
		groupBy = "Unit"
		break
	default:
		break
	}

	csr, e := c.NewQuery().
		Where(query...).
		Aggr(dbox.AggrSum, "$NetGeneration", "NetGeneration").
		Aggr(dbox.AggrSum, "$PrctWAF", "PrctWAF").
		From(ve.TableName()).Group(groupBy).Cursor(nil)
	e = csr.Fetch(&result, 0, false)
	csr.Close()
	if e != nil {
		return nil, e
	}
	for _, i := range result {
		i.Set("_id", i.Get("_id").(tk.M).Get(groupBy))
	}
	// Remaining : sort by _id. descending for "result"
	return result, e
}

func (m *HistoricalValueEquation) GetAssetWorkData(ctx *orm.DataContext, k *knot.WebContext) (interface{}, error) {
	m.SetPayLoad(k)
	var e error = nil
	result := []tk.M{}
	c := ctx.Connection
	ve := ValueEquation{}
	query := []*dbox.Filter{}
	query = append(query, dbox.Gte("Period.Dates", m.StartPeriod))
	query = append(query, dbox.Lte("Period.Dates", m.EndPeriod))
	groupBy := "Plant"
	switch m.Scope {
	case "Kingdom":
		break
	case "Plant":
		if m.Selected != nil && len(m.Selected) > 0 {
			query = append(query, dbox.In("Plant", m.Selected))
		}
		groupBy = "Unit"
		break
	case "Phase":
		if m.Selected != nil && len(m.Selected) > 0 {
			query = append(query, dbox.In("Phase", m.Selected))
		}
		break
	case "Unit":
		query = append(query, dbox.Eq("Plant", m.SelectedPlant))
		if m.Selected != nil && len(m.Selected) > 0 {
			query = append(query, dbox.In("Unit", m.Selected))
		}
		groupBy = "Unit"
		break
	default:
		break
	}

	csr, e := c.NewQuery().
		Where(query...).
		Aggr(dbox.AggrSum, "$ValueEquationCost", "ValueEquationCost").
		Aggr(dbox.AggrSum, "$MaxPowerGeneration", "MaxPowerGeneration").
		Aggr(dbox.AggrSum, "$PotentialRevenue", "PotentialRevenue").
		Aggr(dbox.AggrSum, "$NetGeneration", "NetGeneration").
		Aggr(dbox.AggrSum, "$Revenue", "Revenue").
		Aggr(dbox.AggrSum, "$ForcedOutages", "ForcedOutages").
		Aggr(dbox.AggrSum, "$ForcedOutagesLoss", "ForcedOutagesLoss").
		Aggr(dbox.AggrSum, "$UnforcedOutages", "UnforcedOutages").
		Aggr(dbox.AggrSum, "$UnforcedOutagesLoss", "UnforcedOutagesLoss").
		Aggr(dbox.AggrSum, "$TotalLabourCost", "TotalLabourCost").
		Aggr(dbox.AggrSum, "$TotalMaterialCost", "TotalMaterialCost").
		Aggr(dbox.AggrSum, "$TotalServicesCost", "TotalServicesCost").
		Aggr(dbox.AggrSum, "$MaintenanceCost", "MaintenanceCost").
		Aggr(dbox.AggrSum, "$OperatingCost", "OperatingCost").
		From(ve.TableName()).Group(groupBy).Cursor(nil)
	e = csr.Fetch(&result, 0, false)
	csr.Close()
	if e != nil {
		return nil, e
	}
	for _, i := range result {
		i.Set("_id", i.Get("_id").(tk.M).Get(groupBy))
	}
	// Remaining : sort by _id. descending for "result"
	return result, e
}
