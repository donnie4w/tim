// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package sys

import (
	"github.com/donnie4w/tim/log"
)

func blankLine() {
	log.Write([]byte("\n"))
}

func timlogo() {
	_r := `
	=================================================================================
	===========        ===============   ===   ======       ======        ===========
	===========              ===               === ===     === ===        ===========
	===========              ===         ===   ===  ===   ===  ===        ===========
	===========              ===         ===   ===   === ===   ===        ===========
	===========              ===         ===   ===    =====    ===        ===========
	=================================================================================
	`
	log.Info(_r)
}
