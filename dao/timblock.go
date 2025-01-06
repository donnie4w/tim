// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdao
//
// datetime :2024-10-06 20:30:14
// gdao version 1.2.0
// dbtype:sqlite ,database:timdb ,tablename:timblock

package dao

import (
	"fmt"
	"github.com/donnie4w/gdao"
	"github.com/donnie4w/gdao/base"
	
)

type Timblock struct {
	gdao.Table[Timblock]

	ID      *base.Field[Timblock]
	UNIKID      *base.Field[Timblock]
	UUID      *base.Field[Timblock]
	TUUID      *base.Field[Timblock]
	TIMESERIES      *base.Field[Timblock]
	_ID      *int64
	_UNIKID      *int64
	_UUID      *int64
	_TUUID      *int64
	_TIMESERIES      *int64
}

var _Timblock_ID = &base.Field[Timblock]{"id"}
var _Timblock_UNIKID = &base.Field[Timblock]{"unikid"}
var _Timblock_UUID = &base.Field[Timblock]{"uuid"}
var _Timblock_TUUID = &base.Field[Timblock]{"tuuid"}
var _Timblock_TIMESERIES = &base.Field[Timblock]{"timeseries"}

func (u *Timblock) GetId() (_r int64){
	if u._ID != nil {
		_r = *u._ID
	}
	return
}

func (u *Timblock) SetId(arg int64) *Timblock{
	u.Put0(u.ID.FieldName, arg)
	u._ID = &arg
	return u
}

func (u *Timblock) GetUnikid() (_r int64){
	if u._UNIKID != nil {
		_r = *u._UNIKID
	}
	return
}

func (u *Timblock) SetUnikid(arg int64) *Timblock{
	u.Put0(u.UNIKID.FieldName, arg)
	u._UNIKID = &arg
	return u
}

func (u *Timblock) GetUuid() (_r int64){
	if u._UUID != nil {
		_r = *u._UUID
	}
	return
}

func (u *Timblock) SetUuid(arg int64) *Timblock{
	u.Put0(u.UUID.FieldName, arg)
	u._UUID = &arg
	return u
}

func (u *Timblock) GetTuuid() (_r int64){
	if u._TUUID != nil {
		_r = *u._TUUID
	}
	return
}

func (u *Timblock) SetTuuid(arg int64) *Timblock{
	u.Put0(u.TUUID.FieldName, arg)
	u._TUUID = &arg
	return u
}

func (u *Timblock) GetTimeseries() (_r int64){
	if u._TIMESERIES != nil {
		_r = *u._TIMESERIES
	}
	return
}

func (u *Timblock) SetTimeseries(arg int64) *Timblock{
	u.Put0(u.TIMESERIES.FieldName, arg)
	u._TIMESERIES = &arg
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
	case "timeseries":
		u.SetTimeseries(base.AsInt64(value))
	}
}

func (t *Timblock) ToGdao() {
	t.init("timblock")
}

func (t *Timblock) Copy(h *Timblock) *Timblock{
	t.SetId(h.GetId())
	t.SetUnikid(h.GetUnikid())
	t.SetUuid(h.GetUuid())
	t.SetTuuid(h.GetTuuid())
	t.SetTimeseries(h.GetTimeseries())
	return t
}

func (t *Timblock) String() string {
	return fmt.Sprint("Id:",t.GetId(), ",","Unikid:",t.GetUnikid(), ",","Uuid:",t.GetUuid(), ",","Tuuid:",t.GetTuuid(), ",","Timeseries:",t.GetTimeseries())
}

func (t *Timblock)init(tablename string) {
	t.ID = _Timblock_ID
	t.UNIKID = _Timblock_UNIKID
	t.UUID = _Timblock_UUID
	t.TUUID = _Timblock_TUUID
	t.TIMESERIES = _Timblock_TIMESERIES
	t.Init(tablename, []base.Column[Timblock]{t.ID,t.UNIKID,t.UUID,t.TUUID,t.TIMESERIES})
}

func NewTimblock(tablename ...string) (_r *Timblock) {
	_r = &Timblock{}
	s := "timblock"
	if len(tablename) > 0 && tablename[0] != "" {
		s = tablename[0]
	}
	_r.init(s)
	return
}

func (t *Timblock) Encode() ([]byte, error) {
	m := make(map[string]any, 0)
	m["id"] = t.GetId()
	m["unikid"] = t.GetUnikid()
	m["uuid"] = t.GetUuid()
	m["tuuid"] = t.GetTuuid()
	m["timeseries"] = t.GetTimeseries()
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

