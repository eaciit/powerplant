package models

import (
	"github.com/eaciit/orm"
	tk "github.com/eaciit/toolkit"
	// "gopkg.in/mgo.v2/bson"
)

type Availability struct {
	orm.ModelBase `bson:"-",json:"-"`
	// Id            int     `bson:"-",json:"_id"`
	PowerPlant string  ` bson:"PowerPlant" , json:"PowerPlant" `
	Turbine    string  ` bson:"Turbine" , json:"Turbine" `
	PrctWUF    float64 ` bson:"PrctWUF" , json:"PrctWUF" `
	PrctWAF    float64 ` bson:"PrctWAF" , json:"PrctWAF" `
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

func (m *Availability) GetData(ID int, SqlCtx *orm.DataContext) (interface{}, error) {
	csr, e := SqlCtx.Connection.NewQuery().Command("procedure", tk.M{}.
		Set("name", "TEST").
		Set("parms", tk.M{}.Set("@ID", ID))).
		Cursor(nil)
	result := []Availability{}
	e = csr.Fetch(&result, 0, false)
	defer csr.Close()
	if e != nil {
		tk.Errorf("Unable to get availability data: %s \n", e.Error())
		return nil, e
	}
	return result, nil
}

// csr, e := c.NewQuery().Command("procedure", toolkit.M{}.
// 	Set("name", "staticproc").
// 	Set("parms", toolkit.M{}.Set("", ""))).
// 	Cursor(nil)

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

// --ALTER PROCEDURE TEST
// --@ID INT
// --AS
// --BEGIN
// --	DECLARE @PowerPlant VARCHAR(50)
// --	DECLARE csr CURSOR FOR
// --	SELECT PowerPlant FROM Availability
// --	OPEN csr
// --	FETCH NEXT FROM csr INTO @PowerPlant

// --	WHILE (@@FETCH_STATUS <> -1)
// --	BEGIN
// --	 PRINT CAST(@PowerPlant AS VARCHAR(50))
// --	   FETCH NEXT FROM csr INTO @PowerPlant
// --	END
// --	CLOSE csr
// --	DEALLOCATE csr
// --END
// --GO

// ALTER PROCEDURE TEST @ID INT AS
// BEGIN
// 	SELECT * FROM Availability WHERE _id = @ID
// END
// GO

// EXEC TEST  @ID = 1
