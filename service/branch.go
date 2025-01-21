// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package service

import (
	"github.com/donnie4w/tim/errs"
	"time"

	goutil "github.com/donnie4w/gofer/util"
	"github.com/donnie4w/tim/data"
	. "github.com/donnie4w/tim/stub"
	"github.com/donnie4w/tim/sys"
	"github.com/donnie4w/tim/util"
	"github.com/donnie4w/tlnet"
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
		switch sys.BUSINESSTYPE(*tr.Rtype) {
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
			return tss.modifyauth(bs, ws)
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
			wsware.SendWs(ws.Id, &TimNodes{Ntype: int32(sys.NODEINFO_ROSTER), Nodelist: _r}, sys.TIMNODES)
		}
	}
	return
}

func (tss *timservice) usergroup(ws *tlnet.Websocket) (err errs.ERROR) {
	if wss, ok := wsware.Get(ws); ok {
		tid := wss.tid
		if _r := data.Service.UserGroup(tid.Node, tid.Domain); len(_r) > 0 {
			wsware.SendWs(ws.Id, &TimNodes{Ntype: int32(sys.NODEINFO_ROOM), Nodelist: _r}, sys.TIMNODES)
		}
	}
	return
}

func (tss *timservice) grouproster(bs []byte, ws *tlnet.Websocket) (err errs.ERROR) {
	if wss, ok := wsware.Get(ws); ok {
		tr := newTimReq(bs)
		node := *tr.Node
		if !checkNode(node) || !existGroup(&Tid{Node: node, Domain: wss.tid.Domain}) {
			return errs.ERR_ACCOUNT
		}
		if ok := AuthGroup(node, wss.tid.Node, wss.tid.Domain); ok {
			if gs := data.Service.GroupRoster(node); len(gs) > 0 {
				wsware.SendWs(ws.Id, &TimNodes{Ntype: int32(sys.NODEINFO_ROOMMEMBER), Nodelist: gs, Node: &node}, sys.TIMNODES)
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
	if !checkNode(*tq.Node) {
		return errs.ERR_ACCOUNT
	}
	if wss, b := wsware.Get(ws); b {
		if wss.tid.Node == *tq.Node {
			return errs.ERR_PERM_DENIED
		}

		id, t := goutil.UUID64(), time.Now().UnixNano()
		tm := &TimMessage{ID: &id, MsType: sys.SOURCE_USER, OdType: sys.ORDER_BUSINESS, ToTid: &Tid{Node: *tq.Node}, FromTid: wss.tid, DataString: tq.ReqStr, Timestamp: &t}

		var status int8
		if status, err = data.Service.Addroster(wss.tid.Node, *tq.Node, wss.tid.Domain); err == nil {
			//if status == 0x10|0x01 {
			//	bt := int32(sys.BUSINESS_FRIEND)
			//	tm.BnType = &bt
			//	wsware.SendWs(ws.Id, tm, sys.TIMMESSAGE)
			//} else {
			bt := int32(sys.BUSINESS_ADDROSTER)
			tm.BnType = &bt
			//}
			sys.TimMessageProcessor(tm, sys.TRANS_SOURCE)
		} else if status == 0x11 {
			bt := int32(sys.BUSINESS_FRIEND)
			tm.BnType = &bt
			wsware.SendWs(ws.Id, tm, sys.TIMMESSAGE)
		}
	}
	return
}

func (tss *timservice) rmroster(bs []byte, ws *tlnet.Websocket) (err errs.ERROR) {
	tr := newTimReq(bs)
	if tr == nil || tr.Node == nil {
		return errs.ERR_PARAMS
	}
	if !checkNode(*tr.Node) {
		return errs.ERR_ACCOUNT
	}
	if wss, b := wsware.Get(ws); b {
		if ms, ok := data.Service.Rmroster(wss.tid.Node, *tr.Node, wss.tid.Domain); !ok {
			return errs.ERR_PARAMS
		} else {
			if ms {
				id, t, bt := goutil.UUID64(), time.Now().UnixNano(), int32(sys.BUSINESS_REMOVEROSTER)
				tm := &TimMessage{ID: &id, MsType: sys.SOURCE_USER, BnType: &bt, OdType: sys.ORDER_BUSINESS, FromTid: wss.tid, ToTid: &Tid{Node: *tr.Node}, Timestamp: &t}
				sys.TimMessageProcessor(tm, sys.TRANS_SOURCE)
			}
			t := int64(sys.BUSINESS_REMOVEROSTER)
			wsware.SendWs(ws.Id, &TimAck{Ok: true, TimType: int8(sys.TIMBUSINESS), N: tr.Node, T: &t}, sys.TIMACK)
		}
	}
	return
}

func (tss *timservice) blockroster(bs []byte, ws *tlnet.Websocket) (err errs.ERROR) {
	tq := newTimReq(bs)
	if tq == nil || tq.Node == nil {
		return errs.ERR_PARAMS
	}
	if !checkNode(*tq.Node) {
		return errs.ERR_ACCOUNT
	}
	if wss, b := wsware.Get(ws); b {
		if wss.tid.Node == *tq.Node {
			return errs.ERR_PERM_DENIED
		}
		if ms, ok := data.Service.Blockroster(wss.tid.Node, *tq.Node, wss.tid.Domain); !ok {
			return errs.ERR_PARAMS
		} else {
			if ms {
				id, _t, bt := goutil.UUID64(), time.Now().UnixNano(), int32(sys.BUSINESS_BLOCKROSTER)
				tm := &TimMessage{ID: &id, MsType: sys.SOURCE_USER, OdType: sys.ORDER_BUSINESS, ToTid: &Tid{Node: *tq.Node}, FromTid: wss.tid, BnType: &bt, Timestamp: &_t}
				sys.TimMessageProcessor(tm, sys.TRANS_SOURCE)
			}
			t := int64(sys.BUSINESS_BLOCKROSTER)
			wsware.SendWs(ws.Id, &TimAck{Ok: true, TimType: int8(sys.TIMBUSINESS), N: tq.Node, T: &t}, sys.TIMACK)
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
		if sys.TIMTYPE(*tr.ReqInt) != rtype {
			rtype = sys.GROUP_OPEN
		}
		var gnode string
		if gnode, err = data.Service.Newgroup(wss.tid.Node, *tr.Node, rtype, wss.tid.Domain); err == nil {
			t := int64(sys.BUSINESS_NEWROOM)
			wsware.SendWs(ws.Id, &TimAck{Ok: true, TimType: int8(sys.TIMBUSINESS), N: &gnode, T: &t}, sys.TIMACK)
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

func (tss *timservice) osModifygroupInfo(unode, gnode string, trb *TimRoomBean) errs.ERROR {
	if !checkNode(gnode) || !checkNode(unode) {
		return errs.ERR_ACCOUNT
	}
	return data.Service.ModifygroupInfo(gnode, unode, trb, true)
}

func (tss *timservice) addgroup(bs []byte, ws *tlnet.Websocket) (err errs.ERROR) {
	tr := newTimReq(bs)
	if tr == nil || tr.Node == nil || !checkNode(*tr.Node) || tr.ReqStr == nil {
		return errs.ERR_PARAMS
	}
	if wss, b := wsware.Get(ws); b {
		var gtype int8
		if gtype, err = data.Service.GroupGtype(*tr.Node, wss.tid.Domain); err == nil {
			switch sys.TIMTYPE(gtype) {
			case sys.GROUP_OPEN:
				if err = data.Service.Addgroup(*tr.Node, wss.tid.Node, wss.tid.Domain); err == nil {
					t := int64(sys.BUSINESS_ADDROOM)
					wsware.SendWs(ws.Id, &TimAck{Ok: true, TimType: int8(sys.TIMBUSINESS), N: tr.Node, T: &t}, sys.TIMACK)
				}
			case sys.GROUP_PRIVATE:
				if err = data.Service.Addgroup(*tr.Node, wss.tid.Node, wss.tid.Domain); err != nil {
					return
				}
				var ms []string
				if ms, err = data.Service.GroupManagers(*tr.Node, wss.tid.Domain); err == nil {
					for _, u := range ms {
						id, _t, bt := goutil.UUID64(), time.Now().UnixNano(), int32(sys.BUSINESS_ADDROOM)
						tm := &TimMessage{ID: &id, FromTid: wss.tid, ToTid: &Tid{Node: u}, BnType: &bt, RoomTid: &Tid{Node: *tr.Node}, MsType: sys.SOURCE_USER, OdType: sys.ORDER_BUSINESS, DataString: tr.ReqStr, Timestamp: &_t}
						sys.TimMessageProcessor(tm, sys.TRANS_SOURCE)
					}
				}
			}
		}
	}
	return
}

func (tss *timservice) pullgroup(bs []byte, ws *tlnet.Websocket) (err errs.ERROR) {
	tr := newTimReq(bs)
	if tr == nil || tr.Node == nil || tr.Node2 == nil || !checkNode(*tr.Node) || !checkNode(*tr.Node2) {
		return errs.ERR_PARAMS
	}
	if wss, b := wsware.Get(ws); b {
		if wss.tid.Node == *tr.Node2 {
			return errs.ERR_PERM_DENIED
		}
		var isreq bool
		if isreq, err = data.Service.Pullgroup(*tr.Node, wss.tid.Node, *tr.Node2, wss.tid.Domain); err == nil {
			id, _t := goutil.UUID64(), time.Now().UnixNano()
			tm := &TimMessage{ID: &id, FromTid: wss.tid, ToTid: &Tid{Node: *tr.Node2}, RoomTid: &Tid{Node: *tr.Node}, MsType: sys.SOURCE_USER, OdType: sys.ORDER_BUSINESS, Timestamp: &_t}
			if isreq {
				bt := int32(sys.BUSINESS_PASSROOM)
				tm.BnType = &bt
			} else {
				bt := int32(sys.BUSINESS_PULLROOM)
				tm.BnType = &bt
			}
			sys.TimMessageProcessor(tm, sys.TRANS_SOURCE)
			t := int64(sys.BUSINESS_PULLROOM)
			wsware.SendWs(ws.Id, &TimAck{Ok: true, TimType: int8(sys.TIMBUSINESS), N: tr.Node, T: &t, N2: tr.Node2}, sys.TIMACK)
		}
	}
	return
}

func (tss *timservice) nopassgroup(bs []byte, ws *tlnet.Websocket) (err errs.ERROR) {
	tr := newTimReq(bs)
	if tr == nil || tr.Node == nil || tr.Node2 == nil || !checkNode(*tr.Node) || !checkNode(*tr.Node2) {
		return errs.ERR_PARAMS
	}
	if wss, b := wsware.Get(ws); b {
		if err = data.Service.Nopassgroup(*tr.Node, wss.tid.Node, *tr.Node2, wss.tid.Domain); err == nil {
			id, _t, bt := goutil.UUID64(), time.Now().UnixNano(), int32(sys.BUSINESS_NOPASSROOM)
			tm := &TimMessage{ID: &id, FromTid: wss.tid, ToTid: &Tid{Node: *tr.Node2}, RoomTid: &Tid{Node: *tr.Node}, MsType: sys.SOURCE_USER, OdType: sys.ORDER_BUSINESS, Timestamp: &_t}
			tm.BnType = &bt
			sys.TimMessageProcessor(tm, sys.TRANS_SOURCE)
			t := int64(sys.BUSINESS_NOPASSROOM)
			wsware.SendWs(ws.Id, &TimAck{Ok: true, TimType: int8(sys.TIMBUSINESS), N: tr.Node, T: &t, N2: tr.Node2}, sys.TIMACK)
		}
	}
	return
}

func (tss *timservice) kickgroup(bs []byte, ws *tlnet.Websocket) (err errs.ERROR) {
	tr := newTimReq(bs)
	if tr == nil || tr.Node == nil || tr.Node2 == nil || !checkNode(*tr.Node) || !checkNode(*tr.Node2) {
		return errs.ERR_PARAMS
	}
	if wss, b := wsware.Get(ws); b {
		if wss.tid.Node == *tr.Node2 {
			return errs.ERR_PERM_DENIED
		}
		if err = data.Service.Kickgroup(*tr.Node, wss.tid.Node, *tr.Node2, wss.tid.Domain); err == nil {
			id, _t, bt := goutil.UUID64(), time.Now().UnixNano(), int32(sys.BUSINESS_KICKROOM)
			tm := &TimMessage{ID: &id, BnType: &bt, FromTid: wss.tid, ToTid: &Tid{Node: *tr.Node2}, RoomTid: &Tid{Node: *tr.Node}, MsType: sys.SOURCE_USER, OdType: sys.ORDER_BUSINESS, Timestamp: &_t}
			sys.TimMessageProcessor(tm, sys.TRANS_SOURCE)
			t := int64(sys.BUSINESS_KICKROOM)
			wsware.SendWs(ws.Id, &TimAck{Ok: true, TimType: int8(sys.TIMBUSINESS), N: tr.Node, T: &t, N2: tr.Node2}, sys.TIMACK)
		}
	}
	return
}

func (tss *timservice) leavegroup(bs []byte, ws *tlnet.Websocket) (err errs.ERROR) {
	tr := newTimReq(bs)
	if tr == nil || tr.Node == nil || !checkNode(*tr.Node) {
		return errs.ERR_PARAMS
	}
	if wss, b := wsware.Get(ws); b {
		if err = data.Service.Leavegroup(*tr.Node, wss.tid.Node, wss.tid.Domain); err == nil {
			t := int64(sys.BUSINESS_LEAVEROOM)
			wsware.SendWs(ws.Id, &TimAck{Ok: true, TimType: int8(sys.TIMBUSINESS), N: tr.Node, T: &t}, sys.TIMACK)
		}
	}
	return
}

func (tss *timservice) cancelgroup(bs []byte, ws *tlnet.Websocket) (err errs.ERROR) {
	tr := newTimReq(bs)
	if tr == nil || tr.Node == nil || !checkNode(*tr.Node) {
		return errs.ERR_PARAMS
	}
	if wss, b := wsware.Get(ws); b {
		if err = data.Service.Cancelgroup(*tr.Node, wss.tid.Node, wss.tid.Domain); err == nil {
			t := int64(sys.BUSINESS_CANCELROOM)
			wsware.SendWs(ws.Id, &TimAck{Ok: true, TimType: int8(sys.TIMBUSINESS), N: tr.Node, T: &t}, sys.TIMACK)
		}
	}
	return
}

func (tss *timservice) blockgroup(bs []byte, ws *tlnet.Websocket) (err errs.ERROR) {
	tr := newTimReq(bs)
	if tr == nil || tr.Node == nil || !checkNode(*tr.Node) {
		return errs.ERR_PARAMS
	}
	if wss, b := wsware.Get(ws); b {
		if err = data.Service.Blockgroup(*tr.Node, wss.tid.Node, wss.tid.Domain); err == nil {
			t := int64(sys.BUSINESS_BLOCKROOM)
			wsware.SendWs(ws.Id, &TimAck{Ok: true, TimType: int8(sys.TIMBUSINESS), N: tr.Node, T: &t}, sys.TIMACK)
		}
	}
	return
}

func (tss *timservice) blockgroupmember(bs []byte, ws *tlnet.Websocket) (err errs.ERROR) {
	tr := newTimReq(bs)
	if tr == nil || tr.Node == nil || !checkNode(*tr.Node) || tr.Node2 == nil || !checkNode(*tr.Node2) {
		return errs.ERR_PARAMS
	}
	if wss, b := wsware.Get(ws); b {
		if wss.tid.Node == *tr.Node2 {
			return errs.ERR_PERM_DENIED
		}
		if err = data.Service.Blockgroupmember(*tr.Node, wss.tid.Node, *tr.Node2, wss.tid.Domain); err == nil {
			t := int64(sys.BUSINESS_BLOCKROOMMEMBER)
			wsware.SendWs(ws.Id, &TimAck{Ok: true, TimType: int8(sys.TIMBUSINESS), N: tr.Node, T: &t, N2: tr.Node2}, sys.TIMACK)
		}
	}
	return
}

func (tss *timservice) blockrosterlist(ws *tlnet.Websocket) (err errs.ERROR) {
	if wss, b := wsware.Get(ws); b {
		if ss := data.Service.Blockrosterlist(wss.tid.Node); len(ss) > 0 {
			wsware.SendWs(ws.Id, &TimNodes{Ntype: int32(sys.NODEINFO_BLOCKROSTERLIST), Nodelist: ss}, sys.TIMNODES)
		}
	}
	return
}

func (tss *timservice) blockroomlist(ws *tlnet.Websocket) (err errs.ERROR) {
	if wss, b := wsware.Get(ws); b {
		if ss := data.Service.Blockroomlist(wss.tid.Node); len(ss) > 0 {
			wsware.SendWs(ws.Id, &TimNodes{Ntype: int32(sys.NODEINFO_BLOCKROOMLIST), Nodelist: ss}, sys.TIMNODES)
		}
	}
	return
}

func (tss *timservice) blockroommemberlist(bs []byte, ws *tlnet.Websocket) (err errs.ERROR) {
	tr := newTimReq(bs)
	if tr == nil || tr.Node == nil || !checkNode(*tr.Node) {
		return errs.ERR_PARAMS
	}
	if wss, b := wsware.Get(ws); b {
		if ss := data.Service.Blockroommemberlist(*tr.Node, wss.tid.Node); len(ss) > 0 {
			wsware.SendWs(ws.Id, &TimNodes{Ntype: int32(sys.NODEINFO_BLOCKROOMMEMBERLIST), Nodelist: ss}, sys.TIMNODES)
		}
	}
	return
}

func (tss *timservice) sysmodify(node string, oldpwd *string, newpwd string, domain *string) errs.ERROR {
	return data.Service.Modify(util.NodeToUUID(node), oldpwd, newpwd, domain)
}

func (tss *timservice) modifyauth(bs []byte, ws *tlnet.Websocket) (err errs.ERROR) {
	tr := newTimReq(bs)
	if tr == nil || tr.ReqStr == nil || tr.ReqStr2 == nil {
		return errs.ERR_MODIFYAUTH
	}
	if wss, b := wsware.Get(ws); b {
		t := int64(sys.BUSINESS_MODIFYAUTH)
		if err = data.Service.Modify(util.NodeToUUID(wss.tid.Node), tr.ReqStr, *tr.ReqStr2, wss.tid.Domain); err == nil {
			wsware.SendWs(ws.Id, &TimAck{Ok: true, TimType: int8(sys.TIMBUSINESS), T: &t}, sys.TIMACK)
		} else {
			wsware.SendWs(ws.Id, &TimAck{Ok: false, Error: errs.ERR_MODIFYAUTH.TimError(), TimType: int8(sys.TIMBUSINESS), T: &t}, sys.TIMACK)
		}
	}
	return
}

func (tss *timservice) osuserbean(node string, tb *TimUserBean) (err errs.ERROR) {
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
			if m, e := data.Service.GetUserInfo(tr.Nodelist); e == nil && m != nil && len(m) > 0 {
				tr.Nodelist = nil
				tr.Usermap = m
				wsware.SendWs(ws.Id, tr, sys.TIMNODES)
			} else {
				err = e
			}
		case sys.NODEINFO_ROOMINFO:
			if m, e := data.Service.GetGroupInfo(tr.Nodelist); e == nil && m != nil && len(m) > 0 {
				tr.Nodelist = nil
				tr.Roommap = m
				wsware.SendWs(ws.Id, tr, sys.TIMNODES)
			} else {
				err = e
			}
		case sys.NODEINFO_MODIFYUSER:
			if tr.Usermap != nil && len(tr.Usermap) == 1 {
				for _, v := range tr.Usermap {
					if err = tss.osuserbean(wss.tid.Node, v); err == nil {
						t := int64(sys.NODEINFO_MODIFYUSER)
						wsware.SendWs(ws.Id, &TimAck{Ok: true, TimType: int8(sys.TIMNODES), T: &t}, sys.TIMACK)
					}
				}
			} else {
				return errs.ERR_PARAMS
			}
		case sys.NODEINFO_MODIFYROOM:
			if tr.Roommap != nil && len(tr.Roommap) == 1 {
				for k, v := range tr.Roommap {
					if err = data.Service.ModifygroupInfo(k, wss.tid.Node, v, false); err == nil {
						t := int64(sys.NODEINFO_MODIFYROOM)
						wsware.SendWs(ws.Id, &TimAck{Ok: true, TimType: int8(sys.TIMNODES), T: &t}, sys.TIMACK)
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
