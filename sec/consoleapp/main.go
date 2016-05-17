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
	defer mongo.Close()

	sql, e := PrepareConnection("mssql")
	if e != nil {
		tk.Println(e)
	}
	defer sql.Close()

	base := new(BaseController)
	base.MongoCtx = orm.New(mongo)
	base.SqlCtx = orm.New(sql)

	base.ConvertMGOToSQLServer(new(WOList))
	base.ConvertMGOToSQLServer(new(AnomaliesWOList))
	base.ConvertMGOToSQLServer(new(Availability))
	base.ConvertMGOToSQLServer(new(Consolidated))
	base.ConvertMGOToSQLServer(new(FuelCost))
	base.ConvertMGOToSQLServer(new(FuelTransport))
	base.ConvertMGOToSQLServer(new(FunctionalLocation))
	base.ConvertMGOToSQLServer(new(AnomaliesFunctionalLocation))
	base.ConvertMGOToSQLServer(new(GenerationAppendix))
	base.ConvertMGOToSQLServer(new(MaintenanceCost))
	base.ConvertMGOToSQLServer(new(MaintenanceCostByHour))
	base.ConvertMGOToSQLServer(new(MaintenancePlan))
	base.ConvertMGOToSQLServer(new(MaintenanceWorkOrder))
	base.ConvertMGOToSQLServer(new(MappedEquipmentType))
	base.ConvertMGOToSQLServer(new(MasterEquipmentType))
	base.ConvertMGOToSQLServer(new(MasterMROElement))
	base.ConvertMGOToSQLServer(new(MasterOrderType))
	base.ConvertMGOToSQLServer(new(MasterPlant))
	base.ConvertMGOToSQLServer(new(NewEquipmentType))
	base.ConvertMGOToSQLServer(new(NotificationFailure))
	base.ConvertMGOToSQLServer(new(OperationalData))
	base.ConvertMGOToSQLServer(new(PerformanceFactors))
	base.ConvertMGOToSQLServer(new(PlannedMaintenance))
	base.ConvertMGOToSQLServer(new(PowerPlantCoordinates))
	base.ConvertMGOToSQLServer(new(PowerPlantInfo))
	base.ConvertMGOToSQLServer(new(PrevMaintenanceValueEquation))
	base.ConvertMGOToSQLServer(new(RPPCloseWO))
	base.ConvertMGOToSQLServer(new(StartupPaymentAndPenalty))
	base.ConvertMGOToSQLServer(new(SyntheticPM))
	base.ConvertMGOToSQLServer(new(UnitCost))
	base.ConvertMGOToSQLServer(new(Vibration))

	// s := Sample{base}
	// s.GetSampleData()
	// s.UpdateSampleData()
	// s.GetSampleData()
	// s.InsertSampleData()
	// s.RemoveSampleData()
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
