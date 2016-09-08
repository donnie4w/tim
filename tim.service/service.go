/**
 * donnie4w@gmail.com  tim server
 */
package service

import (
	"fmt"
	"net"

	"github.com/donnie4w/go-logger/logger"
	. "tim.common"
)

func Start() {
	defer func() {
		if err := recover(); err != nil {
		}
	}()
	service := fmt.Sprint(CF.Addr, ":", CF.Port)
	logger.Debug("listen port:", CF.Port)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	for {
		conn, err := listener.Accept()
		logger.Debug(conn.RemoteAddr().String(), "is linking")
		if err == nil {
			go handler(conn)
		}
	}
}

func handler(conn net.Conn) {
	fmt.Println("handler")
}

func checkError(err error) {
	if err != nil {
		logger.Error(err.Error())
		panic(err.Error())
	}
}
