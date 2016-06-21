package main

import (
	_ "github.com/eaciit/dbox/dbc/mssql"
	"github.com/eaciit/orm"
	. "github.com/eaciit/powerplant/sec/consoleapp/generator/controllers"
	"os"
	// . "github.com/eaciit/powerplant/sec/library/models"
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
	}

	base := new(BaseController)
	base.Ctx = orm.New(sql)
	defer base.Ctx.Close()

	// Generate DataMaster
	Mst := DataMaster{base}
	Mst.Generate()

	tk.Println("Application closed..")

}
