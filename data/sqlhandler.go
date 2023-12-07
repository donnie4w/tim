// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim
//

package data

import (
	. "github.com/donnie4w/gofer/util"
	. "github.com/donnie4w/tim/stub"
	"github.com/donnie4w/tim/sys"
	"github.com/donnie4w/tim/util"
)

type sqlhandler struct{}

func (this *sqlhandler) Register(username, pwd string, domain *string) (node string, e sys.ERROR) {
	e = sys.ERR_INTERFACE
	return
}

func (this *sqlhandler) Login(username, pwd string, domain *string) (_r string, e sys.ERROR) {
	if _r, _ = sqlHandle.login(username, pwd); _r != "" {
		return
	}
	e = sys.ERR_NOPASS
	return
}

func (this *sqlhandler) Modify(uint64, *string, string, *string) sys.ERROR {
	return sys.ERR_AUTH
}

func (this *sqlhandler) Token(username, pwd string, domain *string) (_r string, err sys.ERROR) {
	if _r, _ = sqlHandle.token(username, pwd); _r != "" {
		return
	}
	err = sys.ERR_NOPASS
	return
}

func (this *sqlhandler) SaveMessage(tm *TimMessage) (err error) {
	id := *tm.ID
	tm.ID = nil
	var chatId uint64
	if tm.MsType == sys.SOURCE_ROOM {
		chatId = util.ChatIdByRoom(tm.RoomTid.Node, tm.FromTid.Domain)
	} else {
		chatId = util.ChatIdByNode(tm.FromTid.Node, tm.ToTid.Node, tm.FromTid.Domain)
	}
	fid := tm.FromTid
	tm.FromTid = &Tid{Node: fid.Node}
	stanze := util.Mask(TEncode(tm))
	var mid int64
	mid, err = sqlHandle.saveMessage(chatId, stanze)
	tm.Mid = &mid
	if id == 0 {
		id = RandId()
	}
	tm.ID = &id
	tm.FromTid = fid
	return
}

func (this *sqlhandler) GetMessage(fromTid *Tid, rtype int8, to string, mid, limit int64) (tmList []*TimMessage, err error) {
	chatId := uint64(0)
	if rtype == 1 {
		chatId = util.ChatIdByNode(fromTid.Node, to, fromTid.Domain)
	} else {
		chatId = util.ChatIdByRoom(to, fromTid.Domain)
	}
	if rs, e := sqlHandle.getmessage(chatId, mid, limit); e == nil {
		tmList = make([]*TimMessage, 0)
		for k, v := range rs {
			if tm, err := TDecode(util.Mask(v), &TimMessage{}); err == nil {
				tm.Mid = &k
				tmList = append(tmList, tm)
			}
		}
	} else {
		err = e
	}
	return
}

func (this *sqlhandler) GetMessageByMid(tid uint64, mid int64) (tm *TimMessage, err error) {
	if mid <= 0 {
		err = sys.ERR_PARAMS.Error()
		return
	}
	if bs, e := sqlHandle.getMessageById(mid); e == nil {
		tm, err = TDecode(util.Mask(bs), &TimMessage{})
	} else {
		err = e
	}
	return
}

func (this *sqlhandler) DelMessageByMid(tid uint64, mid int64) (err error) {
	return sqlHandle.delMessageById(mid)
}

// func (this *sqlhandler) ExistOfflineMessage(tm *TimMessage) (exsit bool) {
// 	exsit, _ = sqlHandle.existOfflineMessage(tm.ToTid.Node, *tm.ID)
// 	return
// }

func (this *sqlhandler) SaveOfflineMessage(tm *TimMessage) (err error) {
	node := tm.ToTid.Node
	fid := tm.FromTid
	tm.FromTid = &Tid{Node: fid.Node}
	mid := int64(0)
	if tm.OdType == sys.ORDER_INOF && tm.Mid != nil {
		mid = *tm.Mid
	}
	err = sqlHandle.saveOfflineMessage(node, *tm.ID, util.Mask(TEncode(tm)), mid)
	tm.FromTid = fid
	return
}

func (this *sqlhandler) GetOfflineMessage(tid *Tid, limit int) (oblist []*OfflineBean, err error) {
	if oblist, err = sqlHandle.getOfflineMessage(tid.Node, limit); err == nil {
		for _, ob := range oblist {
			ob.Stanze = util.Mask(ob.Stanze)
		}
	}
	return
}

func (this *sqlhandler) DelOfflineMessage(tid uint64, ids ...int64) (_r int64, err error) {
	return sqlHandle.delOfflineMessage(ids...)
}

func (this *sqlhandler) AuthGroupAndUser(groupnode, usernode string, domain *string) (ok bool, err error) {
	return sqlHandle.authGroup(groupnode, usernode)
}
