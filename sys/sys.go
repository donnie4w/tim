// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package sys

import (
	"github.com/donnie4w/tim/log"
	"os"
)

func init() {
	Service(INIT_SYS, server(0))
}

type server byte

func (s server) Serve() error {
	praseflag()
	log.Info(timlogo)
	return nil
}

func (s server) Close() (err error) {
	service.Descend(func(_ int, s Server) bool {
		s.Close()
		return true
	})
	os.Exit(0)
	return
}

func AddNode(addr string) (err error) {
	//return Client2Serve(addr)
	return nil
}

func UseBuiltInData() bool {
	return Conf.InlineDB != nil || len(Conf.InlineExtent) > 0 || Conf.Tldb != nil || len(Conf.TldbExtent) > 0
}

type Server interface {
	Serve() (err error)
	Close() (err error)
}

type istat interface {
	CReq() int64
	CReqDo()
	CReqDone()

	CPros() int64
	CProsDo()
	CProsDone()

	Tx() int64
	TxDo()
	TxDone()

	Ibs() int64
	Ib(int64)

	Obs() int64
	Ob(int64)
}

func GetDBMOD() DBMOD {
	if Conf.InlineDB != nil || len(Conf.InlineExtent) > 0 {
		return INLINEDB
	}
	if Conf.Tldb != nil || len(Conf.TldbExtent) > 0 {
		return TLDB
	}
	if Conf.ExternalDB != nil {
		return EXTERNALDB
	}
	if Conf.NoDB != nil && *Conf.NoDB {
		return NODB
	}
	return INLINEDB
}

func GetCstype() CSTYPE {
	if Conf.Raftx != nil {
		return CS_RAFTX
	}
	if Conf.Rax != nil {
		return CS_RAX
	}
	if Conf.Redis != nil {
		return CS_REDIS
	}
	if Conf.Etcd != nil {
		return CS_ETCD
	}
	if Conf.ZooKeeper != nil {
		return CS_ZOOKEEPER
	}
	return 0
}
