// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package timnet

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	. "github.com/donnie4w/gofer/buffer"
	. "github.com/donnie4w/gofer/hashmap"
	"github.com/donnie4w/gofer/keystore"
	. "github.com/donnie4w/gofer/util"
	. "github.com/donnie4w/tim/stub"
	"github.com/donnie4w/tim/sys"
	"github.com/donnie4w/tim/util"
	"github.com/donnie4w/tlnet"
)

func init() {
	go expiredTimer()
	sys.Service.Put(3, timservice)
}

type timService struct {
	isClose  bool
	tlnetTim *tlnet.Tlnet
}

var timservice = &timService{false, tlnet.NewTlnet()}

func (t *timService) Serve() (err error) {
	tls := sys.Conf.Ssl_crt != "" && sys.Conf.Ssl_crt_key != ""
	err = t._serve(sys.IMADDR, tls, sys.Conf.Ssl_crt, sys.Conf.Ssl_crt_key)
	return
}

func (t *timService) Close() (err error) {
	defer util.Recover()
	t.isClose = true
	err = t.tlnetTim.Close()
	return
}

func (t *timService) _serve(listen int, TLS bool, serverCrt, serverKey string) (err error) {
	defer util.Recover()
	tlnet.SetLogOFF()
	addr := strings.TrimSpace(fmt.Sprint(":", listen))
	t.tlnetTim.Handle("/tim2", httpHandler)
	t.tlnetTim.HandleWebSocketBindConfig("/tim", wsHandler, wsConfig())
	sys.FmtLog("tim uuid[", sys.UUID, "]")

	if TLS {
		if IsFileExist(serverCrt) && IsFileExist(serverKey) {
			sys.FmtLog("tim listen tls[", addr, "]")
			err = t.tlnetTim.HttpStartTLS(addr, serverCrt, serverKey)
		} else {
			sys.FmtLog("tim listen tls[", addr, "]")
			err = t.tlnetTim.HttpStartTlsBytes(addr, []byte(keystore.ServerCrt), []byte(keystore.ServerKey))
		}
	}
	if !t.isClose {
		sys.FmtLog("tim listen[", addr, "]")
		err = t.tlnetTim.HttpStart(addr)
	}

	if !t.isClose && err != nil {
		sys.FmtLog("tim start failed:", err.Error())
		os.Exit(0)
	}
	return
}

var expiredMap = NewLinkedMap[*tlnet.Websocket, int64]()

func wsConfig() *tlnet.WebsocketConfig {
	wc := &tlnet.WebsocketConfig{}
	wc.Origin = sys.ORIGIN
	wc.OnError = func(self *tlnet.Websocket) {
		sys.DelWs(self)
		rmBigDataId(self.Id)
	}
	wc.OnOpen = func(hc *tlnet.HttpContext) {
		if !handle_connlimit(hc) {
			expiredMap.Put(hc.WS, time.Now().Unix())
		}
	}
	return wc
}

func httpHandler(hc *tlnet.HttpContext) {
	defer util.Recover()
	if sys.ORIGIN != "" && hc.ReqInfo.Header.Get("Origin") != sys.ORIGIN {
		return
	}
	bs := hc.RequestBody()
	if isForBitIface(bs[0]) {
		return
	}
	if overMaxData(nil, int64(len(bs))) {
		return
	}
	if tasklimit() {
		handle_http(hc, bs)
	} else {
		hc.ResponseBytes(http.StatusInternalServerError, nil)
	}
}

func handle_http(hc *tlnet.HttpContext, bs []byte) {
	j := util.JTP(bs[0])
	switch sys.TIMTYPE(bs[0] & 0x7f) {
	case sys.TIMREGISTER:
		if isForBidRegister() {
			hc.ResponseBytes(http.StatusForbidden, nil)
			return
		}
		if node, err := sys.RegisterHandle(bs); err == nil {
			hc.ResponseBytes(0, reTimAck(j, &TimAck{Ok: true, TimType: int8(sys.TIMREGISTER), N: &node}))
		} else {
			hc.ResponseBytes(0, reTimAck(j, &TimAck{Ok: false, TimType: int8(sys.TIMREGISTER), Error: err.TimError()}))
		}
	case sys.TIMTOKEN:
		if isForBidToken() {
			hc.ResponseBytes(http.StatusForbidden, nil)
			return
		}
		if t, err := sys.TokenHandle(bs); err == nil {
			hc.ResponseBytes(0, reTimAck(j, &TimAck{Ok: true, TimType: int8(sys.TIMTOKEN), T: &t}))
		} else {
			hc.ResponseBytes(0, reTimAck(j, &TimAck{Ok: false, TimType: int8(sys.TIMTOKEN), Error: err.TimError()}))
		}
	}
}

func reTimAck(j bool, ta *TimAck) []byte {
	buf := NewBuffer()
	buf.WriteByte(byte(sys.TIMACK))
	if j {
		buf.Write(JsonEncode(ta))
	} else {
		buf.Write(TEncode(ta))
	}
	return buf.Bytes()
}

func tasklimit() bool {
	over := 100
	for over > 0 && sys.Stat.Tx() > int64(*sys.Conf.TaskLimit) {
		over--
		<-time.After(time.Second / time.Duration(over))
	}
	return over > 0
}

func wsHandler(hc *tlnet.HttpContext) {
	defer util.Recover()
	if overMaxData(hc.WS, int64(len(hc.WS.Read()))) {
		return
	}
	bs := make([]byte, len(hc.WS.Read()))
	sys.Stat.Ib(int64(len(bs)))
	if isForBitIface(bs[0]) {
		return
	}
	copy(bs, hc.WS.Read())
	if isBigData(hc.WS.Id) {
		addBigData(hc, bs)
	} else {
		if t := sys.TIMTYPE(bs[0] & 0x7f); t == sys.TIMBIGBINARY || t == sys.TIMBIGSTRING || t == sys.TIMBIGBINARYSTREAM {
			parseBigData(hc, bs)
		} else {
			parseWsData(bs, hc)
		}
	}
}

func parseWsData(bs []byte, hc *tlnet.HttpContext) {
	t := sys.TIMTYPE(bs[0] & 0x7f)
	if bs == nil || (t != sys.TIMAUTH && t != sys.TIMPING && !isAuth(hc.WS)) {
		hc.WS.Close()
		return
	}
	if overHz(hc) {
		return
	}
	switch t {
	case sys.TIMMESSAGE, sys.TIMREVOKEMESSAGE, sys.TIMBURNMESSAGE, sys.TIMPRESENCE, sys.TIMSTREAM, sys.TIMBIGSTRING, sys.TIMBIGBINARY, sys.TIMBIGBINARYSTREAM:
		if tasklimit() {
			go handle_core(hc, bs, t)
		} else {
			go handle_err_overload(hc, t)
		}
	default:
		go handle_business(hc, bs, t)
	}
}

func handle_core(hc *tlnet.HttpContext, bs []byte, t sys.TIMTYPE) {
	defer util.Recover()
	sys.Stat.TxDo()
	defer sys.Stat.TxDone()
	var err sys.ERROR
	switch t {
	case sys.TIMMESSAGE, sys.TIMREVOKEMESSAGE, sys.TIMBURNMESSAGE:
		err = sys.MessageHandle(bs, hc.WS)
	case sys.TIMPRESENCE:
		err = sys.PresenceHandle(bs, hc.WS)
	case sys.TIMSTREAM:
		err = sys.StreamHandle(bs, hc.WS)
	case sys.TIMBIGSTRING:
		err = sys.BigStringHandle(bs, hc.WS)
	case sys.TIMBIGBINARY:
		err = sys.BigBinaryHandle(bs, hc.WS)
	case sys.TIMBIGBINARYSTREAM:
		err = sys.BigBinaryStreamHandle(bs, hc.WS)
	}

	if err != nil {
		sys.SendWs(hc.WS.Id, &TimAck{Ok: false, TimType: int8(t), Error: err.TimError()}, sys.TIMACK)
	}
}

func handle_business(hc *tlnet.HttpContext, bs []byte, t sys.TIMTYPE) {
	defer util.Recover()
	sys.Stat.TxDo()
	defer sys.Stat.TxDone()
	var err sys.ERROR
	switch t {
	case sys.TIMAUTH:
		if !connectAuth(bs) {
			hc.WS.Close()
			return
		}
		if err = sys.AuthHandle(bs, hc.WS); err == nil {
			expiredMap.Del(hc.WS)
		} else {
			hc.WS.Send(reTimAck(util.JTP(bs[0]), &TimAck{Ok: false, TimType: int8(t), Error: err.TimError()}))
			return
		}
	case sys.TIMACK:
		err = sys.AckHandle(bs)
	case sys.TIMPING:
		err = sys.PingHandle(hc.WS)
	case sys.TIMOFFLINEMSG:
		err = sys.OfflinemsgHandle(hc.WS)
	case sys.TIMPULLMESSAGE:
		err = sys.PullMessageHandle(bs, hc.WS)
	case sys.TIMBROADPRESENCE:
		err = sys.BroadpresenceHandle(bs, hc.WS)
	case sys.TIMBUSINESS:
		err = sys.BusinessHandle(bs, hc.WS)
	case sys.TIMVROOM:
		err = sys.VRoomHandle(bs, hc.WS)
	case sys.TIMNODES:
		err = sys.NodeInfoHandle(bs, hc.WS)
	default:
		err = sys.ERR_PARAMS
	}
	if err != nil {
		sys.SendWs(hc.WS.Id, &TimAck{Ok: false, TimType: int8(t), Error: err.TimError()}, sys.TIMACK)
	}
}

func handle_err_overload(hc *tlnet.HttpContext, t sys.TIMTYPE) {
	defer util.Recover()
	sys.SendWs(hc.WS.Id, &TimAck{Ok: false, TimType: int8(t), Error: sys.ERR_OVERLOAD.TimError()}, sys.TIMACK)
}

// func c(hc *tlnet.HttpContext,  t sys.TIMTYPE) {
// 	defer util.Recover()
// 	sys.SendWs(hc.WS.Id, &TimAck{Ok: false, TimType: int8(t), Error: sys.ERR_OVERHZ.TimError()}, sys.TIMACK)
// }

func handle_connlimit(hc *tlnet.HttpContext) (_r bool) {
	defer util.Recover()
	if _r = sys.WssLen() > sys.Conf.ConnectLimit; _r {
		hc.WS.Close()
	}
	return
}

func isAuth(ws *tlnet.Websocket) (_b bool) {
	return sys.HasWs(ws)
}

func expiredTimer() {
	t := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-t.C:
			func() {
				defer util.Recover()
				expiredMap.BackForEach(func(k *tlnet.Websocket, v int64) bool {
					if v+3 < time.Now().Unix() {
						k.Close()
						expiredMap.Del(k)
						return true
					}
					return false
				})
			}()
		}
	}
}
