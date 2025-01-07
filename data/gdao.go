// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package data

import (
	"database/sql"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/donnie4w/gdao"
	"github.com/donnie4w/gdao/base"
	"github.com/donnie4w/gofer/util"
	"github.com/donnie4w/tim/stub"
	"github.com/donnie4w/tim/sys"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/godror/godror"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"sort"
	"strings"
)

var (
	m_user   = fmt.Sprint(util.UUID64())
	m_pwd    = fmt.Sprint(util.UUID64())
	m_host   = fmt.Sprint(util.UUID64())
	m_port   = fmt.Sprint(util.UUID64())
	m_dbname = fmt.Sprint(util.UUID64())

	mysqlMod      = fmt.Sprint(m_user, ":", m_pwd, "@tcp(", m_host, ":", m_port, ")/", m_dbname)
	postgreSQLMod = fmt.Sprint("host=", m_host, " port=", m_port, " user=", m_user, " password=", m_pwd, " dbname=", m_dbname, " sslmode=disable")
	oracleMod     = fmt.Sprint(`user="`, m_user, `" password="`, m_pwd, `" connectstring="`, m_host, ":", m_port, "/", m_dbname, `"`)
	sqlserverMod  = fmt.Sprint("server=", m_host, ";port=", m_port, ";userid=", m_user, ";password=", m_pwd, ";database=", m_dbname)
)

func parseConnection(s string, dbmode string) string {
	if ss := strings.Split(strings.TrimSpace(s), " "); len(ss) == 5 {
		for _, sv := range ss {
			idx := strings.Index(strings.TrimSpace(sv), "=")
			k, v := sv[:idx], sv[idx+1:]
			if k == "host" {
				dbmode = strings.ReplaceAll(dbmode, m_host, v)
			} else if k == "port" {
				dbmode = strings.ReplaceAll(dbmode, m_port, v)
			} else if k == "user" {
				dbmode = strings.ReplaceAll(dbmode, m_user, v)
			} else if k == "password" {
				dbmode = strings.ReplaceAll(dbmode, m_pwd, v)
			} else if k == "dbname" {
				dbmode = strings.ReplaceAll(dbmode, m_dbname, v)
			}
		}
	}
	return dbmode
}

func getDBhandle(driverName, dataSourceName string) (dbhandle base.DBhandle, err error) {
	if db, er := sql.Open(driverName, dataSourceName); er == nil {
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

func (sh *sourceHandle) Add(driverName, dataSourceName string, extent int) error {
	if extent == 0 {
		extent = sys.MB
	}
	if dbhandle, err := getDBhandle(driverName, dataSourceName); err == nil {
		initLocalDB(dbhandle, driverName)
		sh.dbsource = append(sh.dbsource, &sqlsource{extent: extent, dbhandle: dbhandle})
	} else {
		return err
	}
	sort.Slice(sh.dbsource, func(i, j int) bool { return sh.dbsource[i].extent < sh.dbsource[j].extent })
	return nil
}

func (sh *sourceHandle) AddInlineDB(v *stub.InlineDB) error {
	if v.Tim_sqlite_connection != "" {
		return sh.Add(Driver_Sqlite, v.Tim_sqlite_connection, v.ExtentMax)
	} else if v.Tim_postgresql_connection != "" {
		return sh.Add(Driver_Postgres, v.Tim_postgresql_connection, v.ExtentMax)
	} else if v.Tim_mysql_connection != "" {
		return sh.Add(Driver_Mysql, v.Tim_mysql_connection, v.ExtentMax)
	} else if v.Tim_sqlserver_connection != "" {
		return sh.Add(Driver_Sqlserver, v.Tim_sqlserver_connection, v.ExtentMax)
	} else if v.Tim_postgresql_connection_mod != "" {
		return sh.Add(Driver_Postgres, v.Tim_postgresql_connection, v.ExtentMax)
	} else if v.Tim_mysql_connection_mod != "" {
		return sh.Add(Driver_Mysql, v.Tim_mysql_connection, v.ExtentMax)
	} else if v.Tim_sqlserver_connection_mod != "" {
		return sh.Add(Driver_Sqlserver, v.Tim_sqlserver_connection, v.ExtentMax)
	}
	return nil
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

func lastInsertId(dbtype base.DBType, tx base.Transaction) (int64, bool) {
	var s string
	switch dbtype {
	case gdao.SQLITE:
		s = "SELECT LAST_INSERT_ROWID()"
	case gdao.MYSQL, gdao.MARIADB, gdao.TIDB, gdao.OCEANBASE:
		s = "SELECT LAST_INSERT_ID()"
	case gdao.SQLSERVER:
		s = "SELECT SCOPE_IDENTITY()"
	case gdao.POSTGRESQL, gdao.GREENPLUM, gdao.OPENGAUSS:
		return 0, true
	default:
		return 0, false
	}
	return tx.ExecuteQueryBean(s).ToInt64(), false
}
