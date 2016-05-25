package main

import (
	"bufio"
	"os"
	"runtime"
	"strings"

	"github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/mongo"
	_ "github.com/eaciit/dbox/dbc/mssql"
	"github.com/eaciit/orm"
	. "github.com/eaciit/powerplant/sec/consoleapp/controllers"
	. "github.com/eaciit/powerplant/sec/consoleapp/models"
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

	// convert(new(WOList), base) // 						done
	// convert(new(AnomaliesWOList), base) //				done
	//// convert(new(Availability), base) //				done
	// convert(new(Consolidated), base) // 					done
	// convert(new(FuelCost), base) //						done
	////// convert(new(FuelTransport), base) // 			done
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
	// convert(new(NotificationFailure), base)
	// convert(new(OperationalData), base) // 				done
	// convert(new(PerformanceFactors), base) // 			done
	convert(new(PlannedMaintenance), base)
	// convert(new(PowerPlantCoordinates), base) // 		done
	// convert(new(PowerPlantInfo),base) // 				done
	// convert(new(PrevMaintenanceValueEquation), base) // 	done
	// convert(new(RPPCloseWO), base) // 					done
	// convert(new(StartupPaymentAndPenalty), base) // 		done
	// convert(new(SyntheticPM),base) // 					done
	// convert(new(UnitCost), base) // 						done
	// convert(new(Vibration), base) // 					done

	// convert(new(SummaryData), base) // -------------------- error comma
	// convert(new(MORSummary), base) // 					done
	// convert(new(MORCalculationFlatSummary), base) // --------------- gak bisa utk anaknya period.year
	// convert(new(PreventiveCorrectiveSummary), base)
	// convert(new(WODurationSummary), base)
	// convert(new(WOListSummary), base)
	// convert(new(FailureAfterPreventiveSummary), base) // done
	// convert(new(RegenMasterPlant), base)

	defer mongo.Close()
	defer sql.Close()
}

func convert(m orm.IModel, base *BaseController) {
	e := base.ConvertMGOToSQLServer(m)
	if e != nil {
		tk.Printf("\nERROR: %v \n", e.Error())
	}
}

func PrepareConnection(ConnectionType string) (dbox.IConnection, error) {
	config := ReadConfig()
	tk.Println(config["host"])
	ci := &dbox.ConnectionInfo{config["host_"+ConnectionType], config["database_"+ConnectionType], config["username_"+ConnectionType], config["password_"+ConnectionType], nil}
	c, e := dbox.NewConnection(ConnectionType, ci)

	if e != nil {
		return nil, e
	}

	e = c.Connect()
	if e != nil {
		return nil, e
	}

	return c, nil
}

func ReadConfig() map[string]string {
	ret := make(map[string]string)
	file, err := os.Open(wd + "conf/app.conf")
	if err == nil {
		defer file.Close()

		reader := bufio.NewReader(file)
		for {
			line, _, e := reader.ReadLine()
			if e != nil {
				break
			}

			sval := strings.Split(string(line), "=")
			ret[sval[0]] = sval[1]
		}
	} else {
		tk.Println(err.Error())
	}

	return ret
}
