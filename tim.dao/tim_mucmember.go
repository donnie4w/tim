package dao

/**
tablename:tim_mucmember
datetime :2016-06-10 00:59:53
*/
import (
	"github.com/donnie4w/gdao"
	"reflect"
)

type tim_mucmember_Id struct {
	gdao.Field
	fieldName  string
	FieldValue *int32
}

func (c *tim_mucmember_Id) Name() string {
	return c.fieldName
}

func (c *tim_mucmember_Id) Value() interface{} {
	return c.FieldValue
}

type tim_mucmember_Roomtid struct {
	gdao.Field
	fieldName  string
	FieldValue *string
}

func (c *tim_mucmember_Roomtid) Name() string {
	return c.fieldName
}

func (c *tim_mucmember_Roomtid) Value() interface{} {
	return c.FieldValue
}

type tim_mucmember_Domain struct {
	gdao.Field
	fieldName  string
	FieldValue *string
}

func (c *tim_mucmember_Domain) Name() string {
	return c.fieldName
}

func (c *tim_mucmember_Domain) Value() interface{} {
	return c.FieldValue
}

type tim_mucmember_Tidname struct {
	gdao.Field
	fieldName  string
	FieldValue *string
}

func (c *tim_mucmember_Tidname) Name() string {
	return c.fieldName
}

func (c *tim_mucmember_Tidname) Value() interface{} {
	return c.FieldValue
}

type tim_mucmember_Nickname struct {
	gdao.Field
	fieldName  string
	FieldValue *string
}

func (c *tim_mucmember_Nickname) Name() string {
	return c.fieldName
}

func (c *tim_mucmember_Nickname) Value() interface{} {
	return c.FieldValue
}

type tim_mucmember_Affiliation struct {
	gdao.Field
	fieldName  string
	FieldValue *int32
}

func (c *tim_mucmember_Affiliation) Name() string {
	return c.fieldName
}

func (c *tim_mucmember_Affiliation) Value() interface{} {
	return c.FieldValue
}

type tim_mucmember_Updatetime struct {
	gdao.Field
	fieldName  string
	FieldValue *string
}

func (c *tim_mucmember_Updatetime) Name() string {
	return c.fieldName
}

func (c *tim_mucmember_Updatetime) Value() interface{} {
	return c.FieldValue
}

type tim_mucmember_Createtime struct {
	gdao.Field
	fieldName  string
	FieldValue *string
}

func (c *tim_mucmember_Createtime) Name() string {
	return c.fieldName
}

func (c *tim_mucmember_Createtime) Value() interface{} {
	return c.FieldValue
}

type tim_mucmember_Type struct {
	gdao.Field
	fieldName  string
	FieldValue *int32
}

func (c *tim_mucmember_Type) Name() string {
	return c.fieldName
}

func (c *tim_mucmember_Type) Value() interface{} {
	return c.FieldValue
}

type Tim_mucmember struct {
	gdao.Table
	Type *tim_mucmember_Type
	Roomtid *tim_mucmember_Roomtid
	Domain *tim_mucmember_Domain
	Tidname *tim_mucmember_Tidname
	Nickname *tim_mucmember_Nickname
	Affiliation *tim_mucmember_Affiliation
	Updatetime *tim_mucmember_Updatetime
	Createtime *tim_mucmember_Createtime
	Id *tim_mucmember_Id
}

func (u *Tim_mucmember) GetAffiliation() int32 {
	return *u.Affiliation.FieldValue
}

func (u *Tim_mucmember) SetAffiliation(arg int64) {
	u.Table.ModifyMap[u.Affiliation.fieldName] = arg
	v := int32(arg)
	u.Affiliation.FieldValue = &v
}

func (u *Tim_mucmember) GetUpdatetime() string {
	return *u.Updatetime.FieldValue
}

func (u *Tim_mucmember) SetUpdatetime(arg string) {
	u.Table.ModifyMap[u.Updatetime.fieldName] = arg
	v := string(arg)
	u.Updatetime.FieldValue = &v
}

func (u *Tim_mucmember) GetCreatetime() string {
	return *u.Createtime.FieldValue
}

func (u *Tim_mucmember) SetCreatetime(arg string) {
	u.Table.ModifyMap[u.Createtime.fieldName] = arg
	v := string(arg)
	u.Createtime.FieldValue = &v
}

func (u *Tim_mucmember) GetId() int32 {
	return *u.Id.FieldValue
}

func (u *Tim_mucmember) SetId(arg int64) {
	u.Table.ModifyMap[u.Id.fieldName] = arg
	v := int32(arg)
	u.Id.FieldValue = &v
}

func (u *Tim_mucmember) GetRoomtid() string {
	return *u.Roomtid.FieldValue
}

func (u *Tim_mucmember) SetRoomtid(arg string) {
	u.Table.ModifyMap[u.Roomtid.fieldName] = arg
	v := string(arg)
	u.Roomtid.FieldValue = &v
}

func (u *Tim_mucmember) GetDomain() string {
	return *u.Domain.FieldValue
}

func (u *Tim_mucmember) SetDomain(arg string) {
	u.Table.ModifyMap[u.Domain.fieldName] = arg
	v := string(arg)
	u.Domain.FieldValue = &v
}

func (u *Tim_mucmember) GetTidname() string {
	return *u.Tidname.FieldValue
}

func (u *Tim_mucmember) SetTidname(arg string) {
	u.Table.ModifyMap[u.Tidname.fieldName] = arg
	v := string(arg)
	u.Tidname.FieldValue = &v
}

func (u *Tim_mucmember) GetNickname() string {
	return *u.Nickname.FieldValue
}

func (u *Tim_mucmember) SetNickname(arg string) {
	u.Table.ModifyMap[u.Nickname.fieldName] = arg
	v := string(arg)
	u.Nickname.FieldValue = &v
}

func (u *Tim_mucmember) GetType() int32 {
	return *u.Type.FieldValue
}

func (u *Tim_mucmember) SetType(arg int64) {
	u.Table.ModifyMap[u.Type.fieldName] = arg
	v := int32(arg)
	u.Type.FieldValue = &v
}

func (t *Tim_mucmember) Query(columns ...gdao.Column) ([]Tim_mucmember,error) {
	if columns == nil {
		columns = []gdao.Column{ t.Nickname,t.Affiliation,t.Updatetime,t.Createtime,t.Id,t.Roomtid,t.Domain,t.Tidname,t.Type}
	}
	rs,err := t.Table.Query(columns...)
	if rs == nil || err != nil {
		return nil, err
	}
	ts := make([]Tim_mucmember, 0, len(rs))
	c := make(chan int16,len(rs))
	for _, rows := range rs {
		t := NewTim_mucmember()
		go copyTim_mucmember(c, rows, t, columns)
		<-c
		ts = append(ts, *t)
	}
	return ts,nil
}

func copyTim_mucmember(channle chan int16, rows []interface{}, t *Tim_mucmember, columns []gdao.Column) {
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

func (t *Tim_mucmember) QuerySingle(columns ...gdao.Column) (*Tim_mucmember,error) {
	if columns == nil {
		columns = []gdao.Column{ t.Nickname,t.Affiliation,t.Updatetime,t.Createtime,t.Id,t.Roomtid,t.Domain,t.Tidname,t.Type}
	}
	rs,err := t.Table.QuerySingle(columns...)
	if rs == nil || err != nil {
		return nil, err
	}
	rt := NewTim_mucmember()
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

func (t *Tim_mucmember) Select(columns ...gdao.Column) (*Tim_mucmember,error) {
	if columns == nil {
		columns = []gdao.Column{ t.Nickname,t.Affiliation,t.Updatetime,t.Createtime,t.Id,t.Roomtid,t.Domain,t.Tidname,t.Type}
	}
	rows,err := t.Table.Selects(columns...)
	defer rows.Close()
	if err != nil || rows==nil {
		return nil, err
	}
	buff := make([]interface{}, len(columns))
	if rows.Next() {
		n := NewTim_mucmember()
		cpTim_mucmember(buff, n, columns)
		row_err := rows.Scan(buff...)
		if row_err != nil {
			return nil, row_err
		}
		return n, nil
	}
	return nil, nil
}

func (t *Tim_mucmember) Selects(columns ...gdao.Column) ([]*Tim_mucmember,error) {
	if columns == nil {
		columns = []gdao.Column{ t.Nickname,t.Affiliation,t.Updatetime,t.Createtime,t.Id,t.Roomtid,t.Domain,t.Tidname,t.Type}
	}
	rows,err := t.Table.Selects(columns...)
	defer rows.Close()
	if err != nil || rows==nil {
		return nil, err
	}
	ns := make([]*Tim_mucmember, 0)
	buff := make([]interface{}, len(columns))
	for rows.Next() {
		n := NewTim_mucmember()
		cpTim_mucmember(buff, n, columns)
		row_err := rows.Scan(buff...)
		if row_err != nil {
			return nil, row_err
		}
		ns = append(ns, n)
	}
	return ns, nil
}

func  cpTim_mucmember(buff []interface{}, t *Tim_mucmember, columns []gdao.Column) {
	for i, column := range columns {
		field := column.Name()
		switch field {
		case "createtime":
			buff[i] = &t.Createtime.FieldValue
		case "id":
			buff[i] = &t.Id.FieldValue
		case "roomtid":
			buff[i] = &t.Roomtid.FieldValue
		case "domain":
			buff[i] = &t.Domain.FieldValue
		case "tidname":
			buff[i] = &t.Tidname.FieldValue
		case "nickname":
			buff[i] = &t.Nickname.FieldValue
		case "affiliation":
			buff[i] = &t.Affiliation.FieldValue
		case "updatetime":
			buff[i] = &t.Updatetime.FieldValue
		case "type":
			buff[i] = &t.Type.FieldValue
		}
	}
}

func NewTim_mucmember(tableName ...string) *Tim_mucmember {
	type_ := &tim_mucmember_Type{fieldName: "type"}
	type_.Field.FieldName = "type"
	affiliation := &tim_mucmember_Affiliation{fieldName: "affiliation"}
	affiliation.Field.FieldName = "affiliation"
	updatetime := &tim_mucmember_Updatetime{fieldName: "updatetime"}
	updatetime.Field.FieldName = "updatetime"
	createtime := &tim_mucmember_Createtime{fieldName: "createtime"}
	createtime.Field.FieldName = "createtime"
	id := &tim_mucmember_Id{fieldName: "id"}
	id.Field.FieldName = "id"
	roomtid := &tim_mucmember_Roomtid{fieldName: "roomtid"}
	roomtid.Field.FieldName = "roomtid"
	domain := &tim_mucmember_Domain{fieldName: "domain"}
	domain.Field.FieldName = "domain"
	tidname := &tim_mucmember_Tidname{fieldName: "tidname"}
	tidname.Field.FieldName = "tidname"
	nickname := &tim_mucmember_Nickname{fieldName: "nickname"}
	nickname.Field.FieldName = "nickname"
	table := &Tim_mucmember{Type:type_,Nickname:nickname,Affiliation:affiliation,Updatetime:updatetime,Createtime:createtime,Id:id,Roomtid:roomtid,Domain:domain,Tidname:tidname}
	table.Table.ModifyMap = make(map[string]interface{})
	if len(tableName) == 1 {
		table.Table.TableName = tableName[0]
	} else {
		table.Table.TableName = "tim_mucmember"
	}
	return table
}
