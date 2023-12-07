// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim
//

package data

import (
	. "github.com/donnie4w/tim/stub"
	"github.com/donnie4w/tim/sys"
)

func (this *sqlhandler) UserGroup(node string, domain *string) []string {
	return sqlHandle.userGroup(node)
}

func (this *sqlhandler) GroupRoster(groupnode string) []string {
	return sqlHandle.groupRoster(groupnode)
}

func (this *sqlhandler) Roster(node string) []string {
	return sqlHandle.roster(node)
}

func (this *sqlhandler) Blockrosterlist(string) (_r []string) {
	return
}

func (this *sqlhandler) Blockroomlist(string) (_r []string) {
	return
}

func (this *sqlhandler) Blockroommemberlist(string, string) (_r []string) {
	return
}

func (this *sqlhandler) AuthUserAndUser(fnode, tnode string, domain *string) (_r bool) {
	_r, _ = sqlHandle.authUser(fnode, tnode)
	return
}

func (this *sqlhandler) ExistUser(node string) (_r bool) {
	_r, _ = sqlHandle.existUser(node)
	return
}

func (this *sqlhandler) ExistGroup(node string) (_r bool) {
	_r, _ = sqlHandle.existGroup(node)
	return
}

func (this *sqlhandler) Addroster(fnode, tnode string, domain *string) (status int8, err sys.ERROR) {
	err = sys.ERR_INTERFACE
	return
}

func (this *sqlhandler) Rmroster(fnode, tnode string, domain *string) (mstell bool, ok bool) {
	return
}

func (this *sqlhandler) Blockroster(fnode, tnode string, domain *string) (mstell bool, ok bool) {
	return
}

func (this *sqlhandler) GroupGtype(groupnode string, domain *string) (gtype int8, err sys.ERROR) {
	err = sys.ERR_INTERFACE
	return
}

func (this *sqlhandler) GroupManagers(groupnode string, domain *string) (s []string, err sys.ERROR) {
	err = sys.ERR_INTERFACE
	return
}

func (this *sqlhandler) Newgroup(fnode, groupname string, gtype int8, domain *string) (gnode string, err sys.ERROR) {
	err = sys.ERR_INTERFACE
	return
}

func (this *sqlhandler) Addgroup(groupnode, fromnode string, domain *string) (err sys.ERROR) {
	err = sys.ERR_INTERFACE
	return
}

func (this *sqlhandler) Pullgroup(groupnode, fromnode, tonode string, domain *string) (isReq bool, err sys.ERROR) {
	err = sys.ERR_INTERFACE
	return
}

func (this *sqlhandler) Nopassgroup(groupnode, fromnode, tonode string, domain *string) (err sys.ERROR) {
	err = sys.ERR_INTERFACE
	return
}

func (this *sqlhandler) Kickgroup(groupnode, fromnode, tonode string, domain *string) (err sys.ERROR) {
	err = sys.ERR_INTERFACE
	return
}
func (this *sqlhandler) Leavegroup(groupnode, fromnode string, domain *string) (err sys.ERROR) {
	err = sys.ERR_INTERFACE
	return
}

func (this *sqlhandler) Cancelgroup(groupnode, fromnode string, domain *string) (err sys.ERROR) {
	err = sys.ERR_INTERFACE
	return
}

func (this *sqlhandler) Blockgroup(groupnode, fromnode string, domain *string) (err sys.ERROR) {
	err = sys.ERR_INTERFACE
	return
}

func (this *sqlhandler) Blockgroupmember(groupnode, fromnode, tonodme string, domain *string) (err sys.ERROR) {
	err = sys.ERR_INTERFACE
	return
}

func (this *sqlhandler) ModifyUserInfo(node string, tu *TimUserBean) (err sys.ERROR) {
	return
}
func (this *sqlhandler) GetUserInfo(node []string) (m map[string]*TimUserBean, err sys.ERROR) {
	return
}

func (this *sqlhandler) ModifygroupInfo(node, fnode string, tu *TimRoomBean) (err sys.ERROR) {
	return
}
func (this *sqlhandler) GetGroupInfo(node []string) (m map[string]*TimRoomBean, err sys.ERROR) {
	return
}
