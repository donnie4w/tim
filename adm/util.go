// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of t source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package adm

import (
	"context"
	"github.com/donnie4w/tim/data"
	"github.com/donnie4w/tim/errs"
	"github.com/donnie4w/tim/stub"
	"github.com/donnie4w/tim/sys"
	"strings"

	goutil "github.com/donnie4w/gofer/util"
	"github.com/donnie4w/tim/keystore"
)

func ctx2CliContext(ctx context.Context) *pcontext {
	return ctx.Value("CliContext").(*pcontext)
}

func auth(name, pwd, domain string) (b bool) {
	if sys.Conf.TimAdminAuth {
		return data.Service.TimAdminAuth(name, pwd, domain)
	} else if _r, ok := keystore.Admin.GetAdmin(name); ok {
		b = strings.EqualFold(_r.Pwd, goutil.Md5Str(pwd))
	}
	return
}

func noAuthAndClose(cc *pcontext) (err errs.ERROR) {
	if !cc.isAuth {
		cc.tt.Close()
		return errs.ERR_PERM_DENIED
	}
	return nil
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

func admMessageToTimMessage(amb *stub.AdmMessage) *stub.TimMessage {
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

func timMessageToAdmMessage(tm *stub.TimMessage) *stub.AdmMessage {
	am := stub.NewAdmMessage()
	am.MsType = tm.MsType
	am.OdType = tm.OdType
	am.ID = tm.ID
	am.Mid = tm.Mid
	am.BnType = tm.BnType

	if tm.FromTid != nil {
		node := tm.FromTid.Node
		am.FromNode = &node
	}

	if tm.ToTid != nil {
		node := tm.ToTid.Node
		am.ToNode = &node
	}

	if tm.RoomTid != nil {
		node := tm.RoomTid.Node
		am.RoomNode = &node
	}

	am.DataBinary = tm.DataBinary
	am.DataString = tm.DataString
	am.Udtype = tm.Udtype
	am.Udshow = tm.Udshow
	am.Extend = tm.Extend
	am.Extra = tm.Extra
	return am
}

func admPresenceToTimPresence(amp *stub.AdmPresence) *stub.TimPresence {
	tp := stub.NewTimPresence()
	tp.ID = amp.ID
	tp.FromTid = nodeToTid(amp.FromNode)
	tp.ToTid = nodeToTid(amp.ToNode)
	tp.ToList = amp.ToList
	tp.Offline = amp.Offline
	tp.Status = amp.Status
	tp.Show = amp.Show
	tp.Extra = amp.Extra
	tp.Extend = amp.Extend
	tp.Extra = amp.Extra
	return tp
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

func timAckToAdmAck(ta *stub.TimAck) *stub.AdmAck {
	if ta != nil {
		aa := stub.NewAdmAck()
		aa.Ok = &ta.Ok
		aa.T = ta.T
		aa.T2 = ta.T2
		aa.N = ta.N
		aa.N2 = ta.N2
		if ta.TimType != 0 {
			aa.TimType = &ta.TimType
		}
		if ta.Error != nil {
			aa.Errcode = ta.Error.Code
		}
		return aa
	}
	return nil
}
