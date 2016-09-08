/**
 * donnie4w@gmail.com  tim server
 */
package clusterRoute

import (
	"errors"
	"fmt"
	"runtime/debug"

	"github.com/donnie4w/go-logger/logger"
	"tim.cluster"
	. "tim.connect"
	"tim.daoService"
	. "tim.protocol"
	"tim.route"
)

/**********************************************Message***********************************************/
/**Message*/

func ClusterRouteMBean(mbean *TimMBean, beans []*cluster.CluserUserBean) (er error) {
	er = errors.New("ClusterRouteMBean fail")
	if cluster.IsCluster() {
		for _, bean := range beans {
			err := bean.SendMBean(mbean)
			if err == nil {
				er = nil
			}
		}
	}
	return
}

func ClusterRoutePBean(pbean *TimPBean) (er error) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("ClusterRoutePBean:", err)
			logger.Error(string(debug.Stack()))
			er = errors.New(fmt.Sprint(err))
		}
	}()
	if cluster.IsCluster() {
		if pbean.GetToTid() == nil {
			fromtid := pbean.GetFromTid()
			tids := daoService.GetOnlineRoser(fromtid)
			if tids != nil {
				for _, tid := range tids {
					beans := OtherClusterUserBean(tid)
					if beans != nil && len(beans) > 0 {
						for _, bean := range beans {
							bean.SendPBean(pbean)
						}
					}
				}
			}
		} else {
			beans := OtherClusterUserBean(pbean.GetToTid())
			if beans != nil && len(beans) > 0 {
				for _, bean := range beans {
					bean.SendPBean(pbean)
				}
			} else {
				route.RouteSinglePBean(pbean)
			}
		}
	}
	return
}

func OtherClusterUserBean(tid *Tid) (beans []*cluster.CluserUserBean) {
	var err error
	loginname, _ := GetLoginName(tid)
	beans, err = cluster.GetUserBeans(loginname)
	if beans == nil || len(beans) == 0 || err != nil {
		return nil
	}
	return
}
