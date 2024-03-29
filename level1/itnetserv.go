// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim
package level1

import (
	"context"
	"fmt"
	"sort"
	"sync/atomic"
	"time"

	. "github.com/donnie4w/gofer/buffer"
	"github.com/donnie4w/gofer/keystore"
	. "github.com/donnie4w/gofer/lock"
	. "github.com/donnie4w/gofer/util"
	"github.com/donnie4w/simplelog/logging"
	. "github.com/donnie4w/tim/stub"
	"github.com/donnie4w/tim/sys"
	"github.com/donnie4w/tim/util"
)

type itnetServ struct {
	mux *Numlock
}

func ctx2TlContext(ctx context.Context) (tc *tlContext) {
	tc = ctx.Value(tlContextCtx).(*tlContext)
	tc.pingNum = 0
	return
}

func (this *itnetServ) Ping(ctx context.Context, pingBs []byte) (_err error) {
	go func() {
		defer util.Recover()
		tc := ctx2TlContext(ctx)
		if tc.pongNum > 10 && !tc.isAuth {
			tc.CloseAndEnd()
			return
		}
		if tc.isAuth {
			tc.iface.Pong(context.TODO(), poBs(tc))
			if pingBs != nil && len(pingBs) > 0 {
				pn := BytesToInt64(pingBs)
				if pn != pingnum() {
					tc.iface.SyncNode(context.TODO(), nodeWare.GetUUIDNode(), true)
				}
			}
		}
	}()
	return
}

func (this *itnetServ) Pong(ctx context.Context, pongBs []byte) (_err error) {
	go func() {
		defer util.Recover()
		tc := ctx2TlContext(ctx)
		if tc.pingNum > 0 {
			atomic.AddInt64(&tc.pingNum, -1)
		}
		atomic.AddInt64(&tc.pongNum, 1)
		if tc.isAuth {
			if pongBs != nil && len(pongBs) > 0 {
				if d, err := TDecode(pongBs, &Data{}); err == nil {
					if d.OnNum != nil {
						tc.onNum = *d.OnNum
					}
					if d.SyncNum != nil && *d.SyncNum != syncNum() {
						tc.iface.Pong(context.TODO(), cslistBytes())
					}
					if d.CsNum != nil {
						tc.remoteCsNum = *d.CsNum
					}
					if d.Bytes != nil {
						is := BytesToIntArray(d.Bytes)
						ds := util.ArraySub(is, nodeWare.getCsList())
						if ds != nil {
							for _, v := range ds {
								_unaccess.Put(v, 0)
							}
						}
					}
				}
			}
		}
	}()
	return
}

func (this *itnetServ) Chap(ctx context.Context, abs []byte) (_err error) {
	go func() {
		defer util.Recover()
		tc := ctx2TlContext(ctx)
		pass := false
		if ab, err := decodeChapBean(abs); err == nil {
			if !chapTxTemp.Has(ab.TxId) && ab.Time+24*int64(time.Hour) > time.Now().UnixNano() {
				chapTxTemp.Put(ab.TcId, 0)
				switch ab.Stat {
				case 1:
					if ab.Key == sys.Conf.Pwd {
						tc.remoteUuid, tc.verifycode = ab.UUID, RandId()
						ab.Code, ab.TcId, ab.TxId, ab.Stat = tc.verifycode, tc.id, RandId(), ab.Stat+1
						if bs, err := encodeChapBean(ab); err == nil {
							if err = tc.iface.Chap(context.TODO(), bs); err == nil {
								pass = true
							}
						}
					}
				case 2:
					if ab.IDcard == tc.id && ab.UUID == sys.UUID {
						ab.TxId, ab.Code, ab.Stat, ab.IDcard = RandId(), ab.Code+1, ab.Stat+1, ab.IDcard+1
						if bs, err := encodeChapBean(ab); err == nil {
							tc.iface.Chap(context.TODO(), bs)
							pass = true
						}
					}
				case 3:
					if ab.Key == sys.Conf.Pwd && ab.Code == tc.verifycode+1 && ab.UUID == tc.remoteUuid && ab.TcId == tc.id {
						tc.isAuth, ab.Stat, ab.TcId = true, ab.Stat+1, sys.UUID
						availMap.Del(tc)
						if bs, err := encodeChapBean(ab); err == nil {
							tc.iface.Chap(context.TODO(), bs)
							pass = true
						}
					}
				case 4:
					if ab.IDcard == tc.id+1 && ab.UUID == sys.UUID {
						tc.isAuth, pass, tc.remoteUuid = true, true, ab.TcId
						availMap.Del(tc)
						if !sys.LA {
							tc.iface.SyncNode(context.TODO(), nodeWare.GetUUIDNode(), true)
						} else {
							tc.iface.SyncAddr(context.TODO(), "", true)
						}
					}
				}
			}
		}
		if !pass {
			tc.Close()
		}
	}()
	return
}

func (this *itnetServ) Auth2(ctx context.Context, authKey []byte) (_err error) {
	defer util.Recover()
	tc := ctx2TlContext(ctx)
	if authTc(tc, authKey) {
		tc.iface.SyncNode(context.TODO(), nodeWare.GetUUIDNode(), true)
	}
	return
}

func (this *itnetServ) SyncNode(ctx context.Context, node *Node, ir bool) (_err error) {
	go func() {
		defer util.Recover()
		tc := ctx2TlContext(ctx)
		if tc.isAuth && !sys.LA {
			if tc.remoteUuid == 0 || node.UUID == 0 || node.Addr == "" || tc.remoteUuid != node.UUID {
				return
			}
			tc.remoteUuid, tc.remoteAddr = node.UUID, node.Addr
			nodeWare.add(tc)
			if ir {
				tc.iface.SyncNode(context.TODO(), nodeWare.GetUUIDNode(), !ir)
			}

			for k, v := range node.Nodekv {
				if k != sys.UUID && !nodeWare.hasUUID(k) && !clientLinkCache.Has(v) {
					<-time.After(time.Duration(Rand(6)) * time.Second)
					if !nodeWare.hasUUID(k) {
						go func(k int64, v string) {
							if _, err2 := tnetservice.connectLoop(v, true, 10); err2 != nil {
								_unaccess.Put(k, 0)
							} else {
								_unaccess.Del(k)
							}
						}(k, v)
					}
				}
			}
		}
	}()
	return
}

func (this *itnetServ) SyncAddr(ctx context.Context, node string, ir bool) (_err error) {
	go func() {
		defer util.Recover()
		tc := ctx2TlContext(ctx)
		if tc.isAuth {
			if ir {
				tc.iface.SyncAddr(context.TODO(), fmt.Sprint(tc.remoteHost2), !ir)
			} else {
				tc.CloseAndEnd()
				lnetservice._server(node)
			}
		}
	}()
	return
}

func (this *itnetServ) SyncTxMerge(ctx context.Context, syncList map[int64]int8) (_err error) {
	go func() {
		defer util.Recover()
		if ctx2TlContext(ctx).isAuth {
			if syncList != nil {
				for k, v := range syncList {
					this._SyncTx(k, v)
				}
			}
		}
	}()
	return
}

func (this *itnetServ) CsUser(ctx context.Context, sendId int64, cu *CsUser) (_err error) {
	go func() {
		defer util.Recover()
		sys.Stat.CProsDo()
		defer sys.Stat.CProsDone()
		tc := ctx2TlContext(ctx)
		if tc.isAuth {
			syncTxAck(tc, tc.remoteUuid, sendId, 0)
			if cu.Node != nil {
				if cu.Stat == 1 {
					if cu.Node != nil {
						for k := range cu.Node {
							addcsu(k, tc.remoteUuid)
						}
						bkCsuserBatch(cu.Node, tc.remoteUuid, cu.Stat)
					}
					if cu.BkNode != nil {
						for k, v := range cu.BkNode {
							addcsu(k, v)
						}
					}
				} else if cu.Stat == 2 {
					if cu.Node != nil {
						for k := range cu.Node {
							delcsu(k, tc.remoteUuid)
						}
						bkCsuserBatch(cu.Node, tc.remoteUuid, cu.Stat)
					}
					if cu.BkNode != nil {
						for k, v := range cu.BkNode {
							delcsu(k, v)
						}
					}
				}
			}
		}
	}()
	return
}

func (this *itnetServ) CsBs(ctx context.Context, sendId int64, cb *CsBs) (_err error) {
	go func() {
		defer util.Recover()
		sys.Stat.CProsDo()
		defer sys.Stat.CProsDone()
		tc := ctx2TlContext(ctx)
		if tc.isAuth {
			if cb.Cache != nil {
				if hascsu(string(cb.Node)) {
					syncTxAck(tc, tc.remoteUuid, sendId, 1)
				} else {
					syncTxAck(tc, tc.remoteUuid, sendId, 0)
					return
				}
			} else {
				syncTxAck(tc, tc.remoteUuid, sendId, 1)
			}
			switch cb.BsType {
			case sys.CB_MESSAGE:
				if tm, err := TDecode(cb.Bs, &TimMessage{}); err == nil {
					if reTx.Has(*tm.ID) {
						logging.Warn("reTx>>>", tm)
						return
					}
					reTx.Put(*tm.ID, 0)
					sys.TimMessageProcessor(tm, cb.TransType)
				}
			case sys.CB_PRESENCE:
				if tp, err := TDecode(cb.Bs, &TimPresence{}); err == nil {
					if cb.TransType == sys.TRANS_STAFF || cb.TransType == sys.TRANS_GOAL {
						if tp.ToList != nil {
							for _, u := range tp.ToList {
								_tp := &TimPresence{ID: tp.ID, FromTid: tp.FromTid, ToTid: &Tid{Node: u}, SubStatus: tp.SubStatus, Offline: tp.Offline, Show: tp.Show, Status: tp.Status, Extend: tp.Extend, Extra: tp.Extra}
								sys.TimPresenceProcessor(_tp, sys.TRANS_GOAL)
							}
						} else {
							sys.TimPresenceProcessor(tp, sys.TRANS_GOAL)
						}
					} else {
						sys.TimPresenceProcessor(tp, cb.TransType)
					}
				}
			}
		}
	}()
	return
}

func (this *itnetServ) CsReq(ctx context.Context, sendId int64, ack bool, cb *CsBean) (_err error) {
	go func() {
		defer util.Recover()
		sys.Stat.CProsDo()
		defer sys.Stat.CProsDone()
		tc := ctx2TlContext(ctx)
		if tc.isAuth {
			if ack {
				awaitCsBean.DelAndPut(sendId, cb)
			} else {
				switch cb.RType {
				case 1:
					retm := map[string][]int64{}
					for k := range cb.Bsm2 {
						if cm := getcsu(k); cm != nil {
							im := []int64{}
							cm.Range(func(k int64, _ int8) bool {
								im = append(im, k)
								return true
							})
							retm[k] = im
						}
					}
					cb.Bsm2 = retm
					nodeWare.csReqHandle(tc.remoteUuid, sendId, true, cb)
				case 2:
					bsm := make(map[string][]byte, 0)
					for u := range cb.Bsm {
						bsm[u] = sys.WssInfo(u)
					}
					cb.Bsm = bsm
					nodeWare.csReqHandle(tc.remoteUuid, sendId, true, cb)
				}
			}
		}
	}()
	return
}

func (this *itnetServ) CsVr(ctx context.Context, sendId int64, vrb *VBean) (_err error) {
	go func() {
		defer util.Recover()
		sys.Stat.CProsDo()
		defer sys.Stat.CProsDone()
		tc := ctx2TlContext(ctx)
		if tc.isAuth {
			processVBean(vrb, tc.remoteUuid)
			// switch vrb.Rtype {
			// case 1:
			// 	vgate.VGate.Register(*vrb.FoundNode, vrb.Vnode)
			// 	bkCsVr(vrb, tc.remoteUuid)
			// case 2:
			// 	vgate.VGate.Remove(*vrb.FoundNode, vrb.Vnode)
			// 	bkCsVr(vrb, tc.remoteUuid)
			// case 3:
			// 	vgate.VGate.AddAuth(vrb.Vnode, *vrb.FoundNode, *vrb.Rnode)
			// 	bkCsVr(vrb, tc.remoteUuid)
			// case 4:
			// 	vgate.VGate.DelAuth(vrb.Vnode, *vrb.FoundNode, *vrb.Rnode)
			// 	bkCsVr(vrb, tc.remoteUuid)
			// case 5:
			// 	if _, b := vgate.VGate.GetVroom(vrb.Vnode); b {
			// 		logging.Debug(">>>", tc.remoteUuid)
			// 		vgate.VGate.Sub(vrb.Vnode, tc.remoteUuid, 0)
			// 	} else {
			// 		nodeWare.csVrHandle(tc.remoteUuid, 0, &VBean{Rtype: 10, Vnode: vrb.Vnode, Rnode: vrb.Rnode})
			// 	}
			// case 6:
			// 	vgate.VGate.DelUuid(vrb.Vnode, tc.remoteUuid)
			// case 7:
			// 	if vr, ok := vgate.VGate.GetVroom(vrb.Vnode); ok {
			// 		if !vr.Auth(*vrb.Rnode) {
			// 			nodeWare.csVrHandle(tc.remoteUuid, 0, &VBean{Rtype: 9, Vnode: vrb.Vnode, Rnode: vrb.Rnode})
			// 			return
			// 		}
			// 		vr.Updatetime()
			// 		vrb.Rtype = 50 + vrb.Rtype
			// 		m := map[int64]int8{}
			// 		vgate.VGate.GetUUID(vrb.Vnode).Range(func(k int64, _ int8) bool {
			// 			logging.Debug(">>>>", k, "  =====>", tc.remoteUuid)
			// 			if k != tc.remoteUuid {
			// 				m[k] = 0
			// 				nodeWare.csVrHandle(k, 0, vrb)
			// 			}
			// 			return true
			// 		})
			// 		go sys.TimSteamProcessor(vrb)
			// 		if li, ok := nodeWare.bkuuid(vrb.Vnode); ok && len(li) > 0 {
			// 			for _, u := range li {
			// 				if _, ok := m[u]; !ok && u != tc.remoteUuid {
			// 					m[u] = 0
			// 					nodeWare.csVrHandle(u, 0, vrb)
			// 				}
			// 			}
			// 		}
			// 	} else {
			// 		nodeWare.csVrHandle(tc.remoteUuid, 0, &VBean{Rtype: 8, Vnode: vrb.Vnode, Rnode: vrb.Rnode})
			// 	}
			// case 8:
			// 	f := false
			// 	if vr, ok := vgate.VGate.GetVroom(vrb.Vnode); ok {
			// 		if vr.FoundNode != "" {
			// 			nodeWare.csVrHandle(tc.remoteUuid, 0, &VBean{Rtype: 1, Vnode: vrb.Vnode, FoundNode: &vr.FoundNode})
			// 			vr.AuthMap().Range(func(k string, _ int8) bool {
			// 				nodeWare.csVrHandle(tc.remoteUuid, 0, &VBean{Rtype: 3, Vnode: vrb.Vnode, FoundNode: &vr.FoundNode, Rnode: &k})
			// 				return true
			// 			})
			// 			f = true
			// 		}
			// 	}
			// 	if !f {
			// 		sys.SendNode(*vrb.Rnode, &TimAck{Ok: false, TimType: int8(sys.TIMSTREAM), Error: sys.ERR_NOEXIST.TimError()}, sys.TIMACK)
			// 	}
			// case 9:
			// 	sys.SendNode(*vrb.Rnode, &TimAck{Ok: false, TimType: int8(sys.TIMSTREAM), Error: sys.ERR_AUTH.TimError(), N: &vrb.Vnode}, sys.TIMACK)
			// case 10:
			// 	sys.SendNode(*vrb.Rnode, &TimAck{Ok: false, TimType: int8(sys.TIMVROOM), Error: sys.ERR_NOEXIST.TimError(), N: &vrb.Vnode}, sys.TIMACK)
			// case 51:
			// 	vgate.VGate.Register(*vrb.FoundNode, vrb.Vnode)
			// case 52:
			// 	vgate.VGate.Remove(*vrb.FoundNode, vrb.Vnode)
			// case 53:
			// 	vgate.VGate.AddAuth(vrb.Vnode, *vrb.FoundNode, *vrb.Rnode)
			// case 54:
			// 	vgate.VGate.DelAuth(vrb.Vnode, *vrb.FoundNode, *vrb.Rnode)
			// case 57:
			// 	logging.Debug(">>>", vrb.Body)
			// 	if !reStream.Has(*vrb.StreamId) {
			// 		reStream.Put(*vrb.StreamId, 0)
			// 		sys.TimSteamProcessor(vrb)
			// 	} else {
			// 		buf := NewBufferByPool()
			// 		defer buf.Free()
			// 		buf.WriteString(vrb.Vnode)
			// 		buf.WriteString(sys.Conf.Salt)
			// 		buf.Write(Int64ToBytes(tc.remoteUuid))
			// 		f := CRC64(buf.Bytes())
			// 		if i, ok := reStreamUUID.Get(f); ok {
			// 			if i > 5 {
			// 				vb := &VBean{Rtype: 6, Vnode: vrb.Vnode}
			// 				nodeWare.csVrHandle(tc.remoteUuid, 0, vb)
			// 			} else {
			// 				reStreamUUID.Put(f, atomic.AddInt32(&i, 1))
			// 			}
			// 		} else {
			// 			reStreamUUID.Put(f, 1)
			// 		}
			// 	}
			// default:
			// }
		}
	}()
	return
}

func (this *itnetServ) _SyncTx(syncId int64, result int8) {
	defer util.Recover()
	if result == 0 {
		await.DelAndClose(syncId)
	} else {
		await.DelAndPut(syncId, result)
	}
}

func syncTxAck(tc *tlContext, uuid, syncId int64, result int8) (err error) {
	defer util.Recover()
	tc.mergeChan <- &syncBean{syncId, result}
	atomic.AddInt64(&tc.mergeCount, 1)
	if tc.mergemux.TryLock() {
		go syncMerge(tc)
	}
	return
}

func authTc(tc *tlContext, authKey []byte) (b bool) {
	if authKey != nil && len(authKey) > 0 {
		var bs []byte
		var err error
		if sys.OpenSSL.PrivateBytes != nil {
			bs, err = keystore.RsaDecryptByBytes(authKey, sys.OpenSSL.PrivateBytes)
		} else {
			bs, err = keystore.RsaDecrypt(authKey, "")
		}
		if err == nil {
			if r := BytesToInt64(bs[:8]); r > 0 && r != sys.UUID && !nodeWare.hasUUID(r) && string(bs[8:]) == sys.Conf.Pwd {
				tc.remoteUuid = r
				tc.isAuth = true
				b = true
			}
		}
	}
	return
}

func pingnum() (_r int64) {
	all := nodeWare.getCsList()
	if len(all) > 0 && !sys.LA {
		sort.Slice(all, func(i, j int) bool {
			return all[i] < all[j]
		})
		var buf = NewBufferByPool()
		defer buf.Free()
		for _, a := range all {
			buf.Write(Int64ToBytes(a))
		}
		_r = int64(CRC64(buf.Bytes()))
	}
	return
}

func piBs(tc *tlContext) (_r []byte) {
	id := pingnum()
	if tc.selfPing == id {
		_r = nil
	} else {
		_r = Int64ToBytes(id)
		tc.selfPing = id
	}
	return
}

func poBs(tc *tlContext) (_r []byte) {
	d := &Data{}
	sn := syncNum()
	d.SyncNum = &sn
	on := sys.WssLen()
	cn := int32(len(nodeWare.getCsList()))
	d.CsNum = &cn
	d.OnNum = &on
	_r = TEncode(d)
	id := CRC32(_r)
	if tc.selfPong == id {
		_r = nil
	} else {
		tc.selfPong = id
	}
	return
}

func syncNum() (_r int64) {
	all := nodeWare.getCsList()
	if len(all) > 0 && !sys.LA {
		sort.Slice(all, func(i, j int) bool {
			return all[i] < all[j]
		})
		var buf = NewBufferByPool()
		defer buf.Free()
		for _, a := range all {
			buf.Write(Int64ToBytes(a))
		}
		_r = int64(CRC64(buf.Bytes()))
	}
	return
}

func cslistBytes() (_r []byte) {
	d := &Data{}
	all := nodeWare.getCsList()
	d.Bytes = IntArrayToBytes(all)
	_r = TEncode(d)
	return
}

func encodeChapBean(ab *chapBean) (_r []byte, err error) {
	bs, _ := Encode(ab)
	return util.AesEncode(bs)
}

func decodeChapBean(bs []byte) (_r *chapBean, err error) {
	if bs, err = util.AesDecode(bs); err == nil {
		_r, err = Decode[chapBean](bs)
	}
	return
}
