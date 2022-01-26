package login

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	grpc2 "github.com/NpoolPlatform/login-gateway/pkg/grpc"
	npool "github.com/NpoolPlatform/message/npool/logingateway"

	appusermgrpb "github.com/NpoolPlatform/message/npool/appusermgr"

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

	resp, err := grpc2.VerifyAppUserByAppAccountPassword(ctx, &appusermgrpb.VerifyAppUserByAppAccountPasswordRequest{
		AppID:        in.GetAppID(),
		Account:      in.GetAccount(),
		PasswordHash: in.GetPasswordHash(),
	})
	if err != nil {
		return nil, xerrors.Errorf("fail verify username or password: %v", err)
	}

	return &npool.LoginResponse{
		Info: resp.Info,
	}, nil
}
