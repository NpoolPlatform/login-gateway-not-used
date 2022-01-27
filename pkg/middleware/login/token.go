package login

import (
	"context"
	"net"
	"os"

	appusermgrpb "github.com/NpoolPlatform/message/npool/appusermgr"
	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc/metadata"

	"github.com/google/uuid"

	"golang.org/x/xerrors"
)

type Metadata struct {
	AppID     uuid.UUID
	Account   string
	LoginType string
	UserID    uuid.UUID
	ClientIP  net.IP
	UserAgent string
	UserInfo  *appusermgrpb.AppUserInfo
}

func MetadataFromContext(ctx context.Context) (*Metadata, error) {
	meta, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, xerrors.Errorf("fail get metadata")
	}

	clientIP := ""
	if forwards, ok := meta["x-forwarded-for"]; ok {
		if len(forwards) > 0 {
			clientIP = forwards[0]
		}
	}

	userAgent := ""
	if agents, ok := meta["grpcgateway-user-agent"]; ok {
		if len(agents) > 0 {
			userAgent = agents[0]
		}
	}

	return &Metadata{
		ClientIP:  net.ParseIP(clientIP),
		UserAgent: userAgent,
	}, nil
}

func (meta *Metadata) ToJWTClaims() jwt.MapClaims {
	claims := jwt.MapClaims{}

	claims["app_id"] = meta.AppID
	claims["user_id"] = meta.UserID
	claims["account"] = meta.Account
	claims["login_type"] = meta.LoginType
	claims["client_ip"] = meta.ClientIP
	claims["user_agent"] = meta.UserAgent

	return claims
}

func createToken(meta *Metadata) (string, error) {
	tokenAccessSecret := os.Getenv("LOGIN_TOKEN_ACCESS_SECRET")
	if tokenAccessSecret == "" {
		return "", xerrors.Errorf("invalid login token access secret")
	}

	claims := meta.ToJWTClaims()
	candidate := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := candidate.SignedString([]byte(tokenAccessSecret))
	if err != nil {
		return "", xerrors.Errorf("fail sign jwt claims: %v", err)
	}

	return token, nil
}
