package controllers

import (
	"github.com/eaciit/dbox"
	"github.com/eaciit/orm"
	"log"
	"os"
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

func (b *BaseController) Update(m orm.IModel, id interface{}, data interface{}, column_name ...string) error {
	column_id := "Id"
	if column_name != nil && len(column_name) > 0 {
		column_id = column_name[0]
	}
	e := b.SqlCtx.Connection.NewQuery().From(m.(orm.IModel).TableName()).Where(dbox.Eq(column_id, id)).Update().Exec(nil)
	// tk.M{"data", data}
	if e != nil {
		return e
	}
	return nil
}

func (b *BaseController) Turncate(m orm.IModel) error {
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
