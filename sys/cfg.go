// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package sys

import (
	"flag"
	"fmt"
	"github.com/donnie4w/gofer/hashmap"
	"github.com/donnie4w/gofer/util"
	"github.com/donnie4w/tim/log"
	"github.com/donnie4w/tim/stub"
	"os"
	"runtime"
	"runtime/debug"
	"sync/atomic"
	"time"
)

func init() {
	Tim = &server{}
}

func praseflag() {
	flag.StringVar(&DEBUGADDR, "debug", "", "debug address")
	flag.StringVar(&TIMJSON, "c", "tim.json", "configuration file of tim in json")
	flag.StringVar(&ORIGIN, "origin", "", "origin for websocket")
	flag.StringVar(&KEYSTORE, "ks", "", "dir of keystore")
	flag.BoolVar(&LOGDEBUG, "log", false, "debug log on or off")
	flag.IntVar(&GOGC, "gc", -1, "a collection is triggered when the ratio of freshly allocated data")
	flag.Usage = usage
	flag.Parse()
	flag.Usage()
	parsec()

	if Conf.Pwd == "" {
		Conf.Pwd = defaultPwd
	}
	if Conf.EncryptKey == "" {
		Conf.EncryptKey = defaultAesencryptkey
	}
	if Conf.Seed > 0 {
		MaskSeed = util.Int64ToBytes(Conf.Seed)
	}
	if Conf.Salt == "" {
		Conf.Salt = defaultsyssalt
	}
	if Conf.PingTo > 0 {
		PINGTO = Conf.PingTo
	}
	if Conf.MaxBackup != nil {
		MaxBackup = *Conf.MaxBackup
	}
	if Conf.NodeMaxlength != nil {
		NodeMaxlength = *Conf.NodeMaxlength
	}
	if Conf.LimitRate <= 0 {
		Conf.LimitRate = defaultLimitRate
	}
	if Conf.AdminListen != "" {
		WEBADMINADDR = Conf.AdminListen
	}
	if Conf.Public != "" {
		LA, CSADDR = true, ""
	} else if Conf.CsListen != "" {
		CSADDR = Conf.CsListen
	}

	if Conf.AdmListen != nil {
		ADMADDR = *Conf.AdmListen
	}

	if Conf.MaxMessageSize > 0 {
		MaxTransLength = Conf.MaxMessageSize * MB
	}

	if Conf.Memlimit <= 0 {
		Conf.Memlimit = defaultMemlimit
	}

	if Conf.DeviceLimit > 0 {
		DeviceLimit = Conf.DeviceLimit
	}

	if Conf.DevicetypeLimit > 0 {
		DeviceTypeLimit = Conf.DevicetypeLimit
	}

	if Conf.ConnectLimit <= 0 {
		Conf.ConnectLimit = defaultConnectLimit
	}

	if Conf.Listen > 0 {
		IMADDR = Conf.Listen
	}

	if Conf.Bind != nil {
		Bind = *Conf.Bind
	}
	if Conf.Keystore != nil {
		KEYSTORE = *Conf.Keystore
	}
	if Conf.TTL == 0 {
		Conf.TTL = defaultTTL
	}

	if Conf.Security != nil && Conf.BlockAPI != nil {
		BlockApiMap = hashmap.NewMap[TIMTYPE, int8]()
		for _, v := range Conf.BlockAPI {
			BlockApiMap.Put(TIMTYPE(v), 0)
		}
	}

	debug.SetMemoryLimit(int64(Conf.Memlimit) * MB)
	debug.SetGCPercent(GOGC)
	Stat = &stat{}
	//if LOGDEBUG {
	//	logging.SetFormat(logging.FORMAT_DATE|logging.FORMAT_TIME|logging.FORMAT_SHORTFILENAME).SetRollingFile("", "tim.log", 100, logging.MB)
	//} else {
	//	logging.SetLevel(logging.LEVEL_OFF)
	//}
}

func usage() {
	exename := "tim"
	if runtime.GOOS == "windows" {
		exename = "tim.exe"
	}
	fmt.Fprintln(os.Stderr, `tim version: tim/`+VERSION+`
Usage: `+exename+`	
	-c: configuration file  e.g:  -c tim.json
`)
}

func parsec() {
	var err error
	if defaultConf != "" {
		Conf, _ = util.JsonDecode[*stub.ConfBean]([]byte(defaultConf))
	} else if bs, e := util.ReadFile(TIMJSON); e == nil {
		Conf, err = util.JsonDecode[*stub.ConfBean](bs)
		if err != nil {
			log.FmtPrint("configuration file[", TIMJSON, "] parsing error:", err)
		}
	}
	if Conf == nil {
		log.FmtPrint("empty config")
		Conf = &stub.ConfBean{}
	}
}

type stat struct {
	creq  int64
	cpros int64
	tx    int64
	ibs   int64
	obs   int64
}

func (this *stat) CReq() int64 {
	return this.creq
}
func (this *stat) CReqDo() {
	atomic.AddInt64(&this.creq, 1)
}
func (this *stat) CReqDone() {
	atomic.AddInt64(&this.creq, -1)
}

func (this *stat) CPros() int64 {
	return this.cpros
}
func (this *stat) CProsDo() {
	atomic.AddInt64(&this.cpros, 1)
}
func (this *stat) CProsDone() {
	atomic.AddInt64(&this.cpros, -1)
}

func (this *stat) Tx() int64 {
	return this.tx
}
func (this *stat) TxDo() {
	atomic.AddInt64(&this.tx, 1)
}
func (this *stat) TxDone() {
	atomic.AddInt64(&this.tx, -1)
}

func (this *stat) Ibs() int64 {
	return this.ibs
}
func (this *stat) Ib(i int64) {
	atomic.AddInt64(&this.ibs, i)
}
func (this *stat) Obs() int64 {
	return this.obs
}
func (this *stat) Ob(i int64) {
	atomic.AddInt64(&this.obs, i)
}

func tkSqlProperty() {
	tk := time.NewTicker(time.Minute)
	for {
		select {
		case <-tk.C:
			func() {
				defer func() { recover() }()
				now := time.Now()
				t0 := now.Unix()
				t1 := time.Date(now.Year(), now.Month(), now.Day(), now.Hour()+1, 0, 0, 0, now.Location()).Unix()
				<-time.After(time.Duration(t1-t0) * (time.Second))
				if bs, err := util.ReadFile(TIMJSON); err == nil {
					Conf, _ = util.JsonDecode[*stub.ConfBean](bs)
				}
			}()
		}
	}
}
