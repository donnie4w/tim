// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package amr

import (
	"context"
	"errors"
	"github.com/donnie4w/tim/log"
	"github.com/donnie4w/tim/sys"
	"github.com/redis/go-redis/v9"
	"time"
)

type redisAmr struct {
	client redis.Cmdable
	ctx    context.Context
}

func newRedisAmr() amrStore {
	if sys.Conf.Redis == nil {
		log.FmtPrint("Redis configuration is not initialized")
		return nil
	}

	var client redis.Cmdable
	if len(sys.Conf.Redis.Addr) == 1 {
		opt := &redis.Options{
			Addr:     sys.Conf.Redis.Addr[0],
			Username: sys.Conf.Redis.Username,
			Password: sys.Conf.Redis.Password,
			DB:       sys.Conf.Redis.DB,
			Protocol: sys.Conf.Redis.Protocol,
		}
		client = redis.NewClient(opt)
	} else {
		opt := &redis.ClusterOptions{
			Addrs:    sys.Conf.Redis.Addr,
			Username: sys.Conf.Redis.Username,
			Password: sys.Conf.Redis.Password,
			Protocol: sys.Conf.Redis.Protocol,
		}
		client = redis.NewClusterClient(opt)
	}

	if err := client.Ping(context.TODO()).Err(); err != nil {
		log.FmtPrint("redis connect fail:", err)
		return nil
	}
	log.FmtPrint("Redis connection established")
	return &redisAmr{
		client: client,
		ctx:    context.Background(),
	}
}

func (s *redisAmr) put(atype AMRTYPE, key, value []byte, ttl uint64) error {
	return s.client.Set(s.ctx, string(amrKey(atype, key)), value, time.Duration(ttl)*time.Second).Err()
}

func (s *redisAmr) get(atype AMRTYPE, key []byte) ([]byte, error) {
	return s.client.Get(s.ctx, string(amrKey(atype, key))).Bytes()
}

func (s *redisAmr) remove(atype AMRTYPE, key []byte) error {
	return s.client.Del(s.ctx, string(amrKey(atype, key))).Err()
}

func (s *redisAmr) append(atype AMRTYPE, key, value []byte, ttl uint64) (err error) {
	fullKey := string(amrKey(atype, key))
	if err = s.client.SAdd(s.ctx, fullKey, string(value)).Err(); err == nil && ttl > 0 {
		err = s.client.Expire(s.ctx, fullKey, time.Duration(ttl)*time.Second).Err()
	}
	return
}

func (s *redisAmr) getMutil(atype AMRTYPE, key []byte) [][]byte {
	if members, err := s.client.SMembers(s.ctx, string(amrKey(atype, key))).Result(); err == nil && len(members) > 0 {
		results := make([][]byte, len(members))
		for i := range members {
			results[i] = []byte(members[i])
		}
		return results
	}
	return nil
}

func (s *redisAmr) removeKV(atype AMRTYPE, key, value []byte) error {
	return s.client.SRem(s.ctx, string(amrKey(atype, key)), value).Err()
}

func (s *redisAmr) close() error {
	switch v := s.client.(type) {
	case *redis.Client:
		return v.Close()
	case *redis.ClusterClient:
		return v.Close()
	default:
		return errors.New("unknown client type")
	}
}
