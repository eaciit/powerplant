package controllers

import (
	// "github.com/eaciit/crowd"
	// "github.com/eaciit/dbox"
	. "github.com/eaciit/powerplant/sec/consoleapp/generator/helpers"
	. "github.com/eaciit/powerplant/sec/library/models"
	tk "github.com/eaciit/toolkit"
	// "strconv"
	"github.com/tealeg/xlsx"
	"os"
	"strings"
	"time"
)

// REPlannedMaintenance
type REPlannedMaintenance struct {
	*BaseController
}

// Generate
func (d *REPlannedMaintenance) Generate(base *BaseController) {
	var (
		folderName string = "Planned Maintenance"
	)
	if base != nil {
		d.BaseController = base
	}
	ctx := d.BaseController.Ctx
	dataSources, path := base.GetDataSource(folderName)
	tk.Println("Generating Planned Maintenance from Excel File..")
	for _, source := range dataSources {
		if strings.Contains(source.Name(), "Planned Preventive") {
			tk.Println(path + "\\" + source.Name())
			file, e := xlsx.OpenFile(path + "\\" + source.Name())
			if e != nil {
				tk.Println(e)
				os.Exit(0)
			}
			sheet := file.Sheets[0]
			firstCell := ""
			temp := ""
			if !strings.Contains(source.Name(), "PP9") && !strings.Contains(source.Name(), "QPP") {
				for _, row := range sheet.Rows {
					if len(row.Cells) > 0 {
						firstCell, _ = row.Cells[0].String()
					}
					if len(row.Cells) > 0 && strings.Trim(strings.ToLower(firstCell), " ") != "equipment" && strings.Trim(strings.ToLower(firstCell), " ") != "" {

						data := new(PlannedMaintenance)
						data.Equipment, _ = row.Cells[0].String()
						data.MaintenancePlan, _ = row.Cells[1].String()
						data.MainWorkCtr, _ = row.Cells[2].String()
						data.MaintenanceItem, _ = row.Cells[3].String()
						data.CostCtr, _ = row.Cells[4].String()
						data.FunctionalLocation, _ = row.Cells[5].String()
						data.Description, _ = row.Cells[6].String()
						data.MaintenanceItemText, _ = row.Cells[7].String()
						data.CreatedBy, _ = row.Cells[8].String()
						temp, _ = row.Cells[9].String()
						data.CreatedOn, _ = time.Parse(temp, "02/01/2013")
						data.Priority, _ = row.Cells[10].String()
						data.ABCIndicator, _ = row.Cells[11].String()
						temp, _ = row.Cells[12].String()
						data.ChangedOn, _ = time.Parse(temp, "02/01/2013")
						data.ChangedBy, _ = row.Cells[13].String()
						data.Asset, _ = row.Cells[14].String()
						data.SubNumber, _ = row.Cells[15].String()
						data.WorkCtr, _ = row.Cells[16].String()
						data.OrderType, _ = row.Cells[17].String()
						data.OrderNo, _ = row.Cells[18].String()
						data.Assembl, _ = row.Cells[19].String()
						data.PlantSection, _ = row.Cells[20].String()
						data.PurchaseOrder, _ = row.Cells[21].String()
						data.Item1, _ = row.Cells[22].String()
						data.CompanyCode, _ = row.Cells[23].String()
						data.CycleSetSequence, _ = row.Cells[24].String()
						data.StandingOrder, _ = row.Cells[25].String()
						data.ChangeAuthentication, _ = row.Cells[26].String()
						data.LogicalSystem1, _ = row.Cells[27].String()
						data.LogicalSystem2, _ = row.Cells[28].String()
						data.SortField, _ = row.Cells[29].String()
						data.BusinessArea, _ = row.Cells[30].String()
						data.SettlementOrder, _ = row.Cells[31].String()
						data.MaintenanceActiveType, _ = row.Cells[32].String()
						data.ILOAIndividual, _ = row.Cells[33].String()
						data.LOCACCAssmt, _ = row.Cells[34].String()
						data.SettlementRule, _ = row.Cells[35].String()
						data.PlanningPlant, _ = row.Cells[36].String()
						data.SalesDocument, _ = row.Cells[37].String()
						data.Item2, _ = row.Cells[38].String()
						data.COArea, _ = row.Cells[39].String()
						data.Language, _ = row.Cells[40].String()
						data.LastOrder, _ = row.Cells[41].String()
						data.LTextIndicator, _ = row.Cells[42].String()
						data.Client, _ = row.Cells[43].String()
						data.MPlanCategory, _ = row.Cells[44].String()
						data.Room, _ = row.Cells[45].String()
						data.ObjectList = CheckNumberValue(row.Cells[46].Float())
						data.GroupCounter, _ = row.Cells[47].String()
						data.Grup, _ = row.Cells[48].String()
						data.TaskListType, _ = row.Cells[49].String()
						data.FLDescription, _ = row.Cells[50].String()
						data.WBSElement, _ = row.Cells[51].String()
						data.NotificationType, _ = row.Cells[52].String()
						data.SerialNumber, _ = row.Cells[53].String()
						data.Material, _ = row.Cells[54].String()
						data.Division, _ = row.Cells[55].String()
						data.Status, _ = row.Cells[56].String()
						data.Location, _ = row.Cells[57].String()
						data.MaintenancePlant, _ = row.Cells[58].String()
						data.SalesOrg, _ = row.Cells[59].String()
						data.DistributionChannel, _ = row.Cells[60].String()
						data.PlannerGroup, _ = row.Cells[61].String()
						data.ItemNumber, _ = row.Cells[62].String()
						data.Strategy, _ = row.Cells[63].String()

						_, e := ctx.InsertOut(data)
						if e != nil {
							tk.Println("ERR on file :", source.Name(), " | ROW :", idx)
							tk.Println(e)
						}
					}
				}
			} else if strings.Contains(source.Name(), "QPP") {
				for _, row := range sheet.Rows {
					if len(row.Cells) > 0 {
						firstCell, _ = row.Cells[0].String()
					}
					if len(row.Cells) > 0 && strings.Trim(strings.ToLower(firstCell), " ") != "equipment" && strings.Trim(strings.ToLower(firstCell), " ") != "" {

						data := new(PlannedMaintenance)
						data.Equipment, _ = row.Cells[0].String()
						data.MaintenancePlan, _ = row.Cells[1].String()
						data.MainWorkCtr, _ = row.Cells[2].String()
						data.MaintenanceItem, _ = row.Cells[3].String()
						data.CostCtr, _ = row.Cells[4].String()
						data.FunctionalLocation, _ = row.Cells[5].String()
						data.Description, _ = row.Cells[6].String()
						data.MaintenanceItemText, _ = row.Cells[7].String()
						data.CreatedBy, _ = row.Cells[8].String()
						temp, _ = row.Cells[9].String()
						data.CreatedOn, _ = time.Parse("01/02/2006", strings.Replace(temp, " ", "", -1))
						data.Priority, _ = row.Cells[10].String()
						data.ABCIndicator, _ = row.Cells[11].String()

						data.Assembl, _ = row.Cells[12].String()
						data.Asset, _ = row.Cells[13].String()
						data.BusinessArea, _ = row.Cells[14].String()
						data.ChangeAuthentication, _ = row.Cells[15].String()
						data.ChangedBy, _ = row.Cells[16].String()
						temp, _ = row.Cells[17].String()
						data.ChangedOn, _ = time.Parse("01/02/2006", strings.Replace(temp, " ", "", -1))
						data.Client, _ = row.Cells[18].String()
						data.CompanyCode, _ = row.Cells[19].String()
						data.COArea, _ = row.Cells[20].String()
						data.CycleSetSequence, _ = row.Cells[21].String()
						data.FLDescription, _ = row.Cells[22].String()
						data.DistributionChannel, _ = row.Cells[23].String()
						data.Division, _ = row.Cells[24].String()
						data.Grup, _ = row.Cells[25].String()
						data.GroupCounter, _ = row.Cells[26].String()
						data.ILOAIndividual, _ = row.Cells[27].String()
						data.ItemNumber, _ = row.Cells[28].String()
						data.Language, _ = row.Cells[29].String()
						data.LastOrder, _ = row.Cells[30].String()
						data.LOCACCAssmt, _ = row.Cells[31].String()
						data.Location, _ = row.Cells[32].String()
						data.LogicalSystem1, _ = row.Cells[33].String()
						data.LogicalSystem2, _ = row.Cells[34].String()
						data.LTextIndicator, _ = row.Cells[35].String()
						data.MPlanCategory, _ = row.Cells[36].String()
						data.MaintenanceActiveType, _ = row.Cells[37].String()
						data.MaintenancePlant, _ = row.Cells[38].String()
						data.Strategy, _ = row.Cells[39].String()
						data.Material, _ = row.Cells[40].String()
						data.NotificationType, _ = row.Cells[41].String()
						data.ObjectList = CheckNumberValue(row.Cells[42].Float())
						data.OrderNo, _ = row.Cells[43].String()
						data.OrderType, _ = row.Cells[44].String()
						data.PlannerGroup, _ = row.Cells[45].String()
						data.PlanningPlant, _ = row.Cells[46].String()

						data.PlantSection, _ = row.Cells[47].String()
						data.PurchaseOrder, _ = row.Cells[48].String()
						data.Item1, _ = row.Cells[49].String()
						data.Room, _ = row.Cells[50].String()
						data.SalesDocument, _ = row.Cells[51].String()
						data.Item2, _ = row.Cells[52].String()
						data.SalesOrg, _ = row.Cells[53].String()
						data.SerialNumber, _ = row.Cells[54].String()
						data.SettlementOrder, _ = row.Cells[55].String()
						data.SettlementRule, _ = row.Cells[56].String()
						data.SortField, _ = row.Cells[57].String()
						data.StandingOrder, _ = row.Cells[58].String()
						data.Status, _ = row.Cells[59].String()
						data.SubNumber, _ = row.Cells[60].String()
						data.TaskListType, _ = row.Cells[61].String()
						data.WBSElement, _ = row.Cells[62].String()
						data.WorkCtr, _ = row.Cells[63].String()

						_, e := ctx.InsertOut(data)
						if e != nil {
							tk.Println("ERR on file :", source.Name(), " | ROW :", idx)
							tk.Println(e)
						}
					}
				}
			} else if strings.Contains(source.Name(), "PP9") {
				for _, row := range sheet.Rows {
					if len(row.Cells) > 0 {
						firstCell, _ = row.Cells[0].String()
					}
					if len(row.Cells) > 0 && strings.Trim(strings.ToLower(firstCell), " ") != "equipment" && strings.Trim(strings.ToLower(firstCell), " ") != "" {

						data := new(PlannedMaintenance)
						data.ABCIndicator, _ = row.Cells[0].String()
						data.Assembl, _ = row.Cells[1].String()
						data.Asset, _ = row.Cells[2].String()
						data.BusinessArea, _ = row.Cells[3].String()
						data.ChangeAuthentication, _ = row.Cells[4].String()
						data.ChangedBy, _ = row.Cells[5].String()
						temp, _ = row.Cells[6].String()
						data.ChangedOn, _ = time.Parse("01/02/2006", strings.Replace(temp, " ", "", -1))
						data.Client, _ = row.Cells[7].String()
						data.CompanyCode, _ = row.Cells[8].String()
						data.COArea, _ = row.Cells[9].String()
						data.CycleSetSequence, _ = row.Cells[10].String()
						data.Description, _ = row.Cells[11].String()

						data.DistributionChannel, _ = row.Cells[12].String()
						data.Division, _ = row.Cells[13].String()
						data.Grup, _ = row.Cells[14].String()
						data.GroupCounter, _ = row.Cells[15].String()
						data.ILOAIndividual, _ = row.Cells[16].String()
						data.ItemNumber, _ = row.Cells[17].String()
						data.Language, _ = row.Cells[18].String()
						data.LastOrder, _ = row.Cells[19].String()
						data.LOCACCAssmt, _ = row.Cells[20].String()
						data.Location, _ = row.Cells[21].String()
						data.LogicalSystem1, _ = row.Cells[22].String()
						data.LogicalSystem2, _ = row.Cells[23].String()
						data.LTextIndicator, _ = row.Cells[24].String()
						data.MPlanCategory, _ = row.Cells[25].String()
						data.MaintenanceActiveType, _ = row.Cells[26].String()
						data.MaintenancePlant, _ = row.Cells[27].String()
						data.Strategy, _ = row.Cells[28].String()
						data.Material, _ = row.Cells[29].String()
						data.NotificationType, _ = row.Cells[30].String()
						data.ObjectList = CheckNumberValue(row.Cells[31].Float())
						data.OrderNo, _ = row.Cells[32].String()
						data.PlannerGroup, _ = row.Cells[33].String()
						data.PlanningPlant, _ = row.Cells[34].String()
						data.PlantSection, _ = row.Cells[35].String()
						data.PurchaseOrder, _ = row.Cells[36].String()
						data.Item1, _ = row.Cells[37].String()
						data.Room, _ = row.Cells[38].String()
						data.SalesDocument, _ = row.Cells[39].String()
						data.Item2, _ = row.Cells[40].String()
						data.SalesOrg, _ = row.Cells[41].String()
						data.SerialNumber, _ = row.Cells[42].String()
						data.SettlementOrder, _ = row.Cells[43].String()
						data.SettlementRule, _ = row.Cells[44].String()
						data.SortField, _ = row.Cells[45].String()
						data.StandingOrder, _ = row.Cells[46].String()

						data.StandingOrder, _ = row.Cells[47].String()
						data.Status, _ = row.Cells[48].String()
						data.SubNumber, _ = row.Cells[49].String()
						data.TaskListType, _ = row.Cells[50].String()
						data.WBSElement, _ = row.Cells[51].String()
						data.WorkCtr, _ = row.Cells[52].String()
						data.Equipment, _ = row.Cells[53].String()
						data.MaintenancePlan, _ = row.Cells[54].String()
						data.MainWorkCtr, _ = row.Cells[55].String()
						data.MaintenanceItem, _ = row.Cells[56].String()
						data.CostCtr, _ = row.Cells[57].String()
						data.FunctionalLocation, _ = row.Cells[58].String()
						data.FLDescription, _ = row.Cells[59].String()
						data.MaintenanceItemText, _ = row.Cells[60].String()
						data.CreatedBy, _ = row.Cells[61].String()
						temp, _ = row.Cells[62].String()
						data.CreatedOn, _ = time.Parse("01/02/2006", strings.Replace(temp, " ", "", -1))
						data.Priority, _ = row.Cells[63].String()
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
	tk.Println("Planned Maintenance from Excel File : COMPLETE")
}
