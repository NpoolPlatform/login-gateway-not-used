package grpc

import (
	"context"
	"fmt"
	"time"

	grpc2 "github.com/NpoolPlatform/go-service-framework/pkg/grpc"

	appusermgrconst "github.com/NpoolPlatform/appuser-manager/pkg/message/const" //nolint
	appusermgrpb "github.com/NpoolPlatform/message/npool/appusermgr"

	thirdgwpb "github.com/NpoolPlatform/message/npool/thirdgateway"
	thirdgwconst "github.com/NpoolPlatform/third-gateway/pkg/message/const"

	thirdlogingwpb "github.com/NpoolPlatform/message/npool/third-login-gateway"
	thirdlogingwconst "github.com/NpoolPlatform/third-login-gateway/pkg/message/const"
)

const (
	grpcTimeout = 5 * time.Second
)

func VerifyAppUserByAppAccountPassword(ctx context.Context, in *appusermgrpb.VerifyAppUserByAppAccountPasswordRequest) (*appusermgrpb.VerifyAppUserByAppAccountPasswordResponse, error) {
	conn, err := grpc2.GetGRPCConn(appusermgrconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get app user manager connection: %v", err)
	}
	defer conn.Close()

	cli := appusermgrpb.NewAppUserManagerClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	return cli.VerifyAppUserByAppAccountPassword(ctx, in)
}

func GetAppInfo(ctx context.Context, in *appusermgrpb.GetAppInfoRequest) (*appusermgrpb.GetAppInfoResponse, error) {
	conn, err := grpc2.GetGRPCConn(appusermgrconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get app user manager connection: %v", err)
	}
	defer conn.Close()

	cli := appusermgrpb.NewAppUserManagerClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	return cli.GetAppInfo(ctx, in)
}

// =======================================================================================================================

func VerifyGoogleRecaptchaV3(ctx context.Context, in *thirdgwpb.VerifyGoogleRecaptchaV3Request) (*thirdgwpb.VerifyGoogleRecaptchaV3Response, error) {
	conn, err := grpc2.GetGRPCConn(thirdgwconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get third gateway connection: %v", err)
	}

	defer conn.Close()

	client := thirdgwpb.NewThirdGatewayClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	return client.VerifyGoogleRecaptchaV3(ctx, in)
}

func AuthLogin(ctx context.Context, in *thirdlogingwpb.AuthLoginRequest) (*appusermgrpb.AppUserInfo, error) {
	conn, err := grpc2.GetGRPCConn(thirdlogingwconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get third login gateway connection: %v", err)
	}
	defer conn.Close()

	cli := thirdlogingwpb.NewThirdLoginGatewayClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.AuthLogin(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail auth login: %v", err)
	}

	return resp.Info, nil
}
