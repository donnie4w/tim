// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim
package service

import (
	"time"

	. "github.com/donnie4w/gofer/util"
	"github.com/donnie4w/tim/data"
	. "github.com/donnie4w/tim/stub"
	"github.com/donnie4w/tim/sys"
	"github.com/donnie4w/tim/util"
	"github.com/donnie4w/tlnet"
)

func (this *timservice) business(bs []byte, ws *tlnet.Websocket) (err sys.ERROR) {
	defer util.Recover()
	sys.Stat.TxDo()
	defer sys.Stat.TxDone()
	tr := newTimReq(bs)
	if tr == nil || tr.Rtype == nil {
		return sys.ERR_PARAMS
	}
	if _, ok := wsware.Get(ws); ok {
		switch *tr.Rtype {
		case sys.BUSINESS_ROSTER:
			return this.roster(ws)
		case sys.BUSINESS_USERROOM:
			return this.usergroup(ws)
		case sys.BUSINESS_ROOMUSERS:
			return this.grouproster(bs, ws)
		case sys.BUSINESS_ADDROSTER:
			return this.addroster(bs, ws)
		case sys.BUSINESS_REMOVEROSTER:
			return this.rmroster(bs, ws)
		case sys.BUSINESS_BLOCKROSTER:
			return this.blockroster(bs, ws)
		case sys.BUSINESS_NEWROOM:
			return this.newgroup(bs, ws)
		case sys.BUSINESS_ADDROOM:
			return this.addgroup(bs, ws)
		case sys.BUSINESS_PULLROOM:
			return this.pullgroup(bs, ws)
		case sys.BUSINESS_NOPASSROOM:
			return this.nopassgroup(bs, ws)
		case sys.BUSINESS_KICKROOM:
			return this.kickgroup(bs, ws)
		case sys.BUSINESS_LEAVEROOM:
			return this.leavegroup(bs, ws)
		case sys.BUSINESS_CANCELROOM:
			return this.cancelgroup(bs, ws)
		case sys.BUSINESS_BLOCKROOM:
			return this.blockgroup(bs, ws)
		case sys.BUSINESS_BLOCKROOMMEMBER:
			return this.blockgroupmember(bs, ws)
		case sys.BUSINESS_BLOCKROSTERLIST:
			return this.blockrosterlist(ws)
		case sys.BUSINESS_BLOCKROOMLIST:
			return this.blockroomlist(ws)
		case sys.BUSINESS_BLOCKROOMMEMBERLIST:
			return this.blockroommemberlist(bs, ws)
		case sys.BUSINESS_MODIFYAUTH:
			return this.modifyauth(bs, ws)
		default:
			return sys.ERR_PARAMS
		}
	}
	return
}

func (this *timservice) roster(ws *tlnet.Websocket) (err sys.ERROR) {
	if wss, ok := wsware.Get(ws); ok {
		tid := wss.tid
		if _r := data.Handler.Roster(tid.Node); _r != nil && len(_r) > 0 {
			wsware.SendWs(ws.Id, &TimNodes{Ntype: sys.NODEINFO_ROSTER, Nodelist: _r}, sys.TIMNODES)
		}
	}
	return
}

func (this *timservice) usergroup(ws *tlnet.Websocket) (err sys.ERROR) {
	if wss, ok := wsware.Get(ws); ok {
		tid := wss.tid
		if _r := data.Handler.UserGroup(tid.Node, tid.Domain); _r != nil && len(_r) > 0 {
			wsware.SendWs(ws.Id, &TimNodes{Ntype: sys.NODEINFO_ROOM, Nodelist: _r}, sys.TIMNODES)
		}
	}
	return
}

func (this *timservice) grouproster(bs []byte, ws *tlnet.Websocket) (err sys.ERROR) {
	if wss, ok := wsware.Get(ws); ok {
		tr := newTimReq(bs)
		node := *tr.Node
		if !checkNode(node) || !existGroup(&Tid{Node: node, Domain: wss.tid.Domain}) {
			return sys.ERR_ACCOUNT
		}
		if ok := authGroup(node, wss.tid.Node, wss.tid.Domain); ok {
			if gs := data.Handler.GroupRoster(node); gs != nil && len(gs) > 0 {
				wsware.SendWs(ws.Id, &TimNodes{Ntype: sys.NODEINFO_ROOMMEMBER, Nodelist: gs, Node: &node}, sys.TIMNODES)
			}
		} else {
			err = sys.ERR_PARAMS
		}
	}
	return
}

func (this *timservice) addroster(bs []byte, ws *tlnet.Websocket) (err sys.ERROR) {
	tq := newTimReq(bs)
	if tq == nil || tq.Node == nil || tq.ReqStr == nil {
		return sys.ERR_PARAMS
	}
	if !checkNode(*tq.Node) {
		return sys.ERR_ACCOUNT
	}
	if wss, b := wsware.Get(ws); b {
		if wss.tid.Node == *tq.Node {
			return sys.ERR_AUTH
		}
		var status int8
		if status, err = data.Handler.Addroster(wss.tid.Node, *tq.Node, wss.tid.Domain); err == nil {
			id, _t := RandId(), time.Now().UnixNano()
			tm := &TimMessage{ID: &id, MsType: sys.SOURCE_USER, OdType: sys.ORDER_BUSINESS, ToTid: &Tid{Node: *tq.Node}, FromTid: wss.tid, DataString: tq.ReqStr, Timestamp: &_t}
			if status == 0x10|0x01 {
				tm.BnType = &sys.BUSINESS_FRIEND
				wsware.SendWs(ws.Id, tm, sys.TIMMESSAGE)
			} else {
				tm.BnType = &sys.BUSINESS_ADDROSTER
			}
			sys.TimMessageProcessor(tm, sys.TRANS_SOURCE)
		}
	}
	return
}

func (this *timservice) rmroster(bs []byte, ws *tlnet.Websocket) (err sys.ERROR) {
	tr := newTimReq(bs)
	if tr == nil || tr.Node == nil {
		return sys.ERR_PARAMS
	}
	if !checkNode(*tr.Node) {
		return sys.ERR_ACCOUNT
	}
	if wss, b := wsware.Get(ws); b {
		if ms, ok := data.Handler.Rmroster(wss.tid.Node, *tr.Node, wss.tid.Domain); !ok {
			return sys.ERR_PARAMS
		} else {
			if ms {
				id, _t := RandId(), time.Now().UnixNano()
				tm := &TimMessage{ID: &id, MsType: sys.SOURCE_USER, BnType: &sys.BUSINESS_REMOVEROSTER, OdType: sys.ORDER_BUSINESS, FromTid: wss.tid, ToTid: &Tid{Node: *tr.Node}, Timestamp: &_t}
				sys.TimMessageProcessor(tm, sys.TRANS_SOURCE)
			}
			t := int64(sys.BUSINESS_REMOVEROSTER)
			wsware.SendWs(ws.Id, &TimAck{Ok: true, TimType: int8(sys.TIMBUSINESS), N: tr.Node, T: &t}, sys.TIMACK)
		}
	}
	return
}

func (this *timservice) blockroster(bs []byte, ws *tlnet.Websocket) (err sys.ERROR) {
	tq := newTimReq(bs)
	if tq == nil || tq.Node == nil {
		return sys.ERR_PARAMS
	}
	if !checkNode(*tq.Node) {
		return sys.ERR_ACCOUNT
	}
	if wss, b := wsware.Get(ws); b {
		if wss.tid.Node == *tq.Node {
			return sys.ERR_AUTH
		}
		if ms, ok := data.Handler.Blockroster(wss.tid.Node, *tq.Node, wss.tid.Domain); !ok {
			return sys.ERR_PARAMS
		} else {
			if ms {
				id, _t := RandId(), time.Now().UnixNano()
				tm := &TimMessage{ID: &id, MsType: sys.SOURCE_USER, OdType: sys.ORDER_BUSINESS, ToTid: &Tid{Node: *tq.Node}, FromTid: wss.tid, BnType: &sys.BUSINESS_BLOCKROSTER, Timestamp: &_t}
				sys.TimMessageProcessor(tm, sys.TRANS_SOURCE)
			}
			t := int64(sys.BUSINESS_BLOCKROSTER)
			wsware.SendWs(ws.Id, &TimAck{Ok: true, TimType: int8(sys.TIMBUSINESS), N: tq.Node, T: &t}, sys.TIMACK)
		}
	}
	return
}

func (this *timservice) newgroup(bs []byte, ws *tlnet.Websocket) (err sys.ERROR) {
	tr := newTimReq(bs)
	if tr == nil || tr.Node == nil || tr.ReqInt == nil {
		return sys.ERR_PARAMS
	}
	if wss, b := wsware.Get(ws); b {
		rtype := sys.GROUP_PRIVATE
		if int8(*tr.ReqInt) != rtype {
			rtype = sys.GROUP_OPEN
		}
		var gnode string
		if gnode, err = data.Handler.Newgroup(wss.tid.Node, *tr.Node, rtype, wss.tid.Domain); err == nil {
			t := int64(sys.BUSINESS_NEWROOM)
			wsware.SendWs(ws.Id, &TimAck{Ok: true, TimType: int8(sys.TIMBUSINESS), N: &gnode, T: &t}, sys.TIMACK)
		}
	}
	return
}

func (this *timservice) osnewgroup(unode, gname string, domain *string, gtype int8) (gnode string, err sys.ERROR) {
	gnode, err = data.Handler.Newgroup(unode, gname, gtype, domain)
	return
}

func (this *timservice) osModifygroupInfo(unode, gnode string, trb *TimRoomBean) sys.ERROR {
	return data.Handler.ModifygroupInfo(gnode, unode, trb)
}

func (this *timservice) addgroup(bs []byte, ws *tlnet.Websocket) (err sys.ERROR) {
	tr := newTimReq(bs)
	if tr == nil || tr.Node == nil || !checkNode(*tr.Node) || tr.ReqStr == nil {
		return sys.ERR_PARAMS
	}
	if wss, b := wsware.Get(ws); b {
		var gtype int8
		if gtype, err = data.Handler.GroupGtype(*tr.Node, wss.tid.Domain); err == nil {
			switch gtype {
			case sys.GROUP_OPEN:
				if err = data.Handler.Addgroup(*tr.Node, wss.tid.Node, wss.tid.Domain); err == nil {
					t := int64(sys.BUSINESS_ADDROOM)
					wsware.SendWs(ws.Id, &TimAck{Ok: true, TimType: int8(sys.TIMBUSINESS), N: tr.Node, T: &t}, sys.TIMACK)
				}
			case sys.GROUP_PRIVATE:
				if err = data.Handler.Addgroup(*tr.Node, wss.tid.Node, wss.tid.Domain); err != nil {
					return
				}
				var ms []string
				if ms, err = data.Handler.GroupManagers(*tr.Node, wss.tid.Domain); err == nil {
					for _, u := range ms {
						id, _t := RandId(), time.Now().UnixNano()
						tm := &TimMessage{ID: &id, FromTid: wss.tid, ToTid: &Tid{Node: u}, BnType: &sys.BUSINESS_ADDROOM, RoomTid: &Tid{Node: *tr.Node}, MsType: sys.SOURCE_USER, OdType: sys.ORDER_BUSINESS, DataString: tr.ReqStr, Timestamp: &_t}
						sys.TimMessageProcessor(tm, sys.TRANS_SOURCE)
					}
				}
			}
		}
	}
	return
}

func (this *timservice) pullgroup(bs []byte, ws *tlnet.Websocket) (err sys.ERROR) {
	tr := newTimReq(bs)
	if tr == nil || tr.Node == nil || tr.Node2 == nil || !checkNode(*tr.Node) || !checkNode(*tr.Node2) {
		return sys.ERR_PARAMS
	}
	if wss, b := wsware.Get(ws); b {
		if wss.tid.Node == *tr.Node2 {
			return sys.ERR_AUTH
		}
		var isreq bool
		if isreq, err = data.Handler.Pullgroup(*tr.Node, wss.tid.Node, *tr.Node2, wss.tid.Domain); err == nil {
			id, _t := RandId(), time.Now().UnixNano()
			tm := &TimMessage{ID: &id, FromTid: wss.tid, ToTid: &Tid{Node: *tr.Node2}, RoomTid: &Tid{Node: *tr.Node}, MsType: sys.SOURCE_USER, OdType: sys.ORDER_BUSINESS, Timestamp: &_t}
			if isreq {
				tm.BnType = &sys.BUSINESS_PASSROOM
			} else {
				tm.BnType = &sys.BUSINESS_PULLROOM
			}
			sys.TimMessageProcessor(tm, sys.TRANS_SOURCE)
			t := int64(sys.BUSINESS_PULLROOM)
			wsware.SendWs(ws.Id, &TimAck{Ok: true, TimType: int8(sys.TIMBUSINESS), N: tr.Node, T: &t, N2: tr.Node2}, sys.TIMACK)
		}
	}
	return
}

func (this *timservice) nopassgroup(bs []byte, ws *tlnet.Websocket) (err sys.ERROR) {
	tr := newTimReq(bs)
	if tr == nil || tr.Node == nil || tr.Node2 == nil || !checkNode(*tr.Node) || !checkNode(*tr.Node2) {
		return sys.ERR_PARAMS
	}
	if wss, b := wsware.Get(ws); b {
		if err = data.Handler.Nopassgroup(*tr.Node, wss.tid.Node, *tr.Node2, wss.tid.Domain); err == nil {
			id, _t := RandId(), time.Now().UnixNano()
			tm := &TimMessage{ID: &id, FromTid: wss.tid, ToTid: &Tid{Node: *tr.Node2}, RoomTid: &Tid{Node: *tr.Node}, MsType: sys.SOURCE_USER, OdType: sys.ORDER_BUSINESS, Timestamp: &_t}
			tm.BnType = &sys.BUSINESS_NOPASSROOM
			sys.TimMessageProcessor(tm, sys.TRANS_SOURCE)
			t := int64(sys.BUSINESS_NOPASSROOM)
			wsware.SendWs(ws.Id, &TimAck{Ok: true, TimType: int8(sys.TIMBUSINESS), N: tr.Node, T: &t, N2: tr.Node2}, sys.TIMACK)
		}
	}
	return
}

func (this *timservice) kickgroup(bs []byte, ws *tlnet.Websocket) (err sys.ERROR) {
	tr := newTimReq(bs)
	if tr == nil || tr.Node == nil || tr.Node2 == nil || !checkNode(*tr.Node) || !checkNode(*tr.Node2) {
		return sys.ERR_PARAMS
	}
	if wss, b := wsware.Get(ws); b {
		if wss.tid.Node == *tr.Node2 {
			return sys.ERR_AUTH
		}
		if err = data.Handler.Kickgroup(*tr.Node, wss.tid.Node, *tr.Node2, wss.tid.Domain); err == nil {
			id, _t := RandId(), time.Now().UnixNano()
			tm := &TimMessage{ID: &id, BnType: &sys.BUSINESS_KICKROOM, FromTid: wss.tid, ToTid: &Tid{Node: *tr.Node2}, RoomTid: &Tid{Node: *tr.Node}, MsType: sys.SOURCE_USER, OdType: sys.ORDER_BUSINESS, Timestamp: &_t}
			sys.TimMessageProcessor(tm, sys.TRANS_SOURCE)
			t := int64(sys.BUSINESS_KICKROOM)
			wsware.SendWs(ws.Id, &TimAck{Ok: true, TimType: int8(sys.TIMBUSINESS), N: tr.Node, T: &t, N2: tr.Node2}, sys.TIMACK)
		}
	}
	return
}

func (this *timservice) leavegroup(bs []byte, ws *tlnet.Websocket) (err sys.ERROR) {
	tr := newTimReq(bs)
	if tr == nil || tr.Node == nil || !checkNode(*tr.Node) {
		return sys.ERR_PARAMS
	}
	if wss, b := wsware.Get(ws); b {
		if err = data.Handler.Leavegroup(*tr.Node, wss.tid.Node, wss.tid.Domain); err == nil {
			t := int64(sys.BUSINESS_LEAVEROOM)
			wsware.SendWs(ws.Id, &TimAck{Ok: true, TimType: int8(sys.TIMBUSINESS), N: tr.Node, T: &t}, sys.TIMACK)
		}
	}
	return
}

func (this *timservice) cancelgroup(bs []byte, ws *tlnet.Websocket) (err sys.ERROR) {
	tr := newTimReq(bs)
	if tr == nil || tr.Node == nil || !checkNode(*tr.Node) {
		return sys.ERR_PARAMS
	}
	if wss, b := wsware.Get(ws); b {
		if err = data.Handler.Cancelgroup(*tr.Node, wss.tid.Node, wss.tid.Domain); err == nil {
			t := int64(sys.BUSINESS_CANCELROOM)
			wsware.SendWs(ws.Id, &TimAck{Ok: true, TimType: int8(sys.TIMBUSINESS), N: tr.Node, T: &t}, sys.TIMACK)
		}
	}
	return
}

func (this *timservice) blockgroup(bs []byte, ws *tlnet.Websocket) (err sys.ERROR) {
	tr := newTimReq(bs)
	if tr == nil || tr.Node == nil || !checkNode(*tr.Node) {
		return sys.ERR_PARAMS
	}
	if wss, b := wsware.Get(ws); b {
		if err = data.Handler.Blockgroup(*tr.Node, wss.tid.Node, wss.tid.Domain); err == nil {
			t := int64(sys.BUSINESS_BLOCKROOM)
			wsware.SendWs(ws.Id, &TimAck{Ok: true, TimType: int8(sys.TIMBUSINESS), N: tr.Node, T: &t}, sys.TIMACK)
		}
	}
	return
}

func (this *timservice) blockgroupmember(bs []byte, ws *tlnet.Websocket) (err sys.ERROR) {
	tr := newTimReq(bs)
	if tr == nil || tr.Node == nil || !checkNode(*tr.Node) || tr.Node2 == nil || !checkNode(*tr.Node2) {
		return sys.ERR_PARAMS
	}
	if wss, b := wsware.Get(ws); b {
		if wss.tid.Node == *tr.Node2 {
			return sys.ERR_AUTH
		}
		if err = data.Handler.Blockgroupmember(*tr.Node, wss.tid.Node, *tr.Node2, wss.tid.Domain); err == nil {
			t := int64(sys.BUSINESS_BLOCKROOMMEMBER)
			wsware.SendWs(ws.Id, &TimAck{Ok: true, TimType: int8(sys.TIMBUSINESS), N: tr.Node, T: &t, N2: tr.Node2}, sys.TIMACK)
		}
	}
	return
}

func (this *timservice) blockrosterlist(ws *tlnet.Websocket) (err sys.ERROR) {
	if wss, b := wsware.Get(ws); b {
		if ss := data.Handler.Blockrosterlist(wss.tid.Node); ss != nil && len(ss) > 0 {
			wsware.SendWs(ws.Id, &TimNodes{Ntype: sys.NODEINFO_BLOCKROSTERLIST, Nodelist: ss}, sys.TIMNODES)
		}
	}
	return
}

func (this *timservice) blockroomlist(ws *tlnet.Websocket) (err sys.ERROR) {
	if wss, b := wsware.Get(ws); b {
		if ss := data.Handler.Blockroomlist(wss.tid.Node); ss != nil && len(ss) > 0 {
			wsware.SendWs(ws.Id, &TimNodes{Ntype: sys.NODEINFO_BLOCKROOMLIST, Nodelist: ss}, sys.TIMNODES)
		}
	}
	return
}

func (this *timservice) blockroommemberlist(bs []byte, ws *tlnet.Websocket) (err sys.ERROR) {
	tr := newTimReq(bs)
	if tr == nil || tr.Node == nil || !checkNode(*tr.Node) {
		return sys.ERR_PARAMS
	}
	if wss, b := wsware.Get(ws); b {
		if ss := data.Handler.Blockroommemberlist(*tr.Node, wss.tid.Node); ss != nil && len(ss) > 0 {
			wsware.SendWs(ws.Id, &TimNodes{Ntype: sys.NODEINFO_BLOCKROOMMEMBERLIST, Nodelist: ss}, sys.TIMNODES)
		}
	}
	return
}

func (this *timservice) sysmodifyauth(account, pwd string, domain *string) sys.ERROR {
	return data.Handler.Modify(util.CreateUUID(account, domain), nil, pwd, domain)
}

func (this *timservice) modifyauth(bs []byte, ws *tlnet.Websocket) (err sys.ERROR) {
	tr := newTimReq(bs)
	if tr == nil || tr.ReqStr == nil || tr.ReqStr2 == nil {
		return sys.ERR_MODIFYAUTH
	}
	if wss, b := wsware.Get(ws); b {
		t := int64(sys.BUSINESS_MODIFYAUTH)
		if err = data.Handler.Modify(util.NodeToUUID(wss.tid.Node), tr.ReqStr, *tr.ReqStr2, wss.tid.Domain); err == nil {
			wsware.SendWs(ws.Id, &TimAck{Ok: true, TimType: int8(sys.TIMBUSINESS), T: &t}, sys.TIMACK)
		} else {
			wsware.SendWs(ws.Id, &TimAck{Ok: false, Error: sys.ERR_MODIFYAUTH.TimError(), TimType: int8(sys.TIMBUSINESS), T: &t}, sys.TIMACK)
		}
	}
	return
}

func (this *timservice) osuserbean(node string, tb *TimUserBean) (err sys.ERROR) {
	return data.Handler.ModifyUserInfo(node, tb)
}

func (this *timservice) nodeinfo(bs []byte, ws *tlnet.Websocket) (err sys.ERROR) {
	defer util.Recover()
	sys.Stat.TxDo()
	defer sys.Stat.TxDone()
	tr := newTimNodes(bs)
	if tr == nil {
		return sys.ERR_PARAMS
	}
	if wss, b := wsware.Get(ws); b {
		switch tr.Ntype {
		case sys.NODEINFO_USERINFO:
			if m, e := data.Handler.GetUserInfo(tr.Nodelist); e == nil && m != nil && len(m) > 0 {
				tr.Nodelist = nil
				tr.Usermap = m
				wsware.SendWs(ws.Id, tr, sys.TIMNODES)
			} else {
				err = e
			}
		case sys.NODEINFO_ROOMINFO:
			if m, e := data.Handler.GetGroupInfo(tr.Nodelist); e == nil && m != nil && len(m) > 0 {
				tr.Nodelist = nil
				tr.Roommap = m
				wsware.SendWs(ws.Id, tr, sys.TIMNODES)
			} else {
				err = e
			}
		case sys.NODEINFO_MODIFYUSER:
			if tr.Usermap != nil && len(tr.Usermap) == 1 {
				for _, v := range tr.Usermap {
					if err = this.osuserbean(wss.tid.Node, v); err == nil {
						t := int64(sys.NODEINFO_MODIFYUSER)
						wsware.SendWs(ws.Id, &TimAck{Ok: true, TimType: int8(sys.TIMNODES), T: &t}, sys.TIMACK)
					}
				}
			} else {
				return sys.ERR_PARAMS
			}
		case sys.NODEINFO_MODIFYROOM:
			if tr.Roommap != nil && len(tr.Roommap) == 1 {
				for k, v := range tr.Roommap {
					if err = data.Handler.ModifygroupInfo(k, wss.tid.Node, v); err == nil {
						t := int64(sys.NODEINFO_MODIFYROOM)
						wsware.SendWs(ws.Id, &TimAck{Ok: true, TimType: int8(sys.TIMNODES), T: &t}, sys.TIMACK)
					}
				}
			} else {
				return sys.ERR_PARAMS
			}
		default:
			err = sys.ERR_PARAMS
		}
	}
	return
}
