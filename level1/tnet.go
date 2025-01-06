// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package level1

import (
	"errors"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	. "github.com/donnie4w/gofer/hashmap"
	. "github.com/donnie4w/gofer/lock"
	"github.com/donnie4w/simplelog/logging"
	. "github.com/donnie4w/tim/stub"
	"github.com/donnie4w/tim/sys"
)

var clientLinkCache = NewMapL[string, int8]()
var await = NewAwait[int8](1 << 8)
var awaitCsBean = NewAwait[*CsBean](1 << 8)
var strLock = NewStrlock(1 << 8)
var once = &sync.Once{}
var tnetservice = &tnetService{ok: []byte{5}}
var chapTxTemp = NewLimitHashMap[int64, int8](1 << 15)
var reTx = NewLimitHashMap[int64, int8](1 << 18)
var reStream = NewLimitHashMap[int64, int8](1 << 18)
var reStreamUUID = NewLimitHashMap[uint64, int32](1 << 18)

var nodeCache = NewLimitHashMap[string, []int64](1 << 15)

func init() {
	sys.Client2Serve = client2Serve
	sys.GetALLUUIDS = nodeWare.GetALLUUID
	sys.GetRemoteNode = getRemoteNode
	sys.Csuser = nodeWare.csuser
	sys.WssTt = nodeWare.wsstt
	sys.Unaccess = nodeWare.unaccess
	//sys.CsMessage = nodeWare.csmessage
	//sys.CsPresence = nodeWare.cspresence
	sys.CsWssInfo = nodeWare.cswssinfo
	//sys.CsVBean = nodeWare.csVbean
	sys.CsNode = nodeWare.csnode
}

type tnetService struct {
	ok []byte
}

func (this *tnetService) Serve() (err error) {
	if sys.CSADDR != "" {
		err = this._serve(sys.CSADDR)
	}
	return
}

func (this *tnetService) _serve(addr string) (err error) {
	tnetserver := &tnetServer{ok: this.ok}
	tnetserver.handle(&itnetServ{NewNumLock(64)}, myServer2ClientHandler, mySecvErrorHandler)
	go once.Do(heardbeat)
	return tnetserver.Serve(addr)
}

func (this *tnetService) Connect(addr string, async bool) (err1, err2 error) {
	if nodeWare.LenByAddr(addr) > 0 {
		err1 = errors.New("addr:" + addr + " already existed")
	}
	if addr == sys.CSADDR {
		err1 = errors.New("addr:" + addr + " is local addr")
	}
	if addr = strings.Trim(addr, " "); addr == "" {
		err1 = errors.New("addr is bad")
	}
	if err1 != nil {
		return
	}
	logging.Debug("conn:", addr)
	tnetserver := new(tnetServer)
	tnetserver.handle(&itnetServ{NewNumLock(64)}, myClient2ServerHandler, myCliErrorHandler)
	err2 = tnetserver.Connect(addr, async)
	return
}

func (this *tnetService) connectLoop(addr string, async bool, loopNum int) (err1, err2 error) {
	i := 0
	for i < loopNum {
		i++
		err1, err2 = this.Connect(addr, async)
		if err1 != nil {
			return
		}
		if err1 == nil && err2 == nil {
			return
		}
		<-time.After(time.Second)
	}
	return
}

func (this *tnetService) Close() (err error) {
	if sys.CSADDR != "" {
		tsfclientserver.close()
		nodeWare.close()
	}
	return
}

func (this *tnetService) Ok() bool {
	return this.ok[0] == 1
}

func client2Serve(addr string) (err error) {
	if sys.CSADDR == "" {
		return errors.New("node is stand-alone")
	}
	if nodeWare.LenByAddr(addr) > 0 {
		return errors.New("addr:" + addr + " already existed")
	}
	if addr == sys.CSADDR {
		return errors.New("addr:" + addr + " is local addr")
	}
	if addr = strings.Trim(addr, " "); addr == "" {
		return errors.New("addr is bad")
	}
	err1, err2 := tnetservice.Connect(addr, true)
	if err1 != nil {
		err = err1
	}
	if err2 != nil {
		err = err2
	}
	return
}

var lnetservice = &lnetService{}

type lnetService struct {
	v int8
	c int64
}

func (this *lnetService) Serve() (err error) {
	if sys.Conf.Public != "" {
		_, err = tnetservice.connectLoop(sys.Conf.Public, true, 15)
	}
	return
}

func (this *lnetService) _server(node string) (err error) {
	this.v++
	go this.connect(this.v)
	<-time.After(5 * time.Second)
	sys.LA = false
	if err = tnetservice._serve(node); err != nil && this.c < 10 {
		this.c++
		sys.LA = true
		this.Serve()
	} else if this.c >= 10 {
		this.Close()
	}
	return
}

func (this *lnetService) Connect(addr string, async bool) (err1, err2 error) {
	return
}
func (this *lnetService) Close() (err error) {
	os.Exit(0)
	return
}

func (this *lnetService) connect(v int8) {
	i := int32(0)
	for !tnetservice.Ok() && v == this.v {
		if atomic.AddInt32(&i, 1) > 60 {
			break
		}
		<-time.After(1 * time.Second)
	}
	if tnetservice.Ok() && v == this.v {
		this.Serve()
	}
}

func getRemoteNode() []*RemoteNode {
	return nodeWare.getRemoteNodes()
}

var netservice = &netService{}

type netService struct {
	service sys.Server
}

func (this *netService) Serve() (err error) {
	if sys.LA {
		this.service = lnetservice
	} else {
		this.service = tnetservice
	}
	go this.service.Serve()
	return nil
}

func (this *netService) Close() (err error) {
	return this.service.Close()
}
