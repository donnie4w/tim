// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdao
//
// datetime :2024-10-06 20:30:14
// gdao version 1.2.0
// dbtype:sqlite ,database:timdb ,tablename:timdomain

package dao

import (
	"fmt"
	"github.com/donnie4w/gdao"
	"github.com/donnie4w/gdao/base"
	
)

type Timdomain struct {
	gdao.Table[Timdomain]

	ID      *base.Field[Timdomain]
	ADMINACCOUNT      *base.Field[Timdomain]
	ADMINPASSWORD      *base.Field[Timdomain]
	TIMDOMAIN      *base.Field[Timdomain]
	CREATETIME      *base.Field[Timdomain]
	TIMESERIES      *base.Field[Timdomain]
	_ID      *int64
	_ADMINACCOUNT      *string
	_ADMINPASSWORD      *string
	_TIMDOMAIN      *string
	_CREATETIME      *int64
	_TIMESERIES      *int64
}

var _Timdomain_ID = &base.Field[Timdomain]{"id"}
var _Timdomain_ADMINACCOUNT = &base.Field[Timdomain]{"adminaccount"}
var _Timdomain_ADMINPASSWORD = &base.Field[Timdomain]{"adminpassword"}
var _Timdomain_TIMDOMAIN = &base.Field[Timdomain]{"timdomain"}
var _Timdomain_CREATETIME = &base.Field[Timdomain]{"createtime"}
var _Timdomain_TIMESERIES = &base.Field[Timdomain]{"timeseries"}

func (u *Timdomain) GetId() (_r int64){
	if u._ID != nil {
		_r = *u._ID
	}
	return
}

func (u *Timdomain) SetId(arg int64) *Timdomain{
	u.Put0(u.ID.FieldName, arg)
	u._ID = &arg
	return u
}

func (u *Timdomain) GetAdminaccount() (_r string){
	if u._ADMINACCOUNT != nil {
		_r = *u._ADMINACCOUNT
	}
	return
}

func (u *Timdomain) SetAdminaccount(arg string) *Timdomain{
	u.Put0(u.ADMINACCOUNT.FieldName, arg)
	u._ADMINACCOUNT = &arg
	return u
}

func (u *Timdomain) GetAdminpassword() (_r string){
	if u._ADMINPASSWORD != nil {
		_r = *u._ADMINPASSWORD
	}
	return
}

func (u *Timdomain) SetAdminpassword(arg string) *Timdomain{
	u.Put0(u.ADMINPASSWORD.FieldName, arg)
	u._ADMINPASSWORD = &arg
	return u
}

func (u *Timdomain) GetTimdomain() (_r string){
	if u._TIMDOMAIN != nil {
		_r = *u._TIMDOMAIN
	}
	return
}

func (u *Timdomain) SetTimdomain(arg string) *Timdomain{
	u.Put0(u.TIMDOMAIN.FieldName, arg)
	u._TIMDOMAIN = &arg
	return u
}

func (u *Timdomain) GetCreatetime() (_r int64){
	if u._CREATETIME != nil {
		_r = *u._CREATETIME
	}
	return
}

func (u *Timdomain) SetCreatetime(arg int64) *Timdomain{
	u.Put0(u.CREATETIME.FieldName, arg)
	u._CREATETIME = &arg
	return u
}

func (u *Timdomain) GetTimeseries() (_r int64){
	if u._TIMESERIES != nil {
		_r = *u._TIMESERIES
	}
	return
}

func (u *Timdomain) SetTimeseries(arg int64) *Timdomain{
	u.Put0(u.TIMESERIES.FieldName, arg)
	u._TIMESERIES = &arg
	return u
}


func (u *Timdomain) Scan(fieldname string, value any) {
	switch fieldname {
	case "id":
		u.SetId(base.AsInt64(value))
	case "adminaccount":
		u.SetAdminaccount(base.AsString(value))
	case "adminpassword":
		u.SetAdminpassword(base.AsString(value))
	case "timdomain":
		u.SetTimdomain(base.AsString(value))
	case "createtime":
		u.SetCreatetime(base.AsInt64(value))
	case "timeseries":
		u.SetTimeseries(base.AsInt64(value))
	}
}

func (t *Timdomain) ToGdao() {
	t.init("timdomain")
}

func (t *Timdomain) Copy(h *Timdomain) *Timdomain{
	t.SetId(h.GetId())
	t.SetAdminaccount(h.GetAdminaccount())
	t.SetAdminpassword(h.GetAdminpassword())
	t.SetTimdomain(h.GetTimdomain())
	t.SetCreatetime(h.GetCreatetime())
	t.SetTimeseries(h.GetTimeseries())
	return t
}

func (t *Timdomain) String() string {
	return fmt.Sprint("Id:",t.GetId(), ",","Adminaccount:",t.GetAdminaccount(), ",","Adminpassword:",t.GetAdminpassword(), ",","Timdomain:",t.GetTimdomain(), ",","Createtime:",t.GetCreatetime(), ",","Timeseries:",t.GetTimeseries())
}

func (t *Timdomain)init(tablename string) {
	t.ID = _Timdomain_ID
	t.ADMINACCOUNT = _Timdomain_ADMINACCOUNT
	t.ADMINPASSWORD = _Timdomain_ADMINPASSWORD
	t.TIMDOMAIN = _Timdomain_TIMDOMAIN
	t.CREATETIME = _Timdomain_CREATETIME
	t.TIMESERIES = _Timdomain_TIMESERIES
	t.Init(tablename, []base.Column[Timdomain]{t.ID,t.ADMINACCOUNT,t.ADMINPASSWORD,t.TIMDOMAIN,t.CREATETIME,t.TIMESERIES})
}

func NewTimdomain(tablename ...string) (_r *Timdomain) {
	_r = &Timdomain{}
	s := "timdomain"
	if len(tablename) > 0 && tablename[0] != "" {
		s = tablename[0]
	}
	_r.init(s)
	return
}

func (t *Timdomain) Encode() ([]byte, error) {
	m := make(map[string]any, 0)
	m["id"] = t.GetId()
	m["adminaccount"] = t.GetAdminaccount()
	m["adminpassword"] = t.GetAdminpassword()
	m["timdomain"] = t.GetTimdomain()
	m["createtime"] = t.GetCreatetime()
	m["timeseries"] = t.GetTimeseries()
	return t.Table.Encode(m)
}

func (t *Timdomain) Decode(bs []byte) (err error) {
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

