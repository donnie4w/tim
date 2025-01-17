// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package tc

import (
	"bytes"
	"github.com/donnie4w/gofer/hashmap"
	goutil "github.com/donnie4w/gofer/util"
	"github.com/donnie4w/tim/adm"
	"github.com/donnie4w/tim/mq"
	"github.com/donnie4w/tim/stub"
	"github.com/donnie4w/tim/sys"
	"github.com/donnie4w/tim/util"
	"github.com/donnie4w/tlnet"
	"time"
)

func init() {
	go expiredTimer()
}

func wsAdmConfig() *tlnet.WebsocketConfig {
	wc := &tlnet.WebsocketConfig{}
	wc.Origin = sys.ORIGIN
	wc.OnError = func(self *tlnet.Websocket) {
		mq.Unsub(mq.ONLINESTATUS, self.Id)
		admwsware.delws(self)
	}
	wc.OnOpen = func(hc *tlnet.HttpContext) {
		expiredMap.Put(hc.WS, time.Now().Unix())
	}
	return wc
}

func wsAdmHandler(hc *tlnet.HttpContext) {
	defer util.Recover()
	bs := make([]byte, len(hc.WS.Read()))
	sys.Stat.Ib(int64(len(bs)))
	copy(bs, hc.WS.Read())
	if t := sys.TIMTYPE(bs[0] & 0x7f); t != sys.ADMAUTH && t != sys.ADMPING && !isAuth(hc.WS) {
		hc.WS.Close()
		return
	}
	go processor(hc.WS, bs)
}

func processor(ws *tlnet.Websocket, bs []byte) {
	defer util.Recover()
	switch sys.TIMTYPE(bs[0] & 0x7f) {
	case sys.ADMPING:
		admwsware.Ping(ws.Id)
	case sys.ADMAUTH:
		if ab, err := goutil.TDecode(bs[1:], stub.NewAuthBean()); err == nil {
			ack := adm.Auth(ab)
			if ack != nil && ack.GetOk() {
				expiredMap.Delete(ws)
				admwsware.Addws(ws)
			}
			admwsware.SendWs(ws, ack, sys.ADMAUTH)
		}
	case sys.ADMSUB:
		if adb, err := goutil.TDecode(bs[1:], stub.NewAdmSubBean()); err == nil {
			switch mq.TopicType(adb.GetSubType()) {
			case mq.ONLINESTATUS:
				mq.Sub(mq.ONLINESTATUS, ws.Id, func(a any) {
					adb.Bs = a.(*stub.AdmSubBean).GetBs()
					admwsware.Send(ws.Id, adb, sys.ADMSUB)
				})
			}
		}
	case sys.ADMSTREAM:
	case sys.ADMBIGSTRING:
	case sys.ADMBIGBINARY:
	case sys.ADMBIGBINARYSTREAM:
		bs = bs[1:]
		vnodeIdx := bytes.IndexByte(bs, sys.SEP_BIN)
		vnode := string(bs[:vnodeIdx])
		bs = bs[vnodeIdx+1:]
		fnodeIdx := bytes.IndexByte(bs, sys.SEP_BIN)
		fnode := string(bs[:fnodeIdx])
		bs = bs[fnodeIdx+1:]
		sid := goutil.UUID64()
		vb := &stub.VBean{StreamId: &sid, Vnode: vnode, Rnode: &fnode, Body: bs, Rtype: int8(sys.VROOM_MESSAGE)}
		if b, _ := sys.TimSteamProcessor(vb, sys.TRANS_GOAL); !b {
			ok, n, t := false, vnode, int64(int32(goutil.FNVHash32(goutil.Int64ToBytes(sys.UUID))))
			ack := stub.NewAdmAck()
			ack.Ok = &ok
			ack.N = &n
			ack.T = &t
			admwsware.Send(ws.Id, ack, sys.ADMBIGBINARYSTREAM)
		}
	}
}

func isAuth(ws *tlnet.Websocket) (_b bool) {
	return admwsware.hasws(ws)
}

var expiredMap = hashmap.NewLinkedHashMap[*tlnet.Websocket, int64](1<<63 - 1)

func expiredTimer() {
	t := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-t.C:
			func() {
				defer util.Recover()
				iter := expiredMap.Iterator(false)
				for {
					if k, v, ok := iter.Next(); ok {
						if v+5 < time.Now().Unix() {
							k.Close()
							admwsware.delws(k)
							expiredMap.Delete(k)
							continue
						}
					}
					break
				}
			}()
		}
	}
}
