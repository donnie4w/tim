// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package level1

import (
	"context"
	"github.com/donnie4w/tim/errs"
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

func (nd *nodeware) add(tc *tlContext) (err error) {
	nd.mux.Lock(tc.remoteUuid)
	defer nd.mux.Unlock(tc.remoteUuid)
	if nd.tlMap.Has(tc) {
		return errs.ERR_UUID_REUSE.Error()
	}
	if sys.UUID == tc.remoteUuid || nd.hasUUID(tc.remoteUuid) {
		tc.CloseAndEnd()
		return errs.ERR_UUID_REUSE.Error()
	}
	nd.tlMap.Put(tc, 1)
	_unaccess.Del(tc.remoteUuid)
	nd._cacheMap.Put(tc.remoteUuid, 0)
	nd._cacheCsList = nil
	go nd.addCsMapNode(tc.remoteUuid)
	return
}

func (nd *nodeware) has(tc *tlContext) bool {
	return nd.tlMap.Has(tc)
}

func (nd *nodeware) hasUUID(uuid int64) (b bool) {
	return nd._cacheMap.Has(uuid) || uuid == sys.UUID
}

func (nd *nodeware) del(tc *tlContext) {
	if !tc.isClose {
		tc.Close()
	}
	nd.tlMap.Del(tc)
	nd._cacheMap.Del(tc.remoteUuid)
	nd._cacheCsList = nil
	if !nd.hasUUID(tc.remoteUuid) {
		nd.csMap.Del(tc.remoteUuid)
	}
}

func (nd *nodeware) delAndNoReconn(tc *tlContext) {
	tc._do_reconn = true
	nd.del(tc)
}

func (nd *nodeware) LenByAddr(addr string) (i int32) {
	nd.tlMap.Range(func(k *tlContext, _ int8) bool {
		if addr == k.remoteAddr {
			i++
		}
		return true
	})
	return
}

func (nd *nodeware) GetRemoteUUIDS() (_r []int64) {
	_r = make([]int64, 0)
	_m := make(map[int64]int8, 0)
	nd.tlMap.Range(func(k *tlContext, _ int8) bool {
		if _, ok := _m[k.remoteUuid]; !ok {
			_m[k.remoteUuid] = 0
			_r = append(_r, k.remoteUuid)
		}
		return true
	})
	return
}

func (nd *nodeware) GetUUIDNode() (_r *Node) {
	_r = &Node{UUID: sys.UUID, Addr: sys.CSADDR}
	if strings.Index(_r.Addr, ":") == 0 {
		_r.Addr = sys.Bind + _r.Addr
	}
	_r.Nodekv = make(map[int64]string, 0)
	nd.tlMap.Range(func(k *tlContext, _ int8) bool {
		if strings.Index(k.remoteAddr, ":") == 0 {
			_r.Nodekv[k.remoteUuid] = k.remoteIP + k.remoteAddr
		} else {
			_r.Nodekv[k.remoteUuid] = k.remoteAddr
		}
		return true
	})
	return
}

func (nd *nodeware) GetALLUUID() []int64 {
	return append(nd.GetRemoteUUIDS(), sys.UUID)
}

func (nd *nodeware) GetAllTlContext() (tcs []*tlContext) {
	tcs = make([]*tlContext, 0)
	nd.tlMap.Range(func(k *tlContext, _ int8) bool {
		tcs = append(tcs, k)
		return true
	})
	return
}

func (nd *nodeware) GetTlContext(uuid int64) (tc *tlContext) {
	nd.tlMap.Range(func(k *tlContext, _ int8) bool {
		if uuid == k.remoteUuid {
			tc = k
			return false
		}
		return true
	})
	return
}

func (nd *nodeware) getRemoteNodes() (_r []*RemoteNode) {
	_r = make([]*RemoteNode, 0)
	nd.tlMap.Range(func(k *tlContext, _ int8) bool {
		r := &RemoteNode{Addr: k.remoteAddr, UUID: k.remoteUuid, CSNUM: k.remoteCsNum, Host: k.remoteIP}
		_r = append(_r, r)
		return true
	})
	return
}

func (nd *nodeware) IsLocal(v string) bool {
	return nd.csnode(v) == sys.UUID
}

func (nd *nodeware) addCsMapNode(uuid int64) {
	if sys.WssLen != nil {
		if n := sys.WssLen(); n > 1<<10 {
			t := time.Duration(n / (1 << 13))
			<-time.After(t * time.Second)
		}
	}
	nd.csMap.Add(uuid)
	nd.csuserBatch(uuid)
}

func (nd *nodeware) getCsList() (_r []int64) {
	if nd._cacheCsList != nil {
		return nd._cacheCsList
	} else {
		_r = nd.GetALLUUID()
		nd._cacheCsList = _r
	}
	return
}

func (nd *nodeware) csnode(to string) (_r int64) {
	if n, ok := nd.csMap.GetStr(to); ok {
		if nd.hasUUID(n) {
			return n
		} else {
			if sys.MaxBackup > 0 {
				if ns, ok := nd.bkuuid(to); ok {
					for _, n := range ns {
						if nd.hasUUID(n) {
							return n
						}
					}
				}
			}
		}
	}
	return sys.UUID
}

func (nd *nodeware) bkuuid(node string) ([]int64, bool) {
	return nd.csMap.GetNextNodeStr(node, sys.MaxBackup)
}

func (nd *nodeware) getcsnode(_node string) (m map[int64]int8) {
	m = map[int64]int8{}
	if node := nd.csnode(_node); node != sys.UUID {
		if cr, err := nd.csReqHandle(node, 0, false, &CsBean{RType: 1, Bsm2: map[string][]int64{_node: nil}}); err == nil {
			if li, b := cr.Bsm2[_node]; b {
				nodeCache.Put(_node, li)
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

func (nd *nodeware) cswssinfo(node string) (b []byte) {
	li := nd.getcsnode(node)
	b = make([]byte, 0)
	for u := range li {
		if u != sys.UUID {
			if cr, err := nd.csReqHandle(u, 0, false, &CsBean{RType: 2, Bsm: map[string][]byte{node: nil}}); err == nil {
				if cr.Bsm != nil {
					b, _ = cr.Bsm[node]
				}
			}
		}
	}
	return
}

/*********************************************************************/
func (nd *nodeware) csmessage(tm *TimMessage, transType int8) (ok bool) {
	toNode := tm.ToTid.Node
	switch transType {
	case sys.TRANS_SOURCE:
		sl := false
		if node := nd.csnode(toNode); node != sys.UUID {
			if li, b := nodeCache.Get(toNode); b {
				rm := false
				if rm, ok = nd.csbsMessageByCache(tm, toNode, li); ok || (!ok && !rm) {
					return
				}
			}
			if cr, err := nd.csReqHandle(node, 0, false, &CsBean{RType: 1, Bsm2: map[string][]int64{toNode: nil}}); err == nil {
				if li, b := cr.Bsm2[toNode]; b {
					nodeCache.Put(toNode, li)
					m := map[int64]int8{}
					for _, u := range li {
						if _, ok := m[u]; ok {
							continue
						} else {
							m[u] = 0
						}
						if u != sys.UUID {
							if _, err := nd.csbsHandle(u, &CsBs{Bs: TEncode(tm), TransType: sys.TRANS_GOAL, BsType: sys.CB_MESSAGE}); err == nil {
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
					if _, err := nd.csbsHandle(k, &CsBs{Bs: TEncode(tm), TransType: sys.TRANS_GOAL, BsType: sys.CB_MESSAGE}); err == nil {
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

func (nd *nodeware) csbsMessageByCache(tm *TimMessage, tonode string, nodes []int64) (rmc, ok bool) {
	m, cache := map[int64]int8{}, true
	for _, u := range nodes {
		if _, b := m[u]; b {
			continue
		} else {
			m[u] = 0
		}
		if u != sys.UUID {
			if b, err := nd.csbsHandle(u, &CsBs{Bs: TEncode(tm), Node: []byte(tonode), TransType: sys.TRANS_GOAL, BsType: sys.CB_MESSAGE, Cache: &cache}); err == nil {
				if b {
					ok = true
				} else {
					rmc = true
					nodeCache.Del(tonode)
				}
			}
		} else if !ok {
			_, ok = sys.TimMessageProcessor(tm, sys.TRANS_GOAL), true
		}
	}
	return
}

func (nd *nodeware) cspresence(tp *TimPresence, transType int8) (ok bool) {
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
		if node := nd.csnode(toNode); node != sys.UUID {
			if cr, err := nd.csReqHandle(node, 0, false, &CsBean{RType: 1, Bsm2: map[string][]int64{toNode: nil}}); err == nil {
				if li, b := cr.Bsm2[toNode]; b {
					nodeCache.Put(toNode, li)
					m := map[int64]int8{}
					for _, u := range li {
						if _, ok := m[u]; ok {
							continue
						} else {
							m[u] = 0
						}
						if u != sys.UUID {
							if _, err := nd.csbsHandle(u, &CsBs{Bs: TEncode(tp), TransType: sys.TRANS_GOAL, BsType: sys.CB_PRESENCE}); err == nil {
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
					if _, err := nd.csbsHandle(k, &CsBs{Bs: TEncode(tp), TransType: sys.TRANS_GOAL, BsType: sys.CB_PRESENCE}); err == nil {
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
			if len(tp.ToList) > len(nd.csMap.Nodes())/2 && len(nd.GetRemoteUUIDS()) > 0 {
				for _, u := range tp.ToList {
					_tp := &TimPresence{ID: tp.ID, FromTid: tp.FromTid, ToTid: &Tid{Node: u}, SubStatus: tp.SubStatus, Offline: tp.Offline, Show: tp.Show, Status: tp.Status, Extend: tp.Extend, Extra: tp.Extra}
					sys.TimPresenceProcessor(_tp, sys.TRANS_GOAL)
				}
				for _, uuid := range nd.GetRemoteUUIDS() {
					nd.csbsHandle(uuid, &CsBs{Bs: TEncode(tp), TransType: sys.TRANS_STAFF, BsType: sys.CB_PRESENCE})
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

func (nd *nodeware) csuser(_node string, on bool, wsId int64) (err error) {
	stat := int8(2)
	if on {
		stat = 1
		addcsu(_node, sys.UUID)
	} else {
		delcsu(_node, sys.UUID)
		vgate.VGate.UnSubWithWsId(wsId)
	}
	if len(nd.getCsList()) > 1 {
		if node, ok := nd.csMap.GetStr(_node); ok && node != sys.UUID {
			err = nd.csuserHandle(node, &CsUser{Node: map[string]int8{_node: 0}, Stat: stat})
		} else {
			bkCsuserBatch(map[string]int8{_node: 0}, sys.UUID, stat)
		}
	}
	return
}

func (nd *nodeware) csuserBatch(uuid int64) {
	if sys.WssList == nil || uuid == sys.UUID {
		return
	}
	us, _ := sys.WssList(0, 0)
	m := map[string]int8{}
	for _, u := range us {
		if node, _ := nd.csMap.GetStr(u.Node); node == uuid {
			m[u.Node] = 0
		} else if sys.MaxBackup > 0 {
			if bckNode, ok := nd.bkuuid(u.Node); ok && len(bckNode) > 0 {
				if util.ContainInt(bckNode, uuid) {
					m[u.Node] = 0
				}
			}
		}
	}
	if len(m) > 0 {
		nd.csuserHandle(uuid, &CsUser{Node: m, Stat: 1})
	}
	vgate.VGate.Nodes().Range(func(k string, v *vgate.VRoom) bool {
		if !v.Expires() {
			if node, _ := nd.csMap.GetStr(k); node == uuid {
				vb := &VBean{Rtype: 1, Vnode: k, FoundNode: &v.FoundNode}
				nd.csVrHandle(node, 0, vb)
			} else if sys.MaxBackup > 0 {
				if bckNode, ok := nd.bkuuid(k); ok && len(bckNode) > 0 {
					if util.ContainInt(bckNode, uuid) {
						vb := &VBean{Rtype: 1, Vnode: k, FoundNode: &v.FoundNode}
						nd.csVrHandle(node, 0, vb)
					}
				}
			}
		} else {
			vgate.VGate.Remove(v.FoundNode, k)
		}
		return true
	})
}

func (nd *nodeware) close() {
	nd.tlMap.Range(func(k *tlContext, _ int8) bool {
		k._do_reconn = true
		nd.del(k)
		return true
	})
}

func (nd *nodeware) csVbean(vb *VBean) (b bool) {
	if nodeId, _ := nd.csMap.GetStr(vb.Vnode); nodeId != sys.UUID {
		if vb.Rtype == int8(sys.VROOM_MESSAGE) {
			go sys.TimSteamProcessor(vb, sys.TRANS_SOURCE)
		}
		if b = nd.csVrHandle(nodeId, 0, vb) == nil; !b && vb.Rtype == int8(sys.VROOM_SUB) {
			if nodes, ok := nodeWare.bkuuid(vb.Vnode); ok && len(nodes) > 0 {
				for _, v := range nodes {
					if b = nd.csVrHandle(v, 0, vb) == nil; b {
						break
					}
				}
			}
		}
	} else {
		if vb.Rtype != int8(sys.VROOM_SUB) && vb.Rtype != int8(sys.VROOM_MESSAGE) {
			bkCsVr(vb, sys.UUID)
		} else if vb.Rtype == int8(sys.VROOM_MESSAGE) {
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

func (nd *nodeware) csuserHandle(uuid int64, cu *CsUser) (err error) {
	defer util.Recover()
	sys.Stat.CReqDo()
	defer sys.Stat.CReqDone()
	if uuid == sys.UUID {
		return
	}
	limit := 30
	for limit > 0 {
		limit--
		sendId := UUID64()
		ch := await.Get(sendId)
		if tc := nd.GetTlContext(uuid); tc != nil {
			err = tc.csnet.CsUser(context.TODO(), sendId, cu)
		}
		select {
		case <-ch:
			err = nil
			goto END
		case <-time.After(sys.WaitTimeout):
			err = errs.ERR_OVERTIME.Error()
		}
		await.Close(sendId)
	}
END:
	return
}

func (nd *nodeware) csbsHandle(uuid int64, cb *CsBs) (ok bool, err error) {
	defer util.Recover()
	sys.Stat.CReqDo()
	defer sys.Stat.CReqDone()
	if uuid == sys.UUID {
		return
	}
	c := time.Duration(0)
	for c < 30 {
		c++
		sendId := UUID64()
		ch := await.Get(sendId)
		if tc := nd.GetTlContext(uuid); tc != nil {
			err = tc.csnet.CsBs(context.TODO(), sendId, cb)
		}
		select {
		case _r := <-ch:
			ok = _r == 1
			err = nil
			goto END
		case <-time.After(c * sys.WaitTimeout):
			err = errs.ERR_OVERTIME.Error()
		}
		await.Close(sendId)
	}
END:
	return
}

func (nd *nodeware) csReqHandle(uuid, sendId int64, ack bool, cb *CsBean) (_r *CsBean, err error) {
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
			sendId = UUID64()
		}
		if tc := nd.GetTlContext(uuid); tc != nil {
			err = tc.csnet.CsReq(context.TODO(), sendId, ack, cb)
		}
		if !ack {
			ch := awaitCsBean.Get(sendId)
			select {
			case _r = <-ch:
				err = nil
				goto END
			case <-time.After(sys.WaitTimeout):
				err = errs.ERR_OVERTIME.Error()
			}
			awaitCsBean.Close(sendId)
		}
	}
END:
	return
}

func (nd *nodeware) csVrHandle(uuid, sendId int64, vb *VBean) (err error) {
	defer util.Recover()
	sys.Stat.CReqDo()
	defer sys.Stat.CReqDone()
	if uuid == sys.UUID {
		return
	}
	if tc := nd.GetTlContext(uuid); tc != nil {
		err = tc.csnet.CsVr(context.TODO(), sendId, vb)
	} else {
		err = errs.ERR_NOEXIST.Error()
	}
	if sendId > 0 {
		ch := await.Get(sendId)
		select {
		case <-ch:
			err = nil
			goto END
		case <-time.After(sys.WaitTimeout):
			err = errs.ERR_OVERTIME.Error()
		}
		await.Close(sendId)
	}
END:
	return
}

func (nd *nodeware) csuTicker() {
	tk := time.NewTicker(10 * time.Minute)
	maskMap := map[string]int{}
	for {
		select {
		case <-tk.C:
			func() {
				defer util.Recover()
				csuMap.Range(func(k string, cm *MapL[int64, int8]) bool {
					mark := true
					if node, ok := nd.csMap.GetStr(k); ok {
						if node != sys.UUID {
							if sys.MaxBackup > 0 {
								if ns, ok := nd.bkuuid(k); ok {
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
						if !nd.hasUUID(k) {
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
