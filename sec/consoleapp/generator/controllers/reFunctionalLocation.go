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

// REFunctionalLocation
type REFunctionalLocation struct {
	*BaseController
}

// Generate
func (d *REFunctionalLocation) Generate(base *BaseController) {
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
	totalData := 0
	for _, source := range dataSources {
		if strings.Contains(source.Name(), "FLOC Structure") {
			file, e := xlsx.OpenFile(path + "\\" + source.Name())
			if e != nil {
				tk.Println(e)
				os.Exit(0)
			}
			sheet := file.Sheets[0]
			totalDataEachFile := 0
			totalInsertedData := 0
			for _, row := range sheet.Rows {
				firstCell := ""
				if len(row.Cells) > 0 {
					firstCell, _ = row.Cells[1].String()
				}
				if len(row.Cells) > 0 && strings.Trim(strings.ToLower(firstCell), " ") != "functional location" && strings.Trim(strings.ToLower(firstCell), " ") != "" {
					totalDataEachFile++
					totalData++
					Str, _ := row.Cells[2].String()
					Description, _ := row.Cells[3].String()
					SupFunctionalLocation, _ := row.Cells[25].String()
					temp := ""
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
							if len(row.Cells) >= 1 {
								data.FunctionalLocationCode, _ = row.Cells[1].String()
							}
							if len(row.Cells) >= 2 {
								data.Str, _ = row.Cells[2].String()
							}
							if len(row.Cells) >= 3 {
								data.Description, _ = row.Cells[3].String()
							}
							if len(row.Cells) >= 4 {
								data.CostCtr, _ = row.Cells[4].String()
							}
							if len(row.Cells) >= 5 {
								data.Location, _ = row.Cells[5].String()
							}
							if len(row.Cells) >= 6 {
								data.PIPI, _ = row.Cells[6].String()
							}
							if len(row.Cells) >= 7 {
								data.PInt, _ = row.Cells[7].String()
							}
							if len(row.Cells) >= 8 {
								data.MainWorkCtr, _ = row.Cells[8].String()
							}
							if len(row.Cells) >= 9 {
								data.CatProf, _ = row.Cells[9].String()
							}
							if len(row.Cells) >= 10 {
								data.SortField, _ = row.Cells[10].String()
							}
							if len(row.Cells) >= 11 {
								data.ModelNo, _ = row.Cells[11].String()
							}
							if len(row.Cells) >= 12 {
								data.SerNo, _ = row.Cells[12].String()
							}
							if len(row.Cells) >= 13 {
								data.UserStatus, _ = row.Cells[13].String()
							}
							if len(row.Cells) >= 14 {
								data.A, _ = row.Cells[14].String()
							}
							if len(row.Cells) >= 15 {
								data.ObjectType, _ = row.Cells[15].String()
							}
							if len(row.Cells) >= 16 {
								data.PG, _ = row.Cells[16].String()
							}
							if len(row.Cells) >= 17 {
								data.ManParNo, _ = row.Cells[17].String()
							}
							if len(row.Cells) >= 18 {
								data.Asset, _ = row.Cells[18].String()
							}
							if len(row.Cells) >= 19 {
								temp, _ = row.Cells[19].String()
								data.Date, _ = time.Parse(temp, "2013.01.02")
							}
							if len(row.Cells) >= 20 {
								data.AcqValue, _ = row.Cells[20].String()
							}
							if len(row.Cells) >= 22 {
								data.InvNo, _ = row.Cells[21].String()
							}
							if len(row.Cells) >= 22 {
								data.ConstType, _ = row.Cells[22].String()
							}
							if len(row.Cells) >= 23 {
								temp, _ = row.Cells[23].String()
								data.StartFrom, _ = time.Parse(temp, "2013.01.02")
							}
							if len(row.Cells) >= 24 {
								temp, _ = row.Cells[24].String()
								data.CreatedOn, _ = time.Parse(temp, "2013.01.02")
							}
							if len(row.Cells) >= 25 {
								data.SupFunctionalLocation, _ = row.Cells[25].String()
							}
							_, e := ctx.InsertOut(data)
							if e != nil {
								tk.Println("ERR on file :", source.Name(), " | ROW :", idx)
								tk.Println(e)
							} else {
								totalInsertedData++
							}
						}
					} else {
						data := new(AnomaliesFunctionalLocation)
						data.FunctionalLocationCode, _ = row.Cells[1].String()
						data.Str, _ = row.Cells[2].String()
						data.Description, _ = row.Cells[3].String()
						data.CostCtr, _ = row.Cells[4].String()
						data.Location, _ = row.Cells[5].String()
						data.PIPI, _ = row.Cells[6].String()
						data.PInt, _ = row.Cells[7].String()
						data.MainWorkCtr, _ = row.Cells[8].String()
						data.CatProf, _ = row.Cells[9].String()

						data.SortField, _ = row.Cells[10].String()
						data.ModelNo, _ = row.Cells[11].String()
						data.SerNo, _ = row.Cells[12].String()
						data.UserStatus, _ = row.Cells[13].String()
						data.A, _ = row.Cells[14].String()

						data.ObjectType, _ = row.Cells[15].String()
						data.PG, _ = row.Cells[16].String()
						data.ManParNo, _ = row.Cells[17].String()
						data.Asset, _ = row.Cells[18].String()
						temp, _ = row.Cells[19].String()
						data.Date, _ = time.Parse(temp, "2013.01.02")
						data.AcqValue, _ = row.Cells[20].String()
						data.InvNo, _ = row.Cells[21].String()
						data.ConstType, _ = row.Cells[22].String()
						temp, _ = row.Cells[23].String()
						data.StartFrom, _ = time.Parse(temp, "2013.01.02")
						temp, _ = row.Cells[24].String()
						data.CreatedOn, _ = time.Parse(temp, "2013.01.02")
						data.SupFunctionalLocation, _ = row.Cells[25].String()
						_, e := ctx.InsertOut(data)
						if e != nil {
							tk.Println("ERR on file :", source.Name(), " | ROW :", idx)
							tk.Println(e)
						} else {
							totalInsertedData++
						}
					}
				}
			}
			tk.Println(source.Name(), " : ", totalInsertedData, " / ", totalDataEachFile, " | Max Row : ", sheet.MaxRow)
		}
	}
	tk.Println("TOTAL DATA : ", totalData)
	tk.Println("Functional Location from Excel File : COMPLETE")
}
