// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim
package sys

import (
	"os"
	"time"
)

type server struct{}

func (this *server) Serve() error {
	praseflag()
	blankLine()
	timlogo()
	DataInit()
	KeyStoreInit(KEYSTORE)
	Service.BackForEach(func(_ int, s Server) bool {
		defer func() { recover() }()
		go s.Serve()
		<-time.After(time.Millisecond << 9)
		return true
	})
	select {}
}

func (this *server) Close() (err error) {
	Service.FrontForEach(func(_ int, s Server) bool {
		s.Close()
		return true
	})
	os.Exit(0)
	return
}

func AddNode(addr string) (err error) {
	return Client2Serve(addr)
}

func UseDefaultDB() bool {
	return Conf.Tldb != nil || (Conf.TldbExtent != nil && len(Conf.TldbExtent) > 0)
}

func UseTldbExtent() bool {
	return Conf.TldbExtent != nil && len(Conf.TldbExtent) > 0
}
