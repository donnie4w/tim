// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim
//

package level1

import (
	"context"
	"strings"
	"sync"
	"time"

	. "github.com/donnie4w/gofer/hashmap"
	. "github.com/donnie4w/gofer/util"
	"github.com/donnie4w/gothrift/thrift"
	. "github.com/donnie4w/tim/stub"
	"github.com/donnie4w/tim/sys"
	"github.com/donnie4w/tim/util"
	. "github.com/donnie4w/tsf/server"
)

const tlContextCtx = "tlContext"

type tlContext struct {
	id          int64
	iface       Itnet
	remoteAddr  string
	remoteUuid  int64
	remoteIP    string
	remoteHost2 string
	remoteCsNum int32
	verifycode  int64
	selfPing    int64
	selfPong    uint32

	transport  thrift.TTransport
	mux        *sync.Mutex
	defaultCtx context.Context
	cancleChan chan byte
	mergeChan  chan *syncBean
	mergeCount int64
	mergemux   *sync.Mutex

	isServer   bool
	isClose    bool
	pingNum    int64
	pongNum    int64
	isAuth     bool
	onNum      int64
	_do_reconn bool
}

type syncBean struct {
	SyncId int64
	Result int8
}

func (this *tlContext) SetId(_id int64) {
	this.id = _id
}

func (this *tlContext) Close() {
	defer util.Recover()
	this.mux.Lock()
	defer this.mux.Unlock()
	if !this.isClose {
		this.isClose = true
		nodeWare.del(this)
		this.transport.Close()
		close(this.cancleChan)
	}
}

func (this *tlContext) CloseAndEnd() (err error) {
	this._do_reconn = true
	this.Close()
	return
}

func newTlContext2(socket *TSocket) (tc *tlContext) {
	tc = &tlContext{id: RandId(), mux: new(sync.Mutex), mergeChan: make(chan *syncBean, 1<<17), mergemux: &sync.Mutex{}}
	tc.transport = socket
	tc.iface = &ItnetImpl{socket}
	go func() {
		availMap.Put(tc, time.Now().Unix())
		if availmux.TryLock() {
			availtk()
		}
	}()
	return
}

func remoteHost(transport thrift.TTransport) (_r string) {
	defer recover()
	if addr := transport.(*thrift.TSocket).Conn().RemoteAddr(); addr != nil {
		if ss := strings.Split(addr.String(), ":"); len(ss) == 2 {
			_r = ss[0]
		}
	}
	return
}

func remoteHost2(tsocket *TSocket) (string, string) {
	defer util.Recover()
	if addr := tsocket.Conn().RemoteAddr(); addr != nil {
		if ss := strings.Split(addr.String(), ":"); len(ss) == 2 {
			return ss[0], ss[1]
		}
	}
	return "", ""
}

var availMap = NewMapL[*tlContext, int64]()
var availmux = &sync.Mutex{}

func availtk() {
	defer availmux.Unlock()
	tk := time.NewTicker(10 * time.Second)
	for {
		if availMap.Len() == 0 {
			break
		}
		select {
		case <-tk.C:
			availMap.Range(func(k *tlContext, v int64) bool {
				if v+30 < time.Now().Unix() {
					k.CloseAndEnd()
					availMap.Del(k)
				}
				return true
			})
		}
	}
}

type chapBean struct {
	Stat   int8
	Code   int64
	Key    string
	TcId   int64
	UUID   int64
	Time   int64
	TxId   int64
	IDcard int64
}

func newchapBean() (a *chapBean) {
	a = &chapBean{Stat: 1, Code: 0, Key: sys.Conf.Pwd, TcId: 0, UUID: sys.UUID, Time: time.Now().UnixNano(), TxId: RandId()}
	return
}
