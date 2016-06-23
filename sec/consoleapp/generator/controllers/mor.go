package controllers

import (
	. "github.com/eaciit/powerplant/sec/consoleapp/generator/models"
	tk "github.com/eaciit/toolkit"
)

type DataMOR struct {
	*BaseController
}

func (c *DataMOR) Generate() {
	var (
		e error
	)
	tk.Println("##Generating MOR Data..")
	mor := new(MOR)
	// e = mor.GenerateMORSummary(c.Ctx)
	// if e != nil {
	// 	tk.Println(e)
	// }
	e = mor.GenerateMORFlatCalculationSummary(c.Ctx)
	if e != nil {
		tk.Println(e)
	}
	tk.Println("##MOR Data : DONE\n")
}
