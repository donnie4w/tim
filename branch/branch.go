// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package branch

import (
	goutil "github.com/donnie4w/gofer/util"
	"github.com/donnie4w/tim/data"
	"github.com/donnie4w/tim/errs"
	"github.com/donnie4w/tim/stub"
	"github.com/donnie4w/tim/sys"
	"github.com/donnie4w/tim/util"
	"github.com/donnie4w/tim/vgate"
	"sort"
	"time"
)

func Roster(fromnode string, domain *string) (_r []string, err errs.ERROR) {
	if !util.CheckNodes(fromnode) {
		return nil, errs.ERR_PARAMS
	}
	_r = data.Service.Roster(fromnode)
	return
}

func Addroster(fromnode string, domain *string, toNode string, msg *string) (err errs.ERROR) {
	if fromnode == toNode {
		return errs.ERR_PERM_DENIED
	}
	if !util.CheckNodes(fromnode, toNode) {
		return errs.ERR_PARAMS
	}
	if _, err = data.Service.Addroster(fromnode, toNode, domain); err == nil {
		id, _t, bt := goutil.UUID64(), time.Now().UnixNano(), int32(sys.BUSINESS_ADDROSTER)
		tm := &stub.TimMessage{ID: &id, MsType: sys.SOURCE_USER, OdType: sys.ORDER_BUSINESS, ToTid: &stub.Tid{Node: toNode}, FromTid: newTid(fromnode, domain), DataString: msg, Timestamp: &_t}
		//if status == 0x10|0x01 {
		//	bt := int32(sys.BUSINESS_FRIEND)
		//	tm.BnType = &bt
		//	sys.SendWs(ws.Id, tm, sys.TIMMESSAGE)
		//} else {
		//	bt := int32(sys.BUSINESS_ADDROSTER)
		tm.BnType = &bt
		//}/
		sys.TimMessageProcessor(tm, sys.TRANS_SOURCE)
	}
	return
}

func Rmroster(fromnode string, domain *string, toNode string) errs.ERROR {
	if fromnode == toNode {
		return errs.ERR_PERM_DENIED
	}
	if !util.CheckNodes(fromnode, toNode) {
		return errs.ERR_PARAMS
	}
	if ms, ok := data.Service.Rmroster(fromnode, toNode, domain); !ok {
		return errs.ERR_PARAMS
	} else if ms {
		id, t, bt := goutil.UUID64(), time.Now().UnixNano(), int32(sys.BUSINESS_REMOVEROSTER)
		tm := &stub.TimMessage{ID: &id, MsType: sys.SOURCE_USER, BnType: &bt, OdType: sys.ORDER_BUSINESS, FromTid: newTid(fromnode, domain), ToTid: &stub.Tid{Node: toNode}, Timestamp: &t}
		sys.TimMessageProcessor(tm, sys.TRANS_SOURCE)
	}
	return nil
}

func Blockroster(fromnode string, domain *string, toNode string) errs.ERROR {
	if fromnode == toNode {
		return errs.ERR_PERM_DENIED
	}
	if !util.CheckNodes(fromnode, toNode) {
		return errs.ERR_PARAMS
	}
	if ms, ok := data.Service.Blockroster(fromnode, toNode, domain); !ok {
		return errs.ERR_PARAMS
	} else if ms {
		id, _t, bt := goutil.UUID64(), time.Now().UnixNano(), int32(sys.BUSINESS_BLOCKROSTER)
		tm := &stub.TimMessage{ID: &id, MsType: sys.SOURCE_USER, OdType: sys.ORDER_BUSINESS, ToTid: &stub.Tid{Node: toNode}, FromTid: newTid(fromnode, domain), BnType: &bt, Timestamp: &_t}
		sys.TimMessageProcessor(tm, sys.TRANS_SOURCE)
	}
	return nil
}

func PullUserMessage(fromnode string, domain *string, toNode string, mid int64, limit int64) []*stub.TimMessage {
	if !util.CheckNodes(fromnode, toNode) {
		return nil
	}
	if oblist, _ := data.Service.GetMessage(fromnode, domain, 1, toNode, mid, limit); len(oblist) > 0 {
		if *oblist[0].Mid == mid {
			oblist = oblist[1:]
		}
		sort.Slice(oblist, func(i, j int) bool { return *oblist[i].Mid > *oblist[j].Mid })
		return oblist
	}
	return nil
}

func PullRoomMessage(fromnode string, domain *string, roomNode string, mid int64, limit int64) []*stub.TimMessage {
	if !util.CheckNodes(fromnode, roomNode) {
		return nil
	}
	if oblist, _ := data.Service.GetMessage(fromnode, domain, 2, roomNode, mid, limit); len(oblist) > 0 {
		if *oblist[0].Mid == mid {
			oblist = oblist[1:]
		}
		sort.Slice(oblist, func(i, j int) bool { return *oblist[i].Mid > *oblist[j].Mid })
		return oblist
	}
	return nil
}

func OfflineMsg(fromnode string, domain *string, limit int) []*stub.TimMessage {
	if !util.CheckNode(fromnode) {
		return nil
	}
	if oblist, _ := data.Service.GetOfflineMessage(fromnode, limit); len(oblist) > 0 {
		tmList := make([]*stub.TimMessage, 0)
		isOff := true
		ids := make([]int64, 0)
		for _, ob := range oblist {
			ids = append(ids, ob.Id)
			if ob.Stanze != nil {
				if tm, err := goutil.TDecode(ob.Stanze, stub.NewTimMessage()); err == nil {
					tm.IsOffline = &isOff
					tmList = append(tmList, tm)
					if ob.Mid > 0 {
						tm.Mid = &ob.Mid
					}
				}
			}
		}
		sort.Slice(tmList, func(i, j int) bool {
			return *tmList[i].Timestamp < *tmList[j].Timestamp
		})
		return tmList
	}
	return nil
}

func DelOfflineMsg(fromNode string, ids []int64) (int64, error) {
	if uuid := util.NodeToUUID(fromNode); uuid > 0 {
		return data.Service.DelOfflineMessage(util.NodeToUUID(fromNode), ids...)
	}
	return 0, errs.ERR_ACCOUNT.Error()
}

func UserRoom(fromnode string, domain *string) []string {
	if !util.CheckNode(fromnode) {
		return nil
	}
	return data.Service.UserGroup(fromnode, domain)
}

func RoomUsers(roomNode string, domain *string) []string {
	if !util.CheckNode(roomNode) {
		return nil
	}
	return data.Service.GroupRoster(roomNode)
}

func AddRoom(fromnode string, domain *string, roomNode string, msg string) (err errs.ERROR) {
	if !util.CheckNodes(fromnode, roomNode) {
		return errs.ERR_PARAMS
	}
	var gtype int8
	if gtype, err = data.Service.GroupGtype(roomNode, domain); err == nil {
		switch sys.TIMTYPE(gtype) {
		case sys.GROUP_OPEN:
			err = data.Service.Addgroup(roomNode, fromnode, domain)
		case sys.GROUP_PRIVATE:
			if err = data.Service.Addgroup(roomNode, fromnode, domain); err != nil {
				return
			}
			var ms []string
			if ms, err = data.Service.GroupManagers(roomNode, domain); err == nil {
				for _, u := range ms {
					id, t, bt := goutil.UUID64(), time.Now().UnixNano(), int32(sys.BUSINESS_ADDROOM)
					tm := &stub.TimMessage{ID: &id, FromTid: newTid(fromnode, domain), ToTid: &stub.Tid{Node: u}, BnType: &bt, RoomTid: &stub.Tid{Node: roomNode}, MsType: sys.SOURCE_USER, OdType: sys.ORDER_BUSINESS, DataString: &msg, Timestamp: &t}
					sys.TimMessageProcessor(tm, sys.TRANS_SOURCE)
				}
			}
		}
	}
	return
}

func PullInRoom(fromnode string, domain *string, roomNode string, toNode string) (err errs.ERROR) {
	if fromnode == toNode {
		return errs.ERR_PERM_DENIED
	}
	if !util.CheckNodes(fromnode, roomNode, toNode) {
		return errs.ERR_PARAMS
	}
	var isreq bool
	if isreq, err = data.Service.Pullgroup(roomNode, fromnode, toNode, domain); err == nil {
		id, _t := goutil.UUID64(), time.Now().UnixNano()
		tm := &stub.TimMessage{ID: &id, FromTid: newTid(fromnode, domain), ToTid: &stub.Tid{Node: toNode}, RoomTid: &stub.Tid{Node: roomNode}, MsType: sys.SOURCE_USER, OdType: sys.ORDER_BUSINESS, Timestamp: &_t}
		if isreq {
			bt := int32(sys.BUSINESS_PASSROOM)
			tm.BnType = &bt
		} else {
			bt := int32(sys.BUSINESS_PULLROOM)
			tm.BnType = &bt
		}
		sys.TimMessageProcessor(tm, sys.TRANS_SOURCE)
	}
	return
}

func RejectRoom(fromnode string, domain *string, roomNode string, toNode string, msg string) (err errs.ERROR) {
	if !util.CheckNodes(fromnode, roomNode, toNode) {
		return errs.ERR_PARAMS
	}
	if err = data.Service.Nopassgroup(roomNode, fromnode, toNode, domain); err == nil {
		id, _t, bt := goutil.UUID64(), time.Now().UnixNano(), int32(sys.BUSINESS_NOPASSROOM)
		tm := &stub.TimMessage{ID: &id, FromTid: newTid(fromnode, domain), ToTid: &stub.Tid{Node: toNode}, RoomTid: &stub.Tid{Node: roomNode}, MsType: sys.SOURCE_USER, OdType: sys.ORDER_BUSINESS, Timestamp: &_t}
		tm.BnType = &bt
		sys.TimMessageProcessor(tm, sys.TRANS_SOURCE)
	}
	return
}

func KickRoom(fromnode string, domain *string, roomNode string, toNode string) (err errs.ERROR) {
	if fromnode == toNode {
		return errs.ERR_PERM_DENIED
	}
	if !util.CheckNodes(fromnode, roomNode, toNode) {
		return errs.ERR_PARAMS
	}
	if err = data.Service.Kickgroup(roomNode, fromnode, toNode, domain); err == nil {
		id, _t, bt := goutil.UUID64(), time.Now().UnixNano(), int32(sys.BUSINESS_KICKROOM)
		tm := &stub.TimMessage{ID: &id, BnType: &bt, FromTid: newTid(fromnode, domain), ToTid: &stub.Tid{Node: toNode}, RoomTid: &stub.Tid{Node: roomNode}, MsType: sys.SOURCE_USER, OdType: sys.ORDER_BUSINESS, Timestamp: &_t}
		sys.TimMessageProcessor(tm, sys.TRANS_SOURCE)
	}
	return
}

func LeaveRoom(fromnode string, domain *string, roomNode string) errs.ERROR {
	if !util.CheckNodes(fromnode, roomNode) {
		return errs.ERR_PARAMS
	}
	return data.Service.Leavegroup(roomNode, fromnode, domain)
}

func CancelRoom(fromnode string, domain *string, roomNode string) errs.ERROR {
	if !util.CheckNodes(fromnode, roomNode) {
		return errs.ERR_PARAMS
	}
	return data.Service.Cancelgroup(roomNode, fromnode, domain)
}

func BlockRoom(fromnode string, domain *string, roomNode string) errs.ERROR {
	if !util.CheckNodes(fromnode, roomNode) {
		return errs.ERR_PARAMS
	}
	return data.Service.Blockgroup(roomNode, fromnode, domain)
}

func BlockRoomMember(fromnode string, domain *string, roomNode string, toNode string) errs.ERROR {
	if fromnode == toNode {
		return errs.ERR_PERM_DENIED
	}
	if !util.CheckNodes(fromnode, roomNode, toNode) {
		return errs.ERR_PARAMS
	}
	return data.Service.Blockgroupmember(roomNode, fromnode, toNode, domain)
}

func BlockRosterList(fromnode string, domain *string) (_r []string, err errs.ERROR) {
	if !util.CheckNodes(fromnode) {
		err = errs.ERR_PARAMS
		return
	}
	_r = data.Service.Blockrosterlist(fromnode)
	return
}

func BlockRoomList(fromnode string, domain *string) (_r []string, err errs.ERROR) {
	if !util.CheckNodes(fromnode) {
		err = errs.ERR_PARAMS
		return
	}
	_r = data.Service.Blockroomlist(fromnode)
	return
}

func BlockRoomMemberlist(fromnode string, domain *string, roomNode string) (_r []string, err errs.ERROR) {
	if !util.CheckNodes(fromnode) {
		err = errs.ERR_PARAMS
		return
	}
	_r = data.Service.Blockroommemberlist(fromnode, roomNode)
	return
}

func VirtualroomRegister(fromnode, domain string) (r *stub.TimAck, err errs.ERROR) {
	if !util.CheckNodes(fromnode) {
		err = errs.ERR_PARAMS
		return
	}
	vnode := vgate.VGate.NewVRoom(fromnode)
	//if sys.CsVBean(&stub.VBean{Rtype: int8(sys.VROOM_NEW), Vnode: vnode, FoundNode: &fromnode}) {
	t := int64(sys.VROOM_NEW)
	r = &stub.TimAck{Ok: true, TimType: int8(sys.TIMVROOM), N: &vnode, T: &t}
	//}
	return
}

func VirtualroomRemove(fromnode, domain string, vNode string) (r *stub.TimAck, err errs.ERROR) {
	if !util.CheckNodes(fromnode, vNode) {
		err = errs.ERR_PARAMS
		return
	}
	if !vgate.VGate.Remove(fromnode, vNode) {
		err = errs.ERR_NOEXIST
		return
	}
	//if sys.CsVBean(&stub.VBean{Rtype: int8(sys.VROOM_REMOVE), Vnode: vNode, FoundNode: &fromnode}) {
	t := int64(sys.VROOM_REMOVE)
	r = &stub.TimAck{Ok: true, TimType: int8(sys.TIMVROOM), N: &vNode, T: &t}
	//}
	return
}

func VirtualroomAddAuth(fromnode, domain string, vNode string, toNode string) (r *stub.TimAck, err errs.ERROR) {
	if !util.CheckNodes(fromnode, vNode, toNode) {
		err = errs.ERR_PARAMS
		return
	}
	//if !vgate.VGate.AddAuth(vNode, fromnode, toNode) {
	//	err = errs.ERR_PERM_DENIED
	//	return
	//}
	//if sys.CsVBean(&stub.VBean{Rtype: int8(sys.VROOM_ADDAUTH), Vnode: vNode, FoundNode: &fromnode, Rnode: &toNode}) {
	//	t := int64(sys.VROOM_ADDAUTH)
	//	r = &stub.TimAck{Ok: true, TimType: int8(sys.TIMVROOM), N: &toNode, T: &t}
	//}
	return
}

func VirtualroomDelAuth(fromnode, domain string, vNode string, toNode string) (r *stub.TimAck, err errs.ERROR) {
	if !util.CheckNodes(fromnode, vNode, toNode) {
		err = errs.ERR_PARAMS
		return
	}
	//vgate.VGate.DelAuth(vNode, fromnode, toNode)
	//if sys.CsVBean(&stub.VBean{Rtype: int8(sys.VROOM_DELAUTH), Vnode: vNode, FoundNode: &fromnode, Rnode: &toNode}) {
	//	t := int64(sys.VROOM_DELAUTH)
	//	r = &stub.TimAck{Ok: true, TimType: int8(sys.TIMVROOM), N: &toNode, T: &t}
	//}
	return
}

func VirtualroomSub(wsId int64, fromnode string, domain string, vNode string, subType int8) (r *stub.TimAck, err errs.ERROR) {
	if !util.CheckNodes(fromnode, vNode) {
		err = errs.ERR_PARAMS
		return
	}
	if _, b := sys.WsById(wsId); !b {
		err = errs.ERR_ACCOUNT
		return
	}
	var success bool
	if subType == 1 {
		success = vgate.VGate.SubBinary(vNode, sys.UUID, wsId)
	} else {
		success = vgate.VGate.Sub(vNode, sys.UUID, wsId)
	}
	if success && sys.TimSteamProcessor(&stub.VBean{Rtype: int8(sys.VROOM_SUB), Vnode: vNode, Rnode: &fromnode}, sys.TRANS_SOURCE) == nil {
		//if sys.CsVBean(&stub.VBean{Rtype: int8(sys.VROOM_SUB), Vnode: vNode, Rnode: &fromnode}) {
		t := int64(sys.VROOM_SUB)
		r = &stub.TimAck{Ok: true, TimType: int8(sys.TIMVROOM), N: &fromnode, T: &t}
	} else {
		err = errs.ERR_UNDEFINED
	}

	return
}

func VirtualroomUnSub(wsId int64, fromnode string, domain string, vNode string) (r *stub.TimAck, err errs.ERROR) {
	if !util.CheckNodes(fromnode, vNode) {
		err = errs.ERR_PARAMS
		return
	}
	if _, b := sys.WsById(wsId); !b {
		err = errs.ERR_ACCOUNT
		return
	}
	if k, b := vgate.VGate.UnSub(vNode, wsId); b && k == 0 {
		t := int64(sys.VROOM_UNSUB)
		//sys.CsVBean(&stub.VBean{Rtype: int8(sys.VROOM_UNSUB), Vnode: vNode})
		sys.TimSteamProcessor(&stub.VBean{Rtype: int8(sys.VROOM_UNSUB), Vnode: vNode}, sys.TRANS_SOURCE)
		r = &stub.TimAck{Ok: true, TimType: int8(sys.TIMVROOM), N: &fromnode, T: &t}
	}
	return
}
