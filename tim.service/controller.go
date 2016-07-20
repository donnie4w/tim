/**
 * donnie4w@gmail.com  tim server
 */
package service

import (
	"errors"
	"fmt"
	"log"
	"runtime/debug"
	"strconv"
	"sync"
	"time"

	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/donnie4w/go-logger/logger"
	"tim.FW"
	. "tim.common"
	. "tim.connect"
	. "tim.impl"
	. "tim.protocol"
	"tim.thriftserver"
	"tim.utils"
)

type Controlloer struct {
	Port int
	Ip   string
}

func ServerStart() {
	go Httpserver()
	s := new(Controlloer)
	s.SetAddr(ConfBean.GetIp())
	s.Server()
}

func (t *Controlloer) SetAddr(port int, ip string) {
	t.Port = port
	t.Ip = ip
}

func (t *Controlloer) ListenAddr() string {
	return fmt.Sprint(t.Ip, ":", t.Port)
}

func (t *Controlloer) Server() {
	transportFactory := thrift.NewTBufferedTransportFactory(1024)
	protocolFactory := thrift.NewTCompactProtocolFactory()
	serverTransport, err := thrift.NewTServerSocket(t.ListenAddr())
	if err != nil {
		logger.Error("server:", err.Error())
		panic(err.Error())
	}
	handler := new(TimImpl)
	processor := NewITimProcessor(handler)
	server := thriftserver.NewTSimpleServer4(processor, serverTransport, transportFactory, protocolFactory)
	logger.Info("server listen:", t.ListenAddr())
	//go TP.Run()
	//time.Sleep(1 * time.Second)
	Listen(server, 100)
	if err == nil {
		for {
			client, err := Accept(server)
			if err == nil {
				go controllerHandler(client)
			}
		}
	}
}

func Listen(server *thriftserver.TSimpleServer, count int) (err error) {
	if count <= 0 {
		err = errors.New("")
		return
	}
	defer func() {
		if err := recover(); err != nil {
			logger.Error("Listen,", err)
			logger.Error(string(debug.Stack()))
			count--
			Listen(server, count)
		}
	}()
	err = server.Listen()
	return
}

func Accept(server *thriftserver.TSimpleServer) (client thrift.TTransport, err error) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("Accept,", err)
			logger.Error(string(debug.Stack()))
		}
	}()
	client, err = server.Accept()
	return
}

func controllerHandler(tt thrift.TTransport) {
	isclose := false
	var gorutineclose *bool = &isclose
	defer func() {
		if err := recover(); err != nil {
			logger.Error("controllerHandler,", err)
			logger.Error(string(debug.Stack()))
			*gorutineclose = true
		}
	}()
	tu := &TimUser{Client: NewTimClient(tt), OverLimit: 3, Fw: FW.CONNECT, IdCardNo: utils.TimeMills(), Sendflag: make(chan string), Sync: new(sync.Mutex)}
	//	TP.Register <- tu
	TP.AddConnect(tu)
	//	defer func() { TP.Unregister <- tu }()
	defer func() { TP.DeleteTimUser(tu) }()
	defer func() { tt.Close() }()
	monitorChan := make(chan string, 2)
	heartbeat := ConfBean.HeartBeat
	if heartbeat > 0 {
		go func() {
			defer func() {
				if err := recover(); err != nil {
					logger.Error(string(debug.Stack()))
				}
			}()
			defer func() {
				if err := recover(); err != nil {
				}
				*gorutineclose = true
				monitorChan <- "ping end"
			}()
			for {
				if *gorutineclose {
					break
				}
				for i := 0; i < heartbeat; i++ {
					time.Sleep(1 * time.Second)
					if *gorutineclose {
						goto END
					}
					if tu.OverLimit <= 0 {
						goto END
					}
					if tu.Fw == FW.CLOSE {
						goto END
					}
				}
				if tu.Fw == FW.AUTH {
					er := tu.Client.TimPing(fmt.Sprint(utils.TimeMills()))
					if er != nil {
						logger.Error("ping err", er.Error())
						panic("ping err")
					}
					tu.OverLimit--
				} else {
					logger.Error("auth over time")
					goto END
				}
				if tu.Fw == FW.CLOSE {
					break
				}
			}
		END:
		}()
	}
	go TimProcessor(tt, tu, gorutineclose, monitorChan)
	errormsg := <-monitorChan
	logger.Error("errormsg:", errormsg)
}

func NewTimClient(tt thrift.TTransport) *ITimClient {
	transportFactory := thrift.NewTBufferedTransportFactory(1024)
	protocolFactory := thrift.NewTCompactProtocolFactory()
	useTransport := transportFactory.GetTransport(tt)
	return NewITimClientFactory(useTransport, protocolFactory)
}

func TimProcessor(client thrift.TTransport, tu *TimUser, gorutineclose *bool, monitorChan chan string) error {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(string(debug.Stack()))
		}
	}()
	defer func() {
		if err := recover(); err != nil {
			logger.Error("TimProcessor error", err)
			logger.Error(string(debug.Stack()))
		}
		*gorutineclose = true
		monitorChan <- "timProcessor end"
	}()
	compactprotocol := thrift.NewTCompactProtocol(client)
	pub := strconv.Itoa(time.Now().Nanosecond())
	handler := &TimImpl{Pub: pub, Client: client, Tu: tu}
	processor := NewITimProcessor(handler)
	for {
		ok, err := processor.Process(compactprotocol, compactprotocol)
		if err, ok := err.(thrift.TTransportException); ok && err.TypeId() == thrift.END_OF_FILE {
			return nil
		} else if err != nil {
			log.Printf("error processing request: %s", err)
			return err
		}
		if !ok {
			logger.Debug("is not ok", ok, err)
			break
		}
		if tu.Fw == FW.CLOSE || tu.OverLimit <= 0 {
			break
		}
	}
	return nil
}
