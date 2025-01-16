// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package mesh

import (
	"context"
	. "github.com/donnie4w/gofer/util"
	"github.com/donnie4w/tim/log"
	. "github.com/donnie4w/tim/stub"
	"github.com/donnie4w/tim/util"
	"sync/atomic"
	"time"
)

type tnetServer struct {
	servermux *serverMux
	ok        []byte
}

func (this *tnetServer) handle(processor Itnet, handler func(tc *tlContext), cliError func(tc *tlContext)) {
	if this.servermux == nil {
		this.servermux = &serverMux{}
		this.servermux.Handle(processor, handler, cliError)
	}
}

func (this *tnetServer) Serve(_addr string) (err error) {
	if _addr, err = util.ParseAddr(_addr); err != nil {
		return
	}
	err = tsfclientserver.server(_addr, this.servermux.processor, this.servermux.handler, this.servermux.cliError, this.ok)
	return
}

func (this *tnetServer) Connect(_addr string, async bool) (err error) {
	if _addr, err = util.ParseAddr(_addr); err != nil {
		return
	}
	return tsfserverclient.server(_addr, this.servermux.processor, this.servermux.handler, this.servermux.cliError, async)
}

type serverMux struct {
	processor Itnet
	handler   func(tc *tlContext)
	cliError  func(tc *tlContext)
}

func (this *serverMux) Handle(processor Itnet, handler func(tc *tlContext), cliError func(tc *tlContext)) {
	this.processor = processor
	this.handler = handler
	this.cliError = cliError
}

func myServer2ClientHandler(tc *tlContext) {
}

func myClient2ServerHandler(tc *tlContext) {
	defer util.Recover()
	if !tc.isClose {
		ab := newchapBean()
		ab.IDcard = tc.id
		if bs, err := encodeChapBean(ab); err == nil {
			if err := tc.csnet.Chap(context.Background(), bs); err != nil {
				log.Error(err)
			}
		}
	}
}

func mySecvErrorHandler(tc *tlContext) {
	go reconn(tc)
}

func myCliErrorHandler(tc *tlContext) {
	clientLinkCache.Del(tc.remoteAddr)
	go reconn(tc)
}

func reconn(tc *tlContext) {
	if tc == nil || tc.remoteAddr == "" || tc.remoteUuid == 0 {
		return
	}
	defer util.Recover()
	log.Info(">>>[", tc.remoteUuid, "][", tc.remoteAddr, "]")
	nodeWare.del(tc)
	if !tc.isServer && !tc._do_reconn {
		tc._do_reconn = true
		i := 0
		for !nodeWare.hasUUID(tc.remoteUuid) {
			if clientLinkCache.Has(tc.remoteAddr) {
				<-time.After(time.Duration(RandUint(6)) * time.Second)
			} else {
				if err1, err2 := tnetservice.Connect(tc.remoteAddr, false); err1 != nil {
					break
				} else if err2 != nil {
					if i < 100 {
						i++
					} else if i > 1<<13 {
						break
					}
					<-time.After(time.Duration(RandUint(uint(6+i))) * time.Second)
				}
			}
		}
	}
}

func heardbeat() {
	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-ticker.C:
			<-time.After(time.Duration(RandUint(5)) * time.Second)
			_heardbeat(nodeWare.GetAllTlContext())
		}
	}
}

func _heardbeat(tcs []*tlContext) {
	defer util.Recover()
	for _, tc := range tcs {
		func(tc *tlContext) {
			defer util.Recover()
			tc.csnet.Ping(context.TODO(), piBs(tc))
			if atomic.AddInt64(&tc.pingNum, 1) > 8 {
				log.Error("ping failed:[", tc.remoteUuid, "][", tc.remoteAddr, "] ping number:", tc.pingNum)
				go reconn(tc)
			}
		}(tc)
	}
}
