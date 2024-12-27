// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim
package tc

import (
	"net/http"
	"strconv"
	"strings"

	. "github.com/donnie4w/gofer/util"
	. "github.com/donnie4w/tim/stub"
	"github.com/donnie4w/tim/sys"
	"github.com/donnie4w/tim/util"
	"github.com/donnie4w/tlnet"
)

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

func timOsMessageHandler(hc *tlnet.HttpContext) {
	defer util.Recover()
	type tk struct {
		Nodes   []string    `json:"nodes"`
		Message *TimMessage `json:"message"`
	}
	var nodes []string
	var message *TimMessage

	if reqform(hc) {
		nodes, _ = JsonDecode[[]string]([]byte(hc.PostParam("nodes")))
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

func timMessageHandler(hc *tlnet.HttpContext) {
	defer util.Recover()
	bs := hc.RequestBody()
	if id, err := strconv.ParseInt(hc.ReqInfo.Header.Get("id"), 10, 64); err == nil {
		if ws, b := sys.WsById(id); b {
			if err := sys.MessageHandle(bs, ws); err == nil {
				tk := &TimAck{Ok: true}
				hc.ResponseString(string(JsonEncode(tk)))
				return
			}
		}
	}
	tk := &TimAck{Ok: false, Error: sys.ERR_PARAMS.TimError()}
	hc.ResponseString(string(JsonEncode(tk)))
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

// list of block
func timBlockListHandler(hc *tlnet.HttpContext) {
	defer util.Recover()
	hc.ResponseString(string(JsonEncode(sys.BlockList())))
}

// list of online users
func timOnlineHandler(hc *tlnet.HttpContext) {
	defer util.Recover()
	tidlist, _ := sys.WssList(0, 0)
	hc.ResponseString(string(JsonEncode(tidlist)))
}

// vroom operation
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

// modify user info
func timModifyUserInfoHnadler(hc *tlnet.HttpContext) {
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

// user online status detection
func timDetectHandler(hc *tlnet.HttpContext) {
	defer util.Recover()
	type tk struct {
		Nodes []string `json:"nodes"`
	}
	var nodes []string
	if reqform(hc) {
		if n := hc.PostParamTrimSpace("nodes"); n != "" {
			nodes = strings.Split(n, ",")
		}
	} else {
		bs := hc.RequestBody()
		if t, err := JsonDecode[tk](bs); err == nil {
			nodes = t.Nodes
		}
	}
	if len(nodes) > 0 {
		sys.Detect(nodes)
	}
}

// add roster
func timAddrosterHandler(hc *tlnet.HttpContext) {}

// remove roster
func timRmrosterHandler(hc *tlnet.HttpContext) {

}

// block roster
func timBlockrosterHandler(hc *tlnet.HttpContext) {

}

// list of  block roster
func timBlockRosterListHandler(hc *tlnet.HttpContext) {

}

// add to room
func timAddRoomHandler(hc *tlnet.HttpContext) {

}

func timPullInRoomHandler(hc *tlnet.HttpContext) {

}

func timRejectRoomHandler(hc *tlnet.HttpContext) {

}

func timKickRoomHandler(hc *tlnet.HttpContext) {

}

func timLeaveRoomHandler(hc *tlnet.HttpContext) {

}

func timCancelRoomHandler(hc *tlnet.HttpContext) {

}

func timBlockRoomHandler(hc *tlnet.HttpContext) {

}

func timBlockRoomMemberHandler(hc *tlnet.HttpContext) {

}

func timBlockRoomListHandler(hc *tlnet.HttpContext) {

}

func timBlockRoomMemberlistHandler(hc *tlnet.HttpContext) {

}
