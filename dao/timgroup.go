// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdao
//
// datetime :2025-01-06 18:18:47
// gdao version 1.2.0
// dbtype:sqlite ,database:timdb ,tablename:timgroup

package dao

import (
	"fmt"
	"github.com/donnie4w/gdao"
	"github.com/donnie4w/gdao/base"
	
)

type Timgroup struct {
	gdao.Table[Timgroup]

	ID      *base.Field[Timgroup]
	GTYPE      *base.Field[Timgroup]
	UUID      *base.Field[Timgroup]
	CREATETIME      *base.Field[Timgroup]
	STATUS      *base.Field[Timgroup]
	RBEAN      *base.Field[Timgroup]
	TIMESERIES      *base.Field[Timgroup]
	_ID      *int64
	_GTYPE      *int64
	_UUID      *int64
	_CREATETIME      *int64
	_STATUS      *int64
	_RBEAN      []byte
	_TIMESERIES      *int64
}

var _Timgroup_ID = &base.Field[Timgroup]{"id"}
var _Timgroup_GTYPE = &base.Field[Timgroup]{"gtype"}
var _Timgroup_UUID = &base.Field[Timgroup]{"uuid"}
var _Timgroup_CREATETIME = &base.Field[Timgroup]{"createtime"}
var _Timgroup_STATUS = &base.Field[Timgroup]{"status"}
var _Timgroup_RBEAN = &base.Field[Timgroup]{"rbean"}
var _Timgroup_TIMESERIES = &base.Field[Timgroup]{"timeseries"}

func (u *Timgroup) GetId() (_r int64){
	if u._ID != nil {
		_r = *u._ID
	}
	return
}

func (u *Timgroup) SetId(arg int64) *Timgroup{
	u.Put0(u.ID.FieldName, arg)
	u._ID = &arg
	return u
}

func (u *Timgroup) GetGtype() (_r int64){
	if u._GTYPE != nil {
		_r = *u._GTYPE
	}
	return
}

func (u *Timgroup) SetGtype(arg int64) *Timgroup{
	u.Put0(u.GTYPE.FieldName, arg)
	u._GTYPE = &arg
	return u
}

func (u *Timgroup) GetUuid() (_r int64){
	if u._UUID != nil {
		_r = *u._UUID
	}
	return
}

func (u *Timgroup) SetUuid(arg int64) *Timgroup{
	u.Put0(u.UUID.FieldName, arg)
	u._UUID = &arg
	return u
}

func (u *Timgroup) GetCreatetime() (_r int64){
	if u._CREATETIME != nil {
		_r = *u._CREATETIME
	}
	return
}

func (u *Timgroup) SetCreatetime(arg int64) *Timgroup{
	u.Put0(u.CREATETIME.FieldName, arg)
	u._CREATETIME = &arg
	return u
}

func (u *Timgroup) GetStatus() (_r int64){
	if u._STATUS != nil {
		_r = *u._STATUS
	}
	return
}

func (u *Timgroup) SetStatus(arg int64) *Timgroup{
	u.Put0(u.STATUS.FieldName, arg)
	u._STATUS = &arg
	return u
}

func (u *Timgroup) GetRbean() (_r []byte){
	_r = u._RBEAN
	return
}

func (u *Timgroup) SetRbean(arg []byte) *Timgroup{
	u.Put0(u.RBEAN.FieldName, arg)
	u._RBEAN = arg
	return u
}

func (u *Timgroup) GetTimeseries() (_r int64){
	if u._TIMESERIES != nil {
		_r = *u._TIMESERIES
	}
	return
}

func (u *Timgroup) SetTimeseries(arg int64) *Timgroup{
	u.Put0(u.TIMESERIES.FieldName, arg)
	u._TIMESERIES = &arg
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
	case "timeseries":
		u.SetTimeseries(base.AsInt64(value))
	}
}

func (t *Timgroup) ToGdao() {
	t.init("timgroup")
}

func (t *Timgroup) Copy(h *Timgroup) *Timgroup{
	t.SetId(h.GetId())
	t.SetGtype(h.GetGtype())
	t.SetUuid(h.GetUuid())
	t.SetCreatetime(h.GetCreatetime())
	t.SetStatus(h.GetStatus())
	t.SetRbean(h.GetRbean())
	t.SetTimeseries(h.GetTimeseries())
	return t
}

func (t *Timgroup) String() string {
	return fmt.Sprint("Id:",t.GetId(), ",","Gtype:",t.GetGtype(), ",","Uuid:",t.GetUuid(), ",","Createtime:",t.GetCreatetime(), ",","Status:",t.GetStatus(), ",","Rbean:",t.GetRbean(), ",","Timeseries:",t.GetTimeseries())
}

func (t *Timgroup)init(tablename string) {
	t.ID = _Timgroup_ID
	t.GTYPE = _Timgroup_GTYPE
	t.UUID = _Timgroup_UUID
	t.CREATETIME = _Timgroup_CREATETIME
	t.STATUS = _Timgroup_STATUS
	t.RBEAN = _Timgroup_RBEAN
	t.TIMESERIES = _Timgroup_TIMESERIES
	t.Init(tablename, []base.Column[Timgroup]{t.ID,t.GTYPE,t.UUID,t.CREATETIME,t.STATUS,t.RBEAN,t.TIMESERIES})
}

func NewTimgroup(tablename ...string) (_r *Timgroup) {
	_r = &Timgroup{}
	s := "timgroup"
	if len(tablename) > 0 && tablename[0] != "" {
		s = tablename[0]
	}
	_r.init(s)
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
	m["timeseries"] = t.GetTimeseries()
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

