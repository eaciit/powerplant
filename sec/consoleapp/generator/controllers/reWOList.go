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
	// ctx := d.BaseController.Ctx
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

			for _, row := range sheet.Rows {
				if len(row.Cells) > 1 && strings.Trim(strings.ToLower(row.Cells[1].String()), " ") != "user status" && strings.Trim(strings.ToLower(row.Cells[1].String()), " ") != "" {

					data := new(WOList)
					data.UserStatus = row.Cells[1].String()
					data.SystemStatus = row.Cells[2].String()
					data.Type = row.Cells[3].String()
					data.OrderCode = row.Cells[4].String()
					data.NotificationCode = row.Cells[5].String()
					data.A = row.Cells[6].String()
					data.NotificationCode = row.Cells[7].String()
					data.EnteredBy = row.Cells[8].String()
					data.Description = row.Cells[9].String()
					data.Plant = row.Cells[10].String()
					data.FunctionalLocation = row.Cells[11].String()
					data.EquipmentCode = row.Cells[12].String()
					data.SortField = row.Cells[13].String()
					data.PriorityText = row.Cells[14].String()
					data.WorkCtr = row.Cells[15].String()
					data.CostCtr = row.Cells[16].String()
					data.MAT = row.Cells[17].String()
					data.CreatedOn, _ = time.Parse(row.Cells[18].String(), "02.01.2013")
					data.BasicStart, _ = time.Parse(row.Cells[19].String(), "02.01.2013")
					data.BasicFinish, _ = time.Parse(row.Cells[21].String(), "02.01.2013")
					data.ScheduledStart, _ = time.Parse(row.Cells[23].String(), "02.01.2013")
					data.ScheduledFinish, _ = time.Parse(row.Cells[25].String(), "02.01.2013")
					data.ActualStart, _ = time.Parse(row.Cells[27].String(), "02.01.2013")
					data.ActualFinish, _ = time.Parse(row.Cells[29].String(), "02.01.2013")
					data.RefDate, _ = time.Parse(row.Cells[31].String(), "02.01.2013")
					data.RespCCTR = row.Cells[33].String()
					data.ReleaseDate, _ = time.Parse(row.Cells[34].String(), "02.01.2013")
					data.AvailFrom, _ = time.Parse(row.Cells[35].String(), "02.01.2013")
					if len(row.Cells) > 37 {
						data.AvailTo, _ = time.Parse(row.Cells[37].String(), "02.01.2013")
					}
					if len(row.Cells) > 39 {
						data.ActualCost, _ = row.Cells[39].Float()
					}
					if len(row.Cells) > 40 {
						data.DF = row.Cells[40].String()
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
