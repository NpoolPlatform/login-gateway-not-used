package loginhistory

import (
	"context"
	"time"

	db "github.com/NpoolPlatform/login-gateway/pkg/db"
	"github.com/NpoolPlatform/login-gateway/pkg/db/ent/loginhistory"
	npool "github.com/NpoolPlatform/message/npool/logingateway"

	"github.com/google/uuid"

	"golang.org/x/xerrors"
)

const (
	dbTimeout = 5 * time.Second
)

func Create(ctx context.Context, in *npool.LoginHistory) error {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	cli, err := db.Client()
	if err != nil {
		return xerrors.Errorf("fail get db client: %v", err)
	}

	appID, err := uuid.Parse(in.GetAppID())
	if err != nil {
		return xerrors.Errorf("invalid app id: %v", err)
	}

	userID, err := uuid.Parse(in.GetUserID())
	if err != nil {
		return xerrors.Errorf("invalid user id: %v", err)
	}

	_, err = cli.
		LoginHistory.
		Create().
		SetAppID(appID).
		SetUserID(userID).
		SetClientIP(in.GetClientIP()).
		SetUserAgent(in.GetUserAgent()).
		Save(ctx)
	if err != nil {
		return xerrors.Errorf("fail create login history: %v", err)
	}

	return nil
}

func GetAll(ctx context.Context, in *npool.GetLoginHistoriesRequest) (*npool.GetLoginHistoriesResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	cli, err := db.Client()
	if err != nil {
		return nil, xerrors.Errorf("fail get db client: %v", err)
	}

	appID, err := uuid.Parse(in.GetAppID())
	if err != nil {
		return nil, xerrors.Errorf("invalid app id: %v", err)
	}

	userID, err := uuid.Parse(in.GetUserID())
	if err != nil {
		return nil, xerrors.Errorf("invalid user id: %v", err)
	}

	infos, err := cli.
		LoginHistory.
		Query().
		Where(
			loginhistory.And(
				loginhistory.AppID(appID),
				loginhistory.UserID(userID),
			),
		).
		All(ctx)
	if err != nil {
		return nil, xerrors.Errorf("fail query login history: %v", err)
	}

	histories := []*npool.LoginHistory{}
	for _, info := range infos {
		histories = append(histories, &npool.LoginHistory{
			ID:        info.ID.String(),
			AppID:     info.AppID.String(),
			UserID:    info.UserID.String(),
			ClientIP:  info.ClientIP,
			UserAgent: info.UserAgent,
			CreateAt:  info.CreateAt,
		})
	}

	return &npool.GetLoginHistoriesResponse{
		Infos: histories,
	}, nil
}
