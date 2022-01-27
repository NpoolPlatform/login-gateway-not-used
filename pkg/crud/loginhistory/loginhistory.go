package loginhistory

import (
	"context"
	"time"

	db "github.com/NpoolPlatform/login-gateway/pkg/db"
	npool "github.com/NpoolPlatform/message/npool/logingateway"

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

	_, err = cli.
		LoginHistory.
		Create().
		SetClientIP(in.GetClientIP()).
		SetUserAgent(in.GetUserAgent()).
		Save(ctx)
	if err != nil {
		return xerrors.Errorf("fail create login history: %v", err)
	}

	return nil
}
