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
	"github.com/donnie4w/tim/errs"
	"github.com/donnie4w/tim/sys"
)

func GetVnode(vnode string) int64 {
	if islocalamr {
		return sys.UUID
	} else {
		fh := util.FNVHash64([]byte(vnode))
		if v, b := cache.VnodeCache.Get(fh); b {
			return v
		}
		if bs, _ := amr.get(VNODE, []byte(vnode)); len(bs) > 0 {
			r := util.BytesToInt64(bs)
			cache.VnodeCache.Put(fh, r)
			return r
		}
	}
	return 0
}

func PutVnode(vnode string) error {
	if islocalamr {
		return nil
	}
	if vnode != "" {
		return amr.put(VNODE, []byte(vnode), util.Int64ToBytes(sys.UUID), 86400)
	} else {
		return errs.ERR_PARAMS.Error()
	}
}

func DelVnode(vnode string) {
	if vnode != "" {
		amr.remove(VNODE, []byte(vnode))
		cache.VnodeCache.Del(util.FNVHash64([]byte(vnode)))
	}
}
