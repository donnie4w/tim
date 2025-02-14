// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of t source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package adm

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/donnie4w/tim/log"
	"github.com/donnie4w/tim/stub"
	"os"
	"strings"
	"sync"
	"time"

	goutil "github.com/donnie4w/gofer/util"
	"github.com/donnie4w/gothrift/thrift"
	"github.com/donnie4w/tim/sys"
	"github.com/donnie4w/tim/util"
)

func init() {
	sys.Service(sys.INIT_ADM, new(service))
}

var transportFactory = thrift.NewTBufferedTransportFactory(1 << 13)
var tcompactProtocolFactory = thrift.NewTCompactProtocolFactoryConf(&thrift.TConfiguration{})
var socketTimeout = 10 * time.Second

type service struct {
	isClose         bool
	serverTransport thrift.TServerTransport
}

func (t *service) _server(addr string, processor thrift.TProcessor, TLS bool, serverCrt, serverKey string) (err error) {
	if addr, err = util.ParseAddr(addr); err != nil {
		log.FmtPrint("Server API Service ParseAddr error:", err.Error())
		os.Exit(1)
	}

	if TLS {
		cfg := &tls.Config{}
		var cert tls.Certificate
		if cert, err = tls.LoadX509KeyPair(serverCrt, serverKey); err == nil {
			cfg.Certificates = append(cfg.Certificates, cert)
			t.serverTransport, err = thrift.NewTSSLServerSocketTimeout(addr, cfg, socketTimeout)
		}
	} else {
		t.serverTransport, err = thrift.NewTServerSocketTimeout(addr, socketTimeout)
	}

	if err == nil && t.serverTransport != nil {
		server := thrift.NewTSimpleServer4(processor, t.serverTransport, nil, nil)
		if err = server.Listen(); err == nil {
			s := fmt.Sprint("Server API Service start[", addr, "]")
			if TLS {
				s = fmt.Sprint("Server API Service start tls[", addr, "]")
			}
			log.FmtPrint(s)
			for {
				if _transport, err := server.ServerTransport().Accept(); err == nil {
					go func() {
						defer util.Recover()
						cc := newCliContext(_transport)
						defer cc.close()
						defaultCtx := context.WithValue(context.Background(), "CliContext", cc)
						if inputTransport, err := transportFactory.GetTransport(_transport); err == nil {
							inputProtocol := tcompactProtocolFactory.GetProtocol(inputTransport)
							for {
								ok, err := processor.Process(defaultCtx, inputProtocol, inputProtocol)
								if errors.Is(err, thrift.ErrAbandonRequest) {
									break
								}
								if errors.As(err, new(thrift.TTransportException)) && err != nil {
									break
								}
								if !ok {
									break
								}
							}
						}
					}()
				}
			}
		}
	}
	if !t.isClose && err != nil {
		fmt.Println("Server API Service start failed:", err)
		os.Exit(1)
	}
	return
}

func (t *service) Close() (err error) {
	defer util.Recover()
	if sys.Conf.AdmListen != nil {
		t.isClose = true
		err = t.serverTransport.Close()
	}
	return
}

func (t *service) Serve() (err error) {
	if sys.Conf.AdmListen != nil {
		tls := false
		if sys.Conf.Ssl_crt != "" && sys.Conf.Ssl_crt_key != "" {
			tls = true
		}
		go t._server(strings.TrimSpace(*sys.Conf.AdmListen), stub.NewAdmifaceProcessor(admProcessor), tls, sys.Conf.Ssl_crt, sys.Conf.Ssl_crt_key)
	} else {
		log.FmtPrint("No Server API Service")
	}
	return
}

type pcontext struct {
	Id       int64
	isAuth   bool
	tt       thrift.TTransport
	mux      *sync.Mutex
	_isClose bool
}

func newCliContext(tt thrift.TTransport) (cc *pcontext) {
	cc = &pcontext{goutil.UUID64(), false, tt, &sync.Mutex{}, false}
	return
}

func (t *pcontext) close() {
	defer util.Recover()
	t.mux.Lock()
	defer t.mux.Unlock()
	if !t._isClose {
		t._isClose = true
		t.tt.Close()
	}
}
