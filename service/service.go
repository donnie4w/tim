// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package service

import (
	"bytes"
	"github.com/donnie4w/tim/amr"
	"github.com/donnie4w/tim/errs"
	"sort"
	"strings"
	"time"

	. "github.com/donnie4w/gofer/util"
	"github.com/donnie4w/tim/cache"
	"github.com/donnie4w/tim/data"
	. "github.com/donnie4w/tim/stub"
	"github.com/donnie4w/tim/sys"
	"github.com/donnie4w/tim/util"
	"github.com/donnie4w/tim/vgate"
	"github.com/donnie4w/tlnet"
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
	var ta *TimAuth
	var err error
	if util.JTP(bs[0]) {
		ta, err = JsonDecode[*TimAuth](bs[1:])
	} else {
		ta, err = TDecode(bs[1:], &TimAuth{})
	}
	if err == nil {
		node, e = data.Service.Register(*ta.Name, *ta.Pwd, ta.Domain)
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

func (tss *timservice) ostoken(nodeorname string, password, domain *string) (_r int64, _n string, e errs.ERROR) {
	defer util.Recover()
	sys.Stat.TxDo()
	defer sys.Stat.TxDone()
	if password != nil {
		var node string
		if node, e = data.Service.AuthNode(nodeorname, *password, domain); e == nil {
			tid := &Tid{Node: node, Domain: domain}
			_r, _n = token(), node
			//tokenTempCache.add(_r, tid)
			cache.TokenCache.Put(_r, tid)
		}
	} else {
		switch sys.GetDBMOD() {
		case sys.NODB, sys.EXTERNALDB:
			_r, _n = token(), nodeorname
		case sys.TLDB, sys.INLINEDB:
			if !existUser(&Tid{Node: nodeorname, Domain: domain}) {
				return _r, "", errs.ERR_NOEXIST
			}
			_r, _n = token(), util.UUIDToNode(util.CreateUUID(nodeorname, domain))
		default:
			e = errs.ERR_DATABASE
		}
		tid := &Tid{Node: _n, Domain: domain}
		cache.TokenCache.Put(_r, tid)
	}
	return
}

func (tss *timservice) token(bs []byte) (_r int64, e errs.ERROR) {
	defer util.Recover()
	sys.Stat.TxDo()
	defer sys.Stat.TxDone()
	ta := newAuth(bs)
	if ta == nil {
		return _r, errs.ERR_PARAMS
	}
	var node string
	if node, e = data.Service.AuthNode(*ta.Name, *ta.Pwd, ta.Domain); e == nil {
		tid := &Tid{Node: node, Domain: ta.Domain, Extend: ta.Extend}
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
		wsware.SendWs(ws.Id, &TimAck{Ok: true, TimType: int8(sys.TIMAUTH)}, sys.TIMACK)
		return
	}
	ta := newAuth(bs)
	if ta == nil {
		return errs.ERR_PARAMS
	}
	isAuth := false
	var tid *Tid
	if ta.Token != nil {
		if tid, _ = cache.TokenCache.Get(*ta.Token); tid != nil {
			tid.Resource, tid.Termtyp = ta.Resource, ta.Termtyp
			cache.TokenCache.Del(*ta.Token)
			if !isblock(tid.Node) {
				isAuth = true
			}
		} else {
			return errs.ERR_TOKEN
		}
	} else if ta.Name != nil && ta.Pwd != nil && !isblock(*ta.Name) {
		if _r, err := data.Service.Login(*ta.Name, *ta.Pwd, ta.Domain); err == nil {
			tid = &Tid{Node: _r, Domain: ta.Domain, Extend: ta.Extend, Resource: ta.Resource, Termtyp: ta.Termtyp}
			if !isblock(_r) {
				isAuth = true
			}
		}
	}
	if isAuth {
		overentry := true
		if wsware.GetUserDeviceLen(tid.Node) < sys.DeviceLimit {
			wis := sys.CsWssInfo(tid.Node)
			if len(wis)+wsware.GetUserDeviceLen(tid.Node) < sys.DeviceLimit {
				overentry = false
				if tid.Termtyp != nil {
					typebs := wsware.GetUserDeviceTypeLen(tid.Node)
					c := 0
					for _, u := range append(wis, typebs...) {
						if u == byte(*tid.Termtyp) {
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
		wsware.SendWs(ws.Id, &TimAck{Ok: true, TimType: int8(sys.TIMAUTH), N: &tid.Node}, sys.TIMACK)
	} else {
		e = errs.ERR_PERM_DENIED
	}
	return
}

func sysMessage(nodelist []string, tm *TimMessage) (err errs.ERROR) {
	if len(nodelist) == 0 && tm == nil {
		return errs.ERR_PARAMS
	}
	if checkList(nodelist) {
		t := time.Now().UnixNano()
		for _, u := range nodelist {
			tm.ToTid = &Tid{Node: u}
			tm.Timestamp = &t
			service.osmessage(tm)
		}
	} else {
		return errs.ERR_ACCOUNT
	}
	return
}

func (tss *timservice) osmessage(tm *TimMessage) (err errs.ERROR) {
	defer util.Recover()
	sys.Stat.TxDo()
	defer sys.Stat.TxDone()
	if tm == nil || tm.ToTid == nil {
		return errs.ERR_PARAMS
	}
	tm.MsType, tm.OdType = sys.SOURCE_OS, sys.ORDER_INOF
	return sys.TimMessageProcessor(tm, sys.TRANS_SOURCE)
}

func sysPresence(nodelist []string, tm *TimPresence) (err errs.ERROR) {
	if len(nodelist) == 0 && tm == nil {
		return errs.ERR_PARAMS
	}
	if checkList(nodelist) {
		tm.ToList = nodelist
		tm.Offline = nil
		sys.TimPresenceProcessor(tm, sys.TRANS_STAFF)
	} else {
		return errs.ERR_ACCOUNT
	}
	return
}

func (tss *timservice) pxmessage(connectid int64, tm *TimMessage) (err errs.ERROR) {
	if ws, b := sys.WsById(connectid); b {
		err = tss.message(tm, ws)
	} else {
		err = errs.ERR_NOEXIST
	}
	return
}

func (tss *timservice) bigString(bs []byte, ws *tlnet.Websocket) (_r errs.ERROR) {
	bigString := string(bs[5:])
	idx := strings.Index(bigString, sys.SEP_STR)
	dataString := bigString[idx+1:]
	if wss, b := wsware.Get(ws); b {
		_r = tss.messagehandler(&TimMessage{MsType: 2, OdType: sys.ORDER_BIGSTRING, DataString: &dataString, FromTid: wss.tid, ToTid: &Tid{Node: bigString[:idx]}})
	}
	return
}

func (tss *timservice) bigBinary(bs []byte, ws *tlnet.Websocket) (_r errs.ERROR) {
	idx := bytes.IndexByte(bs[5:], sys.SEP_BIN)
	if wss, b := wsware.Get(ws); b {
		_r = tss.messagehandler(&TimMessage{MsType: 2, OdType: sys.ORDER_BIGBINARY, DataBinary: bs[5:][idx+1:], FromTid: wss.tid, ToTid: &Tid{Node: string(bs[5:][:idx])}})
	}
	return
}

func (tss *timservice) bigBinaryStreamHandle(bs []byte, ws *tlnet.Websocket) (_r errs.ERROR) {
	idx := bytes.IndexByte(bs[5:], sys.SEP_BIN)
	if wss, b := wsware.Get(ws); b {
		t := &TimStream{ID: UUID64(), VNode: string(bs[5:][:idx]), Body: bs[5:][idx+1:], FromNode: wss.tid.Node}
		return tss.streamhandler(t, ws)
	}
	return
}

func (tss *timservice) messageHandle(bs []byte, ws *tlnet.Websocket) (_r errs.ERROR) {
	return tss.message(newTimMessage(bs), ws)
}

func (tss *timservice) message(tm *TimMessage, ws *tlnet.Websocket) (_r errs.ERROR) {
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
		if !existUser(tm.ToTid) || !existGroup(tm.RoomTid) {
			return errs.ERR_ACCOUNT
		}
		if tm.MsType == sys.SOURCE_ROOM {
			if tm.RoomTid != nil && AuthGroup(tm.RoomTid.Node, tm.FromTid.Node, tm.FromTid.Domain) {
				var err error
				switch tm.OdType {
				case sys.ORDER_INOF:
					err = data.Service.SaveMessage(tm)
				case sys.ORDER_REVOKE:
					if tm.Mid == nil || *tm.Mid == 0 {
						return errs.ERR_PARAMS
					}
					tid := util.ChatIdByRoom(tm.RoomTid.Node, wss.tid.Domain)
					if _t, err := data.Service.GetMessageByMid(tid, *tm.Mid); err == nil && _t != nil {
						if _t.FromTid.Node == tm.FromTid.Node && _t.RoomTid.Node == tm.RoomTid.Node {
							if err = data.Service.DelMessageByMid(tid, *tm.Mid); err == nil {
								t := int64(sys.SOURCE_ROOM)
								wsware.SendNode(tm.FromTid.Node, &TimAck{Ok: true, TimType: int8(sys.TIMREVOKEMESSAGE), N: &tm.RoomTid.Node, T: &t, T2: tm.Mid}, sys.TIMACK)
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
					if tm.OdType == sys.ORDER_INOF {
						wsware.SendNode(tm.FromTid.Node, tm, sys.TIMMESSAGE)
					}
					if rs := data.Service.GroupRoster(tm.RoomTid.Node); len(rs) > 0 {
						//for _, u := range rs {
						//	if u != wss.tid.Node {
						//		t := shallowCloneTimMessageData(tm)
						//		t.FromTid, t.ToTid = wss.tid, &Tid{Node: u}
						//		util.GoPoolTx.Go(func() { sys.TimMessageProcessor(t, sys.TRANS_SOURCE) })
						//	}
						//}
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
					return tss.messagehandler(tm)
				} else {
					return errs.ERR_PERM_DENIED
				}
			} else if tm.ToTid != nil && AuthUser(tm.FromTid, tm.ToTid, false) {
				return tss.messagehandler(tm)
			} else {
				return errs.ERR_PERM_DENIED
			}
		} else {
			return errs.ERR_PARAMS
		}
	}
	return
}

func (tss *timservice) messagehandler(tm *TimMessage) (_r errs.ERROR) {
	ok := true
	switch tm.OdType {
	case sys.ORDER_INOF:
		if err := data.Service.SaveMessage(tm); err == nil {
			wsware.SendNode(tm.FromTid.Node, tm, sys.TIMMESSAGE)
		} else {
			return errs.ERR_DATABASE
		}
	case sys.ORDER_REVOKE:
		if tm.Mid == nil || *tm.Mid == 0 {
			return errs.ERR_PARAMS
		}
		if !AuthUser(tm.FromTid, tm.ToTid, true) {
			return errs.ERR_PERM_DENIED
		}
		tid := util.ChatIdByNode(tm.FromTid.Node, tm.ToTid.Node, tm.FromTid.Domain)
		if _t, err := data.Service.GetMessageByMid(tid, *tm.Mid); err == nil && _t != nil {
			if _t.FromTid.Node == tm.FromTid.Node && _t.ToTid.Node == tm.ToTid.Node {
				if err = data.Service.DelMessageByMid(tid, *tm.Mid); err == nil {
					t := int64(sys.SOURCE_USER)
					wsware.SendNode(tm.FromTid.Node, &TimAck{Ok: true, TimType: int8(sys.TIMREVOKEMESSAGE), N: &tm.ToTid.Node, T: &t, T2: tm.Mid}, sys.TIMACK)
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
		if tm.Mid == nil || *tm.Mid == 0 {
			return errs.ERR_PARAMS
		}
		if !AuthUser(tm.FromTid, tm.ToTid, true) {
			return errs.ERR_PERM_DENIED
		}
		tid := util.ChatIdByNode(tm.FromTid.Node, tm.ToTid.Node, tm.FromTid.Domain)
		if _t, err := data.Service.GetMessageByMid(tid, *tm.Mid); err == nil && _t != nil {
			if _t.FromTid.Node == tm.ToTid.Node && _t.ToTid.Node == tm.FromTid.Node {
				if err = data.Service.DelMessageByMid(tid, *tm.Mid); err == nil {
					t := int64(sys.SOURCE_USER)
					wsware.SendNode(tm.FromTid.Node, &TimAck{Ok: true, TimType: int8(sys.TIMBURNMESSAGE), N: &tm.ToTid.Node, T: &t, T2: tm.Mid}, sys.TIMACK)
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
		if !existUser(tp.ToTid) || !existList(tp.ToList, tp.FromTid.Domain) {
			return errs.ERR_ACCOUNT
		}
		if tp.ToList != nil {
			sys.TimPresenceProcessor(tp, sys.TRANS_STAFF)
		} else if tp.ToTid != nil {
			sys.TimPresenceProcessor(tp, sys.TRANS_SOURCE)
		}
	}
	return
}

func (tss *timservice) interrupt(tid *Tid) (err errs.ERROR) {
	defer util.Recover()
	sys.Stat.TxDo()
	defer sys.Stat.TxDone()
	if !wsware.hasUser(tid.Node) {
		if !sys.Conf.PresenceOfflineBlock {
			a := true
			if rs := data.Service.Roster(tid.Node); len(rs) > 0 {
				rid := UUID64()
				tp := &TimPresence{ID: &rid, FromTid: tid, ToList: rs, Offline: &a}
				sys.TimPresenceProcessor(tp, sys.TRANS_STAFF)
			}
			rid := UUID64()
			tp := &TimPresence{ID: &rid, FromTid: tid, ToTid: tid, Offline: &a}
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
		if oblist, err := data.Service.GetOfflineMessage(wss.tid.Node, 10); err == nil && oblist != nil && len(oblist) > 0 {
			tmList := make([]*TimMessage, 0)
			isOff := true
			ids := make([]int64, 0)
			for _, ob := range oblist {
				ids = append(ids, ob.Id)
				if ob.Stanze != nil {
					if tm, err := TDecode(ob.Stanze, &TimMessage{}); err == nil {
						tm.IsOffline = &isOff
						tmList = append(tmList, tm)
						if ob.Mid > 0 {
							tm.Mid = &ob.Mid
						}
					}
				}
			}
			sort.Slice(tmList, func(i, j int) bool {
				return *tmList[i].Timestamp < *tmList[j].Timestamp
			})
			id := UUID64()
			if wsware.SendWsWithAck(ws.Id, &TimMessageList{MessageList: tmList, ID: &id}, sys.TIMOFFLINEMSG) {
				if _r, err := data.Service.DelOfflineMessage(util.NodeToUUID(wss.tid.Node), ids...); err == nil && _r > 0 {
					tss.offlineMsg(ws)
				}
			}
		} else if err == nil {
			wsware.SendWs(ws.Id, nil, sys.TIMOFFLINEMSGEND)
		} else {
			return errs.ERR_DATABASE
		}
	}
	return
}

func (tss *timservice) broadpresence(bs []byte, ws *tlnet.Websocket) (e errs.ERROR) {
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
				for i := 0; i < len(rs); i++ {
					t := &TimPresence{FromTid: fid, ToTid: &Tid{Node: rs[i]}, SubStatus: tp.SubStatus, Show: tp.Show, Status: tp.Status, Extend: tp.Extend, Extra: tp.Extra}
					sys.TimPresenceProcessor(t, sys.TRANS_SOURCE)
				}
			}
		} else {
			tp.FromTid = fid
			if tp.ToTid != nil {
				tp.ToList = nil
				sys.TimPresenceProcessor(tp, sys.TRANS_SOURCE)
			} else if tp.ToList != nil {
				tp.ToTid = nil
				for _, u := range tp.ToList {
					t := &TimPresence{FromTid: fid, ToTid: &Tid{Node: u}, SubStatus: tp.SubStatus, Show: tp.Show, Status: tp.Status, Extend: tp.Extend, Extra: tp.Extra}
					sys.TimPresenceProcessor(t, sys.TRANS_SOURCE)
				}
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
	if tr == nil || tr.Rtype == nil || tr.Node == nil || tr.ReqInt == nil || tr.ReqInt2 == nil || !checkNode(*tr.Node) {
		return errs.ERR_PARAMS
	}
	if wss, b := wsware.Get(ws); b {
		if *tr.Rtype == 2 {
			if !AuthGroup(*tr.Node, wss.tid.Node, wss.tid.Domain) {
				return errs.ERR_PERM_DENIED
			}
		}
		if oblist, err := data.Service.GetMessage(wss.tid.Node, wss.tid.Domain, int8(*tr.Rtype), *tr.Node, *tr.ReqInt, *tr.ReqInt2); err == nil && len(oblist) > 0 {
			if *oblist[0].Mid == *tr.ReqInt {
				oblist = oblist[1:]
			}
			sort.Slice(oblist, func(i, j int) bool { return *oblist[i].Mid > *oblist[j].Mid })
			wsware.SendWs(ws.Id, &TimMessageList{MessageList: oblist}, sys.TIMPULLMESSAGE)
		}
	}
	return
}

func (tss *timservice) osvroomprocess(node string, rtype int8) (_r string) {
	switch rtype {
	case 1:
		_r = vgate.VGate.NewVRoom(node)
		//if sys.CsVBean(&VBean{Rtype: 1, Vnode: vnode, FoundNode: &node}) {
		//	_r = vnode
		//}
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
		t := int64(*tr.Rtype)
		switch sys.TIMTYPE(*tr.Rtype) {
		case sys.VROOM_NEW:
			vnode := vgate.VGate.NewVRoom(wss.tid.Node)
			//if sys.CsVBean(&VBean{Rtype: 1, Vnode: vnode, FoundNode: &wss.tid.Node}) {
			wsware.SendWs(ws.Id, &TimAck{Ok: true, TimType: int8(sys.TIMVROOM), N: &vnode, T: &t}, sys.TIMACK)
			//} else {
			//	return errs.ERR_UNDEFINED
			//}
		case sys.VROOM_REMOVE:
			if tr.Node == nil || !util.CheckNode(*tr.Node) {
				return errs.ERR_PARAMS
			}
			vgate.VGate.Remove(wss.tid.Node, *tr.Node)
			//sys.CsVBean(&VBean{Rtype: 2, Vnode: *tr.Node, FoundNode: &wss.tid.Node})
			wsware.SendWs(ws.Id, &TimAck{Ok: true, TimType: int8(sys.TIMVROOM), N: tr.Node, T: &t}, sys.TIMACK)
		case sys.VROOM_ADDAUTH:
			if tr.Node == nil || tr.Node2 == nil || !util.CheckNode(*tr.Node) || !checkNode(*tr.Node2) {
				return errs.ERR_PARAMS
			}
			//vgate.VGate.AddAuth(*tr.Node, wss.tid.Node, *tr.Node2)
			//sys.CsVBean(&VBean{Rtype: 3, Vnode: *tr.Node, FoundNode: &wss.tid.Node, Rnode: tr.Node2})
			wsware.SendWs(ws.Id, &TimAck{Ok: false, TimType: int8(sys.TIMVROOM), N: tr.Node2, T: &t}, sys.TIMACK)
		case sys.VROOM_DELAUTH:
			if tr.Node == nil || tr.Node2 == nil || !util.CheckNode(*tr.Node) || !util.CheckNode(*tr.Node2) {
				return errs.ERR_PARAMS
			}
			//vgate.VGate.DelAuth(*tr.Node, wss.tid.Node, *tr.Node2)
			//sys.CsVBean(&VBean{Rtype: 4, Vnode: *tr.Node, FoundNode: &wss.tid.Node, Rnode: tr.Node2})
			wsware.SendWs(ws.Id, &TimAck{Ok: false, TimType: int8(sys.TIMVROOM), N: tr.Node2, T: &t}, sys.TIMACK)
		case sys.VROOM_SUB:
			if tr.Node == nil || !util.CheckNode(*tr.Node) {
				return errs.ERR_PARAMS
			}
			var success = false
			if tr.ReqInt == nil {
				success = vgate.VGate.Sub(*tr.Node, sys.UUID, wss.ws.Id)
			} else if *tr.ReqInt == 1 {
				success = vgate.VGate.SubBinary(*tr.Node, sys.UUID, wss.ws.Id)
			}
			//if sys.CsVBean(&VBean{Rtype: 5, Vnode: *tr.Node, Rnode: &wss.tid.Node}) {
			if success && sys.TimSteamProcessor(&VBean{Rtype: int8(sys.VROOM_SUB), Vnode: *tr.Node, Rnode: &wss.tid.Node}, sys.TRANS_SOURCE) == nil {
				wsware.SendWs(ws.Id, &TimAck{Ok: true, TimType: int8(sys.TIMVROOM), N: tr.Node, T: &t}, sys.TIMACK)
			} else {
				return errs.ERR_UNDEFINED
			}
		case sys.VROOM_UNSUB:
			if tr.Node == nil || !util.CheckNode(*tr.Node) {
				return errs.ERR_PARAMS
			}
			if r, b := vgate.VGate.UnSub(*tr.Node, wss.ws.Id); b && r == 0 {
				//sys.CsVBean(&VBean{Rtype: 6, Vnode: *tr.Node})
				sys.TimSteamProcessor(&VBean{Rtype: int8(sys.VROOM_UNSUB), Vnode: *tr.Node, Rnode: &wss.tid.Node}, sys.TRANS_SOURCE)
			}
			wsware.SendWs(ws.Id, &TimAck{Ok: true, TimType: int8(sys.TIMVROOM), N: tr.Node, T: &t}, sys.TIMACK)
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

func (tss *timservice) streamhandler(t *TimStream, ws *tlnet.Websocket) (err errs.ERROR) {
	if tss == nil || !util.CheckNode(t.VNode) {
		return errs.ERR_PARAMS
	}
	if wss, b := wsware.Get(ws); b {
		if uuid := amr.GetVnode(t.VNode); uuid == sys.UUID {
			//if sys.CsNode(t.VNode) == sys.UUID {
			auth := false
			if vr, ok := vgate.VGate.GetVroom(t.VNode); ok {
				auth = vr.AuthStream(wss.tid.Node)
			}
			if !auth {
				wsware.SendWs(ws.Id, &TimAck{Ok: false, TimType: int8(sys.TIMSTREAM), Error: errs.ERR_PERM_DENIED.TimError(), N: &t.VNode}, sys.TIMACK)
				return
			}
			csvb := &VBean{Vnode: t.VNode, Rnode: &wss.tid.Node, Body: t.Body, Dtype: t.Dtype, Rtype: int8(sys.VROOM_MESSAGE), StreamId: &t.ID}
			//sys.CsVBean(csvb)
			sys.TimSteamProcessor(csvb, sys.TRANS_SOURCE)
		}
	}
	return
}
