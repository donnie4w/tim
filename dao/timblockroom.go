// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdao
//
// datetime :2024-10-06 20:30:14
// gdao version 1.2.0
// dbtype:sqlite ,database:timdb ,tablename:timblockroom

package dao

import (
	"fmt"
	"github.com/donnie4w/gdao"
	"github.com/donnie4w/gdao/base"
	
)

type Timblockroom struct {
	gdao.Table[Timblockroom]

	ID      *base.Field[Timblockroom]
	UNIKID      *base.Field[Timblockroom]
	UUID      *base.Field[Timblockroom]
	TUUID      *base.Field[Timblockroom]
	TIMESERIES      *base.Field[Timblockroom]
	_ID      *int64
	_UNIKID      *int64
	_UUID      *int64
	_TUUID      *int64
	_TIMESERIES      *int64
}

var _Timblockroom_ID = &base.Field[Timblockroom]{"id"}
var _Timblockroom_UNIKID = &base.Field[Timblockroom]{"unikid"}
var _Timblockroom_UUID = &base.Field[Timblockroom]{"uuid"}
var _Timblockroom_TUUID = &base.Field[Timblockroom]{"tuuid"}
var _Timblockroom_TIMESERIES = &base.Field[Timblockroom]{"timeseries"}

func (u *Timblockroom) GetId() (_r int64){
	if u._ID != nil {
		_r = *u._ID
	}
	return
}

func (u *Timblockroom) SetId(arg int64) *Timblockroom{
	u.Put0(u.ID.FieldName, arg)
	u._ID = &arg
	return u
}

func (u *Timblockroom) GetUnikid() (_r int64){
	if u._UNIKID != nil {
		_r = *u._UNIKID
	}
	return
}

func (u *Timblockroom) SetUnikid(arg int64) *Timblockroom{
	u.Put0(u.UNIKID.FieldName, arg)
	u._UNIKID = &arg
	return u
}

func (u *Timblockroom) GetUuid() (_r int64){
	if u._UUID != nil {
		_r = *u._UUID
	}
	return
}

func (u *Timblockroom) SetUuid(arg int64) *Timblockroom{
	u.Put0(u.UUID.FieldName, arg)
	u._UUID = &arg
	return u
}

func (u *Timblockroom) GetTuuid() (_r int64){
	if u._TUUID != nil {
		_r = *u._TUUID
	}
	return
}

func (u *Timblockroom) SetTuuid(arg int64) *Timblockroom{
	u.Put0(u.TUUID.FieldName, arg)
	u._TUUID = &arg
	return u
}

func (u *Timblockroom) GetTimeseries() (_r int64){
	if u._TIMESERIES != nil {
		_r = *u._TIMESERIES
	}
	return
}

func (u *Timblockroom) SetTimeseries(arg int64) *Timblockroom{
	u.Put0(u.TIMESERIES.FieldName, arg)
	u._TIMESERIES = &arg
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
	case "timeseries":
		u.SetTimeseries(base.AsInt64(value))
	}
}

func (t *Timblockroom) ToGdao() {
	t.init("timblockroom")
}

func (t *Timblockroom) Copy(h *Timblockroom) *Timblockroom{
	t.SetId(h.GetId())
	t.SetUnikid(h.GetUnikid())
	t.SetUuid(h.GetUuid())
	t.SetTuuid(h.GetTuuid())
	t.SetTimeseries(h.GetTimeseries())
	return t
}

func (t *Timblockroom) String() string {
	return fmt.Sprint("Id:",t.GetId(), ",","Unikid:",t.GetUnikid(), ",","Uuid:",t.GetUuid(), ",","Tuuid:",t.GetTuuid(), ",","Timeseries:",t.GetTimeseries())
}

func (t *Timblockroom)init(tablename string) {
	t.ID = _Timblockroom_ID
	t.UNIKID = _Timblockroom_UNIKID
	t.UUID = _Timblockroom_UUID
	t.TUUID = _Timblockroom_TUUID
	t.TIMESERIES = _Timblockroom_TIMESERIES
	t.Init(tablename, []base.Column[Timblockroom]{t.ID,t.UNIKID,t.UUID,t.TUUID,t.TIMESERIES})
}

func NewTimblockroom(tablename ...string) (_r *Timblockroom) {
	_r = &Timblockroom{}
	s := "timblockroom"
	if len(tablename) > 0 && tablename[0] != "" {
		s = tablename[0]
	}
	_r.init(s)
	return
}

func (t *Timblockroom) Encode() ([]byte, error) {
	m := make(map[string]any, 0)
	m["id"] = t.GetId()
	m["unikid"] = t.GetUnikid()
	m["uuid"] = t.GetUuid()
	m["tuuid"] = t.GetTuuid()
	m["timeseries"] = t.GetTimeseries()
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

