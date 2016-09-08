package dao

/**
tablename:tim_mucroom
datetime :2016-09-07 11:32:22
*/
import (
	"github.com/donnie4w/gdao"
	"reflect"
)

type tim_mucroom_Description struct {
	gdao.Field
	fieldName  string
	FieldValue *string
}

func (c *tim_mucroom_Description) Name() string {
	return c.fieldName
}

func (c *tim_mucroom_Description) Value() interface{} {
	return c.FieldValue
}

type tim_mucroom_Updatetime struct {
	gdao.Field
	fieldName  string
	FieldValue *string
}

func (c *tim_mucroom_Updatetime) Name() string {
	return c.fieldName
}

func (c *tim_mucroom_Updatetime) Value() interface{} {
	return c.FieldValue
}

type tim_mucroom_Domain struct {
	gdao.Field
	fieldName  string
	FieldValue *string
}

func (c *tim_mucroom_Domain) Name() string {
	return c.fieldName
}

func (c *tim_mucroom_Domain) Value() interface{} {
	return c.FieldValue
}

type tim_mucroom_Maxusers struct {
	gdao.Field
	fieldName  string
	FieldValue *int32
}

func (c *tim_mucroom_Maxusers) Name() string {
	return c.fieldName
}

func (c *tim_mucroom_Maxusers) Value() interface{} {
	return c.FieldValue
}

type tim_mucroom_Theme struct {
	gdao.Field
	fieldName  string
	FieldValue *string
}

func (c *tim_mucroom_Theme) Name() string {
	return c.fieldName
}

func (c *tim_mucroom_Theme) Value() interface{} {
	return c.FieldValue
}

type tim_mucroom_Name struct {
	gdao.Field
	fieldName  string
	FieldValue *string
}

func (c *tim_mucroom_Name) Name() string {
	return c.fieldName
}

func (c *tim_mucroom_Name) Value() interface{} {
	return c.FieldValue
}

type tim_mucroom_Password struct {
	gdao.Field
	fieldName  string
	FieldValue *string
}

func (c *tim_mucroom_Password) Name() string {
	return c.fieldName
}

func (c *tim_mucroom_Password) Value() interface{} {
	return c.FieldValue
}

type tim_mucroom_Createtime struct {
	gdao.Field
	fieldName  string
	FieldValue *string
}

func (c *tim_mucroom_Createtime) Name() string {
	return c.fieldName
}

func (c *tim_mucroom_Createtime) Value() interface{} {
	return c.FieldValue
}

type tim_mucroom_Id struct {
	gdao.Field
	fieldName  string
	FieldValue *int32
}

func (c *tim_mucroom_Id) Name() string {
	return c.fieldName
}

func (c *tim_mucroom_Id) Value() interface{} {
	return c.FieldValue
}

type tim_mucroom_Roomtid struct {
	gdao.Field
	fieldName  string
	FieldValue *string
}

func (c *tim_mucroom_Roomtid) Name() string {
	return c.fieldName
}

func (c *tim_mucroom_Roomtid) Value() interface{} {
	return c.FieldValue
}

type Tim_mucroom struct {
	gdao.Table
	Id *tim_mucroom_Id
	Roomtid *tim_mucroom_Roomtid
	Theme *tim_mucroom_Theme
	Name *tim_mucroom_Name
	Password *tim_mucroom_Password
	Createtime *tim_mucroom_Createtime
	Domain *tim_mucroom_Domain
	Maxusers *tim_mucroom_Maxusers
	Description *tim_mucroom_Description
	Updatetime *tim_mucroom_Updatetime
}

func (u *Tim_mucroom) GetPassword() string {
	return *u.Password.FieldValue
}

func (u *Tim_mucroom) SetPassword(arg string) {
	u.Table.ModifyMap[u.Password.fieldName] = arg
	v := string(arg)
	u.Password.FieldValue = &v
}

func (u *Tim_mucroom) GetCreatetime() string {
	return *u.Createtime.FieldValue
}

func (u *Tim_mucroom) SetCreatetime(arg string) {
	u.Table.ModifyMap[u.Createtime.fieldName] = arg
	v := string(arg)
	u.Createtime.FieldValue = &v
}

func (u *Tim_mucroom) GetId() int32 {
	return *u.Id.FieldValue
}

func (u *Tim_mucroom) SetId(arg int64) {
	u.Table.ModifyMap[u.Id.fieldName] = arg
	v := int32(arg)
	u.Id.FieldValue = &v
}

func (u *Tim_mucroom) GetRoomtid() string {
	return *u.Roomtid.FieldValue
}

func (u *Tim_mucroom) SetRoomtid(arg string) {
	u.Table.ModifyMap[u.Roomtid.fieldName] = arg
	v := string(arg)
	u.Roomtid.FieldValue = &v
}

func (u *Tim_mucroom) GetTheme() string {
	return *u.Theme.FieldValue
}

func (u *Tim_mucroom) SetTheme(arg string) {
	u.Table.ModifyMap[u.Theme.fieldName] = arg
	v := string(arg)
	u.Theme.FieldValue = &v
}

func (u *Tim_mucroom) GetName() string {
	return *u.Name.FieldValue
}

func (u *Tim_mucroom) SetName(arg string) {
	u.Table.ModifyMap[u.Name.fieldName] = arg
	v := string(arg)
	u.Name.FieldValue = &v
}

func (u *Tim_mucroom) GetDomain() string {
	return *u.Domain.FieldValue
}

func (u *Tim_mucroom) SetDomain(arg string) {
	u.Table.ModifyMap[u.Domain.fieldName] = arg
	v := string(arg)
	u.Domain.FieldValue = &v
}

func (u *Tim_mucroom) GetMaxusers() int32 {
	return *u.Maxusers.FieldValue
}

func (u *Tim_mucroom) SetMaxusers(arg int64) {
	u.Table.ModifyMap[u.Maxusers.fieldName] = arg
	v := int32(arg)
	u.Maxusers.FieldValue = &v
}

func (u *Tim_mucroom) GetDescription() string {
	return *u.Description.FieldValue
}

func (u *Tim_mucroom) SetDescription(arg string) {
	u.Table.ModifyMap[u.Description.fieldName] = arg
	v := string(arg)
	u.Description.FieldValue = &v
}

func (u *Tim_mucroom) GetUpdatetime() string {
	return *u.Updatetime.FieldValue
}

func (u *Tim_mucroom) SetUpdatetime(arg string) {
	u.Table.ModifyMap[u.Updatetime.fieldName] = arg
	v := string(arg)
	u.Updatetime.FieldValue = &v
}

func (t *Tim_mucroom) Query(columns ...gdao.Column) ([]Tim_mucroom,error) {
	if columns == nil {
		columns = []gdao.Column{ t.Id,t.Roomtid,t.Theme,t.Name,t.Password,t.Createtime,t.Domain,t.Maxusers,t.Description,t.Updatetime}
	}
	rs,err := t.Table.Query(columns...)
	if rs == nil || err != nil {
		return nil, err
	}
	ts := make([]Tim_mucroom, 0, len(rs))
	c := make(chan int16,len(rs))
	for _, rows := range rs {
		t := NewTim_mucroom()
		go copyTim_mucroom(c, rows, t, columns)
		<-c
		ts = append(ts, *t)
	}
	return ts,nil
}

func copyTim_mucroom(channle chan int16, rows []interface{}, t *Tim_mucroom, columns []gdao.Column) {
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

func (t *Tim_mucroom) QuerySingle(columns ...gdao.Column) (*Tim_mucroom,error) {
	if columns == nil {
		columns = []gdao.Column{ t.Id,t.Roomtid,t.Theme,t.Name,t.Password,t.Createtime,t.Domain,t.Maxusers,t.Description,t.Updatetime}
	}
	rs,err := t.Table.QuerySingle(columns...)
	if rs == nil || err != nil {
		return nil, err
	}
	rt := NewTim_mucroom()
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

func (t *Tim_mucroom) Select(columns ...gdao.Column) (*Tim_mucroom,error) {
	if columns == nil {
		columns = []gdao.Column{ t.Id,t.Roomtid,t.Theme,t.Name,t.Password,t.Createtime,t.Domain,t.Maxusers,t.Description,t.Updatetime}
	}
	rows,err := t.Table.Selects(columns...)
	defer rows.Close()
	if err != nil || rows==nil {
		return nil, err
	}
	buff := make([]interface{}, len(columns))
	if rows.Next() {
		n := NewTim_mucroom()
		cpTim_mucroom(buff, n, columns)
		row_err := rows.Scan(buff...)
		if row_err != nil {
			return nil, row_err
		}
		return n, nil
	}
	return nil, nil
}

func (t *Tim_mucroom) Selects(columns ...gdao.Column) ([]*Tim_mucroom,error) {
	if columns == nil {
		columns = []gdao.Column{ t.Id,t.Roomtid,t.Theme,t.Name,t.Password,t.Createtime,t.Domain,t.Maxusers,t.Description,t.Updatetime}
	}
	rows,err := t.Table.Selects(columns...)
	defer rows.Close()
	if err != nil || rows==nil {
		return nil, err
	}
	ns := make([]*Tim_mucroom, 0)
	buff := make([]interface{}, len(columns))
	for rows.Next() {
		n := NewTim_mucroom()
		cpTim_mucroom(buff, n, columns)
		row_err := rows.Scan(buff...)
		if row_err != nil {
			return nil, row_err
		}
		ns = append(ns, n)
	}
	return ns, nil
}

func  cpTim_mucroom(buff []interface{}, t *Tim_mucroom, columns []gdao.Column) {
	for i, column := range columns {
		field := column.Name()
		switch field {
		case "maxusers":
			buff[i] = &t.Maxusers.FieldValue
		case "description":
			buff[i] = &t.Description.FieldValue
		case "updatetime":
			buff[i] = &t.Updatetime.FieldValue
		case "domain":
			buff[i] = &t.Domain.FieldValue
		case "roomtid":
			buff[i] = &t.Roomtid.FieldValue
		case "theme":
			buff[i] = &t.Theme.FieldValue
		case "name":
			buff[i] = &t.Name.FieldValue
		case "password":
			buff[i] = &t.Password.FieldValue
		case "createtime":
			buff[i] = &t.Createtime.FieldValue
		case "id":
			buff[i] = &t.Id.FieldValue
		}
	}
}

func NewTim_mucroom(tableName ...string) *Tim_mucroom {
	name := &tim_mucroom_Name{fieldName: "name"}
	name.Field.FieldName = "name"
	password := &tim_mucroom_Password{fieldName: "password"}
	password.Field.FieldName = "password"
	createtime := &tim_mucroom_Createtime{fieldName: "createtime"}
	createtime.Field.FieldName = "createtime"
	id := &tim_mucroom_Id{fieldName: "id"}
	id.Field.FieldName = "id"
	roomtid := &tim_mucroom_Roomtid{fieldName: "roomtid"}
	roomtid.Field.FieldName = "roomtid"
	theme := &tim_mucroom_Theme{fieldName: "theme"}
	theme.Field.FieldName = "theme"
	updatetime := &tim_mucroom_Updatetime{fieldName: "updatetime"}
	updatetime.Field.FieldName = "updatetime"
	domain := &tim_mucroom_Domain{fieldName: "domain"}
	domain.Field.FieldName = "domain"
	maxusers := &tim_mucroom_Maxusers{fieldName: "maxusers"}
	maxusers.Field.FieldName = "maxusers"
	description := &tim_mucroom_Description{fieldName: "description"}
	description.Field.FieldName = "description"
	table := &Tim_mucroom{Domain:domain,Maxusers:maxusers,Description:description,Updatetime:updatetime,Id:id,Roomtid:roomtid,Theme:theme,Name:name,Password:password,Createtime:createtime}
	table.Table.ModifyMap = make(map[string]interface{})
	if len(tableName) == 1 {
		table.Table.TableName = tableName[0]
	} else {
		table.Table.TableName = "tim_mucroom"
	}
	return table
}
