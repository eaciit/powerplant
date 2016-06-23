package controllers

import (
	"bufio"
	"fmt"
	"os"

	"github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	. "github.com/eaciit/powerplant/sec/library/models"
	tk "github.com/eaciit/toolkit"
	//"gopkg.in/mgo.v2/bson"
	"strconv"
	"strings"
	"time"

	"github.com/tealeg/xlsx"
)

var (
	SP_ALIAS  = "RES."
	SQLSCRIPT = "sqlscript"
)

type DataBrowserController struct {
	*BaseController
}

type Result struct {
	PageId             string
	HypothesisCategory string
	HypothesisId       string
	PageTitle          string
	DBFields           interface{}
	SelectedFields     interface{}
}

func (this *DataBrowserController) Default(k *knot.WebContext) interface{} {
	this.LoadPartial(k, "shared/databrowser.html")

	k.Config.OutputType = knot.OutputTemplate

	result := make([]tk.M, 0)
	cursor, _ := this.DB().Connection.NewQuery().From(new(DataBrowserSelectedFields).TableName()).Where(dbox.Eq("Hypothesis", "H3")).Cursor(nil)

	_ = cursor.Fetch(&result, 0, true)

	result1 := &Result{}

	if len(result) == 0 {
		result1.DBFields = ""
		result1.SelectedFields = ""
	} else {
		DBFields := make([]tk.M, 0)

		cursor = nil
		cursor, _ = this.DB().Connection.NewQuery().From(new(DataBrowserFields).TableName()).Where(dbox.Eq("FieldsReference", result[0].GetString("fieldsreference"))).Cursor(nil)

		cursor.Fetch(&DBFields, 0, true)

		defer cursor.Close()

		var selectedFields []string

		for _, str := range result {
			selectedFields = append(selectedFields, str.GetString("selectedfields"))
		}

		result1.DBFields = DBFields
		result1.SelectedFields = selectedFields
	}

	result1.PageId = "DataBrowser"
	result1.HypothesisCategory = ""
	result1.HypothesisId = "H3"
	result1.PageTitle = "Data Browser"

	return result1
}

func constructOperator(operator string) string {
	switch operator {
	case "eq":
		return "="
	case "neq":
		return "<>"
	case "startswith":
		return "like"
	case "contains":
		return "like"
	case "endswith":
		return "like"
	case "doesnotcontain":
		return "not like"
	case "gt":
		return ">"
	case "gte":
		return ">="
	case "lte":
		return "<="
	case "lt":
		return "<"
	}

	return ""
}

func (this *DataBrowserController) GetGridDb(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	//var e error

	params := tk.M{}
	r := new(tk.Result)

	d := struct {
		EQType       []string
		FailureCode  []string
		Fields       []string
		Fieldsdouble []string
		Hypoid       string
		OrderType    []string
		Page         int
		PageSize     int
		Period       string
		Plant        []string
		Skip         int
		Take         int
		Top          int
		Filter       tk.M
		Sort         []tk.M
	}{}

	_ = k.GetPayload(&d)

	r.Run(func(in interface{}) (interface{}, error) {
		hypoid := d.Hypoid
		procedure := tk.M{}
		orderParam := []string{}
		procedureParam := tk.M{}

		ret := tk.M{}

		if d.Period == "" {
			params.Set("@PeriodYear", " AND PeriodYear = 2014")
			params.Set("@PeriodFrom", "'1/1/2014  00:00:00.000'")
			params.Set("@PeriodTo", "'1/1/2015  00:00:00.000'")

		} else if d.Period != "" {
			selectedPeriod, _ := strconv.Atoi(d.Period)
			params.Set("@PeriodYear", " AND PeriodYear = "+d.Period)
			params.Set("@PeriodFrom", "'1/1/"+d.Period+"  00:00:00.000'")
			params.Set("@PeriodTo", "'1/1/"+strconv.Itoa(selectedPeriod+1)+"  00:00:00.000'")
		}

		if len(d.EQType) > 0 {
			str := strings.Join(d.EQType, ",")
			params.Set("@EquipmentType", " AND EquipmentType IN ("+str+")")
		} else {
			params.Set("@EquipmentType", " AND EquipmentType <> 'xxx'")
		}

		if len(d.Plant) > 0 {
			str := strings.Join(d.Plant, ",")
			params.Set("@PlantName", " AND PlantName IN ("+str+")")
		} else {
			params.Set("@PlantName", " AND PlantName <> ''")
		}

		// fields := d.Fields
		fieldsdouble := d.Fieldsdouble

		if len(d.OrderType) > 0 {
			str := strings.Join(d.OrderType, ",")
			params.Set("@WorkOrderType", " AND WO.Type IN ("+str+")")
		} else {
			params.Set("@WorkOrderType", " 1=1")
		}

		/*var filterindex = 0
		var filterclause = ""
		for r.FormValue("filter[filters]["+strconv.Itoa(filterindex)+"][field]") != "" {
			var filteroperator = r.FormValue("filter[filters][" + strconv.Itoa(filterindex) + "][operator]")

			if filteroperator != "" && filteroperator == "eq" {
				filterclause += r.FormValue("filter[filters]["+strconv.Itoa(filterindex)+"][field]") + " : " + r.FormValue("filter[filters]["+strconv.Itoa(filterindex)+"][value]") + ", "
				pipesgrid = append(pipesgrid, bson.M{"$match": bson.M{r.FormValue("filter[filters][" + strconv.Itoa(filterindex) + "][field]"): bson.M{"$eq": r.FormValue("filter[filters][" + strconv.Itoa(filterindex) + "][value]")}}})
			} else if filteroperator != "" && filteroperator == "neq" {
				filterclause += r.FormValue("filter[filters]["+strconv.Itoa(filterindex)+"][field]") + " : " + r.FormValue("filter[filters]["+strconv.Itoa(filterindex)+"][value]") + ", "
				pipesgrid = append(pipesgrid, bson.M{"$match": bson.M{r.FormValue("filter[filters][" + strconv.Itoa(filterindex) + "][field]"): bson.M{"$ne": r.FormValue("filter[filters][" + strconv.Itoa(filterindex) + "][value]")}}})
			} else if filteroperator != "" && (filteroperator == "startswith" || filteroperator == "contains" || filteroperator == "endswith") {
				filterclause += r.FormValue("filter[filters]["+strconv.Itoa(filterindex)+"][field]") + " : " + r.FormValue("filter[filters]["+strconv.Itoa(filterindex)+"][value]") + ", "
				pipesgrid = append(pipesgrid, bson.M{"$match": bson.M{r.FormValue("filter[filters][" + strconv.Itoa(filterindex) + "][field]"): bson.M{"$regex": r.FormValue("filter[filters][" + strconv.Itoa(filterindex) + "][value]")}}})
			} else if filteroperator != "" && filteroperator == "doesnotcontain" {
				filterclause += r.FormValue("filter[filters]["+strconv.Itoa(filterindex)+"][field]") + " : " + r.FormValue("filter[filters]["+strconv.Itoa(filterindex)+"][value]") + ", "
				pipesgrid = append(pipesgrid, bson.M{"$match": bson.M{r.FormValue("filter[filters][" + strconv.Itoa(filterindex) + "][field]"): bson.M{"not": bson.M{"$eq": "/." + r.FormValue("filter[filters]["+strconv.Itoa(filterindex)+"][value]"+"./")}}}})
			} else if filteroperator != "" && filteroperator == "gt" {
				var val, _ = strconv.ParseFloat(r.FormValue("filter[filters]["+strconv.Itoa(filterindex)+"][value]"), 64)
				filterclause += r.FormValue("filter[filters]["+strconv.Itoa(filterindex)+"][field]") + " : " + r.FormValue("filter[filters]["+strconv.Itoa(filterindex)+"][value]") + ", "
				pipesgrid = append(pipesgrid, bson.M{"$match": bson.M{r.FormValue("filter[filters][" + strconv.Itoa(filterindex) + "][field]"): bson.M{"$gt": val}}})
			} else if filteroperator != "" && filteroperator == "gte" {
				var val, _ = strconv.ParseFloat(r.FormValue("filter[filters]["+strconv.Itoa(filterindex)+"][value]"), 64)

				filterclause += r.FormValue("filter[filters]["+strconv.Itoa(filterindex)+"][field]") + " : " + r.FormValue("filter[filters]["+strconv.Itoa(filterindex)+"][value]") + ", "
				pipesgrid = append(pipesgrid, bson.M{"$match": bson.M{r.FormValue("filter[filters][" + strconv.Itoa(filterindex) + "][field]"): bson.M{"$gte": val}}})
			} else if filteroperator != "" && filteroperator == "lte" {
				var val, _ = strconv.ParseFloat(r.FormValue("filter[filters]["+strconv.Itoa(filterindex)+"][value]"), 64)

				filterclause += r.FormValue("filter[filters]["+strconv.Itoa(filterindex)+"][field]") + " : " + r.FormValue("filter[filters]["+strconv.Itoa(filterindex)+"][value]") + ", "
				pipesgrid = append(pipesgrid, bson.M{"$match": bson.M{r.FormValue("filter[filters][" + strconv.Itoa(filterindex) + "][field]"): bson.M{"$lte": val}}})
			} else if filteroperator != "" && filteroperator == "lt" {
				var val, _ = strconv.ParseFloat(r.FormValue("filter[filters]["+strconv.Itoa(filterindex)+"][value]"), 64)

				filterclause += r.FormValue("filter[filters]["+strconv.Itoa(filterindex)+"][field]") + " : " + r.FormValue("filter[filters]["+strconv.Itoa(filterindex)+"][value]") + ", "
				pipesgrid = append(pipesgrid, bson.M{"$match": bson.M{r.FormValue("filter[filters][" + strconv.Itoa(filterindex) + "][field]"): bson.M{"$lt": val}}})
			}

			filterindex += 1
		}*/

		/*pipescount = pipesgrid
		pipes = pipesgrid

		take := d.Take
		skip := d.Skip

		pbuild.Set("_id", "")
		pipes = append(pipes, tk.M{"$group": pbuild})
		p.Set("pipe", pipes)
		result := make([]tk.M, 0)

		if len(fieldsdouble) == 0 {
			ret.Set("Summary", result)
		} else if skip == 0 {
			curr, _ := this.DB().Connection.NewQuery().From("DataBrowser").Command("pipe", pipes).Cursor(nil)
			defer curr.Close()
			_ = curr.Fetch(&result, 0, true)

			k.SetSession(hypoid+"summary", result)
			ret.Set("Summary", result)
		} else {
			ret.Set("Summary", k.Session(hypoid+"summary", nil))
		}*/

		// ret.Set("Summary", []tk.M{})

		//sortfield := r.FormValue("sort[0][field]")
		//dir := r.FormValue("sort[0][dir]")
		/*if sortfield == "" {
			sortfield = "Plant.PlantName"
		}
		var sort int
		sort = 1
		if dir != "" && dir != "asc" {
			sort = -1
		}*/

		/*var sortindex = 0
		var sortclause = ""
		for r.FormValue("sort["+strconv.Itoa(sortindex)+"][field]") != "" {
			var sorttextdir = r.FormValue("sort[" + strconv.Itoa(sortindex) + "][dir]")
			var sortdir = 1

			if sorttextdir != "" && sorttextdir != "asc" {
				sortdir = -1
			}

			sortclause += r.FormValue("sort["+strconv.Itoa(sortindex)+"][field]") + " : " + strconv.Itoa(sortdir) + ", "
			pipesgrid = append(pipesgrid, bson.M{"$sort": bson.M{r.FormValue("sort[" + strconv.Itoa(sortindex) + "][field]"): sortdir}})
			sortindex += 1
		}

		if sortclause == "" {
			pipesgrid = append(pipesgrid, bson.M{"$sort": bson.M{sortfield: sort}})
		}*/

		/*pipesgrid = append(pipesgrid, tk.M{"$skip": skip})
		pipesgrid = append(pipesgrid, tk.M{"$limit": take})

		pgrid.Set("pipe", pipesgrid)
		datas := make([]tk.M, 0)
		fmt.Println(pipesgrid)

		curr1, _ := this.DB().Connection.NewQuery().From("DataBrowser").Command("pipe", pipesgrid).Cursor(nil)
		defer curr1.Close()
		_ = curr1.Fetch(&datas, 0, true)

		ret.Set("Datas", datas)

		pipescount = append(pipescount, tk.M{"$group": tk.M{"_id": "", "count": tk.M{"$sum": 1}}})
		pcount.Set("pipe", pipescount)
		tot := make([]tk.M, 0)

		if skip == 0 {
			curr1, _ := this.DB().Connection.NewQuery().From("DataBrowser").Command("pipe", pipescount).Cursor(nil)
			defer curr1.Close()
			_ = curr1.Fetch(&tot, 0, true)

			if len(tot) == 0 {
				ret.Set("Total", 0)
			} else {
				ret.Set("Total", tot[0].GetInt("count"))
				k.SetSession(hypoid+"count", tot[0].GetInt("count"))
			}
		} else {
			ret.Set("Total", k.Session(hypoid+"count", nil))
		}

		return ret, e

		_ = hypoid
		_ = pipesgrid
		_ = pipes
		_ = query
		_ = pipescount
		_ = p
		_ = pgrid
		_ = pcount
		_ = ret
		_ = pbuildgrid
		_ = pbuild

		result = []tk.M{}

		cursor, _ := this.DB().Connection.NewQuery().From("DataBrowser").Command("pipe", pipesgrid).Cursor(nil)

		defer cursor.Close()

		e = cursor.Fetch(&result, 0, true)

		return result, e*/

		if d.Page != 0 {
			if d.Page != 1 {
				params.Set("@Offset", strconv.Itoa((d.Page-1)*d.PageSize))
			} else {
				params.Set("@Offset", strconv.Itoa(d.Page-1))
			}

		} else {
			params.Set("@Offset", "0")
		}

		if d.PageSize != 0 {
			params.Set("@Limit", strconv.Itoa(d.PageSize))
		} else {
			params.Set("@Limit", "10")
		}

		// filter

		filtersParam := []string{}

		// filtersParam = append(filtersParam, "and")

		if d.Filter != nil {
			filters := d.Filter.Get("filters").([]interface{})
			logic := d.Filter.GetString("logic")

			for _, fil := range filters {
				// tk.Printf(": %#v %#v \n", idx, filter)

				val := fil.(map[string]interface{})

				// tk.Printf("val: %#v \n", val)

				filters2 := val["filters"]

				if filters2 != nil {
					listFilters2 := filters2.([]interface{})
					logic2 := val["logic"].(string)

					for _, fil2 := range listFilters2 {
						// tk.Printf("fil2: %#v \n", fil2)
						val2 := fil2.(map[string]interface{})

						field2 := val2["field"].(string)
						operator2 := val2["operator"].(string)
						value2 := val2["value"]

						filtersParam = append(filtersParam, SP_ALIAS+field2+" "+constructOperator(operator2)+" "+tk.ToString(value2)+" "+logic2)
					}
				} else {
					field := val["field"].(string)
					operator := val["operator"].(string)
					value := val["value"]

					filtersParam = append(filtersParam, SP_ALIAS+field+" "+constructOperator(operator)+" "+tk.ToString(value)+" "+logic)
				}
			}
		}

		if len(filtersParam) > 0 {
			fmt.Printf("len: %#v \n", len(filtersParam))
			strFilterParam := strings.Join(filtersParam, " ")
			strFilterParam = " WHERE 1=1 " + strFilterParam + " 1=1 "

			orderParam = append(orderParam, "@FILTERS")
			procedureParam.Set("@FILTERS", strFilterParam)
		}

		if len(d.Sort) > 0 {
			orderByStr := ""
			for _, val := range d.Sort {
				orderByStr += " ," + val.GetString("field") + " " + val.GetString("dir") + " "
			}

			orderByStr = strings.Replace(orderByStr, ",", "", 1)

			params.Set("@ORDERBY", " "+orderByStr+" ")
		} else {
			params.Set("@ORDERBY", " PlantPlantName asc ")
		}

		// get datas

		script := getSQLScript(SQLSCRIPT+"/databrowser_h3.sql", params)

		tk.Printf("---\n%#v \n----\n", script)

		cursor, e := this.DB().Connection.NewQuery().
			Command("freequery", tk.M{}.
				Set("syntax", script)).
			Cursor(nil)

		defer cursor.Close()

		datas := []SPDataBrowser{}

		e = cursor.Fetch(&datas, 0, true)

		if e != nil {
			tk.Println(e.Error())
		}

		ret.Set("Datas", datas)

		// get total and summary

		total := k.Session(hypoid+"Total", nil)
		summary := k.Session(hypoid+"Summary", nil)

		if total == nil || summary == nil || len(fieldsdouble) > 0 {
			summaryStr := " count(*) as Total"

			for _, val := range fieldsdouble {
				strSum := ",(Select CAST((Sum(result." + val + ")) as float)) as " + val + "sum"
				strAvg := ",(Select CAST((Avg(result." + val + ")) as float)) as " + val + "avg"

				summaryStr += strSum
				summaryStr += strAvg
			}

			params.Set("@Summary", summaryStr)

			script = getSQLScript(SQLSCRIPT+"/databrowser_h3_summary.sql", params)
			tk.Printf("---\n%#v \n----\n", script)

			cursorTotal, e := this.DB().Connection.NewQuery().
				Command("freequery", tk.M{}.
					Set("syntax", script)).
				Cursor(nil)

			_ = e

			defer cursorTotal.Close()

			resSum := []tk.M{}

			e = cursorTotal.Fetch(&resSum, 0, true)

			if len(resSum) > 0 {
				tmp := resSum[0]
				total = tmp.GetInt("total")
				tmpSummary := tk.M{}

				for _, val := range fieldsdouble {
					tmpSummary.Set(val+"avg", tmp.GetFloat64(strings.ToLower(val+"avg")))
					tmpSummary.Set(val+"sum", tmp.GetFloat64(strings.ToLower(val+"sum")))
				}

				summary = []tk.M{tmpSummary}

				k.SetSession(hypoid+"Total", total)
				k.SetSession(hypoid+"Summary", summary)

				ret.Set("Total", total)
				ret.Set("Summary", summary)
			}
		} else {
			ret.Set("Total", total)
			ret.Set("Summary", summary)
		}

		_ = procedure
		_ = fieldsdouble
		_ = summary

		return ret, e
	}, nil)

	return r
}

type SumList struct {
	field string
	tipe  string
}

func (this *DataBrowserController) SaveExcel(k *knot.WebContext) interface{} {

	d := struct {
		DisplayTypeCount int
		DisplayTypeList  []tk.M
		EQType           []string
		FailureCode      []string
		Fields           []string
		Fieldsdouble     []string
		Hypoid           string
		OrderType        []string
		Page             int
		PageSize         int
		Period           string
		Plant            []string
		Skip             int
		Take             int
		Top              int
		HeaderList       []string
	}{}

	_ = k.GetPayload(&d)

	var (
		selectedColumn []string
		DisplaySumList []SumList
	)
	r := new(tk.Result)
	r.Run(func(in interface{}) (interface{}, error) {

		hypoid := d.Hypoid

		var (
			pipesgrid  []tk.M
			pipes      []tk.M
			query      []tk.M
			pipescount []tk.M
		)

		p := tk.M{}
		pgrid := tk.M{}
		pcount := tk.M{}
		ret := tk.M{}
		pbuildgrid := tk.M{}
		pbuild := tk.M{}

		/*if hypoid == "H16" {
			FromPeriod := r.FormValue("From")
			ToPeriod := r.FormValue("To")
			var DFrom time.Time
			var DTo time.Time
			DFrom, _ = fmtdate.Parse("DD-MMM-YYYY hh:mm:ss", FromPeriod+" 00:00:00")
			DTo, _ = fmtdate.Parse("DD-MMM-YYYY hh:mm:ss", ToPeriod+" 00:00:00")

			y, _ := DFrom.ISOWeek()
			yy, _ := DTo.ISOWeek()

			var wy []bson.M
			wy = append(wy, bson.M{"Period.Year": bson.M{"$eq": y}})
			wy = append(wy, bson.M{"Period.Year": bson.M{"$eq": yy}})

			query = append(query, bson.M{"$or": wy})
		} */if d.Period == "" {
			query = append(query, tk.M{"Period.Year": tk.M{"$eq": 2014}})
		} else {
			selectedPeriod, _ := strconv.Atoi(d.Period)
			query = append(query, tk.M{"Period.Year": tk.M{"$eq": selectedPeriod}})
		}

		if len(d.EQType) > 0 {
			query = append(query, tk.M{"EquipmentType": tk.M{"$in": d.EQType}})
		} else {
			query = append(query, tk.M{"EquipmentType": tk.M{"$ne": "xxx"}})
		}

		if len(d.Plant) > 0 {
			query = append(query, tk.M{"Plant.PlantName": tk.M{"$in": d.Plant}})
		} else {
			query = append(query, tk.M{"Plant.PlantName": tk.M{"$ne": ""}})
		}

		//Cek Hypo Where
		if hypoid == "H2" {
			query = append(query, tk.M{"Maintenance": tk.M{"$ne": nil}})
			query = append(query, tk.M{"AssetType": tk.M{"$eq": "Steam"}})
		} else if hypoid == "H3" || hypoid == "H6" || hypoid == "H15" || hypoid == "H18" || hypoid == "H1" || hypoid == "H7" || hypoid == "H4" {
			query = append(query, tk.M{"Maintenance": tk.M{"$ne": nil}})
		} else if hypoid == "H8" || hypoid == "H10" {
			query = append(query, tk.M{"MROElement": tk.M{"$ne": nil}})
		} else if hypoid == "H17" {
			query = append(query, tk.M{"FailureNotification": tk.M{"$ne": nil}})
		} else if hypoid == "H16" {
			query = append(query, tk.M{"TurbineVibrations": tk.M{"$ne": nil}})
		}

		if query != nil && len(query) > 0 {
			pipesgrid = append(pipesgrid, tk.M{"$match": tk.M{"$and": query}})
			pipes = append(pipes, tk.M{"$match": tk.M{"$and": query}})
			pipescount = append(pipescount, tk.M{"$match": tk.M{"$and": query}})
		}

		fields := d.Fields
		fieldsdouble := d.Fieldsdouble

		for _, fi := range fields {
			pbuildgrid.Set(fi, 1)
			selectedColumn = append(selectedColumn, fi)
		}

		headerList := d.HeaderList
		for _, fi := range fieldsdouble {
			pbuild.Set(strings.Replace(fi, ".", "", -1)+"sum", tk.M{"$sum": "$" + fi})
			pbuild.Set(strings.Replace(fi, ".", "", -1)+"avg", tk.M{"$avg": "$" + fi})
		}

		//Cek Hypo Unwind

		if hypoid == "H2" || hypoid == "H3" || hypoid == "H6" || hypoid == "H15" || hypoid == "H18" || hypoid == "H1" || hypoid == "H4" {

			pipesgrid = append(pipesgrid, tk.M{"$unwind": "$Maintenance"})

			pipescount = append(pipescount, tk.M{"$unwind": "$Maintenance"})

			pipes = append(pipes, tk.M{"$unwind": "$Maintenance"})

			//where after unwind
			if len(d.OrderType) > 0 {
				pipes = append(pipes, tk.M{"$match": tk.M{"Maintenance.WorkOrderType": tk.M{"$in": d.OrderType}}})
				pipesgrid = append(pipesgrid, tk.M{"$match": tk.M{"Maintenance.WorkOrderType": tk.M{"$in": d.OrderType}}})
				pipescount = append(pipescount, tk.M{"$match": tk.M{"Maintenance.WorkOrderType": tk.M{"$in": d.OrderType}}})
			}

		} else if hypoid == "H8" || hypoid == "H10" {

			pipesgrid = append(pipesgrid, tk.M{"$unwind": "$MROElement"})

			pipescount = append(pipescount, tk.M{"$unwind": "$MROElement"})

			pipes = append(pipes, tk.M{"$unwind": "$MROElement"})

			//where after unwind
			if len(d.OrderType) > 0 {
				pipes = append(pipes, tk.M{"$match": tk.M{"MROElement.MROOrderType": tk.M{"$in": d.OrderType}}})
				pipesgrid = append(pipesgrid, tk.M{"$match": tk.M{"MROElement.MROOrderType": tk.M{"$in": d.OrderType}}})
				pipescount = append(pipescount, tk.M{"$match": tk.M{"MROElement.MROOrderType": tk.M{"$in": d.OrderType}}})
			}
		} else if hypoid == "H17" {
			pipesgrid = append(pipesgrid, tk.M{"$unwind": "$FailureNotification"})
			pipescount = append(pipescount, tk.M{"$unwind": "$FailureNotification"})
			pipes = append(pipes, tk.M{"$unwind": "$FailureNotification"})

			//Where After Unwind
			if len(d.FailureCode) > 0 {
				pipes = append(pipes, tk.M{"$match": tk.M{"FailureNotification.FailureCode": tk.M{"$in": d.FailureCode}}})
				pipesgrid = append(pipesgrid, tk.M{"$match": tk.M{"FailureNotification.FailureCode": tk.M{"$in": d.FailureCode}}})
				pipescount = append(pipescount, tk.M{"$match": tk.M{"FailureNotification.FailureCode": tk.M{"$in": d.FailureCode}}})
			}

		} else if hypoid == "H16" {
			pipesgrid = append(pipesgrid, tk.M{"$unwind": "$TurbineVibrations"})
			pipescount = append(pipescount, tk.M{"$unwind": "$TurbineVibrations"})
			pipes = append(pipes, tk.M{"$unwind": "$TurbineVibrations"})

			ppr := tk.M{}
			ppr.Set("Plant", 1)
			ppr.Set("TurbineVibrations", 1)

			pipesgrid = append(pipesgrid, tk.M{"$project": ppr})

			/*if r.FormValue("From") != "" && r.FormValue("To") != "" {
				FromPeriod := r.FormValue("From")
				ToPeriod := r.FormValue("To")
				var DFrom time.Time
				var DTo time.Time
				DFrom, _ = fmtdate.Parse("DD-MMM-YYYY hh:mm:ss", FromPeriod+" 00:00:00")
				DTo, _ = fmtdate.Parse("DD-MMM-YYYY hh:mm:ss", ToPeriod+" 00:00:00")

				delta := DTo.Sub(DFrom)
				daydiff := delta.Hours() / 24

				var queryweek []bson.M
				for ; daydiff > 0; daydiff-- {

					var queryAt []bson.M

					d := DFrom.AddDate(0, 0, int(daydiff))
					isoYear, isoWeek := d.ISOWeek()

					queryAt = append(queryAt, bson.M{"TurbineVibrations.WeekNo": bson.M{"$eq": isoWeek}})
					queryAt = append(queryAt, bson.M{"TurbineVibrations.Year": bson.M{"$eq": isoYear}})

					queryweek = append(queryweek, bson.M{"$and": queryAt})

				}

				pipesgrid = append(pipesgrid, bson.M{"$match": bson.M{"$or": queryweek}})
				pipescount = append(pipescount, bson.M{"$match": bson.M{"$or": queryweek}})
				pipes = append(pipes, bson.M{"$match": bson.M{"$or": queryweek}})
			}

			if r.Form["UnitNo[]"] != nil {
				old := r.Form["UnitNo[]"]
				newi := make([]interface{}, len(old))
				for i, v := range old {
					newi[i] = v
				}
				pipesgrid = append(pipesgrid, bson.M{"$match": bson.M{"TurbineVibrations.UnitNo": bson.M{"$in": newi}}})
				pipescount = append(pipescount, bson.M{"$match": bson.M{"TurbineVibrations.UnitNo": bson.M{"$in": newi}}})
				pipes = append(pipes, bson.M{"$match": bson.M{"TurbineVibrations.UnitNo": bson.M{"$in": newi}}})

			}*/

		}

		/*var filterindex = 0
		var filterclause = ""
		for r.FormValue("filter[filters]["+strconv.Itoa(filterindex)+"][field]") != "" {
			var filteroperator = r.FormValue("filter[filters][" + strconv.Itoa(filterindex) + "][operator]")

			if filteroperator != "" && filteroperator == "eq" {
				filterclause += r.FormValue("filter[filters]["+strconv.Itoa(filterindex)+"][field]") + " : " + r.FormValue("filter[filters]["+strconv.Itoa(filterindex)+"][value]") + ", "
				pipesgrid = append(pipesgrid, bson.M{"$match": bson.M{r.FormValue("filter[filters][" + strconv.Itoa(filterindex) + "][field]"): bson.M{"$eq": r.FormValue("filter[filters][" + strconv.Itoa(filterindex) + "][value]")}}})
			} else if filteroperator != "" && filteroperator == "neq" {
				filterclause += r.FormValue("filter[filters]["+strconv.Itoa(filterindex)+"][field]") + " : " + r.FormValue("filter[filters]["+strconv.Itoa(filterindex)+"][value]") + ", "
				pipesgrid = append(pipesgrid, bson.M{"$match": bson.M{r.FormValue("filter[filters][" + strconv.Itoa(filterindex) + "][field]"): bson.M{"$ne": r.FormValue("filter[filters][" + strconv.Itoa(filterindex) + "][value]")}}})
			} else if filteroperator != "" && (filteroperator == "startswith" || filteroperator == "contains" || filteroperator == "endswith") {
				filterclause += r.FormValue("filter[filters]["+strconv.Itoa(filterindex)+"][field]") + " : " + r.FormValue("filter[filters]["+strconv.Itoa(filterindex)+"][value]") + ", "
				pipesgrid = append(pipesgrid, bson.M{"$match": bson.M{r.FormValue("filter[filters][" + strconv.Itoa(filterindex) + "][field]"): bson.M{"$regex": r.FormValue("filter[filters][" + strconv.Itoa(filterindex) + "][value]")}}})
			} else if filteroperator != "" && filteroperator == "doesnotcontain" {
				filterclause += r.FormValue("filter[filters]["+strconv.Itoa(filterindex)+"][field]") + " : " + r.FormValue("filter[filters]["+strconv.Itoa(filterindex)+"][value]") + ", "
				pipesgrid = append(pipesgrid, bson.M{"$match": bson.M{r.FormValue("filter[filters][" + strconv.Itoa(filterindex) + "][field]"): bson.M{"not": bson.M{"$eq": "/." + r.FormValue("filter[filters]["+strconv.Itoa(filterindex)+"][value]"+"./")}}}})
			} else if filteroperator != "" && filteroperator == "gt" {
				var val, _ = strconv.ParseFloat(r.FormValue("filter[filters]["+strconv.Itoa(filterindex)+"][value]"), 64)
				filterclause += r.FormValue("filter[filters]["+strconv.Itoa(filterindex)+"][field]") + " : " + r.FormValue("filter[filters]["+strconv.Itoa(filterindex)+"][value]") + ", "
				pipesgrid = append(pipesgrid, bson.M{"$match": bson.M{r.FormValue("filter[filters][" + strconv.Itoa(filterindex) + "][field]"): bson.M{"$gt": val}}})
			} else if filteroperator != "" && filteroperator == "gte" {
				var val, _ = strconv.ParseFloat(r.FormValue("filter[filters]["+strconv.Itoa(filterindex)+"][value]"), 64)

				filterclause += r.FormValue("filter[filters]["+strconv.Itoa(filterindex)+"][field]") + " : " + r.FormValue("filter[filters]["+strconv.Itoa(filterindex)+"][value]") + ", "
				pipesgrid = append(pipesgrid, bson.M{"$match": bson.M{r.FormValue("filter[filters][" + strconv.Itoa(filterindex) + "][field]"): bson.M{"$gte": val}}})
			} else if filteroperator != "" && filteroperator == "lte" {
				var val, _ = strconv.ParseFloat(r.FormValue("filter[filters]["+strconv.Itoa(filterindex)+"][value]"), 64)

				filterclause += r.FormValue("filter[filters]["+strconv.Itoa(filterindex)+"][field]") + " : " + r.FormValue("filter[filters]["+strconv.Itoa(filterindex)+"][value]") + ", "
				pipesgrid = append(pipesgrid, bson.M{"$match": bson.M{r.FormValue("filter[filters][" + strconv.Itoa(filterindex) + "][field]"): bson.M{"$lte": val}}})
			} else if filteroperator != "" && filteroperator == "lt" {
				var val, _ = strconv.ParseFloat(r.FormValue("filter[filters]["+strconv.Itoa(filterindex)+"][value]"), 64)

				filterclause += r.FormValue("filter[filters]["+strconv.Itoa(filterindex)+"][field]") + " : " + r.FormValue("filter[filters]["+strconv.Itoa(filterindex)+"][value]") + ", "
				pipesgrid = append(pipesgrid, bson.M{"$match": bson.M{r.FormValue("filter[filters][" + strconv.Itoa(filterindex) + "][field]"): bson.M{"$lt": val}}})
			}

			filterindex += 1
		}*/
		pipescount = pipesgrid

		// take, _ := strconv.Atoi(r.FormValue("take"))
		skip := 1

		pbuild.Set("_id", "")
		pipes = append(pipes, tk.M{"$group": pbuild})
		p.Set("pipe", pipes)
		result := make([]tk.M, 0)

		if len(fieldsdouble) == 0 {
			ret.Set("Summary", result)
		} else if skip == 0 {
			curr, _ := this.DB().Connection.NewQuery().From("DataBrowser").Command("pipe", pipes).Cursor(nil)
			defer curr.Close()

			curr.Fetch(&result, 0, true)

			k.SetSession(hypoid+"summary", result)
			ret.Set("Summary", result)
		} else {
			ret.Set("Summary", k.Session(hypoid+"summary", nil))
		}

		/*sortfield := r.FormValue("sort[0][field]")
		dir := r.FormValue("sort[0][dir]")
		if sortfield == "" {
			sortfield = "Plant.PlantName"
		}
		var sort int
		sort = 1
		if dir != "" && dir != "asc" {
			sort = -1
		}

		var sortindex = 0
		var sortclause = ""
		for r.FormValue("sort["+strconv.Itoa(sortindex)+"][field]") != "" {
			var sorttextdir = r.FormValue("sort[" + strconv.Itoa(sortindex) + "][dir]")
			var sortdir = 1

			if sorttextdir != "" && sorttextdir != "asc" {
				sortdir = -1
			}

			sortclause += r.FormValue("sort["+strconv.Itoa(sortindex)+"][field]") + " : " + strconv.Itoa(sortdir) + ", "
			pipesgrid = append(pipesgrid, bson.M{"$sort": bson.M{r.FormValue("sort[" + strconv.Itoa(sortindex) + "][field]"): sortdir}})
			sortindex += 1
		}

		if sortclause == "" {
			pipesgrid = append(pipesgrid, bson.M{"$sort": bson.M{sortfield: sort}})
		}*/

		// pipesgrid = append(pipesgrid, bson.M{"$skip": skip})
		//pipesgrid = append(pipesgrid, bson.M{"$limit": 10})

		pgrid.Set("pipe", pipesgrid)
		datas := make([]tk.M, 0)

		curr, _ := this.DB().Connection.NewQuery().From("DataBrowser").Command("pipe", pipesgrid).Cursor(nil)
		defer curr.Close()

		curr.Fetch(&datas, 0, true)
		ret.Set("Datas", datas)

		pipescount = append(pipescount, tk.M{"$group": tk.M{"_id": "", "count": tk.M{"$sum": 1}}})
		pcount.Set("pipe", pipescount)
		tot := make([]tk.M, 0)

		if skip == 0 {
			curr, _ := this.DB().Connection.NewQuery().From("DataBrowser").Command("pipe", pipescount).Cursor(nil)
			defer curr.Close()

			curr.Fetch(&tot, 0, true)

			if len(tot) == 0 {
				ret.Set("Total", 0)
			} else {
				ret.Set("Total", tot[0].GetInt("count"))
				k.SetSession(hypoid+"count", tot[0].GetInt("count"))
			}
		} else {
			ret.Set("Total", k.Session(hypoid+"count", nil))
		}
		DisplayTypeCount := d.DisplayTypeCount
		DisplaySumList = []SumList{}
		for i := 0; i < DisplayTypeCount; i++ {
			sumData := SumList{}
			sumData.field = d.DisplayTypeList[i].Get("field").(string)
			sumData.tipe = d.DisplayTypeList[i].Get("tipe").(string)
			DisplaySumList = append(DisplaySumList, sumData)
		}

		excelFile, e := this.GenExcelFile(headerList, selectedColumn, datas, result, DisplaySumList)
		return excelFile, e
	}, nil)

	return r
}

func (this *DataBrowserController) GenExcelFile(header []string, selectedColumn []string, datas []tk.M, dataSummary []tk.M, DisplaySumList []SumList) (string, error) {
	today := time.Now().UTC()
	fileName := "static/files/databrowser_" + today.Format("2006-01-02T150405") + ".xlsx"
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	var err error

	file = xlsx.NewFile()
	sheet, _ = file.AddSheet("Sheet1")

	for i, data := range datas {
		if i == 0 {
			row = sheet.AddRow()
			for _, hdr := range header {
				cell = row.AddCell()
				cell.Value = hdr
			}
		}
		row = sheet.AddRow()
		for _, field := range selectedColumn {
			cell = row.AddCell()
			cell.SetValue(this.GetExcelValue(data, field))
		}
	}
	if DisplaySumList != nil && len(DisplaySumList) > 0 {
		var summary = dataSummary[0]

		row = sheet.AddRow()
		cell = row.AddCell()
		cell.SetValue("Average")
		cell.Merge(len(DisplaySumList)-1, 0)
		row = sheet.AddRow()
		for _, i := range DisplaySumList {
			cell = row.AddCell()
			if i.tipe == "string" || i.tipe == "date" {
				cell.SetValue("-")
			} else {
				field := strings.Replace(i.field, ".", "", -1) + "avg"
				cell.SetValue(summary.Get(field))
			}
		}
		row = sheet.AddRow()
		cell = row.AddCell()
		cell.Merge(len(DisplaySumList)-1, 0)
		cell.SetValue("Total")
		row = sheet.AddRow()
		for _, i := range DisplaySumList {
			cell = row.AddCell()
			if i.tipe == "string" || i.tipe == "date" {
				cell.SetValue("-")
			} else {
				field := strings.Replace(i.field, ".", "", -1) + "sum"
				cell.SetValue(summary.Get(field))
			}
		}
	}
	err = file.Save(fileName)
	// file := xlsx.NewFile()
	// sheet := file.AddSheet("Sheet1")
	// // header := []string{"Matnr", "Matkl"}
	// // for i, _ := range datas {
	// // 	if i == 0 {
	// // 		rowHeader := sheet.AddRow()
	// // 		for _, hdr := range header {
	// // 			cell := rowHeader.AddCell()
	// // 			cell.Value = hdr
	// // 		}
	// // 	}
	// // }
	// err := file.Save(fileName)
	return fileName, err
}

func (this *DataBrowserController) GetExcelValue(data tk.M, field string) interface{} {
	numberOfDot := strings.Count(field, ".")
	var result interface{}
	if numberOfDot > 0 {
		d := data.Get(field[0:strings.Index(field, ".")]).(tk.M)
		new_field := field[strings.Index(field, ".")+1 : len(field)]
		result = this.GetExcelValue(d, new_field)
	} else {
		result = data.Get(field)
	}
	if result == nil {
		result = ""
	}
	return result
}

func getSQLScript(path string, params tk.M) (script string) {

	file, err := os.Open(wd + path)
	if err == nil {
		defer file.Close()

		reader := bufio.NewReader(file)

		for {
			line, _, e := reader.ReadLine()
			if e != nil {
				break
			}

			script += string(line[:len(line)])
		}
	} else {
		fmt.Println(err.Error())
	}

	for idx, val := range params {
		script = strings.Replace(script, idx, val.(string), -1)
	}

	script = strings.Replace(script, "\t", "", -1)

	return
}

/*func (this *DataBrowserController) GetFilter() {
	var (
		e              error
		selectedColumn []string
	)
	r := new(tk.Result)
	r.Run(func(in interface{}) (interface{}, error) {

		r := this.Ctx.Request
		hypoid := r.FormValue("hypoid")

		var (
			pipesgrid []bson.M
			pipes     []bson.M
			query     []bson.M
		)

		pgrid := tk.M{}
		ret := tk.M{}
		pbuildgrid := tk.M{}
		pbuild := tk.M{}

		if hypoid == "H16" {
			FromPeriod := r.FormValue("From")
			ToPeriod := r.FormValue("To")
			var DFrom time.Time
			var DTo time.Time
			DFrom, _ = fmtdate.Parse("DD-MMM-YYYY hh:mm:ss", FromPeriod+" 00:00:00")
			DTo, _ = fmtdate.Parse("DD-MMM-YYYY hh:mm:ss", ToPeriod+" 00:00:00")

			y, _ := DFrom.ISOWeek()
			yy, _ := DTo.ISOWeek()

			var wy []bson.M
			wy = append(wy, bson.M{"Period.Year": bson.M{"$eq": y}})
			wy = append(wy, bson.M{"Period.Year": bson.M{"$eq": yy}})

			query = append(query, bson.M{"$or": wy})
		} else if r.FormValue("Period") == "" {
			query = append(query, bson.M{"Period.Year": bson.M{"$eq": 2014}})
		} else if r.FormValue("Period") != "" {
			selectedPeriod, _ := strconv.Atoi(r.FormValue("Period"))
			query = append(query, bson.M{"Period.Year": bson.M{"$eq": selectedPeriod}})
		} else {
			PeriodFrom, _ := strconv.Atoi(r.FormValue("PeriodFrom"))
			PeriodTo, _ := strconv.Atoi(r.FormValue("PeriodTo"))
			query = append(query, bson.M{"Period.Year": bson.M{"$gte": PeriodFrom}})
			query = append(query, bson.M{"Period.Year": bson.M{"$lte": PeriodTo}})
		}

		if r.Form["EQType[]"] != nil {
			query = append(query, bson.M{"EquipmentType": bson.M{"$in": r.Form["EQType[]"]}})
		} else {
			query = append(query, bson.M{"EquipmentType": bson.M{"$ne": "xxx"}})
			// query = append(query, bson.M{"isTurbine": bson.M{"$eq": true}})
		}

		if r.Form["Plant[]"] != nil {
			query = append(query, bson.M{"Plant.PlantName": bson.M{"$in": r.Form["Plant[]"]}})
		} else {
			query = append(query, bson.M{"Plant.PlantName": bson.M{"$ne": ""}})
		}

		//Cek Hypo Where
		if hypoid == "H2" {
			query = append(query, bson.M{"Maintenance": bson.M{"$ne": nil}})
			query = append(query, bson.M{"AssetType": bson.M{"$eq": "Steam"}})
		} else if hypoid == "H3" || hypoid == "H6" || hypoid == "H15" || hypoid == "H18" || hypoid == "H1" || hypoid == "H7" || hypoid == "H4" {
			query = append(query, bson.M{"Maintenance": bson.M{"$ne": nil}})
		} else if hypoid == "H8" || hypoid == "H10" {
			query = append(query, bson.M{"MROElement": bson.M{"$ne": nil}})
		} else if hypoid == "H17" {
			query = append(query, bson.M{"FailureNotification": bson.M{"$ne": nil}})
		} else if hypoid == "H16" {
			query = append(query, bson.M{"TurbineVibrations": bson.M{"$ne": nil}})
		}

		if query != nil && len(query) > 0 {
			pipesgrid = append(pipesgrid, bson.M{"$match": bson.M{"$and": query}})
			pipes = append(pipes, bson.M{"$match": bson.M{"$and": query}})
		}

		fields := r.Form["fields[]"]
		fieldsdouble := r.Form["fieldsdouble[]"]

		for _, fi := range fields {
			pbuildgrid.Set(fi, 1)
			selectedColumn = append(selectedColumn, fi)
		}

		for _, fi := range fieldsdouble {
			pbuild.Set(strings.Replace(fi, ".", "", -1)+"sum", bson.M{"$sum": "$" + fi})
			pbuild.Set(strings.Replace(fi, ".", "", -1)+"avg", bson.M{"$avg": "$" + fi})
		}

		//Cek Hypo Unwind

		if hypoid == "H2" || hypoid == "H3" || hypoid == "H6" || hypoid == "H15" || hypoid == "H18" || hypoid == "H1" || hypoid == "H4" {

			pipesgrid = append(pipesgrid, bson.M{"$unwind": "$Maintenance"})

			pipes = append(pipes, bson.M{"$unwind": "$Maintenance"})

			//where after unwind
			if r.Form["OrderType[]"] != nil {
				pipes = append(pipes, bson.M{"$match": bson.M{"Maintenance.WorkOrderType": bson.M{"$in": r.Form["OrderType[]"]}}})
				pipesgrid = append(pipesgrid, bson.M{"$match": bson.M{"Maintenance.WorkOrderType": bson.M{"$in": r.Form["OrderType[]"]}}})

			}

		} else if hypoid == "H8" || hypoid == "H10" {

			pipesgrid = append(pipesgrid, bson.M{"$unwind": "$MROElement"})

			pipes = append(pipes, bson.M{"$unwind": "$MROElement"})

			//where after unwind
			if r.Form["OrderType[]"] != nil {
				pipes = append(pipes, bson.M{"$match": bson.M{"MROElement.MROOrderType": bson.M{"$in": r.Form["OrderType[]"]}}})
				pipesgrid = append(pipesgrid, bson.M{"$match": bson.M{"MROElement.MROOrderType": bson.M{"$in": r.Form["OrderType[]"]}}})
			}
		} else if hypoid == "H17" {
			pipesgrid = append(pipesgrid, bson.M{"$unwind": "$FailureNotification"})
			pipes = append(pipes, bson.M{"$unwind": "$FailureNotification"})

			//Where After Unwind
			if r.Form["FailureCode[]"] != nil {
				pipes = append(pipes, bson.M{"$match": bson.M{"FailureNotification.FailureCode": bson.M{"$in": r.Form["FailureCode[]"]}}})
				pipesgrid = append(pipesgrid, bson.M{"$match": bson.M{"FailureNotification.FailureCode": bson.M{"$in": r.Form["FailureCode[]"]}}})
			}

		} else if hypoid == "H16" {
			pipesgrid = append(pipesgrid, bson.M{"$unwind": "$TurbineVibrations"})
			pipes = append(pipes, bson.M{"$unwind": "$TurbineVibrations"})

			ppr := tk.M{}
			ppr.Set("Plant", 1)
			ppr.Set("TurbineVibrations", 1)

			pipesgrid = append(pipesgrid, bson.M{"$project": ppr})

			if r.FormValue("From") != "" && r.FormValue("To") != "" {
				FromPeriod := r.FormValue("From")
				ToPeriod := r.FormValue("To")
				var DFrom time.Time
				var DTo time.Time
				DFrom, _ = fmtdate.Parse("DD-MMM-YYYY hh:mm:ss", FromPeriod+" 00:00:00")
				DTo, _ = fmtdate.Parse("DD-MMM-YYYY hh:mm:ss", ToPeriod+" 00:00:00")

				delta := DTo.Sub(DFrom)
				daydiff := delta.Hours() / 24

				var queryweek []bson.M
				for ; daydiff > 0; daydiff-- {

					var queryAt []bson.M

					d := DFrom.AddDate(0, 0, int(daydiff))
					isoYear, isoWeek := d.ISOWeek()

					queryAt = append(queryAt, bson.M{"TurbineVibrations.WeekNo": bson.M{"$eq": isoWeek}})
					queryAt = append(queryAt, bson.M{"TurbineVibrations.Year": bson.M{"$eq": isoYear}})

					queryweek = append(queryweek, bson.M{"$and": queryAt})

				}

				pipesgrid = append(pipesgrid, bson.M{"$match": bson.M{"$or": queryweek}})
				pipes = append(pipes, bson.M{"$match": bson.M{"$or": queryweek}})
			}

			if r.Form["UnitNo[]"] != nil {
				old := r.Form["UnitNo[]"]
				newi := make([]interface{}, len(old))
				for i, v := range old {
					newi[i] = v
				}
				pipesgrid = append(pipesgrid, bson.M{"$match": bson.M{"TurbineVibrations.UnitNo": bson.M{"$in": newi}}})
				pipes = append(pipes, bson.M{"$match": bson.M{"TurbineVibrations.UnitNo": bson.M{"$in": newi}}})

			}

		}

		var filterindex = 0
		var filterclause = ""
		for r.FormValue("filter[filters]["+strconv.Itoa(filterindex)+"][field]") != "" {
			var filteroperator = r.FormValue("filter[filters][" + strconv.Itoa(filterindex) + "][operator]")

			if filteroperator != "" && filteroperator == "eq" {
				filterclause += r.FormValue("filter[filters]["+strconv.Itoa(filterindex)+"][field]") + " : " + r.FormValue("filter[filters]["+strconv.Itoa(filterindex)+"][value]") + ", "
				pipesgrid = append(pipesgrid, bson.M{"$match": bson.M{r.FormValue("filter[filters][" + strconv.Itoa(filterindex) + "][field]"): bson.M{"$eq": r.FormValue("filter[filters][" + strconv.Itoa(filterindex) + "][value]")}}})
			} else if filteroperator != "" && filteroperator == "neq" {
				filterclause += r.FormValue("filter[filters]["+strconv.Itoa(filterindex)+"][field]") + " : " + r.FormValue("filter[filters]["+strconv.Itoa(filterindex)+"][value]") + ", "
				pipesgrid = append(pipesgrid, bson.M{"$match": bson.M{r.FormValue("filter[filters][" + strconv.Itoa(filterindex) + "][field]"): bson.M{"$ne": r.FormValue("filter[filters][" + strconv.Itoa(filterindex) + "][value]")}}})
			} else if filteroperator != "" && (filteroperator == "startswith" || filteroperator == "contains" || filteroperator == "endswith") {
				filterclause += r.FormValue("filter[filters]["+strconv.Itoa(filterindex)+"][field]") + " : " + r.FormValue("filter[filters]["+strconv.Itoa(filterindex)+"][value]") + ", "
				pipesgrid = append(pipesgrid, bson.M{"$match": bson.M{r.FormValue("filter[filters][" + strconv.Itoa(filterindex) + "][field]"): bson.M{"$regex": r.FormValue("filter[filters][" + strconv.Itoa(filterindex) + "][value]")}}})
			} else if filteroperator != "" && filteroperator == "doesnotcontain" {
				filterclause += r.FormValue("filter[filters]["+strconv.Itoa(filterindex)+"][field]") + " : " + r.FormValue("filter[filters]["+strconv.Itoa(filterindex)+"][value]") + ", "
				pipesgrid = append(pipesgrid, bson.M{"$match": bson.M{r.FormValue("filter[filters][" + strconv.Itoa(filterindex) + "][field]"): bson.M{"not": bson.M{"$eq": "/." + r.FormValue("filter[filters]["+strconv.Itoa(filterindex)+"][value]"+"./")}}}})
			} else if filteroperator != "" && filteroperator == "gt" {
				var val, _ = strconv.ParseFloat(r.FormValue("filter[filters]["+strconv.Itoa(filterindex)+"][value]"), 64)
				filterclause += r.FormValue("filter[filters]["+strconv.Itoa(filterindex)+"][field]") + " : " + r.FormValue("filter[filters]["+strconv.Itoa(filterindex)+"][value]") + ", "
				pipesgrid = append(pipesgrid, bson.M{"$match": bson.M{r.FormValue("filter[filters][" + strconv.Itoa(filterindex) + "][field]"): bson.M{"$gt": val}}})
			} else if filteroperator != "" && filteroperator == "gte" {
				var val, _ = strconv.ParseFloat(r.FormValue("filter[filters]["+strconv.Itoa(filterindex)+"][value]"), 64)

				filterclause += r.FormValue("filter[filters]["+strconv.Itoa(filterindex)+"][field]") + " : " + r.FormValue("filter[filters]["+strconv.Itoa(filterindex)+"][value]") + ", "
				pipesgrid = append(pipesgrid, bson.M{"$match": bson.M{r.FormValue("filter[filters][" + strconv.Itoa(filterindex) + "][field]"): bson.M{"$gte": val}}})
			} else if filteroperator != "" && filteroperator == "lte" {
				var val, _ = strconv.ParseFloat(r.FormValue("filter[filters]["+strconv.Itoa(filterindex)+"][value]"), 64)

				filterclause += r.FormValue("filter[filters]["+strconv.Itoa(filterindex)+"][field]") + " : " + r.FormValue("filter[filters]["+strconv.Itoa(filterindex)+"][value]") + ", "
				pipesgrid = append(pipesgrid, bson.M{"$match": bson.M{r.FormValue("filter[filters][" + strconv.Itoa(filterindex) + "][field]"): bson.M{"$lte": val}}})
			} else if filteroperator != "" && filteroperator == "lt" {
				var val, _ = strconv.ParseFloat(r.FormValue("filter[filters]["+strconv.Itoa(filterindex)+"][value]"), 64)

				filterclause += r.FormValue("filter[filters]["+strconv.Itoa(filterindex)+"][field]") + " : " + r.FormValue("filter[filters]["+strconv.Itoa(filterindex)+"][value]") + ", "
				pipesgrid = append(pipesgrid, bson.M{"$match": bson.M{r.FormValue("filter[filters][" + strconv.Itoa(filterindex) + "][field]"): bson.M{"$lt": val}}})
			}

			filterindex += 1
		}

		active_field := r.FormValue("active_field")

		// switch active_field {
		// case "Plant.PlantName":
		// 	pipesgrid = append(pipesgrid, bson.M{"$group": bson.M{"_id": "$" + active_field}})
		// 	break
		// default:
		// 	break
		// }
		pipesgrid = append(pipesgrid, bson.M{"$group": bson.M{"_id": "$" + active_field}})
		pgrid.Set("pipe", pipesgrid)
		datas := make([]tk.M, 0)
		csr := this.Db.Table("DataBrowser", pgrid)
		_ = csr.FetchAll(&datas, true)
		csr.Close()
		switch active_field {
		case "Plant.PlantName":
			for _, i := range datas {
				plant := make(tk.M)
				i.Set("Plant", plant.Set("PlantName", i.Get("_id")))
			}
			break
		case "Maintenance.WorkOrderType":
			for _, i := range datas {
				wotype := make(tk.M)
				i.Set("Maintenance", wotype.Set("WorkOrderType", i.Get("_id")))
			}
			break
		// case "Maintenance.MaintenanceOrder":
		// 	for _, i := range datas {
		// 		mo := make(tk.M)
		// 		i.Set("Maintenance", mo.Set("MaintenanceOrder", i.Get("_id")))
		// 	}
		// 	break
		// case "Maintenance.MaintenanceDescription":
		// 	for _, i := range datas {
		// 		md := make(tk.M)
		// 		i.Set("Maintenance", md.Set("MaintenanceDescription", i.Get("_id")))
		// 	}
		// 	break

		// case "EquipmentType":
		// 	for _, i := range datas {
		// 		i.Set("EquipmentType", i.Get("_id"))
		// 	}
		// 	break
		// case "EquipmentTypeDescription":
		// 	for _, i := range datas {
		// 		i.Set("EquipmentTypeDescription", i.Get("_id"))
		// 	}
		// 	break
		// case "FLDescription":
		// 	for _, i := range datas {
		// 		i.Set("FLDescription", i.Get("_id"))
		// 	}
		// 	break
		default:
			break
		}
		ret.Set("Datas", datas)

		return datas, e
	}, nil)
	this.Json(r)

}*/
