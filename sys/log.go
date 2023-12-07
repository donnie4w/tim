// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim
//

package sys

import (
	"fmt"

	"github.com/donnie4w/simplelog/logging"
)

var log = logging.NewLogger().SetFormat(logging.FORMAT_DATE | logging.FORMAT_TIME).SetLevel(logging.LEVEL_INFO)

func FmtLog(v ...any) {
	info := fmt.Sprint(v...)
	a, b := "", ""
	ll := 80
	if ll >= len(info) {
		for i := 0; i < (ll-len(info))/2; i++ {
			a = a + "="
		}
		b = a
		if ll > len(info)+len(a)*2 {
			b = a + "="
		}
	}
	log.Info(a, info, b)
}

func blankLine() {
	log.Write([]byte("\n"))
}

func timlogo(){
	_r :=`
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