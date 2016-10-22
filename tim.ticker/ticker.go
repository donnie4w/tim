/**
 * donnie4w@gmail.com  tim server
 */
package ticker

import (
	"runtime/debug"
	"time"

	"github.com/donnie4w/go-logger/logger"
	. "tim.common"
	"tim.daoService"
)

func TickerStart() {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("tickerStart", err)
			logger.Error(string(debug.Stack()))
		}
	}()
	logger.Debug("tickerStart>>>>")
	go Ticker4Second(CF.GetConfLoad(600), daoService.AddConf)
}

func Ticker4Second(second int, function func()) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("Ticker4Second error :", err)
			logger.Error(string(debug.Stack()))
		}
	}()
	time.Sleep(time.Duration(second) * time.Second)
	timer := time.NewTicker(time.Duration(second) * time.Second)
	for {
		select {
		case <-timer.C:
			go function()
		}
	}
}
