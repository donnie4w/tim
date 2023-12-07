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

type engine interface {
	Login(string, string, *string) (string, sys.ERROR)
	Register(string, string, *string) (string, sys.ERROR)
	Token(string, string, *string) (string, sys.ERROR)
	Modify(uint64, *string, string, *string) sys.ERROR
	SaveMessage(*TimMessage) error
	SaveOfflineMessage(*TimMessage) error
	DelMessageByMid(uint64, int64) error
	DelOfflineMessage(uint64, ...int64) (int64, error)

	ExistGroup(string) bool
	ExistUser(string) bool
	AuthGroupAndUser(string, string, *string) (bool, error)
	AuthUserAndUser(string, string, *string) bool

	GetMessage(*Tid, int8, string, int64, int64) ([]*TimMessage, error)
	GetMessageByMid(uint64, int64) (*TimMessage, error)
	GetOfflineMessage(*Tid, int) ([]*OfflineBean, error)

	Addroster(string, string, *string) (int8, sys.ERROR)
	Blockrosterlist(string) []string
	Rmroster(string, string, *string) (bool, bool)
	Roster(string) []string
	UserGroup(string, *string) []string
	Blockroster(string, string, *string) (bool, bool)
	ModifyUserInfo(string, *TimUserBean) sys.ERROR
	GetUserInfo([]string) (map[string]*TimUserBean, sys.ERROR)

	Newgroup(string, string, int8, *string) (string, sys.ERROR)
	Addgroup(string, string, *string) sys.ERROR
	Blockroomlist(string) []string
	Blockroommemberlist(string, string) []string
	GroupManagers(string, *string) ([]string, sys.ERROR)
	GroupRoster(string) []string
	GroupGtype(string, *string) (int8, sys.ERROR)
	Kickgroup(string, string, string, *string) sys.ERROR
	Leavegroup(string, string, *string) (err sys.ERROR)
	Nopassgroup(string, string, string, *string) sys.ERROR
	Pullgroup(string, string, string, *string) (bool, sys.ERROR)
	Blockgroup(string, string, *string) sys.ERROR
	Blockgroupmember(string, string, string, *string) sys.ERROR
	Cancelgroup(string, string, *string) sys.ERROR
	ModifygroupInfo(string, string, *TimRoomBean) sys.ERROR
	GetGroupInfo([]string) (map[string]*TimRoomBean, sys.ERROR)
}
