// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim
package data

import "github.com/donnie4w/tim/sys"

func (this *sqlhandle) roster(username string) (_r []string) {
	if sys.Conf.Property.Tim_sql_roster == "" {
		return
	}
	if rs, _ := this.query(sys.Conf.Property.Tim_sql_roster, username); rs != nil && len(rs) > 0 {
		_r = make([]string, 0)
		for _, bs := range rs {
			_r = append(_r, _getString(bs[0]))
		}
	}
	return
}

func (this *sqlhandle) userGroup(username string) (_r []string) {
	if sys.Conf.Property.Tim_sql_userroom == "" {
		return
	}
	if rs, _ := this.query(sys.Conf.Property.Tim_sql_userroom, username); rs != nil && len(rs) > 0 {
		_r = make([]string, 0)
		for _, bs := range rs {
			_r = append(_r, _getString(bs[0]))
		}
	}
	return
}

func (this *sqlhandle) groupRoster(groupname string) (_r []string) {
	if sys.Conf.Property.Tim_sql_roomroster == "" {
		return
	}
	if rs, _ := this.query(sys.Conf.Property.Tim_sql_roomroster, groupname); rs != nil && len(rs) > 0 {
		_r = make([]string, 0)
		for _, bs := range rs {
			_r = append(_r, _getString(bs[0]))
		}
	}
	return
}

func (this *sqlhandle) addroster(fromname, toname string) (ok bool) {
	if sys.Conf.Property.Tim_sql_roster_add != "" {
		if id, err := this.exec(sys.Conf.Property.Tim_sql_roster_add, fromname, toname); err == nil && id > 0 {
			ok = true
		}
	}
	return
}

func (this *sqlhandle) rmroster(fromname, toname string) (ok bool) {
	if sys.Conf.Property.Tim_sql_roster_rm != "" {
		if id, err := this.exec(sys.Conf.Property.Tim_sql_roster_rm, fromname, toname); err == nil && id > 0 {
			ok = true
		}
	}
	return
}

func (this *sqlhandle) blockroster(fromname, toname string) (ok bool) {
	if sys.Conf.Property.Tim_sql_roster_block != "" {
		if id, err := this.exec(sys.Conf.Property.Tim_sql_roster_block, fromname, toname); err == nil && id > 0 {
			ok = true
		}
	}
	return
}