package controllers

import (
	"bufio"
	//. "github.com/eaciit/powerplant/sec/models"
	"fmt"

	"github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	"github.com/eaciit/orm"
	// tk "github.com/eaciit/toolkit"
	//"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"os"
	"strings"
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
	base IBaseController
	Ctx  *orm.DataContext
}

type PageInfo struct {
	PageId       string
	PageTitle    string
	SelectedMenu string
	Breadcrumbs  map[string]string
}

func (b *BaseController) ReadConfig() map[string]string {
	ret := make(map[string]string)
	file, err := os.Open(wd + "conf/app.conf")
	if err == nil {
		defer file.Close()

		reader := bufio.NewReader(file)
		for {
			line, _, e := reader.ReadLine()
			if e != nil {
				break
			}

			sval := strings.Split(string(line), "=")
			ret[sval[0]] = sval[1]
		}
	} else {
		fmt.Println(err.Error())
	}

	return ret
}

func (b *BaseController) SetDb(conn dbox.IConnection) error {
	b.CloseDb()
	b.Ctx = orm.New(conn)
	return nil
}

func (b *BaseController) DB() *orm.DataContext {
	return b.Ctx
}

func (b *BaseController) CloseDb() {
	if b.Ctx != nil {
		b.Ctx.Close()
	}
}

func (b *BaseController) LoadBase(k *knot.WebContext) {
	k.Config.NoLog = true
	b.IsAuthenticate(k)
}

func (b *BaseController) IsAuthenticate(k *knot.WebContext) {
	if k.Session("userid") == nil {
		b.Redirect(k, "login", "default")
	}
	return
}

func (b *BaseController) LoadPartial(k *knot.WebContext, tpls ...string) {
	defaultTpls := []string{"shared/processing.html"}
	if len(tpls) > 0 {
		defaultTpls = append(defaultTpls, tpls...)
	}
	k.Config.IncludeFiles = defaultTpls
}

func (b *BaseController) WriteLog(msg interface{}) {
	log.Printf("%#v\n\r", msg)
	return
}

func (b *BaseController) Redirect(k *knot.WebContext, controller string, action string) {
	log.Println("invalid session , redirecting to " + controller + "/" + action)
	http.Redirect(k.Writer, k.Request, "/"+controller+"/"+action, http.StatusTemporaryRedirect)
}

type ResultInformation struct {
	IsError bool
	Message string
	Data    interface{}
}

func ResultInfo(data interface{}, e error) ResultInformation {
	ri := ResultInformation{}
	if e != nil {
		ri.IsError = true
		ri.Message = e.Error()
	}
	ri.Data = data
	return ri
}
