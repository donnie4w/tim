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

	"github.com/donnie4w/gofer/util"
	"github.com/donnie4w/tim/sys"
	"github.com/donnie4w/tlcli-go/tlcli"
	"github.com/donnie4w/tlorm-go/orm"
)

func Create[T any]() (err error) {
	return orm.Create[T]()
}

func Insert(a timstruct) (seq int64, err error) {
	return orm.Table[byte](tldbCli(a.Tid())).Insert(a)
}

// update data for non nil
// func Update(a timstruct) (err error) {
// 	return orm.Table[byte](tldbCli(a.Tid())).Update(a)
// }

// Delete data for non nil
func Delete[T any](tid uint64, id int64) (err error) {
	return orm.Table[T](tldbCli(tid)).Delete(id)
}

// UpdateNonzero data for non zero
func UpdateNonzero(a timstruct) (err error) {
	return orm.Table[byte](tldbCli(a.Tid())).UpdateNonzero(a)
}

func SelectIdByIdx[T timstruct](columnName string, columnValue []byte) (id int64, err error) {
	return orm.Table[T](tldbCli(util.FNVHash64(columnValue))).SelectIdByIdx(columnName, columnValue)
}

func SelectById[T timstruct](tid []byte, id int64) (a *T, err error) {
	return orm.Table[T](tldbCli(util.FNVHash64(tid))).SelectById(id)
}

func SelectByIdx[T any](columnName string, columnValue []byte) (a *T, err error) {
	return orm.Table[T](tldbCli(util.FNVHash64(columnValue))).SelectByIdx(columnName, columnValue)
}

func SelectByIdxWithInt[T any](columnName string, columnValue uint64) (a *T, err error) {
	return orm.Table[T](tldbCli(columnValue)).SelectByIdx(columnName, columnValue)
}

func SelectByIdxWithTid[T any](tid uint64, columnName string, columnValue []byte) (a *T, err error) {
	return orm.Table[T](tldbCli(tid)).SelectByIdx(columnName, columnValue)
}

func SelectAllByIdx[T any](columnName string, columnValue uint64) (as []*T, err error) {
	return orm.Table[T](tldbCli(columnValue)).SelectAllByIdx(columnName, columnValue)
}

func SelectAllByIdxWithTid[T any](tid uint64, columnName string, columnValue []byte) (as []*T, err error) {
	return orm.Table[T](tldbCli(tid)).SelectAllByIdx(columnName, columnValue)
}

func SelectByIdxLimit[T any](startId, limit int64, columnName string, columnValue uint64) (as []*T, err error) {
	return orm.Table[T](tldbCli(columnValue)).SelectByIdxLimit(startId, limit, columnName, columnValue)
}

func SelectByIdxDescLimit[T any](id, limit int64, columnName string, columnValue []byte) (as []*T, err error) {
	return orm.Table[T](tldbCli(util.FNVHash64(columnValue))).SelectByIdxDescLimit(columnName, columnValue, id, limit)
}

func SelectIdByIdxSeq[T any](columnName string, columnValue []byte, id int64) (_r int64, err error) {
	return orm.Table[T](tldbCli(util.FNVHash64(columnValue))).SelectIdByIdxSeq(columnName, columnValue, id)
}

var tldbsources = make([]*tldbsource, 0)

type tldbsource struct {
	extent int
	client *tlcli.Client
}

func tldbInit() error {
	if len(sys.Conf.TldbExtent) > 0 {
		for _, te := range sys.Conf.TldbExtent {
			if cli, e := orm.NewConn(te.Tls, te.Addr, te.Auth); e == nil {
				extent := sys.MB
				if te.ExtentMax > 0 && te.ExtentMax < extent {
					extent = te.ExtentMax
				}
				tldbsources = append(tldbsources, &tldbsource{extent: extent, client: cli})
				initTable(cli)
			} else {
				return fmt.Errorf("%s", "tldb init error:"+e.Error())
			}
		}
		sort.Slice(tldbsources, func(i, j int) bool { return tldbsources[i].extent < tldbsources[j].extent })
	} else if sys.Conf.Tldb != nil {
		if cli, e := orm.NewConn(sys.Conf.Tldb.Tls, sys.Conf.Tldb.Addr, sys.Conf.Tldb.Auth); e == nil {
			tldbsources = append(tldbsources, &tldbsource{extent: sys.MB, client: cli})
			initTable(cli)
		} else {
			return fmt.Errorf("%s", "tldb init error:"+e.Error())
		}
	} else {
		panic("tldb init error")
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

func tldbCli(tid uint64) *tlcli.Client {
	idx := 0
	if idx = sort.Search(len(tldbsources), func(i int) bool { return tldbsources[i].extent >= int(tid%sys.MB) }); idx >= len(tldbsources) {
		idx = 0
	}
	return tldbsources[idx].client
}
