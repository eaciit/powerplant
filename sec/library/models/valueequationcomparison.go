package models

import (
	"strconv"
	"time"

	"github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	"github.com/eaciit/orm"
	tk "github.com/eaciit/toolkit"
)

type ValueEquationComparison struct {
	StartPeriod       time.Time
	EndPeriod         time.Time
	SelectedPlant     []string
	SelectedUnit      []string
	SelectedUnitGroup string
	Index             int
}

func (m *ValueEquationComparison) SetPayLoad(k *knot.WebContext) error {
	d := struct {
		StartPeriod       string
		EndPeriod         string
		SelectedPlant     []string
		SelectedUnit      []string
		SelectedUnitGroup string
		Index             string
	}{}
	e := k.GetPayload(&d)
	m.StartPeriod, _ = time.Parse(time.RFC3339, d.StartPeriod)
	m.EndPeriod, _ = time.Parse(time.RFC3339, d.EndPeriod)
	m.SelectedPlant = d.SelectedPlant
	m.SelectedUnit = d.SelectedUnit
	m.SelectedUnitGroup = d.SelectedUnitGroup
	m.Index, _ = strconv.Atoi(d.Index)
	return e
}

func (m *ValueEquationComparison) GetData(ctx *orm.DataContext, k *knot.WebContext) (interface{}, error) {
	m.SetPayLoad(k)
	type DataValue struct {
		Revenue         float64 `json:'Revenue'`
		MaintenanceCost float64 `json:'MaintenanceCost'`
		OperatingCost   float64 `json:'OperatingCost'`
		NetGeneration   float64 `json:'NetGeneration'`
	}
	result, Result := tk.M{}, []DataValue{}
	c := ctx.Connection
	ve := ValueEquation{}
	query := []*dbox.Filter{}
	query = append(query, dbox.Gte("Dates", m.StartPeriod))
	query = append(query, dbox.Lte("Dates", m.EndPeriod))
	if m.SelectedPlant != nil && len(m.SelectedPlant) > 0 {
		query = append(query, dbox.In("Plant", m.SelectedPlant))
	}
	if m.SelectedUnit != nil && len(m.SelectedUnit) > 0 {
		query = append(query, dbox.In("Unit", m.SelectedUnit))
	}
	if m.SelectedUnitGroup != "ALL" {
		query = append(query, dbox.In("UnitGroup", m.SelectedUnitGroup))
	}

	csr, e := c.NewQuery().
		Where(query...).
		Aggr(dbox.AggrSum, "Revenue", "Revenue").
		Aggr(dbox.AggrSum, "MaintenanceCost", "MaintenanceCost").
		Aggr(dbox.AggrSum, "OperatingCost", "OperatingCost").
		Aggr(dbox.AggrSum, "NetGeneration", "NetGeneration").
		From(ve.TableName()).Cursor(nil)
	e = csr.Fetch(&Result, 0, false)
	if e != nil {
		return nil, e
	}
	csr.Close()
	result.Set("Index", m.Index)
	result.Set("DataValue", Result)
	return result, nil
}
func (m *ValueEquationComparison) GetUnitList(ctx *orm.DataContext, k *knot.WebContext) (interface{}, error) {
	m.SetPayLoad(k)
	mup := MasterUnitPlant{}
	result, UnitData := tk.M{}, []MasterUnitPlant{}
	c := ctx.Connection
	query := []*dbox.Filter{}
	if m.SelectedPlant != nil && len(m.SelectedPlant) > 0 {
		query = append(query, dbox.In("Plant", m.SelectedPlant))
	}

	csr, e := c.NewQuery().
		Where(query...).Select("Unit").
		From(mup.TableName()).Group("Unit").Cursor(nil)
	e = csr.Fetch(&UnitData, 0, false)
	if e != nil {
		return nil, e
	}
	csr.Close()
	result.Set("Index", m.Index)
	result.Set("UnitData", UnitData)
	return result, e
}
