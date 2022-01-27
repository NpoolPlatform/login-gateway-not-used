package login

import (
	"context"
	"fmt"
	"time"

	"github.com/NpoolPlatform/go-service-framework/pkg/redis"

	"github.com/google/uuid"

	"golang.org/x/xerrors"
)

const (
	redisTimeout    = 5 * time.Second
	loginExpiration = 4 * time.Hour
)

func metaToAccountKey(meta *Metadata) string {
	return fmt.Sprintf("%v:%v:%v", meta.AppID, meta.Account, meta.LoginType)
}

func metaToUserKey(meta *Metadata) string {
	return fmt.Sprintf("%v:%v", meta.AppID, meta.UserID)
}

type valAppUser struct {
	Account   string
	LoginType string
}

func CreateCache(ctx context.Context, meta *Metadata) error {
	cli, err := redis.GetClient()
	if err != nil {
		return xerrors.Errorf("fail get redis client: %v", err)
	}

	ctx, cancel := context.WithTimeout(ctx, redisTimeout)
	defer cancel()

	err = cli.Set(ctx, metaToAccountKey(meta), meta, loginExpiration).Err()
	if err != nil {
		return xerrors.Errorf("fail create login account cache: %v", err)
	}

	err = cli.Set(ctx, metaToUserKey(meta), &valAppUser{
		Account:   meta.Account,
		LoginType: meta.LoginType,
	}, loginExpiration).Err()
	if err != nil {
		return xerrors.Errorf("fail create login user cache: %v", err)
	}

	return nil
}

func QueryByAppAccount(ctx context.Context, appID uuid.UUID, account string) (*Metadata, error) {
	return nil, nil
}

func QueryByAppUser(ctx context.Context, appID, userID uuid.UUID) (*Metadata, error) {
	return nil, nil
}
