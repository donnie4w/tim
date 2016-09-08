package dao

/**
tablename:tim_message
datetime :2016-09-07 11:32:22
*/
import (
	"github.com/donnie4w/gdao"
	"reflect"
)

type tim_message_Stanza struct {
	gdao.Field
	fieldName  string
	FieldValue *string
}

func (c *tim_message_Stanza) Name() string {
	return c.fieldName
}

func (c *tim_message_Stanza) Value() interface{} {
	return c.FieldValue
}

type tim_message_Createtime struct {
	gdao.Field
	fieldName  string
	FieldValue *string
}

func (c *tim_message_Createtime) Name() string {
	return c.fieldName
}

func (c *tim_message_Createtime) Value() interface{} {
	return c.FieldValue
}

type tim_message_Chatid struct {
	gdao.Field
	fieldName  string
	FieldValue *string
}

func (c *tim_message_Chatid) Name() string {
	return c.fieldName
}

func (c *tim_message_Chatid) Value() interface{} {
	return c.FieldValue
}

type tim_message_Fromuser struct {
	gdao.Field
	fieldName  string
	FieldValue *string
}

func (c *tim_message_Fromuser) Name() string {
	return c.fieldName
}

func (c *tim_message_Fromuser) Value() interface{} {
	return c.FieldValue
}

type tim_message_Msgtype struct {
	gdao.Field
	fieldName  string
	FieldValue *int32
}

func (c *tim_message_Msgtype) Name() string {
	return c.fieldName
}

func (c *tim_message_Msgtype) Value() interface{} {
	return c.FieldValue
}

type tim_message_Gname struct {
	gdao.Field
	fieldName  string
	FieldValue *string
}

func (c *tim_message_Gname) Name() string {
	return c.fieldName
}

func (c *tim_message_Gname) Value() interface{} {
	return c.FieldValue
}

type tim_message_Small struct {
	gdao.Field
	fieldName  string
	FieldValue *int32
}

func (c *tim_message_Small) Name() string {
	return c.fieldName
}

func (c *tim_message_Small) Value() interface{} {
	return c.FieldValue
}

type tim_message_Large struct {
	gdao.Field
	fieldName  string
	FieldValue *int32
}

func (c *tim_message_Large) Name() string {
	return c.fieldName
}

func (c *tim_message_Large) Value() interface{} {
	return c.FieldValue
}

type tim_message_Id struct {
	gdao.Field
	fieldName  string
	FieldValue *int32
}

func (c *tim_message_Id) Name() string {
	return c.fieldName
}

func (c *tim_message_Id) Value() interface{} {
	return c.FieldValue
}

type tim_message_Stamp struct {
	gdao.Field
	fieldName  string
	FieldValue *string
}

func (c *tim_message_Stamp) Name() string {
	return c.fieldName
}

func (c *tim_message_Stamp) Value() interface{} {
	return c.FieldValue
}

type tim_message_Touser struct {
	gdao.Field
	fieldName  string
	FieldValue *string
}

func (c *tim_message_Touser) Name() string {
	return c.fieldName
}

func (c *tim_message_Touser) Value() interface{} {
	return c.FieldValue
}

type tim_message_Msgmode struct {
	gdao.Field
	fieldName  string
	FieldValue *int32
}

func (c *tim_message_Msgmode) Name() string {
	return c.fieldName
}

func (c *tim_message_Msgmode) Value() interface{} {
	return c.FieldValue
}

type Tim_message struct {
	gdao.Table
	Stamp *tim_message_Stamp
	Touser *tim_message_Touser
	Msgmode *tim_message_Msgmode
	Small *tim_message_Small
	Large *tim_message_Large
	Id *tim_message_Id
	Fromuser *tim_message_Fromuser
	Msgtype *tim_message_Msgtype
	Gname *tim_message_Gname
	Stanza *tim_message_Stanza
	Createtime *tim_message_Createtime
	Chatid *tim_message_Chatid
}

func (u *Tim_message) GetCreatetime() string {
	return *u.Createtime.FieldValue
}

func (u *Tim_message) SetCreatetime(arg string) {
	u.Table.ModifyMap[u.Createtime.fieldName] = arg
	v := string(arg)
	u.Createtime.FieldValue = &v
}

func (u *Tim_message) GetChatid() string {
	return *u.Chatid.FieldValue
}

func (u *Tim_message) SetChatid(arg string) {
	u.Table.ModifyMap[u.Chatid.fieldName] = arg
	v := string(arg)
	u.Chatid.FieldValue = &v
}

func (u *Tim_message) GetFromuser() string {
	return *u.Fromuser.FieldValue
}

func (u *Tim_message) SetFromuser(arg string) {
	u.Table.ModifyMap[u.Fromuser.fieldName] = arg
	v := string(arg)
	u.Fromuser.FieldValue = &v
}

func (u *Tim_message) GetMsgtype() int32 {
	return *u.Msgtype.FieldValue
}

func (u *Tim_message) SetMsgtype(arg int64) {
	u.Table.ModifyMap[u.Msgtype.fieldName] = arg
	v := int32(arg)
	u.Msgtype.FieldValue = &v
}

func (u *Tim_message) GetGname() string {
	return *u.Gname.FieldValue
}

func (u *Tim_message) SetGname(arg string) {
	u.Table.ModifyMap[u.Gname.fieldName] = arg
	v := string(arg)
	u.Gname.FieldValue = &v
}

func (u *Tim_message) GetStanza() string {
	return *u.Stanza.FieldValue
}

func (u *Tim_message) SetStanza(arg string) {
	u.Table.ModifyMap[u.Stanza.fieldName] = arg
	v := string(arg)
	u.Stanza.FieldValue = &v
}

func (u *Tim_message) GetLarge() int32 {
	return *u.Large.FieldValue
}

func (u *Tim_message) SetLarge(arg int64) {
	u.Table.ModifyMap[u.Large.fieldName] = arg
	v := int32(arg)
	u.Large.FieldValue = &v
}

func (u *Tim_message) GetId() int32 {
	return *u.Id.FieldValue
}

func (u *Tim_message) SetId(arg int64) {
	u.Table.ModifyMap[u.Id.fieldName] = arg
	v := int32(arg)
	u.Id.FieldValue = &v
}

func (u *Tim_message) GetStamp() string {
	return *u.Stamp.FieldValue
}

func (u *Tim_message) SetStamp(arg string) {
	u.Table.ModifyMap[u.Stamp.fieldName] = arg
	v := string(arg)
	u.Stamp.FieldValue = &v
}

func (u *Tim_message) GetTouser() string {
	return *u.Touser.FieldValue
}

func (u *Tim_message) SetTouser(arg string) {
	u.Table.ModifyMap[u.Touser.fieldName] = arg
	v := string(arg)
	u.Touser.FieldValue = &v
}

func (u *Tim_message) GetMsgmode() int32 {
	return *u.Msgmode.FieldValue
}

func (u *Tim_message) SetMsgmode(arg int64) {
	u.Table.ModifyMap[u.Msgmode.fieldName] = arg
	v := int32(arg)
	u.Msgmode.FieldValue = &v
}

func (u *Tim_message) GetSmall() int32 {
	return *u.Small.FieldValue
}

func (u *Tim_message) SetSmall(arg int64) {
	u.Table.ModifyMap[u.Small.fieldName] = arg
	v := int32(arg)
	u.Small.FieldValue = &v
}

func (t *Tim_message) Query(columns ...gdao.Column) ([]Tim_message,error) {
	if columns == nil {
		columns = []gdao.Column{ t.Stanza,t.Createtime,t.Chatid,t.Fromuser,t.Msgtype,t.Gname,t.Small,t.Large,t.Id,t.Stamp,t.Touser,t.Msgmode}
	}
	rs,err := t.Table.Query(columns...)
	if rs == nil || err != nil {
		return nil, err
	}
	ts := make([]Tim_message, 0, len(rs))
	c := make(chan int16,len(rs))
	for _, rows := range rs {
		t := NewTim_message()
		go copyTim_message(c, rows, t, columns)
		<-c
		ts = append(ts, *t)
	}
	return ts,nil
}

func copyTim_message(channle chan int16, rows []interface{}, t *Tim_message, columns []gdao.Column) {
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

func (t *Tim_message) QuerySingle(columns ...gdao.Column) (*Tim_message,error) {
	if columns == nil {
		columns = []gdao.Column{ t.Stanza,t.Createtime,t.Chatid,t.Fromuser,t.Msgtype,t.Gname,t.Small,t.Large,t.Id,t.Stamp,t.Touser,t.Msgmode}
	}
	rs,err := t.Table.QuerySingle(columns...)
	if rs == nil || err != nil {
		return nil, err
	}
	rt := NewTim_message()
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

func (t *Tim_message) Select(columns ...gdao.Column) (*Tim_message,error) {
	if columns == nil {
		columns = []gdao.Column{ t.Stanza,t.Createtime,t.Chatid,t.Fromuser,t.Msgtype,t.Gname,t.Small,t.Large,t.Id,t.Stamp,t.Touser,t.Msgmode}
	}
	rows,err := t.Table.Selects(columns...)
	defer rows.Close()
	if err != nil || rows==nil {
		return nil, err
	}
	buff := make([]interface{}, len(columns))
	if rows.Next() {
		n := NewTim_message()
		cpTim_message(buff, n, columns)
		row_err := rows.Scan(buff...)
		if row_err != nil {
			return nil, row_err
		}
		return n, nil
	}
	return nil, nil
}

func (t *Tim_message) Selects(columns ...gdao.Column) ([]*Tim_message,error) {
	if columns == nil {
		columns = []gdao.Column{ t.Stanza,t.Createtime,t.Chatid,t.Fromuser,t.Msgtype,t.Gname,t.Small,t.Large,t.Id,t.Stamp,t.Touser,t.Msgmode}
	}
	rows,err := t.Table.Selects(columns...)
	defer rows.Close()
	if err != nil || rows==nil {
		return nil, err
	}
	ns := make([]*Tim_message, 0)
	buff := make([]interface{}, len(columns))
	for rows.Next() {
		n := NewTim_message()
		cpTim_message(buff, n, columns)
		row_err := rows.Scan(buff...)
		if row_err != nil {
			return nil, row_err
		}
		ns = append(ns, n)
	}
	return ns, nil
}

func  cpTim_message(buff []interface{}, t *Tim_message, columns []gdao.Column) {
	for i, column := range columns {
		field := column.Name()
		switch field {
		case "chatid":
			buff[i] = &t.Chatid.FieldValue
		case "fromuser":
			buff[i] = &t.Fromuser.FieldValue
		case "msgtype":
			buff[i] = &t.Msgtype.FieldValue
		case "gname":
			buff[i] = &t.Gname.FieldValue
		case "stanza":
			buff[i] = &t.Stanza.FieldValue
		case "createtime":
			buff[i] = &t.Createtime.FieldValue
		case "id":
			buff[i] = &t.Id.FieldValue
		case "stamp":
			buff[i] = &t.Stamp.FieldValue
		case "touser":
			buff[i] = &t.Touser.FieldValue
		case "msgmode":
			buff[i] = &t.Msgmode.FieldValue
		case "small":
			buff[i] = &t.Small.FieldValue
		case "large":
			buff[i] = &t.Large.FieldValue
		}
	}
}

func NewTim_message(tableName ...string) *Tim_message {
	large := &tim_message_Large{fieldName: "large"}
	large.Field.FieldName = "large"
	id := &tim_message_Id{fieldName: "id"}
	id.Field.FieldName = "id"
	stamp := &tim_message_Stamp{fieldName: "stamp"}
	stamp.Field.FieldName = "stamp"
	touser := &tim_message_Touser{fieldName: "touser"}
	touser.Field.FieldName = "touser"
	msgmode := &tim_message_Msgmode{fieldName: "msgmode"}
	msgmode.Field.FieldName = "msgmode"
	small := &tim_message_Small{fieldName: "small"}
	small.Field.FieldName = "small"
	createtime := &tim_message_Createtime{fieldName: "createtime"}
	createtime.Field.FieldName = "createtime"
	chatid := &tim_message_Chatid{fieldName: "chatid"}
	chatid.Field.FieldName = "chatid"
	fromuser := &tim_message_Fromuser{fieldName: "fromuser"}
	fromuser.Field.FieldName = "fromuser"
	msgtype_ := &tim_message_Msgtype{fieldName: "msgtype"}
	msgtype_.Field.FieldName = "msgtype"
	gname := &tim_message_Gname{fieldName: "gname"}
	gname.Field.FieldName = "gname"
	stanza := &tim_message_Stanza{fieldName: "stanza"}
	stanza.Field.FieldName = "stanza"
	table := &Tim_message{Gname:gname,Stanza:stanza,Createtime:createtime,Chatid:chatid,Fromuser:fromuser,Msgtype:msgtype_,Msgmode:msgmode,Small:small,Large:large,Id:id,Stamp:stamp,Touser:touser}
	table.Table.ModifyMap = make(map[string]interface{})
	if len(tableName) == 1 {
		table.Table.TableName = tableName[0]
	} else {
		table.Table.TableName = "tim_message"
	}
	return table
}
