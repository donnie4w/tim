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

	. "github.com/donnie4w/gofer/util"
	. "github.com/donnie4w/tim/stub"
	"github.com/donnie4w/tim/sys"
	"github.com/donnie4w/tim/util"
)

func (this *tldbhandler) Roster(node string) (_r []string) {
	_r = make([]string, 0)
	if tr, err := SelectAllByIdx[timroster]("UUID", util.NodeToUUID(node)); err == nil {
		for _, a := range tr {
			_r = append(_r, util.UUIDToNode(a.TUUID))
		}
	}
	return
}

func (this *tldbhandler) Blockrosterlist(node string) (_r []string) {
	_r = make([]string, 0)
	if tr, err := SelectAllByIdx[timblock]("UUID", util.NodeToUUID(node)); err == nil {
		for _, a := range tr {
			_r = append(_r, util.UUIDToNode(a.TUUID))
		}
	}
	return
}

func (this *tldbhandler) Blockroomlist(node string) (_r []string) {
	_r = make([]string, 0)
	if tr, err := SelectAllByIdx[timblockroom]("UUID", util.NodeToUUID(node)); err == nil {
		for _, a := range tr {
			_r = append(_r, util.UUIDToNode(a.TUUID))
		}
	}
	return
}

func (this *tldbhandler) Blockroommemberlist(node string, fnode string) (_r []string) {
	if checkAdmin(node, fnode, "") != nil {
		return
	}
	_r = make([]string, 0)
	if tr, err := SelectAllByIdx[timblockroom]("UUID", util.NodeToUUID(node)); err == nil {
		for _, a := range tr {
			_r = append(_r, util.UUIDToNode(a.TUUID))
		}
	}
	return
}

func (this *tldbhandler) UserGroup(node string, domain *string) (_r []string) {
	if tr, _ := SelectAllByIdx[timmucroster]("UUID", util.NodeToUUID(node)); tr != nil && len(tr) > 0 {
		_r = make([]string, 0)
		for _, a := range tr {
			_r = append(_r, util.UUIDToNode(a.TUUID))
		}
	}
	return
}

func (this *tldbhandler) GroupRoster(groupnode string) (_r []string) {
	if tr, _ := SelectAllByIdx[timmucroster]("UUID", util.NodeToUUID(groupnode)); tr != nil && len(tr) > 0 {
		_r = make([]string, 0)
		for _, a := range tr {
			_r = append(_r, util.UUIDToNode(a.TUUID))
		}
	}
	return
}

func (this *tldbhandler) Addroster(fnode, tnode string, domain *string) (status int8, err sys.ERROR) {
	cid := util.ChatIdByNode(fnode, tnode, domain)
	uuid1, uuid2 := util.NodeToUUID(fnode), util.NodeToUUID(tnode)
	numlock.Lock(int64(cid))
	defer numlock.Unlock(int64(cid))
	if a, _ := SelectByIdx[timrelate]("UUID", cid); a != nil {
		stat := a.Status
		if uuid1 > uuid2 {
			if a.Status&0x0f == 0x02 {
				err = sys.ERR_BLOCK
				return
			}
			a.Status = 0x10 | (a.Status & 0x0f)
		} else {
			if a.Status&0xf0 == 0x20 {
				err = sys.ERR_BLOCK
				return
			}
			a.Status = (a.Status & 0xf0) | 0x01
		}
		if stat != a.Status {
			UpdateNonzero(a)
		}
		status = int8(a.Status)
		if stat != 0x11 && status == 0x11 {
			Insert(&timroster{Relate: cid, UUID: uuid1, TUUID: uuid2})
			Insert(&timroster{Relate: cid, UUID: uuid2, TUUID: uuid1})
		}
	} else {
		status = 0x10
		if uuid1 < uuid2 {
			status = 0x01
		}
		if id, _ := Insert(&timrelate{UUID: cid, Status: uint8(status)}); id == 0 {
			err = sys.ERR_DATABASE
		}
	}
	return
}

func (this *tldbhandler) Rmroster(fnode, tnode string, domain *string) (mstell bool, ok bool) {
	cid := util.ChatIdByNode(fnode, tnode, domain)
	uuid1, uuid2 := util.NodeToUUID(fnode), util.NodeToUUID(tnode)
	numlock.Lock(int64(cid))
	defer numlock.Unlock(int64(cid))
	if as, _ := SelectAllByIdxWithTid[timroster](uuid1, "Relate", cid); as != nil {
		for _, a := range as {
			Delete[timroster](uuid1, a.Id)
		}
	}
	if as, _ := SelectAllByIdxWithTid[timroster](uuid2, "Relate", cid); as != nil {
		for _, a := range as {
			Delete[timroster](uuid2, a.Id)
		}
	}

	ukid := util.UnikId(uuid1, uuid2)
	if as, _ := SelectAllByIdxWithTid[timblock](uuid1, "UnikId", ukid); as != nil {
		for _, a := range as {
			Delete[timblock](uuid1, a.Id)
		}
	}

	if a, _ := SelectByIdx[timrelate]("UUID", cid); a != nil {
		if uuid1 > uuid2 {
			if a.Status = a.Status & 0x0f; a.Status == 0x02 {
				UpdateNonzero(a)
			} else {
				mstell = true
				Delete[timrelate](a.tid(), a.Id)
			}
		} else {
			if a.Status = a.Status & 0xf0; a.Status == 0x20 {
				UpdateNonzero(a)
			} else {
				mstell = true
				Delete[timrelate](a.tid(), a.Id)
			}
		}
		ok = true
	}
	return
}

func (this *tldbhandler) Blockroster(fnode, tnode string, domain *string) (mstell bool, ok bool) {
	cid := util.ChatIdByNode(fnode, tnode, domain)
	uuid1, uuid2 := util.NodeToUUID(fnode), util.NodeToUUID(tnode)
	numlock.Lock(int64(cid))
	defer numlock.Unlock(int64(cid))

	if as, _ := SelectAllByIdxWithTid[timroster](uuid1, "Relate", cid); as != nil {
		for _, a := range as {
			Delete[timroster](uuid1, a.Id)
		}
	}
	if as, _ := SelectAllByIdxWithTid[timroster](uuid2, "Relate", cid); as != nil {
		for _, a := range as {
			Delete[timroster](uuid2, a.Id)
		}
	}

	ukid := util.UnikId(uuid1, uuid2)
	if a, _ := SelectByIdxWithTid[timblock](uuid1, "UnikId", ukid); a == nil {
		Insert(&timblock{UnikId: ukid, UUID: uuid1, TUUID: uuid2})
	}

	if a, _ := SelectByIdx[timrelate]("UUID", cid); a != nil {
		stat := a.Status
		if uuid1 > uuid2 {
			a.Status = 0x20 | (a.Status & 0x0f)
		} else {
			a.Status = (a.Status & 0xf0) | 0x02
		}
		if stat != a.Status {
			if err := UpdateNonzero(a); err == nil {
				mstell = true
			}
		}
		ok = true
	} else {
		status := uint8(0x20)
		if uuid1 < uuid2 {
			status = 0x02
		}
		if id, _ := Insert(&timrelate{UUID: cid, Status: status}); id > 0 {
			ok, mstell = true, true
		}
	}
	return
}

func (this *tldbhandler) GroupGtype(groupnode string, domain *string) (gtype int8, err sys.ERROR) {
	if guuid := util.NodeToUUID(groupnode); guuid > 0 {
		if g, _ := SelectByIdx[timgroup]("UUID", guuid); g != nil {
			if g.Status == sys.GROUP_STATUS_CANCELLED {
				return 0, sys.ERR_CANCEL
			}
			gtype = g.Gtype
		} else {
			return 0, sys.ERR_PARAMS
		}
	} else {
		err = sys.ERR_NOEXIST
	}
	return
}

func (this *tldbhandler) GroupManagers(groupnode string, domain *string) (s []string, err sys.ERROR) {
	if guuid := util.NodeToUUID(groupnode); guuid > 0 {
		if g, _ := SelectByIdx[timgroup]("UUID", guuid); g != nil {
			if g.Status == sys.GROUP_STATUS_CANCELLED {
				return nil, sys.ERR_CANCEL
			}
			if tr, _ := TDecode(util.Mask(g.RBean), &TimRoomBean{}); tr != nil {
				s = tr.Managers
				if s == nil {
					s = []string{*tr.Founder}
				}
				if !util.ContainStrings(s, *tr.Founder) {
					s = append(s, *tr.Founder)
				}
			}
		} else {
			return nil, sys.ERR_PARAMS
		}
	} else {
		err = sys.ERR_NOEXIST
	}
	return
}

func (this *tldbhandler) Newgroup(fnode, groupname string, gtype int8, domain *string) (gnode string, err sys.ERROR) {
	UUID := util.CreateUUID(string(Int64ToBytes(RandId())), domain)
	tg := &timgroup{Gtype: gtype, UUID: UUID, Createtime: time.Now().UnixNano(), Status: sys.GROUP_STATUS_ALIVE}
	ubean := &TimRoomBean{Founder: &fnode, Topic: &groupname, Createtime: &tg.Createtime, Gtype: &gtype}
	tg.RBean = util.Mask(TEncode(ubean))
	if id, _ := Insert(tg); id == 0 {
		return "", sys.ERR_DATABASE
	}
	gnode = util.UUIDToNode(UUID)
	rid := util.RelateIdForGroup(gnode, fnode, domain)
	tuuid := util.NodeToUUID(fnode)
	if id, _ := Insert(&timrelate{UUID: rid, Status: 0x11}); id == 0 {
		return "", sys.ERR_DATABASE
	}
	if id, _ := Insert(&timmucroster{Relate: rid, UUID: UUID, TUUID: tuuid}); id == 0 {
		return "", sys.ERR_DATABASE
	}
	if id, _ := Insert(&timmucroster{Relate: rid, UUID: tuuid, TUUID: UUID}); id == 0 {
		return "", sys.ERR_DATABASE
	}
	return
}

func (this *tldbhandler) Addgroup(groupnode, fromnode string, domain *string) (err sys.ERROR) {
	if guuid := util.NodeToUUID(groupnode); guuid > 0 {
		rid := util.RelateIdForGroup(groupnode, fromnode, domain)
		if g, _ := SelectByIdx[timgroup]("UUID", guuid); g != nil {
			if g.Status == sys.GROUP_STATUS_CANCELLED {
				return sys.ERR_CANCEL
			}
			numlock.Lock(int64(rid))
			defer numlock.Unlock(int64(rid))
			if a, _ := SelectByIdx[timrelate]("UUID", rid); a != nil {
				if a.Status&0xf0 == 0x20 {
					return sys.ERR_BLOCK
				}
				if a.Status == 0x11 {
					return sys.ERR_HASEXIST
				}
				if g.Gtype == sys.GROUP_OPEN && a.Status != 0x11 {
					a.Status = 0x11
					UpdateNonzero(a)
					Insert(&timmucroster{Relate: rid, UUID: guuid, TUUID: util.NodeToUUID(fromnode)})
					Insert(&timmucroster{Relate: rid, UUID: util.NodeToUUID(fromnode), TUUID: guuid})
				} else if g.Gtype == sys.GROUP_PRIVATE && a.Status != 0x01 {
					a.Status = 0x01
					UpdateNonzero(a)
				}
				return
			} else {
				if g.Gtype == sys.GROUP_PRIVATE {
					if id, _ := Insert(&timrelate{UUID: rid, Status: 0x01}); id == 0 {
						return sys.ERR_DATABASE
					}
				} else if g.Gtype == sys.GROUP_OPEN {
					if id, _ := Insert(&timrelate{UUID: rid, Status: 0x11}); id == 0 {
						return sys.ERR_DATABASE
					}
					Insert(&timmucroster{Relate: rid, UUID: guuid, TUUID: util.NodeToUUID(fromnode)})
					Insert(&timmucroster{Relate: rid, UUID: util.NodeToUUID(fromnode), TUUID: guuid})
				} else {
					return sys.ERR_AUTH
				}
			}
		}
	} else {
		err = sys.ERR_NOEXIST
	}
	return
}

func (this *tldbhandler) Pullgroup(groupnode, fromnode, tonode string, domain *string) (isReq bool, err sys.ERROR) {
	if guuid := util.NodeToUUID(groupnode); guuid > 0 {
		if g, _ := SelectByIdx[timgroup]("UUID", guuid); g != nil {
			if g.Status == sys.GROUP_STATUS_CANCELLED {
				return isReq, sys.ERR_CANCEL
			}
			if g.Gtype == sys.GROUP_PRIVATE {
				if tr, _ := TDecode(util.Mask(g.RBean), &TimRoomBean{}); tr != nil {
					if *tr.Founder != fromnode && !util.ContainStrings(tr.Managers, fromnode) {
						err = sys.ERR_AUTH
						return
					}
				}
			}
			rid := util.RelateIdForGroup(groupnode, tonode, domain)
			numlock.Lock(int64(rid))
			defer numlock.Unlock(int64(rid))
			if a, _ := SelectByIdx[timrelate]("UUID", rid); a != nil {
				if a.Status&0x0f == 0x02 {
					return isReq, sys.ERR_BLOCK
				}
				if a.Status == 0x11 {
					return isReq, sys.ERR_HASEXIST
				}
				isReq = a.Status == 0x01
				if a.Status != 0x11 {
					a.Status = 0x11
					UpdateNonzero(a)
				}
			} else if id, _ := Insert(&timrelate{UUID: rid, Status: 0x11}); id == 0 {
				return isReq, sys.ERR_DATABASE
			}
			Insert(&timmucroster{Relate: rid, UUID: guuid, TUUID: util.NodeToUUID(tonode)})
			Insert(&timmucroster{Relate: rid, UUID: util.NodeToUUID(tonode), TUUID: guuid})
		} else {
			return isReq, sys.ERR_NOEXIST
		}
	} else {
		err = sys.ERR_NOEXIST
	}
	return
}

func (this *tldbhandler) Nopassgroup(groupnode, fromnode, tonode string, domain *string) (err sys.ERROR) {
	if guuid := util.NodeToUUID(groupnode); guuid > 0 {
		if g, _ := SelectByIdx[timgroup]("UUID", guuid); g != nil {
			if g.Status == sys.GROUP_STATUS_CANCELLED {
				return sys.ERR_CANCEL
			}
			if tr, _ := TDecode(util.Mask(g.RBean), &TimRoomBean{}); tr != nil {
				if *tr.Founder == fromnode || util.ContainStrings(tr.Managers, fromnode) {
					rid := util.RelateIdForGroup(groupnode, tonode, domain)
					numlock.Lock(int64(rid))
					defer numlock.Unlock(int64(rid))
					if a, _ := SelectByIdx[timrelate]("UUID", rid); a != nil {
						if a.Status == 0x01 {
							Delete[timrelate](a.tid(), a.Id)
							return
						}
						if a.Status|0xf0 != 0 {
							return sys.ERR_EXPIREOP
						}
					} else {
						return sys.ERR_NOEXIST
					}
				} else {
					return sys.ERR_AUTH
				}
			} else {
				return sys.ERR_UNDEFINED
			}
		} else {
			return sys.ERR_NOEXIST
		}
	} else {
		err = sys.ERR_NOEXIST
	}
	return
}

func (this *tldbhandler) Kickgroup(groupnode, fromnode, tonode string, domain *string) (err sys.ERROR) {
	if guuid := util.NodeToUUID(groupnode); guuid > 0 {
		if fromnode == tonode {
			return sys.ERR_AUTH
		}
		if g, _ := SelectByIdx[timgroup]("UUID", guuid); g != nil && g.Status != sys.GROUP_STATUS_CANCELLED {
			if tr, _ := TDecode(util.Mask(g.RBean), &TimRoomBean{}); tr != nil {
				if *tr.Founder == fromnode || util.ContainStrings(tr.Managers, fromnode) {
					if *tr.Founder != fromnode && util.ContainStrings(tr.Managers, tonode) {
						return sys.ERR_AUTH
					}
					if util.ContainStrings(tr.Managers, tonode) {
						tr.Managers = util.ArraySub(tr.Managers, []string{tonode})
						UpdateNonzero(&timgroup{Id: g.Id, UUID: guuid, RBean: util.Mask(TEncode(tr))})
					}
					rid := util.RelateIdForGroup(groupnode, tonode, domain)
					numlock.Lock(int64(rid))
					defer numlock.Unlock(int64(rid))
					if as, _ := SelectAllByIdxWithTid[timmucroster](guuid, "Relate", rid); as != nil {
						for _, a := range as {
							Delete[timmucroster](guuid, a.Id)
						}
					}
					tuuid := util.NodeToUUID(tonode)
					if as, _ := SelectAllByIdxWithTid[timmucroster](tuuid, "Relate", rid); as != nil {
						for _, a := range as {
							Delete[timmucroster](tuuid, a.Id)
						}
					}

					ukid := util.UnikId(guuid, tuuid)
					if as, _ := SelectAllByIdxWithTid[timblockroom](guuid, "UnikId", ukid); as != nil {
						for _, a := range as {
							Delete[timblockroom](guuid, a.Id)
						}
					}

					if a, _ := SelectByIdx[timrelate]("UUID", rid); a != nil {
						if a.Status|0x0f == 0x02 {
							a.Status = 0x02
							if UpdateNonzero(a) != nil {
								err = sys.ERR_DATABASE
							}
						} else {
							if Delete[timrelate](a.tid(), a.Id) != nil {
								err = sys.ERR_DATABASE
							}
						}
					} else {
						err = sys.ERR_NOEXIST
					}
				} else {
					return sys.ERR_AUTH
				}
			} else {
				return sys.ERR_UNDEFINED
			}
		}
	} else {
		err = sys.ERR_NOEXIST
	}
	return
}

func (this *tldbhandler) Leavegroup(groupnode, fromnode string, domain *string) (err sys.ERROR) {
	guuid, tuuid := util.NodeToUUID(groupnode), util.NodeToUUID(fromnode)
	if guuid > 0 {
		if err = func() (err sys.ERROR) {
			numlock.Lock(int64(guuid))
			defer numlock.Unlock(int64(guuid))
			if g, _ := SelectByIdx[timgroup]("UUID", guuid); g != nil {
				if g.Status == sys.GROUP_STATUS_CANCELLED {
					return sys.ERR_CANCEL
				}
				if tr, _ := TDecode(util.Mask(g.RBean), &TimRoomBean{}); tr != nil {
					if *tr.Founder == fromnode {
						return sys.ERR_AUTH
					}
					if util.ContainStrings(tr.Managers, fromnode) {
						tr.Managers = util.ArraySub(tr.Managers, []string{fromnode})
						UpdateNonzero(&timgroup{Id: g.Id, UUID: guuid, RBean: util.Mask(TEncode(tr))})
					}
				} else {
					return sys.ERR_UNDEFINED
				}
			} else {
				err = sys.ERR_NOEXIST
			}
			return
		}(); err != nil {
			return
		}
	} else {
		err = sys.ERR_PARAMS
	}
	rid := util.RelateIdForGroup(groupnode, fromnode, domain)
	if as, _ := SelectAllByIdxWithTid[timmucroster](guuid, "Relate", rid); as != nil {
		for _, a := range as {
			Delete[timmucroster](guuid, a.Id)
		}
	}
	if as, _ := SelectAllByIdxWithTid[timmucroster](tuuid, "Relate", rid); as != nil {
		for _, a := range as {
			Delete[timmucroster](tuuid, a.Id)
		}
	}
	ukid := util.UnikId(tuuid, guuid)
	if as, _ := SelectAllByIdxWithTid[timblockroom](tuuid, "UnikId", ukid); as != nil {
		for _, a := range as {
			Delete[timblockroom](tuuid, a.Id)
		}
	}
	numlock.Lock(int64(rid))
	defer numlock.Unlock(int64(rid))
	if a, _ := SelectByIdx[timrelate]("UUID", rid); a != nil {
		if a.Status&0xf0 == 0x20 {
			a.Status = 0x20
			if UpdateNonzero(a) != nil {
				err = sys.ERR_DATABASE
			}
		} else if Delete[timrelate](a.tid(), a.Id) != nil {
			err = sys.ERR_DATABASE
		}
	} else {
		err = sys.ERR_NOEXIST
	}
	return
}

func (this *tldbhandler) Cancelgroup(groupnode, fromnode string, domain *string) (err sys.ERROR) {
	if guuid := util.NodeToUUID(groupnode); guuid > 0 {
		numlock.Lock(int64(guuid))
		defer numlock.Unlock(int64(guuid))
		if g, _ := SelectByIdx[timgroup]("UUID", guuid); g != nil {
			if g.Status == sys.GROUP_STATUS_CANCELLED {
				return sys.ERR_CANCEL
			}
			if tr, _ := TDecode(util.Mask(g.RBean), &TimRoomBean{}); tr != nil {
				if *tr.Founder == fromnode {
					if tus, _ := SelectAllByIdx[timmucroster]("UUID", guuid); tus != nil {
						if len(tus) > 0 {
							for _, tu := range tus {
								if tu.TUUID != util.NodeToUUID(fromnode) {
									return sys.ERR_AUTH
								}
							}
						}
						for _, tu := range tus {
							Delete[timmucroster](guuid, tu.Id)
						}
					}
					g.Status = sys.GROUP_STATUS_CANCELLED
					if UpdateNonzero(g) != nil {
						return sys.ERR_DATABASE
					}
				} else {
					return sys.ERR_AUTH
				}
			} else {
				return sys.ERR_UNDEFINED
			}
		} else {
			err = sys.ERR_NOEXIST
		}
	} else {
		err = sys.ERR_PARAMS
	}
	return
}

func (this *tldbhandler) Blockgroup(groupnode, fromnode string, domain *string) (err sys.ERROR) {
	rid := util.RelateIdForGroup(groupnode, fromnode, domain)
	guuid, tuuid := util.NodeToUUID(groupnode), util.NodeToUUID(fromnode)
	if as, _ := SelectAllByIdxWithTid[timmucroster](guuid, "Relate", rid); as != nil {
		for _, a := range as {
			Delete[timmucroster](guuid, a.Id)
		}
	}
	if as, _ := SelectAllByIdxWithTid[timmucroster](tuuid, "Relate", rid); as != nil {
		for _, a := range as {
			Delete[timmucroster](tuuid, a.Id)
		}
	}
	numlock.Lock(int64(rid))
	defer numlock.Unlock(int64(rid))

	ukid := util.UnikId(tuuid, guuid)
	if a, _ := SelectByIdxWithTid[timblockroom](guuid, "UnikId", ukid); a == nil {
		Insert(&timblockroom{UnikId: ukid, UUID: tuuid, TUUID: guuid})
	}

	if a, _ := SelectByIdx[timrelate]("UUID", rid); a != nil {
		if state := a.Status; state&0x0f != 0x02 {
			a.Status = state&0xf0 | 0x02
			if UpdateNonzero(a) != nil {
				err = sys.ERR_DATABASE
			}
		}
	} else {
		Insert(&timrelate{UUID: rid, Status: 0x02})
	}
	return
}

func (this *tldbhandler) Blockgroupmember(groupnode, fromnode, tonode string, domain *string) (err sys.ERROR) {
	if err = checkAdmin(groupnode, fromnode, tonode); err != nil {
		return
	}
	rid := util.RelateIdForGroup(groupnode, tonode, domain)
	guuid, tuuid := util.NodeToUUID(groupnode), util.NodeToUUID(tonode)
	if as, _ := SelectAllByIdxWithTid[timmucroster](guuid, "Relate", rid); as != nil {
		for _, a := range as {
			Delete[timmucroster](guuid, a.Id)
		}
	}
	if as, _ := SelectAllByIdxWithTid[timmucroster](tuuid, "Relate", rid); as != nil {
		for _, a := range as {
			Delete[timmucroster](tuuid, a.Id)
		}
	}
	numlock.Lock(int64(rid))
	defer numlock.Unlock(int64(rid))

	ukid := util.UnikId(guuid, tuuid)
	if a, _ := SelectByIdxWithTid[timblockroom](guuid, "UnikId", ukid); a == nil {
		Insert(&timblockroom{UnikId: ukid, UUID: guuid, TUUID: tuuid})
	}

	if a, _ := SelectByIdx[timrelate]("UUID", rid); a != nil {
		if state := a.Status; state&0xf0 != 0x20 {
			a.Status = state&0x0f | 0x20
			if UpdateNonzero(a) != nil {
				err = sys.ERR_DATABASE
			}
		}
	} else {
		Insert(&timrelate{UUID: rid, Status: 0x20})
	}
	return
}

func checkAdmin(groupnode, fromnode, tonode string) (err sys.ERROR) {
	if guuid := util.NodeToUUID(groupnode); guuid > 0 {
		if fromnode == tonode {
			return sys.ERR_AUTH
		}
		if g, _ := SelectByIdx[timgroup]("UUID", guuid); g != nil {
			if g.Status == sys.GROUP_STATUS_CANCELLED {
				return sys.ERR_CANCEL
			}
			if tr, _ := TDecode(util.Mask(g.RBean), &TimRoomBean{}); tr != nil {
				if *tr.Founder == fromnode || util.ContainStrings(tr.Managers, fromnode) {
					if *tr.Founder != fromnode && util.ContainStrings(tr.Managers, tonode) {
						return sys.ERR_AUTH
					}
				} else {
					return sys.ERR_AUTH
				}
			} else {
				return sys.ERR_UNDEFINED
			}
		} else {
			return sys.ERR_NOEXIST
		}
	} else {
		err = sys.ERR_ACCOUNT
	}
	return
}

func (this *tldbhandler) ModifyUserInfo(node string, tu *TimUserBean) (err sys.ERROR) {
	uuid := util.NodeToUUID(node)
	if a, _ := SelectByIdx[timuser]("UUID", uuid); a != nil {
		if a.UBean != nil {
			if ub, _ := TDecode(util.Mask(a.UBean), &TimUserBean{}); ub != nil {
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
		UpdateNonzero(&timuser{Id: a.Id, UUID: a.UUID, UBean: util.Mask(TEncode(tu))})
	} else {
		err = sys.ERR_NOEXIST
	}
	return
}
func (this *tldbhandler) GetUserInfo(nodes []string) (m map[string]*TimUserBean, err sys.ERROR) {
	if nodes != nil {
		m = make(map[string]*TimUserBean, 0)
		for _, node := range nodes {
			uuid := util.NodeToUUID(node)
			if a, _ := SelectByIdx[timuser]("UUID", uuid); a != nil {
				if a.UBean != nil {
					if tub, _ := TDecode(util.Mask(a.UBean), &TimUserBean{}); tub != nil {
						tub.Createtime = &a.Createtime
						m[node] = tub
					}
				}
			}
		}
	} else {
		err = sys.ERR_PARAMS
	}
	return
}

func (this *tldbhandler) ModifygroupInfo(node, fnode string, tu *TimRoomBean) (err sys.ERROR) {
	if guuid := util.NodeToUUID(node); guuid > 0 {
		if g, _ := SelectByIdx[timgroup]("UUID", guuid); g != nil {
			if g.Status == sys.GROUP_STATUS_CANCELLED {
				return sys.ERR_CANCEL
			}
			if tr, _ := TDecode(util.Mask(g.RBean), &TimRoomBean{}); tr != nil {
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
					UpdateNonzero(&timgroup{Id: g.Id, UUID: guuid, RBean: util.Mask(TEncode(tr))})
				} else {
					return sys.ERR_AUTH
				}
			} else {
				return sys.ERR_UNDEFINED
			}
		} else {
			return sys.ERR_NOEXIST
		}
	} else {
		err = sys.ERR_ACCOUNT
	}
	return
}

func (this *tldbhandler) GetGroupInfo(nodes []string) (m map[string]*TimRoomBean, err sys.ERROR) {
	if nodes != nil {
		m = make(map[string]*TimRoomBean, 0)
		for _, node := range nodes {
			if guuid := util.NodeToUUID(node); guuid > 0 {
				if g, _ := SelectByIdx[timgroup]("UUID", guuid); g != nil && g.Status != sys.GROUP_STATUS_CANCELLED {
					if tr, _ := TDecode(util.Mask(g.RBean), &TimRoomBean{}); tr != nil {
						m[node] = tr
					}
				}
			}
		}
	} else {
		err = sys.ERR_PARAMS
	}
	return
}
