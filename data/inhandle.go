// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package data

import (
	"bytes"
	"database/sql"
	"github.com/donnie4w/gdao"
	goutil "github.com/donnie4w/gofer/util"
	"github.com/donnie4w/tim/errs"
	"github.com/donnie4w/tim/log"
	"github.com/donnie4w/tim/stub"
	"github.com/donnie4w/tim/sys"
	"github.com/donnie4w/tim/util"
	"os"
	"strings"
)

type inlineHandle struct{}

func (h *inlineHandle) init() service {
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
		err = gdaoHandle.AddInlineDB(defaultInlineDB())
	}
	if err != nil {
		log.FmtPrint(err)
		os.Exit(1)
	}
	return h
}

func (h *inlineHandle) Register(username, pwd string, domain *string) (node string, e errs.ERROR) {
	uuid := util.CreateUUID(username, domain)
	tu := newTimuser(uuid)
	tu.Where(tu.UUID.EQ(int64(uuid)))
	if a, _ := tu.Select(tu.ID); a != nil {
		e = errs.ERR_HASEXIST
		return
	}
	if hashPwd, err := util.Password(uuid, pwd, domain); err == nil && hashPwd != "" {
		tu = newTimuser(uuid)
		if _, err = tu.SetCreatetime(TimeNano()).SetUuid(int64(uuid)).SetPwd(hashPwd).SetTimeseries(TimeNano()).Insert(); err == nil {
			node = util.UUIDToNode(uuid)
			return
		}
	}
	return "", errs.ERR_DATABASE
}

func (h *inlineHandle) Login(username, pwd string, domain *string) (_r string, e errs.ERROR) {
	uuid := util.CreateUUID(username, domain)
	tu := newTimuser(uuid)
	tu.Where(tu.UUID.EQ(int64(uuid)))
	if a, _ := tu.Select(tu.PWD); a != nil {
		if util.CheckPasswordHash(uuid, pwd, domain, a.GetPwd()) {
			_r = util.UUIDToNode(uuid)
			return
		}
	}
	e = errs.ERR_NOPASS
	return
}

func (h *inlineHandle) Modify(uuid uint64, oldpwd *string, newpwd string, domain *string) (e errs.ERROR) {
	if uuid == 0 {
		return errs.ERR_ACCOUNT
	}
	tu := newTimuser(uuid)
	tu.Where(tu.UUID.EQ(int64(uuid)))
	if a, _ := tu.Select(tu.ID, tu.PWD); a != nil {
		if oldpwd != nil {
			if !util.CheckPasswordHash(uuid, *oldpwd, domain, a.GetPwd()) {
				return errs.ERR_PERM_DENIED
			}
		} else if newpwd == "" {
			return errs.ERR_PARAMS
		}
		if hashpwd, err := util.Password(uuid, newpwd, domain); err == nil && hashpwd != "" {
			tuu := newTimuser(uuid)
			tuu.SetPwd(hashpwd).Where(tuu.ID.EQ(a.GetId())).Update()
		} else {
			e = errs.ERR_ACCOUNT
		}
	} else {
		e = errs.ERR_ACCOUNT
	}
	return
}

func (h *inlineHandle) AuthNode(username, pwd string, domain *string) (node string, err errs.ERROR) {
	uuid := util.CreateUUID(username, domain)
	tu := newTimuser(uuid)
	tu.Where(tu.UUID.EQ(int64(uuid)))
	if a, _ := tu.Select(tu.PWD); a != nil {
		if util.CheckPasswordHash(uuid, pwd, domain, a.GetPwd()) {
			node = util.UUIDToNode(uuid)
			return
		}
	}
	err = errs.ERR_NOPASS
	return
}

func (h *inlineHandle) SaveMessage(tm *stub.TimMessage) (err error) {
	id := tm.GetID()
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
	dbhandle := gdaoHandle.GetDBHandle(goutil.FNVHash64(chatId))
	var mid int64
	switch dbhandle.GetDBType() {
	case gdao.MYSQL, gdao.MARIADB, gdao.SQLITE, gdao.OCEANBASE, gdao.TIDB:
		sqls := "INSERT INTO timmessage (chatid,fid,stanza,timeseries) VALUES(?,?,?,?)"
		var rs sql.Result
		if rs, err = dbhandle.ExecuteUpdate(sqls, chatId, int64(int32(goutil.FNVHash32([]byte(fid.Node)))), stanze, TimeNano()); err == nil {
			mid, err = rs.LastInsertId()
			tm.Mid = &mid
		}
	case gdao.POSTGRESQL, gdao.GREENPLUM, gdao.OPENGAUSS, gdao.COCKROACHDB, gdao.ENTERPRISEDB:
		sqls := "INSERT INTO timmessage (chatid,fid,stanza,timeseries) VALUES($1, $2, $3, $4) RETURNING id"
		if err = dbhandle.GetDB().QueryRow(sqls, chatId, int64(int32(goutil.FNVHash32([]byte(fid.Node)))), stanze, TimeNano()).Scan(&mid); err == nil {
			tm.Mid = &mid
		}
	case gdao.SQLSERVER:
		sqls := "INSERT INTO timmessage (chatid,fid,stanza,timeseries) OUTPUT INSERTED.id VALUES(@P1, @P2, @P3, @P4)"
		if err = dbhandle.GetDB().QueryRow(sqls, chatId, int64(int32(goutil.FNVHash32([]byte(fid.Node)))), stanze, TimeNano()).Scan(&mid); err == nil {
			tm.Mid = &mid
		}
	case gdao.ORACLE:
		sqls := "INSERT INTO timmessage (chatid, fid, stanza, timeseries) VALUES (:1, :2, :3, :4) RETURNING id INTO :5"
		var id64 sql.NullInt64
		if _, err = dbhandle.GetDB().Exec(sqls, chatId, int64(int32(goutil.FNVHash32([]byte(fid.Node)))), stanze, TimeNano(), sql.Named("id", sql.Out{Dest: &id64})); err == nil {
			if id64.Valid {
				mid = id64.Int64
				tm.Mid = &mid
			}
		}
	default:
		err = errs.ERR_DATABASE.Error()
	}
	if err == nil {
		if id == 0 {
			id = goutil.UUID64()
		}
		tm.ID = &id
		tm.FromTid = fid
	}
	return
}

func (h *inlineHandle) GetMessage(fromNode string, domain *string, rtype int8, to string, mid, timeseries, limit int64) (tmList []*stub.TimMessage, err error) {
	var chatId []byte
	if rtype == 1 {
		chatId = util.ChatIdByNode(fromNode, to, domain)
	} else {
		chatId = util.ChatIdByRoom(to, domain)
	}
	tmg := newTimmessage(goutil.FNVHash64(chatId))
	if mid > 0 {
		tmg.Where(tmg.CHATID.EQ(chatId), tmg.ID.LE(mid))
	} else if timeseries > 0 {
		tmg.Where(tmg.CHATID.EQ(chatId), tmg.TIMESERIES.LE(timeseries))
	} else {
		tmg.Where(tmg.CHATID.EQ(chatId))
	}
	tmg.OrderBy(tmg.ID.Desc())
	tmg.Limit(limit)
	if list, err := tmg.Selects(tmg.ID, tmg.STANZA); err == nil {
		tmList = make([]*stub.TimMessage, 0)
		for _, a := range list {
			if tm, err := goutil.TDecode(util.Mask(a.GetStanza()), &stub.TimMessage{}); err == nil {
				id := a.GetId()
				tm.Mid = &id
				tmList = append(tmList, tm)
			}
		}
	}

	return
}

func (h *inlineHandle) GetFidByMid(tid []byte, mid int64) (int64, error) {
	if mid <= 0 {
		return 0, errs.ERR_PARAMS.Error()
	}
	tmg := newTimmessage(goutil.FNVHash64(tid))
	tmg.Where(tmg.ID.EQ(mid))
	if t, err := tmg.Select(tmg.FID, tmg.CHATID); t != nil {
		if !bytes.Equal(t.GetChatid(), tid) {
			return 0, nil
		} else {
			return t.GetFid(), nil
		}
	} else {
		return 0, err
	}
}

func (h *inlineHandle) DelMessageByMid(tid []byte, mid int64) (err error) {
	tmg := newTimmessage(goutil.FNVHash64(tid))
	_, err = tmg.Where(tmg.ID.EQ(mid)).Delete()
	return
}

func (h *inlineHandle) SaveOfflineMessage(tnode string, tm *stub.TimMessage) (err error) {
	t := TimeNano()
	uuid := util.NodeToUUID(tnode)
	if tm.OdType == sys.ORDER_INOF && tm.GetMid() > 0 {
		var chatId []byte
		if tm.MsType == sys.SOURCE_ROOM {
			chatId = util.ChatIdByRoom(tm.RoomTid.Node, tm.FromTid.Domain)
		} else {
			chatId = util.ChatIdByNode(tm.FromTid.Node, tm.ToTid.Node, tm.FromTid.Domain)
		}
		_, err = newTimoffline(uuid).SetChatid(chatId).SetUuid(int64(uuid)).SetMid(int64(tm.GetMid())).SetTimeseries(t).Insert()
	} else {
		if tm.Timestamp == nil {
			tm.Timestamp = &t
		}
		fid := tm.FromTid
		if fid != nil {
			tm.FromTid = &stub.Tid{Node: fid.Node}
		}
		_, err = newTimoffline(uuid).SetStanza(util.Mask(goutil.TEncode(tm))).SetUuid(int64(uuid)).SetTimeseries(t).Insert()
		tm.FromTid = fid
	}
	return
}

func (h *inlineHandle) GetOfflineMessage(node string, limit int) (oblist []*OfflineBean, err error) {
	uuid := util.NodeToUUID(node)
	if uuid == 0 {
		return
	}
	tff := newTimoffline(uuid)
	tff.Where(tff.UUID.EQ(int64(uuid)))
	tff.Limit(int64(limit))
	if as, _ := tff.Selects(tff.ID, tff.MID, tff.STANZA, tff.CHATID); len(as) > 0 {
		oblist = make([]*OfflineBean, 0)
		idmap := make(map[uint64]map[int64]int64)
		for _, tf := range as {
			if tf.GetMid() == 0 {
				ob := &OfflineBean{Id: tf.GetId(), Mid: tf.GetMid()}
				if tf.GetStanza() != nil {
					ob.Stanze = util.Mask(tf.GetStanza())
				}
				oblist = append(oblist, ob)
			} else {
				//tmg := newTimmessage(goutil.FNVHash64(tf.GetChatid()))
				//tmg.Where(tmg.ID.EQ(tf.GetMid()))
				//if a, _ := tmg.Select(tmg.STANZA); a != nil {
				//	ob.Stanze, ob.Mid = util.Mask(a.GetStanza()), tf.GetMid()
				//}
				cid := goutil.FNVHash64(tf.GetChatid())
				mm, b := idmap[cid]
				if !b {
					mm = make(map[int64]int64)
					idmap[cid] = mm
				}
				mm[tf.GetMid()] = tf.GetId()
			}
		}
		if len(idmap) > 0 {
			for cid, mm := range idmap {
				ids := make([]any, 0)
				for k := range mm {
					ids = append(ids, k)
				}
				tmg := newTimmessage(cid)
				tmg.Where(tmg.ID.IN(ids...))
				if list, _ := tmg.Selects(tmg.STANZA, tmg.ID); len(list) > 0 {
					for _, a := range list {
						oblist = append(oblist, &OfflineBean{Id: mm[a.GetId()], Mid: a.GetId(), Stanze: util.Mask(a.GetStanza())})
					}
				}
			}
		}
	}
	return
}

func (h *inlineHandle) DelOfflineMessage(tid uint64, ids ...any) (_r int64, err error) {
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
		uuid := util.NodeToUUID(groupnode)
		tg := newTimgroup(uuid)
		tg.Where(tg.UUID.EQ(int64(uuid)))
		if a, e := tg.Select(tg.STATUS); a != nil && sys.TIMTYPE(a.GetStatus()) == sys.GROUP_STATUS_ALIVE {
			relateid := util.RelateIdForGroup(groupnode, usernode, domain)
			tr := newTimrelate(goutil.FNVHash64(relateid))
			t, _ := tr.Where(tr.UUID.EQ(relateid)).Select(tr.STATUS)
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
		tr := newTimrelate(goutil.FNVHash64(cid))
		a, _ := tr.Where(tr.UUID.EQ(cid)).Select(tr.STATUS)
		_r = a != nil && a.GetStatus() == 0x10|0x1
	}
	return
}

func (h *inlineHandle) ExistUser(node string) bool {
	if uuid := util.NodeToUUID(node); uuid > 0 {
		tu := newTimuser(uuid)
		a, _ := tu.Where(tu.UUID.EQ(int64(uuid))).Select(tu.ID)
		return a != nil
	}
	return false
}

func (h *inlineHandle) ExistGroup(node string) bool {
	if uuid := util.NodeToUUID(node); uuid > 0 {
		tg := newTimgroup(uuid)
		a, _ := tg.Where(tg.UUID.EQ(int64(uuid))).Select(tg.ID)
		return a != nil
	}
	return false
}

/*********************************************************************************************************/

func (h *inlineHandle) Roster(node string) (_r []string) {
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
	if !checkuseruuidGdao(uuid1, uuid2) {
		err = errs.ERR_ACCOUNT
		return
	}
	lock := strlock.Lock(string(cid))
	defer lock.Unlock()

	tr := newTimrelate(goutil.FNVHash64(cid))
	tr.Where(tr.UUID.EQ(cid))
	if a, _ := tr.Select(tr.ID, tr.STATUS); a != nil {
		if a.GetStatus() == 0x11 {
			return 0x11, errs.ERR_REPEAT
		}
		if uuid1 > uuid2 {
			if uint8(a.GetStatus())&0x0f == 0x02 {
				err = errs.ERR_BLOCK
				return
			}
			status = 0x10 | (int8(a.GetStatus()) & 0x0f)
		} else {
			if uint8(a.GetStatus())&0xf0 == 0x20 {
				err = errs.ERR_BLOCK
				return
			}
			status = int8((uint8(a.GetStatus()) & 0xf0) | 0x01)
		}
		tru := newTimrelate(goutil.FNVHash64(cid))
		tru.SetStatus(int64(status))
		if _, e := tru.Where(tru.ID.EQ(a.GetId())).Update(); e == nil {
			t := TimeNano()
			newTimroster(uuid1).SetUuid(int64(uuid1)).SetUnikid(util.UnikIdByNode(fnode, tnode, domain)).SetTuuid(int64(uuid2)).SetTimeseries(t).Insert()
			newTimroster(uuid2).SetUuid(int64(uuid2)).SetUnikid(util.UnikIdByNode(tnode, fnode, domain)).SetTuuid(int64(uuid1)).SetTimeseries(t).Insert()
		}
	} else {
		status = 0x10
		if uuid1 < uuid2 {
			status = 0x01
		}
		if _, e := newTimrelate(goutil.FNVHash64(cid)).SetUuid(cid).SetStatus(int64(status)).SetTimeseries(TimeNano()).Insert(); e != nil {
			err = errs.ERR_DATABASE
		}
	}
	return
}

func (h *inlineHandle) Rmroster(fnode, tnode string, domain *string) (mustTell bool, ok bool) {
	cid := util.ChatIdByNode(fnode, tnode, domain)
	uuid1, uuid2 := util.NodeToUUID(fnode), util.NodeToUUID(tnode)
	lock := strlock.Lock(string(cid))
	defer lock.Unlock()

	tr := newTimroster(uuid1)
	tr.Where(tr.UNIKID.EQ(util.UnikIdByNode(fnode, tnode, domain))).Delete()

	tr = newTimroster(uuid2)
	tr.Where(tr.UNIKID.EQ(util.UnikIdByNode(tnode, fnode, domain))).Delete()

	ukid := util.UnikId(uuid1, uuid2)

	tb := newTimblock(uuid1)
	tb.Where(tb.UNIKID.EQ(ukid)).Delete()

	tl := newTimrelate(goutil.FNVHash64(cid))
	tl.Where(tl.UUID.EQ(cid))
	if a, _ := tl.Select(tl.ID, tl.STATUS); a != nil {
		if uuid1 > uuid2 {
			if aStatus := uint8(a.GetStatus()) & 0x0f; aStatus == 0x02 {

				alu := newTimrelate(goutil.FNVHash64(cid))
				alu.SetStatus(int64(aStatus)).Where(alu.ID.EQ(a.GetId())).Update()
			} else {
				mustTell = true

				alu := newTimrelate(goutil.FNVHash64(cid))
				alu.Where(alu.ID.EQ(a.GetId())).Delete()
			}
		} else {
			if aStatus := uint8(a.GetStatus()) & 0xf0; aStatus == 0x20 {
				alu := newTimrelate(goutil.FNVHash64(cid))
				alu.SetStatus(int64(aStatus)).Where(alu.ID.EQ(a.GetId())).Update()
			} else {
				mustTell = true
				alu := newTimrelate(goutil.FNVHash64(cid))
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
	lock := strlock.Lock(string(cid))
	defer lock.Unlock()

	tr := newTimroster(uuid1)
	tr.Where(tr.UNIKID.EQ(util.UnikIdByNode(fnode, tnode, domain))).Delete()

	tr = newTimroster(uuid2)
	tr.Where(tr.UNIKID.EQ(util.UnikIdByNode(tnode, fnode, domain))).Delete()

	ukid := util.UnikId(uuid1, uuid2)
	tb := newTimblock(uuid1)
	tb.Where(tb.UNIKID.EQ(ukid))
	if a, _ := tb.Select(tb.ID); a == nil {
		newTimblock(uuid1).SetUuid(int64(uuid1)).SetTuuid(int64(uuid2)).SetUnikid(ukid).SetTimeseries(TimeNano()).Insert()
	}

	tl := newTimrelate(goutil.FNVHash64(cid))
	tl.Where(tl.UUID.EQ(cid))
	if a, _ := tl.Select(tl.ID, tl.STATUS); a != nil {
		stat := uint8(a.GetStatus())
		var aStatus uint8
		if uuid1 > uuid2 {
			aStatus = 0x20 | (uint8(a.GetStatus()) & 0x0f)
		} else {
			aStatus = (uint8(a.GetStatus()) & 0xf0) | 0x02
		}
		if stat != aStatus {
			tlu := newTimrelate(goutil.FNVHash64(cid))
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
		if _, e := newTimrelate(goutil.FNVHash64(cid)).SetStatus(int64(status)).SetUuid(cid).SetTimeseries(TimeNano()).Insert(); e == nil {
			ok, mustTell = true, true
		}
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

func (h *inlineHandle) GroupManagers(gnode string, domain *string) (s []string, err errs.ERROR) {
	if guuid := util.NodeToUUID(gnode); guuid > 0 {
		tg := newTimgroup(guuid)
		tg.Where(tg.UUID.EQ(int64(guuid)))
		if g, _ := tg.Select(tg.STATUS, tg.RBEAN); g != nil {
			if sys.TIMTYPE(g.GetStatus()) == sys.GROUP_STATUS_CANCELLED {
				return nil, errs.ERR_CANCEL
			}
			mm := make(map[string]struct{})
			if tr, _ := goutil.TDecode(util.Mask(g.GetRbean()), &stub.TimRoomBean{}); tr != nil {
				mm[tr.GetFounder()] = struct{}{}
				for _, manager := range tr.Managers {
					mm[manager] = struct{}{}
				}
				s = make([]string, 0, len(mm))
				for k := range mm {
					s = append(s, k)
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
	if fuuid == 0 || !checkuseruuidGdao(fuuid) {
		err = errs.ERR_ACCOUNT
		return
	}
	guuid := util.CreateUUID(string(goutil.Int64ToBytes(goutil.UUID64())), domain)
	ctime := TimeNano()
	if gtype != sys.GROUP_PRIVATE {
		gtype = sys.GROUP_OPEN
	}
	gt := int8(gtype)
	ubean := &stub.TimRoomBean{Founder: &fnode, Topic: &groupname, Createtime: &ctime, Gtype: &gt}

	if _, err := newTimgroup(guuid).SetUuid(int64(guuid)).SetRbean(util.Mask(goutil.TEncode(ubean))).SetStatus(int64(sys.GROUP_STATUS_ALIVE)).SetGtype(int64(gtype)).SetTimeseries(ctime).SetCreatetime(ctime).Insert(); err != nil {
		return "", errs.ERR_DATABASE
	}

	gnode = util.UUIDToNode(guuid)
	rid := util.RelateIdForGroup(gnode, fnode, domain)
	if _, e := newTimrelate(goutil.FNVHash64(rid)).SetUuid(rid).SetStatus(0x11).SetTimeseries(ctime).Insert(); e != nil {
		return "", errs.ERR_DATABASE
	}

	newTimmucroster(guuid).SetUuid(int64(guuid)).SetTimeseries(ctime).SetTuuid(int64(fuuid)).SetUnikid(util.UnikIdByNode(gnode, fnode, domain)).Insert()
	newTimmucroster(fuuid).SetUuid(int64(fuuid)).SetTimeseries(ctime).SetTuuid(int64(guuid)).SetUnikid(util.UnikIdByNode(fnode, gnode, domain)).Insert()
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
			lock := strlock.Lock(string(rid))
			defer lock.Unlock()

			tr := newTimrelate(goutil.FNVHash64(rid))
			tr.Where(tr.UUID.EQ(rid))
			if a, _ := tr.Select(tr.ID, tr.STATUS); a != nil {
				if uint8(a.GetStatus())&0xf0 == 0x20 {
					return errs.ERR_BLOCK
				}
				if uint8(a.GetStatus()) == 0x11 {
					return errs.ERR_HASEXIST
				}
				if sys.TIMTYPE(g.GetGtype()) == sys.GROUP_OPEN && uint8(a.GetStatus()) != 0x11 {
					ta := newTimrelate(goutil.FNVHash64(rid))
					ta.SetStatus(0x11).Where(ta.ID.EQ(a.GetId())).Update()

					ctime := TimeNano()
					fuuid := util.NodeToUUID(fromnode)
					newTimmucroster(guuid).SetUnikid(util.UnikIdByNode(groupnode, fromnode, domain)).SetTimeseries(ctime).SetTuuid(int64(fuuid)).SetUuid(int64(guuid)).Insert()
					newTimmucroster(fuuid).SetUnikid(util.UnikIdByNode(fromnode, groupnode, domain)).SetTimeseries(ctime).SetTuuid(int64(guuid)).SetUuid(int64(fuuid)).Insert()

				} else if sys.TIMTYPE(g.GetGtype()) == sys.GROUP_PRIVATE && a.GetStatus() != 0x01 {
					ta := newTimrelate(goutil.FNVHash64(rid))
					ta.SetStatus(0x01).Where(ta.ID.EQ(a.GetId())).Update()
				}
				return
			} else {
				if sys.TIMTYPE(g.GetGtype()) == sys.GROUP_PRIVATE {
					if _, err := newTimrelate(goutil.FNVHash64(rid)).SetUuid(rid).SetTimeseries(TimeNano()).SetStatus(0x01).Insert(); err != nil {
						return errs.ERR_DATABASE
					}
				} else if sys.TIMTYPE(g.GetGtype()) == sys.GROUP_OPEN {
					if _, err := newTimrelate(goutil.FNVHash64(rid)).SetUuid(rid).SetTimeseries(TimeNano()).SetStatus(0x11).Insert(); err != nil {
						return errs.ERR_DATABASE
					}

					ctime := TimeNano()
					fuuid := util.NodeToUUID(fromnode)
					newTimmucroster(guuid).SetUnikid(util.UnikIdByNode(groupnode, fromnode, domain)).SetTimeseries(ctime).SetTuuid(int64(fuuid)).SetUuid(int64(guuid)).Insert()
					newTimmucroster(fuuid).SetUnikid(util.UnikIdByNode(fromnode, groupnode, domain)).SetTimeseries(ctime).SetTuuid(int64(guuid)).SetUuid(int64(fuuid)).Insert()

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
				if tr, _ := goutil.TDecode(util.Mask(g.GetRbean()), &stub.TimRoomBean{}); tr != nil {
					if tr.GetFounder() != fromnode && !util.ContainStrings(tr.Managers, fromnode) {
						err = errs.ERR_PERM_DENIED
						return
					}
				}
			}
			rid := util.RelateIdForGroup(groupnode, tonode, domain)
			lock := strlock.Lock(string(rid))
			defer lock.Unlock()

			tr := newTimrelate(goutil.FNVHash64(rid))
			tr.Where(tr.UUID.EQ(rid))
			if a, _ := tr.Select(tr.ID, tr.STATUS); a != nil {
				if uint8(a.GetStatus())&0x0f == 0x02 {
					return isReq, errs.ERR_BLOCK
				}
				if uint8(a.GetStatus()) == 0x11 {
					return isReq, errs.ERR_HASEXIST
				}
				isReq = uint8(a.GetStatus()) == 0x01
				if uint8(a.GetStatus()) != 0x11 {
					ta := newTimrelate(goutil.FNVHash64(rid))
					ta.SetStatus(0x11).Where(ta.ID.EQ(a.GetId())).Update()
				}
			} else {
				tr := newTimrelate(goutil.FNVHash64(rid))
				if _, e := tr.SetUuid(rid).SetStatus(0x11).SetTimeseries(TimeNano()).Insert(); e != nil {
					return isReq, errs.ERR_DATABASE
				}
			}
			ctime := TimeNano()
			fuuid := util.NodeToUUID(tonode)
			newTimmucroster(guuid).SetUnikid(util.UnikIdByNode(groupnode, tonode, domain)).SetTimeseries(ctime).SetTuuid(int64(fuuid)).SetUuid(int64(guuid)).Insert()
			newTimmucroster(fuuid).SetUnikid(util.UnikIdByNode(tonode, groupnode, domain)).SetTimeseries(ctime).SetTuuid(int64(guuid)).SetUuid(int64(fuuid)).Insert()
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
			if tr, _ := goutil.TDecode(util.Mask(g.GetRbean()), &stub.TimRoomBean{}); tr != nil {
				if tr.GetFounder() == fromnode || util.ContainStrings(tr.Managers, fromnode) {
					rid := util.RelateIdForGroup(groupnode, tonode, domain)
					lock := strlock.Lock(string(rid))
					defer lock.Unlock()

					tr := newTimrelate(goutil.FNVHash64(rid))
					tr.Where(tr.UUID.EQ(rid))
					if a, _ := tr.Select(tr.ID, tr.STATUS); a != nil {
						if a.GetStatus() == 0x01 {
							td := newTimrelate(goutil.FNVHash64(rid))
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
		if g, _ := tg.Select(tg.ID, tg.STATUS, tg.RBEAN); g != nil && sys.TIMTYPE(g.GetStatus()) != sys.GROUP_STATUS_CANCELLED {
			if tr, _ := goutil.TDecode(util.Mask(g.GetRbean()), &stub.TimRoomBean{}); tr != nil {
				if tr.GetFounder() == fromnode || util.ContainStrings(tr.Managers, fromnode) {
					if tr.GetFounder() != fromnode && util.ContainStrings(tr.Managers, tonode) {
						return errs.ERR_PERM_DENIED
					}
					if util.ContainStrings(tr.Managers, tonode) {
						tr.Managers = util.ArraySub(tr.Managers, []string{tonode})
						//UpdateNonzero(&timgroup{Id: g.Id, UUID: guuid, RBean: util.Mask(TEncode(tr))})
						tgu := newTimgroup(guuid)
						tgu.SetRbean(util.Mask(goutil.TEncode(tr))).Where(tgu.ID.EQ(g.GetId())).Update()
					}
					rid := util.RelateIdForGroup(groupnode, tonode, domain)
					lock := strlock.Lock(string(rid))
					defer lock.Unlock()

					tmu := newTimmucroster(guuid)
					tmu.Where(tmu.UNIKID.EQ(util.UnikIdByNode(groupnode, tonode, domain))).Delete()

					tuuid := util.NodeToUUID(tonode)
					tmu = newTimmucroster(tuuid)
					tmu.Where(tmu.UNIKID.EQ(util.UnikIdByNode(tonode, groupnode, domain))).Delete()

					ukid := util.UnikId(guuid, tuuid)
					tb := newTimblockroom(guuid)
					tb.Where(tb.UNIKID.EQ(ukid)).Delete()

					tr := newTimrelate(goutil.FNVHash64(rid))
					tr.Where(tr.UUID.EQ(rid))
					if a, _ := tr.Select(tr.ID, tr.STATUS); a != nil {
						if int8(a.GetStatus())|0x0f == 0x02 {
							tru := newTimrelate(goutil.FNVHash64(rid))
							tru.SetStatus(0x02).Where(tru.ID.EQ(a.GetId())).Update()
						} else {
							tdd := newTimrelate(goutil.FNVHash64(rid))
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
			lock := numlock.Lock(int64(guuid))
			defer lock.Unlock()
			tg := newTimgroup(guuid)
			tg.Where(tg.UUID.EQ(int64(guuid)))
			if g, _ := tg.Select(tg.ID, tg.RBEAN, tg.STATUS); g != nil {
				if sys.TIMTYPE(g.GetStatus()) == sys.GROUP_STATUS_CANCELLED {
					return errs.ERR_CANCEL
				}
				if tr, _ := goutil.TDecode(util.Mask(g.GetRbean()), &stub.TimRoomBean{}); tr != nil {
					if tr.GetFounder() == fromnode {
						return errs.ERR_PERM_DENIED
					}
					if util.ContainStrings(tr.Managers, fromnode) {
						tr.Managers = util.ArraySub(tr.Managers, []string{fromnode})
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
	tmu.Where(tmu.UNIKID.EQ(util.UnikIdByNode(groupnode, fromnode, domain))).Delete()

	tmu = newTimmucroster(tuuid)
	tmu.Where(tmu.UNIKID.EQ(util.UnikIdByNode(fromnode, groupnode, domain))).Delete()

	ukid := util.UnikId(tuuid, guuid)
	tbr := newTimblockroom(tuuid)
	tbr.Where(tbr.UNIKID.EQ(ukid)).Delete()

	lock := strlock.Lock(string(rid))
	defer lock.Unlock()

	tr := newTimrelate(goutil.FNVHash64(rid))
	tr.Where(tr.UUID.EQ(rid))
	if a, _ := tr.Select(tr.ID, tr.STATUS); a != nil {
		if uint8(a.GetStatus())&0xf0 == 0x20 {
			tru := newTimrelate(goutil.FNVHash64(rid))
			tru.SetStatus(0x20).Where(tru.ID.EQ(a.GetId())).Update()
		} else {
			trd := newTimrelate(goutil.FNVHash64(rid))
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
		lock := numlock.Lock(int64(guuid))
		defer lock.Unlock()
		tg := newTimgroup(guuid)
		tg.Where(tg.UUID.EQ(int64(guuid)))
		if g, _ := tg.Select(tg.ID, tg.STATUS, tg.RBEAN); g != nil {
			if sys.TIMTYPE(g.GetStatus()) == sys.GROUP_STATUS_CANCELLED {
				return errs.ERR_CANCEL
			}
			if tr, _ := goutil.TDecode(util.Mask(g.GetRbean()), &stub.TimRoomBean{}); tr != nil {
				if tr.GetFounder() == fromnode {
					tu := newTimmucroster(guuid)
					tu.Where(tu.UUID.EQ(int64(guuid)))
					if tus, _ := tu.Selects(tu.TUUID); len(tus) == 1 {
						tuuid := util.NodeToUUID(fromnode)
						if uint64(tus[0].GetTuuid()) != tuuid {
							return errs.ERR_PERM_DENIED
						}
						ids := make([]any, len(tus))
						for i := range tus {
							ids[i] = tus[i].GetId()
						}
						if len(ids) > 0 {
							tmu := newTimmucroster(guuid)
							tmu.Where(tmu.ID.IN(ids...)).Delete()
						}
						tgu := newTimgroup(guuid)
						if _, e := tgu.SetStatus(int64(sys.GROUP_STATUS_CANCELLED)).Where(tgu.ID.EQ(g.GetId())).Update(); e != nil {
							return errs.ERR_DATABASE
						}

						tmu := newTimmucroster(guuid)
						tmu.Where(tmu.UNIKID.EQ(util.UnikIdByNode(groupnode, fromnode, domain))).Delete()

						tmu = newTimmucroster(tuuid)
						tmu.Where(tmu.UNIKID.EQ(util.UnikIdByNode(fromnode, groupnode, domain))).Delete()
					} else {
						return errs.ERR_PERM_DENIED
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
	tmu.Where(tmu.UNIKID.EQ(util.UnikIdByNode(groupnode, fromnode, domain))).Delete()

	tmu = newTimmucroster(tuuid)
	tmu.Where(tmu.UNIKID.EQ(util.UnikIdByNode(fromnode, groupnode, domain))).Delete()

	lock := strlock.Lock(string(rid))
	defer lock.Unlock()

	ukid := util.UnikId(tuuid, guuid)

	newTimblockroom(tuuid).SetTuuid(int64(guuid)).SetUuid(int64(tuuid)).SetUnikid(ukid).SetTimeseries(TimeNano()).Insert()

	uuid := goutil.FNVHash64(rid)
	tr := newTimrelate(uuid)
	tr.Where(tr.UUID.EQ(rid))
	if a, _ := tr.Select(tr.ID, tr.STATUS); a != nil {
		if state := uint8(a.GetStatus()); state&0x0f != 0x02 {
			tru := newTimrelate(uuid)
			if _, e := tru.SetStatus(int64(state&0xf0 | 0x02)).Where(tru.ID.EQ(a.GetId())).Update(); e != nil {
				err = errs.ERR_DATABASE
			}
		}
	} else {
		newTimrelate(uuid).SetStatus(0x02).SetUuid(rid).SetTimeseries(TimeNano()).Insert()
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
	tmu.Where(tmu.UNIKID.EQ(util.UnikIdByNode(groupnode, tonode, domain))).Delete()

	tmu = newTimmucroster(tuuid)
	tmu.Where(tmu.UNIKID.EQ(util.UnikIdByNode(tonode, groupnode, domain))).Delete()

	lock := strlock.Lock(string(rid))
	defer lock.Unlock()

	ukid := util.UnikId(guuid, tuuid)
	tbr := newTimblockroom(guuid)
	tbr.Where(tbr.UNIKID.EQ(ukid))
	if a, _ := tbr.Select(tbr.ID); a == nil {
		newTimblockroom(guuid).SetUnikid(ukid).SetUuid(int64(guuid)).SetTuuid(int64(tuuid)).SetTimeseries(TimeNano()).Insert()
	}

	tr := newTimrelate(goutil.FNVHash64(rid))
	tr.Where(tr.UUID.EQ(rid))
	if a, _ := tr.Select(tr.ID, tr.STATUS); a != nil {
		if state := uint8(a.GetStatus()); state&0xf0 != 0x20 {
			tru := newTimrelate(goutil.FNVHash64(rid))
			if _, e := tru.SetStatus(int64(state&0x0f | 0x20)).Where(tru.ID.EQ(a.GetId())).Update(); e != nil {
				err = errs.ERR_DATABASE
			}
		}
	} else {
		newTimrelate(goutil.FNVHash64(rid)).SetStatus(0x20).SetUuid(rid).SetTimeseries(TimeNano()).Insert()
	}
	return
}

func (h *inlineHandle) checkAdmin(groupnode, fromnode, tonode string) (err errs.ERROR) {
	if guuid := util.NodeToUUID(groupnode); guuid > 0 {
		if fromnode == tonode {
			return errs.ERR_PERM_DENIED
		}
		tg := newTimgroup(guuid)
		tg.Where(tg.UUID.EQ(int64(guuid)))
		if g, _ := tg.Select(tg.STATUS, tg.RBEAN); g != nil {
			if sys.TIMTYPE(g.GetStatus()) == sys.GROUP_STATUS_CANCELLED {
				return errs.ERR_CANCEL
			}
			if tr, _ := goutil.TDecode(util.Mask(g.GetRbean()), &stub.TimRoomBean{}); tr != nil {
				if tr.GetFounder() == fromnode || util.ContainStrings(tr.Managers, fromnode) {
					if tr.GetFounder() != fromnode && util.ContainStrings(tr.Managers, tonode) {
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

func (h *inlineHandle) ModifyUserInfo(node string, tu *stub.TimUserBean) (err errs.ERROR) {
	uuid := util.NodeToUUID(node)
	if uuid == 0 {
		return errs.ERR_ACCOUNT
	}
	tm := newTimuser(uuid)
	tm.Where(tm.UUID.EQ(int64(uuid)))
	if a, _ := tm.Select(tm.ID, tm.UBEAN); a != nil {
		if a.GetUbean() != nil {
			if ub, _ := goutil.TDecode(util.Mask(a.GetUbean()), &stub.TimUserBean{}); ub != nil {
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
func (h *inlineHandle) GetUserInfo(nodes []string) (m map[string]*stub.TimUserBean, err errs.ERROR) {
	if len(nodes) > 0 {
		m = make(map[string]*stub.TimUserBean, 0)
		for _, node := range nodes {
			uuid := util.NodeToUUID(node)
			tm := newTimuser(uuid)
			tm.Where(tm.UUID.EQ(int64(uuid)))
			if a, _ := tm.Select(tm.UBEAN, tm.CREATETIME); a != nil {
				if a.GetUbean() != nil {
					if tub, _ := goutil.TDecode(util.Mask(a.GetUbean()), &stub.TimUserBean{}); tub != nil {
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

func (h *inlineHandle) ModifygroupInfo(node, fnode string, tu *stub.TimRoomBean, admin bool) (err errs.ERROR) {
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
			if tr, _ := goutil.TDecode(util.Mask(g.GetRbean()), &stub.TimRoomBean{}); tr != nil {
				if admin {
					if tu.Founder != nil {
						tr.Founder = tu.Founder
					}
					if tu.Managers != nil {
						tr.Managers = tu.Managers
					}
				} else if tr.GetFounder() != fnode && !util.ContainStrings(tr.Managers, fnode) {
					return errs.ERR_PERM_DENIED
				}
				if !admin {
					if tr.GetFounder() == fnode && tu.Managers != nil {
						tr.Managers = tu.Managers
					}
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
				tmg := newTimgroup(guuid)
				tmg.SetRbean(util.Mask(goutil.TEncode(tr))).Where(tmg.ID.EQ(g.GetId())).Update()
			}
		} else {
			return errs.ERR_NOEXIST
		}
	} else {
		return errs.ERR_ACCOUNT
	}
	return
}

func (h *inlineHandle) GetGroupInfo(nodes []string) (m map[string]*stub.TimRoomBean, err errs.ERROR) {
	if len(nodes) > 0 {
		m = make(map[string]*stub.TimRoomBean, 0)
		for _, node := range nodes {
			if guuid := util.NodeToUUID(node); guuid > 0 {
				tg := newTimgroup(guuid)
				tg.Where(tg.UUID.EQ(int64(guuid)))
				if g, _ := tg.Select(tg.STATUS, tg.RBEAN); g != nil && sys.TIMTYPE(g.GetStatus()) != sys.GROUP_STATUS_CANCELLED {
					if tr, _ := goutil.TDecode(util.Mask(g.GetRbean()), &stub.TimRoomBean{}); tr != nil {
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
