// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdao
//
// datetime :2025-01-06 18:18:47
// gdao version 1.2.0
// dbtype:sqlite ,database:timdb ,tablename:timoffline

package dao

import (
	"fmt"
	"github.com/donnie4w/gdao"
	"github.com/donnie4w/gdao/base"
	
)

type Timoffline struct {
	gdao.Table[Timoffline]

	ID      *base.Field[Timoffline]
	UUID      *base.Field[Timoffline]
	CHATID      *base.Field[Timoffline]
	STANZA      *base.Field[Timoffline]
	MID      *base.Field[Timoffline]
	TIMESERIES      *base.Field[Timoffline]
	_ID      *int64
	_UUID      *int64
	_CHATID      []byte
	_STANZA      []byte
	_MID      *int64
	_TIMESERIES      *int64
}

var _Timoffline_ID = &base.Field[Timoffline]{"id"}
var _Timoffline_UUID = &base.Field[Timoffline]{"uuid"}
var _Timoffline_CHATID = &base.Field[Timoffline]{"chatid"}
var _Timoffline_STANZA = &base.Field[Timoffline]{"stanza"}
var _Timoffline_MID = &base.Field[Timoffline]{"mid"}
var _Timoffline_TIMESERIES = &base.Field[Timoffline]{"timeseries"}

func (u *Timoffline) GetId() (_r int64){
	if u._ID != nil {
		_r = *u._ID
	}
	return
}

func (u *Timoffline) SetId(arg int64) *Timoffline{
	u.Put0(u.ID.FieldName, arg)
	u._ID = &arg
	return u
}

func (u *Timoffline) GetUuid() (_r int64){
	if u._UUID != nil {
		_r = *u._UUID
	}
	return
}

func (u *Timoffline) SetUuid(arg int64) *Timoffline{
	u.Put0(u.UUID.FieldName, arg)
	u._UUID = &arg
	return u
}

func (u *Timoffline) GetChatid() (_r []byte){
	_r = u._CHATID
	return
}

func (u *Timoffline) SetChatid(arg []byte) *Timoffline{
	u.Put0(u.CHATID.FieldName, arg)
	u._CHATID = arg
	return u
}

func (u *Timoffline) GetStanza() (_r []byte){
	_r = u._STANZA
	return
}

func (u *Timoffline) SetStanza(arg []byte) *Timoffline{
	u.Put0(u.STANZA.FieldName, arg)
	u._STANZA = arg
	return u
}

func (u *Timoffline) GetMid() (_r int64){
	if u._MID != nil {
		_r = *u._MID
	}
	return
}

func (u *Timoffline) SetMid(arg int64) *Timoffline{
	u.Put0(u.MID.FieldName, arg)
	u._MID = &arg
	return u
}

func (u *Timoffline) GetTimeseries() (_r int64){
	if u._TIMESERIES != nil {
		_r = *u._TIMESERIES
	}
	return
}

func (u *Timoffline) SetTimeseries(arg int64) *Timoffline{
	u.Put0(u.TIMESERIES.FieldName, arg)
	u._TIMESERIES = &arg
	return u
}


func (u *Timoffline) Scan(fieldname string, value any) {
	switch fieldname {
	case "id":
		u.SetId(base.AsInt64(value))
	case "uuid":
		u.SetUuid(base.AsInt64(value))
	case "chatid":
		u.SetChatid(base.AsBytes(value))
	case "stanza":
		u.SetStanza(base.AsBytes(value))
	case "mid":
		u.SetMid(base.AsInt64(value))
	case "timeseries":
		u.SetTimeseries(base.AsInt64(value))
	}
}

func (t *Timoffline) ToGdao() {
	t.init("timoffline")
}

func (t *Timoffline) Copy(h *Timoffline) *Timoffline{
	t.SetId(h.GetId())
	t.SetUuid(h.GetUuid())
	t.SetChatid(h.GetChatid())
	t.SetStanza(h.GetStanza())
	t.SetMid(h.GetMid())
	t.SetTimeseries(h.GetTimeseries())
	return t
}

func (t *Timoffline) String() string {
	return fmt.Sprint("Id:",t.GetId(), ",","Uuid:",t.GetUuid(), ",","Chatid:",t.GetChatid(), ",","Stanza:",t.GetStanza(), ",","Mid:",t.GetMid(), ",","Timeseries:",t.GetTimeseries())
}

func (t *Timoffline)init(tablename string) {
	t.ID = _Timoffline_ID
	t.UUID = _Timoffline_UUID
	t.CHATID = _Timoffline_CHATID
	t.STANZA = _Timoffline_STANZA
	t.MID = _Timoffline_MID
	t.TIMESERIES = _Timoffline_TIMESERIES
	t.Init(tablename, []base.Column[Timoffline]{t.ID,t.UUID,t.CHATID,t.STANZA,t.MID,t.TIMESERIES})
}

func NewTimoffline(tablename ...string) (_r *Timoffline) {
	_r = &Timoffline{}
	s := "timoffline"
	if len(tablename) > 0 && tablename[0] != "" {
		s = tablename[0]
	}
	_r.init(s)
	return
}

func (t *Timoffline) Encode() ([]byte, error) {
	m := make(map[string]any, 0)
	m["id"] = t.GetId()
	m["uuid"] = t.GetUuid()
	m["chatid"] = t.GetChatid()
	m["stanza"] = t.GetStanza()
	m["mid"] = t.GetMid()
	m["timeseries"] = t.GetTimeseries()
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

