package controllers

import (
	"github.com/eaciit/crowd"
	"github.com/eaciit/dbox"
	//"github.com/eaciit/orm"
	//. "github.com/eaciit/powerplant/sec/consoleapp/generator/helpers"*/
	. "github.com/eaciit/powerplant/sec/library/models"
	tk "github.com/eaciit/toolkit"
	"strconv"
	/*  "strings"
	"time"*/
	"log"
)

type ValueEquationGenerator struct {
	*BaseController
}

func (v *ValueEquationGenerator) Generate() {
	tk.Println("##Generating Value equation Data..")
	e := v.GenerateValueEquation()
	if e != nil {
		tk.Println(e)
	}
	tk.Println("##Value Equation Data : DONE\n")
}

func (v *ValueEquationGenerator) GenerateValueEquation() error {
	var e error
	Year := 2014
	YearFirst := strconv.Itoa(Year) + "-01-01 00:00:00.000"
	YearLast := strconv.Itoa(Year+1) + "-01-01 00:00:00.000"

	Plant := "Qurayyah CC"

	c := v.Ctx.Connection
	query := []*dbox.Filter{}

	query = append(query, dbox.Eq("Plant", Plant))
	csr, e := c.NewQuery().From(new(PerformanceFactors).TableName()).Where(query...).Cursor(nil)

	if e != nil {
		return e
	} else {
		defer csr.Close()
	}

	pfs := []tk.M{}
	e = csr.Fetch(&pfs, 0, false)

	csr1, e := c.NewQuery().From(new(Consolidated).TableName()).Where(query...).Cursor(nil)

	if e != nil {
		return e
	} else {
		defer csr1.Close()
	}

	cons := []tk.M{}
	e = csr1.Fetch(&cons, 0, false)

	query = append(query, dbox.And(dbox.Gte("DatePerformed", YearFirst), dbox.Lt("DatePerformed", YearLast)))
	csr2, e := c.NewQuery().From(new(PrevMaintenanceValueEquation).TableName()).Where(query...).Cursor(nil)

	if e != nil {
		return e
	} else {
		defer csr2.Close()
	}

	lists := []tk.M{}
	e = csr2.Fetch(&lists, 0, false)

	query = nil
	query = append(query, dbox.Eq("Plant", Plant))
	query = append(query, dbox.Eq("Year", Year))
	csr3, e := c.NewQuery().From(new(PowerPlantOutages).TableName()).Where(query...).Cursor(nil)

	if e != nil {
		return e
	} else {
		defer csr3.Close()
	}

	outages := []tk.M{}
	e = csr3.Fetch(&outages, 0, false)

	csr4, e := c.NewQuery().From(new(StartupPaymentAndPenalty).TableName()).Where(query...).Cursor(nil)

	if e != nil {
		return e
	} else {
		defer csr4.Close()
	}

	start := []tk.M{}
	e = csr4.Fetch(&start, 0, false)

	csr5, e := c.NewQuery().From(new(FuelCost).TableName()).Where(query...).Cursor(nil)

	if e != nil {
		return e
	} else {
		defer csr5.Close()
	}

	fuelcosts := []tk.M{}
	e = csr5.Fetch(&fuelcosts, 0, false)

	query = nil
	query = append(query, dbox.Eq("Plant", Plant))
	query = append(query, dbox.And(dbox.Gte("ScheduledStart", YearFirst), dbox.Lt("ScheduledStart", YearLast)))
	csr6, e := c.NewQuery().From(new(SyntheticPM).TableName()).Where(query...).Cursor(nil)

	if e != nil {
		return e
	} else {
		defer csr6.Close()
	}

	syn := []tk.M{}
	e = csr6.Fetch(&syn, 0, false)

	query = nil
	query = append(query, dbox.Eq("Plant", Plant))
	query = append(query, dbox.Eq("Year", Year))
	csr7, e := c.NewQuery().From(new(FuelTransport).TableName()).Where(query...).Cursor(nil)

	if e != nil {
		return e
	} else {
		defer csr7.Close()
	}

	trans := []tk.M{}
	e = csr7.Fetch(&trans, 0, false)

	sintax := "select * from DataBrowser inner join PowerPlantCoordinates on DataBrowser.PlantCode = PowerPlantCoordinates.PlantCode where PeriodYear = " + strconv.Itoa(Year) + " and PowerPlantCoordinates.PlantName = '" + Plant + "'"
	csr8, e := c.NewQuery().Command("freequery", tk.M{}.Set("syntax", sintax)).Cursor(nil)

	if e != nil {
		return e
	} else {
		defer csr8.Close()
	}

	databrowser := []tk.M{}
	e = csr8.Fetch(&databrowser, 0, false)

	csr9, e := c.NewQuery().From(new(GenerationAppendix).TableName()).Cursor(nil)

	if e != nil {
		return e
	} else {
		defer csr9.Close()
	}

	genA := []tk.M{}
	e = csr9.Fetch(&genA, 0, false)

	Units := crowd.From(&pfs).Group(func(x interface{}) interface{} {
		return x.(tk.M).GetString("unit")
	}, nil).Exec().Result.Data().([]crowd.KV)

	log.Println(Units)
	//DieselConsumptions :=
	//double DieselConsumptions = fuelcosts.Where(x => x.PrimaryFuelType == "DIESEL").Sum(x => x.PrimaryFuelConsumed) + fuelcosts.Where(x => x.BackupFuelType == "DIESEL").Sum(x => x.BackupFuelConsumed) * 1000;
	return e
}
