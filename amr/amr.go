// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package amr

import (
	"github.com/donnie4w/gofer/util"
	"github.com/donnie4w/tim/sys"
	"sync"
)

func init() {
	sys.Service.Put(sys.INIT_AMR, (amrservie)(0))
}

type AMR interface {
	Put(key, value []byte, ttl uint64)
	Get(key []byte) []byte
	Remove(key []byte)
	AddAccount(node string, uuid int64)
	RemoveAccount(node string, uuid int64)
	GetAccount(node string) []int64
}

var amr AMR

type amrservie byte

func (amrservie) Serve() error {
	switch sys.GetCstype() {
	case sys.CS_RAFTX:
		amr = newRaftxAmr()
	case sys.CS_RAX:
		panic("unrealized rax")
	case sys.CS_REDIS:
		panic("unrealized redis")
	case sys.CS_ETCD:
		panic("unrealized etcd")
	case sys.CS_ZOOKEEPER:
		panic("unrealized zookeeper")
	default:
		amr = newAMR()
	}
	if amr == nil {
		panic("amr init failed")
	}
	return nil
}

func (amrservie) Close() error {
	return nil
}

func Put(key, value []byte, ttl uint64) {
	amr.Put(key, value, ttl)
}

func Get(key []byte) []byte {
	return amr.Get(key)
}

func Remove(key []byte) {
	amr.Remove(key)
}

func AddAccount(node string, uuid int64) {
	amr.AddAccount(node, uuid)
}

func RemoveAccount(node string, uuid int64) {
	amr.RemoveAccount(node, uuid)
}

func GetAccount(node string) []int64 {
	return amr.GetAccount(node)
}

type simpleAmr struct {
	mux sync.RWMutex
	am  map[uint64]map[int64]byte
	bm  map[uint64][]byte
}

func newAMR() AMR {
	return &simpleAmr{am: make(map[uint64]map[int64]byte)}
}

func (a *simpleAmr) Put(key, value []byte, ttl uint64) {
	a.mux.Lock()
	defer a.mux.Unlock()
	id := util.FNVHash64(key)
	a.bm[id] = value
}

func (a *simpleAmr) Get(key []byte) []byte {
	a.mux.RLock()
	defer a.mux.RUnlock()
	id := util.FNVHash64(key)
	v, _ := a.bm[id]
	return v
}

func (a *simpleAmr) Remove(key []byte) {
	a.mux.Lock()
	defer a.mux.Unlock()
	id := util.FNVHash64(key)
	delete(a.bm, id)
}

func (a *simpleAmr) AddAccount(node string, uuid int64) {
	a.mux.Lock()
	defer a.mux.Unlock()
	id := util.FNVHash64([]byte(node))
	mm, ok := a.am[id]
	if !ok {
		mm = make(map[int64]byte)
	}
	mm[uuid] = 0
	a.am[id] = mm
}

func (a *simpleAmr) RemoveAccount(node string, uuid int64) {
	a.mux.Lock()
	defer a.mux.Unlock()
	id := util.FNVHash64([]byte(node))
	if mm, ok := a.am[id]; ok {
		delete(mm, uuid)
		if len(mm) == 0 {
			delete(a.am, id)
		}
	}
}

func (a *simpleAmr) GetAccount(node string) []int64 {
	a.mux.RLock()
	defer a.mux.RUnlock()
	id := util.FNVHash64([]byte(node))
	if mm, ok := a.am[id]; ok {
		r := make([]int64, 0, len(mm))
		for k := range mm {
			r = append(r, k)
		}
		return r
	}
	return nil
}
