// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim
//

package stub

import (
	"testing"

	. "github.com/donnie4w/gofer/util"
	"github.com/donnie4w/simplelog/logging"
)

func Test_json(t *testing.T) {
	n, p := "helloworld", "123abc"
	ta := &TimAuth{Name: &n, Pwd: &p}
	bs := JsonEncode(ta)
	logging.Debug(len(bs), ">>", string(bs))
	logging.Debug(JsonDecode[*TimAuth](bs))
}

func Test_tcode(t *testing.T) {
	n, p := "helloworld", "123abc"
	ta := &TimAuth{Name: &n, Pwd: &p}
	bs := TEncode(ta)
	logging.Debug(len(bs), ">>", string(bs))
	logging.Debug(TDecode(bs, &TimAuth{}))
}

func Test_tcode2(t *testing.T) {
	s:="123"
	ta := &TimMessage{DataString: &s}
	bs := TEncode(ta)
	logging.Debug(len(bs), ">>", string(bs))
	logging.Debug(TDecode(bs, &TimMessage{}))
}