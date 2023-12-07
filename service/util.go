// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim
package service

import (
	"time"

	"github.com/donnie4w/gofer/cache"
	. "github.com/donnie4w/gofer/hashmap"
	. "github.com/donnie4w/gofer/util"
	"github.com/donnie4w/tim/data"
	. "github.com/donnie4w/tim/stub"
	"github.com/donnie4w/tim/sys"
	"github.com/donnie4w/tim/util"
)

var chatIdTempCache = cache.NewLruCache[int64](1 << 15)
var tokenTempCache = cache.NewLruCache[*Tid](1 << 13)
var userCache = NewLimitMap[uint32, int8](1 << 19)
var groupCache = NewLimitMap[uint32, int8](1 << 19)
var blockm = NewMap[string, int64]()

func init() {
	sys.PingHandle = service.ping
	sys.AckHandle = service.ack
	sys.RegisterHandle = service.register
	sys.TokenHandle = service.token
	sys.AuthHandle = service.auth
	sys.Interrupt = service.interrupt
	sys.MessageHandle = service.message
	sys.PresenceHandle = service.presence
	sys.PullMessageHandle = service.pullmessage
	sys.OfflinemsgHandle = service.offlineMsg
	sys.BroadpresenceHandle = service.broadpresence
	sys.BusinessHandle = service.business
	sys.VRoomHandle = service.vroomprocess
	sys.StreamHandle = service.stream
	sys.NodeInfoHandle = service.nodeinfo
	sys.OsModify = service.sysmodifyauth
	sys.OsMessage = sysMessage
	sys.OsUserBean = service.osuserbean
	sys.OsRoom = service.osnewgroup
	sys.OsRoomBean = service.osModifygroupInfo
	sys.TimMessageProcessor = timMessage
	sys.TimPresenceProcessor = timPresence
	sys.TimSteamProcessor = timStream
	sys.HasNode = wsware.hasUser
	sys.HasWs = wsware.hasws
	sys.DelWs = wsware.delws
	sys.WssLen = wsware.wsLen
	sys.WssList = wsware.wssList
	sys.WssInfo = wsware.wssInfo
	sys.OsToken = service.ostoken
	sys.OsRegister = service.osregister
	sys.SendNode = wsware.SendNode
	sys.SendWs = wsware.SendWs
	sys.BlockUser = blocku
	sys.BlockList = blocklist
	go ticker()
}

func token() (_r int64) {
	return RandId()
}

func existUser(tid *Tid) (_r bool) {
	if tid != nil {
		f := CRC32(Int64ToBytes(int64(util.CreateUUIDByTid(tid))))
		if _r = userCache.Has(f); !_r {
			if _r = data.Handler.ExistUser(tid.Node); _r {
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
			f := CRC32(Int64ToBytes(int64(util.CreateUUID(u, domain))))
			if _r = userCache.Has(f); !_r {
				if _r = data.Handler.ExistUser(u); _r {
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
		if _r = groupCache.Has(f); !_r {
			if _r = data.Handler.ExistGroup(tid.Node); _r {
				groupCache.Put(f, 0)
			}
		}
	} else {
		_r = true
	}
	return
}

func authTidNode(fTid, tTid *Tid) (ok bool) {
	defer util.Recover()
	if sys.Conf.MessageNoauth {
		return true
	}
	if sys.Conf.CacheExpireTime > 0 {
		cid := util.ChatIdByNode(fTid.Node, tTid.Node, fTid.Domain)
		if t, b := chatIdTempCache.Get(cid); !b || t+int64(sys.Conf.CacheExpireTime*int(time.Second)) < time.Now().Unix() {
			if ok = data.Handler.AuthUserAndUser(fTid.Node, tTid.Node, fTid.Domain); ok {
				chatIdTempCache.Add(cid, time.Now().Unix())
			}
		} else {
			chatIdTempCache.Add(cid, time.Now().Unix())
			ok = true
		}
	} else {
		ok = data.Handler.AuthUserAndUser(fTid.Node, tTid.Node, fTid.Domain)
	}
	return
}

func authGroup(gnode, unode string, domain *string) (ok bool) {
	defer util.Recover()
	if sys.Conf.MessageNoauth {
		return true
	}
	if sys.Conf.CacheExpireTime > 0 {
		rid := util.RelateIdForGroup(gnode, unode, domain)
		if t, b := chatIdTempCache.Get(rid); !b || t+int64(sys.Conf.CacheExpireTime*int(time.Second)) < time.Now().Unix() {
			if ok, _ = data.Handler.AuthGroupAndUser(gnode, unode, domain); ok {
				chatIdTempCache.Add(rid, time.Now().Unix())
			}
		} else {
			chatIdTempCache.Add(rid, time.Now().Unix())
			ok = true
		}
	} else {
		ok, _ = data.Handler.AuthGroupAndUser(gnode, unode, domain)
	}
	return
}

func newTimMessage(bs []byte) (tm *TimMessage) {
	if util.JTP(bs[0]) {
		tm, _ = JsonDecode[*TimMessage](bs[1:])
	} else {
		tm, _ = TDecode(bs[1:], &TimMessage{})
	}
	if tm.ID == nil {
		id := RandId()
		tm.ID = &id
	}
	t := time.Now().UnixNano()
	tm.Timestamp = &t
	return
}

func newTimPresence(bs []byte) (tp *TimPresence) {
	if util.JTP(bs[0]) {
		tp, _ = JsonDecode[*TimPresence](bs[1:])
	} else {
		tp, _ = TDecode(bs[1:], &TimPresence{})
	}
	if tp.ID == nil {
		id := RandId()
		tp.ID = &id
	}
	tp.Offline = nil
	return
}

func checkTid(tid *Tid) (_r bool) {
	if sys.UseDefaultDB() && tid != nil {
		return util.CheckNode(tid.Node)
	}
	return true
}

func checkNode(node string) (_r bool) {
	if sys.UseDefaultDB() && node != "" {
		return util.CheckNode(node)
	}
	return true
}

func checkList(ls []string) (_r bool) {
	if sys.UseDefaultDB() && ls != nil {
		for _, u := range ls {
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
	if util.JTP(bs[0]) {
		ts, _ = JsonDecode[*TimStream](bs[1:])
	} else {
		ts, _ = TDecode(bs[1:], &TimStream{})
	}
	ts.ID = RandId()
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

func ticker() {
	tk := time.NewTicker(time.Second << 4)
	for {
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
