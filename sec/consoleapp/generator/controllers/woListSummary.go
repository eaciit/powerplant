package controllers

import (
	"github.com/eaciit/crowd"
	"github.com/eaciit/dbox"
	//"github.com/eaciit/orm"
	//. "github.com/eaciit/powerplant/sec/consoleapp/generator/helpers"
	"strconv"
	"strings"
	"time"

	. "github.com/eaciit/powerplant/sec/library/models"
	tk "github.com/eaciit/toolkit"
)

// GenWOListSummary ...
type GenWOListSummary struct {
	*BaseController
}

// Generate ...
func (d *GenWOListSummary) Generate() {
	tk.Println("##Generating Summary Data..")
	e := d.generateDurationIntervalSummary()
	if e != nil {
		tk.Println(e)
	}
	tk.Println("##Summary Data : DONE\n")
}

// generateDurationIntervalSummary ...
func (d *GenWOListSummary) generateDurationIntervalSummary() error {
	years := [3]int{2013, 2014, 2015}

	c := d.Ctx.Connection

	csr, e := c.NewQuery().From(new(TempMstPlant).TableName()).Cursor(nil)

	if e != nil {
		return e
	} else {
		defer csr.Close()
	}

	MstPlantData := []tk.M{}
	e = csr.Fetch(&MstPlantData, 0, false)

	for _, year := range years {
		query := []*dbox.Filter{}

		yearFirst := strconv.Itoa(year)
		yearFirst = yearFirst + "-01-01 00:00:00.000"

		yearLast := strconv.Itoa(year + 1)
		yearLast = yearLast + "-01-01 00:00:00.000"

		query = append(query, dbox.And(dbox.Gte("ActualStart", yearFirst), dbox.Lte("ActualFinish", yearLast)))

		csr1, e := c.NewQuery().From(new(WOList).TableName()).Where(query...).Order("ActualStart").Cursor(nil)

		if e != nil {
			return e
		} else {
			defer csr1.Close()
		}

		datas := []tk.M{}
		e = csr1.Fetch(&datas, 0, false)

		if len(datas) > 0 {
			for _, data := range datas {
				woles := new(WOListSummary)
				woles.PeriodYear = year
				woles.OrderType = data.GetString("type")
				woles.FunctionalLocation = data.GetString("functionallocation")

				query = nil
				query = append(query, dbox.Eq("FunctionalLocationCode", data.GetString("functionallocation")))

				csr2, e := c.NewQuery().From(new(FunctionalLocation).TableName()).Where(query...).Select("CatProf").Cursor(nil)

				if e != nil {
					return e
				} else {
					defer csr2.Close()
				}

				tempResult := []tk.M{}
				e = csr2.Fetch(&tempResult, 0, false)

				if len(tempResult) > 0 {
					woles.EquipmentType = tempResult[0].GetString("catprof")
				}

				woles.MainenanceOrderCode = data.GetString("ordercode")
				woles.NotificationCode = data.GetString("notificationcode")

				query = nil
				query = append(query, dbox.Eq("FunctionalLocationCode", data.GetString("plant")))

				csr3, e := c.NewQuery().From(new(FunctionalLocation).TableName()).Where(query...).Select("Description").Cursor(nil)

				if e != nil {
					return e
				} else {
					defer csr3.Close()
				}

				tempResult = []tk.M{}
				e = csr3.Fetch(&tempResult, 0, false)

				if len(tempResult) > 0 {
					woles.Plant = tempResult[0].GetString("description")
				}

				woles.PlanStart, e = time.Parse(time.RFC3339, data["scheduledstart"].(string))

				if e != nil {
					return e
				}

				woles.PlanEnd, e = time.Parse(time.RFC3339, data["scheduledfinish"].(string))

				if e != nil {
					return e
				}

				subTime := woles.PlanEnd.Sub(woles.PlanStart)
				woles.PlanDuration = subTime.Hours()

				woles.ActualStart, e = time.Parse(time.RFC3339, data.GetString("actualstart"))
				if e != nil {
					return e
				}

				woles.ActualEnd, e = time.Parse(time.RFC3339, data.GetString("actualfinish"))
				if e != nil {
					return e
				}

				subTime = woles.ActualEnd.Sub(woles.ActualStart)
				woles.ActualDuration = subTime.Hours()

				actualstartPart := strings.Split(woles.ActualStart.String(), " ")
				actualstartPart = []string{actualstartPart[0], actualstartPart[1]}
				query = nil
				query = append(query, dbox.Lt("ActualFinish", strings.Join(actualstartPart, " ")))
				query = append(query, dbox.Eq("FunctionalLocation", woles.FunctionalLocation))
				csr4, e := c.NewQuery().Select("ActualFinish").From(new(WOList).TableName()).Order("-ActualFinish").Where(query...).Cursor(nil)

				if e != nil {
					return e
				} else {
					defer csr4.Close()
				}

				tempResult = []tk.M{}
				e = csr4.Fetch(&tempResult, 0, false)

				if len(tempResult) > 0 {
					woles.LastMaintenanceEnd, e = time.Parse(time.RFC3339, tempResult[0].GetString("actualfinish"))
				}

				if woles.ActualStart.String() != "" && woles.LastMaintenanceEnd.String() != "" {
					subTime = woles.ActualStart.Sub(woles.LastMaintenanceEnd)
					woles.LastMaintenanceInterval = subTime.Seconds() / 86400
				}

				woles.Cost = data.GetFloat64("actualcost")
				plantTypes := crowd.From(&MstPlantData).Where(func(x interface{}) interface{} {
					return x.(tk.M).GetString("plantcode") == data.GetString("plant")
				}).Exec().Result.Data().([]tk.M)

				if len(plantTypes) > 0 {
					woles.PlantType = plantTypes[0].GetString("plantcode")
				}

				for {
					e = d.Ctx.Insert(woles)
					if e == nil {
						break
					} else {
						d.Ctx.Connection.Connect()
					}
				}
			}
		}
	}
	return nil
}
