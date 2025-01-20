// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package cache

import (
	"github.com/donnie4w/gofer/hashmap"
	"time"
)

type blockPool struct {
	mm *hashmap.Map[string, int64]
}

func newBlockPool() *blockPool {
	bp := &blockPool{mm: hashmap.NewMap[string, int64]()}
	go bp.tick()
	return bp
}

func (bp *blockPool) Put(key string, value int64) {
	bp.mm.Put(key, value+time.Now().Unix())
}

func (bp *blockPool) Get(key string) int64 {
	r, _ := bp.mm.Get(key)
	return r
}

func (bp *blockPool) Del(key string) {
	bp.mm.Del(key)
}

func (bp *blockPool) tick() int {
	tk := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-tk.C:
			bp.mm.Range(func(k string, v int64) bool {
				if v < time.Now().Unix() {
					bp.Del(k)
				}
				return true
			})
		}
	}
}
