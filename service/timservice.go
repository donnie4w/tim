// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package service

import (
	. "github.com/donnie4w/tim/stub"
)

/**
func timMessageService(tm *TimMessage) (err errs.ERROR) {
	defer util.Recover()
	sys.Stat.TxDo()
	defer sys.Stat.TxDone()

	ol := false
	node := tm.ToTid.Node
	if wsware.hasUser(node) {
		if tm.OdType != sys.ORDER_STREAM && tm.OdType != sys.ORDER_BIGSTRING && tm.OdType != sys.ORDER_BIGBINARY {
			if wsware.SendNodeWithAck(node, tm, sys.TIMMESSAGE) {
				ol = true
			}
		} else if tm.OdType == sys.ORDER_BIGSTRING {
			wsware.SendBigData(tm.ToTid.Node, []byte(tm.FromTid.Node+sys.SEP_STR+tm.GetDataString()), sys.TIMBIGSTRING)
		} else if tm.OdType == sys.ORDER_BIGBINARY {
			buf := NewBufferWithCapacity(len(tm.FromTid.Node) + 1 + len(tm.GetDataBinary()))
			buf.WriteString(tm.FromTid.Node)
			buf.WriteByte(sys.SEP_BIN)
			buf.Write(tm.GetDataBinary())
			wsware.SendBigData(tm.ToTid.Node, buf.Bytes(), sys.TIMBIGBINARY)
		} else {
			wsware.SendNode(node, tm, sys.TIMMESSAGE)
		}
	}
	if !ol && tm.OdType != sys.ORDER_STREAM && tm.OdType != sys.ORDER_BIGSTRING && tm.OdType != sys.ORDER_BIGBINARY {
		data.Service.SaveOfflineMessage(tm)
	}

	return
}

func timPresenceService(tp *TimPresence, transType int8) (err errs.ERROR) {
	defer util.Recover()
	sys.Stat.TxDo()
	defer sys.Stat.TxDone()
	if tp.ToTid != nil || tp.ToList != nil {
		switch transType {
		case sys.TRANS_CONSISHASH:
			sys.CsPresence(tp, sys.TRANS_CONSISHASH)
		case sys.TRANS_SOURCE:
			sys.CsPresence(tp, sys.TRANS_SOURCE)
		case sys.TRANS_STAFF:
			sys.CsPresence(tp, sys.TRANS_STAFF)
		case sys.TRANS_GOAL:
			wsware.SendNode(tp.ToTid.Node, tp, sys.TIMPRESENCE)
		}
	} else {
		err = errs.ERR_PARAMS
	}
	return
}

func timStreamService(vb *VBean) (err errs.ERROR) {
	defer util.Recover()
	sys.Stat.TxDo()
	defer sys.Stat.TxDone()
	if vr, ok := vgate.VGate.GetVroom(vb.Vnode); ok {
		vr.Updatetime()
		vr.SubMap().Range(func(wsId int64, ib int8) bool {
			if ib == 0 {
				ts := &TimStream{ID: *vb.StreamId, VNode: vb.Vnode, FromNode: *vb.Rnode, Body: vb.Body, Dtype: vb.Dtype}
				util.GoPoolTx2.Go(func() { wsware.SendWs(wsId, ts, sys.TIMSTREAM) })
			} else if ib == 1 {
				buf := NewBufferWithCapacity(len(vb.Vnode) + 1 + len(vb.Body))
				buf.WriteString(vb.Vnode)
				buf.WriteByte(sys.SEP_BIN)
				buf.Write(vb.Body)
				util.GoPoolTx2.Go(func() { wsware.SendBigDataByWs(wsId, buf.Bytes(), sys.TIMBIGBINARYSTREAM) })
			}
			return true
		})
	} else {
		err = errs.ERR_NOEXIST
	}
	return
}
***/
