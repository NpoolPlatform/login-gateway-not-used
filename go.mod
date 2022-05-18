module github.com/NpoolPlatform/login-gateway

go 1.16

require (
	entgo.io/ent v0.10.1
	github.com/NpoolPlatform/api-manager v0.0.0-20220505084652-c0caff45e937
	github.com/NpoolPlatform/appuser-manager v0.0.0-20220503155754-feb61d1897e0
	github.com/NpoolPlatform/go-service-framework v0.0.0-20220404143809-82c40930388a
	github.com/NpoolPlatform/message v0.0.0-20220503165256-540362bedc23
	github.com/NpoolPlatform/third-gateway v0.0.0-20220205094704-a585d53bd025
	github.com/NpoolPlatform/third-login-gateway v0.0.0-20220511110308-4d1de890428c
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-redis/redis/v8 v8.11.4
	github.com/go-resty/resty/v2 v2.7.0
	github.com/google/uuid v1.3.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.10.0
	github.com/streadway/amqp v1.0.0
	github.com/stretchr/testify v1.7.1
	github.com/urfave/cli/v2 v2.5.1
	google.golang.org/grpc v1.46.2
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.2.0
	google.golang.org/protobuf v1.28.0
)

replace google.golang.org/grpc => github.com/grpc/grpc-go v1.41.0

replace (
	github.com/NpoolPlatform/appuser-manager => github.com/NpoolPlatform/appuser-manager v0.0.0-20220518102322-fda5c6f5c73c
	github.com/NpoolPlatform/message => github.com/NpoolPlatform/message v0.0.0-20220518093034-9f64c6f0c028
	github.com/NpoolPlatform/third-login-gateway => github.com/NpoolPlatform/third-login-gateway v0.0.0-20220518103408-fe9e0ed44901
)
