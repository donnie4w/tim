// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package amr

import (
	"errors"
	"fmt"
	"github.com/donnie4w/tim/log"
	"github.com/donnie4w/tim/sys"
	"github.com/go-zookeeper/zk"
	"time"
)

type zkAmr struct {
	conn *zk.Conn
}

func newZkAmr() amrStore {
	if sys.Conf.ZooKeeper == nil {
		log.FmtPrint("Zookeeper configuration is not initialized")
		return nil
	}

	zkConn, _, err := zk.Connect(sys.Conf.ZooKeeper.Servers, time.Duration(sys.Conf.ZooKeeper.SessionTimeout)*time.Second)
	if err != nil {
		log.FmtPrint("failed to connect to Zookeeper: ", err)
		return nil
	}

	if sys.Conf.ZooKeeper.Username != "" && sys.Conf.ZooKeeper.Password != "" {
		auth := fmt.Sprintf("%s:%s", sys.Conf.ZooKeeper.Username, sys.Conf.ZooKeeper.Password)
		if err = zkConn.AddAuth("digest", []byte(auth)); err != nil {
			log.FmtPrint("failed to authenticate with Zookeeper: ", err)
			return nil
		}
	}
	log.FmtPrint("Zookeeper connection established")
	return &zkAmr{
		conn: zkConn,
	}
}

func (s *zkAmr) put(atype AMRTYPE, key, value []byte, ttl uint64) (err error) {
	fullKey := string(amrKey(atype, key))
	_, err = s.conn.Set(fullKey, value, -1)
	if err != nil && errors.Is(err, zk.ErrNoNode) {
		_, err = s.conn.Create(fullKey, value, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	}
	return
}
func (s *zkAmr) get(atype AMRTYPE, key []byte) ([]byte, error) {
	data, _, err := s.conn.Get(string(amrKey(atype, key)))
	return data, err
}

func (s *zkAmr) remove(atype AMRTYPE, key []byte) error {
	return s.conn.Delete(string(amrKey(atype, key)), -1)
}

func (s *zkAmr) append(atype AMRTYPE, key, value []byte, ttl uint64) error {
	_, err := s.conn.Create(string(amrKey(atype, key))+"/"+string(value), nil, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	return err
}

func (s *zkAmr) getMutil(atype AMRTYPE, key []byte) [][]byte {
	if children, _, err := s.conn.Children(string(amrKey(atype, key))); err == nil {
		results := make([][]byte, len(children))
		for i := range children {
			results[i] = []byte(children[i])
		}
		return results
	}
	return nil
}
func (s *zkAmr) removeKV(atype AMRTYPE, key, value []byte) error {
	return s.conn.Delete(string(amrKey(atype, key))+"/"+string(value), -1)
}

func (s *zkAmr) close() error {
	s.conn.Close()
	return nil
}
