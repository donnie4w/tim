// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package tc

import (
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
		mq.Unsub(mq.ONLINE, self.Id)
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
	//case sys.ADMRESETAUTH:
	//	if nab, err := goutil.TDecode(bs[1:], stub.NewAuthBean()); err == nil {
	//		admwsware.SendWs(ws.Id, adm.Admhandler.ModifyPwd(nab), sys.ADMRESETAUTH)
	//	}
	case sys.ADMAUTH:
		if ab, err := goutil.TDecode(bs[1:], stub.NewAuthBean()); err == nil {
			ack := adm.Admhandler.Auth(ab)
			if ack != nil && ack.GetOk() {
				admwsware.Addws(ws)
			}
			admwsware.SendWs(ws.Id, ack, sys.ADMAUTH)
		}
	//case sys.ADMTOKEN:
	//	if at, err := goutil.TDecode(bs[1:], stub.NewAdmToken()); err == nil {
	//		admwsware.SendWs(ws.Id, adm.Admhandler.Token(at), sys.ADMTOKEN)
	//	}
	//case sys.ADMPROXYMESSAGE:
	//	if apm, err := goutil.TDecode(bs[1:], stub.NewAdmProxyMessage()); err == nil {
	//		admwsware.SendWs(ws.Id, adm.Admhandler.ProxyMessage(apm), sys.ADMPROXYMESSAGE)
	//	}
	//case sys.ADMREGISTER:
	//	if am, err := goutil.TDecode(bs[1:], stub.NewAdmMessage()); err == nil {
	//		admwsware.SendWs(ws.Id, adm.Admhandler.TimMessage(am), sys.ADMREGISTER)
	//	}
	//case sys.ADMMODIFYUSERINFO:
	//	if amui, err := goutil.TDecode(bs[1:], stub.NewAdmModifyUserInfo()); err == nil {
	//		admwsware.SendWs(ws.Id, adm.Admhandler.ModifyUserInfo(amui), sys.ADMMODIFYUSERINFO)
	//	}
	//case sys.ADMMODIFYROOMINFO:
	//	if arb, err := goutil.TDecode(bs[1:], stub.NewAdmRoomBean()); err == nil {
	//		admwsware.SendWs(ws.Id, adm.Admhandler.ModifyRoomInfo(arb), sys.ADMMODIFYROOMINFO)
	//	}
	//case sys.ADMBLOCKUSER:
	//	if abu, err := goutil.TDecode(bs[1:], stub.NewAdmBlockUser()); err == nil {
	//		admwsware.SendWs(ws.Id, adm.Admhandler.BlockUser(abu), sys.ADMBLOCKUSER)
	//	}
	//case sys.ADMBLOCKLIST:
	//	if abl := adm.Admhandler.BlockList(); abl != nil {
	//		admwsware.SendWs(ws.Id, abl, sys.ADMBLOCKLIST)
	//	}
	//case sys.ADMONLINEUSER:
	//	if aou, err := goutil.TDecode(bs[1:], stub.NewAdmOnlineUser()); err == nil {
	//		admwsware.SendWs(ws.Id, adm.Admhandler.OnlineUser(aou), sys.ADMONLINEUSER)
	//	}
	//case sys.ADMVROOM:
	//	if avb, err := goutil.TDecode(bs[1:], stub.NewAdmVroomBean()); err == nil {
	//		admwsware.SendWs(ws.Id, adm.Admhandler.Vroom(avb), sys.ADMVROOM)
	//	}
	//case sys.ADMTIMROOM:
	//	if arb, err := goutil.TDecode(bs[1:], stub.NewAdmRoomReq()); err == nil {
	//		admwsware.SendWs(ws.Id, adm.Admhandler.TimRoom(arb), sys.ADMTIMROOM)
	//	}
	//case sys.ADMDETECT:
	//	if adb, err := goutil.TDecode(bs[1:], stub.NewAdmDetectBean()); err == nil {
	//		admwsware.SendWs(ws.Id, adm.Admhandler.Detect(adb), sys.ADMDETECT)
	//	}
	case sys.ADMSUB:
		if adb, err := goutil.TDecode(bs[1:], stub.NewAdmSubBean()); err == nil {
			switch mq.TopicType(adb.GetSubType()) {
			case mq.ONLINE:
				mq.Sub(mq.ONLINE, ws.Id, func(a any) {
					adb.Bs = a.(*stub.AdmSubBean).GetBs()
					admwsware.SendWs(ws.Id, adb, sys.ADMSUB)
				})
			}
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
						} else {
							break
						}
					}
				}
			}()
		}
	}
}
