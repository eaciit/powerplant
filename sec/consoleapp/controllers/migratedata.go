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
	"strconv"
	"sync"
	"time"
)

type MigrateData struct {
	*BaseController
}

func (m *MigrateData) DoDataBrowser() {
	/*tStart := time.Now()
	tk.Println("Starting DoDataBrowser..")
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
	return nil*/
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

		periods := val.Get("Period").(tk.M)
		val.Set("PeriodYear", periods["Year"])
		val.Set("PeriodMonth", periods["Month"])
		val.Set("PeriodDates", periods["Dates"])
		tk.Printf("\n----------- RES -------------- \n %v \n\n %#v \n-------------------------  \n", val)

		for {
			_, e := m.InsertOut(val, new(ValueEquation))
			if e == nil {
				break
			} else {
				m.SqlCtx.Connection.Connect()
			}
		}

		for _, fuel := range fuels {
			f := fuel.(tk.M)
			f.Set("VEId", id)
			for {
				_, e = m.InsertOut(f, new(ValueEquationFuel))
				if e == nil {
					break
				} else {
					m.SqlCtx.Connection.Connect()
				}
			}
		}

		for _, detail := range details {
			d := detail.(tk.M)
			d.Set("VEId", id)
			for {
				_, e = m.InsertOut(d, new(ValueEquationDetails))
				if e == nil {
					break
				} else {
					m.SqlCtx.Connection.Connect()
				}
			}
		}

		for _, top10 := range top10s {
			t := top10.(tk.M)
			t.Set("VEId", id)
			for {
				_, e = m.InsertOut(t, new(ValueEquationTop10))
				if e == nil {
					break
				} else {
					m.SqlCtx.Connection.Connect()
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

func (m *MigrateData) DoValueEquationDashboard() error {
	tStart := time.Now()
	tk.Println("Starting DoValueEquationDashboard..")
	mod := new(ValueEquationDashboard)

	c, e := m.BaseController.MongoCtx.Connection.NewQuery().From(mod.TableName()).Cursor(nil)

	if e != nil {
		return e
	}

	defer c.Close()

	result := []tk.M{}
	e = c.Fetch(&result, 0, false)
	for _, val := range result {

		// _, e := m.InsertOut(val, new(ValueEquationDashboard))
		// if e != nil {
		// 	tk.Printf("\n----------- ERROR -------------- \n %v \n\n %#v \n-------------------------  \n", e.Error(), val)
		// 	return e
		// }
		filter := []*dbox.Filter{}
		filter = append(filter, dbox.Eq("Dates", GetMgoValue(val, "Period.Dates").(time.Time).UTC()))
		filter = append(filter, dbox.Eq("Plant", GetMgoValue(val, "Plant")))
		filter = append(filter, dbox.Eq("Unit", GetMgoValue(val, "Unit")))

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

			// Fuel := val.Get("Fuel").([]interface{})
			// for _, x := range Fuel {
			// 	doc := x.(tk.M).Set("VEId", Id)
			// 	_, e = m.InsertOut(doc, new(VEDFuel))
			// 	if e != nil {
			// 		tk.Printf("\n----------- ERROR -------------- \n %v \n\n %#v \n-------------------------  \n", e.Error(), val)
			// 		return e
			// 	}
			// }

			// Detail := val.Get("Detail")
			// if Detail != nil {
			// 	for _, x := range Detail.([]interface{}) {
			// 		doc := x.(tk.M).Set("VEId", Id)
			// 		_, e = m.InsertOut(doc, new(VEDDetail))
			// 		if e != nil {
			// 			tk.Printf("\n----------- ERROR -------------- \n %v \n\n %#v \n-------------------------  \n", e.Error(), val)
			// 			return e
			// 		}
			// 	}
			// }

			Top10 := val.Get("Top10")
			if Top10 != nil {
				for _, x := range Top10.([]interface{}) {
					doc := x.(tk.M).Set("VEId", Id)
					for index := 0; index < retry; index++ {
						_, e = m.InsertOut(doc, new(VEDTop10))
						if e == nil {
							break
						} else {
							m.SqlCtx.Connection.Connect()
						}
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

func (m *MigrateData) DoGenerateVibration() error {
	tStart := time.Now()
	tk.Println("Starting DoGenerateVibration..")
	mod := new(Vibration)

	c, e := m.BaseController.MongoCtx.Connection.NewQuery().From(mod.TableName()).Cursor(nil)

	if e != nil {
		return e
	}

	defer c.Close()

	result := []tk.M{}
	e = c.Fetch(&result, 0, false)

	for _, val := range result {
		for {
			_, e := m.InsertOut(val, new(Vibration))
			if e == nil {
				break
			} else {
				m.SqlCtx.Connection.Connect()
			}
		}

		id := val.GetString("_id")

		for idx := 0; idx < 6; idx++ {
			idConv := strconv.Itoa(idx)

			bearingString := "Bearing" + idConv

			bear := val.Get(bearingString)
			if bear != nil {
				bear1 := bear.(interface{}).([]interface{})

				for _, bearing := range bear1 {
					bea := bearing.(tk.M)

					bea.Set("VibrationId", id)
					bea.Set("Type", bearingString)
					for {
						_, e = m.InsertOut(bea, new(BearingDetail))
						if e == nil {
							break
						} else {
							m.SqlCtx.Connection.Connect()
						}
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

func (m *MigrateData) DoScenarioSimulation() error {
	tStart := time.Now()
	tk.Println("Starting DoScenarioSimulation..")
	mod := new(ScenarioSimulation)

	c, e := m.BaseController.MongoCtx.Connection.NewQuery().From(mod.TableName()).Cursor(nil)

	if e != nil {
		return e
	}

	defer c.Close()

	result := []tk.M{}
	e = c.Fetch(&result, 0, false)

	for _, val := range result {
		plants := val.Get("SelectedPlant").(interface{}).([]interface{})
		val.Set("SelectedPlant", nil)
		units := val.Get("SelectedUnit").(interface{}).([]interface{})
		val.Set("SelectedUnit", nil)
		scenarios := val.Get("SelectedScenario").(interface{}).([]interface{})
		val.Set("SelectedScenario", nil)

		sid := val["_id"]
		id := sid.(bson.ObjectId).Hex()
		val.Set("_id", id)

		historicresult := val.Get("HistoricResult").(tk.M)
		val.Set("HistoricResultRevenue", historicresult["Revenue"])
		val.Set("HistoricResultLaborCost", historicresult["LaborCost"])
		val.Set("HistoricResultMaterialCost", historicresult["MaterialCost"])
		val.Set("HistoricResultServiceCost", historicresult["ServiceCost"])
		val.Set("HistoricResultOperatingCost", historicresult["OperatingCost"])
		val.Set("HistoricResultMaintenanceCost", historicresult["MaintenanceCost"])
		val.Set("HistoricResultValueEquation", historicresult["ValueEquation"])

		futureresult := val.Get("FutureResult").(tk.M)
		val.Set("FutureResultRevenue", futureresult["Revenue"])
		val.Set("FutureResultLaborCost", futureresult["LaborCost"])
		val.Set("FutureResultMaterialCost", futureresult["MaterialCost"])
		val.Set("FutureResultServiceCost", futureresult["ServiceCost"])
		val.Set("FutureResultOperatingCost", futureresult["OperatingCost"])
		val.Set("FutureResultMaintenanceCost", futureresult["MaintenanceCost"])
		val.Set("FutureResultValueEquation", futureresult["ValueEquation"])

		differential := val.Get("Differential").(tk.M)
		val.Set("DifferentialRevenue", differential["Revenue"])
		val.Set("DifferentialLaborCost", differential["LaborCost"])
		val.Set("DifferentialMaterialCost", differential["MaterialCost"])
		val.Set("DifferentialServiceCost", differential["ServiceCost"])
		val.Set("DifferentialOperatingCost", differential["OperatingCost"])
		val.Set("DifferentialMaintenanceCost", differential["MaintenanceCost"])
		val.Set("DifferentialValueEquation", differential["ValueEquation"])

		_, e := m.InsertOut(val, new(ScenarioSimulation))
		if e != nil {
			tk.Println(e.Error())
		}

		for _, plant := range plants {
			p := tk.M{}
			p.Set("SSId", id)
			p.Set("Plant", plant)
			_, e = m.InsertOut(p, new(ScenarioSimulationSelectedPlant))
			if e != nil {
				tk.Println(e.Error())
			}
		}

		for _, unit := range units {
			u := tk.M{}
			u.Set("SSId", id)
			u.Set("Unit", unit)
			_, e = m.InsertOut(u, new(ScenarioSimulationSelectedUnit))
			if e != nil {
				tk.Println(e.Error())
			}
		}

		for _, scenario := range scenarios {
			s := scenario.(tk.M)
			s.Set("SSId", id)
			_, e = m.InsertOut(s, new(ScenarioSimulationSelectedScenario))
			if e != nil {
				tk.Println(e.Error())
			}
		}

	}

	cr, e := m.BaseController.SqlCtx.Connection.NewQuery().From(mod.TableName()).Cursor(nil)
	ctn := cr.Count()
	cr.Close()

	tk.Printf("Completed Success in %v | %v data(s)\n", time.Since(tStart), ctn)
	return nil
}

func (m *MigrateData) DoGenerateAssetClass() error {
	tStart := time.Now()
	tk.Println("Starting DoGenerateAssetClass..")
	mod := new(SampleAssetClass)

	c, e := m.BaseController.MongoCtx.Connection.NewQuery().From(mod.TableName()).Cursor(nil)

	if e != nil {
		return e
	}

	defer c.Close()

	result := []tk.M{}
	e = c.Fetch(&result, 0, false)

	for _, val := range result {
		for {
			_, e := m.InsertOut(val, new(SampleAssetClass))
			if e == nil {
				break
			} else {
				m.SqlCtx.Connection.Connect()
			}
		}
	}

	cr, e := m.BaseController.SqlCtx.Connection.NewQuery().From(mod.TableName()).Cursor(nil)
	ctn := cr.Count()
	cr.Close()

	tk.Printf("Completed Success in %v | %v data(s)\n", time.Since(tStart), ctn)
	return nil
}

func (m *MigrateData) DoGenerateAssetType() error {
	tStart := time.Now()
	tk.Println("Starting DoGenerateAssetType..")
	mod := new(SampleAssetType)

	c, e := m.BaseController.MongoCtx.Connection.NewQuery().From(mod.TableName()).Cursor(nil)

	if e != nil {
		return e
	}

	defer c.Close()

	result := []tk.M{}
	e = c.Fetch(&result, 0, false)

	for _, val := range result {
		for {
			_, e := m.InsertOut(val, new(SampleAssetType))
			if e == nil {
				break
			} else {
				m.SqlCtx.Connection.Connect()
			}
		}
	}

	cr, e := m.BaseController.SqlCtx.Connection.NewQuery().From(mod.TableName()).Cursor(nil)
	ctn := cr.Count()
	cr.Close()

	tk.Printf("Completed Success in %v | %v data(s)\n", time.Since(tStart), ctn)
	return nil
}

func (m *MigrateData) DoGenerateAssetLevel() error {
	tStart := time.Now()
	tk.Println("Starting DoGenerateAssetLevel..")
	mod := new(SampleAssetLevel)

	c, e := m.BaseController.MongoCtx.Connection.NewQuery().From(mod.TableName()).Cursor(nil)

	if e != nil {
		return e
	}

	defer c.Close()

	result := []tk.M{}
	e = c.Fetch(&result, 0, false)

	for _, val := range result {
		for {
			_, e := m.InsertOut(val, new(SampleAssetLevel))
			if e == nil {
				break
			} else {
				m.SqlCtx.Connection.Connect()
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

	return out, e
}
