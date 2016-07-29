package controllers

import (
	"github.com/eaciit/knot/knot.v1"
	. "github.com/eaciit/powerplant/sec/library/models"
	tk "github.com/eaciit/toolkit"
)

type LoginController struct {
	*BaseController
}

func (c *LoginController) Default(k *knot.WebContext) interface{} {
	c.LoadPartial(k)
	k.Config.NoLog = true
	k.Config.OutputType = knot.OutputTemplate
	k.Config.LayoutTemplate = ""
	return ""
}

func (c *LoginController) Do(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	k.Config.NoLog = true

	msg := ""
	p := struct {
		UserName   string
		Password   string
		RememberMe bool
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		c.WriteLog(e)
	}
	User := new(UserModel)
	User, Found := User.LoginDo(c.Ctx, p.UserName, p.Password)

	if User.Enable == false {
		msg = "Your User is Not Active, Please Contact your Admin!"
	} else if Found {
		k.SetSession("userid", User.UserName)
		k.SetSession("username", User.FullName)
		k.SetSession("userrole", "ADMIN")
	} else {
		msg = "Your User ID or password are not matches!"
	}
	return tk.M{}.Set("IsLogged", Found).Set("Message", msg).Set("success", true)
}

func (c *LoginController) Logout(k *knot.WebContext) interface{} {
	k.Config.NoLog = true
	k.Config.OutputType = knot.OutputNone
	k.SetSession("userid", nil)
	k.SetSession("username", nil)
	k.SetSession("userrole", nil)

	c.Redirect(k, "login", "default")

	return true
}
