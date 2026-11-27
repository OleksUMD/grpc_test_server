package interceptor

import (
	"context"
	"log"

	"github.com/OleksUMD/proto-example/protogen/go/hello"
	"github.com/OleksUMD/proto-example/protogen/go/resiliency"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type InterceptedServerStream struct {
	grpc.ServerStream
}

func (s *InterceptedServerStream) RecvMsg(msg any) error {
	if err := s.ServerStream.RecvMsg(msg); err != nil {
		return err
	}
	switch request := msg.(type) {
	case *hello.HelloRequest:
		request.Name = "[MODIFIED BY SERVER INTERCEPTOR - 4] " + request.Name
	}
	return nil
}

func (s *InterceptedServerStream) SendMsg(msg any) error {
	switch response := msg.(type) {
	case *hello.HelloResponse:
		response.Greet = "[MODIFIED BY SERVER INTERCEPTOR - 5] " + response.Greet
	case *resiliency.ResiliencyResponse:
		response.DummyString = "[MODIFIED BY SERVER INTERCEPTOR - 6] " + response.DummyString
	}
	return s.ServerStream.SendMsg(msg)
}

func LogUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		log.Println("[LOGGED BY SERVER INTERCEPTOR] ", req)
		return handler(ctx, req)
	}
}

func BasicUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		switch request := req.(type) {
		case *hello.HelloRequest:
			request.Name = "[MODIFIED BY SERVER INTERCEPTOR - 1] " + request.Name
		}
		responseMeta, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			responseMeta = metadata.New(nil)
		}
		responseMeta.Append("from-interceptor-key-1", "from-interceptor-value-1")
		responseMeta.Append("from-interceptor-key-2", "from-interceptor-value-2")
		ctx = metadata.NewOutgoingContext(ctx, responseMeta)
		grpc.SetHeader(ctx, responseMeta)
		resp, err = handler(ctx, req)
		if err != nil {
			return resp, err
		}
		switch response := resp.(type) {
		case *hello.HelloResponse:
			response.Greet = "[MODIFIED BY SERVER INTERCEPTOR - 2] " + response.Greet
		case *resiliency.ResiliencyResponse:
			response.DummyString = "[MODIFIED BY SERVER INTERCEPTOR - 3] " + response.DummyString
		}
		return resp, nil
	}
}

func LogStreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(serv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		log.Println("[LOGGED BY SERVER INTERCEPTOR] ", info)
		return handler(serv, ss)
	}
}

func BasicStreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(serv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		interceptedServerStream := &InterceptedServerStream{
			ServerStream: ss,
		}
		responseMeta, ok := metadata.FromOutgoingContext(interceptedServerStream.Context())
		if !ok {
			responseMeta = metadata.New(nil)
		}
		responseMeta.Append("from-interceptor-key-1", "from-interceptor-value-1")
		responseMeta.Append("from-interceptor-key-2", "from-interceptor-value-2")
		interceptedServerStream.SetHeader(responseMeta)

		log.Println("[LOGGED BY SERVER INTERCEPTOR] ", info)
		return handler(serv, interceptedServerStream)
	}
}
