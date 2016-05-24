package controllers

import (
	_ "github.com/eaciit/dbox/dbc/mongo"
	_ "github.com/eaciit/dbox/dbc/mssql"
	// . "github.com/eaciit/powerplant/sec/consoleapp/models"
)

/*
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	tk.Println("Generate Data...")
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

	generateSummaryData(base)
}

func generateSummaryData(base *BaseController) (e error) {

	return
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
}*/
