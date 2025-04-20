// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package data

import (
	"github.com/donnie4w/tim/errs"
	"github.com/donnie4w/tim/stub"
	"github.com/donnie4w/tim/sys"
)

type cassandraHandle struct{}

func (ch *cassandraHandle) init() service {
	if err := camanager.init(); err != nil {
		panic(err)
	}
	return ch
}

func (ch *cassandraHandle) Register(username, pwd string, domain *string) (node string, e errs.ERROR) {
	e = errs.ERR_DATABASE
	return
}

func (ch *cassandraHandle) Login(username, pwd string, domain *string) (_r string, e errs.ERROR) {
	_r = username
	return
}

func (ch *cassandraHandle) Modify(uint64, *string, string, *string) (e errs.ERROR) {
	return
}

func (ch *cassandraHandle) AuthNode(username, pwd string, domain *string) (_r string, err errs.ERROR) {
	_r = username
	return
}

func (ch *cassandraHandle) SaveMessage(tm *stub.TimMessage) (err error) {
	return
}

func (ch *cassandraHandle) GetMessage(fromNode string, domain *string, rtype int8, to string, mid, timestamp, limit int64) (tmList []*stub.TimMessage, err error) {
	return
}

func (ch *cassandraHandle) GetFidByMid(tid []byte, mid int64) (fid int64, err error) {
	return
}

func (ch *cassandraHandle) DelMessageByMid(tid []byte, mid int64) (err error) {
	return
}

func (ch *cassandraHandle) SaveOfflineMessage(string, *stub.TimMessage) (err error) {
	return
}

func (ch *cassandraHandle) GetOfflineMessage(node string, limit int) (oblist []*OfflineBean, err error) {
	return
}

func (ch *cassandraHandle) DelOfflineMessage(tid uint64, ids ...any) (_r int64, err error) {
	return
}

func (ch *cassandraHandle) UserGroup(node string, domain *string) (_r []string) {
	return
}

func (ch *cassandraHandle) GroupRoster(groupnode string) (_r []string) {
	return
}

func (ch *cassandraHandle) AuthGroupAndUser(groupnode, usernode string, domain *string) (ok bool, err error) {
	ok = true
	return
}

/***************************************************************************************************************/

func (ch *cassandraHandle) Roster(node string) (_r []string) {
	return
}

func (ch *cassandraHandle) Blockrosterlist(string) (_r []string) {
	return
}
func (ch *cassandraHandle) Blockroomlist(string) (_r []string) {
	return
}
func (ch *cassandraHandle) Blockroommemberlist(string, string) (_r []string) {
	return
}

func (ch *cassandraHandle) AuthUserAndUser(fnode, tnode string, domain *string) (_r bool) {
	return true
}

func (ch *cassandraHandle) ExistUser(node string) (_r bool) {
	return true
}

func (ch *cassandraHandle) ExistGroup(node string) (_r bool) {
	return true
}

func (ch *cassandraHandle) Addroster(fnode, tnode string, domain *string) (status int8, err errs.ERROR) {
	err = errs.ERR_DATABASE
	return
}

func (ch *cassandraHandle) Rmroster(fnode, tnode string, domain *string) (mstell bool, ok bool) {
	ok = false
	return
}

func (ch *cassandraHandle) Blockroster(fnode, tnode string, domain *string) (mstell bool, ok bool) {
	ok = false
	return
}

func (ch *cassandraHandle) GroupGtype(groupnode string, domain *string) (gtype int8, err errs.ERROR) {
	err = errs.ERR_DATABASE
	return
}

func (ch *cassandraHandle) GroupManagers(groupnode string, domain *string) (s []string, err errs.ERROR) {
	return
}

func (ch *cassandraHandle) Newgroup(fnode, groupname string, gtype sys.TIMTYPE, domain *string) (gnode string, err errs.ERROR) {
	return
}

func (ch *cassandraHandle) Addgroup(groupnode, fromnode string, domain *string) (err errs.ERROR) {
	return errs.ERR_DATABASE
}

func (ch *cassandraHandle) Pullgroup(groupnode, fromnode, tonode string, domain *string) (isReq bool, err errs.ERROR) {
	err = errs.ERR_DATABASE
	return
}

func (ch *cassandraHandle) Nopassgroup(groupnode, fromnode, tonode string, domain *string) (err errs.ERROR) {
	return errs.ERR_DATABASE
}

func (ch *cassandraHandle) Kickgroup(groupnode, fromnode, tonode string, domain *string) (err errs.ERROR) {
	return errs.ERR_DATABASE
}
func (ch *cassandraHandle) Leavegroup(groupnode, fromnode string, domain *string) (err errs.ERROR) {
	return errs.ERR_DATABASE
}

func (ch *cassandraHandle) Cancelgroup(groupnode, fromnode string, domain *string) (err errs.ERROR) {
	return errs.ERR_DATABASE
}

func (ch *cassandraHandle) Blockgroup(groupnode, fromnode string, domain *string) (err errs.ERROR) {
	return errs.ERR_DATABASE
}
func (ch *cassandraHandle) Blockgroupmember(groupnode, fromnode, tonodme string, domain *string) (err errs.ERROR) {
	return errs.ERR_DATABASE
}

func (ch *cassandraHandle) ModifygroupInfo(node, fnode string, tu *stub.TimRoomBean, admin bool) (err errs.ERROR) {
	return errs.ERR_DATABASE
}
func (ch *cassandraHandle) GetGroupInfo(node []string) (m map[string]*stub.TimRoomBean, err errs.ERROR) {
	err = errs.ERR_DATABASE
	return
}
func (ch *cassandraHandle) ModifyUserInfo(node string, tu *stub.TimUserBean) (err errs.ERROR) {
	return errs.ERR_DATABASE
}
func (ch *cassandraHandle) GetUserInfo(node []string) (m map[string]*stub.TimUserBean, err errs.ERROR) {
	err = errs.ERR_DATABASE
	return
}
func (ch *cassandraHandle) TimAdminAuth(account, password, domain string) bool {
	return true
}
