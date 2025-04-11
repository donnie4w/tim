// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package service

import (
	gocache "github.com/donnie4w/gofer/cache"
	goutil "github.com/donnie4w/gofer/util"
	"github.com/donnie4w/gofer/uuid"
	"github.com/donnie4w/tim/amr"
	"github.com/donnie4w/tim/cache"
	"github.com/donnie4w/tim/data"
	. "github.com/donnie4w/tim/stub"
	"github.com/donnie4w/tim/sys"
	"github.com/donnie4w/tim/util"
	"time"
)

var userCache = gocache.NewBloomFilter(1<<21, 0.0001)
var groupCache = gocache.NewBloomFilter(1<<20, 0.0001)

func init() {
	sys.PingHandle = service.ping
	sys.AckHandle = service.ack
	sys.RegisterHandle = service.register
	sys.TokenHandle = service.token
	sys.AuthHandle = service.auth
	sys.Interrupt = service.interrupt
	sys.MessageHandle = service.messageHandle
	sys.PresenceHandle = service.presence
	sys.PullMessageHandle = service.pullmessage
	sys.OfflinemsgHandle = service.offlineMsg
	sys.BroadPresenceHandle = service.broadPresence
	sys.BusinessHandle = service.business
	sys.VRoomHandle = service.vroomHandle
	sys.StreamHandle = service.streamHandle
	sys.BigStringHandle = service.bigString
	sys.BigBinaryHandle = service.bigBinary
	sys.BigBinaryStreamHandle = service.bigBinaryStreamHandle
	sys.NodeInfoHandle = service.nodeinfo
	sys.AuthRoster = authRoster
	sys.AuthGroupuser = authGroupuser
	sys.OsModify = service.sysModify
	sys.OsMessage = sysMessage
	sys.OsPresence = sysPresence
	sys.OsUserBean = service.osuserbean
	sys.OsRoom = service.osnewgroup
	sys.OsRoomBean = service.osModifygroupInfo
	sys.OsVroomprocess = service.osvroomprocess
	sys.PxMessage = service.pxMessage
	sys.TimMessageProcessor = timMessageProcessor
	sys.TimPresenceProcessor = timPresenceProcessor
	sys.TimSteamProcessor = timStreamProcessor
	sys.HasNode = wsware.hasUser
	sys.HasWs = wsware.hasws
	sys.DelWs = wsware.delws
	sys.WsById = wsware.wsById
	sys.WsByNode = wsware.wsByNode
	sys.WssLen = wsware.wsLen
	sys.WssList = wsware.wssList
	sys.DeviceTypeList = wsware.deviceTypeList
	sys.OsToken = service.ostoken
	sys.OsRegister = service.osregister
	sys.SendNode = wsware.SendNode
	sys.SendWs = wsware.SendWs
	sys.OsBlockUser = blockNode
	sys.Detect = detect
}

func token() (r string) {
	return util.UUIDToNode(util.CreateUUID(uuid.NewUUID().String(), nil))
}

func existUser(tid *Tid) (b bool) {
	if tid == nil {
		return false
	}
	if sys.UseBuiltInData() {
		if b = userCache.Contains([]byte(tid.GetNode())); !b {
			if b = data.Service.ExistUser(tid.GetNode()); b {
				userCache.Add([]byte(tid.GetNode()))
			}
		}
	} else {
		b = true
	}
	return
}

func existList(ls []string) (b bool) {
	if ls == nil {
		return false
	}
	if sys.UseBuiltInData() {
		for _, node := range ls {
			if b = userCache.Contains([]byte(node)); !b {
				if b = data.Service.ExistUser(node); b {
					userCache.Add([]byte(node))
				} else {
					return false
				}
			}
		}
	}
	return true
}

func existGroup(tid *Tid) (_r bool) {
	if tid == nil {
		return false
	}
	if _r = groupCache.Contains([]byte(tid.GetNode())); !_r {
		if _r = data.Service.ExistGroup(tid.GetNode()); _r {
			groupCache.Add([]byte(tid.GetNode()))
		}
	}
	return
}

func authUser(fTid, tTid *Tid, readtime bool) (ok bool) {
	if sys.Conf.MessageNoAuth {
		return true
	}
	return authRoster(fTid.Node, tTid.Node, fTid.Domain, false)
}

func authRoster(fnode, tnode string, domain *string, readtime bool) (ok bool) {
	defer util.Recover()
	if sys.Conf.CacheAuthExpire > 0 && !readtime {
		if cache.AuthCache.Has(fnode, tnode, domain, false) {
			return true
		} else {
			if ok = data.Service.AuthUserAndUser(fnode, tnode, domain); ok {
				cache.AuthCache.Put(fnode, tnode, domain, false)
			}
		}
	} else {
		ok = data.Service.AuthUserAndUser(fnode, tnode, domain)
	}
	return
}

func AuthGroup(gnode, unode string, domain *string) (ok bool) {
	if sys.Conf.MessageNoAuth {
		return true
	}
	return authGroupuser(gnode, unode, domain)
}

func authGroupuser(gnode, unode string, domain *string) (ok bool) {
	defer util.Recover()
	if sys.Conf.CacheAuthExpire > 0 {
		if cache.AuthCache.Has(gnode, unode, domain, true) {
			return true
		} else {
			if ok, _ = data.Service.AuthGroupAndUser(gnode, unode, domain); ok {
				cache.AuthCache.Put(gnode, unode, domain, true)
			}
		}
	} else {
		ok, _ = data.Service.AuthGroupAndUser(gnode, unode, domain)
	}
	return
}

func newTimMessage(bs []byte) (tm *TimMessage) {
	var err error
	if util.JTP(bs[0]) {
		tm, err = goutil.JsonDecode[*TimMessage](bs[1:])
	} else {
		tm, err = goutil.TDecode(bs[1:], &TimMessage{})
	}
	if err == nil {
		fullTimMessage(tm)
	}
	return
}

func fullTimMessage(tm *TimMessage) *TimMessage {
	if tm.ID == nil {
		id := goutil.UUID64()
		tm.ID = &id
	}
	t := time.Now().UnixNano()
	tm.Timestamp = &t
	return tm
}

//func shallowCloneTimMessageData(tm *TimMessage) (r *TimMessage) {
//	r = &TimMessage{MsType: tm.MsType, OdType: tm.OdType, BnType: tm.BnType, Mid: tm.Mid, ID: tm.ID}
//	r.DataBinary = tm.DataBinary
//	r.DataString = tm.DataString
//	r.Extend = tm.Extend
//	r.Extra = tm.Extra
//	r.Timestamp = tm.Timestamp
//	r.Udshow = tm.Udshow
//	r.Udtype = tm.Udtype
//	r.FromTid = tm.FromTid
//	r.ToTid = tm.ToTid
//	r.RoomTid = tm.RoomTid
//	return
//}

func newTimPresence(bs []byte) (tp *TimPresence) {
	var err error
	if util.JTP(bs[0]) {
		tp, err = goutil.JsonDecode[*TimPresence](bs[1:])
	} else {
		tp, err = goutil.TDecode(bs[1:], &TimPresence{})
	}
	if err == nil {
		if tp.ID == nil {
			id := goutil.UUID64()
			tp.ID = &id
		}
		tp.Offline = nil
	}
	return
}

func checkTid(tid *Tid) (_r bool) {
	if sys.UseBuiltInData() && tid != nil {
		return util.CheckNode(tid.Node)
	}

	if tid != nil && len(tid.Node) > sys.NodeMaxSize {
		return false
	}
	return true
}

func checkNode(node string) (_r bool) {
	if sys.UseBuiltInData() && node != "" {
		return util.CheckNode(node)
	}
	if len(node) > sys.NodeMaxSize {
		return false
	}
	return true
}

func checkList(ls []string) (_r bool) {
	if sys.UseBuiltInData() && ls != nil {
		for _, u := range ls {
			if len(u) > sys.NodeMaxSize {
				return false
			}
			if !util.CheckNode(u) {
				return false
			}
		}
	}
	return true
}

func newTimReq(bs []byte) (tr *TimReq) {
	if util.JTP(bs[0]) {
		tr, _ = goutil.JsonDecode[*TimReq](bs[1:])
	} else {
		tr, _ = goutil.TDecode(bs[1:], &TimReq{})
	}
	return
}

func newTimNodes(bs []byte) (tr *TimNodes) {
	if util.JTP(bs[0]) {
		tr, _ = goutil.JsonDecode[*TimNodes](bs[1:])
	} else {
		tr, _ = goutil.TDecode(bs[1:], &TimNodes{})
	}
	return
}

func newAuth(bs []byte) (ta *TimAuth) {
	if util.JTP(bs[0]) {
		ta, _ = goutil.JsonDecode[*TimAuth](bs[1:])
	} else {
		ta, _ = goutil.TDecode(bs[1:], &TimAuth{})
	}
	return
}

func newTimStream(bs []byte) (ts *TimStream) {
	var err error
	if util.JTP(bs[0]) {
		ts, err = goutil.JsonDecode[*TimStream](bs[1:])
	} else {
		ts, err = goutil.TDecode(bs[1:], &TimStream{})
	}
	if err == nil {
		ts.ID = goutil.UUID64()
	}
	return
}

func blockNode(node string, t int64) {
	if t < 0 {
		amr.DelBlock(node)
	} else {
		amr.PutBlock(node, t)
		wsware.delnode(node)
	}
}

func isblock(node string) bool {
	return amr.GetBlock(node) > 0
}

func detect(nodes []string) {
	for _, node := range nodes {
		wsware.detect(node)
	}
}

//func loginstat(node string, on bool, tid *Tid, wsId int64) {
//	defer util.Recover()
//	if sys.Conf.Notice != nil && sys.Conf.Notice.Loginstat != nil {
//		type tk struct {
//			Node   string `json:"node"`
//			Active bool   `json:"active"`
//			Tid    *Tid   `json:"tid"`
//			WsId   int64  `json:"wsId"`
//		}
//		httpclient.Post2(JsonEncode(&tk{Node: node, Active: on, Tid: tid, WsId: wsId}), true, *sys.Conf.Notice.Loginstat)
//	}
//}
