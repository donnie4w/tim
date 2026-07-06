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
	"github.com/donnie4w/tim/sys"
)

var (
	TokenCache     *tokenPool
	AuthCache      *auth
	AccountCache   *accountPool
	BlockUserCache *blockPool
	CsAccessCache  *csAccessPool
	TokenUsedCache *gocache.BloomFilter
	VnodeCache     *hashmap.LimitFifoMap[uint64, int64]
)

func init() {
	sys.Service(sys.INIT_CACHE, serv(1))
}

type serv byte

func (serv) Serve() error {
	if sys.Conf.Memlimit >= 1<<10 {
		TokenUsedCache = gocache.NewBloomFilter(1<<19, 0.01)
		VnodeCache = hashmap.NewLimitFifoMap[uint64, int64](1 << 18)
	} else {
		TokenUsedCache = gocache.NewBloomFilter(1<<13, 0.01)
		VnodeCache = hashmap.NewLimitFifoMap[uint64, int64](1 << 13)
	}
	TokenCache = newTokenPool()
	AuthCache = newAuth()
	AccountCache = newAccountPool()
	BlockUserCache = newBlockPool()
	CsAccessCache = newCsAccessPool()
	return nil
}

func (serv) Close() error {
	return nil
}
