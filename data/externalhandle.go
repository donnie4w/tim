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
	. "github.com/donnie4w/tim/stub"
	"github.com/donnie4w/tim/sys"
	"github.com/donnie4w/tim/util"
)

type externalhandle struct {
	externaldb *externaldb
}

func (this *externalhandle) init() *externalhandle {
	this.externaldb = &externaldb{}
	this.externaldb.init()
	return this
}

func (this *externalhandle) Register(username, pwd string, domain *string) (node string, e errs.ERROR) {
	e = errs.ERR_INTERFACE
	return
}

func (this *externalhandle) Login(username, pwd string, domain *string) (_r string, e errs.ERROR) {
	if _r, _ = this.externaldb.login(username, pwd); _r != "" {
		return
	}
	e = errs.ERR_NOPASS
	return
}

func (this *externalhandle) Modify(uint64, *string, string, *string) errs.ERROR {
	return errs.ERR_PERM_DENIED
}

func (this *externalhandle) AuthNode(username, pwd string, domain *string) (node string, err errs.ERROR) {
	if node, _ = this.externaldb.authNode(username, pwd); node != "" {
		return
	}
	err = errs.ERR_NOPASS
	return
}

func (this *externalhandle) SaveMessage(tm *TimMessage) (err error) {
	id := *tm.ID
	tm.ID = nil
	var chatId []byte
	if tm.MsType == sys.SOURCE_ROOM {
		chatId = util.ChatIdByRoom(tm.RoomTid.Node, tm.FromTid.Domain)
	} else {
		chatId = util.ChatIdByNode(tm.FromTid.Node, tm.ToTid.Node, tm.FromTid.Domain)
	}
	fid := tm.FromTid
	tm.FromTid = &Tid{Node: fid.Node}
	stanze := util.Mask(goutil.TEncode(tm))
	var mid int64
	mid, err = this.externaldb.saveMessage(chatId, int32(goutil.FNVHash32([]byte(fid.Node))), stanze)
	tm.Mid = &mid
	if id == 0 {
		id = goutil.UUID64()
	}
	tm.ID = &id
	tm.FromTid = fid
	return
}

func (this *externalhandle) GetMessage(fromNode string, domain *string, rtype int8, to string, mid, limit int64) (tmList []*TimMessage, err error) {
	var chatId []byte
	if rtype == 1 {
		chatId = util.ChatIdByNode(fromNode, to, domain)
	} else {
		chatId = util.ChatIdByRoom(to, domain)
	}
	if rs, e := this.externaldb.getmessage(chatId, mid, limit); e == nil {
		tmList = make([]*TimMessage, 0)
		for k, v := range rs {
			if tm, err := goutil.TDecode(util.Mask(v), &TimMessage{}); err == nil {
				tm.Mid = &k
				tmList = append(tmList, tm)
			}
		}
	} else {
		err = e
	}
	return
}

func (this *externalhandle) GetChatIdByMid(tid []byte, mid int64) ([]byte, int64, error) {
	if mid <= 0 {
		return nil, 0, errs.ERR_PARAMS.Error()
	}
	return this.externaldb.getChatIdById(mid)
}

func (this *externalhandle) DelMessageByMid(tid []byte, mid int64) (err error) {
	return this.externaldb.delMessageById(mid)
}

// func (this *externalhandle) ExistOfflineMessage(tm *TimMessage) (exsit bool) {
// 	exsit, _ = this.externaldb.existOfflineMessage(tm.ToTid.Node, *tm.ID)
// 	return
// }

func (this *externalhandle) SaveOfflineMessage(tm *TimMessage) (err error) {
	node := tm.ToTid.Node
	fid := tm.FromTid
	tm.FromTid = &Tid{Node: fid.Node}
	mid := int64(0)
	if tm.OdType == sys.ORDER_INOF && tm.Mid != nil {
		mid = *tm.Mid
	}
	err = this.externaldb.saveOfflineMessage(node, *tm.ID, util.Mask(goutil.TEncode(tm)), mid)
	tm.FromTid = fid
	return
}

func (this *externalhandle) GetOfflineMessage(node string, limit int) (oblist []*OfflineBean, err error) {
	if oblist, err = this.externaldb.getOfflineMessage(node, limit); err == nil {
		for _, ob := range oblist {
			ob.Stanze = util.Mask(ob.Stanze)
		}
	}
	return
}

func (this *externalhandle) DelOfflineMessage(tid uint64, ids ...int64) (_r int64, err error) {
	return this.externaldb.delOfflineMessage(ids...)
}

func (this *externalhandle) AuthGroupAndUser(groupnode, usernode string, domain *string) (ok bool, err error) {
	return this.externaldb.authGroup(groupnode, usernode)
}
