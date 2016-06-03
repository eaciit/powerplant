package controllers

import (
	"bufio"
	"log"

	"github.com/eaciit/dbox"
	"github.com/eaciit/orm"
	. "github.com/eaciit/powerplant/sec/consoleapp/models"
	tk "github.com/eaciit/toolkit"
	// "math"
	"os"
	"reflect"
	"strings"
	"sync"
	"time"
)

var (
	wd = func() string {
		d, _ := os.Getwd()
		return d + "/"
	}()

	// mu                 = &sync.Mutex{}
	retry              = 10
	worker             = 100
	maxDataEachProcess = 100000
	idx                = 0
	mu                 = &sync.Mutex{}
	muinsert           = &sync.Mutex{}
)

type IBaseController interface {
	// not implemented anything yet
}

type BaseController struct {
	base     IBaseController
	MongoCtx *orm.DataContext
	SqlCtx   *orm.DataContext
}

// func (b *BaseController) ConvertMGOToSQLServer(m orm.IModel) error {
// 	// OLD
// 	tk.Printf("\nConvertMGOToSQLServer: Converting %v \n", m.TableName())
// 	tk.Println("ConvertMGOToSQLServer: Starting to convert...\n")
// 	date := time.Now()
// 	csr, e := b.MongoCtx.Connection.NewQuery().From(m.TableName()).Cursor(nil)
// 	defer csr.Close()
// 	if e != nil {
// 		return e
// 	}
// 	result := []tk.M{}
// 	e = csr.Fetch(&result, 0, false)
// 	if e != nil {
// 		return e
// 	}
// 	query := b.SqlCtx.Connection.NewQuery().SetConfig("multiexec", true).From(m.TableName()).Save()
// 	var wg sync.WaitGroup
// 	take := 10000
// 	diff := float64(len(result) / take)
// 	wgTotal := math.Ceil(diff)
// 	wgTotal++
// 	wg.Add(int(wgTotal))
// 	tk.Println("TOTAL : ", len(result))
// 	for x := 0; x < int(wgTotal); x++ {
// 		go func(swg *sync.WaitGroup, result []tk.M, skip int, max int) {
// 			t := skip + take
// 			if (skip + take) > max {
// 				t = max
// 			}
// 			for _, i := range result[skip:t] {
// 				valueType := reflect.TypeOf(m).Elem()
// 				for f := 0; f < valueType.NumField(); f++ {
// 					field := valueType.Field(f)
// 					bsonField := field.Tag.Get("bson")
// 					jsonField := field.Tag.Get("json")
// 					if jsonField != bsonField && field.Name != "RWMutex" && field.Name != "ModelBase" {
// 						i.Set(field.Name, GetMgoValue(i, bsonField))
// 					}
// 					switch field.Type.Name() {
// 					case "string":
// 						if GetMgoValue(i, bsonField) == nil {
// 							i.Set(field.Name, "")
// 						}
// 						break
// 					case "Time":
// 						if GetMgoValue(i, bsonField) == nil {
// 							i.Set(field.Name, time.Time{})
// 						} else {
// 							i.Set(field.Name, GetMgoValue(i, bsonField).(time.Time).UTC())
// 						}
// 						break
// 					default:
// 						break
// 					}
// 				}
// 				d := getNewPointer(m)
// 				e := tk.Serde(i, &d, "json")

// 				// query := b.SqlCtx.Connection.NewQuery().From(m.TableName()).Save()
// 				e = query.Exec(tk.M{"data": d})
// 				if e != nil {
// 					tk.Println(m)
// 					tk.Println("-------------------------------------")
// 					tk.Println(d)
// 					tk.Println("-------------------------------------")
// 					tk.Println(e)
// 					wg.Done()
// 				}
// 			}
// 			swg.Done()
// 		}(&wg, result, (x * take), len(result))
// 	}
// 	wg.Wait()

// 	tk.Println("\nConvertMGOToSQLServer: Finish.")
// 	tk.Println(time.Since(date))
// 	return nil
// }

func (b *BaseController) ConvertMGOToSQLServer(m orm.IModel) error {
	tStart := time.Now()

	tk.Printf("\nConvertMGOToSQLServer: Converting %v \n", m.TableName())
	tk.Println("ConvertMGOToSQLServer: Starting to convert...\n")

	c, e := b.MongoCtx.Connection.NewQuery().From(m.TableName()).Cursor(nil)

	if e != nil {
		return e
	}

	defer c.Close()

	totalData := c.Count()
	processIter := tk.ToInt(tk.ToFloat64(totalData/maxDataEachProcess, 0, tk.RoundingUp), tk.RoundingUp)

	if maxDataEachProcess == 0 {
		processIter = 0
	}

	for iter := 0; iter < processIter+1; iter++ {

		skip := iter * maxDataEachProcess
		take := maxDataEachProcess

		if maxDataEachProcess == 0 {
			take = totalData
		} else if iter == processIter {
			take = totalData - skip
		}

		csr, e := b.MongoCtx.Connection.NewQuery().From(m.TableName()).Skip(skip).Take(take).Cursor(nil)

		if e != nil {
			return e
		}

		result := []tk.M{}
		e = csr.Fetch(&result, 0, false)
		csr.Close()

		if e != nil {
			return e
		}

		dtLen := len(result)

		resPart := make([][]tk.M, worker)

		if dtLen < worker {
			resPart = make([][]tk.M, 1)
			resPart[0] = result
		} else {
			workerTaskCount := tk.ToInt(tk.ToFloat64(dtLen/worker, 0, tk.RoundingAuto), tk.RoundingAuto)
			count := 0

			for i := 0; i < worker; i++ {
				if i == worker-1 {
					resPart[i] = result[count:]
				} else {
					resPart[i] = result[count : count+workerTaskCount]
				}
				count += workerTaskCount
			}
		}

		wg := &sync.WaitGroup{}
		wg.Add(len(resPart))

		for _, val := range resPart {
			// go b.Insert(val, m, wg)
			go b.InsertBulk(val, m, wg)
		}

		wg.Wait()
	}

	cr, e := b.SqlCtx.Connection.NewQuery().From(m.TableName()).Cursor(nil)
	ctn := cr.Count()
	cr.Close()

	tk.Println("\nConvertMGOToSQLServer: Finish.")
	tk.Printf("Completed Success in %v | %v data(s)\n", time.Since(tStart), ctn)
	return nil
}

func (b *BaseController) InsertBulk(result []tk.M, m orm.IModel, wg *sync.WaitGroup) {
	var datas []orm.IModel
	for _, i := range result {
		valueType := reflect.TypeOf(m).Elem()
		for f := 0; f < valueType.NumField(); f++ {
			field := valueType.Field(f)
			bsonField := field.Tag.Get("bson")
			jsonField := field.Tag.Get("json")

			if jsonField != bsonField && field.Name != "RWMutex" && field.Name != "ModelBase" {
				i.Set(field.Name, GetMgoValue(i, bsonField))
			}
			switch field.Type.Name() {
			case "string":
				if GetMgoValue(i, bsonField) == nil {
					i.Set(field.Name, "")
				}
				break
			case "Time":
				if GetMgoValue(i, bsonField) == nil {
					i.Set(field.Name, time.Time{})
				} else {
					i.Set(field.Name, GetMgoValue(i, bsonField).(time.Time).UTC())
				}
				break
			default:
				break
			}

		}

		newPointer := getNewPointer(m)
		e := tk.Serde(i, newPointer, "json")
		datas = append(datas, newPointer)

		if e != nil {
			tk.Printf("\n----------- ERROR -------------- \n %v \n\n %#v \n\n %#v \n-------------------------  \n", e.Error(), i, newPointer)
			wg.Done()
		}

	}

	if nil != datas {
		muinsert.Lock()
		for {
			e := b.SqlCtx.InsertBulk(datas)
			if e == nil {
				ctn := len(result)
				idx += ctn
				tk.Printf("saved: %v data(s)\n", idx)
				break
			} else {
				b.SqlCtx.Connection.Connect()
			}
		}
		muinsert.Unlock()
	}

	wg.Done()
}

func (b *BaseController) Insert(result []tk.M, m orm.IModel, wg *sync.WaitGroup) {
	// muinsert := &sync.Mutex{}
	for _, i := range result {
		valueType := reflect.TypeOf(m).Elem()
		for f := 0; f < valueType.NumField(); f++ {
			field := valueType.Field(f)
			bsonField := field.Tag.Get("bson")
			jsonField := field.Tag.Get("json")

			if jsonField != bsonField && field.Name != "RWMutex" && field.Name != "ModelBase" {
				i.Set(field.Name, GetMgoValue(i, bsonField))
			}
			switch field.Type.Name() {
			case "string":
				if GetMgoValue(i, bsonField) == nil {
					i.Set(field.Name, "")
				}
				break
			case "Time":
				if GetMgoValue(i, bsonField) == nil {
					i.Set(field.Name, time.Time{})
				} else {
					i.Set(field.Name, GetMgoValue(i, bsonField).(time.Time).UTC())
				}
				break
			default:
				break
			}

		}

		newPointer := getNewPointer(m)
		e := tk.Serde(i, newPointer, "json")
		var newId int64
		for index := 0; index < retry; index++ {
			muinsert.Lock()
			newId, e = b.SqlCtx.InsertOut(newPointer)
			_ = newId
			muinsert.Unlock()
			if e == nil {
				wg.Done()
				break
			} else {
				b.SqlCtx.Connection.Connect()
			}
		}

		if e != nil {
			tk.Printf("\n----------- ERROR -------------- \n %v \n\n %#v \n\n %#v \n-------------------------  \n", e.Error(), i, newPointer)
			wg.Done()
		}

	}
	wg.Done()
}
func GetMgoValue(d tk.M, fieldName string) interface{} {
	index := strings.Index(fieldName, ".")
	if index < 0 {
		return d.Get(fieldName)
	} else {
		data := d.Get(fieldName[0:index])
		if data != nil {
			return GetMgoValue(data.(tk.M), fieldName[(index+1):len(fieldName)])
		} else {
			return nil
		}
	}
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

func getNewPointer(m orm.IModel) orm.IModel {
	switch m.TableName() {
	case "MaintenanceCostByHour":
		return new(MaintenanceCostByHour)
	case "MaintenanceCost":
		return new(MaintenanceCost)
	case "PlannedMaintenance":
		return new(PlannedMaintenance)
	case "SummaryData":
		return new(SummaryData)
	case "DataBrowser":
		return new(DataBrowser)
	case "MORCalculationFlatSummary":
		return new(MORCalculationFlatSummary)
	case "PreventiveCorrectiveSummary":
		return new(PreventiveCorrectiveSummary)
	case "RegenMasterPlant":
		return new(RegenMasterPlant)
	case "NotificationFailure":
		return new(NotificationFailure)
	case "WODurationSummary":
		return new(WODurationSummary)
	case "WOListSummary":
		return new(WOListSummary)
	case "SyntheticPM":
		return new(SyntheticPM)

	default:
		return m
	}

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

func PrepareConnection(ConnectionType string) (dbox.IConnection, error) {
	config := ReadConfig()
	tk.Println(config["host"])
	ci := &dbox.ConnectionInfo{config["host_"+ConnectionType], config["database_"+ConnectionType], config["username_"+ConnectionType], config["password_"+ConnectionType], nil}
	c, e := dbox.NewConnection(ConnectionType, ci)

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
		tk.Println(err.Error())
	}

	return ret
}
