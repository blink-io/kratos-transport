package registrar

import (
	"context"

	kgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	"google.golang.org/grpc"
)

type ServiceRegistrar = grpc.ServiceRegistrar

type RegisterFunc func(context.Context, ServiceRegistrar) error

type WithRegistrar interface {
	GRPCRegistrar(context.Context) RegisterFunc
}

type Func[S any] func(ServiceRegistrar, S)

type FuncWithErr[S any] func(ServiceRegistrar, S) error

type CtxFunc[S any] func(context.Context, ServiceRegistrar, S)

type CtxFuncWithErr[S any] func(context.Context, ServiceRegistrar, S) error

type Registrar interface {
	RegisterToGRPC(context.Context, ServiceRegistrar) error
}

type registrar[S any] struct {
	s S
	f CtxFuncWithErr[S]
}

var _ Registrar = (*registrar[any])(nil)

func New[S any](s S, f Func[S]) Registrar {
	cf := func(ctx context.Context, r ServiceRegistrar, s S) error {
		f(r, s)
		return nil
	}
	return NewCtxWithErr[S](s, cf)
}

func NewWithErr[S any](s S, f FuncWithErr[S]) Registrar {
	cf := func(ctx context.Context, r ServiceRegistrar, s S) error {
		return f(r, s)
	}
	return NewCtxWithErr(s, cf)
}

func NewCtx[S any](s S, f CtxFunc[S]) Registrar {
	cf := func(ctx context.Context, r ServiceRegistrar, s S) error {
		f(ctx, r, s)
		return nil
	}
	return NewCtxWithErr(s, cf)
}

func NewCtxWithErr[S any](s S, f CtxFuncWithErr[S]) Registrar {
	h := &registrar[S]{
		s: s,
		f: f,
	}
	return h
}

func (h registrar[S]) RegisterToGRPC(ctx context.Context, r ServiceRegistrar) error {
	return h.f(ctx, r, h.s)
}

func IsServerThen(s any, f func(server *kgrpc.Server)) {
	if v, ok := s.(*kgrpc.Server); ok {
		f(v)
	}
}
