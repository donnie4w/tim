// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package service

import (
	"github.com/donnie4w/gofer/buffer"
	"github.com/donnie4w/gofer/hashmap"
	"github.com/donnie4w/gofer/lock"
	goutil "github.com/donnie4w/gofer/util"
	"github.com/donnie4w/gothrift/thrift"
	"github.com/donnie4w/tim/amr"
	"github.com/donnie4w/tim/errs"
	"github.com/donnie4w/tim/mq"
	"github.com/donnie4w/tim/stub"
	"github.com/donnie4w/tim/sys"
	"github.com/donnie4w/tim/util"
	"github.com/donnie4w/tlnet"
	"sync/atomic"
	"time"
)

var strLock = lock.NewStrlock(64)
var numLock = lock.NewNumLock(64)
var await = lock.NewFastAwait[int8]()

var wsware = newWsWare()

type wswareHandle struct {
	uMap  *hashmap.Map[string, []int64]
	wsmap *hashmap.MapL[int64, *WsSock]
}

func newWsWare() *wswareHandle {
	ww := &wswareHandle{wsmap: hashmap.NewMapL[int64, *WsSock](), uMap: hashmap.NewMap[string, []int64]()}
	go ww.wsExpriedTicker()
	return ww
}

func (wh *wswareHandle) AddTid(ws *tlnet.Websocket, tid *stub.Tid) {
	lock := numLock.Lock(ws.Id)
	defer lock.Unlock()
	if wh.wsmap.Has(ws.Id) {
		return
	}
	wss := NewWsSock(ws, tid)
	wh.wsmap.Put(ws.Id, wss)
	if a, ok := wh.uMap.Get(tid.Node); !ok {
		wh.uMap.Put(tid.Node, []int64{ws.Id})
	} else {
		wh.uMap.Put(tid.Node, append(a, ws.Id))
	}
	mq.PushOnline(tid.Node, true)
	amr.AddAccount(tid.Node)
}

func (wh *wswareHandle) SetJsonOn(ws *tlnet.Websocket) {
	if wss, ok := wh.wsmap.Get(ws.Id); ok {
		wss.SetJsonOn(true)
	}
}

func (wh *wswareHandle) SendNode(node string, ts thrift.TStruct, tt sys.TIMTYPE) (_r bool) {
	if ids, ok := wh.uMap.Get(node); ok {
		for _, id := range ids {
			if wh.SendWs(id, ts, tt) {
				_r = true
			}
		}
	}
	return
}

func (wh *wswareHandle) SendBigData(node string, data []byte, tt sys.TIMTYPE) (_r bool) {
	if ids, ok := wh.uMap.Get(node); ok {
		for _, id := range ids {
			if wh.SendBigDataByWs(id, data, tt) {
				_r = true
			}
		}
	}
	return
}

func (wh *wswareHandle) SendBigDataByWs(id int64, data []byte, tt sys.TIMTYPE) (_r bool) {
	if wss, ok := wh.wsmap.Get(id); ok {
		if err := wss.sendBigData(data, tt); err == nil {
			_r = ok
		}
	}
	return
}

func (wh *wswareHandle) SendNodeWithAck(node string, ts thrift.TStruct, tt sys.TIMTYPE) (_r bool) {
	if ids, ok := wh.uMap.Get(node); ok {
		for _, id := range ids {
			if wh.SendWsWithAck(id, ts, tt) {
				_r = true
			}
		}
	}
	return
}

func (wh *wswareHandle) SendWs(id int64, ts thrift.TStruct, tt sys.TIMTYPE) (_r bool) {
	if wss, ok := wh.wsmap.Get(id); ok {
		if err := wss.send(ts, tt, false); err == nil {
			_r = ok
		}
	}
	return
}

func (wh *wswareHandle) SendWsWithAck(id int64, ts thrift.TStruct, tt sys.TIMTYPE) (_r bool) {
	if wss, ok := wh.wsmap.Get(id); ok {
		if err := wss.send(ts, tt, true); err == nil {
			_r = ok
		}
	}
	return
}

func (wh *wswareHandle) deviceNums(node string) (_r int) {
	if as, ok := wh.uMap.Get(node); ok {
		_r = len(as)
	}
	return
}

func (wh *wswareHandle) deviceTypeList(node string) (_r []byte) {
	_r = make([]byte, 0)
	if as, ok := wh.uMap.Get(node); ok {
		for _, v := range as {
			if ws, ok := wh.wsmap.Get(v); ok {
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

func (wh *wswareHandle) Ping(id int64) {
	if wss, ok := wh.wsmap.Get(id); ok {
		wss.send(nil, sys.TIMPING, false)
		wss.pTime = time.Now().UnixNano()
	}
}

func (wh *wswareHandle) detect(node string) {
	if ids, ok := wh.uMap.Get(node); ok {
		for _, id := range ids {
			wh.Ping(id)
		}
	}
}

func (wh *wswareHandle) Get(ws *tlnet.Websocket) (*WsSock, bool) {
	return wh.wsmap.Get(ws.Id)
}

func (wh *wswareHandle) hasws(ws *tlnet.Websocket) bool {
	return wh.wsmap.Has(ws.Id)
}

func (wh *wswareHandle) hasUser(node string) bool {
	return wh.uMap.Has(node)
}

func (wh *wswareHandle) wsById(id int64) (*tlnet.Websocket, bool) {
	if ws, b := wh.wsmap.Get(id); b {
		return ws.ws, b
	}
	return nil, false
}

func (wh *wswareHandle) wsByNode(node string) (*tlnet.Websocket, bool) {
	if ids, ok := wh.uMap.Get(node); ok {
		return wh.wsById(ids[0])
	}
	return nil, false
}

func (wh *wswareHandle) wsLen() int64 {
	return wh.wsmap.Len()
}

func (wh *wswareHandle) wssList(index, limit int64) (_r []*stub.Tid, length int64) {
	_r = make([]*stub.Tid, 0)
	count := index
	length = wh.wsmap.Len()
	wh.wsmap.Range(func(_ int64, wss *WsSock) bool {
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

func (wh *wswareHandle) delws(ws *tlnet.Websocket) {
	defer util.Recover()
	wh.delId(ws.Id)
}

func (wh *wswareHandle) delId(id int64) {
	if sk, ok := wh.wsmap.Get(id); ok {
		sk.ws.Close()
		wh.wsmap.Del(id)
		if a, ok := wh.uMap.Get(sk.tid.Node); ok {
			lock := strLock.Lock(sk.tid.Node)
			defer lock.Unlock()
			if _a := util.ArraySub2(a, id); len(_a) > 0 {
				wh.uMap.Put(sk.tid.Node, _a)
			} else {
				wh.uMap.Del(sk.tid.Node)
				amr.RemoveAccount(sk.tid.Node)
				mq.PushOnline(sk.tid.Node, false)
				sys.Interrupt(sk.tid)
			}
		}
	}
}

func (wh *wswareHandle) delnode(node string) {
	if as, ok := wh.uMap.Get(node); ok {
		for _, v := range as {
			wh.delId(v)
		}
	}
}

func (wh *wswareHandle) wsExpriedTicker() {
	tk := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-tk.C:
			func() {
				defer util.Recover()
				wh.wsmap.Range(func(k int64, v *WsSock) bool {
					if v.pingTo() {
						wh.delId(k)
					}
					return true
				})
			}()
		}
	}
}

type WsSock struct {
	ws      *tlnet.Websocket
	tid     *stub.Tid
	jsonOn  bool
	pTime   int64
	amrTime int64
}

func NewWsSock(ws *tlnet.Websocket, tid *stub.Tid) (_r *WsSock) {
	t := time.Now().UnixNano()
	_r = &WsSock{ws: ws, tid: tid, pTime: t, amrTime: t}
	return
}

func (wsk *WsSock) amr() {
	now := time.Now().UnixNano()
	if wsk.amrTime+int64(sys.Conf.TTL*1e9)/2 < now {
		amr.AddAccount(wsk.tid.Node)
		wsk.amrTime = now
	}
}

func (wsk *WsSock) SetJsonOn(on bool) {
	wsk.jsonOn = on
}

func (wsk *WsSock) _send(buf *buffer.Buffer) (err error) {
	if wsk.pingTo() {
		wsk.close()
		return errs.ERR_PING.Error()
	}
	sys.Stat.Ob(int64(buf.Len()))
	return wsk.ws.Send(buf.Bytes())
}

func (wsk *WsSock) pingTo() bool {
	return wsk.pTime+int64(sys.Conf.PingTo*int64(time.Second)) < time.Now().UnixNano()
}

var syncIndex atomic.Int32

func (wsk *WsSock) send(ts thrift.TStruct, tt sys.TIMTYPE, sync bool) (err error) {
	lenght := 1
	if sync {
		lenght = 5
	}
	var bs []byte
	if ts != nil {
		if wsk.jsonOn {
			bs = goutil.JsonEncode(ts)
		} else {
			bs = goutil.TEncode(ts)
		}
		lenght += len(bs)
	}
	resendNum := byte(2)
START:
	buf := buffer.NewBufferWithCapacity(lenght)
	var sendId int32
	if sync {
		sendId = syncIndex.Add(1)
		buf.WriteByte(byte(tt) | 0x80)
		buf.Write(goutil.Int32ToBytes(sendId))
	} else {
		buf.WriteByte(byte(tt))
	}
	if len(bs) > 0 {
		buf.Write(bs)
	}
	if err = wsk._send(buf); err == nil && sync {
		if _, err = await.Wait(int64(sendId), sys.WaitTimeout); err == nil {
			return
		} else if wsk.ws.Error == nil && resendNum > 0 {
			resendNum--
			goto START
		}
	}
	if wsk.ws.Error != nil || err != nil {
		err = errs.ERR_OVERTIME.Error()
	}
	return
}

func (wsk *WsSock) sendBigData(data []byte, tt sys.TIMTYPE) (err error) {
	buf := buffer.NewBufferWithCapacity(1 + len(data))
	buf.WriteByte(byte(tt))
	buf.Write(data)
	return wsk._send(buf)
}

func (wsk *WsSock) close() {
	wsware.delws(wsk.ws)
}

func awaitEnd(bs []byte) {
	await.Close(int64(goutil.BytesToInt32(bs)))
}
