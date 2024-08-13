// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdao
//
// datetime :2024-08-10 01:01:22
// gdao version 1.1.0
// dbtype:sqlite ,database:timdb ,tablename:timoffline

package dao

import (
	"fmt"
	"github.com/donnie4w/gdao"
	"github.com/donnie4w/gdao/base"
	
)

type timoffline_Id[T any] struct {
	base.Field[T]
	fieldName  string
	fieldValue *int64
}

func (t *timoffline_Id[T]) Name() string {
	return t.fieldName
}

func (t *timoffline_Id[T]) Value() any {
	return t.fieldValue
}

type timoffline_Uuid[T any] struct {
	base.Field[T]
	fieldName  string
	fieldValue *int64
}

func (t *timoffline_Uuid[T]) Name() string {
	return t.fieldName
}

func (t *timoffline_Uuid[T]) Value() any {
	return t.fieldValue
}

type timoffline_Chatid[T any] struct {
	base.Field[T]
	fieldName  string
	fieldValue *int64
}

func (t *timoffline_Chatid[T]) Name() string {
	return t.fieldName
}

func (t *timoffline_Chatid[T]) Value() any {
	return t.fieldValue
}

type timoffline_Stanza[T any] struct {
	base.Field[T]
	fieldName  string
	fieldValue []byte
}

func (t *timoffline_Stanza[T]) Name() string {
	return t.fieldName
}

func (t *timoffline_Stanza[T]) Value() any {
	return t.fieldValue
}

type timoffline_Mid[T any] struct {
	base.Field[T]
	fieldName  string
	fieldValue *int64
}

func (t *timoffline_Mid[T]) Name() string {
	return t.fieldName
}

func (t *timoffline_Mid[T]) Value() any {
	return t.fieldValue
}

type Timoffline struct {
	gdao.Table[Timoffline]

	Id		*timoffline_Id[Timoffline]
	Uuid		*timoffline_Uuid[Timoffline]
	Chatid		*timoffline_Chatid[Timoffline]
	Stanza		*timoffline_Stanza[Timoffline]
	Mid		*timoffline_Mid[Timoffline]
}

func (u *Timoffline) GetId() (_r int64){
	if u.Id.fieldValue != nil {
		_r = *u.Id.fieldValue
	}
	return
}

func (u *Timoffline) SetId(arg int64) *Timoffline{
	u.Put0(u.Id.fieldName, arg)
	u.Id.fieldValue = &arg
	return u
}

func (u *Timoffline) GetUuid() (_r int64){
	if u.Uuid.fieldValue != nil {
		_r = *u.Uuid.fieldValue
	}
	return
}

func (u *Timoffline) SetUuid(arg int64) *Timoffline{
	u.Put0(u.Uuid.fieldName, arg)
	u.Uuid.fieldValue = &arg
	return u
}

func (u *Timoffline) GetChatid() (_r int64){
	if u.Chatid.fieldValue != nil {
		_r = *u.Chatid.fieldValue
	}
	return
}

func (u *Timoffline) SetChatid(arg int64) *Timoffline{
	u.Put0(u.Chatid.fieldName, arg)
	u.Chatid.fieldValue = &arg
	return u
}

func (u *Timoffline) GetStanza() (_r []byte){
	_r = u.Stanza.fieldValue
	return
}

func (u *Timoffline) SetStanza(arg []byte) *Timoffline{
	u.Put0(u.Stanza.fieldName, arg)
	u.Stanza.fieldValue = arg
	return u
}

func (u *Timoffline) GetMid() (_r int64){
	if u.Mid.fieldValue != nil {
		_r = *u.Mid.fieldValue
	}
	return
}

func (u *Timoffline) SetMid(arg int64) *Timoffline{
	u.Put0(u.Mid.fieldName, arg)
	u.Mid.fieldValue = &arg
	return u
}


func (u *Timoffline) Scan(fieldname string, value any) {
	switch fieldname {
	case "id":
		u.SetId(base.AsInt64(value))
	case "uuid":
		u.SetUuid(base.AsInt64(value))
	case "chatid":
		u.SetChatid(base.AsInt64(value))
	case "stanza":
		u.SetStanza(base.AsBytes(value))
	case "mid":
		u.SetMid(base.AsInt64(value))
	}
}

func (t *Timoffline) ToGdao() {
	_t := NewTimoffline()
	*t = *_t
}

func (t *Timoffline) Copy(h *Timoffline) *Timoffline{
	t.SetId(h.GetId())
	t.SetUuid(h.GetUuid())
	t.SetChatid(h.GetChatid())
	t.SetStanza(h.GetStanza())
	t.SetMid(h.GetMid())
	return t
}

func (t *Timoffline) String() string {
	return fmt.Sprint("Id:",t.GetId(), ",","Uuid:",t.GetUuid(), ",","Chatid:",t.GetChatid(), ",","Stanza:",t.GetStanza(), ",","Mid:",t.GetMid())
}

func NewTimoffline(tablename ...string) (_r *Timoffline) {

	id := &timoffline_Id[Timoffline]{fieldName: "id"}
	id.Field.FieldName = "id"

	uuid := &timoffline_Uuid[Timoffline]{fieldName: "uuid"}
	uuid.Field.FieldName = "uuid"

	chatid := &timoffline_Chatid[Timoffline]{fieldName: "chatid"}
	chatid.Field.FieldName = "chatid"

	stanza := &timoffline_Stanza[Timoffline]{fieldName: "stanza"}
	stanza.Field.FieldName = "stanza"

	mid := &timoffline_Mid[Timoffline]{fieldName: "mid"}
	mid.Field.FieldName = "mid"

	_r = &Timoffline{Id:id,Uuid:uuid,Chatid:chatid,Stanza:stanza,Mid:mid}
	s := "timoffline"
	if len(tablename) > 0 && tablename[0] != "" {
		s = tablename[0]
	}
	_r.Init(s, []base.Column[Timoffline]{id,uuid,chatid,stanza,mid})
	return
}

func (t *Timoffline) Encode() ([]byte, error) {
	m := make(map[string]any, 0)
	m["id"] = t.GetId()
	m["uuid"] = t.GetUuid()
	m["chatid"] = t.GetChatid()
	m["stanza"] = t.GetStanza()
	m["mid"] = t.GetMid()
	return t.Table.Encode(m)
}

func (t *Timoffline) Decode(bs []byte) (err error) {
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

