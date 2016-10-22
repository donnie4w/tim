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
	. "tim.common"
	. "tim.connect"
	"tim.daoService"
	. "tim.protocol"
	"tim.utils"
)

/**********************************************Message***********************************************/
/**Message*/
func RouteMBean(mbean *TimMBean, isSingle, async bool) (mid string, er error, offline bool) {
	defer func() {
		if err := recover(); err != nil {
			er = errors.New(fmt.Sprint("RouteMBean:", err))
			logger.Error("RouteMBean,", err)
			logger.Error(string(debug.Stack()))
		}
	}()

	loginname, _ := GetLoginName(mbean.GetToTid())
	if CF.Db_Exsit == 0 {
		mid = fmt.Sprint(utils.GetRand(100000000))
	} else {
		if isSingle {
			mid, _, er = daoService.SaveSingleMBean(mbean)
		} else {
			if mbean.GetType() == "groupchat" {
				mid, er = daoService.SaveMucMBean(mbean)
			} else {
				mid, _, er = daoService.SaveMBean(mbean)
			}
		}
	}
	if er != nil {
		return
	}
	if async {
		go func() {
			defer func() {
				if err := recover(); err != nil {
					logger.Error(string(debug.Stack()))
				}
			}()
			mbean.Mid = &mid
			tus := TP.GetLoginUser(loginname)
			if tus != nil {
				if len(tus) > 0 {
					isSendok := false
					for _, tu := range tus {
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
		mbean.Mid = &mid
		tus := TP.GetLoginUser(loginname)
		if tus != nil {
			if len(tus) > 0 {
				isSendok := false
				for _, tu := range tus {
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

func SaveMBean(mbean *TimMBean) {
	daoService.SaveMBean(mbean)
}

/**Message List */
func RouteMBeanList(mbeans []*TimMBean, async bool) {
	if async {
		go _RouteMBeanList(mbeans)
	} else {
		_RouteMBeanList(mbeans)
	}
}

func _RouteMBeanList(mbeans []*TimMBean) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(string(debug.Stack()))
		}
	}()
	if mbeans != nil && len(mbeans) > 0 {
		loginnamemap := make(map[string][]*TimMBean, 0)
		for _, mbean := range mbeans {
			SaveMBean(mbean)
			loginname, _ := GetLoginName(mbean.GetToTid())
			if _, ok := loginnamemap[loginname]; !ok {
				loginnamemap[loginname] = make([]*TimMBean, 0)
			}
			loginnamemap[loginname] = append(loginnamemap[loginname], mbean)
		}
		if len(loginnamemap) > 0 {
			for k, v := range loginnamemap {
				tus := TP.GetLoginUser(k)
				if tus != nil {
					if len(tus) > 0 {
						isSendok := false
						for _, tu := range tus {
							err := tu.SendMBeanList(v)
							if err != nil {
								logger.Error("routemessage :", err)
							} else {
								isSendok = true
							}
						}
						if !isSendok {
							daoService.SaveOfflineMBeanList(v)
						}
					}
				} else {
					daoService.SaveOfflineMBeanList(v)
				}
			}
		}
	}
}

func RouteOffLineMBean(tu *TimUser) (er error) {
	if CF.Db_Exsit == 0 {
		return
	}
	defer func() {
		if err := recover(); err != nil {
			logger.Error("RouteOffLineMBean,", err)
			logger.Error(string(debug.Stack()))
		}
	}()
	time.Sleep(1000 * time.Millisecond)
	mbeans := daoService.LoadOfflineMBean(tu.UserTid)
	if mbeans != nil && len(mbeans) > 0 {
		if tu.Interflow > 0 {
			mids := make([]interface{}, 0)
			for _, mbean := range mbeans {
				mids = append(mids, mbean.GetMid())
			}
			err := tu.SendMBeanList(mbeans)
			if err == nil {
				daoService.DelOfflineMBeanList(mids...)
				daoService.UpdateOffMessageList(mbeans, 1)
			}
		} else {
			for _, mbean := range mbeans {
				err := tu.SendMBean(mbean)
				if err != nil {
					er = err
					break
				} else {
					go daoService.DelOfflineMBean(mbean.Mid)
					go daoService.UpdateOffMessage(mbean, 1)
				}
			}
		}
	}
	mbeans = daoService.LoadOfflineMucMBean(tu.UserTid)
	if mbeans != nil && len(mbeans) > 0 {
		if tu.Interflow > 0 {
			mids := make([]interface{}, 0)
			for _, mbean := range mbeans {
				mids = append(mids, mbean.GetMid())
			}
			err := tu.SendMBeanList(mbeans)
			if err == nil {
				daoService.DelOfflineMucMBeanList(mids...)
			}
		} else {
			for _, mbean := range mbeans {
				err := tu.SendMBean(mbean)
				if err != nil {
					er = err
					break
				} else {
					go daoService.DelOfflineMucMBean(mbean.Mid)
				}
			}
		}
	}
	return
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
			tus := TP.GetLoginUser(loginname)
			if tus != nil {
				if len(tus) > 0 {
					for _, tu := range tus {
						pbean.ToTid = tu.UserTid
						tu.SendPBean(pbean)
					}
				}
			}
		}
	}
}

/**Presence list */
func RoutePBeanList(pbeans []*TimPBean) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("RoutePBeanList:", err)
			logger.Error(string(debug.Stack()))
		}
	}()
	if pbeans != nil && len(pbeans) > 0 {
		loginnamemap := make(map[string][]*TimPBean, 0)
		for _, pbean := range pbeans {
			fromtid := pbean.GetFromTid()
			tids := daoService.GetOnlineRoser(fromtid)
			if tids != nil {
				for _, tid := range tids {
					loginname, _ := GetLoginName(tid)
					if _, ok := loginnamemap[loginname]; !ok {
						loginnamemap[loginname] = make([]*TimPBean, 0)
					}
					loginnamemap[loginname] = append(loginnamemap[loginname], pbean)
				}
			}
		}
		if len(loginnamemap) > 0 {
			for k, v := range loginnamemap {
				tus := TP.GetLoginUser(k)
				if tus != nil {
					if len(tus) > 0 {
						for _, tu := range tus {
							tu.SendPBeanList(v)
						}
					}
				}
			}
		}
	}
}

func RouteSinglePBean(pbean *TimPBean) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("RouteSinglePBean,", err)
			logger.Error(string(debug.Stack()))
		}
	}()
	tid := pbean.GetToTid()
	if tid != nil {
		loginname, _ := GetLoginName(tid)
		tus := TP.GetLoginUser(loginname)
		if tus != nil && len(tus) > 0 {
			for _, tu := range tus {
				pbean.ToTid = tu.UserTid
				tu.SendPBean(pbean)
			}
		}
	}
}
