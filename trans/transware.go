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
	"github.com/donnie4w/tim/log"
	"github.com/donnie4w/tim/stub"
	"github.com/donnie4w/tim/sys"
	"math"
	"os"
	"time"
)

func init() {
	sys.Service(sys.INIT_TRANS, tw)
	sys.CsMessageService = tw.CsMessageService
	sys.CsPresenceService = tw.CsPresenceService
	sys.CsVBeanService = tw.CsVBeanService
	sys.CsDevice = tw.CsDevice
	sys.GetALLUUIDS = tw.getAlluuids
	sys.Unaccess = tw.unAccess
}

var tw = NewTransWare()

type TransWare struct {
	am       *hashmap.Map[int64, string]
	conn     *hashmap.Map[int64, csNet] //sockId->csNet
	csmap    *hashmap.Map[int64, csNet] //uuid->csNet
	unaccess *hashmap.Map[int64, struct{}]
	dataWait *lock.FastAwait[any]
	tserver  *tServer
	numlock  *lock.Numlock
}

func NewTransWare() *TransWare {
	tw := new(TransWare)
	tw.am = hashmap.NewMap[int64, string]()
	tw.conn = hashmap.NewMap[int64, csNet]()
	tw.csmap = hashmap.NewMap[int64, csNet]()
	tw.unaccess = hashmap.NewMap[int64, struct{}]()
	tw.numlock = lock.NewNumLock(1 << 7)
	tw.dataWait = lock.NewFastAwait[any]()
	go tw.twTicker()
	return tw
}

func (tw *TransWare) Serve() error {
	if sys.GetCstype() != 0 {
		if sys.Conf.CsListen != "" && sys.Conf.CsAccess != "" {
			tw.tserver = newTServer(tw)
			go tw.tserver.serve(sys.Conf.CsListen)
			tw.registerUuid()
		} else {
			log.FmtPrint("Cluster configuration is incomplete[cluster_listen:", sys.Conf.CsListen, "][cluster_access:", sys.Conf.CsAccess, "]")
			os.Exit(1)
		}
	}
	return nil
}

func (tw *TransWare) registerUuid() {
	defer util.Recover(nil)
	for {
		if r, err := amr.GetCsAccess(sys.UUID); err == nil && r == "" {
			if amr.PutCsAccess() == nil {
				if list, _ := sys.WssList(0, math.MaxInt64); len(list) > 0 {
					for i, aa := range list {
						if amr.AddAccount(aa.Node) == nil {
							goto END
						}
						time.Sleep(time.Millisecond)
						if i%10000 == 0 {
							amr.PutCsAccess()
						}
					}
				}
			}
		}
		if amr.PutCsAccess() == nil {
			break
		} else {
			time.Sleep(5 * time.Second)
			log.Warn("register uuid fail")
		}
	}
END:
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
		switch sys.TIMTYPE(vb.GetRtype()) {
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
			if addr, _ := amr.GetCsAccess(uuid); addr != "" {
				conn := newConnect(tw)
				if conn.open(addr) == nil {
					tw.csmap.Put(uuid, conn)
					tw.unaccess.Del(uuid)
					return conn
				}
			}
		}
	}
	if cn == nil {
		tw.unaccess.Put(uuid, struct{}{})
	}
	return cn
}

func (tw *TransWare) CsDevice(node string) []byte {
	if uuids := amr.GetAccount(node); len(uuids) > 0 {
		for _, uuid := range uuids {
			if uuid != sys.UUID {
				if cn := tw.getAndSetCsnet(uuid); cn != nil {
					syncId := util.UUID64()
					if cn.TimCsDevice(syncId, &stub.CsDevice{Node: &node}) == nil {
						if v, err := tw.dataWait.Wait(syncId, sys.WaitTimeout); err == nil && v != nil {
							return v.(*stub.CsDevice).GetTypeList()
						} else if err != nil {
							cn.addNoAck()
						}
					} else {
						return []byte{}
					}
				}
			}
		}
	}
	return []byte{}
}

func (tw *TransWare) twTicker() {
	tk := time.NewTicker(time.Duration(sys.UUIDCSTIME*2/3) * time.Second)
	for {
		select {
		case <-tk.C:
			tw.registerUuid()
		}
	}
}

func (tw *TransWare) getAlluuids() (r []int64) {
	r = make([]int64, 0)
	tw.conn.Range(func(k int64, _ csNet) bool {
		r = append(r, k)
		return true
	})
	return
}

func (tw *TransWare) unAccess() (r []int64) {
	r = make([]int64, 0)
	tw.unaccess.Range(func(k int64, _ struct{}) bool {
		r = append(r, k)
		return true
	})
	return
}
