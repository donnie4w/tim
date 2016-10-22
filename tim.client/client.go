/**
 * donnie4w@gmail.com  tim server
 */
package client

import (
	"errors"
	"fmt"
	"os"
	"runtime/debug"
	"sync"
	"time"

	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/donnie4w/go-logger/logger"
	. "tim.protocol"
	"tim.route"
)

type FLOW string

const (
	START  FLOW = "start"
	AUTH   FLOW = "auth"
	NOAUTH FLOW = "noauth"
	CLOSE  FLOW = "close"

	CONNECT_START FLOW = "connect_start"
	CONNECT_RUN   FLOW = "connect_run"
	CONNECT_STOP  FLOW = "connect_stop"
)

type Config struct {
	Ip       string
	Port     string
	Domain   *string
	Resource *string
}

type Connect struct {
	Client      *ITimClient
	FlowConnect FLOW
	Super       *Cli
}

func (this *Connect) Close() {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("sendmsg,", err)
			logger.Error(string(debug.Stack()))
		}
	}()
	if this.Client != nil && this.Client.Transport != nil && this.FlowConnect != CONNECT_STOP {
		this.FlowConnect = CONNECT_STOP
		this.Client.Transport.Close()
	}
}

func (this *Connect) setITimClient(client *ITimClient) {
	this.Client = client
}

//--------------------------------------------------------

type Cli struct {
	Connect     *Connect
	Sync        *sync.Mutex
	Flow        FLOW
	Addr        string
	Name        string
	Pwd         string
	Domain      *string
	Resource    string
	FlowChan    chan int
	FLowValue   int
	ReConnLimit int
}

type ConfirmBean struct {
	isClusterAck   bool
	confirmMessage map[string]chan int
	confirmMBean   map[string]*TimMBean
	lock           *sync.RWMutex
}

var Confirm *ConfirmBean = &ConfirmBean{isClusterAck: false, confirmMessage: make(map[string]chan int, 0), confirmMBean: make(map[string]*TimMBean), lock: new(sync.RWMutex)}

func (this *ConfirmBean) Add(threadId string, mbean *TimMBean) {
	defer func() {
		if er := recover(); er != nil {
		}
	}()
	if !this.isClusterAck {
		return
	}
	this.lock.Lock()
	defer this.lock.Unlock()
	this.confirmMessage[threadId] = make(chan int, 0)
	this.confirmMBean[threadId] = mbean
	go this.monitorMbean(threadId)
}

func (this *ConfirmBean) monitorMbean(threadId string) {
	defer func() {
		if er := recover(); er != nil {
		}
	}()
	if !this.isClusterAck {
		return
	}
	this.lock.Lock()
	defer this.lock.Unlock()
	if chanvalue, ok := this.confirmMessage[threadId]; ok {
		timer := time.NewTicker(5 * time.Second)
		select {
		case <-timer.C:
			if mbean, ok := this.confirmMBean[threadId]; ok {
				route.RouteMBean(mbean, false, true)
			}
		case <-chanvalue:
		}
		this._Del(threadId)
	}
}

func (this *ConfirmBean) ConfirmThreadId(threadId string) {
	defer func() {
		if er := recover(); er != nil {
		}
	}()
	if !this.isClusterAck {
		return
	}
	this.lock.RLocker()
	defer this.lock.RUnlock()
	if chanvalue, ok := this.confirmMessage[threadId]; ok {
		timer := time.NewTicker(5 * time.Second)
		select {
		case <-timer.C:
		case chanvalue <- 1:
		}
	}
}

func (this *ConfirmBean) Del(threadId string) {
	defer func() {
		if er := recover(); er != nil {
		}
	}()
	if !this.isClusterAck {
		return
	}
	this.lock.Lock()
	defer this.lock.Unlock()
	this._Del(threadId)
}

func (this *ConfirmBean) _Del(threadId string) {
	defer func() {
		if er := recover(); er != nil {
		}
	}()
	if !this.isClusterAck {
		return
	}
	delete(this.confirmMessage, threadId)
	delete(this.confirmMBean, threadId)
}

func NewCli(addr, username, pwd string, domain *string) (Timc *Cli) {
	Timc = &Cli{Sync: new(sync.Mutex), Flow: AUTH, Addr: addr, Name: username, Pwd: pwd, FlowChan: make(chan int, 0), Domain: domain}
	return
}

func (this *Cli) SendMBean(mbean *TimMBean) (err error) {
	defer func() {
		if er := recover(); er != nil {
			err = errors.New(fmt.Sprint(er))
			logger.Error("SendMBean,", er)
			logger.Error(string(debug.Stack()))
		}
	}()
	this.Sync.Lock()
	defer this.Sync.Unlock()
	if mbean != nil && mbean.ThreadId != "" {
		if this.Flow == START {
			timer := time.NewTicker(5 * time.Second)
			select {
			case <-timer.C:
				err = errors.New("sendMbean error")
				return
			case this.FLowValue = <-this.FlowChan:
			}
		}
		if this.Flow != AUTH {
			err = errors.New("sendMbean error")
		} else {
			if this.FLowValue == 0 {
				this.FLowValue = 1
			}
		}
		if this.Addr != "" {
			err = this.Connect.Client.TimMessage(mbean)
			Confirm.Add(mbean.GetThreadId(), mbean)
		}
	} else {
		err = errors.New(fmt.Sprint("mbean error:", mbean))
	}
	return
}

func (this *Cli) SendPBean(pbean *TimPBean) (err error) {
	defer func() {
		if er := recover(); er != nil {
			err = errors.New(fmt.Sprint(er))
			logger.Error("SendPBean,", er)
			logger.Error(string(debug.Stack()))
		}
	}()
	this.Sync.Lock()
	defer this.Sync.Unlock()
	//	logger.Debug("SendPBean==========>", pbean)
	if pbean != nil && pbean.ThreadId != "" {
		if this.Flow == START {
			timer := time.NewTicker(5 * time.Second)
			select {
			case <-timer.C:
				err = errors.New("SendPBean error")
				return
			case this.FLowValue = <-this.FlowChan:
			}
		}
		if this.Flow != AUTH {
			err = errors.New("SendPBean error")
		} else {
			if this.FLowValue == 0 {
				this.FLowValue = 1
			}
		}
		if this.Addr != "" {
			err = this.Connect.Client.TimPresence(pbean)
			if err != nil {
				logger.Error("SendPBean error:", err.Error())
			}
		}
	} else {
		err = errors.New(fmt.Sprint("pbean error:", pbean))
	}
	return
}

func (this *Cli) DisConnect() {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("DisConnect,", err)
			logger.Error(string(debug.Stack()))
		}
	}()
	if this != nil && this.Connect != nil {
		this.Connect.Client.Transport.Close()
	}
}

func (this *Cli) Ack(ab *TimAckBean) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("Ack,", err)
			logger.Error(string(debug.Stack()))
		}
	}()
	this.Sync.Lock()
	defer this.Sync.Unlock()
	if this != nil && this.Connect != nil && this.Flow == AUTH {
		this.Connect.Client.TimAck(ab)
	}
}

func (this *Cli) Close() {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("Close,", err)
			logger.Error(string(debug.Stack()))
		}
	}()
	if this != nil && this.Connect != nil {
		this.Flow = CLOSE
		this.Connect.Client.Transport.Close()
	}
}

func (this *Cli) Ping() {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("Ping,", err)
			logger.Error(string(debug.Stack()))
			go this.Login()
		}
	}()
	this.Sync.Lock()
	defer this.Sync.Unlock()
	for {
		if this.ReConnLimit >= 10 {
			this.Connect.FlowConnect = CONNECT_STOP
			break
		}
		for i := 0; i < 20; i++ {
			time.Sleep(1 * time.Second)
		}
		if this == nil || this.Flow == CLOSE {
			logger.Debug("this.Flow====break")
			break
		}
		if this != nil && this.Flow == AUTH {
			if this.Connect == nil {
				continue
			}
			err := func() (err error) {
				defer func() {
					if err := recover(); err != nil {
						logger.Error(string(debug.Stack()))
					}
				}()
				logger.Debug("client ping>>>>>>>>>>>>>>>>>")
				this.Sync.Lock()
				defer this.Sync.Unlock()
				err = this.Connect.Client.TimPing(fmt.Sprint(currentTimeMillis()))
				logger.Debug("client ping<<<<<<<<<<<<<<<<<")
				return
			}()
			if err != nil {
				this.ReConnLimit++
				logger.Error("client ping err>>>>>>>>>>>>>>>>>", err.Error())
				go this.Login()
			}
		}
	}
}

func (this *Cli) Login() (err error) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("controllerHandler,", err)
			logger.Error(string(debug.Stack()))
			this.ReConnLimit++
		}
	}()
	logger.Debug("Login")
	transportFactory := thrift.NewTBufferedTransportFactory(1024)
	protocolFactory := thrift.NewTCompactProtocolFactory()
	transport, err := thrift.NewTSocket(this.Addr)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error resolving address:", err)
		return
	}
	useTransport := transportFactory.GetTransport(transport)
	timclient := NewITimClientFactory(useTransport, protocolFactory)
	if this.Connect != nil {
		this.Connect.Close()
	}
	this.Connect = &Connect{FlowConnect: CONNECT_START}
	this.Connect.setITimClient(timclient)
	this.Connect.Super = this
	if err = transport.Open(); err != nil {
		fmt.Fprintln(os.Stderr, "Error opening socket to 127.0.0.1:9090", " ", err)
	}
	processorchan := make(chan int)
	go this.Connect.processor(processorchan)
	<-processorchan
	tid := new(Tid)
	resource := "goclient"
	tid.Domain, tid.Resource, tid.Name = this.Domain, &resource, this.Name
	err = timclient.TimLogin(tid, this.Pwd)
	if err != nil {
		logger.Error("cluster login err", err)
		this.ReConnLimit++
	} else {
		this.ReConnLimit = 0
	}
	return err
}

func currentTimeMillis() int64 {
	return time.Now().UnixNano() / 1000000
}

func (this *Connect) processor(processorchan chan int) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("controllerHandler,", err)
			logger.Error(string(debug.Stack()))
		}
		if this != nil {
			this.FlowConnect = CONNECT_STOP
		}
	}()
	handler := new(TimImpl)
	//	handler.Client = this.Super
	processor := NewITimProcessor(handler)
	protocol := thrift.NewTCompactProtocol(this.Client.Transport)
	for {
		if this == nil || this.FlowConnect == CONNECT_STOP {
			break
		}
		if this.FlowConnect == CONNECT_START {
			this.FlowConnect = CONNECT_RUN
			processorchan <- 1
		}
		b, err := processor.Process(protocol, protocol)
		if err != nil && !b {
			logger.Error("cluster processor error:", err.Error())
			break
		}
	}
}
