// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package amr

import (
	"github.com/donnie4w/tim/log"
	"github.com/donnie4w/tim/sys"
	"os"
)

type amrStore interface {
	put(atype AMRTYPE, key, value []byte, ttl uint64) error
	get(atype AMRTYPE, key []byte) ([]byte, error)
	remove(atype AMRTYPE, key []byte) error
	append(atype AMRTYPE, key, value []byte, ttl uint64) error
	getMutil(atype AMRTYPE, key []byte) [][]byte
	removeKV(atype AMRTYPE, key, value []byte) error
	close() error
}

var amr amrStore

type amrservie byte

func (amrservie) Serve() error {
	switch sys.GetCstype() {
	case sys.CS_RAFTX:
		amr = newRaftxAmr()
	case sys.CS_RAX:
		log.FmtPrint("unrealized rax")
	case sys.CS_REDIS:
		amr = newRedisAmr()
	case sys.CS_ETCD:
		amr = newEtcdAmr()
	case sys.CS_ZOOKEEPER:
		amr = newZkAmr()
	default:
		log.FmtPrint("No Cluster Service")
		amr = localAmr(1)
		islocalamr = true
	}
	if amr == nil {
		log.FmtPrint("amr init failed")
		os.Exit(1)
	}
	return nil
}

func (amrservie) Close() error {
	if amr == nil {
		return nil
	}
	return amr.close()
}

type localAmr byte

func (a localAmr) put(atype AMRTYPE, key, value []byte, ttl uint64) error { return nil }
func (a localAmr) get(atype AMRTYPE, key []byte) ([]byte, error) {
	return nil, nil
}
func (a localAmr) remove(atype AMRTYPE, key []byte) error { return nil }

func (a localAmr) append(atype AMRTYPE, key, value []byte, ttl uint64) error { return nil }

func (a localAmr) getMutil(atype AMRTYPE, key []byte) [][]byte {
	return [][]byte{}
}

func (a localAmr) removeKV(atype AMRTYPE, key, value []byte) error { return nil }

func (a localAmr) close() error {
	return nil
}
