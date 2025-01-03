// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package amr

import "github.com/donnie4w/gofer/util"

func GetVnode(vnode string) int64 {
	if bs := Get([]byte(vnode)); len(bs) > 0 {
		return util.BytesToInt64(bs)
	}
	return 0
}
func PutVnode(vnode string, uuid int64) {
	if vnode != "" && uuid != 0 {
		Put([]byte(vnode), util.Int64ToBytes(uuid), 365*86400)
	}
}
func DelVnode(vnode string) {
	if vnode != "" {
		Remove([]byte(vnode))
	}
}
