package login

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	redis2 "github.com/NpoolPlatform/go-service-framework/pkg/redis"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

const (
	redisTimeout    = 5 * time.Second
	loginExpiration = 4 * time.Hour
)

func appAccountKey(appID uuid.UUID, account, loginType string) string {
	return fmt.Sprintf("login-%v:%v:%v", appID, account, loginType)
}

func metaToAccountKey(meta *Metadata) string {
	return appAccountKey(meta.AppID, meta.Account, meta.AccountType)
}

func appUserKey(appID, userID uuid.UUID) string {
	return fmt.Sprintf("login-%v:%v", appID, userID)
}

func metaToUserKey(meta *Metadata) string {
	return appUserKey(meta.AppID, meta.UserID)
}

type valAppUser struct {
	Account     string
	AccountType string
}

func createCache(ctx context.Context, meta *Metadata) error {
	cli, err := redis2.GetClient()
	if err != nil {
		return fmt.Errorf("fail get redis client: %v", err)
	}

	ctx, cancel := context.WithTimeout(ctx, redisTimeout)
	defer cancel()

	body, err := json.Marshal(meta)
	if err != nil {
		return fmt.Errorf("fail marshal login account meta: %v", err)
	}

	err = cli.Set(ctx, metaToAccountKey(meta), body, loginExpiration).Err()
	if err != nil {
		return fmt.Errorf("fail create login account cache: %v", err)
	}

	body, err = json.Marshal(&valAppUser{
		Account:     meta.Account,
		AccountType: meta.AccountType,
	})
	if err != nil {
		return fmt.Errorf("fail marshal login user meta: %v", err)
	}

	err = cli.Set(ctx, metaToUserKey(meta), body, loginExpiration).Err()
	if err != nil {
		return fmt.Errorf("fail create login user cache: %v", err)
	}

	return nil
}

func queryByAppAccount(ctx context.Context, appID uuid.UUID, account, loginType string) (*Metadata, error) {
	cli, err := redis2.GetClient()
	if err != nil {
		return nil, fmt.Errorf("fail get redis client: %v", err)
	}

	ctx, cancel := context.WithTimeout(ctx, redisTimeout)
	defer cancel()

	val, err := cli.Get(ctx, appAccountKey(appID, account, loginType)).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("fail get app account: %v", err)
	}

	meta := Metadata{}
	err = json.Unmarshal([]byte(val), &meta)
	if err != nil {
		return nil, fmt.Errorf("fail unmarshal val: %v", err)
	}

	return &meta, nil
}

func queryByAppUser(ctx context.Context, appID, userID uuid.UUID) (*Metadata, error) {
	cli, err := redis2.GetClient()
	if err != nil {
		return nil, fmt.Errorf("fail get redis client: %v", err)
	}

	ctx, cancel := context.WithTimeout(ctx, redisTimeout)
	defer cancel()

	val, err := cli.Get(ctx, appUserKey(appID, userID)).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("fail get app user: %v", err)
	}

	appUser := valAppUser{}
	err = json.Unmarshal([]byte(val), &appUser)
	if err != nil {
		return nil, fmt.Errorf("fail unmarshal val: %v", err)
	}

	val, err = cli.Get(ctx, appAccountKey(appID, appUser.Account, appUser.AccountType)).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("fail get app user: %v", err)
	}

	meta := Metadata{}
	err = json.Unmarshal([]byte(val), &meta)
	if err != nil {
		return nil, fmt.Errorf("fail unmarshal val: %v", err)
	}

	return &meta, nil
}

func deleteCache(ctx context.Context, meta *Metadata) error {
	cli, err := redis2.GetClient()
	if err != nil {
		return fmt.Errorf("fail get redis client: %v", err)
	}

	ctx, cancel := context.WithTimeout(ctx, redisTimeout)
	defer cancel()

	err = cli.Del(ctx, metaToUserKey(meta)).Err()
	if err != nil {
		return fmt.Errorf("fail delete login user cache")
	}

	err = cli.Del(ctx, metaToAccountKey(meta)).Err()
	if err != nil {
		return fmt.Errorf("fail delete login account cache")
	}

	return nil
}
