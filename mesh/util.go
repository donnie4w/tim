// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package mesh

import (
	. "github.com/donnie4w/gofer/buffer"
	. "github.com/donnie4w/gofer/hashmap"
	. "github.com/donnie4w/gofer/util"
	"github.com/donnie4w/tim/errs"
	. "github.com/donnie4w/tim/stub"
	"github.com/donnie4w/tim/sys"
	"github.com/donnie4w/tim/vgate"
	"sync/atomic"
)

func bkCsuserBatch(nodem map[string]int8, srcUuid int64, stat int8) {
	if sys.MaxBackup > 0 {
		bkm := map[int64]map[string]int64{}
		for k := range nodem {
			if bckNode, ok := nodeWare.bkuuid(k); ok && len(bckNode) > 0 {
				for _, n := range bckNode {
					if n != srcUuid {
						if bn, ok := bkm[n]; ok {
							bn[k] = srcUuid
						} else {
							bkm[n] = map[string]int64{k: srcUuid}
						}
					}
				}
			}
		}
		for k, bn := range bkm {
			nodeWare.csuserHandle(k, &CsUser{Stat: stat, BkNode: bn})
		}
	}
}

func bkCsVr(vb *VBean, srcuuid int64) {
	if nodes, ok := nodeWare.bkuuid(vb.Vnode); ok && len(nodes) > 0 {
		vb.Rtype = vb.Rtype + 50
		for _, u := range nodes {
			if u != srcuuid {
				nodeWare.csVrHandle(u, 0, vb)
			}
		}
	}
}

func processVBean(vrb *VBean, srcuuid int64) {
	switch vrb.Rtype {
	case 1:
		vgate.VGate.Register(*vrb.FoundNode, vrb.Vnode)
		bkCsVr(vrb, srcuuid)
	case 2:
		vgate.VGate.Remove(*vrb.FoundNode, vrb.Vnode)
		bkCsVr(vrb, srcuuid)
	case 3:
		//vgate.VGate.AddAuth(vrb.Vnode, *vrb.FoundNode, *vrb.Rnode)
		//bkCsVr(vrb, srcuuid)
	case 4:
		//vgate.VGate.DelAuth(vrb.Vnode, *vrb.FoundNode, *vrb.Rnode)
		//bkCsVr(vrb, srcuuid)
	case 5:
		if _, b := vgate.VGate.GetVroom(vrb.Vnode); b {
			vgate.VGate.Sub(vrb.Vnode, srcuuid, 0)
		} else {
			nodeWare.csVrHandle(srcuuid, 0, &VBean{Rtype: 10, Vnode: vrb.Vnode, Rnode: vrb.Rnode})
		}
	case 6:
		vgate.VGate.UnSubWithUUID(vrb.Vnode, srcuuid)
	case 7:
		go sys.TimSteamProcessor(vrb, sys.TRANS_SOURCE)
		if vr, ok := vgate.VGate.GetVroom(vrb.Vnode); ok {
			if !vr.AuthStream(*vrb.Rnode) {
				nodeWare.csVrHandle(srcuuid, 0, &VBean{Rtype: 9, Vnode: vrb.Vnode, Rnode: vrb.Rnode})
				return
			}
			vr.Updatetime()
			vrb.Rtype = 50 + vrb.Rtype
			m := map[int64]int8{}
			vgate.VGate.GetSubUUID(vrb.Vnode).Range(func(k int64, _ int8) bool {
				if k != srcuuid {
					m[k] = 0
					nodeWare.csVrHandle(k, 0, vrb)
				}
				return true
			})
			if li, ok := nodeWare.bkuuid(vrb.Vnode); ok && len(li) > 0 {
				for _, u := range li {
					if _, ok := m[u]; !ok && u != srcuuid {
						m[u] = 0
						nodeWare.csVrHandle(u, 0, vrb)
					}
				}
			}
		} else {
			nodeWare.csVrHandle(srcuuid, 0, &VBean{Rtype: 8, Vnode: vrb.Vnode, Rnode: vrb.Rnode})
		}
	case 8:
		f := false
		if vr, ok := vgate.VGate.GetVroom(vrb.Vnode); ok {
			if vr.FoundNode != "" {
				nodeWare.csVrHandle(srcuuid, 0, &VBean{Rtype: 1, Vnode: vrb.Vnode, FoundNode: &vr.FoundNode})
				//vr.AuthMap().Range(func(k string, _ int8) bool {
				//	nodeWare.csVrHandle(srcuuid, 0, &VBean{Rtype: 3, Vnode: vrb.Vnode, FoundNode: &vr.FoundNode, Rnode: &k})
				//	return true
				//})
				f = true
			}
		}
		if !f {
			sys.SendNode(*vrb.Rnode, &TimAck{Ok: false, TimType: int8(sys.TIMSTREAM), Error: errs.ERR_NOEXIST.TimError()}, sys.TIMACK)
		}
	case 9:
		sys.SendNode(*vrb.Rnode, &TimAck{Ok: false, TimType: int8(sys.TIMSTREAM), Error: errs.ERR_PERM_DENIED.TimError(), N: &vrb.Vnode}, sys.TIMACK)
	case 10:
		sys.SendNode(*vrb.Rnode, &TimAck{Ok: false, TimType: int8(sys.TIMVROOM), Error: errs.ERR_NOEXIST.TimError(), N: &vrb.Vnode}, sys.TIMACK)
	case 51:
		vgate.VGate.Register(*vrb.FoundNode, vrb.Vnode)
	case 52:
		vgate.VGate.Remove(*vrb.FoundNode, vrb.Vnode)
	case 53:
		//vgate.VGate.AddAuth(vrb.Vnode, *vrb.FoundNode, *vrb.Rnode)
	case 54:
		//vgate.VGate.DelAuth(vrb.Vnode, *vrb.FoundNode, *vrb.Rnode)
	case 57:
		//problem->send repeatedly
		if !reStream.Contains(*vrb.StreamId) {
			reStream.Put(*vrb.StreamId, 0)
			sys.TimSteamProcessor(vrb, sys.TRANS_GOAL)
		} else {
			buf := NewBufferByPool()
			defer buf.Free()
			buf.WriteString(vrb.Vnode)
			buf.WriteString(sys.Conf.Salt)
			buf.Write(Int64ToBytes(srcuuid))
			f := CRC64(buf.Bytes())
			if i, ok := reStreamUUID.Get(f); ok {
				if i > 5 {
					vb := &VBean{Rtype: 6, Vnode: vrb.Vnode}
					nodeWare.csVrHandle(srcuuid, 0, vb)
				} else {
					reStreamUUID.Put(f, atomic.AddInt32(&i, 1))
				}
			} else {
				reStreamUUID.Put(f, 1)
			}
		}
	default:
	}
}

func (nd *nodeware) wsstt() (_r int64) {
	tcc := nd.GetAllTlContext()
	for _, tc := range tcc {
		_r += tc.onNum
	}
	_r += sys.WssLen()
	return
}

var _unaccess = NewMap[int64, int8]()

func (nd *nodeware) unaccess() (_r []int64) {
	_r = make([]int64, 0)
	_unaccess.Range(func(k int64, _ int8) bool {
		_r = append(_r, k)
		return true
	})
	return
}
