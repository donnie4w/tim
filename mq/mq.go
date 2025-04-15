// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package mq

import (
	"errors"
	"github.com/donnie4w/gofer/hashmap"
	goutil "github.com/donnie4w/gofer/util"
	"github.com/donnie4w/tim/stub"
	"github.com/donnie4w/tim/sys"
	"github.com/donnie4w/tim/util"
	"sync"
)

var mq = newMqBean()

func Sub(topicType TopicType, id int64, f func(any)) {
	mq.Sub(topicType, id, f)
}

func Unsub(topicType TopicType, id int64) {
	mq.UnSub(topicType, id)
}

func Push(topicType TopicType, msg any) {
	defer util.Recover()
	mq.Push(topicType, msg)
}

type mqBean struct {
	topicmap *hashmap.MapL[TopicType, *hashmap.Map[int64, func(any)]]
	mux      *sync.RWMutex
}

func newMqBean() *mqBean {
	return &mqBean{topicmap: hashmap.NewMapL[TopicType, *hashmap.Map[int64, func(any)]](), mux: new(sync.RWMutex)}
}

func (m *mqBean) Sub(topicType TopicType, id int64, f func(any)) error {
	if !m.check(topicType) || id == 0 || f == nil {
		return errors.New("invalid param")
	}
	var cm *hashmap.Map[int64, func(any)]
	if hm, b := m.topicmap.Get(topicType); b {
		cm = hm
	} else {
		m.mux.Lock()
		defer m.mux.Unlock()
		if hm, b := m.topicmap.Get(topicType); b {
			cm = hm
		} else {
			cm = hashmap.NewMap[int64, func(any)]()
			m.topicmap.Put(topicType, cm)
		}
	}
	cm.Put(id, f)
	return nil
}

func (m *mqBean) check(topicType TopicType) bool {
	switch topicType {
	case ONLINESTATUS:
		return true
	}
	return false
}

func (m *mqBean) UnSub(topicType TopicType, id int64) {
	if hm, b := m.topicmap.Get(topicType); b {
		hm.Del(id)
	}
}

func (m *mqBean) Len() int64 {
	return m.topicmap.Len()
}

func (m *mqBean) Push(topicType TopicType, bean any) {
	if m.Len() > 0 {
		if hm, b := m.topicmap.Get(topicType); b {
			hm.Range(func(_ int64, v func(any)) bool {
				v(bean)
				return true
			})
		}
	}
}

type TopicType int8

var ONLINESTATUS TopicType = 1

func PushOnline(node string, on bool) {
	ab := stub.NewAdmSubBean()
	st := int8(ONLINESTATUS)
	ab.SubType = &st
	asob := stub.NewAdmSubOnlineBean()
	asob.Node = &node
	if on {
		asob.Status = &sys.ONLINE
	} else {
		asob.Status = &sys.OFFLIINE
	}
	ab.Bs = goutil.TEncode(asob)
	Push(ONLINESTATUS, ab)
}
