package registrar

import (
	"context"
	"net/http"

	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/tx7do/kratos-transport/transport/http3"
)

type ServerRouter interface {
	Handle(path string, h http.Handler)
	HandleFunc(path string, h http.HandlerFunc)
	HandleHeader(key, val string, h http.HandlerFunc)
	HandlePrefix(prefix string, h http.Handler)
	Route(prefix string, filters ...khttp.FilterFunc) *http3.Router
}

var _ ServerRouter = (*http3.Server)(nil)

type RegisterFunc func(context.Context, ServerRouter) error

type WithRegistrar interface {
	HTTP3Registrar(context.Context) RegisterFunc
}

type Func[S any] func(ServerRouter, S)

type FuncWithErr[S any] func(ServerRouter, S) error

type CtxFunc[S any] func(context.Context, ServerRouter, S)

type CtxFuncWithErr[S any] func(context.Context, ServerRouter, S) error

type Registrar interface {
	RegisterToHTTP3(context.Context, ServerRouter) error
}

type registrar[S any] struct {
	s S
	f CtxFuncWithErr[S]
}

var _ Registrar = (*registrar[any])(nil)

func New[S any](s S, f Func[S]) Registrar {
	cf := func(ctx context.Context, r ServerRouter, s S) error {
		f(r, s)
		return nil
	}
	return NewCtxWithErr(s, cf)
}

func NewCtx[S any](s S, f CtxFunc[S]) Registrar {
	cf := func(ctx context.Context, r ServerRouter, s S) error {
		f(ctx, r, s)
		return nil
	}
	return NewCtxWithErr(s, cf)
}

// NewWithErr creates a registrar with returning error.
func NewWithErr[S any](s S, f FuncWithErr[S]) Registrar {
	cf := func(ctx context.Context, r ServerRouter, s S) error {
		return f(r, s)
	}
	return NewCtxWithErr(s, cf)
}

// NewCtxWithErr creates a registrar with a context parameter and returning error.
func NewCtxWithErr[S any](s S, f CtxFuncWithErr[S]) Registrar {
	h := &registrar[S]{
		s: s,
		f: f,
	}
	return h
}

func (h registrar[S]) RegisterToHTTP3(ctx context.Context, r ServerRouter) error {
	return h.f(ctx, r, h.s)
}
