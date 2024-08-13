// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdao
//
// datetime :2024-08-10 01:01:22
// gdao version 1.1.0
// dbtype:sqlite ,database:timdb ,tablename:timgroup

package dao

import (
	"fmt"
	"github.com/donnie4w/gdao"
	"github.com/donnie4w/gdao/base"
	
)

type timgroup_Id[T any] struct {
	base.Field[T]
	fieldName  string
	fieldValue *int64
}

func (t *timgroup_Id[T]) Name() string {
	return t.fieldName
}

func (t *timgroup_Id[T]) Value() any {
	return t.fieldValue
}

type timgroup_Gtype[T any] struct {
	base.Field[T]
	fieldName  string
	fieldValue *int64
}

func (t *timgroup_Gtype[T]) Name() string {
	return t.fieldName
}

func (t *timgroup_Gtype[T]) Value() any {
	return t.fieldValue
}

type timgroup_Uuid[T any] struct {
	base.Field[T]
	fieldName  string
	fieldValue *int64
}

func (t *timgroup_Uuid[T]) Name() string {
	return t.fieldName
}

func (t *timgroup_Uuid[T]) Value() any {
	return t.fieldValue
}

type timgroup_Createtime[T any] struct {
	base.Field[T]
	fieldName  string
	fieldValue *int64
}

func (t *timgroup_Createtime[T]) Name() string {
	return t.fieldName
}

func (t *timgroup_Createtime[T]) Value() any {
	return t.fieldValue
}

type timgroup_Status[T any] struct {
	base.Field[T]
	fieldName  string
	fieldValue *int64
}

func (t *timgroup_Status[T]) Name() string {
	return t.fieldName
}

func (t *timgroup_Status[T]) Value() any {
	return t.fieldValue
}

type timgroup_Rbean[T any] struct {
	base.Field[T]
	fieldName  string
	fieldValue []byte
}

func (t *timgroup_Rbean[T]) Name() string {
	return t.fieldName
}

func (t *timgroup_Rbean[T]) Value() any {
	return t.fieldValue
}

type Timgroup struct {
	gdao.Table[Timgroup]

	Id		*timgroup_Id[Timgroup]
	Gtype		*timgroup_Gtype[Timgroup]
	Uuid		*timgroup_Uuid[Timgroup]
	Createtime		*timgroup_Createtime[Timgroup]
	Status		*timgroup_Status[Timgroup]
	Rbean		*timgroup_Rbean[Timgroup]
}

func (u *Timgroup) GetId() (_r int64){
	if u.Id.fieldValue != nil {
		_r = *u.Id.fieldValue
	}
	return
}

func (u *Timgroup) SetId(arg int64) *Timgroup{
	u.Put0(u.Id.fieldName, arg)
	u.Id.fieldValue = &arg
	return u
}

func (u *Timgroup) GetGtype() (_r int64){
	if u.Gtype.fieldValue != nil {
		_r = *u.Gtype.fieldValue
	}
	return
}

func (u *Timgroup) SetGtype(arg int64) *Timgroup{
	u.Put0(u.Gtype.fieldName, arg)
	u.Gtype.fieldValue = &arg
	return u
}

func (u *Timgroup) GetUuid() (_r int64){
	if u.Uuid.fieldValue != nil {
		_r = *u.Uuid.fieldValue
	}
	return
}

func (u *Timgroup) SetUuid(arg int64) *Timgroup{
	u.Put0(u.Uuid.fieldName, arg)
	u.Uuid.fieldValue = &arg
	return u
}

func (u *Timgroup) GetCreatetime() (_r int64){
	if u.Createtime.fieldValue != nil {
		_r = *u.Createtime.fieldValue
	}
	return
}

func (u *Timgroup) SetCreatetime(arg int64) *Timgroup{
	u.Put0(u.Createtime.fieldName, arg)
	u.Createtime.fieldValue = &arg
	return u
}

func (u *Timgroup) GetStatus() (_r int64){
	if u.Status.fieldValue != nil {
		_r = *u.Status.fieldValue
	}
	return
}

func (u *Timgroup) SetStatus(arg int64) *Timgroup{
	u.Put0(u.Status.fieldName, arg)
	u.Status.fieldValue = &arg
	return u
}

func (u *Timgroup) GetRbean() (_r []byte){
	_r = u.Rbean.fieldValue
	return
}

func (u *Timgroup) SetRbean(arg []byte) *Timgroup{
	u.Put0(u.Rbean.fieldName, arg)
	u.Rbean.fieldValue = arg
	return u
}


func (u *Timgroup) Scan(fieldname string, value any) {
	switch fieldname {
	case "id":
		u.SetId(base.AsInt64(value))
	case "gtype":
		u.SetGtype(base.AsInt64(value))
	case "uuid":
		u.SetUuid(base.AsInt64(value))
	case "createtime":
		u.SetCreatetime(base.AsInt64(value))
	case "status":
		u.SetStatus(base.AsInt64(value))
	case "rbean":
		u.SetRbean(base.AsBytes(value))
	}
}

func (t *Timgroup) ToGdao() {
	_t := NewTimgroup()
	*t = *_t
}

func (t *Timgroup) Copy(h *Timgroup) *Timgroup{
	t.SetId(h.GetId())
	t.SetGtype(h.GetGtype())
	t.SetUuid(h.GetUuid())
	t.SetCreatetime(h.GetCreatetime())
	t.SetStatus(h.GetStatus())
	t.SetRbean(h.GetRbean())
	return t
}

func (t *Timgroup) String() string {
	return fmt.Sprint("Id:",t.GetId(), ",","Gtype:",t.GetGtype(), ",","Uuid:",t.GetUuid(), ",","Createtime:",t.GetCreatetime(), ",","Status:",t.GetStatus(), ",","Rbean:",t.GetRbean())
}

func NewTimgroup(tablename ...string) (_r *Timgroup) {

	id := &timgroup_Id[Timgroup]{fieldName: "id"}
	id.Field.FieldName = "id"

	gtype := &timgroup_Gtype[Timgroup]{fieldName: "gtype"}
	gtype.Field.FieldName = "gtype"

	uuid := &timgroup_Uuid[Timgroup]{fieldName: "uuid"}
	uuid.Field.FieldName = "uuid"

	createtime := &timgroup_Createtime[Timgroup]{fieldName: "createtime"}
	createtime.Field.FieldName = "createtime"

	status := &timgroup_Status[Timgroup]{fieldName: "status"}
	status.Field.FieldName = "status"

	rbean := &timgroup_Rbean[Timgroup]{fieldName: "rbean"}
	rbean.Field.FieldName = "rbean"

	_r = &Timgroup{Id:id,Gtype:gtype,Uuid:uuid,Createtime:createtime,Status:status,Rbean:rbean}
	s := "timgroup"
	if len(tablename) > 0 && tablename[0] != "" {
		s = tablename[0]
	}
	_r.Init(s, []base.Column[Timgroup]{id,gtype,uuid,createtime,status,rbean})
	return
}

func (t *Timgroup) Encode() ([]byte, error) {
	m := make(map[string]any, 0)
	m["id"] = t.GetId()
	m["gtype"] = t.GetGtype()
	m["uuid"] = t.GetUuid()
	m["createtime"] = t.GetCreatetime()
	m["status"] = t.GetStatus()
	m["rbean"] = t.GetRbean()
	return t.Table.Encode(m)
}

func (t *Timgroup) Decode(bs []byte) (err error) {
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

