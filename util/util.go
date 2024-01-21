// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim
//

package util

import (
	"bytes"
	"crypto/tls"
	"errors"
	"io"
	"net/http"
	"runtime"
	"runtime/debug"
	"sort"
	"strings"

	"github.com/donnie4w/gofer/buffer"
	"github.com/donnie4w/gofer/keystore"
	"github.com/donnie4w/gofer/pool/gopool"
	. "github.com/donnie4w/gofer/util"
	"github.com/donnie4w/simplelog/logging"
	. "github.com/donnie4w/tim/stub"
	"github.com/donnie4w/tim/sys"
)

var GoPool = gopool.NewPool(100, 100<<3)
var GoPoolTx = gopool.NewPool(int64(runtime.NumCPU()), int64(runtime.NumCPU())<<2)
var GoPoolTx2 = gopool.NewPool(int64(runtime.NumCPU()), int64(runtime.NumCPU())<<3)

func CreateUUIDByTid(tid *Tid) uint64 {
	return CreateUUID(tid.Node, tid.Domain)
}

func CreateUUID(node string, domain *string) uint64 {
	buf := buffer.NewBufferByPool()
	defer buf.Free()
	buf.WriteString(node)
	buf.WriteString(sys.Conf.Salt)
	if domain != nil && *domain != "" {
		buf.WriteString(*domain)
	}
	u := CRC64(buf.Bytes())
	bs := Int64ToBytes(int64(u))
	_bs := MaskWithSeed(bs, Mask(sys.MaskSeed))
	b8 := CRC8(_bs[:7])
	bs[7] = b8
	return uint64(BytesToInt64(bs))
}

func NewTimUUID() uint64 {
	buf := buffer.NewBuffer()
	buf.WriteString(sys.Conf.Salt)
	buf.Write(Int64ToBytes(RandId()))
	return CreateUUID(string(buf.Bytes()), nil)
}

func NameToNode(name string, domain *string) string {
	return UUIDToNode(CreateUUID(name, domain))
}

func NodeToUUID(node string) (_r uint64) {
	_r, _ = Base58DecodeForInt64([]byte(node))
	return
}

func UUIDToNode(uuid uint64) string {
	return string(Base58EncodeForInt64(uuid))
}

func CheckUUID(uuid uint64) bool {
	bs := Int64ToBytes(int64(uuid))
	_bs := MaskWithSeed(bs, Mask(sys.MaskSeed))
	b8 := CRC8(_bs[:7])
	return b8 == bs[7]
}

func CheckNode(node string) bool {
	if len(node) <= sys.NodeMaxlength {
		if _r := NodeToUUID(node); _r > 0 {
			return CheckUUID(_r)
		}
	}
	return false
}

func ChatIdByRoom(node string, domain *string) uint64 {
	buf := buffer.NewBufferByPool()
	defer buf.Free()
	buf.Write(Int64ToBytes(int64(CreateUUID(node, domain))))
	buf.WriteString(sys.Conf.Salt)
	return CRC64(buf.Bytes())
}

func ChatIdByNode(fromNode, toNode string, domain *string) uint64 {
	buf := buffer.NewBufferByPool()
	defer buf.Free()
	f, t := CreateUUID(fromNode, domain), CreateUUID(toNode, domain)
	if f < t {
		f, t = t, f
	}
	buf.Write(Int64ToBytes(int64(f)))
	buf.Write(Int64ToBytes(int64(t)))
	buf.WriteString(sys.Conf.Salt)
	return CRC64(buf.Bytes())
}

func RelateIdForGroup(groupNode, userNode string, domain *string) uint64 {
	buf := buffer.NewBufferByPool()
	defer buf.Free()
	f, t := CreateUUID(groupNode, domain), CreateUUID(userNode, domain)
	buf.Write(Int64ToBytes(int64(f)))
	buf.WriteString(sys.Conf.Salt)
	buf.Write(Int64ToBytes(int64(t)))
	buf.WriteString(sys.Conf.Salt)
	return CRC64(buf.Bytes())
}

func UnikId(f, t uint64) uint64 {
	buf := buffer.NewBufferByPool()
	defer buf.Free()
	buf.Write(Int64ToBytes(int64(f)))
	buf.WriteString(sys.Conf.Salt)
	buf.Write(Int64ToBytes(int64(t)))
	buf.WriteString(sys.Conf.Salt)
	return CRC64(buf.Bytes())
}

/***********************************************************/

func MaskId(id int64) (_r int64) {
	ids := Int64ToBytes(id)
	return BytesToInt64(Mask(ids))
}

func Mask(bs []byte) (_r []byte) {
	if bs != nil {
		_r = make([]byte, len(bs))
		for i, j := 0, 0; i < len(bs); i++ {
			_r[i] = bs[i] ^ sys.MaskSeed[j]
			j = i % 8
		}
	}
	return
}

func MaskStr(s string) (_r string) {
	return string(Mask([]byte(s)))
}

func MaskTid(tid *Tid) {
	if tid != nil {
		tid.Node = MaskStr(tid.Node)
	}
}

func MaskWithSeed(bs []byte, seed []byte) (_r []byte) {
	_r = make([]byte, len(bs))
	for i, j := 0, 0; i < len(bs); i++ {
		_r[i] = bs[i] ^ seed[j]
		j = i % len(seed)
	}
	return
}

func ParseAddr(addr string) (_r string, err error) {
	if _r = addr; !strings.Contains(_r, ":") {
		if MatchString("^[0-9]{4,5}$", addr) {
			_r = ":" + _r
		} else {
			err = errors.New("error format :" + addr)
		}
	}
	return
}

func JTP(b byte) bool {
	return b&0x80 == 0x80
}

func ArraySub2[K int | int8 | int32 | int64 | string](a []K, k K) (_r []K) {
	if a != nil {
		_r = make([]K, 0)
		for _, v := range a {
			if v != k {
				_r = append(_r, v)
			}
		}
	}
	return
}

func ArraySub[K int | int8 | int32 | int64 | string](a1, a2 []K) (_r []K) {
	_r = make([]K, 0)
	if a1 != nil && a2 != nil {
		m := make(map[K]byte, 0)
		for _, a := range a2 {
			m[a] = 0
		}
		for _, a := range a1 {
			if _, ok := m[a]; !ok {
				_r = append(_r, a)
			}
		}
	} else if a2 == nil {
		return a1
	}
	return
}

func Recover() {
	if err := recover(); err != nil {
		logging.Error(string(debug.Stack()))
	}
}

func AesEncode(bs []byte) ([]byte, error) {
	return keystore.AesEncrypter(bs, sys.Conf.EncryptKey)
}

func AesDecode(bs []byte) ([]byte, error) {
	return keystore.AesDecrypter(bs, sys.Conf.EncryptKey)
}

func ContainStrings(li []string, v string) (b bool) {
	if li == nil {
		return false
	}
	sort.Strings(li)
	idx := sort.SearchStrings(li, v)
	if idx < len(li) {
		b = li[idx] == v
	}
	return
}

func ContainInt[T int64 | uint64 | int | uint | uint32 | int32](li []T, v T) (b bool) {
	if li == nil {
		return false
	}
	sort.Slice(li, func(i, j int) bool { return li[i] < li[j] })
	idx := sort.Search(len(li), func(i int) bool { return li[i] >= v })
	if idx < len(li) {
		b = li[idx] == v
	}
	return
}

func HttpPost(bs []byte, close bool, httpurl string) (_r []byte, err error) {
	tr := &http.Transport{DisableKeepAlives: true}
	if strings.HasPrefix(httpurl, "https:") {
		tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	client := http.Client{Transport: tr}
	var req *http.Request
	if req, err = http.NewRequest(http.MethodPost, httpurl, bytes.NewReader(bs)); err == nil {
		if close {
			req.Close = true
		}
		var resp *http.Response
		if resp, err = client.Do(req); err == nil {
			if close {
				defer resp.Body.Close()
			}
			var body []byte
			if body, err = io.ReadAll(resp.Body); err == nil {
				_r = body
			}
		}
	}
	return
}
