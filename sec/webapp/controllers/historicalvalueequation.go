package controllers

import (
	"github.com/eaciit/knot/knot.v1"
	. "github.com/eaciit/powerplant/sec/library/models"
	// tk "github.com/eaciit/toolkit"
	// "strconv"
	// "gopkg.in/mgo.v2/bson"
	"time"

	"github.com/eaciit/dbox"
	"github.com/eaciit/orm"
	tk "github.com/eaciit/toolkit"
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
	type DataValue struct {
		ID              string
		Plant           string  `json:'Plant'`
		Unit            string  `json:'Unit'`
		Revenue         float64 `json:'Revenue'`
		MaintenanceCost float64 `json:'MaintenanceCost'`
		OperatingCost   float64 `json:'OperatingCost'`
	}
	result, DataChart, DataDetail := tk.M{}, []DataValue{}, []*DataValue{}
	c := ctx.Connection
	ve := new(ValueEquation)
	query := []*dbox.Filter{}
	query = append(query, dbox.Gte("Dates", m.StartPeriod))
	query = append(query, dbox.Lte("Dates", m.EndPeriod))

	csr, e := c.NewQuery().
		Where(query...).
		Select("Plant").
		Aggr(dbox.AggrSum, "Revenue", "Revenue").
		Aggr(dbox.AggrSum, "MaintenanceCost", "MaintenanceCost").
		Aggr(dbox.AggrSum, "OperatingCost", "OperatingCost").
		From(ve.TableName()).Group("Plant").Order("-Plant").Cursor(nil)
	if e != nil {
		return nil, e
	}

	if csr != nil {
		e = csr.Fetch(&DataChart, 0, false)
	}
	if e != nil {
		return nil, e
	}
	csr.Close()
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
			Where(query...).Select(groupBy).
			Aggr(dbox.AggrSum, "Revenue", "Revenue").
			Aggr(dbox.AggrSum, "MaintenanceCost", "MaintenanceCost").
			Aggr(dbox.AggrSum, "OperatingCost", "OperatingCost").
			From(ve.TableName()).Group(groupBy).Order("-" + groupBy).Cursor(nil)
		if csr != nil {
			e = csr.Fetch(&DataDetail, 0, false)
		}
		csr.Close()
		if e != nil {
			return nil, e
		}
		for _, i := range DataDetail {
			if i.Unit != "" {
				i.ID = i.Unit
			} else {
				i.ID = i.Plant
			}
		}
	}

	result.Set("DataDetail", DataDetail)
	return result, nil
}

func (m *HistoricalValueEquation) GetMaintenanceData(ctx *orm.DataContext, k *knot.WebContext) (interface{}, error) {
	m.SetPayLoad(k)
	type DataValue struct {
		ID            string
		DataSource    string  `json:'DataSource'`
		WorkOrderType string  `json:'WorkOrderType'`
		Plant         string  `json:'Plant'`
		Unit          string  `json:'Unit'`
		LaborCost     float64 `json:'LaborCost'`
		MaterialCost  float64 `json:'MaterialCost'`
		ServiceCost   float64 `json:'ServiceCost'`
	}

	result, DataMainEx, DataOrder, DataChart, DataTable, Temp, IDList := tk.M{}, []*DataValue{}, []*DataValue{}, []tk.M{}, []DataValue{}, []tk.M{}, []interface{}{}
	c := ctx.Connection
	query := []*dbox.Filter{}
	query = append(query, dbox.Gte("Dates", m.StartPeriod))
	query = append(query, dbox.Lte("Dates", m.EndPeriod))
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

	// Get DataMainEx
	csr, e := c.NewQuery().
		Where(query...).
		Select(groupBy).
		Aggr(dbox.AggrSum, "TotalLabourCost", "LaborCost").
		Aggr(dbox.AggrSum, "TotalMaterialCost", "MaterialCost").
		Aggr(dbox.AggrSum, "TotalServicesCost", "ServiceCost").
		From(new(ValueEquation).TableName()).Group(groupBy).Order("-" + groupBy).
		Order("DataSource").
		Cursor(nil)
	if e != nil {
		return nil, e
	}
	if csr != nil {
		e = csr.Fetch(&DataMainEx, 0, false)
	}
	csr.Close()
	if e != nil {
		return nil, e
	}
	for _, i := range DataMainEx {
		if i.Unit != "" {
			i.ID = i.Unit
		} else {
			i.ID = i.Plant
		}
	}

	csr, e = c.NewQuery().
		Where(query...).Select("Id").
		From(new(ValueEquation).TableName()).Cursor(nil)
	if e != nil {
		return nil, e
	}
	if csr != nil {
		e = csr.Fetch(&Temp, 0, false)
	}
	csr.Close()
	if e != nil {
		return nil, e
	}
	for _, i := range Temp {
		IDList = append(IDList, i.GetInt("id"))
	}

	query = []*dbox.Filter{}
	query = append(query, dbox.In("VEId", IDList...))

	// Get DataTable
	csr, e = c.NewQuery().
		Where(query...).
		Select("DataSource", "WorkOrderType").
		Aggr(dbox.AggrSum, "LaborCost", "LaborCost").
		Aggr(dbox.AggrSum, "MaterialCost", "MaterialCost").
		Aggr(dbox.AggrSum, "ServiceCost", "ServiceCost").
		From(new(ValueEquationDetails).TableName()).Group("DataSource", "WorkOrderType").
		Order("DataSource").
		Cursor(nil)
	if e != nil {
		return nil, e
	}
	if csr != nil {
		e = csr.Fetch(&DataTable, 0, false)
	}
	csr.Close()
	if e != nil {
		return nil, e
	}

	// Get DataOrder - For Visualisation
	csr, e = c.NewQuery().
		Where(query...).
		Select("WorkOrderType").
		From(new(ValueEquationTop10).TableName()).Group("WorkOrderType").Cursor(nil)
	if e != nil {
		return nil, e
	}
	if csr != nil {
		e = csr.Fetch(&DataOrder, 0, false)
	}
	csr.Close()
	if e != nil {
		return nil, e
	}
	for _, i := range DataOrder {
		if i.Unit != "" {
			i.ID = i.Unit
		} else {
			i.ID = i.Plant
		}
	}

	// Data Chart - set for empty, if you want to add it. its available on the previous version while its still using beego.
	result.Set("DataMainEx", DataMainEx).Set("DataOrder", DataOrder).Set("DataChart", DataChart).Set("DataTable", DataTable)
	return result, nil
}

func (m *HistoricalValueEquation) GetOperatingData(ctx *orm.DataContext, k *knot.WebContext) (interface{}, error) {
	type DataValue struct {
		ID                string
		Plant             string  `json:'Plant'`
		Unit              string  `json:'Unit'`
		OperatingCost     float64 `json:'OperatingCost'`
		FuelTransportCost float64 `json:'FuelTransportCost'`
		Capacity          float64 `json:'Capacity'`
		NetGeneration     float64 `json:'NetGeneration'`
	}

	m.SetPayLoad(k)
	result, DataChart, DataTable, DataTotal, DataDetail := tk.M{}, []tk.M{}, []tk.M{}, []*DataValue{}, []*DataValue{}
	c := ctx.Connection
	ve := ValueEquation{}
	query := []*dbox.Filter{}
	query = append(query, dbox.Gte("Dates", m.StartPeriod))
	query = append(query, dbox.Lte("Dates", m.EndPeriod))
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
		Aggr(dbox.AggrSum, "OperatingCost", "OperatingCost").
		Aggr(dbox.AggrSum, "FuelTransportCost", "FuelTransportCost").
		From(ve.TableName()).Cursor(nil)
	if csr != nil {
		e = csr.Fetch(&DataTotal, 0, false)
	}
	csr.Close()
	if e != nil {
		return nil, e
	}

	if m.Scope != "Unit" {
		csr, e = c.NewQuery().
			Where(query...).Select(groupBy).
			Aggr(dbox.AggrSum, "Capacity", "Capacity").
			Aggr(dbox.AggrSum, "NetGeneration", "NetGeneration").
			Aggr(dbox.AggrSum, "OperatingCost", "OperatingCost").
			From(ve.TableName()).Group(groupBy).Order(groupBy).Cursor(nil)
		if csr != nil {
			e = csr.Fetch(&DataDetail, 0, false)
		}
		csr.Close()
		if e != nil {
			return nil, e
		}
		for _, i := range DataDetail {
			if i.Unit != "" {
				i.ID = i.Unit
			} else {
				i.ID = i.Plant
			}
		}
	}

	result.Set("DataChart", DataChart).Set("DataTable", DataTable).Set("DataTotal", DataTotal).Set("DataDetail", DataDetail)
	return result, e
}

func (m *HistoricalValueEquation) GetRevenueData(ctx *orm.DataContext, k *knot.WebContext) (interface{}, error) {
	type DataValue struct {
		ID              string
		Plant           string  `json:'Plant'`
		Unit            string  `json:'Unit'`
		CapacityPayment float64 `json:'CapacityPayment'`
		EnergyPayment   float64 `json:'EnergyPayment'`
		StartupPayment  float64 `json:'StartupPayment'`
		PenaltyAmount   float64 `json:'PenaltyAmount'`
		Incentive       float64 `json:'Incentive'`
		Revenue         float64 `json:'Revenue'`
	}

	m.SetPayLoad(k)
	result, DataChartRevenue, DataChartRevenueEx := tk.M{}, []*DataValue{}, []*DataValue{}
	c := ctx.Connection
	ve := ValueEquation{}
	query := []*dbox.Filter{}
	query = append(query, dbox.Gte("Dates", m.StartPeriod))
	query = append(query, dbox.Lte("Dates", m.EndPeriod))
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
		Aggr(dbox.AggrSum, "CapacityPayment", "CapacityPayment").
		Aggr(dbox.AggrSum, "EnergyPayment", "EnergyPayment").
		Aggr(dbox.AggrSum, "StartupPayment", "StartupPayment").
		Aggr(dbox.AggrSum, "PenaltyAmount", "PenaltyAmount").
		Aggr(dbox.AggrSum, "Incentive", "Incentive").
		Aggr(dbox.AggrSum, "Revenue", "Revenue").
		From(ve.TableName()).Cursor(nil)
	if csr != nil {
		e = csr.Fetch(&DataChartRevenue, 0, false)
	}
	csr.Close()
	if e != nil {
		return nil, e
	}
	csr, e = c.NewQuery().
		Where(query...).Select(groupBy).
		Aggr(dbox.AggrSum, "CapacityPayment", "CapacityPayment").
		Aggr(dbox.AggrSum, "EnergyPayment", "EnergyPayment").
		Aggr(dbox.AggrSum, "StartupPayment", "StartupPayment").
		Aggr(dbox.AggrSum, "PenaltyAmount", "PenaltyAmount").
		Aggr(dbox.AggrSum, "Incentive", "Incentive").
		Aggr(dbox.AggrSum, "Revenue", "Revenue").
		From(ve.TableName()).Group(groupBy).Order("-" + groupBy).Cursor(nil)
	if csr != nil {
		e = csr.Fetch(&DataChartRevenueEx, 0, false)
	}
	csr.Close()
	if e != nil {
		return nil, e
	}
	for _, i := range DataChartRevenueEx {
		if i.Unit != "" {
			i.ID = i.Unit
		} else {
			i.ID = i.Plant
		}
	}
	return result.Set("DataChartRevenue", DataChartRevenue).Set("DataChartRevenueEx", DataChartRevenueEx), e
}

func (m *HistoricalValueEquation) GetDataQuality(ctx *orm.DataContext, k *knot.WebContext) (interface{}, error) {
	type DataValue struct {
		ID                       string
		Plant                    string  `json:'Plant'`
		Unit                     string  `json:'Unit'`
		Count                    float64 `json:'Count'`
		CapacityPayment_Data     float64 `json:'CapacityPayment_Data'`
		EnergyPayment_Data       float64 `json:'EnergyPayment_Data'`
		StartupPayment_Data      float64 `json:'StartupPayment_Data'`
		Penalty_Data             float64 `json:'Penalty_Data'`
		Incentive_Data           float64 `json:'Incentive_Data'`
		MaintenanceCost_Data     float64 `json:'MaintenanceCost_Data'`
		MaintenanceDuration_Data float64 `json:'MaintenanceDuration_Data'`
		PrimaryFuel1st_Data      float64 `json:'PrimaryFuel1st_Data'`
		PrimaryFuel2nd_Data      float64 `json:'PrimaryFuel2nd_Data'`
		BackupFuel_Data          float64 `json:'BackupFuel_Data'`
		FuelTransport_Data       float64 `json:'FuelTransport_Data'`
	}

	m.SetPayLoad(k)
	var e error = nil
	result := []*DataValue{}
	c := ctx.Connection
	vedq := ValueEquationDataQuality{}
	query := []*dbox.Filter{}
	query = append(query, dbox.Gte("Dates", m.StartPeriod))
	query = append(query, dbox.Lte("Dates", m.EndPeriod))
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
		query = append(query, dbox.Gte("Dates", m.StartPeriod))
		query = append(query, dbox.Lte("Dates", m.EndPeriod))
		if m.Scope == "Unit" {
			query = append(query, dbox.Eq("Plant", m.SelectedPlant))
			if m.Selected != nil && len(m.Selected) > 0 {
				query = append(query, dbox.In("Unit", m.Selected))
			}
		} else {
			query = append(query, dbox.Eq("Plant", m.Selected[0]))
		}
		temp := []*ValueEquationDataQuality{}
		csr, e := ctx.Find(new(ValueEquationDataQuality), tk.M{}.Set("where", dbox.And(query...)))
		if csr != nil {
			e = csr.Fetch(&temp, 0, false)
		}
		csr.Close()
		if e != nil {
			return nil, e
		} else {
			m.GetValueEquationDocument(ctx, temp)
			return temp, e
		}
	} else {
		csr, e := c.NewQuery().
			Where(query...).Select(groupBy).
			Aggr(dbox.AggrSum, "1", "Count").
			Aggr(dbox.AggrSum, "CapacityPayment_Data", "CapacityPayment_Data").
			Aggr(dbox.AggrSum, "EnergyPayment_Data", "EnergyPayment_Data").
			Aggr(dbox.AggrSum, "StartupPayment_Data", "StartupPayment_Data").
			Aggr(dbox.AggrSum, "Penalty_Data", "Penalty_Data").
			Aggr(dbox.AggrSum, "Incentive_Data", "Incentive_Data").
			Aggr(dbox.AggrSum, "MaintenanceCost_Data", "MaintenanceCost_Data").
			Aggr(dbox.AggrSum, "MaintenanceDuration_Data", "MaintenanceDuration_Data").
			Aggr(dbox.AggrSum, "PrimaryFuel1st_Data", "PrimaryFuel1st_Data").
			Aggr(dbox.AggrSum, "PrimaryFuel2nd_Data", "PrimaryFuel2nd_Data").
			Aggr(dbox.AggrSum, "BackupFuel_Data", "BackupFuel_Data").
			Aggr(dbox.AggrSum, "FuelTransport_Data", "FuelTransport_Data").
			From(vedq.TableName()).Group(groupBy).Order("-" + groupBy).Cursor(nil)
		if csr != nil {
			e = csr.Fetch(&result, 0, false)
		}
		csr.Close()
		if e != nil {
			return nil, e
		}
		for _, i := range result {
			if i.Unit != "" {
				i.ID = i.Unit
			} else {
				i.ID = i.Plant
			}
		}
	}
	return result, e
}

func (m *HistoricalValueEquation) GetValueEquationDocument(ctx *orm.DataContext, VEDQList []*ValueEquationDataQuality) {
	for _, i := range VEDQList {
		query := []*dbox.Filter{}
		query = append(query, dbox.Eq("VEId", int(i.Id)))

		CapacityPaymentDocuments := []VEDQCapacityPaymentDocuments{}
		csr, _ := ctx.Find(new(VEDQCapacityPaymentDocuments), tk.M{}.Set("where", dbox.And(query...)))
		if csr != nil {
			csr.Fetch(&CapacityPaymentDocuments, 0, false)
			i.CapacityPaymentDocuments = CapacityPaymentDocuments
		}
		csr.Close()

		BackupFuelDocuments := []VEDQBackupFuelDocuments{}
		csr, _ = ctx.Find(new(VEDQBackupFuelDocuments), tk.M{}.Set("where", dbox.And(query...)))
		if csr != nil {
			csr.Fetch(&BackupFuelDocuments, 0, false)
			i.BackupFuelDocuments = BackupFuelDocuments
		}
		csr.Close()

		EnergyPaymentDocuments := []VEDQEnergyPaymentDocuments{}
		csr, _ = ctx.Find(new(VEDQEnergyPaymentDocuments), tk.M{}.Set("where", dbox.And(query...)))
		if csr != nil {
			csr.Fetch(&EnergyPaymentDocuments, 0, false)
			i.EnergyPaymentDocuments = EnergyPaymentDocuments
		}
		csr.Close()

		MaintenanceCostDocuments := []VEDQMaintenanceCostDocuments{}
		csr, _ = ctx.Find(new(VEDQMaintenanceCostDocuments), tk.M{}.Set("where", dbox.And(query...)))
		if csr != nil {
			csr.Fetch(&MaintenanceCostDocuments, 0, false)
			i.MaintenanceCostDocuments = MaintenanceCostDocuments
		}
		csr.Close()

		MaintenanceDurationDocuments := []VEDQMaintenanceDurationDocuments{}
		csr, _ = ctx.Find(new(VEDQMaintenanceDurationDocuments), tk.M{}.Set("where", dbox.And(query...)))
		if csr != nil {
			csr.Fetch(&MaintenanceDurationDocuments, 0, false)
			i.MaintenanceDurationDocuments = MaintenanceDurationDocuments
		}
		csr.Close()

		PenaltyDocuments := []VEDQPenaltyDocuments{}
		csr, _ = ctx.Find(new(VEDQPenaltyDocuments), tk.M{}.Set("where", dbox.And(query...)))
		if csr != nil {
			csr.Fetch(&PenaltyDocuments, 0, false)
			i.PenaltyDocuments = PenaltyDocuments
		}
		csr.Close()

		IncentiveDocuments := []VEDQIncentiveDocuments{}
		csr, _ = ctx.Find(new(VEDQIncentiveDocuments), tk.M{}.Set("where", dbox.And(query...)))
		if csr != nil {
			csr.Fetch(&IncentiveDocuments, 0, false)
			i.IncentiveDocuments = IncentiveDocuments
		}
		csr.Close()

		PrimaryFuel1stDocuments := []VEDQPrimaryFuel1stDocuments{}
		csr, _ = ctx.Find(new(VEDQPrimaryFuel1stDocuments), tk.M{}.Set("where", dbox.And(query...)))
		if csr != nil {
			csr.Fetch(&PrimaryFuel1stDocuments, 0, false)
			i.PrimaryFuel1stDocuments = PrimaryFuel1stDocuments
		}
		csr.Close()

		PrimaryFuel2ndDocuments := []VEDQPrimaryFuel2ndDocuments{}
		csr, _ = ctx.Find(new(VEDQPrimaryFuel2ndDocuments), tk.M{}.Set("where", dbox.And(query...)))
		if csr != nil {
			csr.Fetch(&PrimaryFuel2ndDocuments, 0, false)
			i.PrimaryFuel2ndDocuments = PrimaryFuel2ndDocuments
		}
		csr.Close()

		StartupPaymentDocuments := []VEDQStartupPaymentDocuments{}
		csr, _ = ctx.Find(new(VEDQStartupPaymentDocuments), tk.M{}.Set("where", dbox.And(query...)))
		if csr != nil {
			csr.Fetch(&StartupPaymentDocuments, 0, false)
			i.StartupPaymentDocuments = StartupPaymentDocuments
		}
		csr.Close()
	}
}

func (m *HistoricalValueEquation) GetPerformanceData(ctx *orm.DataContext, k *knot.WebContext) (interface{}, error) {
	type DataValue struct {
		ID            string
		Plant         string  `json:'Plant'`
		Unit          string  `json:'Unit'`
		NetGeneration float64 `json:'NetGeneration'`
		PrctWAF       float64 `json:'PrctWAF'`
	}
	m.SetPayLoad(k)
	var e error = nil
	result := []*DataValue{}
	c := ctx.Connection
	ve := ValueEquation{}
	query := []*dbox.Filter{}
	query = append(query, dbox.Gte("Dates", m.StartPeriod))
	query = append(query, dbox.Lte("Dates", m.EndPeriod))
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
		Where(query...).Select(groupBy).
		Aggr(dbox.AggrSum, "NetGeneration", "NetGeneration").
		Aggr(dbox.AggrSum, "PrctWAF", "PrctWAF").
		From(ve.TableName()).Group(groupBy).Order(groupBy).Cursor(nil)
	if csr != nil {
		e = csr.Fetch(&result, 0, false)
	}
	csr.Close()
	if e != nil {
		return nil, e
	}
	for _, i := range result {
		if i.Unit != "" {
			i.ID = i.Unit
		} else {
			i.ID = i.Plant
		}
	}
	return result, e
}

func (m *HistoricalValueEquation) GetAssetWorkData(ctx *orm.DataContext, k *knot.WebContext) (interface{}, error) {
	type DataValue struct {
		ID                 string
		Plant              string  `json:'Plant'`
		Unit               string  `json:'Unit'`
		ValueEquationCost  float64 `json:'ValueEquationCost'`
		MaxPowerGeneration float64 `json:'MaxPowerGeneration'`
		PotentialRevenue   float64 `json:'PotentialRevenue'`
		NetGeneration      float64 `json:'NetGeneration'`
		Revenue            float64 `json:'Revenue'`
		ForcedOutages      float64 `json:'ForcedOutages'`
		ForcedOutagesLoss  float64 `json:'ForcedOutagesLoss'`
		UnforcedOutages    float64 `json:'UnforcedOutages'`

		UnforcedOutagesLoss float64 `json:'UnforcedOutagesLoss'`
		TotalLabourCost     float64 `json:'TotalLabourCost'`
		TotalMaterialCost   float64 `json:'TotalMaterialCost'`
		TotalServicesCost   float64 `json:'TotalServicesCost'`
		MaintenanceCost     float64 `json:'MaintenanceCost'`
		OperatingCost       float64 `json:'OperatingCost'`
	}

	m.SetPayLoad(k)
	var e error = nil
	result := []*DataValue{}
	c := ctx.Connection
	ve := ValueEquation{}
	query := []*dbox.Filter{}
	query = append(query, dbox.Gte("Dates", m.StartPeriod))
	query = append(query, dbox.Lte("Dates", m.EndPeriod))
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
		Where(query...).Select(groupBy).
		Aggr(dbox.AggrSum, "ValueEquationCost", "ValueEquationCost").
		Aggr(dbox.AggrSum, "MaxPowerGeneration", "MaxPowerGeneration").
		Aggr(dbox.AggrSum, "PotentialRevenue", "PotentialRevenue").
		Aggr(dbox.AggrSum, "NetGeneration", "NetGeneration").
		Aggr(dbox.AggrSum, "Revenue", "Revenue").
		Aggr(dbox.AggrSum, "ForcedOutages", "ForcedOutages").
		Aggr(dbox.AggrSum, "ForcedOutagesLoss", "ForcedOutagesLoss").
		Aggr(dbox.AggrSum, "UnforcedOutages", "UnforcedOutages").
		Aggr(dbox.AggrSum, "UnforcedOutagesLoss", "UnforcedOutagesLoss").
		Aggr(dbox.AggrSum, "TotalLabourCost", "TotalLabourCost").
		Aggr(dbox.AggrSum, "TotalMaterialCost", "TotalMaterialCost").
		Aggr(dbox.AggrSum, "TotalServicesCost", "TotalServicesCost").
		Aggr(dbox.AggrSum, "MaintenanceCost", "MaintenanceCost").
		Aggr(dbox.AggrSum, "OperatingCost", "OperatingCost").
		From(ve.TableName()).Group(groupBy).Order(groupBy).Cursor(nil)
	if csr != nil {
		e = csr.Fetch(&result, 0, false)
	}
	csr.Close()
	if e != nil {
		return nil, e
	}
	for _, i := range result {
		if i.Unit != "" {
			i.ID = i.Unit
		} else {
			i.ID = i.Plant
		}
	}
	return result, e
}
