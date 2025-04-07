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
	"github.com/donnie4w/tim/sys"
)

func GetCsAccess(uuid int64) (r string, err error) {
	if r = cache.CsAccessCache.Get(uuid); r != "" {
		return
	}
	if bs, err := amr.get(UUID, util.Int64ToBytes(uuid)); err == nil && len(bs) > 0 {
		r = string(bs)
		cache.CsAccessCache.Put(uuid, r)
	} else {
		return r, err
	}
	return
}

func PutCsAccess() (err error) {
	defer util.Recover(&err)
	return amr.put(UUID, util.Int64ToBytes(sys.UUID), []byte(sys.Conf.CsAccess), uint64(sys.UUIDCSTIME))
}
