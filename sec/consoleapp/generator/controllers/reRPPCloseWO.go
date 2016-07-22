package controllers

import (
	// "github.com/eaciit/crowd"
	// "github.com/eaciit/dbox"
	// . "github.com/eaciit/powerplant/sec/consoleapp/generator/helpers"
	. "github.com/eaciit/powerplant/sec/library/models"
	tk "github.com/eaciit/toolkit"
	// "strconv"
	"github.com/tealeg/xlsx"
	"os"
	"strings"
	"time"
)

// RERPPCloseWO
type RERPPCloseWO struct {
	*BaseController
}

// Generate
func (d *RERPPCloseWO) Generate(base *BaseController) {
	var (
		folderName string = "RPPClose WO"
	)
	if base != nil {
		d.BaseController = base
	}
	ctx := d.BaseController.Ctx
	dataSources, path := base.GetDataSource(folderName)
	tk.Println("Generating RPPClose WO from Excel File..")
	for _, source := range dataSources {
		if strings.Contains(source.Name(), "RPPClose WO") {
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
				temp := ""
				if len(row.Cells) > 0 && strings.Trim(strings.ToLower(firstCell), " ") != "notification" && strings.Trim(strings.ToLower(firstCell), " ") != "" {
					data := new(RPPCloseWO)
					data.Notification, _ = row.Cells[0].String()
					data.OrderCode, _ = row.Cells[1].String()
					data.UserStatus, _ = row.Cells[2].String()
					data.FunctionalLocation, _ = row.Cells[3].String()
					data.MainWorkCtr, _ = row.Cells[4].String()
					data.Description, _ = row.Cells[5].String()
					data.OrderType, _ = row.Cells[6].String()
					temp, _ = row.Cells[7].String()
					data.ReferenceDate, _ = time.Parse(temp, "02.01.2013")
					data.WorkCenter, _ = row.Cells[8].String()
					data.Equipment, _ = row.Cells[9].String()
					data.TotalActCost, _ = row.Cells[10].Float()

					temp, _ = row.Cells[11].String()
					data.ActualStart, _ = time.Parse(temp, "02.01.2013")
					data.CostCenter, _ = row.Cells[12].String()
					data.TotalSettlement, _ = row.Cells[13].Float()
					data.PlannerGroup, _ = row.Cells[14].String()
					data.WBSElement, _ = row.Cells[13].String()
					data.SystemStatus, _ = row.Cells[14].String()

					_, e := ctx.InsertOut(data)
					if e != nil {
						tk.Println("ERR on file :", source.Name(), " | ROW :", idx)
						tk.Println(e)
					}

				}
			}
		}
	}
	tk.Println("RPPClose WO from Excel File : COMPLETE")
}
