package controllers

import (
	// "github.com/eaciit/crowd"
	// "github.com/eaciit/dbox"
	. "github.com/eaciit/powerplant/sec/consoleapp/generator/helpers"
	. "github.com/eaciit/powerplant/sec/library/models"
	tk "github.com/eaciit/toolkit"
	// "strconv"
	"github.com/tealeg/xlsx"
	"os"
	"strings"
	// "time"
)

// REFuelCost
type REFuelCost struct {
	*BaseController
}

// Generate
func (d *REFuelCost) Generate(base *BaseController) {
	var (
		folderName string = "Fuel Cost"
	)
	if base != nil {
		d.BaseController = base
	}
	ctx := d.BaseController.Ctx
	dataSources, path := base.GetDataSource(folderName)
	tk.Println("Generating Fuel Cost Data from Excel File..")
	for _, source := range dataSources {
		if strings.Contains(strings.ToLower(source.Name()), "gen_reports_public_2_allplants") {
			tk.Println(path + "\\" + source.Name())
			file, e := xlsx.OpenFile(path + "\\" + source.Name())
			if e != nil {
				tk.Println(e)
				os.Exit(0)
			}
			sheet := file.Sheet["Units Monthly"]
			for _, row := range sheet.Rows {
				firstCell := ""
				if len(row.Cells) > 1 {
					firstCell, _ = row.Cells[1].String()
				}
				if len(row.Cells) > 1 && strings.Trim(strings.ToLower(firstCell), " ") != "plant_name_en" && strings.Trim(strings.ToLower(firstCell), " ") != "" {

					data := new(FuelCost)

					data.Plant, _ = row.Cells[1].String()
					data.UnitId, _ = row.Cells[5].String()
					data.Year, _ = row.Cells[6].Int()
					data.Month, _ = row.Cells[7].Int()

					data.DependableCapacity = CheckNumberValue(row.Cells[8].Float())
					data.EnergyGrossProduction = CheckNumberValue(row.Cells[9].Float())
					data.EnergyNetProduction = CheckNumberValue(row.Cells[10].Float())
					data.PrimaryFuelType, _ = row.Cells[11].String()
					data.PrimaryFuelGCV = CheckNumberValue(row.Cells[12].Float())
					data.PrimaryFuelConsumed = CheckNumberValue(row.Cells[13].Float())
					data.Primary2FuelType, _ = row.Cells[14].String()
					data.Primary2FuelGCV = CheckNumberValue(row.Cells[15].Float())
					data.Primary2FuelConsumed = CheckNumberValue(row.Cells[16].Float())

					data.BackupFuelType, _ = row.Cells[17].String()
					data.BackupFuelGCV = CheckNumberValue(row.Cells[18].Float())
					data.BackupFuelConsumed = CheckNumberValue(row.Cells[19].Float())
					data.AverageHeatRate = CheckNumberValue(row.Cells[20].Float())

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
	tk.Println("Fuel Cost Data from Excel File : COMPLETE")
}
