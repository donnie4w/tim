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
	"github.com/donnie4w/gdao"
	"github.com/donnie4w/gdao/base"
	"regexp"
	"strings"

	"github.com/donnie4w/tim/sys"
)

type externdb struct {
	db base.DBhandle
}

func (eh *externdb) init() (err error) {
	if sys.Conf.ExternalDB.MYSQL != nil {
		err = eh.connect(sys.Conf.ExternalDB.MYSQL.DSN(gdao.MYSQL))
	} else if sys.Conf.ExternalDB.POSTGRESQL != nil {
		err = eh.connect(sys.Conf.ExternalDB.POSTGRESQL.DSN(gdao.POSTGRESQL))
	} else if sys.Conf.ExternalDB.ORACLE != nil {
		err = eh.connect(sys.Conf.ExternalDB.ORACLE.DSN(gdao.ORACLE))
	} else if sys.Conf.ExternalDB.SQLSERVER != nil {
		err = eh.connect(sys.Conf.ExternalDB.SQLSERVER.DSN(gdao.SQLSERVER))
	} else if sys.Conf.ExternalDB.MARIADB != nil {
		err = eh.connect(sys.Conf.ExternalDB.MARIADB.DSN(gdao.MARIADB))
	} else if sys.Conf.ExternalDB.OCEANBASE != nil {
		err = eh.connect(sys.Conf.ExternalDB.OCEANBASE.DSN(gdao.OCEANBASE))
	} else if sys.Conf.ExternalDB.TIDB != nil {
		err = eh.connect(sys.Conf.ExternalDB.TIDB.DSN(gdao.TIDB))
	} else if sys.Conf.ExternalDB.GREENPLUM != nil {
		err = eh.connect(sys.Conf.ExternalDB.GREENPLUM.DSN(gdao.GREENPLUM))
	} else if sys.Conf.ExternalDB.OPENGAUSS != nil {
		err = eh.connect(sys.Conf.ExternalDB.OPENGAUSS.DSN(gdao.OPENGAUSS))
	} else if sys.Conf.ExternalDB.COCKROACHDB != nil {
		err = eh.connect(sys.Conf.ExternalDB.COCKROACHDB.DSN(gdao.COCKROACHDB))
	} else if sys.Conf.ExternalDB.ENTERPRISEDB != nil {
		err = eh.connect(sys.Conf.ExternalDB.ENTERPRISEDB.DSN(gdao.ENTERPRISEDB))
	} else {
		err = errors.New("no dataSource provided")
	}
	if err != nil {
		panic(fmt.Sprint("%s", "databean connect error:"+err.Error()))
	}
	return
}

func (eh *externdb) connect(driverName, dataSourceName string) (err error) {
	eh.db, err = getDBhandle(driverName, dataSourceName)
	return
}

func (eh *externdb) Close() error {
	return eh.db.Close()
}

func (eh *externdb) insert(query string, args ...any) (id int64, err error) {
	var r sql.Result
	if r, err = eh.db.ExecuteUpdate(query, args...); err == nil {
		return r.LastInsertId()
	}
	return
}

func (eh *externdb) exec(query string, args ...any) (_r int64, err error) {
	var r sql.Result
	if r, err = eh.db.ExecuteUpdate(query, args...); err == nil {
		return r.RowsAffected()
	}
	return
}

// ////////////////////////////////////////////////////////////////////////////////////////////////////
func (eh *externdb) getmessage(chatid []byte, mid int64, timeseries int64, limit int64) (_r map[int64][]byte, err error) {
	//if mid <= 0 {
	//	mid = 1<<63 - 1
	//}
	databeans := eh.db.ExecuteQueryBeans(sys.Conf.ExternalDB.Sql_message_get, chatid, mid, timeseries, limit)
	if err = databeans.GetError(); err == nil {
		_r = map[int64][]byte{}
		for _, databean := range databeans.Beans {
			_r[databean.FieldByIndex(0).ValueInt64()] = databean.FieldByIndex(1).ValueBytes()
		}
	}

	return
}

func (eh *externdb) authNode(username, pwd string) (_r string, err error) {
	databean := eh.db.ExecuteQueryBean(sys.Conf.ExternalDB.Sql_token, username, pwd)
	if err = databean.GetError(); err == nil {
		_r = databean.FieldByIndex(0).ValueString()
	}
	return
}

func (eh *externdb) login(username, pwd string) (_r string, err error) {
	if sys.Conf.ExternalDB.Sql_login == "" {
		err = errors.New("empty login sql")
		return
	}
	databean := eh.db.ExecuteQueryBean(sys.Conf.ExternalDB.Sql_login, username, pwd)
	if err = databean.GetError(); err == nil {
		_r = databean.FieldByIndex(0).ValueString()
	}
	return
}

func (eh *externdb) saveMessage(chatId []byte, fid int32, stanza []byte) (mid int64, err error) {
	return eh.insert(sys.Conf.ExternalDB.Sql_message_save, chatId, fid, stanza)
}

func (eh *externdb) getFidById(chatId []byte, id int64) (fid int64, err error) {
	databean := eh.db.ExecuteQueryBean(sys.Conf.ExternalDB.Sql_message_fid, chatId, id)
	if err = databean.GetError(); err == nil {
		fid = databean.FieldByIndex(0).ValueInt64()
	}
	return
}

func (eh *externdb) delMessageById(id int64) (err error) {
	_, err = eh.exec(sys.Conf.ExternalDB.Sql_message_del, id)
	return
}

func (eh *externdb) saveOfflineMessage(node string, uniqueid int64, stanza []byte, mid int64) (err error) {
	if sys.Conf.ExternalDB.Sql_offlinemsg_save != "" {
		_, err = eh.insert(sys.Conf.ExternalDB.Sql_offlinemsg_save, node, uniqueid, stanza, mid)
	} else if sys.Conf.ExternalDB.Sql_offlinemsg_save_mid != "" && mid > 0 {
		_, err = eh.insert(sys.Conf.ExternalDB.Sql_offlinemsg_save_mid, node, uniqueid, mid)
	} else if sys.Conf.ExternalDB.Sql_offlinemsg_save_nomid != "" && mid == 0 {
		_, err = eh.insert(sys.Conf.ExternalDB.Sql_offlinemsg_save_nomid, node, uniqueid, stanza)
	}
	return
}

func (eh *externdb) getOfflineMessage(username string, limit int) (obList []*OfflineBean, err error) {
	databeans := eh.db.ExecuteQueryBeans(sys.Conf.ExternalDB.Sql_offlinemsg_get, username, limit)
	if err = databeans.GetError(); err == nil {
		obList = []*OfflineBean{}
		for _, databean := range databeans.Beans {
			obList = append(obList, &OfflineBean{Id: databean.FieldByIndex(0).ValueInt64(), Stanze: databean.FieldByIndex(1).ValueBytes()})
		}
	}
	return
}

func (eh *externdb) delOfflineMessage(ids ...any) (_r int64, err error) {
	_ids := make([]any, 0)
	if sys.Conf.ExternalDB.Sql_offlinemsg_del != "" {
		for _, id := range ids {
			_r, err = eh.exec(sys.Conf.ExternalDB.Sql_offlinemsg_del, id)
		}
	} else if sys.Conf.ExternalDB.Sql_offlinemsg_delin != "" {
		sb := &strings.Builder{}
		sb.WriteString("in(")
		for i := 0; i < len(ids)-1; i++ {
			sb.WriteString("?,")
		}
		sb.WriteString("?)")
		re := regexp.MustCompile(`in\s{0,}\(.*?\)`)
		s := re.ReplaceAll([]byte(sys.Conf.ExternalDB.Sql_offlinemsg_delin), []byte(sb.String()))
		for _, i := range ids {
			_ids = append(_ids, i)
		}
		_r, err = eh.exec(string(s), _ids...)
	}
	return
}

func (eh *externdb) authUser(fromnode, tonode string) (_r bool, err error) {
	if sys.Conf.ExternalDB.Sql_authuser == "" {
		return true, nil
	}
	databean := eh.db.ExecuteQueryBean(sys.Conf.ExternalDB.Sql_authuser, fromnode, tonode)
	if err = databean.GetError(); err == nil {
		_r = databean.FieldByIndex(0).ValueInt64() > 0
	}
	return
}

func (eh *externdb) authGroup(groupnode, usernode string) (_r bool, err error) {
	if sys.Conf.ExternalDB.Sql_authroom == "" {
		return true, nil
	}
	databean := eh.db.ExecuteQueryBean(sys.Conf.ExternalDB.Sql_authroom, groupnode, usernode)
	if err = databean.GetError(); err == nil {
		_r = databean.FieldByIndex(0).ValueInt64() > 0
	}
	return
}

func (eh *externdb) existUser(node string) (_r bool, err error) {
	if sys.Conf.ExternalDB.Sql_existuser == "" {
		return true, nil
	}
	databean := eh.db.ExecuteQueryBean(sys.Conf.ExternalDB.Sql_existuser, node)
	if err = databean.GetError(); err == nil {
		_r = databean.FieldByIndex(0).ValueInt64() > 0
	}
	return
}

func (eh *externdb) existGroup(node string) (_r bool, err error) {
	if sys.Conf.ExternalDB.Sql_existroom == "" {
		return true, nil
	}
	databean := eh.db.ExecuteQueryBean(sys.Conf.ExternalDB.Sql_existroom, node)
	if err = databean.GetError(); err == nil {
		_r = databean.FieldByIndex(0).ValueInt64() > 0
	}
	return
}

/*****************************************************************************************************************/

func (eh *externdb) roster(username string) (_r []string) {
	if sys.Conf.ExternalDB.Sql_roster == "" {
		return
	}
	databeans := eh.db.ExecuteQueryBeans(sys.Conf.ExternalDB.Sql_roster, username)
	if databeans != nil {
		for _, databean := range databeans.Beans {
			_r = append(_r, databean.FieldByIndex(0).ValueString())
		}
	}
	return
}

func (eh *externdb) userGroup(username string) (_r []string) {
	if sys.Conf.ExternalDB.Sql_userroom == "" {
		return
	}
	databeans := eh.db.ExecuteQueryBeans(sys.Conf.ExternalDB.Sql_userroom, username)
	if databeans != nil {
		for _, databean := range databeans.Beans {
			_r = append(_r, databean.FieldByIndex(0).ValueString())
		}
	}
	return
}

func (eh *externdb) groupRoster(groupname string) (_r []string) {
	if sys.Conf.ExternalDB.Sql_roomroster == "" {
		return
	}
	databeans := eh.db.ExecuteQueryBeans(sys.Conf.ExternalDB.Sql_roomroster, groupname)
	if databeans != nil {
		for _, databean := range databeans.Beans {
			_r = append(_r, databean.FieldByIndex(0).ValueString())
		}
	}
	return
}

func (eh *externdb) addroster(fromname, toname string) (ok bool) {
	if sys.Conf.ExternalDB.Sql_roster_add != "" {
		if id, err := eh.exec(sys.Conf.ExternalDB.Sql_roster_add, fromname, toname); err == nil && id > 0 {
			ok = true
		}
	}
	return
}

func (eh *externdb) rmroster(fromname, toname string) (ok bool) {
	if sys.Conf.ExternalDB.Sql_roster_rm != "" {
		if id, err := eh.exec(sys.Conf.ExternalDB.Sql_roster_rm, fromname, toname); err == nil && id > 0 {
			ok = true
		}
	}
	return
}

func (eh *externdb) blockroster(fromname, toname string) (ok bool) {
	if sys.Conf.ExternalDB.Sql_roster_block != "" {
		if id, err := eh.exec(sys.Conf.ExternalDB.Sql_roster_block, fromname, toname); err == nil && id > 0 {
			ok = true
		}
	}
	return
}
