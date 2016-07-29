package controllers

import (
	"time"

	"github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	. "github.com/eaciit/powerplant/sec/library/models"
	. "github.com/eaciit/powerplant/sec/webapp/helpers"
	"github.com/eaciit/toolkit"
)

type AdministrationController struct {
	*BaseController
}

func (c *AdministrationController) UserManagement(k *knot.WebContext) interface{} {
	if k.Session("userid") == nil {
		c.Redirect(k, "login", "default")
	}
	c.LoadPartial(k, "administration/usermanagement.html")
	k.Config.OutputType = knot.OutputTemplate

	infos := PageInfo{}
	infos.PageId = "UserManagement"
	infos.PageTitle = "User Management"
	infos.Breadcrumbs = make(map[string]string, 0)

	return infos
}

func (c *AdministrationController) GetUserList(k *knot.WebContext) interface{} {
	csr, e := c.Ctx.Find(new(UserModel), nil)
	defer csr.Close()
	UserList := make([]UserModel, 0)
	e = csr.Fetch(&UserList, 0, false)
	if e != nil {
		return e.Error()
	}
	return ResultInfo(UserList, e)
}

func (c *AdministrationController) SaveUser(k *knot.WebContext) interface{} {
	d := struct {
		Id       string
		UserName string
		FullName string
		Password string
		Email    string
		Enable   bool
		ADUser   bool
	}{}
	e := k.GetPayload(&d)
	if e != nil {
		return ResultInfo(nil, e)
	}

	data := new(UserModel)
	data.UserName = d.UserName
	data.FullName = d.FullName
	data.PasswordHash = GetMD5Hash(d.Password)
	data.Email = d.Email
	data.Enable = d.Enable
	data.ADUser = d.ADUser
	data.SecurityStamp = time.Now().UTC()
	data.ConfirmedAtUtc = time.Now().UTC()

	e = c.Ctx.Connection.NewQuery().From(data.TableName()).Where(dbox.Eq("Email", d.Id)).Delete().Exec(nil)
	if e != nil {
		return ResultInfo(nil, e)
	}
	_, e = c.Ctx.InsertOut(data)
	if e != nil {
		return ResultInfo(nil, e)
	}
	return ResultInfo(data, e)
}

func (c *AdministrationController) DeactivateUser(k *knot.WebContext) interface{} {
	d := struct {
		Email string
	}{}

	e := k.GetPayload(&d)

	if e != nil {
		return ResultInfo(nil, e)
	}

	/*csr, e := c.Ctx.Connection.NewQuery().From(new(UserModel).TableName()).Where(dbox.Eq("Email", d.Email)).Cursor(nil)
	result := new(UserModel)
	e = csr.Fetch(&result, 1, false)
	defer csr.Close()

	e = c.Ctx.Connection.NewQuery().From(new(UserModel).TableName()).Where(dbox.Eq("Email", d.Email)).Delete().Exec(nil)
	if e != nil {
		return ResultInfo(nil, e)
	}

	data := result
	data.Enable = false

	e = c.Ctx.Save(data)*/

	e = c.Ctx.Connection.NewQuery().
		Update().
		From(new(UserModel).TableName()).
		Where(dbox.Contains("email", d.Email)).
		Exec(toolkit.M{}.Set("data", toolkit.M{}.Set("enable", false)))

	if e != nil {
		return ResultInfo(nil, e)
	}

	return ResultInfo(nil, e)
}

func (c *AdministrationController) ReactivateUser(k *knot.WebContext) interface{} {
	d := struct {
		Email string
	}{}
	e := k.GetPayload(&d)

	if e != nil {
		return ResultInfo(nil, e)
	}

	/*csr, e := c.Ctx.Connection.NewQuery().From(new(UserModel).TableName()).Where(dbox.Eq("Email", d.Email)).Cursor(nil)
	result := new(UserModel)
	e = csr.Fetch(&result, 1, false)
	defer csr.Close()

	e = c.Ctx.Connection.NewQuery().From(new(UserModel).TableName()).Where(dbox.Eq("Email", d.Email)).Delete().Exec(nil)
	if e != nil {
		return ResultInfo(nil, e)
	}

	data := result
	data.Enable = true
	e = c.Ctx.Save(data)*/

	e = c.Ctx.Connection.NewQuery().
		Update().
		From(new(UserModel).TableName()).
		Where(dbox.Contains("email", d.Email)).
		Exec(toolkit.M{}.Set("data", toolkit.M{}.Set("enable", true)))

	if e != nil {
		return ResultInfo(nil, e)
	}

	return ResultInfo(nil, e)
}
