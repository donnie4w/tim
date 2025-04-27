// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of n source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package data

import (
	"github.com/donnie4w/tim/errs"
	"github.com/donnie4w/tim/log"
	. "github.com/donnie4w/tim/stub"
	"github.com/donnie4w/tim/sys"
)

type nodbHandle struct{}

func (n *nodbHandle) init() service {
	log.FmtPrint("No database schema")
	return n
}

func (n *nodbHandle) Register(username, pwd string, domain *string) (node string, e errs.ERROR) {
	e = errs.ERR_DATABASE
	return
}

func (n *nodbHandle) Login(username, pwd string, domain *string) (_r string, e errs.ERROR) {
	_r = username
	return
}

func (n *nodbHandle) Modify(uint64, *string, string, *string) (e errs.ERROR) {
	return
}

func (n *nodbHandle) AuthNode(username, pwd string, domain *string) (_r string, err errs.ERROR) {
	_r = username
	return
}

func (n *nodbHandle) SaveMessage(tm *TimMessage) (err error) {
	return
}

func (n *nodbHandle) GetMessage(fromNode string, domain *string, rtype int8, to string, mid, timestamp, limit int64) (tmList []*TimMessage, err error) {
	return
}

func (n *nodbHandle) GetFidByMid(tid []byte, mid int64) (fid int64, err error) {
	return
}

func (n *nodbHandle) DelMessageByMid(tid []byte, mid int64) (err error) {
	return
}

func (n *nodbHandle) SaveOfflineMessage(string, *TimMessage) (err error) {
	return
}

func (n *nodbHandle) GetOfflineMessage(node string, limit int) (oblist []*OfflineBean, err error) {
	return
}

func (n *nodbHandle) DelOfflineMessage(tid uint64, ids ...any) (_r int64, err error) {
	return
}

func (n *nodbHandle) UserGroup(node string, domain *string) (_r []string) {
	return
}

func (n *nodbHandle) GroupRoster(groupnode string) (_r []string) {
	return
}

func (n *nodbHandle) AuthGroupAndUser(groupnode, usernode string, domain *string) (ok bool, err error) {
	ok = true
	return
}

/***************************************************************************************************************/

func (n *nodbHandle) Roster(node string) (_r []string) {
	return
}

func (n *nodbHandle) Blockrosterlist(string) (_r []string) {
	return
}
func (n *nodbHandle) Blockroomlist(string) (_r []string) {
	return
}
func (n *nodbHandle) Blockroommemberlist(string, string) (_r []string) {
	return
}

func (n *nodbHandle) AuthUserAndUser(fnode, tnode string, domain *string) (_r bool) {
	return true
}

func (n *nodbHandle) ExistUser(node string) (_r bool) {
	return true
}

func (n *nodbHandle) ExistGroup(node string) (_r bool) {
	return true
}

func (n *nodbHandle) Addroster(fnode, tnode string, domain *string) (status int8, err errs.ERROR) {
	err = errs.ERR_DATABASE
	return
}

func (n *nodbHandle) Rmroster(fnode, tnode string, domain *string) (mstell bool, ok bool) {
	ok = false
	return
}

func (n *nodbHandle) Blockroster(fnode, tnode string, domain *string) (mstell bool, ok bool) {
	ok = false
	return
}

func (n *nodbHandle) GroupGtype(groupnode string, domain *string) (gtype int8, err errs.ERROR) {
	err = errs.ERR_DATABASE
	return
}

func (n *nodbHandle) GroupManagers(groupnode string, domain *string) (s []string, err errs.ERROR) {
	return
}

func (n *nodbHandle) Newgroup(fnode, groupname string, gtype sys.TIMTYPE, domain *string) (gnode string, err errs.ERROR) {
	return
}

func (n *nodbHandle) Addgroup(groupnode, fromnode string, domain *string) (err errs.ERROR) {
	return errs.ERR_DATABASE
}

func (n *nodbHandle) Pullgroup(groupnode, fromnode, tonode string, domain *string) (isReq bool, err errs.ERROR) {
	err = errs.ERR_DATABASE
	return
}

func (n *nodbHandle) Nopassgroup(groupnode, fromnode, tonode string, domain *string) (err errs.ERROR) {
	return errs.ERR_DATABASE
}

func (n *nodbHandle) Kickgroup(groupnode, fromnode, tonode string, domain *string) (err errs.ERROR) {
	return errs.ERR_DATABASE
}
func (n *nodbHandle) Leavegroup(groupnode, fromnode string, domain *string) (err errs.ERROR) {
	return errs.ERR_DATABASE
}

func (n *nodbHandle) Cancelgroup(groupnode, fromnode string, domain *string) (err errs.ERROR) {
	return errs.ERR_DATABASE
}

func (n *nodbHandle) Blockgroup(groupnode, fromnode string, domain *string) (err errs.ERROR) {
	return errs.ERR_DATABASE
}
func (n *nodbHandle) Blockgroupmember(groupnode, fromnode, tonodme string, domain *string) (err errs.ERROR) {
	return errs.ERR_DATABASE
}

func (n *nodbHandle) ModifygroupInfo(node, fnode string, tu *TimRoomBean, admin bool) (err errs.ERROR) {
	return errs.ERR_DATABASE
}
func (n *nodbHandle) GetGroupInfo(node []string) (m map[string]*TimRoomBean, err errs.ERROR) {
	err = errs.ERR_DATABASE
	return
}
func (n *nodbHandle) ModifyUserInfo(node string, tu *TimUserBean) (err errs.ERROR) {
	return errs.ERR_DATABASE
}
func (n *nodbHandle) GetUserInfo(node []string) (m map[string]*TimUserBean, err errs.ERROR) {
	err = errs.ERR_DATABASE
	return
}
func (n *nodbHandle) TimAdminAuth(account, password, domain string) bool {
	return true
}
