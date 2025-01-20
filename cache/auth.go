// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package cache

import (
	"github.com/donnie4w/gofer/hashmap"
	"github.com/donnie4w/gofer/lock"
	goutil "github.com/donnie4w/gofer/util"
	"github.com/donnie4w/tim/data"
	"github.com/donnie4w/tim/sys"
	"github.com/donnie4w/tim/util"
	"time"
)

type auth struct {
	mm      *hashmap.LimitHashMap[uint64, int64]
	strlock *lock.Strlock
}

func newAuth() *auth {
	return &auth{mm: hashmap.NewLimitHashMapWithSegment[uint64, int64](1<<20, 1<<10), strlock: lock.NewStrlock(1 << 10)}
}

func (a *auth) Has(fnode, tnode string, domain *string, group bool) bool {
	var idx []byte
	if group {
		idx = util.RelateIdForGroup(fnode, tnode, domain)
	} else {
		idx = util.ChatIdByNode(fnode, tnode, domain)
	}
	if r, _ := a.mm.Get(goutil.FNVHash64(idx)); r > 0 {
		if v := time.Now().Unix() - r; v < sys.Conf.CacheAuthExpire {
			if v < sys.Conf.CacheAuthExpire/2 {
				go a.check(fnode, tnode, domain, group)
			}
			return true
		}
	}
	return false
}

func (a *auth) Put(fnode, tnode string, domain *string, group bool) {
	if group {
		rid := util.RelateIdForGroup(fnode, tnode, domain)
		a.mm.Put(goutil.FNVHash64(rid), time.Now().Unix())
	} else {
		cid := util.ChatIdByNode(fnode, tnode, domain)
		a.mm.Put(goutil.FNVHash64(cid), time.Now().Unix())
	}
}

func (a *auth) Delete(key uint64) {
	a.mm.Del(key)
}

func (a *auth) DeleteWithNode(fnode, tnode string, domain *string, group bool) {
	if group {
		a.Delete(goutil.FNVHash64(util.RelateIdForGroup(fnode, tnode, domain)))
	} else {
		a.Delete(goutil.FNVHash64(util.ChatIdByNode(fnode, tnode, domain)))
	}
}

func (a *auth) check(fnode, tnode string, domain *string, group bool) {
	if lock, ok := a.strlock.TryLock(fnode + tnode); ok {
		defer lock.Unlock()
		var b bool
		if group {
			b, _ = data.Service.AuthGroupAndUser(fnode, tnode, domain)
		} else {
			b = data.Service.AuthUserAndUser(fnode, tnode, domain)
		}
		if b {
			a.Put(fnode, tnode, domain, group)
		} else {
			a.DeleteWithNode(fnode, tnode, domain, group)
		}
	}
}
