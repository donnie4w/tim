// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of t source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim
package adm

import (
	"context"
	. "github.com/donnie4w/tim/stub"
	"github.com/donnie4w/tim/util"
)

type ifacehandle struct {
}

var ifaceProcessor = &ifacehandle{}

// Parameters:
//   - Loginname
//   - Domain
//   - Pwd
func (t *ifacehandle) ResetAuth(ctx context.Context, ab *AuthBean) (_r *AdmAck, _err error) {
	defer util.Recover()
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if noAuthAndClose(cc) {
		_r = newAdmAck(false)
		return
	}
	if ab.Username != nil || ab.Password != nil {
		_r = newAdmAck(false)
		return
	}
	_r = Admhandler.ResetAuth(ab)
	return
}

// Parameters:
//   - Ab
func (t *ifacehandle) Auth(ctx context.Context, ab *AuthBean) (_r *AdmAck, _err error) {
	defer util.Recover()
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if _r = Admhandler.Auth(ab); _r.GetOk() {
		ctx2CliContext(ctx).isAuth = true
	}
	return
}

func (t *ifacehandle) Ping(ctx context.Context) (_r *AdmAck, _err error) {
	defer util.Recover()
	mux := ctx2CliContext(ctx).mux
	defer mux.Unlock()
	mux.Lock()
	_r = newAdmAck(true)
	return
}

// Parameters:
//   - Atoken
func (t *ifacehandle) Token(ctx context.Context, atoken *AdmToken) (_r *AdmAck, _err error) {
	defer util.Recover()
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if noAuthAndClose(cc) {
		_r = newAdmAck(false)
		return
	}
	if atoken.GetName() == "" {
		_r = newAdmAck(false)
		return
	}
	_r = Admhandler.Token(atoken)
	return
}

// Parameters:
//   - Am
func (t *ifacehandle) OsMessage(ctx context.Context, am *AdmMessage) (_r *AdmAck, _err error) {
	defer util.Recover()
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if noAuthAndClose(cc) {
		_r = newAdmAck(false)
		return
	}
	_r = Admhandler.OsMessage(am)
	return
}

// Parameters:
//   - amb
func (t *ifacehandle) ProxyMessage(ctx context.Context, amb *AdmProxyMessage) (_r *AdmAck, _err error) {
	defer util.Recover()

	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if amb.GetConnectid() <= 0 || noAuthAndClose(cc) {
		_r = newAdmAck(false)
		return
	}
	_r = Admhandler.ProxyMessage(amb)
	return
}

// Parameters:
//   - ab
func (t *ifacehandle) Register(ctx context.Context, ab *AuthBean) (_r *AdmAck, _err error) {
	defer util.Recover()
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if noAuthAndClose(cc) {
		_r = newAdmAck(false)
		return
	}
	_r = Admhandler.Register(ab)
	return
}

// Parameters:
//   - Amui
func (t *ifacehandle) ModifyUserInfo(ctx context.Context, amui *AdmModifyUserInfo) (_r *AdmAck, _err error) {
	defer util.Recover()
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if noAuthAndClose(cc) {
		_r = newAdmAck(false)
		return
	}
	_r = Admhandler.ModifyUserInfo(amui)
	return
}

// Parameters:
//   - Abu
func (t *ifacehandle) BlockUser(ctx context.Context, abu *AdmBlockUser) (_r *AdmAck, _err error) {
	defer util.Recover()
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if noAuthAndClose(cc) {
		_r = newAdmAck(false)
		return
	}
	_r = Admhandler.BlockUser(abu)
	return
}

func (t *ifacehandle) BlockList(ctx context.Context) (_r *AdmBlockList, _err error) {
	defer util.Recover()
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if noAuthAndClose(cc) {
		return
	}
	_r = Admhandler.BlockList()
	return
}

// Parameters:
//   - Au
func (t *ifacehandle) OnlineUser(ctx context.Context, au *AdmOnlineUser) (_r *AdmTidList, _err error) {
	defer util.Recover()
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if noAuthAndClose(cc) {
		return
	}
	_r = Admhandler.OnlineUser(au)
	return
}

// Parameters:
//   - Avb
func (t *ifacehandle) Vroom(ctx context.Context, avb *AdmVroomBean) (_r *AdmAck, _err error) {
	defer util.Recover()
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if noAuthAndClose(cc) {
		_r = newAdmAck(false)
		return
	}
	_r = Admhandler.Vroom(avb)
	return
}

// Parameters:
//   - Arb
func (t *ifacehandle) TimRoom(ctx context.Context, arb *AdmRoomReq) (_r *AdmAck, _err error) {
	defer util.Recover()
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if noAuthAndClose(cc) {
		_r = newAdmAck(false)
		return
	}
	_r = Admhandler.TimRoom(arb)
	return
}

// Parameters:
//   - Adb
func (t *ifacehandle) Detect(ctx context.Context, adb *AdmDetectBean) (_r *AdmAck, _err error) {
	defer util.Recover()
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if noAuthAndClose(cc) {
		_r = newAdmAck(false)
		return
	}
	_r = Admhandler.Detect(adb)
	return
}

// Parameters:
//   - Arb
func (t *ifacehandle) ModifyRoomInfo(ctx context.Context, arb *AdmRoomBean) (_r *AdmAck, _err error) {
	defer util.Recover()
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if noAuthAndClose(cc) {
		_r = newAdmAck(false)
		return
	}
	_r = Admhandler.ModifyRoomInfo(arb)
	return
}
