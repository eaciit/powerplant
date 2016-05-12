package controllers

import (
	m "github.com/eaciit/powerplant/sec/consoleapp/models"
	// "github.com/eaciit/dbox"
	tk "github.com/eaciit/toolkit"
	// "gopkg.in/mgo.v2/bson"
)

type Sample struct {
	*BaseController
}

func (c *Sample) InsertSampleData() {
	tk.Println("Starting Insert sample data..")
	availability := m.Availability{}
	availability.ConvertMGOToSQLServer(c.MongoCtx, c.SqlCtx)
	tk.Println("Process Complete")
}

func (c *Sample) RemoveSampleData() {
	// Removing all data.
	// d := new([]m.Availability)
	// // d.Id = 526
	// csr, _ := c.SqlCtx.Connection.NewQuery().From("Availability").Where(dbox.Eq("Id", 526)).Cursor(nil)
	// defer csr.Close()
	// e := csr.Fetch(d, 1, false)
	// if e != nil {
	// 	tk.Errorf("Error : %s \n", e.Error())
	// }
	// tk.Println(d)
	// x := d[0]
	// tk.Println(x)
	// tk.Println(d)

	// _ = c.SqlCtx.Delete()
	// defer csr.Close()
	data, e := c.GetById(new(m.Availability), 526)
	tk.Println(data)
	// e = c.SqlCtx.Delete(&zxc)
	// tk.Println(zxc)
	if e != nil {
		tk.Errorf("Unable to remove: %s \n", e.Error())
	}
	// tk.Println(d)
}
