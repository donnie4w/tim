// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package keystore

import (
	"fmt"
	"github.com/donnie4w/tim/log"
	"os"
	"strconv"
	"time"

	ks "github.com/donnie4w/gofer/keystore"
	"github.com/donnie4w/gofer/util"
	"github.com/donnie4w/tim/sys"
)

func init() {
	sys.Service(sys.INIT_KEYSTORE, serv(1))
}

type serv byte

func (s serv) Serve() error {
	s.init(sys.KEYSTORE)
	return nil
}

func (serv) Close() error {
	return nil
}

func (s serv) init(dir string) {
	if dir == "" {
		dir, _ = os.Getwd()
	}
	if err := s.initAdmin(dir); err != nil {
		log.FmtPrint("keystore init failed")
		os.Exit(0)
	}
	if sys.OpenSSL.PublicPath != "" || sys.OpenSSL.PrivatePath != "" {
		a := fmt.Sprint(time.Now().UnixNano())
		var err error
		var bs []byte
		var ok bool
		if bs, err = ks.RsaEncrypt([]byte(a), sys.OpenSSL.PublicPath); err == nil {
			if bs, err = ks.RsaDecrypt(bs, sys.OpenSSL.PrivatePath); err == nil {
				ok = a == string(bs)
			}
		}
		if err != nil || !ok {
			panic("publickey and privatekey authFailed")
		}
	}
}

func (serv) initAdmin(dir string) (err error) {
	if ks.KeyStore, err = ks.NewKeyStore(dir, "keystore.tdb"); err == nil {
		Admin.Load()
		if v, ok := Admin.GetOther("TIMUUID"); ok {
			id, _ := strconv.ParseUint(v, 10, 64)
			sys.UUID = int64(id)
		} else {
			sys.UUID = int64(util.UUID32())
			Admin.PutOther("TIMUUID", fmt.Sprint(sys.UUID))
		}
	}
	log.FmtPrint("UUID [", sys.UUID, "]")
	return
}
