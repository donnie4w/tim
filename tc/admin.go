// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package tc

import (
	"fmt"
	"github.com/donnie4w/gofer/hashmap"
	KS "github.com/donnie4w/gofer/keystore"
	goutil "github.com/donnie4w/gofer/util"
	"github.com/donnie4w/tim/keystore"
	"github.com/donnie4w/tim/log"
	"github.com/donnie4w/tim/sys"
	"github.com/donnie4w/tim/util"
	"github.com/donnie4w/tlnet"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func init() {
	sys.Service(sys.INIT_TC, newAdminService())
}

type adminService struct {
	isClose bool
	server  *tlnet.Tlnet
}

func newAdminService() *adminService {
	return &adminService{false, tlnet.NewTlnet()}
}

func (t *adminService) Serve() (err error) {
	if strings.TrimSpace(sys.Conf.PprofAddr) != "" {
		go tlDebug()
		<-time.After(500 * time.Millisecond)
	}
	if sys.Conf.WebAdminListen != "" {
		if !sys.Conf.NoInit {
			initAccount()
		}
		go t._serve(strings.TrimSpace(sys.Conf.WebAdminListen), !sys.Conf.WebAdminNoTls, sys.Conf.Ssl_crt, sys.Conf.Ssl_crt_key)
	} else {
		log.FmtPrint("No WebAdmin Service")
	}
	return
}

func (t *adminService) Close() (err error) {
	defer util.Recover()
	if sys.Conf.WebAdminListen != "" {
		t.isClose = true
		err = t.server.Close()
	}
	return
}

func (t *adminService) _serve(addr string, TLS bool, serverCrt, serverKey string) (err error) {
	defer util.Recover()
	if addr, err = util.ParseAddr(addr); err != nil {
		log.FmtPrint("WebAdmin Service ParseAddr error:", err.Error())
		os.Exit(1)
	}
	t.server.Handle("/login", loginHandler)
	t.server.Handle("/init", initHandler)
	t.server.Handle("/lang", langHandler)
	t.server.Handle("/", initHandler)
	t.server.Handle("/bootstrap.css", cssHandler)
	t.server.Handle("/bootstrap.min.js", jsHandler)
	t.server.HandleWithFilter("/sysvar", loginFilter(), sysVarHtml)

	t.server.HandleWithFilter("/timResetAuth", authFilter(), timResetAuthHandler)
	t.server.HandleWithFilter("/timToken", authFilter(), timTokenHandler)
	t.server.HandleWithFilter("/timOsMessage", authFilter(), timOsMessageHandler)
	t.server.HandleWithFilter("/timMessage", authFilter(), timMessageHandler)
	t.server.HandleWithFilter("/timRegister", authFilter(), timRegisterHandler)
	t.server.HandleWithFilter("/timModifyUserInfo", authFilter(), timModifyUserInfoHnadler)
	t.server.HandleWithFilter("/timBlockUser", authFilter(), timBlockUserHandler)
	t.server.HandleWithFilter("/timOnline", authFilter(), timOnlineHandler)
	t.server.HandleWithFilter("/timVroom", authFilter(), timVroomHandler)
	t.server.HandleWithFilter("/timNewRoom", authFilter(), timNewRoomHandler)
	t.server.HandleWithFilter("/timDetect", authFilter(), timDetectHandler)
	t.server.HandleWithFilter("/timModifyRoomInfo", authFilter(), timModifyRoomInfoHandler)
	t.server.HandleWithFilter("/timAddroster", authFilter(), timAddrosterHandler)
	t.server.HandleWithFilter("/timRmroster", authFilter(), timRmrosterHandler)
	t.server.HandleWithFilter("/timBlockroster", authFilter(), timBlockrosterHandler)
	t.server.HandleWithFilter("/timBlockRosterList", authFilter(), timBlockRosterListHandler)
	t.server.HandleWithFilter("/timAddRoom", authFilter(), timAddRoomHandler)
	t.server.HandleWithFilter("/timPullInRoom", authFilter(), timPullInRoomHandler)
	t.server.HandleWithFilter("/timRejectRoom", authFilter(), timRejectRoomHandler)
	t.server.HandleWithFilter("/timKickRoom", authFilter(), timKickRoomHandler)
	t.server.HandleWithFilter("/timLeaveRoom", authFilter(), timLeaveRoomHandler)
	t.server.HandleWithFilter("/timCancelRoom", authFilter(), timCancelRoomHandler)
	t.server.HandleWithFilter("/timBlockRoom", authFilter(), timBlockRoomHandler)
	t.server.HandleWithFilter("/timBlockRoomMember", authFilter(), timBlockRoomMemberHandler)
	t.server.HandleWithFilter("/timBlockRoomList", authFilter(), timBlockRoomListHandler)
	t.server.HandleWithFilter("/timBlockRoomMemberlist", authFilter(), timBlockRoomMemberlistHandler)

	t.server.HandleWithFilter("/monitor", loginFilter(), monitorHtml)
	t.server.HandleWebSocketBindConfig("/monitorData", mntHandler, mntConfig())
	t.server.HandleWithFilter("/data", loginFilter(), dataMonitorHtml)
	t.server.HandleWebSocketBindConfig("/ddmonitorData", ddmntHandler, ddmntConfig())

	t.server.HandleWebSocketBindConfig("/tim", wsAdmHandler, wsAdmConfig())

	if TLS {
		if goutil.IsFileExist(serverCrt) && goutil.IsFileExist(serverKey) {
			log.FmtPrint("WebAdmin Service start tls[", addr, "]")
			err = t.server.HttpsStart(addr, serverCrt, serverKey)
		} else {
			log.FmtPrint("WebAdmin Service start tls by bytes[", addr, "]")
			err = t.server.HttpsStartWithBytes(addr, []byte(KS.ServerCrt), []byte(KS.ServerKey))
		}
	}
	if !t.isClose {
		log.FmtPrint("WebAdmin Service start[", addr, "]")
		err = t.server.HttpStart(addr)
	}
	if !t.isClose && err != nil {
		log.FmtPrint("WebAdmin Service start failed:", err.Error())
		os.Exit(1)
	}
	return
}

var sessionMap = hashmap.NewMapL[string, *KS.UserBean]()

func loginFilter() (f *tlnet.Filter) {
	defer util.Recover()
	f = tlnet.NewFilter()
	f.AddIntercept(".*?", func(hc *tlnet.HttpContext) bool {
		if len(keystore.Admin.AdminList()) > 0 {
			if !isLogin(hc) {
				hc.Redirect("/login")
				return true
			}
		} else {
			hc.Redirect("/init")
			return true
		}
		return false
	})

	f.AddIntercept(`[^\s]+`, func(hc *tlnet.HttpContext) bool {
		if hc.PostParamTrimSpace("atype") != "" && !isAdmin(hc) {
			hc.ResponseString(resultHtml("Permission Denied"))
			return true
		}
		return false
	})
	return
}

func authFilter() (f *tlnet.Filter) {
	defer util.Recover()
	f = tlnet.NewFilter()
	f.AddIntercept(".*?", func(hc *tlnet.HttpContext) bool {
		name := hc.Request().Header.Get("username")
		pwd := hc.Request().Header.Get("password")
		if _r, ok := keystore.Admin.GetAdmin(name); ok {
			if strings.EqualFold(_r.Pwd, goutil.Md5Str(pwd)) {
				return false
			}
		}
		hc.ResponseBytes(http.StatusUnauthorized, nil)
		return true
	})
	return
}

func getSessionid() string {
	return fmt.Sprint("t", goutil.CRC32(goutil.Int64ToBytes(sys.UUID)))
}

func getLangId() string {
	return fmt.Sprint("l", goutil.CRC32(goutil.Int64ToBytes(sys.UUID)))
}

func isLogin(hc *tlnet.HttpContext) (isLogin bool) {
	if len(keystore.Admin.AdminList()) > 0 {
		if _r, err := hc.GetCookie(getSessionid()); err == nil && sessionMap.Has(_r) {
			isLogin = true
		}
	}
	return
}

func isAdmin(hc *tlnet.HttpContext) (_r bool) {
	if c, err := hc.GetCookie(getSessionid()); err == nil {
		if u, ok := sessionMap.Get(c); ok {
			_r = u.Type == 1
		}
	}
	return
}

/***********************************************************************/
func langHandler(hc *tlnet.HttpContext) {
	defer util.Recover()
	lang := hc.GetParamTrimSpace("lang")
	if lang == "en" || lang == "zh" {
		hc.SetCookie(getLangId(), lang, "/", 86400)
	}
	hc.Redirect("/")
}

func getLang(hc *tlnet.HttpContext) LANG {
	if lang, err := hc.GetCookie(getLangId()); err == nil {
		if lang == "zh" {
			return ZH
		} else if lang == "en" {
			return EN
		}
	}
	return ZH
}

func cssHandler(hc *tlnet.HttpContext) {
	hc.Writer().Header().Add("Content-Type", "text/html")
	textTplByText(cssContent(), nil, hc)
}

func jsHandler(hc *tlnet.HttpContext) {
	hc.Writer().Header().Add("Content-Type", "text/html")
	textTplByText(jsContent(), nil, hc)
}

/***********************************************************************/
func initHandler(hc *tlnet.HttpContext) {
	defer func() {
		if err := recover(); err != nil {
			hc.ResponseString(resultHtml("server error:", err))
		}
	}()
	if len(keystore.Admin.AdminList()) > 0 && !isLogin(hc) {
		hc.Redirect("/login")
		return
	}
	if _type := hc.GetParam("type"); _type != "" {
		isadmin := isAdmin(hc)
		if _type == "1" {
			if name, pwd, _type := hc.PostParamTrimSpace("adminName"), hc.PostParamTrimSpace("adminPwd"), hc.PostParamTrimSpace("adminType"); name != "" && pwd != "" {
				if n := len(keystore.Admin.AdminList()); (n > 0 && isadmin) || n == 0 {
					alterType := false
					if t, err := strconv.Atoi(_type); err == nil {
						if _r, err := hc.GetCookie(getSessionid()); err == nil && sessionMap.Has(_r) {
							if u, ok := sessionMap.Get(_r); ok && u.Name == name && t != int(u.Type) {
								alterType = true
							}
						}
						if !alterType {
							keystore.Admin.PutAdmin(name, pwd, int8(t))
						}
					}
				} else {
					goto DENIED
				}
			}
		} else if _type == "2" && isLogin(hc) {
			if isadmin {
				if name := hc.PostParamTrimSpace("adminName"); name != "" {
					if u, ok := keystore.Admin.GetAdmin(name); ok && u.Type == 1 {
						i, j := 0, 0
						for _, s := range keystore.Admin.AdminList() {
							if _u, _ := keystore.Admin.GetAdmin(s); _u.Type == 1 {
								i++
							} else if _u.Type == 2 {
								j++
							}
						}
						if j > 0 && i == 1 {
							hc.ResponseString(resultHtml("failed,There cannot be only Observed users"))
							return
						}
					}
					keystore.Admin.DelAdmin(name)
					sessionMap.Range(func(k string, v *KS.UserBean) bool {
						if v.Name == name {
							sessionMap.Del(k)
						}
						return true
					})
				}
			} else {
				goto DENIED
			}
		}
		hc.Redirect("/init")
		return
	} else {
		initHtml(hc)
		return
	}
DENIED:
	hc.ResponseString(resultHtml("Permission Denied"))
}

func loginHandler(hc *tlnet.HttpContext) {
	defer func() {
		if err := recover(); err != nil {
			hc.ResponseString(resultHtml("server error:", err))
		}
	}()
	if hc.PostParamTrimSpace("type") == "1" {
		name, pwd := hc.PostParamTrimSpace("name"), hc.PostParamTrimSpace("pwd")
		if _r, ok := keystore.Admin.GetAdmin(name); ok {
			if strings.EqualFold(_r.Pwd, goutil.Md5Str(pwd)) {
				sid := goutil.Md5Str(fmt.Sprint(time.Now().UnixNano()))
				sessionMap.Put(sid, _r)
				hc.SetCookie(getSessionid(), sid, "/", 86400)
				hc.Redirect("/")
				return
			}
		}
		hc.ResponseString(resultHtml("Login Failed"))
		return
	}
	loginHtml(hc)
}

/*****************************************************************************/
func initHtml(hc *tlnet.HttpContext) {
	defer func() {
		if err := recover(); err != nil {
			hc.ResponseString(resultHtml("server error:", err))
		}
	}()
	_isAdmin := isAdmin(hc)
	show, init, sc := "", false, _isAdmin
	if len(keystore.Admin.AdminList()) == 0 {
		show, init, sc = "no user is created for admin, create a management user first", true, true
	}
	av := &AdminView{Show: show, Init: init, ShowCreate: sc}
	if isLogin(hc) {
		m := make(map[string]string, 0)
		for _, s := range keystore.Admin.AdminList() {
			if u, ok := keystore.Admin.GetAdmin(s); ok {
				if _isAdmin && u.Type == 1 {
					m[s] = "Admin"
				} else if u.Type == 2 {
					m[s] = "Observed"
				}
			}
		}
		av.AdminUser = m
	}
	tplToHtml(getLang(hc), INIT, av, hc)
}

func loginHtml(hc *tlnet.HttpContext) {
	defer func() {
		if err := recover(); err != nil {
			hc.ResponseString(resultHtml("server error:", err))
		}
	}()
	tplToHtml(getLang(hc), LOGIN, []byte{}, hc)
}

func initAccount() {
	if len(keystore.Admin.AdminList()) == 0 {
		keystore.Admin.PutAdmin(sys.DefaultAccount[0], sys.DefaultAccount[1], 1)
	}
}

/********************************************************/
func sysVarHtml(hc *tlnet.HttpContext) {
	defer func() {
		if err := recover(); err != nil {
			hc.ResponseString(resultHtml("server error:", err))
		}
	}()
	//rn := sys.GetRemoteNode()
	//sort.Slice(rn, func(i, j int) bool { return rn[i].UUID > rn[j].UUID })
	svv := &SysVarView{Show: "", RN: nil}
	if _type := hc.PostParamTrimSpace("atype"); _type != "" {
		if _type == "1" {
			_addr := hc.PostParamTrimSpace("addr")
			if addr := strings.Trim(_addr, " "); addr != "" {
				if err := sys.AddNode(addr); err != nil {
					hc.ResponseString(resultHtml("Failed :", err.Error()))
					return
				} else {
					<-time.After(1000 * time.Millisecond)
					svv.Show = "ADD NODE [ " + _addr + " ]"
				}
			}
		}
	}
	svv.SYS = sysvar()
	tplToHtml(getLang(hc), SYSVAR, svv, hc)
}

func sysvar() (s *SysVar) {
	s = &SysVar{}
	s.StartTime = fmt.Sprint(sys.STARTTIME)
	s.Time = fmt.Sprint(time.Now())
	s.UUID = sys.UUID
	//s.CSNUM = int32(len(sys.GetALLUUIDS()))
	s.ADDR = fmt.Sprint(sys.Conf.CsAccess)
	s.ADMINADDR = sys.Conf.WebAdminListen
	//aus := append([]int64{}, sys.GetALLUUIDS()...)
	//sort.Slice(aus, func(i, j int) bool { return aus[i] > aus[j] })
	//s.ALLUUIDS = fmt.Sprint(aus)
	return
}

/**********************************************************/
func dataMonitorHtml(hc *tlnet.HttpContext) {
	tplToHtml(getLang(hc), DATA, nil, hc)
}

func ddmntConfig() (wc *tlnet.WebsocketConfig) {
	wc = &tlnet.WebsocketConfig{}
	wc.OnOpen = func(hc *tlnet.HttpContext) {
		if !isLogin(hc) {
			hc.WS.Close()
			return
		}
	}
	return
}

func ddmntHandler(hc *tlnet.HttpContext) {
	defer util.Recover()
	s := string(hc.WS.Read())
	if t, err := strconv.Atoi(s); err == nil {
		if t < 1 {
			t = 1
		}
		for hc.WS.Error == nil {
			if j, err := ddmonitorToJson(); err == nil {
				hc.WS.Send(j)
			}
			<-time.After(time.Duration(t) * time.Second)
		}
	}
}

/**********************************************************/
func monitorHtml(hc *tlnet.HttpContext) {
	tplToHtml(getLang(hc), MONITOR, nil, hc)
}

func mntConfig() (wc *tlnet.WebsocketConfig) {
	wc = &tlnet.WebsocketConfig{}
	wc.OnOpen = func(hc *tlnet.HttpContext) {
		if !isLogin(hc) {
			hc.WS.Close()
			return
		}
	}
	return
}

func mntHandler(hc *tlnet.HttpContext) {
	defer util.Recover()
	s := string(hc.WS.Read())
	if t, err := strconv.Atoi(s); err == nil {
		if t < 1 {
			t = 1
		}
		for hc.WS.Error == nil {
			if j, err := monitorToJson(); err == nil {
				hc.WS.Send(j)
			}
			<-time.After(time.Duration(t) * time.Second)
		}
	}
}
