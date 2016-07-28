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

// REVBMOperationalData
type REVBMOperationalData struct {
	*BaseController
}

// Generate
func (d *REVBMOperationalData) Generate(base *BaseController) {
	var (
		folderName string = "VBM Operational Data"
	)
	if base != nil {
		d.BaseController = base
	}
	ctx := d.BaseController.Ctx
	dataSources, path := base.GetDataSource(folderName)
	tk.Println("Generating VBM Operational Data from Excel File..")
	for _, source := range dataSources {
		if strings.Contains(source.Name(), "VBM Operational Data") {
			tk.Println(path + "\\" + source.Name())
			file, e := xlsx.OpenFile(path + "\\" + source.Name())
			if e != nil {
				tk.Println(e)
				os.Exit(0)
			}
			sheet := file.Sheet["Basic Technical Information"]
			for _, row := range sheet.Rows {
				firstCell := ""
				if len(row.Cells) > 0 {
					firstCell, _ = row.Cells[0].String()
				}
				if len(row.Cells) > 0 && strings.Trim(strings.ToLower(firstCell), " ") != "pp name" && strings.Trim(strings.ToLower(firstCell), " ") != "" {
					data := new(PowerPlantInfo)
					data.Name, _ = row.Cells[0].String()
					data.System, _ = row.Cells[1].String()
					data.Province, _ = row.Cells[2].String()
					data.Region, _ = row.Cells[3].String()
					data.City, _ = row.Cells[4].String()
					data.FuelTypes_Crude = row.Cells[5].Bool()
					data.FuelTypes_Diesel = row.Cells[6].Bool()
					data.FuelTypes_Heavy = row.Cells[7].Bool()
					data.FuelTypes_Gas = row.Cells[8].Bool()
					data.GasTurbineUnit, _ = row.Cells[9].Int()
					data.GasTurbineCapacity, _ = row.Cells[10].Float()
					data.SteamUnit, _ = row.Cells[11].Int()
					data.SteamCapacity, _ = row.Cells[12].Float()
					data.DieselUnit, _ = row.Cells[13].Int()
					data.DieselCapacity, _ = row.Cells[14].Float()
					data.CombinedCycleUnit, _ = row.Cells[15].Int()
					data.CombinedCycleCapacity, _ = row.Cells[16].Float()
					_, e := ctx.InsertOut(data)
					if e != nil {
						tk.Println(data)
						tk.Println("ERR on file :", source.Name(), " | ROW :", idx)
						tk.Println(e)
					}

				}
			}

			sheet = file.Sheet["Operational Data"]
			IsDataSource := false
			for _, row := range sheet.Rows {
				firstCell := ""
				if len(row.Cells) > 0 {
					firstCell, _ = row.Cells[0].String()
				}

				if len(row.Cells) > 0 && IsDataSource && strings.Trim(strings.ToLower(firstCell), " ") != "total" && strings.Trim(strings.ToLower(firstCell), " ") != "" {
					data := new(OperationalData)
					data.Plant, _ = row.Cells[0].String()
					data.Year = 2014
					data.GenerationGross = CheckNumberValue(row.Cells[1].Float())
					data.GenerationAux = CheckNumberValue(row.Cells[2].Float())
					data.GenerationNet = CheckNumberValue(row.Cells[3].Float())
					data.ServiceHours = CheckNumberValue(row.Cells[4].Float())
					data.ReserveShutdownHours = CheckNumberValue(row.Cells[5].Float())
					data.MaintenanceOutageHours = CheckNumberValue(row.Cells[6].Float())
					data.ExtendedOutageMaitenanceHours = CheckNumberValue(row.Cells[7].Float())
					data.PlantOutageHours = CheckNumberValue(row.Cells[8].Float())
					data.ExtendedPlanHours = CheckNumberValue(row.Cells[9].Float())
					data.ForcedOutageHours = CheckNumberValue(row.Cells[10].Float())
					data.OutOfManagementControl = CheckNumberValue(row.Cells[11].Float())
					data.MonthBall = CheckNumberValue(row.Cells[12].Float())
					data.UnderCommisionFO = CheckNumberValue(row.Cells[13].Float())
					data.UnderCommisionIS = CheckNumberValue(row.Cells[14].Float())
					data.UnderCommisionMO = CheckNumberValue(row.Cells[15].Float())
					data.UnderCommisionRS = CheckNumberValue(row.Cells[16].Float())
					data.IR = CheckNumberValue(row.Cells[17].Float())
					data.PH = CheckNumberValue(row.Cells[18].Float())
					data.NoOfStartAttempt, _ = row.Cells[19].Int()
					data.NoOfStartActual, _ = row.Cells[20].Int()
					data.FuelDiesel = CheckNumberValue(row.Cells[21].Float())
					data.FuelCrude = CheckNumberValue(row.Cells[22].Float())
					data.FuelHeavy = CheckNumberValue(row.Cells[23].Float())
					data.FuelGas = CheckNumberValue(row.Cells[24].Float())

					_, e := ctx.InsertOut(data)
					if e != nil {
						tk.Println(data)
						tk.Println("ERR on file :", source.Name(), " | ROW :", idx)
						tk.Println(e)
					}

				}

				if strings.ToLower(strings.Replace(firstCell, " ", "", -1)) == "powerplant" {
					IsDataSource = true
				}
			}

			sheet = file.Sheet["Load 2"]
			IsDetail := false
			PlantName := ""
			for _, row := range sheet.Rows {
				firstCell := ""
				if len(row.Cells) > 0 {
					firstCell, _ = row.Cells[0].String()
				}
				if len(row.Cells) > 0 && strings.Trim(strings.ToLower(firstCell), " ") == "plant" && strings.Trim(strings.ToLower(firstCell), " ") != "" {
					PlantName, _ = row.Cells[1].String()
				}
				if firstCell == "" {
					IsDetail = false
				}
				if len(row.Cells) > 0 && IsDetail && strings.Trim(strings.ToLower(firstCell), " ") != "unit" && strings.Trim(strings.ToLower(firstCell), " ") != "" {
					data := new(UnitPower)
					data.Plant = PlantName
					data.Year = 2014
					data.Unit, _ = row.Cells[0].String()
					data.MaxPower, _ = row.Cells[1].Float()
					_, e := ctx.InsertOut(data)
					if e != nil {
						tk.Println(data)
						tk.Println("ERR on file :", source.Name(), " | ROW :", idx)
						tk.Println(e)
					}

				}
				if len(row.Cells) > 0 && strings.Trim(strings.ToLower(firstCell), " ") == "unit" && strings.Trim(strings.ToLower(firstCell), " ") != "" {
					IsDetail = true
				}
			}
		}
	}
	tk.Println("VBM Operational Data from Excel File : COMPLETE")
}
