package http

import (
	"net/http"
)

type keyValue struct {
	key   string
	value string
}

type doOptions struct {
	method      string
	statusCode  int
	ahs         []*keyValue
	shs         map[string]string
	skipQuery   bool
	skipVars    bool
	skipReqBody bool
	skipResBody bool
}

type DoOption func(*doOptions)

func applyDoOptions(ops ...DoOption) *doOptions {
	opt := &doOptions{
		ahs: make([]*keyValue, 0),
		shs: make(map[string]string),
	}
	for _, o := range ops {
		o(opt)
	}
	if len(opt.method) == 0 {
		opt.method = http.MethodGet
	}
	if opt.statusCode == 0 {
		opt.statusCode = http.StatusOK
	}
	return opt
}

func StatusCode(statusCode int) DoOption {
	return func(o *doOptions) {
		o.statusCode = statusCode
	}
}

func AddHeader(key string, value string) DoOption {
	return func(o *doOptions) {
		o.ahs = append(o.ahs, &keyValue{key, value})
	}
}

func SetHeader(key string, value string) DoOption {
	return func(o *doOptions) {
		o.shs[key] = value
	}
}

func SkipVars() DoOption {
	return func(o *doOptions) {
		o.skipVars = true
	}
}

func SkipQuery() DoOption {
	return func(o *doOptions) {
		o.skipQuery = true
	}
}

func SkipReqBody() DoOption {
	return func(o *doOptions) {
		o.skipReqBody = true
	}
}

func SkipResBody() DoOption {
	return func(o *doOptions) {
		o.skipResBody = true
	}
}
