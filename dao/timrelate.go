// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdao
//
// datetime :2024-10-06 20:30:14
// gdao version 1.2.0
// dbtype:sqlite ,database:timdb ,tablename:timrelate

package dao

import (
	"fmt"
	"github.com/donnie4w/gdao"
	"github.com/donnie4w/gdao/base"
	
)

type Timrelate struct {
	gdao.Table[Timrelate]

	ID      *base.Field[Timrelate]
	UUID      *base.Field[Timrelate]
	STATUS      *base.Field[Timrelate]
	TIMESERIES      *base.Field[Timrelate]
	_ID      *int64
	_UUID      *int64
	_STATUS      *int64
	_TIMESERIES      *int64
}

var _Timrelate_ID = &base.Field[Timrelate]{"id"}
var _Timrelate_UUID = &base.Field[Timrelate]{"uuid"}
var _Timrelate_STATUS = &base.Field[Timrelate]{"status"}
var _Timrelate_TIMESERIES = &base.Field[Timrelate]{"timeseries"}

func (u *Timrelate) GetId() (_r int64){
	if u._ID != nil {
		_r = *u._ID
	}
	return
}

func (u *Timrelate) SetId(arg int64) *Timrelate{
	u.Put0(u.ID.FieldName, arg)
	u._ID = &arg
	return u
}

func (u *Timrelate) GetUuid() (_r int64){
	if u._UUID != nil {
		_r = *u._UUID
	}
	return
}

func (u *Timrelate) SetUuid(arg int64) *Timrelate{
	u.Put0(u.UUID.FieldName, arg)
	u._UUID = &arg
	return u
}

func (u *Timrelate) GetStatus() (_r int64){
	if u._STATUS != nil {
		_r = *u._STATUS
	}
	return
}

func (u *Timrelate) SetStatus(arg int64) *Timrelate{
	u.Put0(u.STATUS.FieldName, arg)
	u._STATUS = &arg
	return u
}

func (u *Timrelate) GetTimeseries() (_r int64){
	if u._TIMESERIES != nil {
		_r = *u._TIMESERIES
	}
	return
}

func (u *Timrelate) SetTimeseries(arg int64) *Timrelate{
	u.Put0(u.TIMESERIES.FieldName, arg)
	u._TIMESERIES = &arg
	return u
}


func (u *Timrelate) Scan(fieldname string, value any) {
	switch fieldname {
	case "id":
		u.SetId(base.AsInt64(value))
	case "uuid":
		u.SetUuid(base.AsInt64(value))
	case "status":
		u.SetStatus(base.AsInt64(value))
	case "timeseries":
		u.SetTimeseries(base.AsInt64(value))
	}
}

func (t *Timrelate) ToGdao() {
	t.init("timrelate")
}

func (t *Timrelate) Copy(h *Timrelate) *Timrelate{
	t.SetId(h.GetId())
	t.SetUuid(h.GetUuid())
	t.SetStatus(h.GetStatus())
	t.SetTimeseries(h.GetTimeseries())
	return t
}

func (t *Timrelate) String() string {
	return fmt.Sprint("Id:",t.GetId(), ",","Uuid:",t.GetUuid(), ",","Status:",t.GetStatus(), ",","Timeseries:",t.GetTimeseries())
}

func (t *Timrelate)init(tablename string) {
	t.ID = _Timrelate_ID
	t.UUID = _Timrelate_UUID
	t.STATUS = _Timrelate_STATUS
	t.TIMESERIES = _Timrelate_TIMESERIES
	t.Init(tablename, []base.Column[Timrelate]{t.ID,t.UUID,t.STATUS,t.TIMESERIES})
}

func NewTimrelate(tablename ...string) (_r *Timrelate) {
	_r = &Timrelate{}
	s := "timrelate"
	if len(tablename) > 0 && tablename[0] != "" {
		s = tablename[0]
	}
	_r.init(s)
	return
}

func (t *Timrelate) Encode() ([]byte, error) {
	m := make(map[string]any, 0)
	m["id"] = t.GetId()
	m["uuid"] = t.GetUuid()
	m["status"] = t.GetStatus()
	m["timeseries"] = t.GetTimeseries()
	return t.Table.Encode(m)
}

func (t *Timrelate) Decode(bs []byte) (err error) {
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

