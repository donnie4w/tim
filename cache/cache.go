// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package cache

import (
	gocache "github.com/donnie4w/gofer/cache"
	"github.com/donnie4w/gofer/hashmap"
)

var TokenCache = newTokenPool()

var AuthCache = newAuth()

var AccountCache = newAccountPool()

var TokenUsedCache = gocache.NewBloomFilter(1<<19, 0.01)

var BlockUserCache = newBlockPool()

var VnodeCache = hashmap.NewLimitHashMap[uint64, int64](1 << 18)

var CsAccessCache = newCsAccessPool()
