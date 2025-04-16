package amr

import (
	"github.com/donnie4w/gofer/util"
	"github.com/donnie4w/tim/cache"
	"github.com/donnie4w/tim/stub"
	"github.com/donnie4w/tim/sys"
)

func PutToken(token string, tid *stub.Tid) {
	if islocalamr {
		cache.TokenCache.Put(token, tid)
	} else {
		amr.put(TOKEN, []byte(token), util.TEncode(tid), uint64(sys.Conf.TokenTimeout))
	}
}

func GetToken(token string) *stub.Tid {
	if cache.TokenUsedCache.Contains([]byte(token)) {
		return nil
	}
	if islocalamr {
		return cache.TokenCache.Get(token)
	} else {
		if bs, _ := amr.get(TOKEN, []byte(token)); len(bs) > 0 {
			r, _ := util.TDecode[*stub.Tid](bs, stub.NewTid())
			return r
		}
	}
	return nil
}

func DelToken(token string) {
	if islocalamr {
		cache.TokenCache.Del(token)
	} else {
		amr.remove(TOKEN, []byte(token))
	}
	cache.TokenUsedCache.Add([]byte(token))
}
