/**
 * donnie4w@gmail.com  tim server
 */
package tfClient

import (
	"errors"
	"fmt"
	"runtime/debug"

	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/donnie4w/go-logger/logger"
	//	. "tim.common"
	. "tim.protocol"
)

func HttpClient(f func(*ITimClient) error, urlstr string) (err error) {
	defer func() {
		if er := recover(); er != nil {
			err = errors.New(fmt.Sprint(er))
			logger.Error(er)
			logger.Error(string(debug.Stack()))
		}
	}()
	if urlstr != "" {
		logger.Debug("httpClient url:", urlstr)
		transport, err := thrift.NewTHttpPostClient(urlstr)
		defer transport.Close()
		if err == nil {
			factory := thrift.NewTCompactProtocolFactory()
			transport.Open()
			itimClient := NewITimClientFactory(transport, factory)
			err = f(itimClient)
		}
	}
	return
}

func HttpClient2(f func(*ITimClient) error, user_auth_url string) (err error) {
	defer func() {
		if er := recover(); er != nil {
			err = errors.New(fmt.Sprint(er))
			logger.Error(er)
			logger.Error(string(debug.Stack()))
		}
	}()
	if user_auth_url != "" {
		logger.Debug("httpClient url:", user_auth_url)
		transport, err := thrift.NewTHttpPostClient(user_auth_url)
		defer transport.Close()
		if err == nil {
			factory := thrift.NewTCompactProtocolFactory()
			transport.Open()
			itimClient := NewITimClientFactory(transport, factory)
			err = f(itimClient)
		}
	} else {
		err = errors.New("httpclient url is null")
	}
	return
}

func TcpClient(f func(*ITimClient), urlstr string) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(err)
		}
	}()
	if urlstr != "" {
		logger.Debug("tcpClient addr:", urlstr)
		transport, err := thrift.NewTSocket(urlstr)
		defer transport.Close()
		if err == nil {
			protocolFactory := thrift.NewTCompactProtocolFactory()
			itimClient := NewITimClientFactory(transport, protocolFactory)
			f(itimClient)
		}
	}
}
