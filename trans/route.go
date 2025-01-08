/*
 * Copyright (c) 2024 donnie4w <donnie4w@gmail.com>. All rights reserved.
 * Original source: https://github.com/donnie4w/raftx
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package trans

import (
	"fmt"
	"github.com/donnie4w/gofer/util"
	"github.com/donnie4w/tim/stub"
)

func router(bs []byte, cn csNet) (err error) {
	defer util.Recover(&err)
	switch bs[0] {
	case TIMMESSAGE:
		if len(bs) <= 9 {
			return cn.TimMessage(util.BytesToInt64(bs[1:9]), nil)
		}
		if p, err := util.TDecode[*stub.TimMessage](bs[9:], stub.NewTimMessage()); err == nil {
			return cn.TimMessage(util.BytesToInt64(bs[1:9]), p)
		}
	case TIMPRESENCE:
		if len(bs) <= 9 {
			return cn.TimPresence(util.BytesToInt64(bs[1:9]), nil)
		}
		if p, err := util.TDecode[*stub.TimPresence](bs[9:], stub.NewTimPresence()); err == nil {
			return cn.TimPresence(util.BytesToInt64(bs[1:9]), p)
		}
	case TIMSTREAM:
		if len(bs) <= 9 {
			return cn.TimStream(util.BytesToInt64(bs[1:9]), nil)
		}
		if p, err := util.TDecode[*stub.VBean](bs[9:], stub.NewVBean()); err == nil {
			return cn.TimStream(util.BytesToInt64(bs[1:9]), p)
		}
	case TIMCSVBEAN:
		if len(bs) <= 9 {
			return cn.TimCsVBean(util.BytesToInt64(bs[1:9]), nil)
		}
		if p, err := util.TDecode[*stub.CsVrBean](bs[9:], stub.NewCsVrBean()); err == nil {
			return cn.TimCsVBean(util.BytesToInt64(bs[1:9]), p)
		}
	case TIMACK:
		if len(bs) <= 9 {
			return cn.TimAck(util.BytesToInt64(bs[1:9]))
		}
	default:
		err = fmt.Errorf("unknow command")
	}
	return
}
