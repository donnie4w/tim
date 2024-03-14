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
	"time"

	. "github.com/donnie4w/gofer/hashmap"
	. "github.com/donnie4w/gofer/util"
)

const VERSION = "2.0.0"

const (
	MB = 1 << 20
	GB = 1 << 30
)

const (
	TRANS_SOURCE     int8 = 1
	TRANS_CONSISHASH int8 = 2
	TRANS_STAFF      int8 = 3
	TRANS_GOAL       int8 = 4
)

const (
	ORDER_INOF      int8 = 1
	ORDER_REVOKE    int8 = 2
	ORDER_BURN      int8 = 3
	ORDER_BUSINESS  int8 = 4
	ORDER_STREAM    int8 = 5
	ORDER_BIGSTRING int8 = 6
	ORDER_BIGBINARY int8 = 7
	ORDER_RESERVED  int8 = 30
)

const (
	CB_MESSAGE  int8 = 1
	CB_PRESENCE int8 = 2
)

const (
	SOURCE_OS   int8 = 1
	SOURCE_USER int8 = 2
	SOURCE_ROOM int8 = 3
)

var (
	Service   = NewSortMap[int, Server]()
	Tim       Server
	STARTTIME = time.Now()
	UUID      int64
	LOGDEBUG  bool
	GOGC      int
	ORIGIN    string
	DEBUGADDR string

	LA                   bool
	TIMJSON              string
	KEYSTORE             string
	SEP_BIN              = byte(131)
	SEP_STR              = "|"
	Stat                 istat
	Conf                 *ConfBean
	Bind                 string
	DBtype               byte
	DefaultAccount       = [2]string{"admin", "123"}
	MaskSeed             = Int64ToBytes(int64(1 << 60))
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
	ForBitIfaceMap       *Map[TIMTYPE, int8]
	OpenSSL              = &openssl{}
	defaultPwd           = "tim20171212"
	defaultAesencryptkey = "ie8*&(I984){bW{@a@#￥%H'"
	defaultConnectLimit  = int64(1 << 24)
	defaultMemlimit      = 1 << 11
	defaultsyssalt       = "#@*=+-<>?:|$&()%$#{]aQkLIPM79643028U'TRKF_}"
	defaultTasks         = 1 << 8
	defaultConf          = ""
)
