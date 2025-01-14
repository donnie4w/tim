// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package sys

import (
	"github.com/donnie4w/gothrift/thrift"
	"github.com/donnie4w/tim/errs"
	. "github.com/donnie4w/tim/stub"
	"github.com/donnie4w/tlnet"
)

type TIMTYPE byte
type BUSINESSTYPE int32

const (
	TIMEXP TIMTYPE = 0

	TIMACK             TIMTYPE = 12
	TIMPING            TIMTYPE = 13
	TIMREGISTER        TIMTYPE = 14
	TIMTOKEN           TIMTYPE = 15
	TIMAUTH            TIMTYPE = 16
	TIMOFFLINEMSG      TIMTYPE = 17
	TIMOFFLINEMSGEND   TIMTYPE = 18
	TIMBROADPRESENCE   TIMTYPE = 19
	TIMLOGOUT          TIMTYPE = 20
	TIMPULLMESSAGE     TIMTYPE = 21
	TIMVROOM           TIMTYPE = 22
	TIMBUSINESS        TIMTYPE = 41
	TIMNODES           TIMTYPE = 42
	TIMMESSAGE         TIMTYPE = 90
	TIMPRESENCE        TIMTYPE = 91
	TIMREVOKEMESSAGE   TIMTYPE = 92
	TIMBURNMESSAGE     TIMTYPE = 93
	TIMSTREAM          TIMTYPE = 94
	TIMBIGSTRING       TIMTYPE = 95
	TIMBIGBINARY       TIMTYPE = 96
	TIMBIGBINARYSTREAM TIMTYPE = 97
)

const (
	ADMPING TIMTYPE = 11
	ADMAUTH TIMTYPE = 12
	ADMSUB  TIMTYPE = 13
)

const (
	GROUP_PRIVATE          TIMTYPE = 1
	GROUP_OPEN             TIMTYPE = 2
	GROUP_STATUS_ALIVE     TIMTYPE = 1
	GROUP_STATUS_CANCELLED TIMTYPE = 2
)

const (
	BUSINESS_ROSTER              BUSINESSTYPE = 1
	BUSINESS_USERROOM            BUSINESSTYPE = 2
	BUSINESS_ROOMUSERS           BUSINESSTYPE = 3
	BUSINESS_ADDROSTER           BUSINESSTYPE = 4
	BUSINESS_FRIEND              BUSINESSTYPE = 5
	BUSINESS_REMOVEROSTER        BUSINESSTYPE = 6
	BUSINESS_BLOCKROSTER         BUSINESSTYPE = 7
	BUSINESS_NEWROOM             BUSINESSTYPE = 8
	BUSINESS_ADDROOM             BUSINESSTYPE = 9
	BUSINESS_PASSROOM            BUSINESSTYPE = 10
	BUSINESS_NOPASSROOM          BUSINESSTYPE = 11
	BUSINESS_PULLROOM            BUSINESSTYPE = 12
	BUSINESS_KICKROOM            BUSINESSTYPE = 13
	BUSINESS_BLOCKROOM           BUSINESSTYPE = 14
	BUSINESS_BLOCKROOMMEMBER     BUSINESSTYPE = 15
	BUSINESS_LEAVEROOM           BUSINESSTYPE = 16
	BUSINESS_CANCELROOM          BUSINESSTYPE = 17
	BUSINESS_BLOCKROSTERLIST     BUSINESSTYPE = 18
	BUSINESS_BLOCKROOMLIST       BUSINESSTYPE = 19
	BUSINESS_BLOCKROOMMEMBERLIST BUSINESSTYPE = 20
	BUSINESS_MODIFYAUTH          BUSINESSTYPE = 21
)

const (
	NODEINFO_ROSTER              BUSINESSTYPE = 1
	NODEINFO_ROOM                BUSINESSTYPE = 2
	NODEINFO_ROOMMEMBER          BUSINESSTYPE = 3
	NODEINFO_USERINFO            BUSINESSTYPE = 4
	NODEINFO_ROOMINFO            BUSINESSTYPE = 5
	NODEINFO_MODIFYUSER          BUSINESSTYPE = 6
	NODEINFO_MODIFYROOM          BUSINESSTYPE = 7
	NODEINFO_BLOCKROSTERLIST     BUSINESSTYPE = 8
	NODEINFO_BLOCKROOMLIST       BUSINESSTYPE = 9
	NODEINFO_BLOCKROOMMEMBERLIST BUSINESSTYPE = 10
)

const (
	VROOM_NEW     TIMTYPE = 1
	VROOM_REMOVE  TIMTYPE = 2
	VROOM_ADDAUTH TIMTYPE = 3
	VROOM_DELAUTH TIMTYPE = 4
	VROOM_SUB     TIMTYPE = 5
	VROOM_UNSUB   TIMTYPE = 6
	VROOM_MESSAGE TIMTYPE = 7
)

var (
	AckHandle             func([]byte) (err errs.ERROR)
	PingHandle            func(*tlnet.Websocket) (err errs.ERROR)
	RegisterHandle        func([]byte) (node string, err errs.ERROR)
	TokenHandle           func([]byte) (_r int64, err errs.ERROR)
	AuthHandle            func([]byte, *tlnet.Websocket) (err errs.ERROR)
	OfflinemsgHandle      func(*tlnet.Websocket) (err errs.ERROR)
	BroadpresenceHandle   func([]byte, *tlnet.Websocket) (err errs.ERROR)
	PullMessageHandle     func([]byte, *tlnet.Websocket) (err errs.ERROR)
	VRoomHandle           func([]byte, *tlnet.Websocket) (err errs.ERROR)
	MessageHandle         func([]byte, *tlnet.Websocket) (err errs.ERROR)
	BigStringHandle       func([]byte, *tlnet.Websocket) (err errs.ERROR)
	BigBinaryHandle       func([]byte, *tlnet.Websocket) (err errs.ERROR)
	BigBinaryStreamHandle func([]byte, *tlnet.Websocket) (err errs.ERROR)
	PresenceHandle        func([]byte, *tlnet.Websocket) (err errs.ERROR)
	StreamHandle          func([]byte, *tlnet.Websocket) (err errs.ERROR)
	BusinessHandle        func([]byte, *tlnet.Websocket) (err errs.ERROR)
	NodeInfoHandle        func([]byte, *tlnet.Websocket) (err errs.ERROR)

	OsToken              func(string, *string, *string) (int64, string, errs.ERROR)
	OsRegister           func(string, string, *string) (string, errs.ERROR)
	OsUserBean           func(string, *TimUserBean) errs.ERROR
	OsRoom               func(string, string, *string, int8) (string, errs.ERROR)
	OsRoomBean           func(string, string, *TimRoomBean) errs.ERROR
	OsMessage            func([]string, *TimMessage) (err errs.ERROR)
	OsPresence           func([]string, *TimPresence) (err errs.ERROR)
	OsModify             func(string, *string, string, *string) errs.ERROR
	OsVroomprocess       func(string, int8) string
	OsBlockUser          func(string, int64)
	OsBlockList          func() map[string]int64
	TimMessageProcessor  func(*TimMessage, int8) errs.ERROR
	TimPresenceProcessor func(*TimPresence, int8) errs.ERROR
	TimSteamProcessor    func(*VBean, int8) errs.ERROR
	PxMessage func(int64, *TimMessage) errs.ERROR

	// Deprecated: Use CsMessageService instead.
	//CsMessage func(*TimMessage, int8) bool
	// Deprecated: Use CsPresenceService instead.
	//CsPresence func(*TimPresence, int8) bool
	// Deprecated: Use CsVBeanService instead.
	//CsVBean func(*VBean) bool

	CsNode        func(string) int64
	CsWssInfo     func(string) []byte
	Csuser        func(string, bool, int64) error
	GetALLUUIDS   func() []int64
	Client2Serve  func(string) error
	GetRemoteNode func() []*RemoteNode
	SendNode      func(string, thrift.TStruct, TIMTYPE) bool
	SendWs        func(int64, thrift.TStruct, TIMTYPE) bool
	WssList       func(int64, int64) ([]*Tid, int64)
	WssInfo       func(string) []byte
	WssLen        func() int64
	WssTt         func() int64
	DelWs         func(*tlnet.Websocket)
	WsById        func(int64) (*tlnet.Websocket, bool)
	HasNode       func(string) bool
	HasWs         func(*tlnet.Websocket) bool
	Unaccess      func() []int64
	Interrupt     func(*Tid) errs.ERROR
	Detect        func([]string)

	CsMessageService  func(int64, *TimMessage, bool) bool
	CsPresenceService func(int64, *TimPresence, bool) bool
	CsVBeanService    func(int64, *VBean, bool) bool
)
