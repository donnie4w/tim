/**
 * donnie4w@gmail.com  tim server
 */
package impl

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"
	"time"

	"tim/FW"
	"tim/cluster"
	"tim/clusterRoute"
	. "tim/common"
	. "tim/connect"
	"tim/daoService"
	. "tim/protocol"
	"tim/route"
	"tim/tfClient"
	"tim/utils"

	"github.com/apache/thrift/lib/go/thrift"
	"github.com/donnie4w/go-logger/logger"
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
func (this *TimImpl) TimStream(ctx context.Context, param *TimParam) (err error) {
	if param != nil {
		if param.GetInterflow() == "1" {
			this.Tu.Interflow = 1
		}
		if param.GetTLS() == "1" {
			this.Tu.TLS = 1
		}
		this.Tu.Version = param.GetVersion()
	}
	return
}
func (this *TimImpl) TimStarttls(ctx context.Context) (err error) {
	panic("")
}

// Parameters:
//  - Tid
//  - Pwd
func (this *TimImpl) TimLogin(ctx context.Context, tid *Tid, pwd string) (err error) {
	defer func() {
		if err := recover(); err != nil {
			logger.Warn("Login error", err)
			//			logger.Error(string(debug.Stack()))
		}
	}()
	isAuth := false
	if this.Tu.Fw == FW.AUTH {
		ack := NewTimAckBean()
		status200, typelogin := "200", "login"
		ack.AckStatus, ack.AckType = &status200, &typelogin
		this.Tu.SendAckBean(ack)
		return
	}
	if CF.MustAuth == 0 {
		isAuth = true
	} else {
		user_auth_url := CF.GetKV("user_auth_url", "")
		if len(user_auth_url) > 9 {
			isAuth = httpAuth(tid, pwd, user_auth_url)
		} else {
			b := daoService.Auth(tid, pwd)
			if b {
				isAuth = true
				logger.Debug("login is success:", tid.GetName())
			}
		}
	}

	if isAuth {
		ack := NewTimAckBean()
		this.Tu.UserTid = tid
		this.Tu.Fw = FW.AUTH
		this.Tu.Auth(tid)
		if cluster.IsCluster() {
			loginname, _ := GetLoginName(tid)
			cluster.SetLoginnameToCluster(loginname)
		}
		status200, typelogin := "200", "login"
		ack.AckStatus, ack.AckType = &status200, &typelogin
		this.Tu.SendAckBean(ack)
		_TimPresence(this, OnlinePBean(this.Tu.UserTid), false)
		go route.RouteOffLineMBean(this.Tu)
	} else {
		ack := NewTimAckBean()
		status400, typeType := "400", "login"
		ack.AckStatus, ack.AckType = &status400, &typeType
		this.Tu.SendAckBean(ack)
		panic(fmt.Sprint("loginname or pwd is error:", tid.GetName(), " | ", pwd))
	}
	return
}

// Parameters:
//  - Ab
func (this *TimImpl) TimAck(cx context.Context, ab *TimAckBean) (err error) {
	if this.Tu.Fw != FW.AUTH {
		panic(fmt.Sprint("not auth:", this.Tu.Fw))
	}
	this.Tu.OverLimit = 3
	go func() {
		defer func() {
			if err := recover(); err != nil {
			}
		}()
		if CF.ConfirmAck == 1 && ab != nil && ab.GetID() == this.Tu.LastSyncThreadId {
			timer := time.NewTicker(5 * time.Second)
			select {
			case <-timer.C:
				logger.Error("ack msg threadid over time", ab)
			case this.Tu.Sendflag <- ab.GetID():
			}
		}
	}()
	return
}

// Parameters:
//  - Pbean
func (this *TimImpl) TimPresence(ctx context.Context, pbean *TimPBean) (err error) {
	if CF.Presence != 1 {
		return
	}
	if this.Tu.Fw != FW.AUTH {
		panic("not auth")
	}
	//	logger.Debug("pbean", pbean)
	if this.Tu.UserType == 0 {
		pbean.FromTid = this.Tu.UserTid
		_type := pbean.GetType()
		switch _type {
		case "groupchat":
			pbean.LeaguerTid = this.Tu.UserTid
		}
	}
	return _TimPresence(this, pbean, true)
}

func _TimPresence(this *TimImpl, pbean *TimPBean, isAck bool) (err error) {
	defer func() {
		if er := recover(); er != nil {
			logger.Error(string(debug.Stack()))
			err = errors.New(fmt.Sprint(er))
		}
	}()
	if CF.Presence != 1 {
		return
	}
	if pbean.GetThreadId() == "" {
		pbean.ThreadId = utils.TimeMills()
	}
	//isTotidExist := daoService.IsTidExist(pbean.GetToTid())
	_type := pbean.GetType()
	switch _type {
	case "groupchat":
		pbean.FromTid = pbean.ToTid
		//default:
		//	pbean.ToTid.Domain = pbean.FromTid.Domain //只能发送到相同domain的用户
	}

	mustRoute := false
	if cluster.IsCluster() && this.Tu.UserType == 0 {
		er := clusterRoute.ClusterRoutePBean(pbean)
		if er != nil {
			mustRoute = true
		}
	} else {
		mustRoute = true
	}
	if mustRoute {
		if pbean.GetToTid() == nil {
			route.RoutePBean(pbean)
		} else {
			route.RouteSinglePBean(pbean)
		}
		if isAck {
			ack := NewTimAckBean()
			id := pbean.ThreadId
			ack.ID = &id
			status200, typemessage := "200", "presence"
			ack.AckStatus, ack.AckType = &status200, &typemessage
			this.Tu.SendAckBean(ack)
		}
	}
	return
}

// Parameters:
//  - Mbean
func (this *TimImpl) TimMessage(ctx context.Context, mbean *TimMBean) (err error) {
	if this.Tu.Fw != FW.AUTH {
		panic("not auth")
	}
	//	logger.Debug("TimMessage=====>", mbean)
	if this.Tu.UserType == 0 {
		mbean.FromTid = this.Tu.UserTid
		//		isTotidExist := daoService.IsTidExist(mbean.GetToTid())
		_type := mbean.GetType()
		switch _type {
		case "groupchat":
			b := daoService.AuthMucmember(mbean.GetToTid(), this.Tu.UserTid)
			if !b {
				panic("auth room failed")
			}
		}
	}
	//	if isTotidExist {
	//		id, _, _ := route.RouteMBean(mbean, false, false)
	//		ack := NewTimAckBean()
	//		status200, typemessage := "200", "message"
	//		ack.AckStatus, ack.AckType = &status200, &typemessage
	//		ack.ExtraMap = make(map[string]string, 0)
	//		ack.ExtraMap["mid"] = fmt.Sprint(id)
	//		this.Tu.SendAckBean(ack)
	//	}
	return _TimMessage(this, mbean)
}

func _TimMessage(this *TimImpl, mbean *TimMBean) (err error) {
	if mbean.GetThreadId() == "" {
		mbean.ThreadId = utils.TimeMills()
	}
	isTotidExist := daoService.IsTidExist(mbean.GetToTid())
	_type := mbean.GetType()
	switch _type {
	case "groupchat":
		mbean.FromTid = mbean.ToTid
		mbean.LeaguerTid = this.Tu.UserTid
		mbean.FromTid.Domain = this.Tu.UserTid.Domain
		mbean.ToTid = nil
		timestamp := utils.TimeMills()
		mbean.Timestamp = &timestamp
	default:
		mbean.ToTid.Domain = mbean.FromTid.Domain //只能发送到相同domain的用户
		timestamp := utils.TimeMills()
		mbean.Timestamp = &timestamp
	}
	if isTotidExist && mbean.GetToTid() != nil {
		mustRoute := true
		if cluster.IsCluster() {
			clusterBean := clusterRoute.OtherClusterUserBean(mbean.GetToTid())
			if this.Tu.UserType == 0 && clusterBean != nil {
				er := clusterRoute.ClusterRouteMBean(mbean, clusterBean)
				if er != nil {
					mustRoute = true
				} else {
					mustRoute = false
				}
			} else {
				mustRoute = true
			}
		}

		if mustRoute {
			id, er, _ := route.RouteMBean(mbean, false, true)
			ack := NewTimAckBean()
			thid := mbean.ThreadId
			ack.ID = &thid
			if er == nil {
				status, typemessage := TIM_SC_SUCCESS, "message"
				ack.AckStatus, ack.AckType = &status, &typemessage
				ack.ExtraMap = make(map[string]string, 0)
				ack.ExtraMap["mid"] = fmt.Sprint(id)
			} else {
				status, typemessage := TIM_SC_FAILED, "message"
				ack.AckStatus, ack.AckType = &status, &typemessage
			}
			this.Tu.SendAckBean(ack)
		}
	}
	return
}

// Parameters:
//  - ThreadId
func (this *TimImpl) TimPing(ctx context.Context, threadId string) (err error) {
	if this.Tu.Fw != FW.AUTH {
		panic("not auth")
	}
	//	logger.Debug("ping>>>>>", threadId)
	ab := NewTimAckBean()
	ab.ID = &threadId
	acktype, ackstatus := "ping", "200"
	ab.AckType, ab.AckStatus = &acktype, &ackstatus
	this.Tu.SendAckBean(ab)
	return
}

// Parameters:
//  - E
func (this *TimImpl) TimError(cx context.Context, e *TimError) (err error) {
	panic("TimError")
}
func (this *TimImpl) TimLogout(cx context.Context) (err error) {
	panic("TimLogout")
}

// Parameters:
//  - Tid
//  - Pwd
func (this *TimImpl) TimRegist(ctx context.Context, tid *Tid, pwd string) (err error) {
	panic("error TimRegist")
}

// Parameters:
//  - Tid
//  - Pwd
func (this *TimImpl) TimRemoteUserAuth(ctx context.Context, tid *Tid, pwd string, auth *TimAuth) (r *TimRemoteUserBean, err error) {
	panic("error TimRemoteUserAuth")
}

// Parameters:
//  - Tid
func (this *TimImpl) TimRemoteUserGet(ctx context.Context, tid *Tid, auth *TimAuth) (r *TimRemoteUserBean, err error) {
	panic("error TimRemoteUserGet")
}

// Parameters:
//  - Tid
//  - Ub
func (this *TimImpl) TimRemoteUserEdit(ctx context.Context, tid *Tid, ub *TimUserBean, auth *TimAuth) (r *TimRemoteUserBean, err error) {
	panic("error TimRemoteUserEdit")
}

// Parameters:
//  - Pbean
func (this *TimImpl) TimResponsePresence(ctx context.Context, pbean *TimPBean, auth *TimAuth) (r *TimResponseBean, err error) {
	panic("TimResponsePresence")
}

// Parameters:
//  - Mbean
func (this *TimImpl) TimResponseMessage(ctx context.Context, mbean *TimMBean, auth *TimAuth) (r *TimResponseBean, err error) {
	r = NewTimResponseBean()
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

func (this *TimImpl) TimMessageIq(ctx context.Context, timMsgIq *TimMessageIq, iqType string) (err error) {
	//	logger.Debug("TimMessageIq:", timMsgIq, " ", iqType)
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
				if this.Tu.Interflow > 0 {
					mbeanlist := NewTimMBeanList()
					mbeanlist.ThreadId = utils.TimeMills()
					mbeanlist.TimMBeanList = mbeans
					this.Tu.Client.TimMessageList(context.Background(), mbeanlist)
				} else {
					if mbeans != nil {
						for _, mbean := range mbeans {
							er := this.Tu.Client.TimMessageResult_(context.Background(), mbean)
							if er != nil {
								break
							}
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
func (this *TimImpl) TimMessageResult_(ctx context.Context, mbean *TimMBean) (err error) {
	logger.Debug("TimMessageResult_:", mbean)
	panic("error TimMessageResult_")
}

func (this *TimImpl) TimRoser(ctx context.Context, roster *TimRoster) (err error) {
	logger.Debug("TimRoser:", roster)
	panic("error TimRoser")
}

func (this *TimImpl) TimResponseMessageIq(ctx context.Context, timMsgIq *TimMessageIq, iqType string, auth *TimAuth) (r *TimMBeanList, err error) {
	//	logger.Debug("TimResponseMessageIq:", timMsgIq, iqType, auth)
	user_auth_url := CF.GetKV("user_auth_url", "")
	isAuth := false
	tid := NewTid()
	tid.Domain, tid.Name = auth.Domain, auth.GetUsername()
	pwd := auth.GetPwd()
	if user_auth_url != "" {
		isAuth = httpAuth(tid, pwd, user_auth_url)
	} else {
		b := daoService.Auth(tid, pwd)
		if b {
			isAuth = true
		}
	}
	if !isAuth {
		return
	}
	switch iqType {
	case "offline":
		r = NewTimMBeanList()
		mbeans := daoService.LoadOfflineMBean(tid)
		r.ThreadId = utils.TimeMills()
		r.TimMBeanList = mbeans
		mids := make([]interface{}, 0)
		for _, mbean := range mbeans {
			mids = append(mids, mbean.GetMid())
		}
		daoService.DelOfflineMBeanList(mids...)
		daoService.UpdateOffMessageList(mbeans, 1)
	case "get":
	}
	return
}

func httpAuth(tid *Tid, pwd, user_auth_url string) (isAuth bool) {
	var r *TimRemoteUserBean
	tfClient.HttpClient(func(client *ITimClient) (er error) {
		defer func() {
			if err := recover(); err != nil {
				er = errors.New(fmt.Sprint(err))
				logger.Error(string(debug.Stack()))
			}
		}()
		r, er = client.TimRemoteUserAuth(context.Background(), tid, pwd, nil)
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
		return er
	}, user_auth_url)
	return
}

func (this *TimImpl) TimMessageList(ctx context.Context, mbeanList *TimMBeanList) (err error) {
	logger.Debug("TimMessageList:", mbeanList)
	panic("error TimMessageList")
}

// Parameters:
//  - PbeanList
func (this *TimImpl) TimPresenceList(ctx context.Context, pbeanList *TimPBeanList) (err error) {
	logger.Debug("TimPresenceList:", pbeanList)
	panic("error TimPresenceList")
}

func (this *TimImpl) TimResponsePresenceList(ctx context.Context, pbeanList *TimPBeanList, auth *TimAuth) (r *TimResponseBean, err error) {
	logger.Debug("TimResponsePresenceList:", pbeanList)
	panic("error TimResponsePresenceList")
}

// Parameters:
//  - MbeanList
//  - Auth
func (this *TimImpl) TimResponseMessageList(ctx context.Context, mbeanList *TimMBeanList, auth *TimAuth) (r *TimResponseBean, err error) {
	logger.Debug("TimResponseMessageList:", mbeanList)
	panic("error TimResponseMessageList")
}

func (this *TimImpl) TimProperty(ctx context.Context, tpb *TimPropertyBean) (err error) {
	if this.Tu.Fw != FW.AUTH {
		panic("not auth")
	}
	//	logger.Debug("TimProperty:", tpb)
	interflow := tpb.GetInterflow()
	tls := tpb.GetTLS()
	if interflow == "1" {
		this.Tu.Interflow = 1
	}
	if tls == "1" {
		this.Tu.TLS = 1
	}
	return
}
