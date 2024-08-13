// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdao
//
// datetime :2024-08-10 01:01:22
// gdao version 1.1.0
// dbtype:sqlite ,database:timdb ,tablename:timmessage

package dao

import (
	"fmt"
	"github.com/donnie4w/gdao"
	"github.com/donnie4w/gdao/base"
	
)

type timmessage_Id[T any] struct {
	base.Field[T]
	fieldName  string
	fieldValue *int64
}

func (t *timmessage_Id[T]) Name() string {
	return t.fieldName
}

func (t *timmessage_Id[T]) Value() any {
	return t.fieldValue
}

type timmessage_Chatid[T any] struct {
	base.Field[T]
	fieldName  string
	fieldValue *int64
}

func (t *timmessage_Chatid[T]) Name() string {
	return t.fieldName
}

func (t *timmessage_Chatid[T]) Value() any {
	return t.fieldValue
}

type timmessage_Stanza[T any] struct {
	base.Field[T]
	fieldName  string
	fieldValue []byte
}

func (t *timmessage_Stanza[T]) Name() string {
	return t.fieldName
}

func (t *timmessage_Stanza[T]) Value() any {
	return t.fieldValue
}

type Timmessage struct {
	gdao.Table[Timmessage]

	Id		*timmessage_Id[Timmessage]
	Chatid		*timmessage_Chatid[Timmessage]
	Stanza		*timmessage_Stanza[Timmessage]
}

func (u *Timmessage) GetId() (_r int64){
	if u.Id.fieldValue != nil {
		_r = *u.Id.fieldValue
	}
	return
}

func (u *Timmessage) SetId(arg int64) *Timmessage{
	u.Put0(u.Id.fieldName, arg)
	u.Id.fieldValue = &arg
	return u
}

func (u *Timmessage) GetChatid() (_r int64){
	if u.Chatid.fieldValue != nil {
		_r = *u.Chatid.fieldValue
	}
	return
}

func (u *Timmessage) SetChatid(arg int64) *Timmessage{
	u.Put0(u.Chatid.fieldName, arg)
	u.Chatid.fieldValue = &arg
	return u
}

func (u *Timmessage) GetStanza() (_r []byte){
	_r = u.Stanza.fieldValue
	return
}

func (u *Timmessage) SetStanza(arg []byte) *Timmessage{
	u.Put0(u.Stanza.fieldName, arg)
	u.Stanza.fieldValue = arg
	return u
}


func (u *Timmessage) Scan(fieldname string, value any) {
	switch fieldname {
	case "id":
		u.SetId(base.AsInt64(value))
	case "chatid":
		u.SetChatid(base.AsInt64(value))
	case "stanza":
		u.SetStanza(base.AsBytes(value))
	}
}

func (t *Timmessage) ToGdao() {
	_t := NewTimmessage()
	*t = *_t
}

func (t *Timmessage) Copy(h *Timmessage) *Timmessage{
	t.SetId(h.GetId())
	t.SetChatid(h.GetChatid())
	t.SetStanza(h.GetStanza())
	return t
}

func (t *Timmessage) String() string {
	return fmt.Sprint("Id:",t.GetId(), ",","Chatid:",t.GetChatid(), ",","Stanza:",t.GetStanza())
}

func NewTimmessage(tablename ...string) (_r *Timmessage) {

	id := &timmessage_Id[Timmessage]{fieldName: "id"}
	id.Field.FieldName = "id"

	chatid := &timmessage_Chatid[Timmessage]{fieldName: "chatid"}
	chatid.Field.FieldName = "chatid"

	stanza := &timmessage_Stanza[Timmessage]{fieldName: "stanza"}
	stanza.Field.FieldName = "stanza"

	_r = &Timmessage{Id:id,Chatid:chatid,Stanza:stanza}
	s := "timmessage"
	if len(tablename) > 0 && tablename[0] != "" {
		s = tablename[0]
	}
	_r.Init(s, []base.Column[Timmessage]{id,chatid,stanza})
	return
}

func (t *Timmessage) Encode() ([]byte, error) {
	m := make(map[string]any, 0)
	m["id"] = t.GetId()
	m["chatid"] = t.GetChatid()
	m["stanza"] = t.GetStanza()
	return t.Table.Encode(m)
}

func (t *Timmessage) Decode(bs []byte) (err error) {
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

