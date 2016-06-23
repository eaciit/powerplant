package controllers

import (
	// "time"
	"github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	. "github.com/eaciit/powerplant/sec/webapp/models"
	tk "github.com/eaciit/toolkit"
	"math/rand"
	// "strings"
	. "github.com/eaciit/powerplant/sec/webapp/helpers"
	"gopkg.in/gomail.v2"
)

type AccountController struct {
	*BaseController
}

func (c *AccountController) ForgotPassword(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	var e error
	result := ""
	d := struct {
		UserEmail string
	}{}

	e = k.GetPayload(&d)

	csr, e := c.Ctx.Find(new(UserModel), tk.M{}.Set("where", dbox.Eq("Email", d.UserEmail)))
	defer csr.Close()
	Users := []*UserModel{}
	e = csr.Fetch(&Users, 0, false)
	if e != nil {
		return ResultInfo(result, e)
	}

	if len(Users) > 0 {
		user := Users[0]
		var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

		b := make([]rune, 5)
		for i := range b {
			b[i] = letterRunes[rand.Intn(len(letterRunes))]
		}
		newPassword := string(b)
		user.PasswordHash = GetMD5Hash(newPassword)
		e = c.Ctx.Delete(user)
		if e != nil {
			return ResultInfo(result, e)
		}
		e = c.Ctx.Save(user)
		if e != nil {
			return ResultInfo(result, e)
		}
		conf := gomail.NewDialer("smtp.office365.com", 587, "admin.support@eaciit.com", "B920Support")
		s, e := conf.Dial()
		if e != nil {

			return ResultInfo(result, e)
		}

		mailsubj := tk.Sprintf("%v", "[NOREPLY] Forgot Password SEC Apps")
		mailmsg := tk.Sprintf("Dear %v \n\n This is your new password : %v \n\nThis mail has been sent from forgot password feature and if you did not initiate this change, please contact your System Administrator", user.FullName, newPassword)

		m := gomail.NewMessage()

		m.SetHeader("From", "admin.support@eaciit.com")
		m.SetHeader("To", "ainur.rochman@eaciit.com")
		m.SetHeader("Subject", mailsubj)
		m.SetBody("text/html", mailmsg)

		e = gomail.Send(s, m)
		m.Reset()
		if e != nil {
			return ResultInfo(result, e)
		}
		result = "OK"
	} else {
		result = "NOK"
	}

	return ResultInfo(result, e)
}
