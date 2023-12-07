// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim
//

package data

import (
	"fmt"
	"testing"
)

func TestSqlDB(t *testing.T) {
	fmt.Println(sqlHandle.connect("mysql", "root:123@tcp(127.0.0.1:3306)/timdb"))
	if rs, err := sqlHandle.query("select username,pwd from timuser where `uid` <?", 5); err == nil {
		for _, bs := range rs {
			username := _getString(bs[0])
			pwd := _getString(bs[1])
			fmt.Println(username, " >>>", pwd)
		}
	} else {
		fmt.Println(">>>", err)
	}
}
