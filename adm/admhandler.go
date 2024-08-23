// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of t source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim
package adm

import (
	. "github.com/donnie4w/tim/stub"
	"github.com/donnie4w/tim/sys"
	"github.com/donnie4w/tim/util"
)

type admhandle struct{}

var Admhandler = &admhandle{}

func (t *admhandle) ResetAuth(ab *AuthBean) (_r *AdmAck) {
	_r = newAdmAck(false)
	if sys.OsModify(ab.GetUsername(), ab.GetPassword(), ab.Domain) == nil {
		*_r.Ok = true
	}
	return
}

func (t *admhandle) Auth(ab *AuthBean) (_r *AdmAck) {
	_r = newAdmAck(false)
	if auth(ab.GetUsername(), ab.GetPassword()) {
		*_r.Ok = true
	}
	return
}

func (t *admhandle) Token(atoken *AdmToken) (_r *AdmAck) {
	if atoken.GetName() == "" {
		return
	}
	_r = newAdmAck(false)
	if t, n, err := sys.OsToken(atoken.GetName(), atoken.Password, atoken.Domain); err == nil {
		*_r.Ok = true
		*_r.TimType = int8(sys.TIMTOKEN)
		_r.N = &n
		_r.T = &t
	}
	return
}

func (t *admhandle) OsMessage(am *AdmMessage) (_r *AdmAck) {
	_r = newAdmAck(false)
	nodes := am.GetNodes()
	message := admMessageToTimMessage(am.GetMsgbean())
	if err := sys.OsMessage(nodes, message); err == nil {
		*_r.Ok = true
	}
	return
}

func (t *admhandle) ProxyMessage(amb *AdmProxyMessage) (_r *AdmAck) {
	_r = newAdmAck(false)
	if amb.GetConnectid() <= 0 {
		return
	}
	if err := sys.PxMessage(amb.GetConnectid(), admMessageToTimMessage(amb.GetAmb())); err == nil {
		*_r.Ok = true
	}
	return
}

func (t *admhandle) Register(ab *AuthBean) (_r *AdmAck) {
	defer util.Recover()
	_r = newAdmAck(false)
	if node, err := sys.OsRegister(ab.GetUsername(), ab.GetPassword(), ab.Domain); err == nil {
		*_r.Ok = true
		_r.N = &node
	}
	return
}

func (t *admhandle) ModifyUserInfo(amui *AdmModifyUserInfo) (_r *AdmAck) {
	_r = newAdmAck(false)
	if err := sys.OsUserBean(amui.GetNode(), admUserBeanToTimUserBean(amui.GetUserbean())); err == nil {
		*_r.Ok = true
	}
	return
}

func (t *admhandle) BlockUser(abu *AdmBlockUser) (_r *AdmAck) {
	if sys.HasNode(abu.GetAccount()) {
		sys.SendNode(abu.GetAccount(), &TimAck{Ok: true, TimType: int8(sys.TIMLOGOUT)}, sys.TIMACK)
	}
	sys.BlockUser(abu.GetAccount(), abu.GetBlocktime())
	return newAdmAck(true)
}

func (t *admhandle) BlockList() (_r *AdmBlockList) {
	_r = NewAdmBlockList()
	_r.Usermap = sys.BlockList()
	return
}

func (t *admhandle) OnlineUser(au *AdmOnlineUser) (_r *AdmTidList) {
	_r = NewAdmTidList()
	tids, size := sys.WssList(au.GetIndex(), au.GetLimit())
	_r.Size = &size
	adtids := make([]*AdmTid, 0)
	for _, v := range tids {
		adtids = append(adtids, timTidToAdmTid(v))
	}
	_r.Tidlist = adtids
	return
}

func (t *admhandle) Vroom(avb *AdmVroomBean) (_r *AdmAck) {
	_r = newAdmAck(false)
	if avb.GetNode() != "" && avb.GetRtype() > 0 {
		if s := sys.OsVroomprocess(avb.GetNode(), avb.GetRtype()); s != "" {
			*_r.Ok = true
			_r.N = &s
		}
	}
	return
}

func (t *admhandle) TimRoom(arb *AdmRoomReq) (_r *AdmAck) {
	_r = newAdmAck(false)
	if gnode, err := sys.OsRoom(arb.GetNode(), arb.GetTopic(), arb.Domain, arb.GetGtype()); err == nil {
		*_r.Ok = true
		_r.N = &gnode
	}
	return
}

func (t *admhandle) Detect(adb *AdmDetectBean) (_r *AdmAck) {
	sys.Detect(adb.GetNodes())
	return newAdmAck(true)
}

func (t *admhandle) ModifyRoomInfo(arb *AdmRoomBean) (_r *AdmAck) {
	_r = newAdmAck(false)
	if arb.GetUnode() == "" || arb.GetGnode() == "" {
		return
	}
	if err := sys.OsRoomBean(arb.GetUnode(), arb.GetGnode(), admTimRoomToTimRoom(arb.GetAtr())); err == nil {
		*_r.Ok = true
	}
	return
}
