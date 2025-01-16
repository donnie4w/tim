// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package data

import (
	"github.com/donnie4w/gdao/base"
	"github.com/donnie4w/tim/log"
	"os"
	"strconv"
	"time"

	"github.com/donnie4w/gofer/lock"
	"github.com/donnie4w/tim/sys"
)

var Service service
var numlock = lock.NewNumLock(1 << 9)
var strlock = lock.NewStrlock(1 << 9)

func getService() service {
	switch sys.GetDBMOD() {
	case sys.INLINEDB:
		return new(inlineHandle).init()
	case sys.TLDB:
		return new(tldbhandle).init()
	case sys.EXTERNALDB:
		return new(externalhandle).init()
	case sys.NODB:
		return new(nodbhandle)
	}
	panic("No supported service found")
}

func init() {
	sys.Service(sys.INIT_DATA, serv(1))
}

type serv byte

func (serv) Serve() error {
	defer func() {
		if err := recover(); err != nil {
			log.FmtPrint(err)
			os.Exit(0)
		}
	}()
	Service = getService()
	return nil
}

func (serv) Close() error {
	return nil
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
		_r = a.(int64)
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

const (
	Driver_Sqlite    = "sqlite3"
	Driver_Postgres  = "postgres"
	Driver_Mysql     = "mysql"
	Driver_Sqlserver = "sqlserver"
	Driver_Oracle    = "godror"
)

const localDB = "tim.db"

func initLocalDB(dbhandle base.DBhandle, driverName string) {
	var ss []string
	switch driverName {
	case Driver_Sqlite:
		ss = sys.Sqlite("").CreateSql()
	case Driver_Postgres:
		ss = sys.PostgreSql("").CreateSql()
	case Driver_Mysql:
		ss = sys.Mysql("").CreateSql()
	case Driver_Sqlserver:
		ss = sys.SqlServer("").CreateSql()
	}
	if len(ss) > 0 {
		for _, s := range ss {
			dbhandle.ExecuteUpdate(s)
		}
	}
}

func TimeNano() int64 {
	return time.Now().UnixNano()
}
