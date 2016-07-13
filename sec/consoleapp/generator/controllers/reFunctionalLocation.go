package controllers

import (
	// "github.com/eaciit/crowd"
	// "github.com/eaciit/dbox"
	// . "github.com/eaciit/powerplant/sec/consoleapp/generator/helpers"
	. "github.com/eaciit/powerplant/sec/library/models"
	tk "github.com/eaciit/toolkit"
	// "strconv"
	"strings"
	// "time"
	"github.com/tealeg/xlsx"
	"os"
)

// REFunctionalLocation
type REFunctionalLocation struct {
	*BaseController
}

// Generate
func (d *REFunctionalLocation) Generate(base *BaseController) {
	var (
		folderName string = "Functional Location"
	)
	if base != nil {
		d.BaseController = base
	}
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
			FLCodeColumn := 0
			for _, row := range sheet.Rows {
				d := new(FunctionalLocation)
				for i, cell := range row.Cells {
					if strings.Contains(strings.ToLower(cell.String()), "functional location") && len(strings.Replace(cell.String(), " ", "", -1)) < 20 {
						FLCodeColumn = i
						break
					}
				}
				d.FunctionalLocationCode = row.Cells[FLCodeColumn].String()
				tk.Println(d)
				// Process data for each row (including insert)
			}
		}
	}
}
