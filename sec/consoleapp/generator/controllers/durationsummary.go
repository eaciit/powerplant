package controllers

import (
	//"github.com/eaciit/crowd"
	"github.com/eaciit/dbox"
	//"github.com/eaciit/orm"
	//. "github.com/eaciit/powerplant/sec/consoleapp/generator/helpers"
	. "github.com/eaciit/powerplant/sec/library/models"
	tk "github.com/eaciit/toolkit"
	//"log"
	"strconv"
	//"strings"
)

type DurationIntervalSummary struct {
	*BaseController
}

func (d *DurationIntervalSummary) Generate() {
	tk.Println("##Generating Summary Data..")
	e := d.GenerateDurationIntervalSummary()
	if e != nil {
		tk.Println(e)
	}
	tk.Println("##Summary Data : DONE\n")
}

func (d *DurationIntervalSummary) GenerateDurationIntervalSummary() error {
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
				_ = data
				/*woles := new(WOListSummary)
				  woles.PeriodYear = year
				  woles.OrderType = data.Type
				  woles.FunctionalLocation = data.FunctionalLocation

				  query = nil
				  query = append(query, dbox.Eq("FunctionalLocationCode", data.FunctionalLocation))

				  csr2, e := c.NewQuery().From(new(FunctionalLocation).TableName()).Where(query...).Select("CatProf").Cursor(nil)

				  if e != nil {
				     return e
				  } else {
				     defer csr1.Close()
				  }*/

				//datas := []tk.M{}
				//e = csr1.Fetch(&datas, 0, false)
				//woles.EquipmentType = DataHelper.Populate<FunctionalLocation>("FunctionalLocation", Query.EQ("FunctionalLocationCode", data.FunctionalLocation)).Select(x => x.CatProf).FirstOrDefault();
			}
		}
	}
	return nil
}
