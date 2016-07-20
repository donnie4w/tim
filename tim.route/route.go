/**
 * donnie4w@gmail.com  tim server
 */
package route

import (
	"errors"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/donnie4w/go-logger/logger"
	. "tim.connect"
	"tim.daoService"
	. "tim.protocol"
)

/**********************************************Message***********************************************/
/**Message*/
func RouteMBean(mbean *TimMBean, isSingle, async bool) (id int32, er error, offline bool) {
	defer func() {
		if err := recover(); err != nil {
			er = errors.New(fmt.Sprint("RouteMBean:", err))
			logger.Error("RouteMBean,", err)
			logger.Error(string(debug.Stack()))
		}
	}()
	logger.Debug("RouteMBean:", mbean)
	loginname, _ := GetLoginName(mbean.GetToTid())
	if isSingle {
		id, _ = daoService.SaveSingleMBean(mbean)
	} else {
		id, _ = daoService.SaveMBean(mbean)
	}
	if async {
		go func() {
			defer func() {
				if err := recover(); err != nil {
					logger.Error(string(debug.Stack()))
				}
			}()
			mid := fmt.Sprint(id)
			mbean.Mid = &mid
			//			if tumap, ok := TP.PoolUser[loginname]; ok {
			tumap := TP.GetLoginUser(loginname)
			if tumap != nil {
				if len(tumap) > 0 {
					isSendok := false
					for tu, _ := range tumap {
						err := tu.SendMBean(mbean)
						if err != nil {
							logger.Error("routemessage :", err)
						} else {
							isSendok = true
						}
					}
					if !isSendok {
						daoService.SaveOfflineMBean(mbean)
					}
				}
			} else {
				daoService.SaveOfflineMBean(mbean)
			}
		}()
	} else {
		mid := fmt.Sprint(id)
		mbean.Mid = &mid
		//		if tumap, ok := TP.PoolUser[loginname]; ok {
		tumap := TP.GetLoginUser(loginname)
		if tumap != nil {
			if len(tumap) > 0 {
				isSendok := false
				for tu, _ := range tumap {
					err := tu.SendMBean(mbean)
					if err != nil {
						logger.Error("routemessage :", err)
					} else {
						isSendok = true
					}
				}
				if !isSendok {
					daoService.SaveOfflineMBean(mbean)
					offline = true
				}
			}
		} else {
			daoService.SaveOfflineMBean(mbean)
			offline = true
		}
	}
	return
}

func RouteOffLineMBean(tu *TimUser) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("RouteOffLineMBean,", err)
			logger.Error(string(debug.Stack()))
		}
	}()
	time.Sleep(3 * time.Second)
	mbeans := daoService.LoadOfflineMBean(tu.UserTid)
	if mbeans != nil && len(mbeans) > 0 {
		for _, mbean := range mbeans {
			err := tu.SendMBean(mbean)
			if err != nil {
				break
			} else {
				go daoService.DelteOfflineMBean(mbean.Mid)
				go daoService.UpdateOffMessage(mbean, 1)
			}
		}
	}
}

/**********************************************Presence***********************************************/
/**Presence*/
func RoutePBean(pbean *TimPBean) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("RoutePBean,", err)
			logger.Error(string(debug.Stack()))
		}
	}()
	fromtid := pbean.GetFromTid()
	tids := daoService.GetOnlineRoser(fromtid)
	if tids != nil {
		for _, tid := range tids {
			loginname, _ := GetLoginName(tid)
			tumap := TP.GetLoginUser(loginname)
			if tumap != nil {
				if len(tumap) > 0 {
					for tu, _ := range tumap {
						pbean.ToTid = tu.UserTid
						tu.SendPBean(pbean)
					}
				}
			}
		}
	}
}
