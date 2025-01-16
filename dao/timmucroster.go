// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdao
//
// datetime :2025-01-06 18:18:47
// gdao version 1.2.0
// dbtype:sqlite ,database:timdb ,tablename:timmucroster

package dao

import (
	"fmt"
	"github.com/donnie4w/gdao"
	"github.com/donnie4w/gdao/base"
	
)

type Timmucroster struct {
	gdao.Table[Timmucroster]

	ID      *base.Field[Timmucroster]
	UNIKID      *base.Field[Timmucroster]
	UUID      *base.Field[Timmucroster]
	TUUID      *base.Field[Timmucroster]
	TIMESERIES      *base.Field[Timmucroster]
	_ID      *int64
	_UNIKID      []byte
	_UUID      *int64
	_TUUID      *int64
	_TIMESERIES      *int64
}

var _Timmucroster_ID = &base.Field[Timmucroster]{"id"}
var _Timmucroster_UNIKID = &base.Field[Timmucroster]{"unikid"}
var _Timmucroster_UUID = &base.Field[Timmucroster]{"uuid"}
var _Timmucroster_TUUID = &base.Field[Timmucroster]{"tuuid"}
var _Timmucroster_TIMESERIES = &base.Field[Timmucroster]{"timeseries"}

func (u *Timmucroster) GetId() (_r int64){
	if u._ID != nil {
		_r = *u._ID
	}
	return
}

func (u *Timmucroster) SetId(arg int64) *Timmucroster{
	u.Put0(u.ID.FieldName, arg)
	u._ID = &arg
	return u
}

func (u *Timmucroster) GetUnikid() (_r []byte){
	_r = u._UNIKID
	return
}

func (u *Timmucroster) SetUnikid(arg []byte) *Timmucroster{
	u.Put0(u.UNIKID.FieldName, arg)
	u._UNIKID = arg
	return u
}

func (u *Timmucroster) GetUuid() (_r int64){
	if u._UUID != nil {
		_r = *u._UUID
	}
	return
}

func (u *Timmucroster) SetUuid(arg int64) *Timmucroster{
	u.Put0(u.UUID.FieldName, arg)
	u._UUID = &arg
	return u
}

func (u *Timmucroster) GetTuuid() (_r int64){
	if u._TUUID != nil {
		_r = *u._TUUID
	}
	return
}

func (u *Timmucroster) SetTuuid(arg int64) *Timmucroster{
	u.Put0(u.TUUID.FieldName, arg)
	u._TUUID = &arg
	return u
}

func (u *Timmucroster) GetTimeseries() (_r int64){
	if u._TIMESERIES != nil {
		_r = *u._TIMESERIES
	}
	return
}

func (u *Timmucroster) SetTimeseries(arg int64) *Timmucroster{
	u.Put0(u.TIMESERIES.FieldName, arg)
	u._TIMESERIES = &arg
	return u
}


func (u *Timmucroster) Scan(fieldname string, value any) {
	switch fieldname {
	case "id":
		u.SetId(base.AsInt64(value))
	case "unikid":
		u.SetUnikid(base.AsBytes(value))
	case "uuid":
		u.SetUuid(base.AsInt64(value))
	case "tuuid":
		u.SetTuuid(base.AsInt64(value))
	case "timeseries":
		u.SetTimeseries(base.AsInt64(value))
	}
}

func (t *Timmucroster) ToGdao() {
	t.init("timmucroster")
}

func (t *Timmucroster) Copy(h *Timmucroster) *Timmucroster{
	t.SetId(h.GetId())
	t.SetUnikid(h.GetUnikid())
	t.SetUuid(h.GetUuid())
	t.SetTuuid(h.GetTuuid())
	t.SetTimeseries(h.GetTimeseries())
	return t
}

func (t *Timmucroster) String() string {
	return fmt.Sprint("Id:",t.GetId(), ",","Unikid:",t.GetUnikid(), ",","Uuid:",t.GetUuid(), ",","Tuuid:",t.GetTuuid(), ",","Timeseries:",t.GetTimeseries())
}

func (t *Timmucroster)init(tablename string) {
	t.ID = _Timmucroster_ID
	t.UNIKID = _Timmucroster_UNIKID
	t.UUID = _Timmucroster_UUID
	t.TUUID = _Timmucroster_TUUID
	t.TIMESERIES = _Timmucroster_TIMESERIES
	t.Init(tablename, []base.Column[Timmucroster]{t.ID,t.UNIKID,t.UUID,t.TUUID,t.TIMESERIES})
}

func NewTimmucroster(tablename ...string) (_r *Timmucroster) {
	_r = &Timmucroster{}
	s := "timmucroster"
	if len(tablename) > 0 && tablename[0] != "" {
		s = tablename[0]
	}
	_r.init(s)
	return
}

func (t *Timmucroster) Encode() ([]byte, error) {
	m := make(map[string]any, 0)
	m["id"] = t.GetId()
	m["unikid"] = t.GetUnikid()
	m["uuid"] = t.GetUuid()
	m["tuuid"] = t.GetTuuid()
	m["timeseries"] = t.GetTimeseries()
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

