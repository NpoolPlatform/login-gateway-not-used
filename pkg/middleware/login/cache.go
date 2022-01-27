package login

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	redis2 "github.com/NpoolPlatform/go-service-framework/pkg/redis"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"

	"golang.org/x/xerrors"
)

const (
	redisTimeout    = 5 * time.Second
	loginExpiration = 4 * time.Hour
)

func appAccountKey(appID uuid.UUID, account, loginType string) string {
	return fmt.Sprintf("%v:%v:%v", appID, account, loginType)
}

func metaToAccountKey(meta *Metadata) string {
	return appAccountKey(meta.AppID, meta.Account, meta.LoginType)
}

func appUserKey(appID, userID uuid.UUID) string {
	return fmt.Sprintf("%v:%v", appID, userID)
}

func metaToUserKey(meta *Metadata) string {
	return appUserKey(meta.AppID, meta.UserID)
}

type valAppUser struct {
	Account   string
	LoginType string
}

func createCache(ctx context.Context, meta *Metadata) error {
	cli, err := redis2.GetClient()
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

func queryByAppAccount(ctx context.Context, appID uuid.UUID, account, loginType string) (*Metadata, error) {
	cli, err := redis2.GetClient()
	if err != nil {
		return nil, xerrors.Errorf("fail get redis client: %v", err)
	}

	ctx, cancel := context.WithTimeout(ctx, redisTimeout)
	defer cancel()

	val, err := cli.Get(ctx, appAccountKey(appID, account, loginType)).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, xerrors.Errorf("fail get app account: %v", err)
	}

	meta := Metadata{}
	err = json.Unmarshal([]byte(val), &meta)
	if err != nil {
		return nil, xerrors.Errorf("fail unmarshal val: %v", err)
	}

	return &meta, nil
}

/*
func queryByAppUser(ctx context.Context, appID, userID uuid.UUID) (*Metadata, error) {
	cli, err := redis2.GetClient()
	if err != nil {
		return nil, xerrors.Errorf("fail get redis client: %v", err)
	}

	ctx, cancel := context.WithTimeout(ctx, redisTimeout)
	defer cancel()

	val, err := cli.Get(ctx, appUserKey(appID, userID)).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, xerrors.Errorf("fail get app user: %v", err)
	}

	meta := Metadata{}
	err = json.Unmarshal([]byte(val), &meta)
	if err != nil {
		return nil, xerrors.Errorf("fail unmarshal val: %v", err)
	}

	return &meta, nil
}
*/
