package controllers

import (
	// "github.com/eaciit/crowd"
	// "github.com/eaciit/dbox"
	// . "github.com/eaciit/powerplant/sec/consoleapp/generator/helpers"
	. "github.com/eaciit/powerplant/sec/library/models"
	tk "github.com/eaciit/toolkit"
	"github.com/tealeg/xlsx"
	"os"
	// "strconv"
	"strings"
	// "time"
)

// RENewObjectType
type RENewObjectType struct {
	*BaseController
}

// Generate
func (d *RENewObjectType) Generate(base *BaseController) {
	var (
		folderName string = "New Object Type"
	)
	if base != nil {
		d.BaseController = base
	}
	ctx := d.BaseController.Ctx
	dataSources, path := base.GetDataSource(folderName)
	tk.Println("Generating New Object Type from Excel File..")
	for _, source := range dataSources {
		tk.Println(path + "\\" + source.Name())
		file, e := xlsx.OpenFile(path + "\\" + source.Name())
		if e != nil {
			tk.Println(e)
			os.Exit(0)
		}
		sheet := file.Sheet["Final Object Type Breakdown"]

		for _, row := range sheet.Rows {
			firstCell := ""
			if len(row.Cells) > 1 {
				firstCell, _ = row.Cells[1].String()
			}
			if len(row.Cells) > 0 && strings.Trim(strings.ToLower(firstCell), " ") != "object type" && strings.Trim(strings.ToLower(firstCell), " ") != "" {

				data := new(NewEquipmentType)
				data.EquipmentType, _ = row.Cells[1].String()
				data.EquipmentText, _ = row.Cells[2].String()
				data.NewEquipmentGroup, _ = row.Cells[3].String()
				_, e := ctx.InsertOut(data)
				if e != nil {
					tk.Println("ERR on file :", source.Name(), " | ROW :", idx)
					tk.Println(e)
				}
			}
		}
	}
	tk.Println("New Object Type from Excel File : COMPLETE")
}
