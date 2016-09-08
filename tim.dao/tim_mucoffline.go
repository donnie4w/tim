package dao

/**
tablename:tim_mucoffline
datetime :2016-09-07 11:32:22
*/
import (
	"github.com/donnie4w/gdao"
	"reflect"
)

type tim_mucoffline_Msgtype struct {
	gdao.Field
	fieldName  string
	FieldValue *int32
}

func (c *tim_mucoffline_Msgtype) Name() string {
	return c.fieldName
}

func (c *tim_mucoffline_Msgtype) Value() interface{} {
	return c.FieldValue
}

type tim_mucoffline_Message_size struct {
	gdao.Field
	fieldName  string
	FieldValue *int32
}

func (c *tim_mucoffline_Message_size) Name() string {
	return c.fieldName
}

func (c *tim_mucoffline_Message_size) Value() interface{} {
	return c.FieldValue
}

type tim_mucoffline_Createtime struct {
	gdao.Field
	fieldName  string
	FieldValue *string
}

func (c *tim_mucoffline_Createtime) Name() string {
	return c.fieldName
}

func (c *tim_mucoffline_Createtime) Value() interface{} {
	return c.FieldValue
}

type tim_mucoffline_Id struct {
	gdao.Field
	fieldName  string
	FieldValue *int32
}

func (c *tim_mucoffline_Id) Name() string {
	return c.fieldName
}

func (c *tim_mucoffline_Id) Value() interface{} {
	return c.FieldValue
}

type tim_mucoffline_Mid struct {
	gdao.Field
	fieldName  string
	FieldValue *int32
}

func (c *tim_mucoffline_Mid) Name() string {
	return c.fieldName
}

func (c *tim_mucoffline_Mid) Value() interface{} {
	return c.FieldValue
}

type tim_mucoffline_Domain struct {
	gdao.Field
	fieldName  string
	FieldValue *string
}

func (c *tim_mucoffline_Domain) Name() string {
	return c.fieldName
}

func (c *tim_mucoffline_Domain) Value() interface{} {
	return c.FieldValue
}

type tim_mucoffline_Username struct {
	gdao.Field
	fieldName  string
	FieldValue *string
}

func (c *tim_mucoffline_Username) Name() string {
	return c.fieldName
}

func (c *tim_mucoffline_Username) Value() interface{} {
	return c.FieldValue
}

type tim_mucoffline_Stamp struct {
	gdao.Field
	fieldName  string
	FieldValue *string
}

func (c *tim_mucoffline_Stamp) Name() string {
	return c.fieldName
}

func (c *tim_mucoffline_Stamp) Value() interface{} {
	return c.FieldValue
}

type tim_mucoffline_Roomid struct {
	gdao.Field
	fieldName  string
	FieldValue *string
}

func (c *tim_mucoffline_Roomid) Name() string {
	return c.fieldName
}

func (c *tim_mucoffline_Roomid) Value() interface{} {
	return c.FieldValue
}

type Tim_mucoffline struct {
	gdao.Table
	Msgtype *tim_mucoffline_Msgtype
	Message_size *tim_mucoffline_Message_size
	Domain *tim_mucoffline_Domain
	Username *tim_mucoffline_Username
	Stamp *tim_mucoffline_Stamp
	Roomid *tim_mucoffline_Roomid
	Createtime *tim_mucoffline_Createtime
	Id *tim_mucoffline_Id
	Mid *tim_mucoffline_Mid
}

func (u *Tim_mucoffline) GetMsgtype() int32 {
	return *u.Msgtype.FieldValue
}

func (u *Tim_mucoffline) SetMsgtype(arg int64) {
	u.Table.ModifyMap[u.Msgtype.fieldName] = arg
	v := int32(arg)
	u.Msgtype.FieldValue = &v
}

func (u *Tim_mucoffline) GetMessage_size() int32 {
	return *u.Message_size.FieldValue
}

func (u *Tim_mucoffline) SetMessage_size(arg int64) {
	u.Table.ModifyMap[u.Message_size.fieldName] = arg
	v := int32(arg)
	u.Message_size.FieldValue = &v
}

func (u *Tim_mucoffline) GetDomain() string {
	return *u.Domain.FieldValue
}

func (u *Tim_mucoffline) SetDomain(arg string) {
	u.Table.ModifyMap[u.Domain.fieldName] = arg
	v := string(arg)
	u.Domain.FieldValue = &v
}

func (u *Tim_mucoffline) GetUsername() string {
	return *u.Username.FieldValue
}

func (u *Tim_mucoffline) SetUsername(arg string) {
	u.Table.ModifyMap[u.Username.fieldName] = arg
	v := string(arg)
	u.Username.FieldValue = &v
}

func (u *Tim_mucoffline) GetStamp() string {
	return *u.Stamp.FieldValue
}

func (u *Tim_mucoffline) SetStamp(arg string) {
	u.Table.ModifyMap[u.Stamp.fieldName] = arg
	v := string(arg)
	u.Stamp.FieldValue = &v
}

func (u *Tim_mucoffline) GetRoomid() string {
	return *u.Roomid.FieldValue
}

func (u *Tim_mucoffline) SetRoomid(arg string) {
	u.Table.ModifyMap[u.Roomid.fieldName] = arg
	v := string(arg)
	u.Roomid.FieldValue = &v
}

func (u *Tim_mucoffline) GetCreatetime() string {
	return *u.Createtime.FieldValue
}

func (u *Tim_mucoffline) SetCreatetime(arg string) {
	u.Table.ModifyMap[u.Createtime.fieldName] = arg
	v := string(arg)
	u.Createtime.FieldValue = &v
}

func (u *Tim_mucoffline) GetId() int32 {
	return *u.Id.FieldValue
}

func (u *Tim_mucoffline) SetId(arg int64) {
	u.Table.ModifyMap[u.Id.fieldName] = arg
	v := int32(arg)
	u.Id.FieldValue = &v
}

func (u *Tim_mucoffline) GetMid() int32 {
	return *u.Mid.FieldValue
}

func (u *Tim_mucoffline) SetMid(arg int64) {
	u.Table.ModifyMap[u.Mid.fieldName] = arg
	v := int32(arg)
	u.Mid.FieldValue = &v
}

func (t *Tim_mucoffline) Query(columns ...gdao.Column) ([]Tim_mucoffline,error) {
	if columns == nil {
		columns = []gdao.Column{ t.Msgtype,t.Message_size,t.Stamp,t.Roomid,t.Createtime,t.Id,t.Mid,t.Domain,t.Username}
	}
	rs,err := t.Table.Query(columns...)
	if rs == nil || err != nil {
		return nil, err
	}
	ts := make([]Tim_mucoffline, 0, len(rs))
	c := make(chan int16,len(rs))
	for _, rows := range rs {
		t := NewTim_mucoffline()
		go copyTim_mucoffline(c, rows, t, columns)
		<-c
		ts = append(ts, *t)
	}
	return ts,nil
}

func copyTim_mucoffline(channle chan int16, rows []interface{}, t *Tim_mucoffline, columns []gdao.Column) {
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

func (t *Tim_mucoffline) QuerySingle(columns ...gdao.Column) (*Tim_mucoffline,error) {
	if columns == nil {
		columns = []gdao.Column{ t.Msgtype,t.Message_size,t.Stamp,t.Roomid,t.Createtime,t.Id,t.Mid,t.Domain,t.Username}
	}
	rs,err := t.Table.QuerySingle(columns...)
	if rs == nil || err != nil {
		return nil, err
	}
	rt := NewTim_mucoffline()
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

func (t *Tim_mucoffline) Select(columns ...gdao.Column) (*Tim_mucoffline,error) {
	if columns == nil {
		columns = []gdao.Column{ t.Msgtype,t.Message_size,t.Stamp,t.Roomid,t.Createtime,t.Id,t.Mid,t.Domain,t.Username}
	}
	rows,err := t.Table.Selects(columns...)
	defer rows.Close()
	if err != nil || rows==nil {
		return nil, err
	}
	buff := make([]interface{}, len(columns))
	if rows.Next() {
		n := NewTim_mucoffline()
		cpTim_mucoffline(buff, n, columns)
		row_err := rows.Scan(buff...)
		if row_err != nil {
			return nil, row_err
		}
		return n, nil
	}
	return nil, nil
}

func (t *Tim_mucoffline) Selects(columns ...gdao.Column) ([]*Tim_mucoffline,error) {
	if columns == nil {
		columns = []gdao.Column{ t.Msgtype,t.Message_size,t.Stamp,t.Roomid,t.Createtime,t.Id,t.Mid,t.Domain,t.Username}
	}
	rows,err := t.Table.Selects(columns...)
	defer rows.Close()
	if err != nil || rows==nil {
		return nil, err
	}
	ns := make([]*Tim_mucoffline, 0)
	buff := make([]interface{}, len(columns))
	for rows.Next() {
		n := NewTim_mucoffline()
		cpTim_mucoffline(buff, n, columns)
		row_err := rows.Scan(buff...)
		if row_err != nil {
			return nil, row_err
		}
		ns = append(ns, n)
	}
	return ns, nil
}

func  cpTim_mucoffline(buff []interface{}, t *Tim_mucoffline, columns []gdao.Column) {
	for i, column := range columns {
		field := column.Name()
		switch field {
		case "mid":
			buff[i] = &t.Mid.FieldValue
		case "domain":
			buff[i] = &t.Domain.FieldValue
		case "username":
			buff[i] = &t.Username.FieldValue
		case "stamp":
			buff[i] = &t.Stamp.FieldValue
		case "roomid":
			buff[i] = &t.Roomid.FieldValue
		case "createtime":
			buff[i] = &t.Createtime.FieldValue
		case "id":
			buff[i] = &t.Id.FieldValue
		case "message_size":
			buff[i] = &t.Message_size.FieldValue
		case "msgtype":
			buff[i] = &t.Msgtype.FieldValue
		}
	}
}

func NewTim_mucoffline(tableName ...string) *Tim_mucoffline {
	msgtype_ := &tim_mucoffline_Msgtype{fieldName: "msgtype"}
	msgtype_.Field.FieldName = "msgtype"
	message_size := &tim_mucoffline_Message_size{fieldName: "message_size"}
	message_size.Field.FieldName = "message_size"
	domain := &tim_mucoffline_Domain{fieldName: "domain"}
	domain.Field.FieldName = "domain"
	username := &tim_mucoffline_Username{fieldName: "username"}
	username.Field.FieldName = "username"
	stamp := &tim_mucoffline_Stamp{fieldName: "stamp"}
	stamp.Field.FieldName = "stamp"
	roomid := &tim_mucoffline_Roomid{fieldName: "roomid"}
	roomid.Field.FieldName = "roomid"
	createtime := &tim_mucoffline_Createtime{fieldName: "createtime"}
	createtime.Field.FieldName = "createtime"
	id := &tim_mucoffline_Id{fieldName: "id"}
	id.Field.FieldName = "id"
	mid := &tim_mucoffline_Mid{fieldName: "mid"}
	mid.Field.FieldName = "mid"
	table := &Tim_mucoffline{Msgtype:msgtype_,Message_size:message_size,Roomid:roomid,Createtime:createtime,Id:id,Mid:mid,Domain:domain,Username:username,Stamp:stamp}
	table.Table.ModifyMap = make(map[string]interface{})
	if len(tableName) == 1 {
		table.Table.TableName = tableName[0]
	} else {
		table.Table.TableName = "tim_mucoffline"
	}
	return table
}
