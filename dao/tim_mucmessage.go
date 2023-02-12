package dao

/**
tablename:tim_mucmessage
datetime :2016-09-07 11:32:22
*/
import (
	"github.com/donnie4w/gdao"
	"reflect"
)

type tim_mucmessage_Stanza struct {
	gdao.Field
	fieldName  string
	FieldValue *string
}

func (c *tim_mucmessage_Stanza) Name() string {
	return c.fieldName
}

func (c *tim_mucmessage_Stanza) Value() interface{} {
	return c.FieldValue
}

type tim_mucmessage_Createtime struct {
	gdao.Field
	fieldName  string
	FieldValue *string
}

func (c *tim_mucmessage_Createtime) Name() string {
	return c.fieldName
}

func (c *tim_mucmessage_Createtime) Value() interface{} {
	return c.FieldValue
}

type tim_mucmessage_Id struct {
	gdao.Field
	fieldName  string
	FieldValue *int32
}

func (c *tim_mucmessage_Id) Name() string {
	return c.fieldName
}

func (c *tim_mucmessage_Id) Value() interface{} {
	return c.FieldValue
}

type tim_mucmessage_Stamp struct {
	gdao.Field
	fieldName  string
	FieldValue *string
}

func (c *tim_mucmessage_Stamp) Name() string {
	return c.fieldName
}

func (c *tim_mucmessage_Stamp) Value() interface{} {
	return c.FieldValue
}

type tim_mucmessage_Fromuser struct {
	gdao.Field
	fieldName  string
	FieldValue *string
}

func (c *tim_mucmessage_Fromuser) Name() string {
	return c.fieldName
}

func (c *tim_mucmessage_Fromuser) Value() interface{} {
	return c.FieldValue
}

type tim_mucmessage_Roomtidname struct {
	gdao.Field
	fieldName  string
	FieldValue *string
}

func (c *tim_mucmessage_Roomtidname) Name() string {
	return c.fieldName
}

func (c *tim_mucmessage_Roomtidname) Value() interface{} {
	return c.FieldValue
}

type tim_mucmessage_Domain struct {
	gdao.Field
	fieldName  string
	FieldValue *string
}

func (c *tim_mucmessage_Domain) Name() string {
	return c.fieldName
}

func (c *tim_mucmessage_Domain) Value() interface{} {
	return c.FieldValue
}

type tim_mucmessage_Msgtype struct {
	gdao.Field
	fieldName  string
	FieldValue *int32
}

func (c *tim_mucmessage_Msgtype) Name() string {
	return c.fieldName
}

func (c *tim_mucmessage_Msgtype) Value() interface{} {
	return c.FieldValue
}

type Tim_mucmessage struct {
	gdao.Table
	Roomtidname *tim_mucmessage_Roomtidname
	Domain *tim_mucmessage_Domain
	Msgtype *tim_mucmessage_Msgtype
	Stanza *tim_mucmessage_Stanza
	Createtime *tim_mucmessage_Createtime
	Id *tim_mucmessage_Id
	Stamp *tim_mucmessage_Stamp
	Fromuser *tim_mucmessage_Fromuser
}

func (u *Tim_mucmessage) GetStanza() string {
	return *u.Stanza.FieldValue
}

func (u *Tim_mucmessage) SetStanza(arg string) {
	u.Table.ModifyMap[u.Stanza.fieldName] = arg
	v := string(arg)
	u.Stanza.FieldValue = &v
}

func (u *Tim_mucmessage) GetCreatetime() string {
	return *u.Createtime.FieldValue
}

func (u *Tim_mucmessage) SetCreatetime(arg string) {
	u.Table.ModifyMap[u.Createtime.fieldName] = arg
	v := string(arg)
	u.Createtime.FieldValue = &v
}

func (u *Tim_mucmessage) GetId() int32 {
	return *u.Id.FieldValue
}

func (u *Tim_mucmessage) SetId(arg int64) {
	u.Table.ModifyMap[u.Id.fieldName] = arg
	v := int32(arg)
	u.Id.FieldValue = &v
}

func (u *Tim_mucmessage) GetStamp() string {
	return *u.Stamp.FieldValue
}

func (u *Tim_mucmessage) SetStamp(arg string) {
	u.Table.ModifyMap[u.Stamp.fieldName] = arg
	v := string(arg)
	u.Stamp.FieldValue = &v
}

func (u *Tim_mucmessage) GetFromuser() string {
	return *u.Fromuser.FieldValue
}

func (u *Tim_mucmessage) SetFromuser(arg string) {
	u.Table.ModifyMap[u.Fromuser.fieldName] = arg
	v := string(arg)
	u.Fromuser.FieldValue = &v
}

func (u *Tim_mucmessage) GetRoomtidname() string {
	return *u.Roomtidname.FieldValue
}

func (u *Tim_mucmessage) SetRoomtidname(arg string) {
	u.Table.ModifyMap[u.Roomtidname.fieldName] = arg
	v := string(arg)
	u.Roomtidname.FieldValue = &v
}

func (u *Tim_mucmessage) GetDomain() string {
	return *u.Domain.FieldValue
}

func (u *Tim_mucmessage) SetDomain(arg string) {
	u.Table.ModifyMap[u.Domain.fieldName] = arg
	v := string(arg)
	u.Domain.FieldValue = &v
}

func (u *Tim_mucmessage) GetMsgtype() int32 {
	return *u.Msgtype.FieldValue
}

func (u *Tim_mucmessage) SetMsgtype(arg int64) {
	u.Table.ModifyMap[u.Msgtype.fieldName] = arg
	v := int32(arg)
	u.Msgtype.FieldValue = &v
}

func (t *Tim_mucmessage) Query(columns ...gdao.Column) ([]Tim_mucmessage,error) {
	if columns == nil {
		columns = []gdao.Column{ t.Domain,t.Msgtype,t.Stanza,t.Createtime,t.Id,t.Stamp,t.Fromuser,t.Roomtidname}
	}
	rs,err := t.Table.Query(columns...)
	if rs == nil || err != nil {
		return nil, err
	}
	ts := make([]Tim_mucmessage, 0, len(rs))
	c := make(chan int16,len(rs))
	for _, rows := range rs {
		t := NewTim_mucmessage()
		go copyTim_mucmessage(c, rows, t, columns)
		<-c
		ts = append(ts, *t)
	}
	return ts,nil
}

func copyTim_mucmessage(channle chan int16, rows []interface{}, t *Tim_mucmessage, columns []gdao.Column) {
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

func (t *Tim_mucmessage) QuerySingle(columns ...gdao.Column) (*Tim_mucmessage,error) {
	if columns == nil {
		columns = []gdao.Column{ t.Domain,t.Msgtype,t.Stanza,t.Createtime,t.Id,t.Stamp,t.Fromuser,t.Roomtidname}
	}
	rs,err := t.Table.QuerySingle(columns...)
	if rs == nil || err != nil {
		return nil, err
	}
	rt := NewTim_mucmessage()
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

func (t *Tim_mucmessage) Select(columns ...gdao.Column) (*Tim_mucmessage,error) {
	if columns == nil {
		columns = []gdao.Column{ t.Domain,t.Msgtype,t.Stanza,t.Createtime,t.Id,t.Stamp,t.Fromuser,t.Roomtidname}
	}
	rows,err := t.Table.Selects(columns...)
	defer rows.Close()
	if err != nil || rows==nil {
		return nil, err
	}
	buff := make([]interface{}, len(columns))
	if rows.Next() {
		n := NewTim_mucmessage()
		cpTim_mucmessage(buff, n, columns)
		row_err := rows.Scan(buff...)
		if row_err != nil {
			return nil, row_err
		}
		return n, nil
	}
	return nil, nil
}

func (t *Tim_mucmessage) Selects(columns ...gdao.Column) ([]*Tim_mucmessage,error) {
	if columns == nil {
		columns = []gdao.Column{ t.Domain,t.Msgtype,t.Stanza,t.Createtime,t.Id,t.Stamp,t.Fromuser,t.Roomtidname}
	}
	rows,err := t.Table.Selects(columns...)
	defer rows.Close()
	if err != nil || rows==nil {
		return nil, err
	}
	ns := make([]*Tim_mucmessage, 0)
	buff := make([]interface{}, len(columns))
	for rows.Next() {
		n := NewTim_mucmessage()
		cpTim_mucmessage(buff, n, columns)
		row_err := rows.Scan(buff...)
		if row_err != nil {
			return nil, row_err
		}
		ns = append(ns, n)
	}
	return ns, nil
}

func  cpTim_mucmessage(buff []interface{}, t *Tim_mucmessage, columns []gdao.Column) {
	for i, column := range columns {
		field := column.Name()
		switch field {
		case "createtime":
			buff[i] = &t.Createtime.FieldValue
		case "id":
			buff[i] = &t.Id.FieldValue
		case "stamp":
			buff[i] = &t.Stamp.FieldValue
		case "fromuser":
			buff[i] = &t.Fromuser.FieldValue
		case "roomtidname":
			buff[i] = &t.Roomtidname.FieldValue
		case "domain":
			buff[i] = &t.Domain.FieldValue
		case "msgtype":
			buff[i] = &t.Msgtype.FieldValue
		case "stanza":
			buff[i] = &t.Stanza.FieldValue
		}
	}
}

func NewTim_mucmessage(tableName ...string) *Tim_mucmessage {
	stamp := &tim_mucmessage_Stamp{fieldName: "stamp"}
	stamp.Field.FieldName = "stamp"
	fromuser := &tim_mucmessage_Fromuser{fieldName: "fromuser"}
	fromuser.Field.FieldName = "fromuser"
	roomtidname := &tim_mucmessage_Roomtidname{fieldName: "roomtidname"}
	roomtidname.Field.FieldName = "roomtidname"
	domain := &tim_mucmessage_Domain{fieldName: "domain"}
	domain.Field.FieldName = "domain"
	msgtype_ := &tim_mucmessage_Msgtype{fieldName: "msgtype"}
	msgtype_.Field.FieldName = "msgtype"
	stanza := &tim_mucmessage_Stanza{fieldName: "stanza"}
	stanza.Field.FieldName = "stanza"
	createtime := &tim_mucmessage_Createtime{fieldName: "createtime"}
	createtime.Field.FieldName = "createtime"
	id := &tim_mucmessage_Id{fieldName: "id"}
	id.Field.FieldName = "id"
	table := &Tim_mucmessage{Id:id,Stamp:stamp,Fromuser:fromuser,Roomtidname:roomtidname,Domain:domain,Msgtype:msgtype_,Stanza:stanza,Createtime:createtime}
	table.Table.ModifyMap = make(map[string]interface{})
	if len(tableName) == 1 {
		table.Table.TableName = tableName[0]
	} else {
		table.Table.TableName = "tim_mucmessage"
	}
	return table
}
