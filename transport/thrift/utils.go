package thrift

import (
	"crypto/tls"

	"github.com/apache/thrift/lib/go/thrift"
)

const (
	ProtocolCompact    = "compact"
	ProtocolBinary     = "binary"
	ProtocolSimpleJSON = "simplejson"
	ProtocolJSON       = "json"
	ProtocolDebug      = "debug"
)

func createProtocolFactory(protocol string, conf *thrift.TConfiguration) thrift.TProtocolFactory {
	switch protocol {
	case ProtocolCompact:
		return thrift.NewTCompactProtocolFactoryConf(conf)
	case ProtocolSimpleJSON:
		return thrift.NewTSimpleJSONProtocolFactoryConf(conf)
	case ProtocolJSON:
		return thrift.NewTJSONProtocolFactory()
	case ProtocolBinary, "":
		return thrift.NewTBinaryProtocolFactoryConf(conf)
	default:
		return nil
	}
}

func createTransportFactory(cfg *thrift.TConfiguration, buffered, framed bool, bufferSize int) thrift.TTransportFactory {
	var transportFactory thrift.TTransportFactory

	if buffered {
		transportFactory = thrift.NewTBufferedTransportFactory(bufferSize)
	} else {
		transportFactory = thrift.NewTTransportFactory()
	}

	if framed {
		transportFactory = thrift.NewTFramedTransportFactoryConf(transportFactory, cfg)
	}

	return transportFactory
}

func createServerTransport(address string, tlsConf *tls.Config) (thrift.TServerTransport, error) {
	if tlsConf != nil {
		return thrift.NewTSSLServerSocket(address, tlsConf)
	} else {
		return thrift.NewTServerSocket(address)
	}
}

func createClientTransport(transportFactory thrift.TTransportFactory, address string, secure bool, cfg *thrift.TConfiguration) (thrift.TTransport, error) {
	var transport thrift.TTransport
	if secure {
		transport = thrift.NewTSSLSocketConf(address, cfg)
	} else {
		transport = thrift.NewTSocketConf(address, cfg)
	}
	transport, err := transportFactory.GetTransport(transport)
	if err != nil {
		return nil, err
	}
	if err := transport.Open(); err != nil {
		return nil, err
	}
	return transport, nil
}
