package controllers

import (
	"github.com/eaciit/dbox"
	"github.com/eaciit/orm"
	tk "github.com/eaciit/toolkit"
	"log"
	"os"
	"reflect"
	"time"
)

var (
	wd = func() string {
		d, _ := os.Getwd()
		return d + "/"
	}()
)

type IBaseController interface {
	// not implemented anything yet
}

type BaseController struct {
	base     IBaseController
	MongoCtx *orm.DataContext
	SqlCtx   *orm.DataContext
}

func (b *BaseController) ConvertMGOToSQLServer(m orm.IModel) error {
	tk.Printf("\nConvertMGOToSQLServer: Converting %v \n", m.TableName())
	tk.Println("ConvertMGOToSQLServer: Starting to convert...\n")
	csr, e := b.MongoCtx.Connection.NewQuery().From(m.TableName()).Cursor(nil)
	defer csr.Close()
	if e != nil {
		return e
	}
	result := []tk.M{}
	e = csr.Fetch(&result, 0, false)
	if e != nil {
		return e
	}
	query := b.SqlCtx.Connection.NewQuery().SetConfig("multiexec", true).From(m.TableName()).Save()
	for _, i := range result {
		valueType := reflect.TypeOf(m).Elem()
		for f := 0; f < valueType.NumField(); f++ {
			field := valueType.Field(f)
			bsonField := field.Tag.Get("bson")
			jsonField := field.Tag.Get("json")
			if jsonField != bsonField && field.Name != "RWMutex" && field.Name != "ModelBase" {
				i.Set(field.Name, i.Get(bsonField))
			}
			if field.Type.Name() == "Time" && i.Get(bsonField) == nil {
				i.Set(field.Name, time.Time{})
			}
		}
		e := tk.Serde(i, m, "json")
		e = query.Exec(tk.M{"data": m})
		if e != nil {
			return e
		}
	}
	tk.Println("\nConvertMGOToSQLServer: Finish.")
	return nil
}

func (b *BaseController) GetById(m orm.IModel, id interface{}, column_name ...string) error {
	var e error
	c := b.SqlCtx.Connection
	column_id := "Id"
	if column_name != nil && len(column_name) > 0 {
		column_id = column_name[0]
	}
	csr, e := c.NewQuery().From(m.(orm.IModel).TableName()).Where(dbox.Eq(column_id, id)).Cursor(nil)
	defer csr.Close()
	if e != nil {
		return e
	}
	e = csr.Fetch(m, 1, false)
	if e != nil {
		return e
	}

	return nil
}

func (b *BaseController) Delete(m orm.IModel, id interface{}, column_name ...string) error {
	column_id := "Id"
	if column_name != nil && len(column_name) > 0 {
		column_id = column_name[0]
	}
	e := b.SqlCtx.Connection.NewQuery().From(m.(orm.IModel).TableName()).Where(dbox.Eq(column_id, id)).Delete().Exec(nil)
	if e != nil {
		return e
	}
	return nil
}

func (b *BaseController) Update(m orm.IModel, id interface{}, column_name ...string) error {
	column_id := "Id"
	if column_name != nil && len(column_name) > 0 {
		column_id = column_name[0]
	}
	e := b.SqlCtx.Connection.NewQuery().From(m.(orm.IModel).TableName()).Where(dbox.Eq(column_id, id)).Update().Exec(tk.M{"data": m})
	if e != nil {
		return e
	}
	return nil
}

func (b *BaseController) Truncate(m orm.IModel) error {
	c := b.SqlCtx.Connection
	e := c.NewQuery().From(m.(orm.IModel).TableName()).Delete().Exec(nil)
	if e != nil {
		return e
	}

	return nil
}
func (b *BaseController) CloseDb() {
	if b.MongoCtx != nil {
		b.MongoCtx.Close()
	}
	if b.SqlCtx != nil {
		b.SqlCtx.Close()
	}
}

func (b *BaseController) WriteLog(msg interface{}) {
	log.Printf("%#v\n\r", msg)
	return
}
