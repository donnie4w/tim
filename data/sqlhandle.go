// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim
package data

import (
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strings"

	_ "github.com/denisenkom/go-mssqldb"
	. "github.com/donnie4w/gofer/util"
	"github.com/donnie4w/tim/sys"
	_ "github.com/go-sql-driver/mysql"

	_ "github.com/godror/godror"
	_ "github.com/lib/pq"
)

var sqlHandle = &sqlhandle{}

type sqlhandle struct {
	db *sql.DB
}

func sqlInit() (err error) {
	if sys.Conf.Property.Tim_mysql_connection != "" {
		err = sqlHandle.connect("mysql", sys.Conf.Property.Tim_mysql_connection)
	} else if sys.Conf.Property.Tim_postgreSQL_connection != "" {
		err = sqlHandle.connect("postgres", sys.Conf.Property.Tim_postgreSQL_connection)
	} else if sys.Conf.Property.Tim_oracle_connection != "" {
		err = sqlHandle.connect("godror", sys.Conf.Property.Tim_oracle_connection)
	} else if sys.Conf.Property.Tim_sqlserver_connection != "" {
		err = sqlHandle.connect("sqlserver", sys.Conf.Property.Tim_sqlserver_connection)
	} else if sys.Conf.Property.Tim_mysql_connection_mod != "" {
		err = sqlHandle.connect("mysql", parseConnection(sys.Conf.Property.Tim_mysql_connection_mod, mysqlMod))
	} else if sys.Conf.Property.Tim_postgreSQL_connection_mod != "" {
		err = sqlHandle.connect("postgres", parseConnection(sys.Conf.Property.Tim_postgreSQL_connection_mod, postgreSQLMod))
	} else if sys.Conf.Property.Tim_oracle_connection_mod != "" {
		err = sqlHandle.connect("godror", parseConnection(sys.Conf.Property.Tim_oracle_connection_mod, oracleMod))
	} else if sys.Conf.Property.Tim_sqlserver_connection_mod != "" {
		err = sqlHandle.connect("sqlserver", parseConnection(sys.Conf.Property.Tim_sqlserver_connection_mod, sqlserverMod))
	} else {
		err = errors.New("no dataSource provided")
	}
	if err != nil {
		err = fmt.Errorf("%s", "database connect error:"+err.Error())
		return
	}
	return
}

var M_user = fmt.Sprint(CRC64(Int64ToBytes(RandId() + 1)))
var M_pwd = fmt.Sprint(CRC64(Int64ToBytes(RandId() + 2)))
var M_host = fmt.Sprint(CRC64(Int64ToBytes(RandId() + 3)))
var M_port = fmt.Sprint(CRC64(Int64ToBytes(RandId() + 4)))
var M_dbname = fmt.Sprint(CRC64(Int64ToBytes(RandId() + 5)))

var mysqlMod = fmt.Sprint(M_user, ":", M_pwd, "@tcp(", M_host, ":", M_port, ")/", M_dbname)
var postgreSQLMod = fmt.Sprint("host=", M_host, " port=", M_port, " user=", M_user, " password=", M_pwd, " dbname=", M_dbname, " sslmode=disable")
var oracleMod = fmt.Sprint(`user="`, M_user, `" password="`, M_pwd, `" connectString="`, M_host, ":", M_port, "/", M_dbname, `"`)
var sqlserverMod = fmt.Sprint("server=", M_host, ";port=", M_port, ";userid=", M_user, ";password=", M_pwd, ";database=", M_dbname)

func parseConnection(s string, dbmode string) string {
	if ss := strings.Split(strings.TrimSpace(s), " "); len(ss) == 5 {
		for _, sv := range ss {
			idx := strings.Index(strings.TrimSpace(sv), "=")
			k, v := sv[:idx], sv[idx+1:]
			if k == "host" {
				dbmode = strings.ReplaceAll(dbmode, M_host, v)
			} else if k == "port" {
				dbmode = strings.ReplaceAll(dbmode, M_port, v)
			} else if k == "user" {
				dbmode = strings.ReplaceAll(dbmode, M_user, v)
			} else if k == "password" {
				dbmode = strings.ReplaceAll(dbmode, M_pwd, v)
			} else if k == "dbname" {
				dbmode = strings.ReplaceAll(dbmode, M_dbname, v)
			}
		}
	}
	return dbmode
}

func (this *sqlhandle) connect(dbtype, addr string) (err error) {
	if this.db, err = sql.Open(dbtype, addr); err == nil {
		err = this.db.Ping()
	}
	return
}

func (this *sqlhandle) IsAvail() bool {
	return this.db != nil
}

func (this *sqlhandle) Close() error {
	return this.db.Close()
}

func (this *sqlhandle) query(query string, args ...any) (_rs [][]any, err error) {
	var rs *sql.Rows
	if rs, err = this.db.Query(query, args...); err == nil {
		_rs = make([][]any, 0)
		clos, _ := rs.Columns()
		params := len(clos)
		for rs.Next() {
			ps := make([]any, params)
			for i := range ps {
				ps[i] = new(any)
			}
			if err = rs.Scan(ps...); err == nil {
				_r := make([]any, params)
				for i := range _r {
					_r[i] = *(ps[i].(*any))
				}
				_rs = append(_rs, _r)
			}
		}
	}
	return
}

func (this *sqlhandle) insert(query string, args ...any) (id int64, err error) {
	var r sql.Result
	if r, err = this.db.Exec(query, args...); err == nil {
		return r.LastInsertId()
	}
	return
}

func (this *sqlhandle) exec(query string, args ...any) (_r int64, err error) {
	var r sql.Result
	if r, err = this.db.Exec(query, args...); err == nil {
		return r.RowsAffected()
	}
	return
}

// ////////////////////////////////////////////////////////////////////////////////////////////////////
func (this *sqlhandle) getmessage(chatid uint64, id int64, limit int64) (_r map[int64][]byte, err error) {
	var rs [][]any
	if id <= 0 {
		id = 1<<63 - 1
	}
	if rs, err = this.query(sys.Conf.Property.Tim_sql_getmessage, chatid, id, limit); err == nil {
		_r = map[int64][]byte{}
		for _, bs := range rs {
			_r[_getInt64(bs[0])] = _getBytes(bs[1])
		}
	}
	return
}

func (this *sqlhandle) token(username, pwd string) (_r string, err error) {
	var rs [][]any
	if rs, err = this.query(sys.Conf.Property.Tim_sql_token, username, pwd); err == nil {
		for _, bs := range rs {
			_r = _getString(bs[0])
			break
		}
	}
	return
}

func (this *sqlhandle) login(username, pwd string) (_r string, err error) {
	if sys.Conf.Property.Tim_sql_login == "" {
		err = errors.New("empty login sql")
		return
	}
	var rs [][]any
	if rs, err = this.query(sys.Conf.Property.Tim_sql_login, username, pwd); err == nil {
		for _, bs := range rs {
			_r = _getString(bs[0])
			break
		}
	}
	return
}

func (this *sqlhandle) saveMessage(chatId uint64, stanza []byte) (mid int64, err error) {
	return this.insert(sys.Conf.Property.Tim_sql_savemessage, chatId, stanza)
}

func (this *sqlhandle) getMessageById(id int64) (bs []byte, err error) {
	var rs [][]any
	if rs, err = this.query(sys.Conf.Property.Tim_sql_getmessage_byid, id); err == nil && len(rs) > 0 {
		bs = _getBytes(rs[0][0])
	}
	return
}

func (this *sqlhandle) delMessageById(id int64) (err error) {
	_, err = this.exec(sys.Conf.Property.Tim_sql_delmessage_byid, id)
	return
}

func (this *sqlhandle) saveOfflineMessage(node string, uniqueid int64, stanza []byte, mid int64) (err error) {
	if sys.Conf.Property.Tim_sql_offlinemsg_save != "" {
		_, err = this.insert(sys.Conf.Property.Tim_sql_offlinemsg_save, node, uniqueid, stanza, mid)
	} else if sys.Conf.Property.Tim_sql_offlinemsg_save_mid != "" && mid > 0 {
		_, err = this.insert(sys.Conf.Property.Tim_sql_offlinemsg_save_mid, node, uniqueid, mid)
	} else if sys.Conf.Property.Tim_sql_offlinemsg_save_nomid != "" && mid == 0 {
		_, err = this.insert(sys.Conf.Property.Tim_sql_offlinemsg_save_nomid, node, uniqueid, stanza)
	}
	return
}

func (this *sqlhandle) getOfflineMessage(username string, limit int) (obList []*OfflineBean, err error) {
	var rs [][]any
	if rs, err = this.query(sys.Conf.Property.Tim_sql_offlinemsg_get, username, limit); err == nil {
		obList = make([]*OfflineBean, 0)
		for _, bs := range rs {
			ob := &OfflineBean{Id: _getInt64(bs[0]), Stanze: _getBytes(bs[1])}
			obList = append(obList, ob)
		}
	}
	return
}

func (this *sqlhandle) delOfflineMessage(ids ...int64) (_r int64, err error) {
	_ids := make([]any, 0)
	if sys.Conf.Property.Tim_sql_offlinemsg_del != "" {
		for _, id := range ids {
			_r, err = this.exec(sys.Conf.Property.Tim_sql_offlinemsg_del, id)
		}
	} else if sys.Conf.Property.Tim_sql_offlinemsg_delin != "" {
		sb := &strings.Builder{}
		sb.WriteString("in(")
		for i := 0; i < len(ids)-1; i++ {
			sb.WriteString("?,")
		}
		sb.WriteString("?)")
		re := regexp.MustCompile(`in\s{0,}\(.*?\)`)
		s := re.ReplaceAll([]byte(sys.Conf.Property.Tim_sql_offlinemsg_delin), []byte(sb.String()))
		for _, i := range ids {
			_ids = append(_ids, i)
		}
		_r, err = this.exec(string(s), _ids...)
	}
	return
}

// func (this *sqlhandle) existOfflineMessage(node string, uniqueId int64) (exist bool, err error) {
// 	if sys.Conf.Property.Tim_sql_offlinemsg_exist != "" {
// 		if rs, e := this.query(sys.Conf.Property.Tim_sql_offlinemsg_exist, node, uniqueId); e == nil {
// 			for _, bs := range rs {
// 				exist = _getInt64(bs[0]) > 0
// 				break
// 			}
// 		} else {
// 			err = e
// 		}
// 	}
// 	return
// }

func (this *sqlhandle) authUser(fromnode, tonode string) (_r bool, err error) {
	if sys.Conf.Property.Tim_sql_authuser == "" {
		return true, nil
	}
	var rs [][]any
	if rs, err = this.query(sys.Conf.Property.Tim_sql_authuser, fromnode, tonode); err == nil {
		for _, bs := range rs {
			_r = _getInt64(bs[0]) > 0
			break
		}
	}
	return
}

func (this *sqlhandle) authGroup(groupnode, usernode string) (_r bool, err error) {
	if sys.Conf.Property.Tim_sql_authroom == "" {
		return true, nil
	}
	var rs [][]any
	if rs, err = this.query(sys.Conf.Property.Tim_sql_authroom, groupnode, usernode); err == nil {
		for _, bs := range rs {
			_r = _getInt64(bs[0]) > 0
			break
		}
	}
	return
}

func (this *sqlhandle) existUser(node string) (_r bool, err error) {
	if sys.Conf.Property.Tim_sql_existuser == "" {
		return true, nil
	}
	var rs [][]any
	if rs, err = this.query(sys.Conf.Property.Tim_sql_existuser, node); err == nil {
		for _, bs := range rs {
			_r = _getInt64(bs[0]) > 0
			break
		}
	}
	return
}

func (this *sqlhandle) existGroup(node string) (_r bool, err error) {
	if sys.Conf.Property.Tim_sql_existroom == "" {
		return true, nil
	}
	var rs [][]any
	if rs, err = this.query(sys.Conf.Property.Tim_sql_existroom, node); err == nil {
		for _, bs := range rs {
			_r = _getInt64(bs[0]) > 0
			break
		}
	}
	return
}
