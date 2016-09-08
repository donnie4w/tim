/**
 * donnie4w@gmail.com  tim server
 */
package Map

import (
	"sync"
)

type HashTable struct {
	maptable map[interface{}]interface{}
	lock     *sync.RWMutex
}

func NewHashTable() *HashTable {
	return &HashTable{make(map[interface{}]interface{}, 0), new(sync.RWMutex)}
}

func (this *HashTable) Put(k, v interface{}) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.maptable[k] = v
}

func (this *HashTable) Get(k interface{}) (i interface{}) {
	this.lock.RLock()
	defer this.lock.RUnlock()
	i = this.maptable[k]
	return
}

/**获取所有key*/
func (this *HashTable) GetKeys() (keys []interface{}) {
	this.lock.RLock()
	defer this.lock.RUnlock()
	keys = make([]interface{}, 0)
	for k, _ := range this.maptable {
		keys = append(keys, k)
	}
	return
}

/**获取所有值*/
func (this *HashTable) GetValues() (values []interface{}) {
	this.lock.RLock()
	defer this.lock.RUnlock()
	values = make([]interface{}, 0)
	for _, v := range this.maptable {
		values = append(values, v)
	}
	return
}

/**k不存在时 设置v值*/
func (this *HashTable) Putnx(k, v interface{}) {
	this.lock.Lock()
	defer this.lock.Unlock()
	if _, ok := this.maptable[k]; !ok {
		this.maptable[k] = v
	}
}

func (this *HashTable) Del(k interface{}) {
	this.lock.Lock()
	defer this.lock.Unlock()
	delete(this.maptable, k)
}

/**k,v对应存在时删除*/
func (this *HashTable) Delnx(k, v interface{}) bool {
	this.lock.Lock()
	defer this.lock.Unlock()
	if value, ok := this.maptable[k]; ok {
		if value == v {
			delete(this.maptable, k)
			return true
		}
	}
	return false
}
