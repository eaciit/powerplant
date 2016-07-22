package controllers

import (
	// "github.com/eaciit/crowd"
	// "github.com/eaciit/dbox"
	// . "github.com/eaciit/powerplant/sec/consoleapp/generator/helpers"
	. "github.com/eaciit/powerplant/sec/library/models"
	tk "github.com/eaciit/toolkit"
	"github.com/tealeg/xlsx"
	"os"
	"strconv"
	"strings"
	// "time"
)

// REPerfromanceFactors
type REPerfromanceFactors struct {
	*BaseController
}

// Generate
func (d *REPerfromanceFactors) Generate(base *BaseController, Plant string) {
	var (
		folderName string = "Perfromance Factors"
	)
	if base != nil {
		d.BaseController = base
	}
	ctx := d.BaseController.Ctx
	dataSources, path := base.GetDataSource(folderName)
	tk.Println("Generating Perfromance Factors from Excel File..")
	for _, source := range dataSources {

		tk.Println(path + "\\" + source.Name())
		file, e := xlsx.OpenFile(path + "\\" + source.Name())
		if e != nil {
			tk.Println(e)
			os.Exit(0)
		}
		sheet := file.Sheet[Plant]
		Year := 0
		SequenceData := 0
		for _, row := range sheet.Rows {
			firstCell := ""
			temp := ""
			if len(row.Cells) > 0 {
				firstCell, _ = row.Cells[0].String()
			}
			if strings.Contains(strings.ToLower(firstCell), "for") {
				temp = strings.Replace(strings.Replace(firstCell, " ", "", -1), "for", "", -1)
				Year, e = strconv.Atoi(temp)
				if e != nil {
					Year, _ = strconv.Atoi(temp[(len(temp) - 4):len(temp)])
				}
			}
			if firstCell == "Unit No" {
				SequenceData++
			}
			if len(row.Cells) > 0 && strings.Trim(strings.ToLower(firstCell), " ") != "Unit No" && strings.Trim(strings.ToLower(firstCell), " ") != "" && (strings.Contains(firstCell, "GT") || strings.Contains(firstCell, "ST")) {
				data := new(PerformanceFactors)
				if Plant == "Qurayyah" && SequenceData%2 == 0 {
					data.Plant = Plant + " CC"
				} else {
					data.Plant = Plant
				}
				data.Year = Year

				data.Unit, _ = row.Cells[0].String()
				data.GSHR, _ = row.Cells[1].Float()
				data.NSHR, _ = row.Cells[2].Float()
				data.GTEF, _ = row.Cells[3].Float()
				data.NTEF, _ = row.Cells[4].Float()
				data.GCF, _ = row.Cells[5].Float()
				data.NCF, _ = row.Cells[6].Float()
				data.SRF, _ = row.Cells[7].Float()
				data.ORF, _ = row.Cells[8].Float()
				data.ART, _ = row.Cells[9].Float()
				data.EAF, _ = row.Cells[10].Float()
				data.EFOF, _ = row.Cells[11].Float()
				data.EUOF, _ = row.Cells[12].Float()
				_, e := ctx.InsertOut(data)
				if e != nil {
					tk.Println("ERR on file :", source.Name(), " | ROW :", idx)
					tk.Println(e)
				}

				spp := new(StartupPaymentAndPenalty)

				spp.Plant = data.Plant
				spp.Year = data.Year

				if strings.Contains(data.Unit, "ST") {
					spp.StartupPayment = 18000
					spp.Penalty = 9000
				} else {
					tempUnit, _ := strconv.Atoi(strings.Replace(data.Unit, "GT", "", -1))
					if tempUnit <= 16 {
						spp.StartupPayment = 18000
						spp.Penalty = 9000
					} else if tempUnit <= 24 {
						spp.StartupPayment = 4000
						spp.Penalty = 2000
					} else if tempUnit <= 56 {
						spp.StartupPayment = 6000
						spp.Penalty = 3000
					}
				}
				_, e = ctx.InsertOut(spp)
				if e != nil {
					tk.Println("ERR on file :", source.Name(), " | ROW :", idx)
					tk.Println(e)
				}
			}
		}
	}
	tk.Println("Perfromance Factors from Excel File : COMPLETE")
}
