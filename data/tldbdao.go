// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim
//

package data

type timmessage struct {
	Id     int64
	ChatId uint64 `idx:"1"`
	Stanza []byte
}

func (this timmessage) tid() uint64 { return this.ChatId }

type timuser struct {
	Id         int64
	UUID       uint64 `idx:"1"`
	Pwd        uint64
	Createtime int64
	UBean      []byte
}

func (this timuser) tid() uint64 { return this.UUID }

type timgroup struct {
	Id         int64
	Gtype      int8
	UUID       uint64 `idx:"1"`
	Createtime int64
	Status     int8
	RBean      []byte
}

func (this timgroup) tid() uint64 { return this.UUID }

type timoffline struct {
	Id     int64
	UUID   uint64 `idx:"1"`
	ChatId uint64
	Stanza []byte
	Mid    int64
}

func (this timoffline) tid() uint64 { return this.UUID }

type timrelate struct {
	Id     int64
	UUID   uint64 `idx:"1"`
	Status uint8
}

func (this timrelate) tid() uint64 { return this.UUID }

type timroster struct {
	Id     int64
	Relate uint64 `idx:"1"`
	UUID   uint64 `idx:"1"`
	TUUID  uint64
}

func (this timroster) tid() uint64 { return this.UUID }

type timmucroster struct {
	Id     int64
	Relate uint64 `idx:"1"`
	UUID   uint64 `idx:"1"`
	TUUID  uint64
}

func (this timmucroster) tid() uint64 { return this.UUID }

type timblock struct {
	Id     int64
	UnikId uint64 `idx:"1"`
	UUID   uint64 `idx:"1"`
	TUUID  uint64
}

func (this timblock) tid() uint64 { return this.UUID }

type timblockroom struct {
	Id     int64
	UnikId uint64 `idx:"1"`
	UUID   uint64 `idx:"1"`
	TUUID  uint64
}

func (this timblockroom) tid() uint64 { return this.UUID }

type timstruct interface {
	tid() uint64
}

type OfflineBean struct {
	Id     int64
	Mid    int64
	Stanze []byte
}
