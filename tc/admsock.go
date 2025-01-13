// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package tc

import (
	"github.com/donnie4w/tim/errs"
	"sync/atomic"
	"time"

	. "github.com/donnie4w/gofer/buffer"
	. "github.com/donnie4w/gofer/hashmap"
	. "github.com/donnie4w/gofer/lock"
	. "github.com/donnie4w/gofer/util"
	"github.com/donnie4w/gothrift/thrift"
	"github.com/donnie4w/tim/sys"
	"github.com/donnie4w/tim/util"
	"github.com/donnie4w/tlnet"
)

var numLock = NewNumLock(64)
var await = NewAwait[int8](1 << 10)

var admwsware = &wsware{wsmap: NewMapL[int64, *WsSock]()}

type wsware struct {
	wsmap *MapL[int64, *WsSock]
}

func (t *wsware) Addws(ws *tlnet.Websocket) {
	numLock.Lock(ws.Id)
	defer numLock.Unlock(ws.Id)
	if t.wsmap.Has(ws.Id) {
		return
	}
	wss := NewWsSock(ws)
	t.wsmap.Put(ws.Id, wss)
}

func (t *wsware) SendWs(id int64, ts thrift.TStruct, tt sys.TIMTYPE) (_r bool) {
	if wss, ok := t.wsmap.Get(id); ok {
		if err := wss.send(ts, tt, false); err == nil {
			_r = ok
		}
	}
	return
}

func (t *wsware) SendWsWithAck(id int64, ts thrift.TStruct, tt sys.TIMTYPE) (_r bool) {
	if wss, ok := t.wsmap.Get(id); ok {
		if err := wss.send(ts, tt, true); err == nil {
			_r = ok
		}
	}
	return
}

func (t *wsware) Ping(id int64) {
	if wss, ok := t.wsmap.Get(id); ok {
		wss.send(nil, sys.ADMPING, false)
	}
}

func (t *wsware) Get(ws *tlnet.Websocket) (*WsSock, bool) {
	return t.wsmap.Get(ws.Id)
}

func (t *wsware) hasws(ws *tlnet.Websocket) bool {
	return t.wsmap.Has(ws.Id)
}

func (t *wsware) wsById(id int64) (*tlnet.Websocket, bool) {
	if ws, b := t.wsmap.Get(id); b {
		return ws.ws, b
	}
	return nil, false
}

func (t *wsware) wsLen() int64 {
	return t.wsmap.Len()
}

func (t *wsware) delws(ws *tlnet.Websocket) {
	defer util.Recover()
	t.delId(ws.Id)
}

func (t *wsware) delId(id int64) {
	if sk, ok := t.wsmap.Get(id); ok {
		sk.ws.Close()
		t.wsmap.Del(id)
	}
}

type WsSock struct {
	ws *tlnet.Websocket
}

func NewWsSock(ws *tlnet.Websocket) (_r *WsSock) {
	_r = &WsSock{ws: ws}
	return
}

func (t *WsSock) _send(buf *Buffer) (err error) {
	sys.Stat.Ob(int64(buf.Len()))
	return t.ws.Send(buf.Bytes())
}

var seq int32 = 0

func (t *WsSock) send(ts thrift.TStruct, tt sys.TIMTYPE, sync bool) (err error) {
	length := 1
	if sync {
		length = 5
	}
	var bs []byte
	if ts != nil {
		bs = TEncode(ts)
		length += len(bs)
	}
	buf := NewBufferWithCapacity(length)
	sendId := atomic.AddInt32(&seq, 1)
	var ch chan int8
	if sync {
		buf.WriteByte(byte(tt) | 0x80)
		buf.Write(Int32ToBytes(sendId))
		ch = await.Get(int64(sendId))
	} else {
		buf.WriteByte(byte(tt))
	}
	//if ts != nil {
	//	buf.Write(TEncode(ts))
	//}
	if len(bs) > 0 {
		buf.Write(bs)
	}
	if err = t._send(buf); err == nil && sync {
		i := 0
		for t.ws.Error == nil && i < 100 {
			i++
			select {
			case <-ch:
				err = nil
				goto END
			case <-time.After(time.Second):
				err = errs.ERR_OVERTIME.Error()
			}
		}
		if t.ws.Error != nil || err != nil {
			err = errs.ERR_OVERTIME.Error()
		}
	}
END:
	return
}

//func (t *WsSock) sendBigData(data []byte, tt sys.TIMTYPE) (err error) {
//	buf := NewBuffer()
//	buf.WriteByte(byte(tt))
//	buf.Write(data)
//	return t._send(buf)
//}

func (t *WsSock) close() {
	admwsware.delws(t.ws)
}

func awaitEnd(bs []byte) {
	await.Close(int64(BytesToInt32(bs)))
}
