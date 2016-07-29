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
	"time"
)

// REConsolidated
type REConsolidated struct {
	*BaseController
}

// Generate
func (d *REConsolidated) Generate(base *BaseController) {
	var (
		folderName string = "Consolidated"
	)
	if base != nil {
		d.BaseController = base
	}
	ctx := d.BaseController.Ctx
	dataSources, path := base.GetDataSource(folderName)
	tk.Println("Generating Consolidated Data from Excel File..")
	for _, source := range dataSources {
		if strings.Contains(source.Name(), "Consolidated") {
			tk.Println(path + "\\" + source.Name())
			file, e := xlsx.OpenFile(path + "\\" + source.Name())
			if e != nil {
				tk.Println(e)
				os.Exit(0)
			}
			sheet := file.Sheet["Summary"]

			for _, row := range sheet.Rows {
				firstCell := ""
				if len(row.Cells) > 0 {
					firstCell, _ = row.Cells[0].String()
				}
				temp := ""
				if len(row.Cells) > 0 && strings.Trim(strings.ToLower(firstCell), " ") != "power plant" && strings.Trim(strings.ToLower(firstCell), " ") != "" {
					data := new(Consolidated)
					temp, _ = row.Cells[0].String()
					data.Plant = PlantNormalization(temp)
					data.Unit, _ = row.Cells[1].String()
					temp, _ = row.Cells[2].String()
					temp = strings.Replace(temp, "\\", "", -1)
					data.ConsolidatedDate, _ = time.Parse("02-01-2006", temp)
					data.Load1 = CheckNumberValue(row.Cells[3].Float())
					data.Load2 = CheckNumberValue(row.Cells[4].Float())
					data.Load3 = CheckNumberValue(row.Cells[5].Float())
					data.Load4 = CheckNumberValue(row.Cells[6].Float())
					data.Load5 = CheckNumberValue(row.Cells[7].Float())
					data.Load6 = CheckNumberValue(row.Cells[8].Float())
					data.Load7 = CheckNumberValue(row.Cells[9].Float())
					data.Load8 = CheckNumberValue(row.Cells[10].Float())
					data.Load9 = CheckNumberValue(row.Cells[11].Float())
					data.Load10 = CheckNumberValue(row.Cells[12].Float())
					data.Load11 = CheckNumberValue(row.Cells[13].Float())
					data.Load12 = CheckNumberValue(row.Cells[14].Float())
					data.Load13 = CheckNumberValue(row.Cells[15].Float())
					data.Load14 = CheckNumberValue(row.Cells[16].Float())
					data.Load15 = CheckNumberValue(row.Cells[17].Float())
					data.Load16 = CheckNumberValue(row.Cells[18].Float())
					data.Load17 = CheckNumberValue(row.Cells[19].Float())
					data.Load18 = CheckNumberValue(row.Cells[20].Float())
					data.Load19 = CheckNumberValue(row.Cells[21].Float())
					data.Load20 = CheckNumberValue(row.Cells[22].Float())
					data.Load21 = CheckNumberValue(row.Cells[23].Float())
					data.Load22 = CheckNumberValue(row.Cells[24].Float())
					data.Load23 = CheckNumberValue(row.Cells[25].Float())
					data.Load0 = CheckNumberValue(row.Cells[26].Float())

					data.EnergyGross = CheckNumberValue(row.Cells[27].Float())
					data.EnergyNet = CheckNumberValue(row.Cells[28].Float())
					data.Capacity = CheckNumberValue(row.Cells[29].Float())
					data.FuelType, _ = row.Cells[30].String()

					data.FuelConsumption_Gas = CheckNumberValue(row.Cells[31].Float())
					data.FuelConsumption_Diesel = CheckNumberValue(row.Cells[32].Float())
					data.FuelConsumption_Crude = CheckNumberValue(row.Cells[33].Float())

					data.TotalCapacity = CheckNumberValue(row.Cells[34].Float())
					data.CapacityPayment = CheckNumberValue(row.Cells[35].Float())
					data.EnergyPayment = CheckNumberValue(row.Cells[36].Float())
					data.StartupPayment = CheckNumberValue(row.Cells[37].Float())
					data.Penalty = CheckNumberValue(row.Cells[38].Float())
					data.Incentive = CheckNumberValue(row.Cells[39].Float())

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
	tk.Println("Consolidated Data from Excel File : COMPLETE")
}
