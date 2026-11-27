package grpc

import (
	"context"
	"fmt"
	"io"

	"github.com/OleksUMD/proto-example/protogen/go/hello"
	"google.golang.org/grpc"
)

func (a *GrpcAdapter) SayHello(ctx context.Context, req *hello.HelloRequest) (*hello.HelloResponse, error) {
	greet := a.helloService.GenerateHello(req.Name)
	return &hello.HelloResponse{
		Greet: greet,
	}, nil
}

func (a *GrpcAdapter) ManyHello(req *hello.HelloRequest, stream grpc.ServerStreamingServer[hello.HelloResponse]) error {
	for ix := range 10 {
		resp := &hello.HelloResponse{
			Greet: fmt.Sprintf("[%d] Hello [%s]", ix, req.Name),
		}
		if err := stream.Send(resp); err != nil {
			return err
		}
	}
	return nil
}

func (a *GrpcAdapter) HelloToEveryOne(stream grpc.ClientStreamingServer[hello.HelloRequest, hello.HelloResponse]) error {
	var greet string
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		fmt.Printf("Received request with name: %s\n", req.Name)
		greet += fmt.Sprintf("Hello [%s]; ", req.Name)
	}
	resp := &hello.HelloResponse{
		Greet: greet,
	}
	if err := stream.SendAndClose(resp); err != nil {
		return err
	}
	return nil
}

func (a *GrpcAdapter) HelloContinuous(stream grpc.BidiStreamingServer[hello.HelloRequest, hello.HelloResponse]) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		greet := a.helloService.GenerateHello(req.Name)
		resp := &hello.HelloResponse{
			Greet: greet,
		}
		if err := stream.Send(resp); err != nil {
			return err
		}
	}
}
