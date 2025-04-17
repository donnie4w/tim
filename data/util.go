// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package data

import (
	"github.com/donnie4w/gdao/base"
	gocache "github.com/donnie4w/gofer/cache"
	"github.com/donnie4w/gofer/lock"
	goutil "github.com/donnie4w/gofer/util"
	"github.com/donnie4w/tim/log"
	"github.com/donnie4w/tim/sys"
	"math/rand"
	"os"
	"strconv"
	"time"
)

var Service service
var numlock = lock.NewNumLock(1 << 9)
var strlock = lock.NewStrlock(1 << 9)

func getService() service {
	switch sys.GetDBMOD() {
	case sys.INLINEDB:
		return new(inlineHandle).init()
	case sys.TLDB:
		return new(tldbHandle).init()
	case sys.MONGODB:
		return new(mongoHandle).init()
	case sys.EXTERNALDB:
		return new(externHandle).init()
	case sys.NODB:
		return new(nodbHandle).init()
	case sys.CASSANDRA:
		return new(cassandraHandle).init()
	}
	return nil
}

func init() {
	sys.Service(sys.INIT_DATA, serv(1))
}

type serv byte

func (serv) Serve() error {
	defer func() {
		if err := recover(); err != nil {
			log.FmtPrint(err)
			os.Exit(1)
		}
	}()
	if Service = getService(); Service == nil {
		panic("no database service found")
	}
	return nil
}

func (serv) Close() error {
	return nil
}

var uuidCache = gocache.NewBloomFilter(1<<19, 0.01)

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

const localdb = "tim.db"

func initLocalDB(dbhandle base.DBhandle, driver string) {
	var ss []string
	switch driver {
	case Driver_Sqlite:
		ss = Sqlite("").CreateSql()
	case Driver_Postgres:
		ss = PostgreSql("").CreateSql()
	case Driver_Mysql:
		ss = Mysql("").CreateSql()
	case Driver_Sqlserver:
		ss = SqlServer("").CreateSql()
	case Driver_Oracle:
		ss = Oracle("").CreateSql()
	}
	if len(ss) > 0 {
		for _, s := range ss {
			dbhandle.ExecuteUpdate(s)
		}
	}
}

func TimeNano() int64 {
	return timenano.unix()
}

var timenano = newNano()

type nano struct {
	precision byte
	count     uint64
}

func newNano() *nano {
	t := time.Now()
	if t.UnixNano() == t.UnixNano()/1e3*1e3 {
		return &nano{precision: 3}
	} else if t.UnixNano() == t.UnixNano()/1e2*1e2 {
		return &nano{precision: 2}
	} else if t.UnixNano() == t.UnixNano()/1e1*1e1 {
		return &nano{precision: 1}
	} else {
		return &nano{precision: 0}
	}
}

func (n *nano) unix() (r int64) {
	switch n.precision {
	case 3:
		r = time.Now().UnixNano() + int64(n.count%1e3)
	case 2:
		r = time.Now().UnixNano() + int64(n.count%1e2)
	case 1:
		r = time.Now().UnixNano() + int64(n.count%1e1)
	default:
		r = time.Now().UnixNano()
	}
	n.count++
	return r
}

var r = rand.New(rand.NewSource(goutil.UUID64()))

func midUUID() int64 {
	switch sys.Conf.UuidBits {
	case 32:
		return int64(goutil.UUID32())
	case 53:
		uuid32 := uint32(goutil.UUID32())
		return int64(uuid32)<<21 | r.Int63n(1<<21)
	default:
		return goutil.UUID64()
	}
}
