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

	. "github.com/donnie4w/gofer/keystore"
	. "github.com/donnie4w/gofer/util"
	"github.com/donnie4w/tim/sys"
)

func init() {
	sys.Service.Put(sys.INIT_KEYSTORE, (serv(1)))
}

type serv byte

func (serv) Serve() error {
	Init(sys.KEYSTORE)
	return nil
}

func (serv) Close() error {
	return nil
}

func Init(dir string) {
	if dir == "" {
		dir, _ = os.Getwd()
	}
	if err := InitAdmin(dir); err != nil {
		log.FmtPrint("keystore init failed")
		os.Exit(0)
	}
	if sys.OpenSSL.PublicPath != "" || sys.OpenSSL.PrivatePath != "" {
		a := fmt.Sprint(time.Now().UnixNano())
		var err error
		var bs []byte
		var ok bool
		if bs, err = RsaEncrypt([]byte(a), sys.OpenSSL.PublicPath); err == nil {
			if bs, err = RsaDecrypt(bs, sys.OpenSSL.PrivatePath); err == nil {
				ok = a == string(bs)
			}
		}
		if err != nil || !ok {
			panic("publickey and privatekey authFailed")
		}
	}
}

func InitAdmin(dir string) (err error) {
	if KeyStore, err = NewKeyStore(dir, "keystore.tdb"); err == nil {
		Admin.Load()
		if v, ok := Admin.GetOther("TIMUUID"); ok {
			id, _ := strconv.ParseUint(v, 10, 64)
			sys.UUID = int64(id)
		} else {
			sys.UUID = int64(uuid())
			Admin.PutOther("TIMUUID", fmt.Sprint(sys.UUID))
		}
	}
	return
}

func uuid() uint32 {
	return UUID32()
}
