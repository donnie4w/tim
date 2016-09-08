package clusterClient

import (
	"os"
	"runtime/debug"
	"sync"

	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/donnie4w/go-logger/logger"
	. "tim.protocol"
)

type ClusterClient struct {
	timclient *ITimClient
	lock      *sync.RWMutex
	Weight    int
}

func (this *ClusterClient) Close() {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("Close,", err)
			logger.Error(string(debug.Stack()))
		}
	}()
	this.Weight = 0
	if this.timclient != nil {
		this.timclient.Transport.Flush()
		this.timclient.Transport.Close()
	}
}

func (this *ClusterClient) SendMBean(mbean *TimMBean, auth *TimAuth) (r *TimResponseBean, err error) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("SendMBean,", err)
			logger.Error(string(debug.Stack()))
		}
	}()
	this.lock.Lock()
	defer this.lock.Unlock()
	r, err = this.timclient.TimResponseMessage(mbean, auth)
	return
}

func (this *ClusterClient) SendMBeanList(mbeanList *TimMBeanList, auth *TimAuth) (r *TimResponseBean, err error) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("SendMBeanList,", err)
			logger.Error(string(debug.Stack()))
		}
	}()
	this.lock.Lock()
	defer this.lock.Unlock()
	r, err = this.timclient.TimResponseMessageList(mbeanList, auth)
	return
}

func (this *ClusterClient) SendPBean(pbean *TimPBean, auth *TimAuth) (r *TimResponseBean, err error) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("SendPBean error:", err)
			logger.Error(string(debug.Stack()))
		}
	}()
	this.lock.Lock()
	defer this.lock.Unlock()
	r, err = this.timclient.TimResponsePresence(pbean, auth)
	return
}

func (this *ClusterClient) SendPBeanList(pbeanList *TimPBeanList, auth *TimAuth) (r *TimResponseBean, err error) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("SendPBeanList error:", err)
			logger.Error(string(debug.Stack()))
		}
	}()
	this.lock.Lock()
	defer this.lock.Unlock()
	r, err = this.timclient.TimResponsePresenceList(pbeanList, auth)
	return
}

func NewClusterClient(addr string) (clusterClient *ClusterClient, err error) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("ClusterClient,", err)
			logger.Error(string(debug.Stack()))
		}
	}()
	transportFactory := thrift.NewTBufferedTransportFactory(1024)
	protocolFactory := thrift.NewTCompactProtocolFactory()
	transport, er := thrift.NewTSocket(addr)
	if er != nil {
		logger.Error(os.Stderr, "error resolving address:", err)
		err = er
		return
	}
	useTransport := transportFactory.GetTransport(transport)
	timclient := NewITimClientFactory(useTransport, protocolFactory)
	if er = transport.Open(); er != nil {
		logger.Error(os.Stderr, "Error opening socket to ", addr, " ", er)
		err = er
		return
	}
	clusterClient = new(ClusterClient)
	clusterClient.timclient = timclient
	clusterClient.lock = new(sync.RWMutex)
	clusterClient.Weight = 1
	return
}
