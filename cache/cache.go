// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim
package cache

//import (
//	. "github.com/donnie4w/gofer/cache"
//)
//
//type LruPool[K int | int64 | uint64, V any] struct {
//	lruList []*LruCache[V]
//	router  []int
//}
//
//func NewLruPool[K int | int64 | uint64, V any](poolsize int, LruLimit int) *LruPool[K, V] {
//	p := &LruPool[K, V]{}
//	p.lruList = make([]*LruCache[V], poolsize)
//	p.router = make([]int, poolsize)
//	for i := 0; i < poolsize; i++ {
//		p.lruList[i] = NewLruCache[V](LruLimit)
//		p.router[i] = i
//	}
//	return p
//}
//
//func (this *LruPool[K, V]) Put(t K, v V) {
//	this.lruList[int(t)%len(this.lruList)].Add(t, v)
//}
//
//func (this *LruPool[K, V]) Get(t K) (v V, b bool) {
//	return this.lruList[int(t)%len(this.lruList)].Get(t)
//}
//
//func (this *LruPool[K, V]) Del(t K) {
//	this.lruList[int(t)%len(this.lruList)].Remove(t)
//}
