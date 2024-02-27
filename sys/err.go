// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim
package sys

import (
	"fmt"

	. "github.com/donnie4w/tim/stub"
)

var ERR_HASEXIST = err(4101, "has exist")
var ERR_NOPASS = err(4102, "no pass")
var ERR_EXPIREOP = err(4103, "expire operate")
var ERR_PARAMS = err(4104, "parameter incorrect")
var ERR_AUTH = err(4105, "limited authority")
var ERR_ACCOUNT = err(4106, "account incorrect")
var ERR_INTERFACE = err(4107, "interface incorrect")
var ERR_CANCEL = err(4108, "must not be a cancle object")
var ERR_NOEXIST = err(4109, "must not be a no exist object")
var ERR_BLOCK = err(4110, "blocked object")
var ERR_OVERENTRY = err(4111, "over entry")
var ERR_MODIFYAUTH = err(4112, "modify password failed")
var ERR_FORMAT = err(4113, "format error")
var ERR_BIGDATA = err(4114, "big data error")
var ERR_TOKEN = err(4115, "error token")
var ERR_PING = err(4116, "error ping count")

var ERR_UNDEFINED = err(5101, "undefined error")
var ERR_BLOCKHANDLE = err(5102, "blocking operation")
var ERR_DATABASE = err(5103, "database error")
var ERR_OVERLOAD = err(5104, "heavy server load")
var ERR_OVERHZ = err(5105, "freq out of limit")
var ERR_UUID_REUSE = err(1101, "uuid reuse")
var ERR_OVERTIME = err(1102, "overtime")

type ERROR interface {
	TimError() *TimError
	Error() error
}

type timerror struct {
	code int32
	info string
}

func err(code int32, info string) ERROR {
	return &timerror{code, info}
}

func (this *timerror) TimError() *TimError {
	return &TimError{Code: &this.code, Info: &this.info}
}

func (this *timerror) Error() error {
	return fmt.Errorf("code:%d,info:%s", this.code, this.info)
}
