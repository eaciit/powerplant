package controllers

import (
	. "github.com/eaciit/powerplant/sec/consoleapp/generator/models"
	tk "github.com/eaciit/toolkit"
)

type DataMaster struct {
	*BaseController
}

func (c *DataMaster) Generate() {
	tk.Println("##Generating Master Data..")
	/*mst := new(Master)
	e := mst.GeneratePlantMaster(c.Ctx)
	if e != nil {
		tk.Println(e)
	}*/
	mst := new(PreventiveSummary)
	e := mst.GeneratePreventiveCorrectiveSummary(c.Ctx)
	if e != nil {
		tk.Println(e)
	}

	tk.Println("##Master Data : DONE\n")
}
