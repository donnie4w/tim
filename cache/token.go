// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package cache

import (
	"github.com/donnie4w/gofer/hashmap"
	"github.com/donnie4w/tim/stub"
	"github.com/donnie4w/tim/sys"
	"github.com/donnie4w/tim/util"
	"time"
)

type tokenBean struct {
	tid       *stub.Tid
	timestamp int64
}

type tokenPool struct {
	mm *hashmap.LinkedHashMap[string, *tokenBean]
}

func newTokenPool() *tokenPool {
	tp := &tokenPool{mm: hashmap.NewLinkedHashMap[string, *tokenBean](1 << 16)}
	go tp.ticker()
	return tp
}

func (t *tokenPool) Put(token string, tid *stub.Tid) {
	t.mm.Put(token, &tokenBean{tid: tid, timestamp: time.Now().UnixNano()})
}

func (t *tokenPool) Get(token string) *stub.Tid {
	if r, _ := t.mm.Get(token); r != nil {
		return r.tid
	}
	return nil
}

func (t *tokenPool) Del(token string) {
	t.mm.Delete(token)
}

func (t *tokenPool) ticker() {
	tk := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-tk.C:
			func() {
				defer util.Recover()
				iterator := t.mm.Iterator(false)
				for {
					if k, v, b := iterator.Next(); b {
						if time.Now().UnixNano() > v.timestamp+sys.Conf.TokenTimeout {
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
