package models

import (
	"sync"

	"github.com/eaciit/orm"
	tk "github.com/eaciit/toolkit"
	// "gopkg.in/mgo.v2/bson"
)

type Availability struct {
	sync.RWMutex
	orm.ModelBase `bson:"-" json:"-"`
	// Id            int64   `bson:"id" json:"id"`
	Plant   string  ` bson:"PowerPlant" json:"Plant" `
	Turbine string  ` bson:"Turbine" json:"Turbine" `
	PrctWUF float64 ` bson:"PrctWUF" json:"PrctWUF" `
	PrctWAF float64 ` bson:"PrctWAF" json:"PrctWAF" `
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
// 	SELECT * FROM Availability WHERE id = @ID
// END
// GO

// EXEC TEST  @ID = 1
