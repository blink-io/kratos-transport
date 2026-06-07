module github.com/blink-io/kratos-transport/transport/http3

go 1.25.0

require (
	github.com/blink-io/kratos-transport v0.0.0-20260507153638-31dc78fc0ffb
	github.com/go-kratos/kratos/v2 v2.9.2
	github.com/gorilla/mux v1.8.1
	github.com/quic-go/quic-go v0.60.0
	github.com/stretchr/testify v1.11.1
)

replace github.com/blink-io/kratos-transport => ../../

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-kratos/aegis v0.2.0 // indirect
	github.com/go-playground/form/v4 v4.3.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/quic-go/qpack v0.6.0 // indirect
	go.uber.org/mock v0.6.0 // indirect
	golang.org/x/crypto v0.52.0 // indirect
	golang.org/x/net v0.55.0 // indirect
	golang.org/x/sys v0.45.0 // indirect
	golang.org/x/text v0.37.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260526163538-3dc84a4a5aaa // indirect
	google.golang.org/grpc v1.81.1 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
