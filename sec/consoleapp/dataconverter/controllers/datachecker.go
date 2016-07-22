package controllers

import (
	"github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/mongo"
	_ "github.com/eaciit/dbox/dbc/mssql"
	. "github.com/eaciit/powerplant/sec/library/models"
	tk "github.com/eaciit/toolkit"
	"time"
)

type DataChecker struct {
	*BaseController
}

func (m *DataChecker) CheckDetailData() error {
	tStart := time.Now()
	tk.Println("Start Checking DoValueEquationDashboard..")
	mod := new(ValueEquationDashboard)

	c, e := m.BaseController.MongoCtx.Connection.NewQuery().From(mod.TableName()).Cursor(nil)

	if e != nil {
		return e
	}

	defer c.Close()

	result := []tk.M{}
	e = c.Fetch(&result, 0, false)
	for _, val := range result {

		filter := []*dbox.Filter{}
		filter = append(filter, dbox.Eq("Dates", GetMgoValue(val, "Period.Dates").(time.Time).UTC()))
		filter = append(filter, dbox.Eq("Plant", GetMgoValue(val, "Plant")))
		filter = append(filter, dbox.Eq("Unit", GetMgoValue(val, "Unit")))

		csr, e := m.SqlCtx.Connection.NewQuery().From(mod.TableName()).Where(filter...).Cursor(nil)
		if e != nil {
			return e
		}
		tempResult := []tk.M{}
		e = csr.Fetch(&tempResult, 0, false)
		csr.Close()
		if e != nil {
			return e
		}

		if tempResult != nil && len(tempResult) > 0 {
			Id := tempResult[0].Get("id")

			Fuel := val.Get("Fuel")
			if Fuel != nil && len(Fuel.([]interface{})) > 0 {
				Len := len(Fuel.([]interface{}))
				filterDetail := []*dbox.Filter{}
				filterDetail = append(filterDetail, dbox.Eq("VEId", Id))
				csr, e := m.SqlCtx.Connection.NewQuery().Select("VEId").From(new(VEDFuel).TableName()).Where(filterDetail...).Cursor(nil)
				if e != nil {
					return e
				}
				detailResult := []tk.M{}
				e = csr.Fetch(&detailResult, 0, false)
				csr.Close()
				if e != nil {
					return e
				}

				if detailResult == nil || Len != len(detailResult) {
					tk.Println("FUEL DATA DID NOT MATCH. ID : ", Id)
				}
			}

			Detail := val.Get("Detail")
			if Detail != nil && len(Detail.([]interface{})) > 0 {
				Len := len(Detail.([]interface{}))
				filterDetail := []*dbox.Filter{}
				filterDetail = append(filterDetail, dbox.Eq("VEId", Id))
				csr, e := m.SqlCtx.Connection.NewQuery().Select("VEId").From(new(VEDDetail).TableName()).Where(filterDetail...).Cursor(nil)
				if e != nil {
					return e
				}
				detailResult := []tk.M{}
				e = csr.Fetch(&detailResult, 0, false)
				csr.Close()
				if e != nil {
					return e
				}

				if detailResult == nil || Len != len(detailResult) {
					tk.Println("DETAIL DATA DID NOT MATCH. ID : ", Id)
				}
			}

			Top10 := val.Get("Top10")
			if Top10 != nil && len(Top10.([]interface{})) > 0 {
				Len := len(Top10.([]interface{}))
				filterDetail := []*dbox.Filter{}
				filterDetail = append(filterDetail, dbox.Eq("VEId", Id))
				csr, e := m.SqlCtx.Connection.NewQuery().Select("VEId").From(new(VEDTop10).TableName()).Where(filterDetail...).Cursor(nil)
				if e != nil {
					return e
				}
				detailResult := []tk.M{}
				e = csr.Fetch(&detailResult, 0, false)
				csr.Close()
				if e != nil {
					return e
				}

				if detailResult == nil || Len != len(detailResult) {
					tk.Println("TOP 10 DATA DID NOT MATCH. ID : ", Id)
				}
			}

		}

	}

	tk.Printf("Checking complete in %v\n", time.Since(tStart))
	return nil
}
