/**
 * donnie4w@gmail.com  tim server
 */
package cluster

import (
	"sync"
	"time"

	. "tim.Map"
	. "tim.clusterClient"
)

var MAX_POOL_SIZE int = 100

type ClusterPool struct {
	pool *HashTable
	lock *sync.RWMutex
}

var Pool *ClusterPool = &ClusterPool{pool: NewHashTable(), lock: new(sync.RWMutex)}

func (this *ClusterPool) Get(addr string) (cc *ClusterClient) {
	this.lock.RLock()
	defer this.lock.RUnlock()
	chancc := this.pool.Get(addr)
	if chancc != nil {
		timer := time.NewTicker(2 * time.Second)
		select {
		case <-timer.C:
		case cc = <-chancc.(chan *ClusterClient):
		}
	}
	if cc == nil {
		client, err := NewClusterClient(addr)
		if err == nil {
			cc = client
		}
	}
	return
}

func (this *ClusterPool) Put(addr string, cc *ClusterClient) {
	this.lock.Lock()
	defer this.lock.Unlock()
	chancc := this.pool.Get(addr)
	if chancc != nil {
		go pushChan(chancc.(chan *ClusterClient), cc)
	} else {
		cpchan := make(chan *ClusterClient, MAX_POOL_SIZE)
		this.pool.Put(addr, cpchan)
		go pushChan(cpchan, cc)
	}
}

func pushChan(cpchan chan *ClusterClient, cc *ClusterClient) {
	if len(cpchan) > MAX_POOL_SIZE {
		cc.Close()
	} else {
		cpchan <- cc
	}
}
