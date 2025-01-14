// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package service

import (
	"github.com/donnie4w/tim/errs"
	"github.com/donnie4w/tim/mq"
	"sync/atomic"
	"time"

	. "github.com/donnie4w/gofer/buffer"
	. "github.com/donnie4w/gofer/hashmap"
	. "github.com/donnie4w/gofer/lock"
	. "github.com/donnie4w/gofer/util"
	"github.com/donnie4w/gothrift/thrift"
	. "github.com/donnie4w/tim/stub"
	"github.com/donnie4w/tim/sys"
	"github.com/donnie4w/tim/util"
	"github.com/donnie4w/tlnet"
)

var strLock = NewStrlock(64)
var numLock = NewNumLock(64)
var await = NewAwait[int8](1 << 10)

var wsware = &_wsWare{wsmap: NewMapL[int64, *WsSock](), uMap: NewMap[string, []int64]()}

type _wsWare struct {
	uMap  *Map[string, []int64]
	wsmap *MapL[int64, *WsSock]
}

func (t *_wsWare) AddTid(ws *tlnet.Websocket, tid *Tid) {
	numLock.Lock(ws.Id)
	defer numLock.Unlock(ws.Id)
	if t.wsmap.Has(ws.Id) {
		return
	}
	wss := NewWsSock(ws)
	wss.tid = tid
	t.wsmap.Put(ws.Id, wss)
	if a, ok := t.uMap.Get(tid.Node); !ok {
		t.uMap.Put(tid.Node, []int64{ws.Id})
	} else {
		t.uMap.Put(tid.Node, append(a, ws.Id))
	}
	go sys.Csuser(tid.Node, true, ws.Id)
	//go loginstat(tid.Node, true, tid, ws.Id)
	go mq.PushOnline(tid.Node, true)
}

func (t *_wsWare) SetJsonOn(ws *tlnet.Websocket) {
	if wss, ok := t.wsmap.Get(ws.Id); ok {
		wss.SetJsonOn(true)
	}
}

func (t *_wsWare) SendNode(node string, ts thrift.TStruct, tt sys.TIMTYPE) (_r bool) {
	if ids, ok := t.uMap.Get(node); ok {
		for _, id := range ids {
			if t.SendWs(id, ts, tt) {
				_r = true
			}
		}
	}
	return
}

func (t *_wsWare) SendBigData(node string, data []byte, tt sys.TIMTYPE) (_r bool) {
	if ids, ok := t.uMap.Get(node); ok {
		for _, id := range ids {
			if t.SendBigDataByWs(id, data, tt) {
				_r = true
			}
		}
	}
	return
}

func (t *_wsWare) SendBigDataByWs(id int64, data []byte, tt sys.TIMTYPE) (_r bool) {
	if wss, ok := t.wsmap.Get(id); ok {
		if err := wss.sendBigData(data, tt); err == nil {
			_r = ok
		}
	}
	return
}

func (t *_wsWare) SendNodeWithAck(node string, ts thrift.TStruct, tt sys.TIMTYPE) (_r bool) {
	if ids, ok := t.uMap.Get(node); ok {
		for _, id := range ids {
			if t.SendWsWithAck(id, ts, tt) {
				_r = true
			}
		}
	}
	return
}

func (t *_wsWare) SendWs(id int64, ts thrift.TStruct, tt sys.TIMTYPE) (_r bool) {
	if wss, ok := t.wsmap.Get(id); ok {
		if err := wss.send(ts, tt, false); err == nil {
			_r = ok
		}
	}
	return
}

func (t *_wsWare) SendWsWithAck(id int64, ts thrift.TStruct, tt sys.TIMTYPE) (_r bool) {
	if wss, ok := t.wsmap.Get(id); ok {
		if err := wss.send(ts, tt, true); err == nil {
			_r = ok
		}
	}
	return
}

func (t *_wsWare) GetUserDeviceLen(node string) (_r int) {
	if as, ok := t.uMap.Get(node); ok {
		_r = len(as)
	}
	return
}

func (t *_wsWare) GetUserDeviceTypeLen(node string) (_r []byte) {
	_r = make([]byte, 0)
	if as, ok := t.uMap.Get(node); ok {
		for _, v := range as {
			if tid, ok := t.wsmap.Get(v); ok {
				if tid.tid.Termtyp != nil {
					_r = append(_r, byte(*tid.tid.Termtyp))
				}
			}
		}
	}
	return
}

func (t *_wsWare) Ping(id int64) {
	if wss, ok := t.wsmap.Get(id); ok {
		wss.send(nil, sys.TIMPING, false)
		wss.pingt = sys.InaccurateTime
	}
}

func (t *_wsWare) detect(node string) {
	if ids, ok := t.uMap.Get(node); ok {
		for _, id := range ids {
			t.Ping(id)
		}
	}
}

func (t *_wsWare) Get(ws *tlnet.Websocket) (*WsSock, bool) {
	return t.wsmap.Get(ws.Id)
}

func (t *_wsWare) hasws(ws *tlnet.Websocket) bool {
	return t.wsmap.Has(ws.Id)
}

func (t *_wsWare) hasUser(node string) bool {
	return t.uMap.Has(node)
}

func (t *_wsWare) wsById(id int64) (*tlnet.Websocket, bool) {
	if ws, b := t.wsmap.Get(id); b {
		return ws.ws, b
	}
	return nil, false
}

func (t *_wsWare) wsLen() int64 {
	return t.wsmap.Len()
}

func (t *_wsWare) wssList(index, limit int64) (_r []*Tid, length int64) {
	_r = make([]*Tid, 0)
	count := index
	length = t.wsmap.Len()
	t.wsmap.Range(func(_ int64, wss *WsSock) bool {
		if limit == 0 {
			_r = append(_r, wss.tid)
		} else if count < limit {
			_r = append(_r, wss.tid)
			count++
		}
		return true
	})
	return
}

func (t *_wsWare) wssInfo(node string) (_r []byte) {
	_r = make([]byte, 0)
	if as, ok := t.uMap.Get(node); ok {
		for _, v := range as {
			if ws, ok := t.wsmap.Get(v); ok {
				if ws.tid.Termtyp != nil {
					_r = append(_r, byte(*ws.tid.Termtyp))
				} else {
					_r = append(_r, 0)
				}
			}
		}
	}
	return
}

func (t *_wsWare) delws(ws *tlnet.Websocket) {
	defer util.Recover()
	t.delId(ws.Id)
}

func (t *_wsWare) delId(id int64) {
	if sk, ok := t.wsmap.Get(id); ok {
		sk.ws.Close()
		t.wsmap.Del(id)
		if a, ok := t.uMap.Get(sk.tid.Node); ok {
			strLock.Lock(sk.tid.Node)
			defer strLock.Unlock(sk.tid.Node)
			if _a := util.ArraySub2(a, id); len(_a) > 0 {
				t.uMap.Put(sk.tid.Node, _a)
			} else {
				t.uMap.Del(sk.tid.Node)
				go sys.Csuser(sk.tid.Node, false, id)
				go sys.Interrupt(sk.tid)
				go mq.PushOnline(sk.tid.Node, false)
			}
		}
	}
}

func (t *_wsWare) delnode(node string) {
	if as, ok := t.uMap.Get(node); ok {
		for _, v := range as {
			t.delId(v)
		}
	}
}

type WsSock struct {
	ws     *tlnet.Websocket
	tid    *Tid
	jsonOn bool
	pingt  int64
}

func NewWsSock(ws *tlnet.Websocket) (_r *WsSock) {
	_r = &WsSock{ws: ws, pingt: sys.InaccurateTime}
	return
}

func (t *WsSock) SetJsonOn(on bool) {
	t.jsonOn = on
}

func (t *WsSock) _send(buf *Buffer) (err error) {
	if t.pingt+int64(sys.PINGTO*int64(time.Second)) < sys.InaccurateTime {
		t.close()
		return errs.ERR_PING.Error()
	}
	sys.Stat.Ob(int64(buf.Len()))
	return t.ws.Send(buf.Bytes())
}

var seq int32 = 0

func (t *WsSock) send(ts thrift.TStruct, tt sys.TIMTYPE, sync bool) (err error) {
	lenght := 1
	if sync {
		lenght = 5
	}
	var bs []byte
	if ts != nil {
		if t.jsonOn {
			bs = JsonEncode(ts)
		} else {
			bs = TEncode(ts)
		}
		lenght += len(bs)
	}
	buf := NewBufferWithCapacity(lenght)
	sendId := atomic.AddInt32(&seq, 1)
	var ch chan int8
	if sync {
		buf.WriteByte(byte(tt) | 0x80)
		buf.Write(Int32ToBytes(sendId))
		ch = await.Get(int64(sendId))
	} else {
		buf.WriteByte(byte(tt))
	}
	if len(bs) > 0 {
		buf.Write(bs)
	}
	//if ts != nil {
	//if t.jsonOn {
	//	buf.Write(JsonEncode(ts))
	//} else {
	//	buf.Write(TEncode(ts))
	//}
	//}
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

func (t *WsSock) sendBigData(data []byte, tt sys.TIMTYPE) (err error) {
	buf := NewBufferWithCapacity(1 + len(data))
	buf.WriteByte(byte(tt))
	buf.Write(data)
	return t._send(buf)
}

func (t *WsSock) close() {
	wsware.delws(t.ws)
}

func awaitEnd(bs []byte) {
	await.Close(int64(BytesToInt32(bs)))
}
