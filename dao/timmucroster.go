// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdao
//
// datetime :2024-08-10 01:01:22
// gdao version 1.1.0
// dbtype:sqlite ,database:timdb ,tablename:timmucroster

package dao

import (
	"fmt"
	"github.com/donnie4w/gdao"
	"github.com/donnie4w/gdao/base"
	
)

type timmucroster_Id[T any] struct {
	base.Field[T]
	fieldName  string
	fieldValue *int64
}

func (t *timmucroster_Id[T]) Name() string {
	return t.fieldName
}

func (t *timmucroster_Id[T]) Value() any {
	return t.fieldValue
}

type timmucroster_Relate[T any] struct {
	base.Field[T]
	fieldName  string
	fieldValue *int64
}

func (t *timmucroster_Relate[T]) Name() string {
	return t.fieldName
}

func (t *timmucroster_Relate[T]) Value() any {
	return t.fieldValue
}

type timmucroster_Uuid[T any] struct {
	base.Field[T]
	fieldName  string
	fieldValue *int64
}

func (t *timmucroster_Uuid[T]) Name() string {
	return t.fieldName
}

func (t *timmucroster_Uuid[T]) Value() any {
	return t.fieldValue
}

type timmucroster_Tuuid[T any] struct {
	base.Field[T]
	fieldName  string
	fieldValue *int64
}

func (t *timmucroster_Tuuid[T]) Name() string {
	return t.fieldName
}

func (t *timmucroster_Tuuid[T]) Value() any {
	return t.fieldValue
}

type Timmucroster struct {
	gdao.Table[Timmucroster]

	Id		*timmucroster_Id[Timmucroster]
	Relate		*timmucroster_Relate[Timmucroster]
	Uuid		*timmucroster_Uuid[Timmucroster]
	Tuuid		*timmucroster_Tuuid[Timmucroster]
}

func (u *Timmucroster) GetId() (_r int64){
	if u.Id.fieldValue != nil {
		_r = *u.Id.fieldValue
	}
	return
}

func (u *Timmucroster) SetId(arg int64) *Timmucroster{
	u.Put0(u.Id.fieldName, arg)
	u.Id.fieldValue = &arg
	return u
}

func (u *Timmucroster) GetRelate() (_r int64){
	if u.Relate.fieldValue != nil {
		_r = *u.Relate.fieldValue
	}
	return
}

func (u *Timmucroster) SetRelate(arg int64) *Timmucroster{
	u.Put0(u.Relate.fieldName, arg)
	u.Relate.fieldValue = &arg
	return u
}

func (u *Timmucroster) GetUuid() (_r int64){
	if u.Uuid.fieldValue != nil {
		_r = *u.Uuid.fieldValue
	}
	return
}

func (u *Timmucroster) SetUuid(arg int64) *Timmucroster{
	u.Put0(u.Uuid.fieldName, arg)
	u.Uuid.fieldValue = &arg
	return u
}

func (u *Timmucroster) GetTuuid() (_r int64){
	if u.Tuuid.fieldValue != nil {
		_r = *u.Tuuid.fieldValue
	}
	return
}

func (u *Timmucroster) SetTuuid(arg int64) *Timmucroster{
	u.Put0(u.Tuuid.fieldName, arg)
	u.Tuuid.fieldValue = &arg
	return u
}


func (u *Timmucroster) Scan(fieldname string, value any) {
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

func (t *Timmucroster) ToGdao() {
	_t := NewTimmucroster()
	*t = *_t
}

func (t *Timmucroster) Copy(h *Timmucroster) *Timmucroster{
	t.SetId(h.GetId())
	t.SetRelate(h.GetRelate())
	t.SetUuid(h.GetUuid())
	t.SetTuuid(h.GetTuuid())
	return t
}

func (t *Timmucroster) String() string {
	return fmt.Sprint("Id:",t.GetId(), ",","Relate:",t.GetRelate(), ",","Uuid:",t.GetUuid(), ",","Tuuid:",t.GetTuuid())
}

func NewTimmucroster(tablename ...string) (_r *Timmucroster) {

	id := &timmucroster_Id[Timmucroster]{fieldName: "id"}
	id.Field.FieldName = "id"

	relate := &timmucroster_Relate[Timmucroster]{fieldName: "relate"}
	relate.Field.FieldName = "relate"

	uuid := &timmucroster_Uuid[Timmucroster]{fieldName: "uuid"}
	uuid.Field.FieldName = "uuid"

	tuuid := &timmucroster_Tuuid[Timmucroster]{fieldName: "tuuid"}
	tuuid.Field.FieldName = "tuuid"

	_r = &Timmucroster{Id:id,Relate:relate,Uuid:uuid,Tuuid:tuuid}
	s := "timmucroster"
	if len(tablename) > 0 && tablename[0] != "" {
		s = tablename[0]
	}
	_r.Init(s, []base.Column[Timmucroster]{id,relate,uuid,tuuid})
	return
}

func (t *Timmucroster) Encode() ([]byte, error) {
	m := make(map[string]any, 0)
	m["id"] = t.GetId()
	m["relate"] = t.GetRelate()
	m["uuid"] = t.GetUuid()
	m["tuuid"] = t.GetTuuid()
	return t.Table.Encode(m)
}

func (t *Timmucroster) Decode(bs []byte) (err error) {
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

