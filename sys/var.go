// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package sys

import (
	"github.com/donnie4w/gofer/hashmap"
	"github.com/donnie4w/gofer/util"
	"github.com/donnie4w/tim/stub"
	"time"
)

const VERSION = "2.1.0"

var (
	STARTTIME              = time.Now()
	MaskSeed               = util.Int64ToBytes(int64(1 << 60))
	UUIDCSTIME             = 180 //180 second
	ConnectTimeout         = 10 * time.Second
	WaitTimeout            = 10 * time.Second
	MaxTransLength         = 10 * MB
	DeviceLimit            = 1
	DeviceTypeLimit        = 1
	MaxBackup              = 3
	NodeMaxSize            = 64
	OpenSSL                = &stub.Openssl{}
	UUID                   int64
	GOGC                   int
	ORIGIN                 string
	TIMJSON                string
	SEP_BIN                = byte(131)
	SEP_STR                = "|"
	Stat                   istat
	Conf                   *stub.ConfBean
	service                = hashmap.NewTreeMap[int, Server](5)
	defaultAdminAccount    = &stub.AdminAccount{Username: "admin", Password: "123"}
	defaultPingTimeot      = int64(600) // 600 second
	defaultPwd             = "tim20171212"
	defaultAesencryptkey   = "ie8*&(I984){bW{@a@#￥%H'"
	defaultConnectLimit    = int64(1 << 24)
	defaultMemlimit        = 1 << 10
	defaultSalt            = "#@*=+-<>?:|$&()%$#{]aQkLIPM79643028U'TRKF_}"
	defaultLimitRate       = int64(1 << 8)
	defaultTTL             = uint64(24 * 60 * 60) // 1 day
	defaultTokenTimeout    = 10 * time.Second.Nanoseconds()
	defaultCacheAuthExpire = int64(300) //300 Second
)

var (
	LA   bool
	Bind string
)
