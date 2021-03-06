package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"

	"github.com/golang/glog"

	foo "github.com/AmandaCameron/protoc-gen-gokit/examples/foo-service"
)

func main() {
	go grpcServer()
	ctx := context.Background()

	conn, err := grpc.Dial("localhost:8675", grpc.WithInsecure())
	if err != nil {
		glog.Fatal(err)
	}

	cli := foo.NewFooServiceClient(conn)

	http.ListenAndServe("localhost:8000", foo.MakeMux_FooService(cli, kithttp.Server{
		Context: ctx,
		Logger:  log.NewLogfmtLogger(os.Stderr),
		EncodeResponseFunc: func(wr http.ResponseWriter, data interface{}) error {
			return json.NewEncoder(wr).Encode(data)
		},
	}))
}

func grpcServer() {
	lis, err := net.Listen("tcp", ":8675")
	if err != nil {
		glog.Fatal(err)
	}

	srv := grpc.NewServer()

	foo.RegisterFooServiceServer(srv, &impl{})

	srv.Serve(lis)
}

type impl struct{}

func (i *impl) SayHello(ctx context.Context, req *foo.HelloRequest) (*foo.HelloResponse, error) {
	return &foo.HelloResponse{
		Response: "Hello " + req.Who + "!",
	}, nil
}

func (i *impl) CountTo(ctx context.Context, req *foo.CountToRequest) (*foo.CountToResponse, error) {
	resp := &foo.CountToResponse{}

	for i := int32(0); i < req.Target; i++ {
		resp.Response += fmt.Sprintf(" %d", i+1)
	}

	resp.Response = resp.Response[1:]

	return resp, nil
}
