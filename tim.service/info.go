/**
 * donnie4w@gmail.com  tim server
 */
package service

import (
	"fmt"
	"io"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/donnie4w/go-logger/logger"
	. "tim.connect"
	"tim.hbase"
)

func info(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(debug.Stack())
		}
	}()
	logger.Debug("RemoteAddr:", r.RemoteAddr)
	logger.Debug("X-Forwarded-For:", r.Header.Get("X-Forwarded-For"))
	logger.Debug("ContentLength:", r.ContentLength)
	X_Forwarded_For := r.Header.Get("X-Forwarded-For")
	ss := strings.Split(r.RemoteAddr, ":")
	ipaddr := ss[0]
	if X_Forwarded_For != "" && X_Forwarded_For != "127.0.0.1" {
		ipaddr = X_Forwarded_For
	}
	logger.Debug("ip:", ipaddr)
	if r.ContentLength >= 2*1024*1024 {
		return
	}
	str := fmt.Sprintln("user:", TP.Len4P(), "===>", TP.Len4PU())
	io.WriteString(w, str)
}

func userInfo(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(debug.Stack())
		}
	}()
	logger.Debug("RemoteAddr:", r.RemoteAddr)
	logger.Debug("X-Forwarded-For:", r.Header.Get("X-Forwarded-For"))
	logger.Debug("ContentLength:", r.ContentLength)
	X_Forwarded_For := r.Header.Get("X-Forwarded-For")
	ss := strings.Split(r.RemoteAddr, ":")
	ipaddr := ss[0]
	if X_Forwarded_For != "" && X_Forwarded_For != "127.0.0.1" {
		ipaddr = X_Forwarded_For
	}
	logger.Debug("ip:", ipaddr)
	if r.ContentLength >= 2*1024*1024 {
		return
	}
	io.WriteString(w, TP.PrintUsersInfo())
}

func hbaseclient(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(debug.Stack())
		}
	}()
	logger.Debug("RemoteAddr:", r.RemoteAddr)
	logger.Debug("X-Forwarded-For:", r.Header.Get("X-Forwarded-For"))
	logger.Debug("ContentLength:", r.ContentLength)
	X_Forwarded_For := r.Header.Get("X-Forwarded-For")
	ss := strings.Split(r.RemoteAddr, ":")
	ipaddr := ss[0]
	if X_Forwarded_For != "" && X_Forwarded_For != "127.0.0.1" {
		ipaddr = X_Forwarded_For
	}
	logger.Debug("ip:", ipaddr)
	if r.ContentLength >= 2*1024*1024 {
		return
	}
	io.WriteString(w, hbase.PrintPoolInfo())
}
