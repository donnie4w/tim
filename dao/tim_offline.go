package dao

/**
tablename:tim_offline
datetime :2016-09-07 11:32:22
*/
import (
	"github.com/donnie4w/gdao"
	"reflect"
)

type tim_offline_Domain struct {
	gdao.Field
	fieldName  string
	FieldValue *string
}

func (c *tim_offline_Domain) Name() string {
	return c.fieldName
}

func (c *tim_offline_Domain) Value() interface{} {
	return c.FieldValue
}

type tim_offline_Username struct {
	gdao.Field
	fieldName  string
	FieldValue *string
}

func (c *tim_offline_Username) Name() string {
	return c.fieldName
}

func (c *tim_offline_Username) Value() interface{} {
	return c.FieldValue
}

type tim_offline_Stamp struct {
	gdao.Field
	fieldName  string
	FieldValue *string
}

func (c *tim_offline_Stamp) Name() string {
	return c.fieldName
}

func (c *tim_offline_Stamp) Value() interface{} {
	return c.FieldValue
}

type tim_offline_Fromuser struct {
	gdao.Field
	fieldName  string
	FieldValue *string
}

func (c *tim_offline_Fromuser) Name() string {
	return c.fieldName
}

func (c *tim_offline_Fromuser) Value() interface{} {
	return c.FieldValue
}

type tim_offline_Gname struct {
	gdao.Field
	fieldName  string
	FieldValue *string
}

func (c *tim_offline_Gname) Name() string {
	return c.fieldName
}

func (c *tim_offline_Gname) Value() interface{} {
	return c.FieldValue
}

type tim_offline_Createtime struct {
	gdao.Field
	fieldName  string
	FieldValue *string
}

func (c *tim_offline_Createtime) Name() string {
	return c.fieldName
}

func (c *tim_offline_Createtime) Value() interface{} {
	return c.FieldValue
}

type tim_offline_Id struct {
	gdao.Field
	fieldName  string
	FieldValue *int32
}

func (c *tim_offline_Id) Name() string {
	return c.fieldName
}

func (c *tim_offline_Id) Value() interface{} {
	return c.FieldValue
}

type tim_offline_Mid struct {
	gdao.Field
	fieldName  string
	FieldValue *int32
}

func (c *tim_offline_Mid) Name() string {
	return c.fieldName
}

func (c *tim_offline_Mid) Value() interface{} {
	return c.FieldValue
}

type tim_offline_Msgtype struct {
	gdao.Field
	fieldName  string
	FieldValue *int32
}

func (c *tim_offline_Msgtype) Name() string {
	return c.fieldName
}

func (c *tim_offline_Msgtype) Value() interface{} {
	return c.FieldValue
}

type tim_offline_Msgmode struct {
	gdao.Field
	fieldName  string
	FieldValue *int32
}

func (c *tim_offline_Msgmode) Name() string {
	return c.fieldName
}

func (c *tim_offline_Msgmode) Value() interface{} {
	return c.FieldValue
}

type tim_offline_Message_size struct {
	gdao.Field
	fieldName  string
	FieldValue *int32
}

func (c *tim_offline_Message_size) Name() string {
	return c.fieldName
}

func (c *tim_offline_Message_size) Value() interface{} {
	return c.FieldValue
}

type tim_offline_Stanza struct {
	gdao.Field
	fieldName  string
	FieldValue *string
}

func (c *tim_offline_Stanza) Name() string {
	return c.fieldName
}

func (c *tim_offline_Stanza) Value() interface{} {
	return c.FieldValue
}

type Tim_offline struct {
	gdao.Table
	Domain *tim_offline_Domain
	Username *tim_offline_Username
	Stamp *tim_offline_Stamp
	Fromuser *tim_offline_Fromuser
	Gname *tim_offline_Gname
	Createtime *tim_offline_Createtime
	Id *tim_offline_Id
	Mid *tim_offline_Mid
	Msgtype *tim_offline_Msgtype
	Msgmode *tim_offline_Msgmode
	Message_size *tim_offline_Message_size
	Stanza *tim_offline_Stanza
}

func (u *Tim_offline) GetStanza() string {
	return *u.Stanza.FieldValue
}

func (u *Tim_offline) SetStanza(arg string) {
	u.Table.ModifyMap[u.Stanza.fieldName] = arg
	v := string(arg)
	u.Stanza.FieldValue = &v
}

func (u *Tim_offline) GetId() int32 {
	return *u.Id.FieldValue
}

func (u *Tim_offline) SetId(arg int64) {
	u.Table.ModifyMap[u.Id.fieldName] = arg
	v := int32(arg)
	u.Id.FieldValue = &v
}

func (u *Tim_offline) GetMid() int32 {
	return *u.Mid.FieldValue
}

func (u *Tim_offline) SetMid(arg int64) {
	u.Table.ModifyMap[u.Mid.fieldName] = arg
	v := int32(arg)
	u.Mid.FieldValue = &v
}

func (u *Tim_offline) GetMsgtype() int32 {
	return *u.Msgtype.FieldValue
}

func (u *Tim_offline) SetMsgtype(arg int64) {
	u.Table.ModifyMap[u.Msgtype.fieldName] = arg
	v := int32(arg)
	u.Msgtype.FieldValue = &v
}

func (u *Tim_offline) GetMsgmode() int32 {
	return *u.Msgmode.FieldValue
}

func (u *Tim_offline) SetMsgmode(arg int64) {
	u.Table.ModifyMap[u.Msgmode.fieldName] = arg
	v := int32(arg)
	u.Msgmode.FieldValue = &v
}

func (u *Tim_offline) GetMessage_size() int32 {
	return *u.Message_size.FieldValue
}

func (u *Tim_offline) SetMessage_size(arg int64) {
	u.Table.ModifyMap[u.Message_size.fieldName] = arg
	v := int32(arg)
	u.Message_size.FieldValue = &v
}

func (u *Tim_offline) GetCreatetime() string {
	return *u.Createtime.FieldValue
}

func (u *Tim_offline) SetCreatetime(arg string) {
	u.Table.ModifyMap[u.Createtime.fieldName] = arg
	v := string(arg)
	u.Createtime.FieldValue = &v
}

func (u *Tim_offline) GetDomain() string {
	return *u.Domain.FieldValue
}

func (u *Tim_offline) SetDomain(arg string) {
	u.Table.ModifyMap[u.Domain.fieldName] = arg
	v := string(arg)
	u.Domain.FieldValue = &v
}

func (u *Tim_offline) GetUsername() string {
	return *u.Username.FieldValue
}

func (u *Tim_offline) SetUsername(arg string) {
	u.Table.ModifyMap[u.Username.fieldName] = arg
	v := string(arg)
	u.Username.FieldValue = &v
}

func (u *Tim_offline) GetStamp() string {
	return *u.Stamp.FieldValue
}

func (u *Tim_offline) SetStamp(arg string) {
	u.Table.ModifyMap[u.Stamp.fieldName] = arg
	v := string(arg)
	u.Stamp.FieldValue = &v
}

func (u *Tim_offline) GetFromuser() string {
	return *u.Fromuser.FieldValue
}

func (u *Tim_offline) SetFromuser(arg string) {
	u.Table.ModifyMap[u.Fromuser.fieldName] = arg
	v := string(arg)
	u.Fromuser.FieldValue = &v
}

func (u *Tim_offline) GetGname() string {
	return *u.Gname.FieldValue
}

func (u *Tim_offline) SetGname(arg string) {
	u.Table.ModifyMap[u.Gname.fieldName] = arg
	v := string(arg)
	u.Gname.FieldValue = &v
}

func (t *Tim_offline) Query(columns ...gdao.Column) ([]Tim_offline,error) {
	if columns == nil {
		columns = []gdao.Column{ t.Message_size,t.Stanza,t.Id,t.Mid,t.Msgtype,t.Msgmode,t.Gname,t.Createtime,t.Domain,t.Username,t.Stamp,t.Fromuser}
	}
	rs,err := t.Table.Query(columns...)
	if rs == nil || err != nil {
		return nil, err
	}
	ts := make([]Tim_offline, 0, len(rs))
	c := make(chan int16,len(rs))
	for _, rows := range rs {
		t := NewTim_offline()
		go copyTim_offline(c, rows, t, columns)
		<-c
		ts = append(ts, *t)
	}
	return ts,nil
}

func copyTim_offline(channle chan int16, rows []interface{}, t *Tim_offline, columns []gdao.Column) {
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

func (t *Tim_offline) QuerySingle(columns ...gdao.Column) (*Tim_offline,error) {
	if columns == nil {
		columns = []gdao.Column{ t.Message_size,t.Stanza,t.Id,t.Mid,t.Msgtype,t.Msgmode,t.Gname,t.Createtime,t.Domain,t.Username,t.Stamp,t.Fromuser}
	}
	rs,err := t.Table.QuerySingle(columns...)
	if rs == nil || err != nil {
		return nil, err
	}
	rt := NewTim_offline()
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

func (t *Tim_offline) Select(columns ...gdao.Column) (*Tim_offline,error) {
	if columns == nil {
		columns = []gdao.Column{ t.Message_size,t.Stanza,t.Id,t.Mid,t.Msgtype,t.Msgmode,t.Gname,t.Createtime,t.Domain,t.Username,t.Stamp,t.Fromuser}
	}
	rows,err := t.Table.Selects(columns...)
	defer rows.Close()
	if err != nil || rows==nil {
		return nil, err
	}
	buff := make([]interface{}, len(columns))
	if rows.Next() {
		n := NewTim_offline()
		cpTim_offline(buff, n, columns)
		row_err := rows.Scan(buff...)
		if row_err != nil {
			return nil, row_err
		}
		return n, nil
	}
	return nil, nil
}

func (t *Tim_offline) Selects(columns ...gdao.Column) ([]*Tim_offline,error) {
	if columns == nil {
		columns = []gdao.Column{ t.Message_size,t.Stanza,t.Id,t.Mid,t.Msgtype,t.Msgmode,t.Gname,t.Createtime,t.Domain,t.Username,t.Stamp,t.Fromuser}
	}
	rows,err := t.Table.Selects(columns...)
	defer rows.Close()
	if err != nil || rows==nil {
		return nil, err
	}
	ns := make([]*Tim_offline, 0)
	buff := make([]interface{}, len(columns))
	for rows.Next() {
		n := NewTim_offline()
		cpTim_offline(buff, n, columns)
		row_err := rows.Scan(buff...)
		if row_err != nil {
			return nil, row_err
		}
		ns = append(ns, n)
	}
	return ns, nil
}

func  cpTim_offline(buff []interface{}, t *Tim_offline, columns []gdao.Column) {
	for i, column := range columns {
		field := column.Name()
		switch field {
		case "msgtype":
			buff[i] = &t.Msgtype.FieldValue
		case "msgmode":
			buff[i] = &t.Msgmode.FieldValue
		case "message_size":
			buff[i] = &t.Message_size.FieldValue
		case "stanza":
			buff[i] = &t.Stanza.FieldValue
		case "id":
			buff[i] = &t.Id.FieldValue
		case "mid":
			buff[i] = &t.Mid.FieldValue
		case "stamp":
			buff[i] = &t.Stamp.FieldValue
		case "fromuser":
			buff[i] = &t.Fromuser.FieldValue
		case "gname":
			buff[i] = &t.Gname.FieldValue
		case "createtime":
			buff[i] = &t.Createtime.FieldValue
		case "domain":
			buff[i] = &t.Domain.FieldValue
		case "username":
			buff[i] = &t.Username.FieldValue
		}
	}
}

func NewTim_offline(tableName ...string) *Tim_offline {
	message_size := &tim_offline_Message_size{fieldName: "message_size"}
	message_size.Field.FieldName = "message_size"
	stanza := &tim_offline_Stanza{fieldName: "stanza"}
	stanza.Field.FieldName = "stanza"
	id := &tim_offline_Id{fieldName: "id"}
	id.Field.FieldName = "id"
	mid := &tim_offline_Mid{fieldName: "mid"}
	mid.Field.FieldName = "mid"
	msgtype_ := &tim_offline_Msgtype{fieldName: "msgtype"}
	msgtype_.Field.FieldName = "msgtype"
	msgmode := &tim_offline_Msgmode{fieldName: "msgmode"}
	msgmode.Field.FieldName = "msgmode"
	gname := &tim_offline_Gname{fieldName: "gname"}
	gname.Field.FieldName = "gname"
	createtime := &tim_offline_Createtime{fieldName: "createtime"}
	createtime.Field.FieldName = "createtime"
	domain := &tim_offline_Domain{fieldName: "domain"}
	domain.Field.FieldName = "domain"
	username := &tim_offline_Username{fieldName: "username"}
	username.Field.FieldName = "username"
	stamp := &tim_offline_Stamp{fieldName: "stamp"}
	stamp.Field.FieldName = "stamp"
	fromuser := &tim_offline_Fromuser{fieldName: "fromuser"}
	fromuser.Field.FieldName = "fromuser"
	table := &Tim_offline{Stamp:stamp,Fromuser:fromuser,Gname:gname,Createtime:createtime,Domain:domain,Username:username,Msgtype:msgtype_,Msgmode:msgmode,Message_size:message_size,Stanza:stanza,Id:id,Mid:mid}
	table.Table.ModifyMap = make(map[string]interface{})
	if len(tableName) == 1 {
		table.Table.TableName = tableName[0]
	} else {
		table.Table.TableName = "tim_offline"
	}
	return table
}
