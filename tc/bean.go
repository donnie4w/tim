// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim
//

package tc

import (
	"github.com/donnie4w/tim/sys"
)

type SysVar struct {
	StartTime      string
	Time           string
	UUID           int64
	CSNUM          int32
	ALLUUIDS       string
	ADDR           string
	ADMINADDR      string
}

type SysVarView struct {
	Show string
	SYS  *SysVar
	RN   []*sys.RemoteNode
}

type AdminView struct {
	Show       string
	AdminUser  map[string]string
	Init       bool
	ShowCreate bool
}

type Tables struct {
	Name    string
	Columns []string
	Idxs    []string
	Seq     int64
	Sub     int64
}

type TData struct {
	Name    string
	Id      int64
	Columns map[string]string
}

type SelectBean struct {
	Name        string
	Id          string
	ColumnName  string
	ColumnValue string
	StartId     string
	Limit       string
}

type DataView struct {
	Tb      []*Tables
	Tds     []*TData
	ColName map[string][]byte
	Sb      *SelectBean
	Stat    bool
}

/**********************************************************************************/
type SysParam struct {
	DBFILEDIR         string
	MQTLS             bool
	ADMINTLS          bool
	CLITLS            bool
	CLICRT            string
	CLIKEY            string
	MQCRT             string
	MQKEY             string
	ADMINCRT          string
	ADMINKEY          string
	COCURRENT_PUT     int64
	COCURRENT_GET     int64
	DBMode            int
	NAMESPACE         string
	VERSION           string
	BINLOGSIZE        int64
	ADDR              string
	CLIADDR           string
	MQADDR            string
	WEBADMINADDR      string
	CLUSTER_NUM       int
	PWD               string
	PUBLICKEY         string
	PRIVATEKEY        string
	CLUSTER_NUM_FINAL bool
}

type SysParamView struct {
	SYS  *SysParam
	Stat bool
}

/**********************************************************************************/
type AlterTable struct {
	TableName   string
	ID          int64
	Columns     map[string]*FieldInfo
	ColumnValue map[string]string
}

type FieldInfo struct {
	Idx   bool
	Type  string
	Tname string
}
