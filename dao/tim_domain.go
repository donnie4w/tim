package dao

/**
tablename:tim_domain
datetime :2016-09-07 11:32:22
*/
import (
	"github.com/donnie4w/gdao"
	"reflect"
)

type tim_domain_Remark struct {
	gdao.Field
	fieldName  string
	FieldValue *string
}

func (c *tim_domain_Remark) Name() string {
	return c.fieldName
}

func (c *tim_domain_Remark) Value() interface{} {
	return c.FieldValue
}

type tim_domain_Id struct {
	gdao.Field
	fieldName  string
	FieldValue *int32
}

func (c *tim_domain_Id) Name() string {
	return c.fieldName
}

func (c *tim_domain_Id) Value() interface{} {
	return c.FieldValue
}

type tim_domain_Domain struct {
	gdao.Field
	fieldName  string
	FieldValue *string
}

func (c *tim_domain_Domain) Name() string {
	return c.fieldName
}

func (c *tim_domain_Domain) Value() interface{} {
	return c.FieldValue
}

type tim_domain_Createtime struct {
	gdao.Field
	fieldName  string
	FieldValue *string
}

func (c *tim_domain_Createtime) Name() string {
	return c.fieldName
}

func (c *tim_domain_Createtime) Value() interface{} {
	return c.FieldValue
}

type Tim_domain struct {
	gdao.Table
	Domain *tim_domain_Domain
	Createtime *tim_domain_Createtime
	Remark *tim_domain_Remark
	Id *tim_domain_Id
}

func (u *Tim_domain) GetCreatetime() string {
	return *u.Createtime.FieldValue
}

func (u *Tim_domain) SetCreatetime(arg string) {
	u.Table.ModifyMap[u.Createtime.fieldName] = arg
	v := string(arg)
	u.Createtime.FieldValue = &v
}

func (u *Tim_domain) GetRemark() string {
	return *u.Remark.FieldValue
}

func (u *Tim_domain) SetRemark(arg string) {
	u.Table.ModifyMap[u.Remark.fieldName] = arg
	v := string(arg)
	u.Remark.FieldValue = &v
}

func (u *Tim_domain) GetId() int32 {
	return *u.Id.FieldValue
}

func (u *Tim_domain) SetId(arg int64) {
	u.Table.ModifyMap[u.Id.fieldName] = arg
	v := int32(arg)
	u.Id.FieldValue = &v
}

func (u *Tim_domain) GetDomain() string {
	return *u.Domain.FieldValue
}

func (u *Tim_domain) SetDomain(arg string) {
	u.Table.ModifyMap[u.Domain.fieldName] = arg
	v := string(arg)
	u.Domain.FieldValue = &v
}

func (t *Tim_domain) Query(columns ...gdao.Column) ([]Tim_domain,error) {
	if columns == nil {
		columns = []gdao.Column{ t.Remark,t.Id,t.Domain,t.Createtime}
	}
	rs,err := t.Table.Query(columns...)
	if rs == nil || err != nil {
		return nil, err
	}
	ts := make([]Tim_domain, 0, len(rs))
	c := make(chan int16,len(rs))
	for _, rows := range rs {
		t := NewTim_domain()
		go copyTim_domain(c, rows, t, columns)
		<-c
		ts = append(ts, *t)
	}
	return ts,nil
}

func copyTim_domain(channle chan int16, rows []interface{}, t *Tim_domain, columns []gdao.Column) {
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

func (t *Tim_domain) QuerySingle(columns ...gdao.Column) (*Tim_domain,error) {
	if columns == nil {
		columns = []gdao.Column{ t.Remark,t.Id,t.Domain,t.Createtime}
	}
	rs,err := t.Table.QuerySingle(columns...)
	if rs == nil || err != nil {
		return nil, err
	}
	rt := NewTim_domain()
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

func (t *Tim_domain) Select(columns ...gdao.Column) (*Tim_domain,error) {
	if columns == nil {
		columns = []gdao.Column{ t.Remark,t.Id,t.Domain,t.Createtime}
	}
	rows,err := t.Table.Selects(columns...)
	defer rows.Close()
	if err != nil || rows==nil {
		return nil, err
	}
	buff := make([]interface{}, len(columns))
	if rows.Next() {
		n := NewTim_domain()
		cpTim_domain(buff, n, columns)
		row_err := rows.Scan(buff...)
		if row_err != nil {
			return nil, row_err
		}
		return n, nil
	}
	return nil, nil
}

func (t *Tim_domain) Selects(columns ...gdao.Column) ([]*Tim_domain,error) {
	if columns == nil {
		columns = []gdao.Column{ t.Remark,t.Id,t.Domain,t.Createtime}
	}
	rows,err := t.Table.Selects(columns...)
	defer rows.Close()
	if err != nil || rows==nil {
		return nil, err
	}
	ns := make([]*Tim_domain, 0)
	buff := make([]interface{}, len(columns))
	for rows.Next() {
		n := NewTim_domain()
		cpTim_domain(buff, n, columns)
		row_err := rows.Scan(buff...)
		if row_err != nil {
			return nil, row_err
		}
		ns = append(ns, n)
	}
	return ns, nil
}

func  cpTim_domain(buff []interface{}, t *Tim_domain, columns []gdao.Column) {
	for i, column := range columns {
		field := column.Name()
		switch field {
		case "id":
			buff[i] = &t.Id.FieldValue
		case "domain":
			buff[i] = &t.Domain.FieldValue
		case "createtime":
			buff[i] = &t.Createtime.FieldValue
		case "remark":
			buff[i] = &t.Remark.FieldValue
		}
	}
}

func NewTim_domain(tableName ...string) *Tim_domain {
	id := &tim_domain_Id{fieldName: "id"}
	id.Field.FieldName = "id"
	domain := &tim_domain_Domain{fieldName: "domain"}
	domain.Field.FieldName = "domain"
	createtime := &tim_domain_Createtime{fieldName: "createtime"}
	createtime.Field.FieldName = "createtime"
	remark := &tim_domain_Remark{fieldName: "remark"}
	remark.Field.FieldName = "remark"
	table := &Tim_domain{Createtime:createtime,Remark:remark,Id:id,Domain:domain}
	table.Table.ModifyMap = make(map[string]interface{})
	if len(tableName) == 1 {
		table.Table.TableName = tableName[0]
	} else {
		table.Table.TableName = "tim_domain"
	}
	return table
}
