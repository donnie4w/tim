// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of t source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package adm

import (
	"context"
	"github.com/donnie4w/tim/errs"
	. "github.com/donnie4w/tim/stub"
	"github.com/donnie4w/tim/util"
)

type ifacehandle struct{}

var ifaceProcessor = &ifacehandle{}

func (t *ifacehandle) ModifyPwd(ctx context.Context, fromnode, oldpwd, newpwd, domain string) (_r *AdmAck, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if e := noAuthAndClose(cc); e != nil {
		_r = newAdmAck(false)
		_r.Errcode = e.TimError().Code
		return
	}
	if fromnode == "" || newpwd == "" {
		_r = newAdmAck(false)
		_r.Errcode = errs.ERR_PARAMS.TimError().Code
		return
	}
	_r = Admhandler.ModifyPwd(fromnode, oldpwd, newpwd, domain)
	return
}

func (t *ifacehandle) Auth(ctx context.Context, ab *AuthBean) (_r *AdmAck, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if _r = Admhandler.Auth(ab); _r.GetOk() {
		ctx2CliContext(ctx).isAuth = true
	}
	return
}

func (t *ifacehandle) Ping(ctx context.Context) (_r *AdmAck, _err error) {
	defer util.Recover2(&_err)
	mux := ctx2CliContext(ctx).mux
	defer mux.Unlock()
	mux.Lock()
	_r = newAdmAck(true)
	return
}

func (t *ifacehandle) Token(ctx context.Context, atoken *AdmToken) (_r *AdmAck, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if e := noAuthAndClose(cc); e != nil {
		_r = newAdmAck(false)
		_r.Errcode = e.TimError().Code
		return
	}
	_r = Admhandler.Token(atoken)
	return
}

func (t *ifacehandle) TimMessageBroadcast(ctx context.Context, amb *AdmMessageBroadcast) (_r *AdmAck, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if e := noAuthAndClose(cc); e != nil {
		_r = newAdmAck(false)
		_r.Errcode = e.TimError().Code
		return
	}
	_r = Admhandler.TimMessageBroadcast(amb)
	return
}

func (t *ifacehandle) TimPresenceBroadcast(ctx context.Context, apb *AdmPresenceBroadcast) (_r *AdmAck, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if e := noAuthAndClose(cc); e != nil {
		_r = newAdmAck(false)
		_r.Errcode = e.TimError().Code
		return
	}
	_r = Admhandler.TimPresenceBroadcast(apb)
	return
}

func (t *ifacehandle) ProxyMessage(ctx context.Context, amb *AdmProxyMessage) (_r *AdmAck, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if e := noAuthAndClose(cc); e != nil {
		_r = newAdmAck(false)
		_r.Errcode = e.TimError().Code
		return
	}
	_r = Admhandler.ProxyMessage(amb)
	return
}

func (t *ifacehandle) Register(ctx context.Context, ab *AuthBean) (_r *AdmAck, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if e := noAuthAndClose(cc); e != nil {
		_r = newAdmAck(false)
		_r.Errcode = e.TimError().Code
		return
	}
	_r = Admhandler.Register(ab)
	return
}

func (t *ifacehandle) ModifyUserInfo(ctx context.Context, amui *AdmModifyUserInfo) (_r *AdmAck, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if e := noAuthAndClose(cc); e != nil {
		_r = newAdmAck(false)
		_r.Errcode = e.TimError().Code
		return
	}
	_r = Admhandler.ModifyUserInfo(amui)
	return
}

func (t *ifacehandle) SysBlockUser(ctx context.Context, abu *AdmSysBlockUser) (_r *AdmAck, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if e := noAuthAndClose(cc); e != nil {
		_r = newAdmAck(false)
		_r.Errcode = e.TimError().Code
		return
	}
	_r = Admhandler.SysBlockUser(abu)
	return
}

func (t *ifacehandle) SysBlockList(ctx context.Context) (_r *AdmSysBlockList, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if noAuthAndClose(cc) != nil {
		return
	}
	_r = Admhandler.SysBlockList()
	return
}

func (t *ifacehandle) OnlineUser(ctx context.Context, au *AdmOnlineUser) (_r *AdmTidList, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if noAuthAndClose(cc) != nil {
		return
	}
	_r = Admhandler.OnlineUser(au)
	return
}

func (t *ifacehandle) Vroom(ctx context.Context, avb *AdmVroomBean) (_r *AdmAck, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if e := noAuthAndClose(cc); e != nil {
		_r = newAdmAck(false)
		_r.Errcode = e.TimError().Code
		return
	}
	_r = Admhandler.Vroom(avb)
	return
}

func (t *ifacehandle) Detect(ctx context.Context, adb *AdmDetectBean) (_r *AdmAck, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if e := noAuthAndClose(cc); e != nil {
		_r = newAdmAck(false)
		_r.Errcode = e.TimError().Code
		return
	}
	_r = Admhandler.Detect(adb)
	return
}

func (t *ifacehandle) ModifyRoomInfo(ctx context.Context, arb *AdmRoomBean) (_r *AdmAck, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if e := noAuthAndClose(cc); e != nil {
		_r = newAdmAck(false)
		_r.Errcode = e.TimError().Code
		return
	}
	_r = Admhandler.ModifyRoomInfo(arb)
	return
}

func (t *ifacehandle) Roster(ctx context.Context, fromnode, domain string) (_r *AdmNodeList, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if e := noAuthAndClose(cc); e == nil {
		_r = Admhandler.Roster(fromnode, domain)
	}
	return
}

func (t *ifacehandle) Addroster(ctx context.Context, fromnode, domain string, tonode string, msg string) (_r *AdmAck, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if e := noAuthAndClose(cc); e != nil {
		_r = newAdmAck(false)
		_r.Errcode = e.TimError().Code
		return
	}
	_r = Admhandler.Addroster(fromnode, domain, tonode, msg)
	return
}

func (t *ifacehandle) Rmroster(ctx context.Context, fromnode, domain string, tonode string) (_r *AdmAck, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if e := noAuthAndClose(cc); e != nil {
		_r = newAdmAck(false)
		_r.Errcode = e.TimError().Code
		return
	}
	_r = Admhandler.Rmroster(fromnode, domain, tonode)
	return
}

func (t *ifacehandle) Blockroster(ctx context.Context, fromnode, domain string, tonode string) (_r *AdmAck, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if e := noAuthAndClose(cc); e != nil {
		_r = newAdmAck(false)
		_r.Errcode = e.TimError().Code
		return
	}
	_r = Admhandler.Blockroster(fromnode, domain, tonode)
	return
}

func (t *ifacehandle) PullUserMessage(ctx context.Context, fromnode, domain string, tonode string, mid int64, limit int64) (_r *AdmMessageList, _err error) {
	return
}

func (t *ifacehandle) PullRoomMessage(ctx context.Context, fromnode, domain string, tonode string, mid int64, limit int64) (_r *AdmMessageList, _err error) {
	return
}

func (t *ifacehandle) OfflineMsg(ctx context.Context, fromnode, domain string, limit int64) (_r *AdmMessageList, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if e := noAuthAndClose(cc); e != nil {
		return
	}
	_r = Admhandler.OfflineMsg(fromnode, domain, limit)
	return
}

func (ah *ifacehandle) DelOfflineMsg(ctx context.Context, fromnode string, domain string, ids []int64) (_r *AdmAck, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if e := noAuthAndClose(cc); e != nil {
		return
	}
	_r = Admhandler.DelOfflineMsg(fromnode, domain, ids)
	return
}

func (t *ifacehandle) UserRoom(ctx context.Context, fromnode, domain string) (_r *AdmNodeList, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if e := noAuthAndClose(cc); e != nil {
		return
	}
	_r = Admhandler.UserRoom(fromnode, domain)
	return
}

func (t *ifacehandle) RoomUsers(ctx context.Context, fromnode, domain string) (_r *AdmNodeList, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if e := noAuthAndClose(cc); e != nil {
		return
	}
	_r = Admhandler.RoomUsers(fromnode, domain)
	return
}

func (t *ifacehandle) CreateRoom(ctx context.Context, fromnode, domain string, topic string, gtype int8) (_r *AdmAck, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if e := noAuthAndClose(cc); e != nil {
		_r = newAdmAck(false)
		_r.Errcode = e.TimError().Code
		return
	}
	_r = Admhandler.CreateRoom(fromnode, domain, topic, gtype)
	return
}

func (t *ifacehandle) AddRoom(ctx context.Context, fromnode, domain string, roomNode string, msg string) (_r *AdmAck, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if e := noAuthAndClose(cc); e != nil {
		_r = newAdmAck(false)
		_r.Errcode = e.TimError().Code
		return
	}
	_r = Admhandler.AddRoom(fromnode, domain, roomNode, msg)
	return
}

func (t *ifacehandle) PullInRoom(ctx context.Context, fromnode, domain string, roomNode string, toNode string) (_r *AdmAck, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if e := noAuthAndClose(cc); e != nil {
		_r = newAdmAck(false)
		_r.Errcode = e.TimError().Code
		return
	}
	_r = Admhandler.PullInRoom(fromnode, domain, roomNode, toNode)
	return
}

func (t *ifacehandle) RejectRoom(ctx context.Context, fromnode, domain string, roomNode string, toNode string, msg string) (_r *AdmAck, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if e := noAuthAndClose(cc); e != nil {
		_r = newAdmAck(false)
		_r.Errcode = e.TimError().Code
		return
	}
	_r = Admhandler.RejectRoom(fromnode, domain, roomNode, toNode, msg)
	return
}

func (t *ifacehandle) KickRoom(ctx context.Context, fromnode, domain string, roomNode string, toNode string) (_r *AdmAck, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if e := noAuthAndClose(cc); e != nil {
		_r = newAdmAck(false)
		_r.Errcode = e.TimError().Code
		return
	}
	_r = Admhandler.KickRoom(fromnode, domain, roomNode, toNode)
	return
}

func (t *ifacehandle) LeaveRoom(ctx context.Context, fromnode, domain string, roomNode string) (_r *AdmAck, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if e := noAuthAndClose(cc); e != nil {
		_r = newAdmAck(false)
		_r.Errcode = e.TimError().Code
		return
	}
	_r = Admhandler.LeaveRoom(fromnode, domain, roomNode)
	return
}

func (t *ifacehandle) CancelRoom(ctx context.Context, fromnode, domain string, roomNode string) (_r *AdmAck, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if e := noAuthAndClose(cc); e != nil {
		_r = newAdmAck(false)
		_r.Errcode = e.TimError().Code
		return
	}
	_r = Admhandler.CancelRoom(fromnode, domain, roomNode)
	return
}

func (t *ifacehandle) BlockRoom(ctx context.Context, fromnode, domain string, roomNode string) (_r *AdmAck, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if e := noAuthAndClose(cc); e != nil {
		_r = newAdmAck(false)
		_r.Errcode = e.TimError().Code
		return
	}
	_r = Admhandler.BlockRoom(fromnode, domain, roomNode)
	return
}

func (t *ifacehandle) BlockRoomMember(ctx context.Context, fromnode, domain string, roomNode string, toNode string) (_r *AdmAck, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if e := noAuthAndClose(cc); e != nil {
		_r = newAdmAck(false)
		_r.Errcode = e.TimError().Code
		return
	}
	_r = Admhandler.BlockRoomMember(fromnode, domain, roomNode, toNode)
	return
}

func (t *ifacehandle) BlockRosterList(ctx context.Context, fromnode, domain string) (_r *AdmNodeList, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if e := noAuthAndClose(cc); e == nil {
		_r = Admhandler.BlockRosterList(fromnode, domain)
	}
	return
}

func (t *ifacehandle) BlockRoomList(ctx context.Context, fromnode, domain string) (_r *AdmNodeList, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if e := noAuthAndClose(cc); e == nil {
		_r = Admhandler.BlockRoomList(fromnode, domain)
	}
	return
}

func (t *ifacehandle) BlockRoomMemberlist(ctx context.Context, fromnode, domain string, roomNode string) (_r *AdmNodeList, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if e := noAuthAndClose(cc); e == nil {
		_r = Admhandler.BlockRoomMemberlist(fromnode, domain, roomNode)
	}
	return
}

func (t *ifacehandle) VirtualroomRegister(ctx context.Context, fromnode, domain string) (_r *AdmAck, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if e := noAuthAndClose(cc); e != nil {
		_r = newAdmAck(false)
		_r.Errcode = e.TimError().Code
		return
	}
	_r = Admhandler.VirtualroomRegister(fromnode, domain)
	return
}

func (t *ifacehandle) VirtualroomRemove(ctx context.Context, fromnode, domain string, roomNode string) (_r *AdmAck, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if e := noAuthAndClose(cc); e != nil {
		_r = newAdmAck(false)
		_r.Errcode = e.TimError().Code
		return
	}
	_r = Admhandler.VirtualroomRemove(fromnode, domain, roomNode)
	return
}

func (t *ifacehandle) VirtualroomAddAuth(ctx context.Context, fromnode, domain string, roomNode string, toNode string) (_r *AdmAck, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if e := noAuthAndClose(cc); e != nil {
		_r = newAdmAck(false)
		_r.Errcode = e.TimError().Code
		return
	}
	_r = Admhandler.VirtualroomAddAuth(fromnode, domain, roomNode, toNode)
	return
}

func (t *ifacehandle) VirtualroomDelAuth(ctx context.Context, fromnode, domain string, roomNode string, toNode string) (_r *AdmAck, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if e := noAuthAndClose(cc); e != nil {
		_r = newAdmAck(false)
		_r.Errcode = e.TimError().Code
		return
	}
	_r = Admhandler.VirtualroomDelAuth(fromnode, domain, roomNode, toNode)
	return
}

func (t *ifacehandle) VirtualroomSub(ctx context.Context, wsid int64, fromnode string, domain string, vNode string, subType int8) (_r *AdmAck, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if e := noAuthAndClose(cc); e != nil {
		_r = newAdmAck(false)
		_r.Errcode = e.TimError().Code
		return
	}
	_r = Admhandler.VirtualroomSub(wsid, fromnode, domain, vNode, subType)
	return
}

func (t *ifacehandle) VirtualroomUnSub(ctx context.Context, wsid int64, fromnode string, domain string, vNode string) (_r *AdmAck, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if e := noAuthAndClose(cc); e != nil {
		_r = newAdmAck(false)
		_r.Errcode = e.TimError().Code
		return
	}
	_r = Admhandler.VirtualroomUnSub(wsid, fromnode, domain, vNode)
	return
}
