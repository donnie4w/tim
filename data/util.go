// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim
//

package data

import (
	"os"
	"strconv"

	. "github.com/donnie4w/gofer/lock"
	"github.com/donnie4w/tim/sys"
)

var Handler engine
var numlock = NewNumLock(1 << 8)

func newEngine() engine {
	if sys.UseDefaultDB() {
		sys.DBtype = 1
		return &tldbhandler{}
	} else if sqlHandle.IsAvail() {
		sys.DBtype = 2
		return &sqlhandler{}
	}
	sys.DBtype = 0
	return &nilprocess{}
}

func init() {
	sys.DataInit = Init
}

func Init() (err error) {
	defer func() {
		if err != nil {
			sys.FmtLog(err)
			os.Exit(0)
		}
	}()
	if sys.UseDefaultDB() {
		err = tlormInit()
	} else if sys.Conf.Property != nil {
		err = sqlInit()
	} else {
		sys.FmtLog("no data source provided")
	}
	Handler = newEngine()
	return
}

func _getString(a any) (_r string) {
	switch a.(type) {
	case string:
		_r = a.(string)
	case []uint8:
		_r = string(a.([]uint8))
	case int8, int16, int32, int64, int, uint, uint8, uint16, uint32, uint64:
		_r = strconv.FormatInt(_getInt64(a), 10)
	}
	return
}

func _getBytes(a any) (_r []byte) {
	switch a.(type) {
	case string:
		_r = []byte(a.(string))
	case []uint8:
		_r = a.([]uint8)
	}
	return
}

func _getInt64(a any) (_r int64) {
	switch a.(type) {
	case int8:
		_r = int64(a.(int8))
	case int16:
		_r = int64(a.(int16))
	case int32:
		_r = int64(a.(int32))
	case int64:
		_r = int64(a.(int64))
	case int:
		_r = int64(a.(int))
	case uint:
		_r = int64(a.(uint))
	case uint8:
		_r = int64(a.(uint8))
	case uint16:
		_r = int64(a.(uint16))
	case uint32:
		_r = int64(a.(uint32))
	case uint64:
		_r = int64(a.(uint64))
	}
	return
}
