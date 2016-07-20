/**
 * donnie4w@gmail.com  tim server
 */
package impl

import (
	"fmt"
	"runtime/debug"
	"time"

	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/donnie4w/go-logger/logger"
	"tim.FW"
	. "tim.common"
	. "tim.connect"
	"tim.daoService"
	. "tim.protocol"
	"tim.route"
	"tim.tfClient"
	"tim.utils"
)

type TimImpl struct {
	Ip     string
	Port   int
	Pub    string //发布id
	Tu     *TimUser
	Client thrift.TTransport
}

// Parameters:
//  - Param
func (this *TimImpl) TimStream(param *TimParam) (err error) {
	panic("")
	return
}
func (this *TimImpl) TimStarttls() (err error) {
	panic("")
	return
}

// Parameters:
//  - Tid
//  - Pwd
func (this *TimImpl) TimLogin(tid *Tid, pwd string) (err error) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("Login error", err)
			logger.Error(string(debug.Stack()))
		}
	}()
	logger.Debug("Login:", tid, pwd)
	user_auth_url := ConfBean.GetKV("user_auth_url", "")
	var er error
	isAuth := false
	if this.Tu.Fw == FW.AUTH {
		ack := NewTimAckBean()
		status200, typelogin := "200", "login"
		ack.AckStatus, ack.AckType = &status200, &typelogin
		this.Tu.SendAckBean(ack)
		return
	}
	if user_auth_url != "" {
		var r *TimRemoteUserBean
		tfClient.HttpClient(func(client *ITimClient) {
			r, er = client.TimRemoteUserAuth(tid, pwd)
			if er == nil && r != nil {
				logger.Debug(r)
				if r.ExtraMap != nil {
					if password, ok := r.ExtraMap["password"]; ok {
						if pwd == password {
							isAuth = true
						}
					}
					if extraAuth, ok := r.ExtraMap["extraAuth"]; ok {
						if pwd == extraAuth {
							isAuth = true
						}
					}
				}
			}
		})
	} else {
		b := daoService.Auth(tid.GetDomain(), tid.GetName(), pwd)
		if b {
			isAuth = true
			logger.Debug("login is success:", tid.GetName(), "/", pwd)
		}
	}
	if isAuth {
		ack := NewTimAckBean()
		this.Tu.UserTid = tid
		this.Tu.Fw = FW.AUTH
		this.Tu.Auth(tid)
		status200, typelogin := "200", "login"
		ack.AckStatus, ack.AckType = &status200, &typelogin
		this.Tu.SendAckBean(ack)
		go route.RouteOffLineMBean(this.Tu)
	} else {
		ack := NewTimAckBean()
		status400, typelogin := "400", "login"
		ack.AckStatus, ack.AckType = &status400, &typelogin
		this.Tu.SendAckBean(ack)
		panic("loginname or pwd is error")
	}
	return
}

// Parameters:
//  - Ab
func (this *TimImpl) TimAck(ab *TimAckBean) (err error) {
	if this.Tu.Fw != FW.AUTH {
		panic("not auth")
	}
	logger.Debug("ack", ab)
	if ConfBean.ConfirmAck == 1 && ab != nil && ab.GetID() == this.Tu.LastSyncThreadId {
		timer := time.NewTicker(10 * time.Second)
		select {
		case <-timer.C:
			panic("ack msg threadid over time")
		case this.Tu.Sendflag <- ab.GetID():
		}
	}
	this.Tu.OverLimit = 3
	return
}

// Parameters:
//  - Pbean
func (this *TimImpl) TimPresence(pbean *TimPBean) (err error) {
	if ConfBean.Presence != 1 {
		return
	}
	if this.Tu.Fw != FW.AUTH {
		panic("not auth")
	}
	logger.Debug("pbean", pbean)
	pbean.FromTid = this.Tu.UserTid
	isTotidExist := daoService.IsTidExist(pbean.GetToTid())
	if isTotidExist {
		go route.RoutePBean(pbean)
	}
	return
}

// Parameters:
//  - Mbean
func (this *TimImpl) TimMessage(mbean *TimMBean) (err error) {
	if this.Tu.Fw != FW.AUTH {
		panic("not auth")
	}
	logger.Debug("mbean", mbean)
	mbean.FromTid = this.Tu.UserTid
	isTotidExist := daoService.IsTidExist(mbean.GetToTid())
	_type := mbean.GetType()
	switch _type {
	case "groupchat":
		mbean.LeaguerTid = this.Tu.UserTid
		timestamp := utils.TimeMills()
		mbean.Timestamp = &timestamp
		mbean.FromTid = mbean.ToTid
	default:
		mbean.ToTid.Domain = mbean.FromTid.Domain //只能发送到相同domain的用户
		timestamp := utils.TimeMills()
		mbean.Timestamp = &timestamp
	}
	if isTotidExist {
		id, _, _ := route.RouteMBean(mbean, false, false)
		ack := NewTimAckBean()
		status200, typemessage := "200", "message"
		ack.AckStatus, ack.AckType = &status200, &typemessage
		ack.ExtraMap = make(map[string]string, 0)
		ack.ExtraMap["mid"] = fmt.Sprint(id)
		this.Tu.SendAckBean(ack)
	}
	return
}

// Parameters:
//  - ThreadId
func (this *TimImpl) TimPing(threadId string) (err error) {
	if this.Tu.Fw != FW.AUTH {
		panic("not auth")
	}
	//logger.Debug("ping>>>>>", threadId)
	ab := NewTimAckBean()
	acktype, ackstatus := "ping", "200"
	ab.AckType, ab.AckStatus = &acktype, &ackstatus
	this.Tu.SendAckBean(ab)
	return
}

// Parameters:
//  - E
func (this *TimImpl) TimError(e *TimError) (err error) {
	panic("timimpl")
	return
}
func (this *TimImpl) TimLogout() (err error) {
	panic("timimpl")
	return
}

// Parameters:
//  - Tid
//  - Pwd
func (this *TimImpl) TimRegist(tid *Tid, pwd string) (err error) {
	panic("timimpl")
	return
}

// Parameters:
//  - Tid
//  - Pwd
func (this *TimImpl) TimRemoteUserAuth(tid *Tid, pwd string) (r *TimRemoteUserBean, err error) {
	panic("error process")
	return
}

// Parameters:
//  - Tid
func (this *TimImpl) TimRemoteUserGet(tid *Tid) (r *TimRemoteUserBean, err error) {
	panic("error process")
	return
}

// Parameters:
//  - Tid
//  - Ub
func (this *TimImpl) TimRemoteUserEdit(tid *Tid, ub *TimUserBean) (r *TimRemoteUserBean, err error) {
	panic("TimRemoteUserEdit")
	return
}

// Parameters:
//  - Pbean
func (this *TimImpl) TimResponsePresence(pbean *TimPBean) (r *TimResponseBean, err error) {
	panic("TimResponsePresence")
	return
}

// Parameters:
//  - Mbean
func (this *TimImpl) TimResponseMessage(mbean *TimMBean) (r *TimResponseBean, err error) {
	r = NewTimResponseBean()
	if timresponseauth, ok := mbean.ExtraMap["timresponseauth"]; ok {
		logger.Debug("timresponseauth:", timresponseauth)
	}
	fromDomain := mbean.GetFromTid().GetDomain()
	toDomain := mbean.GetToTid().GetDomain()
	if fromDomain == toDomain {
		if !daoService.CheckDomain(fromDomain) {
			logger.Error("domain check fail:", fromDomain)
			return
		}
	} else {
		logger.Error("fromDomain != toDomain", fromDomain, " ", toDomain)
		return
	}
	isTotidExist := daoService.IsTidExist(mbean.GetToTid())
	mbean.ToTid.Domain = mbean.FromTid.Domain //只能发送到相同domain的用户
	timestamp := utils.TimeMills()
	mbean.Timestamp = &timestamp
	if isTotidExist {
		isSinglePush := false
		if mbean.ExtraMap != nil {
			if tim_pushType, ok := mbean.ExtraMap["tim_pushType"]; ok {
				if tim_pushType == "single" {
					isSinglePush = true
				}
				delete(mbean.ExtraMap, "tim_pushType")
			}
		}

		id, er, offline := route.RouteMBean(mbean, isSinglePush, false)
		if er == nil {
			r.ExtraMap = make(map[string]string, 0)
			r.ExtraMap["mid"] = fmt.Sprint(id)
			r.ExtraMap["timestamp"] = timestamp
			if offline {
				r.ExtraMap["offline"] = "1"
			} else {
				r.ExtraMap["offline"] = "0"
			}
		}
	}
	return
}

func (this *TimImpl) TimMessageIq(timMsgIq *TimMessageIq, iqType string) (err error) {
	logger.Debug("TimMessageIq:", timMsgIq, " ", iqType)
	switch iqType {
	case "get":
		fidname := this.Tu.UserTid.GetName()
		domain := this.Tu.UserTid.Domain
		tidnames := timMsgIq.Tidlist
		limitcount := timMsgIq.TimPage.LimitCount
		fromstamp := timMsgIq.TimPage.FromTimeStamp
		tostamp := timMsgIq.TimPage.ToTimeStamp
		if tidnames != nil {
			for _, tidname := range tidnames {
				mbeans := daoService.LoadMBean(fidname, tidname, *domain, fromstamp, tostamp, *limitcount)
				if mbeans != nil {
					for _, mbean := range mbeans {
						er := this.Tu.Client.TimMessageResult_(mbean)
						if er != nil {
							break
						}
					}
				}
			}
		}
	case "del":
		fidname := this.Tu.UserTid.GetName()
		domain := this.Tu.UserTid.Domain
		tidnames := timMsgIq.Tidlist
		mids := timMsgIq.Midlist
		if len(tidnames) == 1 && len(mids) == 1 {
			daoService.DelMBean(fidname, tidnames[0], *domain, mids[0])
		}
	case "delAll":
		fidname := this.Tu.UserTid.GetName()
		domain := this.Tu.UserTid.Domain
		tidnames := timMsgIq.Tidlist
		if len(tidnames) == 1 {
			daoService.DelAllMBean(fidname, tidnames[0], *domain)
		}
	default:
		panic("error iqType")
	}
	return
}

// Parameters:
//  - Mbean
func (this *TimImpl) TimMessageResult_(mbean *TimMBean) (err error) {
	logger.Debug("TimMessageResult_:", mbean)

	return
}

func (this *TimImpl) TimRoser(roster *TimRoster) (err error) {
	logger.Debug("TimRoser:", roster)
	return
}
