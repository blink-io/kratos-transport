package registrar

import (
	"fmt"
	"testing"

	"github.com/go-kratos/kratos/v2/transport/grpc"
	kgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
)

func TestIsServerThen_1(t *testing.T) {
	var s *kgrpc.Server
	IsServerThen(s, func(s *grpc.Server) {
		fmt.Println("KKK")
	})
}
