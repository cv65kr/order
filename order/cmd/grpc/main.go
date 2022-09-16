package main

import (
	"context"
	pb_order "github.com/cv65kr/order/api/orderservice/v1"
	"github.com/cv65kr/order/internal"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

func main() {
	fx.New(
		internal.Module,
		fx.Invoke(GrpcServerHooks),
	).Run()
}

func GrpcServerHooks(lifecycle fx.Lifecycle, log *zap.Logger, handler *internal.WorkflowHandler) {
	lifecycle.Append(
		fx.Hook{
			OnStart: func(context.Context) error {
				lis, err := net.Listen("tcp", ":8081")
				if err != nil {
					log.Fatal("failed to listen: %v", zap.Error(err))
				}
				var opts []grpc.ServerOption
				grpcServer := grpc.NewServer(opts...)
				pb_order.RegisterOrderServiceServer(grpcServer, handler)

				reflection.Register(grpcServer)
				go grpcServer.Serve(lis)

				return nil
			},
			OnStop: func(context.Context) error {
				return log.Sync()
			},
		},
	)
}
