// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package amr

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"github.com/donnie4w/tim/log"
	"github.com/donnie4w/tim/sys"
	etcd "go.etcd.io/etcd/client/v3"
	"os"
	"time"
)

type etcdAmr struct {
	client *etcd.Client
	ctx    context.Context
}

func newEtcdAmr() amrStore {
	if sys.Conf.Etcd == nil {
		log.FmtPrint("Etcd configuration is not initialized")
		return nil
	}

	cfg := etcd.Config{
		Endpoints:   sys.Conf.Etcd.Endpoints,
		DialTimeout: time.Second * 5,
	}

	if sys.Conf.Etcd.Username != "" && sys.Conf.Etcd.Password != "" {
		cfg.Username = sys.Conf.Etcd.Username
		cfg.Password = sys.Conf.Etcd.Password
	}

	tlsConfig := &tls.Config{}
	if sys.Conf.Etcd.CAFile != "" || sys.Conf.Etcd.CertFile != "" || sys.Conf.Etcd.KeyFile != "" {
		if sys.Conf.Etcd.CAFile != "" {
			caCert, err := os.ReadFile(sys.Conf.Etcd.CAFile)
			if err != nil {
				log.FmtPrint("etcd:failed to read CA file:", err)
				return nil
			}
			caCertPool := x509.NewCertPool()
			if ok := caCertPool.AppendCertsFromPEM(caCert); !ok {
				log.FmtPrint("etcd:failed to append CA certs")
				return nil
			}
			tlsConfig.RootCAs = caCertPool
		}

		if sys.Conf.Etcd.CertFile != "" && sys.Conf.Etcd.KeyFile != "" {
			cert, err := tls.LoadX509KeyPair(sys.Conf.Etcd.CertFile, sys.Conf.Etcd.KeyFile)
			if err != nil {
				log.FmtPrint("etcd:failed to load client cert and key: ", err)
				return nil
			}
			tlsConfig.Certificates = []tls.Certificate{cert}
		}

		tlsConfig.InsecureSkipVerify = false
		cfg.TLS = tlsConfig
	}

	client, err := etcd.New(cfg)
	if err != nil {
		log.FmtPrint("failed to connect to etcd:", err)
		return nil
	}

	if _, err = client.Status(context.Background(), sys.Conf.Etcd.Endpoints[0]); err != nil {
		log.FmtPrint("failed to connect to etcd: ", err)
		return nil
	}
	log.FmtPrint("Etcd connection established")
	return &etcdAmr{
		client: client,
		ctx:    context.Background(),
	}
}

func (s *etcdAmr) put(atype AMRTYPE, key, value []byte, ttl uint64) error {
	resp, err := s.client.Grant(s.ctx, int64(ttl))
	if err != nil {
		return err
	}
	_, err = s.client.Put(s.ctx, string(amrKey(atype, key)), string(value), etcd.WithLease(resp.ID))
	return err
}

func (s *etcdAmr) get(atype AMRTYPE, key []byte) ([]byte, error) {
	resp, err := s.client.Get(s.ctx, string(amrKey(atype, key)))
	if err != nil {
		return nil, err
	}
	if len(resp.Kvs) == 0 {
		return nil, nil
	}
	return resp.Kvs[0].Value, nil
}

func (s *etcdAmr) remove(atype AMRTYPE, key []byte) error {
	_, err := s.client.Delete(s.ctx, string(amrKey(atype, key)))
	return err
}

func (s *etcdAmr) append(atype AMRTYPE, key, value []byte, ttl uint64) error {
	setKey := string(amrKey(atype, key)) + "/" + string(value)
	var leaseID etcd.LeaseID
	if ttl > 0 {
		resp, err := s.client.Grant(s.ctx, int64(ttl))
		if err != nil {
			return err
		}
		leaseID = resp.ID
	}
	_, err := s.client.Put(s.ctx, setKey, "", etcd.WithLease(leaseID))
	return err
}

func (s *etcdAmr) getMutil(atype AMRTYPE, key []byte) [][]byte {
	prefix := string(amrKey(atype, key)) + "/"
	resp, err := s.client.Get(s.ctx, prefix, etcd.WithPrefix())
	if err != nil {
		return nil
	}
	results := make([][]byte, 0, len(resp.Kvs))
	for _, kv := range resp.Kvs {
		results = append(results, kv.Key[len(prefix):])
	}
	return results
}

func (s *etcdAmr) removeKV(atype AMRTYPE, key, value []byte) error {
	setKey := string(amrKey(atype, key)) + "/" + string(value)
	_, err := s.client.Delete(s.ctx, setKey)
	return err
}

func (s *etcdAmr) close() error {
	return s.client.Close()
}
