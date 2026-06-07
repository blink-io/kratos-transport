module github.com/blink-io/kratos-transport/transport/thrift

go 1.25.0

require (
	github.com/apache/thrift v0.23.0
	github.com/blink-io/kratos-transport v0.0.0-20260507153638-31dc78fc0ffb
	github.com/go-kratos/kratos/v2 v2.9.2
)

replace github.com/blink-io/kratos-transport => ../../

require (
	github.com/go-playground/form/v4 v4.3.0 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
