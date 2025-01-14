// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package sys

type DBMOD byte

const (
	NODB       DBMOD = 1
	TLDB       DBMOD = 2
	INLINEDB   DBMOD = 3
	EXTERNALDB DBMOD = 4
	MONGODB    DBMOD = 5
	CASSANDRA  DBMOD = 6
)

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

type CSTYPE byte

const (
	CS_RAFTX     CSTYPE = 1
	CS_REDIS     CSTYPE = 2
	CS_ETCD      CSTYPE = 3
	CS_ZOOKEEPER CSTYPE = 4
	CS_RAX       CSTYPE = 5
)

const (
	SOURCE_OS   int8 = 1
	SOURCE_USER int8 = 2
	SOURCE_ROOM int8 = 3
)

const (
	_ = iota
	INIT_DATA
	INIT_KEYSTORE
	INIT_TIMNET
	INIT_TC
	INIT_TNET
	INIT_ADM
	INIT_AMR
	INIT_TRANS
)
