/**
 * donnie4w@gmail.com  tim server
 */
package tfClient

import (
	"runtime/debug"

	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/donnie4w/go-logger/logger"
	. "tim.common"
	. "tim.protocol"
)

func HttpClient(f func(*ITimClient)) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(err)
			logger.Error(string(debug.Stack()))
		}
	}()
	user_auth_url := ConfBean.GetKV("user_auth_url", "")
	if user_auth_url != "" {
		logger.Debug("httpClient url:", user_auth_url)
		transport, err := thrift.NewTHttpPostClient(user_auth_url)
		defer transport.Close()
		if err == nil {
			factory := thrift.NewTCompactProtocolFactory()
			transport.Open()
			itimClient := NewITimClientFactory(transport, factory)
			f(itimClient)
		}
	}
}

func TcpClient(f func(*ITimClient)) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(err)
		}
	}()
	user_auth_addr := ConfBean.GetKV("user_auth_addr", "")
	if user_auth_addr != "" {
		logger.Debug("tcpClient addr:", user_auth_addr)
		transport, err := thrift.NewTSocket(user_auth_addr)
		defer transport.Close()
		if err == nil {
			protocolFactory := thrift.NewTCompactProtocolFactory()
			itimClient := NewITimClientFactory(transport, protocolFactory)
			f(itimClient)
		}
	}
}
