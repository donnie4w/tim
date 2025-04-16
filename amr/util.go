// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package amr

import (
	"github.com/donnie4w/tim/sys"
)

type AMRTYPE byte

const (
	ACCOUNT AMRTYPE = 1
	TOKEN   AMRTYPE = 2
	VNODE   AMRTYPE = 3
	BLOCK   AMRTYPE = 4
	UUID    AMRTYPE = 5
)

var islocalamr = false

func init() {
	sys.Service(sys.INIT_AMR, (amrservie)(0))
}

func amrKey(atype AMRTYPE, key []byte) []byte {
	bs := make([]byte, len(key)+1)
	bs[0] = byte(atype)
	copy(bs[1:], key)
	return bs
}
