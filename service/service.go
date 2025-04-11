// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package service

import (
	"bytes"
	"github.com/donnie4w/gofer/base58"
	goutil "github.com/donnie4w/gofer/util"
	"github.com/donnie4w/tim/amr"
	"github.com/donnie4w/tim/cache"
	"github.com/donnie4w/tim/data"
	"github.com/donnie4w/tim/errs"
	"github.com/donnie4w/tim/stub"
	"github.com/donnie4w/tim/sys"
	"github.com/donnie4w/tim/util"
	"github.com/donnie4w/tim/vgate"
	"github.com/donnie4w/tlnet"
	"sort"
	"strconv"
	"strings"
	"time"
)

var service = &timservice{}

type timservice struct{}

func (tss *timservice) osregister(name, pwd string, domain *string) (node string, e errs.ERROR) {
	defer util.Recover()
	sys.Stat.TxDo()
	defer sys.Stat.TxDone()
	return data.Service.Register(name, pwd, domain)
}

func (tss *timservice) register(bs []byte) (node string, e errs.ERROR) {
	defer util.Recover()
	sys.Stat.TxDo()
	defer sys.Stat.TxDone()
	var ta *stub.TimAuth
	var err error
	if util.JTP(bs[0]) {
		ta, err = goutil.JsonDecode[*stub.TimAuth](bs[1:])
	} else {
		ta, err = goutil.TDecode(bs[1:], &stub.TimAuth{})
	}
	if err == nil {
		node, e = data.Service.Register(ta.GetName(), ta.GetPwd(), ta.Domain)
	} else {
		e = errs.ERR_PARAMS
	}
	return
}

func (tss *timservice) ping(ws *tlnet.Websocket) (err errs.ERROR) {
	defer util.Recover()
	sys.Stat.TxDo()
	defer sys.Stat.TxDone()
	wsware.Ping(ws.Id)
	return
}

func (tss *timservice) ack(bs []byte) (err errs.ERROR) {
	defer util.Recover()
	sys.Stat.TxDo()
	defer sys.Stat.TxDone()
	awaitEnd(bs[1:])
	return
}

func (tss *timservice) ostoken(nodeorname string, password, domain *string) (_r string, _n string, e errs.ERROR) {
	defer util.Recover()
	sys.Stat.TxDo()
	defer sys.Stat.TxDone()
	if password != nil {
		var node string
		if node, e = data.Service.AuthNode(nodeorname, *password, domain); e == nil {
			tid := &stub.Tid{Node: node, Domain: domain}
			_r, _n = token(), node
			amr.PutToken(_r, tid)
		}
	} else {
		switch sys.GetDBMOD() {
		case sys.NODB, sys.EXTERNALDB:
			_r, _n = token(), nodeorname
		case sys.TLDB, sys.INLINEDB:
			if !existUser(&stub.Tid{Node: nodeorname, Domain: domain}) {
				return _r, "", errs.ERR_NOEXIST
			}
			_r, _n = token(), util.UUIDToNode(util.CreateUUID(nodeorname, domain))
		default:
			e = errs.ERR_DATABASE
		}
		tid := &stub.Tid{Node: _n, Domain: domain}
		cache.TokenCache.Put(_r, tid)
	}
	return
}

func (tss *timservice) token(bs []byte) (_r string, e errs.ERROR) {
	defer util.Recover()
	sys.Stat.TxDo()
	defer sys.Stat.TxDone()
	ta := newAuth(bs)
	if ta == nil {
		return _r, errs.ERR_PARAMS
	}
	var node string
	if node, e = data.Service.AuthNode(ta.GetName(), ta.GetPwd(), ta.Domain); e == nil {
		tid := &stub.Tid{Node: node, Domain: ta.Domain, Extend: ta.Extend}
		_r = token()
		cache.TokenCache.Put(_r, tid)
	}
	return
}

func (tss *timservice) auth(bs []byte, ws *tlnet.Websocket) (e errs.ERROR) {
	defer util.Recover()
	sys.Stat.TxDo()
	defer sys.Stat.TxDone()
	if wsware.wsmap.Has(ws.Id) {
		wsware.SendWs(ws.Id, &stub.TimAck{Ok: true, TimType: int8(sys.TIMAUTH)}, sys.TIMACK)
		return
	}
	ta := newAuth(bs)
	if ta == nil {
		return errs.ERR_PARAMS
	}
	isAuth := false
	var tid *stub.Tid
	if ta.GetToken() != "" {
		if !util.CheckNode(ta.GetToken()) {
			return errs.ERR_TOKEN
		}
		if tid = amr.GetToken(ta.GetToken()); tid != nil {
			tid.Resource, tid.Termtyp = ta.Resource, ta.Termtyp
			amr.DelToken(ta.GetToken())
			if !isblock(tid.Node) {
				isAuth = true
			}
		} else {
			return errs.ERR_TOKEN
		}
	} else if ta.GetName() != "" && ta.GetPwd() != "" && !isblock(ta.GetName()) {
		if _r, err := data.Service.Login(ta.GetName(), ta.GetPwd(), ta.Domain); err == nil {
			tid = &stub.Tid{Node: _r, Domain: ta.Domain, Extend: ta.Extend, Resource: ta.Resource, Termtyp: ta.Termtyp}
			if !isblock(_r) {
				isAuth = true
			}
		}
	}
	if isAuth {
		overentry := true
		if deviceNums := wsware.deviceNums(tid.Node); deviceNums < sys.DeviceLimit {
			dtl := sys.CsDevice(tid.Node)
			if len(dtl)+deviceNums < sys.DeviceLimit {
				overentry = false
				if tid.Termtyp != nil {
					typebs := sys.DeviceTypeList(tid.Node)
					c := 0
					for _, u := range append(dtl, typebs...) {
						if u == byte(tid.GetTermtyp()) {
							c++
						}
					}
					if c > sys.DeviceTypeLimit {
						overentry = true
					}
				}
			}
		}
		if overentry {
			return errs.ERR_OVERENTRY
		}
		wsware.AddTid(ws, tid)
		if util.JTP(bs[0]) {
			wsware.SetJsonOn(ws)
		}
		maskuuid := int64(int32(goutil.FNVHash32(goutil.Int64ToBytes(sys.UUID))))
		maskWsId := string(base58.EncodeForInt64(uint64(util.MaskId(ws.Id))))
		wsware.SendWs(ws.Id, &stub.TimAck{Ok: true, N2: &maskWsId, T: &maskuuid, TimType: int8(sys.TIMAUTH), N: &tid.Node}, sys.TIMACK)
	} else {
		e = errs.ERR_PERM_DENIED
	}
	return
}

func sysMessage(nodelist []string, tm *stub.TimMessage) (err errs.ERROR) {
	defer util.Recover()
	sys.Stat.TxDo()
	defer sys.Stat.TxDone()
	if len(nodelist) == 0 && tm == nil {
		return errs.ERR_PARAMS
	}
	if checkList(nodelist) {
		tm.MsType, tm.OdType = sys.SOURCE_OS, sys.ORDER_INOF
		t := time.Now().UnixNano()
		for _, node := range nodelist {
			tm.ToTid = &stub.Tid{Node: node}
			if existUser(tm.ToTid) {
				tm.Timestamp = &t
				sys.TimMessageProcessor(tm, sys.TRANS_SOURCE)
			}
		}
	} else {
		return errs.ERR_ACCOUNT
	}
	return
}

func sysPresence(nodelist []string, tm *stub.TimPresence) (err errs.ERROR) {
	if len(nodelist) == 0 && tm == nil {
		return errs.ERR_PARAMS
	}
	if checkList(nodelist) {
		tm.ToList = nodelist
		tm.Offline = nil
		sys.TimPresenceProcessor(tm, sys.TRANS_SOURCE)
	} else {
		return errs.ERR_ACCOUNT
	}
	return
}

func (tss *timservice) pxMessage(tm *stub.TimMessage) (err errs.ERROR) {
	if tm.FromTid == nil || tm.FromTid.Node == "" {
		return errs.ERR_PARAMS
	}
	if tm.ToTid == nil && tm.ToList == nil && tm.RoomTid == nil {
		return errs.ERR_PARAMS
	}
	switch tm.MsType {
	case sys.SOURCE_USER:
		err = tss.messagehandler(fullTimMessage(tm), false)
	case sys.SOURCE_ROOM:
		err = tss.messagehandle(fullTimMessage(tm), false)
	default:
		err = errs.ERR_PARAMS
	}
	return
}

func (tss *timservice) bigString(bs []byte, ws *tlnet.Websocket) (_r errs.ERROR) {
	bigString := string(bs[5:])
	idx := strings.Index(bigString, sys.SEP_STR)
	dataString := bigString[idx+1:]
	if wss, b := wsware.Get(ws); b {
		_r = tss.messagehandler(&stub.TimMessage{MsType: 2, OdType: sys.ORDER_BIGSTRING, DataString: &dataString, FromTid: wss.tid, ToTid: &stub.Tid{Node: bigString[:idx]}}, true)
	}
	return
}

func (tss *timservice) bigBinary(bs []byte, ws *tlnet.Websocket) (_r errs.ERROR) {
	idx := bytes.IndexByte(bs[5:], sys.SEP_BIN)
	if wss, b := wsware.Get(ws); b {
		_r = tss.messagehandler(&stub.TimMessage{MsType: 2, OdType: sys.ORDER_BIGBINARY, DataBinary: bs[5:][idx+1:], FromTid: wss.tid, ToTid: &stub.Tid{Node: string(bs[5:][:idx])}}, true)
	}
	return
}

func (tss *timservice) bigBinaryStreamHandle(bs []byte, ws *tlnet.Websocket) (_r errs.ERROR) {
	idx := bytes.IndexByte(bs[5:], sys.SEP_BIN)
	if wss, b := wsware.Get(ws); b {
		t := &stub.TimStream{ID: goutil.UUID64(), VNode: string(bs[5:][:idx]), Body: bs[5:][idx+1:], FromNode: wss.tid.Node}
		return tss.streamhandler(t, ws)
	}
	return
}

func (tss *timservice) messageHandle(bs []byte, ws *tlnet.Websocket) (_r errs.ERROR) {
	return tss.wssMessage(newTimMessage(bs), ws)
}

func (tss *timservice) wssMessage(tm *stub.TimMessage, ws *tlnet.Websocket) (_r errs.ERROR) {
	defer util.Recover()
	sys.Stat.TxDo()
	defer sys.Stat.TxDone()
	if tm == nil {
		return errs.ERR_FORMAT
	}
	if tm.MsType == sys.SOURCE_OS {
		return errs.ERR_PARAMS
	}
	if !checkTid(tm.ToTid) || !checkTid(tm.RoomTid) {
		return errs.ERR_ACCOUNT
	}
	if wss, b := wsware.Get(ws); b {
		tm.FromTid = wss.tid
		return tss.messagehandle(tm, true)
	}
	return
}

func (tss *timservice) messagehandle(tm *stub.TimMessage, auth bool) (_r errs.ERROR) {
	if tm.ToTid != nil {
		tm.ToTid.Domain = tm.FromTid.Domain
	}
	if !existUser(tm.ToTid) && !existGroup(tm.RoomTid) {
		return errs.ERR_ACCOUNT
	}
	if tm.MsType == sys.SOURCE_ROOM {
		if tm.RoomTid != nil {
			if auth && !AuthGroup(tm.RoomTid.Node, tm.FromTid.Node, tm.FromTid.Domain) {
				return errs.ERR_PERM_DENIED
			}
			var err error
			switch tm.OdType {
			case sys.ORDER_INOF:
				err = data.Service.SaveMessage(tm)
			case sys.ORDER_REVOKE:
				if tm.GetMid() == 0 {
					return errs.ERR_PARAMS
				}
				tid := util.ChatIdByRoom(tm.RoomTid.Node, tm.FromTid.Domain)
				if fid, err := data.Service.GetFidByMid(tid, tm.GetMid()); err == nil && fid != 0 {
					b1 := uint64(goutil.BytesToInt64(tid[0:8])) == util.CreateUUID(tm.RoomTid.Node, tm.FromTid.Domain)
					b2 := int32(fid) == int32(goutil.FNVHash32([]byte(tm.FromTid.Node)))
					if b1 && b2 {
						if err = data.Service.DelMessageByMid(tid, tm.GetMid()); err == nil {
							t := int64(sys.SOURCE_ROOM)
							wsware.SendNode(tm.FromTid.Node, &stub.TimAck{Ok: true, TimType: int8(sys.TIMREVOKEMESSAGE), N: &tm.RoomTid.Node, T: &t, T2: tm.Mid}, sys.TIMACK)
						} else {
							return errs.ERR_DATABASE
						}
					} else {
						return errs.ERR_PERM_DENIED
					}
				} else {
					return errs.ERR_PARAMS
				}
			case sys.ORDER_STREAM, sys.ORDER_BIGSTRING, sys.ORDER_BIGBINARY:
			default:
				return errs.ERR_PARAMS
			}
			if err == nil {
				//if tm.OdType == sys.ORDER_INOF {
				//	if !wsware.SendNode(tm.FromTid.Node, tm, sys.TIMMESSAGE) {
				//		timMessage4goal(tm.FromTid.Node, tm)
				//	}
				//}
				if rs := data.Service.GroupRoster(tm.RoomTid.Node); len(rs) > 0 {
					tm.ToList = rs
					sys.TimMessageProcessor(tm, sys.TRANS_SOURCE)
				}
			}
		} else {
			return errs.ERR_PERM_DENIED
		}
	} else if tm.MsType == sys.SOURCE_USER {
		if tm.RoomTid != nil && tm.ToTid != nil {
			if AuthGroup(tm.RoomTid.Node, tm.FromTid.Node, tm.FromTid.Domain) && AuthGroup(tm.RoomTid.Node, tm.ToTid.Node, tm.FromTid.Domain) {
				return tss.messagehandler(tm, auth)
			} else {
				return errs.ERR_PERM_DENIED
			}
		} else if tm.ToTid != nil && authUser(tm.FromTid, tm.ToTid, false) {
			return tss.messagehandler(tm, auth)
		} else {
			return errs.ERR_PERM_DENIED
		}
	} else {
		return errs.ERR_PARAMS
	}
	return
}

func (tss *timservice) messagehandler(tm *stub.TimMessage, auth bool) (_r errs.ERROR) {
	ok := true
	switch tm.OdType {
	case sys.ORDER_INOF:
		if err := data.Service.SaveMessage(tm); err == nil {
			if !wsware.SendNode(tm.FromTid.Node, tm, sys.TIMMESSAGE) {
				timMessage4goal(tm.FromTid.Node, tm)
			}
		} else {
			return errs.ERR_DATABASE
		}
	case sys.ORDER_REVOKE:
		if tm.GetMid() == 0 {
			return errs.ERR_PARAMS
		}
		if auth && !authUser(tm.FromTid, tm.ToTid, true) {
			return errs.ERR_PERM_DENIED
		}
		tid := util.ChatIdByNode(tm.FromTid.Node, tm.ToTid.Node, tm.FromTid.Domain)
		if fid, err := data.Service.GetFidByMid(tid, tm.GetMid()); err == nil && fid != 0 {
			if int32(fid) == int32(goutil.FNVHash32([]byte(tm.FromTid.Node))) {
				if err = data.Service.DelMessageByMid(tid, tm.GetMid()); err == nil {
					//t := int64(sys.SOURCE_USER)
					//wsware.SendNode(tm.FromTid.Node, &stub.TimAck{Ok: true, TimType: int8(sys.TIMREVOKEMESSAGE), N: &tm.ToTid.Node, T: &t, T2: tm.Mid}, sys.TIMACK)
					timMessage4goal(tm.FromTid.Node, tm)
				} else {
					return errs.ERR_DATABASE
				}
				ok = true
			} else {
				return errs.ERR_PERM_DENIED
			}
		} else {
			return errs.ERR_PARAMS
		}
	case sys.ORDER_BURN:
		if tm.GetMid() == 0 {
			return errs.ERR_PARAMS
		}
		if auth && !authUser(tm.FromTid, tm.ToTid, true) {
			return errs.ERR_PERM_DENIED
		}
		tid := util.ChatIdByNode(tm.FromTid.Node, tm.ToTid.Node, tm.FromTid.Domain)
		if fid, err := data.Service.GetFidByMid(tid, tm.GetMid()); err == nil && fid != 0 {
			if int32(fid) == int32(goutil.FNVHash32([]byte(tm.ToTid.Node))) {
				if err = data.Service.DelMessageByMid(tid, tm.GetMid()); err == nil {
					timMessage4goal(tm.FromTid.Node, tm)
					//t := int64(sys.SOURCE_USER)
					//wsware.SendNode(tm.FromTid.Node, &stub.TimAck{Ok: true, TimType: int8(sys.TIMBURNMESSAGE), N: &tm.ToTid.Node, T: &t, T2: tm.Mid}, sys.TIMACK)
				} else {
					return errs.ERR_DATABASE
				}
				ok = true
			} else {
				return errs.ERR_PERM_DENIED
			}
		} else {
			return errs.ERR_PARAMS
		}
	case sys.ORDER_BUSINESS, sys.ORDER_STREAM, sys.ORDER_BIGSTRING, sys.ORDER_BIGBINARY:
	default:
		if tm.OdType <= sys.ORDER_RESERVED {
			return errs.ERR_PARAMS
		}
	}
	if ok {
		_r = sys.TimMessageProcessor(tm, sys.TRANS_SOURCE)
	} else {
		return errs.ERR_PARAMS
	}
	return
}

func (tss *timservice) presence(bs []byte, ws *tlnet.Websocket) (_r errs.ERROR) {
	defer util.Recover()
	sys.Stat.TxDo()
	defer sys.Stat.TxDone()
	tp := newTimPresence(bs)
	if tp == nil {
		return errs.ERR_FORMAT
	}
	if !checkTid(tp.ToTid) || !checkList(tp.ToList) {
		return errs.ERR_ACCOUNT
	}
	if wss, b := wsware.Get(ws); b {
		tp.FromTid = wss.tid
		if !existUser(tp.ToTid) && !existList(tp.ToList) {
			return errs.ERR_ACCOUNT
		}
		sys.TimPresenceProcessor(tp, sys.TRANS_SOURCE)
	}
	return
}

func (tss *timservice) interrupt(tid *stub.Tid) (err errs.ERROR) {
	defer util.Recover()
	sys.Stat.TxDo()
	defer sys.Stat.TxDone()
	if !wsware.hasUser(tid.Node) {
		if !sys.Conf.PresenceOfflineBlock {
			a := true
			if rs := data.Service.Roster(tid.Node); len(rs) > 0 {
				rid := goutil.UUID64()
				tp := &stub.TimPresence{ID: &rid, FromTid: tid, ToList: rs, Offline: &a}
				sys.TimPresenceProcessor(tp, sys.TRANS_SOURCE)
			}
			rid := goutil.UUID64()
			tp := &stub.TimPresence{ID: &rid, FromTid: tid, ToTid: tid, Offline: &a}
			sys.TimPresenceProcessor(tp, sys.TRANS_SOURCE)
		}
	}
	return
}

func (tss *timservice) offlineMsg(ws *tlnet.Websocket) (err errs.ERROR) {
	defer util.Recover()
	sys.Stat.TxDo()
	defer sys.Stat.TxDone()
	if wss, ok := wsware.Get(ws); ok {
		if oblist, _ := data.Service.GetOfflineMessage(wss.tid.Node, 10); len(oblist) > 0 {
			tmList := make([]*stub.TimMessage, 0)
			isOff := true
			ids := make([]any, 0)
			for _, ob := range oblist {
				ids = append(ids, ob.Id)
				if len(ob.Stanze) > 0 {
					if tm, err := goutil.TDecode(ob.Stanze, &stub.TimMessage{}); err == nil && tm != nil {
						tm.IsOffline = &isOff
						tmList = append(tmList, tm)
						tm.Mid = &ob.Mid
					}
				}
			}
			sort.Slice(tmList, func(i, j int) bool {
				return tmList[i].GetTimestamp() < tmList[j].GetTimestamp()
			})
			id := goutil.UUID64()
			if wsware.SendWsWithAck(ws.Id, &stub.TimMessageList{MessageList: tmList, ID: &id}, sys.TIMOFFLINEMSG) {
				if _r, err := data.Service.DelOfflineMessage(util.NodeToUUID(wss.tid.Node), ids...); err == nil && _r > 0 {
					tss.offlineMsg(ws)
				}
			}
		} else if err == nil {
			wss.Send(nil, sys.TIMOFFLINEMSGEND)
		} else {
			return errs.ERR_DATABASE
		}
	}
	return
}

func (tss *timservice) broadPresence(bs []byte, ws *tlnet.Websocket) (e errs.ERROR) {
	defer util.Recover()
	sys.Stat.TxDo()
	defer sys.Stat.TxDone()
	tp := newTimPresence(bs)
	if tp == nil {
		return errs.ERR_FORMAT
	}
	if wss, ok := wsware.Get(ws); ok {
		fid := wss.tid
		if tp.ToTid == nil && tp.ToList == nil {
			if rs := data.Service.Roster(wss.tid.Node); len(rs) > 0 {
				t := &stub.TimPresence{FromTid: fid, ToList: rs, SubStatus: tp.SubStatus, Show: tp.Show, Status: tp.Status, Extend: tp.Extend, Extra: tp.Extra}
				sys.TimPresenceProcessor(t, sys.TRANS_SOURCE)
			}
		} else {
			tp.FromTid = fid
			if tp.ToTid != nil {
				tp.ToList = nil
				sys.TimPresenceProcessor(tp, sys.TRANS_SOURCE)
			} else if tp.ToList != nil {
				t := &stub.TimPresence{FromTid: fid, ToList: tp.ToList, SubStatus: tp.SubStatus, Show: tp.Show, Status: tp.Status, Extend: tp.Extend, Extra: tp.Extra}
				sys.TimPresenceProcessor(t, sys.TRANS_SOURCE)
			}
		}
	}
	return
}

func (tss *timservice) pullmessage(bs []byte, ws *tlnet.Websocket) (err errs.ERROR) {
	defer util.Recover()
	sys.Stat.TxDo()
	defer sys.Stat.TxDone()
	tr := newTimReq(bs)
	if tr == nil || tr.Rtype == nil || tr.Node == nil || tr.ReqInt == nil || tr.ReqInt2 == nil || !checkNode(tr.GetNode()) {
		return errs.ERR_PARAMS
	}
	if wss, b := wsware.Get(ws); b {
		if tr.GetRtype() == 2 {
			if !AuthGroup(tr.GetNode(), wss.tid.Node, wss.tid.Domain) {
				return errs.ERR_PERM_DENIED
			}
		}
		timestamp := int64(0)
		mid := tr.GetReqInt()
		if tr.GetNode2() != "" {
			timestamp, _ = strconv.ParseInt(tr.GetNode2(), 10, 64)
		}
		if oblist, _ := data.Service.GetMessage(wss.tid.Node, wss.tid.Domain, int8(tr.GetRtype()), tr.GetNode(), mid, timestamp, tr.GetReqInt2()); len(oblist) > 0 {
			//if oblist[0].GetMid() == tr.GetReqInt() {
			//	oblist = oblist[1:]
			//}
			//sort.Slice(oblist, func(i, j int) bool { return oblist[i].GetTimestamp() > oblist[j].GetTimestamp() })
			wss.Send(&stub.TimMessageList{MessageList: oblist}, sys.TIMPULLMESSAGE)
		}
	}
	return
}

func (tss *timservice) osvroomprocess(node string, rtype int8) (_r string) {
	switch rtype {
	case 1:
		_r = vgate.VGate.NewVRoom(node)
	}
	return
}

func (tss *timservice) vroomHandle(bs []byte, ws *tlnet.Websocket) (err errs.ERROR) {
	defer util.Recover()
	sys.Stat.TxDo()
	defer sys.Stat.TxDone()
	tr := newTimReq(bs)
	if tr == nil || tr.Rtype == nil {
		return errs.ERR_PARAMS
	}
	if wss, b := wsware.Get(ws); b {
		t := int64(tr.GetRtype())
		switch sys.TIMTYPE(tr.GetRtype()) {
		case sys.VROOM_NEW:
			vnode := vgate.VGate.NewVRoom(wss.tid.Node)
			wss.Send(&stub.TimAck{Ok: true, TimType: int8(sys.TIMVROOM), N: &vnode, T: &t}, sys.TIMACK)
		case sys.VROOM_REMOVE:
			if tr.Node == nil || !util.CheckNode(tr.GetNode()) {
				return errs.ERR_PARAMS
			}
			vgate.VGate.Remove(wss.tid.Node, tr.GetNode())
			wss.Send(&stub.TimAck{Ok: true, TimType: int8(sys.TIMVROOM), N: tr.Node, T: &t}, sys.TIMACK)
		case sys.VROOM_ADDAUTH:
			if tr.Node == nil || tr.Node2 == nil || !util.CheckNode(tr.GetNode()) || !checkNode(tr.GetNode2()) {
				return errs.ERR_PARAMS
			}
			wss.Send(&stub.TimAck{Ok: false, TimType: int8(sys.TIMVROOM), N: tr.Node2, T: &t}, sys.TIMACK)
		case sys.VROOM_DELAUTH:
			if tr.Node == nil || tr.Node2 == nil || !util.CheckNode(tr.GetNode()) || !util.CheckNode(tr.GetNode2()) {
				return errs.ERR_PARAMS
			}
			wss.Send(&stub.TimAck{Ok: false, TimType: int8(sys.TIMVROOM), N: tr.Node2, T: &t}, sys.TIMACK)
		case sys.VROOM_SUB:
			if tr.Node == nil || !util.CheckNode(tr.GetNode()) {
				return errs.ERR_PARAMS
			}
			var success = false
			if tr.ReqInt == nil {
				success = vgate.VGate.Sub(tr.GetNode(), sys.UUID, wss.ws.Id)
			} else if tr.GetReqInt() == 1 {
				success = vgate.VGate.SubBinary(tr.GetNode(), sys.UUID, wss.ws.Id)
			}
			if success {
				if _, err = sys.TimSteamProcessor(&stub.VBean{Rtype: int8(sys.VROOM_SUB), Vnode: tr.GetNode(), Rnode: &wss.tid.Node}, sys.TRANS_SOURCE); err == nil {
					wss.Send(&stub.TimAck{Ok: true, TimType: int8(sys.TIMVROOM), N: tr.Node, T: &t}, sys.TIMACK)
				}
			} else {
				return errs.ERR_UNDEFINED
			}
		case sys.VROOM_UNSUB:
			if tr.Node == nil || !util.CheckNode(tr.GetNode()) {
				return errs.ERR_PARAMS
			}
			if r, b := vgate.VGate.UnSub(tr.GetNode(), wss.ws.Id); b && r == 0 {
				sys.TimSteamProcessor(&stub.VBean{Rtype: int8(sys.VROOM_UNSUB), Vnode: tr.GetNode(), Rnode: &wss.tid.Node}, sys.TRANS_SOURCE)
			}
			wss.Send(&stub.TimAck{Ok: true, TimType: int8(sys.TIMVROOM), N: tr.Node, T: &t}, sys.TIMACK)
		default:
			return errs.ERR_PARAMS
		}
	}
	return
}

func (tss *timservice) streamHandle(bs []byte, ws *tlnet.Websocket) (err errs.ERROR) {
	defer util.Recover()
	sys.Stat.TxDo()
	defer sys.Stat.TxDone()
	t := newTimStream(bs)
	return tss.streamhandler(t, ws)
}

func (tss *timservice) streamhandler(t *stub.TimStream, ws *tlnet.Websocket) (err errs.ERROR) {
	if tss == nil || !util.CheckNode(t.VNode) {
		return errs.ERR_PARAMS
	}
	if wss, b := wsware.Get(ws); b {
		if uuid := amr.GetVnode(t.VNode); uuid == sys.UUID {
			auth := false
			if vr, ok := vgate.VGate.GetVroom(t.VNode); ok {
				auth = vr.AuthStream(wss.tid.Node)
			}
			if !auth {
				wss.Send(&stub.TimAck{Ok: false, TimType: int8(sys.TIMSTREAM), Error: errs.ERR_PERM_DENIED.TimError(), N: &t.VNode}, sys.TIMACK)
				return
			}
			csvb := &stub.VBean{Vnode: t.VNode, Rnode: &wss.tid.Node, Body: t.Body, Dtype: t.Dtype, Rtype: int8(sys.VROOM_MESSAGE), StreamId: &t.ID}
			sys.TimSteamProcessor(csvb, sys.TRANS_SOURCE)
		}
	}
	return
}
