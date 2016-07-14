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

		// new(GenSummaryData).Generate(base)
		// new(GenMOR).Generate(base)
		// new(GenPreventiveCorrectiveSummary).Generate(base)
		// new(GenWODurationSummary).Generate(base)
		// new(GenWOListSummary).Generate(base)

		// new(GenPlantMaster).Generate(base)

		// new(GenValueEquation).Generate(base)
		new(REFunctionalLocation).Generate(base)

	}

	tk.Println("Application Close..")
}
