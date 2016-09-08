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
	//	. "tim.Map"
	. "tim.common"
	. "tim.protocol"
	. "tim.utils"
)

type TimPool struct {
	// 注册了的连接器
	pool map[*TimUser]bool
	//	pool *HashTable
	// 在线用户  domain+name
	poolUser map[string]map[*TimUser]bool
	//	poolUser *HashTable
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

func (this *TimPool) AddSingleAuthUser(tu *TimUser) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("AddSingleAuthUser,", err)
			logger.Error(string(debug.Stack()))
		}
	}()
	this.rwLock.Lock()
	defer this.rwLock.Unlock()
	loginname, _ := GetLoginName(tu.UserTid)
	if pu, ok := this.poolUser[loginname]; ok {
		if pu != nil && len(pu) > 0 {
			for k, _ := range pu {
				k.Close()
				delete(pu, k)
				delete(this.pool, k)
			}
		}
	}
	if _, ok := this.pool[tu]; ok {
		if this.poolUser[loginname] == nil {
			this.poolUser[loginname] = make(map[*TimUser]bool)
		}
		this.poolUser[loginname][tu] = true
	}
}

func (this *TimPool) GetAllLoginName() (ss []string) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("GetAllLoginName,", err)
			logger.Error(string(debug.Stack()))
		}
	}()
	//	this.rwLock.RLock()
	//	defer this.rwLock.RUnlock()
	ss = make([]string, 0)
	for k, _ := range this.poolUser {
		ss = append(ss, k)
	}
	return
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
	tu.Close()
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
	for _, vmap := range t.poolUser {
		for tu, _ := range vmap {
			str = fmt.Sprintln(str, tu.UserTid.GetName(), " # ", TimeMills2TimeFormat(tu.IdCardNo), " # ", tu.UserTid.GetResource())
		}
	}
	return str
}

func (t *TimPool) GetLoginUser(loginname string) (tus []*TimUser) {
	t.rwLock.RLock()
	defer t.rwLock.RUnlock()
	if tumap, ok := t.poolUser[loginname]; ok {
		tus = make([]*TimUser, 0)
		for tu, _ := range tumap {
			tus = append(tus, tu)
		}
	}
	return
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
	UserType         int // 0 client  1 cluster client
	TLS              int //
	Interflow        int
	Version          int16
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
	if CF.SingleClient == 1 {
		TP.AddSingleAuthUser(t)
	} else {
		TP.AddAuthUser(t)
	}
	return
}

func (t *TimUser) SendMBean(mbean *TimMBean) (er error) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("sendMBean error:", err)
			logger.Error(string(debug.Stack()))
			er = errors.New("sendMBean err")
		}
	}()
	t.Sync.Lock()
	defer t.Sync.Unlock()
	if t.IsClose {
		er = errors.New("timuser is close")
		return
	}
	if CF.ConfirmAck == 1 {
		timer := time.NewTicker(3 * time.Second)
		t.LastSyncThreadId = mbean.GetThreadId()
		er = t.Client.TimMessage(mbean)
		if er == nil {
			select {
			case <-timer.C:
				er = errors.New(fmt.Sprint(t.UserTid.GetName(), ", send ack overtime:", mbean.GetThreadId(), "  ", t))
				logger.Error("send ack overtime:", mbean.GetThreadId())
			case threadId := <-t.Sendflag:
				if t.LastSyncThreadId != threadId {
					er = errors.New(fmt.Sprint("error msg ack threadid:", t.LastSyncThreadId, "!=", threadId))
					logger.Error("error msg ack threadid:", t.LastSyncThreadId, "!=", threadId)
				}
			}
		} else {
			logger.Error("sendMBean:", er.Error())
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

func (t *TimUser) SendMBeanList(mbeans []*TimMBean) (er error) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("SendMBeanList error:", err)
			logger.Error(string(debug.Stack()))
			er = errors.New("SendMBeanList err")
		}
	}()
	t.Sync.Lock()
	defer t.Sync.Unlock()
	if t.IsClose {
		er = errors.New("timuser is close")
		return
	}
	mbeanList := NewTimMBeanList()
	mbeanList.TimMBeanList = mbeans
	mbeanList.ThreadId = TimeMills()
	if CF.ConfirmAck == 1 {
		timer := time.NewTicker(5 * time.Second)
		t.LastSyncThreadId = mbeanList.GetThreadId()
		er = t.Client.TimMessageList(mbeanList)
		select {
		case <-timer.C:
			er = errors.New(fmt.Sprint("send ack overtime:", mbeanList.GetThreadId()))
			logger.Error("send ack overtime:", mbeanList.GetThreadId())
		case threadId := <-t.Sendflag:
			if mbeanList.GetThreadId() != threadId {
				er = errors.New(fmt.Sprint("error msg ack threadid:", mbeanList.GetThreadId(), "!=", threadId))
				logger.Error("error msg ack threadid:", mbeanList.GetThreadId(), "!=", threadId)
			}
		}
	} else {
		er = t.Client.TimMessageList(mbeanList)
	}
	if er != nil {
		t.IsClose = true
		t.Fw = CLOSE
	}
	return
}

func (t *TimUser) Ping() (er error) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("Ping,", err)
			logger.Error(string(debug.Stack()))
			er = errors.New("Ping err")
		}
	}()
	t.Sync.Lock()
	defer t.Sync.Unlock()
	er = t.Client.TimPing(TimeMills())
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
	if CF.Presence != 1 {
		return
	}
	t.Sync.Lock()
	defer t.Sync.Unlock()
	er = t.Client.TimPresence(pbean)
	return
}

func (t *TimUser) SendPBeanList(pbean []*TimPBean) (er error) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("SendPBeanList,", err)
			logger.Error(string(debug.Stack()))
			er = errors.New("SendPBeanList err")
		}
	}()
	if CF.Presence != 1 {
		return
	}
	t.Sync.Lock()
	defer t.Sync.Unlock()
	pbeanList := NewTimPBeanList()
	pbeanList.ThreadId = TimeMills()
	pbeanList.TimPBeanList = pbean
	er = t.Client.TimPresenceList(pbeanList)
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

func (t *TimUser) Close() (er error) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("Close,", err)
			logger.Error(string(debug.Stack()))
			er = errors.New("Close err")
		}
	}()
	t.Sync.Lock()
	defer t.Sync.Unlock()
	er = t.Client.Transport.Close()
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
	if tid == nil || tid.GetDomain() == "" || tid.GetName() == "" {
		return "", errors.New("tid error")
	}
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
