package http3

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"testing"

	api "github.com/blink-io/kratos-transport/testing/api/protobuf"
	"github.com/blink-io/kratos-transport/testing/tlsutil"

	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-kratos/kratos/v2/transport/http/binding"
	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
	"github.com/stretchr/testify/assert"
)

func HygrothermographHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("HygrothermographHandler [%s] [%s] [%s]\n", r.Proto, r.Method, r.RequestURI)

	if r.Method == "POST" {
		var in api.Hygrothermograph
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			fmt.Printf("decode error: %s\n", err.Error())
		}
		fmt.Printf("Humidity: %s Temperature: %s \n", in.Humidity, in.Temperature)
	}

	var out api.Hygrothermograph
	out.Humidity = strconv.FormatInt(int64(rand.Intn(100)), 10)
	out.Temperature = strconv.FormatInt(int64(rand.Intn(100)), 10)
	_ = json.NewEncoder(w).Encode(&out)
}

type MyInfo2Req struct {
	Action string `json:"action"`
}

type MyInfo2Res struct {
	Action     string `json:"action"`
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}

//func info2() khttp.HandlerFunc {
//
//	hdlr := func(ctx context.Context, req *MyInfo2Req) (*MyInfo2Res, error) {
//		res := &MyInfo2Res{
//			Action:     req.Action,
//			Message:    "You are testing my info2",
//			StatusCode: 200,
//		}
//		return res, nil
//	}
//	return kthttp.GET[MyInfo2Req, MyInfo2Res]("info", hdlr)
//}

func TestServer(t *testing.T) {
	ctx := context.Background()

	srv := NewServer(
		Address(":8800"),
		TLSConfig(tlsutil.GenerateTLSConfig()),
	)

	srv.HandleFunc("/hygrothermograph", HygrothermographHandler)

	//sr := srv.Route("/my")
	//sr.GET("/info2", info2())

	if err := srv.Start(ctx); err != nil {
		panic(err)
	}

	defer func() {
		if err := srv.Stop(ctx); err != nil {
			t.Errorf("expected nil got %v", err)
		}
	}()
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

func CreateHygrothermograph(ctx context.Context, cli *khttp.Client, in *api.Hygrothermograph, opts ...khttp.CallOption) (*api.Hygrothermograph, error) {
	var out api.Hygrothermograph

	pattern := "/hygrothermograph"
	path := binding.EncodeURL(pattern, in, false)

	opts = append(opts, khttp.Operation("/CreateHygrothermograph"))
	opts = append(opts, khttp.PathTemplate(pattern))

	err := cli.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}

	return &out, nil
}

func GetMyInfo2(ctx context.Context, cli *khttp.Client, in *MyInfo2Req, opts ...khttp.CallOption) (*MyInfo2Res, error) {
	var out MyInfo2Res

	pattern := "/my/info2"
	path := binding.EncodeURL(pattern, in, false)

	opts = append(opts, khttp.Operation("/my/info2"))
	opts = append(opts, khttp.PathTemplate(pattern))

	err := cli.Invoke(ctx, "GET", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}

	return &out, nil
}

func TestClient(t *testing.T) {
	ctx := context.Background()

	var qconf quic.Config

	tlsConf := tlsutil.MustInsecureTLSConfig()
	cli, err := khttp.NewClient(ctx,
		khttp.WithEndpoint("127.0.0.1:8800"),
		khttp.WithTLSConfig(tlsConf),
		khttp.WithTransport(&http3.RoundTripper{TLSClientConfig: tlsConf, QUICConfig: &qconf}),
	)
	assert.Nil(t, err)
	assert.NotNil(t, cli)

	var req api.Hygrothermograph
	req.Humidity = strconv.FormatInt(int64(rand.Intn(100)), 10)
	req.Temperature = strconv.FormatInt(int64(rand.Intn(100)), 10)

	resp, err := GetHygrothermograph(ctx, cli, &req, khttp.EmptyCallOption{})
	assert.Nil(t, err)
	t.Log(resp)

	resp, err = CreateHygrothermograph(ctx, cli, &req, khttp.EmptyCallOption{})
	assert.Nil(t, err)
	t.Log(resp)

	iresp, ierr := GetMyInfo2(ctx, cli, &MyInfo2Req{
		Action: "Test My Info2",
	}, khttp.EmptyCallOption{})
	assert.Nil(t, ierr)
	t.Log(iresp)
}
