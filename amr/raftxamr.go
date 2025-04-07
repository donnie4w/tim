// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package amr

import (
	"github.com/donnie4w/raftx"
	"github.com/donnie4w/raftx/raft"
	"github.com/donnie4w/tim/log"
	"github.com/donnie4w/tim/sys"
	"time"
)

type raftxAmr struct {
	raft     raftx.Raftx
	ctime    int64
	getlocal bool
}

func newRaftxAmr() amrStore {
	if sys.Conf.Raftx == nil || len(sys.Conf.Raftx.Peers) <= 1 {
		log.Error("Cluster raftx Service start failed:", "no raftx peers")
		return nil
	}
	rx := raftx.NewRaftx(&raft.Config{ListenAddr: sys.Conf.Raftx.ListenAddr, PeerAddr: sys.Conf.Raftx.Peers})
	if err := rx.Open(); err != nil {
		log.Error("Raftx Service start failed:", err)
		return nil
	} else {
		log.FmtPrint("Raftx Service start: [", sys.Conf.Raftx.ListenAddr, "]")
	}
	log.FmtPrint("Raftx Wait init...")
	rx.WaitRun()
	return &raftxAmr{raft: rx, ctime: time.Now().UnixNano()}
}

func (ra *raftxAmr) put(atype AMRTYPE, key, value []byte, ttl uint64) error {
	return ra.raft.MemCommand(amrKey(atype, key), value, ttl, raft.MEM_PUT)
}

func (ra *raftxAmr) get(atype AMRTYPE, key []byte) ([]byte, error) {
	return ra.raft.GetMemValue(amrKey(atype, key))
}

func (ra *raftxAmr) remove(atype AMRTYPE, key []byte) error {
	return ra.raft.MemCommand(amrKey(atype, key), nil, 0, raft.MEM_DEL)
}

func (ra *raftxAmr) append(atype AMRTYPE, key, value []byte, ttl uint64) error {
	return ra.raft.MemCommand(amrKey(atype, key), value, ttl, raft.MEM_APPEND)
}

func (ra *raftxAmr) getMutil(atype AMRTYPE, key []byte) [][]byte {
	key = amrKey(atype, key)
	if ra.getlocal {
		return ra.raft.GetLocalMemMultiValue(key)
	} else {
		if ra.ctime+int64(sys.Conf.TTL*1e9) > time.Now().UnixNano() {
			v, _ := ra.raft.GetMemMultiValue(key)
			return v
		} else {
			ra.getlocal = true
			return ra.raft.GetLocalMemMultiValue(key)
		}
	}
}

func (ra *raftxAmr) removeKV(atype AMRTYPE, key, value []byte) error {
	return ra.raft.MemCommand(amrKey(atype, key), value, 0, raft.MEM_DELKV)
}

func (ra *raftxAmr) close() error {
	return ra.raft.Stop()
}
