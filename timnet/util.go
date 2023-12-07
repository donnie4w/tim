// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim
//

package timnet

import (
	"sync/atomic"
	"time"

	. "github.com/donnie4w/gofer/hashmap"
	"github.com/donnie4w/gofer/httputil"
	. "github.com/donnie4w/gofer/util"
	. "github.com/donnie4w/tim/stub"
	"github.com/donnie4w/tim/sys"
	"github.com/donnie4w/tim/util"
	"github.com/donnie4w/tlnet"
)

func isForBidRegister() bool {
	return sys.Conf.Security != nil && sys.Conf.Security.ForBidRegister
}

func isForBidToken() bool {
	return sys.Conf.Security != nil && sys.Conf.Security.ForBidToken
}

func reqHzSecond() int {
	if sys.Conf.Security != nil {
		return sys.Conf.Security.ReqHzSecond
	}
	return 0
}

var hzmap = NewLimitMap[int64, []int64](1 << 17)

func overHz(hc *tlnet.HttpContext) (b bool) {
	defer util.Recover()
	if reqHzSecond() > 0 {
		if t, ok := hzmap.Get(hc.WS.Id); ok {
			if time.Now().UnixNano()-t[0] < int64(int(time.Second/time.Nanosecond)/reqHzSecond()) {
				hzmap.Put(hc.WS.Id, []int64{time.Now().UnixNano(), atomic.AddInt64(&t[1], 1)})
				if b = t[1] > 5; b {
					hc.WS.Close()
				}
			}
		}
		hzmap.Put(hc.WS.Id, []int64{time.Now().UnixNano(), 0})
	}
	return
}

func overMaxData(ws *tlnet.Websocket, length int64) (_r bool) {
	if sys.Conf.Security != nil && sys.Conf.Security.MaxDatalimit > 0 {
		if sys.Conf.Security.MaxDatalimit > length {
			if ws != nil {
				ws.Close()
			}
			_r = true
		}
	}
	return
}

func newAuth(bs []byte) (ta *TimAuth) {
	if util.JTP(bs[0]) {
		ta, _ = JsonDecode[*TimAuth](bs[1:])
	} else {
		ta, _ = TDecode(bs[1:], &TimAuth{})
	}
	return
}

func connectAuth(bs []byte) (_r bool) {
	defer util.Recover()
	if sys.Conf.Security != nil && sys.Conf.Security.ConnectAuthUrl != nil {
		if ta := newAuth(bs); ta != nil {
			if bs, err := httputil.HttpPost(JsonEncode(ta), true, *sys.Conf.Security.ConnectAuthUrl); err == nil && bs != nil {
				_r = bs[0] == 1
			}
		}
	} else {
		return true
	}
	return
}
