/**
 * donnie4w@gmail.com  tim server
 */
package cluster

import (
	"errors"
	"fmt"
	"os"
	"runtime/debug"
	"sync"
	"time"

	"github.com/donnie4w/go-logger/logger"

	. "tim.Map"
	. "tim.common"
	"tim.daoService"
	. "tim.protocol"
	"tim.utils"
)

var flag bool = false
var tableString string = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+#=/"

type CluserUserBean struct {
	Tcp_addr string
	Lock     *sync.Mutex
}

var sha1Addcmd *ScriptCmd
var sha1Getcmd *ScriptCmd

func InitCluster(filexml string) {
	if ClusterConf.Init(filexml) {
		Redis.initPool()
		var err error
		flag, err = Redis.Ping()
		if !flag && err != nil {
			logger.Error("redis connect failed:", err.Error())
			os.Exit(1)
		}
		sha1Addcmd = Redis.NewScript(scriptAddCmd, 2)
		sha1Getcmd = Redis.NewScript(scriptGetCmd, 1)
	}
}

func IsCluster() bool {
	return flag && ClusterConf.IsCluster == 1
}

/**
 *
 */
func GetUserBeans(loginname string) (beans []*CluserUserBean, err error) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(string(debug.Stack()))
		}
	}()
	var addrs []string
	addrs, err = GetLoginnameFromCluter(loginname)
	if err == nil {
		beans = make([]*CluserUserBean, 0)
		for _, addr := range addrs {
			localAddr := parseAddr(ClusterConf.RequestAddr)
			if localAddr != addr {
				cu := new(CluserUserBean)
				cu.Tcp_addr = formatAddr(addr)
				cu.Lock = new(sync.Mutex)
				beans = append(beans, cu)
			}
		}
	}
	return
}

/**
 *   key               field             value
 *   loginname        ip:tcp_port		http_port
 */
func SetLoginnameToCluster(loginname string) (err error) {
	_, err = sha1Addcmd.EvalSha(loginname, utils.MD5(fmt.Sprint(loginname, parseAddr(ClusterConf.RequestAddr))), fmt.Sprint(ClusterConf.Keytimeout), parseAddr(ClusterConf.RequestAddr))
	return
}

func DelLoginnameFromCluter(loginname string) (err error) {
	err = Redis.Del(utils.MD5(fmt.Sprint(loginname, parseAddr(ClusterConf.RequestAddr))))
	return
}

func GetLoginnameFromCluter(loginname string) (addrs []string, err error) {
	addrs, err = sha1Getcmd.EvalShaStrings(loginname)
	return
}

func (this *CluserUserBean) SendMBean(tmb *TimMBean) (err error) {
	defer func() {
		if er := recover(); er != nil {
			err = errors.New(fmt.Sprint("SendMBean:", er))
			logger.Error("SendMBean,", er)
			logger.Error(string(debug.Stack()))
		}
	}()
	this.Lock.Lock()
	defer this.Lock.Unlock()
	if ClusterConf.Interflow > 0 {
		er := InterflowSendMBean(this.Tcp_addr, tmb)
		if er != nil {
			return _sendMBean(this.Tcp_addr, tmb, 3)
		}
	} else {
		return _sendMBean(this.Tcp_addr, tmb, 3)
	}
	return
}

func _sendMBean(addr string, tmb *TimMBean, count int) (err error) {
	if count <= 0 {
		err = errors.New("over limitcount")
		daoService.SaveOfflineMBean(tmb)
	} else {
		count--
		client := Pool.Get(addr)
		_, err = client.SendMBean(tmb, NewAuth())
		if err != nil {
			_sendMBean(addr, tmb, count)
		} else {
			Pool.Put(addr, client)
		}
	}
	return
}

func (this *CluserUserBean) SendPBean(tmpb *TimPBean) (err error) {
	defer func() {
		if er := recover(); er != nil {
			err = errors.New(fmt.Sprint("SendPBean:", er))
			logger.Error("SendPBean,", er)
			logger.Error(string(debug.Stack()))
		}
	}()
	this.Lock.Lock()
	defer this.Lock.Unlock()
	if ClusterConf.Interflow > 0 {
		er := InterflowSendPBean(this.Tcp_addr, tmpb)
		if er != nil {
			client := Pool.Get(this.Tcp_addr)
			_, err := client.SendPBean(tmpb, NewAuth())
			if err == nil {
				Pool.Put(this.Tcp_addr, client)
			}
		}
	} else {
		client := Pool.Get(this.Tcp_addr)
		_, err := client.SendPBean(tmpb, NewAuth())
		if err == nil {
			Pool.Put(this.Tcp_addr, client)
		}
	}
	return
}

func GetRandomTBString(length int) (s string) {
	l := len(tableString)
	rand := utils.GetRand(l)
	for i := 0; i < length; i++ {
		if rand > 0 {
			s = fmt.Sprint(s, tableString[rand-1:rand])
		} else {
			s = fmt.Sprint(s, tableString[0:1])
		}
	}
	return
}

func NewAuth() (ta *TimAuth) {
	ta = NewTimAuth()
	domain, name, pwd := ClusterConf.Domain, ClusterConf.Username, ClusterConf.Password
	ta.Domain, ta.Username, ta.Pwd = &domain, &name, &pwd
	return
}

//---------------------------------------------------------------------------------------------

type InterflowPool struct {
	//	pool  map[string]*InterflowBean
	//	pool2 map[string]*InterflowBean
	pool  *HashTable
	pool2 *HashTable
	//	lock  *sync.RWMutex
}

var IFPool = &InterflowPool{pool: NewHashTable(), pool2: NewHashTable()}

func (this *InterflowPool) add(addr string, ib *InterflowBean, _type int) {
	//	this.lock.Lock()
	//	defer this.lock.Unlock()
	if _type == 1 {
		this.pool.Put(addr, ib)
	} else if _type == 2 {
		//		this.pool2[addr] = ib
		this.pool2.Put(addr, ib)
	}
}

func (this *InterflowPool) get(addr string, _type int) *InterflowBean {
	//	this.lock.RLock()
	//	defer this.lock.RUnlock()
	if _type == 1 {
		i := this.pool.Get(addr)
		if i != nil {
			return i.(*InterflowBean)
		}
	} else if _type == 2 {
		i := this.pool2.Get(addr)
		if i != nil {
			return i.(*InterflowBean)
		}
	}
	return nil
}

func (this *InterflowPool) del(addr string, ifb *InterflowBean, _type int) {
	//	this.lock.Lock()
	//	defer this.lock.Unlock()
	if _type == 1 {
		this.pool.Delnx(addr, ifb)
	} else if _type == 2 {
		this.pool2.Delnx(addr, ifb)
	}
}

type InterflowBean struct {
	Addr          string
	TimMBeanList  []*TimMBean
	TimPBeanList  []*TimPBean
	Lock          *sync.Mutex
	sendTimestamp int64
	sendFlag      chan int
	status        int
}

func NewInterflowBean(addr string, _type int) (ib *InterflowBean) {
	ib = &InterflowBean{Addr: addr, Lock: new(sync.Mutex), sendTimestamp: utils.Atoi64(utils.TimeMills()) + int64(ClusterConf.Interflow*1000)}
	if _type == 1 {
		go ib._sendMBean()
	}
	if _type == 2 {
		go ib._sendPBean()
	}
	return
}

func (this *InterflowBean) _sendMBean() {
	defer func() {
		if er := recover(); er != nil {
			logger.Error(string(debug.Stack()))
		}
	}()
	for {
		timer := time.NewTicker(time.Duration(ClusterConf.Interflow) * time.Millisecond)
		select {
		case <-timer.C:
			if this.status == 1 || this.TimMBeanList != nil {
				goto ENDPOOL
			}
		}
	}
ENDPOOL:
	this.status = 1
	go IFPool.del(this.Addr, this, 1)
	if this.TimMBeanList != nil {
		client := Pool.Get(this.Addr)
		timMBeanList := NewTimMBeanList()
		timMBeanList.TimMBeanList = this.TimMBeanList
		timMBeanList.ThreadId = utils.TimeMills()
		_, err := client.SendMBeanList(timMBeanList, NewAuth())
		if err != nil {
			daoService.SaveOfflineMBeanList(timMBeanList.GetTimMBeanList())
		} else {
			defer Pool.Put(this.Addr, client)
		}
	}
}

func (this *InterflowBean) _sendPBean() {
	defer func() {
		if er := recover(); er != nil {
			logger.Error(string(debug.Stack()))
		}
	}()
	for {
		timer := time.NewTicker(time.Duration(ClusterConf.Interflow) * time.Millisecond)
		select {
		case <-timer.C:
			if this.status == 1 || this.TimPBeanList != nil {
				goto ENDPOOL
			}
		}
	}
ENDPOOL:
	this.status = 1
	go IFPool.del(this.Addr, this, 2)
	if this.TimPBeanList != nil {
		client := Pool.Get(this.Addr)
		timPBeanList := NewTimPBeanList()
		timPBeanList.TimPBeanList = this.TimPBeanList
		timPBeanList.ThreadId = utils.TimeMills()
		_, err := client.SendPBeanList(timPBeanList, NewAuth())
		if err == nil {
			defer Pool.Put(this.Addr, client)
		}
	}
}

func (this *InterflowBean) SendMBean(mbean *TimMBean) (err error) {
	if this.status == 1 {
		err = errors.New("interflow SendMBean fail")
		return
	}
	this.Lock.Lock()
	defer this.Lock.Unlock()
	timestamp := utils.Atoi64(utils.TimeMills())
	if timestamp >= this.sendTimestamp || this.status == 1 {
		this.status = 1
		err = errors.New("interflow SendMBean fail")
		return
	} else {
		if this.TimMBeanList == nil {
			this.TimMBeanList = make([]*TimMBean, 0)
		}
		this.TimMBeanList = append(this.TimMBeanList, mbean)
	}
	return
}

func (this *InterflowBean) SendPBean(pbean *TimPBean) (err error) {
	if this.status == 1 {
		err = errors.New("interflow SendPBean fail")
		return
	}
	this.Lock.Lock()
	defer this.Lock.Unlock()
	timestamp := utils.Atoi64(utils.TimeMills())
	if timestamp >= this.sendTimestamp || this.status == 1 {
		this.status = 1
		err = errors.New("interflow SendPBean fail")
		return
	} else {
		if this.TimPBeanList == nil {
			this.TimPBeanList = make([]*TimPBean, 0)
		}
		this.TimPBeanList = append(this.TimPBeanList, pbean)
	}
	return
}

func InterflowSendMBean(addr string, mbean *TimMBean) (err error) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(string(debug.Stack()))
		}
	}()
	return _InterflowSendMBean(addr, mbean, 5)
}

func _InterflowSendMBean(addr string, mbean *TimMBean, count int) (err error) {
	if count <= 0 {
		return errors.New("InterflowSendMBean overtime")
	}
	count--
	ifb := IFPool.get(addr, 1)
	if ifb == nil {
		ifb = NewInterflowBean(addr, 1)
		IFPool.add(addr, ifb, 1)
	}
	err = ifb.SendMBean(mbean)
	if err != nil {
		return _InterflowSendMBean(addr, mbean, count)
	}
	return
}

func InterflowSendPBean(addr string, pbean *TimPBean) (err error) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(string(debug.Stack()))
		}
	}()
	return _InterflowSendPBean(addr, pbean, 5)
}

func _InterflowSendPBean(addr string, pbean *TimPBean, count int) (err error) {
	if count <= 0 {
		return errors.New("_InterflowSendPBean overtime")
	}
	count--
	ifb := IFPool.get(addr, 2)
	if ifb == nil {
		ifb = NewInterflowBean(addr, 2)
		IFPool.add(addr, ifb, 2)
	}
	err = ifb.SendPBean(pbean)
	if err != nil {
		return _InterflowSendPBean(addr, pbean, count)
	}
	return
}
