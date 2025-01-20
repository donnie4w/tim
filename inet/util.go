// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package inet

import (
	"github.com/donnie4w/gofer/httpclient"
	"github.com/donnie4w/tim/errs"
	"sync/atomic"
	"time"

	"github.com/donnie4w/gofer/buffer"
	"github.com/donnie4w/gofer/hashmap"
	goutil "github.com/donnie4w/gofer/util"
	"github.com/donnie4w/tim/stub"
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

var hzmap = hashmap.NewLimitMap[int64, []int64](1 << 17)

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

func newAuth(bs []byte) (ta *stub.TimAuth) {
	if util.JTP(bs[0]) {
		ta, _ = goutil.JsonDecode[*stub.TimAuth](bs[1:])
	} else {
		ta, _ = goutil.TDecode(bs[1:], &stub.TimAuth{})
	}
	return
}

func connectAuth(bs []byte) (_r bool) {
	defer util.Recover()
	if sys.Conf.Security != nil && sys.Conf.Security.ConnectAuthUrl != nil {
		if ta := newAuth(bs); ta != nil {
			if bs, err := httpclient.Post2(goutil.JsonEncode(ta), true, *sys.Conf.Security.ConnectAuthUrl); err == nil && bs != nil {
				_r = bs[0] == 1
			}
		}
	} else {
		return true
	}
	return
}

var bigMap = hashmap.NewMap[int64, *bm]()

type bm struct {
	buf    *buffer.Buffer
	length int
}

func (this *bm) addData(bs []byte) (_bs []byte, ok bool) {
	if l := this.length - this.buf.Len() - len(bs); l > 0 {
		this.buf.Write(bs)
		return nil, false
	} else if l == 0 {
		this.buf.Write(bs)
		return this.buf.Bytes(), true
	} else if l < 0 {
		return bs, true
	}
	return
}

func addBigData(hc *tlnet.HttpContext, bs []byte) {
	if m, ok := bigMap.Get(hc.WS.Id); ok {
		if _bs, b := m.addData(bs); b {
			bigMap.Del(hc.WS.Id)
			parseWsData(_bs, hc)
		}
	}
}

func newBigData(hc *tlnet.HttpContext, bs []byte, len int) {
	bigMap.Put(hc.WS.Id, &bm{buf: buffer.NewBufferBySlice(bs), length: len})
}

func parseBigData(hc *tlnet.HttpContext, bs []byte) {
	if len(bs) < 6 {
		return
	}
	length := int(goutil.BytesToInt32(bs[1:5]))
	if fg := length - len(bs); fg == 0 {
		parseWsData(bs, hc)
	} else if fg > 0 {
		newBigData(hc, bs, length)
	} else {
		go sys.SendWs(hc.WS.Id, &stub.TimAck{Ok: false, TimType: int8(sys.TIMBIGBINARY), Error: errs.ERR_BIGDATA.TimError()}, sys.TIMACK)
	}
}

func rmBigDataId(id int64) {
	bigMap.Del(id)
}

func isBigData(id int64) bool {
	return bigMap.Has(id)
}

func isForBitIface(b byte) bool {
	if sys.BlockApiMap != nil {
		return sys.BlockApiMap.Has(sys.TIMTYPE(b & 0x7f))
	}
	return false
}
