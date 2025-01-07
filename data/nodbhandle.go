// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of n source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package data

import (
	"github.com/donnie4w/tim/errs"
	. "github.com/donnie4w/tim/stub"
	"github.com/donnie4w/tim/sys"
)

type nodbhandle struct{}

func (n *nodbhandle) Register(username, pwd string, domain *string) (node string, e errs.ERROR) {
	e = errs.ERR_DATABASE
	return
}

func (n *nodbhandle) Login(username, pwd string, domain *string) (_r string, e errs.ERROR) {
	_r = username
	return
}

func (n *nodbhandle) Modify(uint64, *string, string, *string) (e errs.ERROR) {
	return
}

func (n *nodbhandle) AuthNode(username, pwd string, domain *string) (_r string, err errs.ERROR) {
	_r = username
	return
}

func (n *nodbhandle) SaveMessage(tm *TimMessage) (err error) {
	return
}

func (n *nodbhandle) GetMessage(fromNode string, domain *string, rtype int8, to string, mid, limit int64) (tmList []*TimMessage, err error) {
	return
}

func (n *nodbhandle) GetMessageByMid(tid uint64, mid int64) (tm *TimMessage, err error) {
	return
}

func (n *nodbhandle) DelMessageByMid(tid uint64, mid int64) (err error) {
	return
}

func (n *nodbhandle) SaveOfflineMessage(tm *TimMessage) (err error) {
	return
}

func (n *nodbhandle) GetOfflineMessage(node string, limit int) (oblist []*OfflineBean, err error) {
	return
}

func (n *nodbhandle) DelOfflineMessage(tid uint64, ids ...int64) (_r int64, err error) {
	return
}

func (n *nodbhandle) UserGroup(node string, domain *string) (_r []string) {
	return
}

func (n *nodbhandle) GroupRoster(groupnode string) (_r []string) {
	return
}

func (n *nodbhandle) AuthGroupAndUser(groupnode, usernode string, domain *string) (ok bool, err error) {
	ok = true
	return
}

/***************************************************************************************************************/

func (n *nodbhandle) Roster(node string) (_r []string) {
	return
}

func (n *nodbhandle) Blockrosterlist(string) (_r []string) {
	return
}
func (n *nodbhandle) Blockroomlist(string) (_r []string) {
	return
}
func (n *nodbhandle) Blockroommemberlist(string, string) (_r []string) {
	return
}

func (n *nodbhandle) AuthUserAndUser(fnode, tnode string, domain *string) (_r bool) {
	return true
}

func (n *nodbhandle) ExistUser(node string) (_r bool) {
	return true
}

func (n *nodbhandle) ExistGroup(node string) (_r bool) {
	return true
}

func (n *nodbhandle) Addroster(fnode, tnode string, domain *string) (status int8, err errs.ERROR) {
	err = errs.ERR_DATABASE
	return
}

func (n *nodbhandle) Rmroster(fnode, tnode string, domain *string) (mstell bool, ok bool) {
	ok = false
	return
}

func (n *nodbhandle) Blockroster(fnode, tnode string, domain *string) (mstell bool, ok bool) {
	ok = false
	return
}

func (n *nodbhandle) GroupGtype(groupnode string, domain *string) (gtype int8, err errs.ERROR) {
	err = errs.ERR_DATABASE
	return
}

func (n *nodbhandle) GroupManagers(groupnode string, domain *string) (s []string, err errs.ERROR) {
	return
}

func (n *nodbhandle) Newgroup(fnode, groupname string, gtype sys.TIMTYPE, domain *string) (gnode string, err errs.ERROR) {
	return
}

func (n *nodbhandle) Addgroup(groupnode, fromnode string, domain *string) (err errs.ERROR) {
	return errs.ERR_DATABASE
}

func (n *nodbhandle) Pullgroup(groupnode, fromnode, tonode string, domain *string) (isReq bool, err errs.ERROR) {
	err = errs.ERR_DATABASE
	return
}

func (n *nodbhandle) Nopassgroup(groupnode, fromnode, tonode string, domain *string) (err errs.ERROR) {
	return errs.ERR_DATABASE
}

func (n *nodbhandle) Kickgroup(groupnode, fromnode, tonode string, domain *string) (err errs.ERROR) {
	return errs.ERR_DATABASE
}
func (n *nodbhandle) Leavegroup(groupnode, fromnode string, domain *string) (err errs.ERROR) {
	return errs.ERR_DATABASE
}

func (n *nodbhandle) Cancelgroup(groupnode, fromnode string, domain *string) (err errs.ERROR) {
	return errs.ERR_DATABASE
}

func (n *nodbhandle) Blockgroup(groupnode, fromnode string, domain *string) (err errs.ERROR) {
	return errs.ERR_DATABASE
}
func (n *nodbhandle) Blockgroupmember(groupnode, fromnode, tonodme string, domain *string) (err errs.ERROR) {
	return errs.ERR_DATABASE
}

func (n *nodbhandle) ModifygroupInfo(node, fnode string, tu *TimRoomBean) (err errs.ERROR) {
	return errs.ERR_DATABASE
}
func (n *nodbhandle) GetGroupInfo(node []string) (m map[string]*TimRoomBean, err errs.ERROR) {
	err = errs.ERR_DATABASE
	return
}
func (n *nodbhandle) ModifyUserInfo(node string, tu *TimUserBean) (err errs.ERROR) {
	return errs.ERR_DATABASE
}
func (n *nodbhandle) GetUserInfo(node []string) (m map[string]*TimUserBean, err errs.ERROR) {
	err = errs.ERR_DATABASE
	return
}
func (n *nodbhandle) TimAdminAuth(account, password, domain string) bool {
	return true
}
