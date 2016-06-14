package controllers

import (
	m "github.com/eaciit/powerplant/sec/library/models"
	// "github.com/eaciit/dbox"
	tk "github.com/eaciit/toolkit"
	// "gopkg.in/mgo.v2/bson"
)

type Sample struct {
	*BaseController
}

func (c *Sample) InsertSampleData() {
	tk.Println("Starting Insert sample data..")
	e := c.ConvertMGOToSQLServer(new(m.Availability))
	if e != nil {
		tk.Println(e)
	}
	tk.Println("Process Complete")
}

func (c *Sample) GetSampleData() {
	tk.Println("Getting sample data..")
	availability := m.Availability{}
	data, _ := availability.GetData(176, c.SqlCtx)
	tk.Println(data)
	tk.Println("Process Complete")
}

func (c *Sample) UpdateSampleData() {
	id := 1
	data := new(m.Availability)
	e := c.GetById(data, id, "id")
	data.PrctWUF = 1.2
	c.Update(data, id, "id")
	if e != nil {
		tk.Errorf("Unable to remove: %s \n", e.Error())
	}
}

func (c *Sample) RemoveSampleData() {
	data := new(m.Availability)
	e := c.GetById(data, 526, "id")
	tk.Println(data)
	e = c.SqlCtx.Delete(data)
	if e != nil {
		tk.Errorf("Unable to remove: %s \n", e.Error())
	}
}
