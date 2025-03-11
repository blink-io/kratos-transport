module github.com/blink-io/kratos-transport/transport/thrift

go 1.23.0

toolchain go1.24.1

require (
	github.com/apache/thrift v0.21.0
	github.com/go-kratos/kratos/v2 v2.8.4
	github.com/tx7do/kratos-transport v1.1.8
)

require (
	github.com/go-playground/form/v4 v4.2.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	google.golang.org/protobuf v1.36.5 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/tx7do/kratos-transport => ../../
