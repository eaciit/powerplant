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

// REFunctionalLocation
type REFunctionalLocation struct {
	*BaseController
}
*/
// Generate
/*func (d *REFunctionalLocation) Generate(base *BaseController) {
	var (
		folderName string   = "Functional Location"
		StrInds    []string = []string{"FMS", "DISTD", "DISTM", "TRNS", "GN-01", "GN-02", "GN-03", "GN-04"}
	)
	if base != nil {
		d.BaseController = base
	}
	ctx := d.BaseController.Ctx
	dataSources, path := base.GetDataSource(folderName)
	tk.Println("Generating Functional Location from Excel File..")
	for _, source := range dataSources {
		if strings.Contains(source.Name(), "FLOC Structure") {
			file, e := xlsx.OpenFile(path + "\\" + source.Name())
			if e != nil {
				tk.Println(e)
				os.Exit(0)
			}
			sheet := file.Sheets[0]

			for idx, row := range sheet.Rows {
				if strings.Trim(strings.ToLower(row.Cells[1].String()), " ") != "functional location" && strings.Trim(strings.ToLower(row.Cells[1].String()), " ") != "" {
					Str := row.Cells[2].String()
					Description := row.Cells[3].String()
					SupFunctionalLocation := row.Cells[25].String()
					if len(Description) <= 200 && SupFunctionalLocation != "" {
						isMatch := false
						for _, s := range StrInds {
							if strings.Contains(s, strings.Trim(Str, " ")) {
								isMatch = true
								break
							}
						}
						if isMatch {
							data := new(FunctionalLocation)
							data.FunctionalLocationCode = row.Cells[1].String()
							data.Str = row.Cells[2].String()
							data.Description = row.Cells[3].String()
							data.CostCtr = row.Cells[4].String()
							data.Location = row.Cells[5].String()
							data.PIPI = row.Cells[6].String()
							data.PInt = row.Cells[7].String()
							data.MainWorkCtr = row.Cells[8].String()
							data.CatProf = row.Cells[9].String()

							data.SortField = row.Cells[10].String()
							data.ModelNo = row.Cells[11].String()
							data.SerNo = row.Cells[12].String()
							data.UserStatus = row.Cells[13].String()
							data.A = row.Cells[14].String()

							data.ObjectType = row.Cells[15].String()
							data.PG = row.Cells[16].String()
							data.ManParNo = row.Cells[17].String()
							data.Asset = row.Cells[18].String()
							data.Date, _ = time.Parse(row.Cells[19].String(), "2013.01.02")
							data.AcqValue = row.Cells[20].String()
							data.InvNo = row.Cells[21].String()
							data.ConstType = row.Cells[22].String()
							data.StartFrom, _ = time.Parse(row.Cells[23].String(), "2013.01.02")
							data.CreatedOn, _ = time.Parse(row.Cells[24].String(), "2013.01.02")
							data.SupFunctionalLocation = row.Cells[25].String()
							_, e := ctx.InsertOut(data)
							if e != nil {
								tk.Println("ERR on file :", source.Name(), " | ROW :", idx)
								tk.Println(e)
							}
						}
					} else {
						data := new(AnomaliesFunctionalLocation)
						data.FunctionalLocationCode = row.Cells[1].String()
						data.Str = row.Cells[2].String()
						data.Description = row.Cells[3].String()
						data.CostCtr = row.Cells[4].String()
						data.Location = row.Cells[5].String()
						data.PIPI = row.Cells[6].String()
						data.PInt = row.Cells[7].String()
						data.MainWorkCtr = row.Cells[8].String()
						data.CatProf = row.Cells[9].String()

						data.SortField = row.Cells[10].String()
						data.ModelNo = row.Cells[11].String()
						data.SerNo = row.Cells[12].String()
						data.UserStatus = row.Cells[13].String()
						data.A = row.Cells[14].String()

						data.ObjectType = row.Cells[15].String()
						data.PG = row.Cells[16].String()
						data.ManParNo = row.Cells[17].String()
						data.Asset = row.Cells[18].String()
						data.Date, _ = time.Parse(row.Cells[19].String(), "2013.01.02")
						data.AcqValue = row.Cells[20].String()
						data.InvNo = row.Cells[21].String()
						data.ConstType = row.Cells[22].String()
						data.StartFrom, _ = time.Parse(row.Cells[23].String(), "2013.01.02")
						data.CreatedOn, _ = time.Parse(row.Cells[24].String(), "2013.01.02")
						data.SupFunctionalLocation = row.Cells[25].String()
						_, e := ctx.InsertOut(data)
						if e != nil {
							tk.Println("ERR on file :", source.Name(), " | ROW :", idx)
							tk.Println(e)
						}
					}
				}
			}
		}
	}
	tk.Println("Functional Location from Excel File : COMPLETE")
}*/
