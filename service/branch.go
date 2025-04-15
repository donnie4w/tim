// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package service

import (
	goutil "github.com/donnie4w/gofer/util"
	"github.com/donnie4w/tim/data"
	"github.com/donnie4w/tim/errs"
	"github.com/donnie4w/tim/stub"
	"github.com/donnie4w/tim/sys"
	"github.com/donnie4w/tim/util"
	"github.com/donnie4w/tlnet"
	"time"
)

func (tss *timservice) business(bs []byte, ws *tlnet.Websocket) (err errs.ERROR) {
	defer util.Recover()
	sys.Stat.TxDo()
	defer sys.Stat.TxDone()
	tr := newTimReq(bs)
	if tr == nil || tr.Rtype == nil {
		return errs.ERR_PARAMS
	}
	if _, ok := wsware.Get(ws); ok {
		switch sys.BUSINESSTYPE(tr.GetRtype()) {
		case sys.BUSINESS_ROSTER:
			return tss.roster(ws)
		case sys.BUSINESS_USERROOM:
			return tss.usergroup(ws)
		case sys.BUSINESS_ROOMUSERS:
			return tss.grouproster(bs, ws)
		case sys.BUSINESS_ADDROSTER:
			return tss.addroster(bs, ws)
		case sys.BUSINESS_REMOVEROSTER:
			return tss.rmroster(bs, ws)
		case sys.BUSINESS_BLOCKROSTER:
			return tss.blockroster(bs, ws)
		case sys.BUSINESS_NEWROOM:
			return tss.newgroup(bs, ws)
		case sys.BUSINESS_ADDROOM:
			return tss.addgroup(bs, ws)
		case sys.BUSINESS_PULLROOM:
			return tss.pullgroup(bs, ws)
		case sys.BUSINESS_NOPASSROOM:
			return tss.nopassgroup(bs, ws)
		case sys.BUSINESS_KICKROOM:
			return tss.kickgroup(bs, ws)
		case sys.BUSINESS_LEAVEROOM:
			return tss.leavegroup(bs, ws)
		case sys.BUSINESS_CANCELROOM:
			return tss.cancelgroup(bs, ws)
		case sys.BUSINESS_BLOCKROOM:
			return tss.blockgroup(bs, ws)
		case sys.BUSINESS_BLOCKROOMMEMBER:
			return tss.blockgroupmember(bs, ws)
		case sys.BUSINESS_BLOCKROSTERLIST:
			return tss.blockrosterlist(ws)
		case sys.BUSINESS_BLOCKROOMLIST:
			return tss.blockroomlist(ws)
		case sys.BUSINESS_BLOCKROOMMEMBERLIST:
			return tss.blockroommemberlist(bs, ws)
		case sys.BUSINESS_MODIFYAUTH:
			return tss.modifyAuth(bs, ws)
		default:
			return errs.ERR_PARAMS
		}
	}
	return
}

func (tss *timservice) roster(ws *tlnet.Websocket) (err errs.ERROR) {
	if wss, ok := wsware.Get(ws); ok {
		tid := wss.tid
		if _r := data.Service.Roster(tid.Node); len(_r) > 0 {
			wss.Send(&stub.TimNodes{Ntype: int32(sys.NODEINFO_ROSTER), Nodelist: _r}, sys.TIMNODES)
		}
	}
	return
}

func (tss *timservice) usergroup(ws *tlnet.Websocket) (err errs.ERROR) {
	if wss, ok := wsware.Get(ws); ok {
		tid := wss.tid
		if _r := data.Service.UserGroup(tid.Node, tid.Domain); len(_r) > 0 {
			wss.Send(&stub.TimNodes{Ntype: int32(sys.NODEINFO_ROOM), Nodelist: _r}, sys.TIMNODES)
		}
	}
	return
}

func (tss *timservice) grouproster(bs []byte, ws *tlnet.Websocket) (err errs.ERROR) {
	if wss, ok := wsware.Get(ws); ok {
		tr := newTimReq(bs)
		node := tr.GetNode()
		if !checkNode(node) || !existGroup(&stub.Tid{Node: node, Domain: wss.tid.Domain}) {
			return errs.ERR_ACCOUNT
		}
		if ok := AuthGroup(node, wss.tid.Node, wss.tid.Domain); ok {
			if gs := data.Service.GroupRoster(node); len(gs) > 0 {
				wss.Send(&stub.TimNodes{Ntype: int32(sys.NODEINFO_ROOMMEMBER), Nodelist: gs, Node: &node}, sys.TIMNODES)
			}
		} else {
			err = errs.ERR_PARAMS
		}
	}
	return
}

func (tss *timservice) addroster(bs []byte, ws *tlnet.Websocket) (err errs.ERROR) {
	tq := newTimReq(bs)
	if tq == nil || tq.Node == nil || tq.ReqStr == nil {
		return errs.ERR_PARAMS
	}
	if !checkNode(tq.GetNode()) {
		return errs.ERR_ACCOUNT
	}
	if wss, b := wsware.Get(ws); b {
		if wss.tid.Node == tq.GetNode() {
			return errs.ERR_PERM_DENIED
		}
		id, t := goutil.UUID64(), time.Now().UnixNano()
		tm := &stub.TimMessage{ID: &id, MsType: sys.SOURCE_USER, OdType: sys.ORDER_BUSINESS, ToTid: &stub.Tid{Node: tq.GetNode()}, FromTid: wss.tid, DataString: tq.ReqStr, Timestamp: &t}
		var status int8
		if status, err = data.Service.Addroster(wss.tid.Node, tq.GetNode(), wss.tid.Domain); err == nil {
			if status == 0x11 {
				bt := int32(sys.BUSINESS_FRIEND)
				tm.BnType = &bt

			} else {
				bt := int32(sys.BUSINESS_ADDROSTER)
				tm.BnType = &bt
			}
			sys.TimMessageProcessor(tm, sys.TRANS_SOURCE)
		} else if status == 0x11 {
			bt := int32(sys.BUSINESS_FRIEND)
			tm.BnType = &bt
			sys.TimMessageProcessor(tm, sys.TRANS_SOURCE)
		}
	}
	return
}

func (tss *timservice) rmroster(bs []byte, ws *tlnet.Websocket) (err errs.ERROR) {
	tr := newTimReq(bs)
	if tr == nil || tr.Node == nil {
		return errs.ERR_PARAMS
	}
	if !checkNode(tr.GetNode()) {
		return errs.ERR_ACCOUNT
	}
	if wss, b := wsware.Get(ws); b {
		if ms, ok := data.Service.Rmroster(wss.tid.Node, tr.GetNode(), wss.tid.Domain); !ok {
			return errs.ERR_PARAMS
		} else {
			if ms {
				id, t, bt := goutil.UUID64(), time.Now().UnixNano(), int32(sys.BUSINESS_REMOVEROSTER)
				tm := &stub.TimMessage{ID: &id, MsType: sys.SOURCE_USER, BnType: &bt, OdType: sys.ORDER_BUSINESS, FromTid: wss.tid, ToTid: &stub.Tid{Node: tr.GetNode()}, Timestamp: &t}
				sys.TimMessageProcessor(tm, sys.TRANS_SOURCE)
			}
			t := int64(sys.BUSINESS_REMOVEROSTER)
			wss.Send(&stub.TimAck{Ok: true, TimType: int8(sys.TIMBUSINESS), N: tr.Node, T: &t}, sys.TIMACK)
		}
	}
	return
}

func (tss *timservice) blockroster(bs []byte, ws *tlnet.Websocket) (err errs.ERROR) {
	tq := newTimReq(bs)
	if tq == nil || tq.Node == nil {
		return errs.ERR_PARAMS
	}
	if !checkNode(tq.GetNode()) {
		return errs.ERR_ACCOUNT
	}
	if wss, b := wsware.Get(ws); b {
		if wss.tid.Node == tq.GetNode() {
			return errs.ERR_PERM_DENIED
		}
		if ms, ok := data.Service.Blockroster(wss.tid.Node, tq.GetNode(), wss.tid.Domain); !ok {
			return errs.ERR_PARAMS
		} else {
			if ms {
				id, _t, bt := goutil.UUID64(), time.Now().UnixNano(), int32(sys.BUSINESS_BLOCKROSTER)
				tm := &stub.TimMessage{ID: &id, MsType: sys.SOURCE_USER, OdType: sys.ORDER_BUSINESS, ToTid: &stub.Tid{Node: tq.GetNode()}, FromTid: wss.tid, BnType: &bt, Timestamp: &_t}
				sys.TimMessageProcessor(tm, sys.TRANS_SOURCE)
			}
			t := int64(sys.BUSINESS_BLOCKROSTER)
			wss.Send(&stub.TimAck{Ok: true, TimType: int8(sys.TIMBUSINESS), N: tq.Node, T: &t}, sys.TIMACK)
		}
	}
	return
}

func (tss *timservice) newgroup(bs []byte, ws *tlnet.Websocket) (err errs.ERROR) {
	tr := newTimReq(bs)
	if tr == nil || tr.Node == nil || tr.ReqInt == nil {
		return errs.ERR_PARAMS
	}
	if wss, b := wsware.Get(ws); b {
		rtype := sys.GROUP_PRIVATE
		if sys.TIMTYPE(tr.GetReqInt()) != rtype {
			rtype = sys.GROUP_OPEN
		}
		gnode, e := data.Service.Newgroup(wss.tid.Node, tr.GetNode(), rtype, wss.tid.Domain)
		t := int64(sys.BUSINESS_NEWROOM)
		if e == nil {
			//id, unix := goutil.UUID64(), time.Now().UnixNano()
			//bnType := int32(sys.BUSINESS_NEWROOM)
			//tm := &stub.TimMessage{ID: &id, BnType: &bnType, RoomTid: &stub.Tid{Node: tr.GetNode()}, MsType: sys.SOURCE_USER, OdType: sys.ORDER_BUSINESS, Timestamp: &unix}
			//timMessage4goal(wss.tid.Node, tm)
			wss.Send(&stub.TimAck{Ok: true, TimType: int8(sys.TIMBUSINESS), N: &gnode, T: &t}, sys.TIMACK)
		} else {
			wss.Send(&stub.TimAck{Ok: false, TimType: int8(sys.TIMBUSINESS), N: &gnode, T: &t, Error: err.TimError()}, sys.TIMACK)
		}

	}
	return
}

func (tss *timservice) osnewgroup(unode, gname string, domain *string, gtype int8) (string, errs.ERROR) {
	if !checkNode(unode) {
		return "", errs.ERR_ACCOUNT
	}
	return data.Service.Newgroup(unode, gname, sys.TIMTYPE(gtype), domain)
}

func (tss *timservice) osModifygroupInfo(unode, gnode string, trb *stub.TimRoomBean) errs.ERROR {
	if !checkNode(gnode) || !checkNode(unode) {
		return errs.ERR_ACCOUNT
	}
	return data.Service.ModifygroupInfo(gnode, unode, trb, true)
}

func (tss *timservice) addgroup(bs []byte, ws *tlnet.Websocket) (err errs.ERROR) {
	tr := newTimReq(bs)
	if tr == nil || tr.Node == nil || !checkNode(tr.GetNode()) || tr.ReqStr == nil {
		return errs.ERR_PARAMS
	}
	if wss, b := wsware.Get(ws); b {
		var gtype int8
		if gtype, err = data.Service.GroupGtype(tr.GetNode(), wss.tid.Domain); err == nil {
			unix := time.Now().UnixNano()
			switch sys.TIMTYPE(gtype) {
			case sys.GROUP_OPEN:
				if err = data.Service.Addgroup(tr.GetNode(), wss.tid.Node, wss.tid.Domain); err == nil {
					id := goutil.UUID64()
					bnType := int32(sys.BUSINESS_PASSROOM)
					tm := &stub.TimMessage{ID: &id, BnType: &bnType, RoomTid: &stub.Tid{Node: tr.GetNode()}, MsType: sys.SOURCE_USER, OdType: sys.ORDER_BUSINESS, Timestamp: &unix}
					timMessage4goal(wss.tid.Node, tm)
				}
			case sys.GROUP_PRIVATE:
				if err = data.Service.Addgroup(tr.GetNode(), wss.tid.Node, wss.tid.Domain); err != nil {
					return
				}
				var ms []string
				if ms, err = data.Service.GroupManagers(tr.GetNode(), wss.tid.Domain); err == nil {
					bnType, roomtid := int32(sys.BUSINESS_ADDROOM), &stub.Tid{Node: tr.GetNode()}
					for _, u := range ms {
						id := goutil.UUID64()
						tm := &stub.TimMessage{ID: &id, FromTid: wss.tid, ToTid: &stub.Tid{Node: u}, BnType: &bnType, RoomTid: roomtid, MsType: sys.SOURCE_USER, OdType: sys.ORDER_BUSINESS, DataString: tr.ReqStr, Timestamp: &unix}
						sys.TimMessageProcessor(tm, sys.TRANS_SOURCE)
					}
				}
			}
		}
		if err != nil {
			t := int64(sys.BUSINESS_ADDROOM)
			wss.Send(&stub.TimAck{Ok: false, TimType: int8(sys.TIMBUSINESS), N: tr.Node, T: &t, Error: err.TimError()}, sys.TIMACK)
		}
	}
	return
}

func (tss *timservice) pullgroup(bs []byte, ws *tlnet.Websocket) (err errs.ERROR) {
	tr := newTimReq(bs)
	if tr == nil || tr.Node == nil || tr.Node2 == nil || !checkNode(tr.GetNode()) || !checkNode(tr.GetNode2()) {
		return errs.ERR_PARAMS
	}
	if wss, b := wsware.Get(ws); b {
		if wss.tid.Node == tr.GetNode2() {
			return errs.ERR_PERM_DENIED
		}
		var isreq bool
		if isreq, err = data.Service.Pullgroup(tr.GetNode(), wss.tid.Node, tr.GetNode2(), wss.tid.Domain); err == nil {
			id, unix := goutil.UUID64(), time.Now().UnixNano()
			tm := &stub.TimMessage{ID: &id, FromTid: wss.tid, ToTid: &stub.Tid{Node: tr.GetNode2()}, RoomTid: &stub.Tid{Node: tr.GetNode()}, MsType: sys.SOURCE_USER, OdType: sys.ORDER_BUSINESS, Timestamp: &unix}
			if isreq {
				bt := int32(sys.BUSINESS_PASSROOM)
				tm.BnType = &bt
			} else {
				bt := int32(sys.BUSINESS_PULLROOM)
				tm.BnType = &bt
			}
			sys.TimMessageProcessor(tm, sys.TRANS_SOURCE)
		}
		t := int64(sys.BUSINESS_PULLROOM)
		if err == nil {
			wss.Send(&stub.TimAck{Ok: true, TimType: int8(sys.TIMBUSINESS), N: tr.Node, T: &t, N2: tr.Node2}, sys.TIMACK)
		} else {
			wss.Send(&stub.TimAck{Ok: false, TimType: int8(sys.TIMBUSINESS), N: tr.Node, T: &t, Error: err.TimError()}, sys.TIMACK)
		}
	}
	return
}

func (tss *timservice) nopassgroup(bs []byte, ws *tlnet.Websocket) (err errs.ERROR) {
	tr := newTimReq(bs)
	if tr == nil || tr.Node == nil || tr.Node2 == nil || !checkNode(tr.GetNode()) || !checkNode(tr.GetNode2()) {
		return errs.ERR_PARAMS
	}
	if wss, b := wsware.Get(ws); b {
		if err = data.Service.Nopassgroup(tr.GetNode(), wss.tid.Node, tr.GetNode2(), wss.tid.Domain); err == nil {
			id, unix, bt := goutil.UUID64(), time.Now().UnixNano(), int32(sys.BUSINESS_NOPASSROOM)
			tm := &stub.TimMessage{ID: &id, FromTid: wss.tid, ToTid: &stub.Tid{Node: tr.GetNode2()}, RoomTid: &stub.Tid{Node: tr.GetNode()}, MsType: sys.SOURCE_USER, OdType: sys.ORDER_BUSINESS, Timestamp: &unix}
			tm.BnType = &bt
			sys.TimMessageProcessor(tm, sys.TRANS_SOURCE)
		}
		t := int64(sys.BUSINESS_NOPASSROOM)
		if err == nil {
			wss.Send(&stub.TimAck{Ok: true, TimType: int8(sys.TIMBUSINESS), N: tr.Node, N2: tr.Node2, T: &t}, sys.TIMACK)
		} else {
			wss.Send(&stub.TimAck{Ok: false, TimType: int8(sys.TIMBUSINESS), N: tr.Node, N2: tr.Node2, T: &t, Error: err.TimError()}, sys.TIMACK)
		}
	}
	return
}

func (tss *timservice) kickgroup(bs []byte, ws *tlnet.Websocket) (err errs.ERROR) {
	tr := newTimReq(bs)
	if tr == nil || tr.Node == nil || tr.Node2 == nil || !checkNode(tr.GetNode()) || !checkNode(tr.GetNode2()) {
		return errs.ERR_PARAMS
	}
	if wss, b := wsware.Get(ws); b {
		if wss.tid.Node == tr.GetNode2() {
			return errs.ERR_PERM_DENIED
		}
		if err = data.Service.Kickgroup(tr.GetNode(), wss.tid.Node, tr.GetNode2(), wss.tid.Domain); err == nil {
			id, unix, bnType := goutil.UUID64(), time.Now().UnixNano(), int32(sys.BUSINESS_KICKROOM)
			tm := &stub.TimMessage{ID: &id, BnType: &bnType, FromTid: wss.tid, ToTid: &stub.Tid{Node: tr.GetNode2()}, RoomTid: &stub.Tid{Node: tr.GetNode()}, MsType: sys.SOURCE_USER, OdType: sys.ORDER_BUSINESS, Timestamp: &unix}
			sys.TimMessageProcessor(tm, sys.TRANS_SOURCE)
		}
		t := int64(sys.BUSINESS_KICKROOM)
		if err == nil {
			wss.Send(&stub.TimAck{Ok: true, TimType: int8(sys.TIMBUSINESS), N: tr.Node, N2: tr.Node2, T: &t}, sys.TIMACK)
		} else {
			wss.Send(&stub.TimAck{Ok: false, TimType: int8(sys.TIMBUSINESS), N: tr.Node, N2: tr.Node2, T: &t, Error: err.TimError()}, sys.TIMACK)
		}
	}
	return
}

func (tss *timservice) leavegroup(bs []byte, ws *tlnet.Websocket) (err errs.ERROR) {
	tr := newTimReq(bs)
	if tr == nil || tr.Node == nil || !checkNode(tr.GetNode()) {
		return errs.ERR_PARAMS
	}
	if wss, b := wsware.Get(ws); b {
		err = data.Service.Leavegroup(tr.GetNode(), wss.tid.Node, wss.tid.Domain)
		t := int64(sys.BUSINESS_LEAVEROOM)
		if err == nil {
			//id, unix, bnType := goutil.UUID64(), time.Now().UnixNano(), int32(sys.BUSINESS_LEAVEROOM)
			//tm := &stub.TimMessage{ID: &id, BnType: &bnType, RoomTid: &stub.Tid{Node: tr.GetNode()}, MsType: sys.SOURCE_USER, OdType: sys.ORDER_BUSINESS, Timestamp: &unix}
			//timMessage4goal(wss.tid.Node, tm)
			wss.Send(&stub.TimAck{Ok: true, TimType: int8(sys.TIMBUSINESS), N: tr.Node, T: &t}, sys.TIMACK)
		} else {
			wss.Send(&stub.TimAck{Ok: false, TimType: int8(sys.TIMBUSINESS), N: tr.Node, T: &t, Error: err.TimError()}, sys.TIMACK)
		}
	}
	return
}

func (tss *timservice) cancelgroup(bs []byte, ws *tlnet.Websocket) (err errs.ERROR) {
	tr := newTimReq(bs)
	if tr == nil || tr.Node == nil || !checkNode(tr.GetNode()) {
		return errs.ERR_PARAMS
	}
	if wss, b := wsware.Get(ws); b {
		err = data.Service.Cancelgroup(tr.GetNode(), wss.tid.Node, wss.tid.Domain)

		t := int64(sys.BUSINESS_CANCELROOM)
		if err == nil {
			//id, unix, bnType := goutil.UUID64(), time.Now().UnixNano(), int32(sys.BUSINESS_CANCELROOM)
			//tm := &stub.TimMessage{ID: &id, BnType: &bnType, RoomTid: &stub.Tid{Node: tr.GetNode()}, MsType: sys.SOURCE_USER, OdType: sys.ORDER_BUSINESS, Timestamp: &unix}
			//timMessage4goal(wss.tid.Node, tm)
			wss.Send(&stub.TimAck{Ok: true, TimType: int8(sys.TIMBUSINESS), N: tr.Node, T: &t}, sys.TIMACK)
		} else {
			wss.Send(&stub.TimAck{Ok: false, TimType: int8(sys.TIMBUSINESS), N: tr.Node, T: &t, Error: err.TimError()}, sys.TIMACK)
		}
	}
	return
}

func (tss *timservice) blockgroup(bs []byte, ws *tlnet.Websocket) (err errs.ERROR) {
	tr := newTimReq(bs)
	if tr == nil || tr.Node == nil || !checkNode(tr.GetNode()) {
		return errs.ERR_PARAMS
	}
	if wss, b := wsware.Get(ws); b {
		err = data.Service.Blockgroup(tr.GetNode(), wss.tid.Node, wss.tid.Domain)

		t := int64(sys.BUSINESS_BLOCKROOM)
		if err == nil {
			//id, unix, bnType := goutil.UUID64(), time.Now().UnixNano(), int32(sys.BUSINESS_BLOCKROOM)
			//tm := &stub.TimMessage{ID: &id, BnType: &bnType, RoomTid: &stub.Tid{Node: tr.GetNode()}, MsType: sys.SOURCE_USER, OdType: sys.ORDER_BUSINESS, Timestamp: &unix}
			//timMessage4goal(wss.tid.Node, tm)
			wss.Send(&stub.TimAck{Ok: true, TimType: int8(sys.TIMBUSINESS), N: tr.Node, T: &t}, sys.TIMACK)
		} else {
			wss.Send(&stub.TimAck{Ok: false, TimType: int8(sys.TIMBUSINESS), N: tr.Node, T: &t, Error: err.TimError()}, sys.TIMACK)
		}
	}
	return
}

func (tss *timservice) blockgroupmember(bs []byte, ws *tlnet.Websocket) (err errs.ERROR) {
	tr := newTimReq(bs)
	if tr == nil || tr.Node == nil || !checkNode(tr.GetNode()) || tr.Node2 == nil || !checkNode(tr.GetNode2()) {
		return errs.ERR_PARAMS
	}
	if wss, b := wsware.Get(ws); b {
		if wss.tid.Node == tr.GetNode2() {
			return errs.ERR_PERM_DENIED
		}
		err = data.Service.Blockgroupmember(tr.GetNode(), wss.tid.Node, tr.GetNode2(), wss.tid.Domain)

		t := int64(sys.BUSINESS_BLOCKROOMMEMBER)
		if err == nil {
			//id, unix, bt := goutil.UUID64(), time.Now().UnixNano(), int32(sys.BUSINESS_BLOCKROOMMEMBER)
			//tm := &stub.TimMessage{ID: &id, BnType: &bt, FromTid: wss.tid, ToTid: &stub.Tid{Node: tr.GetNode2()}, RoomTid: &stub.Tid{Node: tr.GetNode()}, MsType: sys.SOURCE_USER, OdType: sys.ORDER_BUSINESS, Timestamp: &unix}
			//sys.TimMessageProcessor(tm, sys.TRANS_SOURCE)
			wss.Send(&stub.TimAck{Ok: true, TimType: int8(sys.TIMBUSINESS), N: tr.Node, T: &t, N2: tr.Node2}, sys.TIMACK)
		} else {
			wss.Send(&stub.TimAck{Ok: false, TimType: int8(sys.TIMBUSINESS), N: tr.Node, T: &t, N2: tr.Node2, Error: err.TimError()}, sys.TIMACK)
		}
	}
	return
}

func (tss *timservice) blockrosterlist(ws *tlnet.Websocket) (err errs.ERROR) {
	if wss, b := wsware.Get(ws); b {
		if ss := data.Service.Blockrosterlist(wss.tid.Node); len(ss) > 0 {
			wss.Send(&stub.TimNodes{Ntype: int32(sys.NODEINFO_BLOCKROSTERLIST), Nodelist: ss}, sys.TIMNODES)
		}
	}
	return
}

func (tss *timservice) blockroomlist(ws *tlnet.Websocket) (err errs.ERROR) {
	if wss, b := wsware.Get(ws); b {
		if ss := data.Service.Blockroomlist(wss.tid.Node); len(ss) > 0 {
			wss.Send(&stub.TimNodes{Ntype: int32(sys.NODEINFO_BLOCKROOMLIST), Nodelist: ss}, sys.TIMNODES)
		}
	}
	return
}

func (tss *timservice) blockroommemberlist(bs []byte, ws *tlnet.Websocket) (err errs.ERROR) {
	tr := newTimReq(bs)
	if tr == nil || tr.Node == nil || !checkNode(tr.GetNode()) {
		return errs.ERR_PARAMS
	}
	if wss, b := wsware.Get(ws); b {
		if ss := data.Service.Blockroommemberlist(tr.GetNode(), wss.tid.Node); len(ss) > 0 {
			wss.Send(&stub.TimNodes{Ntype: int32(sys.NODEINFO_BLOCKROOMMEMBERLIST), Nodelist: ss}, sys.TIMNODES)
		}
	}
	return
}

func (tss *timservice) sysModify(node string, oldpwd *string, newpwd string, domain *string) errs.ERROR {
	return data.Service.Modify(util.NodeToUUID(node), oldpwd, newpwd, domain)
}

func (tss *timservice) modifyAuth(bs []byte, ws *tlnet.Websocket) (err errs.ERROR) {
	tr := newTimReq(bs)
	if tr == nil || tr.ReqStr == nil || tr.ReqStr2 == nil {
		return errs.ERR_MODIFYAUTH
	}
	if wss, b := wsware.Get(ws); b {
		t := int64(sys.BUSINESS_MODIFYAUTH)
		if err = data.Service.Modify(util.NodeToUUID(wss.tid.Node), tr.ReqStr, tr.GetReqStr2(), wss.tid.Domain); err == nil {
			wss.Send(&stub.TimAck{Ok: true, TimType: int8(sys.TIMBUSINESS), T: &t}, sys.TIMACK)
		} else {
			wss.Send(&stub.TimAck{Ok: false, Error: errs.ERR_MODIFYAUTH.TimError(), TimType: int8(sys.TIMBUSINESS), T: &t}, sys.TIMACK)
		}
	}
	return
}

func (tss *timservice) osuserbean(node string, tb *stub.TimUserBean) (err errs.ERROR) {
	return data.Service.ModifyUserInfo(node, tb)
}

func (tss *timservice) nodeinfo(bs []byte, ws *tlnet.Websocket) (err errs.ERROR) {
	defer util.Recover()
	sys.Stat.TxDo()
	defer sys.Stat.TxDone()
	tr := newTimNodes(bs)
	if tr == nil {
		return errs.ERR_PARAMS
	}
	if wss, b := wsware.Get(ws); b {
		switch sys.BUSINESSTYPE(tr.Ntype) {
		case sys.NODEINFO_USERINFO:
			if m, e := data.Service.GetUserInfo(tr.Nodelist); e == nil && len(m) > 0 {
				tr.Nodelist, tr.Usermap = nil, m
				wss.Send(tr, sys.TIMNODES)
			} else {
				err = e
			}
		case sys.NODEINFO_ROOMINFO:
			if m, e := data.Service.GetGroupInfo(tr.Nodelist); e == nil && len(m) > 0 {
				tr.Nodelist, tr.Roommap = nil, m
				wss.Send(tr, sys.TIMNODES)
			} else {
				err = e
			}
		case sys.NODEINFO_MODIFYUSER:
			if len(tr.Usermap) == 1 {
				for _, v := range tr.Usermap {
					if err = tss.osuserbean(wss.tid.Node, v); err == nil {
						t := int64(sys.NODEINFO_MODIFYUSER)
						wss.Send(&stub.TimAck{Ok: true, TimType: int8(sys.TIMNODES), T: &t}, sys.TIMACK)
					}
				}
			} else {
				return errs.ERR_PARAMS
			}
		case sys.NODEINFO_MODIFYROOM:
			if len(tr.Roommap) == 1 {
				for k, v := range tr.Roommap {
					if err = data.Service.ModifygroupInfo(k, wss.tid.Node, v, false); err == nil {
						t := int64(sys.NODEINFO_MODIFYROOM)
						wss.Send(&stub.TimAck{Ok: true, TimType: int8(sys.TIMNODES), T: &t}, sys.TIMACK)
					}
				}
			} else {
				return errs.ERR_PARAMS
			}
		default:
			err = errs.ERR_PARAMS
		}
	}
	return
}
