// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim
//

package data

import (
	. "github.com/donnie4w/tim/stub"
	"github.com/donnie4w/tim/sys"
)

type nilprocess struct{}

func (this *nilprocess) Register(username, pwd string, domain *string) (node string, e sys.ERROR) {
	e = sys.ERR_DATABASE
	return
}

func (this *nilprocess) Login(username, pwd string, domain *string) (_r string, e sys.ERROR) {
	_r = username
	return
}

func (this *nilprocess) Modify(uint64, *string, string, *string) (e sys.ERROR) {
	return
}

func (this *nilprocess) Token(username, pwd string, domain *string) (_r string, err sys.ERROR) {
	_r = username
	return
}

func (this *nilprocess) SaveMessage(tm *TimMessage) (err error) {
	return
}

func (this *nilprocess) GetMessage(fromTid *Tid, rtype int8, to string, mid, limit int64) (tmList []*TimMessage, err error) {
	return
}

func (this *nilprocess) GetMessageByMid(tid uint64, mid int64) (tm *TimMessage, err error) {
	return
}

func (this *nilprocess) DelMessageByMid(tid uint64, mid int64) (err error) {
	return
}

func (this *nilprocess) SaveOfflineMessage(tm *TimMessage) (err error) {
	return
}

func (this *nilprocess) GetOfflineMessage(tid *Tid, limit int) (oblist []*OfflineBean, err error) {
	return
}

func (this *nilprocess) DelOfflineMessage(tid uint64, ids ...int64) (_r int64, err error) {
	return
}

func (this *nilprocess) UserGroup(node string, domain *string) (_r []string) {
	return
}

func (this *nilprocess) GroupRoster(groupnode string) (_r []string) {
	return
}

func (this *nilprocess) AuthGroupAndUser(groupnode, usernode string, domain *string) (ok bool, err error) {
	ok = true
	return
}
