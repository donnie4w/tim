// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdao
//
// datetime :2024-08-10 01:01:22
// gdao version 1.1.0
// dbtype:sqlite ,database:timdb ,tablename:timblock

package dao

import (
	"fmt"
	"github.com/donnie4w/gdao"
	"github.com/donnie4w/gdao/base"
	
)

type timblock_Id[T any] struct {
	base.Field[T]
	fieldName  string
	fieldValue *int64
}

func (t *timblock_Id[T]) Name() string {
	return t.fieldName
}

func (t *timblock_Id[T]) Value() any {
	return t.fieldValue
}

type timblock_Unikid[T any] struct {
	base.Field[T]
	fieldName  string
	fieldValue *int64
}

func (t *timblock_Unikid[T]) Name() string {
	return t.fieldName
}

func (t *timblock_Unikid[T]) Value() any {
	return t.fieldValue
}

type timblock_Uuid[T any] struct {
	base.Field[T]
	fieldName  string
	fieldValue *int64
}

func (t *timblock_Uuid[T]) Name() string {
	return t.fieldName
}

func (t *timblock_Uuid[T]) Value() any {
	return t.fieldValue
}

type timblock_Tuuid[T any] struct {
	base.Field[T]
	fieldName  string
	fieldValue *int64
}

func (t *timblock_Tuuid[T]) Name() string {
	return t.fieldName
}

func (t *timblock_Tuuid[T]) Value() any {
	return t.fieldValue
}

type Timblock struct {
	gdao.Table[Timblock]

	Id		*timblock_Id[Timblock]
	Unikid		*timblock_Unikid[Timblock]
	Uuid		*timblock_Uuid[Timblock]
	Tuuid		*timblock_Tuuid[Timblock]
}

func (u *Timblock) GetId() (_r int64){
	if u.Id.fieldValue != nil {
		_r = *u.Id.fieldValue
	}
	return
}

func (u *Timblock) SetId(arg int64) *Timblock{
	u.Put0(u.Id.fieldName, arg)
	u.Id.fieldValue = &arg
	return u
}

func (u *Timblock) GetUnikid() (_r int64){
	if u.Unikid.fieldValue != nil {
		_r = *u.Unikid.fieldValue
	}
	return
}

func (u *Timblock) SetUnikid(arg int64) *Timblock{
	u.Put0(u.Unikid.fieldName, arg)
	u.Unikid.fieldValue = &arg
	return u
}

func (u *Timblock) GetUuid() (_r int64){
	if u.Uuid.fieldValue != nil {
		_r = *u.Uuid.fieldValue
	}
	return
}

func (u *Timblock) SetUuid(arg int64) *Timblock{
	u.Put0(u.Uuid.fieldName, arg)
	u.Uuid.fieldValue = &arg
	return u
}

func (u *Timblock) GetTuuid() (_r int64){
	if u.Tuuid.fieldValue != nil {
		_r = *u.Tuuid.fieldValue
	}
	return
}

func (u *Timblock) SetTuuid(arg int64) *Timblock{
	u.Put0(u.Tuuid.fieldName, arg)
	u.Tuuid.fieldValue = &arg
	return u
}


func (u *Timblock) Scan(fieldname string, value any) {
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

func (t *Timblock) ToGdao() {
	_t := NewTimblock()
	*t = *_t
}

func (t *Timblock) Copy(h *Timblock) *Timblock{
	t.SetId(h.GetId())
	t.SetUnikid(h.GetUnikid())
	t.SetUuid(h.GetUuid())
	t.SetTuuid(h.GetTuuid())
	return t
}

func (t *Timblock) String() string {
	return fmt.Sprint("Id:",t.GetId(), ",","Unikid:",t.GetUnikid(), ",","Uuid:",t.GetUuid(), ",","Tuuid:",t.GetTuuid())
}

func NewTimblock(tablename ...string) (_r *Timblock) {

	id := &timblock_Id[Timblock]{fieldName: "id"}
	id.Field.FieldName = "id"

	unikid := &timblock_Unikid[Timblock]{fieldName: "unikid"}
	unikid.Field.FieldName = "unikid"

	uuid := &timblock_Uuid[Timblock]{fieldName: "uuid"}
	uuid.Field.FieldName = "uuid"

	tuuid := &timblock_Tuuid[Timblock]{fieldName: "tuuid"}
	tuuid.Field.FieldName = "tuuid"

	_r = &Timblock{Id:id,Unikid:unikid,Uuid:uuid,Tuuid:tuuid}
	s := "timblock"
	if len(tablename) > 0 && tablename[0] != "" {
		s = tablename[0]
	}
	_r.Init(s, []base.Column[Timblock]{id,unikid,uuid,tuuid})
	return
}

func (t *Timblock) Encode() ([]byte, error) {
	m := make(map[string]any, 0)
	m["id"] = t.GetId()
	m["unikid"] = t.GetUnikid()
	m["uuid"] = t.GetUuid()
	m["tuuid"] = t.GetTuuid()
	return t.Table.Encode(m)
}

func (t *Timblock) Decode(bs []byte) (err error) {
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

