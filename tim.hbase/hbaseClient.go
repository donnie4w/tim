/**
 * donnie4w@gmail.com  tim server
 */
package hbase

import (
	"fmt"
	"os"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"

	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/donnie4w/go-logger/logger"
	. "tim.common"
	"tim.utils"
)

var maxOpenConns int = 100
var maxIdleConns int = 50
var minOpenConns int = 10
var timeoutConns int = 5  // second
var IdleTimeOut int = 180 // second

var pool chan *clientBean
var once sync.Once

func Init() {
	once.Do(func() {
		maxOpenConns, maxIdleConns, minOpenConns, timeoutConns, IdleTimeOut = CF.GetHbaseArgs(100, 50, 10, 5, 180)
		pool = make(chan *clientBean, maxOpenConns)
		initmin := minOpenConns - len(pool)
		if initmin > 0 {
			for i := 1; i <= initmin; i++ {
				go putPool(newClientBean())
			}
		}
		go ClientPool.monitor()
	})
}

func putPool(t *clientBean) {
	t.createtime = utils.TimeMillsInt64()
	pool <- t
}

type clientBean struct {
	tsclient   *THBaseServiceClient
	createtime int64
}

func newClientBean() (cb *clientBean) {
	i := 5
	for i > 0 {
		cb = new(clientBean)
		cb.tsclient = _NewHbaseClient()
		cb.createtime = utils.TimeMillsInt64()
		if cb.tsclient != nil {
			break
		} else {
			cb = nil
			i--
			time.Sleep(100 * time.Millisecond)
		}
	}
	return
}

type HbaseClientPool struct {
	lock         *sync.Mutex
	clientLength int32
}

var ClientPool *HbaseClientPool = &HbaseClientPool{lock: new(sync.Mutex), clientLength: 0}

func (this *HbaseClientPool) put(t *clientBean) {
	defer func() {
		if err := recover(); err != nil {
		}
	}()
	this.lock.Lock()
	defer this.lock.Unlock()
	if t != nil {
		if len(pool) >= maxOpenConns {
			closeTHBaseServiceClient(t)
			this.subLength()
		} else {
			go putPool(t)
		}
	}
}

func (this *HbaseClientPool) del(t *clientBean) {
	if t != nil {
		closeTHBaseServiceClient(t)
		this.subLength()
	}
}

func (this *HbaseClientPool) get() (client *clientBean) {
	defer func() {
		if err := recover(); err != nil {
		}
	}()
	this.lock.Lock()
	defer this.lock.Unlock()
	if len(pool) > 0 {
		client = <-pool
	} else if this.getLength() < int32(maxOpenConns) {
		client = newClientBean()
	} else {
		timer := time.NewTicker(5 * time.Second)
		select {
		case <-timer.C:
		case client = <-pool:
		}
	}
	return
}

func (this *HbaseClientPool) getLength() int32 {
	return this.clientLength
}

func (this *HbaseClientPool) addLength() {
	atomic.AddInt32(&this.clientLength, 1)
}

func (this *HbaseClientPool) subLength() {
	atomic.AddInt32(&this.clientLength, -1)
}

func closeTHBaseServiceClient(t *clientBean) {
	defer func() {
		if err := recover(); err != nil {
		}
	}()
	if t != nil {
		t.tsclient.Transport.Close()
	}
}

func _NewHbaseClient() (hbaseClient *THBaseServiceClient) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("ClusterClient,", err)
			logger.Error(string(debug.Stack()))
		}
	}()
	addr := CF.HbaseAddr
	if addr == "" {
		addr = "127.0.0.1:9090"
	}
	//	tsocket, er := thrift.NewTSocket(addr)
	tsocket, er := thrift.NewTSocketTimeout(addr, time.Duration(timeoutConns)*time.Second)
	tProtocol := thrift.NewTBinaryProtocol(tsocket, true, true)
	if er != nil {
		logger.Error(os.Stderr, "error resolving address:", er.Error())
		return
	}
	hbaseClient = NewTHBaseServiceClientProtocol(tsocket, tProtocol, tProtocol)
	if er = tsocket.Open(); er != nil {
		logger.Error(os.Stderr, "Error opening socket to ", addr, " ", er)
		return nil
	}
	ClientPool.addLength()
	return
}

func (this *HbaseClientPool) monitor() {
	go monitor(300, this._monitor)
	go monitor(IdleTimeOut, this._monitor4IdleTimeOut)
}

func monitor(second int, f func()) {
	timer := time.NewTicker(time.Duration(second) * time.Second)
	for {
		select {
		case <-timer.C:
			f()
		}
	}
}

func (this *HbaseClientPool) _monitor() {
	this.lock.Lock()
	defer this.lock.Unlock()
	if len(pool) > maxIdleConns {
		for i := 0; i < len(pool)-maxIdleConns; i++ {
			c := <-pool
			this.del(c)
		}
	}
}

func (this *HbaseClientPool) _monitor4IdleTimeOut() {
	this.lock.Lock()
	defer this.lock.Unlock()
	if len(pool) > 0 {
		for i := 0; i < len(pool); i++ {
			c := <-pool
			if (utils.TimeMillsInt64()/1000 - c.createtime/1000) > int64(IdleTimeOut) {
				go this.del(c)
			} else {
				pool <- c
			}
		}
	}
}

func PrintPoolInfo() (s string) {
	s = fmt.Sprint("pool length:", len(pool), " | ClientPool'clientLength:", ClientPool.clientLength)
	return
}
