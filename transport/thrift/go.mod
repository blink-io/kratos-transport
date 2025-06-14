module github.com/blink-io/kratos-transport/transport/thrift

go 1.23.0

toolchain go1.24.1

require (
	github.com/apache/thrift v0.22.0
	github.com/blink-io/kratos-transport v0.0.0-20250612063448-9d75c0c37c8e
	github.com/go-kratos/kratos/v2 v2.8.4
)

replace github.com/blink-io/kratos-transport => ../../

require (
	github.com/go-playground/form/v4 v4.2.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
