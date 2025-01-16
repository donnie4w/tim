// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package dao

import "github.com/donnie4w/gofer/util"

func (u *Timblock) Tid() uint64 {
	return uint64(u.GetUuid())
}

func (u *Timuser) Tid() uint64 {
	return uint64(u.GetUuid())
}

func (u *Timblockroom) Tid() uint64 {
	return uint64(u.GetUuid())
}

func (u *Timroster) Tid() uint64 {
	return uint64(u.GetUuid())
}

func (u *Timmessage) Tid() uint64 {
	return util.FNVHash64(u.GetChatid())
}

func (u *Timrelate) Tid() uint64 {
	return util.FNVHash64(u.GetUuid())
}

func (u *Timmucroster) Tid() uint64 {
	return uint64(u.GetUuid())
}

func (u *Timoffline) Tid() uint64 {
	return uint64(u.GetUuid())
}

func (u *Timdomain) Tid() uint64 {
	return uint64(u.GetId())
}
