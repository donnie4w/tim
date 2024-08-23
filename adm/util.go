// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of t source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim
package adm

import (
	"context"
	"github.com/donnie4w/tim/stub"
	"strings"

	goutil "github.com/donnie4w/gofer/util"
	. "github.com/donnie4w/tim/keystore"
)

func ctx2CliContext(ctx context.Context) *pcontext {
	return ctx.Value("CliContext").(*pcontext)
}

func auth(name, pwd string) (b bool) {
	if _r, ok := Admin.GetAdmin(name); ok {
		b = strings.EqualFold(_r.Pwd, goutil.Md5Str(pwd))
	}
	return
}

func noAuthAndClose(cc *pcontext) (b bool) {
	if !cc.isAuth {
		cc.tt.Close()
		b = true
	}
	return
}

func newAdmAck(ok bool) *stub.AdmAck {
	return &stub.AdmAck{Ok: &ok}
}

func nodeToTid(node *string) (_r *stub.Tid) {
	if node != nil {
		_r = &stub.Tid{Node: *node}
	}
	return
}

func admMessageToTimMessage(amb *stub.AdmMessageBean) *stub.TimMessage {
	tm := stub.NewTimMessage()
	tm.MsType = amb.MsType
	tm.OdType = amb.OdType
	tm.ID = amb.ID
	tm.Mid = amb.Mid
	tm.BnType = amb.BnType
	tm.FromTid = nodeToTid(amb.FromNode)
	tm.ToTid = nodeToTid(amb.ToNode)
	tm.RoomTid = nodeToTid(amb.RoomNode)
	tm.DataBinary = amb.DataBinary
	tm.DataString = amb.DataString
	tm.Udtype = amb.Udtype
	tm.Udshow = amb.Udshow
	tm.Extend = amb.Extend
	tm.Extra = amb.Extra
	return tm
}

func admUserBeanToTimUserBean(aub *stub.AdmUserBean) *stub.TimUserBean {
	ub := stub.NewTimUserBean()
	ub.Name = aub.Name
	ub.NickName = aub.NickName
	ub.Brithday = aub.Brithday
	ub.Gender = aub.Gender
	ub.Cover = aub.Cover
	ub.Area = aub.Area
	ub.Createtime = aub.Createtime
	ub.PhotoTidAlbum = aub.PhotoTidAlbum
	ub.Extend = aub.Extend
	ub.Extra = aub.Extra
	return ub
}

func timTidToAdmTid(tid *stub.Tid) *stub.AdmTid {
	at := stub.NewAdmTid()
	at.Node = tid.Node
	at.Domain = tid.Domain
	at.Resource = tid.Resource
	at.Termtyp = tid.Termtyp
	at.Extend = tid.Extend
	return at
}

func admTimRoomToTimRoom(atr *stub.AdmTimRoom) *stub.TimRoomBean {
	rb := stub.NewTimRoomBean()
	rb.Founder = atr.Founder
	rb.Managers = atr.Managers
	rb.Cover = atr.Cover
	rb.Topic = atr.Topic
	rb.Label = atr.Label
	rb.Gtype = atr.Gtype
	rb.Kind = atr.Kind
	rb.Createtime = atr.Createtime
	rb.Extend = atr.Extend
	rb.Extra = atr.Extra
	return rb
}
