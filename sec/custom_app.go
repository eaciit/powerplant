package main

import (
	"bufio"
	"os"
	"reflect"
	"strings"

	"github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/mongo"
	"github.com/eaciit/orm"
	tk "github.com/eaciit/toolkit"
)

var (
	wd = func() string {
		d, _ := os.Getwd()
		return d + "/"
	}()
	ctx *orm.DataContext
)

func main() {
	tk.Println("Starting the app..")
	mongo, e := PrepareConnection("mongo")
	if e != nil {
		tk.Println(e)
	}
	ctx = orm.New(mongo)
	GetStructFromCollection("ValueEquation_Dashboard")
	defer ctx.Close()
}

func GetStructFromCollection(CollectionName string) {
	tk.Println("\nGetting Struct From Collection [WITHOUT DETAIL].. ")
	tk.Println("Collection Name : ", CollectionName)
	tk.Println("#############################################################")
	csr, e := ctx.Connection.NewQuery().From(CollectionName).Take(1).Cursor(nil)
	defer csr.Close()
	if e != nil {
		tk.Println(e.Error())
	}
	result := []tk.M{}
	e = csr.Fetch(&result, 1, false)
	if e != nil {
		tk.Println(e.Error())
	}
	csr.Close()
	tk.Println("type ", CollectionName, " struct {")
	for _, obj := range result {
		for column, i := range obj {
			datatype := reflect.TypeOf(i)
			// rv := reflect.ValueOf(i)

			if datatype != nil {
				switch datatype.Kind() {
				case reflect.Map:
					for c, x := range i.(map[string]interface{}) {
						if x == nil {
							tk.Println(c, "	unidentified `bson:'", column, ".", c, "' json:'", c, "'`")
						} else {
							tk.Println(column, c, "	", reflect.TypeOf(x), " `bson:'", column, ".", c, "' json:'", c, "'`")
						}
					}
					break
				case reflect.Slice:
					break
				default:
					tk.Println(column, "	", datatype, " `bson:'", column, "' json:'", column, "'`")
					break
				}

			} else {
				tk.Println(column, "	unidentified `bson:'", column, "' json:'", column, "'`")
			}
			// tk.Println(column, "	datatype `bson:'' json:'", column, "'`")
		}
	}
	tk.Println("}")
}

func ReadConfig() map[string]string {
	ret := make(map[string]string)
	file, err := os.Open(wd + "/consoleapp/conf/app.conf")
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
		tk.Println(err.Error())
	}

	return ret
}
func PrepareConnection(ConnectionType string) (dbox.IConnection, error) {
	config := ReadConfig()
	tk.Println("Connecting database..")
	ci := &dbox.ConnectionInfo{config["host_"+ConnectionType], config["database_"+ConnectionType], config["username_"+ConnectionType], config["password_"+ConnectionType], nil}
	c, e := dbox.NewConnection(ConnectionType, ci)

	if e != nil {
		return nil, e
	}

	e = c.Connect()
	if e != nil {
		return nil, e
	}
	tk.Println("Database connected..")
	return c, nil
}
