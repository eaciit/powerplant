package models

import (
	"github.com/eaciit/orm"
	tk "github.com/eaciit/toolkit"
	// "gopkg.in/mgo.v2/bson"
)

type Availability struct {
	orm.ModelBase `bson:"-",json:"-"`
	Id            int     `bson:"Id",json:"Id"`
	PowerPlant    string  ` bson:"PowerPlant" , json:"PowerPlant" `
	Turbine       string  ` bson:"Turbine" , json:"Turbine" `
	PrctWUF       float64 ` bson:"PrctWUF" , json:"PrctWUF" `
	PrctWAF       float64 ` bson:"PrctWAF" , json:"PrctWAF" `
}

func (m *Availability) ConvertMGOToSQLServer(MongoCtx *orm.DataContext, SqlCtx *orm.DataContext) {
	csr, e := MongoCtx.Connection.NewQuery().From(m.TableName()).Cursor(nil)
	result := []Availability{}
	e = csr.Fetch(&result, 0, false)
	defer csr.Close()
	if e != nil {
		tk.Errorf("Unable to save: %s \n", e.Error())
		return
	}
	query := SqlCtx.Connection.NewQuery().SetConfig("multiexec", true).From(m.TableName()).Save()
	for _, i := range result {
		e = query.Exec(tk.M{"data": i})
		if e != nil {
			tk.Errorf("Unable to save: %s \n", e.Error())
			return
		}
	}
}
func (m *Availability) TableName() string {
	return "Availability"
}

// CREATE TABLE "Availability"(
// 	id int NOT NULL  IDENTITY(1,1) PRIMARY KEY,
// 	PowerPlant VARCHAR(50) NOT NULL,
// 	Turbine VARCHAR(25) NOT NULL,
// 	PrctWUF FLOAT NOT NULL,
// 	PrctWAF FLOAT NOT NULL,
// )
