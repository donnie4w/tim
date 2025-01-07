// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package branch

import "github.com/donnie4w/tim/stub"

func newTid(node string, domain *string) *stub.Tid {
	tid := stub.NewTid()
	tid.Node = node
	tid.Domain = domain
	return tid
}
