// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim
//

package level1

import (
	"context"

	. "github.com/donnie4w/gofer/buffer"
	. "github.com/donnie4w/gofer/util"
	. "github.com/donnie4w/tim/stub"
	"github.com/donnie4w/tim/sys"
	"github.com/donnie4w/tim/util"
	. "github.com/donnie4w/tsf/packet"
	. "github.com/donnie4w/tsf/server"
)

type tsfClientServer struct {
	isClose         bool
	serverTransport *TServerSocket
}

var tsfclientserver = &tsfClientServer{}

func (this *tsfClientServer) server( _addr string, processor Itnet, handler func(tc *tlContext), cliErrorHandler func(tc *tlContext), ok []byte) (err error) {
	if this.serverTransport, err = NewTServerSocketTimeout(_addr, sys.ConnectTimeout); err == nil {
		if err = this.serverTransport.Listen(); err == nil {
			sys.CSADDR, ok[0] = _addr, 1
			nodeWare.addCsMapNode(sys.UUID)
			sys.FmtLog("cluster listen [", sys.CSADDR, "]")
			for {
				if transport, err := this.serverTransport.Accept(); err == nil {
					transport.SetTConfiguration(&TConfiguration{MaxMessageSize: int32(sys.MaxTransLength), SnappyMergeData: true})
					go func() {
						defer util.Recover()
						twork := &sockWork{transport: transport, processor: processor, handler: handler, cliError: cliErrorHandler, isServer: true}
						twork.work("")
						defer twork.final()
						transport.ProcessMerge(func(pkt *Packet) error { return twork.processRequests(pkt, processor) })
					}()
				}
			}
		}
	}
	if !this.isClose && err != nil {
		ok[0] = 0
		sys.FmtLog("cluster start failed:", err)
	}
	return
}

func (this *tsfClientServer) close() {
	defer recover()
	this.isClose = true
	this.serverTransport.Close()
}

type sockWork struct {
	transport *TSocket
	processor Itnet
	handler   func(tc *tlContext)
	cliError  func(tc *tlContext)
	isServer  bool
	tlcontext *tlContext
}

func (this *sockWork) work(addr string) (err error) {
	defer util.Recover()
	tlcontext := newTlContext2(this.transport)
	tlcontext.remoteAddr = addr
	tlcontext.isServer = this.isServer
	tlcontext.remoteIP, tlcontext.remoteHost2 = remoteHost2(this.transport)
	go this.handler(tlcontext)
	tlcontext.defaultCtx = context.WithValue(context.Background(), tlContextCtx, tlcontext)
	tlcontext.cancleChan = make(chan byte)
	this.tlcontext = tlcontext
	return
}

func (this *sockWork) final() {
	this.tlcontext.Close()
	this.cliError(this.tlcontext)
}

func (this *sockWork) processRequests(packet *Packet, processor Itnet) (err error) {
	defer util.Recover()
	bs := packet.ToBytes()
	sys.Stat.Ib(int64(packet.Len))
	t := bs[0]
	if !this.tlcontext.isAuth {
		if t != _Chap && t != _Ping && t != _Pong {
			return this.tlcontext.CloseAndEnd()
		}
	}
	ctx := this.tlcontext.defaultCtx
	bs = bs[1:]
	switch t {
	case _Ping:
		err = processor.Ping(ctx, util.Mask(bs))
	case _Pong:
		err = processor.Pong(ctx, util.Mask(bs))
	case _Chap:
		err = processor.Chap(ctx, util.Mask(bs))
	case _Auth2:
		err = processor.Auth2(ctx, util.Mask(bs))
	case _SyncNode:
		if node, err := Decode[Node](util.Mask(bs[1:])); err == nil {
			err = processor.SyncNode(ctx, node, byteToBool(bs[0]))
		}
	case _SyncAddr:
		err = processor.SyncAddr(ctx, string(util.Mask(bs[1:])), byteToBool(bs[0]))
	case _SyncTxMerge:
		err = processor.SyncTxMerge(ctx, bytesToMap(bs))
	case _CsUser:
		cu, _ := TDecode(bs[8:], &CsUser{})
		err = processor.CsUser(ctx, BytesToInt64(bs[:8]), cu)
	case _CsBs:
		cb, _ := TDecode(bs[8:], &CsBs{})
		err = processor.CsBs(ctx, BytesToInt64(bs[:8]), cb)
	case _CsReq:
		csbean, _ := TDecode(bs[9:], &CsBean{})
		err = processor.CsReq(ctx, BytesToInt64(bs[:8]), byteToBool(bs[8]), csbean)
	case _CsVr:
		vcr, _ := TDecode(bs[8:], &VBean{})
		err = processor.CsVr(ctx, BytesToInt64(bs[:8]), vcr)
	}
	return

}

type tsfServerClient struct {
	transport *TSocket
}

var tsfserverclient = &tsfServerClient{}

func (this *tsfServerClient) server(addr string, processor Itnet, handler func(tc *tlContext), cliErrorHandler func(tc *tlContext), async bool) (err error) {
	defer util.Recover()
	clientLinkCache.Put(addr, 0)
	defer clientLinkCache.Del(addr)
	this.transport = NewTSocketConf(addr, &TConfiguration{ConnectTimeout: sys.ConnectTimeout, MaxMessageSize: int32(sys.MaxTransLength), SnappyMergeData: true})
	if err = this.transport.Open(); err == nil {
		twork := &sockWork{transport: this.transport, processor: processor, handler: handler, cliError: cliErrorHandler, isServer: false}
		twork.work(addr)
		if async {
			go func() {
				defer twork.final()
				this.transport.ProcessMerge(func(pkt *Packet) error { return twork.processRequests(pkt, processor) })
			}()
		} else {
			defer twork.final()
			this.transport.ProcessMerge(func(pkt *Packet) error { return twork.processRequests(pkt, processor) })
		}
	} 
	return
}

func (this *tsfServerClient) Close() (err error) {
	defer util.Recover()
	return this.transport.Close()
}

func process(packet *Packet) {
	packet.ToBytes()
}

type ItnetImpl struct {
	ts *TSocket
}

func method(name byte) (buf *Buffer) {
	buf = NewBuffer()
	buf.WriteByte(name)
	return
}

func (this *ItnetImpl) Ping(ctx context.Context, pingBs []byte) (_err error) {
	buf := method(_Ping)
	if pingBs != nil && len(pingBs) > 0 {
		buf.Write(util.Mask(pingBs))
	}
	return this.write(buf)
}

func (this *ItnetImpl) Pong(ctx context.Context, pongBs []byte) (_err error) {
	buf := method(_Pong)
	if pongBs != nil && len(pongBs) > 0 {
		buf.Write(util.Mask(pongBs))
	}
	return this.write(buf)
}

func (this *ItnetImpl) Chap(ctx context.Context, bss []byte) (_err error) {
	buf := method(_Chap)
	buf.Write(util.Mask(bss))
	return this.write(buf)
}

func (this *ItnetImpl) Auth2(ctx context.Context, authKey []byte) (_err error) {
	buf := method(_Auth2)
	buf.Write(util.Mask(authKey))
	return this.write(buf)
}

func (this *ItnetImpl) SyncNode(ctx context.Context, node *Node, ir bool) (_err error) {
	buf := method(_SyncNode)
	buf.WriteByte(boolToByte(ir))
	nodebs, _ := Encode(node)
	buf.Write(util.Mask(nodebs))
	return this.write(buf)
}

func (this *ItnetImpl) SyncAddr(ctx context.Context, node string, ir bool) (_err error) {
	buf := method(_SyncAddr)
	buf.WriteByte(boolToByte(ir))
	buf.Write(util.Mask([]byte(node)))
	return this.write(buf)
}

func (this *ItnetImpl) SyncTxMerge(ctx context.Context, syncList map[int64]int8) (_err error) {
	buf := method(_SyncTxMerge)
	syncListBuf := mapTobytes(syncList)
	buf.Write(syncListBuf.Bytes())
	return this.write(buf)
}

func (this *ItnetImpl) CsUser(ctx context.Context, sendId int64, cu *CsUser) (_err error) {
	buf := method(_CsUser)
	buf.Write(Int64ToBytes(sendId))
	buf.Write(TEncode(cu))
	return this.write(buf)
}

func (this *ItnetImpl) CsBs(ctx context.Context, sendId int64, cb *CsBs) (_err error) {
	buf := method(_CsBs)
	buf.Write(Int64ToBytes(sendId))
	buf.Write(TEncode(cb))
	return this.write(buf)
}

func (this *ItnetImpl) CsReq(ctx context.Context, sendId int64, ack bool, cb *CsBean) (_err error) {
	buf := method(_CsReq)
	buf.Write(Int64ToBytes(sendId))
	buf.WriteByte(boolToByte(ack))
	buf.Write(TEncode(cb))
	return this.write(buf)
}

func (this *ItnetImpl) CsVr(ctx context.Context, sendId int64, vrb *VBean) (_err error) {
	buf := method(_CsVr)
	buf.Write(Int64ToBytes(sendId))
	buf.Write(TEncode(vrb))
	return this.write(buf)
}

func (this *ItnetImpl) write(buf *Buffer) (err error) {
	sys.Stat.Ob(int64(buf.Len()))
	_, err = this.ts.WriteWithMerge(buf.Bytes())
	return
}

func mapTobytes(syncList map[int64]int8) (buf *Buffer) {
	buf = NewBuffer()
	for k, v := range syncList {
		buf.Write(Int64ToBytes(k))
		buf.WriteByte(byte(v))
	}
	return
}

func bytesToMap(bs []byte) (syncList map[int64]int8) {
	syncList = make(map[int64]int8, 0)
	for i := 0; i < len(bs); i += 9 {
		k := BytesToInt64(bs[i : i+8])
		syncList[k] = int8(bs[i+8:][0])
	}
	return
}

func boolToByte(b bool) (_r byte) {
	if b {
		_r = 1
	}
	return
}
func byteToBool(b byte) (_r bool) {
	if b == 1 {
		_r = true
	}
	return
}

const (
	_Ping        byte = 1
	_Pong        byte = 2
	_Chap        byte = 3
	_Auth2       byte = 4
	_SyncNode    byte = 5
	_SyncAddr    byte = 6
	_SyncTxMerge byte = 7
	_CsUser      byte = 8
	_CsBs        byte = 9
	_CsReq       byte = 10
	_CsVr        byte = 11
)
