// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdao
//
// datetime :2024-10-06 20:30:14
// gdao version 1.2.0
// dbtype:sqlite ,database:timdb ,tablename:timuser

package dao

import (
	"fmt"
	"github.com/donnie4w/gdao"
	"github.com/donnie4w/gdao/base"
	
)

type Timuser struct {
	gdao.Table[Timuser]

	ID      *base.Field[Timuser]
	UUID      *base.Field[Timuser]
	PWD      *base.Field[Timuser]
	CREATETIME      *base.Field[Timuser]
	UBEAN      *base.Field[Timuser]
	TIMESERIES      *base.Field[Timuser]
	_ID      *int64
	_UUID      *int64
	_PWD      *int64
	_CREATETIME      *int64
	_UBEAN      []byte
	_TIMESERIES      *int64
}

var _Timuser_ID = &base.Field[Timuser]{"id"}
var _Timuser_UUID = &base.Field[Timuser]{"uuid"}
var _Timuser_PWD = &base.Field[Timuser]{"pwd"}
var _Timuser_CREATETIME = &base.Field[Timuser]{"createtime"}
var _Timuser_UBEAN = &base.Field[Timuser]{"ubean"}
var _Timuser_TIMESERIES = &base.Field[Timuser]{"timeseries"}

func (u *Timuser) GetId() (_r int64){
	if u._ID != nil {
		_r = *u._ID
	}
	return
}

func (u *Timuser) SetId(arg int64) *Timuser{
	u.Put0(u.ID.FieldName, arg)
	u._ID = &arg
	return u
}

func (u *Timuser) GetUuid() (_r int64){
	if u._UUID != nil {
		_r = *u._UUID
	}
	return
}

func (u *Timuser) SetUuid(arg int64) *Timuser{
	u.Put0(u.UUID.FieldName, arg)
	u._UUID = &arg
	return u
}

func (u *Timuser) GetPwd() (_r int64){
	if u._PWD != nil {
		_r = *u._PWD
	}
	return
}

func (u *Timuser) SetPwd(arg int64) *Timuser{
	u.Put0(u.PWD.FieldName, arg)
	u._PWD = &arg
	return u
}

func (u *Timuser) GetCreatetime() (_r int64){
	if u._CREATETIME != nil {
		_r = *u._CREATETIME
	}
	return
}

func (u *Timuser) SetCreatetime(arg int64) *Timuser{
	u.Put0(u.CREATETIME.FieldName, arg)
	u._CREATETIME = &arg
	return u
}

func (u *Timuser) GetUbean() (_r []byte){
	_r = u._UBEAN
	return
}

func (u *Timuser) SetUbean(arg []byte) *Timuser{
	u.Put0(u.UBEAN.FieldName, arg)
	u._UBEAN = arg
	return u
}

func (u *Timuser) GetTimeseries() (_r int64){
	if u._TIMESERIES != nil {
		_r = *u._TIMESERIES
	}
	return
}

func (u *Timuser) SetTimeseries(arg int64) *Timuser{
	u.Put0(u.TIMESERIES.FieldName, arg)
	u._TIMESERIES = &arg
	return u
}


func (u *Timuser) Scan(fieldname string, value any) {
	switch fieldname {
	case "id":
		u.SetId(base.AsInt64(value))
	case "uuid":
		u.SetUuid(base.AsInt64(value))
	case "pwd":
		u.SetPwd(base.AsInt64(value))
	case "createtime":
		u.SetCreatetime(base.AsInt64(value))
	case "ubean":
		u.SetUbean(base.AsBytes(value))
	case "timeseries":
		u.SetTimeseries(base.AsInt64(value))
	}
}

func (t *Timuser) ToGdao() {
	t.init("timuser")
}

func (t *Timuser) Copy(h *Timuser) *Timuser{
	t.SetId(h.GetId())
	t.SetUuid(h.GetUuid())
	t.SetPwd(h.GetPwd())
	t.SetCreatetime(h.GetCreatetime())
	t.SetUbean(h.GetUbean())
	t.SetTimeseries(h.GetTimeseries())
	return t
}

func (t *Timuser) String() string {
	return fmt.Sprint("Id:",t.GetId(), ",","Uuid:",t.GetUuid(), ",","Pwd:",t.GetPwd(), ",","Createtime:",t.GetCreatetime(), ",","Ubean:",t.GetUbean(), ",","Timeseries:",t.GetTimeseries())
}

func (t *Timuser)init(tablename string) {
	t.ID = _Timuser_ID
	t.UUID = _Timuser_UUID
	t.PWD = _Timuser_PWD
	t.CREATETIME = _Timuser_CREATETIME
	t.UBEAN = _Timuser_UBEAN
	t.TIMESERIES = _Timuser_TIMESERIES
	t.Init(tablename, []base.Column[Timuser]{t.ID,t.UUID,t.PWD,t.CREATETIME,t.UBEAN,t.TIMESERIES})
}

func NewTimuser(tablename ...string) (_r *Timuser) {
	_r = &Timuser{}
	s := "timuser"
	if len(tablename) > 0 && tablename[0] != "" {
		s = tablename[0]
	}
	_r.init(s)
	return
}

func (t *Timuser) Encode() ([]byte, error) {
	m := make(map[string]any, 0)
	m["id"] = t.GetId()
	m["uuid"] = t.GetUuid()
	m["pwd"] = t.GetPwd()
	m["createtime"] = t.GetCreatetime()
	m["ubean"] = t.GetUbean()
	m["timeseries"] = t.GetTimeseries()
	return t.Table.Encode(m)
}

func (t *Timuser) Decode(bs []byte) (err error) {
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

