// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim
package data

import (
	. "github.com/donnie4w/tim/stub"
	"github.com/donnie4w/tim/sys"
)

func (this *nilprocess) Roster(node string) (_r []string) {
	return
}

func (this *nilprocess) Blockrosterlist(string) (_r []string) {
	return
}
func (this *nilprocess) Blockroomlist(string) (_r []string) {
	return
}
func (this *nilprocess) Blockroommemberlist(string, string) (_r []string) {
	return
}

func (this *nilprocess) AuthUserAndUser(fnode, tnode string, domain *string) (_r bool) {
	return true
}

func (this *nilprocess) ExistUser(node string) (_r bool) {
	return true
}

func (this *nilprocess) ExistGroup(node string) (_r bool) {
	return true
}

func (this *nilprocess) Addroster(fnode, tnode string, domain *string) (status int8, err sys.ERROR) {
	err = sys.ERR_DATABASE
	return
}

func (this *nilprocess) Rmroster(fnode, tnode string, domain *string) (mstell bool, ok bool) {
	ok = false
	return
}

func (this *nilprocess) Blockroster(fnode, tnode string, domain *string) (mstell bool, ok bool) {
	ok = false
	return
}

func (this *nilprocess) GroupGtype(groupnode string, domain *string) (gtype int8, err sys.ERROR) {
	err = sys.ERR_DATABASE
	return
}

func (this *nilprocess) GroupManagers(groupnode string, domain *string) (s []string, err sys.ERROR) {
	return
}

func (this *nilprocess) Newgroup(fnode, groupname string, gtype int8, domain *string) (gnode string, err sys.ERROR) {
	return
}

func (this *nilprocess) Addgroup(groupnode, fromnode string, domain *string) (err sys.ERROR) {
	return sys.ERR_DATABASE
}

func (this *nilprocess) Pullgroup(groupnode, fromnode, tonode string, domain *string) (isReq bool, err sys.ERROR) {
	err = sys.ERR_DATABASE
	return
}

func (this *nilprocess) Nopassgroup(groupnode, fromnode, tonode string, domain *string) (err sys.ERROR) {
	return sys.ERR_DATABASE
}

func (this *nilprocess) Kickgroup(groupnode, fromnode, tonode string, domain *string) (err sys.ERROR) {
	return sys.ERR_DATABASE
}
func (this *nilprocess) Leavegroup(groupnode, fromnode string, domain *string) (err sys.ERROR) {
	return sys.ERR_DATABASE
}

func (this *nilprocess) Cancelgroup(groupnode, fromnode string, domain *string) (err sys.ERROR) {
	return sys.ERR_DATABASE
}

func (this *nilprocess) Blockgroup(groupnode, fromnode string, domain *string) (err sys.ERROR) {
	return sys.ERR_DATABASE
}
func (this *nilprocess) Blockgroupmember(groupnode, fromnode, tonodme string, domain *string) (err sys.ERROR) {
	return sys.ERR_DATABASE
}

func (this *nilprocess) ModifygroupInfo(node, fnode string, tu *TimRoomBean) (err sys.ERROR) {
	return sys.ERR_DATABASE
}
func (this *nilprocess) GetGroupInfo(node []string) (m map[string]*TimRoomBean, err sys.ERROR) {
	err = sys.ERR_DATABASE
	return
}
func (this *nilprocess) ModifyUserInfo(node string, tu *TimUserBean) (err sys.ERROR) {
	return sys.ERR_DATABASE
}
func (this *nilprocess) GetUserInfo(node []string) (m map[string]*TimUserBean, err sys.ERROR) {
	err = sys.ERR_DATABASE
	return
}
