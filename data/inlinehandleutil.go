// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of t source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package data

func checkuseruuid(uuids ...uint64) bool {
	for _, uuid := range uuids {
		tu := newTimuser(uuid)
		tu.Where(tu.UUID.EQ(int64(uuid)))
		if tu, _ := tu.Select(tu.ID); tu != nil && tu.GetId() > 0 {
			continue
		} else {
			return false
		}
	}
	return true
}

func checkgroupuuid(uuids ...uint64) bool {
	for _, uuid := range uuids {
		tu := newTimgroup(uuid)
		tu.Where(tu.UUID.EQ(int64(uuid)))
		if tu, _ := tu.Select(tu.ID); tu != nil && tu.GetId() > 0 {
			continue
		} else {
			return false
		}
	}
	return true
}
