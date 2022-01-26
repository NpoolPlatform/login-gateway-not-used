package api

import (
	"context"

	npool "github.com/NpoolPlatform/message/npool/logingateway"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	npool.UnimplementedLoginGatewayServer
}

func Register(server grpc.ServiceRegistrar) {
	npool.RegisterLoginGatewayServer(server, &Server{})
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return npool.RegisterLoginGatewayHandlerFromEndpoint(context.Background(), mux, endpoint, opts)
}
