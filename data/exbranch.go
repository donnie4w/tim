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

func (eh *externHandle) UserGroup(node string, domain *string) []string {
	return eh.externdb.userGroup(node)
}

func (eh *externHandle) GroupRoster(groupnode string) []string {
	return eh.externdb.groupRoster(groupnode)
}

func (eh *externHandle) Roster(node string) []string {
	return eh.externdb.roster(node)
}

func (eh *externHandle) Blockrosterlist(string) (_r []string) {
	return
}

func (eh *externHandle) Blockroomlist(string) (_r []string) {
	return
}

func (eh *externHandle) Blockroommemberlist(string, string) (_r []string) {
	return
}

func (eh *externHandle) AuthUserAndUser(fnode, tnode string, domain *string) (_r bool) {
	_r, _ = eh.externdb.authUser(fnode, tnode)
	return
}

func (eh *externHandle) ExistUser(node string) (_r bool) {
	_r, _ = eh.externdb.existUser(node)
	return
}

func (eh *externHandle) ExistGroup(node string) (_r bool) {
	_r, _ = eh.externdb.existGroup(node)
	return
}

func (eh *externHandle) Addroster(fnode, tnode string, domain *string) (status int8, err errs.ERROR) {
	err = errs.ERR_INTERFACE
	return
}

func (eh *externHandle) Rmroster(fnode, tnode string, domain *string) (mstell bool, ok bool) {
	return
}

func (eh *externHandle) Blockroster(fnode, tnode string, domain *string) (mstell bool, ok bool) {
	return
}

func (eh *externHandle) GroupGtype(groupnode string, domain *string) (gtype int8, err errs.ERROR) {
	err = errs.ERR_INTERFACE
	return
}

func (eh *externHandle) GroupManagers(groupnode string, domain *string) (s []string, err errs.ERROR) {
	err = errs.ERR_INTERFACE
	return
}

func (eh *externHandle) Newgroup(fnode, groupname string, gtype sys.TIMTYPE, domain *string) (gnode string, err errs.ERROR) {
	err = errs.ERR_INTERFACE
	return
}

func (eh *externHandle) Addgroup(groupnode, fromnode string, domain *string) (err errs.ERROR) {
	err = errs.ERR_INTERFACE
	return
}

func (eh *externHandle) Pullgroup(groupnode, fromnode, tonode string, domain *string) (isReq bool, err errs.ERROR) {
	err = errs.ERR_INTERFACE
	return
}

func (eh *externHandle) Nopassgroup(groupnode, fromnode, tonode string, domain *string) (err errs.ERROR) {
	err = errs.ERR_INTERFACE
	return
}

func (eh *externHandle) Kickgroup(groupnode, fromnode, tonode string, domain *string) (err errs.ERROR) {
	err = errs.ERR_INTERFACE
	return
}
func (eh *externHandle) Leavegroup(groupnode, fromnode string, domain *string) (err errs.ERROR) {
	err = errs.ERR_INTERFACE
	return
}

func (eh *externHandle) Cancelgroup(groupnode, fromnode string, domain *string) (err errs.ERROR) {
	err = errs.ERR_INTERFACE
	return
}

func (eh *externHandle) Blockgroup(groupnode, fromnode string, domain *string) (err errs.ERROR) {
	err = errs.ERR_INTERFACE
	return
}

func (eh *externHandle) Blockgroupmember(groupnode, fromnode, tonodme string, domain *string) (err errs.ERROR) {
	err = errs.ERR_INTERFACE
	return
}

func (eh *externHandle) ModifyUserInfo(node string, tu *stub.TimUserBean) (err errs.ERROR) {
	return
}
func (eh *externHandle) GetUserInfo(node []string) (m map[string]*stub.TimUserBean, err errs.ERROR) {
	return
}

func (eh *externHandle) ModifygroupInfo(node, fnode string, tu *stub.TimRoomBean, admin bool) (err errs.ERROR) {
	return
}
func (eh *externHandle) GetGroupInfo(node []string) (m map[string]*stub.TimRoomBean, err errs.ERROR) {
	return
}

func (eh *externHandle) TimAdminAuth(account, password, domain string) bool {
	return false
}
