package grpc

import (
	"fmt"
	"log"
	"net"

	"github.com/OleksUMD/grpc_test_server/internal/interceptor"
	"github.com/OleksUMD/grpc_test_server/internal/port"
	"github.com/OleksUMD/proto-example/protogen/go/bank"
	"github.com/OleksUMD/proto-example/protogen/go/hello"
	"github.com/OleksUMD/proto-example/protogen/go/resiliency"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type GrpcAdapter struct {
	helloService      port.HelloServicePort
	bankService       port.BankServicePort
	resiliencyService port.ResiliencyServicePort
	grpcPort          int
	server            *grpc.Server
	hello.HelloServiceServer
	bank.BankServiceServer
	resiliency.ResiliencyServiceServer
}

func NewGrpcAdapter(
	helloService port.HelloServicePort,
	bankService port.BankServicePort,
	resiliencyService port.ResiliencyServicePort,
	grpcPort int,
) *GrpcAdapter {
	return &GrpcAdapter{
		helloService:      helloService,
		bankService:       bankService,
		resiliencyService: resiliencyService,
		grpcPort:          grpcPort,
	}
}

func (a *GrpcAdapter) Run() {
	var err error

	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", a.grpcPort))
	if err != nil {
		log.Fatalf("Failed to start grpc server on port %d: %v\n", a.grpcPort, err)
	}
	log.Printf("Server listening on port %d\n", a.grpcPort)
	tlsCreds, err := credentials.NewServerTLSFromFile("./certs/server.crt", "./certs/server.pem")
	if err != nil {
		log.Fatalln("Failed to start initialize TLS server", err)
	}
	grpcServer := grpc.NewServer(
		grpc.Creds(tlsCreds),
		grpc.ChainUnaryInterceptor(
			interceptor.LogUnaryServerInterceptor(),
			interceptor.BasicUnaryServerInterceptor(),
		),
		grpc.ChainStreamInterceptor(
			interceptor.LogStreamServerInterceptor(),
			interceptor.BasicStreamServerInterceptor(),
		),
	)
	a.server = grpcServer
	hello.RegisterHelloServiceServer(grpcServer, a)
	bank.RegisterBankServiceServer(grpcServer, a)
	resiliency.RegisterResiliencyServiceServer(grpcServer, a)
	if err = grpcServer.Serve(listen); err != nil {
		log.Fatalf("Failed to start grpc server on port %d: %v\n", a.grpcPort, err)
	}
}

func (a *GrpcAdapter) Stop() {
	a.server.Stop()
}
