package main

import (
	"os"
	"runtime"

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
	runtime.GOMAXPROCS(runtime.NumCPU())
	tk.Println("Starting the app..\n")

	sql, e := PrepareConnection()
	if e != nil {
		tk.Println(e)
	} else {
		base := new(BaseController)
		base.Ctx = orm.New(sql)
		defer base.Ctx.Close()

		/*new(GenSummaryData).Generate(base)
		new(GenMOR).Generate(base)
		new(GenPreventiveCorrectiveSummary).Generate(base)
		new(GenWODurationSummary).Generate(base)
		new(GenWOListSummary).Generate(base)

		new(GenPlantMaster).Generate(base)*/

		// new(GenDataBrowser).Generate(base)
		// new(REFunctionalLocation).Generate(base)
		//new(GenValueEquation).Generate(base)
		//new(GenValueEquation).Generate(base)
		new(GenValueEquationDashboard).Generate(base)
		// new(REFunctionalLocation).Generate(base)
		// new(REWOList).Generate(base)
		// new(RERPPCloseWO).Generate(base)
		// new(REFuelCargo).Generate(base)
		// new(REMaintenanceWO).Generate(base)
		// new(RESyntheticPM).Generate(base, 2014)
		// new(REPerfromanceFactors).Generate(base, "Qurayyah")
		// new(REConsolidated).Generate(base)
		// new(REVBMOperationalData).Generate(base) //Revenue
		// new(REFuelCost).Generate(base)
		// new(REPreventiveMaintenance).Generate(base)
		// new(REPlannedMaintenance).Generate(base)
		// new(REMaintenancePlan).Generate(base)
		// new(RENewObjectType).Generate(base)

	}

	tk.Println("Application Close..")
}
