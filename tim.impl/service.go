/**
 * donnie4w@gmail.com  tim server
 */
package impl

import (
	. "tim.protocol"
	"tim.utils"
)

func newTid(name string, domain, resource *string) *Tid {
	tid := NewTid()
	tid.Domain = domain
	tid.Name = name
	tid.Resource = resource
	return tid
}

func OnlinePBean(tid *Tid) (pbean *TimPBean) {
	pbean = NewTimPBean()
	pbean.ThreadId = utils.TimeMills()
	pbean.FromTid = tid
	show, status := "online", "probe"
	pbean.Show, pbean.Status = &show, &status
	return
}

func OfflinePBean(tid *Tid) (pbean *TimPBean) {
	pbean = NewTimPBean()
	pbean.ThreadId = utils.TimeMills()
	pbean.FromTid = tid
	show, status := "offline", "unavailable"
	pbean.Show, pbean.Status = &show, &status
	return
}
