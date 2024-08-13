// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdao
//
// datetime :2024-08-10 01:01:22
// gdao version 1.1.0
// dbtype:sqlite ,database:timdb ,tablename:timrelate

package dao

import (
	"fmt"
	"github.com/donnie4w/gdao"
	"github.com/donnie4w/gdao/base"
	
)

type timrelate_Id[T any] struct {
	base.Field[T]
	fieldName  string
	fieldValue *int64
}

func (t *timrelate_Id[T]) Name() string {
	return t.fieldName
}

func (t *timrelate_Id[T]) Value() any {
	return t.fieldValue
}

type timrelate_Uuid[T any] struct {
	base.Field[T]
	fieldName  string
	fieldValue *int64
}

func (t *timrelate_Uuid[T]) Name() string {
	return t.fieldName
}

func (t *timrelate_Uuid[T]) Value() any {
	return t.fieldValue
}

type timrelate_Status[T any] struct {
	base.Field[T]
	fieldName  string
	fieldValue *int64
}

func (t *timrelate_Status[T]) Name() string {
	return t.fieldName
}

func (t *timrelate_Status[T]) Value() any {
	return t.fieldValue
}

type Timrelate struct {
	gdao.Table[Timrelate]

	Id		*timrelate_Id[Timrelate]
	Uuid		*timrelate_Uuid[Timrelate]
	Status		*timrelate_Status[Timrelate]
}

func (u *Timrelate) GetId() (_r int64){
	if u.Id.fieldValue != nil {
		_r = *u.Id.fieldValue
	}
	return
}

func (u *Timrelate) SetId(arg int64) *Timrelate{
	u.Put0(u.Id.fieldName, arg)
	u.Id.fieldValue = &arg
	return u
}

func (u *Timrelate) GetUuid() (_r int64){
	if u.Uuid.fieldValue != nil {
		_r = *u.Uuid.fieldValue
	}
	return
}

func (u *Timrelate) SetUuid(arg int64) *Timrelate{
	u.Put0(u.Uuid.fieldName, arg)
	u.Uuid.fieldValue = &arg
	return u
}

func (u *Timrelate) GetStatus() (_r int64){
	if u.Status.fieldValue != nil {
		_r = *u.Status.fieldValue
	}
	return
}

func (u *Timrelate) SetStatus(arg int64) *Timrelate{
	u.Put0(u.Status.fieldName, arg)
	u.Status.fieldValue = &arg
	return u
}


func (u *Timrelate) Scan(fieldname string, value any) {
	switch fieldname {
	case "id":
		u.SetId(base.AsInt64(value))
	case "uuid":
		u.SetUuid(base.AsInt64(value))
	case "status":
		u.SetStatus(base.AsInt64(value))
	}
}

func (t *Timrelate) ToGdao() {
	_t := NewTimrelate()
	*t = *_t
}

func (t *Timrelate) Copy(h *Timrelate) *Timrelate{
	t.SetId(h.GetId())
	t.SetUuid(h.GetUuid())
	t.SetStatus(h.GetStatus())
	return t
}

func (t *Timrelate) String() string {
	return fmt.Sprint("Id:",t.GetId(), ",","Uuid:",t.GetUuid(), ",","Status:",t.GetStatus())
}

func NewTimrelate(tablename ...string) (_r *Timrelate) {

	id := &timrelate_Id[Timrelate]{fieldName: "id"}
	id.Field.FieldName = "id"

	uuid := &timrelate_Uuid[Timrelate]{fieldName: "uuid"}
	uuid.Field.FieldName = "uuid"

	status := &timrelate_Status[Timrelate]{fieldName: "status"}
	status.Field.FieldName = "status"

	_r = &Timrelate{Id:id,Uuid:uuid,Status:status}
	s := "timrelate"
	if len(tablename) > 0 && tablename[0] != "" {
		s = tablename[0]
	}
	_r.Init(s, []base.Column[Timrelate]{id,uuid,status})
	return
}

func (t *Timrelate) Encode() ([]byte, error) {
	m := make(map[string]any, 0)
	m["id"] = t.GetId()
	m["uuid"] = t.GetUuid()
	m["status"] = t.GetStatus()
	return t.Table.Encode(m)
}

func (t *Timrelate) Decode(bs []byte) (err error) {
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

