// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim
//

package sys

import (
	"github.com/donnie4w/gothrift/thrift"
	. "github.com/donnie4w/tim/stub"
	"github.com/donnie4w/tlnet"
)

type TIMTYPE byte

const (
	TIMEX TIMTYPE = 0

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
	GROUP_PRIVATE          int8 = 1
	GROUP_OPEN             int8 = 2
	GROUP_STATUS_ALIVE     int8 = 1
	GROUP_STATUS_CANCELLED int8 = 2
)

var (
	BUSINESS_ROSTER              int32 = 1
	BUSINESS_USERROOM            int32 = 2
	BUSINESS_ROOMUSERS           int32 = 3
	BUSINESS_ADDROSTER           int32 = 4
	BUSINESS_FRIEND              int32 = 5
	BUSINESS_REMOVEROSTER        int32 = 6
	BUSINESS_BLOCKROSTER         int32 = 7
	BUSINESS_NEWROOM             int32 = 8
	BUSINESS_ADDROOM             int32 = 9
	BUSINESS_PASSROOM            int32 = 10
	BUSINESS_NOPASSROOM          int32 = 11
	BUSINESS_PULLROOM            int32 = 12
	BUSINESS_KICKROOM            int32 = 13
	BUSINESS_BLOCKROOM           int32 = 14
	BUSINESS_BLOCKROOMMEMBER     int32 = 15
	BUSINESS_LEAVEROOM           int32 = 16
	BUSINESS_CANCELROOM          int32 = 17
	BUSINESS_BLOCKROSTERLIST     int32 = 18
	BUSINESS_BLOCKROOMLIST       int32 = 19
	BUSINESS_BLOCKROOMMEMBERLIST int32 = 20
	BUSINESS_MODIFYAUTH          int32 = 21
)

var (
	NODEINFO_ROSTER              int32 = 1
	NODEINFO_ROOM                int32 = 2
	NODEINFO_ROOMMEMBER          int32 = 3
	NODEINFO_USERINFO            int32 = 4
	NODEINFO_ROOMINFO            int32 = 5
	NODEINFO_MODIFYUSER          int32 = 6
	NODEINFO_MODIFYROOM          int32 = 7
	NODEINFO_BLOCKROSTERLIST     int32 = 8
	NODEINFO_BLOCKROOMLIST       int32 = 9
	NODEINFO_BLOCKROOMMEMBERLIST int32 = 10
)

var (
	AckHandle             func([]byte) (err ERROR)
	PingHandle            func(*tlnet.Websocket) (err ERROR)
	RegisterHandle        func([]byte) (node string, err ERROR)
	TokenHandle           func([]byte) (_r int64, err ERROR)
	AuthHandle            func([]byte, *tlnet.Websocket) (err ERROR)
	OfflinemsgHandle      func(*tlnet.Websocket) (err ERROR)
	BroadpresenceHandle   func([]byte, *tlnet.Websocket) (err ERROR)
	PullMessageHandle     func([]byte, *tlnet.Websocket) (err ERROR)
	VRoomHandle           func([]byte, *tlnet.Websocket) (err ERROR)
	MessageHandle         func([]byte, *tlnet.Websocket) (err ERROR)
	BigStringHandle       func([]byte, *tlnet.Websocket) (err ERROR)
	BigBinaryHandle       func([]byte, *tlnet.Websocket) (err ERROR)
	BigBinaryStreamHandle func([]byte, *tlnet.Websocket) (err ERROR)
	RevokemessageHandle   func([]byte, *tlnet.Websocket) (err ERROR)
	BurnmessageHandle     func([]byte, *tlnet.Websocket) (err ERROR)
	PresenceHandle        func([]byte, *tlnet.Websocket) (err ERROR)
	StreamHandle          func([]byte, *tlnet.Websocket) (err ERROR)
	BusinessHandle        func([]byte, *tlnet.Websocket) (err ERROR)
	NodeInfoHandle        func([]byte, *tlnet.Websocket) (err ERROR)

	DataInit             func() error
	OsToken              func(string, *string, *string) (int64, string, ERROR)
	OsRegister           func(string, string, *string) (string, ERROR)
	OsUserBean           func(string, *TimUserBean) ERROR
	OsRoom               func(string, string, *string, int8) (string, ERROR)
	OsRoomBean           func(string, string, *TimRoomBean) ERROR
	KeyStoreInit         func(string)
	TimMessageProcessor  func(*TimMessage, int8) ERROR
	TimPresenceProcessor func(*TimPresence, int8) ERROR
	TimSteamProcessor    func(*VBean) ERROR
	OsMessage            func(*TimNodes, *TimMessage) (err ERROR)
	OsModify             func(string, string, *string) ERROR
	OsVroomprocess       func(string, int8) string
	CsMessage            func(*TimMessage, int8) bool
	CsPresence           func(*TimPresence, int8) bool
	CsVBean              func(*VBean) bool
	CsNode               func(string) int64
	GetALLUUIDS          func() []int64
	Client2Serve         func(string) error
	BroadRmNode          func() error
	GetRemoteNode        func() []*RemoteNode
	NodeInfo             func(string) ([]byte, ERROR)
	SendNode             func(string, thrift.TStruct, TIMTYPE) bool
	SendWs               func(int64, thrift.TStruct, TIMTYPE) bool
	BlockUser            func(string, int64)
	BlockList            func() map[string]int64
	WssList              func() []*Tid
	WssInfo              func(string) []byte
	CsWssInfo            func(string) []byte
	WssLen               func() int64
	WssTt                func() int64
	DelWs                func(*tlnet.Websocket)
	HasNode              func(string) bool
	HasWs                func(*tlnet.Websocket) bool
	Unaccess             func() []int64
	Interrupt            func(*Tid) ERROR
	Csuser               func(string, bool, int64) error
)
