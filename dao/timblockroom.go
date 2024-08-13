// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdao
//
// datetime :2024-08-10 01:01:22
// gdao version 1.1.0
// dbtype:sqlite ,database:timdb ,tablename:timblockroom

package dao

import (
	"fmt"
	"github.com/donnie4w/gdao"
	"github.com/donnie4w/gdao/base"
	
)

type timblockroom_Id[T any] struct {
	base.Field[T]
	fieldName  string
	fieldValue *int64
}

func (t *timblockroom_Id[T]) Name() string {
	return t.fieldName
}

func (t *timblockroom_Id[T]) Value() any {
	return t.fieldValue
}

type timblockroom_Unikid[T any] struct {
	base.Field[T]
	fieldName  string
	fieldValue *int64
}

func (t *timblockroom_Unikid[T]) Name() string {
	return t.fieldName
}

func (t *timblockroom_Unikid[T]) Value() any {
	return t.fieldValue
}

type timblockroom_Uuid[T any] struct {
	base.Field[T]
	fieldName  string
	fieldValue *int64
}

func (t *timblockroom_Uuid[T]) Name() string {
	return t.fieldName
}

func (t *timblockroom_Uuid[T]) Value() any {
	return t.fieldValue
}

type timblockroom_Tuuid[T any] struct {
	base.Field[T]
	fieldName  string
	fieldValue *int64
}

func (t *timblockroom_Tuuid[T]) Name() string {
	return t.fieldName
}

func (t *timblockroom_Tuuid[T]) Value() any {
	return t.fieldValue
}

type Timblockroom struct {
	gdao.Table[Timblockroom]

	Id		*timblockroom_Id[Timblockroom]
	Unikid		*timblockroom_Unikid[Timblockroom]
	Uuid		*timblockroom_Uuid[Timblockroom]
	Tuuid		*timblockroom_Tuuid[Timblockroom]
}

func (u *Timblockroom) GetId() (_r int64){
	if u.Id.fieldValue != nil {
		_r = *u.Id.fieldValue
	}
	return
}

func (u *Timblockroom) SetId(arg int64) *Timblockroom{
	u.Put0(u.Id.fieldName, arg)
	u.Id.fieldValue = &arg
	return u
}

func (u *Timblockroom) GetUnikid() (_r int64){
	if u.Unikid.fieldValue != nil {
		_r = *u.Unikid.fieldValue
	}
	return
}

func (u *Timblockroom) SetUnikid(arg int64) *Timblockroom{
	u.Put0(u.Unikid.fieldName, arg)
	u.Unikid.fieldValue = &arg
	return u
}

func (u *Timblockroom) GetUuid() (_r int64){
	if u.Uuid.fieldValue != nil {
		_r = *u.Uuid.fieldValue
	}
	return
}

func (u *Timblockroom) SetUuid(arg int64) *Timblockroom{
	u.Put0(u.Uuid.fieldName, arg)
	u.Uuid.fieldValue = &arg
	return u
}

func (u *Timblockroom) GetTuuid() (_r int64){
	if u.Tuuid.fieldValue != nil {
		_r = *u.Tuuid.fieldValue
	}
	return
}

func (u *Timblockroom) SetTuuid(arg int64) *Timblockroom{
	u.Put0(u.Tuuid.fieldName, arg)
	u.Tuuid.fieldValue = &arg
	return u
}


func (u *Timblockroom) Scan(fieldname string, value any) {
	switch fieldname {
	case "id":
		u.SetId(base.AsInt64(value))
	case "unikid":
		u.SetUnikid(base.AsInt64(value))
	case "uuid":
		u.SetUuid(base.AsInt64(value))
	case "tuuid":
		u.SetTuuid(base.AsInt64(value))
	}
}

func (t *Timblockroom) ToGdao() {
	_t := NewTimblockroom()
	*t = *_t
}

func (t *Timblockroom) Copy(h *Timblockroom) *Timblockroom{
	t.SetId(h.GetId())
	t.SetUnikid(h.GetUnikid())
	t.SetUuid(h.GetUuid())
	t.SetTuuid(h.GetTuuid())
	return t
}

func (t *Timblockroom) String() string {
	return fmt.Sprint("Id:",t.GetId(), ",","Unikid:",t.GetUnikid(), ",","Uuid:",t.GetUuid(), ",","Tuuid:",t.GetTuuid())
}

func NewTimblockroom(tablename ...string) (_r *Timblockroom) {

	id := &timblockroom_Id[Timblockroom]{fieldName: "id"}
	id.Field.FieldName = "id"

	unikid := &timblockroom_Unikid[Timblockroom]{fieldName: "unikid"}
	unikid.Field.FieldName = "unikid"

	uuid := &timblockroom_Uuid[Timblockroom]{fieldName: "uuid"}
	uuid.Field.FieldName = "uuid"

	tuuid := &timblockroom_Tuuid[Timblockroom]{fieldName: "tuuid"}
	tuuid.Field.FieldName = "tuuid"

	_r = &Timblockroom{Id:id,Unikid:unikid,Uuid:uuid,Tuuid:tuuid}
	s := "timblockroom"
	if len(tablename) > 0 && tablename[0] != "" {
		s = tablename[0]
	}
	_r.Init(s, []base.Column[Timblockroom]{id,unikid,uuid,tuuid})
	return
}

func (t *Timblockroom) Encode() ([]byte, error) {
	m := make(map[string]any, 0)
	m["id"] = t.GetId()
	m["unikid"] = t.GetUnikid()
	m["uuid"] = t.GetUuid()
	m["tuuid"] = t.GetTuuid()
	return t.Table.Encode(m)
}

func (t *Timblockroom) Decode(bs []byte) (err error) {
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

