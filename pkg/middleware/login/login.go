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

	cached := false
	token := ""

	meta, err := queryByAppAccount(ctx, appID, in.GetAccount(), in.GetLoginType())
	if err != nil {
		return nil, xerrors.Errorf("fail query login cache by app acount: %v", err)
	}
	if meta != nil {
		if in.GetToken() != "" {
			token = in.GetToken()

			err := verifyToken(meta, in.GetToken())
			if err == nil {
				cached = true
			}
		}
	}

	if !cached {
		meta, err = MetadataFromContext(ctx)
		if err != nil {
			return nil, xerrors.Errorf("fail create login metadata: %v", err)
		}
		meta.AppID = appID
		meta.Account = in.GetAccount()
		meta.LoginType = in.GetLoginType()

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
		UserID:    meta.UserInfo.User.ID,
		ClientIP:  meta.ClientIP.String(),
		UserAgent: meta.UserAgent,
	})
	if err != nil {
		return nil, xerrors.Errorf("fail create login history: %v", err)
	}
	// TODO: check login type of app

	return &npool.LoginResponse{
		Info:  meta.UserInfo,
		Token: token,
	}, nil
}

func Logined(ctx context.Context, in *npool.LoginedRequest) (*npool.LoginedResponse, error) {
	appID, err := uuid.Parse(in.GetAppID())
	if err != nil {
		return nil, xerrors.Errorf("invalid app id: %v", err)
	}

	userID, err := uuid.Parse(in.GetUserID())
	if err != nil {
		return nil, xerrors.Errorf("invalid user id: %v", err)
	}

	if in.GetToken() == "" {
		return nil, xerrors.Errorf("invalid token")
	}

	meta, err := queryByAppUser(ctx, appID, userID)
	if err != nil {
		return nil, xerrors.Errorf("fail query login cache by app user: %v", err)
	}
	if meta == nil {
		return &npool.LoginedResponse{}, nil
	}

	err = verifyToken(meta, in.GetToken())
	if err != nil {
		return &npool.LoginedResponse{}, nil
	}

	err = createCache(ctx, meta)
	if err != nil {
		return nil, xerrors.Errorf("fail create cache: %v", err)
	}

	return &npool.LoginedResponse{
		Info: meta.UserInfo,
	}, nil
}

func Logout(ctx context.Context, in *npool.LogoutRequest) (*npool.LogoutResponse, error) {
	appID, err := uuid.Parse(in.GetAppID())
	if err != nil {
		return nil, xerrors.Errorf("invalid app id: %v", err)
	}

	userID, err := uuid.Parse(in.GetUserID())
	if err != nil {
		return nil, xerrors.Errorf("invalid user id: %v", err)
	}

	meta, err := queryByAppUser(ctx, appID, userID)
	if err != nil {
		return nil, xerrors.Errorf("fail query login cache by app user: %v", err)
	}

	err = deleteCache(ctx, meta)
	if err != nil {
		return nil, xerrors.Errorf("fail delete login cache: %v", err)
	}

	return &npool.LogoutResponse{
		Info: meta.UserInfo,
	}, nil
}
