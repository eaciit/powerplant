package controllers

import (
	"time"

	"github.com/eaciit/crowd"
	"github.com/eaciit/dbox"
	. "github.com/eaciit/powerplant/sec/consoleapp/generator/helpers"
	. "github.com/eaciit/powerplant/sec/library/models"
	tk "github.com/eaciit/toolkit"
)

// GenSummaryData
type GenWODurationSummary struct {
	*BaseController
}

// Generate ...
func (s *GenWODurationSummary) Generate(base *BaseController) {
	var (
		e error
	)
	if base != nil {
		s.BaseController = base
	}

	tk.Println("##Generating WODurationSummary Data..")

	e = s.generateDurationWorkOrderSummary()
	if e != nil {
		tk.Println(e)
	}

	e = s.generateDurationCostWorkOrderSummary()
	if e != nil {
		tk.Println(e)
	}
	tk.Println("##Summary Data : DONE\n")
}

func (s *GenWODurationSummary) generateDurationWorkOrderSummary() (e error) {

	// waiting for dhira, onprogress

	return
}

func (s *GenWODurationSummary) generateDurationCostWorkOrderSummary() error {
	ctx := s.BaseController.Ctx
	c := ctx.Connection
	var (
		query []*dbox.Filter
	)
	tk.Println("Generating Duration Cost Work Order Summary..")
	Years := []int{2013, 2014, 2015}
	query = []*dbox.Filter{}

	EqTypes := []MappedEquipmentType{}
	csr, e := c.NewQuery().From(new(MappedEquipmentType).TableName()).Cursor(nil)
	defer csr.Close()

	if e != nil {
		return e
	}

	e = csr.Fetch(&EqTypes, 0, false)
	if e != nil {
		return e
	}
	for _, year := range Years {
		query = []*dbox.Filter{}
		query = append(query, dbox.Gte("Period", time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)))
		query = append(query, dbox.Lt("Period", time.Date((year+1), 1, 1, 0, 0, 0, 0, time.UTC)))

		MaintenanceCostByHours := []MaintenanceCostByHour{}
		csr, e := c.NewQuery().From(new(MaintenanceCostByHour).TableName()).Where(query...).Cursor(nil)
		if e != nil {
			return e
		}
		e = csr.Fetch(&MaintenanceCostByHours, 0, false)
		if e != nil {
			return e
		}
		csr.Close()

		MaintenanceCostList := []MaintenanceCost{}
		csr, e = c.NewQuery().From(new(MaintenanceCost).TableName()).Where(query...).Cursor(nil)
		if e != nil {
			return e
		}
		e = csr.Fetch(&MaintenanceCostList, 0, false)
		if e != nil {
			return e
		}
		csr.Close()

		Plants := crowd.From(&MaintenanceCostByHours).Group(func(x interface{}) interface{} {
			return x.(MaintenanceCostByHour).Plant
		}, nil).Exec().Result.Data().([]crowd.KV)

		for _, p := range Plants {
			plant := p.Key.(string)
			EqType := crowd.From(&MaintenanceCostByHours).Where(func(x interface{}) interface{} {
				return x.(MaintenanceCostByHour).Plant == plant
			}).Group(func(x interface{}) interface{} {
				return x.(MaintenanceCostByHour).EquipmentType
			}, nil).Exec().Result.Data().([]crowd.KV)

			for _, eqt := range EqType {
				eq := eqt.Key.(string)
				ActType := crowd.From(&MaintenanceCostByHours).Where(func(x interface{}) interface{} {
					o := x.(MaintenanceCostByHour)
					return o.Plant == plant && o.EquipmentType == eq
				}).Group(func(x interface{}) interface{} {
					return x.(MaintenanceCostByHour).MaintenanceActivityType
				}, nil).Exec().Result.Data().([]crowd.KV)

				for _, a := range ActType {
					act := a.Key.(string)
					OrderType := crowd.From(&MaintenanceCostByHours).Where(func(x interface{}) interface{} {
						o := x.(MaintenanceCostByHour)
						return o.Plant == plant && o.EquipmentType == eq && o.MaintenanceActivityType == act
					}).Group(func(x interface{}) interface{} {
						return x.(MaintenanceCostByHour).OrderType
					}, nil).Exec().Result.Data().([]crowd.KV)

					for _, o := range OrderType {
						order := o.Key.(string)
						d := new(WODurationSummary)
						d.PeriodYear = year
						d.OrderType = order
						if len(eq) == 1 {
							d.EquipmentType = "Other"
						} else {
							d.EquipmentType = eq
						}
						if len(eq) == 1 {
							d.EquipmentTypeDescription = "Other"
						} else {
							EqTypeDesc := crowd.From(&MaintenanceCostByHours).Where(func(x interface{}) interface{} {
								o := x.(MaintenanceCostByHour)
								return o.Plant == plant && o.EquipmentType == eq
							}).Exec().Result.Data().([]MaintenanceCostByHour)
							if len(EqTypeDesc) > 0 {
								d.EquipmentTypeDescription = EqTypeDesc[0].EquipmentTypeDesc
							}
						}

						d.ActivityType = act
						d.Plant = PlantNormalization(plant)
						d.PlanValue = crowd.From(&MaintenanceCostByHours).Where(func(x interface{}) interface{} {
							o := x.(MaintenanceCostByHour)
							return o.Plant == plant && o.EquipmentType == eq && o.MaintenanceActivityType == act
						}).Sum(func(x interface{}) interface{} {
							return x.(MaintenanceCostByHour).PlanVal
						}).Exec().Result.Sum
						d.ActualValue = crowd.From(&MaintenanceCostByHours).Where(func(x interface{}) interface{} {
							o := x.(MaintenanceCostByHour)
							return o.Plant == plant && o.EquipmentType == eq && o.MaintenanceActivityType == act
						}).Sum(func(x interface{}) interface{} {
							return x.(MaintenanceCostByHour).Actual
						}).Exec().Result.Sum
						d.WOCount = len(OrderType)
						d.Cost = crowd.From(&MaintenanceCostList).Where(func(x interface{}) interface{} {
							o := x.(MaintenanceCost)
							return o.Plant == plant && o.EquipmentType == eq && o.MaintenanceActivityType == act
						}).Sum(func(x interface{}) interface{} {
							return x.(MaintenanceCost).PeriodTotalActual
						}).Exec().Result.Sum
						_, e := ctx.InsertOut(d)
						tk.Println("#")
						if e != nil {
							tk.Println(e)
							break
						}

					}

				}
			}
		}
	}
	return e
}
