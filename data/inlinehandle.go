// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package data

import (
	"github.com/donnie4w/gofer/buffer"
	goutil "github.com/donnie4w/gofer/util"
	"github.com/donnie4w/tim/dao"
	"github.com/donnie4w/tim/errs"
	"github.com/donnie4w/tim/log"
	. "github.com/donnie4w/tim/stub"
	"github.com/donnie4w/tim/sys"
	"github.com/donnie4w/tim/util"
	"os"
	"strings"
)

type inlineHandle struct{}

func (h *inlineHandle) init() *inlineHandle {
	var err error
	if len(sys.Conf.InlineExtent) > 0 {
		for _, v := range sys.Conf.InlineExtent {
			if err = gdaoHandle.AddInlineDB(v); err != nil {
				break
			}
		}
	} else if sys.Conf.InlineDB != nil {
		err = gdaoHandle.AddInlineDB(sys.Conf.InlineDB)
	} else {
		err = gdaoHandle.Add(Driver_Sqlite, localDB, 0)
	}
	if err != nil {
		log.FmtPrint(err)
		os.Exit(0)
	}
	return h
}

func (h *inlineHandle) Register(username, pwd string, domain *string) (node string, e errs.ERROR) {
	uuid := util.CreateUUID(username, domain)
	//if a, _ := SelectByIdx[timuser]("UUID", uuid); a != nil {
	//	e = sys.ERR_HASEXIST
	//	return
	//}
	//tu := &timuser{UUID: uuid, Pwd: h.password(uuid, pwd, domain), Createtime: TimeNano()}
	//if _, err := Insert(tu); err == nil {
	//	node = util.UUIDToNode(uuid)
	//} else {
	//	e = sys.ERR_DATABASE
	//}
	tu := newTimuser(uuid)
	tu.Where(tu.UUID.EQ(int64(uuid)))
	if a, _ := tu.Select(tu.ID); a != nil {
		e = errs.ERR_HASEXIST
		return
	}
	tu = newTimuser(uuid)
	if _, err := tu.SetCreatetime(TimeNano()).SetUuid(int64(uuid)).SetPwd(int64(h.password(uuid, pwd, domain))).SetTimeseries(TimeNano()).Insert(); err == nil {
		node = util.UUIDToNode(uuid)
	} else {
		e = errs.ERR_DATABASE
	}
	return
}

func (h *inlineHandle) Login(username, pwd string, domain *string) (_r string, e errs.ERROR) {
	uuid := util.CreateUUID(username, domain)
	//if a, err := SelectByIdx[timuser]("UUID", uuid); err == nil && a != nil {
	//	if a.Pwd == h.password(uuid, pwd, domain) {
	//		_r = util.UUIDToNode(uuid)
	//		return
	//	}
	//}
	tu := newTimuser(uuid)
	tu.Where(tu.UUID.EQ(int64(uuid)))
	if a, _ := tu.Select(tu.PWD); a != nil {
		if uint64(a.GetPwd()) == h.password(uuid, pwd, domain) {
			_r = util.UUIDToNode(uuid)
			return
		}
	}
	e = errs.ERR_NOPASS
	return
}

func (h *inlineHandle) Modify(uuid uint64, oldpwd *string, newpwd string, domain *string) (e errs.ERROR) {
	//if a, err := SelectByIdx[timuser]("UUID", uuid); err == nil && a != nil {
	//	if pwd != nil {
	//		if h.password(uuid, *pwd, domain) != a.Pwd {
	//			return sys.ERR_PERM_DENIED
	//		}
	//	} else if pwdLast == "" {
	//		return sys.ERR_PARAMS
	//	}
	//	UpdateNonzero(&timuser{UUID: uuid, Id: a.Id, Pwd: h.password(uuid, pwdLast, domain)})
	//} else {
	//	e = sys.ERR_NOEXIST
	//}
	if uuid == 0 {
		return errs.ERR_ACCOUNT
	}
	tu := newTimuser(uuid)
	tu.Where(tu.UUID.EQ(int64(uuid)))
	if a, _ := tu.Select(tu.ID, tu.PWD); a != nil {
		if oldpwd != nil {
			if h.password(uuid, *oldpwd, domain) != uint64(a.GetPwd()) {
				return errs.ERR_PERM_DENIED
			}
		} else if newpwd == "" {
			return errs.ERR_PARAMS
		}
		tuu := newTimuser(uuid)
		tuu.SetPwd(int64(h.password(uuid, newpwd, domain))).Where(tuu.ID.EQ(a.GetId())).Update()
	} else {
		e = errs.ERR_ACCOUNT
	}
	return
}

func (h *inlineHandle) AuthNode(username, pwd string, domain *string) (node string, err errs.ERROR) {
	uuid := util.CreateUUID(username, domain)
	//if a, _ := SelectByIdx[timuser]("UUID", uuid); a != nil {
	//	if a.Pwd == h.password(uuid, pwd, domain) {
	//		_r = util.UUIDToNode(uuid)
	//		return
	//	}
	//}
	tu := newTimuser(uuid)
	tu.Where(tu.UUID.EQ(int64(uuid)))
	if a, _ := tu.Select(tu.PWD); a != nil {
		if uint64(a.GetPwd()) == h.password(uuid, pwd, domain) {
			node = util.UUIDToNode(uuid)
			return
		}
	}
	err = errs.ERR_NOPASS
	return
}

func (h *inlineHandle) password(uuid uint64, pwd string, domain *string) uint64 {
	buf := buffer.NewBufferByPool()
	defer buf.Free()
	buf.Write(goutil.Int64ToBytes(int64(uuid)))
	buf.WriteString(sys.Conf.Salt)
	if domain != nil {
		buf.WriteString(*domain)
	}
	buf.Write([]byte(pwd))
	return goutil.CRC64(goutil.MD5(buf.Bytes()))
}

func (h *inlineHandle) SaveMessage(tm *TimMessage) (err error) {
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
	stanze := util.Mask(goutil.TEncode(tm))
	//var mid int64
	//mid, err = Insert(&timmessage{ChatId: chatId, Stanza: stanze})
	//tm.Mid = &mid

	dbhandle := gdaoHandle.GetDBHandle(chatId)
	if tx, e := dbhandle.GetTransaction(); e == nil {
		tmg := dao.NewTimmessage()
		tmg.UseTransaction(tx)
		tmg.SetChatid(int64(chatId)).SetStanza(stanze).SetTimeseries(TimeNano())
		if rs, er := tmg.Insert(); er == nil {
			if lid, ok := lastInsertId(dbhandle.GetDBType(), tx); ok {
				lid, _ = rs.LastInsertId()
				tm.Mid = &lid
			} else if lid > 0 {
				tm.Mid = &lid
			}
		}
		tx.Commit()
	} else {
		tx.Rollback()
		return e
	}
	if id == 0 {
		id = goutil.UUID64()
	}
	tm.ID = &id
	tm.FromTid = fid
	return
}

func (h *inlineHandle) GetMessage(fromNode string, domain *string, rtype int8, to string, mid, limit int64) (tmList []*TimMessage, err error) {
	chatId := uint64(0)
	if rtype == 1 {
		chatId = util.ChatIdByNode(fromNode, to, domain)
	} else {
		chatId = util.ChatIdByRoom(to, domain)
	}
	//if mid == 0 {
	//	mxid, _ := SelectIdByIdx[timmessage]("ChatId", chatId)
	//	mid, _ = SelectIdByIdxSeq[timmessage]("ChatId", chatId, mxid)
	//}
	//if as, err := SelectByIdxDescLimit[timmessage](mid, limit, "ChatId", chatId); err == nil {
	//	tmList = make([]*TimMessage, 0)
	//	for _, a := range as {
	//		if tm, err := TDecode(util.Mask(a.Stanza), &TimMessage{}); err == nil {
	//			tm.Mid = &a.Id
	//			tmList = append(tmList, tm)
	//		}
	//	}
	//}

	tmg := newTimmessage(chatId)
	if mid > 0 {
		tmg.Where(tmg.CHATID.EQ(int64(chatId)), tmg.ID.LE(mid))
	} else {
		tmg.Where(tmg.CHATID.EQ(int64(chatId)))
	}
	tmg.OrderBy(tmg.ID.Desc())
	tmg.Limit(limit)
	if list, err := tmg.Selects(tmg.ID, tmg.STANZA); err == nil {
		tmList = make([]*TimMessage, 0)
		for _, a := range list {
			if tm, err := goutil.TDecode(util.Mask(a.GetStanza()), &TimMessage{}); err == nil {
				id := a.GetId()
				tm.Mid = &id
				tmList = append(tmList, tm)
			}
		}
	}

	return
}

func (h *inlineHandle) GetMessageByMid(tid uint64, mid int64) (tm *TimMessage, err error) {
	if mid <= 0 {
		err = errs.ERR_PARAMS.Error()
		return
	}
	//if a, e := SelectById[timmessage](tid, mid); e == nil && a != nil {
	//	tm, err = TDecode(util.Mask(a.Stanza), &TimMessage{})
	//} else {
	//	err = e
	//}
	tmg := newTimmessage(tid)
	tmg.Where(tmg.ID.EQ(mid))
	if t, e := tmg.Select(tmg.STANZA); e == nil {
		tm, err = goutil.TDecode(util.Mask(t.GetStanza()), &TimMessage{})
	} else {
		err = e
	}
	return
}

func (h *inlineHandle) DelMessageByMid(tid uint64, mid int64) (err error) {
	//return Delete[timmessage](tid, mid)
	tmg := newTimmessage(tid)
	_, err = tmg.Where(tmg.ID.EQ(mid)).Delete()
	return
}

// func (h *inlineHandle) ExistOfflineMessage(tm *TimMessage) (exsit bool) {
// 	uuid := util.NodeToUUID(tm.ToTid.Node)
// 	if a, _ := SelectByIdx[timoffline]("Unik", *tm.ID); a != nil {
// 		exsit = true
// 	}
// 	return
// }

func (h *inlineHandle) SaveOfflineMessage(tm *TimMessage) (err error) {
	node := tm.ToTid.Node
	fid := tm.FromTid
	t := TimeNano()
	if tm.Timestamp != nil {
		tm.Timestamp = &t
	}
	tm.FromTid = &Tid{Node: fid.Node}
	uuid := util.NodeToUUID(node)
	if tm.OdType == sys.ORDER_INOF && tm.Mid != nil && *tm.Mid > 0 {
		cid := util.ChatIdByNode(fid.Node, tm.ToTid.Node, fid.Domain)
		//Insert(&timoffline{UUID: util.NodeToUUID(node), Mid: *tm.Mid, ChatId: cid})
		tf := newTimoffline(uuid)
		_, err = tf.SetChatid(int64(cid)).SetUuid(int64(uuid)).SetMid(int64(*tm.Mid)).SetTimeseries(t).Insert()
	} else {
		//Insert(&timoffline{UUID: util.NodeToUUID(node), Stanza: util.Mask(TEncode(tm))})
		tf := newTimoffline(uuid)
		_, err = tf.SetStanza(util.Mask(goutil.TEncode(tm))).SetUuid(int64(uuid)).SetTimeseries(t).Insert()
	}
	tm.FromTid = fid
	return
}

func (h *inlineHandle) GetOfflineMessage(node string, limit int) (oblist []*OfflineBean, err error) {
	uuid := util.NodeToUUID(node)
	if uuid == 0 {
		return
	}
	//if tfs, _ := SelectByIdxLimit[timoffline](0, int64(limit), "UUID", uuid); tfs != nil {
	//	oblist = make([]*OfflineBean, 0)
	//	for _, tf := range tfs {
	//		ob := &OfflineBean{Id: tf.Id, Mid: tf.Mid}
	//		if tf.Mid == 0 {
	//			if tf.Stanza != nil {
	//				ob.Stanze = util.Mask(tf.Stanza)
	//			}
	//		} else {
	//			if a, _ := SelectById[timmessage](tf.ChatId, tf.Mid); a != nil {
	//				ob.Stanze, ob.Mid = util.Mask(a.Stanza), tf.Mid
	//			}
	//		}
	//		oblist = append(oblist, ob)
	//	}
	//}
	tff := newTimoffline(uuid)
	tff.Where(tff.UUID.EQ(int64(uuid)))
	tff.Limit(int64(limit))
	if as, _ := tff.Selects(tff.ID, tff.MID, tff.STANZA, tff.CHATID); len(as) > 0 {
		oblist = make([]*OfflineBean, 0)
		for _, tf := range as {
			ob := &OfflineBean{Id: tf.GetId(), Mid: tf.GetMid()}
			if tf.GetMid() == 0 {
				if tf.GetStanza() != nil {
					ob.Stanze = util.Mask(tf.GetStanza())
				}
			} else {
				tmg := newTimmessage(uint64(tf.GetChatid()))
				tmg.Where(tmg.ID.EQ(tf.GetMid()))
				if a, _ := tmg.Select(tmg.STANZA); a != nil {
					ob.Stanze, ob.Mid = util.Mask(a.GetStanza()), tf.GetMid()
				}
			}
			oblist = append(oblist, ob)
		}
	}
	return
}

func (h *inlineHandle) DelOfflineMessage(tid uint64, ids ...int64) (_r int64, err error) {
	//for _, id := range ids {
	//	if id > 0 {
	//		_r, err = 1, Delete[timoffline](tid, id)
	//	}
	//}
	if size := len(ids); size > 0 {
		idsy := make([]any, size, size)
		for i := range ids {
			idsy[i] = ids[i]
		}
		tf := newTimoffline(tid)
		_, err = tf.Where(tf.ID.IN(idsy...)).Delete()
	}
	return
}

func (h *inlineHandle) AuthGroupAndUser(groupnode, usernode string, domain *string) (ok bool, err error) {
	if util.CheckNodes(groupnode, usernode) {
		//if a, _ := SelectByIdx[timgroup]("UUID", util.NodeToUUID(groupnode)); a != nil && a.Status == 1 {
		//	relateid := util.RelateIdForGroup(groupnode, usernode, domain)
		//	if a, _ := SelectByIdx[timrelate]("UUID", relateid); a != nil {
		//		ok = a.Status == 0x11
		//	}
		//}
		uuid := util.NodeToUUID(groupnode)
		tg := newTimgroup(uuid)
		tg.Where(tg.UUID.EQ(int64(uuid)))
		if a, e := tg.Select(tg.STATUS); a != nil && sys.TIMTYPE(a.GetStatus()) == sys.GROUP_STATUS_ALIVE {
			relateid := util.RelateIdForGroup(groupnode, usernode, domain)
			tr := newTimrelate(relateid)
			t, _ := tr.Where(tr.UUID.EQ(int64(relateid))).Select(tr.STATUS)
			ok = t != nil && t.GetStatus() == 0x11
		} else {
			err = e
		}
	}
	return
}

func (h *inlineHandle) AuthUserAndUser(fnode, tnode string, domain *string) (_r bool) {
	if util.CheckNodes(fnode, tnode) {
		cid := util.ChatIdByNode(fnode, tnode, domain)
		//if a, _ := SelectByIdx[timrelate]("UUID", cid); a != nil {
		//	if a.Status == 0x10|0x1 {
		//		return true
		//	}
		//}
		tr := newTimrelate(cid)
		a, _ := tr.Where(tr.UUID.EQ(int64(cid))).Select(tr.STATUS)
		_r = a != nil && a.GetStatus() == 0x10|0x1
	}
	return
}

func (h *inlineHandle) ExistUser(node string) bool {
	//if a, _ := SelectByIdx[timuser]("UUID", util.NodeToUUID(node)); a != nil {
	//	_r = true
	//}
	if uuid := util.NodeToUUID(node); uuid > 0 {
		tu := newTimuser(uuid)
		a, _ := tu.Where(tu.UUID.EQ(int64(uuid))).Select(tu.ID)
		return a != nil
	}
	return false
}

func (h *inlineHandle) ExistGroup(node string) bool {
	//if a, _ := SelectByIdx[timgroup]("UUID", util.NodeToUUID(node)); a != nil {
	//	_r = true
	//}
	if uuid := util.NodeToUUID(node); uuid > 0 {
		tg := newTimgroup(uuid)
		a, _ := tg.Where(tg.UUID.EQ(int64(uuid))).Select(tg.ID)
		return a != nil
	}
	return false
}

/*********************************************************************************************************/

func (h *inlineHandle) Roster(node string) (_r []string) {
	//if tr, err := SelectAllByIdx[timroster]("UUID", util.NodeToUUID(node)); err == nil {
	//	for _, a := range tr {
	//		_r = append(_r, util.UUIDToNode(a.TUUID))
	//	}
	//}
	if uuid := util.NodeToUUID(node); uuid > 0 {
		tr := newTimroster(uuid)
		tr.Where(tr.UUID.EQ(int64(uuid)))
		if as, _ := tr.Selects(tr.TUUID); len(as) > 0 {
			_r = make([]string, len(as))
			for i := range as {
				_r[i] = util.UUIDToNode(uint64(as[i].GetTuuid()))
			}
		}
	}
	return
}

func (h *inlineHandle) Blockrosterlist(node string) (_r []string) {
	//if tr, err := SelectAllByIdx[timblock]("UUID", util.NodeToUUID(node)); err == nil {
	//	for _, a := range tr {
	//		_r = append(_r, util.UUIDToNode(a.TUUID))
	//	}
	//}
	if uuid := util.NodeToUUID(node); uuid > 0 {
		tb := newTimblock(uuid)
		tb.Where(tb.UUID.EQ(int64(uuid)))
		if as, _ := tb.Selects(tb.TUUID); len(as) > 0 {
			_r = make([]string, len(as))
			for i := range as {
				_r[i] = util.UUIDToNode(uint64(as[i].GetTuuid()))
			}
		}
	}
	return
}

func (h *inlineHandle) Blockroomlist(node string) (_r []string) {

	//if tr, err := SelectAllByIdx[timblockroom]("UUID", util.NodeToUUID(node)); err == nil {
	//	for _, a := range tr {
	//		_r = append(_r, util.UUIDToNode(a.TUUID))
	//	}
	//}
	if uuid := util.NodeToUUID(node); uuid > 0 {
		tb := newTimblockroom(uuid)
		tb.Where(tb.UUID.EQ(int64(uuid)))
		if as, _ := tb.Selects(tb.TUUID); len(as) > 0 {
			_r = make([]string, len(as))
			for i := range as {
				_r[i] = util.UUIDToNode(uint64(as[i].GetTuuid()))
			}
		}
	}
	return
}

func (h *inlineHandle) Blockroommemberlist(node string, fnode string) (_r []string) {
	if h.checkAdmin(node, fnode, "") != nil {
		return
	}
	//if tr, err := SelectAllByIdx[timblockroom]("UUID", util.NodeToUUID(node)); err == nil {
	//	for _, a := range tr {
	//		_r = append(_r, util.UUIDToNode(a.TUUID))
	//	}
	//}
	if uuid := util.NodeToUUID(node); uuid > 0 {
		tbm := newTimblockroom(uuid)
		tbm.Where(tbm.UUID.EQ(int64(uuid)))
		if as, _ := tbm.Selects(tbm.TUUID); len(as) > 0 {
			_r = make([]string, len(as))
			for i := range as {
				_r[i] = util.UUIDToNode(uint64(as[i].GetTuuid()))
			}
		}
	}
	return
}

func (h *inlineHandle) UserGroup(node string, domain *string) (_r []string) {
	//if tr, _ := SelectAllByIdx[timmucroster]("UUID", util.NodeToUUID(node)); tr != nil && len(tr) > 0 {
	//	_r = make([]string, 0)
	//	for _, a := range tr {
	//		_r = append(_r, util.UUIDToNode(a.TUUID))
	//	}
	//}
	if uuid := util.NodeToUUID(node); uuid > 0 {
		tmu := newTimmucroster(uuid)
		tmu.Where(tmu.UUID.EQ(int64(uuid)))
		if as, _ := tmu.Selects(tmu.TUUID); len(as) > 0 {
			_r = make([]string, len(as))
			for i := range as {
				_r[i] = util.UUIDToNode(uint64(as[i].GetTuuid()))
			}
		}
	}
	return
}

func (h *inlineHandle) GroupRoster(groupnode string) (_r []string) {
	//if tr, _ := SelectAllByIdx[timmucroster]("UUID", util.NodeToUUID(groupnode)); tr != nil && len(tr) > 0 {
	//	_r = make([]string, 0)
	//	for _, a := range tr {
	//		_r = append(_r, util.UUIDToNode(a.TUUID))
	//	}
	//}
	if uuid := util.NodeToUUID(groupnode); uuid > 0 {
		tmu := newTimmucroster(uuid)
		tmu.Where(tmu.UUID.EQ(int64(uuid)))
		if as, _ := tmu.Selects(tmu.TUUID); len(as) > 0 {
			_r = make([]string, len(as))
			for i := range as {
				_r[i] = util.UUIDToNode(uint64(as[i].GetTuuid()))
			}
		}
	}
	return
}

func (h *inlineHandle) Addroster(fnode, tnode string, domain *string) (status int8, err errs.ERROR) {
	cid := util.ChatIdByNode(fnode, tnode, domain)
	uuid1, uuid2 := util.NodeToUUID(fnode), util.NodeToUUID(tnode)
	if !checkuseruuid(uuid1, uuid2) {
		err = errs.ERR_ACCOUNT
		return
	}
	numlock.Lock(int64(cid))
	defer numlock.Unlock(int64(cid))

	tr := newTimrelate(cid)
	tr.Where(tr.UUID.EQ(int64(cid)))
	if a, _ := tr.Select(tr.ID, tr.STATUS); a != nil {
		if a.GetStatus() == 0x11 {
			return 0x11, errs.ERR_REPEAT
		}
		//stat := uint8(a.GetStatus())
		tru := newTimrelate(cid) //for update
		if uuid1 > uuid2 {
			if uint8(a.GetStatus())&0x0f == 0x02 {
				err = errs.ERR_BLOCK
				return
			}
			//a.Status = 0x10 | (int8(a.GetStatus()) & 0x0f)
			tru.SetStatus(int64(0x10 | (int8(a.GetStatus()) & 0x0f)))
		} else {
			if uint8(a.GetStatus())&0xf0 == 0x20 {
				err = errs.ERR_BLOCK
				return
			}
			//a.Status = (a.Status & 0xf0) | 0x01
			tru.SetStatus(int64((uint8(a.GetStatus()) & 0xf0) | 0x01))
		}
		//if stat != uint8(tru.GetStatus()) {
		//UpdateNonzero(a)
		tru.Where(tru.ID.EQ(a.GetId())).Update()
		//}
		//status = int8(tru.GetStatus())
		//if stat != 0x11 && status == 0x11 {
		t := TimeNano()
		tri1 := newTimroster(uuid1)
		tri1.SetUuid(int64(uuid1)).SetUnikid(int64(util.UnikIdByNode(fnode, tnode, domain))).SetTuuid(int64(uuid2)).SetTimeseries(t).Insert()
		tri2 := newTimroster(uuid2)
		tri2.SetUuid(int64(uuid2)).SetUnikid(int64(util.UnikIdByNode(tnode, fnode, domain))).SetTuuid(int64(uuid1)).SetTimeseries(t).Insert()
		//}
	} else {
		status = 0x10
		if uuid1 < uuid2 {
			status = 0x01
		}
		//if id, _ := Insert(&timrelate{UUID: cid, Status: uint8(status)}); id == 0 {
		//	err = sys.ERR_DATABASE
		//}
		tr := newTimrelate(cid)
		if _, e := tr.SetUuid(int64(cid)).SetStatus(int64(status)).SetTimeseries(TimeNano()).Insert(); e != nil {
			err = errs.ERR_DATABASE
		}
	}
	return
}

func (h *inlineHandle) Rmroster(fnode, tnode string, domain *string) (mustTell bool, ok bool) {
	cid := util.ChatIdByNode(fnode, tnode, domain)
	uuid1, uuid2 := util.NodeToUUID(fnode), util.NodeToUUID(tnode)
	numlock.Lock(int64(cid))
	defer numlock.Unlock(int64(cid))
	//if as, _ := SelectAllByIdxWithTid[timroster](uuid1, "Relate", cid); as != nil {
	//	for _, a := range as {
	//		Delete[timroster](uuid1, a.Id)
	//	}
	//}
	tr := newTimroster(uuid1)
	tr.Where(tr.UNIKID.EQ(int64(util.UnikIdByNode(fnode, tnode, domain)))).Delete()

	//if as, _ := SelectAllByIdxWithTid[timroster](uuid2, "Relate", cid); as != nil {
	//	for _, a := range as {
	//		Delete[timroster](uuid2, a.Id)
	//	}
	//}
	tr = newTimroster(uuid2)
	tr.Where(tr.UNIKID.EQ(int64(util.UnikIdByNode(tnode, fnode, domain)))).Delete()

	ukid := util.UnikId(uuid1, uuid2)
	//if as, _ := SelectAllByIdxWithTid[timblock](uuid1, "UnikId", ukid); as != nil {
	//	for _, a := range as {
	//		Delete[timblock](uuid1, a.Id)
	//	}
	//}
	tb := newTimblock(uuid1)
	tb.Where(tb.UNIKID.EQ(int64(ukid))).Delete()

	tl := newTimrelate(cid)
	tl.Where(tl.UUID.EQ(int64(cid)))
	//if a, _ := SelectByIdx[timrelate]("UUID", cid); a != nil {
	if a, _ := tl.Select(tl.ID, tl.STATUS); a != nil {
		if uuid1 > uuid2 {
			if aStatus := uint8(a.GetStatus()) & 0x0f; aStatus == 0x02 {
				//UpdateNonzero(a)
				alu := newTimrelate(cid)
				alu.SetStatus(int64(aStatus)).Where(alu.ID.EQ(a.GetId())).Update()
			} else {
				mustTell = true
				//Delete[timrelate](a.Tid(), a.Id)
				alu := newTimrelate(cid)
				alu.Where(alu.ID.EQ(a.GetId())).Delete()
			}
		} else {
			if aStatus := uint8(a.GetStatus()) & 0xf0; aStatus == 0x20 {
				//UpdateNonzero(a)
				alu := newTimrelate(cid)
				alu.SetStatus(int64(aStatus)).Where(alu.ID.EQ(a.GetId())).Update()
			} else {
				mustTell = true
				//Delete[timrelate](a.Tid(), a.Id)
				alu := newTimrelate(cid)
				alu.Where(alu.ID.EQ(a.GetId())).Delete()
			}
		}
		ok = true
	}
	return
}

func (h *inlineHandle) Blockroster(fnode, tnode string, domain *string) (mustTell bool, ok bool) {
	cid := util.ChatIdByNode(fnode, tnode, domain)
	uuid1, uuid2 := util.NodeToUUID(fnode), util.NodeToUUID(tnode)
	numlock.Lock(int64(cid))
	defer numlock.Unlock(int64(cid))

	tr := newTimroster(uuid1)
	tr.Where(tr.UNIKID.EQ(int64(util.UnikIdByNode(fnode, tnode, domain)))).Delete()

	tr = newTimroster(uuid2)
	tr.Where(tr.UNIKID.EQ(int64(util.UnikIdByNode(tnode, fnode, domain)))).Delete()

	ukid := util.UnikId(uuid1, uuid2)
	tb := newTimblock(uuid1)
	tb.Where(tb.UNIKID.EQ(int64(ukid)))
	if a, _ := tb.Select(tb.ID); a == nil {
		//Insert(&timblock{UnikId: ukid, UUID: uuid1, TUUID: uuid2})
		tbi := newTimblock(uuid1)
		tbi.SetUuid(int64(uuid1)).SetTuuid(int64(uuid2)).SetUnikid(int64(ukid)).SetTimeseries(TimeNano()).Insert()
	}

	tl := newTimrelate(cid)
	tl.Where(tl.UUID.EQ(int64(cid)))
	if a, _ := tl.Select(tl.ID, tl.STATUS); a != nil {
		stat := uint8(a.GetStatus())
		var aStatus uint8
		if uuid1 > uuid2 {
			aStatus = 0x20 | (uint8(a.GetStatus()) & 0x0f)
		} else {
			aStatus = (uint8(a.GetStatus()) & 0xf0) | 0x02
		}
		if stat != aStatus {
			//if err := UpdateNonzero(a); err == nil {
			//	mstell = true
			//}
			tlu := newTimrelate(cid)
			tlu.SetStatus(int64(aStatus))
			if _, e := tlu.Where(tlu.ID.EQ(a.GetId())).Update(); e != nil {
				mustTell = true
			}
		}
		ok = true
	} else {
		status := uint8(0x20)
		if uuid1 < uuid2 {
			status = 0x02
		}
		tli := newTimrelate(cid)
		if _, e := tli.SetStatus(int64(status)).SetUuid(int64(cid)).SetTimeseries(TimeNano()).Insert(); e == nil {
			ok, mustTell = true, true
		}
		//if id, _ := Insert(&timrelate{UUID: cid, Status: status}); id > 0 {
		//	ok, mstell = true, true
		//}
	}
	return
}

func (h *inlineHandle) GroupGtype(groupnode string, domain *string) (gtype int8, err errs.ERROR) {
	if guuid := util.NodeToUUID(groupnode); guuid > 0 {
		tg := newTimgroup(guuid)
		tg.Where(tg.UUID.EQ(int64(guuid)))
		if g, _ := tg.Select(tg.GTYPE, tg.STATUS); g != nil {
			if sys.TIMTYPE(g.GetStatus()) == sys.GROUP_STATUS_CANCELLED {
				return 0, errs.ERR_CANCEL
			}
			gtype = int8(g.GetGtype())
		} else {
			return 0, errs.ERR_PARAMS
		}
	} else {
		err = errs.ERR_NOEXIST
	}
	return
}

func (h *inlineHandle) GroupManagers(groupnode string, domain *string) (s []string, err errs.ERROR) {
	if guuid := util.NodeToUUID(groupnode); guuid > 0 {
		tg := newTimgroup(guuid)
		tg.Where(tg.UUID.EQ(int64(guuid)))
		if g, _ := tg.Select(tg.STATUS, tg.RBEAN); g != nil {
			if sys.TIMTYPE(g.GetStatus()) == sys.GROUP_STATUS_CANCELLED {
				return nil, errs.ERR_CANCEL
			}
			if tr, _ := goutil.TDecode(util.Mask(g.GetRbean()), &TimRoomBean{}); tr != nil {
				s = tr.Managers
				if s == nil {
					s = []string{*tr.Founder}
				}
				if !util.ContainStrings(s, *tr.Founder) {
					s = append(s, *tr.Founder)
				}
			}
		} else {
			return nil, errs.ERR_PARAMS
		}
	} else {
		err = errs.ERR_NOEXIST
	}
	return
}

func (h *inlineHandle) Newgroup(fnode, groupname string, gtype sys.TIMTYPE, domain *string) (gnode string, err errs.ERROR) {
	fuuid := util.NodeToUUID(fnode)
	if fuuid == 0 || !checkuseruuid(fuuid) {
		err = errs.ERR_ACCOUNT
		return
	}

	guuid := util.CreateUUID(string(goutil.Int64ToBytes(goutil.UUID64())), domain)
	ctime := TimeNano()
	//tg := &timgroup{Gtype: gtype, UUID: guuid, Createtime: TimeNano(), Status: sys.GROUP_STATUS_ALIVE}
	if gtype != sys.GROUP_PRIVATE {
		gtype = sys.GROUP_OPEN
	}
	gt := int8(gtype)
	ubean := &TimRoomBean{Founder: &fnode, Topic: &groupname, Createtime: &ctime, Gtype: &gt}
	//tg.RBean = util.Mask(TEncode(ubean))
	//if id, _ := Insert(tg); id == 0 {
	//	return "", sys.ERR_DATABASE
	//}

	tg := newTimgroup(guuid)
	if _, err := tg.SetUuid(int64(guuid)).SetRbean(util.Mask(goutil.TEncode(ubean))).SetStatus(int64(sys.GROUP_STATUS_ALIVE)).SetGtype(int64(gtype)).SetTimeseries(ctime).SetCreatetime(ctime).Insert(); err != nil {
		return "", errs.ERR_DATABASE
	}

	gnode = util.UUIDToNode(guuid)
	rid := util.RelateIdForGroup(gnode, fnode, domain)

	//if id, _ := Insert(&timrelate{UUID: rid, Status: 0x11}); id == 0 {
	//	return "", sys.ERR_DATABASE
	//}
	tr := newTimrelate(rid)
	if _, e := tr.SetUuid(int64(rid)).SetStatus(0x11).SetTimeseries(ctime).Insert(); e != nil {
		return "", errs.ERR_DATABASE
	}

	//if id, _ := Insert(&timmucroster{Relate: rid, UUID: UUID, TUUID: tuuid}); id == 0 {
	//	return "", sys.ERR_DATABASE
	//}
	tu := newTimmucroster(guuid)
	tu.SetUuid(int64(guuid)).SetTimeseries(ctime).SetTuuid(int64(fuuid)).SetUnikid(int64(util.UnikIdByNode(gnode, fnode, domain))).Insert()

	//if id, _ := Insert(&timmucroster{Relate: rid, UUID: tuuid, TUUID: UUID}); id == 0 {
	//	return "", sys.ERR_DATABASE
	//}
	tu = newTimmucroster(fuuid)
	tu.SetUuid(int64(fuuid)).SetTimeseries(ctime).SetTuuid(int64(guuid)).SetUnikid(int64(util.UnikIdByNode(fnode, gnode, domain))).Insert()

	return
}

func (h *inlineHandle) Addgroup(groupnode, fromnode string, domain *string) (err errs.ERROR) {
	if guuid := util.NodeToUUID(groupnode); guuid > 0 {
		rid := util.RelateIdForGroup(groupnode, fromnode, domain)
		tg := newTimgroup(guuid)
		tg.Where(tg.UUID.EQ(int64(guuid)))
		if g, _ := tg.Select(tg.GTYPE, tg.STATUS); g != nil {
			if sys.TIMTYPE(g.GetStatus()) == sys.GROUP_STATUS_CANCELLED {
				return errs.ERR_CANCEL
			}
			numlock.Lock(int64(rid))
			defer numlock.Unlock(int64(rid))

			tr := newTimrelate(rid)
			tr.Where(tr.UUID.EQ(int64(rid)))
			if a, _ := tr.Select(tr.ID, tr.STATUS); a != nil {
				if uint8(a.GetStatus())&0xf0 == 0x20 {
					return errs.ERR_BLOCK
				}
				if uint8(a.GetStatus()) == 0x11 {
					return errs.ERR_HASEXIST
				}
				if sys.TIMTYPE(g.GetGtype()) == sys.GROUP_OPEN && uint8(a.GetStatus()) != 0x11 {
					//a.Status = 0x11
					//UpdateNonzero(a)
					ta := newTimrelate(rid)
					ta.SetStatus(0x11).Where(ta.ID.EQ(a.GetId())).Update()

					//Insert(&timmucroster{Relate: rid, UUID: guuid, TUUID: util.NodeToUUID(fromnode)})
					//Insert(&timmucroster{Relate: rid, UUID: util.NodeToUUID(fromnode), TUUID: guuid})

					ctime := TimeNano()
					fuuid := util.NodeToUUID(fromnode)
					tu1 := newTimmucroster(guuid)
					tu1.SetUnikid(int64(util.UnikIdByNode(groupnode, fromnode, domain))).SetTimeseries(ctime).SetTuuid(int64(fuuid)).SetUuid(int64(guuid)).Insert()
					tu2 := newTimmucroster(fuuid)
					tu2.SetUnikid(int64(util.UnikIdByNode(fromnode, groupnode, domain))).SetTimeseries(ctime).SetTuuid(int64(guuid)).SetUuid(int64(fuuid)).Insert()

				} else if sys.TIMTYPE(g.GetGtype()) == sys.GROUP_PRIVATE && a.GetStatus() != 0x01 {
					//a.Status = 0x01
					//UpdateNonzero(a)
					ta := newTimrelate(rid)
					ta.SetStatus(0x01).Where(ta.ID.EQ(a.GetId())).Update()
				}
				return
			} else {
				if sys.TIMTYPE(g.GetGtype()) == sys.GROUP_PRIVATE {
					//if id, _ := Insert(&timrelate{UUID: rid, Status: 0x01}); id == 0 {
					//	return sys.ERR_DATABASE
					//}
					tr := newTimrelate(rid)
					if _, err := tr.SetUuid(int64(rid)).SetTimeseries(TimeNano()).SetStatus(0x01).Insert(); err != nil {
						return errs.ERR_DATABASE
					}
				} else if sys.TIMTYPE(g.GetGtype()) == sys.GROUP_OPEN {
					//if id, _ := Insert(&timrelate{UUID: rid, Status: 0x11}); id == 0 {
					//	return sys.ERR_DATABASE
					//}
					//Insert(&timmucroster{Relate: rid, UUID: guuid, TUUID: util.NodeToUUID(fromnode)})
					//Insert(&timmucroster{Relate: rid, UUID: util.NodeToUUID(fromnode), TUUID: guuid})
					tr := newTimrelate(rid)
					if _, err := tr.SetUuid(int64(rid)).SetTimeseries(TimeNano()).SetStatus(0x11).Insert(); err != nil {
						return errs.ERR_DATABASE
					}

					ctime := TimeNano()
					fuuid := util.NodeToUUID(fromnode)
					tu1 := newTimmucroster(guuid)
					tu1.SetUnikid(int64(util.UnikIdByNode(groupnode, fromnode, domain))).SetTimeseries(ctime).SetTuuid(int64(fuuid)).SetUuid(int64(guuid)).Insert()
					tu2 := newTimmucroster(fuuid)
					tu2.SetUnikid(int64(util.UnikIdByNode(fromnode, groupnode, domain))).SetTimeseries(ctime).SetTuuid(int64(guuid)).SetUuid(int64(fuuid)).Insert()

				} else {
					return errs.ERR_PERM_DENIED
				}
			}
		}
	} else {
		err = errs.ERR_NOEXIST
	}
	return
}

func (h *inlineHandle) Pullgroup(groupnode, fromnode, tonode string, domain *string) (isReq bool, err errs.ERROR) {
	if guuid := util.NodeToUUID(groupnode); guuid > 0 {
		tg := newTimgroup(guuid)
		tg.Where(tg.UUID.EQ(int64(guuid)))
		if g, _ := tg.Select(tg.GTYPE, tg.STATUS, tg.RBEAN); g != nil {
			if sys.TIMTYPE(g.GetStatus()) == sys.GROUP_STATUS_CANCELLED {
				return isReq, errs.ERR_CANCEL
			}
			if sys.TIMTYPE(g.GetGtype()) == sys.GROUP_PRIVATE {
				if tr, _ := goutil.TDecode(util.Mask(g.GetRbean()), &TimRoomBean{}); tr != nil {
					if *tr.Founder != fromnode && !util.ContainStrings(tr.Managers, fromnode) {
						err = errs.ERR_PERM_DENIED
						return
					}
				}
			}
			rid := util.RelateIdForGroup(groupnode, tonode, domain)
			numlock.Lock(int64(rid))
			defer numlock.Unlock(int64(rid))

			tr := newTimrelate(rid)
			tr.Where(tr.UUID.EQ(int64(rid)))
			if a, _ := tr.Select(tr.ID, tr.STATUS); a != nil {
				if uint8(a.GetStatus())&0x0f == 0x02 {
					return isReq, errs.ERR_BLOCK
				}
				if uint8(a.GetStatus()) == 0x11 {
					return isReq, errs.ERR_HASEXIST
				}
				isReq = uint8(a.GetStatus()) == 0x01
				if uint8(a.GetStatus()) != 0x11 {
					//a.Status = 0x11
					//UpdateNonzero(a)
					ta := newTimrelate(rid)
					ta.SetStatus(0x11).Where(ta.ID.EQ(a.GetId())).Update()
				}
			} else {
				//if id, _ := Insert(&timrelate{UUID: rid, Status: 0x11}); id == 0{
				//	return isReq, sys.ERR_DATABASE
				//}
				tr := newTimrelate(rid)
				if _, e := tr.SetUuid(int64(rid)).SetStatus(0x11).SetTimeseries(TimeNano()).Insert(); e != nil {
					return isReq, errs.ERR_DATABASE
				}
			}
			//Insert(&timmucroster{Relate: rid, UUID: guuid, TUUID: util.NodeToUUID(tonode)})
			//Insert(&timmucroster{Relate: rid, UUID: util.NodeToUUID(tonode), TUUID: guuid})

			ctime := TimeNano()
			fuuid := util.NodeToUUID(tonode)
			tu1 := newTimmucroster(guuid)
			tu1.SetUnikid(int64(util.UnikIdByNode(groupnode, tonode, domain))).SetTimeseries(ctime).SetTuuid(int64(fuuid)).SetUuid(int64(guuid)).Insert()
			tu2 := newTimmucroster(fuuid)
			tu2.SetUnikid(int64(util.UnikIdByNode(tonode, groupnode, domain))).SetTimeseries(ctime).SetTuuid(int64(guuid)).SetUuid(int64(fuuid)).Insert()
		} else {
			return isReq, errs.ERR_NOEXIST
		}
	} else {
		err = errs.ERR_NOEXIST
	}
	return
}

func (h *inlineHandle) Nopassgroup(groupnode, fromnode, tonode string, domain *string) (err errs.ERROR) {
	if guuid := util.NodeToUUID(groupnode); guuid > 0 {
		tg := newTimgroup(guuid)
		tg.Where(tg.UUID.EQ(int64(guuid)))
		if g, _ := tg.Select(tg.STATUS, tg.RBEAN); g != nil {
			if sys.TIMTYPE(g.GetStatus()) == sys.GROUP_STATUS_CANCELLED {
				return errs.ERR_CANCEL
			}
			if tr, _ := goutil.TDecode(util.Mask(g.GetRbean()), &TimRoomBean{}); tr != nil {
				if *tr.Founder == fromnode || util.ContainStrings(tr.Managers, fromnode) {
					rid := util.RelateIdForGroup(groupnode, tonode, domain)
					numlock.Lock(int64(rid))
					defer numlock.Unlock(int64(rid))

					tr := newTimrelate(rid)
					tr.Where(tr.UUID.EQ(int64(rid)))
					if a, _ := tr.Select(tr.ID, tr.STATUS); a != nil {
						if a.GetStatus() == 0x01 {
							//Delete[timrelate](a.Tid(), a.Id)
							td := newTimrelate(rid)
							td.Where(td.ID.EQ(a.GetId())).Delete()
							return
						}
						if uint8(a.GetStatus())|0xf0 != 0 {
							return errs.ERR_EXPIREOP
						}
					} else {
						return errs.ERR_NOEXIST
					}
				} else {
					return errs.ERR_PERM_DENIED
				}
			} else {
				return errs.ERR_UNDEFINED
			}
		} else {
			return errs.ERR_NOEXIST
		}
	} else {
		err = errs.ERR_NOEXIST
	}
	return
}

func (h *inlineHandle) Kickgroup(groupnode, fromnode, tonode string, domain *string) (err errs.ERROR) {
	if guuid := util.NodeToUUID(groupnode); guuid > 0 {
		if fromnode == tonode {
			return errs.ERR_PERM_DENIED
		}
		tg := newTimgroup(guuid)
		tg.Where(tg.UUID.EQ(int64(guuid)))
		if g, _ := tg.Select(); g != nil && sys.TIMTYPE(g.GetStatus()) != sys.GROUP_STATUS_CANCELLED {
			if tr, _ := goutil.TDecode(util.Mask(g.GetRbean()), &TimRoomBean{}); tr != nil {
				if *tr.Founder == fromnode || util.ContainStrings(tr.Managers, fromnode) {
					if *tr.Founder != fromnode && util.ContainStrings(tr.Managers, tonode) {
						return errs.ERR_PERM_DENIED
					}
					if util.ContainStrings(tr.Managers, tonode) {
						tr.Managers = util.ArraySub(tr.Managers, []string{tonode})
						//UpdateNonzero(&timgroup{Id: g.Id, UUID: guuid, RBean: util.Mask(TEncode(tr))})
						tgu := newTimgroup(guuid)
						tgu.SetRbean(util.Mask(goutil.TEncode(tr))).Where(tgu.ID.EQ(g.GetId())).Update()
					}
					rid := util.RelateIdForGroup(groupnode, tonode, domain)
					numlock.Lock(int64(rid))
					defer numlock.Unlock(int64(rid))

					tmu := newTimmucroster(guuid)
					tmu.Where(tmu.UNIKID.EQ(int64(util.UnikIdByNode(groupnode, tonode, domain)))).Delete()

					tuuid := util.NodeToUUID(tonode)
					tmu = newTimmucroster(tuuid)
					tmu.Where(tmu.UNIKID.EQ(int64(util.UnikIdByNode(tonode, groupnode, domain)))).Delete()

					ukid := util.UnikId(guuid, tuuid)
					tb := newTimblockroom(guuid)
					tb.Where(tb.UNIKID.EQ(int64(ukid))).Delete()

					tr := newTimrelate(rid)
					tr.Where(tr.UUID.EQ(int64(rid)))
					if a, _ := tr.Select(tr.ID, tr.STATUS); a != nil {
						if int8(a.GetStatus())|0x0f == 0x02 {
							//a.Status = 0x02
							//if UpdateNonzero(a) != nil {
							//	err = sys.ERR_DATABASE
							//}
							tru := newTimrelate(rid)
							tru.SetStatus(0x02).Where(tru.ID.EQ(a.GetId())).Update()
						} else {
							//if Delete[timrelate](a.Tid(), a.Id) != nil {
							//	err = sys.ERR_DATABASE
							//}
							tdd := newTimrelate(rid)
							if _, e := tdd.Where(tdd.ID.EQ(a.GetId())).Delete(); e != nil {
								err = errs.ERR_DATABASE
							}
						}
					} else {
						err = errs.ERR_NOEXIST
					}
				} else {
					return errs.ERR_PERM_DENIED
				}
			} else {
				return errs.ERR_UNDEFINED
			}
		}
	} else {
		err = errs.ERR_NOEXIST
	}
	return
}

func (h *inlineHandle) Leavegroup(groupnode, fromnode string, domain *string) (err errs.ERROR) {
	guuid, tuuid := util.NodeToUUID(groupnode), util.NodeToUUID(fromnode)
	if guuid > 0 {
		if err = func() (err errs.ERROR) {
			numlock.Lock(int64(guuid))
			defer numlock.Unlock(int64(guuid))
			tg := newTimgroup(guuid)
			tg.Where(tg.UUID.EQ(int64(guuid)))
			if g, _ := tg.Select(tg.ID, tg.RBEAN, tg.STATUS); g != nil {
				if sys.TIMTYPE(g.GetStatus()) == sys.GROUP_STATUS_CANCELLED {
					return errs.ERR_CANCEL
				}
				if tr, _ := goutil.TDecode(util.Mask(g.GetRbean()), &TimRoomBean{}); tr != nil {
					if *tr.Founder == fromnode {
						return errs.ERR_PERM_DENIED
					}
					if util.ContainStrings(tr.Managers, fromnode) {
						tr.Managers = util.ArraySub(tr.Managers, []string{fromnode})
						//UpdateNonzero(&timgroup{Id: g.Id, UUID: guuid, RBean: util.Mask(TEncode(tr))})
						tgu := newTimgroup(guuid)
						tgu.SetRbean(util.Mask(goutil.TEncode(tr))).Where(tgu.ID.EQ(g.GetId())).Update()
					}
				} else {
					return errs.ERR_UNDEFINED
				}
			} else {
				err = errs.ERR_NOEXIST
			}
			return
		}(); err != nil {
			return
		}
	} else {
		err = errs.ERR_PARAMS
	}
	rid := util.RelateIdForGroup(groupnode, fromnode, domain)

	tmu := newTimmucroster(guuid)
	tmu.Where(tmu.UNIKID.EQ(int64(util.UnikIdByNode(groupnode, fromnode, domain)))).Delete()

	tmu = newTimmucroster(tuuid)
	tmu.Where(tmu.UNIKID.EQ(int64(util.UnikIdByNode(fromnode, groupnode, domain)))).Delete()

	ukid := util.UnikId(tuuid, guuid)
	tbr := newTimblockroom(tuuid)
	tbr.Where(tbr.UNIKID.EQ(int64(ukid))).Delete()

	numlock.Lock(int64(rid))
	defer numlock.Unlock(int64(rid))

	tr := newTimrelate(rid)
	tr.Where(tr.UUID.EQ(int64(rid)))
	if a, _ := tr.Select(tr.ID, tr.STATUS); a != nil {
		if uint8(a.GetStatus())&0xf0 == 0x20 {
			//a.Status = 0x20
			//if UpdateNonzero(a) != nil {
			//	err = sys.ERR_DATABASE
			//}
			tru := newTimrelate(rid)
			tru.SetStatus(0x20).Where(tru.ID.EQ(a.GetId())).Update()
		} else {
			//if Delete[timrelate](a.Tid(), a.Id) != nil {
			//	err = sys.ERR_DATABASE
			//}
			trd := newTimrelate(rid)
			if _, e := trd.Where(trd.ID.EQ(a.GetId())).Delete(); e != nil {
				err = errs.ERR_DATABASE
			}
		}
	} else {
		err = errs.ERR_NOEXIST
	}
	return
}

func (h *inlineHandle) Cancelgroup(groupnode, fromnode string, domain *string) (err errs.ERROR) {
	if guuid := util.NodeToUUID(groupnode); guuid > 0 {
		numlock.Lock(int64(guuid))
		defer numlock.Unlock(int64(guuid))
		tg := newTimgroup(guuid)
		tg.Where(tg.UUID.EQ(int64(guuid)))
		if g, _ := tg.Select(tg.ID, tg.STATUS, tg.RBEAN); g != nil {
			if sys.TIMTYPE(g.GetStatus()) == sys.GROUP_STATUS_CANCELLED {
				return errs.ERR_CANCEL
			}
			if tr, _ := goutil.TDecode(util.Mask(g.GetRbean()), &TimRoomBean{}); tr != nil {
				if *tr.Founder == fromnode {
					tu := newTimmucroster(guuid)
					tu.Where(tu.UUID.EQ(int64(guuid)))
					if tus, _ := tu.Selects(tu.ID); len(tus) > 0 {
						for _, v := range tus {
							if uint64(v.GetTuuid()) != util.NodeToUUID(fromnode) {
								return errs.ERR_PERM_DENIED
							}
						}
						ids := make([]any, len(tus))
						for i := range tus {
							//Delete[timmucroster](guuid, v.Id)
							ids[i] = tus[i].GetId()
						}
						if len(ids) > 0 {
							tmu := newTimmucroster(guuid)
							tmu.Where(tmu.ID.IN(ids...)).Delete()
						}
					}
					//g.Status = sys.GROUP_STATUS_CANCELLED
					//if UpdateNonzero(g) != nil {
					//	return sys.ERR_DATABASE
					//}
					tgu := newTimgroup(guuid)
					if _, e := tgu.SetStatus(int64(sys.GROUP_STATUS_CANCELLED)).Where(tgu.ID.EQ(g.GetId())).Update(); e != nil {
						return errs.ERR_DATABASE
					}
				} else {
					return errs.ERR_PERM_DENIED
				}
			} else {
				return errs.ERR_UNDEFINED
			}
		} else {
			err = errs.ERR_NOEXIST
		}
	} else {
		err = errs.ERR_PARAMS
	}
	return
}

func (h *inlineHandle) Blockgroup(groupnode, fromnode string, domain *string) (err errs.ERROR) {
	rid := util.RelateIdForGroup(groupnode, fromnode, domain)
	guuid, tuuid := util.NodeToUUID(groupnode), util.NodeToUUID(fromnode)

	tmu := newTimmucroster(guuid)
	tmu.Where(tmu.UNIKID.EQ(int64(util.UnikIdByNode(groupnode, fromnode, domain)))).Delete()

	tmu = newTimmucroster(tuuid)
	tmu.Where(tmu.UNIKID.EQ(int64(util.UnikIdByNode(fromnode, groupnode, domain)))).Delete()

	numlock.Lock(int64(rid))
	defer numlock.Unlock(int64(rid))

	ukid := util.UnikId(tuuid, guuid)
	//Insert(&timblockroom{UnikId: ukid, UUID: tuuid, TUUID: guuid})
	tbr := newTimblockroom(tuuid)
	tbr.SetTuuid(int64(guuid)).SetUuid(int64(tuuid)).SetUnikid(int64(ukid)).SetTimeseries(TimeNano()).Insert()

	tr := newTimrelate(rid)
	tr.Where(tr.UUID.EQ(int64(rid)))
	if a, _ := tr.Select(tr.ID, tr.STATUS); a != nil {
		if state := uint8(a.GetStatus()); state&0x0f != 0x02 {
			//a.Status = state&0xf0 | 0x02
			//if UpdateNonzero(a) != nil {
			//	err = sys.ERR_DATABASE
			//}
			tru := newTimrelate(rid)
			if _, e := tru.SetStatus(int64(state&0xf0 | 0x02)).Where(tru.ID.EQ(a.GetId())).Update(); e != nil {
				err = errs.ERR_DATABASE
			}
		}
	} else {
		//Insert(&timrelate{UUID: rid, Status: 0x02})
		tri := newTimrelate(rid)
		tri.SetStatus(0x02).SetUuid(int64(rid)).SetTimeseries(TimeNano()).Insert()
	}
	return
}

func (h *inlineHandle) Blockgroupmember(groupnode, fromnode, tonode string, domain *string) (err errs.ERROR) {
	if err = h.checkAdmin(groupnode, fromnode, tonode); err != nil {
		return
	}
	rid := util.RelateIdForGroup(groupnode, tonode, domain)
	guuid, tuuid := util.NodeToUUID(groupnode), util.NodeToUUID(tonode)

	tmu := newTimmucroster(guuid)
	tmu.Where(tmu.UNIKID.EQ(int64(util.UnikIdByNode(groupnode, tonode, domain)))).Delete()

	tmu = newTimmucroster(tuuid)
	tmu.Where(tmu.UNIKID.EQ(int64(util.UnikIdByNode(tonode, groupnode, domain)))).Delete()

	numlock.Lock(int64(rid))
	defer numlock.Unlock(int64(rid))

	ukid := util.UnikId(guuid, tuuid)
	tbr := newTimblockroom(guuid)
	tbr.Where(tbr.UNIKID.EQ(int64(ukid)))
	if a, _ := tbr.Select(tbr.ID); a == nil {
		//Insert(&timblockroom{UnikId: ukid, UUID: guuid, TUUID: tuuid})
		tbri := newTimblockroom(guuid)
		tbri.SetUnikid(int64(ukid)).SetUuid(int64(guuid)).SetTuuid(int64(tuuid)).SetTimeseries(TimeNano()).Insert()
	}

	tr := newTimrelate(rid)
	tr.Where(tr.UUID.EQ(int64(rid)))
	if a, _ := tr.Select(tr.ID, tr.STATUS); a != nil {
		if state := uint8(a.GetStatus()); state&0xf0 != 0x20 {
			//a.Status = state&0x0f | 0x20
			//if UpdateNonzero(a) != nil {
			//	err = sys.ERR_DATABASE
			//}
			tru := newTimrelate(rid)
			if _, e := tru.SetStatus(int64(state&0x0f | 0x20)).Where(tru.ID.EQ(a.GetId())).Update(); e != nil {
				err = errs.ERR_DATABASE
			}
		}
	} else {
		//Insert(&timrelate{UUID: rid, Status: 0x20})
		tri := newTimrelate(rid)
		tri.SetStatus(0x20).SetUuid(int64(rid)).SetTimeseries(TimeNano()).Insert()
	}
	return
}

func (h *inlineHandle) checkAdmin(groupnode, fromnode, tonode string) (err errs.ERROR) {
	if guuid := util.NodeToUUID(groupnode); guuid > 0 {
		if fromnode == tonode {
			return errs.ERR_PERM_DENIED
		}
		//if g, _ := SelectByIdx[timgroup]("UUID", guuid); g != nil {
		//	if g.Status == sys.GROUP_STATUS_CANCELLED {
		//		return sys.ERR_CANCEL
		//	}
		//	if tr, _ := TDecode(util.Mask(g.RBean), &TimRoomBean{}); tr != nil {
		//		if *tr.Founder == fromnode || util.ContainStrings(tr.Managers, fromnode) {
		//			if *tr.Founder != fromnode && util.ContainStrings(tr.Managers, tonode) {
		//				return sys.ERR_PERM_DENIED
		//			}
		//		} else {
		//			return sys.ERR_PERM_DENIED
		//		}
		//	} else {
		//		return sys.ERR_UNDEFINED
		//	}
		//} else {
		//	return sys.ERR_NOEXIST
		//}
		tg := newTimgroup(guuid)
		tg.Where(tg.UUID.EQ(int64(guuid)))
		if g, _ := tg.Select(tg.STATUS, tg.RBEAN); g != nil {
			if sys.TIMTYPE(g.GetStatus()) == sys.GROUP_STATUS_CANCELLED {
				return errs.ERR_CANCEL
			}
			if tr, _ := goutil.TDecode(util.Mask(g.GetRbean()), &TimRoomBean{}); tr != nil {
				if *tr.Founder == fromnode || util.ContainStrings(tr.Managers, fromnode) {
					if *tr.Founder != fromnode && util.ContainStrings(tr.Managers, tonode) {
						return errs.ERR_PERM_DENIED
					}
				} else {
					return errs.ERR_PERM_DENIED
				}
			} else {
				return errs.ERR_UNDEFINED
			}
		} else {
			return errs.ERR_NOEXIST
		}
	} else {
		err = errs.ERR_ACCOUNT
	}
	return
}

func (h *inlineHandle) ModifyUserInfo(node string, tu *TimUserBean) (err errs.ERROR) {
	uuid := util.NodeToUUID(node)
	if uuid == 0 {
		return errs.ERR_ACCOUNT
	}
	//if a, _ := SelectByIdx[timuser]("UUID", uuid); a != nil {
	//	if a.UBean != nil {
	//		if ub, _ := TDecode(util.Mask(a.UBean), &TimUserBean{}); ub != nil {
	//			if tu.Area != nil {
	//				ub.Area = tu.Area
	//			}
	//			if tu.Brithday != nil {
	//				ub.Brithday = tu.Brithday
	//			}
	//			if tu.Cover != nil {
	//				ub.Cover = tu.Cover
	//			}
	//			if tu.Extend != nil {
	//				ub.Extend = tu.Extend
	//			}
	//			if tu.Extra != nil {
	//				ub.Extra = tu.Extra
	//			}
	//			if tu.Gender != nil {
	//				ub.Gender = tu.Gender
	//			}
	//			if tu.Name != nil {
	//				ub.Name = tu.Name
	//			}
	//			if tu.NickName != nil {
	//				ub.NickName = tu.NickName
	//			}
	//			if tu.PhotoTidAlbum != nil {
	//				ub.PhotoTidAlbum = tu.PhotoTidAlbum
	//			}
	//			tu = ub
	//		}
	//	}
	//	UpdateNonzero(&timuser{Id: a.Id, UUID: a.UUID, UBean: util.Mask(TEncode(tu))})
	//} else {
	//	err = sys.ERR_NOEXIST
	//}
	tm := newTimuser(uuid)
	tm.Where(tm.UUID.EQ(int64(uuid)))
	if a, _ := tm.Select(tm.ID, tm.UBEAN); a != nil {
		if a.GetUbean() != nil {
			if ub, _ := goutil.TDecode(util.Mask(a.GetUbean()), &TimUserBean{}); ub != nil {
				if tu.Area != nil {
					ub.Area = tu.Area
				}
				if tu.Brithday != nil {
					ub.Brithday = tu.Brithday
				}
				if tu.Cover != nil {
					ub.Cover = tu.Cover
				}
				if tu.Extend != nil {
					ub.Extend = tu.Extend
				}
				if tu.Extra != nil {
					ub.Extra = tu.Extra
				}
				if tu.Gender != nil {
					ub.Gender = tu.Gender
				}
				if tu.Name != nil {
					ub.Name = tu.Name
				}
				if tu.NickName != nil {
					ub.NickName = tu.NickName
				}
				if tu.PhotoTidAlbum != nil {
					ub.PhotoTidAlbum = tu.PhotoTidAlbum
				}
				tu = ub
			}
		}
		tmu := newTimuser(uuid)
		tmu.SetUbean(util.Mask(goutil.TEncode(tu))).Where(tmu.ID.EQ(a.GetId())).Update()
	} else {
		err = errs.ERR_NOEXIST
	}
	return
}
func (h *inlineHandle) GetUserInfo(nodes []string) (m map[string]*TimUserBean, err errs.ERROR) {
	if len(nodes) > 0 {
		m = make(map[string]*TimUserBean, 0)
		for _, node := range nodes {
			uuid := util.NodeToUUID(node)
			//if a, _ := SelectByIdx[timuser]("UUID", uuid); a != nil {
			//	if a.UBean != nil {
			//		if tub, _ := TDecode(util.Mask(a.UBean), &TimUserBean{}); tub != nil {
			//			tub.Createtime = &a.Createtime
			//			m[node] = tub
			//		}
			//	}
			//}
			tm := newTimuser(uuid)
			tm.Where(tm.UUID.EQ(int64(uuid)))
			if a, _ := tm.Select(tm.UBEAN, tm.CREATETIME); a != nil {
				if a.GetUbean() != nil {
					if tub, _ := goutil.TDecode(util.Mask(a.GetUbean()), &TimUserBean{}); tub != nil {
						ct := a.GetCreatetime()
						tub.Createtime = &ct
						m[node] = tub
					}
				}
			}
		}
	} else {
		err = errs.ERR_PARAMS
	}
	return
}

func (h *inlineHandle) ModifygroupInfo(node, fnode string, tu *TimRoomBean) (err errs.ERROR) {
	if tu == nil {
		return
	}
	if guuid := util.NodeToUUID(node); guuid > 0 {
		tg := newTimgroup(guuid)
		tg.Where(tg.UUID.EQ(int64(guuid)))
		if g, _ := tg.Select(tg.ID, tg.STATUS, tg.RBEAN); g != nil {
			if sys.TIMTYPE(g.GetStatus()) == sys.GROUP_STATUS_CANCELLED {
				return errs.ERR_CANCEL
			}
			if tr, _ := goutil.TDecode(util.Mask(g.GetRbean()), &TimRoomBean{}); tr != nil {
				if *tr.Founder == fnode || util.ContainStrings(tr.Managers, fnode) {
					if *tr.Founder == fnode && tu.Managers != nil {
						tr.Managers = tu.Managers
					}
					if tu.Cover != nil {
						tr.Cover = tu.Cover
					}
					if tu.Topic != nil {
						tr.Topic = tu.Topic
					}
					if tu.Kind != nil {
						tr.Kind = tu.Kind
					}
					if tu.Label != nil {
						tr.Label = tu.Label
					}
					if tu.Extend != nil {
						tr.Extend = tu.Extend
					}
					if tu.Extra != nil {
						tr.Extra = tu.Extra
					}
					//UpdateNonzero(&timgroup{Id: g.Id, UUID: guuid, RBean: util.Mask(TEncode(tr))})
					tmg := newTimgroup(guuid)
					tmg.SetRbean(util.Mask(goutil.TEncode(tr))).Where(tmg.ID.EQ(g.GetId())).Update()
				} else {
					return errs.ERR_PERM_DENIED
				}
			} else {
				return errs.ERR_UNDEFINED
			}
		} else {
			return errs.ERR_NOEXIST
		}
	} else {
		err = errs.ERR_ACCOUNT
	}
	return
}

func (h *inlineHandle) GetGroupInfo(nodes []string) (m map[string]*TimRoomBean, err errs.ERROR) {
	if len(nodes) > 0 {
		m = make(map[string]*TimRoomBean, 0)
		for _, node := range nodes {
			if guuid := util.NodeToUUID(node); guuid > 0 {
				//if g, _ := SelectByIdx[timgroup]("UUID", guuid); g != nil && g.Status != sys.GROUP_STATUS_CANCELLED {
				//	if tr, _ := TDecode(util.Mask(g.RBean), &TimRoomBean{}); tr != nil {
				//		m[node] = tr
				//	}
				//}
				tg := newTimgroup(guuid)
				tg.Where(tg.UUID.EQ(int64(guuid)))
				if g, _ := tg.Select(tg.STATUS, tg.RBEAN); g != nil && sys.TIMTYPE(g.GetStatus()) != sys.GROUP_STATUS_CANCELLED {
					if tr, _ := goutil.TDecode(util.Mask(g.GetRbean()), &TimRoomBean{}); tr != nil {
						m[node] = tr
					}
				}
			}
		}
	} else {
		err = errs.ERR_PARAMS
	}
	return
}

func (h *inlineHandle) TimAdminAuth(account, password, domain string) bool {
	td := newTimdomain()
	td.Where(td.ADMINACCOUNT.EQ(account), td.TIMDOMAIN.EQ(domain))
	if t, _ := td.Select(td.ADMINPASSWORD); t != nil {
		return strings.EqualFold(t.GetAdminpassword(), goutil.Md5Str(password))
	}
	return false
}
