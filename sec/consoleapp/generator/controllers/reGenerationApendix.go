package controllers

import (
	// "github.com/eaciit/crowd"
	// "github.com/eaciit/dbox"
	. "github.com/eaciit/powerplant/sec/consoleapp/generator/helpers"
	. "github.com/eaciit/powerplant/sec/library/models"
	tk "github.com/eaciit/toolkit"
	"github.com/tealeg/xlsx"
	"os"
	// "strconv"
	"strings"
	// "time"
)

// REGenerationApendix
type REGenerationApendix struct {
	*BaseController
}

// Generate
func (d *REGenerationApendix) Generate(base *BaseController) {
	var (
		folderName string = "Generation Apendix"
	)
	if base != nil {
		d.BaseController = base
	}
	ctx := d.BaseController.Ctx
	dataSources, path := base.GetDataSource(folderName)
	tk.Println("Generating Generation Apendix from Excel File..")
	for _, source := range dataSources {
		tk.Println(path + "\\" + source.Name())
		file, e := xlsx.OpenFile(path + "\\" + source.Name())
		if e != nil {
			tk.Println(e)
			os.Exit(0)
		}
		sheet := file.Sheet["Central"]

		for _, row := range sheet.Rows {
			firstCell := ""
			if len(row.Cells) > 1 {
				firstCell, _ = row.Cells[1].String()
			}
			if len(row.Cells) > 0 && strings.Trim(strings.ToLower(firstCell), " ") != "power plant" && strings.Trim(strings.ToLower(firstCell), " ") != "" {

				data := new(GenerationAppendix)
				data.Plant, _ = row.Cells[1].String()
				data.Type, _ = row.Cells[2].String()
				data.Units, _ = row.Cells[3].Int()
				data.ContractedCapacity = CheckNumberValue(row.Cells[4].Float())
				data.CCR = CheckNumberValue(row.Cells[5].Float())
				data.FOMR = CheckNumberValue(row.Cells[6].Float())
				data.VOMR = CheckNumberValue(row.Cells[7].Float())
				data.AGP = CheckNumberValue(row.Cells[8].Float())
				data.LCSummer = CheckNumberValue(row.Cells[9].Float())
				data.LCWinter = CheckNumberValue(row.Cells[10].Float())
				data.LCTotal = CheckNumberValue(row.Cells[11].Float())
				data.Startup = CheckNumberValue(row.Cells[12].Float())
				data.Deduct = CheckNumberValue(row.Cells[13].Float())

				_, e := ctx.InsertOut(data)
				if e != nil {
					tk.Println("ERR on file :", source.Name(), " | ROW :", idx)
					tk.Println(e)
				}
			}
		}
	}
	tk.Println("Generation Apendix from Excel File : COMPLETE")
}
