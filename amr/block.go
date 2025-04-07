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
)

func PutBlock(node string, expried int64) {
	if islocalamr {
		cache.BlockUserCache.Put(node, expried)
	} else {
		amr.put(BLOCK, []byte(node), util.Int64ToBytes(expried), uint64(expried))
	}
}

func GetBlock(node string) int64 {
	if islocalamr {
		return cache.BlockUserCache.Get(node)
	} else {
		if bs, _ := amr.get(BLOCK, []byte(node)); len(bs) > 0 {
			return util.BytesToInt64(bs)
		}
	}
	return 0
}

func DelBlock(node string) {
	if islocalamr {
		cache.BlockUserCache.Del(node)
	} else {
		amr.remove(BLOCK, []byte(node))
	}
}
