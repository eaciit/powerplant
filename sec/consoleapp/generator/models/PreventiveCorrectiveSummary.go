package models

import (
	"github.com/eaciit/crowd"
	"github.com/eaciit/dbox"
	"github.com/eaciit/orm"
	. "github.com/eaciit/powerplant/sec/library/models"
	tk "github.com/eaciit/toolkit"
	"log"
	"math"
	"strconv"
	"strings"
)

type PreventiveSummary struct {
}

func (r *PreventiveSummary) GeneratePreventiveCorrectiveSummary(ctx *orm.DataContext) error {
	var e error
	c := ctx.Connection

	years := [3]int{2013, 2014, 2015}
	sintax := "select Distinct(Element) from MORSummary"
	csr, e := c.NewQuery().Command("freequery", tk.M{}.Set("syntax", sintax)).Cursor(nil)

	if e != nil {
		return e
	} else {
		defer csr.Close()
	}

	MROElements := []tk.M{}
	e = csr.Fetch(&MROElements, 0, false)

	csr1, e := c.NewQuery().From(new(MasterEquipmentType).TableName()).Cursor(nil)

	if e != nil {
		return e
	} else {
		defer csr1.Close()
	}

	query := []*dbox.Filter{}

	for _, year := range years {
		yearFirst := strconv.Itoa(year)
		yearFirst = yearFirst + "-01-01 00:00:00.000"

		yearLast := strconv.Itoa(year + 1)
		yearLast = yearLast + "-01-01 00:00:00.000"

		query = append(query, dbox.And(dbox.Gte("Period", yearFirst), dbox.Lte("Period", yearLast)))

		csr2, e := c.NewQuery().From(new(MaintenanceCost).TableName()).Where(query...).Cursor(nil)

		if e != nil {
			return e
		} else {
			defer csr2.Close()
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
								return x.(tk.M).GetString("maintactivitytype")
							}, nil).Exec().Result.Data().([]crowd.KV)

							for _, act := range MaintActivityTypes {
								MaintActivityType := act.Key.(string)
								OrderType := crowd.From(&ActType).Where(func(x interface{}) interface{} {
									return x.(tk.M).GetString("maintactivitytype") == MaintActivityType
								}).Exec().Result.Data().([]tk.M)

								if len(OrderType) > 0 {
									OrderTypes := crowd.From(&OrderType).Group(func(x interface{}) interface{} {
										return x.(tk.M).GetString("ordertype")
									}, nil).Exec().Result.Data().([]crowd.KV)

									for _, order := range OrderTypes {
										OrderType := order.Key.(string)
										OrderNo := crowd.From(&OrderTypes).Where(func(x interface{}) interface{} {
											return x.(tk.M).GetString("ordertype") == OrderType
										}).Exec().Result.Data().([]tk.M)

										if len(OrderNo) > 0 {
											Equipment := crowd.From(&OrderNo).Group(func(x interface{}) interface{} {
												return x.(tk.M).GetString("equipment")
											}, nil).Exec().Result.Data().([]crowd.KV)

											for _, eqNo := range Equipment {
												eqNoString := eqNo.Key.(string)
												for _, element := range MROElements {
													pcs := new(PreventiveCorrectiveSummary)
													pcs.PeriodYear = year
													pcs.OrderType = OrderType
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
													pcs.Plant = r.PlantNormalization(plant)
													pcs.Element = element.GetString("element")

													result := float64(len(equipmentDescription[0].GetString("maintenanceorder")) / len(MROElements))
													pcs.MOCount = int(r.Round(result, .5, 2))

													switch element {
													case "Internal Labor":
														pcs.Value = crowd.From(&equipmentDescription).Sum(func(x interface{}) interface{} {
															return x.(tk.M).GetString("internallaboractual")
														}).Exec().Result.Sum
													case "Internal Material":
														pcs.Value = crowd.From(&equipmentDescription).Sum(func(x interface{}) interface{} {
															return x.(tk.M).GetString("internalmaterialactual")
														}).Exec().Result.Sum
													case "Direct Material":
														pcs.Value = crowd.From(&equipmentDescription).Sum(func(x interface{}) interface{} {
															return x.(tk.M).GetString("directmaterialactual")
														}).Exec().Result.Sum
													case "External Service":
														pcs.Value = crowd.From(&equipmentDescription).Sum(func(x interface{}) interface{} {
															return x.(tk.M).GetString("externalserviceactual")
														}).Exec().Result.Sum
													}

													_, e := ctx.InsertOut(pcs)

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

func (r *PreventiveSummary) Round(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}

func (r *PreventiveSummary) PlantNormalization(PlantName string) string {
	retVal := ""
	//switch (PlantName)
	//{
	//    case "Rabigh PP": retVal = "Rabigh"; break;
	//    case "QURAYYAH": retVal = "Qurayyah"; break;
	//    case "GHZLAN": retVal = "Ghazlan"; break;
	//    case "QURAYYAH CC": retVal = "Qurayyah CC"; break;
	//    case "Shuaiba Power Plant": retVal = "Shoaiba"; break;
	//    case "RABIGH POWER PLANT": retVal = "Rabigh"; break;
	//    case "Qurayyah Power Plant": retVal = "Qurayyah"; break;
	//    case "Qurayyah Steam": retVal = "Qurayyah"; break;
	//    case "GHAZLAN POWER PLANT": retVal = "Ghazlan"; break;
	//    default: retVal = PlantName; break;
	//}

	switch PlantName {
	case "POWER PLANT #9":
		retVal = "PP9"
	case "RABIGH POWER PLANT":
		retVal = "Rabigh"
	case "Rabigh 2":
		retVal = "Rabigh"
	case "Rabigh PP":
		retVal = "Rabigh"
	case "Shuaiba Power Plant":
		retVal = "Shoaiba"
	case "Sha'iba (CC)":
		retVal = "Shoaiba"
	case "Sha'iba (SEC)":
		retVal = "Shoaiba"
	case "GHAZLAN POWER PLANT":
		retVal = "Ghazlan"
	case "GHZLAN":
		retVal = "Ghazlan"
	case "Qurayyah Power Plant":
		retVal = "Qurayyah"
	case "Qurayyah -Steam":
		retVal = "Qurayyah"
	case "Qurayyah Combined Cycle Power Plant":
		retVal = "Qurayyah CC"
	case "Qurayyah- Combined Cycle":
		retVal = "Qurayyah CC"
	case "QurayyahCC":
		retVal = "Qurayyah CC"
	case "QURAYYAH CC":
		retVal = "Qurayyah CC"
	case "QurayyahPP":
		retVal = "Qurayyah"
	case "Qurayyah Steam":
		retVal = "Qurayyah"
	case "QURAYYAH":
		retVal = "Qurayyah"
	default:
		retVal = PlantName
	}

	return retVal
}
