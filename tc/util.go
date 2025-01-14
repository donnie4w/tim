// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package tc

import (
	"github.com/donnie4w/tlnet"
)

func reqjson(hc *tlnet.HttpContext) bool {
	return hc.Request().Header.Get("content-type") == "application/json"
}

func reqform(hc *tlnet.HttpContext) bool {
	return hc.Request().Header.Get("content-type") == "application/x-www-form-urlencoded"
}
