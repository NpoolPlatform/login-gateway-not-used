package login

import (
	"context"

	loginhistorycrud "github.com/NpoolPlatform/login-gateway/pkg/crud/loginhistory"
	grpc2 "github.com/NpoolPlatform/login-gateway/pkg/grpc"
	npool "github.com/NpoolPlatform/message/npool/logingateway"

	appusermgrpb "github.com/NpoolPlatform/message/npool/appusermgr"

	"github.com/google/uuid"

	"golang.org/x/xerrors"
)

func Login(ctx context.Context, in *npool.LoginRequest) (*npool.LoginResponse, error) {
	// TODO: check man machine spec
	// TODO: check environment, if safe, just login

	appID, err := uuid.Parse(in.GetAppID())
	if err != nil {
		return nil, xerrors.Errorf("invalid app id: %v", err)
	}

	if in.GetAccount() == "" {
		return nil, xerrors.Errorf("invalid account: %v", err)
	}

	meta, err := MetadataFromContext(ctx)
	if err != nil {
		return nil, xerrors.Errorf("fail create login metadata: %v", err)
	}

	meta.AppID = appID
	meta.Account = in.GetAccount()
	meta.LoginType = in.GetLoginType()

	// TODO: check if cached
	// TODO: if verify is not OK, login again

	resp, err := grpc2.VerifyAppUserByAppAccountPassword(ctx, &appusermgrpb.VerifyAppUserByAppAccountPasswordRequest{
		AppID:        in.GetAppID(),
		Account:      in.GetAccount(),
		PasswordHash: in.GetPasswordHash(),
	})
	if err != nil {
		return nil, xerrors.Errorf("fail verify username or password: %v", err)
	}

	meta.UserInfo = resp.Info
	meta.UserID = uuid.MustParse(resp.Info.User.ID)

	token, err := createToken(meta)
	if err != nil {
		return nil, xerrors.Errorf("fail create token: %v", err)
	}

	// TODO: add to redis

	err = loginhistorycrud.Create(ctx, &npool.LoginHistory{
		AppID:     in.GetAppID(),
		UserID:    resp.Info.User.ID,
		ClientIP:  meta.ClientIP.String(),
		UserAgent: meta.UserAgent,
	})
	if err != nil {
		return nil, xerrors.Errorf("fail create login history: %v", err)
	}
	// TODO: check login type of app

	return &npool.LoginResponse{
		Info:  resp.Info,
		Token: token,
	}, nil
}
