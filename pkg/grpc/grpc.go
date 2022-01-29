package grpc

import (
	"context"
	"time"

	grpc2 "github.com/NpoolPlatform/go-service-framework/pkg/grpc"

	appusermgrconst "github.com/NpoolPlatform/appuser-manager/pkg/message/const" //nolint
	appusermgrpb "github.com/NpoolPlatform/message/npool/appusermgr"

	verificationpb "github.com/NpoolPlatform/message/npool/verification"
	verificationconst "github.com/NpoolPlatform/verification-door/pkg/message/const"

	"golang.org/x/xerrors"
)

const (
	grpcTimeout = 5 * time.Second
)

func VerifyAppUserByAppAccountPassword(ctx context.Context, in *appusermgrpb.VerifyAppUserByAppAccountPasswordRequest) (*appusermgrpb.VerifyAppUserByAppAccountPasswordResponse, error) {
	conn, err := grpc2.GetGRPCConn(appusermgrconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, xerrors.Errorf("fail get app user manager connection: %v", err)
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
		return nil, xerrors.Errorf("fail get app user manager connection: %v", err)
	}
	defer conn.Close()

	cli := appusermgrpb.NewAppUserManagerClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	return cli.GetAppInfo(ctx, in)
}

// =======================================================================================================================

func VerifyGoogleRecaptcha(ctx context.Context, in *verificationpb.VerifyGoogleRecaptchaRequest) (*verificationpb.VerifyGoogleRecaptchaResponse, error) {
	conn, err := grpc2.GetGRPCConn(verificationconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, xerrors.Errorf("fail get verification connection: %v", err)
	}

	defer conn.Close()

	client := verificationpb.NewVerificationDoorClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	return client.VerifyGoogleRecaptcha(ctx, in)
}
