// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package data

import "github.com/donnie4w/tim/dao"

func newTimuser(uuid uint64) *dao.Timuser {
	tu := dao.NewTimuser()
	tu.UseDBHandle(gdaoHandle.GetDBHandle(uuid))
	return tu
}

func newTimmessage(chatId uint64) *dao.Timmessage {
	tm := dao.NewTimmessage()
	tm.UseDBHandle(gdaoHandle.GetDBHandle(chatId))
	return tm
}

func newTimblock(uuid uint64) *dao.Timblock {
	tm := dao.NewTimblock()
	tm.UseDBHandle(gdaoHandle.GetDBHandle(uuid))
	return tm
}

func newTimblockroom(uuid uint64) *dao.Timblockroom {
	tm := dao.NewTimblockroom()
	tm.UseDBHandle(gdaoHandle.GetDBHandle(uuid))
	return tm
}

func newTimoffline(uuid uint64) *dao.Timoffline {
	tm := dao.NewTimoffline()
	tm.UseDBHandle(gdaoHandle.GetDBHandle(uuid))
	return tm
}

func newTimgroup(uuid uint64) *dao.Timgroup {
	tm := dao.NewTimgroup()
	tm.UseDBHandle(gdaoHandle.GetDBHandle(uuid))
	return tm
}

func newTimrelate(uuid uint64) *dao.Timrelate {
	tm := dao.NewTimrelate()
	tm.UseDBHandle(gdaoHandle.GetDBHandle(uuid))
	return tm
}

func newTimroster(uuid uint64) *dao.Timroster {
	tm := dao.NewTimroster()
	tm.UseDBHandle(gdaoHandle.GetDBHandle(uuid))
	return tm
}

func newTimmucroster(uuid uint64) *dao.Timmucroster {
	tm := dao.NewTimmucroster()
	tm.UseDBHandle(gdaoHandle.GetDBHandle(uuid))
	return tm
}

func newTimdomain() *dao.Timdomain {
	tm := dao.NewTimdomain()
	tm.UseDBHandle(gdaoHandle.GetDBHandle(0))
	return tm
}
