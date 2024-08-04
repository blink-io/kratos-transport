package registrar

import (
	"context"

	khttp "github.com/go-kratos/kratos/v2/transport/http"
	httprgy "github.com/tx7do/kratos-transport/transport/http/registrar"
	"github.com/tx7do/kratos-transport/transport/http3"
)

type RouterRegistrar = httprgy.RouterRegistrar

var _ RouterRegistrar = (*http3.Server)(nil)

type RegisterFunc func(context.Context, RouterRegistrar) error

type WithRegistrar interface {
	HTTP3Registrar(context.Context) RegisterFunc
}

type Func[S any] func(RouterRegistrar, S)

type FuncWithErr[S any] func(RouterRegistrar, S) error

type CtxFunc[S any] func(context.Context, RouterRegistrar, S)

type CtxFuncWithErr[S any] func(context.Context, RouterRegistrar, S) error

type Registrar interface {
	RegisterToHTTP3(context.Context, RouterRegistrar) error
}

type registrar[S any] struct {
	s S
	f CtxFuncWithErr[S]
}

var _ Registrar = (*registrar[any])(nil)

func New[S any](s S, f Func[S]) Registrar {
	cf := func(ctx context.Context, r RouterRegistrar, s S) error {
		f(r, s)
		return nil
	}
	return NewCtxWithErr(s, cf)
}

func NewCtx[S any](s S, f CtxFunc[S]) Registrar {
	cf := func(ctx context.Context, r RouterRegistrar, s S) error {
		f(ctx, r, s)
		return nil
	}
	return NewCtxWithErr(s, cf)
}

// NewWithErr creates a registrar with returning error.
func NewWithErr[S any](s S, f FuncWithErr[S]) Registrar {
	cf := func(ctx context.Context, r RouterRegistrar, s S) error {
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

func (h registrar[S]) RegisterToHTTP3(ctx context.Context, r RouterRegistrar) error {
	return h.f(ctx, r, h.s)
}

type RouteSupport interface {
	Route(prefix string, filters ...khttp.FilterFunc) *http3.Router
}

func SupportsRouteThen(s any, f func(RouteSupport)) {
	if v, ok := s.(RouteSupport); ok {
		f(v)
	}
}

func IsServerThen(s any, f func(*http3.Server)) {
	if v, ok := s.(*http3.Server); ok {
		f(v)
	}
}
