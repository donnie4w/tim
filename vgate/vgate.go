// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim
//

package vgate

import (
	"sync"
	"time"

	. "github.com/donnie4w/gofer/hashmap"
	. "github.com/donnie4w/gofer/lock"
	"github.com/donnie4w/tim/sys"
	"github.com/donnie4w/tim/util"
)

var VGate = &vgate{vmap: NewMap[string, *VRoom](), mux: &sync.Mutex{}, umap: NewMap[int64, *Map[string, int8]](), strLock: NewStrlock(1 << 7), numLock: NewNumLock(1 << 7)}

func init() {
	go VGate.clearTk()
}

type vgate struct {
	vmap    *Map[string, *VRoom]
	umap    *Map[int64, *Map[string, int8]]
	mux     *sync.Mutex
	strLock *Strlock
	numLock *Numlock
}

func (this *vgate) NewVRoom(fnode string) string {
	// v := &VRoom{vnode: util.UUIDToNode(util.NewTimUUID()), subuuid: NewMap[int64, int8](), subnode: NewMapL[string, int8](), pushAuth: NewMap[string, int8](), FoundNode: fnode}
	// if fnode != "" {
	// 	v.pushAuth.Put(fnode, 0)
	// }
	// this.vmap.Put(v.vnode, v)
	// // go this.umapSub(fnode, v.vnode)
	vnode := util.UUIDToNode(util.NewTimUUID())
	this.Register(fnode, vnode)
	return vnode
}

func (this *vgate) umapSub(wsId int64, vnode string) {
	if m, ok := this.umap.Get(wsId); ok {
		m.Put(vnode, 0)
	} else {
		this.numLock.Lock(wsId)
		defer this.numLock.Unlock(wsId)
		if !this.umap.Has(wsId) {
			this.umap.Put(wsId, NewMap[string, int8]())
		}
		go this.umapSub(wsId, vnode)
	}
}

func (this *vgate) Register(fnode, vnode string) (_r bool) {
	this.mux.Lock()
	defer this.mux.Unlock()
	if !util.CheckNode(vnode) {
		return
	}
	if !this.vmap.Has(vnode) {
		v := &VRoom{vnode: vnode, subuuid: NewMap[int64, int8](), subnode: NewMapL[int64, int8](), pushAuth: NewMap[string, int8](), FoundNode: fnode, lastupdatetime: time.Now().UnixNano()}
		if fnode != "" {
			v.pushAuth.Put(fnode, 0)
		}
		this.vmap.Put(vnode, v)
		_r = true
	}
	return
}

func (this *vgate) Sub(vnode string, uuid int64, wsId int64) {
	if vr, ok := this.vmap.Get(vnode); ok {
		vr.sub(uuid, wsId)
		if wsId > 0 {
			this.umapSub(wsId, vnode)
		}
	} else {
		this.strLock.Lock(vnode)
		defer this.strLock.Unlock(vnode)
		if !this.vmap.Has(vnode) {
			v := &VRoom{vnode: vnode, subuuid: NewMap[int64, int8](), subnode: NewMapL[int64, int8](), pushAuth: NewMap[string, int8](), lastupdatetime: time.Now().UnixNano()}
			this.vmap.Put(vnode, v)
			go this.Sub(vnode, uuid, wsId)
		}
	}
}

func (this *vgate) DelNode(vnode string, wsId int64) (_r int64) {
	if vnode != "" {
		if vr, ok := this.vmap.Get(vnode); ok {
			vr.delnode(wsId)
			_r = vr.subnode.Len()
		}
		this.umap.Del(wsId)
	} else {
		if vm, ok := this.umap.Get(wsId); ok {
			vm.Range(func(k string, _ int8) bool {
				this.DelNode(k, wsId)
				return true
			})
			this.umap.Del(wsId)
		}
	}
	return
}

func (this *vgate) DelUuid(vnode string, uuid int64) {
	if vr, ok := this.vmap.Get(vnode); ok {
		vr.deluuid(uuid)
	}
}

func (this *vgate) GetUUID(vnode string) (m *Map[int64, int8]) {
	if vr, ok := this.vmap.Get(vnode); ok {
		return vr.subuuid
	}
	return
}

func (this *vgate) GetNodes(vnode string) (m *MapL[int64, int8]) {
	if vr, ok := this.vmap.Get(vnode); ok {
		return vr.subnode
	}
	return
}

func (this *vgate) Remove(fnode, vnode string) bool {
	if vr, ok := this.vmap.Get(vnode); ok {
		if vr.FoundNode == fnode {
			this.vmap.Del(vnode)
			return true
		}
	}
	return false
}

func (this *vgate) GetVroom(vnode string) (*VRoom, bool) {
	return this.vmap.Get(vnode)
}

func (this *vgate) AddAuth(vnode, fnode, tnode string) bool {
	if vr, ok := this.vmap.Get(vnode); ok {
		return vr.addauth(fnode, tnode)
	}
	return false
}

func (this *vgate) DelAuth(vnode, fnode, tnode string) {
	if vr, ok := this.vmap.Get(vnode); ok {
		vr.delauth(fnode, tnode)
	}
}

func (this *vgate) Nodes() *Map[string, *VRoom] {
	return this.vmap
}

func (this *vgate) clearTk() {
	tk := time.NewTicker(10 * time.Minute)
	for {
		select {
		case <-tk.C:
			func() {
				defer util.Recover()
				this.vmap.Range(func(k string, v *VRoom) bool {
					if v.Expires() {
						this.vmap.Del(k)
					}
					return true
				})
			}()
		}
	}
}

type VRoom struct {
	vnode          string
	subuuid        *Map[int64, int8]
	subnode        *MapL[int64, int8]
	pushAuth       *Map[string, int8]
	FoundNode      string
	lastupdatetime int64
}

func (this *VRoom) sub(uuid int64, wsId int64) {
	if uuid != sys.UUID {
		this.subuuid.Put(uuid, 0)
	} else {
		this.subnode.Put(wsId, 0)
	}
}

func (this *VRoom) deluuid(uuid int64) {
	this.subuuid.Del(uuid)
}

func (this *VRoom) delnode(wsId int64) {
	this.subnode.Del(wsId)
}

func (this *VRoom) addauth(fnode, node string) bool {
	if fnode == this.FoundNode {
		this.pushAuth.Put(node, 0)
		return true
	}
	return false
}

func (this *VRoom) delauth(fnode, node string) {
	if fnode == this.FoundNode {
		this.pushAuth.Del(node)
	}
}

func (this *VRoom) AuthMap() *Map[string, int8] {
	return this.pushAuth
}

func (this *VRoom) Expires() (_r bool) {
	if (this.lastupdatetime + int64(10*time.Minute)) < time.Now().UnixNano() {
		_r = true
	}
	return
}

func (this *VRoom) Updatetime() {
	this.lastupdatetime = time.Now().UnixNano()
}

func (this *VRoom) Nodes() *MapL[int64, int8] {
	return this.subnode
}

func (this *VRoom) Auth(node string) bool {
	if this.pushAuth.Has(node) || this.FoundNode == node {
		return true
	}
	return false
}
