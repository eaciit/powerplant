package main

import (
	_ "github.com/eaciit/dbox/dbc/mssql"
	"github.com/eaciit/orm"
	. "github.com/eaciit/powerplant/sec/consoleapp/generator/controllers"
	tk "github.com/eaciit/toolkit"
	"os"
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
	}

	base := new(BaseController)
	base.Ctx = orm.New(sql)
	defer base.Ctx.Close()

	// Generate DataMaster
	// Mst := DataMaster{base}
	// Mst.Generate()

	MOR := DataMOR{base}
	MOR.Generate()

	tk.Println("Application closed..")

}