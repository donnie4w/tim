// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package trans

import (
	"github.com/donnie4w/gofer/hashmap"
	"github.com/donnie4w/gofer/lock"
	"github.com/donnie4w/gofer/util"
	"github.com/donnie4w/tim/amr"
	"github.com/donnie4w/tim/stub"
	"github.com/donnie4w/tim/sys"
	"time"
)

func init() {
	sys.Service.Put(sys.INIT_TRANS, tw)
	sys.CsMessageService = tw.CsMessageService
	sys.CsPresenceService = tw.CsPresenceService
	sys.CsVBeanService = tw.CsVBeanService
}

var tw = NewTransWare()

type TransWare struct {
	am       *hashmap.Map[int64, string]
	conn     *hashmap.Map[int64, csNet] //sockId->csNet
	csmap    *hashmap.Map[int64, csNet] //uuid->csNet
	dataWait *lock.FastAwait[any]
	tserver  *tServer
	numlock  *lock.Numlock
}

func NewTransWare() *TransWare {
	tw := new(TransWare)
	tw.am = hashmap.NewMap[int64, string]()
	tw.conn = hashmap.NewMap[int64, csNet]()
	tw.csmap = hashmap.NewMap[int64, csNet]()
	tw.numlock = lock.NewNumLock(1 << 7)
	tw.dataWait = lock.NewFastAwait[any]()
	return tw
}

func (tw *TransWare) Serve() error {
	tw.tserver = newTServer(tw)
	go tw.tserver.serve(sys.Conf.CsListen)
	tw.register()
	return nil
}

func (tw *TransWare) register() {
	amr.Put(util.Int64ToBytes(sys.UUID), []byte(sys.Conf.CsAddr), 60)
}

func (tw *TransWare) Close() error {
	return tw.tserver.close()
}

func (tw *TransWare) Add(id int64, cs csNet) {
	tw.conn.Put(id, cs)
}

func (tw *TransWare) romove(id int64) {
	tw.conn.Del(id)
	tw.csmap.Range(func(k int64, v csNet) bool {
		if v.Id() == id {
			tw.csmap.Del(k)
			return false
		}
		return true
	})
}

func (tw *TransWare) CsMessageService(uuid int64, tm *stub.TimMessage, sync bool) bool {
	if cn := tw.getAndSetCsnet(uuid); cn != nil {
		syncId := int64(0)
		if sync {
			syncId = util.UUID64()
		}
		if cn.TimMessage(syncId, tm) == nil {
			if sync {
				if _, err := tw.dataWait.Wait(syncId, sys.WaitTimeout); err == nil {
					return true
				} else {
					cn.addNoAck()
				}
			} else {
				return true
			}
		}
	}
	return false
}

func (tw *TransWare) CsPresenceService(uuid int64, tp *stub.TimPresence, sync bool) bool {
	if cn := tw.getAndSetCsnet(uuid); cn != nil {
		syncId := int64(0)
		if sync {
			syncId = util.UUID64()
		}
		if cn.TimPresence(syncId, tp) == nil {
			if sync {
				if _, err := tw.dataWait.Wait(syncId, sys.WaitTimeout); err == nil {
					return true
				} else {
					cn.addNoAck()
				}
			} else {
				return true
			}
		}
	}
	return false
}

func (tw *TransWare) CsVBeanService(uuid int64, vb *stub.VBean, sync bool) bool {
	if cn := tw.getAndSetCsnet(uuid); cn != nil {
		syncId := int64(0)
		if sync {
			syncId = util.UUID64()
		}
		switch sys.TIMTYPE(vb.GetDtype()) {
		case sys.VROOM_MESSAGE:
			if cn.TimStream(syncId, vb) == nil {
				if sync {
					if _, err := tw.dataWait.Wait(syncId, sys.WaitTimeout); err == nil {
						return true
					} else {
						cn.addNoAck()
					}
				} else {
					return true
				}
			}
		default:
			cvb := stub.NewCsVrBean()
			cvb.Srcuuid = &sys.UUID
			cvb.Vbean = vb
			if cn.TimCsVBean(syncId, cvb) == nil {
				if sync {
					if _, err := tw.dataWait.Wait(syncId, sys.WaitTimeout); err == nil {
						return true
					} else {
						cn.addNoAck()
					}
				} else {
					return true
				}
			}
		}

	}
	return false
}

func (tw *TransWare) getAndSetCsnet(uuid int64) csNet {
	cn, b := tw.csmap.Get(uuid)
	if !b {
		tw.numlock.Lock(uuid)
		defer tw.numlock.Unlock(uuid)
		if cn, b = tw.csmap.Get(uuid); !b {
			if addr := string(amr.Get(util.Int64ToBytes(uuid))); addr != "" {
				conn := newConnect(tw)
				if conn.open(addr) == nil {
					tw.csmap.Put(uuid, conn)
					return conn
				}
			}
		}
	}
	return cn
}

func (tw *TransWare) twTicker() {
	tk := time.NewTicker(30 * time.Second)
	for {
		select {
		case <-tk.C:
			tw.register()
		}
	}
}
