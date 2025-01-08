// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package trans

import (
	"github.com/donnie4w/tim/stub"
)

type csNet interface {
	addNoAck() int32
	Id() int64
	IsValid() bool
	Close() (err error)
	TimMessage(syncId int64, tm *stub.TimMessage) (err error)
	TimPresence(syncId int64, tp *stub.TimPresence) (err error)
	TimStream(syncId int64, vb *stub.VBean) (err error)
	TimCsVBean(syncId int64, vb *stub.CsVrBean) (err error)
	TimAck(syncId int64) (err error)
}

const (
	TIMMESSAGE  byte = 1
	TIMPRESENCE byte = 2
	TIMSTREAM   byte = 3
	TIMCSVBEAN  byte = 4
	TIMACK      byte = 5
)
