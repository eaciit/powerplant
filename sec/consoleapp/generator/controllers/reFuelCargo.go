package controllers

import (
	// "github.com/eaciit/crowd"
	// "github.com/eaciit/dbox"
	// . "github.com/eaciit/powerplant/sec/consoleapp/generator/helpers"
	. "github.com/eaciit/powerplant/sec/library/models"
	tk "github.com/eaciit/toolkit"
	// "strconv"
	"github.com/tealeg/xlsx"
	"math"
	"os"
	"strings"
)

// REFuelCargo
type REFuelCargo struct {
	*BaseController
}

// Generate
func (d *REFuelCargo) Generate(base *BaseController) {
	var (
		folderName string = "Fuel Cargo"
	)
	if base != nil {
		d.BaseController = base
	}
	ctx := d.BaseController.Ctx
	dataSources, path := base.GetDataSource(folderName)
	tk.Println("Generating Fuel Cargo from Excel File..")
	for _, source := range dataSources {
		if strings.Contains(source.Name(), "Fuel Cargo") {
			tk.Println(path + "\\" + source.Name())
			file, e := xlsx.OpenFile(path + "\\" + source.Name())
			if e != nil {
				tk.Println(e)
				os.Exit(0)
			}
			sheet := file.Sheets[0]

			for _, row := range sheet.Rows {
				firstCell := ""
				if len(row.Cells) > 0 {
					firstCell, _ = row.Cells[0].String()
				}
				if len(row.Cells) > 0 && strings.Trim(strings.ToLower(firstCell), " ") != "plant" && strings.Trim(strings.ToLower(firstCell), " ") != "" {
					for i, cell := range row.Cells {
						if i > 1 {
							data := new(FuelTransport)
							data.Plant, _ = row.Cells[0].String()
							data.Year, e = sheet.Rows[1].Cells[i].Int()
							data.TransportCost, _ = cell.Float()
							if math.IsNaN(data.TransportCost) {
								data.TransportCost = 0
							}
							_, e := ctx.InsertOut(data)
							if e != nil {
								tk.Println(data)
								tk.Println("ERR on file :", source.Name(), " | ROW :", idx)
								tk.Println(e)
							}
						}

					}
				}
			}
		}
	}
	tk.Println("Fuel Cargo from Excel File : COMPLETE")
}
