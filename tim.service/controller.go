/**
 * donnie4w@gmail.com  tim server
 */
package service

import (
	"crypto/tls"
	"errors"
	"fmt"
	"runtime/debug"
	"strconv"
	"sync"
	"time"

	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/donnie4w/go-logger/logger"
	"tim.FW"
	"tim.cluster"
	"tim.clusterRoute"
	"tim.clusterServer"
	. "tim.common"
	. "tim.connect"
	. "tim.impl"
	. "tim.protocol"
	"tim.route"
	"tim.thriftserver"
	"tim.utils"
)

type Controlloer struct {
	Port int
	Ip   string
}

func ServerStart() {
	go Httpserver()
	go tsslServer()
	s := new(Controlloer)
	s.SetAddr(CF.GetIp())
	if cluster.IsCluster() {
		go clusterServer.ServerStart()
	}
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
	fmt.Println("server listen:", t.ListenAddr())
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

func tsslServer() {
	if CF.TLSPort <= 0 && CF.TLSServerPem == "" && CF.TLSServerKey == "" {
		return
	}
	cer, err := tls.LoadX509KeyPair(CF.TLSServerPem, CF.TLSServerKey)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	config := &tls.Config{Certificates: []tls.Certificate{cer}}
	tsslServerSocket, err := thrift.NewTSSLServerSocket(fmt.Sprint(CF.Addr, ":", CF.TLSPort), config)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	err = tsslServerSocket.Listen()
	if err != nil {
		logger.Error(err.Error())
		return
	}
	fmt.Println("tls server listen:", CF.TLSPort)
	for {
		client, err := tsslServerSocket.Accept()
		if err == nil && client != nil {
			go controllerHandler(client)
		}
	}
}

func controllerHandler(tt thrift.TTransport) {
	isclose := false
	var gorutineclose *bool = &isclose
	defer func() {
		if err := recover(); err != nil {
			logger.Error("controllerHandler,", err)
			*gorutineclose = true
		}
	}()
	tu := &TimUser{Client: NewTimClient(tt), OverLimit: 3, Fw: FW.CONNECT, IdCardNo: utils.TimeMills(), Sendflag: make(chan string, 0), Sync: new(sync.Mutex)}
	TP.AddConnect(tu)
	defer func() {
		if cluster.IsCluster() && tu.UserTid != nil {
			loginname, err := GetLoginName(tu.UserTid)
			if loginname != "" && err == nil {
				cluster.DelLoginnameFromCluter(loginname)
			}
			if CF.Presence == 1 {
				p := OfflinePBean(tu.UserTid)
				go clusterRoute.ClusterRoutePBean(p)
			} else {
				go route.RoutePBean(OfflinePBean(tu.UserTid))
			}
		} else if tu.UserTid != nil && CF.Presence == 1 {
			go route.RoutePBean(OfflinePBean(tu.UserTid))
		}
		TP.DeleteTimUser(tu)
	}()
	defer func() { tt.Close() }()
	monitorChan := make(chan string, 2)
	heartbeat := CF.HeartBeat
	if heartbeat == 0 {
		heartbeat = 30 * 60
	}
	//	if heartbeat > 0 {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				//				logger.Error(string(debug.Stack()))
				logger.Error(err)
			}
		}()
		defer func() {
			if err := recover(); err != nil {
			}
			*gorutineclose = true
			monitorChan <- "ping end"
		}()
		checkinCluster := 0
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
				checkinCluster++
				if checkinCluster >= ClusterConf.Keytimeout/3 {
					checkinCluster = 0
					if cluster.IsCluster() && tu.UserTid != nil {
						loginname, err := GetLoginName(tu.UserTid)
						if loginname != "" && err == nil {
							cluster.SetLoginnameToCluster(loginname)
						}
					}
				}
			}
			if tu.Fw == FW.AUTH {
				er := tu.Ping()
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
	//	}
	go TimProcessor(tt, tu, gorutineclose, monitorChan)
	<-monitorChan
	//	errormsg := <-monitorChan
	//	logger.Error("errormsg:", errormsg)
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
			//			logger.Error(string(debug.Stack()))
			logger.Warn("processor:", err)
		}
	}()
	defer func() {
		if err := recover(); err != nil {
			//			logger.Error("TimProcessor error", err)
			//			logger.Error(string(debug.Stack()))
			logger.Warn("processor:", err)
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
			return err
		}
		if !ok {
			logger.Error("Processor error:", err)
			break
		}
		if tu.Fw == FW.CLOSE || tu.OverLimit <= 0 {
			break
		}
	}
	return nil
}
