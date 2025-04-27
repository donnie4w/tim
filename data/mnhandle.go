// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package data

import (
	"bytes"
	goutil "github.com/donnie4w/gofer/util"
	"github.com/donnie4w/tim/errs"
	"github.com/donnie4w/tim/stub"
	"github.com/donnie4w/tim/sys"
	"github.com/donnie4w/tim/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
)

type mongoHandle struct {
}

func (h *mongoHandle) init() service {
	if err := manager.init(); err != nil {
		panic(err)
	}
	return h
}

func (h *mongoHandle) Register(username, pwd string, domain *string) (node string, e errs.ERROR) {
	uuid := util.CreateUUID(username, domain)
	tu := newMnTimUser(uuid)
	if a, _ := tu.Get(bson.M{"uuid": int64(uuid)}, bson.M{"_id": 1}); a != nil {
		e = errs.ERR_HASEXIST
		return
	}
	if hashPwd, err := util.Password(uuid, pwd, domain); err == nil && hashPwd != "" {
		tu = newMnTimUser(uuid)
		tu.Createtime = TimeNano()
		tu.UUID = int64(uuid)
		tu.Pwd = hashPwd
		tu.Timeseries = TimeNano()
		if _, err = tu.Create(); err == nil {
			node = util.UUIDToNode(uuid)
			return
		}
	}
	return "", errs.ERR_DATABASE
}

func (h *mongoHandle) Login(username, pwd string, domain *string) (_r string, e errs.ERROR) {
	uuid := util.CreateUUID(username, domain)
	tu := newMnTimUser(uuid)
	if a, _ := tu.Get(bson.M{"uuid": int64(uuid)}, bson.M{"pwd": 1}); a != nil {
		if util.CheckPasswordHash(uuid, pwd, domain, a.Pwd) {
			_r = util.UUIDToNode(uuid)
			return
		}
	}
	e = errs.ERR_NOPASS
	return
}

func (h *mongoHandle) Modify(uuid uint64, oldpwd *string, newpwd string, domain *string) (e errs.ERROR) {
	if uuid == 0 {
		return errs.ERR_ACCOUNT
	}
	tu := newMnTimUser(uuid)
	if a, _ := tu.Get(bson.M{"uuid": int64(uuid)}, bson.M{"pwd": 1, "_id": 1}); a != nil {
		if oldpwd != nil {
			if !util.CheckPasswordHash(uuid, *oldpwd, domain, a.Pwd) {
				return errs.ERR_PERM_DENIED
			}
		} else if newpwd == "" {
			return errs.ERR_PARAMS
		}
		if hashpwd, err := util.Password(uuid, newpwd, domain); err == nil && hashpwd != "" {
			newMnTimUser(uuid).Update(a.ID, bson.M{"pwd": hashpwd})
		} else {
			e = errs.ERR_ACCOUNT
		}
	} else {
		e = errs.ERR_ACCOUNT
	}
	return
}

func (h *mongoHandle) AuthNode(username, pwd string, domain *string) (node string, err errs.ERROR) {
	uuid := util.CreateUUID(username, domain)
	tu := newMnTimUser(uuid)
	if a, _ := tu.Get(bson.M{"uuid": int64(uuid)}, bson.M{"pwd": 1}); a != nil {
		if util.CheckPasswordHash(uuid, pwd, domain, a.Pwd) {
			node = util.UUIDToNode(uuid)
			return
		}
	}
	err = errs.ERR_NOPASS
	return
}

func (h *mongoHandle) SaveMessage(tm *stub.TimMessage) (err error) {
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

	tmg := newMnTimMessage(goutil.FNVHash64(chatId))
	tmg.ChatID = binary(chatId)
	tmg.Mid = midUUID()
	tmg.FID = int32(goutil.FNVHash32([]byte(fid.Node)))
	tmg.Stanza = binary(util.Mask(goutil.TEncode(tm)))
	tmg.Timeseries = TimeNano()
	if _, err = tmg.Create(); err != nil {
		return
	}
	tm.Mid = &tmg.Mid
	if id == 0 {
		id = goutil.UUID64()
	}

	tm.ID = &id
	tm.FromTid = fid
	return
}

func (h *mongoHandle) GetMessage(fromNode string, domain *string, rtype int8, to string, mid, timeseries, limit int64) (tmList []*stub.TimMessage, err error) {
	var chatId []byte
	if rtype == 1 {
		chatId = util.ChatIdByNode(fromNode, to, domain)
	} else {
		chatId = util.ChatIdByRoom(to, domain)
	}
	tmg := newMnTimMessage(goutil.FNVHash64(chatId))
	opts := options.Find().SetProjection(bson.M{"stanza": 1, "mid": 1}).SetSort(bson.M{"timeseries": -1}).SetLimit(limit)
	var list []*timMessage
	if timeseries > 0 {
		list, err = tmg.ListOptions(bson.M{"chatid": binary(chatId), "timeseries": bson.M{"$lte": timeseries}}, opts)
	} else {
		list, err = tmg.ListOptions(bson.M{"chatid": binary(chatId)}, opts)
	}
	if err == nil {
		tmList = make([]*stub.TimMessage, 0)
		for _, a := range list {
			if tm, err := goutil.TDecode(util.Mask(a.Stanza.Data), &stub.TimMessage{}); err == nil {
				tm.Mid = &a.Mid
				tmList = append(tmList, tm)
			}
		}
	}
	return
}

func (h *mongoHandle) GetFidByMid(tid []byte, mid int64) (int64, error) {
	if mid == 0 {
		return 0, errs.ERR_PARAMS.Error()
	}
	if t, err := newMnTimMessage(goutil.FNVHash64(tid)).Get(bson.M{"mid": mid}, bson.M{"chatid": 1, "fid": 1}); t != nil {
		if !bytes.Equal(t.ChatID.Data, tid) {
			return 0, nil
		} else {
			return int64(t.FID), nil
		}
	} else {
		return 0, err
	}
}

func (h *mongoHandle) DelMessageByMid(tid []byte, mid int64) (err error) {
	tmm := newMnTimMessage(goutil.FNVHash64(tid))
	_, err = tmm.Delete(tid, mid)
	return
}

func (h *mongoHandle) SaveOfflineMessage(tnode string, tm *stub.TimMessage) (err error) {
	t := TimeNano()
	uuid := util.NodeToUUID(tnode)
	tf := newMnTimOffline(uuid)
	if tm.OdType == sys.ORDER_INOF && tm.GetMid() > 0 {
		var chatId []byte
		if tm.MsType == sys.SOURCE_ROOM {
			chatId = util.ChatIdByRoom(tm.RoomTid.Node, tm.FromTid.Domain)
		} else {
			chatId = util.ChatIdByNode(tm.FromTid.Node, tm.ToTid.Node, tm.FromTid.Domain)
		}
		tf.Chatid = binary(chatId)
		tf.UUID = int64(uuid)
		tf.Mid = tm.GetMid()
		tf.Timeseries = t
		_, err = tf.Create()
	} else {
		fid := tm.FromTid
		if fid != nil {
			tm.FromTid = &stub.Tid{Node: fid.Node}
		}
		if tm.Timestamp == nil {
			tm.Timestamp = &t
		}
		tf.Stanza = binary(util.Mask(goutil.TEncode(tm)))
		tf.UUID = int64(uuid)
		tf.Timeseries = t
		_, err = tf.Create()
		tm.FromTid = fid
	}
	return
}

func (h *mongoHandle) GetOfflineMessage(node string, limit int) (oblist []*OfflineBean, err error) {
	uuid := util.NodeToUUID(node)
	if uuid == 0 {
		return
	}
	opts := options.Find().SetProjection(bson.M{"_id": 1, "mid": 1, "stanza": 1, "chatid": 1}).SetSort(bson.M{"timeseries": 1}).SetLimit(int64(limit))
	if as, _ := newMnTimOffline(uuid).ListOptions(bson.M{"uuid": int64(uuid)}, opts); len(as) > 0 {
		oblist = make([]*OfflineBean, 0)
		idmap := make(map[uint64]map[int64]primitive.ObjectID)
		for _, tf := range as {
			if tf.Mid == 0 {
				ob := &OfflineBean{Id: tf.ID, Mid: tf.Mid}
				if tf.Stanza.Data != nil {
					ob.Stanze = util.Mask(tf.Stanza.Data)
				}
				oblist = append(oblist, ob)
			} else {
				//if a, _ := newMnTimMessage(goutil.FNVHash64(tf.Chatid.Data)).Get(bson.M{"chatid": tf.Chatid, "mid": tf.Mid}, bson.M{"stanza": 1}); a != nil {
				//	ob.Stanze, ob.Mid = util.Mask(a.Stanza.Data), tf.Mid
				//}
				cid := goutil.FNVHash64(tf.Chatid.Data)
				mm, b := idmap[cid]
				if !b {
					mm = make(map[int64]primitive.ObjectID)
					idmap[cid] = mm
				}
				mm[tf.Mid] = tf.ID
			}
		}
		if len(idmap) > 0 {
			for cid, mm := range idmap {
				ids := make([]int64, 0)
				for k := range mm {
					ids = append(ids, k)
				}
				if list, _ := newMnTimMessage(cid).List(bson.M{"mid": bson.M{"$in": ids}}, bson.M{"mid": 1, "stanza": 1, "chatid": 1}); len(list) > 0 {
					for _, tm := range list {
						if goutil.FNVHash64(tm.ChatID.Data) == cid {
							oblist = append(oblist, &OfflineBean{Id: mm[tm.Mid], Mid: tm.Mid, Stanze: util.Mask(tm.Stanza.Data)})
						}
					}
				}
			}
		}
	}
	return
}

func (h *mongoHandle) DelOfflineMessage(tid uint64, ids ...any) (_r int64, err error) {
	if size := len(ids); size > 0 {
		newMnTimOffline(tid).DeleteMany(bson.M{"_id": bson.M{"$in": ids}})
	}
	return
}

func (h *mongoHandle) AuthGroupAndUser(groupnode, usernode string, domain *string) (ok bool, err error) {
	if util.CheckNodes(groupnode, usernode) {
		uuid := util.NodeToUUID(groupnode)
		if a, e := newMnTimGroup(uuid).Get(bson.M{"uuid": int64(uuid)}, bson.M{"status": 1}); a != nil && sys.TIMTYPE(a.Status) == sys.GROUP_STATUS_ALIVE {
			relateid := util.RelateIdForGroup(groupnode, usernode, domain)
			if t, _ := newMnTimRelate(goutil.FNVHash64(relateid)).Get(bson.M{"uuid": binary(relateid)}, bson.M{"status": 1}); t != nil {
				ok = t != nil && t.Status == 0x11
			}
		} else {
			err = e
		}

	}
	return
}

func (h *mongoHandle) AuthUserAndUser(fnode, tnode string, domain *string) (_r bool) {
	if util.CheckNodes(fnode, tnode) {
		cid := util.ChatIdByNode(fnode, tnode, domain)
		a, _ := newMnTimRelate(goutil.FNVHash64(cid)).Get(bson.M{"uuid": cid}, bson.M{"status": 1})
		_r = a != nil && a.Status == 0x10|0x1
	}
	return
}

func (h *mongoHandle) ExistUser(node string) bool {
	if uuid := util.NodeToUUID(node); uuid > 0 {
		a, _ := newMnTimUser(uuid).Get(bson.M{"uuid": int64(uuid)}, bson.M{"_id": 1})
		return a != nil
	}
	return false
}

func (h *mongoHandle) ExistGroup(node string) bool {
	if uuid := util.NodeToUUID(node); uuid > 0 {
		a, _ := newMnTimGroup(uuid).Get(bson.M{"uuid": int64(uuid)}, bson.M{"_id": 1})
		return a != nil
	}
	return false
}

/*********************************************************************************************************/

func (h *mongoHandle) Roster(node string) (_r []string) {
	if uuid := util.NodeToUUID(node); uuid > 0 {
		if as, _ := newMnTimRoster(uuid).List(bson.M{"uuid": int64(uuid)}, bson.M{"tuuid": 1}); len(as) > 0 {
			_r = make([]string, len(as))
			for i := range as {
				_r[i] = util.UUIDToNode(uint64(as[i].Tuuid))
			}
		}
	}
	return
}

func (h *mongoHandle) Blockrosterlist(node string) (_r []string) {
	if uuid := util.NodeToUUID(node); uuid > 0 {
		if as, _ := newMnTimBlock(uuid).List(bson.M{"uuid": int64(uuid)}, bson.M{"tuuid": 1}); len(as) > 0 {
			_r = make([]string, len(as))
			for i := range as {
				_r[i] = util.UUIDToNode(uint64(as[i].Tuuid))
			}
		}
	}
	return
}

func (h *mongoHandle) Blockroomlist(node string) (_r []string) {
	if uuid := util.NodeToUUID(node); uuid > 0 {
		if as, _ := newMnTimBlock(uuid).List(bson.M{"uuid": int64(uuid)}, bson.M{"tuuid": 1}); len(as) > 0 {
			_r = make([]string, len(as))
			for i := range as {
				_r[i] = util.UUIDToNode(uint64(as[i].Tuuid))
			}
		}
	}
	return
}

func (h *mongoHandle) Blockroommemberlist(node string, fnode string) (_r []string) {
	if h.checkAdmin(node, fnode, "") != nil {
		return
	}
	if uuid := util.NodeToUUID(node); uuid > 0 {
		if as, _ := newMnTimBlockroom(uuid).List(bson.M{"uuid": int64(uuid)}, bson.M{"tuuid": 1}); len(as) > 0 {
			_r = make([]string, len(as))
			for i := range as {
				_r[i] = util.UUIDToNode(uint64(as[i].Tuuid))
			}
		}
	}
	return
}

func (h *mongoHandle) UserGroup(node string, domain *string) (_r []string) {
	if uuid := util.NodeToUUID(node); uuid > 0 {
		if as, _ := newMnTimMucroster(uuid).List(bson.M{"uuid": int64(uuid)}, bson.M{"tuuid": 1}); len(as) > 0 {
			_r = make([]string, len(as))
			for i := range as {
				_r[i] = util.UUIDToNode(uint64(as[i].Tuuid))
			}
		}
	}
	return
}

func (h *mongoHandle) GroupRoster(groupnode string) (_r []string) {
	if uuid := util.NodeToUUID(groupnode); uuid > 0 {
		if as, _ := newMnTimMucroster(uuid).List(bson.M{"uuid": int64(uuid)}, bson.M{"tuuid": 1}); len(as) > 0 {
			_r = make([]string, len(as))
			for i := range as {
				_r[i] = util.UUIDToNode(uint64(as[i].Tuuid))
			}
		}
	}
	return
}

func (h *mongoHandle) Addroster(fnode, tnode string, domain *string) (status int8, err errs.ERROR) {
	cid := util.ChatIdByNode(fnode, tnode, domain)
	uuid1, uuid2 := util.NodeToUUID(fnode), util.NodeToUUID(tnode)
	if !checkuseruuidMn(uuid1, uuid2) {
		err = errs.ERR_ACCOUNT
		return
	}
	lock := strlock.Lock(string(cid))
	defer lock.Unlock()

	if a, _ := newMnTimRelate(goutil.FNVHash64(cid)).Get(bson.M{"uuid": binary(cid)}, bson.M{"_id": 1, "status": 1}); a != nil {
		if a.Status == 0x11 {
			return 0x11, errs.ERR_REPEAT
		}
		if uuid1 > uuid2 {
			if uint8(a.Status)&0x0f == 0x02 {
				err = errs.ERR_BLOCK
				return
			}
			status = 0x10 | (int8(a.Status) & 0x0f)
		} else {
			if uint8(a.Status)&0xf0 == 0x20 {
				err = errs.ERR_BLOCK
				return
			}
			status = int8((uint8(a.Status) & 0xf0) | 0x01)
		}

		if _, e := newMnTimRelate(goutil.FNVHash64(cid)).Update(a.ID, bson.M{"status": int64(status)}); e == nil {
			t := TimeNano()
			tri1 := newMnTimRoster(uuid1)
			tri1.UUID = int64(uuid1)
			tri1.Unikid = binary(util.UnikIdByNode(fnode, tnode, domain))
			tri1.Tuuid = int64(uuid2)
			tri1.Timeseries = t
			tri1.Create()

			tri2 := newMnTimRoster(uuid2)
			tri2.UUID = int64(uuid2)
			tri2.Unikid = binary(util.UnikIdByNode(tnode, fnode, domain))
			tri2.Tuuid = int64(uuid1)
			tri2.Timeseries = t
			tri2.Create()
		}
	} else {
		status = 0x10
		if uuid1 < uuid2 {
			status = 0x01
		}
		mnt := newMnTimRelate(goutil.FNVHash64(cid))
		mnt.UUID = binary(cid)
		mnt.Status = int32(status)
		mnt.Timeseries = TimeNano()
		if _, e := mnt.Create(); e != nil {
			err = errs.ERR_DATABASE
		}
	}
	return
}

func (h *mongoHandle) Rmroster(fnode, tnode string, domain *string) (mustTell bool, ok bool) {
	cid := util.ChatIdByNode(fnode, tnode, domain)
	uuid1, uuid2 := util.NodeToUUID(fnode), util.NodeToUUID(tnode)
	lock := strlock.Lock(string(cid))
	defer lock.Unlock()

	newMnTimRoster(uuid1).Delete(bson.M{"unikid": binary(util.UnikIdByNode(fnode, tnode, domain))})
	newMnTimRoster(uuid2).Delete(bson.M{"unikid": binary(util.UnikIdByNode(tnode, fnode, domain))})

	ukid := util.UnikId(uuid1, uuid2)
	newMnTimBlock(uuid1).Delete(bson.M{"unikid": binary(ukid)})

	tl := newMnTimRelate(goutil.FNVHash64(cid))
	if a, _ := tl.Get(bson.M{"uuid": binary(cid)}, bson.M{"_id": 1, "status": 1}); a != nil {
		if uuid1 > uuid2 {
			if aStatus := uint8(a.Status) & 0x0f; aStatus == 0x02 {
				newMnTimRelate(goutil.FNVHash64(cid)).Update(a.ID, bson.M{"status": int64(a.Status)})
			} else {
				mustTell = true
				newMnTimRelate(goutil.FNVHash64(cid)).Delete(a.ID)
			}
		} else {
			if aStatus := uint8(a.Status) & 0xf0; aStatus == 0x20 {
				newMnTimRelate(goutil.FNVHash64(cid)).Update(a.ID, bson.M{"status": int32(aStatus)})
			} else {
				mustTell = true
				newMnTimRelate(goutil.FNVHash64(cid)).Delete(a.ID)
			}
		}
		ok = true
	}
	return
}

func (h *mongoHandle) Blockroster(fnode, tnode string, domain *string) (mustTell bool, ok bool) {
	cid := util.ChatIdByNode(fnode, tnode, domain)
	uuid1, uuid2 := util.NodeToUUID(fnode), util.NodeToUUID(tnode)
	lock := strlock.Lock(string(cid))
	defer lock.Unlock()

	newMnTimRoster(uuid1).Delete(bson.M{"unikid": binary(util.UnikIdByNode(fnode, tnode, domain))})
	newMnTimRoster(uuid2).Delete(bson.M{"unikid": binary(util.UnikIdByNode(tnode, fnode, domain))})

	ukid := util.UnikId(uuid1, uuid2)
	tb := newMnTimBlock(uuid1)
	if a, _ := tb.Get(bson.M{"unikid": binary(ukid)}, bson.M{"_id": 1}); a == nil {
		tbi := newMnTimBlock(uuid1)
		tbi.UUID = int64(uuid1)
		tbi.Tuuid = int64(uuid2)
		tbi.Unikid = binary(ukid)
		tbi.Timeseries = TimeNano()
		tbi.Create()
	}

	tl := newMnTimRelate(goutil.FNVHash64(cid))
	if a, _ := tl.Get(bson.M{"uuid": binary(cid)}, bson.M{"_id": 1, "status": 1}); a != nil {
		stat := uint8(a.Status)
		var aStatus uint8
		if uuid1 > uuid2 {
			aStatus = 0x20 | (uint8(a.Status) & 0x0f)
		} else {
			aStatus = (uint8(a.Status) & 0xf0) | 0x02
		}
		if stat != aStatus {

			tlu := newMnTimRelate(goutil.FNVHash64(cid))
			if _, e := tlu.Update(a.ID, bson.M{"status": int32(aStatus)}); e != nil {
				mustTell = true
			}
		}
		ok = true
	} else {
		status := uint8(0x20)
		if uuid1 < uuid2 {
			status = 0x02
		}

		tli := newMnTimRelate(goutil.FNVHash64(cid))
		tli.Status = int32(status)
		tli.UUID = binary(cid)
		tli.Timeseries = TimeNano()
		if _, e := tli.Create(); e == nil {
			ok, mustTell = true, true
		}
	}
	return
}

func (h *mongoHandle) GroupGtype(groupnode string, domain *string) (gtype int8, err errs.ERROR) {
	if guuid := util.NodeToUUID(groupnode); guuid > 0 {
		if g, _ := newMnTimGroup(guuid).Get(bson.M{"uuid": int64(guuid)}, bson.M{"gtype": 1, "status": 1}); g != nil {
			if sys.TIMTYPE(g.Status) == sys.GROUP_STATUS_CANCELLED {
				return 0, errs.ERR_CANCEL
			}
			gtype = int8(g.Gtype)
		} else {
			return 0, errs.ERR_PARAMS
		}
	} else {
		err = errs.ERR_NOEXIST
	}
	return
}

func (h *mongoHandle) GroupManagers(gnode string, domain *string) (s []string, err errs.ERROR) {
	if guuid := util.NodeToUUID(gnode); guuid > 0 {

		if g, _ := newMnTimGroup(guuid).Get(bson.M{"uuid": int64(guuid)}, bson.M{"status": 1, "rbean": 1}); g != nil {
			if sys.TIMTYPE(g.Status) == sys.GROUP_STATUS_CANCELLED {
				return nil, errs.ERR_CANCEL
			}
			mm := make(map[string]struct{})
			if tr, _ := goutil.TDecode(util.Mask(g.Rbean.Data), &stub.TimRoomBean{}); tr != nil {
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

func (h *mongoHandle) Newgroup(fnode, groupname string, gtype sys.TIMTYPE, domain *string) (gnode string, err errs.ERROR) {
	fuuid := util.NodeToUUID(fnode)
	if fuuid == 0 || !checkuseruuidMn(fuuid) {
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
	tg := newMnTimGroup(guuid)
	tg.UUID = int64(guuid)
	tg.Rbean = binary(util.Mask(goutil.TEncode(ubean)))
	tg.Status = int32(sys.GROUP_STATUS_ALIVE)
	tg.Gtype = int32(gtype)
	tg.Timeseries = ctime
	tg.Createtime = ctime
	if _, err := tg.Create(); err != nil {
		return "", errs.ERR_DATABASE
	}

	gnode = util.UUIDToNode(guuid)
	rid := util.RelateIdForGroup(gnode, fnode, domain)

	tr := newMnTimRelate(goutil.FNVHash64(rid))
	tr.UUID = binary(rid)
	tr.Status = int32(0x11)
	tr.Timeseries = ctime
	if _, e := tr.Create(); e != nil {
		return "", errs.ERR_DATABASE
	}

	tu := newMnTimMucroster(guuid)
	tu.UUID = int64(guuid)
	tu.Timeseries = ctime
	tu.Tuuid = int64(fuuid)
	tu.Unikid = binary(util.UnikIdByNode(gnode, fnode, domain))
	tu.Create()

	tu = newMnTimMucroster(fuuid)
	tu.UUID = int64(fuuid)
	tu.Timeseries = ctime
	tu.Tuuid = int64(guuid)
	tu.Unikid = binary(util.UnikIdByNode(fnode, gnode, domain))
	tu.Create()

	return
}

func (h *mongoHandle) Addgroup(groupnode, fromnode string, domain *string) (err errs.ERROR) {
	if guuid := util.NodeToUUID(groupnode); guuid > 0 {
		rid := util.RelateIdForGroup(groupnode, fromnode, domain)

		if g, _ := newMnTimGroup(guuid).Get(bson.M{"uuid": int64(guuid)}, bson.M{"gtype": 1, "status": 1}); g != nil {
			if sys.TIMTYPE(g.Status) == sys.GROUP_STATUS_CANCELLED {
				return errs.ERR_CANCEL
			}
			lock := strlock.Lock(string(rid))
			defer lock.Unlock()

			if a, _ := newMnTimRelate(goutil.FNVHash64(rid)).Get(bson.M{"uuid": binary(rid)}, bson.M{"_id": 1, "status": 1}); a != nil {
				if uint8(a.Status)&0xf0 == 0x20 {
					return errs.ERR_BLOCK
				}
				if uint8(a.Status) == 0x11 {
					return errs.ERR_HASEXIST
				}
				if sys.TIMTYPE(g.Gtype) == sys.GROUP_OPEN && uint8(a.Status) != 0x11 {
					newMnTimRelate(goutil.FNVHash64(rid)).Update(a.ID, bson.M{"status": int32(0x11)})

					ctime := TimeNano()
					fuuid := util.NodeToUUID(fromnode)
					tu := newMnTimMucroster(guuid)
					tu.Unikid = binary(util.UnikIdByNode(groupnode, fromnode, domain))
					tu.Timeseries = ctime
					tu.Tuuid = int64(fuuid)
					tu.UUID = int64(guuid)
					tu.Create()

					tu = newMnTimMucroster(fuuid)
					tu.Unikid = binary(util.UnikIdByNode(fromnode, groupnode, domain))
					tu.Timeseries = ctime
					tu.Tuuid = int64(guuid)
					tu.UUID = int64(fuuid)
					tu.Create()

				} else if sys.TIMTYPE(g.Gtype) == sys.GROUP_PRIVATE && a.Status != 0x01 {
					newMnTimRelate(goutil.FNVHash64(rid)).Update(a.ID, bson.M{"status": int32(0x01)})
				}
				return
			} else {
				if sys.TIMTYPE(g.Gtype) == sys.GROUP_PRIVATE {
					tr := newMnTimRelate(goutil.FNVHash64(rid))
					tr.UUID = binary(rid)
					tr.Timeseries = TimeNano()
					tr.Status = int32(0x01)
					if _, err := tr.Create(); err != nil {
						return errs.ERR_DATABASE
					}
				} else if sys.TIMTYPE(g.Gtype) == sys.GROUP_OPEN {
					tr := newMnTimRelate(goutil.FNVHash64(rid))
					tr.UUID = binary(rid)
					tr.Timeseries = TimeNano()
					tr.Status = int32(0x11)
					if _, err := tr.Create(); err != nil {
						return errs.ERR_DATABASE
					}

					ctime := TimeNano()
					fuuid := util.NodeToUUID(fromnode)
					tu1 := newMnTimMucroster(guuid)
					tu1.Unikid = binary(util.UnikIdByNode(groupnode, fromnode, domain))
					tu1.Timeseries = ctime
					tu1.Tuuid = int64(fuuid)
					tu1.UUID = int64(guuid)
					tu1.Create()

					tu1 = newMnTimMucroster(fuuid)
					tu1.Unikid = binary(util.UnikIdByNode(fromnode, groupnode, domain))
					tu1.Timeseries = ctime
					tu1.Tuuid = int64(guuid)
					tu1.UUID = int64(fuuid)
					tu1.Create()

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

func (h *mongoHandle) Pullgroup(groupnode, fromnode, tonode string, domain *string) (isReq bool, err errs.ERROR) {
	if guuid := util.NodeToUUID(groupnode); guuid > 0 {

		tg := newMnTimGroup(guuid)
		if g, _ := tg.Get(bson.M{"uuid": int64(guuid)}, bson.M{"gtype": 1, "status": 1, "rbean": 1}); g != nil {
			if sys.TIMTYPE(g.Status) == sys.GROUP_STATUS_CANCELLED {
				return isReq, errs.ERR_CANCEL
			}
			if sys.TIMTYPE(g.Gtype) == sys.GROUP_PRIVATE {
				if tr, _ := goutil.TDecode(util.Mask(g.Rbean.Data), &stub.TimRoomBean{}); tr != nil {
					if tr.GetFounder() != fromnode && !util.ContainStrings(tr.Managers, fromnode) {
						err = errs.ERR_PERM_DENIED
						return
					}
				}
			}
			rid := util.RelateIdForGroup(groupnode, tonode, domain)
			lock := strlock.Lock(string(rid))
			defer lock.Unlock()

			tr := newMnTimRelate(goutil.FNVHash64(rid))
			if a, _ := tr.Get(bson.M{"uuid": binary(rid)}, bson.M{"_id": 1, "status": 1}); a != nil {
				if uint8(a.Status)&0x0f == 0x02 {
					return isReq, errs.ERR_BLOCK
				}
				if uint8(a.Status) == 0x11 {
					return isReq, errs.ERR_HASEXIST
				}
				isReq = uint8(a.Status) == 0x01
				if uint8(a.Status) != 0x11 {
					newMnTimRelate(goutil.FNVHash64(rid)).Update(a.ID, bson.M{"status": int32(0x11)})
				}
			} else {

				tr := newMnTimRelate(goutil.FNVHash64(rid))
				tr.UUID = binary(rid)
				tr.Status = int32(0x11)
				tr.Timeseries = TimeNano()
				if _, e := tr.Create(); e != nil {
					return isReq, errs.ERR_DATABASE
				}
			}
			ctime := TimeNano()
			fuuid := util.NodeToUUID(tonode)

			tu := newMnTimMucroster(guuid)
			tu.Unikid = binary(util.UnikIdByNode(groupnode, tonode, domain))
			tu.Timeseries = ctime
			tu.Tuuid = int64(fuuid)
			tu.UUID = int64(guuid)
			tu.Create()

			tu = newMnTimMucroster(fuuid)
			tu.Unikid = binary(util.UnikIdByNode(tonode, groupnode, domain))
			tu.Timeseries = ctime
			tu.Tuuid = int64(guuid)
			tu.UUID = int64(fuuid)
			tu.Create()

		} else {
			return isReq, errs.ERR_NOEXIST
		}
	} else {
		err = errs.ERR_NOEXIST
	}
	return
}

func (h *mongoHandle) Nopassgroup(groupnode, fromnode, tonode string, domain *string) (err errs.ERROR) {
	if guuid := util.NodeToUUID(groupnode); guuid > 0 {
		//tg := newTimgroup(guuid)
		//tg.Where(tg.UUID.EQ(int64(guuid)))
		tg := newMnTimGroup(guuid)
		if g, _ := tg.Get(bson.M{"uuid": int64(guuid)}, bson.M{"status": 1, "rbean": 1}); g != nil {
			if sys.TIMTYPE(g.Status) == sys.GROUP_STATUS_CANCELLED {
				return errs.ERR_CANCEL
			}
			if tr, _ := goutil.TDecode(util.Mask(g.Rbean.Data), &stub.TimRoomBean{}); tr != nil {
				if tr.GetFounder() == fromnode || util.ContainStrings(tr.Managers, fromnode) {
					rid := util.RelateIdForGroup(groupnode, tonode, domain)
					lock := strlock.Lock(string(rid))
					defer lock.Unlock()
					//tr := newTimrelate(goutil.FNVHash64(rid))
					//tr.Where(tr.UUID.EQ(rid))
					if a, _ := newMnTimRelate(goutil.FNVHash64(rid)).Get(bson.M{"uuid": binary(rid)}, bson.M{"_id": 1, "status": 1}); a != nil {
						if a.Status == 0x01 {
							//td := newTimrelate(goutil.FNVHash64(rid))
							//td.Where(td.ID.EQ(a.GetId())).Delete()
							newMnTimRelate(goutil.FNVHash64(rid)).Delete(a.ID)
							return
						}
						if uint8(a.Status)|0xf0 != 0 {
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

func (h *mongoHandle) Kickgroup(groupnode, fromnode, tonode string, domain *string) (err errs.ERROR) {
	if guuid := util.NodeToUUID(groupnode); guuid > 0 {
		if fromnode == tonode {
			return errs.ERR_PERM_DENIED
		}
		if g, _ := newMnTimGroup(guuid).Get(bson.M{"uuid": int64(guuid)}, bson.M{"_id": 1, "status": 1, "rbean": 1}); g != nil && sys.TIMTYPE(g.Status) != sys.GROUP_STATUS_CANCELLED {
			if tr, _ := goutil.TDecode(util.Mask(g.Rbean.Data), &stub.TimRoomBean{}); tr != nil {
				if tr.GetFounder() == fromnode || util.ContainStrings(tr.Managers, fromnode) {
					if tr.GetFounder() != fromnode && util.ContainStrings(tr.Managers, tonode) {
						return errs.ERR_PERM_DENIED
					}
					if util.ContainStrings(tr.Managers, tonode) {
						tr.Managers = util.ArraySub(tr.Managers, []string{tonode})
						newMnTimGroup(guuid).Update(g.ID, bson.M{"rbean": binary(util.Mask(goutil.TEncode(tr)))})
					}
					rid := util.RelateIdForGroup(groupnode, tonode, domain)
					lock := strlock.Lock(string(rid))
					defer lock.Unlock()

					//tmu := newTimmucroster(guuid)
					//tmu.Where(tmu.UNIKID.EQ(util.UnikIdByNode(groupnode, tonode, domain))).Delete()
					newMnTimMucroster(guuid).DeleteOption(bson.M{"unikid": binary(util.UnikIdByNode(groupnode, tonode, domain))})

					tuuid := util.NodeToUUID(tonode)
					//tmu = newTimmucroster(tuuid)
					//tmu.Where(tmu.UNIKID.EQ(util.UnikIdByNode(tonode, groupnode, domain))).Delete()
					newMnTimMucroster(tuuid).DeleteOption(bson.M{"unikid": binary(util.UnikIdByNode(tonode, groupnode, domain))})

					ukid := util.UnikId(guuid, tuuid)
					//tb := newTimblockroom(guuid)
					//tb.Where(tb.UNIKID.EQ(ukid)).Delete()
					newMnTimBlockroom(guuid).DeleteOption(bson.M{"unikid": binary(ukid)})

					//tr := newTimrelate(goutil.FNVHash64(rid))
					//tr.Where(tr.UUID.EQ(rid))
					if a, _ := newMnTimRelate(goutil.FNVHash64(rid)).Get(bson.M{"uuid": binary(rid)}, bson.M{"_id": 1, "status": 1}); a != nil {
						if int8(a.Status)|0x0f == 0x02 {
							//tru := newTimrelate(goutil.FNVHash64(rid))
							//tru.SetStatus(0x02).Where(tru.ID.EQ(a.GetId())).Update()
							tru := newMnTimRelate(goutil.FNVHash64(rid))
							tru.Update(a.ID, bson.M{"status": int32(0x02)})
						} else {
							//tdd := newTimrelate(goutil.FNVHash64(rid))
							//if _, e := tdd.Where(tdd.ID.EQ(a.GetId())).Delete(); e != nil {
							//	err = errs.ERR_DATABASE
							//}
							if _, e := newMnTimRelate(goutil.FNVHash64(rid)).Delete(a.ID); e != nil {
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

func (h *mongoHandle) Leavegroup(groupnode, fromnode string, domain *string) (err errs.ERROR) {
	guuid, tuuid := util.NodeToUUID(groupnode), util.NodeToUUID(fromnode)
	if guuid > 0 {
		if err = func() (err errs.ERROR) {
			lock := numlock.Lock(int64(guuid))
			defer lock.Unlock()

			if g, _ := newMnTimGroup(guuid).Get(bson.M{"uuid": int64(guuid)}, bson.M{"_id": 1, "rbean": 1, "status": 1}); g != nil {
				if sys.TIMTYPE(g.Status) == sys.GROUP_STATUS_CANCELLED {
					return errs.ERR_CANCEL
				}
				if tr, _ := goutil.TDecode(util.Mask(g.Rbean.Data), &stub.TimRoomBean{}); tr != nil {
					if tr.GetFounder() == fromnode {
						return errs.ERR_PERM_DENIED
					}
					if util.ContainStrings(tr.Managers, fromnode) {
						tr.Managers = util.ArraySub(tr.Managers, []string{fromnode})
						//tgu := newTimgroup(guuid)
						//tgu.SetRbean(util.Mask(goutil.TEncode(tr))).Where(tgu.ID.EQ(g.GetId())).Update()
						newMnTimGroup(guuid).Update(g.ID, bson.M{"rbean": binary(util.Mask(goutil.TEncode(tr)))})

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

	//tmu := newTimmucroster(guuid)
	//tmu.Where(tmu.UNIKID.EQ(util.UnikIdByNode(groupnode, fromnode, domain))).Delete()
	newMnTimMucroster(guuid).DeleteOption(bson.M{"unikid": binary(util.UnikIdByNode(groupnode, fromnode, domain))})

	//tmu = newTimmucroster(tuuid)
	//tmu.Where(tmu.UNIKID.EQ(util.UnikIdByNode(fromnode, groupnode, domain))).Delete()
	newMnTimMucroster(tuuid).DeleteOption(bson.M{"unikid": binary(util.UnikIdByNode(fromnode, groupnode, domain))})

	ukid := util.UnikId(tuuid, guuid)
	//tbr := newTimblockroom(tuuid)
	//tbr.Where(tbr.UNIKID.EQ(ukid)).Delete()
	newMnTimBlockroom(tuuid).DeleteOption(bson.M{"unikid": binary(ukid)})

	lock := strlock.Lock(string(rid))
	defer lock.Unlock()

	//tr := newTimrelate(goutil.FNVHash64(rid))
	//tr.Where(tr.UUID.EQ(rid))
	if a, _ := newMnTimRelate(goutil.FNVHash64(rid)).Get(bson.M{"uuid": binary(rid)}, bson.M{"_id": 1, "status": 1}); a != nil {
		if uint8(a.Status)&0xf0 == 0x20 {
			//tru := newTimrelate(goutil.FNVHash64(rid))
			//tru.SetStatus(0x20).Where(tru.ID.EQ(a.GetId())).Update()
			newMnTimRelate(goutil.FNVHash64(rid)).Update(a.ID, bson.M{"status": int32(0x20)})
		} else {
			//trd := newTimrelate(goutil.FNVHash64(rid))
			if _, e := newMnTimRelate(goutil.FNVHash64(rid)).Delete(a.ID); e != nil {
				err = errs.ERR_DATABASE
			}
		}
	} else {
		err = errs.ERR_NOEXIST
	}
	return
}

func (h *mongoHandle) Cancelgroup(groupnode, fromnode string, domain *string) (err errs.ERROR) {
	if guuid := util.NodeToUUID(groupnode); guuid > 0 {
		lock := numlock.Lock(int64(guuid))
		defer lock.Unlock()
		//tg := newTimgroup(guuid)
		//tg.Where(tg.UUID.EQ(int64(guuid)))
		if g, _ := newMnTimGroup(guuid).Get(bson.M{"uuid": int64(guuid)}, bson.M{"_id": 1, "status": 1, "rbean": 1}); g != nil {
			if sys.TIMTYPE(g.Status) == sys.GROUP_STATUS_CANCELLED {
				return errs.ERR_CANCEL
			}
			if tr, _ := goutil.TDecode(util.Mask(g.Rbean.Data), &stub.TimRoomBean{}); tr != nil {
				if tr.GetFounder() == fromnode {
					//tu := newTimmucroster(guuid)
					//tu.Where(tu.UUID.EQ(int64(guuid)))
					if tus, _ := newMnTimMucroster(guuid).List(bson.M{"uuid": int64(guuid)}, bson.M{"tuuid": 1}); len(tus) == 1 {
						tuuid := util.NodeToUUID(fromnode)
						if uint64(tus[0].Tuuid) != tuuid {
							return errs.ERR_PERM_DENIED
						}
						ids := make([]any, len(tus))
						for i := range tus {
							ids[i] = tus[i].ID
						}
						if len(ids) > 0 {
							//tmu := newTimmucroster(guuid)
							//tmu.Where(tmu.ID.IN(ids...)).Delete()
							newMnTimMucroster(guuid).DeleteOption(bson.M{"_id": bson.M{"$in": ids}})
						}
						//tgu := newTimgroup(guuid)
						//if _, e := tgu.SetStatus(int64(sys.GROUP_STATUS_CANCELLED)).Where(tgu.ID.EQ(g.GetId())).Update(); e != nil {
						//	return errs.ERR_DATABASE
						//}
						if _, e := newMnTimGroup(guuid).Update(g.ID, bson.M{"status": int32(sys.GROUP_STATUS_CANCELLED)}); e != nil {
							return errs.ERR_DATABASE
						}

						//tmu := newTimmucroster(guuid)
						//tmu.Where(tmu.UNIKID.EQ(util.UnikIdByNode(groupnode, fromnode, domain))).Delete()
						newMnTimMucroster(guuid).DeleteOption(bson.M{"unikid": binary(util.UnikIdByNode(groupnode, fromnode, domain))})

						//tmu = newTimmucroster(tuuid)
						//tmu.Where(tmu.UNIKID.EQ(util.UnikIdByNode(fromnode, groupnode, domain))).Delete()
						newMnTimMucroster(tuuid).DeleteOption(bson.M{"unikid": binary(util.UnikIdByNode(fromnode, groupnode, domain))})
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

func (h *mongoHandle) Blockgroup(groupnode, fromnode string, domain *string) (err errs.ERROR) {
	rid := util.RelateIdForGroup(groupnode, fromnode, domain)
	guuid, tuuid := util.NodeToUUID(groupnode), util.NodeToUUID(fromnode)

	//tmu := newTimmucroster(guuid)
	//tmu.Where(tmu.UNIKID.EQ(util.UnikIdByNode(groupnode, fromnode, domain))).Delete()
	newMnTimMucroster(guuid).DeleteOption(bson.M{"unikid": binary(util.UnikIdByNode(groupnode, fromnode, domain))})

	//tmu = newTimmucroster(tuuid)
	//tmu.Where(tmu.UNIKID.EQ(util.UnikIdByNode(fromnode, groupnode, domain))).Delete()
	newMnTimMucroster(tuuid).DeleteOption(bson.M{"unikid": binary(util.UnikIdByNode(fromnode, groupnode, domain))})

	lock := strlock.Lock(string(rid))
	defer lock.Unlock()

	ukid := util.UnikId(tuuid, guuid)

	//newTimblockroom(tuuid).SetTuuid(int64(guuid)).SetUuid(int64(tuuid)).SetUnikid(ukid).SetTimeseries(TimeNano()).Insert()
	tb := newMnTimBlockroom(tuuid)
	tb.Tuuid = int64(guuid)
	tb.UUID = int64(tuuid)
	tb.Unikid = binary(ukid)
	tb.Timeseries = TimeNano()
	tb.Create()

	uuid := goutil.FNVHash64(rid)
	//tr := newTimrelate(uuid)
	//tr.Where(tr.UUID.EQ(rid))
	if a, _ := newMnTimRelate(uuid).Get(bson.M{"uuid": binary(rid)}, bson.M{"_id": 1, "status": 1}); a != nil {
		if state := uint8(a.Status); state&0x0f != 0x02 {
			//tru := newTimrelate(uuid)
			//if _, e := tru.SetStatus(int64(state&0xf0 | 0x02)).Where(tru.ID.EQ(a.GetId())).Update(); e != nil {
			//	err = errs.ERR_DATABASE
			//}
			if _, e := newMnTimRelate(uuid).Update(a.ID, bson.M{"status": int32(state&0xf0 | 0x02)}); e != nil {
				err = errs.ERR_DATABASE
			}
		}
	} else {
		//newTimrelate(uuid).SetStatus(0x02).SetUuid(rid).SetTimeseries(TimeNano()).Insert()
		tl := newMnTimRelate(uuid)
		tl.Status = int32(0x02)
		tl.UUID = binary(rid)
		tl.Timeseries = TimeNano()
		tl.Create()
	}
	return
}

func (h *mongoHandle) Blockgroupmember(groupnode, fromnode, tonode string, domain *string) (err errs.ERROR) {
	if err = h.checkAdmin(groupnode, fromnode, tonode); err != nil {
		return
	}
	rid := util.RelateIdForGroup(groupnode, tonode, domain)
	guuid, tuuid := util.NodeToUUID(groupnode), util.NodeToUUID(tonode)

	//tmu := newTimmucroster(guuid)
	//tmu.Where(tmu.UNIKID.EQ(util.UnikIdByNode(groupnode, tonode, domain))).Delete()
	newMnTimMucroster(guuid).DeleteOption(bson.M{"unikid": binary(util.UnikIdByNode(groupnode, tonode, domain))})

	//tmu = newTimmucroster(tuuid)
	//tmu.Where(tmu.UNIKID.EQ(util.UnikIdByNode(tonode, groupnode, domain))).Delete()
	newMnTimMucroster(tuuid).DeleteOption(bson.M{"unikid": binary(util.UnikIdByNode(tonode, groupnode, domain))})

	lock := strlock.Lock(string(rid))
	defer lock.Unlock()

	ukid := util.UnikId(guuid, tuuid)
	//tbr := newTimblockroom(guuid)
	//tbr.Where(tbr.UNIKID.EQ(ukid))
	if a, _ := newMnTimBlockroom(guuid).Get(bson.M{"unikid": binary(ukid)}, bson.M{"_id": 1}); a == nil {
		//newTimblockroom(guuid).SetUnikid(ukid).SetUuid(int64(guuid)).SetTuuid(int64(tuuid)).SetTimeseries(TimeNano()).Insert()
		tl := newMnTimBlockroom(guuid)
		tl.Unikid = binary(ukid)
		tl.UUID = int64(guuid)
		tl.Tuuid = int64(tuuid)
		tl.Timeseries = TimeNano()
		tl.Create()
	}

	//tr := newTimrelate(goutil.FNVHash64(rid))
	//tr.Where(tr.UUID.EQ(rid))
	if a, _ := newMnTimRelate(goutil.FNVHash64(rid)).Get(bson.M{"uuid": binary(rid)}, bson.M{"_id": 1, "status": 1}); a != nil {
		if state := uint8(a.Status); state&0xf0 != 0x20 {
			//tru := newTimrelate(goutil.FNVHash64(rid))
			//if _, e := tru.SetStatus(int64(state&0x0f | 0x20)).Where(tru.ID.EQ(a.GetId())).Update(); e != nil {
			//	err = errs.ERR_DATABASE
			//}
			if _, e := newMnTimRelate(goutil.FNVHash64(rid)).Update(a.ID, bson.M{"status": int32(state&0x0f | 0x20)}); e != nil {
				err = errs.ERR_DATABASE
			}
		}
	} else {
		//newTimrelate(goutil.FNVHash64(rid)).SetStatus(0x20).SetUuid(rid).SetTimeseries(TimeNano()).Insert()
		tl := newMnTimRelate(goutil.FNVHash64(rid))
		tl.Status = int32(0x20)
		tl.UUID = binary(rid)
		tl.Timeseries = TimeNano()
		tl.Create()
	}
	return
}

func (h *mongoHandle) checkAdmin(groupnode, fromnode, tonode string) (err errs.ERROR) {
	if guuid := util.NodeToUUID(groupnode); guuid > 0 {
		if fromnode == tonode {
			return errs.ERR_PERM_DENIED
		}
		//tg := newTimgroup(guuid)
		//tg.Where(tg.UUID.EQ(int64(guuid)))
		if g, _ := newMnTimGroup(guuid).Get(bson.M{"uuid": int64(guuid)}, bson.M{"status": 1, "rbean": 1}); g != nil {
			if sys.TIMTYPE(g.Status) == sys.GROUP_STATUS_CANCELLED {
				return errs.ERR_CANCEL
			}
			if tr, _ := goutil.TDecode(util.Mask(g.Rbean.Data), &stub.TimRoomBean{}); tr != nil {
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

func (h *mongoHandle) ModifyUserInfo(node string, tu *stub.TimUserBean) (err errs.ERROR) {
	uuid := util.NodeToUUID(node)
	if uuid == 0 {
		return errs.ERR_ACCOUNT
	}
	if a, _ := newMnTimUser(uuid).Get(bson.M{"uuid": int64(uuid)}, bson.M{"_id": 1, "ubean": 1}); a != nil {
		if a.Ubean.Data != nil {
			if ub, _ := goutil.TDecode(util.Mask(a.Ubean.Data), &stub.TimUserBean{}); ub != nil {
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
		newMnTimUser(uuid).Update(a.ID, bson.M{"ubean": binary(util.Mask(goutil.TEncode(tu)))})
	} else {
		err = errs.ERR_NOEXIST
	}
	return
}
func (h *mongoHandle) GetUserInfo(nodes []string) (m map[string]*stub.TimUserBean, err errs.ERROR) {
	if len(nodes) > 0 {
		m = make(map[string]*stub.TimUserBean, 0)
		for _, node := range nodes {
			uuid := util.NodeToUUID(node)
			if a, _ := newMnTimUser(uuid).Get(bson.M{"uuid": int64(uuid)}, bson.M{"ubean": 1, "createtime": 1}); a != nil {
				if a.Ubean.Data != nil {
					if tub, _ := goutil.TDecode(util.Mask(a.Ubean.Data), &stub.TimUserBean{}); tub != nil {
						tub.Createtime = &a.Createtime
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

func (h *mongoHandle) ModifygroupInfo(node, fnode string, tu *stub.TimRoomBean, admin bool) (err errs.ERROR) {
	if tu == nil {
		return
	}
	if guuid := util.NodeToUUID(node); guuid > 0 {
		if g, _ := newMnTimGroup(guuid).Get(bson.M{"uuid": int64(guuid)}, bson.M{"_id": 1, "status": 1, "rbean": 1}); g != nil {
			if sys.TIMTYPE(g.Status) == sys.GROUP_STATUS_CANCELLED {
				return errs.ERR_CANCEL
			}
			if tr, _ := goutil.TDecode(util.Mask(g.Rbean.Data), &stub.TimRoomBean{}); tr != nil {
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
				newMnTimGroup(guuid).Update(g.ID, bson.M{"rbean": binary(util.Mask(goutil.TEncode(tr)))})
			}
		} else {
			return errs.ERR_NOEXIST
		}
	} else {
		return errs.ERR_ACCOUNT
	}
	return
}

func (h *mongoHandle) GetGroupInfo(nodes []string) (m map[string]*stub.TimRoomBean, err errs.ERROR) {
	if len(nodes) > 0 {
		m = make(map[string]*stub.TimRoomBean, 0)
		for _, node := range nodes {
			if guuid := util.NodeToUUID(node); guuid > 0 {
				if g, _ := newMnTimGroup(guuid).Get(bson.M{"uuid": int64(guuid)}, bson.M{"status": 1, "rbean": 1}); g != nil && sys.TIMTYPE(g.Status) != sys.GROUP_STATUS_CANCELLED {
					if tr, _ := goutil.TDecode(util.Mask(g.Rbean.Data), &stub.TimRoomBean{}); tr != nil {
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

func (h *mongoHandle) TimAdminAuth(account, password, domain string) bool {
	if t, _ := newMnTimDomain(0).Get(bson.M{"adminaccount": account, "timdomain": domain}, bson.M{"adminpassword": 1}); t != nil {
		return strings.EqualFold(t.Adminpassword, goutil.Md5Str(password))
	}
	return false
}
