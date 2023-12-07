// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim
//

package cache

import (
	"fmt"
	"testing"
)

func Test_cache(t *testing.T) {
	chatIdCache := NewLruPool[int, any](32, 100)

	for i := 0; i < 1<<10; i++ {
		chatIdCache.Put(i, nil)
	}

	for i := 0; i < 1<<10; i++ {
		if _, ok := chatIdCache.Get(i); !ok {
			fmt.Println("err>>", i)
		}
	}
}

func BenchmarkParallel_cache(b *testing.B) {
	chatIdCache := NewLruPool[int, any](32, 100)
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			i++
			chatIdCache.Put(i, nil)
		}
	})
}
