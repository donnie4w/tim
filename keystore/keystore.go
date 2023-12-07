// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim
//

package keystore

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"time"

	. "github.com/donnie4w/gofer/keystore"
	. "github.com/donnie4w/gofer/util"
	"github.com/donnie4w/tim/sys"
)

func init() {
	sys.KeyStoreInit = Init
}

func Init(dir string) {
	if dir == "" {
		dir, _ = os.Getwd()
	}
	if err := InitAdmin(dir); err != nil {
		sys.FmtLog("keystore init failed")
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
	buf := bytes.NewBuffer([]byte{})
	buf.Write(Int64ToBytes(int64(os.Getpid())))
	for i := 0; i < 100; i++ {
		buf.Write(Int64ToBytes(RandId()))
	}
	if _r, err := RandStrict(1 << 31); err == nil && _r > 0 {
		buf.Write(Int64ToBytes(_r))
	}
	return CRC32(buf.Bytes())
}