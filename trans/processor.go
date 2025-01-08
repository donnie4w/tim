package trans

import (
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
	if len(tm.GetToList()) > 0 {
		if sys.TimMessageProcessor(tm, sys.TRANS_GOAL) == nil && syncId != 0 {
			p.tran.TimAck(syncId)
		}
	}
	return nil
}

func (p *processor) TimPresence(syncId int64, tp *stub.TimPresence) (err error) {
	if len(tp.GetToList()) > 0 {
		if sys.TimPresenceProcessor(tp, sys.TRANS_GOAL) == nil && syncId != 0 {
			p.tran.TimAck(syncId)
		}
	}
	return nil
}

func (p *processor) TimStream(syncId int64, vb *stub.VBean) (err error) {
	if sys.TimSteamProcessor(vb, sys.TRANS_GOAL) == nil && syncId != 0 {
		p.tran.TimAck(syncId)
	}
	return nil
}

func (p *processor) TimCsVBean(syncId int64, vb *stub.CsVrBean) (err error) {
	switch sys.TIMTYPE(vb.GetVbean().GetDtype()) {
	case sys.VROOM_NEW:
	case sys.VROOM_REMOVE:
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
	if syncId != 0 {
		p.tran.TimAck(syncId)
	}
	return nil
}

func (p *processor) TimAck(syncId int64) (err error) {
	p.tw.dataWait.Close(syncId)
	return nil
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
