// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package sys

import (
	"fmt"
	"github.com/donnie4w/gofer/hashmap"
	"github.com/donnie4w/gofer/util"
	"github.com/donnie4w/tim/stub"
	"time"
)

const VERSION = "2.0.2"

var (
	Service   = hashmap.NewTreeMap[int, Server](5)
	Tim       Server
	STARTTIME = time.Now()
	UUID      int64
	LOGDEBUG  bool
	GOGC      int
	ORIGIN    string
	DEBUGADDR string

	LA       bool
	TIMJSON  string
	KEYSTORE string
	SEP_BIN  = byte(131)
	SEP_STR  = "|"
	Stat     istat
	Conf     *stub.ConfBean
	Bind     string

	DefaultAccount       = [2]string{"admin", "123"}
	MaskSeed             = util.Int64ToBytes(int64(1 << 60))
	ADMADDR              = ""
	WEBADMINADDR         = fmt.Sprint(6 << 10)
	CSADDR               = fmt.Sprint(7 << 10)
	IMADDR               = 5 << 10
	InaccurateTime       = time.Now().UnixNano()
	PINGTO               = int64(500)
	ConnectTimeout       = 10 * time.Second
	WaitTimeout          = 10 * time.Second
	MaxTransLength       = 10 * MB
	DeviceLimit          = 1
	DeviceTypeLimit      = 1
	MaxBackup            = 3
	NodeMaxlength        = 64
	BlockApiMap          *hashmap.Map[TIMTYPE, int8]
	OpenSSL              = &stub.Openssl{}
	defaultPwd           = "tim20171212"
	defaultAesencryptkey = "ie8*&(I984){bW{@a@#￥%H'"
	defaultConnectLimit  = int64(1 << 24)
	defaultMemlimit      = 1 << 11
	defaultsyssalt       = "#@*=+-<>?:|$&()%$#{]aQkLIPM79643028U'TRKF_}"
	defaultLimitRate     = int64(1 << 8)
	defaultConf          = ""
	defaultTTL           = uint64(24 * 60 * 60) // 1 day
)
