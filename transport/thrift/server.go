package thrift

import (
	"context"
	"crypto/tls"
	"errors"
	"net/url"

	"github.com/apache/thrift/lib/go/thrift"
	klog "github.com/go-kratos/kratos/v3/log"
	"github.com/go-kratos/kratos/v3/transport"
)

var (
	_ transport.Server     = (*Server)(nil)
	_ transport.Endpointer = (*Server)(nil)
)

var (
	ErrInvalidProtocol  = errors.New("invalid protocol")
	ErrInvalidTransport = errors.New("invalid transport")
)

type Server struct {
	Server     *thrift.TSimpleServer
	tlsConf    *tls.Config
	address    string
	protocol   string
	buffered   bool
	framed     bool
	bufferSize int
	err        error
	processor  thrift.TProcessor
	tconf      *thrift.TConfiguration
}

func NewServer(opts ...ServerOption) *Server {
	srv := &Server{
		bufferSize: 8192,
		buffered:   false,
		framed:     false,
		protocol:   ProtocolBinary,
		tconf:      &thrift.TConfiguration{},
	}
	srv.init(opts...)
	return srv
}

func (s *Server) Name() string {
	return string(KindThrift)
}

func (s *Server) init(opts ...ServerOption) {
	for _, o := range opts {
		o(s)
	}
}

func (s *Server) Endpoint() (*url.URL, error) {
	return url.Parse("tcp://" + s.address)
}

func (s *Server) Start(ctx context.Context) error {
	if s.err != nil {
		return s.err
	}

	protocolFactory := createProtocolFactory(s.protocol, s.tconf)
	if protocolFactory == nil {
		return ErrInvalidProtocol
	}

	transportFactory := createTransportFactory(s.tconf, s.buffered, s.framed, s.bufferSize)
	if transportFactory == nil {
		return ErrInvalidTransport
	}

	serverTransport, serverTransportErr := createServerTransport(s.address, s.tlsConf)
	if serverTransportErr != nil {
		return serverTransportErr
	}

	klog.Info("[Thrift] server listening", "addr", s.address)

	s.Server = thrift.NewTSimpleServer4(s.processor, serverTransport, transportFactory, protocolFactory)
	go func() {
		_ = s.Server.Serve()
	}()

	go func() {
		<-ctx.Done()
		_ = s.Server.Stop()
	}()

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	klog.Info("[Thrift] server stopping")

	if s.Server != nil {
		stopCh := make(chan error, 1)
		go func() {
			stopCh <- s.Server.Stop()
		}()

		select {
		case stopErr := <-stopCh:
			return stopErr
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return nil
}
