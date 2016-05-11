package controllers

import (
	//. "github.com/eaciit/powerplant/sec/webapp/models"
	"strconv"
	"time"

	"github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
	"github.com/metakeule/fmtdate"
	"gopkg.in/mgo.v2/bson"
)

type DataBrowserVEController struct {
	*BaseController
}

func (c DataBrowserVEController) Default(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputTemplate
	c.LoadPartial(k, "shared/databrowser.html")

	q := struct {
		DBFields       interface{}
		SelectedFields interface{}
		PageId         string
		PageTitle      string
		Breadcrumbs    map[string]string
	}{}

	result := make([]tk.M, 0)

	cursor, e := c.DB().Connection.NewQuery().From("DataBrowserSelectedFields").Where(dbox.Eq("H3", "ValueEquation")).Cursor(nil)
	defer cursor.Close()

	cursor.Fetch(&result, 0, true)

	if len(result) == 0 {
		q.DBFields = ""
		q.SelectedFields = ""
	} else {
		DBFields := make([]tk.M, 0)

		cursor = nil
		cursor, e = c.DB().Connection.NewQuery().From("DataBrowserFields").Where(dbox.Eq("_id", result[0].Get("FieldsReference"))).Cursor(nil)
		cursor.Fetch(&DBFields, 0, true)

		q.DBFields = DBFields[0].Get("Fields")
		q.SelectedFields = result[0].Get("SelectedFields")
	}

	q.PageId = "Dashboard"
	q.PageTitle = "Dashboard"
	q.Breadcrumbs = make(map[string]string, 0)

	_ = e
	return q
}

func (c *DataBrowserVEController) GetGridDB(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	var e error
	r := new(tk.Result)

	d := struct {
		hypoid     string
		From       string
		To         string
		Period     string
		PeriodFrom string
		PeriodTo   string
		EQType     []string
	}{}

	r.Run(func(in interface{}) (interface{}, error) {

		//r := this.Ctx.Request
		hypoid := d.hypoid

		var (
			pipesgrid  []bson.M
			pipes      []bson.M
			query      []bson.M
			pipescount []bson.M
		)

		p := tk.M{}
		pgrid := tk.M{}
		pcount := tk.M{}
		ret := tk.M{}
		pbuildgrid := tk.M{}
		pbuild := tk.M{}

		if hypoid == "H16" {
			FromPeriod := d.From
			ToPeriod := d.To

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
		} else if d.Period == "" {
			query = append(query, bson.M{"Period.Year": bson.M{"$eq": 2014}})
		} else if d.Period != "" {
			selectedPeriod, _ := strconv.Atoi(d.Period)
			query = append(query, bson.M{"Period.Year": bson.M{"$eq": selectedPeriod}})
		} else {
			PeriodFrom, _ := strconv.Atoi(d.PeriodFrom)
			PeriodTo, _ := strconv.Atoi(d.PeriodTo)
			query = append(query, bson.M{"Period.Year": bson.M{"$gte": PeriodFrom}})
			query = append(query, bson.M{"Period.Year": bson.M{"$lte": PeriodTo}})
		}

		if len(d.EQType) > 0 {
			query = append(query, bson.M{"EquipmentType": bson.M{"$in": d.EQType}})
		} else {
			query = append(query, bson.M{"EquipmentType": bson.M{"$ne": "xxx"}})
			// query = append(query, bson.M{"isTurbine": bson.M{"$eq": true}})
		}

		/*if r.Form["Plant[]"] != nil {
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
			pipescount = append(pipescount, bson.M{"$match": bson.M{"$and": query}})
		}

		fields := r.Form["fields[]"]
		fieldsdouble := r.Form["fieldsdouble[]"]

		for _, fi := range fields {
			pbuildgrid.Set(fi, 1)
		}

		for _, fi := range fieldsdouble {
			pbuild.Set(strings.Replace(fi, ".", "", -1)+"sum", bson.M{"$sum": "$" + fi})
			pbuild.Set(strings.Replace(fi, ".", "", -1)+"avg", bson.M{"$avg": "$" + fi})
		}

		//Cek Hypo Unwind

		if hypoid == "H2" || hypoid == "H3" || hypoid == "H6" || hypoid == "H15" || hypoid == "H18" || hypoid == "H1" || hypoid == "H4" {

			pipesgrid = append(pipesgrid, bson.M{"$unwind": "$Maintenance"})

			pipescount = append(pipescount, bson.M{"$unwind": "$Maintenance"})

			pipes = append(pipes, bson.M{"$unwind": "$Maintenance"})

			//where after unwind
			if r.Form["OrderType[]"] != nil {
				pipes = append(pipes, bson.M{"$match": bson.M{"Maintenance.WorkOrderType": bson.M{"$in": r.Form["OrderType[]"]}}})
				pipesgrid = append(pipesgrid, bson.M{"$match": bson.M{"Maintenance.WorkOrderType": bson.M{"$in": r.Form["OrderType[]"]}}})
				pipescount = append(pipescount, bson.M{"$match": bson.M{"Maintenance.WorkOrderType": bson.M{"$in": r.Form["OrderType[]"]}}})
			}

		} else if hypoid == "H8" || hypoid == "H10" {

			pipesgrid = append(pipesgrid, bson.M{"$unwind": "$MROElement"})

			pipescount = append(pipescount, bson.M{"$unwind": "$MROElement"})

			pipes = append(pipes, bson.M{"$unwind": "$MROElement"})

			//where after unwind
			if r.Form["OrderType[]"] != nil {
				pipes = append(pipes, bson.M{"$match": bson.M{"MROElement.MROOrderType": bson.M{"$in": r.Form["OrderType[]"]}}})
				pipesgrid = append(pipesgrid, bson.M{"$match": bson.M{"MROElement.MROOrderType": bson.M{"$in": r.Form["OrderType[]"]}}})
				pipescount = append(pipescount, bson.M{"$match": bson.M{"MROElement.MROOrderType": bson.M{"$in": r.Form["OrderType[]"]}}})
			}
		} else if hypoid == "H17" {
			pipesgrid = append(pipesgrid, bson.M{"$unwind": "$FailureNotification"})
			pipescount = append(pipescount, bson.M{"$unwind": "$FailureNotification"})
			pipes = append(pipes, bson.M{"$unwind": "$FailureNotification"})

			//Where After Unwind
			if r.Form["FailureCode[]"] != nil {
				pipes = append(pipes, bson.M{"$match": bson.M{"FailureNotification.FailureCode": bson.M{"$in": r.Form["FailureCode[]"]}}})
				pipesgrid = append(pipesgrid, bson.M{"$match": bson.M{"FailureNotification.FailureCode": bson.M{"$in": r.Form["FailureCode[]"]}}})
				pipescount = append(pipescount, bson.M{"$match": bson.M{"FailureNotification.FailureCode": bson.M{"$in": r.Form["FailureCode[]"]}}})
			}

		} else if hypoid == "H16" {
			pipesgrid = append(pipesgrid, bson.M{"$unwind": "$TurbineVibrations"})
			pipescount = append(pipescount, bson.M{"$unwind": "$TurbineVibrations"})
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
		pipescount = pipesgrid
		pipes = pipesgrid

		take, _ := strconv.Atoi(r.FormValue("take"))
		skip, _ := strconv.Atoi(r.FormValue("skip"))

		pbuild.Set("_id", "")
		pipes = append(pipes, bson.M{"$group": pbuild})
		p.Set("pipe", pipes)
		result := make([]tk.M, 0)

		if len(fieldsdouble) == 0 {
			ret.Set("Summary", result)
		} else if skip == 0 {
			e = this.Db.Table("DataBrowser", p).FetchAll(&result, true)
			this.SetSession(hypoid+"summary", result)
			ret.Set("Summary", result)
		} else {
			ret.Set("Summary", this.GetSession(hypoid+"summary", nil))
		}

		sortfield := r.FormValue("sort[0][field]")
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
		}

		pipesgrid = append(pipesgrid, bson.M{"$skip": skip})
		pipesgrid = append(pipesgrid, bson.M{"$limit": take})

		pgrid.Set("pipe", pipesgrid)
		datas := make([]tk.M, 0)
		fmt.Println(pipesgrid)
		e = this.Db.Table("DataBrowser", pgrid).FetchAll(&datas, true)
		ret.Set("Datas", datas)

		pipescount = append(pipescount, bson.M{"$group": bson.M{"_id": "", "count": bson.M{"$sum": 1}}})
		pcount.Set("pipe", pipescount)
		tot := make([]tk.M, 0)

		if skip == 0 {
			e = this.Db.Table("DataBrowser", pcount).FetchAll(&tot, true)
			if len(tot) == 0 {
				ret.Set("Total", 0)
			} else {
				ret.Set("Total", tot[0].GetInt("count"))
				this.SetSession(hypoid+"count", tot[0].GetInt("count"))
			}
		} else {
			ret.Set("Total", this.GetSession(hypoid+"count", nil))
		}*/

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

		return "", e
	}, nil)
	//this.Json(r)
	return ""
}
