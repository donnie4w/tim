// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package tc

import (
	"fmt"
	"github.com/donnie4w/tim/log"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"strings"

	goutil "github.com/donnie4w/gofer/util"
	"github.com/donnie4w/tim/sys"
	"github.com/donnie4w/tim/util"
)

func tlDebug() {
	defer util.Recover()
	if sys.DEBUGADDR != "" {
		runtime.SetMutexProfileFraction(1)
		runtime.SetBlockProfileRate(1)
		if !strings.Contains(sys.DEBUGADDR, ":") && goutil.MatchString("^[0-9]{4,5}$", sys.DEBUGADDR) {
			sys.DEBUGADDR = fmt.Sprint(":", sys.DEBUGADDR)
		}
		log.FmtPrint("Debug start[", sys.DEBUGADDR, "]")
		if err := http.ListenAndServe(sys.DEBUGADDR, nil); err != nil {
			log.FmtPrint("debug  start failed:" + err.Error())
		}
	}
}
