// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdao
//
// datetime :2024-10-06 20:30:14
// gdao version 1.2.0
// dbtype:sqlite ,database:timdb ,tablename:timroster

package dao

import (
	"fmt"
	"github.com/donnie4w/gdao"
	"github.com/donnie4w/gdao/base"
	
)

type Timroster struct {
	gdao.Table[Timroster]

	ID      *base.Field[Timroster]
	UNIKID      *base.Field[Timroster]
	UUID      *base.Field[Timroster]
	TUUID      *base.Field[Timroster]
	TIMESERIES      *base.Field[Timroster]
	_ID      *int64
	_UNIKID      *int64
	_UUID      *int64
	_TUUID      *int64
	_TIMESERIES      *int64
}

var _Timroster_ID = &base.Field[Timroster]{"id"}
var _Timroster_UNIKID = &base.Field[Timroster]{"unikid"}
var _Timroster_UUID = &base.Field[Timroster]{"uuid"}
var _Timroster_TUUID = &base.Field[Timroster]{"tuuid"}
var _Timroster_TIMESERIES = &base.Field[Timroster]{"timeseries"}

func (u *Timroster) GetId() (_r int64){
	if u._ID != nil {
		_r = *u._ID
	}
	return
}

func (u *Timroster) SetId(arg int64) *Timroster{
	u.Put0(u.ID.FieldName, arg)
	u._ID = &arg
	return u
}

func (u *Timroster) GetUnikid() (_r int64){
	if u._UNIKID != nil {
		_r = *u._UNIKID
	}
	return
}

func (u *Timroster) SetUnikid(arg int64) *Timroster{
	u.Put0(u.UNIKID.FieldName, arg)
	u._UNIKID = &arg
	return u
}

func (u *Timroster) GetUuid() (_r int64){
	if u._UUID != nil {
		_r = *u._UUID
	}
	return
}

func (u *Timroster) SetUuid(arg int64) *Timroster{
	u.Put0(u.UUID.FieldName, arg)
	u._UUID = &arg
	return u
}

func (u *Timroster) GetTuuid() (_r int64){
	if u._TUUID != nil {
		_r = *u._TUUID
	}
	return
}

func (u *Timroster) SetTuuid(arg int64) *Timroster{
	u.Put0(u.TUUID.FieldName, arg)
	u._TUUID = &arg
	return u
}

func (u *Timroster) GetTimeseries() (_r int64){
	if u._TIMESERIES != nil {
		_r = *u._TIMESERIES
	}
	return
}

func (u *Timroster) SetTimeseries(arg int64) *Timroster{
	u.Put0(u.TIMESERIES.FieldName, arg)
	u._TIMESERIES = &arg
	return u
}


func (u *Timroster) Scan(fieldname string, value any) {
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

func (t *Timroster) ToGdao() {
	t.init("timroster")
}

func (t *Timroster) Copy(h *Timroster) *Timroster{
	t.SetId(h.GetId())
	t.SetUnikid(h.GetUnikid())
	t.SetUuid(h.GetUuid())
	t.SetTuuid(h.GetTuuid())
	t.SetTimeseries(h.GetTimeseries())
	return t
}

func (t *Timroster) String() string {
	return fmt.Sprint("Id:",t.GetId(), ",","Unikid:",t.GetUnikid(), ",","Uuid:",t.GetUuid(), ",","Tuuid:",t.GetTuuid(), ",","Timeseries:",t.GetTimeseries())
}

func (t *Timroster)init(tablename string) {
	t.ID = _Timroster_ID
	t.UNIKID = _Timroster_UNIKID
	t.UUID = _Timroster_UUID
	t.TUUID = _Timroster_TUUID
	t.TIMESERIES = _Timroster_TIMESERIES
	t.Init(tablename, []base.Column[Timroster]{t.ID,t.UNIKID,t.UUID,t.TUUID,t.TIMESERIES})
}

func NewTimroster(tablename ...string) (_r *Timroster) {
	_r = &Timroster{}
	s := "timroster"
	if len(tablename) > 0 && tablename[0] != "" {
		s = tablename[0]
	}
	_r.init(s)
	return
}

func (t *Timroster) Encode() ([]byte, error) {
	m := make(map[string]any, 0)
	m["id"] = t.GetId()
	m["unikid"] = t.GetUnikid()
	m["uuid"] = t.GetUuid()
	m["tuuid"] = t.GetTuuid()
	m["timeseries"] = t.GetTimeseries()
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

