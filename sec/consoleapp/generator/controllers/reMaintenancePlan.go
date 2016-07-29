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

// REMaintenancePlan
type REMaintenancePlan struct {
	*BaseController
}

// Generate
func (d *REMaintenancePlan) Generate(base *BaseController) {
	var (
		folderName string = "Maintenance Plan"
	)
	if base != nil {
		d.BaseController = base
	}
	ctx := d.BaseController.Ctx
	dataSources, path := base.GetDataSource(folderName)
	tk.Println("Generating Maintenance Plan from Excel File..")
	for _, source := range dataSources {
		if strings.Contains(source.Name(), "2110-PP9-Maintenance Plans") {
			tk.Println(path + "\\" + source.Name())
			file, e := xlsx.OpenFile(path + "\\" + source.Name())
			if e != nil {
				tk.Println(e)
				os.Exit(0)
			}
			sheet := file.Sheets[0]
			firstCell := ""
			temp := ""
			for _, row := range sheet.Rows {
				if len(row.Cells) > 0 {
					firstCell, _ = row.Cells[0].String()
				}
				if len(row.Cells) > 0 && strings.Trim(strings.ToLower(firstCell), " ") != "maintenance item" && strings.Trim(strings.ToLower(firstCell), " ") != "" {
					data := new(MaintenancePlan)
					data.MaintenanceItem, _ = row.Cells[0].String()
					data.MaintenancePlanCode, _ = row.Cells[1].String()
					data.FunctionalLocation, _ = row.Cells[2].String()
					data.FLDescription, _ = row.Cells[3].String()
					data.PWC_TLO, _ = row.Cells[4].String()
					data.PG_PLN, _ = row.Cells[5].String()
					data.TaskListType, _ = row.Cells[6].String()
					data.TaskListGroup, _ = row.Cells[7].String()
					data.TaskListCounter, _ = row.Cells[8].String()
					data.TLCounterDescription, _ = row.Cells[8].String()
					data.Package1 = CheckNumberValue(row.Cells[10].Float())
					data.Package1Text, _ = row.Cells[11].String()
					temp, _ = row.Cells[12].String()
					data.Next1stScheduledDate, _ = time.Parse(temp, "02/01/2013")
					temp, _ = row.Cells[13].String()
					data.Next2ndScheduledDate, _ = time.Parse(temp, "02/01/2013")
					temp, _ = row.Cells[14].String()
					data.Next3rdScheduledDate, _ = time.Parse(temp, "02/01/2013")
					temp, _ = row.Cells[15].String()
					data.Next4thScheduledDate, _ = time.Parse(temp, "02/01/2013")
					data.Equipment, _ = row.Cells[16].String()
					data.EquipmentDescription, _ = row.Cells[17].String()
					data.EquipmentType, _ = row.Cells[18].String()
					data.EquipmentStatus, _ = row.Cells[19].String()
					data.EquipmentStatusText, _ = row.Cells[20].String()
					data.Class, _ = row.Cells[21].String()
					data.Building, _ = row.Cells[22].String()
					data.BuildingText, _ = row.Cells[23].String()
					data.PlanningPlant, _ = row.Cells[24].String()
					data.MaintenancePlant, _ = row.Cells[25].String()
					data.Location, _ = row.Cells[26].String()
					data.PlantSection, _ = row.Cells[27].String()
					data.CompanyCode, _ = row.Cells[28].String()
					data.CostCtr, _ = row.Cells[29].String()
					data.MWC_PLN, _ = row.Cells[30].String()
					data.MWC_PLNDescription, _ = row.Cells[31].String()
					data.MWC_TLH, _ = row.Cells[32].String()
					data.MWC_TLHDescription, _ = row.Cells[33].String()
					data.PWC_TLODescription, _ = row.Cells[34].String()
					data.MWC_EQ, _ = row.Cells[35].String()
					data.MWC_FL, _ = row.Cells[36].String()
					data.PG_PLNDescription, _ = row.Cells[37].String()
					data.PG_FL, _ = row.Cells[38].String()
					data.PG_EQ, _ = row.Cells[39].String()
					data.StdTextKey, _ = row.Cells[40].String()
					data.StdTextKeyDescription, _ = row.Cells[41].String()
					data.OperationShortText, _ = row.Cells[42].String()
					data.DurationOfActivity, _ = row.Cells[43].Int()
					data.UoMDuration, _ = row.Cells[44].String()
					data.TotalWorkInActivity, _ = row.Cells[45].Int()
					data.UoMWork, _ = row.Cells[46].String()
					data.NoOfPerson, _ = row.Cells[47].Int()
					data.ActivityType, _ = row.Cells[48].String()
					data.MaintenanceStrategy, _ = row.Cells[49].String()
					data.Package2 = CheckNumberValue(row.Cells[50].Float())
					data.Package2Text, _ = row.Cells[51].String()
					data.Package3 = CheckNumberValue(row.Cells[52].Float())
					data.Package3Text, _ = row.Cells[53].String()
					data.Package4 = CheckNumberValue(row.Cells[54].Float())
					data.Package4Text, _ = row.Cells[55].String()
					data.Package5 = CheckNumberValue(row.Cells[56].Float())
					data.Package5Text, _ = row.Cells[57].String()
					data.Package6 = CheckNumberValue(row.Cells[58].Float())
					data.Package6Text, _ = row.Cells[59].String()
					data.Package7 = CheckNumberValue(row.Cells[60].Float())
					data.Package7Text, _ = row.Cells[61].String()
					data.PlanCreatedBy, _ = row.Cells[62].String()
					temp, _ = row.Cells[63].String()
					data.PlanCreatedOn, _ = time.Parse(temp, "02/01/2013")
					data.PlanChangedBy, _ = row.Cells[64].String()
					temp, _ = row.Cells[65].String()
					data.PlanChangedOn, _ = time.Parse(temp, "02/01/2013")
					data.Manufacturer, _ = row.Cells[66].String()
					data.ManufacturerSerialNo, _ = row.Cells[67].String()
					data.ManufacturerModelNo, _ = row.Cells[68].String()
					data.ManufacturerPartNo, _ = row.Cells[69].String()

					_, e := ctx.InsertOut(data)
					if e != nil {
						tk.Println("ERR on file :", source.Name(), " | ROW :", idx)
						tk.Println(e)
					}
				}
			}

		}
	}
	tk.Println("Maintenance Plan from Excel File : COMPLETE")
}
