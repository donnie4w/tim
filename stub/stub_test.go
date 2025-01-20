// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim
//

package stub

import (
	"github.com/donnie4w/gofer/util"
	"testing"
)

func Test_json(t *testing.T) {
	n, p := "helloworld", "123abc"
	ta := &TimAuth{Name: &n, Pwd: &p}
	bs := util.JsonEncode(ta)
	t.Log(len(bs), ">>", string(bs))
	t.Log(util.JsonDecode[*TimAuth](bs))
}

func Test_tcode(t *testing.T) {
	n, p := "helloworld", "123abc"
	ta := &TimAuth{Name: &n, Pwd: &p}
	bs := util.TEncode(ta)
	t.Log(len(bs), ">>", string(bs))
	t.Log(util.TDecode(bs, &TimAuth{}))
}

func Test_tcode2(t *testing.T) {
	s := "123"
	ta := &TimMessage{DataString: &s}
	bs := util.TEncode(ta)
	t.Log(len(bs), ">>", string(bs))
	t.Log(util.TDecode(bs, &TimMessage{}))
}
