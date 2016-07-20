/**
 * donnie4w@gmail.com  tim server
 */
package connect

import (
	"errors"
	"fmt"
	"runtime/debug"
	"sync"
	"time"

	"github.com/donnie4w/go-logger/logger"
	. "tim.FW"
	. "tim.common"
	. "tim.protocol"
	. "tim.utils"
)

type PoolMap struct {
	PoolBean map[*TimUser]bool
}

type PoolUserMap struct {
	PoolBean map[string]map[*TimUser]bool
}

type TimPool struct {
	// 注册了的连接器
	pool map[*TimUser]bool
	// 在线用户  domain+name
	poolUser map[string]map[*TimUser]bool
	// 从连接器中注册请求
	register chan *TimUser
	// 从连接器中注销请求
	unregister chan *TimUser
	// 锁
	rwLock *sync.RWMutex
}

func (this *TimPool) AddAuthUser(tu *TimUser) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("AddAuthUser,", err)
			logger.Error(string(debug.Stack()))
		}
	}()
	this.rwLock.Lock()
	defer this.rwLock.Unlock()
	if _, ok := this.pool[tu]; ok {
		loginname, _ := GetLoginName(tu.UserTid)
		if len(this.poolUser[loginname]) == 0 {
			this.poolUser[loginname] = make(map[*TimUser]bool)
		}
		this.poolUser[loginname][tu] = true
	}
}

func (this *TimPool) DeleteTimUser(tu *TimUser) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("DeleteTimUser,", err)
			logger.Error(string(debug.Stack()))
		}
	}()
	this.rwLock.Lock()
	defer this.rwLock.Unlock()
	if tu.UserTid != nil {
		loginname, _ := GetLoginName(tu.UserTid)
		if tumap, ok := this.poolUser[loginname]; ok {
			if _, ok := tumap[tu]; ok {
				delete(tumap, tu)
				if len(tumap) == 0 {
					delete(this.poolUser, loginname)
				}
			}
		}
	}
	if _, ok := this.pool[tu]; ok {
		delete(this.pool, tu)
	}
}

func (t *TimPool) Len4PU() int {
	t.rwLock.RLock()
	defer t.rwLock.RUnlock()
	return len(t.poolUser)
}

func (t *TimPool) Len4P() int {
	t.rwLock.RLock()
	defer t.rwLock.RUnlock()
	return len(t.pool)
}

func (t *TimPool) PrintUsersInfo() string {
	t.rwLock.RLock()
	defer t.rwLock.RUnlock()
	str := fmt.Sprintln(len(t.pool), "======>", len(t.poolUser))
	for _, vmap := range TP.poolUser {
		for tu, _ := range vmap {
			str = fmt.Sprintln(str, tu.UserTid.GetName(), " # ", TimeMills2TimeFormat(tu.IdCardNo), " # ", tu.UserTid.GetResource())
		}
	}
	return str
}

func (t *TimPool) GetLoginUser(loginname string) map[*TimUser]bool {
	t.rwLock.RLock()
	defer t.rwLock.RUnlock()
	if tumap, ok := TP.poolUser[loginname]; ok {
		return tumap
	} else {
		return nil
	}
}

func (t *TimPool) AddConnect(c *TimUser) {
	t.rwLock.Lock()
	defer t.rwLock.Unlock()
	t.pool[c] = true
}

type TimUser struct {
	UserTid          *Tid
	Client           *ITimClient
	Fw               FLOW
	OverLimit        int
	IdCardNo         string
	IsClose          bool
	Sendflag         chan string
	Sync             *sync.Mutex
	LastSyncThreadId string
}

func (t *TimUser) Auth(tid *Tid) (er error) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("Auth,", err)
			logger.Error(string(debug.Stack()))
			er = errors.New("Auth err")
		}
	}()
	t.Sync.Lock()
	defer t.Sync.Unlock()
	t.UserTid = tid
	TP.AddAuthUser(t)
	return
}

func (t *TimUser) SendMBean(mbean *TimMBean) (er error) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("sendMBean,", err)
			logger.Error(string(debug.Stack()))
			er = errors.New("sendMBean err")
		}
	}()
	logger.Debug("SendMBean", mbean)
	t.Sync.Lock()
	defer t.Sync.Unlock()
	if t.IsClose {
		er = errors.New("timuser is close")
		return
	}
	if ConfBean.ConfirmAck == 1 {
		timer := time.NewTicker(3 * time.Second)
		t.LastSyncThreadId = mbean.GetThreadId()
		er = t.Client.TimMessage(mbean)
		select {
		case <-timer.C:
			er = errors.New("send ack overtime")
		case threadId := <-t.Sendflag:
			if mbean.GetThreadId() != threadId {
				er = errors.New("error msg ack threadid")
			}
		}
	} else {
		er = t.Client.TimMessage(mbean)
	}
	if er != nil {
		t.IsClose = true
		t.Fw = CLOSE
	}
	return
}

func (t *TimUser) SendPBean(pbean *TimPBean) (er error) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("sendPBean,", err)
			logger.Error(string(debug.Stack()))
			er = errors.New("sendPBean err")
		}
	}()
	if ConfBean.Presence != 1 {
		return
	}
	t.Sync.Lock()
	defer t.Sync.Unlock()
	er = t.Client.TimPresence(pbean)
	return
}

func (t *TimUser) SendAckBean(ackBean *TimAckBean) (er error) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("sendAckBean,", err)
			logger.Error(string(debug.Stack()))
			er = errors.New("sendAckBean err")
		}
	}()
	t.Sync.Lock()
	defer t.Sync.Unlock()
	er = t.Client.TimAck(ackBean)
	return
}

func GetLoginName(tid *Tid) (loginname string, er error) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("GetLoginName,", err)
			logger.Error(string(debug.Stack()))
			er = errors.New("GetLoginName err")
		}
	}()
	domain := tid.GetDomain()
	name := tid.GetName()
	loginname = MD5(fmt.Sprint(domain, "tim", name))
	return
}

var TP = TimPool{
	register:   make(chan *TimUser),
	unregister: make(chan *TimUser),
	pool:       make(map[*TimUser]bool),
	poolUser:   make(map[string]map[*TimUser]bool),
	rwLock:     new(sync.RWMutex),
}

//func (t *TimPool) Run() {
//	defer func() {
//		if err := recover(); err != nil {
//			logger.Error("Run,", err)
//			logger.Error(string(debug.Stack()))
//		}
//	}()
//	for {
//		select {
//		case c := <-t.Register:
//			t.AddConnect(c)
//		case c := <-t.Unregister:
//			t.DeleteTimUser(c)
//		}
//		//logger.Debug("pool===>", t.Len4P())
//	}
//}
