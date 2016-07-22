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
	"sync"
	"time"
)

// REMaintenanceWO
type REMaintenanceWO struct {
	*BaseController
}

// Generate
func (d *REMaintenanceWO) Generate(base *BaseController) {
	var (
		folderName string = "Maintenance WO"
	)
	if base != nil {
		d.BaseController = base
	}
	ctx := d.BaseController.Ctx
	dataSources, path := base.GetDataSource(folderName)
	tk.Println("Generating Maintenance WO from Excel File..")
	temp := ""
	for _, source := range dataSources {
		if strings.Contains(source.Name(), "Historical Maintenance Plans") || strings.Contains(source.Name(), "Historical Work Order") {
			tk.Println(path + "\\" + source.Name())
			file, e := xlsx.OpenFile(path + "\\" + source.Name())
			if e != nil {
				tk.Println(e)
				os.Exit(0)
			}
			sheet := file.Sheets[0]
			wg := sync.WaitGroup{}
			totalData := float64(len(sheet.Rows))
			dataEach := totalData
			if totalData > 10000 {
				dataEach = math.Ceil(float64(totalData / 10))
			}
			if dataEach != totalData {
				wg.Add(10)
				for i := 0; i < 10; i++ {
					go func(wg *sync.WaitGroup, i int, dataEach int) {
						startWith := i * dataEach
						for x := startWith; x <= startWith+dataEach; x++ {
							row := sheet.Rows[x]
							firstCell := ""
							if len(row.Cells) > 0 {
								firstCell, _ = row.Cells[0].String()
							}
							if len(row.Cells) > 0 && strings.Trim(strings.ToLower(firstCell), " ") != "notification" && strings.Trim(strings.ToLower(firstCell), " ") != "" {
								data := new(MaintenanceWorkOrder)
								data.Notification, _ = row.Cells[0].String()
								data.OrderCode, _ = row.Cells[1].String()
								data.UserStatus, _ = row.Cells[2].String()
								data.FunctionalLocation, _ = row.Cells[3].String()
								data.MainWorkCtr, _ = row.Cells[4].String()
								data.Description, _ = row.Cells[5].String()
								data.OrderType, _ = row.Cells[6].String()
								temp, _ = row.Cells[7].String()
								data.ActualRelease, _ = time.Parse(temp, "02.01.2013")
								data.SystemStatus, _ = row.Cells[8].String()
								temp, _ = row.Cells[9].String()
								data.CreatedOn, _ = time.Parse(temp, "02.01.2013")

								_, e := ctx.InsertOut(data)
								if e != nil {
									tk.Println("ERR on file :", source.Name(), " | ROW :", idx)
									tk.Println(e)
								}

							}
						}
						wg.Done()
					}(&wg, i, int(dataEach))
				}
			}
			wg.Wait()
			os.Exit(0)
			// for _, row := range sheet.Rows {
			// 	firstCell := ""
			// 	if len(row.Cells) > 0 {
			// 		firstCell, _ = row.Cells[0].String()
			// 	}
			// 	if len(row.Cells) > 0 && strings.Trim(strings.ToLower(firstCell), " ") != "notification" && strings.Trim(strings.ToLower(firstCell), " ") != "" {
			// 		data := new(MaintenanceWorkOrder)
			// 		data.Notification, _ = row.Cells[0].String()
			// 		data.OrderCode, _ = row.Cells[1].String()
			// 		data.UserStatus, _ = row.Cells[2].String()
			// 		data.FunctionalLocation, _ = row.Cells[3].String()
			// 		data.MainWorkCtr, _ = row.Cells[4].String()
			// 		data.Description, _ = row.Cells[5].String()
			// 		data.OrderType, _ = row.Cells[6].String()
			// 		temp, _ = row.Cells[7].String()
			// 		data.ActualRelease, _ = time.Parse(temp, "02.01.2013")
			// 		data.SystemStatus, _ = row.Cells[8].String()
			// 		temp, _ = row.Cells[9].String()
			// 		data.CreatedOn, _ = time.Parse(temp, "02.01.2013")

			// 		_, e := ctx.InsertOut(data)
			// 		if e != nil {
			// 			tk.Println("ERR on file :", source.Name(), " | ROW :", idx)
			// 			tk.Println(e)
			// 		}

			// 	}
			// }
		}

		if strings.Contains(source.Name(), "maintenance work order PP9") {
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
				if len(row.Cells) > 0 && strings.Trim(strings.ToLower(firstCell), " ") != "notification" && strings.Trim(strings.ToLower(firstCell), " ") != "" {
					data := new(MaintenanceWorkOrder)

					data.OrderCode, _ = row.Cells[0].String()
					data.OrderDescription, _ = row.Cells[1].String()
					data.FunctionalLocation, _ = row.Cells[2].String()
					data.UserStatus, _ = row.Cells[3].String()
					temp, _ = row.Cells[4].String()
					data.CreatedOn, _ = time.Parse(temp, "02.01.2013")
					temp, _ = row.Cells[5].String()
					data.BasStartDate, _ = time.Parse(temp, "02.01.2013")
					data.OrderType, _ = row.Cells[6].String()
					data.PlantWorkCtr, _ = row.Cells[7].String()
					data.CompanyCode, _ = row.Cells[8].String()
					data.SortField, _ = row.Cells[9].String()
					data.Description, _ = row.Cells[10].String()
					data.Equipment, _ = row.Cells[11].String()
					data.MainWorkCtr, _ = row.Cells[12].String()
					data.PlanningPlant, _ = row.Cells[13].String()
					data.CostCtr, _ = row.Cells[14].String()
					data.RespCostCtr, _ = row.Cells[15].String()
					data.ObjectNumber, _ = row.Cells[16].String()
					data.ProfitCtr, _ = row.Cells[17].String()
					data.Priority, _ = row.Cells[18].String()
					data.PriorityDescription, _ = row.Cells[19].String()
					data.Notification, _ = row.Cells[20].String()
					data.Location, _ = row.Cells[21].String()
					data.SystemStatus, _ = row.Cells[22].String()
					data.MainPlant, _ = row.Cells[23].String()

					_, e := ctx.InsertOut(data)
					if e != nil {
						tk.Println("ERR on file :", source.Name(), " | ROW :", idx)
						tk.Println(e)
					}

				}
			}
		}
	}
	tk.Println("Maintenance WO from Excel File : COMPLETE")
}
