/**
 * donnie4w@gmail.com  tim server
 */
package ticker

import (
	"runtime/debug"
	"time"

	"github.com/donnie4w/go-logger/logger"
	"tim.daoService"
	//	"tim.utils"
)

func TickerStart() {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("tickerStart", err)
			logger.Error(string(debug.Stack()))
		}
	}()
	logger.Debug("-------------tickerStart----------------")
	go Ticker4min(5, daoService.AddConf)
}

func Ticker4min(min int, function func()) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("Ticker4min error :", err)
			logger.Error(string(debug.Stack()))
		}
	}()
	time.Sleep(time.Duration(min) * time.Minute)
	timer := time.NewTicker(time.Duration(min) * time.Minute)
	for {
		select {
		case <-timer.C:
			go function()
		}
	}
}
