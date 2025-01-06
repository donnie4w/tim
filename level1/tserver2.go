// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package level1

import (
	"context"
	"github.com/donnie4w/tim/log"

	. "github.com/donnie4w/gofer/buffer"
	. "github.com/donnie4w/gofer/util"
	. "github.com/donnie4w/tim/stub"
	"github.com/donnie4w/tim/sys"
	"github.com/donnie4w/tim/util"
	. "github.com/donnie4w/tsf"
)

type tsfClientServer struct {
	isClose         bool
	serverTransport *TServerSocket
}

var tsfclientserver = &tsfClientServer{}

func (tcs *tsfClientServer) server(_addr string, processor Itnet, handler func(tc *tlContext), cliErrorHandler func(tc *tlContext), ok []byte) (err error) {
	if tcs.serverTransport, err = NewTServerSocketTimeout(_addr, sys.ConnectTimeout); err == nil {
		if err = tcs.serverTransport.Listen(); err == nil {
			sys.CSADDR, ok[0] = _addr, 1
			nodeWare.addCsMapNode(sys.UUID)
			log.FmtPrint("cluster listen [", sys.CSADDR, "]")
			for {
				if transport, err := tcs.serverTransport.Accept(); err == nil {
					transport.SetTConfiguration(&TConfiguration{MaxMessageSize: int32(sys.MaxTransLength), Snappy: true})
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
	if !tcs.isClose && err != nil {
		ok[0] = 0
		log.FmtPrint("cluster start failed:", err)
	}
	return
}

func (tcs *tsfClientServer) close() {
	defer util.Recover()
	tcs.isClose = true
	tcs.serverTransport.Close()
}

type sockWork struct {
	transport TsfSocket
	processor Itnet
	handler   func(tc *tlContext)
	cliError  func(tc *tlContext)
	isServer  bool
	tlcontext *tlContext
}

func (sw *sockWork) work(addr string) (err error) {
	defer util.Recover2(&err)
	tlcontext := newTlContext2(sw.transport)
	tlcontext.remoteAddr = addr
	tlcontext.isServer = sw.isServer
	tlcontext.remoteIP, tlcontext.remoteHost2 = remoteHost2(sw.transport)
	go sw.handler(tlcontext)
	tlcontext.defaultCtx = context.WithValue(context.Background(), tlContextCtx, tlcontext)
	tlcontext.cancleChan = make(chan byte)
	sw.tlcontext = tlcontext
	return
}

func (sw *sockWork) final() {
	sw.tlcontext.Close()
	sw.cliError(sw.tlcontext)
}

func (sw *sockWork) processRequests(packet *Packet, processor Itnet) (err error) {
	defer util.Recover2(&err)
	bs := packet.ToBytes()
	sys.Stat.Ib(int64(packet.Len))
	t := bs[0]
	if !sw.tlcontext.isAuth {
		if t != CHAP && t != PING && t != PONG {
			return sw.tlcontext.CloseAndEnd()
		}
	}
	ctx := sw.tlcontext.defaultCtx
	bs = bs[1:]
	switch t {
	case PING:
		err = processor.Ping(ctx, util.Mask(bs))
	case PONG:
		err = processor.Pong(ctx, util.Mask(bs))
	case CHAP:
		err = processor.Chap(ctx, util.Mask(bs))
	case AUTH2:
		err = processor.Auth2(ctx, util.Mask(bs))
	case SYNCNODE:
		if node, err := Decode[Node](util.Mask(bs[1:])); err == nil {
			err = processor.SyncNode(ctx, node, byteToBool(bs[0]))
		}
	case SYNCADDR:
		err = processor.SyncAddr(ctx, string(util.Mask(bs[1:])), byteToBool(bs[0]))
	case SYNCTXMERGE:
		err = processor.SyncTxMerge(ctx, bytesToMap(bs))
	case CSUSER:
		cu, _ := TDecode(bs[8:], &CsUser{})
		err = processor.CsUser(ctx, BytesToInt64(bs[:8]), cu)
	case CSBS:
		cb, _ := TDecode(bs[8:], &CsBs{})
		err = processor.CsBs(ctx, BytesToInt64(bs[:8]), cb)
	case CSREQ:
		csbean, _ := TDecode(bs[9:], &CsBean{})
		err = processor.CsReq(ctx, BytesToInt64(bs[:8]), byteToBool(bs[8]), csbean)
	case CSVR:
		vcr, _ := TDecode(bs[8:], &VBean{})
		err = processor.CsVr(ctx, BytesToInt64(bs[:8]), vcr)
	}
	return

}

type tsfServerClient struct {
	transport TsfSocket
}

var tsfserverclient = &tsfServerClient{}

func (tsc *tsfServerClient) server(addr string, processor Itnet, handler func(tc *tlContext), cliErrorHandler func(tc *tlContext), async bool) (err error) {
	defer util.Recover()
	clientLinkCache.Put(addr, 0)
	defer clientLinkCache.Del(addr)
	tsc.transport = NewTSocketConf(addr, &TConfiguration{ConnectTimeout: sys.ConnectTimeout, MaxMessageSize: int32(sys.MaxTransLength), Snappy: true})
	if err = tsc.transport.Open(); err == nil {
		twork := &sockWork{transport: tsc.transport, processor: processor, handler: handler, cliError: cliErrorHandler, isServer: false}
		twork.work(addr)
		if async {
			go func() {
				defer twork.final()
				tsc.transport.ProcessMerge(func(pkt *Packet) error { return twork.processRequests(pkt, processor) })
			}()
		} else {
			defer twork.final()
			tsc.transport.ProcessMerge(func(pkt *Packet) error { return twork.processRequests(pkt, processor) })
		}
	}
	return
}

func (tsc *tsfServerClient) Close() (err error) {
	defer util.Recover2(&err)
	return tsc.transport.Close()
}

//func process(packet *Packet) {
//	packet.ToBytes()
//}

type ItnetImpl struct {
	ts TsfSocket
}

func method(name byte) *Buffer {
	buf := NewBuffer()
	buf.WriteByte(name)
	return buf
}

func (tl *ItnetImpl) Ping(ctx context.Context, pingBs []byte) (_err error) {
	buf := method(PING)
	if pingBs != nil && len(pingBs) > 0 {
		buf.Write(util.Mask(pingBs))
	}
	return tl.write(buf)
}

func (tl *ItnetImpl) Pong(ctx context.Context, pongBs []byte) (_err error) {
	buf := method(PONG)
	if pongBs != nil && len(pongBs) > 0 {
		buf.Write(util.Mask(pongBs))
	}
	return tl.write(buf)
}

func (tl *ItnetImpl) Chap(ctx context.Context, bss []byte) (_err error) {
	buf := method(CHAP)
	buf.Write(util.Mask(bss))
	return tl.write(buf)
}

func (tl *ItnetImpl) Auth2(ctx context.Context, authKey []byte) (_err error) {
	buf := method(AUTH2)
	buf.Write(util.Mask(authKey))
	return tl.write(buf)
}

func (tl *ItnetImpl) SyncNode(ctx context.Context, node *Node, ir bool) (_err error) {
	buf := method(SYNCNODE)
	buf.WriteByte(boolToByte(ir))
	nodebs, _ := Encode(node)
	buf.Write(util.Mask(nodebs))
	return tl.write(buf)
}

func (tl *ItnetImpl) SyncAddr(ctx context.Context, node string, ir bool) (_err error) {
	buf := method(SYNCADDR)
	buf.WriteByte(boolToByte(ir))
	buf.Write(util.Mask([]byte(node)))
	return tl.write(buf)
}

func (tl *ItnetImpl) SyncTxMerge(ctx context.Context, syncList map[int64]int8) (_err error) {
	buf := method(SYNCTXMERGE)
	syncListBuf := mapTobytes(syncList)
	buf.Write(syncListBuf.Bytes())
	return tl.write(buf)
}

func (tl *ItnetImpl) CsUser(ctx context.Context, sendId int64, cu *CsUser) (_err error) {
	buf := method(CSUSER)
	buf.Write(Int64ToBytes(sendId))
	buf.Write(TEncode(cu))
	return tl.write(buf)
}

func (tl *ItnetImpl) CsBs(ctx context.Context, sendId int64, cb *CsBs) (_err error) {
	buf := method(CSBS)
	buf.Write(Int64ToBytes(sendId))
	buf.Write(TEncode(cb))
	return tl.write(buf)
}

func (tl *ItnetImpl) CsReq(ctx context.Context, sendId int64, ack bool, cb *CsBean) (_err error) {
	buf := method(CSREQ)
	buf.Write(Int64ToBytes(sendId))
	buf.WriteByte(boolToByte(ack))
	buf.Write(TEncode(cb))
	return tl.write(buf)
}

func (tl *ItnetImpl) CsVr(ctx context.Context, sendId int64, vrb *VBean) (_err error) {
	buf := method(CSVR)
	buf.Write(Int64ToBytes(sendId))
	buf.Write(TEncode(vrb))
	return tl.write(buf)
}

func (tl *ItnetImpl) write(buf *Buffer) (err error) {
	sys.Stat.Ob(int64(buf.Len()))
	_, err = tl.ts.WriteWithMerge(buf.Bytes())
	return
}

func mapTobytes(syncList map[int64]int8) (buf *Buffer) {
	buf = NewBufferWithCapacity(9 * len(syncList))
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
	PING        byte = 1
	PONG        byte = 2
	CHAP        byte = 3
	AUTH2       byte = 4
	SYNCNODE    byte = 5
	SYNCADDR    byte = 6
	SYNCTXMERGE byte = 7
	CSUSER      byte = 8
	CSBS        byte = 9
	CSREQ       byte = 10
	CSVR        byte = 11
)
