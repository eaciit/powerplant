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

// REPreventiveMaintenance
type REPreventiveMaintenance struct {
	*BaseController
}

// Generate
func (d *REPreventiveMaintenance) Generate(base *BaseController) {
	var (
		folderName string = "Preventive Maintenance"
	)
	if base != nil {
		d.BaseController = base
	}
	ctx := d.BaseController.Ctx
	dataSources, path := base.GetDataSource(folderName)
	tk.Println("Generating Preventive Maintenance from Excel File..")
	for _, source := range dataSources {
		if strings.Contains(source.Name(), "Preventive Maintenance") {
			tk.Println(path + "\\" + source.Name())
			file, e := xlsx.OpenFile(path + "\\" + source.Name())
			if e != nil {
				tk.Println(e)
				os.Exit(0)
			}
			sheet := file.Sheet["PM"]
			filenames := strings.Split(source.Name(), " ")
			PlantName := ""
			if len(filenames) > 0 {
				PlantName = filenames[0]
			}
			temp := ""
			Phase := ""
			Block := ""
			Unit := ""
			if strings.Contains(source.Name(), "GT") {
				for _, row := range sheet.Rows {
					firstCell := ""
					if len(row.Cells) > 0 {
						firstCell, _ = row.Cells[2].String()
					}
					if len(row.Cells) > 0 && strings.Trim(strings.ToLower(firstCell), " ") != "unit" && strings.Trim(strings.ToLower(firstCell), " ") != "" {

						data := new(PrevMaintenanceValueEquation)
						data.Plant = PlantName
						data.Phase, _ = row.Cells[0].String()
						data.Block, _ = row.Cells[1].String()

						if data.Phase != "" {
							Phase = data.Phase
						} else {
							data.Phase = Phase
						}

						if data.Block != "" {
							Block = data.Block
						} else {
							data.Block = Block
						}

						data.Unit, _ = row.Cells[2].String()
						data.Id2, _ = row.Cells[3].String()
						temp, _ = row.Cells[4].String()
						data.DatePerformed, _ = time.Parse("02/01/2006", strings.Replace(temp, " ", "", -1))
						data.WOType, _ = row.Cells[5].String()
						data.UserStatus, _ = row.Cells[6].String()
						data.Description, _ = row.Cells[7].String()
						data.Days, _ = row.Cells[8].Int()
						data.Materials = CheckNumberValue(row.Cells[9].Float())
						data.SkilledLabour = CheckNumberValue(row.Cells[10].Float())
						data.UnSkilledLabour = CheckNumberValue(row.Cells[11].Float())
						data.ExtraCost = 0
						data.ContractMaintenance = CheckNumberValue(row.Cells[12].Float())
						data.TotalCost = CheckNumberValue(row.Cells[13].Float())

						_, e := ctx.InsertOut(data)
						if e != nil {
							tk.Println(data)
							tk.Println("ERR on file :", source.Name(), " | ROW :", idx)
							tk.Println(e)
						}

					}
				}
			} else if strings.Contains(source.Name(), "ST") {
				for _, row := range sheet.Rows {
					firstCell := ""
					if len(row.Cells) > 0 {
						firstCell, _ = row.Cells[1].String()
					}

					if len(row.Cells) > 0 && ((strings.Trim(strings.ToLower(firstCell), " ") != "phase" && strings.Trim(strings.ToLower(firstCell), " ") != "") || Phase != "") {
						data := new(PrevMaintenanceValueEquation)
						data.Plant = PlantName
						data.Phase, _ = row.Cells[1].String()
						data.Block, _ = row.Cells[2].String()

						if data.Phase != "" {
							Phase = data.Phase
						} else {
							data.Phase = Phase
						}

						if data.Block != "" {
							Block = data.Block
						} else {
							data.Block = Block
						}
						tempUnit, _ := row.Cells[3].String()
						tempCombined, _ := row.Cells[4].String()
						if tempUnit != "" && tempCombined == "" {
							Unit = tempUnit
						} else {
							data.Unit = Unit
							data.DatePerformed = time.Date(2014, 1, 1, 0, 0, 0, 0, time.UTC)
							data.Id2, _ = row.Cells[5].String()
							data.WOType, _ = row.Cells[6].String()
							data.UserStatus, _ = row.Cells[7].String()
							data.Description, _ = row.Cells[8].String()
							data.Days, _ = row.Cells[9].Int()
							data.Materials = CheckNumberValue(row.Cells[10].Float())
							data.SkilledLabour = CheckNumberValue(row.Cells[11].Float())
							data.ExtraCost = CheckNumberValue(row.Cells[12].Float())
							data.ContractMaintenance = CheckNumberValue(row.Cells[13].Float())
							data.TotalCost = CheckNumberValue(row.Cells[14].Float())
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
	tk.Println("Preventive Maintenance from Excel File : COMPLETE")
}
