// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package cache

import (
	"github.com/donnie4w/gofer/hashmap"
	"github.com/donnie4w/gofer/util"
	"time"
)

type accountBean struct {
	uuids     []int64
	timestamp int64
}

type accountPool struct {
	mm *hashmap.LinkedHashMap[int64, *accountBean]
}

func newAccountPool() *accountPool {
	ap := &accountPool{mm: hashmap.NewLinkedHashMap[int64, *accountBean](1 << 17)}
	go ap.ticker()
	return ap
}

func (t *accountPool) Put(node string, uuids []int64) {
	t.mm.Put(int64(util.FNVHash64([]byte(node))), &accountBean{uuids: uuids, timestamp: time.Now().UnixNano()})
}

func (t *accountPool) Get(node string) []int64 {
	if r, _ := t.mm.Get(int64(util.FNVHash64([]byte(node)))); r != nil {
		return r.uuids
	}
	return nil
}

func (t *accountPool) Del(node string) {
	t.mm.Delete(int64(util.FNVHash64([]byte(node))))
}

func (t *accountPool) ticker() {
	tk := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-tk.C:
			func() {
				defer util.Recover(nil)
				iterator := t.mm.Iterator(false)
				for {
					if k, v, b := iterator.Next(); b {
						if time.Now().UnixNano() > v.timestamp+time.Minute.Nanoseconds() {
							t.mm.Delete(k)
							continue
						}
					}
					break
				}
			}()
		}
	}
}
