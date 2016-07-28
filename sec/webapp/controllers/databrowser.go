package controllers

import (
	"bufio"
	"os"
	"sort"

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
	// SpAlias     = "RES."
	SQLScript    = "sqlscript"
	FiltersMain  = []string{"MaintenanceInterval"}
	FiltersPlant = []string{"PlantPlantName"}
	// FiltersWO    = []string{"ActualDuration", "PlanDuration", "WorkOrderType", "ActualFinish", "ActualStart", "MaintenanceCost", "MaintenanceDescription", "MaintenanceOrder", "PlanFinish", "PlanStart"}
	FiltersWO = []string{"ActualDuration", "WorkOrderType", "ActualFinish", "ActualStart", "MaintenanceCost", "MaintenanceDescription", "MaintenanceOrder"}
	FiltersDB = []string{"EquipmentType", "EquipmentTypeDescription", "FLDescription", "FunctionalLocation"}
	FilWO     = []string{}
	FilDB     = []string{}
	FilMain   = []string{}
	FilPlant  = []string{}
	FieldStr  = []string{"PlantPlantName", "WorkOrderType"}
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

type DataBrowserInput struct {
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
	Filter           tk.M
	Sort             []tk.M
	DisplayTypeCount int
	DisplayTypeList  []tk.M
	HeaderList       []string
}

func (this *DataBrowserController) Default(k *knot.WebContext) interface{} {
	if k.Session("userid") == nil {
		this.Redirect(k, "login", "default")
	}

	this.LoadPartial(k, "shared/databrowser.html")

	k.Config.OutputType = knot.OutputTemplate

	result := make([]tk.M, 0)
	cursor, _ := this.DB().Connection.NewQuery().From(new(DataBrowserSelectedFields).TableName()).Where(dbox.Eq("Hypothesis", "H3")).Order("Orders").Cursor(nil)
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

func constructFilters(operator string, field string, value interface{}, logic string) {
	condition := ""

	switch operator {
	case "eq":
		condition = "="
	case "neq":
		condition = "<>"
	case "startswith":
		condition = "like"
	case "contains":
		condition = "like"
	case "endswith":
		condition = "like"
	case "doesnotcontain":
		condition = "not like"
	case "gt":
		condition = ">"
	case "gte":
		condition = ">="
	case "lte":
		condition = "<="
	case "lt":
		condition = "<"
	}

	sort.Strings(FiltersDB)
	sort.Strings(FiltersWO)
	sort.Strings(FiltersMain)
	sort.Strings(FiltersPlant)
	sort.Strings(FieldStr)

	val := tk.ToString(value)

	str := sort.SearchStrings(FieldStr, field)

	if str < len(FieldStr) && FieldStr[str] == field {
		val = "'" + val + "'"
	}

	db := sort.SearchStrings(FiltersDB, field)
	wo := sort.SearchStrings(FiltersWO, field)
	main := sort.SearchStrings(FiltersMain, field)
	plant := sort.SearchStrings(FiltersPlant, field)

	if db < len(FiltersDB) && FiltersDB[db] == field {
		FilDB = append(FilDB, " "+logic+" "+field+" "+condition+" "+val)
	} else if wo < len(FiltersWO) && FiltersWO[wo] == field {
		FilWO = append(FilWO, " "+logic+" "+"RESULT."+field+" "+condition+" "+val+" ")
	} else if main < len(FiltersMain) && FiltersMain[main] == field {
		FilMain = append(FilMain, " "+logic+" "+field+" "+condition+" "+val+" ")
	} else if plant < len(FiltersPlant) && FiltersPlant[plant] == field {
		FilPlant = append(FilPlant, " "+logic+" "+field+" "+condition+" "+val+" ")
	}
}

func getDataBrowser(d DataBrowserInput) (params tk.M, e error) {

	FilWO = []string{}
	FilDB = []string{}
	FilMain = []string{}
	FilPlant = []string{}
	params = tk.M{}

	if d.Period == "" {
		// params.Set("@PeriodYear", " PeriodYear = 2014")
		params.Set("@PeriodFrom", "'1/1/2014  00:00:00.000'")
		params.Set("@PeriodTo", "'1/1/2015  00:00:00.000'")
	} else if d.Period != "" {
		selectedPeriod, _ := strconv.Atoi(d.Period)
		// params.Set("@PeriodYear", " PeriodYear = "+d.Period)
		params.Set("@PeriodFrom", "'1/1/"+d.Period+"  00:00:00.000'")
		params.Set("@PeriodTo", "'1/1/"+strconv.Itoa(selectedPeriod+1)+"  00:00:00.000'")
	}

	if len(d.EQType) > 0 {
		str := strings.Join(d.EQType, "','")
		str = "'" + str + "'"
		params.Set("@EquipmentType", " EquipmentType IN ("+str+")")
	} else {
		params.Set("@EquipmentType", " EquipmentType <> 'xxx'")
	}

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

	if d.Filter != nil {
		filters := d.Filter.Get("filters").([]interface{})
		logic := d.Filter.GetString("logic")

		for _, fil := range filters {
			val := fil.(map[string]interface{})
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

					constructFilters(operator2, field2, value2, logic2)
				}
			} else {
				field := val["field"].(string)
				operator := val["operator"].(string)
				value := val["value"]

				constructFilters(operator, field, value, logic)
			}
		}
	}

	if len(FilDB) > 0 {
		strFilter := strings.Join(FilDB, " ")
		params.Set("@FILTERS_DB", strFilter)
	} else {
		params.Set("@FILTERS_DB", "")
	}

	if len(FilWO) > 0 {
		tmpFilter := strings.Split(FilWO[0], " ")
		FilWO[0] = strings.Join(tmpFilter[2:], " ")
		strFilter := strings.Join(FilWO, " ")
		params.Set("@FILTERS_WO", " WHERE "+strFilter)
	} else {
		params.Set("@FILTERS_WO", "")
	}

	if len(FilMain) > 0 {
		tmpFilter := strings.Split(FilMain[0], " ")
		FilMain[0] = strings.Join(tmpFilter[2:], " ")
		strFilter := strings.Join(FilMain, " ")
		params.Set("@FILTERS_MAIN", " WHERE "+strFilter)
	} else {
		params.Set("@FILTERS_MAIN", "")
	}

	if len(d.Plant) > 0 {
		str := strings.Join(d.Plant, "','")
		str = "'" + str + "'"
		params.Set("@PlantName", " WHERE PlantPlantName IN ("+str+")")
	} else if len(FilPlant) > 0 {
		tmpFilter := strings.Split(FilPlant[0], " ")
		FilPlant[0] = strings.Join(tmpFilter[2:], " ")
		strFilter := strings.Join(FilPlant, " ")
		params.Set("@PlantName", " WHERE "+strFilter)
	} else {
		params.Set("@PlantName", " WHERE PlantPlantName <> ''")
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

	return
}

func (this *DataBrowserController) GetGridDb(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	r := new(tk.Result)
	d := DataBrowserInput{}
	_ = k.GetPayload(&d)
	ret := tk.M{}

	r.Run(func(in interface{}) (interface{}, error) {

		params, e := getDataBrowser(d)

		// get datas

		script := getSQLScript(SQLScript+"/databrowser_h3.sql", params)

		// tk.Printf("---\n%#v \n----\n", script)
		datas := []SPDataBrowser{}
		cursor, e := this.DB().Connection.NewQuery().
			Command("freequery", tk.M{}.Set("syntax", script)).
			Cursor(nil)
		e = cursor.Fetch(&datas, 0, true)

		cursor.Close()

		if e != nil && e.Error() == "No more data to fetched!" {
			e = nil
		}

		ret.Set("Datas", datas)

		// get total and summary

		total := k.Session(d.Hypoid+"Total", nil)
		summary := k.Session(d.Hypoid+"Summary", nil)

		// if total == nil || summary == nil || len(fieldsdouble) > 0 {
		summaryStr := " count(*) as Total"

		for _, val := range d.Fieldsdouble {
			strSum := ",(Select CAST((Sum(" + val + ")) as float)) as " + val + "sum"
			strAvg := ",(Select CAST((Avg(" + val + ")) as float)) as " + val + "avg"

			summaryStr += strSum
			summaryStr += strAvg
		}

		params.Set("@Summary", summaryStr)

		script = getSQLScript(SQLScript+"/databrowser_h3_summary.sql", params)
		// tk.Printf("---\n%#v \n----\n", script)
		resSum := []tk.M{}
		cursorTotal, e := this.DB().Connection.NewQuery().
			Command("freequery", tk.M{}.Set("syntax", script)).
			Cursor(nil)
		e = cursorTotal.Fetch(&resSum, 0, true)
		cursorTotal.Close()

		if e != nil && e.Error() == "No more data to fetched!" {
			e = nil
		}

		if len(resSum) > 0 {
			tmp := resSum[0]
			total = tmp.GetInt("total")
			tmpSummary := tk.M{}

			for _, val := range d.Fieldsdouble {
				tmpSummary.Set(val+"avg", tmp.GetFloat64(strings.ToLower(val+"avg")))
				tmpSummary.Set(val+"sum", tmp.GetFloat64(strings.ToLower(val+"sum")))
			}

			summary = []tk.M{tmpSummary}

			/*k.SetSession(d.Hypoid+"Total", total)
			k.SetSession(d.Hypoid+"Summary", summary)*/

			ret.Set("Total", total)
			ret.Set("Summary", summary)
		}
		/*} else {
			ret.Set("Total", total)
			ret.Set("Summary", summary)
		}*/
		return ret, e
	}, nil)

	return r
}

func (this *DataBrowserController) GetFilter(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	var e error

	r := new(tk.Result)
	d := DataBrowserInput{}
	f := tk.M{}
	ret := tk.M{}

	_ = k.GetForms(&f)
	_ = k.GetPayload(&d)

	// tk.Printf("%#v \n", f)

	r.Run(func(in interface{}) (interface{}, error) {

		activeField := f.GetString("active_field")

		if activeField != "" {
			params, err := getDataBrowser(d)
			// get datas
			params.Set("@GROUP", activeField)

			script := getSQLScript(SQLScript+"/databrowser_h3_filter.sql", params)

			// tk.Printf("---\n%#v \n----\n", script)
			cursor, err := this.DB().Connection.NewQuery().
				Command("freequery", tk.M{}.Set("syntax", script)).
				Cursor(nil)

			defer cursor.Close()

			// datas := []SPDataBrowser{}
			tmpDatas := []tk.M{}
			datas := []tk.M{}

			err = cursor.Fetch(&tmpDatas, 0, true)

			if e != nil && e.Error() == "No more data to fetched!" {
				e = nil
			}

			if len(tmpDatas) > 0 {
				for _, val := range tmpDatas {
					tmp := tk.M{}
					tmp.Set("_id", val.Get(strings.ToLower(activeField)))
					tmp.Set(activeField, val.Get(strings.ToLower(activeField)))

					datas = append(datas, tmp)
				}
			}

			ret.Set("Data", datas)

			e = err
		}

		return ret, e
	}, nil)

	return r
}

type SumList struct {
	field string
	tipe  string
}

func (this *DataBrowserController) SaveExcel(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	d := DataBrowserInput{}

	_ = k.GetPayload(&d)

	var (
		DisplaySumList []SumList
	)

	r := new(tk.Result)

	r.Run(func(in interface{}) (interface{}, error) {

		params, e := getDataBrowser(d)

		// ret.Set("Datas", datas)

		// get total and summary

		total := 0          //k.Session(d.Hypoid+"Total", nil)
		summary := []tk.M{} //k.Session(d.Hypoid+"Summary", nil)

		// if total == nil || summary == nil || len(fieldsdouble) > 0 {
		summaryStr := " count(*) as Total"

		for _, val := range d.Fieldsdouble {
			strSum := ",(Select CAST((Sum(" + val + ")) as float)) as " + val + "sum"
			strAvg := ",(Select CAST((Avg(" + val + ")) as float)) as " + val + "avg"

			summaryStr += strSum
			summaryStr += strAvg
		}

		params.Set("@Summary", summaryStr)

		script := getSQLScript(SQLScript+"/databrowser_h3_summary.sql", params)
		// tk.Printf("---\n%#v \n----\n", script)
		cursorTotal, e := this.DB().Connection.NewQuery().
			Command("freequery", tk.M{}.Set("syntax", script)).
			Cursor(nil)

		defer cursorTotal.Close()

		resSum := []tk.M{}

		e = cursorTotal.Fetch(&resSum, 0, true)

		if e != nil && e.Error() == "No more data to fetched!" {
			e = nil
		}

		if len(resSum) > 0 {
			tmp := resSum[0]
			total = tmp.GetInt("total")
			tmpSummary := tk.M{}

			for _, val := range d.Fieldsdouble {
				tmpSummary.Set(val+"avg", tmp.GetFloat64(strings.ToLower(val+"avg")))
				tmpSummary.Set(val+"sum", tmp.GetFloat64(strings.ToLower(val+"sum")))
			}

			summary = []tk.M{tmpSummary}
		}

		params.Set("@Offset", 0)
		params.Set("@Limit", total)

		script = getSQLScript(SQLScript+"/databrowser_h3.sql", params)

		// tk.Printf("---\n%#v \n----\n", script)
		cursor, e := this.DB().Connection.NewQuery().
			Command("freequery", tk.M{}.Set("syntax", script)).
			Cursor(nil)

		defer cursor.Close()

		// datas := []SPDataBrowser{}
		datas := make([]tk.M, 0)

		e = cursor.Fetch(&datas, 0, true)

		if e != nil && e.Error() == "No more data to fetched!" {
			e = nil
		}

		DisplayTypeCount := d.DisplayTypeCount
		DisplaySumList = []SumList{}
		for i := 0; i < DisplayTypeCount; i++ {
			sumData := SumList{}
			sumData.field = d.DisplayTypeList[i].Get("field").(string)
			sumData.tipe = d.DisplayTypeList[i].Get("type").(string)
			DisplaySumList = append(DisplaySumList, sumData)
		}

		excelFile, e := this.genExcelFile(d.HeaderList, d.Fields, datas, summary, DisplaySumList)
		return "../" + excelFile, e
	}, nil)

	tk.Printf("%#v \n", r)

	return r
}

func (this *DataBrowserController) genExcelFile(header []string, selectedColumn []string, datas []tk.M, dataSummary []tk.M, DisplaySumList []SumList) (string, error) {
	today := time.Now().UTC()
	fileName := "files/databrowser_" + today.Format("2006-01-02T150405") + ".xlsx"
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
			cell.SetValue(this.getExcelValue(data, field))
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

func (this *DataBrowserController) getExcelValue(data tk.M, field string) (result interface{}) {
	field = strings.ToLower(field)
	numberOfDot := strings.Count(field, ".")
	if numberOfDot > 0 {
		d := data.Get(field[0:strings.Index(field, ".")]).(tk.M)
		newField := field[strings.Index(field, ".")+1 : len(field)]
		result = this.getExcelValue(d, newField)
	} else {
		result = data.Get(field)
	}
	if result == nil {
		result = ""
	}
	return
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
		tk.Println(err.Error())
	}

	for idx, val := range params {
		script = strings.Replace(script, idx, tk.ToString(val), -1)
	}

	script = strings.Replace(script, "\t", "", -1)

	return
}
