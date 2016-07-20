/**
 * donnie4w@gmail.com  tim server
 */
package thriftserver

import (
	"fmt"
	"log"
	"runtime/debug"

	thrift "git.apache.org/thrift.git/lib/go/thrift"
)

// Simple, non-concurrent server for testing.
type TSimpleServer struct {
	quit chan struct{}

	processorFactory       thrift.TProcessorFactory
	serverTransport        thrift.TServerTransport
	inputTransportFactory  thrift.TTransportFactory
	outputTransportFactory thrift.TTransportFactory
	inputProtocolFactory   thrift.TProtocolFactory
	outputProtocolFactory  thrift.TProtocolFactory
}

func NewTSimpleServer2(processor thrift.TProcessor, serverTransport thrift.TServerTransport) *TSimpleServer {
	return NewTSimpleServerFactory2(thrift.NewTProcessorFactory(processor), serverTransport)
}

func NewTSimpleServer4(processor thrift.TProcessor, serverTransport thrift.TServerTransport, transportFactory thrift.TTransportFactory, protocolFactory thrift.TProtocolFactory) *TSimpleServer {
	return NewTSimpleServerFactory4(thrift.NewTProcessorFactory(processor),
		serverTransport,
		transportFactory,
		protocolFactory,
	)
}

func NewTSimpleServer6(processor thrift.TProcessor, serverTransport thrift.TServerTransport, inputTransportFactory thrift.TTransportFactory, outputTransportFactory thrift.TTransportFactory, inputProtocolFactory thrift.TProtocolFactory, outputProtocolFactory thrift.TProtocolFactory) *TSimpleServer {
	return NewTSimpleServerFactory6(thrift.NewTProcessorFactory(processor),
		serverTransport,
		inputTransportFactory,
		outputTransportFactory,
		inputProtocolFactory,
		outputProtocolFactory,
	)
}

func NewTSimpleServerFactory2(processorFactory thrift.TProcessorFactory, serverTransport thrift.TServerTransport) *TSimpleServer {
	return NewTSimpleServerFactory6(processorFactory,
		serverTransport,
		thrift.NewTTransportFactory(),
		thrift.NewTTransportFactory(),
		thrift.NewTBinaryProtocolFactoryDefault(),
		thrift.NewTBinaryProtocolFactoryDefault(),
	)
}

func NewTSimpleServerFactory4(processorFactory thrift.TProcessorFactory, serverTransport thrift.TServerTransport, transportFactory thrift.TTransportFactory, protocolFactory thrift.TProtocolFactory) *TSimpleServer {
	return NewTSimpleServerFactory6(processorFactory,
		serverTransport,
		transportFactory,
		transportFactory,
		protocolFactory,
		protocolFactory,
	)
}

func NewTSimpleServerFactory6(processorFactory thrift.TProcessorFactory, serverTransport thrift.TServerTransport, inputTransportFactory thrift.TTransportFactory, outputTransportFactory thrift.TTransportFactory, inputProtocolFactory thrift.TProtocolFactory, outputProtocolFactory thrift.TProtocolFactory) *TSimpleServer {
	return &TSimpleServer{
		processorFactory:       processorFactory,
		serverTransport:        serverTransport,
		inputTransportFactory:  inputTransportFactory,
		outputTransportFactory: outputTransportFactory,
		inputProtocolFactory:   inputProtocolFactory,
		outputProtocolFactory:  outputProtocolFactory,
		quit: make(chan struct{}, 1),
	}
}

func (p *TSimpleServer) ProcessorFactory() thrift.TProcessorFactory {
	return p.processorFactory
}

func (p *TSimpleServer) ServerTransport() thrift.TServerTransport {
	return p.serverTransport
}

func (p *TSimpleServer) InputTransportFactory() thrift.TTransportFactory {
	return p.inputTransportFactory
}

func (p *TSimpleServer) OutputTransportFactory() thrift.TTransportFactory {
	return p.outputTransportFactory
}

func (p *TSimpleServer) InputProtocolFactory() thrift.TProtocolFactory {
	return p.inputProtocolFactory
}

func (p *TSimpleServer) OutputProtocolFactory() thrift.TProtocolFactory {
	return p.outputProtocolFactory
}

func (p *TSimpleServer) Listen() error {
	return p.serverTransport.Listen()
}

func (p *TSimpleServer) AcceptLoop() error {
	for {
		client, err := p.serverTransport.Accept()
		if err != nil {
			select {
			case <-p.quit:
				return nil
			default:
			}
			return err
		}
		if client != nil {
			go func() {
				if err := p.processRequests(client); err != nil {
					log.Println("error processing request:", err)
				}
			}()
		}
	}
}

func (p *TSimpleServer) Serve() error {
	err := p.Listen()
	if err != nil {
		return err
	}
	p.AcceptLoop()
	return nil
}

func (p *TSimpleServer) Stop() error {
	p.quit <- struct{}{}
	p.serverTransport.Interrupt()
	return nil
}

func (p *TSimpleServer) processRequests(client thrift.TTransport) error {
	processor := p.processorFactory.GetProcessor(client)
	inputTransport := p.inputTransportFactory.GetTransport(client)
	outputTransport := p.outputTransportFactory.GetTransport(client)
	inputProtocol := p.inputProtocolFactory.GetProtocol(inputTransport)
	outputProtocol := p.outputProtocolFactory.GetProtocol(outputTransport)
	defer func() {
		if e := recover(); e != nil {
			log.Printf("panic in processor: %s: %s", e, debug.Stack())
		}
	}()
	if inputTransport != nil {
		defer inputTransport.Close()
	}
	if outputTransport != nil {
		defer outputTransport.Close()
	}
	for {
		ok, err := processor.Process(inputProtocol, outputProtocol)
		if err, ok := err.(thrift.TTransportException); ok && err.TypeId() == thrift.END_OF_FILE {
			return nil
		} else if err != nil {
			log.Printf("error processing request: %s", err)
			return err
		}
		if !ok {
			break
		}
	}
	return nil
}

/**wuxiaodong*/
func (p *TSimpleServer) Serve2(f func(thrift.TTransport)) error {
	err := p.Listen()
	if err != nil {
		return err
	}
	p.AcceptLoop2(f)
	return nil
}

/**wuxiaodong*/
func (p *TSimpleServer) AcceptLoop2(f func(thrift.TTransport)) error {
	for {
		client, err := p.serverTransport.Accept()
		if err != nil {
			select {
			case <-p.quit:
				return nil
			default:
			}
			return err
		}
		if client != nil {
			go func() {
				if err := p.processRequests(client); err != nil {
					log.Println("error processing request:", err)
				}
			}()
			go f(client)
		}
	}
}

/**wuxiaodong*/
func (p *TSimpleServer) Accept() (client thrift.TTransport, err error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("simple_server Accept() >>>", err)
		}
	}()
	client, err = p.serverTransport.Accept()
	if err != nil {
		select {
		case <-p.quit:
			return
		default:
		}
		return
	}
	return
}
