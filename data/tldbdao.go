// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package data

type timmessage struct {
	Id         int64
	ChatId     uint64 `idx:"1"`
	Stanza     []byte
	Timeseries int64 `idx:"1"`
}

func (this timmessage) Tid() uint64 { return this.ChatId }

type timuser struct {
	Id         int64
	UUID       uint64 `idx:"1"`
	Pwd        uint64
	Createtime int64
	UBean      []byte
	Timeseries int64 `idx:"1"`
}

func (this timuser) Tid() uint64 { return this.UUID }

type timgroup struct {
	Id         int64
	Gtype      int8
	UUID       uint64 `idx:"1"`
	Createtime int64
	Status     int8
	RBean      []byte
	Timeseries int64 `idx:"1"`
}

func (this timgroup) Tid() uint64 { return this.UUID }

type timoffline struct {
	Id         int64
	UUID       uint64 `idx:"1"`
	ChatId     uint64
	Stanza     []byte
	Mid        int64
	Timeseries int64 `idx:"1"`
}

func (this timoffline) Tid() uint64 { return this.UUID }

type timrelate struct {
	Id         int64
	UUID       uint64 `idx:"1"`
	Status     uint8
	Timeseries int64 `idx:"1"`
}

func (this timrelate) Tid() uint64 { return this.UUID }

type timroster struct {
	Id         int64
	Unikid     uint64 `idx:"1"`
	UUID       uint64 `idx:"1"`
	TUUID      uint64
	Timeseries int64 `idx:"1"`
}

func (this timroster) Tid() uint64 { return this.UUID }

type timmucroster struct {
	Id         int64
	Unikid     uint64 `idx:"1"`
	UUID       uint64 `idx:"1"`
	TUUID      uint64
	Timeseries int64 `idx:"1"`
}

func (this timmucroster) Tid() uint64 { return this.UUID }

type timblock struct {
	Id         int64
	UnikId     uint64 `idx:"1"`
	UUID       uint64 `idx:"1"`
	TUUID      uint64
	Timeseries int64 `idx:"1"`
}

func (this timblock) Tid() uint64 { return this.UUID }

type timblockroom struct {
	Id         int64
	UnikId     uint64 `idx:"1"`
	UUID       uint64 `idx:"1"`
	TUUID      uint64
	Timeseries int64 `idx:"1"`
}

func (this timblockroom) Tid() uint64 { return this.UUID }

type timstruct interface {
	Tid() uint64
}

type OfflineBean struct {
	Id         int64
	Mid        int64
	Stanze     []byte
	Timeseries int64 `idx:"1"`
}
