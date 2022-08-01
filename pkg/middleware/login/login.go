package login

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	loginhistorycrud "github.com/NpoolPlatform/login-gateway/pkg/crud/loginhistory"
	grpc2 "github.com/NpoolPlatform/login-gateway/pkg/grpc"
	npool "github.com/NpoolPlatform/message/npool/logingateway"

	appusermgrconst "github.com/NpoolPlatform/appuser-manager/pkg/const"
	appusermgrpb "github.com/NpoolPlatform/message/npool/appuser/mgr/v1"

	thirdgwpb "github.com/NpoolPlatform/message/npool/thirdgateway"

	"github.com/google/uuid"

	"golang.org/x/xerrors"
)

func Login(ctx context.Context, in *npool.LoginRequest) (*npool.LoginResponse, error) { //nolint
	resp, err := grpc2.GetAppInfo(ctx, &appusermgrpb.GetAppInfoRequest{
		ID: in.GetAppID(),
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get app info: %v", err)
	}

	if false && resp.Info.Ctrl != nil && resp.Info.Ctrl.RecaptchaMethod == appusermgrconst.RecaptchaGoogleV3 {
		if in.GetManMachineSpec() == "" {
			return nil, xerrors.Errorf("miss recaptcha")
		}

		resp, err := grpc2.VerifyGoogleRecaptchaV3(ctx, &thirdgwpb.VerifyGoogleRecaptchaV3Request{
			RecaptchaToken: in.GetManMachineSpec(),
		})
		if err != nil {
			return nil, xerrors.Errorf("fail verify google recaptcha: %v", err)
		}

		if resp.Code < 0 {
			return nil, xerrors.Errorf("invalid google recaptcha response")
		}
	}

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

	meta, err := queryByAppAccount(ctx, appID, in.GetAccount(), in.GetAccountType())
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
		meta.AccountType = in.GetAccountType()

		resp, err := grpc2.VerifyAppUserByAppAccountPassword(ctx, &appusermgrpb.VerifyAppUserByAppAccountPasswordRequest{
			AppID:        in.GetAppID(),
			Account:      in.GetAccount(),
			PasswordHash: in.GetPasswordHash(),
		})
		if err != nil {
			return nil, xerrors.Errorf("fail verify username or password: %v", err)
		}

		// TODO: correct login type according to account match

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
		logger.Sugar().Warnf("user %v not in cache", in)
		return &npool.LoginedResponse{}, nil
	}

	err = verifyToken(meta, in.GetToken())
	if err != nil {
		logger.Sugar().Warnf("user %v token not in cache: %v", in, err)
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

func UpdateCache(ctx context.Context, in *npool.UpdateCacheRequest) (*npool.UpdateCacheResponse, error) {
	appID, err := uuid.Parse(in.GetInfo().GetUser().GetAppID())
	if err != nil {
		return nil, xerrors.Errorf("invalid app id: %v", err)
	}

	userID, err := uuid.Parse(in.GetInfo().GetUser().GetID())
	if err != nil {
		return nil, xerrors.Errorf("invalid user id: %v", err)
	}

	meta, err := queryByAppUser(ctx, appID, userID)
	if err != nil {
		return nil, xerrors.Errorf("fail query login cache by app user: %v", err)
	}

	meta.UserInfo = in.GetInfo()
	err = createCache(ctx, meta)
	if err != nil {
		return nil, xerrors.Errorf("fail delete login cache: %v", err)
	}

	return &npool.UpdateCacheResponse{
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
