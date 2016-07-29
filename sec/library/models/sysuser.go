package models

import (
	"github.com/eaciit/orm"
	// "gopkg.in/mgo.v2/bson"
	"github.com/eaciit/dbox"
	. "github.com/eaciit/powerplant/sec/webapp/helpers"
	"sync"
	"time"
)

type UserModel struct {
	sync.RWMutex
	orm.ModelBase `bson:"-",json:"-"`
	// Id                int       `json:Id`
	UserName          string    `json:"UserName"`
	FullName          string    `json:"FullName"`
	Enable            bool      `json:"Enable"`
	ADUser            bool      `json:"ADUser"`
	SecurityStamp     time.Time `json:"SecurityStamp"`
	ConfirmedAtUtc    time.Time `json:"ConfirmedAtUtc"`
	Email             string    `json:"Email"`
	PasswordHash      string    `json:"PasswordHash"`
	HasChangePassword bool      `json:"HasChangePassword"`
	SecretToken       string    `json:"SecretToken`
}

func (u *UserModel) TableName() string {
	return "Users"
}

func (u *UserModel) LoginDo(ctx *orm.DataContext, username string, password string) (*UserModel, bool) {
	PasswordHash := GetMD5Hash(password)
	found := false
	query := []*dbox.Filter{}
	query = append(query, dbox.In("UserName", username))
	query = append(query, dbox.In("PasswordHash", PasswordHash))
	csr, e := ctx.Connection.NewQuery().From(u.TableName()).Where(query...).Cursor(nil)
	if e != nil {
		return nil, found
	}
	result := []*UserModel{}
	csr.Fetch(&result, 0, false)
	csr.Close()
	if len(result) > 0 {
		found = true
		for _, i := range result {
			u = i
		}
	}
	return u, found

}
