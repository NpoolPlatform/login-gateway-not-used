package login

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/logingateway"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"

	"golang.org/x/xerrors"
)

func Login(ctx context.Context, in *npool.LoginRequest) (*npool.LoginResponse, error) {
	myPeer, ok := peer.FromContext(ctx)
	if !ok {
		return nil, xerrors.Errorf("fail get peer")
	}

	meta, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, xerrors.Errorf("fail get metadata")
	}

	logger.Sugar().Infof("addr: %v auth: %v", myPeer.Addr, myPeer.AuthInfo)
	for k, v := range meta {
		logger.Sugar().Infof("key: %v value: %v", k, v)
	}

	return &npool.LoginResponse{}, nil
}
