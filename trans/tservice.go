// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package trans

import (
	"context"
	"fmt"
	"github.com/donnie4w/tim/log"
	"github.com/donnie4w/tim/sys"
	"github.com/donnie4w/tim/util"
	"github.com/donnie4w/tsf"
	"os"
	"time"
)

type tServer struct {
	isClose bool
	server  *tsf.Tsf
	tw      *TransWare
}

func newTServer(transWare *TransWare) *tServer {
	return &tServer{tw: transWare}
}

func (t *tServer) serve(listenAddr string) (err error) {
	if listenAddr, err = util.ParseAddr(listenAddr); err != nil {
		log.FmtPrint("Cluster tim Service ParseAddr error:", err.Error())
		os.Exit(1)
	}
	defer util.Recover2(&err)
	cfg := &tsf.TsfConfig{ListenAddr: listenAddr, TConfiguration: &tsf.TConfiguration{ProcessMerge: true}}
	tc := &tsf.TContext{}
	tc.OnClose = func(ts tsf.TsfSocket) error {
		t.tw.romove(ts.ID())
		return nil
	}
	tc.OnOpenSync = func(socket tsf.TsfSocket) error {
		cn := newTrans(socket)
		socket.SetContext(context.WithValue(context.Background(), true, newProcessor(t.tw, cn, true)))
		t.tw.Add(socket.ID(), cn)
		return nil
	}
	tc.Handler = func(socket tsf.TsfSocket, packet *tsf.Packet) error {
		return router(packet.ToBytes(), socket.GetContext().Value(true).(csNet))
	}
	if t.server, err = tsf.NewTsf(cfg, tc); err == nil {
		if err = t.server.Listen(); err == nil {
			log.FmtPrint("Cluster tim service start [", listenAddr, "] ")
			err = t.server.AcceptLoop()
		}
	}
	if !t.isClose && err != nil {
		log.Error("Cluster tim service failed:", err)
		os.Exit(0)
	}
	return
}

func (t *tServer) close() error {
	t.isClose = true
	return t.server.Close()
}

type connect struct {
	csNet
	tw *TransWare
}

func newConnect(tw *TransWare) *connect {
	return &connect{tw: tw}
}

func (c *connect) open(addr string) (err error) {
	defer util.Recover2(&err)
	tx := &tsf.TContext{}
	wait := make(chan struct{})
	tx.OnOpenSync = func(socket tsf.TsfSocket) error {
		log.Debug("connect open:", sys.Conf.CsListen, "->", addr)
		defer close(wait)
		c.csNet = newTrans(socket)
		socket.SetContext(context.WithValue(context.Background(), true, newProcessor(c.tw, c.csNet, false)))
		c.tw.Add(socket.ID(), c.csNet)
		return nil
	}
	tx.OnClose = func(ts tsf.TsfSocket) error {
		log.Info("OnClose:", ts.ID())
		defer util.Recover()
		defer c.Close()
		c.tw.romove(ts.ID())
		return nil
	}
	tx.Handler = func(socket tsf.TsfSocket, packet *tsf.Packet) error {
		return router(packet.ToBytes(), socket.GetContext().Value(true).(csNet))
	}
	conn := tsf.NewTsfSocketConf(addr, &tsf.TConfiguration{ProcessMerge: true, ConnectTimeout: sys.ConnectTimeout})
	if err = conn.Open(); err == nil {
		go conn.On(tx)
	} else {
		return
	}
	to := time.After(time.Second)
	select {
	case <-wait:
	case <-to:
		conn.Close()
		log.Errorf("connect timeout: %s -> %s ", sys.Conf.CsListen, addr)
		err = fmt.Errorf(fmt.Sprint("connect timeout: %s -> %s ", sys.Conf.CsListen, addr))
	}
	return
}
