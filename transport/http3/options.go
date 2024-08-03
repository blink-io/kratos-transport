package http3

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/go-kratos/kratos/v2/middleware"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gorilla/mux"
)

type ServerOption func(*Server)

func TLSConfig(c *tls.Config) ServerOption {
	return func(o *Server) {
		o.tlsConf = c
	}
}

func Address(addr string) ServerOption {
	return func(s *Server) {
		s.Addr = addr
	}
}

func Timeout(timeout time.Duration) ServerOption {
	return func(s *Server) {
		s.timeout = timeout
	}
}

func Middleware(m ...middleware.Middleware) ServerOption {
	return func(o *Server) {
		o.ms = m
	}
}

func Filter(filters ...khttp.FilterFunc) ServerOption {
	return func(o *Server) {
		o.filters = filters
	}
}

func RequestDecoder(dec khttp.DecodeRequestFunc) ServerOption {
	return func(o *Server) {
		o.decBody = dec
	}
}

func WithResponseEncoder(en khttp.EncodeResponseFunc) ServerOption {
	return func(o *Server) {
		o.enc = en
	}
}

func ErrorEncoder(en khttp.EncodeErrorFunc) ServerOption {
	return func(o *Server) {
		o.ene = en
	}
}

func StrictSlash(strictSlash bool) ServerOption {
	return func(o *Server) {
		o.strictSlash = strictSlash
	}
}

// PathPrefix with mux's PathPrefix, router will be replaced by a subrouter that start with prefix.
func PathPrefix(prefix string) ServerOption {
	return func(s *Server) {
		s.router = s.router.PathPrefix(prefix).Subrouter()
	}
}

// WalkRoute walks the router and all its sub-routers, calling walkFn for each route in the tree.
func (s *Server) WalkRoute(fn khttp.WalkRouteFunc) error {
	return s.router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		methods, err := route.GetMethods()
		if err != nil {
			return nil // ignore no methods
		}
		path, err := route.GetPathTemplate()
		if err != nil {
			return err
		}
		for _, method := range methods {
			if err := fn(khttp.RouteInfo{Method: method, Path: path}); err != nil {
				return err
			}
		}
		return nil
	})
}

// WalkHandle walks the router and all its sub-routers, calling walkFn for each route in the tree.
func (s *Server) WalkHandle(handle func(method, path string, handler http.HandlerFunc)) error {
	return s.WalkRoute(func(r khttp.RouteInfo) error {
		handle(r.Method, r.Path, s.ServeHTTP)
		return nil
	})
}

func NotFoundHandler(handler http.Handler) ServerOption {
	return func(s *Server) {
		s.router.NotFoundHandler = handler
	}
}

func MethodNotAllowedHandler(handler http.Handler) ServerOption {
	return func(s *Server) {
		s.router.MethodNotAllowedHandler = handler
	}
}
