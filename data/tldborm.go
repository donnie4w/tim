// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim
package data

import (
	"fmt"
	"sort"

	"github.com/donnie4w/tim/sys"
	"github.com/donnie4w/tlcli-go/tlcli"
	"github.com/donnie4w/tlorm-go/orm"
)

func Create[T any]() (err error) {
	return orm.Create[T]()
}

func Insert(a timstruct) (seq int64, err error) {
	return orm.Table[byte](cli(a.tid())).Insert(a)
}

// update data for non nil
// func Update(a timstruct) (err error) {
// 	return orm.Table[byte](cli(a.tid())).Update(a)
// }

// update data for non nil
func Delete[T any](tid uint64, id int64) (err error) {
	return orm.Table[T](cli(tid)).Delete(id)
}

// update data for non zero
func UpdateNonzero(a timstruct) (err error) {
	return orm.Table[byte](cli(a.tid())).UpdateNonzero(a)
}

func SelectIdByIdx[T timstruct](columnName string, columnValue uint64) (id int64, err error) {
	return orm.Table[T](cli(columnValue)).SelectIdByIdx(columnName, columnValue)
}

func SelectById[T timstruct](tid uint64, id int64) (a *T, err error) {
	return orm.Table[T](cli(tid)).SelectById(id)
}

func SelectByIdx[T any](columnName string, columnValue uint64) (a *T, err error) {
	return orm.Table[T](cli(columnValue)).SelectByIdx(columnName, columnValue)
}

func SelectByIdxWithTid[T any](tid uint64, columnName string, columnValue uint64) (a *T, err error) {
	return orm.Table[T](cli(tid)).SelectByIdx(columnName, columnValue)
}

func SelectAllByIdx[T any](columnName string, columnValue uint64) (as []*T, err error) {
	return orm.Table[T](cli(columnValue)).SelectAllByIdx(columnName, columnValue)
}

func SelectAllByIdxWithTid[T any](tid uint64, columnName string, columnValue uint64) (as []*T, err error) {
	return orm.Table[T](cli(tid)).SelectAllByIdx(columnName, columnValue)
}

func SelectByIdxLimit[T any](startId, limit int64, columnName string, columnValue uint64) (as []*T, err error) {
	return orm.Table[T](cli(columnValue)).SelectByIdxLimit(startId, limit, columnName, columnValue)
}

func SelectByIdxDescLimit[T any](id, limit int64, columnName string, columnValue uint64) (as []*T, err error) {
	return orm.Table[T](cli(columnValue)).SelectByIdxDescLimit(columnName, columnValue, id, limit)
}

func SelectIdByIdxSeq[T any](columnName string, columnValue uint64, id int64) (_r int64, err error) {
	return orm.Table[T](cli(columnValue)).SelectIdByIdxSeq(columnName, columnValue, id)
}

var dbsources = make([]*dbsource, 0)

type dbsource struct {
	extent int
	client *tlcli.Client
}

func tlormInit() error {
	if sys.UseTldbExtent() {
		for _, te := range sys.Conf.TldbExtent {
			if cli, e := orm.NewConn(te.Tls, te.Addr, te.Auth); e == nil {
				extent := sys.MB
				if te.ExtentMax > 0 && te.ExtentMax < extent {
					extent = te.ExtentMax
				}
				dbsources = append(dbsources, &dbsource{extent: extent, client: cli})
				initTable(cli)
			} else {
				return fmt.Errorf("%s", "tldb init error:"+e.Error())
			}
		}
		sort.Slice(dbsources, func(i, j int) bool { return dbsources[i].extent < dbsources[j].extent })
	} else {
		if cli, e := orm.NewConn(sys.Conf.Tldb.Tls, sys.Conf.Tldb.Addr, sys.Conf.Tldb.Auth); e == nil {
			dbsources = append(dbsources, &dbsource{extent: sys.MB, client: cli})
			initTable(cli)
		} else {
			return fmt.Errorf("%s", "tldb init error:"+e.Error())
		}
	}
	return nil
}

func initTable(cli *tlcli.Client) {
	orm.Table[timgroup](cli).Create()
	orm.Table[timmessage](cli).Create()
	orm.Table[timoffline](cli).Create()
	orm.Table[timuser](cli).Create()
	orm.Table[timrelate](cli).Create()
	orm.Table[timmucroster](cli).Create()
	orm.Table[timroster](cli).Create()
	orm.Table[timblock](cli).Create()
	orm.Table[timblockroom](cli).Create()
}

func cli(tid uint64) *tlcli.Client {
	idx := 0
	if idx = sort.Search(len(dbsources), func(i int) bool { return dbsources[i].extent >= int(tid%sys.MB) }); idx >= len(dbsources) {
		idx = 0
	}
	return dbsources[idx].client
}