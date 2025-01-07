// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package vgate

import (
	"github.com/donnie4w/tim/amr"
	"sync"
	"time"

	"github.com/donnie4w/gofer/hashmap"
	"github.com/donnie4w/gofer/lock"
	"github.com/donnie4w/tim/sys"
	"github.com/donnie4w/tim/util"
)

var VGate = &vgate{vmap: hashmap.NewMap[string, *VRoom](), mux: &sync.Mutex{}, umap: hashmap.NewMap[int64, *hashmap.Map[string, int8]](), strLock: lock.NewStrlock(1 << 7), numLock: lock.NewNumLock(1 << 7)}

func init() {
	go VGate.clearTk()
}

type vgate struct {
	vmap    *hashmap.Map[string, *VRoom]
	umap    *hashmap.Map[int64, *hashmap.Map[string, int8]]
	mux     *sync.Mutex
	strLock *lock.Strlock
	numLock *lock.Numlock
}

func (vg *vgate) NewVRoom(fnode string) string {
	vnode := util.UUIDToNode(util.NewTimUUID())
	vg.Register(fnode, vnode)
	return vnode
}

func (vg *vgate) umapSub(wsId int64, vnode string) {
	if m, ok := vg.umap.Get(wsId); ok {
		m.Put(vnode, 0)
	} else {
		vg.numLock.Lock(wsId)
		defer vg.numLock.Unlock(wsId)
		if !vg.umap.Has(wsId) {
			vg.umap.Put(wsId, hashmap.NewMap[string, int8]())
		}
		go vg.umapSub(wsId, vnode)
	}
}

func (vg *vgate) Register(fnode, vnode string) (_r bool) {
	vg.mux.Lock()
	defer vg.mux.Unlock()
	if !util.CheckNode(vnode) {
		return
	}
	if !vg.vmap.Has(vnode) {
		v := &VRoom{vnode: vnode, subuuid: hashmap.NewMap[int64, int8](), subnode: hashmap.NewMapL[int64, int8](), FoundNode: fnode, lastupdatetime: time.Now().UnixNano()}
		//if fnode != "" {
		//	v.pushAuth.Put(fnode, 0)
		//}
		vg.vmap.Put(vnode, v)
		_r = true
		vg.putAmr(vnode)
	}
	return
}

func (vg *vgate) rmVnode(vnode string) {
	if vg.vmap.Del(vnode) {
		vg.delAmr(vnode)
	}
}

func (vg *vgate) putAmr(vnode string) {
	amr.PutVnode(vnode, sys.UUID)
}

func (vg *vgate) delAmr(vnode string) {
	amr.DelVnode(vnode)
}

func (vg *vgate) Sub(vnode string, uuid int64, wsId int64) bool {
	if amr.GetVnode(vnode) == 0 {
		return false
	}
	vg._sub(vnode, uuid, wsId, false)
	return true
}

func (vg *vgate) SubBinary(vnode string, uuid int64, wsId int64) bool {
	if amr.GetVnode(vnode) == 0 {
		return false
	}
	vg._sub(vnode, uuid, wsId, true)
	return true
}

func (vg *vgate) _sub(vnode string, srcUuid int64, wsId int64, isBinary bool) {
	if vr, ok := vg.vmap.Get(vnode); ok {
		vr.sub(srcUuid, wsId, isBinary)
		if wsId > 0 {
			vg.umapSub(wsId, vnode)
		}
	} else {
		vg.strLock.Lock(vnode)
		defer vg.strLock.Unlock(vnode)
		if !vg.vmap.Has(vnode) {
			v := &VRoom{vnode: vnode, subuuid: hashmap.NewMap[int64, int8](), subnode: hashmap.NewMapL[int64, int8](), lastupdatetime: time.Now().UnixNano()}
			vg.vmap.Put(vnode, v)
			go vg._sub(vnode, srcUuid, wsId, isBinary)
		}
	}
}

func (vg *vgate) UnSub(vnode string, wsId int64) (r int64, ok bool) {
	if amr.GetVnode(vnode) == 0 {
		return 0, false
	}
	if vnode != "" {
		if vr, ok := vg.vmap.Get(vnode); ok {
			vr.delnode(wsId)
			r = vr.subnode.Len()
		}
		vg.umap.Del(wsId)
		ok = true
	}
	return
}

func (vg *vgate) UnSubWithWsId(wsId int64) {
	if vm, ok := vg.umap.Get(wsId); ok {
		vm.Range(func(k string, _ int8) bool {
			vg.UnSub(k, wsId)
			return true
		})
		vg.umap.Del(wsId)
	}
}

func (vg *vgate) UnSubWithUUID(vnode string, uuid int64) {
	if vr, ok := vg.vmap.Get(vnode); ok {
		vr.deluuid(uuid)
	}
}

func (vg *vgate) GetSubUUID(vnode string) (m *hashmap.Map[int64, int8]) {
	if vr, ok := vg.vmap.Get(vnode); ok {
		return vr.subuuid
	}
	return
}

func (vg *vgate) GetNodes(vnode string) (m *hashmap.MapL[int64, int8]) {
	if vr, ok := vg.vmap.Get(vnode); ok {
		return vr.subnode
	}
	return
}

func (vg *vgate) Remove(fnode, vnode string) bool {
	if vr, ok := vg.vmap.Get(vnode); ok {
		if vr.FoundNode == fnode {
			vg.rmVnode(vnode)
			return true
		}
	}
	return false
}

func (vg *vgate) GetVroom(vnode string) (*VRoom, bool) {
	return vg.vmap.Get(vnode)
}

//func (vg *vgate) AddAuth(vnode, fnode, tnode string) bool {
//	if vr, ok := vg.vmap.Get(vnode); ok {
//		return vr.addauth(fnode, tnode)
//	}
//	return false
//}

//func (vg *vgate) DelAuth(vnode, fnode, tnode string) {
//	if vr, ok := vg.vmap.Get(vnode); ok {
//		vr.delauth(fnode, tnode)
//	}
//}

func (vg *vgate) Nodes() *hashmap.Map[string, *VRoom] {
	return vg.vmap
}

func (vg *vgate) clearTk() {
	tk := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-tk.C:
			func() {
				defer util.Recover()
				vg.vmap.Range(func(k string, v *VRoom) bool {
					if v.Expires() {
						vg.rmVnode(k)
					}
					return true
				})
			}()
		}
	}
}

type VRoom struct {
	vnode   string
	subuuid *hashmap.Map[int64, int8]
	subnode *hashmap.MapL[int64, int8]
	//pushAuth       *hashmap.Map[string, int8]
	FoundNode      string
	lastupdatetime int64
}

func (vr *VRoom) sub(uuid int64, wsId int64, isBinary bool) {
	if uuid != sys.UUID {
		vr.subuuid.Put(uuid, 0)
	} else {
		if isBinary {
			vr.subnode.Put(wsId, 1)
		} else {
			vr.subnode.Put(wsId, 0)
		}
	}
}

func (vr *VRoom) deluuid(uuid int64) {
	vr.subuuid.Del(uuid)
}

func (vr *VRoom) delnode(wsId int64) {
	vr.subnode.Del(wsId)
}

//func (vr *VRoom) addauth(fnode, node string) bool {
//	if fnode == vr.FoundNode {
//		vr.pushAuth.Put(node, 0)
//		return true
//	}
//	return false
//}

//func (vr *VRoom) delauth(fnode, node string) {
//	if fnode == vr.FoundNode {
//		vr.pushAuth.Del(node)
//	}
//}

//func (vr *VRoom) AuthMap() *hashmap.Map[string, int8] {
//	return vr.pushAuth
//}

func (vr *VRoom) Expires() bool {
	return vr.lastupdatetime+int64(time.Minute) < time.Now().UnixNano()
}

func (vr *VRoom) Updatetime() {
	vr.lastupdatetime = time.Now().UnixNano()
}

func (vr *VRoom) SubMap() *hashmap.MapL[int64, int8] {
	return vr.subnode
}

func (vr *VRoom) AuthStream(node string) bool {
	//if vr.pushAuth.Has(node) || vr.FoundNode == node {
	//	return true
	//}
	return vr.FoundNode == node
}
