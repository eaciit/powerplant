package controllers

import (
	"github.com/eaciit/knot/knot.v1"
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

	isLogged := false
	msg := ""
	p := struct {
		UserName   string
		Password   string
		RememberMe bool
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		c.WriteLog(e)
		// msg = "Error: " + e.Error()
	}

	tk.Println(p)
	// // temporary
	// if p.UserName == "eaciit" && p.Password == "Password.1" {
	k.SetSession("userid", p.UserName)
	k.SetSession("username", "EACIIT Administrator")
	k.SetSession("userrole", "ADMIN")
	isLogged = true
	// 	msg = ""
	// } else {
	// 	msg = "Your User ID or password are not matches!"
	// }
	return tk.M{}.Set("IsLogged", isLogged).Set("Message", msg).Set("success", true)
}

func (c *LoginController) DoLogout(k *knot.WebContext) interface{} {
	k.Config.NoLog = true
	k.Config.OutputType = knot.OutputNone
	k.SetSession("userid", nil)
	k.SetSession("username", nil)
	k.SetSession("userrole", nil)

	c.Redirect(k, "login", "default")

	return true
}
