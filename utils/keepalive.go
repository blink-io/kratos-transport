package utils

import (
	"context"
	"crypto/tls"
	"fmt"
	"math/rand"
	"net"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/go-kratos/kratos/v2/log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type KeepAliveService struct {
	*grpc.Server
	health *health.Server

	lis      net.Listener
	tlsConf  *tls.Config
	endpoint *url.URL

	mu sync.Mutex
}

func NewKeepAliveService(tlsConf *tls.Config) *KeepAliveService {
	srv := &KeepAliveService{
		tlsConf: tlsConf,
		health:  health.NewServer(),
	}

	var grpcOpts []grpc.ServerOption
	if srv.tlsConf != nil {
		grpcOpts = append(grpcOpts, grpc.Creds(credentials.NewTLS(srv.tlsConf)))
	}

	srv.Server = grpc.NewServer(grpcOpts...)

	grpc_health_v1.RegisterHealthServer(srv.Server, srv.health)

	return srv
}

func (s *KeepAliveService) Start() error {
	if err := s.generateEndpoint(); err != nil {
		return err
	}

	s.health.Resume()

	log.Debugf("keep alive service started at %s", s.lis.Addr().String())

	go func() {
		_ = s.Serve(s.lis)
	}()

	return nil
}

func (s *KeepAliveService) Stop(ctx context.Context) error {
	s.health.Shutdown()

	done := make(chan struct{})
	go func() {
		s.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
		log.Debug("keep alive service stopped")
		return nil
	case <-ctx.Done():
		s.Server.Stop()
		return ctx.Err()
	}
}

func (s *KeepAliveService) generatePort(min, max int) int {
	return rand.Intn(max-min) + min
}

func (s *KeepAliveService) generateEndpoint() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.endpoint != nil {
		return nil
	}

	const maxRetries = 100
	for i := 0; i < maxRetries; i++ {
		port := s.generatePort(10000, 65535)
		host := ""
		if itf, ok := os.LookupEnv("KRATOS_TRANSPORT_KEEPALIVE_INTERFACE"); ok {
			h, err := getIPAddress(itf)
			if err != nil {
				return err
			}
			host = h
		}

		addr := fmt.Sprintf("%s:%d", host, port)
		lis, err := net.Listen("tcp", addr)
		if err == nil && lis != nil {
			s.lis = lis
			endpoint, parseErr := url.Parse("tcp://" + addr)
			if parseErr != nil {
				lis.Close()
				return parseErr
			}
			s.endpoint = endpoint
			return nil
		}
		time.Sleep(10 * time.Millisecond)
	}
	return fmt.Errorf("failed to find available port after %d attempts", maxRetries)
}

func (s *KeepAliveService) Endpoint() (*url.URL, error) {
	if err := s.generateEndpoint(); err != nil {
		return nil, err
	}
	return s.endpoint, nil
}

func getIPAddress(interfaceName string) (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range interfaces {
		if iface.Name == interfaceName {
			addrs, err := iface.Addrs()
			if err != nil {
				return "", err
			}

			// Find the first IPv4 address
			for _, addr := range addrs {
				if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
					return ipnet.IP.String(), nil
				}
			}
			return "", fmt.Errorf("No IPv4 address found for interface %s", interfaceName)
		}
	}

	return "", fmt.Errorf("Interface %s not found", interfaceName)
}
