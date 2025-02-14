// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of t source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package adm

import (
	"context"
	goutil "github.com/donnie4w/gofer/util"
	"github.com/donnie4w/tim/errs"
	"github.com/donnie4w/tim/stub"
	"github.com/donnie4w/tim/sys"
	"github.com/donnie4w/tim/util"
)

type admprocessor struct{}

var admProcessor = newAdmProcessor()

func newAdmProcessor() stub.Admiface {
	return &admprocessor{}
}

func (t *admprocessor) ModifyPwd(ctx context.Context, fromnode, oldpwd, newpwd, domain string) (_r *stub.AdmAck, _err error) {
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
	_r = admHandler.ModifyPwd(fromnode, oldpwd, newpwd, domain)
	return
}

func (t *admprocessor) Auth(ctx context.Context, ab *stub.AuthBean) (_r *stub.AdmAck, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if _r = admHandler.Auth(ab); _r.GetOk() {
		maskuuid := int64(int32(goutil.FNVHash32(goutil.Int64ToBytes(sys.UUID))))
		_r.T = &maskuuid
		ctx2CliContext(ctx).isAuth = true
	}
	return
}

func (t *admprocessor) Ping(ctx context.Context) (_r *stub.AdmAck, _err error) {
	defer util.Recover2(&_err)
	mux := ctx2CliContext(ctx).mux
	defer mux.Unlock()
	mux.Lock()
	_r = newAdmAck(true)
	return
}

func (t *admprocessor) Token(ctx context.Context, atoken *stub.AdmToken) (_r *stub.AdmAck, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if e := noAuthAndClose(cc); e != nil {
		_r = newAdmAck(false)
		_r.Errcode = e.TimError().Code
		return
	}
	_r = admHandler.Token(atoken)
	return
}

func (t *admprocessor) TimMessageBroadcast(ctx context.Context, amb *stub.AdmMessageBroadcast) (_r *stub.AdmAck, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if e := noAuthAndClose(cc); e != nil {
		_r = newAdmAck(false)
		_r.Errcode = e.TimError().Code
		return
	}
	_r = admHandler.TimMessageBroadcast(amb)
	return
}

func (t *admprocessor) TimPresenceBroadcast(ctx context.Context, apb *stub.AdmPresenceBroadcast) (_r *stub.AdmAck, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if e := noAuthAndClose(cc); e != nil {
		_r = newAdmAck(false)
		_r.Errcode = e.TimError().Code
		return
	}
	_r = admHandler.TimPresenceBroadcast(apb)
	return
}

func (t *admprocessor) ProxyMessage(ctx context.Context, amb *stub.AdmProxyMessage) (_r *stub.AdmAck, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if e := noAuthAndClose(cc); e != nil {
		_r = newAdmAck(false)
		_r.Errcode = e.TimError().Code
		return
	}
	_r = admHandler.ProxyMessage(amb)
	return
}

func (t *admprocessor) Register(ctx context.Context, ab *stub.AuthBean) (_r *stub.AdmAck, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if e := noAuthAndClose(cc); e != nil {
		_r = newAdmAck(false)
		_r.Errcode = e.TimError().Code
		return
	}
	_r = admHandler.Register(ab)
	return
}

func (t *admprocessor) ModifyUserInfo(ctx context.Context, amui *stub.AdmModifyUserInfo) (_r *stub.AdmAck, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if e := noAuthAndClose(cc); e != nil {
		_r = newAdmAck(false)
		_r.Errcode = e.TimError().Code
		return
	}
	_r = admHandler.ModifyUserInfo(amui)
	return
}

func (t *admprocessor) SysBlockUser(ctx context.Context, abu *stub.AdmSysBlockUser) (_r *stub.AdmAck, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if e := noAuthAndClose(cc); e != nil {
		_r = newAdmAck(false)
		_r.Errcode = e.TimError().Code
		return
	}
	_r = admHandler.SysBlockUser(abu)
	return
}

func (t *admprocessor) OnlineUser(ctx context.Context, au *stub.AdmOnlineUser) (_r *stub.AdmTidList, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if noAuthAndClose(cc) != nil {
		return
	}
	_r = admHandler.OnlineUser(au)
	return
}

func (t *admprocessor) Vroom(ctx context.Context, avb *stub.AdmVroomBean) (_r *stub.AdmAck, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if e := noAuthAndClose(cc); e != nil {
		_r = newAdmAck(false)
		_r.Errcode = e.TimError().Code
		return
	}
	_r = admHandler.Vroom(avb)
	return
}

func (t *admprocessor) Detect(ctx context.Context, adb *stub.AdmDetectBean) (_r *stub.AdmAck, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if e := noAuthAndClose(cc); e != nil {
		_r = newAdmAck(false)
		_r.Errcode = e.TimError().Code
		return
	}
	_r = admHandler.Detect(adb)
	return
}

func (t *admprocessor) ModifyRoomInfo(ctx context.Context, arb *stub.AdmRoomBean) (_r *stub.AdmAck, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if e := noAuthAndClose(cc); e != nil {
		_r = newAdmAck(false)
		_r.Errcode = e.TimError().Code
		return
	}
	_r = admHandler.ModifyRoomInfo(arb)
	return
}

func (t *admprocessor) Roster(ctx context.Context, fromnode, domain string) (_r *stub.AdmNodeList, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if e := noAuthAndClose(cc); e == nil {
		_r = admHandler.Roster(fromnode, domain)
	}
	return
}

func (t *admprocessor) Addroster(ctx context.Context, fromnode, domain string, tonode string, msg string) (_r *stub.AdmAck, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if e := noAuthAndClose(cc); e != nil {
		_r = newAdmAck(false)
		_r.Errcode = e.TimError().Code
		return
	}
	_r = admHandler.Addroster(fromnode, domain, tonode, msg)
	return
}

func (t *admprocessor) Rmroster(ctx context.Context, fromnode, domain string, tonode string) (_r *stub.AdmAck, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if e := noAuthAndClose(cc); e != nil {
		_r = newAdmAck(false)
		_r.Errcode = e.TimError().Code
		return
	}
	_r = admHandler.Rmroster(fromnode, domain, tonode)
	return
}

func (t *admprocessor) Blockroster(ctx context.Context, fromnode, domain string, tonode string) (_r *stub.AdmAck, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if e := noAuthAndClose(cc); e != nil {
		_r = newAdmAck(false)
		_r.Errcode = e.TimError().Code
		return
	}
	_r = admHandler.Blockroster(fromnode, domain, tonode)
	return
}

func (t *admprocessor) PullUserMessage(ctx context.Context, fromnode, domain string, tonode string, mid, timeseries, limit int64) (_r *stub.AdmMessageList, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if e := noAuthAndClose(cc); e != nil {
		return
	}
	_r = admHandler.PullUserMessage(fromnode, domain, tonode, mid, timeseries, limit)
	return
}

func (t *admprocessor) PullRoomMessage(ctx context.Context, fromnode, domain string, tonode string, mid, timeseries, limit int64) (_r *stub.AdmMessageList, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if e := noAuthAndClose(cc); e != nil {
		return
	}
	return
}

func (t *admprocessor) OfflineMsg(ctx context.Context, fromnode, domain string, limit int64) (_r *stub.AdmMessageList, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if e := noAuthAndClose(cc); e != nil {
		return
	}
	_r = admHandler.OfflineMsg(fromnode, domain, limit)
	return
}

func (ah *admprocessor) DelOfflineMsg(ctx context.Context, fromnode string, domain string, ids []int64) (_r *stub.AdmAck, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if e := noAuthAndClose(cc); e != nil {
		return
	}
	_r = admHandler.DelOfflineMsg(fromnode, domain, ids)
	return
}

func (t *admprocessor) UserRoom(ctx context.Context, fromnode, domain string) (_r *stub.AdmNodeList, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if e := noAuthAndClose(cc); e != nil {
		return
	}
	_r = admHandler.UserRoom(fromnode, domain)
	return
}

func (t *admprocessor) RoomUsers(ctx context.Context, fromnode, domain string) (_r *stub.AdmNodeList, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if e := noAuthAndClose(cc); e != nil {
		return
	}
	_r = admHandler.RoomUsers(fromnode, domain)
	return
}

func (t *admprocessor) CreateRoom(ctx context.Context, fromnode, domain string, topic string, gtype int8) (_r *stub.AdmAck, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if e := noAuthAndClose(cc); e != nil {
		_r = newAdmAck(false)
		_r.Errcode = e.TimError().Code
		return
	}
	_r = admHandler.CreateRoom(fromnode, domain, topic, gtype)
	return
}

func (t *admprocessor) AddRoom(ctx context.Context, fromnode, domain string, roomNode string, msg string) (_r *stub.AdmAck, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if e := noAuthAndClose(cc); e != nil {
		_r = newAdmAck(false)
		_r.Errcode = e.TimError().Code
		return
	}
	_r = admHandler.AddRoom(fromnode, domain, roomNode, msg)
	return
}

func (t *admprocessor) PullInRoom(ctx context.Context, fromnode, domain string, roomNode string, toNode string) (_r *stub.AdmAck, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if e := noAuthAndClose(cc); e != nil {
		_r = newAdmAck(false)
		_r.Errcode = e.TimError().Code
		return
	}
	_r = admHandler.PullInRoom(fromnode, domain, roomNode, toNode)
	return
}

func (t *admprocessor) RejectRoom(ctx context.Context, fromnode, domain string, roomNode string, toNode string, msg string) (_r *stub.AdmAck, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if e := noAuthAndClose(cc); e != nil {
		_r = newAdmAck(false)
		_r.Errcode = e.TimError().Code
		return
	}
	_r = admHandler.RejectRoom(fromnode, domain, roomNode, toNode, msg)
	return
}

func (t *admprocessor) KickRoom(ctx context.Context, fromnode, domain string, roomNode string, toNode string) (_r *stub.AdmAck, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if e := noAuthAndClose(cc); e != nil {
		_r = newAdmAck(false)
		_r.Errcode = e.TimError().Code
		return
	}
	_r = admHandler.KickRoom(fromnode, domain, roomNode, toNode)
	return
}

func (t *admprocessor) LeaveRoom(ctx context.Context, fromnode, domain string, roomNode string) (_r *stub.AdmAck, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if e := noAuthAndClose(cc); e != nil {
		_r = newAdmAck(false)
		_r.Errcode = e.TimError().Code
		return
	}
	_r = admHandler.LeaveRoom(fromnode, domain, roomNode)
	return
}

func (t *admprocessor) CancelRoom(ctx context.Context, fromnode, domain string, roomNode string) (_r *stub.AdmAck, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if e := noAuthAndClose(cc); e != nil {
		_r = newAdmAck(false)
		_r.Errcode = e.TimError().Code
		return
	}
	_r = admHandler.CancelRoom(fromnode, domain, roomNode)
	return
}

func (t *admprocessor) BlockRoom(ctx context.Context, fromnode, domain string, roomNode string) (_r *stub.AdmAck, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if e := noAuthAndClose(cc); e != nil {
		_r = newAdmAck(false)
		_r.Errcode = e.TimError().Code
		return
	}
	_r = admHandler.BlockRoom(fromnode, domain, roomNode)
	return
}

func (t *admprocessor) BlockRoomMember(ctx context.Context, fromnode, domain string, roomNode string, toNode string) (_r *stub.AdmAck, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if e := noAuthAndClose(cc); e != nil {
		_r = newAdmAck(false)
		_r.Errcode = e.TimError().Code
		return
	}
	_r = admHandler.BlockRoomMember(fromnode, domain, roomNode, toNode)
	return
}

func (t *admprocessor) BlockRosterList(ctx context.Context, fromnode, domain string) (_r *stub.AdmNodeList, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if e := noAuthAndClose(cc); e == nil {
		_r = admHandler.BlockRosterList(fromnode, domain)
	}
	return
}

func (t *admprocessor) BlockRoomList(ctx context.Context, fromnode, domain string) (_r *stub.AdmNodeList, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if e := noAuthAndClose(cc); e == nil {
		_r = admHandler.BlockRoomList(fromnode, domain)
	}
	return
}

func (t *admprocessor) BlockRoomMemberlist(ctx context.Context, fromnode, domain string, roomNode string) (_r *stub.AdmNodeList, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if e := noAuthAndClose(cc); e == nil {
		_r = admHandler.BlockRoomMemberlist(fromnode, domain, roomNode)
	}
	return
}

func (t *admprocessor) VirtualroomRegister(ctx context.Context, fromnode, domain string) (_r *stub.AdmAck, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if e := noAuthAndClose(cc); e != nil {
		_r = newAdmAck(false)
		_r.Errcode = e.TimError().Code
		return
	}
	_r = admHandler.VirtualroomRegister(fromnode, domain)
	return
}

func (t *admprocessor) VirtualroomRemove(ctx context.Context, fromnode, domain string, roomNode string) (_r *stub.AdmAck, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if e := noAuthAndClose(cc); e != nil {
		_r = newAdmAck(false)
		_r.Errcode = e.TimError().Code
		return
	}
	_r = admHandler.VirtualroomRemove(fromnode, domain, roomNode)
	return
}

func (t *admprocessor) VirtualroomAddAuth(ctx context.Context, fromnode, domain string, roomNode string, toNode string) (_r *stub.AdmAck, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if e := noAuthAndClose(cc); e != nil {
		_r = newAdmAck(false)
		_r.Errcode = e.TimError().Code
		return
	}
	_r = admHandler.VirtualroomAddAuth(fromnode, domain, roomNode, toNode)
	return
}

func (t *admprocessor) VirtualroomDelAuth(ctx context.Context, fromnode, domain string, roomNode string, toNode string) (_r *stub.AdmAck, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if e := noAuthAndClose(cc); e != nil {
		_r = newAdmAck(false)
		_r.Errcode = e.TimError().Code
		return
	}
	_r = admHandler.VirtualroomDelAuth(fromnode, domain, roomNode, toNode)
	return
}

func (t *admprocessor) VirtualroomSub(ctx context.Context, wsid int64, fromnode string, domain string, vNode string, subType int8) (_r *stub.AdmAck, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if e := noAuthAndClose(cc); e != nil {
		_r = newAdmAck(false)
		_r.Errcode = e.TimError().Code
		return
	}
	_r = admHandler.VirtualroomSub(wsid, fromnode, domain, vNode, subType)
	return
}

func (t *admprocessor) VirtualroomUnSub(ctx context.Context, wsid int64, fromnode string, domain string, vNode string) (_r *stub.AdmAck, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if e := noAuthAndClose(cc); e != nil {
		_r = newAdmAck(false)
		_r.Errcode = e.TimError().Code
		return
	}
	_r = admHandler.VirtualroomUnSub(wsid, fromnode, domain, vNode)
	return
}

func (t *admprocessor) Authroster(ctx context.Context, fromnode string, domain string, tonode string) (_r *stub.AdmAck, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if e := noAuthAndClose(cc); e != nil {
		_r = newAdmAck(false)
		_r.Errcode = e.TimError().Code
		return
	}
	_r = admHandler.Authroster(fromnode, domain, tonode)
	return
}

func (t *admprocessor) Authgroupuser(ctx context.Context, fromnode string, domain string, roomNode string) (_r *stub.AdmAck, _err error) {
	defer util.Recover2(&_err)
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if e := noAuthAndClose(cc); e != nil {
		_r = newAdmAck(false)
		_r.Errcode = e.TimError().Code
		return
	}
	_r = admHandler.Authgroupuser(fromnode, domain, roomNode)
	return
}
