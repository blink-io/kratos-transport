package thrift

import (
	"crypto/tls"

	"github.com/apache/thrift/lib/go/thrift"
	"github.com/go-kratos/kratos/v2/registry"
)

type clientOptions struct {
	Client *thrift.TStandardClient

	discovery registry.Discovery
	tlsConf   *tls.Config

	endpoint string

	protocol string

	buffered   bool
	framed     bool
	bufferSize int

	secure bool

	tconf *thrift.TConfiguration
}

type Connection struct {
	Client    *thrift.TStandardClient
	Transport thrift.TTransport
}

func (c *Connection) Close() error {
	if c.Transport == nil {
		return nil
	}
	return c.Transport.Close()
}

func Dial(opts ...ClientOption) (*Connection, error) {
	return dial(opts...)
}

func dial(opts ...ClientOption) (*Connection, error) {
	cli := &clientOptions{
		bufferSize: 8192,
		buffered:   false,
		framed:     false,
		protocol:   ProtocolBinary,
		secure:     false,
		tconf: &thrift.TConfiguration{},
	}

	for _, o := range opts {
		o(cli)
	}

	protocolFactory := createProtocolFactory(cli.protocol, cli.tconf)
	if protocolFactory == nil {
		return nil, ErrInvalidProtocol
	}

	transportFactory := createTransportFactory(cli.tconf, cli.buffered, cli.framed, cli.bufferSize)
	if transportFactory == nil {
		return nil, ErrInvalidTransport
	}

	clientTransport, err := createClientTransport(transportFactory, cli.endpoint, cli.secure, cli.tconf)
	if err != nil {
		return nil, err
	}

	inProto := protocolFactory.GetProtocol(clientTransport)
	outProto := protocolFactory.GetProtocol(clientTransport)

	return &Connection{
		Client:    thrift.NewTStandardClient(inProto, outProto),
		Transport: clientTransport,
	}, nil
}
