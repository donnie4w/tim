// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdao
//
// datetime :2024-08-10 01:01:22
// gdao version 1.1.0
// dbtype:sqlite ,database:timdb ,tablename:timroster

package dao

import (
	"fmt"
	"github.com/donnie4w/gdao"
	"github.com/donnie4w/gdao/base"
	
)

type timroster_Id[T any] struct {
	base.Field[T]
	fieldName  string
	fieldValue *int64
}

func (t *timroster_Id[T]) Name() string {
	return t.fieldName
}

func (t *timroster_Id[T]) Value() any {
	return t.fieldValue
}

type timroster_Relate[T any] struct {
	base.Field[T]
	fieldName  string
	fieldValue *int64
}

func (t *timroster_Relate[T]) Name() string {
	return t.fieldName
}

func (t *timroster_Relate[T]) Value() any {
	return t.fieldValue
}

type timroster_Uuid[T any] struct {
	base.Field[T]
	fieldName  string
	fieldValue *int64
}

func (t *timroster_Uuid[T]) Name() string {
	return t.fieldName
}

func (t *timroster_Uuid[T]) Value() any {
	return t.fieldValue
}

type timroster_Tuuid[T any] struct {
	base.Field[T]
	fieldName  string
	fieldValue *int64
}

func (t *timroster_Tuuid[T]) Name() string {
	return t.fieldName
}

func (t *timroster_Tuuid[T]) Value() any {
	return t.fieldValue
}

type Timroster struct {
	gdao.Table[Timroster]

	Id		*timroster_Id[Timroster]
	Relate		*timroster_Relate[Timroster]
	Uuid		*timroster_Uuid[Timroster]
	Tuuid		*timroster_Tuuid[Timroster]
}

func (u *Timroster) GetId() (_r int64){
	if u.Id.fieldValue != nil {
		_r = *u.Id.fieldValue
	}
	return
}

func (u *Timroster) SetId(arg int64) *Timroster{
	u.Put0(u.Id.fieldName, arg)
	u.Id.fieldValue = &arg
	return u
}

func (u *Timroster) GetRelate() (_r int64){
	if u.Relate.fieldValue != nil {
		_r = *u.Relate.fieldValue
	}
	return
}

func (u *Timroster) SetRelate(arg int64) *Timroster{
	u.Put0(u.Relate.fieldName, arg)
	u.Relate.fieldValue = &arg
	return u
}

func (u *Timroster) GetUuid() (_r int64){
	if u.Uuid.fieldValue != nil {
		_r = *u.Uuid.fieldValue
	}
	return
}

func (u *Timroster) SetUuid(arg int64) *Timroster{
	u.Put0(u.Uuid.fieldName, arg)
	u.Uuid.fieldValue = &arg
	return u
}

func (u *Timroster) GetTuuid() (_r int64){
	if u.Tuuid.fieldValue != nil {
		_r = *u.Tuuid.fieldValue
	}
	return
}

func (u *Timroster) SetTuuid(arg int64) *Timroster{
	u.Put0(u.Tuuid.fieldName, arg)
	u.Tuuid.fieldValue = &arg
	return u
}


func (u *Timroster) Scan(fieldname string, value any) {
	switch fieldname {
	case "id":
		u.SetId(base.AsInt64(value))
	case "relate":
		u.SetRelate(base.AsInt64(value))
	case "uuid":
		u.SetUuid(base.AsInt64(value))
	case "tuuid":
		u.SetTuuid(base.AsInt64(value))
	}
}

func (t *Timroster) ToGdao() {
	_t := NewTimroster()
	*t = *_t
}

func (t *Timroster) Copy(h *Timroster) *Timroster{
	t.SetId(h.GetId())
	t.SetRelate(h.GetRelate())
	t.SetUuid(h.GetUuid())
	t.SetTuuid(h.GetTuuid())
	return t
}

func (t *Timroster) String() string {
	return fmt.Sprint("Id:",t.GetId(), ",","Relate:",t.GetRelate(), ",","Uuid:",t.GetUuid(), ",","Tuuid:",t.GetTuuid())
}

func NewTimroster(tablename ...string) (_r *Timroster) {

	id := &timroster_Id[Timroster]{fieldName: "id"}
	id.Field.FieldName = "id"

	relate := &timroster_Relate[Timroster]{fieldName: "relate"}
	relate.Field.FieldName = "relate"

	uuid := &timroster_Uuid[Timroster]{fieldName: "uuid"}
	uuid.Field.FieldName = "uuid"

	tuuid := &timroster_Tuuid[Timroster]{fieldName: "tuuid"}
	tuuid.Field.FieldName = "tuuid"

	_r = &Timroster{Id:id,Relate:relate,Uuid:uuid,Tuuid:tuuid}
	s := "timroster"
	if len(tablename) > 0 && tablename[0] != "" {
		s = tablename[0]
	}
	_r.Init(s, []base.Column[Timroster]{id,relate,uuid,tuuid})
	return
}

func (t *Timroster) Encode() ([]byte, error) {
	m := make(map[string]any, 0)
	m["id"] = t.GetId()
	m["relate"] = t.GetRelate()
	m["uuid"] = t.GetUuid()
	m["tuuid"] = t.GetTuuid()
	return t.Table.Encode(m)
}

func (t *Timroster) Decode(bs []byte) (err error) {
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

