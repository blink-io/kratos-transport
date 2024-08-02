package hertz

import (
	"context"
	"math/rand"
	"strconv"
	"testing"

	"github.com/cloudwego/hertz/pkg/app"

	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-kratos/kratos/v2/transport/http/binding"

	"github.com/stretchr/testify/assert"

	api "github.com/tx7do/kratos-transport/testing/api/protobuf"
)

func TestServer(t *testing.T) {
	ctx := context.Background()

	srv := NewServer(
		WithAddress("127.0.0.1:8800"),
	)

	srv.GET("/login/*param", func(ctx context.Context, c *app.RequestContext) {
		if len(c.Params.ByName("param")) > 1 {
			c.AbortWithStatus(404)
			return
		}
		c.String(200, "Hello World!")
	})

	srv.GET("/hygrothermograph", func(ctx context.Context, c *app.RequestContext) {
		var out api.Hygrothermograph
		out.Humidity = strconv.FormatInt(int64(rand.Intn(100)), 10)
		out.Temperature = strconv.FormatInt(int64(rand.Intn(100)), 10)
		c.JSON(200, &out)
	})

	if err := srv.Start(ctx); err != nil {
		panic(err)
	}

	defer func() {
		if err := srv.Stop(ctx); err != nil {
			t.Errorf("expected nil got %v", err)
		}
	}()
}

func TestClient(t *testing.T) {
	ctx := context.Background()

	cli, err := khttp.NewClient(ctx,
		khttp.WithEndpoint("127.0.0.1:8800"),
	)
	assert.Nil(t, err)
	assert.NotNil(t, cli)

	resp, err := GetHygrothermograph(ctx, cli, nil, khttp.EmptyCallOption{})
	assert.Nil(t, err)
	t.Log(resp)
}

func GetHygrothermograph(ctx context.Context, cli *khttp.Client, in *api.Hygrothermograph, opts ...khttp.CallOption) (*api.Hygrothermograph, error) {
	var out api.Hygrothermograph

	pattern := "/hygrothermograph"
	path := binding.EncodeURL(pattern, in, true)

	opts = append(opts, khttp.Operation("/GetHygrothermograph"))
	opts = append(opts, khttp.PathTemplate(pattern))

	err := cli.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}

	return &out, nil
}
