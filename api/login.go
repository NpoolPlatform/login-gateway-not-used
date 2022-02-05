package api

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	mw "github.com/NpoolPlatform/login-gateway/pkg/middleware/login"
	loginhistorymw "github.com/NpoolPlatform/login-gateway/pkg/middleware/loginhistory"
	npool "github.com/NpoolPlatform/message/npool/logingateway"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) Login(ctx context.Context, in *npool.LoginRequest) (*npool.LoginResponse, error) {
	resp, err := mw.Login(ctx, in)
	if err != nil {
		logger.Sugar().Errorf("fail login: %v", err)
		return &npool.LoginResponse{}, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}

func (s *Server) Logined(ctx context.Context, in *npool.LoginedRequest) (*npool.LoginedResponse, error) {
	resp, err := mw.Logined(ctx, in)
	if err != nil {
		logger.Sugar().Errorf("fail logined: %v", err)
		return &npool.LoginedResponse{}, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}

func (s *Server) Logout(ctx context.Context, in *npool.LogoutRequest) (*npool.LogoutResponse, error) {
	resp, err := mw.Logout(ctx, in)
	if err != nil {
		logger.Sugar().Errorf("fail logout: %v", err)
		return &npool.LogoutResponse{}, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}

func (s *Server) GetLoginHistories(ctx context.Context, in *npool.GetLoginHistoriesRequest) (*npool.GetLoginHistoriesResponse, error) {
	resp, err := loginhistorymw.GetByAppUser(ctx, in)
	if err != nil {
		logger.Sugar().Errorf("fail get login history: %v", err)
		return &npool.GetLoginHistoriesResponse{}, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}
