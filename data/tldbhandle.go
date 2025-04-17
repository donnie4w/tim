// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package data

import (
	"bytes"
	"github.com/donnie4w/gofer/uuid"
	"github.com/donnie4w/tim/errs"
	"time"

	"github.com/donnie4w/gofer/buffer"
	. "github.com/donnie4w/gofer/util"
	goutil "github.com/donnie4w/gofer/util"
	. "github.com/donnie4w/tim/stub"
	"github.com/donnie4w/tim/sys"
	"github.com/donnie4w/tim/util"
)

type tldbHandle struct{}

func (th *tldbHandle) init() service {
	tldbInit()
	return th
}

func (th *tldbHandle) Register(username, pwd string, domain *string) (node string, e errs.ERROR) {
	uuid := util.CreateUUID(username, domain)
	if a, _ := SelectByIdxWithInt[timuser]("UUID", uuid); a != nil {
		e = errs.ERR_HASEXIST
		return
	}
	if hs, err := util.Password(uuid, pwd, domain); err != nil {
		tu := &timuser{UUID: uuid, Pwd: hs, Createtime: time.Now().UnixNano(), Timeseries: TimeNano()}
		if _, err := Insert(tu); err == nil {
			node = util.UUIDToNode(uuid)
		}
		return
	}
	return "", errs.ERR_DATABASE
}

func (th *tldbHandle) Login(username, pwd string, domain *string) (_r string, e errs.ERROR) {
	uuid := util.CreateUUID(username, domain)
	if a, err := SelectByIdxWithInt[timuser]("UUID", uuid); err == nil && a != nil {
		if util.CheckPasswordHash(uuid, pwd, domain, a.Pwd) {
			_r = util.UUIDToNode(uuid)
			return
		}
	}
	e = errs.ERR_NOPASS
	return
}

func (th *tldbHandle) Modify(uuid uint64, pwd *string, pwdLast string, domain *string) (e errs.ERROR) {
	if uuid == 0 {
		return errs.ERR_ACCOUNT
	}
	if a, err := SelectByIdxWithInt[timuser]("UUID", uuid); err == nil && a != nil {
		if pwd != nil {
			if util.CheckPasswordHash(uuid, *pwd, domain, a.Pwd) {
				return errs.ERR_PERM_DENIED
			}
		} else if pwdLast == "" {
			return errs.ERR_PARAMS
		}
		if hs, err := util.Password(uuid, pwdLast, domain); err != nil {
			UpdateNonzero(&timuser{UUID: uuid, Id: a.Id, Pwd: hs})
			return
		}
	}
	return errs.ERR_NOEXIST
}

func (th *tldbHandle) AuthNode(username, pwd string, domain *string) (_r string, err errs.ERROR) {
	uuid := util.CreateUUID(username, domain)
	if a, _ := SelectByIdxWithInt[timuser]("UUID", uuid); a != nil {
		if util.CheckPasswordHash(uuid, pwd, domain, a.Pwd) {
			_r = util.UUIDToNode(uuid)
			return
		}
	}
	err = errs.ERR_NOPASS
	return
}

func (th *tldbHandle) _pwd(uuid uint64, pwd string, domain *string) uint64 {
	buf := buffer.NewBufferByPool()
	defer buf.Free()
	buf.Write(Int64ToBytes(int64(uuid)))
	buf.WriteString(sys.Conf.Salt)
	buf.Write([]byte(pwd))
	return CRC64(MD5(buf.Bytes()))
}

func (th *tldbHandle) SaveMessage(tm *TimMessage) (err error) {
	id := tm.GetID()
	tm.ID = nil
	var chatId []byte
	if tm.MsType == sys.SOURCE_ROOM {
		chatId = util.ChatIdByRoom(tm.RoomTid.Node, tm.FromTid.Domain)
	} else {
		chatId = util.ChatIdByNode(tm.FromTid.Node, tm.ToTid.Node, tm.FromTid.Domain)
	}
	fid := tm.FromTid
	tm.FromTid = &Tid{Node: fid.Node}
	stanze := util.Mask(TEncode(tm))
	var mid int64
	mid, err = Insert(&timmessage{ChatId: chatId, Fid: int64(goutil.FNVHash32([]byte(fid.Node))), Stanza: stanze, Timeseries: TimeNano()})
	tm.Mid = &mid
	if id == 0 {
		id = uuid.NewUUID().Int64()
	}
	tm.ID = &id
	tm.FromTid = fid
	return
}

func (th *tldbHandle) GetMessage(fromNode string, domain *string, rtype int8, to string, mid, timestamp, limit int64) (tmList []*TimMessage, err error) {
	var chatId []byte
	if rtype == 1 {
		chatId = util.ChatIdByNode(fromNode, to, domain)
	} else {
		chatId = util.ChatIdByRoom(to, domain)
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

func (th *tldbHandle) GetFidByMid(tid []byte, mid int64) (int64, error) {
	if mid <= 0 {
		return 0, errs.ERR_PARAMS.Error()
	}
	if a, err := SelectById[timmessage](tid, mid); err == nil && a != nil {
		if !bytes.Equal(a.ChatId, tid) {
			return 0, nil
		} else {
			return a.Fid, nil
		}
	} else {
		return 0, err
	}
}

func (th *tldbHandle) DelMessageByMid(tid []byte, mid int64) (err error) {
	return Delete[timmessage](goutil.FNVHash64(tid), mid)
}

// func (th *tldbHandle) ExistOfflineMessage(tm *TimMessage) (exsit bool) {
// 	uuid := util.NodeToUUID(tm.ToTid.Node)
// 	if a, _ := SelectByIdx[timoffline]("Unik", *tm.ID); a != nil {
// 		exsit = true
// 	}
// 	return
// }

func (th *tldbHandle) SaveOfflineMessage(tnode string, tm *TimMessage) (err error) {
	fid := tm.FromTid
	if tm.Timestamp != nil {
		t := TimeNano()
		tm.Timestamp = &t
	}
	if fid != nil {
		tm.FromTid = &Tid{Node: fid.Node}
	}

	if tm.OdType == sys.ORDER_INOF && tm.Mid != nil && *tm.Mid > 0 {
		cid := util.ChatIdByNode(fid.Node, tm.ToTid.Node, fid.Domain)
		Insert(&timoffline{UUID: util.NodeToUUID(tnode), Mid: tm.GetMid(), ChatId: cid, Timeseries: TimeNano()})
	} else {
		Insert(&timoffline{UUID: util.NodeToUUID(tnode), Stanza: util.Mask(TEncode(tm)), Timeseries: TimeNano()})
	}
	tm.FromTid = fid
	return
}

func (th *tldbHandle) GetOfflineMessage(node string, limit int) (oblist []*OfflineBean, err error) {
	uuid := util.NodeToUUID(node)
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

func (th *tldbHandle) DelOfflineMessage(tid uint64, ids ...any) (_r int64, err error) {
	for _, id := range ids {
		if v := id.(int64); v > 0 {
			_r, err = 1, Delete[timoffline](tid, v)
		}
	}
	return
}

func (th *tldbHandle) AuthGroupAndUser(groupnode, usernode string, domain *string) (ok bool, err error) {
	if util.CheckNode(groupnode) && util.CheckNode(usernode) {
		if a, _ := SelectByIdxWithInt[timgroup]("UUID", util.NodeToUUID(groupnode)); a != nil && a.Status == 1 {
			relateid := util.RelateIdForGroup(groupnode, usernode, domain)
			if a, _ := SelectByIdx[timrelate]("UUID", relateid); a != nil {
				ok = a.Status == 0x11
			}
		}
	}
	return
}

func (th *tldbHandle) AuthUserAndUser(fnode, tnode string, domain *string) (_r bool) {
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

func (th *tldbHandle) ExistUser(node string) (_r bool) {
	if a, _ := SelectByIdxWithInt[timuser]("UUID", util.NodeToUUID(node)); a != nil {
		_r = true
	}
	return
}

func (th *tldbHandle) ExistGroup(node string) (_r bool) {
	if a, _ := SelectByIdxWithInt[timgroup]("UUID", util.NodeToUUID(node)); a != nil {
		_r = true
	}
	return
}
