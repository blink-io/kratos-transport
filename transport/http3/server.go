package http3

import (
	"context"
	"crypto/tls"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	khttp "github.com/go-kratos/kratos/v2/transport/http"

	"github.com/go-kratos/kratos/v2/transport"
	"github.com/gorilla/mux"
	"github.com/quic-go/quic-go/http3"
)

const (
	SupportPackageIsVersion1 = khttp.SupportPackageIsVersion1
)

var (
	_ transport.Server     = (*Server)(nil)
	_ transport.Endpointer = (*Server)(nil)
)

type Server struct {
	*http3.Server

	tlsConf  *tls.Config
	timeout  time.Duration
	endpoint *url.URL

	err error

	filters []khttp.FilterFunc
	ms      []middleware.Middleware
	dec     khttp.DecodeRequestFunc
	enc     khttp.EncodeResponseFunc
	ene     khttp.EncodeErrorFunc

	router      *mux.Router
	strictSlash bool
}

func NewServer(opts ...ServerOption) *Server {
	srv := &Server{
		timeout:     1 * time.Second,
		dec:         khttp.DefaultRequestDecoder,
		enc:         khttp.DefaultResponseEncoder,
		ene:         khttp.DefaultErrorEncoder,
		strictSlash: true,
	}

	srv.init(opts...)

	return srv
}

func (s *Server) init(opts ...ServerOption) {
	s.Server = &http3.Server{
		Addr: ":8443",
	}

	for _, o := range opts {
		o(s)
	}

	s.Server.TLSConfig = s.tlsConf

	s.router = mux.NewRouter().StrictSlash(s.strictSlash)
	s.router.NotFoundHandler = http.DefaultServeMux
	s.router.MethodNotAllowedHandler = http.DefaultServeMux

	s.Server.Handler = khttp.FilterChain(s.filters...)(s.router)

	_, _ = s.Endpoint()
}

func (s *Server) Endpoint() (*url.URL, error) {
	addr := s.Addr

	if s.tlsConf == nil {
		return nil, errors.New("http3: no TLS configured")
	}

	var prefix string
	if !strings.HasPrefix(addr, "https://") {
		prefix = "https://"
	}
	addr = prefix + addr

	s.endpoint, s.err = url.Parse(addr)

	return s.endpoint, s.err
}

func (s *Server) Start(ctx context.Context) error {
	if s.tlsConf == nil {
		return errors.New("http3: no TLS configured")
	}

	log.Infof("[HTTP3] server listening on: %s", s.Addr)

	if err := s.ListenAndServe(); err != nil {
		log.Errorf("[HTTP3] start server failed: %s", err.Error())
		return err
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	log.Info("[HTTP3] server stopping")
	return s.Close()
}

func (s *Server) Route(prefix string, filters ...khttp.FilterFunc) *Router {
	return newRouter(prefix, s, filters...)
}

func (s *Server) Handle(path string, h http.Handler) {
	s.router.Handle(path, h)
}

func (s *Server) HandlePrefix(prefix string, h http.Handler) {
	s.router.PathPrefix(prefix).Handler(h)
}

func (s *Server) HandleFunc(path string, h http.HandlerFunc) {
	s.router.HandleFunc(path, h)
}

func (s *Server) HandleHeader(key, val string, h http.HandlerFunc) {
	s.router.Headers(key, val).Handler(h)
}

func (s *Server) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	s.Handler.ServeHTTP(res, req)
}

func (s *Server) filter() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			var (
				ctx    context.Context
				cancel context.CancelFunc
			)
			if s.timeout > 0 {
				ctx, cancel = context.WithTimeout(req.Context(), s.timeout)
			} else {
				ctx, cancel = context.WithCancel(req.Context())
			}
			defer cancel()

			pathTemplate := req.URL.Path
			if route := mux.CurrentRoute(req); route != nil {
				// /path/123 -> /path/{id}
				pathTemplate, _ = route.GetPathTemplate()
			}

			tr := &Transport{
				endpoint:     s.endpoint.String(),
				operation:    pathTemplate,
				reqHeader:    headerCarrier(req.Header),
				replyHeader:  headerCarrier(w.Header()),
				request:      req,
				pathTemplate: pathTemplate,
			}

			tr.request = req.WithContext(transport.NewServerContext(ctx, tr))
			next.ServeHTTP(w, tr.request)
		})
	}
}
