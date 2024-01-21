// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim
//

package tc

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	. "github.com/donnie4w/gofer/hashmap"
	tldbKs "github.com/donnie4w/gofer/keystore"
	. "github.com/donnie4w/gofer/util"
	. "github.com/donnie4w/tim/keystore"
	. "github.com/donnie4w/tim/stub"
	"github.com/donnie4w/tim/sys"
	"github.com/donnie4w/tim/util"
	"github.com/donnie4w/tlnet"
)

func init() {
	sys.Service.Put(4, adminservice)
}

type adminService struct {
	isClose bool
	tlAdmin *tlnet.Tlnet
}

var adminservice = &adminService{false, tlnet.NewTlnet()}

func (this *adminService) Serve() (err error) {
	if strings.TrimSpace(sys.DEBUGADDR) != "" {
		go tlDebug()
		<-time.After(500 * time.Millisecond)
	}
	if strings.TrimSpace(sys.WEBADMINADDR) != "" {
		err = this._serve(strings.TrimSpace(sys.WEBADMINADDR), sys.Conf.AdminTls, sys.Conf.Ssl_crt, sys.Conf.Ssl_crt_key)
	}
	return
}

func (this *adminService) Close() (err error) {
	defer util.Recover()
	if strings.TrimSpace(sys.WEBADMINADDR) != "" {
		this.isClose = true
		err = this.tlAdmin.Close()
	}
	return
}

func (this *adminService) _serve(addr string, TLS bool, serverCrt, serverKey string) (err error) {
	defer util.Recover()
	if addr, err = util.ParseAddr(addr); err != nil {
		return
	}
	sys.WEBADMINADDR = addr
	tlnet.SetLogOFF()
	this.tlAdmin.Handle("/login", loginHandler)
	this.tlAdmin.Handle("/init", initHandler)
	this.tlAdmin.Handle("/lang", langHandler)
	this.tlAdmin.Handle("/", initHandler)
	this.tlAdmin.Handle("/bootstrap.css", cssHandler)
	this.tlAdmin.Handle("/bootstrap.min.js", jsHandler)
	this.tlAdmin.HandleWithFilter("/sysvar", loginFilter(), sysVarHtml)
	this.tlAdmin.HandleWithFilter("/timResetAuth", authFilter(), timResetAuthHandler)
	this.tlAdmin.HandleWithFilter("/timToken", authFilter(), timTokenHandler)
	this.tlAdmin.HandleWithFilter("/timMessage", authFilter(), timMessageHandler)
	this.tlAdmin.HandleWithFilter("/timRegister", authFilter(), timRegisterHandler)
	this.tlAdmin.HandleWithFilter("/timModifyUserInfo", authFilter(), timModifyUserInfo)
	this.tlAdmin.HandleWithFilter("/timBlockUser", authFilter(), timBlockUserHandler)
	this.tlAdmin.HandleWithFilter("/timBlockList", authFilter(), timBlockListHandler)
	this.tlAdmin.HandleWithFilter("/timOnline", authFilter(), timOnlineHandler)
	this.tlAdmin.HandleWithFilter("/timVroom", authFilter(), timVroomHandler)
	this.tlAdmin.HandleWithFilter("/timNewRoom", authFilter(), timNewRoomHandler)
	this.tlAdmin.HandleWithFilter("/timModifyRoomInfo", authFilter(), timModifyRoomInfoHandler)
	this.tlAdmin.HandleWithFilter("/monitor", loginFilter(), monitorHtml)
	this.tlAdmin.HandleWebSocketBindConfig("/monitorData", mntHandler, mntConfig())
	this.tlAdmin.HandleWithFilter("/data", loginFilter(), dataMonitorHtml)
	this.tlAdmin.HandleWebSocketBindConfig("/ddmonitorData", ddmntHandler, ddmntConfig())

	if TLS {
		if IsFileExist(serverCrt) && IsFileExist(serverKey) {
			sys.FmtLog("webAdmin start tls [", addr, "]")
			err = this.tlAdmin.HttpStartTLS(addr, serverCrt, serverKey)
		} else {
			sys.FmtLog("webAdmin start tls by bytes [", addr, "]")
			err = this.tlAdmin.HttpStartTlsBytes(addr, []byte(tldbKs.ServerCrt), []byte(tldbKs.ServerKey))
		}
	}
	if !this.isClose {
		sys.FmtLog("webAdmin start [", addr, "]")
		err = this.tlAdmin.HttpStart(addr)
	}
	if !this.isClose && err != nil {
		sys.FmtLog("webAdmin start failed:", err.Error())
	}
	return
}

var sessionMap = NewMapL[string, *tldbKs.UserBean]()

func loginFilter() (f *tlnet.Filter) {
	defer util.Recover()
	f = tlnet.NewFilter()
	f.AddIntercept(".*?", func(hc *tlnet.HttpContext) bool {
		if len(Admin.AdminList()) > 0 {
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
		if _r, ok := Admin.GetAdmin(name); ok {
			if strings.ToLower(_r.Pwd) == strings.ToLower(Md5Str(pwd)) {
				return false
			}
		}
		hc.ResponseBytes(http.StatusMethodNotAllowed, nil)
		return true
	})
	return
}

func getSessionid() string {
	return fmt.Sprint("t", CRC32(Int64ToBytes(sys.UUID)))
}

func getLangId() string {
	return fmt.Sprint("l", CRC32(Int64ToBytes(sys.UUID)))
}

func isLogin(hc *tlnet.HttpContext) (isLogin bool) {
	if len(Admin.AdminList()) > 0 {
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

func timTokenHandler(hc *tlnet.HttpContext) {
	type tk struct {
		Name     string `json:"name"`
		Password string `json:"password"`
		Domain   string `json:"domain"`
	}
	defer util.Recover()
	var nodeOrName string
	var password *string
	var domain *string

	if reqform(hc) {
		nodeOrName = hc.PostParamTrimSpace("name")
		_domain := hc.PostParamTrimSpace("domain")
		_password := hc.PostParamTrimSpace("password")

		if _domain != "" {
			domain = &_domain
		}
		if _password != "" {
			password = &_password
		}
	} else {
		bs := hc.RequestBody()
		if t, err := JsonDecode[tk](bs); err == nil {
			nodeOrName = t.Name
			if t.Domain != "" {
				domain = &t.Domain
			}
			if t.Password != "" {
				password = &t.Password
			}
		}
	}

	var ta *TimAck
	if nodeOrName != "" {
		if t, n, err := sys.OsToken(nodeOrName, password, domain); err == nil {
			ta = &TimAck{Ok: true, TimType: int8(sys.TIMTOKEN), N: &n, T: &t}
		} else {
			ta = &TimAck{Ok: false, TimType: int8(sys.TIMTOKEN), Error: err.TimError()}
		}
	} else {
		ta = &TimAck{Ok: false, TimType: int8(sys.TIMTOKEN), Error: sys.ERR_PARAMS.TimError()}
	}
	hc.ResponseBytes(http.StatusOK, JsonEncode(ta))
}

func timMessageHandler(hc *tlnet.HttpContext) {
	defer util.Recover()
	type tk struct {
		Nodes   *TimNodes   `json:"nodes"`
		Message *TimMessage `json:"message"`
	}
	var nodes *TimNodes
	var message *TimMessage

	if reqform(hc) {
		nodes, _ = JsonDecode[*TimNodes]([]byte(hc.PostParam("nodes")))
		message, _ = JsonDecode[*TimMessage]([]byte(hc.PostParam("message")))
	} else {
		bs := hc.RequestBody()
		if t, err := JsonDecode[tk](bs); err == nil {
			nodes = t.Nodes
			message = t.Message
		}
	}

	if err := sys.OsMessage(nodes, message); err == nil {
		tk := &TimAck{Ok: true}
		hc.ResponseString(string(JsonEncode(tk)))
	} else {
		tk := &TimAck{Ok: false, Error: err.TimError()}
		hc.ResponseString(string(JsonEncode(tk)))
	}
}

func timRegisterHandler(hc *tlnet.HttpContext) {
	defer util.Recover()
	var username string
	var password string
	var domain *string
	type tk struct {
		Username string  `json:"username"`
		Password string  `json:"password"`
		Domain   *string `json:"domain"`
	}

	if reqform(hc) {
		username = hc.PostParamTrimSpace("username")
		password = hc.PostParam("password")
		if d := hc.PostParam("domain"); d != "" {
			domain = &d
		}
	} else {
		bs := hc.RequestBody()
		if t, err := JsonDecode[tk](bs); err == nil {
			username = t.Username
			password = t.Password
			domain = t.Domain
		}
	}

	if node, err := sys.OsRegister(username, password, domain); err == nil {
		tk := &TimAck{Ok: true, N: &node}
		hc.ResponseString(string(JsonEncode(tk)))
	} else {
		tk := &TimAck{Ok: false, Error: err.TimError()}
		hc.ResponseString(string(JsonEncode(tk)))
	}
}

func timBlockUserHandler(hc *tlnet.HttpContext) {
	defer util.Recover()
	var account string
	var _time int64
	type tk struct {
		Account string `json:"account"`
		Time    int64  `json:"time"`
	}

	if reqform(hc) {
		account = hc.PostParam("account")
		t := hc.PostParam("time")
		if i, e := strconv.Atoi(t); e == nil {
			_time = int64(i)
		} else {
			hc.ResponseString(string(JsonEncode(&TimAck{Ok: false, N: &account, Error: sys.ERR_PARAMS.TimError()})))
		}
	} else {
		bs := hc.RequestBody()
		if t, err := JsonDecode[tk](bs); err == nil {
			account = t.Account
			_time = t.Time
		}
	}

	if sys.HasNode(account) {
		sys.SendNode(account, &TimAck{Ok: true, TimType: int8(sys.TIMLOGOUT)}, sys.TIMACK)
	}
	sys.BlockUser(account, int64(_time))
	hc.ResponseString(string(JsonEncode(&TimAck{Ok: true, N: &account})))
}

func timResetAuthHandler(hc *tlnet.HttpContext) {
	defer util.Recover()
	var loginname string
	var domain *string
	var pwd string
	type tk struct {
		Loginname string  `json:"loginname"`
		Domain    *string `json:"domain"`
		Pwd       string  `json:"pwd"`
	}
	if reqform(hc) {
		loginname = hc.PostParam("loginname")
		if d := hc.PostParam("domain"); d != "" {
			domain = &d
		}
		pwd = hc.PostParam("pwd")
	} else {
		bs := hc.RequestBody()
		if t, err := JsonDecode[tk](bs); err == nil {
			loginname = t.Loginname
			domain = t.Domain
			pwd = t.Pwd
		}
	}
	if loginname == "" || pwd == "" {
		hc.ResponseString(string(JsonEncode(&TimAck{Ok: false, N: &loginname, Error: sys.ERR_PARAMS.TimError()})))
		return
	}
	if err := sys.OsModify(loginname, pwd, domain); err == nil {
		hc.ResponseString(string(JsonEncode(&TimAck{Ok: true, N: &loginname})))
	} else {
		hc.ResponseString(string(JsonEncode(&TimAck{Ok: false, Error: err.TimError()})))
	}
}

func timNewRoomHandler(hc *tlnet.HttpContext) {
	defer util.Recover()
	var node string
	var domain *string
	var topic string
	var gtype int8
	type tk struct {
		Node   string  `json:"node"`
		Topic  string  `json:"topic"`
		Domain *string `json:"domain"`
		Gtype  int8    `json:"gtype"`
	}
	if reqform(hc) {
		node, topic = hc.PostParam("node"), hc.PostParam("topic")
		if d := hc.PostParam("domain"); d != "" {
			domain = &d
		}
		if hc.PostParam("gtype") == "1" {
			gtype = 1
		} else {
			gtype = 2
		}
	} else {
		bs := hc.RequestBody()
		if t, err := JsonDecode[tk](bs); err == nil {
			node = t.Node
			domain = t.Domain
			gtype = t.Gtype
			topic = t.Topic
		}
	}
	if gnode, err := sys.OsRoom(node, topic, domain, gtype); err == nil {
		hc.ResponseString(string(JsonEncode(&TimAck{Ok: true, N: &gnode})))
	} else {
		hc.ResponseString(string(JsonEncode(&TimAck{Ok: false, Error: err.TimError()})))
	}
}

func timModifyRoomInfoHandler(hc *tlnet.HttpContext) {
	defer util.Recover()
	var unode string
	var gnode string
	var trb *TimRoomBean
	type tk struct {
		Unode    string       `json:"unode"`
		Gnode    string       `json:"gnode"`
		RoomBean *TimRoomBean `json:"roombean"`
	}
	if reqform(hc) {
		unode, gnode = hc.PostParam("unode"), hc.PostParam("gnode")
		trb, _ = JsonDecode[*TimRoomBean]([]byte(hc.PostParam("roombean")))
	} else {
		bs := hc.RequestBody()
		if t, err := JsonDecode[tk](bs); err == nil {
			unode = t.Unode
			gnode = t.Gnode
			trb = t.RoomBean
		}
	}
	if unode == "" || gnode == "" {
		hc.ResponseString(string(JsonEncode(&TimAck{Ok: false, Error: sys.ERR_PARAMS.TimError()})))
		return
	}
	if err := sys.OsRoomBean(unode, gnode, trb); err == nil {
		hc.ResponseString(string(JsonEncode(&TimAck{Ok: true})))
	} else {
		hc.ResponseString(string(JsonEncode(&TimAck{Ok: false, Error: err.TimError()})))
	}
}

func timBlockListHandler(hc *tlnet.HttpContext) {
	defer util.Recover()
	hc.ResponseString(string(JsonEncode(sys.BlockList())))
}

func timOnlineHandler(hc *tlnet.HttpContext) {
	defer util.Recover()
	hc.ResponseString(string(JsonEncode(sys.WssList())))
}

func timVroomHandler(hc *tlnet.HttpContext) {
	defer util.Recover()
	type tk struct {
		Node  string `json:"node"`
		Rtype int8   `json:"rtype"`
	}
	var node string
	var rtype int8
	if reqform(hc) {
		node = hc.PostParamTrimSpace("node")
		if r, err := strconv.Atoi(hc.PostParamTrimSpace("rtype")); err == nil {
			rtype = int8(r)
		}
	} else {
		bs := hc.RequestBody()
		if t, err := JsonDecode[tk](bs); err == nil {
			node = t.Node
			rtype = t.Rtype
		}
	}
	if node != "" && rtype > 0 {
		if _r := sys.OsVroomprocess(node, rtype); _r != "" {
			hc.ResponseString(string(JsonEncode(&TimAck{Ok: true, N: &_r})))
			return
		}
	}
	hc.ResponseString(string(JsonEncode(&TimAck{Ok: false})))
}

func timModifyUserInfo(hc *tlnet.HttpContext) {
	defer util.Recover()
	type tk struct {
		Node     string       `json:"node"`
		UserBean *TimUserBean `json:"userbean"`
	}
	var node string
	var userBean *TimUserBean
	if reqform(hc) {
		node = hc.PostParamTrimSpace("node")
		userBean, _ = JsonDecode[*TimUserBean]([]byte(hc.PostParam("userbean")))
	} else {
		bs := hc.RequestBody()
		if t, err := JsonDecode[tk](bs); err == nil {
			node = t.Node
			userBean = t.UserBean
		}
	}
	if err := sys.OsUserBean(node, userBean); err == nil {
		tk := &TimAck{Ok: true}
		hc.ResponseString(string(JsonEncode(tk)))
	} else {
		tk := &TimAck{Ok: false, Error: err.TimError()}
		hc.ResponseString(string(JsonEncode(tk)))
	}
}

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
	if len(Admin.AdminList()) > 0 && !isLogin(hc) {
		hc.Redirect("/login")
		return
	}
	if _type := hc.GetParam("type"); _type != "" {
		isadmin := isAdmin(hc)
		if _type == "1" {
			if name, pwd, _type := hc.PostParamTrimSpace("adminName"), hc.PostParamTrimSpace("adminPwd"), hc.PostParamTrimSpace("adminType"); name != "" && pwd != "" {
				if n := len(Admin.AdminList()); (n > 0 && isadmin) || n == 0 {
					alterType := false
					if t, err := strconv.Atoi(_type); err == nil {
						if _r, err := hc.GetCookie(getSessionid()); err == nil && sessionMap.Has(_r) {
							if u, ok := sessionMap.Get(_r); ok && u.Name == name && t != int(u.Type) {
								alterType = true
							}
						}
						if !alterType {
							Admin.PutAdmin(name, pwd, int8(t))
						}
					}
				} else {
					goto DENIED
				}
			}
		} else if _type == "2" && isLogin(hc) {
			if isadmin {
				if name := hc.PostParamTrimSpace("adminName"); name != "" {
					if u, ok := Admin.GetAdmin(name); ok && u.Type == 1 {
						i, j := 0, 0
						for _, s := range Admin.AdminList() {
							if _u, _ := Admin.GetAdmin(s); _u.Type == 1 {
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
					Admin.DelAdmin(name)
					sessionMap.Range(func(k string, v *tldbKs.UserBean) bool {
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
		if _r, ok := Admin.GetAdmin(name); ok {
			if strings.ToLower(_r.Pwd) == strings.ToLower(Md5Str(pwd)) {
				sid := Md5Str(fmt.Sprint(time.Now().UnixNano()))
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
	if len(Admin.AdminList()) == 0 {
		show, init, sc = "no user is created for admin, create a management user first", true, true
	}
	av := &AdminView{Show: show, Init: init, ShowCreate: sc}
	if isLogin(hc) {
		m := make(map[string]string, 0)
		for _, s := range Admin.AdminList() {
			if u, ok := Admin.GetAdmin(s); ok {
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

// func initAccount() {
// 	if len(Admin.AdminList()) == 0 {
// 		Admin.PutAdmin("admin", "123", 1)
// 	}
// }

/********************************************************/
func sysVarHtml(hc *tlnet.HttpContext) {
	defer func() {
		if err := recover(); err != nil {
			hc.ResponseString(resultHtml("server error:", err))
		}
	}()
	rn := sys.GetRemoteNode()
	sort.Slice(rn, func(i, j int) bool { return rn[i].UUID > rn[j].UUID })
	svv := &SysVarView{Show: "", RN: rn}
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
	s.CSNUM = int32(len(sys.GetALLUUIDS()))
	s.ADDR = fmt.Sprint(sys.CSADDR)
	s.ADMINADDR = sys.WEBADMINADDR
	aus := make([]int64, 0)
	for _, v := range sys.GetALLUUIDS() {
		aus = append(aus, v)
	}
	sort.Slice(aus, func(i, j int) bool { return aus[i] > aus[j] })
	s.ALLUUIDS = fmt.Sprint(aus)
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
