package controllers

import (
	. "github.com/eaciit/powerplant/sec/consoleapp/generator/models"
	tk "github.com/eaciit/toolkit"
)

type DataSummary struct {
	*BaseController
}

func (s *DataSummary) GenerateSummary() {
	tk.Println("##Generating Summary Data..")
	sum := new(Summary)
	e := sum.GenerateSummaryData(s.Ctx)
	if e != nil {
		tk.Println(e)
	}
	tk.Println("##Summary Data : DONE\n")
}
