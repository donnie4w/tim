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
	"github.com/donnie4w/gdao/base"
	"regexp"
	"strings"

	"github.com/donnie4w/tim/sys"
)

type externaldb struct {
	db base.DBhandle
}

func (eh *externaldb) init() {
	var err error
	if sys.Conf.ExternalDB.Tim_mysql_connection != "" {
		err = eh.connect(Driver_Mysql, sys.Conf.ExternalDB.Tim_mysql_connection)
	} else if sys.Conf.ExternalDB.Tim_postgresql_connection != "" {
		err = eh.connect(Driver_Postgres, sys.Conf.ExternalDB.Tim_postgresql_connection)
	} else if sys.Conf.ExternalDB.Tim_oracle_connection != "" {
		err = eh.connect(Driver_Oracle, sys.Conf.ExternalDB.Tim_oracle_connection)
	} else if sys.Conf.ExternalDB.Tim_sqlserver_connection != "" {
		err = eh.connect(Driver_Sqlserver, sys.Conf.ExternalDB.Tim_sqlserver_connection)
	} else if sys.Conf.ExternalDB.Tim_mysql_connection_mod != "" {
		err = eh.connect(Driver_Mysql, parseConnection(sys.Conf.ExternalDB.Tim_mysql_connection_mod, mysqlMod))
	} else if sys.Conf.ExternalDB.Tim_postgresql_connection_mod != "" {
		err = eh.connect(Driver_Postgres, parseConnection(sys.Conf.ExternalDB.Tim_postgresql_connection_mod, postgreSQLMod))
	} else if sys.Conf.ExternalDB.Tim_oracle_connection_mod != "" {
		err = eh.connect(Driver_Oracle, parseConnection(sys.Conf.ExternalDB.Tim_oracle_connection_mod, oracleMod))
	} else if sys.Conf.ExternalDB.Tim_sqlserver_connection_mod != "" {
		err = eh.connect(Driver_Sqlserver, parseConnection(sys.Conf.ExternalDB.Tim_sqlserver_connection_mod, sqlserverMod))
	} else {
		err = errors.New("no dataSource provided")
	}
	if err != nil {
		panic(fmt.Sprint("%s", "databean connect error:"+err.Error()))
	}
	return
}

func (eh *externaldb) connect(driverName, dataSourceName string) (err error) {
	eh.db, err = getDBhandle(driverName, dataSourceName)
	return
}

//func (eh *externaldb) IsAvail() bool {
//	return eh.db != nil
//}

func (eh *externaldb) Close() error {
	return eh.db.Close()
}

//func (eh *externaldb) query(query string, args ...any) (_rs [][]any, err error) {
//	var rs *sql.Rows
//	if rs, err = eh.db.Query(query, args...); err == nil {
//		_rs = make([][]any, 0)
//		clos, _ := rs.Columns()
//		params := len(clos)
//		for rs.Next() {
//			ps := make([]any, params)
//			for i := range ps {
//				ps[i] = new(any)
//			}
//			if err = rs.Scan(ps...); err == nil {
//				_r := make([]any, params)
//				for i := range _r {
//					_r[i] = *(ps[i].(*any))
//				}
//				_rs = append(_rs, _r)
//			}
//		}
//	}
//	return
//}

func (eh *externaldb) insert(query string, args ...any) (id int64, err error) {
	var r sql.Result
	if r, err = eh.db.ExecuteUpdate(query, args...); err == nil {
		return r.LastInsertId()
	}
	return
}

func (eh *externaldb) exec(query string, args ...any) (_r int64, err error) {
	var r sql.Result
	if r, err = eh.db.ExecuteUpdate(query, args...); err == nil {
		return r.RowsAffected()
	}
	return
}

// ////////////////////////////////////////////////////////////////////////////////////////////////////
func (eh *externaldb) getmessage(chatid []byte, id int64, limit int64) (_r map[int64][]byte, err error) {
	//var rs [][]any
	if id <= 0 {
		id = 1<<63 - 1
	}
	//if rs, err = eh.query(sys.Conf.ExternalDB.Tim_sql_getmessage, chatid, id, limit); err == nil {
	//	_r = map[int64][]byte{}
	//	for _, bs := range rs {
	//		_r[data._getInt64(bs[0])] = data._getBytes(bs[1])
	//	}
	//}
	databeans := eh.db.ExecuteQueryBeans(sys.Conf.ExternalDB.Tim_sql_getmessage, chatid, id, limit)
	if err = databeans.GetError(); err == nil {
		_r = map[int64][]byte{}
		for _, databean := range databeans.Beans {
			//_r[data._getInt64(bs[0])] = data._getBytes(bs[1])
			_r[databean.FieldByIndex(0).ValueInt64()] = databean.FieldByIndex(1).ValueBytes()
		}
	}

	return
}

func (eh *externaldb) authNode(username, pwd string) (_r string, err error) {
	//var rs [][]any
	//if rs, err = eh.query(sys.Conf.ExternalDB.Tim_sql_token, username, pwd); err == nil {
	//	for _, bs := range rs {
	//		_r = data._getString(bs[0])
	//		break
	//	}
	//}
	databean := eh.db.ExecuteQueryBean(sys.Conf.ExternalDB.Tim_sql_token, username, pwd)
	if err = databean.GetError(); err == nil {
		_r = databean.FieldByIndex(0).ValueString()
	}
	return
}

func (eh *externaldb) login(username, pwd string) (_r string, err error) {
	if sys.Conf.ExternalDB.Tim_sql_login == "" {
		err = errors.New("empty login sql")
		return
	}
	//var rs [][]any
	//if rs, err = eh.query(sys.Conf.ExternalDB.Tim_sql_login, username, pwd); err == nil {
	//	for _, bs := range rs {
	//		_r = data._getString(bs[0])
	//		break
	//	}
	//}
	databean := eh.db.ExecuteQueryBean(sys.Conf.ExternalDB.Tim_sql_login, username, pwd)
	if err = databean.GetError(); err == nil {
		_r = databean.FieldByIndex(0).ValueString()
	}
	return
}

func (eh *externaldb) saveMessage(chatId []byte, fid int32, stanza []byte) (mid int64, err error) {
	return eh.insert(sys.Conf.ExternalDB.Tim_sql_savemessage, chatId, fid, stanza)
}

func (eh *externaldb) getChatIdById(id int64) (bs []byte, fid int64, err error) {
	databean := eh.db.ExecuteQueryBean(sys.Conf.ExternalDB.Tim_sql_getchatid_byid, id)
	if err = databean.GetError(); err == nil {
		bs = databean.FieldByIndex(0).ValueBytes()
		fid = databean.FieldByIndex(1).ValueInt64()
	}
	return
}

func (eh *externaldb) delMessageById(id int64) (err error) {
	_, err = eh.exec(sys.Conf.ExternalDB.Tim_sql_delmessage_byid, id)
	return
}

func (eh *externaldb) saveOfflineMessage(node string, uniqueid int64, stanza []byte, mid int64) (err error) {
	if sys.Conf.ExternalDB.Tim_sql_offlinemsg_save != "" {
		_, err = eh.insert(sys.Conf.ExternalDB.Tim_sql_offlinemsg_save, node, uniqueid, stanza, mid)
	} else if sys.Conf.ExternalDB.Tim_sql_offlinemsg_save_mid != "" && mid > 0 {
		_, err = eh.insert(sys.Conf.ExternalDB.Tim_sql_offlinemsg_save_mid, node, uniqueid, mid)
	} else if sys.Conf.ExternalDB.Tim_sql_offlinemsg_save_nomid != "" && mid == 0 {
		_, err = eh.insert(sys.Conf.ExternalDB.Tim_sql_offlinemsg_save_nomid, node, uniqueid, stanza)
	}
	return
}

func (eh *externaldb) getOfflineMessage(username string, limit int) (obList []*OfflineBean, err error) {
	//var rs [][]any
	//if rs, err = eh.query(sys.Conf.Property.Tim_sql_offlinemsg_get, username, limit); err == nil {
	//	obList = make([]*data.OfflineBean, 0)
	//	for _, bs := range rs {
	//		ob := &data.OfflineBean{Id: data._getInt64(bs[0]), Stanze: data._getBytes(bs[1])}
	//		obList = append(obList, ob)
	//	}
	//}
	databeans := eh.db.ExecuteQueryBeans(sys.Conf.ExternalDB.Tim_sql_offlinemsg_get, username, limit)
	if err = databeans.GetError(); err == nil {
		obList = []*OfflineBean{}
		for _, databean := range databeans.Beans {
			obList = append(obList, &OfflineBean{Id: databean.FieldByIndex(0).ValueInt64(), Stanze: databean.FieldByIndex(1).ValueBytes()})
		}
	}
	return
}

func (eh *externaldb) delOfflineMessage(ids ...int64) (_r int64, err error) {
	_ids := make([]any, 0)
	if sys.Conf.ExternalDB.Tim_sql_offlinemsg_del != "" {
		for _, id := range ids {
			_r, err = eh.exec(sys.Conf.ExternalDB.Tim_sql_offlinemsg_del, id)
		}
	} else if sys.Conf.ExternalDB.Tim_sql_offlinemsg_delin != "" {
		sb := &strings.Builder{}
		sb.WriteString("in(")
		for i := 0; i < len(ids)-1; i++ {
			sb.WriteString("?,")
		}
		sb.WriteString("?)")
		re := regexp.MustCompile(`in\s{0,}\(.*?\)`)
		s := re.ReplaceAll([]byte(sys.Conf.ExternalDB.Tim_sql_offlinemsg_delin), []byte(sb.String()))
		for _, i := range ids {
			_ids = append(_ids, i)
		}
		_r, err = eh.exec(string(s), _ids...)
	}
	return
}

// func (eh *externaldb) existOfflineMessage(node string, uniqueId int64) (exist bool, err error) {
// 	if sys.Conf.Property.Tim_sql_offlinemsg_exist != "" {
// 		if rs, e := eh.query(sys.Conf.Property.Tim_sql_offlinemsg_exist, node, uniqueId); e == nil {
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

func (eh *externaldb) authUser(fromnode, tonode string) (_r bool, err error) {
	if sys.Conf.ExternalDB.Tim_sql_authuser == "" {
		return true, nil
	}
	//var rs [][]any
	//if rs, err = eh.query(sys.Conf.Property.Tim_sql_authuser, fromnode, tonode); err == nil {
	//	for _, bs := range rs {
	//		_r = data._getInt64(bs[0]) > 0
	//		break
	//	}
	//}
	databean := eh.db.ExecuteQueryBean(sys.Conf.ExternalDB.Tim_sql_authuser, fromnode, tonode)
	if err = databean.GetError(); err == nil {
		_r = databean.FieldByIndex(0).ValueInt64() > 0
	}
	return
}

func (eh *externaldb) authGroup(groupnode, usernode string) (_r bool, err error) {
	if sys.Conf.ExternalDB.Tim_sql_authroom == "" {
		return true, nil
	}
	//var rs [][]any
	//if rs, err = eh.query(sys.Conf.Property.Tim_sql_authroom, groupnode, usernode); err == nil {
	//	for _, bs := range rs {
	//		_r = data._getInt64(bs[0]) > 0
	//		break
	//	}
	//}
	databean := eh.db.ExecuteQueryBean(sys.Conf.ExternalDB.Tim_sql_authroom, groupnode, usernode)
	if err = databean.GetError(); err == nil {
		_r = databean.FieldByIndex(0).ValueInt64() > 0
	}
	return
}

func (eh *externaldb) existUser(node string) (_r bool, err error) {
	if sys.Conf.ExternalDB.Tim_sql_existuser == "" {
		return true, nil
	}
	//var rs [][]any
	//if rs, err = eh.query(sys.Conf.Property.Tim_sql_existuser, node); err == nil {
	//	for _, bs := range rs {
	//		_r = data._getInt64(bs[0]) > 0
	//		break
	//	}
	//}
	databean := eh.db.ExecuteQueryBean(sys.Conf.ExternalDB.Tim_sql_existuser, node)
	if err = databean.GetError(); err == nil {
		_r = databean.FieldByIndex(0).ValueInt64() > 0
	}
	return
}

func (eh *externaldb) existGroup(node string) (_r bool, err error) {
	if sys.Conf.ExternalDB.Tim_sql_existroom == "" {
		return true, nil
	}
	//var rs [][]any
	//if rs, err = eh.query(sys.Conf.Property.Tim_sql_existroom, node); err == nil {
	//	for _, bs := range rs {
	//		_r = data._getInt64(bs[0]) > 0
	//		break
	//	}
	//}
	databean := eh.db.ExecuteQueryBean(sys.Conf.ExternalDB.Tim_sql_existroom, node)
	if err = databean.GetError(); err == nil {
		_r = databean.FieldByIndex(0).ValueInt64() > 0
	}
	return
}

/*****************************************************************************************************************/

func (eh *externaldb) roster(username string) (_r []string) {
	if sys.Conf.ExternalDB.Tim_sql_roster == "" {
		return
	}
	//if rs, _ := eh.query(sys.Conf.Property.Tim_sql_roster, username); rs != nil && len(rs) > 0 {
	//	_r = make([]string, 0)
	//	for _, bs := range rs {
	//		_r = append(_r, data._getString(bs[0]))
	//	}
	//}
	databeans := eh.db.ExecuteQueryBeans(sys.Conf.ExternalDB.Tim_sql_roster, username)
	if databeans != nil {
		for _, databean := range databeans.Beans {
			_r = append(_r, databean.FieldByIndex(0).ValueString())
		}
	}
	return
}

func (eh *externaldb) userGroup(username string) (_r []string) {
	if sys.Conf.ExternalDB.Tim_sql_userroom == "" {
		return
	}
	//if rs, _ := eh.query(sys.Conf.Property.Tim_sql_userroom, username); rs != nil && len(rs) > 0 {
	//	_r = make([]string, 0)
	//	for _, bs := range rs {
	//		_r = append(_r, data._getString(bs[0]))
	//	}
	//}
	databeans := eh.db.ExecuteQueryBeans(sys.Conf.ExternalDB.Tim_sql_userroom, username)
	if databeans != nil {
		for _, databean := range databeans.Beans {
			_r = append(_r, databean.FieldByIndex(0).ValueString())
		}
	}
	return
}

func (eh *externaldb) groupRoster(groupname string) (_r []string) {
	if sys.Conf.ExternalDB.Tim_sql_roomroster == "" {
		return
	}
	//if rs, _ := eh.query(sys.Conf.Property.Tim_sql_roomroster, groupname); rs != nil && len(rs) > 0 {
	//	_r = make([]string, 0)
	//	for _, bs := range rs {
	//		_r = append(_r, data._getString(bs[0]))
	//	}
	//}
	databeans := eh.db.ExecuteQueryBeans(sys.Conf.ExternalDB.Tim_sql_roomroster, groupname)
	if databeans != nil {
		for _, databean := range databeans.Beans {
			_r = append(_r, databean.FieldByIndex(0).ValueString())
		}
	}
	return
}

func (eh *externaldb) addroster(fromname, toname string) (ok bool) {
	if sys.Conf.ExternalDB.Tim_sql_roster_add != "" {
		if id, err := eh.exec(sys.Conf.ExternalDB.Tim_sql_roster_add, fromname, toname); err == nil && id > 0 {
			ok = true
		}
	}
	return
}

func (eh *externaldb) rmroster(fromname, toname string) (ok bool) {
	if sys.Conf.ExternalDB.Tim_sql_roster_rm != "" {
		if id, err := eh.exec(sys.Conf.ExternalDB.Tim_sql_roster_rm, fromname, toname); err == nil && id > 0 {
			ok = true
		}
	}
	return
}

func (eh *externaldb) blockroster(fromname, toname string) (ok bool) {
	if sys.Conf.ExternalDB.Tim_sql_roster_block != "" {
		if id, err := eh.exec(sys.Conf.ExternalDB.Tim_sql_roster_block, fromname, toname); err == nil && id > 0 {
			ok = true
		}
	}
	return
}
