package thrift

import (
	"context"
	"math/rand"
	"testing"
	"time"

	api "github.com/blink-io/kratos-transport/testing/api/thrift/gen-go/hygrothermograph"
)

type HygrothermographHandler struct {
}

func NewHygrothermographHandler() *HygrothermographHandler {
	return &HygrothermographHandler{}
}

func (p *HygrothermographHandler) GetHygrothermograph(_ context.Context) (_r *api.Hygrothermograph, _err error) {
	Humidity := float64(rand.Intn(100))
	Temperature := float64(rand.Intn(100))
	_r = &api.Hygrothermograph{
		Humidity:    &Humidity,
		Temperature: &Temperature,
	}
	return
}

func TestServer(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	srv := NewServer(
		WithAddress(":7700"),
		WithProcessor(api.NewHygrothermographServiceProcessor(NewHygrothermographHandler())),
	)

	if err := srv.Start(ctx); err != nil {
		cancel()
		panic(err)
	}

	time.Sleep(100 * time.Millisecond)

	if err := srv.Stop(ctx); err != nil {
		t.Errorf("expected nil got %v", err)
	}
}
