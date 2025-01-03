// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package cache

import (
	"github.com/donnie4w/gofer/cache"
	"github.com/donnie4w/gofer/hashmap"
	"github.com/donnie4w/tim/stub"
)

var TokenCache = hashmap.NewLimitHashMap[int64, *stub.Tid](1 << 13)

var AuthCache = hashmap.NewLimitHashMap[uint64, int64](1 << 14)

var AmrCache = cache.NewLruCache[uint64](1 << 20)
