package http3

import (
	khttp "github.com/go-kratos/kratos/v2/transport/http"
)

const (
	SupportPackageIsVersion1 = khttp.SupportPackageIsVersion1
)

var (
	FilterChain = khttp.FilterChain

	DefaultRequestDecoder = khttp.DefaultRequestDecoder
	DefaultRequestEncoder = khttp.DefaultRequestEncoder
	DefaultRequestQuery   = khttp.DefaultRequestQuery
	DefaultRequestVars    = khttp.DefaultRequestVars

	DefaultResponseEncoder = khttp.DefaultResponseEncoder
	DefaultResponseDecoder = khttp.DefaultResponseDecoder

	DefaultErrorEncoder = khttp.DefaultErrorEncoder
	DefaultErrorDecoder = khttp.DefaultErrorDecoder

	CodecForResponse = khttp.CodecForResponse
)

type (
	Context = khttp.Context

	HandlerFunc = khttp.HandlerFunc

	FilterFunc = khttp.FilterFunc

	RouteInfo = khttp.RouteInfo

	WalkRouteFunc = khttp.WalkRouteFunc

	DecodeRequestFunc  = khttp.DecodeRequestFunc
	DecodeResponseFunc = khttp.DecodeResponseFunc
	DecodeErrorFunc    = khttp.DecodeErrorFunc

	EncodeRequestFunc  = khttp.EncodeRequestFunc
	EncodeResponseFunc = khttp.EncodeResponseFunc
	EncodeErrorFunc    = khttp.EncodeErrorFunc
)
