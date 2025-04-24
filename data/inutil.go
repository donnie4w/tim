// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of t source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package data

import (
	goutil "github.com/donnie4w/gofer/util"
	"github.com/donnie4w/tim/dao"
	"github.com/donnie4w/tim/stub"
)

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

func checkuseruuidGdao(uuids ...uint64) bool {
	for _, uuid := range uuids {
		idbs := goutil.Int64ToBytes(int64(uuid))
		if uuidCache.Contains(idbs) {
			continue
		}
		tu := newTimuser(uuid)
		tu.Where(tu.UUID.EQ(int64(uuid)))
		if tu, _ := tu.Select(tu.ID); tu != nil && tu.GetId() > 0 {
			uuidCache.Add(idbs)
			continue
		} else {
			return false
		}
	}
	return true
}

func checkgroupuuid(uuids ...uint64) bool {
	for _, uuid := range uuids {
		idbs := goutil.Int64ToBytes(int64(uuid))
		if uuidCache.Contains(idbs) {
			continue
		}
		tu := newTimgroup(uuid)
		tu.Where(tu.UUID.EQ(int64(uuid)))
		if tu, _ := tu.Select(tu.ID); tu != nil && tu.GetId() > 0 {
			uuidCache.Add(idbs)
			continue
		} else {
			return false
		}
	}
	return true
}

func defaultInlineDB() *stub.InlineDB {
	return &stub.InlineDB{SQLITE: &stub.Connect{DBname: "tim.db"}}
}
