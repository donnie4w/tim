package dao

/**
tablename:tim_roster
datetime :2016-09-07 11:32:22
*/
import (
	"github.com/donnie4w/gdao"
	"reflect"
)

type tim_roster_Loginname struct {
	gdao.Field
	fieldName  string
	FieldValue *string
}

func (c *tim_roster_Loginname) Name() string {
	return c.fieldName
}

func (c *tim_roster_Loginname) Value() interface{} {
	return c.FieldValue
}

type tim_roster_Username struct {
	gdao.Field
	fieldName  string
	FieldValue *string
}

func (c *tim_roster_Username) Name() string {
	return c.fieldName
}

func (c *tim_roster_Username) Value() interface{} {
	return c.FieldValue
}

type tim_roster_Rostername struct {
	gdao.Field
	fieldName  string
	FieldValue *string
}

func (c *tim_roster_Rostername) Name() string {
	return c.fieldName
}

func (c *tim_roster_Rostername) Value() interface{} {
	return c.FieldValue
}

type tim_roster_Rostertype struct {
	gdao.Field
	fieldName  string
	FieldValue *string
}

func (c *tim_roster_Rostertype) Name() string {
	return c.fieldName
}

func (c *tim_roster_Rostertype) Value() interface{} {
	return c.FieldValue
}

type tim_roster_Createtime struct {
	gdao.Field
	fieldName  string
	FieldValue *string
}

func (c *tim_roster_Createtime) Name() string {
	return c.fieldName
}

func (c *tim_roster_Createtime) Value() interface{} {
	return c.FieldValue
}

type tim_roster_Remarknick struct {
	gdao.Field
	fieldName  string
	FieldValue *string
}

func (c *tim_roster_Remarknick) Name() string {
	return c.fieldName
}

func (c *tim_roster_Remarknick) Value() interface{} {
	return c.FieldValue
}

type tim_roster_Id struct {
	gdao.Field
	fieldName  string
	FieldValue *int32
}

func (c *tim_roster_Id) Name() string {
	return c.fieldName
}

func (c *tim_roster_Id) Value() interface{} {
	return c.FieldValue
}

type Tim_roster struct {
	gdao.Table
	Username *tim_roster_Username
	Rostername *tim_roster_Rostername
	Rostertype *tim_roster_Rostertype
	Createtime *tim_roster_Createtime
	Remarknick *tim_roster_Remarknick
	Id *tim_roster_Id
	Loginname *tim_roster_Loginname
}

func (u *Tim_roster) GetRostername() string {
	return *u.Rostername.FieldValue
}

func (u *Tim_roster) SetRostername(arg string) {
	u.Table.ModifyMap[u.Rostername.fieldName] = arg
	v := string(arg)
	u.Rostername.FieldValue = &v
}

func (u *Tim_roster) GetRostertype() string {
	return *u.Rostertype.FieldValue
}

func (u *Tim_roster) SetRostertype(arg string) {
	u.Table.ModifyMap[u.Rostertype.fieldName] = arg
	v := string(arg)
	u.Rostertype.FieldValue = &v
}

func (u *Tim_roster) GetCreatetime() string {
	return *u.Createtime.FieldValue
}

func (u *Tim_roster) SetCreatetime(arg string) {
	u.Table.ModifyMap[u.Createtime.fieldName] = arg
	v := string(arg)
	u.Createtime.FieldValue = &v
}

func (u *Tim_roster) GetRemarknick() string {
	return *u.Remarknick.FieldValue
}

func (u *Tim_roster) SetRemarknick(arg string) {
	u.Table.ModifyMap[u.Remarknick.fieldName] = arg
	v := string(arg)
	u.Remarknick.FieldValue = &v
}

func (u *Tim_roster) GetId() int32 {
	return *u.Id.FieldValue
}

func (u *Tim_roster) SetId(arg int64) {
	u.Table.ModifyMap[u.Id.fieldName] = arg
	v := int32(arg)
	u.Id.FieldValue = &v
}

func (u *Tim_roster) GetLoginname() string {
	return *u.Loginname.FieldValue
}

func (u *Tim_roster) SetLoginname(arg string) {
	u.Table.ModifyMap[u.Loginname.fieldName] = arg
	v := string(arg)
	u.Loginname.FieldValue = &v
}

func (u *Tim_roster) GetUsername() string {
	return *u.Username.FieldValue
}

func (u *Tim_roster) SetUsername(arg string) {
	u.Table.ModifyMap[u.Username.fieldName] = arg
	v := string(arg)
	u.Username.FieldValue = &v
}

func (t *Tim_roster) Query(columns ...gdao.Column) ([]Tim_roster,error) {
	if columns == nil {
		columns = []gdao.Column{ t.Loginname,t.Username,t.Rostername,t.Rostertype,t.Createtime,t.Remarknick,t.Id}
	}
	rs,err := t.Table.Query(columns...)
	if rs == nil || err != nil {
		return nil, err
	}
	ts := make([]Tim_roster, 0, len(rs))
	c := make(chan int16,len(rs))
	for _, rows := range rs {
		t := NewTim_roster()
		go copyTim_roster(c, rows, t, columns)
		<-c
		ts = append(ts, *t)
	}
	return ts,nil
}

func copyTim_roster(channle chan int16, rows []interface{}, t *Tim_roster, columns []gdao.Column) {
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

func (t *Tim_roster) QuerySingle(columns ...gdao.Column) (*Tim_roster,error) {
	if columns == nil {
		columns = []gdao.Column{ t.Loginname,t.Username,t.Rostername,t.Rostertype,t.Createtime,t.Remarknick,t.Id}
	}
	rs,err := t.Table.QuerySingle(columns...)
	if rs == nil || err != nil {
		return nil, err
	}
	rt := NewTim_roster()
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

func (t *Tim_roster) Select(columns ...gdao.Column) (*Tim_roster,error) {
	if columns == nil {
		columns = []gdao.Column{ t.Loginname,t.Username,t.Rostername,t.Rostertype,t.Createtime,t.Remarknick,t.Id}
	}
	rows,err := t.Table.Selects(columns...)
	defer rows.Close()
	if err != nil || rows==nil {
		return nil, err
	}
	buff := make([]interface{}, len(columns))
	if rows.Next() {
		n := NewTim_roster()
		cpTim_roster(buff, n, columns)
		row_err := rows.Scan(buff...)
		if row_err != nil {
			return nil, row_err
		}
		return n, nil
	}
	return nil, nil
}

func (t *Tim_roster) Selects(columns ...gdao.Column) ([]*Tim_roster,error) {
	if columns == nil {
		columns = []gdao.Column{ t.Loginname,t.Username,t.Rostername,t.Rostertype,t.Createtime,t.Remarknick,t.Id}
	}
	rows,err := t.Table.Selects(columns...)
	defer rows.Close()
	if err != nil || rows==nil {
		return nil, err
	}
	ns := make([]*Tim_roster, 0)
	buff := make([]interface{}, len(columns))
	for rows.Next() {
		n := NewTim_roster()
		cpTim_roster(buff, n, columns)
		row_err := rows.Scan(buff...)
		if row_err != nil {
			return nil, row_err
		}
		ns = append(ns, n)
	}
	return ns, nil
}

func  cpTim_roster(buff []interface{}, t *Tim_roster, columns []gdao.Column) {
	for i, column := range columns {
		field := column.Name()
		switch field {
		case "rostername":
			buff[i] = &t.Rostername.FieldValue
		case "rostertype":
			buff[i] = &t.Rostertype.FieldValue
		case "createtime":
			buff[i] = &t.Createtime.FieldValue
		case "remarknick":
			buff[i] = &t.Remarknick.FieldValue
		case "id":
			buff[i] = &t.Id.FieldValue
		case "loginname":
			buff[i] = &t.Loginname.FieldValue
		case "username":
			buff[i] = &t.Username.FieldValue
		}
	}
}

func NewTim_roster(tableName ...string) *Tim_roster {
	remarknick := &tim_roster_Remarknick{fieldName: "remarknick"}
	remarknick.Field.FieldName = "remarknick"
	id := &tim_roster_Id{fieldName: "id"}
	id.Field.FieldName = "id"
	loginname := &tim_roster_Loginname{fieldName: "loginname"}
	loginname.Field.FieldName = "loginname"
	username := &tim_roster_Username{fieldName: "username"}
	username.Field.FieldName = "username"
	rostername := &tim_roster_Rostername{fieldName: "rostername"}
	rostername.Field.FieldName = "rostername"
	rostertype_ := &tim_roster_Rostertype{fieldName: "rostertype"}
	rostertype_.Field.FieldName = "rostertype"
	createtime := &tim_roster_Createtime{fieldName: "createtime"}
	createtime.Field.FieldName = "createtime"
	table := &Tim_roster{Username:username,Rostername:rostername,Rostertype:rostertype_,Createtime:createtime,Remarknick:remarknick,Id:id,Loginname:loginname}
	table.Table.ModifyMap = make(map[string]interface{})
	if len(tableName) == 1 {
		table.Table.TableName = tableName[0]
	} else {
		table.Table.TableName = "tim_roster"
	}
	return table
}
