package dao

/**
tablename:tim_config
datetime :2016-09-07 11:32:22
*/
import (
	"github.com/donnie4w/gdao"
	"reflect"
)

type tim_config_Valuestr struct {
	gdao.Field
	fieldName  string
	FieldValue *string
}

func (c *tim_config_Valuestr) Name() string {
	return c.fieldName
}

func (c *tim_config_Valuestr) Value() interface{} {
	return c.FieldValue
}

type tim_config_Createtime struct {
	gdao.Field
	fieldName  string
	FieldValue *string
}

func (c *tim_config_Createtime) Name() string {
	return c.fieldName
}

func (c *tim_config_Createtime) Value() interface{} {
	return c.FieldValue
}

type tim_config_Remark struct {
	gdao.Field
	fieldName  string
	FieldValue *string
}

func (c *tim_config_Remark) Name() string {
	return c.fieldName
}

func (c *tim_config_Remark) Value() interface{} {
	return c.FieldValue
}

type tim_config_Id struct {
	gdao.Field
	fieldName  string
	FieldValue *int32
}

func (c *tim_config_Id) Name() string {
	return c.fieldName
}

func (c *tim_config_Id) Value() interface{} {
	return c.FieldValue
}

type tim_config_Keyword struct {
	gdao.Field
	fieldName  string
	FieldValue *string
}

func (c *tim_config_Keyword) Name() string {
	return c.fieldName
}

func (c *tim_config_Keyword) Value() interface{} {
	return c.FieldValue
}

type Tim_config struct {
	gdao.Table
	Id *tim_config_Id
	Keyword *tim_config_Keyword
	Valuestr *tim_config_Valuestr
	Createtime *tim_config_Createtime
	Remark *tim_config_Remark
}

func (u *Tim_config) GetId() int32 {
	return *u.Id.FieldValue
}

func (u *Tim_config) SetId(arg int64) {
	u.Table.ModifyMap[u.Id.fieldName] = arg
	v := int32(arg)
	u.Id.FieldValue = &v
}

func (u *Tim_config) GetKeyword() string {
	return *u.Keyword.FieldValue
}

func (u *Tim_config) SetKeyword(arg string) {
	u.Table.ModifyMap[u.Keyword.fieldName] = arg
	v := string(arg)
	u.Keyword.FieldValue = &v
}

func (u *Tim_config) GetValuestr() string {
	return *u.Valuestr.FieldValue
}

func (u *Tim_config) SetValuestr(arg string) {
	u.Table.ModifyMap[u.Valuestr.fieldName] = arg
	v := string(arg)
	u.Valuestr.FieldValue = &v
}

func (u *Tim_config) GetCreatetime() string {
	return *u.Createtime.FieldValue
}

func (u *Tim_config) SetCreatetime(arg string) {
	u.Table.ModifyMap[u.Createtime.fieldName] = arg
	v := string(arg)
	u.Createtime.FieldValue = &v
}

func (u *Tim_config) GetRemark() string {
	return *u.Remark.FieldValue
}

func (u *Tim_config) SetRemark(arg string) {
	u.Table.ModifyMap[u.Remark.fieldName] = arg
	v := string(arg)
	u.Remark.FieldValue = &v
}

func (t *Tim_config) Query(columns ...gdao.Column) ([]Tim_config,error) {
	if columns == nil {
		columns = []gdao.Column{ t.Createtime,t.Remark,t.Id,t.Keyword,t.Valuestr}
	}
	rs,err := t.Table.Query(columns...)
	if rs == nil || err != nil {
		return nil, err
	}
	ts := make([]Tim_config, 0, len(rs))
	c := make(chan int16,len(rs))
	for _, rows := range rs {
		t := NewTim_config()
		go copyTim_config(c, rows, t, columns)
		<-c
		ts = append(ts, *t)
	}
	return ts,nil
}

func copyTim_config(channle chan int16, rows []interface{}, t *Tim_config, columns []gdao.Column) {
	defer func() { channle <- 1 }()
	for j, core := range rows {
		if core == nil {
			continue
		}
		field := columns[j].Name()
		setfield := "Set" + gdao.ToUpperFirstLetter(field)
		reflect.ValueOf(t).MethodByName(setfield).Call([]reflect.Value{reflect.ValueOf(gdao.GetValue(&core))})
	}
}

func (t *Tim_config) QuerySingle(columns ...gdao.Column) (*Tim_config,error) {
	if columns == nil {
		columns = []gdao.Column{ t.Createtime,t.Remark,t.Id,t.Keyword,t.Valuestr}
	}
	rs,err := t.Table.QuerySingle(columns...)
	if rs == nil || err != nil {
		return nil, err
	}
	rt := NewTim_config()
	for j, core := range rs {
		if core == nil {
			continue
		}
		field := columns[j].Name()
		setfield := "Set" + gdao.ToUpperFirstLetter(field)
		reflect.ValueOf(rt).MethodByName(setfield).Call([]reflect.Value{reflect.ValueOf(gdao.GetValue(&core))})
	}
	return rt,nil
}

func (t *Tim_config) Select(columns ...gdao.Column) (*Tim_config,error) {
	if columns == nil {
		columns = []gdao.Column{ t.Createtime,t.Remark,t.Id,t.Keyword,t.Valuestr}
	}
	rows,err := t.Table.Selects(columns...)
	defer rows.Close()
	if err != nil || rows==nil {
		return nil, err
	}
	buff := make([]interface{}, len(columns))
	if rows.Next() {
		n := NewTim_config()
		cpTim_config(buff, n, columns)
		row_err := rows.Scan(buff...)
		if row_err != nil {
			return nil, row_err
		}
		return n, nil
	}
	return nil, nil
}

func (t *Tim_config) Selects(columns ...gdao.Column) ([]*Tim_config,error) {
	if columns == nil {
		columns = []gdao.Column{ t.Createtime,t.Remark,t.Id,t.Keyword,t.Valuestr}
	}
	rows,err := t.Table.Selects(columns...)
	defer rows.Close()
	if err != nil || rows==nil {
		return nil, err
	}
	ns := make([]*Tim_config, 0)
	buff := make([]interface{}, len(columns))
	for rows.Next() {
		n := NewTim_config()
		cpTim_config(buff, n, columns)
		row_err := rows.Scan(buff...)
		if row_err != nil {
			return nil, row_err
		}
		ns = append(ns, n)
	}
	return ns, nil
}

func  cpTim_config(buff []interface{}, t *Tim_config, columns []gdao.Column) {
	for i, column := range columns {
		field := column.Name()
		switch field {
		case "id":
			buff[i] = &t.Id.FieldValue
		case "keyword":
			buff[i] = &t.Keyword.FieldValue
		case "valuestr":
			buff[i] = &t.Valuestr.FieldValue
		case "createtime":
			buff[i] = &t.Createtime.FieldValue
		case "remark":
			buff[i] = &t.Remark.FieldValue
		}
	}
}

func NewTim_config(tableName ...string) *Tim_config {
	createtime := &tim_config_Createtime{fieldName: "createtime"}
	createtime.Field.FieldName = "createtime"
	remark := &tim_config_Remark{fieldName: "remark"}
	remark.Field.FieldName = "remark"
	id := &tim_config_Id{fieldName: "id"}
	id.Field.FieldName = "id"
	keyword := &tim_config_Keyword{fieldName: "keyword"}
	keyword.Field.FieldName = "keyword"
	valuestr := &tim_config_Valuestr{fieldName: "valuestr"}
	valuestr.Field.FieldName = "valuestr"
	table := &Tim_config{Keyword:keyword,Valuestr:valuestr,Createtime:createtime,Remark:remark,Id:id}
	table.Table.ModifyMap = make(map[string]interface{})
	if len(tableName) == 1 {
		table.Table.TableName = tableName[0]
	} else {
		table.Table.TableName = "tim_config"
	}
	return table
}
