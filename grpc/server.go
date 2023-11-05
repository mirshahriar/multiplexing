package grpc

import (
	"context"
	"fmt"

	"github.com/mirshahriar/multiplexing/grpc/proto"
	"google.golang.org/grpc"
)

type grpcServer struct{}

func (s *grpcServer) EchoMessage(ctx context.Context, req *proto.EchoRequest) (*proto.EchoResponse, error) {
	fmt.Println("echo message", req.Message)
	return &proto.EchoResponse{Message: fmt.Sprintf("echo %s from grpc", req.Message)}, nil
}

func NewGRPCServer() *grpc.Server {
	server := grpc.NewServer()
	proto.RegisterEchoServiceServer(server, &grpcServer{})
	return server
}
