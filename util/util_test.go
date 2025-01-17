// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package util

import (
	"testing"
)

func TestMarkId(t *testing.T) {
	var i int64 = 1 << 50
	id := MaskId(i)
	id2 := MaskId(id)
	t.Log(i)
	t.Log(id)
	t.Log(id2)
}

func TestMark(t *testing.T) {
	bs := []byte("hello world")
	bs1 := Mask(bs)
	bs2 := Mask(bs1)
	t.Log(string(bs1))
	t.Log(string(bs2))
}

func BenchmarkNodeName(b *testing.B) {
	domain := "tt"
	u := CreateUUID("aiaeinf22ienfefne1f", &domain)
	b.Log(u)
	b.Log(CheckUUID(2790553438565061983))
}

func BenchmarkUUIDByNode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		domain := "tt"
		CreateUUID("aiaeinfienfefne1f", &domain)
	}
}

func BenchmarkSearchString(b *testing.B) {
	b.Log(ContainStrings([]string{"ab", "b", "c"}, "ab"))
}

func BenchmarkSearchInt(b *testing.B) {
	b.Log(ContainInt([]int{11, 22, 33, 44}, 33))
}
