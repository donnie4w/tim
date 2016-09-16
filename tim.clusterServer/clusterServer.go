package clusterServer

/**
 * donnie4w@gmail.com  tim server
 */
import (
	"fmt"

	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/donnie4w/go-logger/logger"
	"tim.client"
	. "tim.common"
	. "tim.protocol"
)

type Controlloer struct {
	addr string
}

func ServerStart() {
	s := new(Controlloer)
	s.SetAddr(ClusterConf.RequestAddr)
	s.Server()
}

func (t *Controlloer) SetAddr(addr string) {
	t.addr = addr
}

func (t *Controlloer) ListenAddr() string {
	return t.addr
}

func (t *Controlloer) Server() {
	transportFactory := thrift.NewTBufferedTransportFactory(1024)
	protocolFactory := thrift.NewTCompactProtocolFactory()
	serverTransport, err := thrift.NewTServerSocket(t.ListenAddr())
	if err != nil {
		logger.Error("server:", err.Error())
		panic(err.Error())
	}
	handler := new(client.TimImpl)
	processor := NewITimProcessor(handler)
	server := thrift.NewTSimpleServer4(processor, serverTransport, transportFactory, protocolFactory)
	fmt.Println("cluster server listen:", t.ListenAddr())
	server.Serve()
}
