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

	var userInfo *appusermgrpb.AppUserInfo
	cached := false
	token := ""

	query, err := queryByAppAccount(ctx, appID, in.GetAccount(), in.GetLoginType())
	if err != nil {
		return nil, xerrors.Errorf("fail query login cache by app acount: %v", err)
	}
	if query != nil {
		// TODO: verify login info and token
		userInfo = query.UserInfo
		cached = true
	}

	if !cached {
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

		userInfo = resp.Info
	}

	meta.UserInfo = userInfo
	meta.UserID = uuid.MustParse(userInfo.User.ID)

	if !cached {
		token, err = createToken(meta)
		if err != nil {
			return nil, xerrors.Errorf("fail create token: %v", err)
		}
	}

	err = createCache(ctx, meta)
	if err != nil {
		return nil, xerrors.Errorf("fail create cache: %v", err)
	}

	err = loginhistorycrud.Create(ctx, &npool.LoginHistory{
		AppID:     in.GetAppID(),
		UserID:    userInfo.User.ID,
		ClientIP:  meta.ClientIP.String(),
		UserAgent: meta.UserAgent,
	})
	if err != nil {
		return nil, xerrors.Errorf("fail create login history: %v", err)
	}
	// TODO: check login type of app

	return &npool.LoginResponse{
		Info:  userInfo,
		Token: token,
	}, nil
}
