// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package service

import (
	"github.com/donnie4w/tim/cache"
	"time"

	"github.com/donnie4w/gofer/hashmap"
	. "github.com/donnie4w/gofer/util"
	"github.com/donnie4w/tim/data"
	. "github.com/donnie4w/tim/stub"
	"github.com/donnie4w/tim/sys"
	"github.com/donnie4w/tim/util"
)

var userCache = hashmap.NewLimitHashMap[uint32, int8](1 << 19)
var groupCache = hashmap.NewLimitHashMap[uint32, int8](1 << 19)
var blockm = hashmap.NewMap[string, int64]()

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
	sys.BroadpresenceHandle = service.broadpresence
	sys.BusinessHandle = service.business
	sys.VRoomHandle = service.vroomHandle
	sys.StreamHandle = service.streamHandle
	sys.BigStringHandle = service.bigString
	sys.BigBinaryHandle = service.bigBinary
	sys.BigBinaryStreamHandle = service.bigBinaryStreamHandle
	sys.NodeInfoHandle = service.nodeinfo
	sys.OsModify = service.sysmodify
	sys.OsMessage = sysMessage
	sys.OsUserBean = service.osuserbean
	sys.OsRoom = service.osnewgroup
	sys.OsRoomBean = service.osModifygroupInfo
	sys.OsVroomprocess = service.osvroomprocess
	sys.PxMessage = service.pxmessage
	sys.TimMessageProcessor = timMessageProcessor
	sys.TimPresenceProcessor = timPresenceProcessor
	sys.TimSteamProcessor = timStreamProcessor
	sys.HasNode = wsware.hasUser
	sys.HasWs = wsware.hasws
	sys.DelWs = wsware.delws
	sys.WsById = wsware.wsById
	sys.WssLen = wsware.wsLen
	sys.WssList = wsware.wssList
	sys.WssInfo = wsware.wssInfo
	sys.OsToken = service.ostoken
	sys.OsRegister = service.osregister
	sys.SendNode = wsware.SendNode
	sys.SendWs = wsware.SendWs
	sys.OsBlockUser = blocku
	sys.OsBlockList = blocklist
	sys.Detect = detect
	go ticker()
}

// token The effective length is only 32 bits
func token() (_r int64) {
	return int64(UUID32())
}

func existUser(tid *Tid) (_r bool) {
	if tid != nil {
		uuid := util.CreateUUIDByTid(tid)
		f := CRC32(Int64ToBytes(int64(uuid)))
		if _r = userCache.Contains(f); !_r {
			if _r = data.Service.ExistUser(util.UUIDToNode(uuid)); _r {
				userCache.Put(f, 0)
			}
		}
	} else {
		_r = true
	}
	return
}

func existList(ls []string, domain *string) (_r bool) {
	if ls != nil {
		for _, u := range ls {
			uuid := util.CreateUUID(u, domain)
			f := CRC32(Int64ToBytes(int64(uuid)))
			if _r = userCache.Contains(f); !_r {
				if _r = data.Service.ExistUser(util.UUIDToNode(uuid)); _r {
					userCache.Put(f, 0)
				} else {
					_r = false
					break
				}
			}
		}
	} else {
		_r = true
	}
	return
}

func existGroup(tid *Tid) (_r bool) {
	if tid != nil {
		f := CRC32(append([]byte{1}, Int64ToBytes(int64(util.CreateUUIDByTid(tid)))...))
		if _r = groupCache.Contains(f); !_r {
			if _r = data.Service.ExistGroup(tid.Node); _r {
				groupCache.Put(f, 0)
			}
		}
	} else {
		_r = true
	}
	return
}

func AuthUser(fTid, tTid *Tid, readtime bool) (ok bool) {
	defer util.Recover()
	if sys.Conf.MessageNoAuth {
		return true
	}
	if sys.Conf.CacheAuthExpire > 0 && !readtime {
		nano := time.Now().UnixNano()
		cid := util.ChatIdByNode(fTid.Node, tTid.Node, fTid.Domain)
		if t, b := cache.AuthCache.Get(cid); !b || t+int64(sys.Conf.CacheAuthExpire*int(time.Second)) < nano {
			if data.Service.AuthUserAndUser(fTid.Node, tTid.Node, fTid.Domain) {
				cache.AuthCache.Put(cid, nano)
			}
		} else {
			cache.AuthCache.Put(cid, nano)
			ok = true
		}
	} else {
		ok = data.Service.AuthUserAndUser(fTid.Node, tTid.Node, fTid.Domain)
	}
	return
}

func AuthGroup(gnode, unode string, domain *string) (ok bool) {
	defer util.Recover()
	if sys.Conf.MessageNoAuth {
		return true
	}
	if sys.Conf.CacheAuthExpire > 0 {
		rid := util.RelateIdForGroup(gnode, unode, domain)
		if t, b := cache.AuthCache.Get(rid); !b || t+int64(sys.Conf.CacheAuthExpire*int(time.Second)) < time.Now().UnixNano() {
			if ok, _ = data.Service.AuthGroupAndUser(gnode, unode, domain); ok {
				cache.AuthCache.Put(rid, time.Now().UnixNano())
			}
		} else {
			cache.AuthCache.Put(rid, time.Now().UnixNano())
			ok = true
		}
	} else {
		ok, _ = data.Service.AuthGroupAndUser(gnode, unode, domain)
	}
	return
}

func newTimMessage(bs []byte) (tm *TimMessage) {
	var err error
	if util.JTP(bs[0]) {
		tm, err = JsonDecode[*TimMessage](bs[1:])
	} else {
		tm, err = TDecode(bs[1:], &TimMessage{})
	}
	if err == nil {
		if tm.ID == nil {
			id := UUID64()
			tm.ID = &id
		}
		t := time.Now().UnixNano()
		tm.Timestamp = &t
	}
	return
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
		tp, err = JsonDecode[*TimPresence](bs[1:])
	} else {
		tp, err = TDecode(bs[1:], &TimPresence{})
	}
	if err == nil {
		if tp.ID == nil {
			id := UUID64()
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

	if tid != nil && len(tid.Node) > sys.NodeMaxlength {
		return false
	}
	return true
}

func checkNode(node string) (_r bool) {
	if sys.UseBuiltInData() && node != "" {
		return util.CheckNode(node)
	}
	if len(node) > sys.NodeMaxlength {
		return false
	}
	return true
}

func checkList(ls []string) (_r bool) {
	if sys.UseBuiltInData() && ls != nil {
		for _, u := range ls {
			if len(u) > sys.NodeMaxlength {
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
		tr, _ = JsonDecode[*TimReq](bs[1:])
	} else {
		tr, _ = TDecode(bs[1:], &TimReq{})
	}
	return
}

func newTimNodes(bs []byte) (tr *TimNodes) {
	if util.JTP(bs[0]) {
		tr, _ = JsonDecode[*TimNodes](bs[1:])
	} else {
		tr, _ = TDecode(bs[1:], &TimNodes{})
	}
	return
}

func newAuth(bs []byte) (ta *TimAuth) {
	if util.JTP(bs[0]) {
		ta, _ = JsonDecode[*TimAuth](bs[1:])
	} else {
		ta, _ = TDecode(bs[1:], &TimAuth{})
	}
	return
}

func newTimStream(bs []byte) (ts *TimStream) {
	var err error
	if util.JTP(bs[0]) {
		ts, err = JsonDecode[*TimStream](bs[1:])
	} else {
		ts, err = TDecode(bs[1:], &TimStream{})
	}
	if err == nil {
		ts.ID = UUID64()
	}
	return
}

func blocku(node string, t int64) {
	if t < 0 {
		blockm.Del(node)
	} else {
		blockm.Put(node, time.Now().Unix()+t)
		wsware.delnode(node)
	}
}

func isblock(node string) bool {
	if t, ok := blockm.Get(node); ok {
		return t > time.Now().Unix()
	}
	return false
}

func blocklist() map[string]int64 {
	m := map[string]int64{}
	blockm.Range(func(k string, v int64) bool {
		m[k] = v
		return true
	})
	return m
}

func detect(nodes []string) {
	for _, node := range nodes {
		wsware.detect(node)
	}
}

func ticker() {
	tk := time.NewTicker(time.Second << 4)
	for {
		sys.InaccurateTime = time.Now().UnixNano()
		select {
		case <-tk.C:
			func() {
				defer util.Recover()
				blockm.Range(func(k string, v int64) bool {
					if v > 0 && v < time.Now().Unix() {
						blockm.Del(k)
					}
					return true
				})
			}()
		}
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
