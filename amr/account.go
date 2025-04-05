// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package amr

import (
	"github.com/donnie4w/gofer/util"
	"github.com/donnie4w/tim/cache"
	"github.com/donnie4w/tim/sys"
)

func AddAccount(node string) error {
	if islocalamr {
		return nil
	}
	return amr.append(ACCOUNT, util.Int64ToBytes(int64(util.FNVHash64([]byte(node)))), util.Int64ToBytes(sys.UUID), sys.Conf.TTL)
}

func RemoveAccount(node string) {
	if islocalamr {
		return
	}
	amr.removeKV(ACCOUNT, util.Int64ToBytes(int64(util.FNVHash64([]byte(node)))), util.Int64ToBytes(sys.UUID))
	cache.AccountCache.Del(node)
}

func GetAccount(node string) (r []int64) {
	if islocalamr {
		return []int64{sys.UUID}
	}
	if r = cache.AccountCache.Get(node); len(r) > 0 {
		return
	}
	if vs := amr.getMutil(ACCOUNT, util.Int64ToBytes(int64(util.FNVHash64([]byte(node))))); len(vs) > 0 {
		if len(vs) > 1 {
			r = make([]int64, 0, len(vs))
			m := make(map[int64]bool)
			for _, v := range vs {
				k := util.BytesToInt64(v)
				if _, b := m[k]; !b {
					r = append(r, k)
					m[k] = true
				}
			}
		} else {
			r = []int64{util.BytesToInt64(vs[0])}
		}
		cache.AccountCache.Put(node, r)
	}
	return
}
