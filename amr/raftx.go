// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package amr

import (
	"github.com/donnie4w/gofer/util"
	"github.com/donnie4w/raftx"
	"github.com/donnie4w/raftx/raft"
	"github.com/donnie4w/tim/log"
	"github.com/donnie4w/tim/sys"
)

type raftxAmr struct {
	raft raftx.Raftx
}

func newRaftxAmr() *raftxAmr {
	if sys.Conf.Raftx == nil || len(sys.Conf.Raftx.Peers) <= 1 {
		panic("no raftx peers")
	}
	r := raftx.NewRaftx(&raft.Config{ListenAddr: sys.Conf.Raftx.ListenAddr, PeerAddr: sys.Conf.Raftx.Peers})
	go func() {
		if err := r.Open(); err != nil {
			log.Error("open raftx err:", err)
			panic(err)
		}
	}()
	return &raftxAmr{r}
}

func (ra *raftxAmr) Put(key, value []byte, ttl uint64) {
	ra.raft.MemCommand(key, value, ttl, raft.MEM_PUT)
}

func (ra *raftxAmr) Get(key []byte) []byte {
	v, _ := ra.raft.GetMemValue(key)
	return v
}

func (ra *raftxAmr) Remove(key []byte) {
	ra.raft.MemCommand(key, nil, 0, raft.MEM_DEL)
}

func (ra *raftxAmr) AddAccount(node string, uuid int64) {
	idbs := util.Int64ToBytes(int64(util.FNVHash64([]byte(node))))
	ra.raft.MemCommand(idbs, util.Int64ToBytes(uuid), sys.Conf.TTL, raft.MEM_APPEND)
}

func (ra *raftxAmr) RemoveAccount(node string, uuid int64) {
	idbs := util.Int64ToBytes(int64(util.FNVHash64([]byte(node))))
	ra.raft.MemCommand(idbs, util.Int64ToBytes(uuid), sys.Conf.TTL, raft.MEM_DELKV)
}

func (ra *raftxAmr) GetAccount(node string) (r []int64) {
	idbs := util.Int64ToBytes(int64(util.FNVHash64([]byte(node))))
	if vs := ra.raft.GetLocalMemMultiValue(idbs); len(vs) > 0 {
		if len(vs) > 1 {
			r = make([]int64, 0, len(vs))
			m := make(map[int64]bool)
			for _, v := range vs {
				k := util.BytesToInt64(v)
				if _, b := m[k]; !b {
					r = append(r, k)
					m[k] = true
				}
			}
		} else {
			r = []int64{util.BytesToInt64(vs[0])}
		}
	}
	return
}
