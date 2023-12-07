// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim
//

package data

import (
	"time"

	"github.com/donnie4w/gofer/buffer"
	. "github.com/donnie4w/gofer/util"
	. "github.com/donnie4w/tim/stub"
	"github.com/donnie4w/tim/sys"
	"github.com/donnie4w/tim/util"
)

type tldbhandler struct{}

func (this *tldbhandler) Register(username, pwd string, domain *string) (node string, e sys.ERROR) {
	uuid := util.CreateUUID(username, domain)
	if a, _ := SelectByIdx[timuser]("UUID", uuid); a != nil {
		e = sys.ERR_HASEXIST
		return
	}
	tu := &timuser{UUID: uuid, Pwd: this._pwd(uuid, pwd, domain), Createtime: time.Now().UnixNano()}
	if _, err := Insert(tu); err == nil {
		node = util.UUIDToNode(uuid)
	} else {
		e = sys.ERR_DATABASE
	}
	return
}

func (this *tldbhandler) Login(username, pwd string, domain *string) (_r string, e sys.ERROR) {
	uuid := util.CreateUUID(username, domain)
	if a, err := SelectByIdx[timuser]("UUID", uuid); err == nil && a != nil {
		if a.Pwd == this._pwd(uuid, pwd, domain) {
			_r = util.UUIDToNode(uuid)
			return
		}
	}
	e = sys.ERR_NOPASS
	return
}

func (this *tldbhandler) Modify(uuid uint64, pwd *string, pwdLast string, domain *string) (e sys.ERROR) {
	if a, err := SelectByIdx[timuser]("UUID", uuid); err == nil && a != nil {
		if pwd != nil {
			if this._pwd(uuid, *pwd, domain) != a.Pwd {
				return sys.ERR_AUTH
			}
		} else if pwdLast == "" {
			return sys.ERR_PARAMS
		}
		UpdateNonzero(&timuser{UUID: uuid, Id: a.Id, Pwd: this._pwd(uuid, pwdLast, domain)})
	} else {
		e = sys.ERR_NOEXIST
	}
	return
}

func (this *tldbhandler) Token(username, pwd string, domain *string) (_r string, err sys.ERROR) {
	uuid := util.CreateUUID(username, domain)
	if a, _ := SelectByIdx[timuser]("UUID", uuid); a != nil {
		if a.Pwd == this._pwd(uuid, pwd, domain) {
			_r = util.UUIDToNode(uuid)
			return
		}
	}
	err = sys.ERR_NOPASS
	return
}

func (this *tldbhandler) _pwd(uuid uint64, pwd string, domain *string) uint64 {
	buf := buffer.NewBufferByPool()
	defer buf.Free()
	buf.Write(Int64ToBytes(int64(uuid)))
	buf.WriteString(sys.Conf.Salt)
	buf.Write([]byte(pwd))
	return CRC64(MD5(buf.Bytes()))
}

func (this *tldbhandler) SaveMessage(tm *TimMessage) (err error) {
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
	mid, err = Insert(&timmessage{ChatId: chatId, Stanza: stanze})
	tm.Mid = &mid
	if id == 0 {
		id = RandId()
	}
	tm.ID = &id
	tm.FromTid = fid
	return
}

func (this *tldbhandler) GetMessage(fromTid *Tid, rtype int8, to string, mid, limit int64) (tmList []*TimMessage, err error) {
	chatId := uint64(0)
	if rtype == 1 {
		chatId = util.ChatIdByNode(fromTid.Node, to, fromTid.Domain)
	} else {
		chatId = util.ChatIdByRoom(to, fromTid.Domain)
	}
	if mid == 0 {
		mxid, _ := SelectIdByIdx[timmessage]("ChatId", chatId)
		mid, _ = SelectIdByIdxSeq[timmessage]("ChatId", chatId, mxid)
	}
	if as, err := SelectByIdxDescLimit[timmessage](mid, limit, "ChatId", chatId); err == nil {
		tmList = make([]*TimMessage, 0)
		for _, a := range as {
			if tm, err := TDecode(util.Mask(a.Stanza), &TimMessage{}); err == nil {
				tm.Mid = &a.Id
				tmList = append(tmList, tm)
			}
		}
	}
	return
}

func (this *tldbhandler) GetMessageByMid(tid uint64, mid int64) (tm *TimMessage, err error) {
	if mid <= 0 {
		err = sys.ERR_PARAMS.Error()
		return
	}
	if a, e := SelectById[timmessage](tid, mid); e == nil && a != nil {
		tm, err = TDecode(util.Mask(a.Stanza), &TimMessage{})
	} else {
		err = e
	}
	return
}

func (this *tldbhandler) DelMessageByMid(tid uint64, mid int64) (err error) {
	return Delete[timmessage](tid, mid)
}

// func (this *tldbhandler) ExistOfflineMessage(tm *TimMessage) (exsit bool) {
// 	uuid := util.NodeToUUID(tm.ToTid.Node)
// 	if a, _ := SelectByIdx[timoffline]("Unik", *tm.ID); a != nil {
// 		exsit = true
// 	}
// 	return
// }

func (this *tldbhandler) SaveOfflineMessage(tm *TimMessage) (err error) {
	node := tm.ToTid.Node
	fid := tm.FromTid
	if tm.Timestamp != nil {
		t := time.Now().UnixNano()
		tm.Timestamp = &t
	}
	tm.FromTid = &Tid{Node: fid.Node}
	if tm.OdType == sys.ORDER_INOF && tm.Mid != nil && *tm.Mid > 0 {
		cid := util.ChatIdByNode(fid.Node, tm.ToTid.Node, fid.Domain)
		Insert(&timoffline{UUID: util.NodeToUUID(node), Mid: *tm.Mid, ChatId: cid})
	} else {
		Insert(&timoffline{UUID: util.NodeToUUID(node), Stanza: util.Mask(TEncode(tm))})
	}
	tm.FromTid = fid
	return
}

func (this *tldbhandler) GetOfflineMessage(tid *Tid, limit int) (oblist []*OfflineBean, err error) {
	uuid := util.NodeToUUID(tid.Node)
	if tfs, _ := SelectByIdxLimit[timoffline](0, int64(limit), "UUID", uuid); tfs != nil {
		oblist = make([]*OfflineBean, 0)
		for _, tf := range tfs {
			ob := &OfflineBean{Id: tf.Id, Mid: tf.Mid}
			if tf.Mid == 0 {
				if tf.Stanza != nil {
					ob.Stanze = util.Mask(tf.Stanza)
				}
			} else {
				if a, _ := SelectById[timmessage](tf.ChatId, tf.Mid); a != nil {
					ob.Stanze, ob.Mid = util.Mask(a.Stanza), tf.Mid
				}
			}
			oblist = append(oblist, ob)
		}
	}
	return
}

func (this *tldbhandler) DelOfflineMessage(tid uint64, ids ...int64) (_r int64, err error) {
	for _, id := range ids {
		if id > 0 {
			_r, err = 1, Delete[timoffline](tid, id)
		}
	}
	return
}

func (this *tldbhandler) AuthGroupAndUser(groupnode, usernode string, domain *string) (ok bool, err error) {
	if util.CheckNode(groupnode) && util.CheckNode(usernode) {
		if a, _ := SelectByIdx[timgroup]("UUID", util.NodeToUUID(groupnode)); a != nil && a.Status == 1 {
			relateid := util.RelateIdForGroup(groupnode, usernode, domain)
			if a, _ := SelectByIdx[timrelate]("UUID", relateid); a != nil {
				ok = a.Status == 0x11
			}
		}
	}
	return
}

func (this *tldbhandler) AuthUserAndUser(fnode, tnode string, domain *string) (_r bool) {
	if util.CheckNode(fnode) && util.CheckNode(tnode) {
		cid := util.ChatIdByNode(fnode, tnode, domain)
		if a, _ := SelectByIdx[timrelate]("UUID", cid); a != nil {
			if a.Status == 0x10|0x1 {
				return true
			}
		}
	}
	return
}

func (this *tldbhandler) ExistUser(node string) (_r bool) {
	if a, _ := SelectByIdx[timuser]("UUID", util.NodeToUUID(node)); a != nil {
		_r = true
	}
	return
}

func (this *tldbhandler) ExistGroup(node string) (_r bool) {
	if a, _ := SelectByIdx[timgroup]("UUID", util.NodeToUUID(node)); a != nil {
		_r = true
	}
	return
}
