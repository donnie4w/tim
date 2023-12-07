// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim
package service

import (
	"github.com/donnie4w/tim/data"
	. "github.com/donnie4w/tim/stub"
	"github.com/donnie4w/tim/sys"
	"github.com/donnie4w/tim/util"
	"github.com/donnie4w/tim/vgate"
)

func timMessage(tm *TimMessage, transType int8) (err sys.ERROR) {
	defer util.Recover()
	sys.Stat.TxDo()
	defer sys.Stat.TxDone()
	if toTid := tm.ToTid; toTid != nil {
		switch transType {
		case sys.TRANS_CONSISHASH:
			sys.CsMessage(tm, sys.TRANS_CONSISHASH)
		case sys.TRANS_SOURCE:
			sys.CsMessage(tm, sys.TRANS_SOURCE)
		case sys.TRANS_STAFF:
			sys.CsMessage(tm, sys.TRANS_STAFF)
		case sys.TRANS_GOAL:
			ol := false
			node := tm.ToTid.Node
			if wsware.hasUser(node) {
				if tm.OdType != sys.ORDER_STREAM {
					if wsware.SendNodeWithAck(node, tm, sys.TIMMESSAGE) {
						ol = true
					}
				} else {
					wsware.SendNode(node, tm, sys.TIMMESSAGE)
				}
			}
			if !ol && tm.OdType != sys.ORDER_STREAM {
				data.Handler.SaveOfflineMessage(tm)
			}
		}
	} else {
		err = sys.ERR_PARAMS
	}
	return
}

func timPresence(tp *TimPresence, transType int8) (err sys.ERROR) {
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
		err = sys.ERR_PARAMS
	}
	return
}

func timStream(vb *VBean) (err sys.ERROR) {
	defer util.Recover()
	sys.Stat.TxDo()
	defer sys.Stat.TxDone()
	if vr, ok := vgate.VGate.GetVroom(vb.Vnode); ok {
		vr.Updatetime()
		vr.Nodes().Range(func(wsId int64, _ int8) bool {
			ts := &TimStream{ID: *vb.StreamId, VNode: vb.Vnode, FromNode: *vb.Rnode, Body: vb.Body, Dtype: vb.Dtype}
			util.GoPoolTx2.Go(func() { wsware.SendWs(wsId, ts, sys.TIMSTREAM) })
			return true
		})
	} else {
		err = sys.ERR_NOEXIST
	}
	return
}
