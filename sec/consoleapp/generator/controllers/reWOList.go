package controllers

/*import (
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

// REWOList
type REWOList struct {
	*BaseController
}

// Generate
func (d *REWOList) Generate(base *BaseController) {
	var (
		folderName string = "WO List"
	)
	if base != nil {
		d.BaseController = base
	}
	ctx := d.BaseController.Ctx
	dataSources, path := base.GetDataSource(folderName)
	tk.Println("Generating WO List from Excel File..")
	for _, source := range dataSources {
		if strings.Contains(source.Name(), "WO List") {
			tk.Println(path + "\\" + source.Name())
			file, e := xlsx.OpenFile(path + "\\" + source.Name())
			if e != nil {
				tk.Println(e)
				os.Exit(0)
			}
			sheet := file.Sheets[0]
			temp := ""
			for _, row := range sheet.Rows {
				firstCell := ""
				if len(row.Cells) > 1 {
					firstCell, _ = row.Cells[1].String()
				}
				if len(row.Cells) > 1 && strings.Trim(strings.ToLower(firstCell), " ") != "user status" && strings.Trim(strings.ToLower(firstCell), " ") != "" {

					data := new(WOList)
					data.UserStatus, _ = row.Cells[1].String()
					data.SystemStatus, _ = row.Cells[2].String()
					data.Type, _ = row.Cells[3].String()
					data.OrderCode, _ = row.Cells[4].String()
					data.NotificationCode, _ = row.Cells[5].String()
					data.A, _ = row.Cells[6].String()
					data.NotificationCode, _ = row.Cells[7].String()
					data.EnteredBy, _ = row.Cells[8].String()
					data.Description, _ = row.Cells[9].String()
					data.Plant, _ = row.Cells[10].String()
					data.FunctionalLocation, _ = row.Cells[11].String()
					data.EquipmentCode, _ = row.Cells[12].String()
					data.SortField, _ = row.Cells[13].String()
					data.PriorityText, _ = row.Cells[14].String()
					data.WorkCtr, _ = row.Cells[15].String()
					data.CostCtr, _ = row.Cells[16].String()
					data.MAT, _ = row.Cells[17].String()
					temp, _ = row.Cells[18].String()
					data.CreatedOn, _ = time.Parse(temp, "02.01.2013")
					data.BasicStart, _ = time.Parse(temp, "02.01.2013")
					temp, _ = row.Cells[21].String()
					data.BasicFinish, _ = time.Parse(temp, "02.01.2013")
					temp, _ = row.Cells[23].String()
					data.ScheduledStart, _ = time.Parse(temp, "02.01.2013")
					temp, _ = row.Cells[25].String()
					data.ScheduledFinish, _ = time.Parse(temp, "02.01.2013")
					temp, _ = row.Cells[27].String()
					data.ActualStart, _ = time.Parse(temp, "02.01.2013")
					temp, _ = row.Cells[29].String()
					data.ActualFinish, _ = time.Parse(temp, "02.01.2013")
					temp, _ = row.Cells[31].String()
					data.RefDate, _ = time.Parse(temp, "02.01.2013")
					data.RespCCTR, _ = row.Cells[33].String()
					data.ReleaseDate, _ = time.Parse(temp, "02.01.2013")
					data.RespCCTR, _ = row.Cells[35].String()
					data.AvailFrom, _ = time.Parse(temp, "02.01.2013")
					if len(row.Cells) > 37 {
						data.RespCCTR, _ = row.Cells[37].String()
						data.AvailTo, _ = time.Parse(temp, "02.01.2013")
					}
					if len(row.Cells) > 39 {
						data.ActualCost, _ = row.Cells[39].Float()
					}
					if len(row.Cells) > 40 {
						data.DF, _ = row.Cells[40].String()
					}

					_, e := ctx.InsertOut(data)
					if e != nil {
						tk.Println("ERR on file :", source.Name(), " | ROW :", idx)
						tk.Println(e)
					}

				}
			}
		}
	}
	tk.Println("WO List from Excel File : COMPLETE")
}
*/
