package trans

import (
	"github.com/donnie4w/tim/cache"
	"github.com/donnie4w/tim/stub"
	"github.com/donnie4w/tim/sys"
	"github.com/donnie4w/tim/vgate"
)

type processor struct {
	tran   csNet
	isServ bool
	tw     *TransWare
}

func (p *processor) Close() (err error) {
	return nil
}

func (p *processor) addNoAck() int32 {
	return 0
}

func (p *processor) TimMessage(syncId int64, tm *stub.TimMessage) (err error) {
	if len(tm.GetToList()) > 0 || tm.GetToTid() != nil {
		go sys.TimMessageProcessor(tm, sys.TRANS_GOAL)
	}
	if syncId != 0 {
		p.tran.TimAck(syncId)
	}
	return nil
}

func (p *processor) TimPresence(syncId int64, tp *stub.TimPresence) (err error) {
	if len(tp.GetToList()) > 0 || tp.GetToTid() != nil {
		if tp.GetOffline() && tp.FromTid != nil {
			cache.AccountCache.Del(tp.FromTid.GetNode())
		}
		go sys.TimPresenceProcessor(tp, sys.TRANS_GOAL)
	}
	if syncId != 0 {
		p.tran.TimAck(syncId)
	}
	return nil
}

func (p *processor) TimStream(syncId int64, vb *stub.VBean) (err error) {
	go sys.TimSteamProcessor(vb, sys.TRANS_GOAL)
	if syncId != 0 {
		p.tran.TimAck(syncId)
	}
	return nil
}

func (p *processor) TimCsVBean(syncId int64, vb *stub.CsVrBean) (err error) {
	go func() {
		switch sys.TIMTYPE(vb.GetVbean().GetRtype()) {
		case sys.VROOM_SUB:
			if vb.GetVbean().GetRnode() != "" && vb.GetSrcuuid() != 0 && vb.GetSrcuuid() != sys.UUID {
				vgate.VGate.Sub(vb.GetVbean().GetVnode(), vb.GetSrcuuid(), 0)
			}
		case sys.VROOM_UNSUB:
			if vb.GetSrcuuid() != 0 && vb.GetSrcuuid() != sys.UUID {
				vgate.VGate.UnSubWithUUID(vb.GetVbean().GetVnode(), vb.GetSrcuuid())
			}
		case sys.VROOM_MESSAGE:
			sys.TimSteamProcessor(vb.GetVbean(), sys.TRANS_GOAL)
		}
	}()
	if syncId != 0 {
		p.tran.TimAck(syncId)
	}
	return nil
}

func (p *processor) TimAck(syncId int64) (err error) {
	p.tw.dataWait.Close(syncId)
	return nil
}

func (p *processor) TimCsDevice(syncId int64, cd *stub.CsDevice) (err error) {
	bs := sys.DeviceTypeList(cd.GetNode())
	cd.TypeList = bs
	p.tran.TimCsDeviceAck(syncId, cd)
	return
}

func (p *processor) TimCsDeviceAck(syncId int64, cd *stub.CsDevice) (err error) {
	p.tw.dataWait.CloseAndPut(syncId, cd)
	return
}

func (p *processor) Id() int64 {
	return p.tran.Id()
}

func (p *processor) IsValid() bool {
	return p.tran.IsValid()
}

func newProcessor(tw *TransWare, c csNet, serv bool) csNet {
	return &processor{tw: tw, tran: c, isServ: serv}
}
