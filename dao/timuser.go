// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdao
//
// datetime :2024-08-10 01:01:22
// gdao version 1.1.0
// dbtype:sqlite ,database:timdb ,tablename:timuser

package dao

import (
	"fmt"
	"github.com/donnie4w/gdao"
	"github.com/donnie4w/gdao/base"
	
)

type timuser_Id[T any] struct {
	base.Field[T]
	fieldName  string
	fieldValue *int64
}

func (t *timuser_Id[T]) Name() string {
	return t.fieldName
}

func (t *timuser_Id[T]) Value() any {
	return t.fieldValue
}

type timuser_Uuid[T any] struct {
	base.Field[T]
	fieldName  string
	fieldValue *int64
}

func (t *timuser_Uuid[T]) Name() string {
	return t.fieldName
}

func (t *timuser_Uuid[T]) Value() any {
	return t.fieldValue
}

type timuser_Pwd[T any] struct {
	base.Field[T]
	fieldName  string
	fieldValue *int64
}

func (t *timuser_Pwd[T]) Name() string {
	return t.fieldName
}

func (t *timuser_Pwd[T]) Value() any {
	return t.fieldValue
}

type timuser_Createtime[T any] struct {
	base.Field[T]
	fieldName  string
	fieldValue *int64
}

func (t *timuser_Createtime[T]) Name() string {
	return t.fieldName
}

func (t *timuser_Createtime[T]) Value() any {
	return t.fieldValue
}

type timuser_Ubean[T any] struct {
	base.Field[T]
	fieldName  string
	fieldValue []byte
}

func (t *timuser_Ubean[T]) Name() string {
	return t.fieldName
}

func (t *timuser_Ubean[T]) Value() any {
	return t.fieldValue
}

type Timuser struct {
	gdao.Table[Timuser]

	Id		*timuser_Id[Timuser]
	Uuid		*timuser_Uuid[Timuser]
	Pwd		*timuser_Pwd[Timuser]
	Createtime		*timuser_Createtime[Timuser]
	Ubean		*timuser_Ubean[Timuser]
}

func (u *Timuser) GetId() (_r int64){
	if u.Id.fieldValue != nil {
		_r = *u.Id.fieldValue
	}
	return
}

func (u *Timuser) SetId(arg int64) *Timuser{
	u.Put0(u.Id.fieldName, arg)
	u.Id.fieldValue = &arg
	return u
}

func (u *Timuser) GetUuid() (_r int64){
	if u.Uuid.fieldValue != nil {
		_r = *u.Uuid.fieldValue
	}
	return
}

func (u *Timuser) SetUuid(arg int64) *Timuser{
	u.Put0(u.Uuid.fieldName, arg)
	u.Uuid.fieldValue = &arg
	return u
}

func (u *Timuser) GetPwd() (_r int64){
	if u.Pwd.fieldValue != nil {
		_r = *u.Pwd.fieldValue
	}
	return
}

func (u *Timuser) SetPwd(arg int64) *Timuser{
	u.Put0(u.Pwd.fieldName, arg)
	u.Pwd.fieldValue = &arg
	return u
}

func (u *Timuser) GetCreatetime() (_r int64){
	if u.Createtime.fieldValue != nil {
		_r = *u.Createtime.fieldValue
	}
	return
}

func (u *Timuser) SetCreatetime(arg int64) *Timuser{
	u.Put0(u.Createtime.fieldName, arg)
	u.Createtime.fieldValue = &arg
	return u
}

func (u *Timuser) GetUbean() (_r []byte){
	_r = u.Ubean.fieldValue
	return
}

func (u *Timuser) SetUbean(arg []byte) *Timuser{
	u.Put0(u.Ubean.fieldName, arg)
	u.Ubean.fieldValue = arg
	return u
}


func (u *Timuser) Scan(fieldname string, value any) {
	switch fieldname {
	case "id":
		u.SetId(base.AsInt64(value))
	case "uuid":
		u.SetUuid(base.AsInt64(value))
	case "pwd":
		u.SetPwd(base.AsInt64(value))
	case "createtime":
		u.SetCreatetime(base.AsInt64(value))
	case "ubean":
		u.SetUbean(base.AsBytes(value))
	}
}

func (t *Timuser) ToGdao() {
	_t := NewTimuser()
	*t = *_t
}

func (t *Timuser) Copy(h *Timuser) *Timuser{
	t.SetId(h.GetId())
	t.SetUuid(h.GetUuid())
	t.SetPwd(h.GetPwd())
	t.SetCreatetime(h.GetCreatetime())
	t.SetUbean(h.GetUbean())
	return t
}

func (t *Timuser) String() string {
	return fmt.Sprint("Id:",t.GetId(), ",","Uuid:",t.GetUuid(), ",","Pwd:",t.GetPwd(), ",","Createtime:",t.GetCreatetime(), ",","Ubean:",t.GetUbean())
}

func NewTimuser(tablename ...string) (_r *Timuser) {

	id := &timuser_Id[Timuser]{fieldName: "id"}
	id.Field.FieldName = "id"

	uuid := &timuser_Uuid[Timuser]{fieldName: "uuid"}
	uuid.Field.FieldName = "uuid"

	pwd := &timuser_Pwd[Timuser]{fieldName: "pwd"}
	pwd.Field.FieldName = "pwd"

	createtime := &timuser_Createtime[Timuser]{fieldName: "createtime"}
	createtime.Field.FieldName = "createtime"

	ubean := &timuser_Ubean[Timuser]{fieldName: "ubean"}
	ubean.Field.FieldName = "ubean"

	_r = &Timuser{Id:id,Uuid:uuid,Pwd:pwd,Createtime:createtime,Ubean:ubean}
	s := "timuser"
	if len(tablename) > 0 && tablename[0] != "" {
		s = tablename[0]
	}
	_r.Init(s, []base.Column[Timuser]{id,uuid,pwd,createtime,ubean})
	return
}

func (t *Timuser) Encode() ([]byte, error) {
	m := make(map[string]any, 0)
	m["id"] = t.GetId()
	m["uuid"] = t.GetUuid()
	m["pwd"] = t.GetPwd()
	m["createtime"] = t.GetCreatetime()
	m["ubean"] = t.GetUbean()
	return t.Table.Encode(m)
}

func (t *Timuser) Decode(bs []byte) (err error) {
	var m map[string]any
	if m, err = t.Table.Decode(bs); err == nil {
		if !t.IsInit() {
			t.ToGdao()
		}
		for name, bean := range m {
			t.Scan(name, bean)
		}
	}
	return
}

