package controllers

import (
	// "github.com/eaciit/crowd"
	// "github.com/eaciit/dbox"
	// . "github.com/eaciit/powerplant/sec/consoleapp/generator/helpers"
	. "github.com/eaciit/powerplant/sec/library/models"
	tk "github.com/eaciit/toolkit"
	"github.com/tealeg/xlsx"
	"os"
	"strconv"
	"strings"
	"time"
)

// RESyntheticPM
type RESyntheticPM struct {
	*BaseController
}

// Generate
func (d *RESyntheticPM) Generate(base *BaseController, year int) {
	var (
		folderName string = "Synthetic"
	)
	if base != nil {
		d.BaseController = base
	}
	ctx := d.BaseController.Ctx
	dataSources, path := base.GetDataSource(folderName)
	tk.Println("Generating Synthetic PM from Excel File..")
	for _, source := range dataSources {

		tk.Println(path + "\\" + source.Name())
		file, e := xlsx.OpenFile(path + "\\" + source.Name())
		if e != nil {
			tk.Println(e)
			os.Exit(0)
		}
		sheet := file.Sheet[strconv.Itoa(year)]

		for _, row := range sheet.Rows {
			firstCell := ""
			if len(row.Cells) > 0 {
				firstCell, _ = row.Cells[0].String()
			}
			temp := ""
			if len(row.Cells) > 0 && strings.Trim(strings.ToLower(firstCell), " ") != "plant" && strings.Trim(strings.ToLower(firstCell), " ") != "" {
				data := new(SyntheticPM)
				data.Plant, _ = row.Cells[0].String()
				data.Unit, _ = row.Cells[1].String()
				temp, _ = row.Cells[2].String()
				data.ScheduledStart, _ = time.Parse(temp, "02.01.2013")
				data.WOID, _ = row.Cells[3].String()
				data.WOType, _ = row.Cells[4].String()
				data.Description, _ = row.Cells[5].String()
				data.PlannedLaborHours, _ = row.Cells[6].Int()
				data.PlannedLaborCost, _ = row.Cells[7].Float()
				data.ActualMaterialCost, _ = row.Cells[8].Float()
				data.Total, _ = row.Cells[9].Float()

				_, e := ctx.InsertOut(data)
				if e != nil {
					tk.Println("ERR on file :", source.Name(), " | ROW :", idx)
					tk.Println(e)
				}
			}
		}
	}
	tk.Println("Synthetic PM from Excel File : COMPLETE")
}
