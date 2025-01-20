// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package cache

import (
	"github.com/donnie4w/gofer/hashmap"
	"github.com/donnie4w/tim/util"
	"time"
)

type csAccessBean struct {
	csaddr    string
	timestamp int64
}

type csAccessPool struct {
	mm *hashmap.LinkedHashMap[int64, *csAccessBean]
}

func newCsAccessPool() *csAccessPool {
	tp := &csAccessPool{mm: hashmap.NewLinkedHashMap[int64, *csAccessBean](1 << 10)}
	go tp.ticker()
	return tp
}

func (t *csAccessPool) Put(uuid int64, csAddr string) {
	t.mm.Put(uuid, &csAccessBean{csaddr: csAddr, timestamp: time.Now().UnixNano()})
}

func (t *csAccessPool) Get(uuid int64) string {
	if r, _ := t.mm.Get(uuid); r != nil {
		return r.csaddr
	}
	return ""
}

func (t *csAccessPool) Del(uuid int64) {
	t.mm.Delete(uuid)
}

func (t *csAccessPool) ticker() {
	tk := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-tk.C:
			func() {
				defer util.Recover()
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
