package controllers

import (
	"log"
	"strconv"
	"strings"

	"github.com/eaciit/crowd"
	"github.com/eaciit/dbox"
	. "github.com/eaciit/powerplant/sec/consoleapp/generator/helpers"
	. "github.com/eaciit/powerplant/sec/library/models"
	tk "github.com/eaciit/toolkit"
)

// GenPreventiveCorrectiveSummary
type GenPreventiveCorrectiveSummary struct {
	*BaseController
}

// Generate
func (s *GenPreventiveCorrectiveSummary) Generate(base *BaseController) {
	if base != nil {
		s.BaseController = base
	}

	tk.Println("##Generating PreventiveCorrectiveSummary..")
	e := s.generatePreventiveCorrectiveSummary()
	if e != nil {
		tk.Println(e)
	}
	tk.Println("##PreventiveCorrectiveSummary : DONE\n")
}

// generatePreventiveCorrectiveSummary
func (s *GenPreventiveCorrectiveSummary) generatePreventiveCorrectiveSummary() error {
	var e error
	ctx := s.BaseController.Ctx
	c := ctx.Connection

	years := [3]int{2013, 2014, 2015}
	sintax := "select Distinct(Element) from MORSummary"
	csr, e := c.NewQuery().Command("freequery", tk.M{}.Set("syntax", sintax)).Cursor(nil)
	defer csr.Close()

	if e != nil {
		return e
	}

	MROElements := []tk.M{}
	e = csr.Fetch(&MROElements, 0, false)

	csr1, e := c.NewQuery().From(new(MasterEquipmentType).TableName()).Cursor(nil)
	defer csr1.Close()

	if e != nil {
		return e
	}

	query := []*dbox.Filter{}

	for _, year := range years {
		yearFirst := strconv.Itoa(year)
		yearFirst = yearFirst + "-01-01 00:00:00.000"

		yearLast := strconv.Itoa(year + 1)
		yearLast = yearLast + "-01-01 00:00:00.000"

		query = append(query, dbox.And(dbox.Gte("Period", yearFirst), dbox.Lte("Period", yearLast)))

		csr2, e := c.NewQuery().From(new(MaintenanceCost).TableName()).Where(query...).Cursor(nil)
		defer csr2.Close()

		if e != nil {
			return e
		}

		datas := []tk.M{}
		e = csr2.Fetch(&datas, 0, false)

		Plants := crowd.From(&datas).Group(func(x interface{}) interface{} {
			return x.(tk.M).GetString("plant")
		}, nil).Exec().Result.Data().([]crowd.KV)

		if len(Plants) > 0 {
			for _, p := range Plants {
				plant := p.Key.(string)
				EqType := crowd.From(&datas).Where(func(x interface{}) interface{} {
					period := x.(tk.M).GetString("period")
					return strings.Contains(period, strconv.Itoa(year)) && x.(tk.M).GetString("plant") == plant
				}).Exec().Result.Data().([]tk.M)

				if len(EqType) > 0 {
					EquipmentTypes := crowd.From(&EqType).Group(func(x interface{}) interface{} {
						return x.(tk.M).GetString("equipmenttype")
					}, nil).Exec().Result.Data().([]crowd.KV)

					for _, eq := range EquipmentTypes {
						EquipmentType := eq.Key.(string)
						ActType := crowd.From(&EqType).Where(func(x interface{}) interface{} {
							return x.(tk.M).GetString("equipmenttype") == EquipmentType
						}).Exec().Result.Data().([]tk.M)

						if len(ActType) > 0 {
							MaintActivityTypes := crowd.From(&ActType).Group(func(x interface{}) interface{} {
								return x.(tk.M).GetString("maintenanceactivitytype")
							}, nil).Exec().Result.Data().([]crowd.KV)

							for _, act := range MaintActivityTypes {
								MaintActivityType := act.Key.(string)

								OrderType := crowd.From(&ActType).Where(func(x interface{}) interface{} {
									return x.(tk.M).GetString("maintenanceactivitytype") == MaintActivityType
								}).Exec().Result.Data().([]tk.M)

								if len(OrderType) > 0 {
									OrderTypes := crowd.From(&OrderType).Group(func(x interface{}) interface{} {
										return x.(tk.M).GetString("ordertype")
									}, nil).Exec().Result.Data().([]crowd.KV)

									for _, order := range OrderTypes {
										OrderTypeString := order.Key.(string)
										OrderNo := crowd.From(&OrderType).Where(func(x interface{}) interface{} {
											return x.(tk.M).GetString("ordertype") == OrderTypeString
										}).Exec().Result.Data().([]tk.M)
										if len(OrderNo) > 0 {
											Equipment := crowd.From(&OrderNo).Group(func(x interface{}) interface{} {
												return x.(tk.M).GetString("equipment")
											}, nil).Exec().Result.Data().([]crowd.KV)

											for _, eqNo := range Equipment {
												eqNoString := eqNo.Key.(string)
												for _, element := range MROElements {
													_ = element
													pcs := new(PreventiveCorrectiveSummary)
													pcs.PeriodYear = year
													pcs.OrderType = OrderTypeString
													pcs.EquipmentNo = eqNoString

													equipmentDescription := crowd.From(&OrderNo).Where(func(x interface{}) interface{} {
														return x.(tk.M).GetString("equipment") == eqNoString
													}).Exec().Result.Data().([]tk.M)

													if len(equipmentDescription) > 0 {
														pcs.EquipmentDescription = equipmentDescription[0].GetString("equipmentdesc")
													}

													if len(EquipmentTypes) == 1 {
														pcs.EquipmentType = "Other"
														pcs.EquipmentTypeDescription = "Other"
													} else {
														pcs.EquipmentType = EquipmentType
														pcs.EquipmentTypeDescription = equipmentDescription[0].GetString("equipmenttypedesc")
													}

													pcs.ActivityType = MaintActivityType
													pcs.Plant = PlantNormalization(plant)
													pcs.Element = element.GetString("element")

													result := float64(len(equipmentDescription) / len(MROElements))
													pcs.MOCount = int(Round(result, .5, 2))

													switch element.GetString("element") {
													case "Internal Labor":
														pcs.Value = crowd.From(&equipmentDescription).Sum(func(x interface{}) interface{} {
															return x.(tk.M).GetFloat64("internallaboractual")
														}).Exec().Result.Sum
														break
													case "Internal Material":
														pcs.Value = crowd.From(&equipmentDescription).Sum(func(x interface{}) interface{} {
															return x.(tk.M).GetFloat64("internalmaterialactual")
														}).Exec().Result.Sum
														break
													case "Direct Material":
														pcs.Value = crowd.From(&equipmentDescription).Sum(func(x interface{}) interface{} {
															return x.(tk.M).GetFloat64("directmaterialactual")
														}).Exec().Result.Sum
														break
													case "External Service":
														pcs.Value = crowd.From(&equipmentDescription).Sum(func(x interface{}) interface{} {
															return x.(tk.M).GetFloat64("externalserviceactual")
														}).Exec().Result.Sum
														break
													}

													e := ctx.InsertOut(pcs)

													if e != nil {
														log.Println(e.Error())
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}

		/*_ = year
		_ = csr2*/
	}

	return e
}
