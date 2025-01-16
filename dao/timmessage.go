// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdao
//
// datetime :2025-01-06 18:18:47
// gdao version 1.2.0
// dbtype:sqlite ,database:timdb ,tablename:timmessage

package dao

import (
	"fmt"
	"github.com/donnie4w/gdao"
	"github.com/donnie4w/gdao/base"
	
)

type Timmessage struct {
	gdao.Table[Timmessage]

	ID      *base.Field[Timmessage]
	CHATID      *base.Field[Timmessage]
	FID      *base.Field[Timmessage]
	STANZA      *base.Field[Timmessage]
	TIMESERIES      *base.Field[Timmessage]
	_ID      *int64
	_CHATID      []byte
	_FID      *int64
	_STANZA      []byte
	_TIMESERIES      *int64
}

var _Timmessage_ID = &base.Field[Timmessage]{"id"}
var _Timmessage_CHATID = &base.Field[Timmessage]{"chatid"}
var _Timmessage_FID = &base.Field[Timmessage]{"fid"}
var _Timmessage_STANZA = &base.Field[Timmessage]{"stanza"}
var _Timmessage_TIMESERIES = &base.Field[Timmessage]{"timeseries"}

func (u *Timmessage) GetId() (_r int64){
	if u._ID != nil {
		_r = *u._ID
	}
	return
}

func (u *Timmessage) SetId(arg int64) *Timmessage{
	u.Put0(u.ID.FieldName, arg)
	u._ID = &arg
	return u
}

func (u *Timmessage) GetChatid() (_r []byte){
	_r = u._CHATID
	return
}

func (u *Timmessage) SetChatid(arg []byte) *Timmessage{
	u.Put0(u.CHATID.FieldName, arg)
	u._CHATID = arg
	return u
}

func (u *Timmessage) GetFid() (_r int64){
	if u._FID != nil {
		_r = *u._FID
	}
	return
}

func (u *Timmessage) SetFid(arg int64) *Timmessage{
	u.Put0(u.FID.FieldName, arg)
	u._FID = &arg
	return u
}

func (u *Timmessage) GetStanza() (_r []byte){
	_r = u._STANZA
	return
}

func (u *Timmessage) SetStanza(arg []byte) *Timmessage{
	u.Put0(u.STANZA.FieldName, arg)
	u._STANZA = arg
	return u
}

func (u *Timmessage) GetTimeseries() (_r int64){
	if u._TIMESERIES != nil {
		_r = *u._TIMESERIES
	}
	return
}

func (u *Timmessage) SetTimeseries(arg int64) *Timmessage{
	u.Put0(u.TIMESERIES.FieldName, arg)
	u._TIMESERIES = &arg
	return u
}


func (u *Timmessage) Scan(fieldname string, value any) {
	switch fieldname {
	case "id":
		u.SetId(base.AsInt64(value))
	case "chatid":
		u.SetChatid(base.AsBytes(value))
	case "fid":
		u.SetFid(base.AsInt64(value))
	case "stanza":
		u.SetStanza(base.AsBytes(value))
	case "timeseries":
		u.SetTimeseries(base.AsInt64(value))
	}
}

func (t *Timmessage) ToGdao() {
	t.init("timmessage")
}

func (t *Timmessage) Copy(h *Timmessage) *Timmessage{
	t.SetId(h.GetId())
	t.SetChatid(h.GetChatid())
	t.SetFid(h.GetFid())
	t.SetStanza(h.GetStanza())
	t.SetTimeseries(h.GetTimeseries())
	return t
}

func (t *Timmessage) String() string {
	return fmt.Sprint("Id:",t.GetId(), ",","Chatid:",t.GetChatid(), ",","Fid:",t.GetFid(), ",","Stanza:",t.GetStanza(), ",","Timeseries:",t.GetTimeseries())
}

func (t *Timmessage)init(tablename string) {
	t.ID = _Timmessage_ID
	t.CHATID = _Timmessage_CHATID
	t.FID = _Timmessage_FID
	t.STANZA = _Timmessage_STANZA
	t.TIMESERIES = _Timmessage_TIMESERIES
	t.Init(tablename, []base.Column[Timmessage]{t.ID,t.CHATID,t.FID,t.STANZA,t.TIMESERIES})
}

func NewTimmessage(tablename ...string) (_r *Timmessage) {
	_r = &Timmessage{}
	s := "timmessage"
	if len(tablename) > 0 && tablename[0] != "" {
		s = tablename[0]
	}
	_r.init(s)
	return
}

func (t *Timmessage) Encode() ([]byte, error) {
	m := make(map[string]any, 0)
	m["id"] = t.GetId()
	m["chatid"] = t.GetChatid()
	m["fid"] = t.GetFid()
	m["stanza"] = t.GetStanza()
	m["timeseries"] = t.GetTimeseries()
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

