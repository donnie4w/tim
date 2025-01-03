// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of t source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package adm

import (
	"github.com/donnie4w/tim/branch"
	"github.com/donnie4w/tim/errs"
	. "github.com/donnie4w/tim/stub"
	"github.com/donnie4w/tim/sys"
	"github.com/donnie4w/tim/util"
)

type admhandle struct {
	domain string
}

var Admhandler = new(admhandle)

func (ah *admhandle) ModifyPwd(fromnode, oldpwd, newpwd, domain string) (_r *AdmAck) {
	_r = newAdmAck(false)
	var old *string
	if oldpwd != "" {
		old = &oldpwd
	}
	var dom *string
	if domain != "" {
		dom = &domain
	}
	if err := sys.OsModify(fromnode, old, newpwd, dom); err == nil {
		*_r.Ok = true
	} else {
		_r.Errcode = err.TimError().Code
	}
	return
}

func (ah *admhandle) Auth(ab *AuthBean) (_r *AdmAck) {
	_r = newAdmAck(false)
	if auth(ab.GetUsername(), ab.GetPassword(), ab.GetDomain()) {
		*_r.Ok = true
		ah.domain = ab.GetDomain()
	} else {
		_r.Errcode = errs.ERR_NOPASS.TimError().Code
	}
	return
}

func (ah *admhandle) Token(atoken *AdmToken) (_r *AdmAck) {
	_r = newAdmAck(false)
	if atoken.GetName() == "" {
		_r.Errcode = errs.ERR_ACCOUNT.TimError().Code
		return
	}
	if tkn, n, err := sys.OsToken(atoken.GetName(), atoken.Password, atoken.Domain); err == nil {
		*_r.Ok = true
		tt := int8(sys.TIMTOKEN)
		_r.TimType = &tt
		_r.N = &n
		_r.T = &tkn
	} else {
		_r.Errcode = err.TimError().Code
	}
	return
}

func (ah *admhandle) TimMessageBroadcast(am *AdmMessageBroadcast) (_r *AdmAck) {
	_r = newAdmAck(false)
	nodes := am.GetNodes()
	message := admMessageToTimMessage(am.GetMessage())
	if err := sys.OsMessage(nodes, message); err == nil {
		*_r.Ok = true
	} else {
		_r.Errcode = err.TimError().Code
	}
	return
}

func (ah *admhandle) TimPresenceBroadcast(apb *AdmPresenceBroadcast) (_r *AdmAck) {
	_r = newAdmAck(false)
	nodes := apb.GetNodes()
	presence := admPresenceToTimPresence(apb.GetPresence())
	if err := sys.OsPresence(nodes, presence); err == nil {
		*_r.Ok = true
	} else {
		_r.Errcode = err.TimError().Code
	}
	return
}

func (ah *admhandle) ProxyMessage(amb *AdmProxyMessage) (_r *AdmAck) {
	_r = newAdmAck(false)
	if amb.GetConnectid() <= 0 {
		_r.Errcode = errs.ERR_ACCOUNT.TimError().Code
		return
	}
	if err := sys.PxMessage(amb.GetConnectid(), admMessageToTimMessage(amb.GetMessage())); err == nil {
		*_r.Ok = true
	} else {
		_r.Errcode = err.TimError().Code
	}
	return
}

func (ah *admhandle) Register(ab *AuthBean) (_r *AdmAck) {
	defer util.Recover()
	_r = newAdmAck(false)
	if node, err := sys.OsRegister(ab.GetUsername(), ab.GetPassword(), ab.Domain); err == nil {
		*_r.Ok = true
		_r.N = &node
	} else {
		_r.Errcode = err.TimError().Code
	}
	return
}

func (ah *admhandle) ModifyUserInfo(amui *AdmModifyUserInfo) (_r *AdmAck) {
	_r = newAdmAck(false)
	if err := sys.OsUserBean(amui.GetNode(), admUserBeanToTimUserBean(amui.GetUserbean())); err == nil {
		*_r.Ok = true
	} else {
		_r.Errcode = err.TimError().Code
	}
	return
}

func (ah *admhandle) ModifyRoomInfo(arb *AdmRoomBean) (_r *AdmAck) {
	_r = newAdmAck(false)
	if arb.GetUnode() == "" || arb.GetGnode() == "" {
		_r.Errcode = errs.ERR_ACCOUNT.TimError().Code
		return
	}
	if err := sys.OsRoomBean(arb.GetUnode(), arb.GetGnode(), admTimRoomToTimRoom(arb.GetAtr())); err == nil {
		*_r.Ok = true
	} else {
		_r.Errcode = err.TimError().Code
	}
	return
}

func (ah *admhandle) SysBlockUser(abu *AdmSysBlockUser) (_r *AdmAck) {
	if sys.HasNode(abu.GetAccount()) {
		sys.SendNode(abu.GetAccount(), &TimAck{Ok: true, TimType: int8(sys.TIMLOGOUT)}, sys.TIMACK)
	}
	sys.OsBlockUser(abu.GetAccount(), abu.GetBlocktime())
	return newAdmAck(true)
}

func (ah *admhandle) SysBlockList() (_r *AdmSysBlockList) {
	_r = NewAdmSysBlockList()
	_r.Usermap = sys.OsBlockList()
	return
}

func (ah *admhandle) OnlineUser(au *AdmOnlineUser) (_r *AdmTidList) {
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

func (ah *admhandle) Vroom(avb *AdmVroomBean) (_r *AdmAck) {
	_r = newAdmAck(false)
	if avb.GetNode() != "" && avb.GetRtype() > 0 {
		if s := sys.OsVroomprocess(avb.GetNode(), avb.GetRtype()); s != "" {
			*_r.Ok = true
			_r.N = &s
		}
	}
	return
}

func (ah *admhandle) Detect(adb *AdmDetectBean) (_r *AdmAck) {
	sys.Detect(adb.GetNodes())
	return newAdmAck(true)
}

func (ah *admhandle) Roster(fromnode, domain string) (_r *AdmNodeList) {
	_r = NewAdmNodeList()
	_r.Nodelist, _ = branch.Roster(fromnode, &domain)
	return
}

func (ah *admhandle) Addroster(fromnode, domain string, tonode string, msg string) (_r *AdmAck) {
	_r = newAdmAck(false)
	if err := branch.Addroster(fromnode, &domain, tonode, &msg); err == nil {
		*_r.Ok = true
	}
	return
}

func (ah *admhandle) Rmroster(fromnode, domain string, tonode string) (_r *AdmAck) {
	_r = newAdmAck(false)
	if err := branch.Rmroster(fromnode, &domain, tonode); err == nil {
		*_r.Ok = true
	}
	return
}

func (ah *admhandle) Blockroster(fromnode, domain string, tonode string) (_r *AdmAck) {
	_r = newAdmAck(false)
	if err := branch.Blockroster(fromnode, &domain, tonode); err == nil {
		*_r.Ok = true
	}
	return
}

func (ah *admhandle) PullUserMessage(fromnode, domain string, tonode string, mid int64, limit int64) (_r *AdmMessageList) {
	if tms := branch.PullUserMessage(fromnode, &domain, tonode, mid, limit); len(tms) > 0 {
		size := int64(len(tms))
		_r = NewAdmMessageList()
		_r.Totalcount = &size
		_r.Msglist = make([]*AdmMessage, size)
		for i, v := range tms {
			_r.Msglist[i] = timMessageToAdmMessage(v)
		}
	}
	return
}

func (ah *admhandle) PullRoomMessage(fromnode, domain string, tonode string, mid int64, limit int64) (_r *AdmMessageList) {
	if tms := branch.PullRoomMessage(fromnode, &domain, tonode, mid, limit); len(tms) > 0 {
		size := int64(len(tms))
		_r = NewAdmMessageList()
		_r.Totalcount = &size
		_r.Msglist = make([]*AdmMessage, size)
		for i, v := range tms {
			_r.Msglist[i] = timMessageToAdmMessage(v)
		}
	}
	return
}

func (ah *admhandle) OfflineMsg(fromnode, domain string, limit int64) (_r *AdmMessageList) {
	if tms := branch.OfflineMsg(fromnode, &domain, int(limit)); len(tms) > 0 {
		size := int64(len(tms))
		_r = NewAdmMessageList()
		_r.Totalcount = &size
		_r.Msglist = make([]*AdmMessage, size)
		for i, v := range tms {
			_r.Msglist[i] = timMessageToAdmMessage(v)
		}
	}
	return
}

func (ah *admhandle) DelOfflineMsg(fromnode string, domain string, ids []int64) (_r *AdmAck) {
	_r = newAdmAck(false)
	if _, err := branch.DelOfflineMsg(fromnode, ids); err == nil {
		*_r.Ok = true
	} else {
		_r.Errcode = errs.ERR_PARAMS.TimError().Code
	}
	return
}

func (ah *admhandle) UserRoom(fromnode, domain string) (_r *AdmNodeList) {
	if ss := branch.UserRoom(fromnode, &domain); len(ss) > 0 {
		_r = NewAdmNodeList()
		_r.Nodelist = ss
	}
	return
}

func (ah *admhandle) RoomUsers(fromnode, domain string) (_r *AdmNodeList) {
	if ss := branch.RoomUsers(fromnode, &domain); len(ss) > 0 {
		_r = NewAdmNodeList()
		_r.Nodelist = ss
	}
	return
}

func (ah *admhandle) CreateRoom(fromnode, domain string, topic string, gtype int8) (_r *AdmAck) {
	_r = newAdmAck(false)
	if !util.CheckNode(fromnode) {
		_r.Errcode = errs.ERR_PARAMS.TimError().Code
		return
	}
	if gnode, err := sys.OsRoom(fromnode, topic, &domain, gtype); err == nil {
		*_r.Ok = true
		_r.N = &gnode
	} else {
		_r.Errcode = err.TimError().Code
	}
	return
}

func (ah *admhandle) AddRoom(fromnode, domain string, roomNode string, msg string) (_r *AdmAck) {
	_r = newAdmAck(false)

	if !util.CheckNodes(fromnode, roomNode) {
		_r.Errcode = errs.ERR_PARAMS.TimError().Code
		return
	}

	if err := branch.AddRoom(fromnode, &domain, roomNode, msg); err == nil {
		*_r.Ok = true
	} else {
		_r.Errcode = err.TimError().Code
	}
	return
}

func (ah *admhandle) PullInRoom(fromnode, domain string, roomNode string, toNode string) (_r *AdmAck) {
	_r = newAdmAck(false)

	if !util.CheckNodes(fromnode, roomNode, toNode) {
		_r.Errcode = errs.ERR_PARAMS.TimError().Code
		return
	}

	if err := branch.PullInRoom(fromnode, &domain, roomNode, toNode); err == nil {
		*_r.Ok = true
	} else {
		_r.Errcode = err.TimError().Code
	}
	return
}

func (ah *admhandle) RejectRoom(fromnode, domain string, roomNode string, toNode string, msg string) (_r *AdmAck) {
	_r = newAdmAck(false)

	if !util.CheckNodes(fromnode, roomNode, toNode) {
		_r.Errcode = errs.ERR_PARAMS.TimError().Code
		return
	}

	if err := branch.RejectRoom(fromnode, &domain, roomNode, toNode, msg); err == nil {
		*_r.Ok = true
	} else {
		_r.Errcode = err.TimError().Code
	}
	return
}

func (ah *admhandle) KickRoom(fromnode, domain string, roomNode string, toNode string) (_r *AdmAck) {
	_r = newAdmAck(false)

	if !util.CheckNodes(fromnode, roomNode, toNode) {
		_r.Errcode = errs.ERR_PARAMS.TimError().Code
		return
	}

	if err := branch.KickRoom(fromnode, &domain, roomNode, toNode); err == nil {
		*_r.Ok = true
	} else {
		_r.Errcode = err.TimError().Code
	}
	return
}

func (ah *admhandle) LeaveRoom(fromnode, domain string, roomNode string) (_r *AdmAck) {
	_r = newAdmAck(false)

	if !util.CheckNodes(fromnode, roomNode) {
		_r.Errcode = errs.ERR_PARAMS.TimError().Code
		return
	}

	if err := branch.LeaveRoom(fromnode, &domain, roomNode); err == nil {
		*_r.Ok = true
	} else {
		_r.Errcode = err.TimError().Code
	}
	return
}

func (ah *admhandle) CancelRoom(fromnode, domain string, roomNode string) (_r *AdmAck) {
	_r = newAdmAck(false)

	if !util.CheckNodes(fromnode, roomNode) {
		_r.Errcode = errs.ERR_PARAMS.TimError().Code
		return
	}

	if err := branch.CancelRoom(fromnode, &domain, roomNode); err == nil {
		*_r.Ok = true
	} else {
		_r.Errcode = err.TimError().Code
	}
	return
}

func (ah *admhandle) BlockRoom(fromnode, domain string, roomNode string) (_r *AdmAck) {
	_r = newAdmAck(false)
	if !util.CheckNodes(fromnode, roomNode) {
		_r.Errcode = errs.ERR_PARAMS.TimError().Code
		return
	}
	if err := branch.BlockRoom(fromnode, &domain, roomNode); err == nil {
		*_r.Ok = true
	} else {
		_r.Errcode = err.TimError().Code
	}
	return
}

func (ah *admhandle) BlockRoomMember(fromnode, domain string, roomNode string, toNode string) (_r *AdmAck) {
	_r = newAdmAck(false)

	if !util.CheckNodes(fromnode, roomNode, toNode) {
		_r.Errcode = errs.ERR_PARAMS.TimError().Code
		return
	}

	if err := branch.BlockRoomMember(fromnode, &domain, roomNode, toNode); err == nil {
		*_r.Ok = true
	} else {
		_r.Errcode = err.TimError().Code
	}
	return
}

func (ah *admhandle) BlockRosterList(fromnode, domain string) (_r *AdmNodeList) {
	if ss, _ := branch.BlockRosterList(fromnode, &domain); len(ss) > 0 {
		_r = NewAdmNodeList()
		_r.Nodelist = ss
	}
	return
}

func (ah *admhandle) BlockRoomList(fromnode, domain string) (_r *AdmNodeList) {
	if ss, _ := branch.BlockRoomList(fromnode, &domain); len(ss) > 0 {
		_r = NewAdmNodeList()
		_r.Nodelist = ss
	}
	return
}

func (ah *admhandle) BlockRoomMemberlist(fromnode, domain string, vNode string) (_r *AdmNodeList) {
	if ss, _ := branch.BlockRoomMemberlist(fromnode, &domain, vNode); len(ss) > 0 {
		_r = NewAdmNodeList()
		_r.Nodelist = ss
	}
	return
}

func (ah *admhandle) VirtualroomRegister(fromnode, domain string) (_r *AdmAck) {
	if ack, err := branch.VirtualroomRegister(fromnode, domain); err == nil {
		_r = timAckToAdmAck(ack)
	} else {
		_r = newAdmAck(false)
		_r.Errcode = err.TimError().Code
	}
	return
}

func (ah *admhandle) VirtualroomRemove(fromnode, domain string, vNode string) (_r *AdmAck) {
	if ack, err := branch.VirtualroomRemove(fromnode, domain, vNode); err == nil {
		_r = timAckToAdmAck(ack)
	} else {
		_r = newAdmAck(false)
		_r.Errcode = err.TimError().Code
	}
	return
}

func (ah *admhandle) VirtualroomAddAuth(fromnode, domain string, vNode string, toNode string) (_r *AdmAck) {
	if ack, err := branch.VirtualroomAddAuth(fromnode, domain, vNode, toNode); err == nil {
		_r = timAckToAdmAck(ack)
	} else {
		_r = newAdmAck(false)
		_r.Errcode = err.TimError().Code
	}
	return
}

func (ah *admhandle) VirtualroomDelAuth(fromnode, domain string, vNode string, toNode string) (_r *AdmAck) {
	if ack, err := branch.VirtualroomDelAuth(fromnode, domain, vNode, toNode); err == nil {
		_r = timAckToAdmAck(ack)
	} else {
		_r = newAdmAck(false)
		_r.Errcode = err.TimError().Code
	}
	return
}

func (ah *admhandle) VirtualroomSub(wsid int64, fromnode string, domain string, vNode string, subType int8) (_r *AdmAck) {
	if ack, err := branch.VirtualroomSub(wsid, fromnode, domain, vNode, subType); err == nil {
		_r = timAckToAdmAck(ack)
	} else {
		_r = newAdmAck(false)
		_r.Errcode = err.TimError().Code
	}
	return
}

func (ah *admhandle) VirtualroomUnSub(wsid int64, fromnode string, domain string, vNode string) (_r *AdmAck) {
	if ack, err := branch.VirtualroomUnSub(wsid, fromnode, domain, vNode); err == nil {
		_r = timAckToAdmAck(ack)
	} else {
		_r = newAdmAck(false)
		_r.Errcode = err.TimError().Code
	}
	return
}
