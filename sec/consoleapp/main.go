package main

import (
	"os"

	_ "github.com/eaciit/dbox/dbc/mongo"
	_ "github.com/eaciit/dbox/dbc/mssql"
	"github.com/eaciit/orm"
	. "github.com/eaciit/powerplant/sec/consoleapp/controllers"
	// . "github.com/eaciit/powerplant/sec/consoleapp/models"
	tk "github.com/eaciit/toolkit"
)

var (
	wd = func() string {
		d, _ := os.Getwd()
		return d + "/"
	}()
)

func main() {
	tk.Println("Starting the app..")
	mongo, e := PrepareConnection("mongo")
	if e != nil {
		tk.Println(e)
	}

	sql, e := PrepareConnection("mssql")
	if e != nil {
		tk.Println(e)
	}

	base := new(BaseController)
	base.MongoCtx = orm.New(mongo)
	base.SqlCtx = orm.New(sql)

	defer base.MongoCtx.Close()
	defer base.SqlCtx.Close()
	// convert(new(MasterUnitNoTurbineParent), base) //		done
	// convert(new(MasterFailureCode), base) //				done
	// convert(new(WOList), base) // 						done
	// convert(new(AnomaliesWOList), base) //				done
	// convert(new(Availability), base) //					done
	// convert(new(Consolidated), base) // 					done
	// convert(new(FuelCost), base) //						done
	// convert(new(FuelTransport), base) // 				done
	// convert(new(FunctionalLocation), base) // 			done
	// convert(new(AnomaliesFunctionalLocation), base) //	done
	// convert(new(GenerationAppendix), base) //			done
	// convert(new(MaintenanceCost), base) //				done
	// convert(new(MaintenanceCostFL), base) //				done
	// convert(new(MaintenanceCostByHour), base) //			done
	// convert(new(MaintenancePlan), base) // 				done
	// convert(new(MaintenanceWorkOrder), base) // 			done
	// convert(new(MappedEquipmentType), base) // 			done
	// convert(new(MasterEquipmentType), base) // 			done
	// convert(new(MasterMROElement), base) // 				done
	// convert(new(MasterOrderType), base) // 				done
	// convert(new(MasterPlant), base) // 					done
	// convert(new(NewEquipmentType), base) // 				done
	// convert(new(NotificationFailure), base) // 			done
	// convert(new(NotificationFailureNoYear), base) // 	done
	// convert(new(OperationalData), base) // 				done
	// convert(new(PerformanceFactors), base) // 			done
	// convert(new(PlannedMaintenance), base) // 			done
	// convert(new(PowerPlantCoordinates), base) // 		done
	// convert(new(PowerPlantInfo),base) // 				done
	// convert(new(PrevMaintenanceValueEquation), base) // 	done
	// convert(new(RPPCloseWO), base) // 					done
	// convert(new(StartupPaymentAndPenalty), base) // 		done
	// convert(new(SyntheticPM), base) // 					done
	// convert(new(UnitCost), base) // 						done
	// convert(new(Vibration), base) // 					done

	// convert(new(MORSummary), base) // 					done
	// convert(new(MORCalculationFlatSummary), base)   // 	done
	// convert(new(PreventiveCorrectiveSummary), base) //	done
	// convert(new(WODurationSummary), base) // 			done
	// convert(new(WOListSummary), base) // 				done

	// convert(new(FailureAfterPreventiveSummary), base)// 	done
	// convert(new(RegenMasterPlant), base) // 				done
	// convert(new(MasterFailureCode), base) // 			done
	// convert(new(MasterUnitNoTurbineParent), base) // 	done
	// convert(new(DataTempMaintenance), base) // 			done

	// convert(new(SummaryData), base) // 					done
	// convert(new(DataBrowser), base)
}

func convert(m orm.IModel, base *BaseController) {
	e := base.ConvertMGOToSQLServer(m)
	if e != nil {
		tk.Printf("\nERROR: %v \n", e.Error())
	}
}
