// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package service

import (
	"github.com/donnie4w/gofer/buffer"
	"github.com/donnie4w/tim/amr"
	"github.com/donnie4w/tim/data"
	"github.com/donnie4w/tim/errs"
	"github.com/donnie4w/tim/stub"
	"github.com/donnie4w/tim/sys"
	"github.com/donnie4w/tim/util"
	"github.com/donnie4w/tim/vgate"
)

func timMessageProcessor(tm *stub.TimMessage, transType int8) (err errs.ERROR) {
	defer util.Recover()
	sys.Stat.TxDo()
	defer sys.Stat.TxDone()
	if tm.ToTid == nil && len(tm.ToList) == 0 {
		return errs.ERR_PARAMS
	}
	switch transType {
	case sys.TRANS_SOURCE:
		if tolist := tm.ToList; len(tolist) > 0 {
			tm.ToList = nil
			uuidmap := make(map[int64][]string)
			for _, node := range tolist {
				if uuids := amr.GetAccount(node); len(uuids) > 0 {
					for _, uuid := range uuids {
						if ll, b := uuidmap[uuid]; b {
							uuidmap[uuid] = append(ll, node)
						} else {
							uuidmap[uuid] = []string{node}
						}
					}
				} else {
					tm.ToTid = &stub.Tid{Node: node}
					timMessageProcessor(tm, sys.TRANS_GOAL)
				}
			}
			for uuid, list := range uuidmap {
				tm.ToTid, tm.ToList = nil, list
				if uuid == sys.UUID {
					timMessageProcessor(tm, sys.TRANS_GOAL)
				} else {
					if !sys.CsMessageService(uuid, tm, true) {
						timMessageProcessor(tm, sys.TRANS_GOAL)
					}
				}
			}
		} else if tm.ToTid != nil {
			if uuids := amr.GetAccount(tm.ToTid.GetNode()); len(uuids) > 0 {
				for _, uuid := range uuids {
					if uuid != sys.UUID {
						sys.CsMessageService(uuid, tm, true)
					} else {
						timMessageProcessor(tm, sys.TRANS_GOAL)
					}
				}
			} else {
				timMessageProcessor(tm, sys.TRANS_GOAL)
			}
		}
	case sys.TRANS_GOAL:
		if len(tm.ToList) > 0 {
			list := tm.ToList
			tm.ToList = nil
			for _, node := range list {
				timMessage4goal(node, tm)
			}
		} else if tm.ToTid != nil {
			timMessage4goal(tm.ToTid.GetNode(), tm)
		}
	}
	return
}

func timMessage4goal(node string, tm *stub.TimMessage) {
	ol := false
	if wsware.hasUser(node) {
		if tm.OdType != sys.ORDER_STREAM && tm.OdType != sys.ORDER_BIGSTRING && tm.OdType != sys.ORDER_BIGBINARY {
			if wsware.SendNodeWithAck(node, tm, sys.TIMMESSAGE) {
				ol = true
			}
		} else if tm.OdType == sys.ORDER_BIGSTRING {
			wsware.SendBigData(node, []byte(tm.FromTid.Node+sys.SEP_STR+tm.GetDataString()), sys.TIMBIGSTRING)
		} else if tm.OdType == sys.ORDER_BIGBINARY {
			buf := buffer.NewBufferWithCapacity(len(tm.FromTid.Node) + 1 + len(tm.GetDataBinary()))
			buf.WriteString(tm.FromTid.Node)
			buf.WriteByte(sys.SEP_BIN)
			buf.Write(tm.GetDataBinary())
			wsware.SendBigData(node, buf.Bytes(), sys.TIMBIGBINARY)
		} else {
			wsware.SendNode(node, tm, sys.TIMMESSAGE)
		}
	}
	if !ol && tm.OdType != sys.ORDER_STREAM && tm.OdType != sys.ORDER_BIGSTRING && tm.OdType != sys.ORDER_BIGBINARY {
		data.Service.SaveOfflineMessage(node, tm)
	}
}

func timPresenceProcessor(tp *stub.TimPresence, transType int8) (err errs.ERROR) {
	defer util.Recover()
	sys.Stat.TxDo()
	defer sys.Stat.TxDone()
	if tp.ToTid != nil || tp.ToList != nil {
		switch transType {
		case sys.TRANS_SOURCE:
			ms := make(map[int64][]string, 0)
			if len(tp.ToList) > 0 {
				for _, node := range tp.ToList {
					uuids := amr.GetAccount(node)
					for _, uuid := range uuids {
						if us, b := ms[uuid]; b {
							ms[uuid] = append(us, node)
						} else {
							ms[uuid] = []string{node}
						}
					}
				}
				for uuid, nodes := range ms {
					tp.ToList = nodes
					if uuid == sys.UUID {
						timPresenceProcessor(tp, sys.TRANS_GOAL)
					} else {
						sys.CsPresenceService(uuid, tp, false)
					}
				}
			} else if tp.ToTid != nil {
				if uuids := amr.GetAccount(tp.ToTid.GetNode()); len(uuids) > 0 {
					for _, uuid := range uuids {
						if uuid == sys.UUID {
							timPresenceProcessor(tp, sys.TRANS_GOAL)
						} else {
							sys.CsPresenceService(uuid, tp, false)
						}
					}
				}
			}
		case sys.TRANS_GOAL:
			if len(tp.ToList) > 0 {
				list := tp.ToList
				tp.ToList = nil
				if tp.ToTid == nil {
					tp.ToTid = stub.NewTid()
				}
				for _, node := range list {
					tp.ToTid.Node = node
					wsware.SendNode(node, tp, sys.TIMPRESENCE)
				}
			} else if tp.ToTid != nil {
				wsware.SendNode(tp.ToTid.Node, tp, sys.TIMPRESENCE)
			}
		}
	} else {
		err = errs.ERR_PARAMS
	}
	return
}

func timStreamProcessor(vb *stub.VBean, transType int8) (ok bool, err errs.ERROR) {
	defer util.Recover()
	sys.Stat.TxDo()
	defer sys.Stat.TxDone()
	switch transType {
	case sys.TRANS_SOURCE:
		uuid := amr.GetVnode(vb.GetVnode())
		if uuid != 0 && uuid != sys.UUID {
			switch sys.TIMTYPE(vb.GetRtype()) {
			case sys.VROOM_SUB, sys.VROOM_UNSUB:
				sys.CsVBeanService(uuid, vb, false)
			}
		} else if uuid == sys.UUID {
			if sys.TIMTYPE(vb.GetRtype()) == sys.VROOM_MESSAGE {
				sendVStream(vb)
				vgate.VGate.GetSubUUID(vb.GetVnode()).Range(func(k int64, _ int8) bool {
					if k != sys.UUID {
						sys.CsVBeanService(k, vb, false)
					}
					return true
				})
			}
		}
	case sys.TRANS_GOAL:
		switch sys.TIMTYPE(vb.GetRtype()) {
		case sys.VROOM_MESSAGE:
			return sendVStream(vb)
		}
	}
	return
}

func sendVStream(vb *stub.VBean) (b bool, err errs.ERROR) {
	defer util.Recover()
	sys.Stat.TxDo()
	defer sys.Stat.TxDone()
	if vr, ok := vgate.VGate.GetVroom(vb.Vnode); ok {
		vr.Updatetime()
		vr.SubMap().Range(func(wsId int64, ib int8) bool {
			var ok bool
			if ib == 0 {
				ts := &stub.TimStream{ID: vb.GetStreamId(), VNode: vb.Vnode, FromNode: vb.GetRnode(), Body: vb.Body, Dtype: vb.Dtype}
				ok = wsware.SendWs(wsId, ts, sys.TIMSTREAM)
			} else if ib == 1 {
				buf := buffer.NewBufferWithCapacity(len(vb.Vnode) + 1 + len(vb.Body))
				buf.WriteString(vb.Vnode)
				buf.WriteByte(sys.SEP_BIN)
				buf.Write(vb.Body)
				ok = wsware.SendBigDataByWs(wsId, buf.Bytes(), sys.TIMBIGBINARYSTREAM)
			}
			if !ok {
				if _, b := sys.WsById(wsId); !b {
					go vgate.VGate.UnSub(vb.Vnode, wsId)
				}
			}
			if !b {
				b = true
			}
			return true
		})
	} else {
		err = errs.ERR_NOEXIST
	}
	return
}
