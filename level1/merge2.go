// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package level1

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/donnie4w/tim/util"
)

func syncMerge(tlcontext *tlContext) {
	defer util.Recover()
	defer tlcontext.mergemux.Unlock()
	_syncMerge(tlcontext)
}

func _syncMerge(tlcontext *tlContext) {
	m := make(map[int64]int8, 0)
	for sb := range tlcontext.mergeChan {
		m[sb.SyncId] = sb.Result
		if atomic.AddInt64(&tlcontext.mergeCount, -1) <= 0 {
			break
		}
	}
	syncTxMerge(m, tlcontext.remoteUuid)
	if tlcontext.mergeCount > 0 {
		_syncMerge(tlcontext)
	}
}

func syncTxMerge(syncList map[int64]int8, uuid int64) (err error) {
	i := 10
	for i > 0 {
		i--
		if tc := nodeWare.GetTlContext(uuid); tc != nil {
			if err = tc.csnet.SyncTxMerge(context.Background(), syncList); err == nil {
				break
			}
		} else {
			<-time.After(1 * time.Second)
		}
	}
	return
}
