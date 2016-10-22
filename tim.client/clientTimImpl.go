/**
 * donnie4w@gmail.com  tim server
 */
package client

import (
	"errors"
	"fmt"
	//	"time"
	//	"runtime/debug"

	"github.com/donnie4w/go-logger/logger"
	. "tim.common"
	"tim.daoService"
	. "tim.protocol"
	"tim.route"
	"tim.utils"
)

type TimImpl struct {
	Ip   string
	Port int
	Pub  string //发布id
}

// Parameters:
//  - Param
func (this *TimImpl) TimStream(param *TimParam) (err error) {
	panic("error")
	return
}
func (this *TimImpl) TimStarttls() (err error) {
	panic("error")
	return
}

// Parameters:
//  - Tid
//  - Pwd
func (this *TimImpl) TimLogin(tid *Tid, pwd string) (err error) {
	logger.Debug("Login:", tid, pwd)
	panic("error")
	return
}

// Parameters:
//  - Ab
func (this *TimImpl) TimAck(ab *TimAckBean) (err error) {
	logger.Debug("TimAck=========>", ab)
	panic("error")
	return
}

// Parameters:
//  - Pbean
func (this *TimImpl) TimPresence(pbean *TimPBean) (err error) {
	logger.Debug(pbean)
	panic("error")
	return
}

// Parameters:
//  - Mbean
func (this *TimImpl) TimMessage(mbean *TimMBean) (err error) {
	logger.Debug(mbean)
	panic("error")
	return
}

// Parameters:
//  - ThreadId
func (this *TimImpl) TimPing(threadId string) (err error) {
	panic("error")
	return
}

// Parameters:
//  - E
func (this *TimImpl) TimError(e *TimError) (err error) {
	panic("error")
	return
}
func (this *TimImpl) TimLogout() (err error) {
	panic("error")
	return
}

// Parameters:
//  - Tid
//  - Pwd
func (this *TimImpl) TimRegist(tid *Tid, pwd string) (err error) {
	panic("error")
	return
}

// Parameters:
//  - Tid
//  - Pwd
func (this *TimImpl) TimRemoteUserAuth(tid *Tid, pwd string, auth *TimAuth) (r *TimRemoteUserBean, err error) {
	panic("error")
	return
}

// Parameters:
//  - Tid
func (this *TimImpl) TimRemoteUserGet(tid *Tid, auth *TimAuth) (r *TimRemoteUserBean, err error) {
	panic("error")
	return
}

// Parameters:
//  - Tid
//  - Ub
func (this *TimImpl) TimRemoteUserEdit(tid *Tid, ub *TimUserBean, auth *TimAuth) (r *TimRemoteUserBean, err error) {
	panic("error")
	return
}

// Parameters:
//  - Pbean
func (this *TimImpl) TimResponsePresence(pbean *TimPBean, auth *TimAuth) (r *TimResponseBean, err error) {
	logger.Debug("TimResponsePresence", pbean, auth)
	if !checkAuth(auth) {
		err = errors.New(fmt.Sprint("cluster auth fail:", auth))
		return
	}
	go _TimResponsePresence(pbean, auth)
	return
}

func _TimResponsePresence(pbean *TimPBean, auth *TimAuth) (r *TimResponseBean, err error) {
	defer func() {
		if er := recover(); er != nil {
			logger.Error("error:", er)
		}
	}()
	if pbean.GetToTid() == nil {
		route.RoutePBean(pbean)
	} else {
		route.RouteSinglePBean(pbean)
	}
	return
}

// Parameters:
//  - Mbean
func (this *TimImpl) TimResponseMessage(mbean *TimMBean, auth *TimAuth) (r *TimResponseBean, err error) {
	logger.Debug("TimResponseMessage", mbean, auth)
	if !checkAuth(auth) {
		err = errors.New(fmt.Sprint("cluster auth fail:", auth))
		return
	}
	go _TimResponseMessage(mbean, auth)
	return
}

func _TimResponseMessage(mbean *TimMBean, auth *TimAuth) (r *TimResponseBean, err error) {
	defer func() {
		if er := recover(); er != nil {
			logger.Error("error:", er)
		}
	}()
	//	r = NewTimResponseBean()
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
		_, err, _ = route.RouteMBean(mbean, false, false)
	} else {
		err = errors.New("TimResponseMessage totid not exist")
	}
	return
}

func (this *TimImpl) TimMessageIq(timMsgIq *TimMessageIq, iqType string) (err error) {
	logger.Debug("TimMessageIq:", timMsgIq, " ", iqType)
	panic("error")
	return
}

// Parameters:
//  - Mbean
func (this *TimImpl) TimMessageResult_(mbean *TimMBean) (err error) {
	logger.Debug("TimMessageResult_:", mbean)
	panic("error")
	return
}

func (this *TimImpl) TimRoser(roster *TimRoster) (err error) {
	logger.Debug("TimRoser:", roster)
	panic("error")
	return
}

func checkAuth(a *TimAuth) bool {
	if a.GetDomain() == ClusterConf.Domain && a.GetUsername() == ClusterConf.Username && a.GetPwd() == ClusterConf.Password {
		return true
	}
	return false
}

func (this *TimImpl) TimResponseMessageIq(timMsgIq *TimMessageIq, iqType string, auth *TimAuth) (r *TimMBeanList, err error) {
	logger.Debug("TimResponseMessageIq:", timMsgIq, iqType, auth)
	panic("error TimResponseMessageIq")
	return
}

func (this *TimImpl) TimMessageList(mbeanList *TimMBeanList) (err error) {
	logger.Debug("TimMessageList:", mbeanList)
	panic("error TimMessageList")
	return
}

// Parameters:
//  - PbeanList
func (this *TimImpl) TimPresenceList(pbeanList *TimPBeanList) (err error) {
	logger.Debug("TimPresenceList:", pbeanList)
	panic("error TimPresenceList")
	return
}

func (this *TimImpl) TimResponsePresenceList(pbeanList *TimPBeanList, auth *TimAuth) (r *TimResponseBean, err error) {
	if !checkAuth(auth) {
		err = errors.New(fmt.Sprint("cluster TimResponsePresenceList fail:", auth))
		return
	}
	defer func() {
		if er := recover(); er != nil {
			logger.Error("error:", er)
		}
	}()
	go _TimResponsePresenceList(pbeanList, auth)
	return
}

func _TimResponsePresenceList(pbeanList *TimPBeanList, auth *TimAuth) (r *TimResponseBean, err error) {
	defer func() {
		if er := recover(); er != nil {
			logger.Error("error:", er)
		}
	}()
	if pbeanList != nil && pbeanList.GetTimPBeanList() != nil && len(pbeanList.GetTimPBeanList()) > 0 {
		if ClusterConf.Interflow > 0 && len(pbeanList.GetTimPBeanList()) > 1 {
			route.RoutePBeanList(pbeanList.GetTimPBeanList())
		} else {
			for _, pbean := range pbeanList.GetTimPBeanList() {
				_TimResponsePresence(pbean, auth)
			}
		}
	}
	return
}

// Parameters:
//  - MbeanList
//  - Auth
func (this *TimImpl) TimResponseMessageList(mbeanList *TimMBeanList, auth *TimAuth) (r *TimResponseBean, err error) {
	if !checkAuth(auth) {
		err = errors.New(fmt.Sprint("cluster TimResponseMessageList fail:", auth))
		return
	}
	defer func() {
		if er := recover(); er != nil {
			logger.Error("error:", er)
		}
	}()
	go _TimResponseMessageList(mbeanList, auth)
	return
}

func _TimResponseMessageList(mbeanList *TimMBeanList, auth *TimAuth) (r *TimResponseBean, err error) {
	defer func() {
		if er := recover(); er != nil {
			logger.Error("error:", er)
		}
	}()
	if mbeanList != nil && mbeanList.GetTimMBeanList() != nil && len(mbeanList.GetTimMBeanList()) > 0 {
		if ClusterConf.Interflow > 0 && len(mbeanList.GetTimMBeanList()) > 1 {
			route.RouteMBeanList(mbeanList.GetTimMBeanList(), true)
		} else {
			for _, mbean := range mbeanList.GetTimMBeanList() {
				route.RouteMBean(mbean, false, true)
			}
		}
	}
	return
}

func (this *TimImpl) TimProperty(tpb *TimPropertyBean) (err error) {
	logger.Debug("TimProperty:", tpb)
	return
}
