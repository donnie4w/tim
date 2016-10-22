/**
 * donnie4w@gmail.com  tim server
 */
package hbase

import (
	"fmt"
	"testing"

	"github.com/donnie4w/go-logger/logger"
	"tim.utils"
)

func TestClient(t *testing.T) {
	Init()
	//testscan()
	//testInsert()
	//testdeleterow()
}

func testscan() {
	bean := new(Bean)
	bean.Family = "index"
	bean.Qualifier = "2690A7FC70FE18FF2637763D910DB840"
	beans := []*Bean{bean}
	results, err := ScansFromRow("tim_offline", beans, 0, true)
	if err == nil {
		for _, result := range results {
			tim_offline := new(Tim_offline)
			Result2object(result, tim_offline)
			logger.Debug("==========>", tim_offline.Mid, " ", tim_offline.Createtime, " ", tim_offline.Fromuser, " ")
		}
	} else {
		logger.Error("error:", err.Error())
	}
}

func testInsert() {
	fmt.Println("------------>testInsert")
	for i := 0; i < 1; i++ {
		tim_message := new(Tim_message)
		tim_message.Chatid = utils.TimeMills()
		tim_message.Createtime = utils.NowTime()
		tim_message.Stamp = tim_message.Createtime
		tim_message.Fromuser = fmt.Sprint("wuxiaodong_", i)
		tim_message.Gname = fmt.Sprint("wu_", i)
		tim_message.Large = "1"
		tim_message.Msgmode = "1"
		tim_message.Msgtype = "1"
		tim_message.Small = "0"
		tim_message.Stanza = fmt.Sprint("aaaaaaaaaaaaaaaaaa_", i)
		tim_message.Touser = fmt.Sprint("dong_", i)
		tim_message.IndexChatid = fmt.Sprint(tim_message.Chatid)
		row, err := tim_message.Insert()
		fmt.Println("timdomain=========>", row, err)
	}
}

func testdeleterow() {
	DeleteRow("tim_message", 1)
}

func _Benchmark_client(b *testing.B) {
	Init()
	b.SetBytes(1024 * 1024 * 100)
	b.SetParallelism(8)
	for i := 0; i < b.N; i++ {
		testscan()
	}
}

func _BenchmarkClientParallel(b *testing.B) {
	b.SetBytes(1024 * 1024 * 100)
	b.SetParallelism(8)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			testscan()
		}
	})
}
