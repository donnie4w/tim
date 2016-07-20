/**
 * donnie4w@gmail.com  tim server
 */
package main

import (
	"flag"
	"fmt"
	"os"

	gdao "github.com/donnie4w/gdao"
	"github.com/donnie4w/go-logger/logger"
	"tim.DB"
	"tim.cluster"
	. "tim.common"
	"tim.daoService"
	"tim.service"
	"tim.ticker"
)

func init() {
	fmt.Println("--------------------------------------------------------")
	fmt.Println("-------------------- tim1.0 server ---------------------")
	fmt.Println("--------------------------------------------------------")
	fmt.Println("------------------ donnie4w@gmail.com ------------------")
	fmt.Println("--------------------------------------------------------")
}

func initGdao() {
	logger.Debug("initGdao")
	DB.Init()
	gdao.SetDB(DB.Master)
	gdao.SetAdapterType(gdao.MYSQL)
	gbs, err := gdao.ExecuteQuery("select 1")
	if err == nil {
		logger.Debug("test db ok", gbs[0].MapIndex(1).Value())
	}
}

func initLog(loglevel string) {
	logger.SetConsole(true)
	logger.SetRollingDaily(ConfBean.GetLog())
	switch loglevel {
	case "debug":
		logger.SetLevel(logger.DEBUG)
	case "info":
		logger.SetLevel(logger.INFO)
	case "warn":
		logger.SetLevel(logger.WARN)
	case "error":
		logger.SetLevel(logger.ERROR)
	default:
		logger.SetLevel(logger.DEBUG)
	}
}

func main() {
	flag.Parse()
	wd, _ := os.Getwd()
	for i := 0; i < flag.NArg(); i++ {
		fmt.Println(flag.Arg(i))
	}
	if flag.NArg() > 2 {
		fmt.Println("error:", "flag's length is", flag.NArg())
		os.Exit(1)
	}
	if flag.NArg() >= 1 {
		ConfBean.Init(flag.Arg(0))
	} else {
		ConfBean.Init(fmt.Sprint(wd, "/tim.xml"))
	}
	if flag.NArg() == 2 {
		initLog(flag.Arg(1))
	} else {
		initLog("")
	}
	cluster.InitCluster(fmt.Sprint(wd, "/cluster.xml"))
	initGdao()
	daoService.AddConf()
	ticker.TickerStart()
	service.ServerStart()
}
