package http

import (
	"net/http"

	khttp "github.com/go-kratos/kratos/v2/transport/http"
)

func StdHandlerFunc(h http.HandlerFunc) khttp.HandlerFunc {
	return func(ctx khttp.Context) error {
		w := ctx.Response()
		r := ctx.Request()
		h(w, r)
		return nil
	}
}

func StdHandler(h http.Handler) khttp.HandlerFunc {
	return func(ctx khttp.Context) error {
		w := ctx.Response()
		r := ctx.Request()
		h.ServeHTTP(w, r)
		return nil
	}
}
