package controllers

import (
	"github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/mongo"
	_ "github.com/eaciit/dbox/dbc/mssql"
	"github.com/eaciit/orm"
	. "github.com/eaciit/powerplant/sec/consoleapp/models"
	tk "github.com/eaciit/toolkit"
	"gopkg.in/mgo.v2/bson"
	"reflect"
	"sync"
	"time"
)

type MigrateData struct {
	*BaseController
}

func (m *MigrateData) DoDataBrowser() {

}

func (m *MigrateData) DoValueEquation() error {
	tStart := time.Now()
	tk.Println("Starting DoValueEquation..")
	mod := new(ValueEquation)

	c, e := m.BaseController.MongoCtx.Connection.NewQuery().From(mod.TableName()).Cursor(nil)
	defer c.Close()
	if e != nil {
		return e
	}

	result := []tk.M{}
	e = c.Fetch(&result, 0, false)

	for _, val := range result {
		fuels := val.Get("Fuel").(interface{}).([]interface{})
		val.Set("Fuel", nil)
		details := val.Get("Detail").(interface{}).([]interface{})
		val.Set("Detail", nil)
		top10s := val.Get("Top10").(interface{}).([]interface{})
		val.Set("Top10", nil)

		sid := val["_id"]
		id := sid.(bson.ObjectId).Hex()
		val.Set("_id", id)

		_, e := m.InsertOut(val, new(ValueEquation))
		if e != nil {
			tk.Printf("\n----------- ERROR -------------- \n %v \n\n %#v \n-------------------------  \n", e.Error(), val)
			return e
		}

		periods := val.Get("Period").(tk.M)
		periods.Set("Id", id)
		_, e = m.InsertOut(periods, new(ValueEquationPeriod))
		if e != nil {
			tk.Printf("\n----------- ERROR -------------- \n %v \n\n %#v \n-------------------------  \n", e.Error(), periods)
			return e
		}

		for _, fuel := range fuels {
			f := fuel.(tk.M)
			f.Set("VEId", id)
			_, e = m.InsertOut(f, new(ValueEquationFuel))
			if e != nil {
				tk.Printf("\n----------- ERROR -------------- \n %v \n\n %#v \n-------------------------  \n", e.Error(), f)
				return e
			}
		}

		for _, detail := range details {
			d := detail.(tk.M)
			d.Set("VEId", id)
			_, e = m.InsertOut(d, new(ValueEquationDetails))
			if e != nil {
				tk.Printf("\n----------- ERROR -------------- \n %v \n\n %#v \n-------------------------  \n", e.Error(), d)
				return e
			}
		}

		for _, top10 := range top10s {
			t := top10.(tk.M)
			t.Set("VEId", id)
			_, e = m.InsertOut(t, new(ValueEquationTop10))
			tk.Println(t.Get("VEId"))
			if e != nil {
				tk.Printf("\n----------- ERROR -------------- \n %v \n\n %#v \n-------------------------  \n", e.Error(), t)
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

func (m *MigrateData) DoValueEquationDashboard() {

}

func (m *MigrateData) DoValueEquationDataQuality() error {
	tStart := time.Now()
	tk.Println("Starting ValueEquationDataQuality..")
	mod := new(ValueEquationDataQuality)

	c, e := m.BaseController.MongoCtx.Connection.NewQuery().From(mod.TableName()).Cursor(nil)

	if e != nil {
		return e
	}

	defer c.Close()

	result := []tk.M{}
	e = c.Fetch(&result, 0, false)

	for _, val := range result {

		// _, e := m.InsertOut(val, new(ValueEquationDataQuality))
		// if e != nil {
		// 	tk.Printf("\n----------- ERROR -------------- \n %v \n\n %#v \n-------------------------  \n", e.Error(), val)
		// 	return e
		// }
		filter := []*dbox.Filter{}
		filter = append(filter, dbox.Eq("Dates", GetMgoValue(val, "Period.Dates").(time.Time).UTC()))
		filter = append(filter, dbox.Eq("Plant", GetMgoValue(val, "Plant")))
		filter = append(filter, dbox.Eq("Unit", GetMgoValue(val, "Unit")))
		// tk.M{}.Set("where", dbox.And(filter...)
		csr, e := m.SqlCtx.Connection.NewQuery().From(mod.TableName()).Where(filter...).Cursor(nil)
		csr.Close()
		if e != nil {
			return e
		}
		result := []tk.M{}
		e = csr.Fetch(&result, 0, false)
		if e != nil {
			return e
		}
		if len(result) > 0 {
			Id := result[0].Get("id")

			CapacityPaymentDocuments := val.Get("CapacityPaymentDocuments").([]interface{})
			for _, x := range CapacityPaymentDocuments {
				doc := x.(tk.M).Set("VEId", Id)
				_, e = m.InsertOut(doc, new(VEDQCapacityPaymentDocuments))
				if e != nil {
					tk.Printf("\n----------- ERROR -------------- \n %v \n\n %#v \n-------------------------  \n", e.Error(), val)
					return e
				}
			}

			EnergyPaymentDocuments := val.Get("EnergyPaymentDocuments").([]interface{})
			for _, x := range EnergyPaymentDocuments {
				doc := x.(tk.M).Set("VEId", Id)
				_, e = m.InsertOut(doc, new(VEDQEnergyPaymentDocuments))
				if e != nil {
					tk.Printf("\n----------- ERROR -------------- \n %v \n\n %#v \n-------------------------  \n", e.Error(), val)
					return e
				}
			}

			StartupPaymentDocuments := val.Get("StartupPaymentDocuments").([]interface{})
			for _, x := range StartupPaymentDocuments {
				doc := x.(tk.M).Set("VEId", Id)
				_, e = m.InsertOut(doc, new(VEDQStartupPaymentDocuments))
				if e != nil {
					tk.Printf("\n----------- ERROR -------------- \n %v \n\n %#v \n-------------------------  \n", e.Error(), val)
					return e
				}
			}

			PenaltyDocuments := val.Get("PenaltyDocuments").([]interface{})
			for _, x := range PenaltyDocuments {
				doc := x.(tk.M).Set("VEId", Id)
				_, e = m.InsertOut(doc, new(VEDQPenaltyDocuments))
				if e != nil {
					tk.Printf("\n----------- ERROR -------------- \n %v \n\n %#v \n-------------------------  \n", e.Error(), val)
					return e
				}
			}

			PrimaryFuel1stDocuments := val.Get("PrimaryFuel1stDocuments").([]interface{})
			for _, x := range PrimaryFuel1stDocuments {
				doc := x.(tk.M).Set("VEId", Id)
				_, e = m.InsertOut(doc, new(VEDQPrimaryFuel1stDocuments))
				if e != nil {
					tk.Printf("\n----------- ERROR -------------- \n %v \n\n %#v \n-------------------------  \n", e.Error(), val)
					return e
				}
			}

			PrimaryFuel2ndDocuments := val.Get("PrimaryFuel2ndDocuments").([]interface{})
			for _, x := range PrimaryFuel2ndDocuments {
				doc := x.(tk.M).Set("VEId", Id)
				_, e = m.InsertOut(doc, new(VEDQPrimaryFuel2ndDocuments))
				if e != nil {
					tk.Printf("\n----------- ERROR -------------- \n %v \n\n %#v \n-------------------------  \n", e.Error(), val)
					return e
				}
			}

			BackupFuelDocuments := val.Get("BackupFuelDocuments").([]interface{})
			for _, x := range BackupFuelDocuments {
				doc := x.(tk.M).Set("VEId", Id)
				_, e = m.InsertOut(doc, new(VEDQBackupFuelDocuments))
				if e != nil {
					tk.Printf("\n----------- ERROR -------------- \n %v \n\n %#v \n-------------------------  \n", e.Error(), val)
					return e
				}
			}

			MaintenanceCostDocuments := val.Get("MaintenanceCostDocuments")
			if MaintenanceCostDocuments != nil {
				for _, x := range MaintenanceCostDocuments.([]interface{}) {
					doc := x.(tk.M).Set("VEId", Id)
					_, e = m.InsertOut(doc, new(VEDQMaintenanceCostDocuments))
					tk.Println(doc)
					if e != nil {
						tk.Printf("\n----------- ERROR -------------- \n %v \n\n %#v \n-------------------------  \n", e.Error(), val)
						return e
					}
				}
			}

			MaintenanceDurationDocuments := val.Get("MaintenanceDurationDocuments")
			if MaintenanceDurationDocuments != nil {
				for _, x := range MaintenanceDurationDocuments.([]interface{}) {
					doc := x.(tk.M).Set("VEId", Id)
					_, e = m.InsertOut(doc, new(VEDQMaintenanceDurationDocuments))
					if e != nil {
						tk.Printf("\n----------- ERROR -------------- \n %v \n\n %#v \n-------------------------  \n", e.Error(), val)
						return e
					}
				}
			}

		}

	}

	cr, e := m.BaseController.SqlCtx.Connection.NewQuery().From(mod.TableName()).Cursor(nil)
	ctn := cr.Count()
	cr.Close()

	tk.Printf("Completed Success in %v | %v data(s)\n", time.Since(tStart), ctn)
	return nil
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

func (m *MigrateData) DoPowerPlantOutages() error {
	tStart := time.Now()
	tk.Println("Starting DoPowerPlantOutages..")
	mod := new(PowerPlantOutages)

	c, e := m.BaseController.MongoCtx.Connection.NewQuery().From(mod.TableName()).Cursor(nil)

	if e != nil {
		return e
	}

	defer c.Close()

	result := []tk.M{}
	e = c.Fetch(&result, 0, false)

	for _, val := range result {
		details := val.Get("Details").(interface{}).([]interface{})
		val.Set("Details", nil)

		_, e := m.InsertOut(val, new(PowerPlantOutages))
		if e != nil {
			tk.Printf("\n----------- ERROR -------------- \n %v \n\n %#v \n-------------------------  \n", e.Error(), val)
			return e
		}
		id := val.GetString("_id")
		tk.Printf("%#v \n\n", id)

		for _, detail := range details {
			det := detail.(tk.M)
			det.Set("POId", id)
			tk.Println(det)
			_, e = m.InsertOut(det, new(PowerPlantOutagesDetails))
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

func (m *MigrateData) DoGeneralInfo() error {
	tStart := time.Now()
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
		_, e := m.InsertOut(val, new(GeneralInfo))
		if e != nil {
			tk.Printf("\n----------- ERROR -------------- \n %v \n\n %#v \n-------------------------  \n", e.Error(), val)
			return e
		}
		id := val.GetString("_id")
		tk.Printf("%#v \n\n", id)

		if nil != val.Get("InstalledMWh") {
			installedMWH := val.Get("InstalledMWh").(interface{}).([]interface{})
			for _, detail := range installedMWH {
				det := detail.(tk.M)
				det.Set("GenID", id)
				det.Set("Type", "InstalledMWh")
				_, e = m.InsertOut(det, new(GeneralInfoDetails))
				if e != nil {
					tk.Printf("\n----------- ERROR -------------- \n %v \n\n %#v \n-------------------------  \n", e.Error(), det)
					return e
				}
			}
		}

		if nil != val.Get("ActualEnergyGeneration") {
			actualEnergyGeneration := val.Get("ActualEnergyGeneration").(interface{}).([]interface{})
			for _, detail := range actualEnergyGeneration {
				det := detail.(tk.M)
				det.Set("GenID", id)
				det.Set("Type", "ActualEnergyGeneration")
				_, e = m.InsertOut(det, new(GeneralInfoDetails))
				if e != nil {
					tk.Printf("\n----------- ERROR -------------- \n %v \n\n %#v \n-------------------------  \n", e.Error(), det)
					return e
				}
			}
		}

		if nil != val.Get("ActualFuelConsumption") {
			actualFuelConsumption := val.Get("ActualFuelConsumption").(interface{}).([]interface{})
			for _, detail := range actualFuelConsumption {
				det := detail.(tk.M)
				det.Set("GenID", id)
				_, e = m.InsertOut(det, new(GeneralInfoActualFuelConsumption))
				if e != nil {
					tk.Printf("\n----------- ERROR -------------- \n %v \n\n %#v \n-------------------------  \n", e.Error(), det)
					return e
				}
			}
		}

		if nil != val.Get("CapacityFactor") {
			capacityFactor := val.Get("CapacityFactor").(interface{}).([]interface{})
			for _, detail := range capacityFactor {
				det := detail.(tk.M)
				det.Set("GenID", id)
				det.Set("Type", "CapacityFactor")
				_, e = m.InsertOut(det, new(GeneralInfoDetails))
				if e != nil {
					tk.Printf("\n----------- ERROR -------------- \n %v \n\n %#v \n-------------------------  \n", e.Error(), det)
					return e
				}
			}
		}
	}

	cr, e := m.BaseController.SqlCtx.Connection.NewQuery().From(mod.TableName()).Cursor(nil)
	ctn := cr.Count()
	cr.Close()

	tk.Printf("Completed Success in %v | %v data(s)\n", time.Since(tStart), ctn)
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
