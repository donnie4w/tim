// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim
//

package service

import (
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

func (this *_wsWare) AddTid(ws *tlnet.Websocket, tid *Tid) {
	numLock.Lock(ws.Id)
	defer numLock.Unlock(ws.Id)
	if this.wsmap.Has(ws.Id) {
		return
	}
	wss := NewWsSock(ws)
	wss.tid = tid
	this.wsmap.Put(ws.Id, wss)
	if a, ok := this.uMap.Get(tid.Node); !ok {
		this.uMap.Put(tid.Node, []int64{ws.Id})
	} else {
		this.uMap.Put(tid.Node, append(a, ws.Id))
	}
	go sys.Csuser(tid.Node, true, ws.Id)
}

func (this *_wsWare) SetJsonOn(ws *tlnet.Websocket) {
	if wss, ok := this.wsmap.Get(ws.Id); ok {
		wss.SetJsonOn(true)
	}
}

func (this *_wsWare) SendNode(node string, ts thrift.TStruct, tt sys.TIMTYPE) (_r bool) {
	if ids, ok := this.uMap.Get(node); ok {
		for _, id := range ids {
			if this.SendWs(id, ts, tt) {
				_r = true
			}
		}
	}
	return
}

func (this *_wsWare) SendNodeWithAck(node string, ts thrift.TStruct, tt sys.TIMTYPE) (_r bool) {
	if ids, ok := this.uMap.Get(node); ok {
		for _, id := range ids {
			if this.SendWsWithAck(id, ts, tt) {
				_r = true
			}
		}
	}
	return
}

func (this *_wsWare) SendWs(id int64, ts thrift.TStruct, tt sys.TIMTYPE) (_r bool) {
	if wss, ok := this.wsmap.Get(id); ok {
		if err := wss.send(ts, tt, false); err == nil {
			_r = ok
		}
	}
	return
}

func (this *_wsWare) SendWsWithAck(id int64, ts thrift.TStruct, tt sys.TIMTYPE) (_r bool) {
	if wss, ok := this.wsmap.Get(id); ok {
		if err := wss.send(ts, tt, true); err == nil {
			_r = ok
		}
	}
	return
}

func (this *_wsWare) GetUserDeviceLen(node string) (_r int) {
	if as, ok := this.uMap.Get(node); ok {
		_r = len(as)
	}
	return
}

func (this *_wsWare) GetUserDeviceTypeLen(node string) (_r []byte) {
	_r = make([]byte, 0)
	if as, ok := this.uMap.Get(node); ok {
		for _, v := range as {
			if tid, ok := this.wsmap.Get(v); ok {
				if tid.tid.Termtyp != nil {
					_r = append(_r, byte(*tid.tid.Termtyp))
				}
			}
		}
	}
	return
}

func (this *_wsWare) Get(ws *tlnet.Websocket) (*WsSock, bool) {
	return this.wsmap.Get(ws.Id)
}

func (this *_wsWare) hasws(ws *tlnet.Websocket) bool {
	return this.wsmap.Has(ws.Id)
}

func (this *_wsWare) hasUser(node string) bool {
	return this.uMap.Has(node)
}

func (this *_wsWare) wsLen() int64 {
	return this.wsmap.Len()
}

func (this *_wsWare) wssList() (_r []*Tid) {
	_r = make([]*Tid, 0)
	this.wsmap.Range(func(_ int64, wss *WsSock) bool {
		_r = append(_r, wss.tid)
		return true
	})
	return
}

func (this *_wsWare) wssInfo(node string) (_r []byte) {
	_r = make([]byte, 0)
	if as, ok := this.uMap.Get(node); ok {
		for _, v := range as {
			if ws, ok := this.wsmap.Get(v); ok {
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

func (this *_wsWare) delws(ws *tlnet.Websocket) {
	defer util.Recover()
	this.delId(ws.Id)
}

func (this *_wsWare) delId(id int64) {
	if sk, ok := this.wsmap.Get(id); ok {
		sk.ws.Close()
		this.wsmap.Del(id)
		if a, ok := this.uMap.Get(sk.tid.Node); ok {
			strLock.Lock(sk.tid.Node)
			defer strLock.Unlock(sk.tid.Node)
			if _a := util.ArraySub2(a, id); len(_a) > 0 {
				this.uMap.Put(sk.tid.Node, _a)
			} else {
				this.uMap.Del(sk.tid.Node)
				go sys.Csuser(sk.tid.Node, false, id)
				go sys.Interrupt(sk.tid)
			}
		}
	}
}

func (this *_wsWare) delnode(node string) {
	if as, ok := this.uMap.Get(node); ok {
		for _, v := range as {
			this.delId(v)
		}
	}
}

type WsSock struct {
	ws     *tlnet.Websocket
	tid    *Tid
	jsonOn bool
}

func NewWsSock(ws *tlnet.Websocket) (_r *WsSock) {
	_r = &WsSock{ws: ws}
	return
}

func (this *WsSock) SetJsonOn(on bool) {
	this.jsonOn = on
}

func (this *WsSock) _send(buf *Buffer) (err error) {
	sys.Stat.Ob(int64(buf.Len()))
	return this.ws.Send(buf.Bytes())
}

var seq int32 = 0

func (this *WsSock) send(ts thrift.TStruct, t sys.TIMTYPE, sync bool) (err error) {
	buf := NewBuffer()
	sendId := atomic.AddInt32(&seq, 1)
	var ch chan int8
	if sync {
		buf.WriteByte(byte(t) | 0x80)
		buf.Write(Int32ToBytes(sendId))
		ch = await.Get(int64(sendId))
	} else {
		buf.WriteByte(byte(t))
	}
	if ts != nil {
		if this.jsonOn {
			buf.Write(JsonEncode(ts))
		} else {
			buf.Write(TEncode(ts))
		}
	}
	if err = this._send(buf); err == nil && sync {
		i := 0
		for this.ws.Error == nil && i < 100 {
			i++
			select {
			case <-ch:
				err = nil
				goto END
			case <-time.After(time.Second):
				err = sys.ERR_OVERTIME.Error()
			}
		}
		if this.ws.Error != nil || err != nil {
			err = sys.ERR_OVERTIME.Error()
		}
	}
END:
	return
}

func (this *WsSock) close() {
	wsware.delws(this.ws)
}

func awaitEnd(bs []byte) {
	await.DelAndClose(int64(BytesToInt32(bs)))
}
