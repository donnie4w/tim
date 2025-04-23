// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package data

import (
	"database/sql"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/donnie4w/gdao"
	"github.com/donnie4w/gdao/base"
	"github.com/donnie4w/tim/stub"
	"github.com/donnie4w/tim/sys"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/godror/godror"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"sort"
)

func getDBhandle(driverName, dsn string) (dbhandle base.DBhandle, err error) {
	if db, er := sql.Open(driverName, dsn); er == nil {
		if err = db.Ping(); err == nil {
			dbhandle = gdao.NewDBHandle(db, getDBtype(driverName))
		}
	} else {
		err = er
	}
	return
}

func getDBtype(driverName string) (dbType base.DBType) {
	switch driverName {
	case Driver_Mysql:
		dbType = gdao.MYSQL
	case Driver_Postgres:
		dbType = gdao.POSTGRESQL
	case Driver_Sqlite:
		dbType = gdao.SQLITE
	case Driver_Oracle:
		dbType = gdao.ORACLE
	case Driver_Sqlserver:
		dbType = gdao.SQLSERVER
	}
	return
}

var gdaoHandle = &sourceHandle{dbsource: make([]*sqlsource, 0)}

type sqlsource struct {
	extent   int
	dbhandle base.DBhandle
}

type sourceHandle struct {
	dbsource []*sqlsource
}

func (sh *sourceHandle) AddConnect(dbType base.DBType, connect *stub.Connect, extent int) error {
	if extent == 0 {
		extent = sys.MB
	}
	driver, dsn := connect.DSN(dbType)
	if dbhandle, err := getDBhandle(driver, dsn); err == nil {
		initLocalDB(dbhandle, driver)
		sh.dbsource = append(sh.dbsource, &sqlsource{extent: extent, dbhandle: dbhandle})
	} else {
		return err
	}
	sort.Slice(sh.dbsource, func(i, j int) bool { return sh.dbsource[i].extent < sh.dbsource[j].extent })
	return nil
}

func (sh *sourceHandle) AddInlineDB(v *stub.InlineDB) error {
	if v.SQLITE != nil {
		return sh.AddConnect(gdao.SQLITE, v.SQLITE, v.ExtentMax)
	} else if v.POSTGRESQL != nil {
		return sh.AddConnect(gdao.POSTGRESQL, v.POSTGRESQL, v.ExtentMax)
	} else if v.GREENPLUM != nil {
		return sh.AddConnect(gdao.GREENPLUM, v.GREENPLUM, v.ExtentMax)
	} else if v.OPENGAUSS != nil {
		return sh.AddConnect(gdao.OPENGAUSS, v.OPENGAUSS, v.ExtentMax)
	} else if v.COCKROACHDB != nil {
		return sh.AddConnect(gdao.COCKROACHDB, v.COCKROACHDB, v.ExtentMax)
	} else if v.ENTERPRISEDB != nil {
		return sh.AddConnect(gdao.ENTERPRISEDB, v.ENTERPRISEDB, v.ExtentMax)
	} else if v.MYSQL != nil {
		return sh.AddConnect(gdao.MYSQL, v.MYSQL, v.ExtentMax)
	} else if v.MARIADB != nil {
		return sh.AddConnect(gdao.MARIADB, v.MARIADB, v.ExtentMax)
	} else if v.OCEANBASE != nil {
		return sh.AddConnect(gdao.OCEANBASE, v.OCEANBASE, v.ExtentMax)
	} else if v.TIDB != nil {
		return sh.AddConnect(gdao.TIDB, v.TIDB, v.ExtentMax)
	} else if v.SQLSERVER != nil {
		return sh.AddConnect(gdao.SQLSERVER, v.SQLSERVER, v.ExtentMax)
	} else if v.ORACLE != nil {
		return sh.AddConnect(gdao.ORACLE, v.ORACLE, v.ExtentMax)
	} else {
		return sh.AddConnect(gdao.SQLITE, &stub.Connect{DBname: localdb}, v.ExtentMax)
	}
}

func (sh *sourceHandle) GetDBHandle(tid uint64) base.DBhandle {
	if len(sh.dbsource) == 0 {
		return nil
	}
	if len(sh.dbsource) == 1 || (tid == 0) {
		return sh.dbsource[0].dbhandle
	}
	idx := 0
	if idx = sort.Search(len(sh.dbsource), func(i int) bool { return sh.dbsource[i].extent >= int(tid%sys.MB) }); idx >= len(sh.dbsource) {
		idx = 0
	}
	return sh.dbsource[idx].dbhandle
}

func (sh *sourceHandle) FirstDBHandle() base.DBhandle {
	if len(sh.dbsource) > 0 {
		return sh.dbsource[0].dbhandle
	}
	return nil
}
