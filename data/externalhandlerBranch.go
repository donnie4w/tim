// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package data

import (
	"github.com/donnie4w/tim/errs"
	. "github.com/donnie4w/tim/stub"
	"github.com/donnie4w/tim/sys"
)

func (eh *externalhandle) UserGroup(node string, domain *string) []string {
	return eh.externaldb.userGroup(node)
}

func (eh *externalhandle) GroupRoster(groupnode string) []string {
	return eh.externaldb.groupRoster(groupnode)
}

func (eh *externalhandle) Roster(node string) []string {
	return eh.externaldb.roster(node)
}

func (eh *externalhandle) Blockrosterlist(string) (_r []string) {
	return
}

func (eh *externalhandle) Blockroomlist(string) (_r []string) {
	return
}

func (eh *externalhandle) Blockroommemberlist(string, string) (_r []string) {
	return
}

func (eh *externalhandle) AuthUserAndUser(fnode, tnode string, domain *string) (_r bool) {
	_r, _ = eh.externaldb.authUser(fnode, tnode)
	return
}

func (eh *externalhandle) ExistUser(node string) (_r bool) {
	_r, _ = eh.externaldb.existUser(node)
	return
}

func (eh *externalhandle) ExistGroup(node string) (_r bool) {
	_r, _ = eh.externaldb.existGroup(node)
	return
}

func (eh *externalhandle) Addroster(fnode, tnode string, domain *string) (status int8, err errs.ERROR) {
	err = errs.ERR_INTERFACE
	return
}

func (eh *externalhandle) Rmroster(fnode, tnode string, domain *string) (mstell bool, ok bool) {
	return
}

func (eh *externalhandle) Blockroster(fnode, tnode string, domain *string) (mstell bool, ok bool) {
	return
}

func (eh *externalhandle) GroupGtype(groupnode string, domain *string) (gtype int8, err errs.ERROR) {
	err = errs.ERR_INTERFACE
	return
}

func (eh *externalhandle) GroupManagers(groupnode string, domain *string) (s []string, err errs.ERROR) {
	err = errs.ERR_INTERFACE
	return
}

func (eh *externalhandle) Newgroup(fnode, groupname string, gtype sys.TIMTYPE, domain *string) (gnode string, err errs.ERROR) {
	err = errs.ERR_INTERFACE
	return
}

func (eh *externalhandle) Addgroup(groupnode, fromnode string, domain *string) (err errs.ERROR) {
	err = errs.ERR_INTERFACE
	return
}

func (eh *externalhandle) Pullgroup(groupnode, fromnode, tonode string, domain *string) (isReq bool, err errs.ERROR) {
	err = errs.ERR_INTERFACE
	return
}

func (eh *externalhandle) Nopassgroup(groupnode, fromnode, tonode string, domain *string) (err errs.ERROR) {
	err = errs.ERR_INTERFACE
	return
}

func (eh *externalhandle) Kickgroup(groupnode, fromnode, tonode string, domain *string) (err errs.ERROR) {
	err = errs.ERR_INTERFACE
	return
}
func (eh *externalhandle) Leavegroup(groupnode, fromnode string, domain *string) (err errs.ERROR) {
	err = errs.ERR_INTERFACE
	return
}

func (eh *externalhandle) Cancelgroup(groupnode, fromnode string, domain *string) (err errs.ERROR) {
	err = errs.ERR_INTERFACE
	return
}

func (eh *externalhandle) Blockgroup(groupnode, fromnode string, domain *string) (err errs.ERROR) {
	err = errs.ERR_INTERFACE
	return
}

func (eh *externalhandle) Blockgroupmember(groupnode, fromnode, tonodme string, domain *string) (err errs.ERROR) {
	err = errs.ERR_INTERFACE
	return
}

func (eh *externalhandle) ModifyUserInfo(node string, tu *TimUserBean) (err errs.ERROR) {
	return
}
func (eh *externalhandle) GetUserInfo(node []string) (m map[string]*TimUserBean, err errs.ERROR) {
	return
}

func (eh *externalhandle) ModifygroupInfo(node, fnode string, tu *TimRoomBean) (err errs.ERROR) {
	return
}
func (eh *externalhandle) GetGroupInfo(node []string) (m map[string]*TimRoomBean, err errs.ERROR) {
	return
}

func (eh *externalhandle) TimAdminAuth(account, password, domain string) bool {
	return false
}
