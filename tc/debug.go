// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package tc

import (
	"github.com/donnie4w/tim/log"
	"github.com/donnie4w/tim/sys"
	"github.com/donnie4w/tim/util"
	"net/http"
	_ "net/http/pprof"
	"runtime"
)

func tlDebug() {
	defer util.Recover()
	if sys.Conf.PprofAddr != "" {
		runtime.SetMutexProfileFraction(1)
		runtime.SetBlockProfileRate(1)
		var err error
		if sys.Conf.PprofAddr, err = util.ParseAddr(sys.Conf.PprofAddr); err == nil {
			log.FmtPrint("Http pprof Service start[", sys.Conf.PprofAddr, "]")
			if err := http.ListenAndServe(sys.Conf.PprofAddr, nil); err != nil {
				log.FmtPrint("Http pprof Service start failed:" + err.Error())
			}
		}
	}
}
