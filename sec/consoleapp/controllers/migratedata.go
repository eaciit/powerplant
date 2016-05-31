package controllers

import (
	"reflect"
	"sync"
	"time"

	"github.com/eaciit/orm"

	_ "github.com/eaciit/dbox/dbc/mongo"
	_ "github.com/eaciit/dbox/dbc/mssql"
	. "github.com/eaciit/powerplant/sec/consoleapp/models"
	tk "github.com/eaciit/toolkit"
)

type MigrateData struct {
	*BaseController
}

func (m *MigrateData) DoDataBrowser() {

}

func (m *MigrateData) DoValueEquation() {

}

func (m *MigrateData) DoValueEquationDashboard() {

}

func (m *MigrateData) DoValueEquationDataQuality() {

}

func (m *MigrateData) DoCostSheet() error {
	tStart := time.Now()
	tk.Println("Starting DoCostSheet..")
	mod := new(CostSheet)

	c, e := m.BaseController.MongoCtx.Connection.NewQuery().From(mod.TableName()).Cursor(nil)

	if e != nil {
		return e
	}

	defer c.Close()

	result := []tk.M{}
	e = c.Fetch(&result, 0, false)

	for _, val := range result {
		_, e := m.InsertOut(val, new(CostSheet))
		if e != nil {
			tk.Printf("\n----------- ERROR -------------- \n %v \n\n %#v \n-------------------------  \n", e.Error(), val)
			return e
		}
		id := val.GetString("_id")
		details := val.Get("Details").(interface{}).([]interface{})
		tk.Printf("%#v \n\n", id)

		for _, detail := range details {
			det := detail.(tk.M)
			det.Set("CostSheet", id)

			_, e = m.InsertOut(det, new(CostSheetDetails))
			if e != nil {
				tk.Printf("\n----------- ERROR -------------- \n %v \n\n %#v \n-------------------------  \n", e.Error(), det)
				return e
			}
		}
	}

	cr, e := m.BaseController.SqlCtx.Connection.NewQuery().From(mod.TableName()).Cursor(nil)
	ctn := cr.Count()
	cr.Close()

	tk.Printf("Completed Success in %v | %v data(s)\n", time.Since(tStart), ctn)
	return nil
}

func (m *MigrateData) DoPowerPlantOutages() {

}

func (m *MigrateData) DoGeneralInfo() error {
	/*tStart := time.Now()
	tk.Println("Starting DoGeneralInfo..")
	mod := new(GeneralInfo)

	c, e := m.BaseController.MongoCtx.Connection.NewQuery().From(mod.TableName()).Cursor(nil)

	if e != nil {
		return e
	}

	defer c.Close()

	result := []tk.M{}
	e = c.Fetch(&result, 0, false)

	for _, val := range result {
		_, e := m.InsertOut(val, new(CostSheet))
		if e != nil {
			tk.Printf("\n----------- ERROR -------------- \n %v \n\n %#v \n-------------------------  \n", e.Error(), val)
			return e
		}
		id := val.GetString("_id")
		details := val.Get("Details").(interface{}).([]interface{})
		tk.Printf("%#v \n\n", id)

		for _, detail := range details {
			det := detail.(tk.M)
			det.Set("CostSheet", id)

			_, e = m.InsertOut(det, new(CostSheetDetails))
			if e != nil {
				tk.Printf("\n----------- ERROR -------------- \n %v \n\n %#v \n-------------------------  \n", e.Error(), det)
				return e
			}
		}
	}

	cr, e := m.BaseController.SqlCtx.Connection.NewQuery().From(mod.TableName()).Cursor(nil)
	ctn := cr.Count()
	cr.Close()

	tk.Printf("Completed Success in %v | %v data(s)\n", time.Since(tStart), ctn)*/
	return nil
}

func (m *MigrateData) InsertOut(in tk.M, mod orm.IModel) (out int64, e error) {
	muinsert := &sync.Mutex{}
	valueType := reflect.TypeOf(mod).Elem()
	for f := 0; f < valueType.NumField(); f++ {
		field := valueType.Field(f)
		bsonField := field.Tag.Get("bson")
		jsonField := field.Tag.Get("json")

		if jsonField != bsonField && field.Name != "RWMutex" && field.Name != "ModelBase" {
			in.Set(field.Name, GetMgoValue(in, bsonField))
		}
		switch field.Type.Name() {
		case "string":
			if GetMgoValue(in, bsonField) == nil {
				in.Set(field.Name, "")
			}
			break
		case "Time":
			if GetMgoValue(in, bsonField) == nil {
				in.Set(field.Name, time.Time{})
			} else {
				in.Set(field.Name, GetMgoValue(in, bsonField).(time.Time).UTC())
			}
			break
		default:
			break
		}

	}

	e = tk.Serde(in, mod, "json")

	if e != nil {
		return
	}

	muinsert.Lock()
	out, e = m.BaseController.SqlCtx.InsertOut(mod)
	muinsert.Unlock()
	return
}
