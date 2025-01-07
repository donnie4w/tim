// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package data

import (
	"github.com/donnie4w/tim/errs"
	"time"

	. "github.com/donnie4w/gofer/util"
	. "github.com/donnie4w/tim/stub"
	"github.com/donnie4w/tim/sys"
	"github.com/donnie4w/tim/util"
)

func (th *tldbhandle) Roster(node string) (_r []string) {
	_r = make([]string, 0)
	if tr, err := SelectAllByIdx[timroster]("UUID", util.NodeToUUID(node)); err == nil {
		for _, a := range tr {
			_r = append(_r, util.UUIDToNode(a.TUUID))
		}
	}
	return
}

func (th *tldbhandle) Blockrosterlist(node string) (_r []string) {
	_r = make([]string, 0)
	if tr, err := SelectAllByIdx[timblock]("UUID", util.NodeToUUID(node)); err == nil {
		for _, a := range tr {
			_r = append(_r, util.UUIDToNode(a.TUUID))
		}
	}
	return
}

func (th *tldbhandle) Blockroomlist(node string) (_r []string) {
	_r = make([]string, 0)
	if tr, err := SelectAllByIdx[timblockroom]("UUID", util.NodeToUUID(node)); err == nil {
		for _, a := range tr {
			_r = append(_r, util.UUIDToNode(a.TUUID))
		}
	}
	return
}

func (th *tldbhandle) Blockroommemberlist(node string, fnode string) (_r []string) {
	if th.checkAdmin(node, fnode, "") != nil {
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

func (th *tldbhandle) UserGroup(node string, domain *string) (_r []string) {
	if tr, _ := SelectAllByIdx[timmucroster]("UUID", util.NodeToUUID(node)); tr != nil && len(tr) > 0 {
		_r = make([]string, 0)
		for _, a := range tr {
			_r = append(_r, util.UUIDToNode(a.TUUID))
		}
	}
	return
}

func (th *tldbhandle) GroupRoster(groupnode string) (_r []string) {
	if tr, _ := SelectAllByIdx[timmucroster]("UUID", util.NodeToUUID(groupnode)); tr != nil && len(tr) > 0 {
		_r = make([]string, 0)
		for _, a := range tr {
			_r = append(_r, util.UUIDToNode(a.TUUID))
		}
	}
	return
}

func (th *tldbhandle) Addroster(fnode, tnode string, domain *string) (status int8, err errs.ERROR) {
	cid := util.ChatIdByNode(fnode, tnode, domain)
	uuid1, uuid2 := util.NodeToUUID(fnode), util.NodeToUUID(tnode)
	numlock.Lock(int64(cid))
	defer numlock.Unlock(int64(cid))
	if a, _ := SelectByIdx[timrelate]("UUID", cid); a != nil {
		if a.Status == 0x11 {
			return 0x11, errs.ERR_REPEAT
		}

		stat := a.Status
		if uuid1 > uuid2 {
			if a.Status&0x0f == 0x02 {
				err = errs.ERR_BLOCK
				return
			}
			a.Status = 0x10 | (a.Status & 0x0f)
		} else {
			if a.Status&0xf0 == 0x20 {
				err = errs.ERR_BLOCK
				return
			}
			a.Status = (a.Status & 0xf0) | 0x01
		}
		if stat != a.Status {
			UpdateNonzero(a)
		}
		status = int8(a.Status)
		if stat != 0x11 && status == 0x11 {
			Insert(&timroster{Unikid: util.UnikIdByNode(fnode, tnode, domain), UUID: uuid1, TUUID: uuid2})
			Insert(&timroster{Unikid: util.UnikIdByNode(tnode, fnode, domain), UUID: uuid2, TUUID: uuid1})
		}
	} else {
		status = 0x10
		if uuid1 < uuid2 {
			status = 0x01
		}
		if id, _ := Insert(&timrelate{UUID: cid, Status: uint8(status)}); id == 0 {
			err = errs.ERR_DATABASE
		}
	}
	return
}

func (th *tldbhandle) Rmroster(fnode, tnode string, domain *string) (mstell bool, ok bool) {
	cid := util.ChatIdByNode(fnode, tnode, domain)
	uuid1, uuid2 := util.NodeToUUID(fnode), util.NodeToUUID(tnode)
	numlock.Lock(int64(cid))
	defer numlock.Unlock(int64(cid))
	if as, _ := SelectAllByIdxWithTid[timroster](uuid1, "Unikid", util.UnikIdByNode(fnode, tnode, domain)); as != nil {
		for _, a := range as {
			Delete[timroster](uuid1, a.Id)
		}
	}
	if as, _ := SelectAllByIdxWithTid[timroster](uuid2, "Unikid", util.UnikIdByNode(tnode, fnode, domain)); as != nil {
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
				Delete[timrelate](a.Tid(), a.Id)
			}
		} else {
			if a.Status = a.Status & 0xf0; a.Status == 0x20 {
				UpdateNonzero(a)
			} else {
				mstell = true
				Delete[timrelate](a.Tid(), a.Id)
			}
		}
		ok = true
	}
	return
}

func (th *tldbhandle) Blockroster(fnode, tnode string, domain *string) (mstell bool, ok bool) {
	cid := util.ChatIdByNode(fnode, tnode, domain)
	uuid1, uuid2 := util.NodeToUUID(fnode), util.NodeToUUID(tnode)
	numlock.Lock(int64(cid))
	defer numlock.Unlock(int64(cid))

	if as, _ := SelectAllByIdxWithTid[timroster](uuid1, "Unikid", util.UnikIdByNode(fnode, tnode, domain)); as != nil {
		for _, a := range as {
			Delete[timroster](uuid1, a.Id)
		}
	}
	if as, _ := SelectAllByIdxWithTid[timroster](uuid2, "Unikid", util.UnikIdByNode(tnode, fnode, domain)); as != nil {
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

func (th *tldbhandle) GroupGtype(groupnode string, domain *string) (gtype int8, err errs.ERROR) {
	if guuid := util.NodeToUUID(groupnode); guuid > 0 {
		if g, _ := SelectByIdx[timgroup]("UUID", guuid); g != nil {
			if sys.TIMTYPE(g.Status) == sys.GROUP_STATUS_CANCELLED {
				return 0, errs.ERR_CANCEL
			}
			gtype = g.Gtype
		} else {
			return 0, errs.ERR_PARAMS
		}
	} else {
		err = errs.ERR_NOEXIST
	}
	return
}

func (th *tldbhandle) GroupManagers(groupnode string, domain *string) (s []string, err errs.ERROR) {
	if guuid := util.NodeToUUID(groupnode); guuid > 0 {
		if g, _ := SelectByIdx[timgroup]("UUID", guuid); g != nil {
			if sys.TIMTYPE(g.Status) == sys.GROUP_STATUS_CANCELLED {
				return nil, errs.ERR_CANCEL
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
			return nil, errs.ERR_PARAMS
		}
	} else {
		err = errs.ERR_NOEXIST
	}
	return
}

func (th *tldbhandle) Newgroup(fnode, groupname string, gtype sys.TIMTYPE, domain *string) (gnode string, err errs.ERROR) {
	fuuid := util.NodeToUUID(fnode)
	if fuuid == 0 || !checkuseruuid(fuuid) {
		err = errs.ERR_ACCOUNT
		return
	}

	UUID := util.CreateUUID(string(Int64ToBytes(UUID64())), domain)
	tg := &timgroup{Gtype: int8(gtype), UUID: UUID, Createtime: time.Now().UnixNano(), Status: int8(sys.GROUP_STATUS_ALIVE)}
	if gtype != sys.GROUP_PRIVATE {
		gtype = sys.GROUP_OPEN
	}
	gt := int8(gtype)
	ubean := &TimRoomBean{Founder: &fnode, Topic: &groupname, Createtime: &tg.Createtime, Gtype: &gt}
	tg.RBean = util.Mask(TEncode(ubean))
	if id, _ := Insert(tg); id == 0 {
		return "", errs.ERR_DATABASE
	}
	gnode = util.UUIDToNode(UUID)
	rid := util.RelateIdForGroup(gnode, fnode, domain)

	if id, _ := Insert(&timrelate{UUID: rid, Status: 0x11}); id == 0 {
		return "", errs.ERR_DATABASE
	}

	if id, _ := Insert(&timmucroster{Unikid: util.UnikIdByNode(gnode, fnode, domain), UUID: UUID, TUUID: fuuid}); id == 0 {
		return "", errs.ERR_DATABASE
	}
	if id, _ := Insert(&timmucroster{Unikid: util.UnikIdByNode(fnode, gnode, domain), UUID: fuuid, TUUID: UUID}); id == 0 {
		return "", errs.ERR_DATABASE
	}
	return
}

func (th *tldbhandle) Addgroup(groupnode, fromnode string, domain *string) (err errs.ERROR) {
	if guuid := util.NodeToUUID(groupnode); guuid > 0 {
		rid := util.RelateIdForGroup(groupnode, fromnode, domain)
		if g, _ := SelectByIdx[timgroup]("UUID", guuid); g != nil {
			if sys.TIMTYPE(g.Status) == sys.GROUP_STATUS_CANCELLED {
				return errs.ERR_CANCEL
			}
			numlock.Lock(int64(rid))
			defer numlock.Unlock(int64(rid))
			if a, _ := SelectByIdx[timrelate]("UUID", rid); a != nil {
				if a.Status&0xf0 == 0x20 {
					return errs.ERR_BLOCK
				}
				if a.Status == 0x11 {
					return errs.ERR_HASEXIST
				}
				if sys.TIMTYPE(g.Gtype) == sys.GROUP_OPEN && a.Status != 0x11 {
					a.Status = 0x11
					UpdateNonzero(a)
					Insert(&timmucroster{Unikid: util.UnikIdByNode(groupnode, fromnode, domain), UUID: guuid, TUUID: util.NodeToUUID(fromnode)})
					Insert(&timmucroster{Unikid: util.UnikIdByNode(fromnode, groupnode, domain), UUID: util.NodeToUUID(fromnode), TUUID: guuid})
				} else if sys.TIMTYPE(g.Gtype) == sys.GROUP_PRIVATE && a.Status != 0x01 {
					a.Status = 0x01
					UpdateNonzero(a)
				}
				return
			} else {
				if sys.TIMTYPE(g.Gtype) == sys.GROUP_PRIVATE {
					if id, _ := Insert(&timrelate{UUID: rid, Status: 0x01}); id == 0 {
						return errs.ERR_DATABASE
					}
				} else if sys.TIMTYPE(g.Gtype) == sys.GROUP_OPEN {
					if id, _ := Insert(&timrelate{UUID: rid, Status: 0x11}); id == 0 {
						return errs.ERR_DATABASE
					}
					Insert(&timmucroster{Unikid: util.UnikIdByNode(groupnode, fromnode, domain), UUID: guuid, TUUID: util.NodeToUUID(fromnode)})
					Insert(&timmucroster{Unikid: util.UnikIdByNode(fromnode, groupnode, domain), UUID: util.NodeToUUID(fromnode), TUUID: guuid})
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

func (th *tldbhandle) Pullgroup(groupnode, fromnode, tonode string, domain *string) (isReq bool, err errs.ERROR) {
	if guuid := util.NodeToUUID(groupnode); guuid > 0 {
		if g, _ := SelectByIdx[timgroup]("UUID", guuid); g != nil {
			if sys.TIMTYPE(g.Status) == sys.GROUP_STATUS_CANCELLED {
				return isReq, errs.ERR_CANCEL
			}
			if sys.TIMTYPE(g.Gtype) == sys.GROUP_PRIVATE {
				if tr, _ := TDecode(util.Mask(g.RBean), &TimRoomBean{}); tr != nil {
					if *tr.Founder != fromnode && !util.ContainStrings(tr.Managers, fromnode) {
						err = errs.ERR_PERM_DENIED
						return
					}
				}
			}
			rid := util.RelateIdForGroup(groupnode, tonode, domain)
			numlock.Lock(int64(rid))
			defer numlock.Unlock(int64(rid))
			if a, _ := SelectByIdx[timrelate]("UUID", rid); a != nil {
				if a.Status&0x0f == 0x02 {
					return isReq, errs.ERR_BLOCK
				}
				if a.Status == 0x11 {
					return isReq, errs.ERR_HASEXIST
				}
				isReq = a.Status == 0x01
				if a.Status != 0x11 {
					a.Status = 0x11
					UpdateNonzero(a)
				}
			} else if id, _ := Insert(&timrelate{UUID: rid, Status: 0x11}); id == 0 {
				return isReq, errs.ERR_DATABASE
			}
			Insert(&timmucroster{Unikid: util.UnikIdByNode(groupnode, tonode, domain), UUID: guuid, TUUID: util.NodeToUUID(tonode)})
			Insert(&timmucroster{Unikid: util.UnikIdByNode(tonode, groupnode, domain), UUID: util.NodeToUUID(tonode), TUUID: guuid})
		} else {
			return isReq, errs.ERR_NOEXIST
		}
	} else {
		err = errs.ERR_NOEXIST
	}
	return
}

func (th *tldbhandle) Nopassgroup(groupnode, fromnode, tonode string, domain *string) (err errs.ERROR) {
	if guuid := util.NodeToUUID(groupnode); guuid > 0 {
		if g, _ := SelectByIdx[timgroup]("UUID", guuid); g != nil {
			if sys.TIMTYPE(g.Status) == sys.GROUP_STATUS_CANCELLED {
				return errs.ERR_CANCEL
			}
			if tr, _ := TDecode(util.Mask(g.RBean), &TimRoomBean{}); tr != nil {
				if *tr.Founder == fromnode || util.ContainStrings(tr.Managers, fromnode) {
					rid := util.RelateIdForGroup(groupnode, tonode, domain)
					numlock.Lock(int64(rid))
					defer numlock.Unlock(int64(rid))
					if a, _ := SelectByIdx[timrelate]("UUID", rid); a != nil {
						if a.Status == 0x01 {
							Delete[timrelate](a.Tid(), a.Id)
							return
						}
						if a.Status|0xf0 != 0 {
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

func (th *tldbhandle) Kickgroup(groupnode, fromnode, tonode string, domain *string) (err errs.ERROR) {
	if guuid := util.NodeToUUID(groupnode); guuid > 0 {
		if fromnode == tonode {
			return errs.ERR_PERM_DENIED
		}
		if g, _ := SelectByIdx[timgroup]("UUID", guuid); g != nil && sys.TIMTYPE(g.Status) != sys.GROUP_STATUS_CANCELLED {
			if tr, _ := TDecode(util.Mask(g.RBean), &TimRoomBean{}); tr != nil {
				if *tr.Founder == fromnode || util.ContainStrings(tr.Managers, fromnode) {
					if *tr.Founder != fromnode && util.ContainStrings(tr.Managers, tonode) {
						return errs.ERR_PERM_DENIED
					}
					if util.ContainStrings(tr.Managers, tonode) {
						tr.Managers = util.ArraySub(tr.Managers, []string{tonode})
						UpdateNonzero(&timgroup{Id: g.Id, UUID: guuid, RBean: util.Mask(TEncode(tr))})
					}
					rid := util.RelateIdForGroup(groupnode, tonode, domain)
					numlock.Lock(int64(rid))
					defer numlock.Unlock(int64(rid))
					if as, _ := SelectAllByIdxWithTid[timmucroster](guuid, "Unikid", util.UnikIdByNode(groupnode, tonode, domain)); as != nil {
						for _, a := range as {
							Delete[timmucroster](guuid, a.Id)
						}
					}
					tuuid := util.NodeToUUID(tonode)
					if as, _ := SelectAllByIdxWithTid[timmucroster](tuuid, "Unikid", util.UnikIdByNode(tonode, groupnode, domain)); as != nil {
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
								err = errs.ERR_DATABASE
							}
						} else {
							if Delete[timrelate](a.Tid(), a.Id) != nil {
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

func (th *tldbhandle) Leavegroup(groupnode, fromnode string, domain *string) (err errs.ERROR) {
	guuid, tuuid := util.NodeToUUID(groupnode), util.NodeToUUID(fromnode)
	if guuid > 0 {
		if err = func() (err errs.ERROR) {
			numlock.Lock(int64(guuid))
			defer numlock.Unlock(int64(guuid))
			if g, _ := SelectByIdx[timgroup]("UUID", guuid); g != nil {
				if sys.TIMTYPE(g.Status) == sys.GROUP_STATUS_CANCELLED {
					return errs.ERR_CANCEL
				}
				if tr, _ := TDecode(util.Mask(g.RBean), &TimRoomBean{}); tr != nil {
					if *tr.Founder == fromnode {
						return errs.ERR_PERM_DENIED
					}
					if util.ContainStrings(tr.Managers, fromnode) {
						tr.Managers = util.ArraySub(tr.Managers, []string{fromnode})
						UpdateNonzero(&timgroup{Id: g.Id, UUID: guuid, RBean: util.Mask(TEncode(tr))})
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
	if as, _ := SelectAllByIdxWithTid[timmucroster](guuid, "Unikid", util.UnikIdByNode(groupnode, fromnode, domain)); as != nil {
		for _, a := range as {
			Delete[timmucroster](guuid, a.Id)
		}
	}
	if as, _ := SelectAllByIdxWithTid[timmucroster](tuuid, "Unikid", util.UnikIdByNode(fromnode, groupnode, domain)); as != nil {
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
				err = errs.ERR_DATABASE
			}
		} else if Delete[timrelate](a.Tid(), a.Id) != nil {
			err = errs.ERR_DATABASE
		}
	} else {
		err = errs.ERR_NOEXIST
	}
	return
}

func (th *tldbhandle) Cancelgroup(groupnode, fromnode string, domain *string) (err errs.ERROR) {
	if guuid := util.NodeToUUID(groupnode); guuid > 0 {
		numlock.Lock(int64(guuid))
		defer numlock.Unlock(int64(guuid))
		if g, _ := SelectByIdx[timgroup]("UUID", guuid); g != nil {
			if sys.TIMTYPE(g.Status) == sys.GROUP_STATUS_CANCELLED {
				return errs.ERR_CANCEL
			}
			if tr, _ := TDecode(util.Mask(g.RBean), &TimRoomBean{}); tr != nil {
				if *tr.Founder == fromnode {
					if tus, _ := SelectAllByIdx[timmucroster]("UUID", guuid); tus != nil {
						if len(tus) > 0 {
							for _, tu := range tus {
								if tu.TUUID != util.NodeToUUID(fromnode) {
									return errs.ERR_PERM_DENIED
								}
							}
						}
						for _, tu := range tus {
							Delete[timmucroster](guuid, tu.Id)
						}
					}
					g.Status = int8(sys.GROUP_STATUS_CANCELLED)
					if UpdateNonzero(g) != nil {
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

func (th *tldbhandle) Blockgroup(groupnode, fromnode string, domain *string) (err errs.ERROR) {
	rid := util.RelateIdForGroup(groupnode, fromnode, domain)
	guuid, tuuid := util.NodeToUUID(groupnode), util.NodeToUUID(fromnode)
	if as, _ := SelectAllByIdxWithTid[timmucroster](guuid, "Unikid", util.UnikIdByNode(groupnode, fromnode, domain)); as != nil {
		for _, a := range as {
			Delete[timmucroster](guuid, a.Id)
		}
	}
	if as, _ := SelectAllByIdxWithTid[timmucroster](tuuid, "Unikid", util.UnikIdByNode(fromnode, groupnode, domain)); as != nil {
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
				err = errs.ERR_DATABASE
			}
		}
	} else {
		Insert(&timrelate{UUID: rid, Status: 0x02})
	}
	return
}

func (th *tldbhandle) Blockgroupmember(groupnode, fromnode, tonode string, domain *string) (err errs.ERROR) {
	if err = th.checkAdmin(groupnode, fromnode, tonode); err != nil {
		return
	}
	rid := util.RelateIdForGroup(groupnode, tonode, domain)
	guuid, tuuid := util.NodeToUUID(groupnode), util.NodeToUUID(tonode)
	if as, _ := SelectAllByIdxWithTid[timmucroster](guuid, "Unikid", util.UnikIdByNode(groupnode, fromnode, domain)); as != nil {
		for _, a := range as {
			Delete[timmucroster](guuid, a.Id)
		}
	}
	if as, _ := SelectAllByIdxWithTid[timmucroster](tuuid, "Unikid", util.UnikIdByNode(fromnode, groupnode, domain)); as != nil {
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
				err = errs.ERR_DATABASE
			}
		}
	} else {
		Insert(&timrelate{UUID: rid, Status: 0x20})
	}
	return
}

func (th *tldbhandle) checkAdmin(groupnode, fromnode, tonode string) (err errs.ERROR) {
	if guuid := util.NodeToUUID(groupnode); guuid > 0 {
		if fromnode == tonode {
			return errs.ERR_PERM_DENIED
		}
		if g, _ := SelectByIdx[timgroup]("UUID", guuid); g != nil {
			if sys.TIMTYPE(g.Status) == sys.GROUP_STATUS_CANCELLED {
				return errs.ERR_CANCEL
			}
			if tr, _ := TDecode(util.Mask(g.RBean), &TimRoomBean{}); tr != nil {
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

func (th *tldbhandle) ModifyUserInfo(node string, tu *TimUserBean) (err errs.ERROR) {
	uuid := util.NodeToUUID(node)
	if uuid == 0 {
		return errs.ERR_ACCOUNT
	}
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
		err = errs.ERR_NOEXIST
	}
	return
}
func (th *tldbhandle) GetUserInfo(nodes []string) (m map[string]*TimUserBean, err errs.ERROR) {
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
		err = errs.ERR_PARAMS
	}
	return
}

func (th *tldbhandle) ModifygroupInfo(node, fnode string, tu *TimRoomBean) (err errs.ERROR) {
	if guuid := util.NodeToUUID(node); guuid > 0 {
		if g, _ := SelectByIdx[timgroup]("UUID", guuid); g != nil {
			if sys.TIMTYPE(g.Status) == sys.GROUP_STATUS_CANCELLED {
				return errs.ERR_CANCEL
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

func (th *tldbhandle) GetGroupInfo(nodes []string) (m map[string]*TimRoomBean, err errs.ERROR) {
	if nodes != nil {
		m = make(map[string]*TimRoomBean, 0)
		for _, node := range nodes {
			if guuid := util.NodeToUUID(node); guuid > 0 {
				if g, _ := SelectByIdx[timgroup]("UUID", guuid); g != nil && sys.TIMTYPE(g.Status) != sys.GROUP_STATUS_CANCELLED {
					if tr, _ := TDecode(util.Mask(g.RBean), &TimRoomBean{}); tr != nil {
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

func (th *tldbhandle) TimAdminAuth(account, password, domain string) bool {
	return false
}
