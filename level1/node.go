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
	"strings"
	"time"

	. "github.com/donnie4w/gofer/hashmap"
	. "github.com/donnie4w/gofer/lock"
	. "github.com/donnie4w/gofer/util"
	. "github.com/donnie4w/tim/stub"
	"github.com/donnie4w/tim/sys"
	"github.com/donnie4w/tim/util"
	"github.com/donnie4w/tim/vgate"
)

type nodeware struct {
	tlMap        *MapL[*tlContext, int8]
	csMap        *Consistenthash
	mux          *Numlock
	_cacheCsList []int64
	_cacheMap    *Map[int64, int8]
}

var nodeWare = newNodeWare()

func newNodeWare() (_r *nodeware) {
	_r = new(nodeware)
	_r.tlMap = NewMapL[*tlContext, int8]()
	_r.csMap = NewConsistenthash(1 << 10)
	_r._cacheMap = NewMap[int64, int8]()
	_r.mux = NewNumLock(1 << 7)
	go _r.csuTicker()
	return
}

func (this *nodeware) add(tc *tlContext) (err error) {
	this.mux.Lock(tc.remoteUuid)
	defer this.mux.Unlock(tc.remoteUuid)
	if this.tlMap.Has(tc) {
		return sys.ERR_UUID_REUSE.Error()
	}
	if sys.UUID == tc.remoteUuid || this.hasUUID(tc.remoteUuid) {
		tc.CloseAndEnd()
		return sys.ERR_UUID_REUSE.Error()
	}
	this.tlMap.Put(tc, 1)
	_unaccess.Del(tc.remoteUuid)
	this._cacheMap.Put(tc.remoteUuid, 0)
	this._cacheCsList = nil
	go this.addCsMapNode(tc.remoteUuid)
	return
}

func (this *nodeware) has(tc *tlContext) bool {
	return this.tlMap.Has(tc)
}

func (this *nodeware) hasUUID(uuid int64) (b bool) {
	return this._cacheMap.Has(uuid) || uuid == sys.UUID
}

func (this *nodeware) del(tc *tlContext) {
	if !tc.isClose {
		tc.Close()
	}
	this.tlMap.Del(tc)
	this._cacheMap.Del(tc.remoteUuid)
	this._cacheCsList = nil
	if !this.hasUUID(tc.remoteUuid) {
		this.csMap.Del(tc.remoteUuid)
	}
}

func (this *nodeware) delAndNoReconn(tc *tlContext) {
	tc._do_reconn = true
	this.del(tc)
}

func (this *nodeware) LenByAddr(addr string) (i int32) {
	this.tlMap.Range(func(k *tlContext, _ int8) bool {
		if addr == k.remoteAddr {
			i++
		}
		return true
	})
	return
}

func (this *nodeware) GetRemoteUUIDS() (_r []int64) {
	_r = make([]int64, 0)
	_m := make(map[int64]int8, 0)
	this.tlMap.Range(func(k *tlContext, _ int8) bool {
		if _, ok := _m[k.remoteUuid]; !ok {
			_m[k.remoteUuid] = 0
			_r = append(_r, k.remoteUuid)
		}
		return true
	})
	return
}

func (this *nodeware) GetUUIDNode() (_r *Node) {
	_r = &Node{UUID: sys.UUID, Addr: sys.CSADDR}
	if strings.Index(_r.Addr, ":") == 0 {
		_r.Addr = sys.Bind + _r.Addr
	}
	_r.Nodekv = make(map[int64]string, 0)
	this.tlMap.Range(func(k *tlContext, _ int8) bool {
		if strings.Index(k.remoteAddr, ":") == 0 {
			_r.Nodekv[k.remoteUuid] = k.remoteIP + k.remoteAddr
		} else {
			_r.Nodekv[k.remoteUuid] = k.remoteAddr
		}
		return true
	})
	return
}

func (this *nodeware) GetALLUUID() []int64 {
	return append(this.GetRemoteUUIDS(), sys.UUID)
}

func (this *nodeware) GetAllTlContext() (tcs []*tlContext) {
	tcs = make([]*tlContext, 0)
	this.tlMap.Range(func(k *tlContext, _ int8) bool {
		tcs = append(tcs, k)
		return true
	})
	return
}

func (this *nodeware) GetTlContext(uuid int64) (tc *tlContext) {
	this.tlMap.Range(func(k *tlContext, _ int8) bool {
		if uuid == k.remoteUuid {
			tc = k
			return false
		}
		return true
	})
	return
}

func (this *nodeware) getRemoteNodes() (_r []*sys.RemoteNode) {
	_r = make([]*sys.RemoteNode, 0)
	this.tlMap.Range(func(k *tlContext, _ int8) bool {
		r := &sys.RemoteNode{Addr: k.remoteAddr, UUID: k.remoteUuid, CSNUM: k.remoteCsNum, Host: k.remoteIP}
		_r = append(_r, r)
		return true
	})
	return
}

func (this *nodeware) IsLocal(v string) bool {
	return this.csnode(v) == sys.UUID
}

func (this *nodeware) addCsMapNode(uuid int64) {
	if sys.WssLen != nil {
		if n := sys.WssLen(); n > 1<<10 {
			t := time.Duration(n / (1 << 13))
			<-time.After(t * time.Second)
		}
	}
	this.csMap.Add(uuid)
	this.csuserBatch(uuid)
}

func (this *nodeware) getCsList() (_r []int64) {
	if this._cacheCsList != nil {
		return this._cacheCsList
	} else {
		_r = this.GetALLUUID()
		this._cacheCsList = _r
	}
	return
}

func (this *nodeware) csnode(to string) (_r int64) {
	if n, ok := this.csMap.GetStr(to); ok {
		if this.hasUUID(n) {
			return n
		} else {
			if sys.MaxBackup > 0 {
				if ns, ok := this.bkuuid(to); ok {
					for _, n := range ns {
						if this.hasUUID(n) {
							return n
						}
					}
				}
			}
		}
	}
	return sys.UUID
}

func (this *nodeware) bkuuid(node string) ([]int64, bool) {
	return this.csMap.GetNextNodeStr(node, sys.MaxBackup)
}

func (this *nodeware) getcsnode(_node string) (m map[int64]int8) {
	m = map[int64]int8{}
	if node := this.csnode(_node); node != sys.UUID {
		if cr, err := this.csReqHandle(node, 0, false, &CsBean{RType: 1, Bsm2: map[string][]int64{_node: nil}}); err == nil {
			if li, b := cr.Bsm2[_node]; b {
				nodeCache.Add(_node, li)
				for _, u := range li {
					if _, ok := m[u]; ok {
						continue
					} else {
						m[u] = 0
					}
				}
			}
		}
	} else {
		if um, ok := csuMap.Get(_node); ok {
			um.Range(func(k int64, _ int8) bool {
				if !nodeWare.hasUUID(k) {
					um.Del(k)
				} else {
					m[k] = 0
				}
				return true
			})
			if um.Len() == 0 {
				csuMap.Del(_node)
			}
		}
	}
	return
}

func (this *nodeware) cswssinfo(node string) (b []byte) {
	li := this.getcsnode(node)
	b = make([]byte, 0)
	for u := range li {
		if u != sys.UUID {
			if cr, err := this.csReqHandle(u, 0, false, &CsBean{RType: 2, Bsm: map[string][]byte{node: nil}}); err == nil {
				if cr.Bsm != nil {
					b, _ = cr.Bsm[node]
				}
			}
		}
	}
	return
}

/*********************************************************************/
func (this *nodeware) csmessage(tm *TimMessage, transType int8) (ok bool) {
	toNode := tm.ToTid.Node
	switch transType {
	case sys.TRANS_SOURCE:
		sl := false
		if node := this.csnode(toNode); node != sys.UUID {
			if li, b := nodeCache.Get(toNode); b {
				rm := false
				if rm, ok = this.csbsMessageByCache(tm, toNode, li); ok || (!ok && !rm) {
					return
				}
			}
			if cr, err := this.csReqHandle(node, 0, false, &CsBean{RType: 1, Bsm2: map[string][]int64{toNode: nil}}); err == nil {
				if li, b := cr.Bsm2[toNode]; b {
					nodeCache.Add(toNode, li)
					m := map[int64]int8{}
					for _, u := range li {
						if _, ok := m[u]; ok {
							continue
						} else {
							m[u] = 0
						}
						if u != sys.UUID {
							if _, err := this.csbsHandle(u, &CsBs{Bs: TEncode(tm), TransType: sys.TRANS_GOAL, BsType: sys.CB_MESSAGE}); err == nil {
								ok = true
							}
						} else if !sl {
							_, sl = sys.TimMessageProcessor(tm, sys.TRANS_GOAL), true
						}
					}
				}
			}
			if !ok && !sl {
				_, sl = sys.TimMessageProcessor(tm, sys.TRANS_GOAL), true
			}
		} else {
			sys.TimMessageProcessor(tm, sys.TRANS_CONSISHASH)
		}
	case sys.TRANS_CONSISHASH:
		sl := false
		if cm := getcsu(toNode); cm != nil {
			cm.Range(func(k int64, _ int8) bool {
				if k != sys.UUID {
					if _, err := this.csbsHandle(k, &CsBs{Bs: TEncode(tm), TransType: sys.TRANS_GOAL, BsType: sys.CB_MESSAGE}); err == nil {
						ok = true
					}
				} else {
					_, sl = sys.TimMessageProcessor(tm, sys.TRANS_GOAL), true
				}
				return true
			})
		}
		if !ok && !sl {
			sys.TimMessageProcessor(tm, sys.TRANS_GOAL)
		}
	case sys.TRANS_STAFF:
		sys.TimMessageProcessor(tm, sys.TRANS_GOAL)
	}
	return
}

func (this *nodeware) csbsMessageByCache(tm *TimMessage, tonode string, nodes []int64) (rmc, ok bool) {
	m, cache := map[int64]int8{}, true
	for _, u := range nodes {
		if _, b := m[u]; b {
			continue
		} else {
			m[u] = 0
		}
		if u != sys.UUID {
			if b, err := this.csbsHandle(u, &CsBs{Bs: TEncode(tm), Node: []byte(tonode), TransType: sys.TRANS_GOAL, BsType: sys.CB_MESSAGE, Cache: &cache}); err == nil {
				if b {
					ok = true
				} else {
					rmc = true
					nodeCache.Remove(tonode)
				}
			}
		} else if !ok {
			_, ok = sys.TimMessageProcessor(tm, sys.TRANS_GOAL), true
		}
	}
	return
}

func (this *nodeware) cspresence(tp *TimPresence, transType int8) (ok bool) {
	if tp.ToTid == nil && tp.ToList == nil {
		return
	}
	switch transType {
	case sys.TRANS_SOURCE:
		if tp.ToTid == nil {
			return
		}
		toNode := tp.ToTid.Node
		sl := false
		if node := this.csnode(toNode); node != sys.UUID {
			if cr, err := this.csReqHandle(node, 0, false, &CsBean{RType: 1, Bsm2: map[string][]int64{toNode: nil}}); err == nil {
				if li, b := cr.Bsm2[toNode]; b {
					nodeCache.Add(toNode, li)
					m := map[int64]int8{}
					for _, u := range li {
						if _, ok := m[u]; ok {
							continue
						} else {
							m[u] = 0
						}
						if u != sys.UUID {
							if _, err := this.csbsHandle(u, &CsBs{Bs: TEncode(tp), TransType: sys.TRANS_GOAL, BsType: sys.CB_PRESENCE}); err == nil {
								ok = true
							}
						} else if !sl {
							_, sl = sys.TimPresenceProcessor(tp, sys.TRANS_GOAL), true
						}
					}
				}
			}
			if !ok && !sl {
				_, sl = sys.TimPresenceProcessor(tp, sys.TRANS_GOAL), true
			}
		} else {
			sys.TimPresenceProcessor(tp, sys.TRANS_CONSISHASH)
		}

	case sys.TRANS_CONSISHASH:
		if tp.ToTid == nil {
			return
		}
		toNode := tp.ToTid.Node
		// sl := false
		if cm := getcsu(toNode); cm != nil {
			cm.Range(func(k int64, _ int8) bool {
				if k != sys.UUID {
					if _, err := this.csbsHandle(k, &CsBs{Bs: TEncode(tp), TransType: sys.TRANS_GOAL, BsType: sys.CB_PRESENCE}); err == nil {
						ok = true
					}
				} else {
					sys.TimPresenceProcessor(tp, sys.TRANS_GOAL)
				}
				return true
			})
		}
		// if !ok && !sl {
		// 	sys.TimPresenceProcessor(tp, sys.TRANS_GOAL)
		// }
	case sys.TRANS_STAFF:
		if tp.ToList != nil {
			if len(tp.ToList) > len(this.csMap.Nodes())/2 && len(this.GetRemoteUUIDS()) > 0 {
				for _, u := range tp.ToList {
					_tp := &TimPresence{ID: tp.ID, FromTid: tp.FromTid, ToTid: &Tid{Node: u}, SubStatus: tp.SubStatus, Offline: tp.Offline, Show: tp.Show, Status: tp.Status, Extend: tp.Extend, Extra: tp.Extra}
					sys.TimPresenceProcessor(_tp, sys.TRANS_GOAL)
				}
				for _, uuid := range this.GetRemoteUUIDS() {
					this.csbsHandle(uuid, &CsBs{Bs: TEncode(tp), TransType: sys.TRANS_STAFF, BsType: sys.CB_PRESENCE})
				}
			} else {
				for _, u := range tp.ToList {
					_tp := &TimPresence{ID: tp.ID, FromTid: tp.FromTid, ToTid: &Tid{Node: u}, SubStatus: tp.SubStatus, Offline: tp.Offline, Show: tp.Show, Status: tp.Status, Extend: tp.Extend, Extra: tp.Extra}
					sys.TimPresenceProcessor(_tp, sys.TRANS_SOURCE)
				}
			}
		} else if tp.ToTid != nil {
			sys.TimPresenceProcessor(tp, sys.TRANS_SOURCE)
		}
	}
	return
}

func (this *nodeware) csuser(_node string, on bool, wsId int64) (err error) {
	stat := int8(2)
	if on {
		stat = 1
		addcsu(_node, sys.UUID)
	} else {
		delcsu(_node, sys.UUID)
		vgate.VGate.DelNode("", wsId)
	}
	if len(this.getCsList()) > 1 {
		if node, ok := this.csMap.GetStr(_node); ok && node != sys.UUID {
			err = this.csuserHandle(node, &CsUser{Node: map[string]int8{_node: 0}, Stat: stat})
		} else {
			bkCsuserBatch(map[string]int8{_node: 0}, sys.UUID, stat)
		}
	}
	return
}

func (this *nodeware) csuserBatch(uuid int64) {
	if sys.WssList == nil || uuid == sys.UUID {
		return
	}
	us, _ := sys.WssList(0, 0)
	m := map[string]int8{}
	for _, u := range us {
		if node, _ := this.csMap.GetStr(u.Node); node == uuid {
			m[u.Node] = 0
		} else if sys.MaxBackup > 0 {
			if bckNode, ok := this.bkuuid(u.Node); ok && len(bckNode) > 0 {
				if util.ContainInt(bckNode, uuid) {
					m[u.Node] = 0
				}
			}
		}
	}
	if len(m) > 0 {
		this.csuserHandle(uuid, &CsUser{Node: m, Stat: 1})
	}
	vgate.VGate.Nodes().Range(func(k string, v *vgate.VRoom) bool {
		if !v.Expires() {
			if node, _ := this.csMap.GetStr(k); node == uuid {
				vb := &VBean{Rtype: 1, Vnode: k, FoundNode: &v.FoundNode}
				this.csVrHandle(node, 0, vb)
			} else if sys.MaxBackup > 0 {
				if bckNode, ok := this.bkuuid(k); ok && len(bckNode) > 0 {
					if util.ContainInt(bckNode, uuid) {
						vb := &VBean{Rtype: 1, Vnode: k, FoundNode: &v.FoundNode}
						this.csVrHandle(node, 0, vb)
					}
				}
			}
		} else {
			vgate.VGate.Remove(v.FoundNode, k)
		}
		return true
	})
}

func (this *nodeware) close() {
	this.tlMap.Range(func(k *tlContext, _ int8) bool {
		k._do_reconn = true
		this.del(k)
		return true
	})
}

func (this *nodeware) csVbean(vb *VBean) (b bool) {
	if node, _ := this.csMap.GetStr(vb.Vnode); node != sys.UUID {
		if vb.Rtype == 7 {
			go sys.TimSteamProcessor(vb)
		}
		if b = this.csVrHandle(node, 0, vb) == nil; !b && vb.Rtype == 5 {
			if nodes, ok := nodeWare.bkuuid(vb.Vnode); ok && len(nodes) > 0 {
				for _, v := range nodes {
					if b = this.csVrHandle(v, 0, vb) == nil; b {
						break
					}
				}
			}
		}
	} else {
		if vb.Rtype != 5 && vb.Rtype != 7 {
			bkCsVr(vb, sys.UUID)
		} else if vb.Rtype == 7 {
			processVBean(vb, sys.UUID)
		}
		b = true
	}
	return
}

var csuMap = NewMapL[string, *MapL[int64, int8]]()

func hascsu(node string) (_r bool) {
	if m, ok := csuMap.Get(node); ok {
		_r = m.Has(sys.UUID)
	}
	return
}

func addcsu(node string, uuid int64) {
	strLock.Lock(node)
	defer strLock.Unlock(node)
	if m, ok := csuMap.Get(node); ok {
		m.Put(uuid, 0)
	} else {
		m = NewMapL[int64, int8]()
		m.Put(uuid, 0)
		csuMap.Put(node, m)
	}
}

func delcsu(node string, uuid int64) {
	strLock.Lock(node)
	defer strLock.Unlock(node)
	if m, ok := csuMap.Get(node); ok {
		m.Del(uuid)
		if m.Len() == 0 {
			csuMap.Del(node)
		}
	}
}

func getcsu(node string) *MapL[int64, int8] {
	if m, ok := csuMap.Get(node); ok {
		m.Range(func(k int64, _ int8) bool {
			if !nodeWare.hasUUID(k) {
				m.Del(k)
			}
			return true
		})
		if m.Len() == 0 {
			csuMap.Del(node)
		} else {
			return m
		}
	}
	return nil
}

func (this *nodeware) csuserHandle(uuid int64, cu *CsUser) (err error) {
	defer util.Recover()
	sys.Stat.CReqDo()
	defer sys.Stat.CReqDone()
	if uuid == sys.UUID {
		return
	}
	limit := 30
	for limit > 0 {
		limit--
		sendId := RandId()
		ch := await.Get(sendId)
		if tc := this.GetTlContext(uuid); tc != nil {
			err = tc.iface.CsUser(context.TODO(), sendId, cu)
		}
		select {
		case <-ch:
			err = nil
			goto END
		case <-time.After(sys.WaitTimeout):
			err = sys.ERR_OVERTIME.Error()
		}
		await.DelAndClose(sendId)
	}
END:
	return
}

func (this *nodeware) csbsHandle(uuid int64, cb *CsBs) (ok bool, err error) {
	defer util.Recover()
	sys.Stat.CReqDo()
	defer sys.Stat.CReqDone()
	if uuid == sys.UUID {
		return
	}
	c := time.Duration(0)
	for c < 30 {
		c++
		sendId := RandId()
		ch := await.Get(sendId)
		if tc := this.GetTlContext(uuid); tc != nil {
			err = tc.iface.CsBs(context.TODO(), sendId, cb)
		}
		select {
		case _r := <-ch:
			ok = _r == 1
			err = nil
			goto END
		case <-time.After(c * sys.WaitTimeout):
			err = sys.ERR_OVERTIME.Error()
		}
		await.DelAndClose(sendId)
	}
END:
	return
}

func (this *nodeware) csReqHandle(uuid, sendId int64, ack bool, cb *CsBean) (_r *CsBean, err error) {
	defer util.Recover()
	sys.Stat.CReqDo()
	defer sys.Stat.CReqDone()
	if uuid == sys.UUID {
		return
	}
	limit := 30
	for limit > 0 {
		limit--
		if !ack {
			sendId = RandId()
		}
		if tc := this.GetTlContext(uuid); tc != nil {
			err = tc.iface.CsReq(context.TODO(), sendId, ack, cb)
		}
		if !ack {
			ch := awaitCsBean.Get(sendId)
			select {
			case _r = <-ch:
				err = nil
				goto END
			case <-time.After(sys.WaitTimeout):
				err = sys.ERR_OVERTIME.Error()
			}
			awaitCsBean.DelAndClose(sendId)
		}
	}
END:
	return
}

func (this *nodeware) csVrHandle(uuid, sendId int64, vb *VBean) (err error) {
	defer util.Recover()
	sys.Stat.CReqDo()
	defer sys.Stat.CReqDone()
	if uuid == sys.UUID {
		return
	}
	if tc := this.GetTlContext(uuid); tc != nil {
		err = tc.iface.CsVr(context.TODO(), sendId, vb)
	} else {
		err = sys.ERR_NOEXIST.Error()
	}
	if sendId > 0 {
		ch := await.Get(sendId)
		select {
		case <-ch:
			err = nil
			goto END
		case <-time.After(sys.WaitTimeout):
			err = sys.ERR_OVERTIME.Error()
		}
		await.DelAndClose(sendId)
	}
END:
	return
}

func (this *nodeware) csuTicker() {
	tk := time.NewTicker(10 * time.Minute)
	maskMap := map[string]int{}
	for {
		select {
		case <-tk.C:
			func() {
				defer util.Recover()
				csuMap.Range(func(k string, cm *MapL[int64, int8]) bool {
					mark := true
					if node, ok := this.csMap.GetStr(k); ok {
						if node != sys.UUID {
							if sys.MaxBackup > 0 {
								if ns, ok := this.bkuuid(k); ok {
									if util.ContainInt(ns, sys.UUID) {
										mark = false
									}
								}
							}
						} else {
							mark = false
						}
					}

					if mark {
						if i, ok := maskMap[k]; ok {
							maskMap[k] = i + 1
						} else {
							maskMap[k] = 1
						}
					}

					cm.Range(func(k int64, _ int8) bool {
						if !this.hasUUID(k) {
							cm.Del(k)
						}
						return true
					})

					if cm.Len() == 0 {
						csuMap.Del(k)
						delete(maskMap, k)
					}

					return true
				})

				for k, v := range maskMap {
					if v >= 3 {
						delete(maskMap, k)
						csuMap.Del(k)
					}
				}

			}()
		}
	}
}
