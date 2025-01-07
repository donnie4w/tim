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

type service interface {
	Login(string, string, *string) (string, errs.ERROR)
	Register(string, string, *string) (string, errs.ERROR)
	AuthNode(string, string, *string) (string, errs.ERROR)
	Modify(uint64, *string, string, *string) errs.ERROR
	SaveMessage(*TimMessage) error
	SaveOfflineMessage(*TimMessage) error
	DelMessageByMid(uint64, int64) error
	DelOfflineMessage(uint64, ...int64) (int64, error)

	ExistGroup(string) bool
	ExistUser(string) bool
	AuthGroupAndUser(string, string, *string) (bool, error)
	AuthUserAndUser(string, string, *string) bool

	GetMessage(string, *string, int8, string, int64, int64) ([]*TimMessage, error)
	GetMessageByMid(uint64, int64) (*TimMessage, error)
	GetOfflineMessage(string, int) ([]*OfflineBean, error)

	Addroster(string, string, *string) (int8, errs.ERROR)
	Blockrosterlist(string) []string
	Rmroster(string, string, *string) (bool, bool)
	Roster(string) []string
	UserGroup(string, *string) []string
	Blockroster(string, string, *string) (bool, bool)
	ModifyUserInfo(string, *TimUserBean) errs.ERROR
	GetUserInfo([]string) (map[string]*TimUserBean, errs.ERROR)

	Newgroup(string, string, sys.TIMTYPE, *string) (string, errs.ERROR)
	Addgroup(string, string, *string) errs.ERROR
	Blockroomlist(string) []string
	Blockroommemberlist(string, string) []string
	GroupManagers(string, *string) ([]string, errs.ERROR)
	GroupRoster(string) []string
	GroupGtype(string, *string) (int8, errs.ERROR)
	Kickgroup(string, string, string, *string) errs.ERROR
	Leavegroup(string, string, *string) (err errs.ERROR)
	Nopassgroup(string, string, string, *string) errs.ERROR
	Pullgroup(string, string, string, *string) (bool, errs.ERROR)
	Blockgroup(string, string, *string) errs.ERROR
	Blockgroupmember(string, string, string, *string) errs.ERROR
	Cancelgroup(string, string, *string) errs.ERROR
	ModifygroupInfo(string, string, *TimRoomBean) errs.ERROR
	GetGroupInfo([]string) (map[string]*TimRoomBean, errs.ERROR)
	TimAdminAuth(account, password, domain string) bool
}
