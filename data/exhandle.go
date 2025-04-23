// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package data

import (
	goutil "github.com/donnie4w/gofer/util"
	"github.com/donnie4w/tim/errs"
	"github.com/donnie4w/tim/stub"
	"github.com/donnie4w/tim/sys"
	"github.com/donnie4w/tim/util"
)

type externHandle struct {
	externdb *externdb
}

func (eh *externHandle) init() service {
	eh.externdb = &externdb{}
	eh.externdb.init()
	return eh
}

func (eh *externHandle) Register(username, pwd string, domain *string) (node string, e errs.ERROR) {
	e = errs.ERR_INTERFACE
	return
}

func (eh *externHandle) Login(username, pwd string, domain *string) (_r string, e errs.ERROR) {
	if _r, _ = eh.externdb.login(username, pwd); _r != "" {
		return
	}
	e = errs.ERR_NOPASS
	return
}

func (eh *externHandle) Modify(uint64, *string, string, *string) errs.ERROR {
	return errs.ERR_PERM_DENIED
}

func (eh *externHandle) AuthNode(username, pwd string, domain *string) (node string, err errs.ERROR) {
	if node, _ = eh.externdb.authNode(username, pwd); node != "" {
		return
	}
	err = errs.ERR_NOPASS
	return
}

func (eh *externHandle) SaveMessage(tm *stub.TimMessage) (err error) {
	id := *tm.ID
	tm.ID = nil
	var chatId []byte
	if tm.MsType == sys.SOURCE_ROOM {
		chatId = util.ChatIdByRoom(tm.RoomTid.Node, tm.FromTid.Domain)
	} else {
		chatId = util.ChatIdByNode(tm.FromTid.Node, tm.ToTid.Node, tm.FromTid.Domain)
	}
	fid := tm.FromTid
	tm.FromTid = &stub.Tid{Node: fid.Node}
	stanze := util.Mask(goutil.TEncode(tm))
	var mid int64
	mid, err = eh.externdb.saveMessage(chatId, int32(goutil.FNVHash32([]byte(fid.Node))), stanze)
	tm.Mid = &mid
	if id == 0 {
		id = goutil.UUID64()
	}
	tm.ID = &id
	tm.FromTid = fid
	return
}

func (eh *externHandle) GetMessage(fromNode string, domain *string, rtype int8, to string, mid, timeseries, limit int64) (tmList []*stub.TimMessage, err error) {
	var chatId []byte
	if rtype == 1 {
		chatId = util.ChatIdByNode(fromNode, to, domain)
	} else {
		chatId = util.ChatIdByRoom(to, domain)
	}
	if rs, e := eh.externdb.getmessage(chatId, mid, timeseries, limit); e == nil {
		tmList = make([]*stub.TimMessage, 0)
		for k, v := range rs {
			if tm, err := goutil.TDecode(util.Mask(v), &stub.TimMessage{}); err == nil {
				tm.Mid = &k
				tmList = append(tmList, tm)
			}
		}
	} else {
		err = e
	}
	return
}

func (eh *externHandle) GetFidByMid(tid []byte, mid int64) (int64, error) {
	if mid <= 0 {
		return 0, errs.ERR_PARAMS.Error()
	}
	return eh.externdb.getFidById(tid, mid)
}

func (eh *externHandle) DelMessageByMid(tid []byte, mid int64) (err error) {
	return eh.externdb.delMessageById(mid)
}

func (eh *externHandle) SaveOfflineMessage(tnode string, tm *stub.TimMessage) (err error) {
	fid := tm.FromTid
	if fid != nil {
		tm.FromTid = &stub.Tid{Node: fid.Node}
	}
	mid := int64(0)
	if tm.OdType == sys.ORDER_INOF {
		mid = tm.GetMid()
	}
	err = eh.externdb.saveOfflineMessage(tnode, tm.GetID(), util.Mask(goutil.TEncode(tm)), mid)
	tm.FromTid = fid
	return
}

func (eh *externHandle) GetOfflineMessage(node string, limit int) (oblist []*OfflineBean, err error) {
	if oblist, err = eh.externdb.getOfflineMessage(node, limit); err == nil {
		for _, ob := range oblist {
			ob.Stanze = util.Mask(ob.Stanze)
		}
	}
	return
}

func (eh *externHandle) DelOfflineMessage(tid uint64, ids ...any) (_r int64, err error) {
	return eh.externdb.delOfflineMessage(ids...)
}

func (eh *externHandle) AuthGroupAndUser(groupnode, usernode string, domain *string) (ok bool, err error) {
	return eh.externdb.authGroup(groupnode, usernode)
}
