package main

import (
	"os"

	_ "github.com/eaciit/dbox/dbc/mssql"
	"github.com/eaciit/orm"
	. "github.com/eaciit/powerplant/sec/consoleapp/generator/controllers"
	tk "github.com/eaciit/toolkit"
)

var (
	wd = func() string {
		d, _ := os.Getwd()
		return d + "/"
	}()
)

func main() {
	tk.Println("Starting the app..\n")

	sql, e := PrepareConnection()
	if e != nil {
		tk.Println(e)
	} else {
		base := new(BaseController)
		base.Ctx = orm.New(sql)
		defer base.Ctx.Close()

		// new(GenPlantMaster).Generate(base)
		// new(GenMOR).Generate(base)
		// new(GenSummaryData).Generate(base)
		// new(GenPreventiveCorrectiveSummary).Generate(base)
	}

	base := new(BaseController)
	base.Ctx = orm.New(sql)
	defer base.Ctx.Close()

	// Generate DataMaster
	//Mst := DataMaster{base}
	//Mst.Generate()

	// Generate MOR
	// MOR := DataMOR{base}
	// MOR.Generate()

	// Generate Summary
	Summary := GenSummaryData{base}
	Summary.Generate(base)
	// Generate DurationSummary
	/*Duration := DurationIntervalSummary{base}
	Duration.Generate()*/
	VE := ValueEquationGenerator{base}
	VE.Generate()

	tk.Println("Application closed..")
}
