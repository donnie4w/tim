package dao

/**
tablename:tim_user
datetime :2016-09-07 11:32:22
*/
import (
	"github.com/donnie4w/gdao"
	"reflect"
)

type tim_user_Loginname struct {
	gdao.Field
	fieldName  string
	FieldValue *string
}

func (c *tim_user_Loginname) Name() string {
	return c.fieldName
}

func (c *tim_user_Loginname) Value() interface{} {
	return c.FieldValue
}

type tim_user_Username struct {
	gdao.Field
	fieldName  string
	FieldValue *string
}

func (c *tim_user_Username) Name() string {
	return c.fieldName
}

func (c *tim_user_Username) Value() interface{} {
	return c.FieldValue
}

type tim_user_Usernick struct {
	gdao.Field
	fieldName  string
	FieldValue *string
}

func (c *tim_user_Usernick) Name() string {
	return c.fieldName
}

func (c *tim_user_Usernick) Value() interface{} {
	return c.FieldValue
}

type tim_user_Plainpassword struct {
	gdao.Field
	fieldName  string
	FieldValue *string
}

func (c *tim_user_Plainpassword) Name() string {
	return c.fieldName
}

func (c *tim_user_Plainpassword) Value() interface{} {
	return c.FieldValue
}

type tim_user_Encryptedpassword struct {
	gdao.Field
	fieldName  string
	FieldValue *string
}

func (c *tim_user_Encryptedpassword) Name() string {
	return c.fieldName
}

func (c *tim_user_Encryptedpassword) Value() interface{} {
	return c.FieldValue
}

type tim_user_Createtime struct {
	gdao.Field
	fieldName  string
	FieldValue *string
}

func (c *tim_user_Createtime) Name() string {
	return c.fieldName
}

func (c *tim_user_Createtime) Value() interface{} {
	return c.FieldValue
}

type tim_user_Updatetime struct {
	gdao.Field
	fieldName  string
	FieldValue *string
}

func (c *tim_user_Updatetime) Name() string {
	return c.fieldName
}

func (c *tim_user_Updatetime) Value() interface{} {
	return c.FieldValue
}

type tim_user_Id struct {
	gdao.Field
	fieldName  string
	FieldValue *int32
}

func (c *tim_user_Id) Name() string {
	return c.fieldName
}

func (c *tim_user_Id) Value() interface{} {
	return c.FieldValue
}

type Tim_user struct {
	gdao.Table
	Username *tim_user_Username
	Usernick *tim_user_Usernick
	Plainpassword *tim_user_Plainpassword
	Encryptedpassword *tim_user_Encryptedpassword
	Createtime *tim_user_Createtime
	Updatetime *tim_user_Updatetime
	Id *tim_user_Id
	Loginname *tim_user_Loginname
}

func (u *Tim_user) GetPlainpassword() string {
	return *u.Plainpassword.FieldValue
}

func (u *Tim_user) SetPlainpassword(arg string) {
	u.Table.ModifyMap[u.Plainpassword.fieldName] = arg
	v := string(arg)
	u.Plainpassword.FieldValue = &v
}

func (u *Tim_user) GetEncryptedpassword() string {
	return *u.Encryptedpassword.FieldValue
}

func (u *Tim_user) SetEncryptedpassword(arg string) {
	u.Table.ModifyMap[u.Encryptedpassword.fieldName] = arg
	v := string(arg)
	u.Encryptedpassword.FieldValue = &v
}

func (u *Tim_user) GetCreatetime() string {
	return *u.Createtime.FieldValue
}

func (u *Tim_user) SetCreatetime(arg string) {
	u.Table.ModifyMap[u.Createtime.fieldName] = arg
	v := string(arg)
	u.Createtime.FieldValue = &v
}

func (u *Tim_user) GetUpdatetime() string {
	return *u.Updatetime.FieldValue
}

func (u *Tim_user) SetUpdatetime(arg string) {
	u.Table.ModifyMap[u.Updatetime.fieldName] = arg
	v := string(arg)
	u.Updatetime.FieldValue = &v
}

func (u *Tim_user) GetId() int32 {
	return *u.Id.FieldValue
}

func (u *Tim_user) SetId(arg int64) {
	u.Table.ModifyMap[u.Id.fieldName] = arg
	v := int32(arg)
	u.Id.FieldValue = &v
}

func (u *Tim_user) GetLoginname() string {
	return *u.Loginname.FieldValue
}

func (u *Tim_user) SetLoginname(arg string) {
	u.Table.ModifyMap[u.Loginname.fieldName] = arg
	v := string(arg)
	u.Loginname.FieldValue = &v
}

func (u *Tim_user) GetUsername() string {
	return *u.Username.FieldValue
}

func (u *Tim_user) SetUsername(arg string) {
	u.Table.ModifyMap[u.Username.fieldName] = arg
	v := string(arg)
	u.Username.FieldValue = &v
}

func (u *Tim_user) GetUsernick() string {
	return *u.Usernick.FieldValue
}

func (u *Tim_user) SetUsernick(arg string) {
	u.Table.ModifyMap[u.Usernick.fieldName] = arg
	v := string(arg)
	u.Usernick.FieldValue = &v
}

func (t *Tim_user) Query(columns ...gdao.Column) ([]Tim_user,error) {
	if columns == nil {
		columns = []gdao.Column{ t.Plainpassword,t.Encryptedpassword,t.Createtime,t.Updatetime,t.Id,t.Loginname,t.Username,t.Usernick}
	}
	rs,err := t.Table.Query(columns...)
	if rs == nil || err != nil {
		return nil, err
	}
	ts := make([]Tim_user, 0, len(rs))
	c := make(chan int16,len(rs))
	for _, rows := range rs {
		t := NewTim_user()
		go copyTim_user(c, rows, t, columns)
		<-c
		ts = append(ts, *t)
	}
	return ts,nil
}

func copyTim_user(channle chan int16, rows []interface{}, t *Tim_user, columns []gdao.Column) {
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

func (t *Tim_user) QuerySingle(columns ...gdao.Column) (*Tim_user,error) {
	if columns == nil {
		columns = []gdao.Column{ t.Plainpassword,t.Encryptedpassword,t.Createtime,t.Updatetime,t.Id,t.Loginname,t.Username,t.Usernick}
	}
	rs,err := t.Table.QuerySingle(columns...)
	if rs == nil || err != nil {
		return nil, err
	}
	rt := NewTim_user()
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

func (t *Tim_user) Select(columns ...gdao.Column) (*Tim_user,error) {
	if columns == nil {
		columns = []gdao.Column{ t.Plainpassword,t.Encryptedpassword,t.Createtime,t.Updatetime,t.Id,t.Loginname,t.Username,t.Usernick}
	}
	rows,err := t.Table.Selects(columns...)
	defer rows.Close()
	if err != nil || rows==nil {
		return nil, err
	}
	buff := make([]interface{}, len(columns))
	if rows.Next() {
		n := NewTim_user()
		cpTim_user(buff, n, columns)
		row_err := rows.Scan(buff...)
		if row_err != nil {
			return nil, row_err
		}
		return n, nil
	}
	return nil, nil
}

func (t *Tim_user) Selects(columns ...gdao.Column) ([]*Tim_user,error) {
	if columns == nil {
		columns = []gdao.Column{ t.Plainpassword,t.Encryptedpassword,t.Createtime,t.Updatetime,t.Id,t.Loginname,t.Username,t.Usernick}
	}
	rows,err := t.Table.Selects(columns...)
	defer rows.Close()
	if err != nil || rows==nil {
		return nil, err
	}
	ns := make([]*Tim_user, 0)
	buff := make([]interface{}, len(columns))
	for rows.Next() {
		n := NewTim_user()
		cpTim_user(buff, n, columns)
		row_err := rows.Scan(buff...)
		if row_err != nil {
			return nil, row_err
		}
		ns = append(ns, n)
	}
	return ns, nil
}

func  cpTim_user(buff []interface{}, t *Tim_user, columns []gdao.Column) {
	for i, column := range columns {
		field := column.Name()
		switch field {
		case "createtime":
			buff[i] = &t.Createtime.FieldValue
		case "updatetime":
			buff[i] = &t.Updatetime.FieldValue
		case "id":
			buff[i] = &t.Id.FieldValue
		case "loginname":
			buff[i] = &t.Loginname.FieldValue
		case "username":
			buff[i] = &t.Username.FieldValue
		case "usernick":
			buff[i] = &t.Usernick.FieldValue
		case "plainpassword":
			buff[i] = &t.Plainpassword.FieldValue
		case "encryptedpassword":
			buff[i] = &t.Encryptedpassword.FieldValue
		}
	}
}

func NewTim_user(tableName ...string) *Tim_user {
	usernick := &tim_user_Usernick{fieldName: "usernick"}
	usernick.Field.FieldName = "usernick"
	plainpassword := &tim_user_Plainpassword{fieldName: "plainpassword"}
	plainpassword.Field.FieldName = "plainpassword"
	encryptedpassword := &tim_user_Encryptedpassword{fieldName: "encryptedpassword"}
	encryptedpassword.Field.FieldName = "encryptedpassword"
	createtime := &tim_user_Createtime{fieldName: "createtime"}
	createtime.Field.FieldName = "createtime"
	updatetime := &tim_user_Updatetime{fieldName: "updatetime"}
	updatetime.Field.FieldName = "updatetime"
	id := &tim_user_Id{fieldName: "id"}
	id.Field.FieldName = "id"
	loginname := &tim_user_Loginname{fieldName: "loginname"}
	loginname.Field.FieldName = "loginname"
	username := &tim_user_Username{fieldName: "username"}
	username.Field.FieldName = "username"
	table := &Tim_user{Createtime:createtime,Updatetime:updatetime,Id:id,Loginname:loginname,Username:username,Usernick:usernick,Plainpassword:plainpassword,Encryptedpassword:encryptedpassword}
	table.Table.ModifyMap = make(map[string]interface{})
	if len(tableName) == 1 {
		table.Table.TableName = tableName[0]
	} else {
		table.Table.TableName = "tim_user"
	}
	return table
}
