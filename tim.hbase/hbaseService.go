/**
 * donnie4w@gmail.com  tim server
 */
package hbase

//	"github.com/donnie4w/go-logger/logger"
import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"tim.utils"
)

var maxcount int32 = (1 << 14)

type Bean struct {
	Row       int64
	Family    string
	Qualifier string
	Value     string
}

func getSerialNo(tablename string, row string) (id int64, err error) {
	client := ClientPool.get()
	defer func() {
		if er := recover(); er != nil {
			ClientPool.del(client)
		} else {
			ClientPool.put(client)
		}
	}()
	if client != nil {
		hbaseClient := client.tsclient
		ti := NewTIncrement()
		ti.Row = []byte(tablename)
		tc := NewTColumnIncrement()
		tc.Family = []byte(row)
		ti.Columns = []*TColumnIncrement{tc}
		var result *TResult_
		result, err = hbaseClient.Increment([]byte("tim_serialno"), ti)
		if err == nil && result != nil && len(result.GetColumnValues()) > 0 {
			bb := result.GetColumnValues()[0].GetValue()
			id = Bytes2hex(bb)
		}
		if err != nil {
			panic(err.Error())
		}
	}
	return
}

//func _insert(tablename string, row int64, family string, qualifier string, value string) (er error) {
//	client := ClientPool.get()
//	defer func() {
//		if err := recover(); err != nil {
//			er = errors.New(fmt.Sprint(err))
//			ClientPool.del(client)
//		} else {
//			ClientPool.put(client)
//		}
//	}()
//	if client != nil {
//		hbaseClient := client.tsclient
//		tput := NewTPut()
//		tput.Row = Hex2bytes(row)
//		tColumnValue := NewTColumnValue()
//		tColumnValue.Family = []byte(family)
//		if value != "" {
//			tColumnValue.Value = []byte(value)
//		}
//		if qualifier != "" {
//			tColumnValue.Qualifier = []byte(qualifier)
//		}
//		tput.ColumnValues = []*TColumnValue{tColumnValue}
//		fmt.Println("=============>", row, " ", family, " ", qualifier, " ", value)
//		er = hbaseClient.Put([]byte(tablename), tput)
//	}
//	return
//}

func UpdateMultiple(tablename string, beans []*Bean) (er error) {
	client := ClientPool.get()
	defer func() {
		if err := recover(); err != nil {
			ClientPool.del(client)
		} else {
			ClientPool.put(client)
		}
	}()
	if client != nil && beans != nil && len(beans) > 0 {
		hbaseClient := client.tsclient
		tputs := make([]*TPut, 0)
		for _, bean := range beans {
			tput := NewTPut()
			tput.Row = Hex2bytes(bean.Row)
			tColumnValue := NewTColumnValue()
			if bean.Family != "" {
				tColumnValue.Family = []byte(bean.Family)
			}
			if bean.Qualifier != "" {
				tColumnValue.Qualifier = []byte(bean.Qualifier)
			}
			if bean.Value != "" {
				tColumnValue.Value = []byte(bean.Value)
			}
			tput.ColumnValues = []*TColumnValue{tColumnValue}
			tputs = append(tputs, tput)
		}
		er = hbaseClient.PutMultiple([]byte(tablename), tputs)
		if er != nil {
			panic(er.Error())
		}
	}
	return
}

func Hex2bytes(row int64) (bs []byte) {
	bs = make([]byte, 0)
	for i := 0; i < 8; i++ {
		r := row >> uint((7-i)*8)
		bs = append(bs, byte(r))
	}
	return
}

func Bytes2hex(bb []byte) (value int64) {
	value = int64(0x00000000)
	for i, b := range bb {
		ii := uint(b) << uint((7-i)*8)
		value = value | int64(ii)
	}
	return
}

func DeleteRow(tablename string, row int64) (er error) {
	client := ClientPool.get()
	defer func() {
		if err := recover(); err != nil {
			ClientPool.del(client)
		} else {
			ClientPool.put(client)
		}
	}()
	if client != nil {
		hbaseClient := client.tsclient
		tdelete := NewTDelete()
		tdelete.Row = Hex2bytes(row)
		er = hbaseClient.DeleteSingle([]byte(tablename), tdelete)
		if er != nil {
			panic(er.Error())
		}
	}
	return
}

//——————————————————————————————————————————————————————————————————————————————————————————————————————
func _DeleteRows(hbaseClient *THBaseServiceClient, tablename string, rows []int64) (err error) {
	tdeletes := make([]*TDelete, 0)
	for _, row := range rows {
		tdelete := NewTDelete()
		tdelete.Row = Hex2bytes(row)
		tdeletes = append(tdeletes, tdelete)
	}
	_, err = hbaseClient.DeleteMultiple([]byte(tablename), tdeletes)
	return
}

func DeleteRows(tablename string, rows []int64) (er error) {
	client := ClientPool.get()
	defer func() {
		if err := recover(); err != nil {
			ClientPool.del(client)
		} else {
			ClientPool.put(client)
		}
	}()
	if client != nil {
		hbaseClient := client.tsclient
		er = _DeleteRows(hbaseClient, tablename, rows)
		if er != nil {
			panic(er.Error())
		}
	}
	return
}

func DeleteFromQualifier(tablename string, beans []*Bean) (err error) {
	client := ClientPool.get()
	defer func() {
		if er := recover(); er != nil {
			err = errors.New(fmt.Sprint(er))
			ClientPool.del(client)
		} else {
			ClientPool.put(client)
		}
	}()
	if client != nil {
		hbaseClient := client.tsclient
		tscan := NewTScan()
		columns := make([]*TColumn, 0)
		for _, bean := range beans {
			family := bean.Family
			qualifier := bean.Qualifier
			if family != "" || qualifier != "" {
				tcolumn := NewTColumn()
				if family != "" {
					tcolumn.Family = []byte(family)
				}
				if qualifier != "" {
					tcolumn.Qualifier = []byte(qualifier)
				}
				columns = append(columns, tcolumn)
			}
		}
		tscan.Columns = columns
		rs, er := hbaseClient.GetScannerResults([]byte(tablename), tscan, maxcount)
		if er != nil {
			panic(er.Error())
		}
		if rs != nil {
			uniqueMap := make(map[int64]byte, 0)
			rows := make([]int64, 0)
			for _, r := range rs {
				row := Bytes2hex(r.GetRow())
				if _, ok := uniqueMap[row]; !ok {
					uniqueMap[row] = 0
					rows = append(rows, row)
				}
			}
			er := _DeleteRows(hbaseClient, tablename, rows)
			if er != nil {
				panic(er.Error())
			}
		}
	}
	return
}

//——————————————————————————————————————————————————————————————————————————————————————————————————————————————————————————————————————————————
func _getResultByBean(hbaseClient *THBaseServiceClient, tablename string, beans []*Bean, count int32, reversed bool) (rs []*TResult_, err error) {
	tscan := NewTScan()
	columns := make([]*TColumn, 0)
	for _, bean := range beans {
		family := bean.Family
		qualifier := bean.Qualifier
		if family != "" || qualifier != "" {
			tcolumn := NewTColumn()
			if family != "" {
				tcolumn.Family = []byte(family)
			}
			if qualifier != "" {
				tcolumn.Qualifier = []byte(qualifier)
			}
			columns = append(columns, tcolumn)
		}
	}
	tscan.Columns = columns
	tscan.Reversed = &reversed
	if count <= 0 {
		count = maxcount
	}
	rs, err = hbaseClient.GetScannerResults([]byte(tablename), tscan, count)
	return
}

func saveObject(o interface{}, tablename string, row int64) (er error) {
	s := reflect.TypeOf(o).Elem()
	beans := make([]*Bean, 0)
	for i := 0; i < s.NumField(); i++ {
		family := fmt.Sprint(s.Field(i).Tag)
		value := ""
		qualifier := ""
		if strings.HasPrefix(family, "#") {
			family = family[1:]
			value = fmt.Sprint(row)
		} else if family == "idx_" {
			family = "index"
			qualifier = reflect.ValueOf(o).Elem().FieldByName(s.Field(i).Name).String()
			value = ""
		} else {
			value = reflect.ValueOf(o).Elem().FieldByName(s.Field(i).Name).String()
		}
		if qualifier == "" && value == "" {
			continue
		}
		bean := new(Bean)
		bean.Family = family
		bean.Row = row
		bean.Qualifier = qualifier
		bean.Value = value
		beans = append(beans, bean)
		//		_insert(tablename, row, family, qualifier, value)
	}
	er = UpdateMultiple(tablename, beans)
	return
}

func saveObjects(o interface{}, tablename string, rows []int64) (er error) {
	s := reflect.TypeOf(o).Elem()
	beans := make([]*Bean, 0)
	for _, row := range rows {
		for i := 0; i < s.NumField(); i++ {
			family := fmt.Sprint(s.Field(i).Tag)
			value := ""
			qualifier := ""
			if strings.HasPrefix(family, "#") {
				family = family[1:]
				value = fmt.Sprint(row)
			} else if family == "idx_" {
				family = "index"
				qualifier = reflect.ValueOf(o).Elem().FieldByName(s.Field(i).Name).String()
				if qualifier == "" {
					continue
				}
				value = ""
			} else {
				value = reflect.ValueOf(o).Elem().FieldByName(s.Field(i).Name).String()
			}
			if value == "" && qualifier == "" {
				continue
			}
			bean := new(Bean)
			bean.Row = row
			bean.Family = family
			bean.Value = value
			bean.Qualifier = qualifier
			beans = append(beans, bean)
		}
	}
	er = UpdateMultiple(tablename, beans)
	return
}

func saveObjectByBeans(o interface{}, tablename string, beans []*Bean) (er error) {
	client := ClientPool.get()
	defer func() {
		if err := recover(); err != nil {
			er = errors.New(fmt.Sprint(err))
			ClientPool.del(client)
		} else {
			ClientPool.put(client)
		}
	}()
	if client != nil {
		hbaseClient := client.tsclient
		s := reflect.TypeOf(o).Elem()
		rs, err := _getResultByBean(hbaseClient, tablename, beans, 0, false)
		if err != nil {
			panic(err.Error())
		}
		if rs != nil {
			rows := make([]int64, 0)
			uniqueMap := make(map[int64]byte, 0)
			for _, r := range rs {
				row := Bytes2hex(r.GetRow())
				if _, ok := uniqueMap[row]; !ok {
					uniqueMap[row] = 0
					rows = append(rows, row)
				}
			}
			beans = make([]*Bean, 0)
			for _, row := range rows {
				for i := 0; i < s.NumField(); i++ {
					family := fmt.Sprint(s.Field(i).Tag)
					value := ""
					qualifier := ""
					if strings.HasPrefix(family, "#") {
						family = family[1:]
						value = fmt.Sprint(row)
					} else if family == "idx_" {
						family = "index"
						qualifier = reflect.ValueOf(o).Elem().FieldByName(s.Field(i).Name).String()
						if qualifier == "" {
							continue
						}
						value = ""
					} else {
						value = reflect.ValueOf(o).Elem().FieldByName(s.Field(i).Name).String()
					}
					if value == "" && qualifier == "" {
						continue
					}
					bean := new(Bean)
					bean.Row = row
					bean.Family = family
					bean.Value = value
					bean.Qualifier = qualifier
					beans = append(beans, bean)
				}
			}
			err := UpdateMultiple(tablename, beans)
			if err != nil {
				panic(err.Error())
			}
		}
	}
	return
}

//func scan(tablename string, family, qualifier string, count int32, reversed *bool) (results []*TResult_, err error) {
//	client := ClientPool.get()
//	defer func() {
//		if err := recover(); err != nil {
//			ClientPool.del(client)
//		} else {
//			ClientPool.put(client)
//		}
//	}()
//	if client != nil {
//		hbaseClient := client.tsclient
//		tscan := NewTScan()
//		if family != "" || qualifier != "" {
//			tcolumn := NewTColumn()
//			if family != "" {
//				tcolumn.Family = []byte(family)
//			}
//			if qualifier != "" {
//				tcolumn.Qualifier = []byte(qualifier)
//			}
//			tscan.Columns = []*TColumn{tcolumn}
//		}
//		tscan.Reversed = reversed
//		results, err = hbaseClient.GetScannerResults([]byte(tablename), tscan, count)
//	}
//	return
//}

func Scans(tablename string, beans []*Bean, count int32, reversed bool) (results []*TResult_, err error) {
	client := ClientPool.get()
	defer func() {
		if err := recover(); err != nil {
			ClientPool.del(client)
		} else {
			ClientPool.put(client)
		}
	}()
	if client != nil {
		hbaseClient := client.tsclient
		results, err = _getResultByBean(hbaseClient, tablename, beans, count, reversed)
		if err != nil {
			panic(err.Error())
		}
	}
	return
}

func ScansFromRow(tablename string, beans []*Bean, count int32, reversed bool) (results []*TResult_, err error) {
	client := ClientPool.get()
	defer func() {
		if err := recover(); err != nil {
			ClientPool.del(client)
		} else {
			ClientPool.put(client)
		}
	}()
	if client != nil {
		hbaseClient := client.tsclient
		rs, err := _getResultByBean(hbaseClient, tablename, beans, 0, false)
		if err == nil && rs != nil {
			tgets := make([]*TGet, 0)
			uniqueMap := make(map[int64]byte, 0)
			for _, r := range rs {
				row := Bytes2hex(r.GetRow())
				if _, ok := uniqueMap[row]; !ok {
					uniqueMap[row] = 0
					tget := NewTGet()
					tget.Row = r.GetRow()
					tgets = append(tgets, tget)
				}
			}
			results, err = hbaseClient.GetMultiple([]byte(tablename), tgets)
			if err != nil {
				panic(err.Error())
			}
		}
	}
	return
}

func SelectByRows(tablename string, rows []int64) (results []*TResult_, err error) {
	client := ClientPool.get()
	defer func() {
		if err := recover(); err != nil {
			ClientPool.del(client)
		} else {
			ClientPool.put(client)
		}
	}()
	if client != nil {
		hbaseClient := client.tsclient
		tgets := make([]*TGet, 0)
		uniqueMap := make(map[int64]byte, 0)
		for _, row := range rows {
			if _, ok := uniqueMap[row]; !ok {
				uniqueMap[row] = 0
				tget := NewTGet()
				tget.Row = Hex2bytes(row)
				tgets = append(tgets, tget)
			}
		}
		results, err = hbaseClient.GetMultiple([]byte(tablename), tgets)
		if err != nil {
			panic(err.Error())
		}
	}
	return
}

func get(tablename string, row int64, family, qualifier string) (result *TResult_, err error) {
	client := ClientPool.get()
	defer func() {
		if err := recover(); err != nil {
			ClientPool.del(client)
		} else {
			ClientPool.put(client)
		}
	}()
	if client != nil {
		hbaseClient := client.tsclient
		tget := NewTGet()
		tget.Row = Hex2bytes(row)
		if family != "" || qualifier != "" {
			tColumn := NewTColumn()
			if family != "" {
				tColumn.Family = []byte(family)
			}
			if qualifier != "" {
				tColumn.Qualifier = []byte(qualifier)
			}
			tget.Columns = []*TColumn{tColumn}
		}
		result, err = hbaseClient.Get([]byte(tablename), tget)
		if err != nil {
			panic(err.Error())
		}
	}
	return
}

func Result2object(result *TResult_, o interface{}) {
	if result != nil {
		resultmap := make(map[string]string, 0)
		for _, resultColumnValue := range result.GetColumnValues() {
			resultmap[string(resultColumnValue.GetFamily())] = string(resultColumnValue.GetValue())
		}
		s := reflect.TypeOf(o).Elem()
		for i := 0; i < s.NumField(); i++ {
			fieldname := fmt.Sprint(s.Field(i).Tag)
			if fieldname == "idx" {
				continue
			}
			if strings.HasPrefix(fieldname, "#") {
				fieldname = fieldname[1:]
				if v, ok := resultmap[fieldname]; ok {
					reflect.ValueOf(o).Elem().FieldByName(s.Field(i).Name).Set(reflect.ValueOf(utils.Atoi64(v)))
				}
			} else {
				if v, ok := resultmap[fieldname]; ok {
					reflect.ValueOf(o).Elem().FieldByName(s.Field(i).Name).SetString(v)
				}
			}
		}
	}
}

func Select(tablename string, row int64, family, qualifier string, o interface{}) error {
	result, err := get(tablename, row, family, qualifier)
	if err == nil {
		Result2object(result, o)
	}
	return err
}

func Selects(tablename string, family, qualifier string, count int32, reversed bool) (results []*TResult_, err error) {
	bean := new(Bean)
	bean.Family = family
	bean.Qualifier = qualifier
	//	results, err = scan(tablename, family, qualifier, count, &reversed)
	results, err = Scans(tablename, []*Bean{bean}, count, reversed)
	return
}
