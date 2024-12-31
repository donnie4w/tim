package tc

import (
	goutil "github.com/donnie4w/gofer/util"
	"github.com/donnie4w/tim/adm"
	"github.com/donnie4w/tim/stub"
	"github.com/donnie4w/tim/sys"
	"github.com/donnie4w/tim/util"
	"github.com/donnie4w/tlnet"
)

func wsAdmConfig() *tlnet.WebsocketConfig {
	wc := &tlnet.WebsocketConfig{}
	wc.Origin = sys.ORIGIN
	wc.OnError = func(self *tlnet.Websocket) {
	}
	wc.OnOpen = func(hc *tlnet.HttpContext) {
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
	case sys.ADMRESETAUTH:
		if nab, err := goutil.TDecode(bs[1:], stub.NewAuthBean()); err == nil {
			adm.Admhandler.ResetAuth(nab)
		}
	case sys.ADMAUTH:
		if ab, err := goutil.TDecode(bs[1:], stub.NewAuthBean()); err == nil {
			adm.Admhandler.Auth(ab)
		}
	case sys.ADMTOKEN:
		if at, err := goutil.TDecode(bs[1:], stub.NewAdmToken()); err == nil {
			adm.Admhandler.Token(at)
		}
	case sys.ADMOSMESSAGE:
		if at, err := goutil.TDecode(bs[1:], stub.NewAdmToken()); err == nil {
			adm.Admhandler.Token(at)
		}
	case sys.ADMPROXYMESSAGE:
		if apm, err := goutil.TDecode(bs[1:], stub.NewAdmProxyMessage()); err == nil {
			adm.Admhandler.ProxyMessage(apm)
		}
	case sys.ADMREGISTER:
		if am, err := goutil.TDecode(bs[1:], stub.NewAdmMessage()); err == nil {
			adm.Admhandler.OsMessage(am)
		}
	case sys.ADMMODIFYUSERINFO:
		if amui, err := goutil.TDecode(bs[1:], stub.NewAdmModifyUserInfo()); err == nil {
			adm.Admhandler.ModifyUserInfo(amui)
		}
	case sys.ADMMODIFYROOMINFO:
		if arb, err := goutil.TDecode(bs[1:], stub.NewAdmRoomBean()); err == nil {
			adm.Admhandler.ModifyRoomInfo(arb)
		}
	case sys.ADMBLOCKUSER:
		if abu, err := goutil.TDecode(bs[1:], stub.NewAdmBlockUser()); err == nil {
			adm.Admhandler.BlockUser(abu)
		}
	case sys.ADMBLOCKLIST:
		if abl := adm.Admhandler.BlockList(); abl != nil {
			admwsware.SendWs(ws.Id, abl, sys.ADMBLOCKLIST)
		}
	case sys.ADMONLINEUSER:
		if aou, err := goutil.TDecode(bs[1:], stub.NewAdmOnlineUser()); err == nil {
			adm.Admhandler.OnlineUser(aou)
		}
	case sys.ADMVROOM:
		if avb, err := goutil.TDecode(bs[1:], stub.NewAdmVroomBean()); err == nil {
			adm.Admhandler.Vroom(avb)
		}
	case sys.ADMTIMROOM:
		if arb, err := goutil.TDecode(bs[1:], stub.NewAdmRoomReq()); err == nil {
			adm.Admhandler.TimRoom(arb)
		}
	case sys.ADMDETECT:
		if adb, err := goutil.TDecode(bs[1:], stub.NewAdmDetectBean()); err == nil {
			adm.Admhandler.Detect(adb)
		}
	}
}

func isAuth(ws *tlnet.Websocket) (_b bool) {
	return admwsware.hasws(ws)
}
