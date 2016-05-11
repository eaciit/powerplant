package webext

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/mongo"
	"github.com/eaciit/knot/knot.v1"
	"github.com/eaciit/orm"
	. "github.com/eaciit/powerplant/sec/webapp/controllers"
)

var (
	wd = func() string {
		d, _ := os.Getwd()
		return d + "/"
	}()
)

func init() {
	conn, err := PrepareConnection()
	if err != nil {
		fmt.Println(err)
	}
	ctx := orm.New(conn)

	baseCont := new(BaseController)
	baseCont.Ctx = ctx

	app := knot.NewApp("sec")
	app.ViewsPath = wd + "views/"

	// register controllers
	app.Register(&LoginController{baseCont})
	app.Register(&DashboardController{baseCont})
	app.Register(&DataBrowserVEController{baseCont})
	app.Register(&ValueEquationController{baseCont})
	app.Register(&ValueEquationComparisonController{baseCont})
	app.Register(&HistoricalValueEquationController{baseCont})
	app.Register(&HypothesisController{baseCont})
	app.Register(&ScenarioSimulation{baseCont})
	app.Register(&UploadDataController{baseCont})
	/*app.Register(&InitController{baseCont})
	app.Register(&OrganizationController{baseCont})
	app.Register(&InventoryController{baseCont})
	app.Register(&UomController{baseCont})*/

	app.Static("static", wd+"assets")
	app.LayoutTemplate = "shared/layout.html"
	app.DefaultOutputType = knot.OutputJson
	knot.RegisterApp(app)
}

func PrepareConnection() (dbox.IConnection, error) {
	config := ReadConfig()
	ci := &dbox.ConnectionInfo{config["host"], config["database"], config["username"], config["password"], nil}
	c, e := dbox.NewConnection("mongo", ci)

	if e != nil {
		return nil, e
	}

	e = c.Connect()
	if e != nil {
		return nil, e
	}

	return c, nil
}

func ReadConfig() map[string]string {
	ret := make(map[string]string)
	fmt.Println(wd, "conf/app.conf")
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
