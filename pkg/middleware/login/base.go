package login

import (
	"context"
	"fmt"

	appusermgrconst "github.com/NpoolPlatform/appuser-manager/pkg/const"
	grpc2 "github.com/NpoolPlatform/login-gateway/pkg/grpc"
	appusermgrpb "github.com/NpoolPlatform/message/npool/appusermgr"
	npool "github.com/NpoolPlatform/message/npool/logingateway"
	thirdlogingwpb "github.com/NpoolPlatform/message/npool/third-login-gateway"
)

var VerifyMap = make(map[string]VerifyMethod)

func init() {
	VerifyMap[appusermgrconst.ThirdGithub] = &ThirdAuthVerify{}
	VerifyMap[appusermgrconst.ThirdGoogle] = &ThirdAuthVerify{}
	VerifyMap[appusermgrconst.SignupByMobile] = &AccountPasswordVerify{}
	VerifyMap[appusermgrconst.SignupByEmail] = &AccountPasswordVerify{}
}

type VerifyMethod interface {
	Verify(ctx context.Context, in *npool.LoginRequest) (*appusermgrpb.AppUserInfo, error)
}

type ThirdAuthVerify struct{}

func (ThirdAuthVerify) Verify(ctx context.Context, in *npool.LoginRequest) (*appusermgrpb.AppUserInfo, error) {
	resp, err := grpc2.AuthLogin(ctx, &thirdlogingwpb.AuthLoginRequest{
		Code:  in.GetAccount(),
		AppID: in.GetAppID(),
		Third: in.GetAccountType(),
	})
	if err != nil {
		return nil, fmt.Errorf("fail auth login: %v", err)
	}
	if resp == nil {
		return nil, fmt.Errorf("fail auth login")
	}
	return resp, nil
}

type AccountPasswordVerify struct{}

func (AccountPasswordVerify) Verify(ctx context.Context, in *npool.LoginRequest) (*appusermgrpb.AppUserInfo, error) {
	resp, err := grpc2.VerifyAppUserByAppAccountPassword(ctx, &appusermgrpb.VerifyAppUserByAppAccountPasswordRequest{
		AppID:        in.GetAppID(),
		Account:      in.GetAccount(),
		PasswordHash: in.GetPasswordHash(),
	})
	if err != nil {
		return nil, fmt.Errorf("fail verify username or password: %v", err)
	}
	if resp == nil {
		return nil, fmt.Errorf("fail verify username or password")
	}
	return resp, err
}
