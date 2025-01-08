// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package trans

import (
	"github.com/donnie4w/gofer/buffer"
	"github.com/donnie4w/gofer/util"
	"github.com/donnie4w/gothrift/thrift"
	"github.com/donnie4w/tim/errs"
	"github.com/donnie4w/tim/stub"

	"github.com/donnie4w/tsf"
	"sync"
	"sync/atomic"
)

type trans struct {
	ts      tsf.TsfSocket
	mux     sync.RWMutex
	noAck   atomic.Int32
	isValid bool
}

func newTrans(ts tsf.TsfSocket) csNet {
	return &trans{ts: ts, isValid: true}
}

func (t *trans) Id() int64 {
	return t.ts.ID()
}

func (t *trans) addNoAck() (r int32) {
	if r = t.noAck.Add(1); r > 3 {
		t.Close()
	}
	return
}

func (t *trans) IsValid() bool { return t.isValid }

func (t *trans) Close() (err error) {
	defer util.Recover(&err)
	if t != nil {
		t.mux.Lock()
		defer t.mux.Unlock()
		if t.isValid {
			t.isValid = false
			return t.ts.Close()
		}
	}
	return nil
}

func (t *trans) write(bs []byte) (err error) {
	if !t.isValid {
		return errs.ERR_CONNECT.Error()
	}
	t.mux.RLock()
	defer t.mux.RUnlock()
	if _, err = t.ts.WriteWithMerge(bs); err != nil {
		t.Close()
	}
	return
}

func (t *trans) writeMessage(tp byte, syncId int64, ts thrift.TStruct) (err error) {
	var bs []byte
	if ts != nil {
		bs = util.TEncode(ts)
	}
	buf := buffer.NewBufferWithCapacity(9 + len(bs))
	buf.WriteByte(tp)
	buf.Write(util.Int64ToBytes(syncId))
	if len(bs) > 0 {
		buf.Write(bs)
	}
	return t.write(buf.Bytes())
}

func (t *trans) writeBytes(tp byte, syncId int64, body []byte) (err error) {
	if syncId != 0 {
		buf := buffer.NewBufferWithCapacity(9 + len(body))
		buf.WriteByte(tp)
		buf.Write(util.Int64ToBytes(syncId))
		buf.Write(body)
		return t.write(buf.Bytes())
	} else {
		buf := buffer.NewBufferWithCapacity(1 + len(body))
		buf.WriteByte(tp)
		buf.Write(body)
		return t.write(buf.Bytes())
	}
}

func (t *trans) TimMessage(syncId int64, tm *stub.TimMessage) (err error) {
	return t.writeMessage(TIMMESSAGE, syncId, tm)
}

func (t *trans) TimPresence(syncId int64, tp *stub.TimPresence) (err error) {
	return t.writeMessage(TIMPRESENCE, syncId, tp)
}

func (t *trans) TimStream(syncId int64, vb *stub.VBean) (err error) {
	return t.writeMessage(TIMSTREAM, syncId, vb)
}

func (t *trans) TimCsVBean(syncId int64, vb *stub.CsVrBean) (err error) {
	return t.writeMessage(TIMCSVBEAN, syncId, vb)
}

func (t *trans) TimAck(syncId int64) (err error) {
	if syncId != 0 {
		t.writeMessage(TIMACK, syncId, nil)
	}
	return nil
}
