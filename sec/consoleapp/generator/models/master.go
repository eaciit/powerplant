package models

import (
	"github.com/eaciit/dbox"
	"github.com/eaciit/orm"
	. "github.com/eaciit/powerplant/sec/library/models"
	tk "github.com/eaciit/toolkit"
)

type Master struct {
}

func (m *Master) GeneratePlantMaster(ctx *orm.DataContext) error {
	tk.Println("Generating Plant Master..")
	c := ctx.Connection
	query := []*dbox.Filter{}
	query = append(query, dbox.Eq("LEN(FunctionalLocationCode)", 4))

	FunctionalLocationList := []FunctionalLocation{}
	csr, e := c.NewQuery().Select("FunctionalLocationCode", "Description").From(new(FunctionalLocation).TableName()).Where(query...).Cursor(nil)
	if e != nil {
		return e
	}
	e = csr.Fetch(&FunctionalLocationList, 0, false)
	csr.Close()
	if e != nil {
		return e
	}

	PowerPlantInfoList := []PowerPlantInfo{}
	csr, e = c.NewQuery().From(new(PowerPlantInfo).TableName()).Cursor(nil)
	if e != nil {
		return e
	}
	e = csr.Fetch(&PowerPlantInfoList, 0, false)
	csr.Close()
	if e != nil {
		return e
	}

	for _, plant := range FunctionalLocationList {
		d := tk.M{}
		e := tk.StructToM(plant, &d)
		if e != nil {
			tk.Println(e)
		}
		tk.Println(d)
		// Generate data as per code
	}
	if e != nil {
		return e
	}
	return nil
}
